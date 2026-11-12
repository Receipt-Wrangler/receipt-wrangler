package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// MaxClaimReceiptIds bounds one add or remove, for the same reason as
// MaxBulkReceiptDeleteIds: the ids are looked up with `id IN ?`, which must stay
// under every supported engine's bind-parameter limit.
const MaxClaimReceiptIds = 500

// ClaimReceiptsCommand is the body of addReceiptsToClaim and
// removeReceiptsFromClaim. Whether each receipt may be added or removed is
// decided per receipt by the service, not here.
type ClaimReceiptsCommand struct {
	ReceiptIds []uint `json:"receiptIds"`
}

func (command *ClaimReceiptsCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, command)
}

// Validate rejects an empty list, one over the cap, and duplicates, which would
// otherwise be reported once as done and once as refused.
func (command *ClaimReceiptsCommand) Validate() structs.ValidatorError {
	errors := structs.ValidatorError{Errors: make(map[string]string)}

	if len(command.ReceiptIds) == 0 {
		errors.Errors["receiptIds"] = "At least one receipt id is required"
	} else if len(command.ReceiptIds) > MaxClaimReceiptIds {
		errors.Errors["receiptIds"] = fmt.Sprintf("At most %d receipts can be changed at once", MaxClaimReceiptIds)
	} else if hasDuplicateIds(command.ReceiptIds) {
		errors.Errors["receiptIds"] = "Duplicate receipt ids are not allowed"
	}

	return errors
}
