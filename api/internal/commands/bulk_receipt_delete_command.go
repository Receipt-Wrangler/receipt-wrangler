package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// MaxBulkReceiptDeleteIds bounds one request. The ids are looked up in a single
// `id IN ?` query, so the bound keeps it under every supported engine's
// bind-parameter limit (999 on the oldest SQLite builds), and it caps how many
// deletes one request runs. The desktop sends at most one table page (100).
const MaxBulkReceiptDeleteIds = 500

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

// Validate rejects an empty list, one longer than MaxBulkReceiptDeleteIds, and
// duplicate ids. A duplicate would be deleted once and then reported as refused,
// which reads as a failure that never happened.
func (command *BulkReceiptDeleteCommand) Validate() structs.ValidatorError {
	errors := structs.ValidatorError{Errors: make(map[string]string)}

	if len(command.ReceiptIds) == 0 {
		errors.Errors["receiptIds"] = "At least one receipt id is required"
	} else if len(command.ReceiptIds) > MaxBulkReceiptDeleteIds {
		errors.Errors["receiptIds"] = fmt.Sprintf("At most %d receipts can be deleted at once", MaxBulkReceiptDeleteIds)
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
