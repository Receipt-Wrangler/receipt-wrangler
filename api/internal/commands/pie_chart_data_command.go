package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type PieChartDataCommand struct {
	ChartGrouping models.ChartGrouping      `json:"chartGrouping"`
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

	initReceiptFilterValues(&command.Filter)
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
