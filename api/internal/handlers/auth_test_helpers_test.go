package handlers

import (
	"fmt"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"testing"
)

// Handlers now enforce modern permissions via PermissionService, which resolves
// the acting user's permissions from the database (the JWT is never trusted).
// These helpers seed a user/member with a modern role so handler tests that
// exercise gated endpoints are authorized. They upsert, so they work whether or
// not the test already created the user/membership, and they clear the role
// permission cache (the test database reuses role ids across truncations).

// grantAppPerms gives userId an app role granting exactly perms, creating the
// user row if it does not already exist. Call it after any user setup.
func grantAppPerms(t *testing.T, userId uint, perms ...string) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	db := repositories.GetDB()

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateAppRole(fmt.Sprintf("Test App Role %d", userId), "", perms)
	if err != nil {
		t.Fatalf("create app role: %v", err)
	}

	var count int64
	db.Model(&models.User{}).Where("id = ?", userId).Count(&count)
	if count == 0 {
		err = db.Create(&models.User{
			BaseModel: models.BaseModel{ID: userId},
			Username:  fmt.Sprintf("perm-user-%d", userId),
			Password:  "password",
			AppRoleID: &role.ID,
		}).Error
	} else {
		err = db.Model(&models.User{}).Where("id = ?", userId).Update("app_role_id", role.ID).Error
	}
	if err != nil {
		t.Fatalf("assign app role: %v", err)
	}
}

// grantGroupPerms gives the (userId, groupId) member a group role granting
// exactly perms, creating the membership if it does not already exist.
func grantGroupPerms(t *testing.T, userId uint, groupId uint, perms ...string) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	// Also drop the group-role grant cache: a truncated test DB reuses role ids,
	// so a fresh (grant-less) role here could otherwise resolve to another case's
	// cached grants — and paid-by visibility now consults that cache on every
	// receipt read (e.g. the pie chart), wrongly restricting this role's members.
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole(fmt.Sprintf("Test Group Role %d-%d", userId, groupId), "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("create group role: %v", err)
	}

	var count int64
	db.Model(&models.GroupMember{}).Where("user_id = ? AND group_id = ?", userId, groupId).Count(&count)
	if count == 0 {
		err = db.Create(&models.GroupMember{UserID: userId, GroupID: groupId, GroupRoleID: &role.ID}).Error
	} else {
		err = db.Model(&models.GroupMember{}).Where("user_id = ? AND group_id = ?", userId, groupId).Update("group_role_id", role.ID).Error
	}
	if err != nil {
		t.Fatalf("assign group role: %v", err)
	}
}

// grantAllAppPerms gives userId every app-scope permission (legacy-admin
// equivalent). Use for tests that exercise an app-gated handler's behavior
// rather than its authorization.
func grantAllAppPerms(t *testing.T, userId uint) {
	t.Helper()
	grantAppPerms(t, userId, permissions.LegacyAppAdminKeys()...)
}

// grantAllGroupPerms gives the (userId, groupId) member every group-scope
// permission (legacy-owner equivalent).
func grantAllGroupPerms(t *testing.T, userId uint, groupId uint) {
	t.Helper()
	grantGroupPerms(t, userId, groupId, permissions.LegacyGroupOwnerKeys()...)
}

// isolateGroupWithSupervisor turns on member-presence isolation for groupId and
// assigns supervisorUserId a SeesAllMembers group role in it (so they remain
// visible to isolated members). Members without such a role become mutually
// invisible. Assumes the membership rows already exist.
func isolateGroupWithSupervisor(t *testing.T, groupId uint, supervisorUserId uint) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	if err := db.Model(&models.Group{}).Where("id = ?", groupId).Update("isolate_members", true).Error; err != nil {
		t.Fatalf("isolate group: %v", err)
	}

	supRole, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		fmt.Sprintf("Iso Supervisor %d-%d", groupId, supervisorUserId), "",
		[]string{permissions.GroupReceiptsRead}, nil, nil, nil, false, true)
	if err != nil {
		t.Fatalf("create supervisor role: %v", err)
	}
	if err := db.Model(&models.GroupMember{}).
		Where("user_id = ? AND group_id = ?", supervisorUserId, groupId).
		Update("group_role_id", supRole.ID).Error; err != nil {
		t.Fatalf("assign supervisor role: %v", err)
	}
}
