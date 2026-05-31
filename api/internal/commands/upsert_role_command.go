package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type UpsertRoleCommand struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Scope       permissions.Scope `json:"scope"`
	Permissions []string          `json:"permissions"`
}

func (command *UpsertRoleCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	return nil
}

func (command *UpsertRoleCommand) Validate() structs.ValidatorError {
	errors := make(map[string]string)
	vErr := structs.ValidatorError{}

	if len(strings.TrimSpace(command.Name)) == 0 {
		errors["name"] = "Name is required"
	}

	if command.Scope != permissions.ScopeApp && command.Scope != permissions.ScopeGroup {
		errors["scope"] = "Scope must be either APP or GROUP"
	}

	seen := make(map[string]bool)
	for _, permission := range command.Permissions {
		if seen[permission] {
			errors["permissions"] = fmt.Sprintf("Duplicate permission: %s", permission)
			break
		}
		seen[permission] = true

		descriptor, ok := permissions.Get(permission)
		if !ok {
			errors["permissions"] = fmt.Sprintf("Unknown permission: %s", permission)
			break
		}

		if descriptor.Scope != command.Scope {
			errors["permissions"] = fmt.Sprintf("Permission %s does not belong to scope %s", permission, command.Scope)
			break
		}
	}

	vErr.Errors = errors
	return vErr
}
