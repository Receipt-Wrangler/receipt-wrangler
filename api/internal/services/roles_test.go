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

func TestCreateGroupRoleWithGrantsExposesGrantIds(t *testing.T) {
	defer repositories.TruncateTestDb()

	category := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&category)
	tag := models.Tag{Name: "Reimbursable"}
	repositories.GetDB().Create(&tag)

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:           "Restricted Group Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{category.ID},
		TagGrants:      []uint{tag.ID},
	}

	roleView, err := service.CreateRole(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roleView.CategoryGrants) != 1 || roleView.CategoryGrants[0] != category.ID {
		utils.PrintTestError(t, roleView.CategoryGrants, []uint{category.ID})
	}
	if len(roleView.TagGrants) != 1 || roleView.TagGrants[0] != tag.ID {
		utils.PrintTestError(t, roleView.TagGrants, []uint{tag.ID})
	}
}

func TestGroupRoleSeesAllMembersRoundTrips(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := NewRoleService(nil)
	created, err := service.CreateRole(commands.UpsertRoleCommand{
		Name:           "Supervisor Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		SeesAllMembers: true,
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !created.SeesAllMembers {
		utils.PrintTestError(t, created.SeesAllMembers, true)
	}

	// It surfaces on the read model too.
	roles, err := service.GetRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	found := false
	for _, role := range roles {
		if role.Id == created.Id {
			found = true
			if !role.SeesAllMembers {
				utils.PrintTestError(t, role.SeesAllMembers, true)
			}
		}
	}
	if !found {
		utils.PrintTestError(t, "role not found in GetRoles", created.Id)
	}

	// Update toggles it off; the returned view reflects false.
	updated, err := service.UpdateRole(created.Id, commands.UpsertRoleCommand{
		Name:           "Supervisor Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		SeesAllMembers: false,
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if updated.SeesAllMembers {
		utils.PrintTestError(t, updated.SeesAllMembers, false)
	}
}

func TestCreateGroupRoleRejectsNonExistentGrant(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:           "Bad Grant Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{999999},
	}

	_, err := service.CreateRole(command)
	if !errors.Is(err, ErrInvalidGrant) {
		utils.PrintTestError(t, err, ErrInvalidGrant)
	}
}

func TestCreateGroupRoleWithPaidByGrantsExposesGrantIds(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "paidby-payer", Password: "x"}
	repositories.GetDB().Create(&payer)

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:                   "Paid-By Group Role",
		Scope:                  permissions.ScopeGroup,
		Permissions:            []string{permissions.GroupReceiptsRead},
		PaidByUserGrants:       []uint{payer.ID},
		IncludeOwnPaidReceipts: true,
	}

	roleView, err := service.CreateRole(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roleView.PaidByUserGrants) != 1 || roleView.PaidByUserGrants[0] != payer.ID {
		utils.PrintTestError(t, roleView.PaidByUserGrants, []uint{payer.ID})
	}
	if !roleView.IncludeOwnPaidReceipts {
		utils.PrintTestError(t, roleView.IncludeOwnPaidReceipts, true)
	}
}

func TestCreateGroupRoleRejectsNonExistentPaidByGrant(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:             "Bad Paid-By Role",
		Scope:            permissions.ScopeGroup,
		Permissions:      []string{permissions.GroupReceiptsRead},
		PaidByUserGrants: []uint{999999},
	}

	_, err := service.CreateRole(command)
	if !errors.Is(err, ErrInvalidGrant) {
		utils.PrintTestError(t, err, ErrInvalidGrant)
	}
}

func TestCreateGroupRoleWithReportTemplateGrantsExposesGrants(t *testing.T) {
	defer repositories.TruncateTestDb()

	template := models.ReportTemplate{Name: "Quarterly", ConfigurationVersion: 1}
	repositories.GetDB().Create(&template)

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:        "Report Restricted Role",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []commands.ReportTemplateGrantCommand{
			{ReportTemplateId: template.ID, Permissions: []string{"read", "generate"}},
		},
	}

	roleView, err := service.CreateRole(command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roleView.ReportTemplateGrants) != 1 {
		utils.PrintTestError(t, len(roleView.ReportTemplateGrants), 1)
		return
	}
	grant := roleView.ReportTemplateGrants[0]
	if grant.ReportTemplateId != template.ID || len(grant.Permissions) != 2 {
		utils.PrintTestError(t, grant, "template with 2 actions")
	}

	// It persists and reads back through GetRoles (grouped per template).
	roles, err := service.GetRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	found := false
	for _, r := range roles {
		if r.Id != roleView.Id {
			continue
		}
		found = true
		if len(r.ReportTemplateGrants) != 1 || r.ReportTemplateGrants[0].ReportTemplateId != template.ID {
			utils.PrintTestError(t, r.ReportTemplateGrants, "one report template grant")
		}
	}
	if !found {
		utils.PrintTestError(t, "role not found in GetRoles", roleView.Id)
	}
}

func TestCreateGroupRoleRejectsNonExistentReportTemplateGrant(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:        "Bad Report Grant Role",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []commands.ReportTemplateGrantCommand{
			{ReportTemplateId: 999999, Permissions: []string{"read"}},
		},
	}

	_, err := service.CreateRole(command)
	if !errors.Is(err, ErrInvalidGrant) {
		utils.PrintTestError(t, err, ErrInvalidGrant)
	}
}

func TestUpdateGroupRoleServiceReplacesGrants(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)

	catA := models.Category{Name: "A"}
	repositories.GetDB().Create(&catA)
	catB := models.Category{Name: "B"}
	repositories.GetDB().Create(&catB)

	created, err := roleRepository.CreateGroupRole("Role", "", []string{permissions.GroupReceiptsRead}, []uint{catA.ID}, nil, nil, false, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	service := NewRoleService(nil)
	command := commands.UpsertRoleCommand{
		Name:           "Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{catB.ID},
	}

	roleView, err := service.UpdateRole(created.ID, command)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roleView.CategoryGrants) != 1 || roleView.CategoryGrants[0] != catB.ID {
		utils.PrintTestError(t, roleView.CategoryGrants, []uint{catB.ID})
	}
}

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
