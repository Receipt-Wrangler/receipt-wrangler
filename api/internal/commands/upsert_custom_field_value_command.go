package commands

import (
	"bytes"
	"encoding/json"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type UpsertCustomFieldValueCommand struct {
	ReceiptId     uint             `json:"receiptId"`
	CustomFieldId uint             `json:"customFieldId"`
	StringValue   *string          `json:"stringValue"`
	DateValue     *time.Time       `json:"dateValue"`
	SelectValue   *uint            `json:"selectValue"`
	CurrencyValue *decimal.Decimal `json:"currencyValue"`
	BooleanValue  *bool            `json:"booleanValue"`
}

// UnmarshalJSON reads a blank currencyValue ("" or whitespace) as no value, the
// same as null. The desktop's masked currency input sends "" for a field the user
// emptied, and decimal.Decimal rejects "", which failed the whole receipt save.
// A blank value is stored like an untouched one, including a group default custom
// field the server adds itself (see api/CLAUDE.md -> "Group Default Custom Fields").
func (command *UpsertCustomFieldValueCommand) UnmarshalJSON(data []byte) error {
	type plain UpsertCustomFieldValueCommand
	wire := struct {
		*plain
		CurrencyValue json.RawMessage `json:"currencyValue"`
	}{plain: (*plain)(command)}

	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	// Key absent: leave the field alone, as the standard decoder would.
	if wire.CurrencyValue == nil {
		return nil
	}

	if isBlankCurrencyValue(wire.CurrencyValue) {
		command.CurrencyValue = nil
		return nil
	}

	value := decimal.Decimal{}
	if err := value.UnmarshalJSON(wire.CurrencyValue); err != nil {
		return err
	}
	command.CurrencyValue = &value

	return nil
}

func isBlankCurrencyValue(raw json.RawMessage) bool {
	if bytes.Equal(raw, []byte("null")) {
		return true
	}

	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return false
	}

	return strings.TrimSpace(text) == ""
}
