// Package receiptsource maps receipts onto the rows the reporting engine
// consumes.
//
// It is the only part of the reporting system that knows what a receipt is.
// Everything downstream — grouping, aggregation, formulas, the report model —
// works on typed values and would serve any other source just as well. Reporting
// at item grain, or a dashboard widget over something else entirely, adds a
// sibling of this package rather than changing the engine.
//
// It fetches nothing either. Receipts arrive already loaded, and the caller is
// responsible for having preloaded the associations the report reads.
package receiptsource

import (
	"sort"
	"strconv"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"

	"github.com/shopspring/decimal"
)

// The built-in fields a report may reference. Keys are plain identifiers so a
// formula can name one without quoting, as in SUM(amount).
const (
	KeyReceiptID    reporting.FieldKey = "receipt_id"
	KeyName         reporting.FieldKey = "name"
	KeyAmount       reporting.FieldKey = "amount"
	KeyDate         reporting.FieldKey = "date"
	KeyResolvedDate reporting.FieldKey = "resolved_date"
	KeyCreatedAt    reporting.FieldKey = "created_at"
	KeyStatus       reporting.FieldKey = "status"
	KeyPaidBy       reporting.FieldKey = "paid_by"
	KeyGroup        reporting.FieldKey = "group"
	KeyCategory     reporting.FieldKey = "category"
	KeyTag          reporting.FieldKey = "tag"
)

// customFieldKeyPrefix builds a custom field's key from its id rather than its
// name, so that renaming a custom field cannot break a saved report.
const customFieldKeyPrefix = "custom_"

// CustomFieldKey returns the field key a custom field is referenced by.
func CustomFieldKey(customFieldID uint) reporting.FieldKey {
	return reporting.FieldKey(customFieldKeyPrefix + strconv.FormatUint(uint64(customFieldID), 10))
}

// builtinFields returns the fields every receipt report may reference.
//
// Categories and tags are multi-valued: grouping on one fans a receipt out into
// every bucket it belongs to, attributing the whole amount to each, exactly as
// the dashboard pie chart does.
func builtinFields() []reporting.FieldRef {
	return []reporting.FieldRef{
		{Key: KeyReceiptID, Label: "Receipt Id", DataType: reporting.TypeNumber},
		{Key: KeyName, Label: "Name", DataType: reporting.TypeString},
		{Key: KeyAmount, Label: "Amount", DataType: reporting.TypeCurrency},
		{Key: KeyDate, Label: "Date", DataType: reporting.TypeDate},
		{Key: KeyResolvedDate, Label: "Resolved Date", DataType: reporting.TypeDate},
		{Key: KeyCreatedAt, Label: "Added At", DataType: reporting.TypeDate},
		{Key: KeyStatus, Label: "Status", DataType: reporting.TypeString},
		{Key: KeyPaidBy, Label: "Paid By", DataType: reporting.TypeString},
		{Key: KeyGroup, Label: "Group", DataType: reporting.TypeString},
		{Key: KeyCategory, Label: "Category", DataType: reporting.TypeString, Multi: true},
		{Key: KeyTag, Label: "Tag", DataType: reporting.TypeString, Multi: true},
	}
}

// Source maps receipts to engine rows against a fixed set of custom fields.
//
// It resolves the custom field definitions once, so the catalog it offers and
// the rows it builds cannot disagree about a field's type or about what a
// select option's value is.
type Source struct {
	catalog reporting.FieldCatalog

	// types maps a custom field's id to its type, which decides which of the
	// five value columns on a CustomFieldValue holds the answer.
	types map[uint]models.CustomFieldType

	// options maps a select field's id to its option values by option id, which
	// is how a stored option id becomes the text a reader sees.
	options map[uint]map[uint]string
}

// New builds a Source over the custom fields a report may reference. Pass them
// with their Options loaded; a select field whose options are missing resolves
// to no value rather than to a bare option id.
//
// Two custom fields sharing an id are rejected, since they would claim the same
// field key.
func New(customFields []models.CustomField) (Source, error) {
	// Sort by id so the catalog does not depend on the order the definitions
	// arrived in.
	sorted := make([]models.CustomField, len(customFields))
	copy(sorted, customFields)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })

	source := Source{
		types:   make(map[uint]models.CustomFieldType, len(sorted)),
		options: make(map[uint]map[uint]string),
	}

	fields := builtinFields()
	for _, customField := range sorted {
		source.types[customField.ID] = customField.Type

		if customField.Type == models.SELECT {
			values := make(map[uint]string, len(customField.Options))
			for _, option := range customField.Options {
				values[option.ID] = option.Value
			}
			source.options[customField.ID] = values
		}

		fields = append(fields, reporting.FieldRef{
			Key:      CustomFieldKey(customField.ID),
			Label:    customField.Name,
			DataType: dataTypeFor(customField.Type),
		})
	}

	catalog, err := reporting.NewFieldCatalog(fields...)
	if err != nil {
		return Source{}, err
	}
	source.catalog = catalog

	return source, nil
}

// Catalog returns the fields a report run against this source may reference:
// every built-in, plus one per custom field.
func (s Source) Catalog() reporting.FieldCatalog { return s.catalog }

