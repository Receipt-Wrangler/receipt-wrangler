package reporting

import (
	"testing"
	"time"
)

// addAll folds a list of values into a fresh accumulator.
func addAll(function AggFunc, values ...Value) accumulator {
	acc := newAccumulator(function)
	for _, value := range values {
		acc.add(value)
	}
	return acc
}

func nums(literals ...string) []Value {
	values := make([]Value, 0, len(literals))
	for _, literal := range literals {
		values = append(values, Num(dec(literal)))
	}
	return values
}

// assertFinalize checks a finalized accumulator against a decimal literal, or
// against null when want is empty. Money is compared by value, never by
// reflect.DeepEqual: 200 and 200.00 are equal but hold different internal
// exponents.
func assertFinalize(t *testing.T, acc accumulator, want string) {
	t.Helper()

	got := acc.finalize(testDivisionScale)
	if want == "" {
		if !got.IsNull() {
			t.Errorf("finalize() = %v, want null", got)
		}
		return
	}

	number, isNumber := got.Decimal()
	if !isNumber {
		t.Fatalf("finalize() = %v, want the number %s", got, want)
	}
	if !number.Equal(dec(want)) {
		t.Errorf("finalize() = %s, want %s", number, want)
	}
}

func TestAccumulator_Sum(t *testing.T) {
	tests := []struct {
		name   string
		values []Value
		want   string
	}{
		{"adds values", nums("120.00", "80.00"), "200.00"},
		{"a single value", nums("15.60"), "15.60"},
		// A category with no HST shows 0.00, not an empty cell.
		{"sum of nothing is zero", nil, "0"},
		{"sum of only nulls is zero", []Value{Null(), Null()}, "0"},
		{"nulls are skipped, not treated as zero", []Value{Num(dec("5")), Null(), Num(dec("5"))}, "10"},
		{"negatives", nums("10", "-4"), "6"},
		{"no float drift", nums("0.1", "0.2"), "0.3"},
		{"cents", nums("0.01", "0.01", "0.01"), "0.03"},
		{"non-numeric values are skipped", []Value{Num(dec("5")), Str("x"), Bool(true)}, "5"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFinalize(t, addAll(AggSum, test.values...), test.want)
		})
	}
}

// COUNT counts rows, not values. A row whose measure is null still happened.
func TestAccumulator_Count(t *testing.T) {
	tests := []struct {
		name   string
		values []Value
		want   string
	}{
		{"counts rows", nums("1", "2", "3"), "3"},
		{"count of nothing is zero", nil, "0"},
		{"counts rows with null measures", []Value{Null(), Null()}, "2"},
		{"counts a mix", []Value{Num(dec("1")), Null(), Num(dec("3"))}, "3"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFinalize(t, addAll(AggCount, test.values...), test.want)
		})
	}
}

func TestAccumulator_Avg(t *testing.T) {
	tests := []struct {
		name   string
		values []Value
		want   string
	}{
		{"averages values", nums("10", "20"), "15"},
		// The divisor is the count of values, not of rows: a null does not drag
		// the average down.
		{"nulls are excluded from the divisor", []Value{Num(dec("10")), Null(), Num(dec("20"))}, "15"},
		{"avg of nothing is null", nil, ""},
		{"avg of only nulls is null", []Value{Null()}, ""},
		// 1 / 3 repeats forever; the division scale bounds it deterministically.
		{"a repeating quotient is bounded by the division scale", nums("1", "0", "0"), "0.333333"},
		{"a single zero averages zero, not null", nums("0"), "0"},
		{"the worked example's grand average", nums("120", "80", "90", "30"), "80"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFinalize(t, addAll(AggAvg, test.values...), test.want)
		})
	}
}

func TestAccumulator_MinMax(t *testing.T) {
	tests := []struct {
		name    string
		values  []Value
		wantMin string
		wantMax string
	}{
		{"ordinary range", nums("5", "1", "9"), "1", "9"},
		{"single value", nums("4.20"), "4.20", "4.20"},
		{"negatives", nums("-5", "3"), "-5", "3"},
		// Zero is a value. Null is the absence of one.
		{"a single zero", nums("0"), "0", "0"},
		{"of nothing is null", nil, "", ""},
		{"of only nulls is null", []Value{Null(), Null()}, "", ""},
		{"nulls are skipped", []Value{Null(), Num(dec("7")), Null()}, "7", "7"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertFinalize(t, addAll(AggMin, test.values...), test.wantMin)
			assertFinalize(t, addAll(AggMax, test.values...), test.wantMax)
		})
	}
}

// The single most important property in this package. Averaging the children's
// averages gives the wrong answer; merging their accumulators gives the right
// one.
func TestAccumulator_AvgRollupIsNotAvgOfAvgs(t *testing.T) {
	// Alex: four receipts totalling 120. Sam: one receipt of 80.
	alex := addAll(AggAvg, nums("30", "30", "30", "30")...)
	sam := addAll(AggAvg, nums("80")...)

	// Averaging the two child averages would give (30 + 80) / 2 = 55.
	assertFinalize(t, alex, "30")
	assertFinalize(t, sam, "80")

	parent := newAccumulator(AggAvg)
	parent.merge(alex)
	parent.merge(sam)

	// The truth is 200 / 5 = 40.
	assertFinalize(t, parent, "40")
}

