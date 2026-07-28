package commands

import (
	"receipt-wrangler/api/internal/permissions"
	"testing"
)

func TestUpsertRoleCommandValidValidAppCommand(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "App Role",
		Description: "An app role",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersCreate, permissions.AppUsersRead},
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandValidGroupCommand(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Group Role",
		Description: "A group role",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReceiptsCreate, permissions.GroupReceiptsQuickScan},
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandMissingName(t *testing.T) {
	command := UpsertRoleCommand{
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersCreate},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["name"]; !ok {
		t.Errorf("expected name error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandBadScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Bad Scope",
		Scope:       permissions.Scope("INVALID"),
		Permissions: []string{},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["scope"]; !ok {
		t.Errorf("expected scope error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandAppPermissionOnGroupScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Mismatch",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.AppUsersCreate},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["permissions"]; !ok {
		t.Errorf("expected permissions error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandUnknownPermission(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Unknown",
		Scope:       permissions.ScopeApp,
		Permissions: []string{"not.a.real.permission"},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["permissions"]; !ok {
		t.Errorf("expected permissions error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandDuplicatePermission(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Duplicate",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersCreate, permissions.AppUsersCreate},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["permissions"]; !ok {
		t.Errorf("expected permissions error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandGrantsOnGroupScopeValid(t *testing.T) {
	command := UpsertRoleCommand{
		Name:           "Restricted Group Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{1, 2},
		TagGrants:      []uint{3},
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandGrantsRejectedOnAppScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:           "App Role With Grants",
		Scope:          permissions.ScopeApp,
		Permissions:    []string{permissions.AppUsersRead},
		CategoryGrants: []uint{1},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["grants"]; !ok {
		t.Errorf("expected grants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandDuplicateCategoryGrant(t *testing.T) {
	command := UpsertRoleCommand{
		Name:           "Dup Category Grant",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{1, 1},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["categoryGrants"]; !ok {
		t.Errorf("expected categoryGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandDuplicateTagGrant(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Dup Tag Grant",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReceiptsRead},
		TagGrants:   []uint{5, 5},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["tagGrants"]; !ok {
		t.Errorf("expected tagGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandPaidByGrantsValidOnGroupScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:                   "Paid-By Group Role",
		Scope:                  permissions.ScopeGroup,
		Permissions:            []string{permissions.GroupReceiptsRead},
		PaidByUserGrants:       []uint{1, 2},
		IncludeOwnPaidReceipts: true,
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandPaidByGrantsRejectedOnAppScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:             "App Role With Paid-By",
		Scope:            permissions.ScopeApp,
		Permissions:      []string{permissions.AppUsersRead},
		PaidByUserGrants: []uint{1},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["grants"]; !ok {
		t.Errorf("expected grants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandIncludeOwnRejectedOnAppScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:                   "App Role With Include Own",
		Scope:                  permissions.ScopeApp,
		Permissions:            []string{permissions.AppUsersRead},
		IncludeOwnPaidReceipts: true,
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["grants"]; !ok {
		t.Errorf("expected grants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandSeesAllMembersValidOnGroupScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:           "Supervisor Group Role",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		SeesAllMembers: true,
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandSeesAllMembersRejectedOnAppScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:           "App Role With SeesAll",
		Scope:          permissions.ScopeApp,
		Permissions:    []string{permissions.AppUsersRead},
		SeesAllMembers: true,
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["grants"]; !ok {
		t.Errorf("expected grants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandDuplicatePaidByGrant(t *testing.T) {
	command := UpsertRoleCommand{
		Name:             "Dup Paid-By Grant",
		Scope:            permissions.ScopeGroup,
		Permissions:      []string{permissions.GroupReceiptsRead},
		PaidByUserGrants: []uint{7, 7},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["paidByUserGrants"]; !ok {
		t.Errorf("expected paidByUserGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandReportTemplateGrantsValidOnGroupScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Report Grant Group Role",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{"read", "generate"}},
			{ReportTemplateId: 2, Permissions: []string{"read"}},
		},
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandReportTemplateGrantsRejectedOnAppScope(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "App Role With Report Grants",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{"read"}},
		},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["reportTemplateGrants"]; !ok {
		t.Errorf("expected reportTemplateGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandDuplicateReportTemplateGrant(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Dup Report Grant",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{"read"}},
			{ReportTemplateId: 1, Permissions: []string{"generate"}},
		},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["reportTemplateGrants"]; !ok {
		t.Errorf("expected reportTemplateGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandReportTemplateGrantUnknownAction(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Bad Action Report Grant",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{"read", "explode"}},
		},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["reportTemplateGrants"]; !ok {
		t.Errorf("expected reportTemplateGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandReportTemplateGrantCreateActionRejected(t *testing.T) {
	// "create" is not scopable per template — there is no existing template to
	// scope it to — so it must be rejected from the matrix.
	command := UpsertRoleCommand{
		Name:        "Create Action Report Grant",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{"create"}},
		},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["reportTemplateGrants"]; !ok {
		t.Errorf("expected reportTemplateGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandReportTemplateGrantEmptyPermissions(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Empty Permissions Report Grant",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{}},
		},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["reportTemplateGrants"]; !ok {
		t.Errorf("expected reportTemplateGrants error, got %+v", vErr.Errors)
	}
}

func TestUpsertRoleCommandDuplicateActionWithinReportTemplateGrant(t *testing.T) {
	command := UpsertRoleCommand{
		Name:        "Dup Action Report Grant",
		Scope:       permissions.ScopeGroup,
		Permissions: []string{permissions.GroupReportsRead},
		ReportTemplateGrants: []ReportTemplateGrantCommand{
			{ReportTemplateId: 1, Permissions: []string{"read", "read"}},
		},
	}

	vErr := command.Validate()
	if _, ok := vErr.Errors["reportTemplateGrants"]; !ok {
		t.Errorf("expected reportTemplateGrants error, got %+v", vErr.Errors)
	}
}
