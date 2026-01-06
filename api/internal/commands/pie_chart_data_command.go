package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type PieChartDataCommand struct {
	ChartGrouping models.ChartGrouping  `json:"chartGrouping"`
	Filter        ReceiptPagedRequestFilter `json:"filter"`
}

func (command *PieChartDataCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	// Initialize filter values if nil
	if command.Filter.Amount.Value == nil || command.Filter.Amount.Value == "" {
		command.Filter.Amount.Value = float64(0)
	}

	if command.Filter.PaidBy.Value == nil || command.Filter.PaidBy.Value == "" {
		command.Filter.PaidBy.Value = make([]interface{}, 0)
	}

	if command.Filter.Categories.Value == nil || command.Filter.Categories.Value == "" {
		command.Filter.Categories.Value = make([]interface{}, 0)
	}

	if command.Filter.Tags.Value == nil || command.Filter.Tags.Value == "" {
		command.Filter.Tags.Value = make([]interface{}, 0)
	}

	if command.Filter.Status.Value == nil || command.Filter.Status.Value == "" {
		command.Filter.Status.Value = make([]interface{}, 0)
	}

	if command.Filter.CreatedAt.Value == nil {
		command.Filter.CreatedAt.Value = ""
	}

	if command.Filter.Date.Value == nil {
		command.Filter.Date.Value = ""
	}

	if command.Filter.ResolvedDate.Value == nil {
		command.Filter.ResolvedDate.Value = ""
	}

	return nil
}

func (command *PieChartDataCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{}
	errorMap := make(map[string]string)

	_, err := command.ChartGrouping.Value()
	if err != nil {
		errorMap["chartGrouping"] = "Invalid chart grouping. Must be CATEGORIES, TAGS, or PAIDBY"
	}

	vErr.Errors = errorMap
	return vErr
}
