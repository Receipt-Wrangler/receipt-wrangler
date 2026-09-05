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
	"strings"
	"time"

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

	// Derived date-period fields. A report groups by one of these to bucket
	// receipts by calendar day, month, or year; the raw date fields above carry
	// the exact instant, which would put every receipt in its own group.
	KeyDateDay   reporting.FieldKey = "date_day"
	KeyDateMonth reporting.FieldKey = "date_month"
	KeyDateYear  reporting.FieldKey = "date_year"

	KeyResolvedDateDay   reporting.FieldKey = "resolved_date_day"
	KeyResolvedDateMonth reporting.FieldKey = "resolved_date_month"
	KeyResolvedDateYear  reporting.FieldKey = "resolved_date_year"

	KeyCreatedAtDay   reporting.FieldKey = "created_at_day"
	KeyCreatedAtMonth reporting.FieldKey = "created_at_month"
	KeyCreatedAtYear  reporting.FieldKey = "created_at_year"
)

// customFieldKeyPrefix builds a custom field's key from its id rather than its
// name, so that renaming a custom field cannot break a saved report.
const customFieldKeyPrefix = "custom_"

// CustomFieldKey returns the field key a custom field is referenced by.
func CustomFieldKey(customFieldID uint) reporting.FieldKey {
	return reporting.FieldKey(customFieldKeyPrefix + strconv.FormatUint(uint64(customFieldID), 10))
}

// ParseCustomFieldKey returns the custom field id a key refers to, and whether
// the key is a custom field key at all. It is the inverse of CustomFieldKey, and
// deliberately strict: only "custom_" followed by digits alone qualifies, so a
// derived period key ("custom_7_month") and anything hand-crafted are rejected.
//
// The receipts table sorts by these same keys, which is why this lives beside
// the key builder rather than being re-derived wherever a key is read.
func ParseCustomFieldKey(key string) (uint, bool) {
	digits, found := strings.CutPrefix(key, customFieldKeyPrefix)
	if !found || len(digits) == 0 {
		return 0, false
	}

	// ParseUint would accept a leading "+", which CustomFieldKey never emits.
	for _, character := range digits {
		if character < '0' || character > '9' {
			return 0, false
		}
	}

	customFieldID, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(customFieldID), true
}

// CustomFieldPeriodKeys returns the day/month/year keys derived from a date
// custom field, the custom-field counterpart of date_day / date_month /
// date_year. A report groups by one of these to bucket by calendar period;
// grouping by the raw field buckets on the exact instant, which puts every
// receipt in its own group.
//
// They cannot collide with another custom field's key: CustomFieldKey is always
// "custom_" followed by digits alone, so no id produces "custom_7_month".
func CustomFieldPeriodKeys(customFieldID uint) (day, month, year reporting.FieldKey) {
	base := string(CustomFieldKey(customFieldID))
	return reporting.FieldKey(base + "_day"),
		reporting.FieldKey(base + "_month"),
		reporting.FieldKey(base + "_year")
}

