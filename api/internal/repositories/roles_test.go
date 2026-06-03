package repositories

import (
	"errors"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"sort"
	"testing"

	"gorm.io/gorm"
)

func TestNewRoleRepository(t *testing.T) {
	repository := NewRoleRepository(nil)

	if repository.DB == nil {
		utils.PrintTestError(t, repository.DB, "a database instance")
	}

	if repository.TX != nil {
		utils.PrintTestError(t, repository.TX, "nil")
	}
}

func TestCreateAppRolePersistsPermissions(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	perms := []string{permissions.AppUsersCreate, permissions.AppUsersRead}
	role, err := repository.CreateAppRole("App Role", "Description", perms)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if role.ID == 0 {
		utils.PrintTestError(t, role.ID, "non-zero id")
	}

	if len(role.Permissions) != 2 {
		utils.PrintTestError(t, len(role.Permissions), 2)
	}
}

func TestCreateGroupRolePersistsPermissions(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	perms := []string{permissions.GroupReceiptsCreate}
	role, err := repository.CreateGroupRole("Group Role", "Description", perms)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if role.ID == 0 {
		utils.PrintTestError(t, role.ID, "non-zero id")
	}

	if len(role.Permissions) != 1 {
		utils.PrintTestError(t, len(role.Permissions), 1)
	}
}

func TestUpdateAppRolePersistsChanges(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	created, err := repository.CreateAppRole("App Role", "Description", []string{permissions.AppUsersCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	updated, err := repository.UpdateAppRole(created.ID, "Renamed Role", "New description", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if updated.Name != "Renamed Role" {
		utils.PrintTestError(t, updated.Name, "Renamed Role")
	}

	if updated.Description != "New description" {
		utils.PrintTestError(t, updated.Description, "New description")
	}
}

func TestUpdateGroupRolePersistsChanges(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	created, err := repository.CreateGroupRole("Group Role", "Description", []string{permissions.GroupReceiptsCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	updated, err := repository.UpdateGroupRole(created.ID, "Renamed Group Role", "New description", []string{permissions.GroupReceiptsRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if updated.Name != "Renamed Group Role" {
		utils.PrintTestError(t, updated.Name, "Renamed Group Role")
	}

	if updated.Description != "New description" {
		utils.PrintTestError(t, updated.Description, "New description")
	}
}

func TestUpdateAppRoleReplacesPermissions(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	created, err := repository.CreateAppRole("App Role", "Description", []string{permissions.AppUsersCreate, permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	updated, err := repository.UpdateAppRole(created.ID, "App Role", "Description", []string{permissions.AppUsersDelete})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(updated.Permissions) != 1 {
		utils.PrintTestError(t, len(updated.Permissions), 1)
		return
	}

	if updated.Permissions[0].Permission != permissions.AppUsersDelete {
		utils.PrintTestError(t, updated.Permissions[0].Permission, permissions.AppUsersDelete)
	}
}

func TestGetAppRoleByIdNotFoundReturnsError(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	_, err := repository.GetAppRoleById(999)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.PrintTestError(t, err, gorm.ErrRecordNotFound)
	}
}

func TestGetAllRolesReturnsBothScopes(t *testing.T) {
	defer TruncateTestDb()
	CreateTestRoles()

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roles) != 2 {
		utils.PrintTestError(t, len(roles), 2)
		return
	}

	// App roles come first.
	appRole := roles[0]
	if appRole.Scope != permissions.ScopeApp {
		utils.PrintTestError(t, appRole.Scope, permissions.ScopeApp)
	}
	if len(appRole.Permissions) != 1 || appRole.Permissions[0] != permissions.AppUsersCreate {
		utils.PrintTestError(t, appRole.Permissions, []string{permissions.AppUsersCreate})
	}

	groupRole := roles[1]
	if groupRole.Scope != permissions.ScopeGroup {
		utils.PrintTestError(t, groupRole.Scope, permissions.ScopeGroup)
	}
	if len(groupRole.Permissions) != 1 || groupRole.Permissions[0] != permissions.GroupReceiptsCreate {
		utils.PrintTestError(t, groupRole.Permissions, []string{permissions.GroupReceiptsCreate})
	}
}

func TestGetAppRolePermissions(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	perms := []string{permissions.AppUsersCreate, permissions.AppUsersRead}
	role, err := repository.CreateAppRole("App Role", "", perms)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	got, err := repository.GetAppRolePermissions(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	expected := []string{permissions.AppUsersCreate, permissions.AppUsersRead}
	sort.Strings(got)
	sort.Strings(expected)
	if !slices.Equal(got, expected) {
		utils.PrintTestError(t, got, expected)
	}

	// A role with no permissions resolves to an empty (non-nil) slice.
	empty, err := repository.CreateAppRole("Empty Role", "", []string{})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	gotEmpty, err := repository.GetAppRolePermissions(empty.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if gotEmpty == nil || len(gotEmpty) != 0 {
		utils.PrintTestError(t, gotEmpty, "empty slice")
	}
}

func TestGetGroupRolePermissions(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	perms := []string{permissions.GroupReceiptsRead}
	role, err := repository.CreateGroupRole("Group Role", "", perms)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	got, err := repository.GetGroupRolePermissions(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(got) != 1 || got[0] != permissions.GroupReceiptsRead {
		utils.PrintTestError(t, got, perms)
	}
}

func TestGetUserAppRoleId(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)
	db := GetDB()

	role, err := repository.CreateAppRole("App Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	withRole := models.User{Username: "with-role", Password: "password", AppRoleID: &role.ID}
	if err := db.Create(&withRole).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	gotId, err := repository.GetUserAppRoleId(withRole.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if gotId == nil || *gotId != role.ID {
		utils.PrintTestError(t, gotId, role.ID)
	}

	noRole := models.User{Username: "no-role", Password: "password"}
	if err := db.Create(&noRole).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	gotNil, err := repository.GetUserAppRoleId(noRole.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if gotNil != nil {
		utils.PrintTestError(t, gotNil, nil)
	}

	if _, err := repository.GetUserAppRoleId(999); !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.PrintTestError(t, err, gorm.ErrRecordNotFound)
	}
}

func TestGetGroupMemberRoleId(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)
	db := GetDB()

	group := models.Group{Name: "role-id-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	role, err := repository.CreateGroupRole("Group Role", "", []string{permissions.GroupReceiptsRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	user := models.User{Username: "member", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRole: models.OWNER, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	gotId, err := repository.GetGroupMemberRoleId(user.ID, group.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if gotId == nil || *gotId != role.ID {
		utils.PrintTestError(t, gotId, role.ID)
	}

	// A user who is not a member of the group is a record-not-found.
	if _, err := repository.GetGroupMemberRoleId(user.ID, 999); !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.PrintTestError(t, err, gorm.ErrRecordNotFound)
	}
}
