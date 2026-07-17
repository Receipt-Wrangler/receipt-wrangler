package reporting

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Cases at the edges of the engine's semantics: values that look like absences,
// dimensions that hold more than one type, degenerate specs, and formulas built
// to exhaust something.

// An empty string is a value. Absence is not. A receipt whose payer is named ""
// belongs to its own bucket, which sorts before every other string, and is not
// the (None) bucket.
func TestRun_EmptyStringIsNotTheNoneBucket(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by"},
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}

	rows := []Row{
		{"paid_by": {Str("Dana")}},
		{"paid_by": {Str("")}},
		{}, // absent
		{"paid_by": {Null()}},
	}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 3 {
		t.Fatalf("got %d buckets, want 3 (empty string, Dana, (None))", len(model.Root.Children))
	}

	empty, dana, none := model.Root.Children[0], model.Root.Children[1], model.Root.Children[2]

	if text, isText := empty.Value.Text(); !isText || text != "" {
		t.Errorf("first bucket = %v, want the empty string", empty.Value)
	}
	if empty.IsNone {
		t.Errorf("the empty string was treated as (None)")
	}
	if text, _ := dana.Value.Text(); text != "Dana" {
		t.Errorf("second bucket = %v, want Dana", dana.Value)
	}
	// An absent value and an explicit null are the same absence, and share a bucket.
	if !none.IsNone || none.RecordCount != 2 {
		t.Errorf("(None) bucket = %+v, want 2 records", none)
	}
}

// A dimension normally holds one type. If a producer mixes them, the report
// must still be totally ordered and stable, because a report that reorders
// itself is worse than one that looks odd.
func TestRun_CrossTypeDimensionIsOrderedAndStable(t *testing.T) {
	catalog, err := NewFieldCatalog(FieldRef{Key: "mixed", DataType: TypeString})
	if err != nil {
		t.Fatalf("NewFieldCatalog() error = %v", err)
	}

	spec := ReportSpec{
		GroupBy: []FieldKey{"mixed"},
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}

	rows := []Row{
		{"mixed": {Bool(true)}},
		{"mixed": {DateVal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))}},
		{"mixed": {Num(dec("5"))}},
		{"mixed": {Str("s")}},
		{},
	}

	model, err := Run(spec, catalog, rows, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	assertModelInvariants(t, spec, model)

	// Values sort by type, and the (None) bucket is last regardless of type.
	wantTypes := []ValueType{ValueString, ValueNumber, ValueDate, ValueBool, ValueNull}
	if len(model.Root.Children) != len(wantTypes) {
		t.Fatalf("got %d buckets, want %d", len(model.Root.Children), len(wantTypes))
	}
	for index, want := range wantTypes {
		if got := model.Root.Children[index].Value.Type(); got != want {
			t.Errorf("bucket %d has type %v, want %v", index, got, want)
		}
	}

	// And the order does not depend on the order the rows arrived in.
	reversed := make([]Row, len(rows))
	for index, row := range rows {
		reversed[len(rows)-1-index] = row
	}
	permuted, err := Run(spec, catalog, reversed, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if serializeModel(permuted) != serializeModel(model) {
		t.Errorf("reversing the rows reordered a mixed-type dimension")
	}
}

// Aggregating the detail by a dimension the report already groups on is
// degenerate but well defined: every group shows every bucket, because the
// fan-out applies at both levels independently.
//
// Pinned rather than forbidden. It is the cross-product rule applied to itself,
// and an author who writes it gets an answer that is consistent, if strange.
func TestRun_DetailDimensionEqualToAGroupByLevel(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"tag"},
		Detail:  DetailSpec{Mode: DetailAggregate, By: "tag"},
		Columns: []Column{
			{Name: "Tag", Kind: ColumnLabel, Field: "tag"},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		GrandTotals: true,
	}

	rows := []Row{{"tag": {Str("Alex"), Str("Sam")}, "amount": {Num(dec("10.00"))}}}
	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d groups, want 2", len(model.Root.Children))
	}
	for _, group := range model.Root.Children {
		if len(group.DetailRows) != 2 {
			t.Errorf("group %v has %d detail rows, want 2", group.Value, len(group.DetailRows))
		}
	}

	// Two groups times two detail buckets: the row is attributed four times.
	assertRow(t, model.GrandTotals, map[string]string{"Subtotal": "40.00"})
}

