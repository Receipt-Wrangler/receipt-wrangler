package repositories

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
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

func TestGetPagedRolesReturnsBothScopesOrderedByName(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	if _, err := repository.CreateAppRole("Beta App", "", []string{permissions.AppUsersRead}); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if _, err := repository.CreateGroupRole("Alpha Group", "", []string{permissions.GroupReceiptsCreate, permissions.GroupReceiptsRead}); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if _, err := repository.CreateAppRole("Zeta App", "", nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	command := commands.PagedRoleRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      25,
			OrderBy:       "name",
			SortDirection: commands.ASCENDING,
		},
	}

	roles, count, err := repository.GetPagedRoles(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if count != 3 {
		utils.PrintTestError(t, count, 3)
	}
	if len(roles) != 3 {
		utils.PrintTestError(t, len(roles), 3)
		return
	}

	// Ordered by name ascending across both scopes (union).
	expectedOrder := []string{"Alpha Group", "Beta App", "Zeta App"}
	for i, name := range expectedOrder {
		if roles[i].Name != name {
			utils.PrintTestError(t, roles[i].Name, name)
		}
	}

	if roles[0].Scope != permissions.ScopeGroup {
		utils.PrintTestError(t, roles[0].Scope, permissions.ScopeGroup)
	}
	if len(roles[0].Permissions) != 2 {
		utils.PrintTestError(t, len(roles[0].Permissions), 2)
	}
	if roles[1].Scope != permissions.ScopeApp {
		utils.PrintTestError(t, roles[1].Scope, permissions.ScopeApp)
	}
	// A role with no permissions yields an empty (non-nil) slice.
	if roles[2].Permissions == nil || len(roles[2].Permissions) != 0 {
		utils.PrintTestError(t, roles[2].Permissions, []string{})
	}
}

func TestGetPagedRolesFiltersByScope(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	if _, err := repository.CreateAppRole("App One", "", nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if _, err := repository.CreateAppRole("App Two", "", nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if _, err := repository.CreateGroupRole("Group One", "", nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	appCommand := commands.PagedRoleRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{Page: 1, PageSize: 25, OrderBy: "name"},
		Filter:              commands.RoleFilter{Scope: permissions.ScopeApp},
	}
	appRoles, appCount, err := repository.GetPagedRoles(appCommand)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if appCount != 2 || len(appRoles) != 2 {
		utils.PrintTestError(t, appCount, 2)
	}
	for _, role := range appRoles {
		if role.Scope != permissions.ScopeApp {
			utils.PrintTestError(t, role.Scope, permissions.ScopeApp)
		}
	}

	groupCommand := commands.PagedRoleRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{Page: 1, PageSize: 25, OrderBy: "name"},
		Filter:              commands.RoleFilter{Scope: permissions.ScopeGroup},
	}
	groupRoles, groupCount, err := repository.GetPagedRoles(groupCommand)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if groupCount != 1 || len(groupRoles) != 1 {
		utils.PrintTestError(t, groupCount, 1)
		return
	}
	if groupRoles[0].Scope != permissions.ScopeGroup {
		utils.PrintTestError(t, groupRoles[0].Scope, permissions.ScopeGroup)
	}
}

func TestGetPagedRolesPaginates(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	for _, name := range []string{"Role A", "Role B", "Role C"} {
		if _, err := repository.CreateAppRole(name, "", nil); err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
	}

	command := commands.PagedRoleRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          2,
			PageSize:      1,
			OrderBy:       "name",
			SortDirection: commands.ASCENDING,
		},
	}

	roles, count, err := repository.GetPagedRoles(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Count is the full total, independent of the page size.
	if count != 3 {
		utils.PrintTestError(t, count, 3)
	}
	if len(roles) != 1 {
		utils.PrintTestError(t, len(roles), 1)
		return
	}
	if roles[0].Name != "Role B" {
		utils.PrintTestError(t, roles[0].Name, "Role B")
	}
}

func TestGetPagedRolesIncludesAssignedCount(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	role, err := repository.CreateAppRole("Assigned Role", "", nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	user := models.User{Username: "assigned-user", Password: "password", DisplayName: "Assigned User", AppRoleID: &role.ID}
	if err := repository.GetDB().Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	command := commands.PagedRoleRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{Page: 1, PageSize: 25, OrderBy: "name"},
		Filter:              commands.RoleFilter{Scope: permissions.ScopeApp},
	}

	roles, _, err := repository.GetPagedRoles(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(roles) != 1 {
		utils.PrintTestError(t, len(roles), 1)
		return
	}
	if roles[0].AssignedCount != 1 {
		utils.PrintTestError(t, roles[0].AssignedCount, 1)
	}
}

func TestGetPagedRolesInvalidOrderByDefaultsToName(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	if _, err := repository.CreateAppRole("Bravo", "", nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if _, err := repository.CreateAppRole("Alpha", "", nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	command := commands.PagedRoleRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      25,
			OrderBy:       "id; DROP TABLE app_roles",
			SortDirection: commands.ASCENDING,
		},
	}

	roles, _, err := repository.GetPagedRoles(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(roles) != 2 {
		utils.PrintTestError(t, len(roles), 2)
		return
	}
	// Falls back to ordering by name.
	if roles[0].Name != "Alpha" {
		utils.PrintTestError(t, roles[0].Name, "Alpha")
	}
}
