package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestReceiptSummaryPosition_Value(t *testing.T) {
	valid := []ReceiptSummaryPosition{
		RECEIPT_SUMMARY_POSITION_TOP,
		RECEIPT_SUMMARY_POSITION_BOTTOM,
	}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted at the DB boundary and normalized to "". It is
	// OrDefault, not the Valuer, that keeps "" off the wire.
	assertValuerValid(t, "empty", ReceiptSummaryPosition(""), "")
}

func TestReceiptSummaryPosition_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", ReceiptSummaryPosition("bogus"))
	assertValuerInvalid(t, "lowercase", ReceiptSummaryPosition("top"))
}

func TestReceiptSummaryPosition_Scan(t *testing.T) {
	var position ReceiptSummaryPosition
	err := position.Scan("TOP")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if position != RECEIPT_SUMMARY_POSITION_TOP {
		utils.PrintTestError(t, position, RECEIPT_SUMMARY_POSITION_TOP)
	}
}

// OrDefault is what guarantees a client never receives "". Every value that is not
// TOP resolves to BOTTOM, which is where the summary rendered before the setting
// existed — so a row that predates the column, and a row holding a value this build
// does not know, both degrade to the historical behaviour rather than to a parse error.
func TestReceiptSummaryPosition_OrDefault(t *testing.T) {
	cases := map[ReceiptSummaryPosition]ReceiptSummaryPosition{
		RECEIPT_SUMMARY_POSITION_TOP:    RECEIPT_SUMMARY_POSITION_TOP,
		RECEIPT_SUMMARY_POSITION_BOTTOM: RECEIPT_SUMMARY_POSITION_BOTTOM,
		"":                              RECEIPT_SUMMARY_POSITION_BOTTOM,
		"bogus":                         RECEIPT_SUMMARY_POSITION_BOTTOM,
	}

	for in, expected := range cases {
		got := in.OrDefault()
		if got != expected {
			utils.PrintTestError(t, got, expected)
		}
	}
}

func TestReceiptSummaryPositions(t *testing.T) {
	positions := ReceiptSummaryPositions()
	if len(positions) != 2 {
		utils.PrintTestError(t, len(positions), 2)
	}
	if positions[0] != RECEIPT_SUMMARY_POSITION_TOP || positions[1] != RECEIPT_SUMMARY_POSITION_BOTTOM {
		utils.PrintTestError(t, positions, "[TOP BOTTOM]")
	}
}
