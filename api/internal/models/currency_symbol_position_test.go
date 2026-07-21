package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestCurrencySymbolPosition_Value(t *testing.T) {
	valid := []CurrencySymbolPosition{START, END}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted and normalized to "".
	assertValuerValid(t, "empty", CurrencySymbolPosition(""), "")
}

func TestCurrencySymbolPosition_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", CurrencySymbolPosition("bogus"))
}

func TestCurrencySymbolPosition_Scan(t *testing.T) {
	var position CurrencySymbolPosition
	err := position.Scan("START")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if position != START {
		utils.PrintTestError(t, position, START)
	}
}
