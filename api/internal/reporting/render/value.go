package render

import (
	"time"

	"receipt-wrangler/api/internal/reporting"
)

// dimensionDateLayout is how a date reads when it names a bucket or fills a
// label column. It matches the derived date-period fields receiptsource emits
// (receiptsource.setDateParts), so grouping by "Date" and grouping by
// "Date (Day)" print the same text for the same day.
const dimensionDateLayout = "2006-01-02"

// formatLabelValue renders a value that names something — a group bucket, or a
// label column's cell — as the text a reader sees.
//
// The engine emits raw typed values and leaves presentation to a renderer, so
// without this a boolean prints "true", a date prints an RFC 3339 instant, and a
// currency amount prints a bare decimal with no symbol. Which of those a value
// is cannot be read off the value alone in every case, so the caller supplies
// the field's declared data type. A date names its calendar day in the zone its
// field declares (location; nil means UTC): a calendar-day field such as a
// receipt's Date is stored as midnight UTC and must be read in UTC, while an
// instant such as Added At is read in the app's time zone. The engine emits
// date buckets in UTC (its canonical form), so the zone is applied here, at
// presentation, and never changes which values share a bucket.
//
// The type is a declaration, not a guarantee: a producer may hand over a value
// whose payload disagrees with it. Each branch therefore checks the payload and
// falls through to the value's own rendering rather than substituting a zero.
func formatLabelValue(
	value reporting.Value,
	dataType reporting.DataType,
	location *time.Location,
	currency *reporting.CurrencyFormat,
	noneLabel string,
) string {
	if value.IsNull() {
		return noneLabel
	}

	switch dataType {
	case reporting.TypeBool:
		if boolean, isBool := value.Boolean(); isBool {
			if boolean {
				return "Yes"
			}
			return "No"
		}

	case reporting.TypeDate:
		if moment, isDate := value.Time(); isDate {
			return moment.In(locationOrUTC(location)).Format(dimensionDateLayout)
		}

	case reporting.TypeCurrency:
		if number, isNumber := value.Decimal(); isNumber {
			if currency != nil {
				return formatCurrency(number, *currency)
			}
			return number.StringFixed(2)
		}
	}

	return value.String()
}

// locationOrUTC treats an undeclared zone as UTC, so a producer that never sets
// one keeps the engine's canonical reading.
func locationOrUTC(location *time.Location) *time.Location {
	if location == nil {
		return time.UTC
	}
	return location
}
