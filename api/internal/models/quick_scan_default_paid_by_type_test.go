package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestQuickScanDefaultPaidByType_Value(t *testing.T) {
	valid := []QuickScanDefaultPaidByType{QUICK_SCAN_PAID_BY_UPLOADER, QUICK_SCAN_PAID_BY_USER}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted and normalized to "".
	assertValuerValid(t, "empty", QuickScanDefaultPaidByType(""), "")
}

func TestQuickScanDefaultPaidByType_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", QuickScanDefaultPaidByType("bogus"))
}

func TestQuickScanDefaultPaidByType_Scan(t *testing.T) {
	var paidByType QuickScanDefaultPaidByType
	err := paidByType.Scan("UPLOADER")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if paidByType != QUICK_SCAN_PAID_BY_UPLOADER {
		utils.PrintTestError(t, paidByType, QUICK_SCAN_PAID_BY_UPLOADER)
	}
}

func TestQuickScanDefaultPaidByTypes(t *testing.T) {
	types := QuickScanDefaultPaidByTypes()

	expected := []interface{}{QUICK_SCAN_PAID_BY_UPLOADER, QUICK_SCAN_PAID_BY_USER}
	if len(types) != len(expected) {
		utils.PrintTestError(t, len(types), len(expected))
	}

	for i, v := range expected {
		if types[i] != v {
			utils.PrintTestError(t, types[i], v)
		}
	}
}
