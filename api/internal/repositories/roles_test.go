package repositories

import (
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"testing"
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
