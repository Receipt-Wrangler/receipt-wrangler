package commands

import (
	"strings"
	"time"
	"unicode/utf8"

	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

// claimAmountLimit is one cent above the largest value a decimal(10,2) column
// holds. Anything at or over it would be rejected by MySQL/Postgres mid-write and
// silently stored by SQLite, so it is refused up front instead.
var claimAmountLimit = decimal.NewFromInt(100_000_000)

// OptionalAmount is a money value that may be absent. It reads a JSON string or
// number; null, "" and whitespace are no value, because the desktop's masked
// currency input sends "" for a field the user emptied (see
// UpsertCustomFieldValueCommand, which has the same rule for currencyValue).
type OptionalAmount struct {
	Value *decimal.Decimal
}

func (amount *OptionalAmount) UnmarshalJSON(data []byte) error {
	if isBlankCurrencyValue(data) {
		amount.Value = nil
		return nil
	}

	value := decimal.Decimal{}
	if err := value.UnmarshalJSON(data); err != nil {
		return err
	}
	amount.Value = &value
	return nil
}

// validateClaimAmount reports what is wrong with an amount, or "" when it is
// acceptable: below the column's limit, at most two decimal places, and zero or
// more (allowZero) or greater than zero.
func validateClaimAmount(value decimal.Decimal, allowZero bool) string {
	if allowZero && value.IsNegative() {
		return "Amount cannot be negative"
	}
	if !allowZero && !value.IsPositive() {
		return "Amount must be greater than zero"
	}
	if !value.Equal(value.Round(2)) {
		return "Amount can have at most two decimal places"
	}
	if value.GreaterThanOrEqual(claimAmountLimit) {
		return "Amount is too large"
	}
	return ""
}

// parseClaimCalendarDay reads a calendar day from a bare yyyy-MM-dd or an RFC
// 3339 string, keeping the day as written (an offset is never converted), and
// returns it as midnight UTC, the form every calendar day is stored in. It is the
// same rule as repositories.CalendarDayOf, which this package cannot import.
func parseClaimCalendarDay(value string) (time.Time, bool) {
	trimmed := strings.TrimSpace(value)
	parsed, err := time.Parse(time.DateOnly, trimmed)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, trimmed)
		if err != nil {
			return time.Time{}, false
		}
	}

	return utils.CalendarDayUTC(parsed), true
}

// optionalClaimDay parses an optional calendar day. A nil or blank value is no
// day (ok with a nil result); anything else must parse.
func optionalClaimDay(value *string) (*time.Time, bool) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, true
	}

	day, ok := parseClaimCalendarDay(*value)
	if !ok {
		return nil, false
	}
	return &day, true
}

// trimmedOrEmpty dereferences an optional string and trims it.
func trimmedOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func tooLong(value string, max int) bool {
	return utf8.RuneCountInString(value) > max
}
