package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// BulkReceiptDeleteCommand is the body of POST /receipt/bulkDelete. Whether the
// caller may delete each receipt is decided per receipt by the service, not here.
type BulkReceiptDeleteCommand struct {
	ReceiptIds []uint `json:"receiptIds"`
}

func (command *BulkReceiptDeleteCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &command)
}

// Validate rejects an empty list and duplicate ids. A duplicate would be deleted
// once and then reported as refused, which reads as a failure that never happened.
func (command *BulkReceiptDeleteCommand) Validate() structs.ValidatorError {
	errors := structs.ValidatorError{Errors: make(map[string]string)}

	if len(command.ReceiptIds) == 0 {
		errors.Errors["receiptIds"] = "At least one receipt id is required"
	} else if hasDuplicateIds(command.ReceiptIds) {
		errors.Errors["receiptIds"] = "Duplicate receipt ids are not allowed"
	}

	return errors
}

func (command *BulkReceiptDeleteCommand) LoadDataFromRequestAndValidate(w http.ResponseWriter, r *http.Request) (structs.ValidatorError, error) {
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		return structs.ValidatorError{}, err
	}

	return command.Validate(), nil
}
