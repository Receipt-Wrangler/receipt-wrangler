package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type UpsertCustomFieldCommand struct {
	Name        string                           `json:"name"`
	Type        models.CustomFieldType           `json:"type"`
	Description string                           `json:"description"`
	Options     []UpsertCustomFieldOptionCommand `json:"options"`
}

func (command *UpsertCustomFieldCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	if command.Type != models.SELECT && len(command.Options) > 0 {
		command.Options = []UpsertCustomFieldOptionCommand{}
	}

	return nil
}

func (command *UpsertCustomFieldCommand) Validate() structs.ValidatorError {
	errors := make(map[string]string)
	vErr := structs.ValidatorError{}

	if len(command.Name) == 0 {
		errors["name"] = "Name is required"
	}

	if len(command.Type) == 0 {
		errors["type"] = "Type is required"
	}

	if command.Type == models.SELECT && len(command.Options) == 0 {
		errors["options"] = "Options are required"
	}

	vErr.Errors = errors
	return vErr
}

// ValidateUpdate applies the rules that only exist for an edit, given the custom
// field as it is currently stored. It is deliberately pure -- the caller loads
// the existing record (options included) and passes it in -- so the update rules
// live beside Validate rather than in the handler.
//
// Two invariants are enforced here:
//
//   - The type is immutable. A CustomFieldValue stores its data in one of five
//     type-specific columns, so re-typing a field would mis-column every value
//     already recorded against it.
//   - An option must belong to this custom field. Options are matched by id so an
//     edit can rename them in place; an id from another field would otherwise let
//     a caller rewrite an unrelated field's option.
//
// Removing an option is not an error -- the repository leaves any option the
// command omits untouched, because CustomFieldValue.SelectValue holds an option
// id that a delete would orphan.
func (command *UpsertCustomFieldCommand) ValidateUpdate(existing models.CustomField) structs.ValidatorError {
	errors := make(map[string]string)
	vErr := structs.ValidatorError{}

	if command.Type != existing.Type {
		errors["type"] = "Type cannot be changed"
	}

	existingOptionIds := make(map[uint]bool, len(existing.Options))
	for _, option := range existing.Options {
		existingOptionIds[option.ID] = true
	}

	for _, optionCommand := range command.Options {
		if optionCommand.Id != 0 && !existingOptionIds[optionCommand.Id] {
			errors["options"] = "One or more options do not belong to this custom field"
			break
		}
	}

	vErr.Errors = errors
	return vErr
}
