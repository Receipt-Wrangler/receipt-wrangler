package reporting

import "github.com/shopspring/decimal"

// accumulator backs one aggregate column at one node of the report tree.
//
// It is a monoid: rows fold into it with add, and a parent folds its children
// into itself with merge. A parent therefore combines its children's
// accumulators, never their finalized values. That distinction is the whole
// point of the type.
//
// Summing a child's SUM gives the parent's SUM, so for SUM the two approaches
// agree. Averaging a child's AVG does not give the parent's AVG. By carrying
// the running sum and count and dividing only in finalize, AVG at every level
// is sum(all descendants) / count(all descendants) — which is the only correct
// answer, and needs no special case.
//
// The zero accumulator is a valid empty SUM. Build others with newAccumulator.
type accumulator struct {
	function AggFunc

	// sum totals the non-null values fed to SUM and AVG.
	sum decimal.Decimal

	// values counts the non-null values fed in, and is AVG's divisor. It is not
	// the record count: a row whose measure is null contributes to records but
	// not to values.
	values int64

	// records counts the rows attributed to this node, and is what COUNT
	// reports. A row that fans out across two buckets is counted once in each,
	// so a multi-value dimension double-counts on purpose.
	records int64

	minimum decimal.Decimal
	maximum decimal.Decimal

	// seen records whether any non-null value arrived, which distinguishes a
	// MIN over nothing (null) from a MIN over a single zero.
	seen bool
}

func newAccumulator(function AggFunc) accumulator {
	return accumulator{function: function}
}

// add folds one row's contribution in. Every row increments the record count,
// whether or not it carries a value, so COUNT() counts rows. A null or
// non-numeric value contributes to nothing else: SUM skips it rather than
// treating it as zero, and AVG excludes it from both its sum and its divisor.
func (a *accumulator) add(value Value) {
	a.records++

	number, isNumber := value.Decimal()
	if !isNumber {
		return
	}

	a.sum = a.sum.Add(number)
	a.values++

	if !a.seen {
		a.minimum = number
		a.maximum = number
		a.seen = true
		return
	}
	if number.LessThan(a.minimum) {
		a.minimum = number
	}
	if number.GreaterThan(a.maximum) {
		a.maximum = number
	}
}

// merge folds a child subtree's accumulator into this one. It is associative
// and commutative, so the order children are merged in cannot affect the result.
func (a *accumulator) merge(other accumulator) {
	a.records += other.records
	a.values += other.values
	a.sum = a.sum.Add(other.sum)

	if !other.seen {
		return
	}
	if !a.seen {
		a.minimum = other.minimum
		a.maximum = other.maximum
		a.seen = true
		return
	}
	if other.minimum.LessThan(a.minimum) {
		a.minimum = other.minimum
	}
	if other.maximum.GreaterThan(a.maximum) {
		a.maximum = other.maximum
	}
}

// finalize reduces the accumulator to the value a cell displays.
//
// SUM and COUNT of nothing are zero: a category with no HST shows 0.00, not an
// empty cell. AVG, MIN and MAX of nothing are null, because there is no value
// to report and zero would be a lie.
func (a accumulator) finalize(divisionScale int32) Value {
	switch a.function {
	case AggSum:
		return Num(a.sum)

	case AggCount:
		return Num(decimal.NewFromInt(a.records))

	case AggAvg:
		if a.values == 0 {
			return Null()
		}
		return Num(a.sum.DivRound(decimal.NewFromInt(a.values), divisionScale))

	case AggMin:
		if !a.seen {
			return Null()
		}
		return Num(a.minimum)

	case AggMax:
		if !a.seen {
			return Null()
		}
		return Num(a.maximum)
	}

	// An unrecognized reduction cannot reach here. Validate rejects one with
	// ErrUnknownAggFunc before a single row is touched, and newAccumulator is
	// called from exactly one place, with a Func that has already passed that
	// check.
	//
	// It is null rather than a panic on purpose. Nothing in this package panics
	// — evalArithmetic guards a zero divisor precisely so shopspring cannot —
	// and a panic here would be a line no test could execute, since every path
	// to it runs through Validate first. The risk a panic would guard against is
	// that valid() and this switch drift apart, which is a mistake made at
	// compile time, not at report time; TestAggFunc_ClosedListsAgree catches it
	// there, exhaustively, over all 256 values rather than the one someone
	// forgot.
	return Null()
}
