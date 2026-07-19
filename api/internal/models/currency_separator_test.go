package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestCurrencySeparator_Value(t *testing.T) {
	valid := []CurrencySeparator{COMMA, DOT}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted and normalized to "".
	assertValuerValid(t, "empty", CurrencySeparator(""), "")
}

func TestCurrencySeparator_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", CurrencySeparator("bogus"))
}

func TestCurrencySeparator_Value_InvalidErrorMessage(t *testing.T) {
	_, err := CurrencySeparator("bogus").Value()
	if err == nil {
		utils.PrintTestError(t, err, "an error")
		return
	}

	expected := "invalid currency separator"
	if err.Error() != expected {
		utils.PrintTestError(t, err.Error(), expected)
	}
}

func TestCurrencySeparator_Scan(t *testing.T) {
	var separator CurrencySeparator
	err := separator.Scan(",")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if separator != COMMA {
		utils.PrintTestError(t, separator, COMMA)
	}
}
