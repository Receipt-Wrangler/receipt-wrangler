package reporting

import (
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

// The expression language caps the size of a tree it will build, so a formula
// cannot be made to exhaust the stack or the heap. Parentheses add no nodes and
// so parse to a shallow tree however deeply they nest.
func TestParseArithmetic_SizeLimits(t *testing.T) {
	t.Run("a very long chain is refused, not fatal", func(t *testing.T) {
		source := "a" + strings.Repeat("+a", 20000)

		node, err := ParseArithmetic(source)
		if err == nil {
			t.Fatalf("a %d-term chain parsed", 20000)
		}
		if node != nil {
			t.Errorf("a rejected formula returned a node")
		}
		assertFormulaError(t, "long chain", err)
	})

	t.Run("deeply nested parentheses parse to a shallow tree", func(t *testing.T) {
		depth := 50000
		source := strings.Repeat("(", depth) + "a" + strings.Repeat(")", depth)

		node, err := ParseArithmetic(source)
		if err != nil {
			t.Fatalf("ParseArithmetic() error = %v", err)
		}
		// Parentheses group; they do not build nodes.
		refs := columnRefs(node)
		if len(refs) != 1 || refs[0] != "a" {
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
