package reporting

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// A bucket key identifies the bucket a value groups into, so two values must
// share a key exactly when they compare equal. Anything else either merges
// distinct values into one bucket or splits one value across two.
//
// This is the single law the key must obey, and it is fuzzed in
// FuzzBucketKeyMatchesCompare.
func assertKeyMatchesCompare(t *testing.T, a, b Value) {
	t.Helper()

	sameKey := bucketKey(a) == bucketKey(b)
	equal := compareValues(a, b) == 0

	if sameKey != equal {
		t.Errorf("bucketKey/compareValues disagree for %v (%v) and %v (%v):\n"+
			"  bucketKey(a) = %q\n  bucketKey(b) = %q\n  same key = %v, but compareValues says equal = %v",
			a, a.Type(), b, b.Type(), bucketKey(a), bucketKey(b), sameKey, equal)
	}
}

// time.Time.UnixNano is documented as undefined for dates before 1678 or after
// 2262: the count overflows an int64 and wraps. Keying a date bucket on it
// silently merges instants exactly 2^64 nanoseconds apart.
//
// A receipt whose CreatedAt was never round-tripped through the database holds
// the zero time, which is one of them.
func TestBucketKey_DatesOutsideUnixNanoRange(t *testing.T) {
	zero := time.Time{}                                          // year 1
	wrapped := zero.Add(math.MaxInt64).Add(math.MaxInt64).Add(2) // year 1 + 2^64ns, mid-585

	if zero.Equal(wrapped) {
		t.Fatalf("the fixture is wrong: these must be different instants")
	}
	if zero.UnixNano() != wrapped.UnixNano() {
		t.Fatalf("the fixture is wrong: these must share a UnixNano")
	}

	assertKeyMatchesCompare(t, DateVal(zero), DateVal(wrapped))

	if bucketKey(DateVal(zero)) == bucketKey(DateVal(wrapped)) {
		t.Errorf("year 1 and year 585 share a bucket key")
	}
}

func TestBucketKey_FarFutureDates(t *testing.T) {
	// Both are outside UnixNano's range and differ by a day.
	first := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	second := first.AddDate(0, 0, 1)

	assertKeyMatchesCompare(t, DateVal(first), DateVal(second))
}

// Two rows dated 584 years apart must not merge, and which of them a merged
// bucket would have displayed must not depend on the order the rows arrived in.
// This is the engine's determinism guarantee, and the UnixNano key breaks it.
func TestRun_DistantDatesDoNotMergeOrDependOnInputOrder(t *testing.T) {
	catalog, err := NewFieldCatalog(
		FieldRef{Key: "date", Label: "Date", DataType: TypeDate},
		FieldRef{Key: "amount", Label: "Amount", DataType: TypeCurrency},
	)
	if err != nil {
		t.Fatalf("NewFieldCatalog() error = %v", err)
	}

	spec := ReportSpec{
		GroupBy: []FieldKey{"date"},
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}

	zero := time.Time{}
	wrapped := zero.Add(math.MaxInt64).Add(math.MaxInt64).Add(2)

	forward := []Row{{"date": {DateVal(zero)}}, {"date": {DateVal(wrapped)}}}
	backward := []Row{{"date": {DateVal(wrapped)}}, {"date": {DateVal(zero)}}}

	first, err := Run(spec, catalog, forward, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	second, err := Run(spec, catalog, backward, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(first.Root.Children) != 2 {
		t.Errorf("two distinct dates produced %d bucket(s), want 2", len(first.Root.Children))
	}

	// Even before the merge is fixed, the report must not change with input order.
	if first.Root.Children[0].Value.String() != second.Root.Children[0].Value.String() {
		t.Errorf("reversing the input changed the first bucket: %v vs %v",
			first.Root.Children[0].Value, second.Root.Children[0].Value)
	}
}

// StringFixed(12) truncates, so a value smaller than a picosecond of a cent
// used to share the zero bucket. decimal.String() is canonical and lossless.
func TestBucketKey_NumbersBelowTwelveDecimalPlaces(t *testing.T) {
	tiny := Num(dec("0.0000000000001")) // 13 decimal places
	zero := Num(dec("0"))

	assertKeyMatchesCompare(t, tiny, zero)

	if bucketKey(tiny) == bucketKey(zero) {
		t.Errorf("0.0000000000001 shares a bucket key with 0")
	}
}

// Scale must not affect the key: 200 and 200.00 are the same number.
func TestBucketKey_NumbersIgnoreScale(t *testing.T) {
	pairs := [][2]string{
		{"200", "200.00"},
		{"0.1", "0.10"},
		{"0", "0.000"},
		{"-1.5", "-1.500"},
		{"1e3", "1000"},
	}

	for _, pair := range pairs {
		t.Run(pair[0]+" vs "+pair[1], func(t *testing.T) {
			a, b := Num(dec(pair[0])), Num(dec(pair[1]))
			assertKeyMatchesCompare(t, a, b)
			if bucketKey(a) != bucketKey(b) {
				t.Errorf("%s and %s should share a bucket: %q vs %q", pair[0], pair[1], bucketKey(a), bucketKey(b))
			}
		})
	}
}

// The key law, across every type and every pair, including the pathological
// instants and magnitudes above.
func TestBucketKey_MatchesCompareAcrossTypes(t *testing.T) {
	zero := time.Time{}
	wrapped := zero.Add(math.MaxInt64).Add(math.MaxInt64).Add(2)
	utc := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)

	values := []Value{
		Null(),
		Str(""), Str("0"), Str("z"), Str("Clothing"),
		Num(dec("0")), Num(dec("0.000")), Num(dec("0.0000000000001")),
		Num(dec("200")), Num(dec("200.00")), Num(dec("-1")),
		DateVal(zero), DateVal(wrapped), DateVal(utc),
		DateVal(utc.In(time.FixedZone("plus-one", 3600))),
		DateVal(utc.Add(time.Nanosecond)),
		DateVal(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)),
		Bool(false), Bool(true),
	}

	for i := range values {
		for j := range values {
			assertKeyMatchesCompare(t, values[i], values[j])
		}
	}
}

