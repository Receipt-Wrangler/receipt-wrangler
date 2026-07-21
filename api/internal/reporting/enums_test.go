package reporting

import "testing"

// An enum the engine switches on in more than one place is only as sound as the
// agreement between those switches, and Go enforces none of it: a switch that
// forgets a value compiles, runs, and returns whatever its fallthrough says.
//
// That is not hypothetical here. Validate once accepted an AggFunc the
// accumulator did not implement, and every cell of the column rendered blank at
// every level with no error anywhere — the switch in finalize simply fell
// through to null. The fix was to add valid(), which closed the hole but created
// a second closed list that can drift from the first.
//
// So the lists are pinned to each other, exhaustively over the whole domain
// rather than over the five values someone remembered. These tests are what make
// the fallthrough in accumulator.finalize safe to leave as null instead of a
// panic: drift is caught here, at `go test` time, rather than in a report
// somebody is waiting on.

// TestAggFunc_ClosedListsAgree pins the four independent switches over AggFunc:
// String, valid, aggFuncFromName, and accumulator.finalize.
func TestAggFunc_ClosedListsAgree(t *testing.T) {
	for candidate := 0; candidate <= 255; candidate++ {
		function := AggFunc(candidate)
		valid := function.valid()

		// String and valid must admit the same set.
		if named := function.String() != "UNKNOWN"; named != valid {
			t.Errorf("AggFunc(%d): valid() = %v but String() = %q", candidate, valid, function.String())
		}

		// finalize must implement exactly the set valid() admits. A populated
		// accumulator finalizes to a number under every known reduction — AVG's
		// empty-divisor guard and MIN/MAX's unseen guard are both satisfied by
		// the single value below — so only the fallthrough can yield null.
		accumulator := newAccumulator(function)
		accumulator.add(Num(dec("1")))

		if implemented := !accumulator.finalize(testDivisionScale).IsNull(); implemented != valid {
			t.Errorf("AggFunc(%d) (%s): valid() = %v but finalize implements it = %v — "+
				"a reduction Validate accepts and the accumulator does not implement renders every cell blank",
				candidate, function, valid, implemented)
		}

		// Only round-trip the valid ones: aggFuncFromName reports failure as
		// (0, false), and 0 is AggSum, so an unknown value would appear to
		// round-trip to SUM.
		if !valid {
			continue
		}
		parsed, known := aggFuncFromName(function.String())
		if !known || parsed != function {
			t.Errorf("AggFunc(%d) (%s): aggFuncFromName(%q) = (%v, %v), want (%v, true)",
				candidate, function, function.String(), parsed, known, function)
		}
	}
}

// TestDetailMode_ClosedListsAgree pins valid against String, the same way.
func TestDetailMode_ClosedListsAgree(t *testing.T) {
	for candidate := 0; candidate <= 255; candidate++ {
		mode := DetailMode(candidate)
		valid := mode.valid()

		if named := mode.String() != "unknown"; named != valid {
			t.Errorf("DetailMode(%d): valid() = %v but String() = %q", candidate, valid, mode.String())
		}
	}
}

// isAggregate is the single spelling of "are the bottom rows summed buckets?",
// and every valid mode must answer it one way or the other. The engine reads it
// as a biconditional — insertLeaf, rollUpLeaf and emitDetailRows all branch on
// its negation — so a mode that is neither records nor aggregate would take the
// records path everywhere while compileDetail called it something else.
func TestDetailMode_IsAggregatePartitionsTheValidModes(t *testing.T) {
	if !DetailAggregate.valid() || !DetailRecords.valid() {
		t.Fatal("the two known modes must be valid")
	}
	if !(DetailSpec{Mode: DetailAggregate}).isAggregate() {
		t.Error("DetailAggregate is not aggregate")
	}
	if (DetailSpec{Mode: DetailRecords}).isAggregate() {
		t.Error("DetailRecords is aggregate")
	}
}
