package commands

import (
	"testing"

	"receipt-wrangler/api/internal/utils"
)

func basePagedGroupCommand(filter AssociatedGroup) PagedGroupRequestCommand {
	command := PagedGroupRequestCommand{}
	command.Page = 1
	command.PageSize = 10
	command.OrderBy = "name"
	command.SortDirection = ASCENDING
	command.GroupFilter.AssociatedGroup = filter
	return command
}

// The associatedGroup filter must be exactly MINE or ALL. Anything else (empty
// or a bogus value) is rejected, because the handler only permission-gates the
// ALL case and the repository only scopes the MINE case — a value that is
// neither would otherwise fall through both and leak every group. Regression
// guard for the fail-open filter bug.
func TestPagedGroupRequestCommand_ValidateAssociatedGroup(t *testing.T) {
	tests := map[string]struct {
		filter    AssociatedGroup
		wantError bool
	}{
		"MINE is accepted":        {ASSOCIATED_GROUP_MINE, false},
		"ALL is accepted":         {ASSOCIATED_GROUP_ALL, false},
		"empty is rejected":       {"", true},
		"bogus value rejected":    {"x", true},
		"lowercase mine rejected": {"mine", true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			command := basePagedGroupCommand(tt.filter)
			vErr := command.Validate(nil)
			_, hasErr := vErr.Errors["associatedGroup"]
			if hasErr != tt.wantError {
				utils.PrintTestError(t, hasErr, tt.wantError)
			}
		})
	}
}