// A date's key must not depend on the location it was expressed in, only on the
// instant it names.
func TestBucketKey_DatesAreLocationIndependent(t *testing.T) {
	utc := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)

	for _, offset := range []int{-43200, -3600, 0, 3600, 50400} {
		zone := time.FixedZone("zone", offset)
		if bucketKey(DateVal(utc)) != bucketKey(DateVal(utc.In(zone))) {
			t.Errorf("offset %d changed the bucket key", offset)
		}
	}

	// A monotonic clock reading must not leak into the key either.
	instant := time.Now()
	if bucketKey(DateVal(instant)) != bucketKey(DateVal(instant.Round(0))) {
		t.Errorf("a monotonic reading changed the bucket key")
	}
}

// FuzzBucketKeyMatchesCompare searches for any pair of same-typed values whose
// bucket key disagrees with their ordering. It is the property that both the
// date and the number defects violated.
func FuzzBucketKeyMatchesCompare(f *testing.F) {
	// The two known-bad pairs, plus the pairs that must stay equal.
	f.Add(int64(-62135596800), 0, int64(-43688852727), 709551616, "0", "0", 0, 0)                 // year 1 vs year 585
	f.Add(int64(0), 0, int64(0), 0, "0", "0.0000000000001", 0, 0)                                 // 0 vs 1e-13
	f.Add(int64(0), 0, int64(0), 0, "200", "200.00", 0, 0)                                        // scale
	f.Add(int64(1777978800), 0, int64(1777978800), 0, "1", "1", 3600, -3600)                      // one instant, two zones
	f.Add(int64(math.MaxInt32), 999999999, int64(math.MinInt32), 1, "-1.5", "1.5", 43200, -43200) // spread
	f.Add(int64(32503680000), 0, int64(32503766400), 0, "1e30", "1e-30", 0, 0)                    // year 3000 +1d, extremes

	f.Fuzz(func(t *testing.T,
		secondsA int64, nanosA int, secondsB int64, nanosB int,
		numberA, numberB string, offsetA, offsetB int,
	) {
		dateA := fuzzDate(secondsA, nanosA, offsetA)
		dateB := fuzzDate(secondsB, nanosB, offsetB)
		assertKeyMatchesCompare(t, DateVal(dateA), DateVal(dateB))

		decimalA, errA := boundedDecimal(numberA)
		decimalB, errB := boundedDecimal(numberB)
		if errA != nil || errB != nil {
			return
		}
		assertKeyMatchesCompare(t, Num(decimalA), Num(decimalB))

		// Cross-type pairs are never equal and must never share a key.
		assertKeyMatchesCompare(t, DateVal(dateA), Num(decimalA))
		assertKeyMatchesCompare(t, Str(numberA), Num(decimalA))
		assertKeyMatchesCompare(t, Null(), DateVal(dateA))
		assertKeyMatchesCompare(t, Str(numberA), Str(numberB))
	})
}

// fuzzDate builds an instant from fuzz input, keeping it inside time.Time's
// representable range so that formatting it for an error message stays sane.
// The zone it is expressed in must not affect the key, so the fuzzer varies it.
func fuzzDate(seconds int64, nanos int, offset int) time.Time {
	const secondsBound = 200_000_000_000 // roughly the years -4300 to 8300

	seconds %= secondsBound
	if nanos < 0 {
		nanos = -nanos
	}
	nanos %= 1_000_000_000

	offset %= 14 * 3600

	return time.Unix(seconds, int64(nanos)).In(time.FixedZone("fuzz", offset))
}

var errDecimalTooLarge = errors.New("decimal is too large to render")

// boundedDecimal parses a fuzzed number, refusing ones whose rendered form
// would be enormous. bucketKey renders a number through decimal.String(), which
// expands the exponent, so "1e999999999" would otherwise allocate a string with
// a billion characters.
func boundedDecimal(literal string) (decimal.Decimal, error) {
	if len(literal) > 64 {
		return decimal.Decimal{}, errDecimalTooLarge
	}

	value, err := decimal.NewFromString(literal)
	if err != nil {
		return decimal.Decimal{}, err
	}
	if value.Exponent() < -1000 || value.Exponent() > 1000 {
		return decimal.Decimal{}, errDecimalTooLarge
	}

	return value, nil
}
