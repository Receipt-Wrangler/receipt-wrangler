package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type CustomFieldValue struct {
	BaseModel
	Receipt Receipt `json:"-"`
	// Indexed together, in this order, for the correlated subquery the receipts
	// table sorts a custom field column with: it filters receipt_id and
	// custom_field_id, then takes the lowest id. Without the index that is a full
	// scan of this table per candidate receipt (MySQL gets one from the FK;
	// SQLite and Postgres do not). Leading with receipt_id also serves the plain
	// "load this receipt's values" reads.
	ReceiptId     uint             `gorm:"index:idx_custom_field_value_lookup,priority:1" json:"receiptId"`
	CustomField   CustomField      `json:"customField"`
	CustomFieldId uint             `gorm:"index:idx_custom_field_value_lookup,priority:2" json:"customFieldId"`
	StringValue   *string          `json:"stringValue"`
	DateValue     *time.Time       `json:"dateValue"`
	SelectValue   *uint            `json:"selectValue"`
	CurrencyValue *decimal.Decimal `json:"currencyValue"`
	BooleanValue  *bool            `json:"booleanValue"`
}
