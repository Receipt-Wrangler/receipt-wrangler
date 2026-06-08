package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpdateRolePersistsChanges(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)

	created, err := roleRepository.CreateAppRole("App Role", "Description", []string{permissions.AppUsersCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:        "Renamed Role",
		Description: "New description",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersRead},
	}

	roleView, err := service.UpdateRole(created.ID, command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if roleView.Name != "Renamed Role" {
		utils.PrintTestError(t, roleView.Name, "Renamed Role")
	}

	if roleView.Scope != permissions.ScopeApp {
		utils.PrintTestError(t, roleView.Scope, permissions.ScopeApp)
	}

	if len(roleView.Permissions) != 1 || roleView.Permissions[0] != permissions.AppUsersRead {
		utils.PrintTestError(t, roleView.Permissions, []string{permissions.AppUsersRead})
	}
}

func TestUpdateRoleBlocksTypeSwitch(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)

	created, err := roleRepository.CreateAppRole("App Role", "Description", []string{permissions.AppUsersCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:        "App Role",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReceiptsCreate},
	}

	_, err = service.UpdateRole(created.ID, command)
	if !errors.Is(err, ErrRoleTypeMismatch) {
		utils.PrintTestError(t, err, ErrRoleTypeMismatch)
	}
}

func TestUpdateRoleBlocksSystemRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	systemRole := models.AppRole{
		Name:        "System Role",
		Description: "system role",
		IsSystem:    true,
		Permissions: []models.AppRolePermission{
			{Permission: permissions.AppUsersCreate},
		},
	}
	if err := db.Create(&systemRole).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:        "Renamed System Role",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersRead},
	}

	_, err := service.UpdateRole(systemRole.ID, command)
	if !errors.Is(err, ErrSystemRoleImmutable) {
		utils.PrintTestError(t, err, ErrSystemRoleImmutable)
	}
}

func TestUpdateRoleNotFound(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:        "Missing Role",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersRead},
	}

	_, err := service.UpdateRole(999, command)
	if !errors.Is(err, ErrRoleNotFound) {
		utils.PrintTestError(t, err, ErrRoleNotFound)
	}
}

func TestSetDefaultRoleApp(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)

	created, err := roleRepository.CreateAppRole("App Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	view, err := service.SetDefaultRole(created.ID, permissions.ScopeApp)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !view.IsDefault {
		utils.PrintTestError(t, view.IsDefault, true)
	}

	defaultId, err := roleRepository.GetDefaultAppRoleId()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if defaultId == nil || *defaultId != created.ID {
		utils.PrintTestError(t, defaultId, created.ID)
	}
}

func TestSetDefaultRoleAllowsSystemRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	systemRole := models.GroupRoleDefinition{
		Name:     "System Group Role",
		IsSystem: true,
		Permissions: []models.GroupRolePermission{
			{Permission: permissions.GroupReceiptsRead},
		},
	}
	if err := db.Create(&systemRole).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	view, err := service.SetDefaultRole(systemRole.ID, permissions.ScopeGroup)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !view.IsDefault {
		utils.PrintTestError(t, view.IsDefault, true)
	}
}

func TestSetDefaultRoleTypeMismatch(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)

	created, err := roleRepository.CreateAppRole("App Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	// The id exists under APP, so requesting it under GROUP is a type mismatch.
	_, err = service.SetDefaultRole(created.ID, permissions.ScopeGroup)
	if !errors.Is(err, ErrRoleTypeMismatch) {
		utils.PrintTestError(t, err, ErrRoleTypeMismatch)
	}
}

func TestSetDefaultRoleNotFound(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := NewRoleService(nil)
	_, err := service.SetDefaultRole(999, permissions.ScopeApp)
	if !errors.Is(err, ErrRoleNotFound) {
		utils.PrintTestError(t, err, ErrRoleNotFound)
	}
}

func TestDeleteRoleRejectsDefault(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)

	created, err := roleRepository.CreateAppRole("Default App Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := roleRepository.SetDefaultAppRole(created.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	err = service.DeleteRole(created.ID, permissions.ScopeApp)
	if !errors.Is(err, ErrRoleIsDefault) {
		utils.PrintTestError(t, err, ErrRoleIsDefault)
	}

	// The default role must be untouched.
	if _, err := roleRepository.GetAppRoleById(created.ID); err != nil {
		utils.PrintTestError(t, err, "default role should be untouched")
	}
}
