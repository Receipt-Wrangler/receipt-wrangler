package reporting

import (
	"cmp"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// bucketNumberScale is the number of decimal places a numeric value is
// normalized to when it is used as a grouping key, so that 200 and 200.00 land
// in the same bucket.
const bucketNumberScale = 12

// ValueType discriminates the Value union.
type ValueType uint8

const (
	// ValueNull is the zero ValueType: an absent or undefined value. It is what
	// a missing field resolves to, what an aggregate over nothing can produce,
	// and what an arithmetic expression yields when an operand is null or a
	// divisor is zero.
	ValueNull ValueType = iota
	ValueString
	ValueNumber
	ValueDate
	ValueBool
)

func (t ValueType) String() string {
	switch t {
	case ValueNull:
		return "null"
	case ValueString:
		return "string"
	case ValueNumber:
		return "number"
	case ValueDate:
		return "date"
	case ValueBool:
		return "bool"
	}
	return "unknown"
}

// Value is an immutable typed scalar. The zero Value is null.
//
// The payload is unexported so a Value can never hold a type tag that disagrees
// with its contents; build one with Null, Str, Num, DateVal, or Bool.
type Value struct {
	valueType ValueType
	str       string
	num       decimal.Decimal
	date      time.Time
	b         bool
}

// Null returns the null Value, equivalent to the zero Value.
func Null() Value { return Value{} }

// Str returns a string Value.
func Str(s string) Value { return Value{valueType: ValueString, str: s} }

// Num returns a numeric Value. All money flows through this constructor.
func Num(d decimal.Decimal) Value { return Value{valueType: ValueNumber, num: d} }

// DateVal returns a date Value.
func DateVal(t time.Time) Value { return Value{valueType: ValueDate, date: t} }

// Bool returns a boolean Value.
func Bool(b bool) Value { return Value{valueType: ValueBool, b: b} }

// Type reports the Value's type tag.
func (v Value) Type() ValueType { return v.valueType }

// IsNull reports whether the Value is null.
func (v Value) IsNull() bool { return v.valueType == ValueNull }

// Decimal returns the numeric payload. The second result is false when the
// Value is not a number, in which case the first result is the zero decimal.
func (v Value) Decimal() (decimal.Decimal, bool) {
	if v.valueType != ValueNumber {
		return decimal.Decimal{}, false
	}
	return v.num, true
}

// Text returns the string payload. The second result is false when the Value is
// not a string.
func (v Value) Text() (string, bool) {
	if v.valueType != ValueString {
		return "", false
	}
	return v.str, true
}

// Time returns the date payload. The second result is false when the Value is
// not a date.
func (v Value) Time() (time.Time, bool) {
	if v.valueType != ValueDate {
		return time.Time{}, false
	}
	return v.date, true
}

// Boolean returns the boolean payload. The second result is false when the
// Value is not a boolean.
func (v Value) Boolean() (bool, bool) {
	if v.valueType != ValueBool {
		return false, false
	}
	return v.b, true
}

// String renders the Value for debugging and error messages. Renderers must
// format from the typed payload instead, since presentation depends on the
// viewer's currency and locale settings.
func (v Value) String() string {
	switch v.valueType {
	case ValueString:
		return v.str
	case ValueNumber:
		return v.num.String()
	case ValueDate:
		return v.date.UTC().Format(time.RFC3339Nano)
	case ValueBool:
		return strconv.FormatBool(v.b)
	}
	return "<null>"
}

// Equal reports whether two Values are the same type and payload. Two nulls are
// equal, and numbers compare by value, so 200 equals 200.00.
func (v Value) Equal(other Value) bool { return compareValues(v, other) == 0 }

// compareValues gives a total order over Values, returning a negative number,
// zero, or a positive number as a sorts before, with, or after b.
//
// Null sorts after every non-null value, which is what puts the (None) bucket
// last. Values within one dimension always share a type; the cross-type
// ordering exists only so the comparison stays total, and therefore the sort
// stays deterministic, if that ever fails to hold.
func compareValues(a, b Value) int {
	if a.IsNull() || b.IsNull() {
		switch {
		case a.IsNull() && b.IsNull():
			return 0
		case a.IsNull():
			return 1
		default:
			return -1
		}
	}

	if a.valueType != b.valueType {
		return cmp.Compare(a.valueType, b.valueType)
	}

	switch a.valueType {
	case ValueString:
		return strings.Compare(a.str, b.str)
	case ValueNumber:
		return a.num.Cmp(b.num)
	case ValueDate:
		return a.date.Compare(b.date)
	case ValueBool:
		return cmp.Compare(boolOrder(a.b), boolOrder(b.b))
	}
	return 0
}

func boolOrder(b bool) int {
	if b {
		return 1
	}
	return 0
}

// bucketKey returns a canonical map key identifying the bucket a Value groups
// into. Keys are prefixed by type so values of different types can never
// collide, numbers are normalized so 200 and 200.00 share a bucket, and dates
// key off the absolute instant so the same moment expressed in two time zones
// does not split into two buckets.
func bucketKey(v Value) string {
	switch v.valueType {
	case ValueString:
		return "s:" + v.str
	case ValueNumber:
		return "n:" + v.num.StringFixed(bucketNumberScale)
	case ValueDate:
		return "d:" + strconv.FormatInt(v.date.UnixNano(), 10)
	case ValueBool:
		if v.b {
			return "b:1"
		}
		return "b:0"
	}
	return "z"
}
