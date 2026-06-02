package commands

import (
	"receipt-wrangler/api/internal/permissions"
	"testing"
)

func validPagedRoleRequestCommand() PagedRoleRequestCommand {
	return PagedRoleRequestCommand{
		PagedRequestCommand: PagedRequestCommand{
			Page:          1,
			PageSize:      50,
			OrderBy:       "name",
			SortDirection: ASCENDING,
		},
	}
}

func TestPagedRoleRequestCommandValidWithoutScope(t *testing.T) {
	command := validPagedRoleRequestCommand()

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", vErr.Errors)
	}
}

func TestPagedRoleRequestCommandValidWithScope(t *testing.T) {
	for _, scope := range []permissions.Scope{permissions.ScopeApp, permissions.ScopeGroup} {
		command := validPagedRoleRequestCommand()
		command.Filter.Scope = scope

		vErr := command.Validate()
		if len(vErr.Errors) > 0 {
			t.Errorf("expected no errors for scope %q, got %+v", scope, vErr.Errors)
		}
	}
}

func TestPagedRoleRequestCommandInvalidScope(t *testing.T) {
	command := validPagedRoleRequestCommand()
	command.Filter.Scope = "NONSENSE"

	vErr := command.Validate()
	if _, ok := vErr.Errors["scope"]; !ok {
		t.Errorf("expected a scope error, got %+v", vErr.Errors)
	}
}

func TestPagedRoleRequestCommandInheritsPagedValidation(t *testing.T) {
	command := validPagedRoleRequestCommand()
	command.Page = 0

	vErr := command.Validate()
	if _, ok := vErr.Errors["page"]; !ok {
		t.Errorf("expected a page error from the embedded command, got %+v", vErr.Errors)
	}
}