// builtinFields returns the fields every receipt report may reference.
//
// Categories and tags are multi-valued: grouping on one fans a receipt out into
// every bucket it belongs to, attributing the whole amount to each, exactly as
// the dashboard pie chart does.
func builtinFields() []reporting.FieldRef {
	fields := []reporting.FieldRef{
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
	fields = append(fields, dateFieldRefs(KeyDateDay, KeyDateMonth, KeyDateYear, "Date")...)
	fields = append(fields, dateFieldRefs(KeyResolvedDateDay, KeyResolvedDateMonth, KeyResolvedDateYear, "Resolved Date")...)
	fields = append(fields, dateFieldRefs(KeyCreatedAtDay, KeyCreatedAtMonth, KeyCreatedAtYear, "Added At")...)
	return fields
}

// dateFieldRefs builds the day/month/year string fields derived from one date
// column, labelled off a shared prefix (e.g. "Date (Month)"). They are strings,
// not dates, so the engine buckets them by their exact ISO text.
func dateFieldRefs(dayKey, monthKey, yearKey reporting.FieldKey, labelPrefix string) []reporting.FieldRef {
	return []reporting.FieldRef{
		{Key: dayKey, Label: labelPrefix + " (Day)", DataType: reporting.TypeString},
		{Key: monthKey, Label: labelPrefix + " (Month)", DataType: reporting.TypeString},
		{Key: yearKey, Label: labelPrefix + " (Year)", DataType: reporting.TypeString},
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

		if customField.Type == models.DATE {
			day, month, year := CustomFieldPeriodKeys(customField.ID)
			fields = append(fields, dateFieldRefs(day, month, year, customField.Name)...)
		}
	}

	catalog, err := reporting.NewFieldCatalog(fields...)
	if err != nil {
		return Source{}, err
	}
	source.catalog = catalog

	return source, nil
}

// Catalog returns the fields a report run against this source may reference:
// every built-in, plus one per custom field — and, for a date custom field,
// its three derived calendar-period fields as well.
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

	setDateParts(row, KeyDateDay, KeyDateMonth, KeyDateYear, receipt.Date)
	setDateParts(row, KeyCreatedAtDay, KeyCreatedAtMonth, KeyCreatedAtYear, receipt.CreatedAt)

	if receipt.ResolvedDate != nil {
		row[KeyResolvedDate] = []reporting.Value{reporting.DateVal(*receipt.ResolvedDate)}
		setDateParts(row, KeyResolvedDateDay, KeyResolvedDateMonth, KeyResolvedDateYear, *receipt.ResolvedDate)
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

// setDateParts writes the day/month/year strings derived from one date onto the
// row. Each is zero-padded ISO in UTC ("2006-01-02" / "2006-01" / "2006"), so
// the values sort chronologically as plain strings and group into calendar
// buckets. The zone is fixed to UTC so a bucket does not depend on where the
// instant is read.
func setDateParts(row reporting.Row, dayKey, monthKey, yearKey reporting.FieldKey, moment time.Time) {
	utc := moment.UTC()
	row[dayKey] = []reporting.Value{reporting.Str(utc.Format("2006-01-02"))}
	row[monthKey] = []reporting.Value{reporting.Str(utc.Format("2006-01"))}
	row[yearKey] = []reporting.Value{reporting.Str(utc.Format("2006"))}
}

// addCustomFields resolves each of a receipt's custom field values against its
// definition, writing a date field's calendar-period parts alongside its value.
// A field the receipt carries no value for simply has no entry, which reads as
// null when measured and as (None) when grouped.
//
// Where a receipt holds several values for one field, the one with the lowest id
// wins. Nothing stops it holding several: custom_field_values carries no unique
// index on (receipt_id, custom_field_id), and a receipt update replaces the
// association with whatever the request contained. Preferring whichever came
// back first would hand the answer to the database, since the association is
// loaded without an ORDER BY — and a report whose numbers depend on row order is
// the one thing this package must not produce.
//
// A value that does not resolve never wins, so an empty low-id row cannot hide a
// real one. Ties are only possible between rows that were never persisted, since
// a saved row has a distinct id; among those the first still wins.
func (s Source) addCustomFields(row reporting.Row, receipt *models.Receipt) {
	// The id of the value currently holding each field.
	winners := make(map[reporting.FieldKey]uint, len(receipt.CustomFields))

	for _, customFieldValue := range receipt.CustomFields {
		key := CustomFieldKey(customFieldValue.CustomFieldId)

		if incumbent, held := winners[key]; held && incumbent <= customFieldValue.ID {
			continue
		}

		value, ok := s.customFieldValue(customFieldValue)
		if !ok {
			continue
		}

		winners[key] = customFieldValue.ID
		row[key] = []reporting.Value{value}

		// A date field also carries its calendar-period fields. This runs again
		// whenever a lower-id value takes over, and setDateParts overwrites, so
		// the parts always describe the value that actually won.
		if moment, isDate := value.Time(); isDate {
			day, month, year := CustomFieldPeriodKeys(customFieldValue.CustomFieldId)
			setDateParts(row, day, month, year, moment)
		}
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
