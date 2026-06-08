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

func TestSetDefaultAppRoleClearsOthers(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	first, err := repository.CreateAppRole("First", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	second, err := repository.CreateAppRole("Second", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := repository.SetDefaultAppRole(first.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := repository.SetDefaultAppRole(second.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	defaultId, err := repository.GetDefaultAppRoleId()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if defaultId == nil || *defaultId != second.ID {
		utils.PrintTestError(t, defaultId, second.ID)
	}

	// Exactly one app role is the default.
	var count int64
	if err := GetDB().Model(&models.AppRole{}).Where("is_default = ?", true).Count(&count).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 1 {
		utils.PrintTestError(t, count, 1)
	}
}

func TestSetDefaultGroupRoleClearsOthers(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	first, err := repository.CreateGroupRole("First", "", []string{permissions.GroupReceiptsRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	second, err := repository.CreateGroupRole("Second", "", []string{permissions.GroupReceiptsRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := repository.SetDefaultGroupRole(first.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := repository.SetDefaultGroupRole(second.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	defaultId, err := repository.GetDefaultGroupRoleId()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if defaultId == nil || *defaultId != second.ID {
		utils.PrintTestError(t, defaultId, second.ID)
	}

	var count int64
	if err := GetDB().Model(&models.GroupRoleDefinition{}).Where("is_default = ?", true).Count(&count).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 1 {
		utils.PrintTestError(t, count, 1)
	}
}

func TestGetDefaultAppRoleIdNilWhenUnset(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	if _, err := repository.CreateAppRole("Role", "", []string{permissions.AppUsersRead}); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	id, err := repository.GetDefaultAppRoleId()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if id != nil {
		utils.PrintTestError(t, id, nil)
	}
}

func TestGetAppRoleIdByName(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	created, err := repository.CreateAppRole("Named Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	id, err := repository.GetAppRoleIdByName("Named Role")
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if id == nil || *id != created.ID {
		utils.PrintTestError(t, id, created.ID)
	}

	// An unknown name resolves to nil, not an error.
	missing, err := repository.GetAppRoleIdByName("Nope")
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if missing != nil {
		utils.PrintTestError(t, missing, nil)
	}
}

func TestGetAllRolesReturnsIsDefault(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	role, err := repository.CreateAppRole("Default App", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := repository.SetDefaultAppRole(role.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	found, ok := findRole(roles, "Default App")
	if !ok || !found.IsDefault {
		utils.PrintTestError(t, "Default App IsDefault", true)
	}
}

func TestDeriveLegacyUserRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)

	legacyAdminId, err := repository.GetAppRoleIdByName(LegacyAdminRoleName)
	if err != nil || legacyAdminId == nil {
		utils.PrintTestError(t, err, "legacy admin id")
		return
	}
	role, err := repository.DeriveLegacyUserRole(*legacyAdminId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if role != models.ADMIN {
		utils.PrintTestError(t, role, models.ADMIN)
	}

	legacyUserId, err := repository.GetAppRoleIdByName(LegacyUserRoleName)
	if err != nil || legacyUserId == nil {
		utils.PrintTestError(t, err, "legacy user id")
		return
	}
	role, err = repository.DeriveLegacyUserRole(*legacyUserId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if role != models.USER {
		utils.PrintTestError(t, role, models.USER)
	}

	// A custom (non-system) app role maps to the least-privilege USER.
	custom, err := repository.CreateAppRole("Auditor", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	role, err = repository.DeriveLegacyUserRole(custom.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if role != models.USER {
		utils.PrintTestError(t, role, models.USER)
	}
}

func TestDeriveLegacyGroupRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)

	cases := []struct {
		roleName string
		expected models.GroupRole
	}{
		{LegacyOwnerRoleName, models.OWNER},
		{LegacyEditorRoleName, models.EDITOR},
		{LegacyViewerRoleName, models.VIEWER},
	}
	for _, c := range cases {
		var groupRole models.GroupRoleDefinition
		if err := GetDB().Select("id").Where("name = ?", c.roleName).First(&groupRole).Error; err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
		derived, err := repository.DeriveLegacyGroupRole(groupRole.ID)
		if err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
		if derived != c.expected {
			utils.PrintTestError(t, derived, c.expected)
		}
	}

	// A custom (non-system) group role maps to the least-privilege VIEWER.
	custom, err := repository.CreateGroupRole("Auditors", "", []string{permissions.GroupReceiptsRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	derived, err := repository.DeriveLegacyGroupRole(custom.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if derived != models.VIEWER {
		utils.PrintTestError(t, derived, models.VIEWER)
	}
}