// dataTypeFor maps a custom field's type onto the engine's. A currency custom
// field is a measure and every other kind is a dimension, which is why a tax
// field needs no special case: it is summable because it is currency.
func dataTypeFor(customFieldType models.CustomFieldType) reporting.DataType {
	switch customFieldType {
	case models.CURRENCY:
		return reporting.TypeCurrency
	case models.DATE:
		return reporting.TypeDate
	case models.BOOLEAN:
		return reporting.TypeBool
	}
	// TEXT and SELECT both read as text.
	return reporting.TypeString
}

// Rows resolves receipts into engine rows, one row per receipt, in the order
// they were given. The engine preserves that order for a records-mode report,
// so the caller's query decides how those rows are sorted.
//
// The caller must have preloaded whatever the report reads: Categories, Tags
// and CustomFields on the receipt, plus PaidByUser and Group for their display
// names. An association that was not loaded resolves to no value, which the
// engine reports as a (None) bucket. A missing preload therefore shows up as an
// empty group rather than as a crash.
func (s Source) Rows(receipts []models.Receipt) []reporting.Row {
	rows := make([]reporting.Row, 0, len(receipts))
	for index := range receipts {
		rows = append(rows, s.row(&receipts[index]))
	}
	return rows
}

func (s Source) row(receipt *models.Receipt) reporting.Row {
	row := reporting.Row{
		KeyReceiptID: {reporting.Num(decimal.NewFromInt(int64(receipt.ID)))},
		KeyName:      {reporting.Str(receipt.Name)},
		KeyAmount:    {reporting.Num(receipt.Amount)},
		KeyDate:      {reporting.DateVal(receipt.Date)},
		KeyCreatedAt: {reporting.DateVal(receipt.CreatedAt)},
		KeyStatus:    {reporting.Str(string(receipt.Status))},
		KeyCategory:  categoryValues(receipt.Categories),
		KeyTag:       tagValues(receipt.Tags),
	}

	if receipt.ResolvedDate != nil {
		row[KeyResolvedDate] = []reporting.Value{reporting.DateVal(*receipt.ResolvedDate)}
	}
	if displayName := userDisplayName(receipt.PaidByUser); len(displayName) > 0 {
		row[KeyPaidBy] = []reporting.Value{reporting.Str(displayName)}
	}
	if len(receipt.Group.Name) > 0 {
		row[KeyGroup] = []reporting.Value{reporting.Str(receipt.Group.Name)}
	}

	s.addCustomFields(row, receipt)

	return row
}

// addCustomFields resolves each of a receipt's custom field values against its
// definition. A field the receipt carries no value for simply has no entry,
// which reads as null when measured and as (None) when grouped. Where a receipt
// holds several values for one field, the first that resolves wins.
func (s Source) addCustomFields(row reporting.Row, receipt *models.Receipt) {
	for _, customFieldValue := range receipt.CustomFields {
		key := CustomFieldKey(customFieldValue.CustomFieldId)
		if _, resolved := row[key]; resolved {
			continue
		}

		value, ok := s.customFieldValue(customFieldValue)
		if !ok {
			continue
		}
		row[key] = []reporting.Value{value}
	}
}

// customFieldValue reads the one column of a value row that its field's type
// says is populated.
//
// The type comes from this source's definitions rather than from the value's
// own embedded CustomField, which the caller need not have preloaded. A value
// whose field is not in the catalog is skipped: the engine could not reference
// it anyway.
func (s Source) customFieldValue(customFieldValue models.CustomFieldValue) (reporting.Value, bool) {
	customFieldType, known := s.types[customFieldValue.CustomFieldId]
	if !known {
		return reporting.Null(), false
	}

	switch customFieldType {
	case models.CURRENCY:
		if customFieldValue.CurrencyValue == nil {
			return reporting.Null(), false
		}
		return reporting.Num(*customFieldValue.CurrencyValue), true

	case models.TEXT:
		if customFieldValue.StringValue == nil {
			return reporting.Null(), false
		}
		return reporting.Str(*customFieldValue.StringValue), true

	case models.SELECT:
		// A select stores an option id, not the text a reader sees.
		if customFieldValue.SelectValue == nil {
			return reporting.Null(), false
		}
		text, found := s.options[customFieldValue.CustomFieldId][*customFieldValue.SelectValue]
		if !found {
			return reporting.Null(), false
		}
		return reporting.Str(text), true

	case models.DATE:
		if customFieldValue.DateValue == nil {
			return reporting.Null(), false
		}
		return reporting.DateVal(*customFieldValue.DateValue), true

	case models.BOOLEAN:
		if customFieldValue.BooleanValue == nil {
			return reporting.Null(), false
		}
		return reporting.Bool(*customFieldValue.BooleanValue), true
	}

	return reporting.Null(), false
}

// userDisplayName prefers what a user chose to be called, falling back to the
// name they log in with, which is how the dashboard names a payer.
func userDisplayName(user models.User) string {
	if len(user.DisplayName) > 0 {
		return user.DisplayName
	}
	return user.Username
}

func categoryValues(categories []models.Category) []reporting.Value {
	values := make([]reporting.Value, 0, len(categories))
	for _, category := range categories {
		values = append(values, reporting.Str(category.Name))
	}
	return values
}

func tagValues(tags []models.Tag) []reporting.Value {
	values := make([]reporting.Value, 0, len(tags))
	for _, tag := range tags {
		values = append(values, reporting.Str(tag.Name))
	}
	return values
}
