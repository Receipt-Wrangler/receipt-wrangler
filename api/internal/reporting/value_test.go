package reporting

import (
	"sort"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func dec(s string) decimal.Decimal {
	return decimal.RequireFromString(s)
}

func TestValue_ZeroValueIsNull(t *testing.T) {
	var zero Value

	if !zero.IsNull() {
		t.Errorf("zero Value.IsNull() = false, want true")
	}
	if zero.Type() != ValueNull {
		t.Errorf("zero Value.Type() = %v, want %v", zero.Type(), ValueNull)
	}
	if !zero.Equal(Null()) {
		t.Errorf("zero Value != Null()")
	}
}

func TestValue_ConstructorsAndAccessors(t *testing.T) {
	instant := time.Date(2026, 5, 1, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    Value
		wantType ValueType
	}{
		{"null", Null(), ValueNull},
		{"string", Str("Clothing"), ValueString},
		{"empty string", Str(""), ValueString},
		{"number", Num(dec("120.00")), ValueNumber},
		{"date", DateVal(instant), ValueDate},
		{"bool", Bool(true), ValueBool},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.value.Type(); got != test.wantType {
				t.Errorf("Type() = %v, want %v", got, test.wantType)
			}
			if got := test.value.IsNull(); got != (test.wantType == ValueNull) {
				t.Errorf("IsNull() = %v, want %v", got, test.wantType == ValueNull)
			}
		})
	}

	t.Run("string payload round trips", func(t *testing.T) {
		got, ok := Str("Clothing").Text()
		if !ok || got != "Clothing" {
			t.Errorf("Text() = %q, %v; want %q, true", got, ok, "Clothing")
		}
	})

	t.Run("number payload round trips", func(t *testing.T) {
		got, ok := Num(dec("120.00")).Decimal()
		if !ok || !got.Equal(dec("120.00")) {
			t.Errorf("Decimal() = %v, %v; want 120.00, true", got, ok)
		}
	})

	t.Run("date payload round trips", func(t *testing.T) {
		got, ok := DateVal(instant).Time()
		if !ok || !got.Equal(instant) {
			t.Errorf("Time() = %v, %v; want %v, true", got, ok, instant)
		}
	})

	t.Run("bool payload round trips", func(t *testing.T) {
		got, ok := Bool(true).Boolean()
		if !ok || !got {
			t.Errorf("Boolean() = %v, %v; want true, true", got, ok)
		}
	})
}

// Accessors must refuse to reinterpret a payload of the wrong type, otherwise a
// null currency field would silently read as a zero amount.
func TestValue_AccessorsRejectWrongType(t *testing.T) {
	tests := []struct {
		name  string
		value Value
	}{
		{"null", Null()},
		{"string", Str("x")},
		{"date", DateVal(time.Now())},
		{"bool", Bool(true)},
	}

	for _, test := range tests {
		t.Run(test.name+" is not a number", func(t *testing.T) {
			got, ok := test.value.Decimal()
			if ok {
				t.Errorf("Decimal() ok = true, want false")
			}
			if !got.IsZero() {
				t.Errorf("Decimal() = %v, want zero decimal", got)
			}
		})
	}

	t.Run("number is not a string", func(t *testing.T) {
		if _, ok := Num(dec("1")).Text(); ok {
			t.Errorf("Text() ok = true, want false")
		}
	})
	t.Run("number is not a date", func(t *testing.T) {
		if _, ok := Num(dec("1")).Time(); ok {
			t.Errorf("Time() ok = true, want false")
		}
	})
	t.Run("number is not a bool", func(t *testing.T) {
		if _, ok := Num(dec("1")).Boolean(); ok {
			t.Errorf("Boolean() ok = true, want false")
		}
	})
}

func TestValue_String(t *testing.T) {
	instant := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{"null", Null(), "<null>"},
		{"string", Str("Clothing"), "Clothing"},
		{"number", Num(dec("120.50")), "120.5"},
		{"date", DateVal(instant), "2026-05-01T12:00:00Z"},
		{"bool", Bool(false), "false"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.value.String(); got != test.want {
				t.Errorf("String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestValue_Equal(t *testing.T) {
	utc := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	elsewhere := utc.In(time.FixedZone("plus-one", 3600))

	tests := []struct {
		name string
		a, b Value
		want bool
	}{
		{"null equals null", Null(), Null(), true},
		{"null does not equal empty string", Null(), Str(""), false},
		{"same string", Str("a"), Str("a"), true},
		{"different string", Str("a"), Str("b"), false},
		{"numbers compare by value not scale", Num(dec("200")), Num(dec("200.00")), true},
		{"different numbers", Num(dec("200")), Num(dec("200.01")), false},
		{"same instant across zones", DateVal(utc), DateVal(elsewhere), true},
		{"different instants", DateVal(utc), DateVal(utc.Add(time.Second)), false},
		{"same bool", Bool(true), Bool(true), true},
		{"different bool", Bool(true), Bool(false), false},
		{"different types", Str("1"), Num(dec("1")), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.a.Equal(test.b); got != test.want {
				t.Errorf("%v.Equal(%v) = %v, want %v", test.a, test.b, got, test.want)
			}
			if got := test.b.Equal(test.a); got != test.want {
				t.Errorf("Equal is not symmetric for %v / %v", test.a, test.b)
			}
		})
	}
}

func TestCompareValues_WithinType(t *testing.T) {
	earlier := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	later := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		a, b Value
		want int
	}{
		{"strings order by bytes", Str("Clothing"), Str("Medical"), -1},
		{"strings equal", Str("Clothing"), Str("Clothing"), 0},
		{"uppercase sorts before lowercase", Str("Z"), Str("a"), -1},
		{"numbers order numerically", Num(dec("9")), Num(dec("10")), -1},
		{"numbers ignore scale", Num(dec("10.0")), Num(dec("10")), 0},
		{"negative numbers", Num(dec("-1")), Num(dec("0")), -1},
		{"dates order chronologically", DateVal(earlier), DateVal(later), -1},
		{"false before true", Bool(false), Bool(true), -1},
		{"bools equal", Bool(true), Bool(true), 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := sign(compareValues(test.a, test.b)); got != test.want {
				t.Errorf("compareValues(%v, %v) = %d, want %d", test.a, test.b, got, test.want)
			}
			if got := sign(compareValues(test.b, test.a)); got != -test.want {
				t.Errorf("compareValues is not antisymmetric for %v / %v", test.a, test.b)
			}
		})
	}
}

