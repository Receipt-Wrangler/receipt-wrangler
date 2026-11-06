package commands

import (
	"encoding/json"
	"testing"
)

// decodeCustomField decodes a receipt body carrying one custom field value whose
// fields are given as raw JSON, the way the handlers decode a request.
func decodeCustomField(t *testing.T, fieldJson string) (UpsertCustomFieldValueCommand, error) {
	t.Helper()
	var command UpsertReceiptCommand
	err := json.Unmarshal([]byte(`{"name":"R","customFields":[`+fieldJson+`]}`), &command)
	if err != nil {
		return UpsertCustomFieldValueCommand{}, err
	}
	if len(command.CustomFields) != 1 {
		t.Fatalf("decoded %d custom fields, want 1", len(command.CustomFields))
	}
	return command.CustomFields[0], nil
}

// The desktop's masked currency input sends "" for a field the user emptied. It
// must save like an untouched field (nil) instead of failing the whole decode.
func TestUpsertCustomFieldValueCommand_BlankCurrencyValueIsNoValue(t *testing.T) {
	tests := map[string]string{
		"empty string":    `{"customFieldId":1,"currencyValue":""}`,
		"whitespace only": `{"customFieldId":1,"currencyValue":"   "}`,
		"null":            `{"customFieldId":1,"currencyValue":null}`,
		"key absent":      `{"customFieldId":1}`,
	}

	for name, fieldJson := range tests {
		t.Run(name, func(t *testing.T) {
			field, err := decodeCustomField(t, fieldJson)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if field.CurrencyValue != nil {
				t.Errorf("CurrencyValue = %s, want nil", field.CurrencyValue)
			}
		})
	}
}

func TestUpsertCustomFieldValueCommand_CurrencyValueDecodes(t *testing.T) {
	tests := map[string]struct {
		fieldJson string
		want      string
	}{
		"numeric string": {fieldJson: `{"currencyValue":"12.34"}`, want: "12.34"},
		"number":         {fieldJson: `{"currencyValue":12.34}`, want: "12.34"},
		"negative":       {fieldJson: `{"currencyValue":"-5.5"}`, want: "-5.5"},
		// An entered zero is a value, not a blank.
		"zero": {fieldJson: `{"currencyValue":"0"}`, want: "0"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			field, err := decodeCustomField(t, test.fieldJson)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if field.CurrencyValue == nil {
				t.Fatalf("CurrencyValue = nil, want %s", test.want)
			}
			if got := field.CurrencyValue.String(); got != test.want {
				t.Errorf("CurrencyValue = %s, want %s", got, test.want)
			}
		})
	}
}

func TestUpsertCustomFieldValueCommand_InvalidCurrencyValueStillErrors(t *testing.T) {
	tests := map[string]string{
		"text":  `{"currencyValue":"abc"}`,
		"minus": `{"currencyValue":"-"}`,
		"bool":  `{"currencyValue":true}`,
	}

	for name, fieldJson := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := decodeCustomField(t, fieldJson); err == nil {
				t.Errorf("want a decode error for %s", fieldJson)
			}
		})
	}
}

// The custom UnmarshalJSON shadows only currencyValue; every other field must
// still decode as the standard decoder would.
func TestUpsertCustomFieldValueCommand_OtherFieldsStillDecode(t *testing.T) {
	field, err := decodeCustomField(t, `{
		"receiptId": 7,
		"customFieldId": 3,
		"stringValue": "text",
		"dateValue": "2026-09-30T00:00:00Z",
		"selectValue": 4,
		"currencyValue": "",
		"booleanValue": true
	}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if field.ReceiptId != 7 {
		t.Errorf("ReceiptId = %d, want 7", field.ReceiptId)
	}
	if field.CustomFieldId != 3 {
		t.Errorf("CustomFieldId = %d, want 3", field.CustomFieldId)
	}
	if field.StringValue == nil || *field.StringValue != "text" {
		t.Errorf("StringValue = %v, want text", field.StringValue)
	}
	if field.DateValue == nil || field.DateValue.Format("2006-01-02") != "2026-09-30" {
		t.Errorf("DateValue = %v, want 2026-09-30", field.DateValue)
	}
	if field.SelectValue == nil || *field.SelectValue != 4 {
		t.Errorf("SelectValue = %v, want 4", field.SelectValue)
	}
	if field.BooleanValue == nil || !*field.BooleanValue {
		t.Errorf("BooleanValue = %v, want true", field.BooleanValue)
	}
	if field.CurrencyValue != nil {
		t.Errorf("CurrencyValue = %s, want nil", field.CurrencyValue)
	}
}
