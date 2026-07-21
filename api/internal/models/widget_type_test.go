package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestWidgetType_Value(t *testing.T) {
	valid := []WidgetType{GROUP_SUMMARY, FILTERED_RECEIPTS, GROUP_ACTIVITY, PIE_CHART, REPORT}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestWidgetType_Value_Invalid(t *testing.T) {
	// This type has no empty-string exception, so an empty value is invalid too.
	assertValuerInvalid(t, "empty", WidgetType(""))
	assertValuerInvalid(t, "bogus", WidgetType("bogus"))
}

func TestWidgetType_Scan(t *testing.T) {
	var widgetType WidgetType
	err := widgetType.Scan("PIE_CHART")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if widgetType != PIE_CHART {
		utils.PrintTestError(t, widgetType, PIE_CHART)
	}
}
