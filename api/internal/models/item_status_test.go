package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestItemStatus_Value(t *testing.T) {
	valid := []ItemStatus{ITEM_OPEN, ITEM_RESOLVED, ITEM_DRAFT}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestItemStatus_Value_Invalid(t *testing.T) {
	// This type has no empty-string exception, so an empty value is invalid too.
	assertValuerInvalid(t, "empty", ItemStatus(""))
	assertValuerInvalid(t, "bogus", ItemStatus("bogus"))
}

func TestItemStatus_Scan(t *testing.T) {
	var status ItemStatus
	err := status.Scan("OPEN")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if status != ITEM_OPEN {
		utils.PrintTestError(t, status, ITEM_OPEN)
	}
}