// MIN and MAX must combine correctly through an internal node, not only within
// a leaf. Everything else about them is unit-tested on the accumulator.
func TestRun_MinMaxRollUpThroughTwoLevels(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by", "tag"},
		Columns: []Column{
			{Name: "Least", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMin, Field: "amount"}},
			{Name: "Most", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMax, Field: "amount"}},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{
		// Dana's extremes are split across her two children.
		receiptRow("Dana", "Alex", "", "5.00", ""),
		receiptRow("Dana", "Alex", "", "9.00", ""),
		receiptRow("Dana", "Sam", "", "1.00", ""),
		receiptRow("Dana", "Sam", "", "3.00", ""),
		// Eve holds the report's maximum, and none of its minimum.
		receiptRow("Eve", "Kim", "", "50.00", ""),
	}

	model := mustRun(t, spec, rows)
	dana, eve := model.Root.Children[0], model.Root.Children[1]

	assertRow(t, dana.Children[0].Subtotals, map[string]string{"Least": "5.00", "Most": "9.00"})
	assertRow(t, dana.Children[1].Subtotals, map[string]string{"Least": "1.00", "Most": "3.00"})
	assertRow(t, dana.Subtotals, map[string]string{"Least": "1.00", "Most": "9.00"})
	assertRow(t, eve.Subtotals, map[string]string{"Least": "50.00", "Most": "50.00"})
	assertRow(t, model.GrandTotals, map[string]string{"Least": "1.00", "Most": "50.00"})
}

