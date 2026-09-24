package commands

import (
	"testing"

	"receipt-wrangler/api/internal/utils"
)

func basePagedApiKeyCommand(filter AssociatedApiKeys) PagedApiKeyRequestCommand {
	command := PagedApiKeyRequestCommand{}
	command.Page = 1
	command.PageSize = 10
	command.OrderBy = "name"
	command.SortDirection = ASCENDING
	command.ApiKeyFilter.AssociatedApiKeys = filter
	return command
}

// The associatedApiKeys filter must be exactly MINE or ALL — same fail-open
// reasoning as the group filter: a value that is neither would leak every user's
// API keys. Regression guard.
func TestPagedApiKeyRequestCommand_ValidateAssociatedApiKeys(t *testing.T) {
	tests := map[string]struct {
		filter    AssociatedApiKeys
		wantError bool
	}{
		"MINE is accepted":     {ASSOCIATED_API_KEYS_MINE, false},
		"ALL is accepted":      {ASSOCIATED_API_KEYS_ALL, false},
		"empty is rejected":    {"", true},
		"bogus value rejected": {"x", true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			command := basePagedApiKeyCommand(tt.filter)
			vErr := command.Validate(nil)
			_, hasErr := vErr.Errors["associatedApiKeys"]
			if hasErr != tt.wantError {
				utils.PrintTestError(t, hasErr, tt.wantError)
			}
		})
	}
}
