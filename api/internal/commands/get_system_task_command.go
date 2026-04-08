package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type GetSystemTaskCommand struct {
	PagedRequestCommand
	AssociatedEntityId   uint                        `json:"associatedEntityId"`
	AssociatedEntityType models.AssociatedEntityType `json:"associatedEntityType"`
	Filter               SystemTaskPagedRequestFilter `json:"filter"`
}

func (command *GetSystemTaskCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	if command.Filter.Type.Value == nil || command.Filter.Type.Value == "" {
		command.Filter.Type.Value = make([]interface{}, 0)
	}

	if command.Filter.Status.Value == nil || command.Filter.Status.Value == "" {
		command.Filter.Status.Value = make([]interface{}, 0)
	}

	if command.Filter.RanByUserId.Value == nil || command.Filter.RanByUserId.Value == "" {
		command.Filter.RanByUserId.Value = make([]interface{}, 0)
	}

	if command.Filter.StartedAt.Value == nil {
		command.Filter.StartedAt.Value = ""
	}

	if command.Filter.EndedAt.Value == nil {
		command.Filter.EndedAt.Value = ""
	}

	return nil
}

func (command *GetSystemTaskCommand) Validate() structs.ValidatorError {
	vErrs := command.PagedRequestCommand.Validate()

	return vErrs
}