// A group whose rows all carry a null measure reports null extremes and a zero
// sum, and that null must survive the merge into a parent that has real values.
func TestRun_NullExtremesMergeWithRealOnes(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by"},
		Columns: []Column{
			{Name: "Least", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMin, Field: "custom_1"}},
			{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{
		{"paid_by": {Str("Dana")}, "custom_1": {Num(dec("4.00"))}},
		{"paid_by": {Str("Eve")}}, // no tax recorded at all
	}

	model := mustRun(t, spec, rows)

	assertRow(t, model.Root.Children[0].Subtotals, map[string]string{"Least": "4.00", "Total": "4.00"})
	assertRow(t, model.Root.Children[1].Subtotals, map[string]string{"Least": "", "Total": "0"})
	assertRow(t, model.GrandTotals, map[string]string{"Least": "4.00", "Total": "4.00"})
}

// A formula's source is bounded before it is parsed, which is the only thing
// that bounds what parsing costs.
//
// It is tempting to rely on the expression language's ten-thousand-node cap
// instead. That cap counts nodes as they are built, and a parenthesis builds
// none — it only groups — while the parser descends several frames to read it.
// Nesting is therefore bounded by the goroutine stack alone: about 640 bytes of
// stack per parenthesis, so Go's default one-gigabyte ceiling falls near 1.6
// million of them. Reaching it is a fatal stack overflow, which recover cannot
// catch and which takes the process, not the request.
//
// An earlier version of this test asserted that 50,000 nested parentheses parse
// fine, and read that as proof the cap protected the stack. It proved only that
// 50,000 frames fit under the ceiling.
func TestParseArithmetic_SizeLimits(t *testing.T) {
	// One long identifier, so the two cases differ by exactly the one byte the
	// bound is stated in.
	t.Run("a formula at exactly the limit is accepted", func(t *testing.T) {
		source := strings.Repeat("a", maxFormulaLength)

		if _, err := ParseArithmetic(source); err != nil {
			t.Errorf("a formula of exactly %d bytes was refused: %v", maxFormulaLength, err)
		}
	})

	t.Run("one byte past the limit is refused", func(t *testing.T) {
		source := strings.Repeat("a", maxFormulaLength+1)

		node, err := ParseArithmetic(source)
		if !errors.Is(err, ErrFormulaTooLong) {
			t.Errorf("error = %v, want ErrFormulaTooLong", err)
		}
		if node != nil {
			t.Errorf("a rejected formula returned a node")
		}
	})

	t.Run("a realistic chain well inside the limit is accepted", func(t *testing.T) {
		source := "a" + strings.Repeat("+a", 100)

		if _, err := ParseArithmetic(source); err != nil {
			t.Errorf("a %d-byte chain was refused: %v", len(source), err)
		}
	})

	t.Run("nesting deep enough to exhaust the stack never reaches the parser", func(t *testing.T) {
		// Far below the ~1.6M that overflows, and far above the limit: the point
		// is that the length check answers before the parser is ever called.
		depth := 100_000
		source := strings.Repeat("(", depth) + "a" + strings.Repeat(")", depth)

		if _, err := ParseArithmetic(source); !errors.Is(err, ErrFormulaTooLong) {
			t.Errorf("error = %v, want ErrFormulaTooLong", err)
		}
	})

	t.Run("a long chain is refused by length, not by the node cap", func(t *testing.T) {
		source := "a" + strings.Repeat("+a", 20000)

		if _, err := ParseArithmetic(source); !errors.Is(err, ErrFormulaTooLong) {
			t.Errorf("error = %v, want ErrFormulaTooLong", err)
		}
	})

	t.Run("nesting within the limit still parses to a shallow tree", func(t *testing.T) {
		// Parentheses group; they do not build nodes. That remains true, and is
		// why the node cap could never have bounded them.
		depth := (maxFormulaLength - 1) / 2
		source := strings.Repeat("(", depth) + "a" + strings.Repeat(")", depth)

		node, err := ParseArithmetic(source)
		if err != nil {
			t.Fatalf("ParseArithmetic() error = %v", err)
		}
		if refs := columnRefs(node); len(refs) != 1 || refs[0] != "a" {
			t.Errorf("columnRefs = %v, want [a]", refs)
		}
		if got := evalArithmetic(node, map[string]Value{"a": Num(dec("7"))}, 6); got.String() != "7" {
			t.Errorf("eval = %v, want 7", got)
		}
	})

	t.Run("deeply nested calls are refused", func(t *testing.T) {
		depth := 5000
		source := strings.Repeat("ROUND(", depth) + "a" + strings.Repeat(", 2)", depth)

		if _, err := ParseArithmetic(source); err == nil {
			t.Fatalf("a %d-deep call chain parsed", depth)
		}
	})

	t.Run("an aggregate source is bounded too", func(t *testing.T) {
		source := "SUM(" + strings.Repeat("a", maxFormulaLength) + ")"

		if _, err := ParseAggregate(source); !errors.Is(err, ErrFormulaTooLong) {
			t.Errorf("error = %v, want ErrFormulaTooLong", err)
		}
	})
}

// Bounding a formula's source (above) bounds what one expression costs to parse;
// it does nothing about what an expression costs to evaluate. Arithmetic columns
// may reference one another, so a chain that each squares the previous — c1 =
// c0 * c0, c2 = c1 * c1, … — makes cN equal 99^(2^N): one decimal whose big.Int
// coefficient doubles in length every column, reaching hundreds of megabytes
// within ~30 columns and exhausting memory, all from a spec a few hundred bytes
// long. evalBinary caps a computed value's magnitude, so an over-large
// intermediate collapses to a null cell (the engine's "correct or null, never
// wrong" contract).
//
// The chain here is kept deliberately short: even with the guard removed the
// deepest column stays a few kilobytes, so a regression fails the assertions
// below rather than OOMing the test process — while still being long enough that
// the guard must null the deep columns for the test to pass.
func TestRun_BoundsRunawayArithmeticGrowth(t *testing.T) {
	const depth = 14

	columns := []Column{
		{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
		{Name: "c0", Kind: ColumnArithmetic, Expr: "99 * 99"},
	}
	for i := 1; i < depth; i++ {
		prev := fmt.Sprintf("c%d", i-1)
		columns = append(columns, Column{
			Name: fmt.Sprintf("c%d", i),
			Kind: ColumnArithmetic,
			Expr: prev + " * " + prev,
		})
	}

	spec := ReportSpec{
		GroupBy:     []FieldKey{"tag"},
		Columns:     columns,
		GrandTotals: true,
	}

	model := mustRun(t, spec, []Row{{"tag": {Str("Alex")}}})

	// The deepest column has to have collapsed to a null cell rather than
	// materialized an astronomically large number.
	deepest := fmt.Sprintf("c%d", depth-1)
	if got := cell(t, model.GrandTotals, deepest); !got.IsNull() {
		t.Errorf("deep arithmetic column %s = %v, want a null (bounded) cell", deepest, got)
	}

	// And nothing that did survive exceeds an independent sanity ceiling (not the
	// production bound, so a weakened bound is caught rather than tracked).
	for _, c := range model.GrandTotals {
		if number, ok := c.Value().Decimal(); ok && number.NumDigits() > runawayDigitCeiling {
			t.Errorf("column %s survived with %d digits, past the independent %d ceiling",
				c.Column, number.NumDigits(), runawayDigitCeiling)
		}
	}
}

// Validate is total: whatever it is handed, it answers, and it never panics.
func TestValidate_IsTotal(t *testing.T) {
	var zeroCatalog FieldCatalog

	specs := []struct {
		name string
		spec ReportSpec
	}{
		{"zero spec", ReportSpec{}},
		{"unknown label field", ReportSpec{Columns: []Column{{Name: "A", Kind: ColumnLabel, Field: "nope"}}}},
		{"out of range kind", ReportSpec{Columns: []Column{{Name: "A", Kind: ColumnKind(200)}}}},
		{"out of range aggregate", ReportSpec{Columns: []Column{{Name: "A", Kind: ColumnAggregate, Agg: Aggregate{Func: AggFunc(99), Field: "x"}}}}},
		{"out of range detail mode", ReportSpec{Detail: DetailSpec{Mode: DetailMode(9)}, Columns: []Column{{Name: "A", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}}},
		{"empty groupBy key", ReportSpec{GroupBy: []FieldKey{""}, Columns: []Column{{Name: "A", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}}},
		{"nil columns", ReportSpec{GroupBy: []FieldKey{"tag"}}},
		{"formula reading itself through a call", ReportSpec{Columns: []Column{{Name: "A", Kind: ColumnArithmetic, Expr: "ROUND(A, 2)"}}}},
		{"negative scale", ReportSpec{Config: EngineConfig{DivisionScale: -5}, Columns: []Column{{Name: "A", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}}},
	}

	for _, test := range specs {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("Validate panicked: %v", recovered)
				}
			}()

			// Against a catalog with no fields, and against a full one.
			if err := Validate(test.spec, zeroCatalog); err == nil {
				t.Errorf("Validate() against an empty catalog returned nil")
			}
			_ = Validate(test.spec, testCatalog(t))
		})
	}
}

// Run rejects what Validate rejects, before touching a single row.
func TestRun_RejectsWhatValidateRejects(t *testing.T) {
	spec := ReportSpec{Columns: []Column{{Name: "Loop", Kind: ColumnArithmetic, Expr: "Loop + 1"}}}

	rows := []Row{{"amount": {Num(dec("1"))}}}
	if _, err := Run(spec, testCatalog(t), rows, MetaInput{}); err == nil {
		t.Fatalf("Run() accepted a spec Validate rejects")
	}
}

// A row may carry fields the catalog does not know, and may omit fields it
// does. Neither is an error: the first is ignored, the second is (None).
func TestRun_RowsNeedNotMatchTheCatalog(t *testing.T) {
	spec := ReportSpec{
		GroupBy:     []FieldKey{"tag"},
		Columns:     []Column{{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}}},
		GrandTotals: true,
	}

	rows := []Row{
		{"tag": {Str("Alex")}, "amount": {Num(dec("1.00"))}, "unknown_field": {Str("ignored")}},
		{"amount": {Num(dec("2.00"))}}, // no tag
	}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d buckets, want 2", len(model.Root.Children))
	}
	if !model.Root.Children[1].IsNone {
		t.Errorf("the row with no tag did not land in (None)")
	}
	assertRow(t, model.GrandTotals, map[string]string{"Total": "3.00"})
}