func TestAccumulator_Merge(t *testing.T) {
	tests := []struct {
		name     string
		function AggFunc
		left     []Value
		right    []Value
		want     string
	}{
		{"sum", AggSum, nums("120.00", "80.00"), nums("90.00", "30.00"), "320.00"},
		{"count", AggCount, nums("1", "2"), []Value{Null()}, "3"},
		{"avg", AggAvg, nums("10", "20"), nums("60"), "30"},
		{"min picks the global minimum", AggMin, nums("5", "9"), nums("1"), "1"},
		{"max picks the global maximum", AggMax, nums("5", "9"), nums("12"), "12"},
		{"min ignores an empty child", AggMin, nums("5"), nil, "5"},
		{"max ignores an empty child", AggMax, nil, nums("5"), "5"},
		{"min of two empty children is null", AggMin, nil, nil, ""},
		{"sum of two empty children is zero", AggSum, nil, nil, "0"},
		{"avg of two empty children is null", AggAvg, nil, nil, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parent := newAccumulator(test.function)
			parent.merge(addAll(test.function, test.left...))
			parent.merge(addAll(test.function, test.right...))
			assertFinalize(t, parent, test.want)
		})
	}
}

// Merging must be associative and commutative, otherwise the order the engine
// walks a subtree in would change the numbers it reports.
func TestAccumulator_MergeIsAssociativeAndCommutative(t *testing.T) {
	functions := []AggFunc{AggSum, AggCount, AggAvg, AggMin, AggMax}

	for _, function := range functions {
		t.Run(function.String(), func(t *testing.T) {
			a := addAll(function, nums("3", "1")...)
			b := addAll(function, nums("9")...)
			c := addAll(function, []Value{Null(), Num(dec("4"))}...)

			// merge(a, merge(b, c))
			right := newAccumulator(function)
			right.merge(b)
			right.merge(c)
			leftAssociated := newAccumulator(function)
			leftAssociated.merge(a)
			leftAssociated.merge(right)

			// merge(merge(a, b), c)
			left := newAccumulator(function)
			left.merge(a)
			left.merge(b)
			rightAssociated := newAccumulator(function)
			rightAssociated.merge(left)
			rightAssociated.merge(c)

			// merge(c, merge(b, a)) — reversed at every step.
			reversed := newAccumulator(function)
			reversed.merge(c)
			reversed.merge(b)
			reversed.merge(a)

			associatedOne := leftAssociated.finalize(testDivisionScale)
			associatedTwo := rightAssociated.finalize(testDivisionScale)
			commuted := reversed.finalize(testDivisionScale)

			if !associatedOne.Equal(associatedTwo) {
				t.Errorf("merge is not associative: %v vs %v", associatedOne, associatedTwo)
			}
			if !associatedOne.Equal(commuted) {
				t.Errorf("merge is not commutative: %v vs %v", associatedOne, commuted)
			}
		})
	}
}

// The zero accumulator must behave as an empty one, since a parent starts from
// it before merging any child.
func TestAccumulator_MergeWithEmptyIsIdentity(t *testing.T) {
	functions := []AggFunc{AggSum, AggCount, AggAvg, AggMin, AggMax}

	for _, function := range functions {
		t.Run(function.String(), func(t *testing.T) {
			populated := addAll(function, nums("2", "8")...)
			want := populated.finalize(testDivisionScale)

			mergedAfter := addAll(function, nums("2", "8")...)
			mergedAfter.merge(newAccumulator(function))

			mergedBefore := newAccumulator(function)
			mergedBefore.merge(addAll(function, nums("2", "8")...))

			if got := mergedAfter.finalize(testDivisionScale); !got.Equal(want) {
				t.Errorf("merging an empty child changed the result: %v, want %v", got, want)
			}
			if got := mergedBefore.finalize(testDivisionScale); !got.Equal(want) {
				t.Errorf("merging into an empty parent changed the result: %v, want %v", got, want)
			}
		})
	}
}

// Deep trees roll up one level at a time, and must land on the same numbers as
// aggregating every leaf at once.
func TestAccumulator_DeepRollupMatchesFlatAggregate(t *testing.T) {
	leaves := [][]Value{
		nums("10", "20"),
		nums("30"),
		{Null(), Num(dec("40"))},
		nil,
	}

	functions := []AggFunc{AggSum, AggCount, AggAvg, AggMin, AggMax}

	for _, function := range functions {
		t.Run(function.String(), func(t *testing.T) {
			// Roll up through an intermediate level.
			branchOne := newAccumulator(function)
			branchOne.merge(addAll(function, leaves[0]...))
			branchOne.merge(addAll(function, leaves[1]...))

			branchTwo := newAccumulator(function)
			branchTwo.merge(addAll(function, leaves[2]...))
			branchTwo.merge(addAll(function, leaves[3]...))

			root := newAccumulator(function)
			root.merge(branchOne)
			root.merge(branchTwo)

			// Aggregate every leaf value directly.
			flat := newAccumulator(function)
			for _, leaf := range leaves {
				for _, value := range leaf {
					flat.add(value)
				}
			}

			nested := root.finalize(testDivisionScale)
			direct := flat.finalize(testDivisionScale)
			if !nested.Equal(direct) {
				t.Errorf("rolled up = %v, aggregated flat = %v", nested, direct)
			}
		})
	}
}

// A record is counted even when its measure is missing, and MIN over a set that
// contains a zero is that zero rather than null.
func TestAccumulator_ZeroIsAValueNullIsNot(t *testing.T) {
	withZero := addAll(AggMin, Num(dec("0")))
	assertFinalize(t, withZero, "0")

	withNull := addAll(AggMin, Null())
	assertFinalize(t, withNull, "")

	counted := addAll(AggCount, Null(), Str("ignored"), DateVal(time.Now()))
	assertFinalize(t, counted, "3")
}
