package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestCustomFieldType_Value(t *testing.T) {
	valid := []CustomFieldType{TEXT, DATE, SELECT, CURRENCY, BOOLEAN}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}

	// An empty value is accepted and normalized to "".
	assertValuerValid(t, "empty", CustomFieldType(""), "")
}

func TestCustomFieldType_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "bogus", CustomFieldType("bogus"))
}

func TestCustomFieldType_Scan(t *testing.T) {
	var fieldType CustomFieldType
	err := fieldType.Scan("TEXT")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if fieldType != TEXT {
		utils.PrintTestError(t, fieldType, TEXT)
	}
}
