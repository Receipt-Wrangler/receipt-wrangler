package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type PagedRoleRequestCommand struct {
	PagedRequestCommand
	Filter RoleFilter `json:"filter"`
}

// RoleFilter narrows the paged role list to a single scope. An empty Scope
// returns both app- and group-scoped roles.
type RoleFilter struct {
	Scope permissions.Scope `json:"scope"`
}

func (command *PagedRoleRequestCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
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

func (command *PagedRoleRequestCommand) Validate() structs.ValidatorError {
	vErrs := command.PagedRequestCommand.Validate()

	if command.Filter.Scope != "" &&
		command.Filter.Scope != permissions.ScopeApp &&
		command.Filter.Scope != permissions.ScopeGroup {
		vErrs.Errors["scope"] = "Scope must be empty, APP, or GROUP"
	}

	return vErrs
}
