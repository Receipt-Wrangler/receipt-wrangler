package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type UpsertRoleCommand struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Scope       permissions.Scope `json:"scope"`
	Permissions []string          `json:"permissions"`
	// CategoryGrants / TagGrants restrict which categories/tags members of a
	// group role may read and use. They are only valid on GROUP-scoped roles. An
	// empty set means "unrestricted" (see every category/tag) — restriction is
	// opt-in. The values are category/tag ids; their existence is validated
	// against the database in the role service.
	CategoryGrants []uint `json:"categoryGrants"`
	TagGrants      []uint `json:"tagGrants"`
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

	if len(command.Name) == 0 {
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

	// Category/tag grants are a group-role concept (they slice the global pool
	// per group role); they make no sense on an app role.
	if command.Scope == permissions.ScopeApp && (len(command.CategoryGrants) > 0 || len(command.TagGrants) > 0) {
		errors["grants"] = "Category and tag grants are only valid on group roles"
	}

	if hasDuplicateUint(command.CategoryGrants) {
		errors["categoryGrants"] = "Duplicate category grant"
	}

	if hasDuplicateUint(command.TagGrants) {
		errors["tagGrants"] = "Duplicate tag grant"
	}

	vErr.Errors = errors
	return vErr
}

// hasDuplicateUint reports whether ids contains the same value more than once.
func hasDuplicateUint(ids []uint) bool {
	seen := make(map[uint]bool, len(ids))
	for _, id := range ids {
		if seen[id] {
			return true
		}
		seen[id] = true
	}
	return false
}
