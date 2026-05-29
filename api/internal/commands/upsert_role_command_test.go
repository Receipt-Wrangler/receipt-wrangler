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
