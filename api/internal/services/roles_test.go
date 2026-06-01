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
