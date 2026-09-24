package commands

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type PagedApiKeyRequestCommand struct {
	PagedRequestCommand
	ApiKeyFilter `json:"filter"`
}

func (command *PagedApiKeyRequestCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
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

func (command *PagedApiKeyRequestCommand) Validate(r *http.Request) structs.ValidatorError {
	vErrs := command.PagedRequestCommand.Validate()

	// Reject anything that is not exactly MINE or ALL. Same fail-open reasoning
	// as PagedGroupRequestCommand: only ALL is permission-gated in the handler
	// and only MINE is scoped to the caller in the repository, so any other
	// value (empty or bogus) would leak every user's API keys. The repository
	// also defaults unknown values to MINE scoping as defense in depth.
	if command.ApiKeyFilter.AssociatedApiKeys != ASSOCIATED_API_KEYS_MINE &&
		command.ApiKeyFilter.AssociatedApiKeys != ASSOCIATED_API_KEYS_ALL {
		vErrs.Errors["associatedApiKeys"] = "Associated API keys must be MINE or ALL"
	}

	return vErrs
}

type ApiKeyFilter struct {
	AssociatedApiKeys AssociatedApiKeys `json:"associatedApiKeys"`
}

type AssociatedApiKeys string

const (
	ASSOCIATED_API_KEYS_MINE AssociatedApiKeys = "MINE"
	ASSOCIATED_API_KEYS_ALL  AssociatedApiKeys = "ALL"
)

func (associatedApiKeys *AssociatedApiKeys) Scan(value string) error {
	*associatedApiKeys = AssociatedApiKeys(value)
	return nil
}

func (associatedApiKeys AssociatedApiKeys) Value() (driver.Value, error) {
	if associatedApiKeys != ASSOCIATED_API_KEYS_MINE && associatedApiKeys != ASSOCIATED_API_KEYS_ALL {
		return nil, errors.New("invalid associatedApiKeys")
	}
	return string(associatedApiKeys), nil
}
