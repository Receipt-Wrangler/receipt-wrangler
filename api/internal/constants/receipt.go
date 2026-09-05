package constants

// DEFAULT_RECEIPT_ORDER_BY is the column a receipt list falls back to when the
// request names no sortable column of its own.
const DEFAULT_RECEIPT_ORDER_BY = "created_at"

// CUSTOM_FIELD_ASSOCIATIONS loads a receipt's custom field values together with
// the definition (and options) each value is read against.
//
// The definition is NOT optional. models.CustomFieldValue.CustomField is a
// non-pointer struct with no omitempty, so loading the values alone still
// serializes a zero-valued definition — including "type":"" — and the mobile
// client's generated CustomFieldType is a closed enum with no empty member.
// A value it cannot deserialize is swallowed by the one_of AnyOfSerializer,
// which would collapse every receipt row in the app rather than error. See
// mobile/CLAUDE.md -> "Serialization contract".
var CUSTOM_FIELD_ASSOCIATIONS = []string{
	"CustomFields",
	"CustomFields.CustomField",
	"CustomFields.CustomField.Options",
}

var FULL_RECEIPT_ASSOCIATIONS = append([]string{
	"PaidByUser",
	"ReceiptItems",
	"ReceiptItems.Categories",
	"ReceiptItems.Tags",
	"ReceiptItems.LinkedItems",
	"ReceiptItems.LinkedItems.Categories",
	"ReceiptItems.LinkedItems.Tags",
	"Comments",
	"ImageFiles",
}, CUSTOM_FIELD_ASSOCIATIONS...)
