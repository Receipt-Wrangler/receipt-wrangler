package structs

import (
	"receipt-wrangler/api/internal/models"

	"github.com/shopspring/decimal"
)

// ReceiptSummaryCustomFieldTotal is one currency custom field's total within a row.
// Name travels with the total so the client needs no second lookup — and so a field
// deleted between two requests simply stops appearing rather than rendering nameless.
type ReceiptSummaryCustomFieldTotal struct {
	CustomFieldId uint            `json:"customFieldId"`
	Name          string          `json:"name"`
	Total         decimal.Decimal `json:"total"`
}

// ReceiptSummaryRow is one line of the summary block.
//
// Status is "" on the overall row and a real status on each breakdown row. The row
// carries no display label: the wording is presentation, and both clients already own
// it (the desktop's formatStatus).
//
// Money is decimal.Decimal, which marshals to a QUOTED STRING — matching Receipt.amount,
// CustomFieldValue.currencyValue and /user/amountOwedForUser. PieChartDataPoint's float64
// is the outlier, forced on it by the charting library.
type ReceiptSummaryRow struct {
	Status            models.ReceiptStatus             `json:"status"`
	ReceiptCount      int64                            `json:"receiptCount"`
	Total             decimal.Decimal                  `json:"total"`
	CustomFieldTotals []ReceiptSummaryCustomFieldTotal `json:"customFieldTotals"`
}

// ReceiptSummary is the whole block.
//
// Enabled false comes back at a normal 200 with zeroed rows rather than a 404 or 400.
// A client whose cached group settings are stale then renders nothing, which is the
// intended off state, instead of surfacing an error toast for a feature the group
// simply has not turned on.
//
// Statuses and every CustomFieldTotals must be non-nil so they serialize as [] rather
// than null — the generated Dart deserializer has no null guard.
// Position rides on the response rather than being read from the client's cached
// GroupReceiptSettings, for the same reason Enabled does: the cache is stale the moment
// an admin changes the configuration, and placement is configuration. It is always one of
// the two real values — never "" — see ReceiptSummaryPosition.OrDefault.
type ReceiptSummary struct {
	Enabled              bool                          `json:"enabled"`
	ConfigurationGroupId uint                          `json:"configurationGroupId"`
	Position             models.ReceiptSummaryPosition `json:"position"`
	Overall              ReceiptSummaryRow             `json:"overall"`
	Statuses             []ReceiptSummaryRow           `json:"statuses"`
}
