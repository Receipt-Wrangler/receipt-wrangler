package reporting

// AggFunc is a vertical reduction: it collapses many rows into one value.
type AggFunc uint8

const (
	AggSum AggFunc = iota
	AggCount
	AggAvg
	AggMin
	AggMax
)

func (f AggFunc) String() string {
	switch f {
	case AggSum:
		return "SUM"
	case AggCount:
		return "COUNT"
	case AggAvg:
		return "AVG"
	case AggMin:
		return "MIN"
	case AggMax:
		return "MAX"
	}
	return "UNKNOWN"
}

// valid reports whether the reduction is one the accumulator implements.
//
// ParseAggregate can only produce a known reduction, but Aggregate.Func is the
// canonical form and a persisted template deserializes an integer straight into
// it. An unknown one would otherwise fall through finalize's switch and render
// every cell of the column null, at every level, without an error anywhere.
func (f AggFunc) valid() bool {
	switch f {
	case AggSum, AggCount, AggAvg, AggMin, AggMax:
		return true
	}
	return false
}

// aggFuncFromName resolves a function name to its reduction. Names are matched
// exactly, so "sum" is not "SUM".
func aggFuncFromName(name string) (AggFunc, bool) {
	switch name {
	case "SUM":
		return AggSum, true
	case "COUNT":
		return AggCount, true
	case "AVG":
		return AggAvg, true
	case "MIN":
		return AggMin, true
	case "MAX":
		return AggMax, true
	}
	return 0, false
}

// Aggregate defines an aggregate column: exactly one reduction over one measure
// field. Holding it structurally rather than as a general expression is what
// lets the engine back each aggregate column with a mergeable accumulator, and
// therefore roll it up by combining children rather than recomputing from
// scratch at every level.
type Aggregate struct {
	Func AggFunc

	// Field is the measure being reduced. It is empty for AggCount, which
	// counts records rather than values.
	Field FieldKey
}

// String renders the aggregate in the source form ParseAggregate accepts.
func (a Aggregate) String() string {
	if a.Func == AggCount {
		return "COUNT()"
	}
	return a.Func.String() + "(" + string(a.Field) + ")"
}
