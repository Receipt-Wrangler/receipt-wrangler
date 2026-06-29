package commands

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type PagedGroupRequestCommand struct {
	PagedRequestCommand
	GroupFilter `json:"filter"`
}

func (command *PagedGroupRequestCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
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

func (command *PagedGroupRequestCommand) Validate(r *http.Request) structs.ValidatorError {
	vErrs := command.PagedRequestCommand.Validate()

	// Authorization for the ALL filter (listing every group) is enforced in the
	// GetPagedGroups handler via the app.groups.read permission, resolved from the
	// database rather than the JWT.
	if command.GroupFilter.AssociatedGroup == "" {
		vErrs.Errors["associatedGroup"] = "Associated group is required"
	}

	return vErrs
}

type GroupFilter struct {
	AssociatedGroup AssociatedGroup `json:"associatedGroup"`
}

type AssociatedGroup string

const (
	ASSOCIATED_GROUP_MINE AssociatedGroup = "MINE"
	ASSOCIATED_GROUP_ALL  AssociatedGroup = "ALL"
)

func (associatedGroup *AssociatedGroup) Scan(value string) error {
	*associatedGroup = AssociatedGroup(value)
	return nil
}

func (associatedGroup AssociatedGroup) Value() (driver.Value, error) {
	if associatedGroup != ASSOCIATED_GROUP_MINE && associatedGroup != ASSOCIATED_GROUP_ALL {
		return nil, errors.New("invalid associatedGroup")
	}
	return string(associatedGroup), nil
}