// Null is the (None) bucket, and it always sorts last regardless of what it is
// compared against.
func TestCompareValues_NullSortsLast(t *testing.T) {
	others := []Value{
		Str(""),
		Str("zzz"),
		Num(dec("-99999")),
		DateVal(time.Time{}),
		Bool(false),
	}

	for _, other := range others {
		t.Run(other.Type().String(), func(t *testing.T) {
			if got := sign(compareValues(Null(), other)); got != 1 {
				t.Errorf("compareValues(Null, %v) = %d, want 1", other, got)
			}
			if got := sign(compareValues(other, Null())); got != -1 {
				t.Errorf("compareValues(%v, Null) = %d, want -1", other, got)
			}
		})
	}

	if compareValues(Null(), Null()) != 0 {
		t.Errorf("compareValues(Null, Null) != 0")
	}
}

// Mixed types never occur within one dimension, but the comparison must stay
// total so that a sort remains deterministic if they ever do.
func TestCompareValues_MixedTypesAreTotallyOrdered(t *testing.T) {
	values := []Value{
		Bool(true),
		DateVal(time.Unix(0, 0)),
		Num(dec("5")),
		Str("s"),
		Null(),
	}

	sorted := make([]Value, len(values))
	copy(sorted, values)
	sort.SliceStable(sorted, func(i, j int) bool { return compareValues(sorted[i], sorted[j]) < 0 })

	// Sorting a reversed copy must land on the same order.
	reversed := make([]Value, len(values))
	for i, value := range values {
		reversed[len(values)-1-i] = value
	}
	sort.SliceStable(reversed, func(i, j int) bool { return compareValues(reversed[i], reversed[j]) < 0 })

	for i := range sorted {
		if !sorted[i].Equal(reversed[i]) {
			t.Fatalf("sort is not order-independent: %v vs %v at index %d", sorted[i], reversed[i], i)
		}
	}

	if !sorted[len(sorted)-1].IsNull() {
		t.Errorf("null did not sort last, got %v", sorted[len(sorted)-1])
	}
}

func TestBucketKey_SameValueSameBucket(t *testing.T) {
	utc := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	elsewhere := utc.In(time.FixedZone("plus-one", 3600))

	tests := []struct {
		name string
		a, b Value
	}{
		{"identical strings", Str("Clothing"), Str("Clothing")},
		{"numbers differing only in scale", Num(dec("200")), Num(dec("200.00"))},
		{"number with trailing zeros", Num(dec("0.10")), Num(dec("0.1"))},
		{"same instant in different zones", DateVal(utc), DateVal(elsewhere)},
		{"same bool", Bool(true), Bool(true)},
		{"two nulls", Null(), Null()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if bucketKey(test.a) != bucketKey(test.b) {
				t.Errorf("bucketKey(%v) = %q, bucketKey(%v) = %q; want equal",
					test.a, bucketKey(test.a), test.b, bucketKey(test.b))
			}
		})
	}
}

func TestBucketKey_DistinctValuesDistinctBuckets(t *testing.T) {
	utc := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)

	// A wall clock reading of 12:00 in a +1 zone is a different instant to
	// 12:00 UTC, and must not share its bucket.
	sameWallClockDifferentInstant := time.Date(2026, 5, 1, 12, 0, 0, 0, time.FixedZone("plus-one", 3600))

	// Every key must be unique: this catches cross-type collisions such as the
	// string "z" against the null sentinel, or the string "1" against true.
	values := []Value{
		Null(),
		Str("z"),
		Str(""),
		Str("1"),
		Str("0"),
		Str("Clothing"),
		Str("Medical"),
		Num(dec("0")),
		Num(dec("1")),
		Num(dec("-1")),
		Num(dec("0.000000000001")),
		DateVal(utc),
		DateVal(sameWallClockDifferentInstant),
		DateVal(utc.Add(time.Nanosecond)),
		Bool(false),
		Bool(true),
	}

	seen := make(map[string]Value, len(values))
	for _, value := range values {
		key := bucketKey(value)
		if previous, exists := seen[key]; exists {
			t.Errorf("bucketKey collision on %q: %v (%v) and %v (%v)",
				key, previous, previous.Type(), value, value.Type())
			continue
		}
		seen[key] = value
	}
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}
