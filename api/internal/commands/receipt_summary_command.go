package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// ReceiptSummaryCommand asks for the block of totals rendered under the receipts
// table: a receipt count and amount total over the WHOLE filtered result set, then
// the same figures per configured status.
//
// It carries the filter but NOT the configuration. Which statuses break out and
// which currency fields are totalled come from the group's GroupReceiptSettings,
// resolved server-side — the configuration belongs to the group and applies to
// every member, so a client must not be able to add a column or opt out of one.
type ReceiptSummaryCommand struct {
	// ConfigurationGroupId names the group whose settings drive the breakdown. It
	// exists for the synthetic "All" group, which spans every group the caller
	// belongs to and has no meaningful settings row of its own, so the client picks
	// which member group's configuration to apply. nil means "the group in the URL",
	// which is the answer for every real group.
	//
	// The DATA is always the URL group's filtered set. This only chooses the shape
	// of the breakdown.
	ConfigurationGroupId *uint                     `json:"configurationGroupId"`
	Filter               ReceiptPagedRequestFilter `json:"filter"`
}

func (command *ReceiptSummaryCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	// Mandatory, exactly as in PieChartDataCommand: BuildGormFilterQuery's type
	// assertions and IntersectReceiptFilterWithGrants both assume the non-nil
	// defaults this seeds. Skipping it turns an omitted filter key into a panic.
	initReceiptFilterValues(&command.Filter)
	return nil
}

// Validate has nothing to check. The filter is seeded by LoadDataFromRequest and the
// configuration group is an authorization question, answered in the service where the
// caller's permissions are available — not a validation one.
func (command *ReceiptSummaryCommand) Validate() structs.ValidatorError {
	return structs.ValidatorError{Errors: make(map[string]string)}
}
