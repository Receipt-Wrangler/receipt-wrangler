package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestChartGrouping_Value(t *testing.T) {
	valid := []ChartGrouping{CHART_GROUPING_CATEGORIES, CHART_GROUPING_TAGS, CHART_GROUPING_PAIDBY}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestChartGrouping_Value_Invalid(t *testing.T) {
	// This type has no empty-string exception, so an empty value is invalid too.
	assertValuerInvalid(t, "empty", ChartGrouping(""))
	assertValuerInvalid(t, "bogus", ChartGrouping("bogus"))
}

func TestChartGrouping_Scan(t *testing.T) {
	var grouping ChartGrouping
	err := grouping.Scan("CATEGORIES")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if grouping != CHART_GROUPING_CATEGORIES {
		utils.PrintTestError(t, grouping, CHART_GROUPING_CATEGORIES)
	}
}
