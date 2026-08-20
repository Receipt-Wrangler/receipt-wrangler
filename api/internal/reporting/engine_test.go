package reporting

import (
	"reflect"
	"testing"
	"time"
)

var testGeneratedAt = time.Date(2026, 6, 1, 9, 30, 0, 0, time.UTC)

func testMeta() MetaInput {
	return MetaInput{
		GeneratedAt: testGeneratedAt,
		Params:      map[string]string{"period": "2026-05-01 TO 2026-05-31"},
		Currency:    &CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: "."},
	}
}

// receiptRow builds a row shaped like a receipt: one payer, one tag, one
// category, an amount and a tax amount.
func receiptRow(paidBy, tag, category, amount, hst string) Row {
	row := Row{
		"paid_by":  {Str(paidBy)},
		"tag":      {Str(tag)},
		"category": {Str(category)},
		"amount":   {Num(dec(amount))},
	}
	if hst != "" {
		row["custom_1"] = []Value{Num(dec(hst))}
	}
	return row
}

// mustRun runs a spec and asserts the model's structural invariants, so every
// engine test below checks them in passing whatever else it is asserting.
func mustRun(t *testing.T, spec ReportSpec, rows []Row) ReportModel {
	t.Helper()

	model, err := Run(spec, testCatalog(t), rows, testMeta())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	assertModelInvariants(t, spec, model)

	return model
}

// cell reads a named column out of a row of cells.
func cell(t *testing.T, cells []Cell, column string) Value {
	t.Helper()

	for _, candidate := range cells {
		if candidate.Column == column {
			return candidate.Value()
		}
	}
	t.Fatalf("no cell for column %q", column)
	return Null()
}

// assertRow checks a whole row against decimal literals, where "" means null.
func assertRow(t *testing.T, cells []Cell, want map[string]string) {
	t.Helper()

	for column, expected := range want {
		got := cell(t, cells, column)
		if expected == "" {
			if !got.IsNull() {
				t.Errorf("column %s = %v, want null", column, got)
			}
			continue
		}

		number, isNumber := got.Decimal()
		if !isNumber {
			// Fall back to a string comparison for label columns.
			if text, isText := got.Text(); isText {
				if text != expected {
					t.Errorf("column %s = %q, want %q", column, text, expected)
				}
				continue
			}
			t.Errorf("column %s = %v, want %s", column, got, expected)
			continue
		}
		if !number.Equal(dec(expected)) {
			t.Errorf("column %s = %s, want %s", column, number, expected)
		}
	}
}

// workedExampleRows is the data behind the design document's worked example:
//
//	Foster Parent: Dana
//	  Child: Alex
//	                     Count   Subtotal    Hst      Total     Avg/Receipt
//	    Clothing           4      120.00    15.60    135.60        33.90
//	    Medical            2       80.00     0.00     80.00        40.00
//	    TOTALS (Alex)      6      200.00    15.60    215.60        35.93
//	  Child: Sam
//	    Clothing           3       90.00    11.70    101.70        33.90
//	    Mileage            1       30.00     0.00     30.00        30.00
//	    TOTALS (Sam)       4      120.00    11.70    131.70        32.93
//	  GRAND TOTALS        10      320.00    27.30    347.30        34.73
func workedExampleRows() []Row {
	return []Row{
		// Alex, Clothing: four receipts totalling 120.00 with 15.60 of tax.
		receiptRow("Dana", "Alex", "Clothing", "30.00", "3.90"),
		receiptRow("Dana", "Alex", "Clothing", "30.00", "3.90"),
		receiptRow("Dana", "Alex", "Clothing", "30.00", "3.90"),
		receiptRow("Dana", "Alex", "Clothing", "30.00", "3.90"),
		// Alex, Medical: two receipts totalling 80.00, no tax recorded at all.
		receiptRow("Dana", "Alex", "Medical", "50.00", ""),
		receiptRow("Dana", "Alex", "Medical", "30.00", ""),
		// Sam, Clothing: three receipts totalling 90.00 with 11.70 of tax.
		receiptRow("Dana", "Sam", "Clothing", "30.00", "3.90"),
		receiptRow("Dana", "Sam", "Clothing", "30.00", "3.90"),
		receiptRow("Dana", "Sam", "Clothing", "30.00", "3.90"),
		// Sam, Mileage: one receipt of 30.00, no tax.
		receiptRow("Dana", "Sam", "Mileage", "30.00", ""),
	}
}

// The golden test. Every number below is read off the design document.
func TestRun_WorkedExample(t *testing.T) {
	model := mustRun(t, workedExampleSpec(), workedExampleRows())

	if model.Meta.Title != "Verified Expenses" {
		t.Errorf("Meta.Title = %q", model.Meta.Title)
	}
	if !model.Meta.GeneratedAt.Equal(testGeneratedAt) {
		t.Errorf("Meta.GeneratedAt = %v, want %v", model.Meta.GeneratedAt, testGeneratedAt)
	}
	if model.Meta.Params["period"] != "2026-05-01 TO 2026-05-31" {
		t.Errorf("Meta.Params = %v", model.Meta.Params)
	}
	if model.Meta.NoneLabel != defaultNoneLabel {
		t.Errorf("Meta.NoneLabel = %q, want %q", model.Meta.NoneLabel, defaultNoneLabel)
	}

	// One foster parent.
	if len(model.Root.Children) != 1 {
		t.Fatalf("root has %d children, want 1", len(model.Root.Children))
	}
	dana := model.Root.Children[0]
	if dana.Dimension != "paid_by" {
		t.Errorf("level 1 dimension = %q, want paid_by", dana.Dimension)
	}
	if text, _ := dana.Value.Text(); text != "Dana" {
		t.Errorf("level 1 value = %v, want Dana", dana.Value)
	}

	// Two children, sorted: Alex before Sam.
	if len(dana.Children) != 2 {
		t.Fatalf("Dana has %d children, want 2", len(dana.Children))
	}
	alex, sam := dana.Children[0], dana.Children[1]
	if text, _ := alex.Value.Text(); text != "Alex" {
		t.Errorf("first child = %v, want Alex", alex.Value)
	}
	if text, _ := sam.Value.Text(); text != "Sam" {
		t.Errorf("second child = %v, want Sam", sam.Value)
	}

	// Alex's detail rows, sorted by category.
	if len(alex.DetailRows) != 2 {
		t.Fatalf("Alex has %d detail rows, want 2", len(alex.DetailRows))
	}
	assertRow(t, alex.DetailRows[0].Cells, map[string]string{
		"Category": "Clothing", "Count": "4", "Subtotal": "120.00", "Hst": "15.60",
		"Total": "135.60", "AvgPerReceipt": "33.9",
	})
	// Medical has no tax at all. SUM of nothing is 0.00, not an empty cell.
	assertRow(t, alex.DetailRows[1].Cells, map[string]string{
		"Category": "Medical", "Count": "2", "Subtotal": "80.00", "Hst": "0",
		"Total": "80.00", "AvgPerReceipt": "40",
	})

	// TOTALS (Alex). Total recomputes to 200 + 15.60, which happens to equal
	// 135.60 + 80.00 because the formula is additive.
	assertRow(t, alex.Subtotals, map[string]string{
		"Count": "6", "Subtotal": "200.00", "Hst": "15.60", "Total": "215.60",
	})
	// Avg/Receipt recomputes to 215.60 / 6. Summing the column would give
	// 73.90 and averaging it 36.95; both are wrong.
	assertAvg(t, alex.Subtotals, "35.933333")
	if alex.RecordCount != 6 {
		t.Errorf("Alex RecordCount = %d, want 6", alex.RecordCount)
	}

	// Sam's detail rows.
	assertRow(t, sam.DetailRows[0].Cells, map[string]string{
		"Category": "Clothing", "Count": "3", "Subtotal": "90.00", "Hst": "11.70",
		"Total": "101.70", "AvgPerReceipt": "33.9",
	})
	assertRow(t, sam.DetailRows[1].Cells, map[string]string{
		"Category": "Mileage", "Count": "1", "Subtotal": "30.00", "Hst": "0",
		"Total": "30.00", "AvgPerReceipt": "30",
	})
	assertRow(t, sam.Subtotals, map[string]string{
		"Count": "4", "Subtotal": "120.00", "Hst": "11.70", "Total": "131.70",
	})
	assertAvg(t, sam.Subtotals, "32.925")

	// Dana's subtotal is the merge of both children.
	assertRow(t, dana.Subtotals, map[string]string{
		"Count": "10", "Subtotal": "320.00", "Hst": "27.30", "Total": "347.30",
	})

	// GRAND TOTALS.
	assertRow(t, model.GrandTotals, map[string]string{
		"Count": "10", "Subtotal": "320.00", "Hst": "27.30", "Total": "347.30",
	})
	assertAvg(t, model.GrandTotals, "34.73")
	if model.Root.RecordCount != 10 {
		t.Errorf("root RecordCount = %d, want 10", model.Root.RecordCount)
	}

	// A label column is empty on a subtotal row: no single category summed it.
	if got := cell(t, alex.Subtotals, "Category"); !got.IsNull() {
		t.Errorf("subtotal Category = %v, want null", got)
	}
	// The root is synthetic: no dimension, no subtotals of its own.
	if model.Root.Dimension != "" || model.Root.Subtotals != nil || model.Root.IsNone {
		t.Errorf("root is not synthetic: %+v", model.Root)
	}
}

func assertAvg(t *testing.T, cells []Cell, want string) {
	t.Helper()

	got := cell(t, cells, "AvgPerReceipt")
	number, isNumber := got.Decimal()
	if !isNumber {
		t.Fatalf("AvgPerReceipt = %v, want a number", got)
	}
	if !number.Equal(dec(want)) {
		t.Errorf("AvgPerReceipt = %s, want %s", number, want)
	}
}

// The document prints the average rounded to two places. Rounding is a
// renderer's job, but a template can ask for it with ROUND.
func TestRun_WorkedExampleRoundedAverage(t *testing.T) {
	spec := workedExampleSpec()
	for index, column := range spec.Columns {
		if column.Name == "AvgPerReceipt" {
			spec.Columns[index].Expr = "ROUND(Total / Count, 2)"
		}
	}

	model := mustRun(t, spec, workedExampleRows())
	dana := model.Root.Children[0]

	assertAvg(t, dana.Children[0].Subtotals, "35.93")
	assertAvg(t, dana.Children[1].Subtotals, "32.93")
	assertAvg(t, model.GrandTotals, "34.73")
}

// Averaging the children's averages would give 36.95 here; summing the column
// would give 73.90. Recomputing from the subtotal row's own inputs gives 35.93.
func TestRun_ArithmeticRecomputesRatherThanRollsUp(t *testing.T) {
	model := mustRun(t, workedExampleSpec(), workedExampleRows())
	alex := model.Root.Children[0].Children[0]

	first, _ := cell(t, alex.DetailRows[0].Cells, "AvgPerReceipt").Decimal()
	second, _ := cell(t, alex.DetailRows[1].Cells, "AvgPerReceipt").Decimal()
	subtotal, _ := cell(t, alex.Subtotals, "AvgPerReceipt").Decimal()

	if subtotal.Equal(first.Add(second)) {
		t.Errorf("subtotal average %s is the sum of the detail averages", subtotal)
	}
	if subtotal.Equal(first.Add(second).Div(dec("2"))) {
		t.Errorf("subtotal average %s is the mean of the detail averages", subtotal)
	}
	if !subtotal.Equal(dec("35.933333")) {
		t.Errorf("subtotal average = %s, want 35.933333", subtotal)
	}
}

// A receipt in two categories is attributed to both in full, so its amount is
// counted twice. This matches the dashboard pie chart, and the double count
// propagates all the way to the grand total. It is intended: do not "fix" it.
func TestRun_MultiValueDimensionsDoubleCount(t *testing.T) {
	spec := ReportSpec{
		Detail: DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns: []Column{
			{Name: "Category", Kind: ColumnLabel, Field: "category"},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		GrandTotals: true,
	}

	rows := []Row{{
		"category": {Str("Clothing"), Str("Medical")},
		"amount":   {Num(dec("100.00"))},
	}}

	model := mustRun(t, spec, rows)

	if len(model.Root.DetailRows) != 2 {
		t.Fatalf("got %d detail rows, want 2", len(model.Root.DetailRows))
	}
	assertRow(t, model.Root.DetailRows[0].Cells, map[string]string{"Category": "Clothing", "Count": "1", "Subtotal": "100.00"})
	assertRow(t, model.Root.DetailRows[1].Cells, map[string]string{"Category": "Medical", "Count": "1", "Subtotal": "100.00"})

	// One receipt of 100, but a grand total of 200 across two counts.
	assertRow(t, model.GrandTotals, map[string]string{"Count": "2", "Subtotal": "200.00"})
	if model.Root.RecordCount != 2 {
		t.Errorf("root RecordCount = %d, want 2", model.Root.RecordCount)
	}
}

// Distinct values fan out; a repeated one does not. A receipt tagged "Alex"
// twice belongs to the Alex bucket once, and its amount is counted once.
//
// This is the boundary of the rule above: fan-out double-counts across
// different buckets, never within one.
func TestRun_RepeatedDimensionValueIsOneAttribution(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"tag"},
		Columns: []Column{
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{{
		"tag":    {Str("Alex"), Str("Alex")},
		"amount": {Num(dec("10.00"))},
	}}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 1 {
		t.Fatalf("got %d tag buckets, want 1", len(model.Root.Children))
	}
	assertRow(t, model.Root.Children[0].Subtotals, map[string]string{"Count": "1", "Subtotal": "10.00"})
	assertRow(t, model.GrandTotals, map[string]string{"Count": "1", "Subtotal": "10.00"})
	if model.Root.RecordCount != 1 {
		t.Errorf("root RecordCount = %d, want 1", model.Root.RecordCount)
	}
}

// The same rule on the detail dimension.
func TestRun_RepeatedDetailValueIsOneRow(t *testing.T) {
	spec := ReportSpec{
		Detail: DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns: []Column{
			{Name: "Category", Kind: ColumnLabel, Field: "category"},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		GrandTotals: true,
	}

	rows := []Row{{
		"category": {Str("Clothing"), Str("Clothing")},
		"amount":   {Num(dec("10.00"))},
	}}

	model := mustRun(t, spec, rows)

	if len(model.Root.DetailRows) != 1 {
		t.Fatalf("got %d detail rows, want 1", len(model.Root.DetailRows))
	}
	assertRow(t, model.Root.DetailRows[0].Cells, map[string]string{"Category": "Clothing", "Count": "1", "Subtotal": "10.00"})
	assertRow(t, model.GrandTotals, map[string]string{"Count": "1", "Subtotal": "10.00"})
}

// A repeat mixed in among distinct values collapses only itself.
func TestRun_RepeatedValueAmongDistinctOnes(t *testing.T) {
	spec := ReportSpec{
		GroupBy:     []FieldKey{"tag"},
		Columns:     []Column{{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}}},
		GrandTotals: true,
	}

	rows := []Row{{
		"tag":    {Str("Alex"), Str("Sam"), Str("Alex")},
		"amount": {Num(dec("10.00"))},
	}}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d tag buckets, want 2", len(model.Root.Children))
	}
	// Two distinct buckets, so the amount is attributed twice -- not three times.
	assertRow(t, model.GrandTotals, map[string]string{"Subtotal": "20.00"})
	if model.Root.RecordCount != 2 {
		t.Errorf("root RecordCount = %d, want 2", model.Root.RecordCount)
	}
}

// Two multi-value grouping levels produce a cross product.
func TestRun_TwoMultiValueLevelsCrossProduct(t *testing.T) {
	spec := ReportSpec{
		GroupBy:     []FieldKey{"tag", "category"},
		Columns:     []Column{{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}}},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{{
		"tag":      {Str("Alex"), Str("Sam")},
		"category": {Str("Clothing"), Str("Medical")},
		"amount":   {Num(dec("10.00"))},
	}}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d tags, want 2", len(model.Root.Children))
	}
	for _, tag := range model.Root.Children {
		if len(tag.Children) != 2 {
			t.Fatalf("tag %v has %d categories, want 2", tag.Value, len(tag.Children))
		}
	}

	// 2 tags x 2 categories = 4 attributions of 10.00.
	assertRow(t, model.GrandTotals, map[string]string{"Subtotal": "40.00"})
	if model.Root.RecordCount != 4 {
		t.Errorf("root RecordCount = %d, want 4", model.Root.RecordCount)
	}
}

// A row with no value for a dimension goes into an explicit (None) bucket, and
// that bucket sorts last.
func TestRun_NoneBuckets(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"tag"},
		Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns: []Column{
			{Name: "Category", Kind: ColumnLabel, Field: "category"},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		GrandTotals: true,
	}

	rows := []Row{
		{"tag": {Str("Alex")}, "category": {Str("Clothing")}, "amount": {Num(dec("10.00"))}},
		// No tag and no category at all.
		{"amount": {Num(dec("5.00"))}},
		// An empty slice reads the same as an absent key.
		{"tag": {}, "category": {}, "amount": {Num(dec("2.00"))}},
	}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d tag buckets, want 2", len(model.Root.Children))
	}

	named, none := model.Root.Children[0], model.Root.Children[1]
	if text, _ := named.Value.Text(); text != "Alex" {
		t.Errorf("first bucket = %v, want Alex", named.Value)
	}
	if !none.IsNone || !none.Value.IsNull() {
		t.Errorf("(None) bucket did not sort last: %+v", none)
	}

	// Both untagged rows landed in the same (None) bucket, and their categories
	// collapsed into a single (None) detail row.
	if len(none.DetailRows) != 1 {
		t.Fatalf("(None) tag has %d detail rows, want 1", len(none.DetailRows))
	}
	if !none.DetailRows[0].IsNone {
		t.Errorf("detail row is not marked (None): %+v", none.DetailRows[0])
	}
	assertRow(t, none.DetailRows[0].Cells, map[string]string{"Subtotal": "7.00"})

	// Nothing was dropped.
	assertRow(t, model.GrandTotals, map[string]string{"Subtotal": "17.00"})
}

func TestRun_NoneSortsLastAmongManyBuckets(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"category"},
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}

	rows := []Row{
		{"category": {Str("Zebra")}},
		{},
		{"category": {Str("Apple")}},
		{"category": {Str("Mango")}},
	}

	model := mustRun(t, spec, rows)

	want := []string{"Apple", "Mango", "Zebra"}
	if len(model.Root.Children) != len(want)+1 {
		t.Fatalf("got %d buckets, want %d", len(model.Root.Children), len(want)+1)
	}
	for index, expected := range want {
		if text, _ := model.Root.Children[index].Value.Text(); text != expected {
			t.Errorf("bucket %d = %v, want %s", index, model.Root.Children[index].Value, expected)
		}
	}
	if !model.Root.Children[len(want)].IsNone {
		t.Errorf("last bucket is not (None)")
	}
}

// AVG must be correct at every level: at a leaf that merges detail buckets, at
// an internal node that merges child groups, and at the grand total.
//
// The tree is deliberately two levels deep. With one level the subtotal would
// come from a leaf, and an implementation that wrongly combined its children's
// finalized averages would still pass.
func TestRun_AvgIsCorrectAtEveryLevel(t *testing.T) {
	spec := ReportSpec{
		GroupBy:     []FieldKey{"paid_by", "tag"},
		Detail:      DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns:     []Column{{Name: "Avg", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: "amount"}}},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{
		// Alex: four receipts of 30 in Clothing and one of 80 in Medical.
		receiptRow("Dana", "Alex", "Clothing", "30.00", ""),
		receiptRow("Dana", "Alex", "Clothing", "30.00", ""),
		receiptRow("Dana", "Alex", "Clothing", "30.00", ""),
		receiptRow("Dana", "Alex", "Clothing", "30.00", ""),
		receiptRow("Dana", "Alex", "Medical", "80.00", ""),
		// Sam: one receipt of 120.
		receiptRow("Dana", "Sam", "Medical", "120.00", ""),
	}

	model := mustRun(t, spec, rows)
	dana := model.Root.Children[0]
	alex, sam := dana.Children[0], dana.Children[1]

	// Detail buckets: 120/4 and 80/1.
	assertRow(t, alex.DetailRows[0].Cells, map[string]string{"Avg": "30"})
	assertRow(t, alex.DetailRows[1].Cells, map[string]string{"Avg": "80"})

	// Leaf, merging its detail buckets. The mean of 30 and 80 is 55; the true
	// average is 200/5 = 40.
	assertRow(t, alex.Subtotals, map[string]string{"Avg": "40"})
	assertRow(t, sam.Subtotals, map[string]string{"Avg": "120"})

	// Internal node, merging its child groups. The mean of 40 and 120 is 80;
	// the true average is 320/6.
	assertRow(t, dana.Subtotals, map[string]string{"Avg": "53.333333"})
	assertRow(t, model.GrandTotals, map[string]string{"Avg": "53.333333"})
}

// A ratio whose denominator is zero renders as an empty cell rather than
// crashing the report.
func TestRun_DivisionByZeroIsAnEmptyCell(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			// custom_1 is null on every row, so this sums to zero.
			{Name: "Tax", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
			{Name: "Ratio", Kind: ColumnArithmetic, Expr: "Subtotal / Tax"},
		},
		GrandTotals: true,
	}

	rows := []Row{{"amount": {Num(dec("100.00"))}}}

	model := mustRun(t, spec, rows)

	assertRow(t, model.GrandTotals, map[string]string{"Subtotal": "100.00", "Tax": "0", "Ratio": ""})
}

// SUM of nothing is zero; AVG, MIN and MAX of nothing are null.
func TestRun_EmptyAndNullAggregates(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
			{Name: "Mean", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: "custom_1"}},
			{Name: "Least", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMin, Field: "custom_1"}},
			{Name: "Most", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMax, Field: "custom_1"}},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
		},
		GrandTotals: true,
	}

	t.Run("rows with no value for the measure", func(t *testing.T) {
		model := mustRun(t, spec, []Row{{"amount": {Num(dec("1"))}}, {"amount": {Num(dec("2"))}}})
		assertRow(t, model.GrandTotals, map[string]string{
			"Total": "0", "Mean": "", "Least": "", "Most": "", "Count": "2",
		})
	})

	t.Run("no rows at all", func(t *testing.T) {
		model := mustRun(t, spec, nil)
		assertRow(t, model.GrandTotals, map[string]string{
			"Total": "0", "Mean": "", "Least": "", "Most": "", "Count": "0",
		})
		if len(model.Root.Children) != 0 || len(model.Root.DetailRows) != 0 {
			t.Errorf("empty input produced rows: %+v", model.Root)
		}
		if model.Root.RecordCount != 0 {
			t.Errorf("root RecordCount = %d, want 0", model.Root.RecordCount)
		}
	})

	t.Run("no rows with grouping", func(t *testing.T) {
		grouped := spec
		grouped.GroupBy = []FieldKey{"tag"}
		grouped.Subtotals = true

		model := mustRun(t, grouped, nil)
		if len(model.Root.Children) != 0 {
			t.Errorf("empty input produced groups: %+v", model.Root.Children)
		}
		assertRow(t, model.GrandTotals, map[string]string{"Total": "0", "Count": "0", "Mean": ""})
	})
}

func TestRun_DeepNesting(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by", "tag", "category"},
		Columns: []Column{
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{
		receiptRow("Dana", "Alex", "Clothing", "10.00", ""),
		receiptRow("Dana", "Alex", "Medical", "20.00", ""),
		receiptRow("Dana", "Sam", "Clothing", "30.00", ""),
		receiptRow("Eve", "Kim", "Clothing", "40.00", ""),
	}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d payers, want 2", len(model.Root.Children))
	}
	dana, eve := model.Root.Children[0], model.Root.Children[1]

	assertRow(t, dana.Subtotals, map[string]string{"Count": "3", "Subtotal": "60.00"})
	assertRow(t, eve.Subtotals, map[string]string{"Count": "1", "Subtotal": "40.00"})

	alex := dana.Children[0]
	assertRow(t, alex.Subtotals, map[string]string{"Count": "2", "Subtotal": "30.00"})

	clothing := alex.Children[0]
	assertRow(t, clothing.Subtotals, map[string]string{"Count": "1", "Subtotal": "10.00"})
	if len(clothing.DetailRows) != 1 {
		t.Errorf("deepest level has %d detail rows, want 1", len(clothing.DetailRows))
	}

	assertRow(t, model.GrandTotals, map[string]string{"Count": "4", "Subtotal": "100.00"})
}

// Records mode emits one row per record, in the order they arrived.
func TestRun_RecordsMode(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by"},
		Detail:  DetailSpec{Mode: DetailRecords},
		Columns: []Column{
			{Name: "Cats", Kind: ColumnLabel, Field: "category"},
			{Name: "Amount", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Doubled", Kind: ColumnArithmetic, Expr: "Amount * 2"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := []Row{
		{"paid_by": {Str("Dana")}, "amount": {Num(dec("30.00"))}, "category": {Str("Zebra")}},
		{"paid_by": {Str("Dana")}, "amount": {Num(dec("10.00"))}, "category": {Str("Apple")}},
		{"paid_by": {Str("Dana")}, "amount": {Num(dec("20.00"))}, "category": {Str("Mango"), Str("Apple")}},
	}

	model := mustRun(t, spec, rows)
	dana := model.Root.Children[0]

	if len(dana.DetailRows) != 3 {
		t.Fatalf("got %d detail rows, want 3", len(dana.DetailRows))
	}

	// Input order, not sorted by anything: the caller's query decides.
	wantAmounts := []string{"30.00", "10.00", "20.00"}
	for index, want := range wantAmounts {
		assertRow(t, dana.DetailRows[index].Cells, map[string]string{"Amount": want, "Count": "1"})
	}

	// An aggregate on a single-record row aggregates that one record.
	assertRow(t, dana.DetailRows[0].Cells, map[string]string{"Doubled": "60.00"})

	// A label cell keeps every value the record carried.
	multi := dana.DetailRows[2].Cells
	for _, candidate := range multi {
		if candidate.Column != "Cats" {
			continue
		}
		if len(candidate.Values) != 2 {
			t.Fatalf("multi-value label cell has %d values, want 2", len(candidate.Values))
		}
		if first, _ := candidate.Values[0].Text(); first != "Mango" {
			t.Errorf("label values were reordered: %v", candidate.Values)
		}
	}

	// A record contributes once, even when it carries two categories: the
	// fan-out only applies to dimensions being grouped on.
	assertRow(t, dana.Subtotals, map[string]string{"Amount": "60.00", "Count": "3", "Doubled": "120.00"})
	assertRow(t, model.GrandTotals, map[string]string{"Amount": "60.00", "Count": "3"})

	// Detail rows in records mode are not keyed on anything.
	if dana.DetailRows[0].Dimension != "" || !dana.DetailRows[0].Value.IsNull() {
		t.Errorf("records-mode detail row is keyed: %+v", dana.DetailRows[0])
	}
}

// Measuring a multi-valued field is refused by Validate; displaying one shows
// every value it holds, and drops nothing.
func TestRun_MultiValuedLabelShowsEveryValue(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{{Name: "Amounts", Kind: ColumnLabel, Field: "item_amounts"}},
	}

	rows := []Row{{"item_amounts": {Num(dec("10.00")), Num(dec("90.00"))}}}
	model := mustRun(t, spec, rows)

	values := model.Root.DetailRows[0].Cells[0].Values
	if len(values) != 2 {
		t.Fatalf("label cell holds %d values, want 2", len(values))
	}
	for index, want := range []string{"10", "90"} {
		number, isNumber := values[index].Decimal()
		if !isNumber || !number.Equal(dec(want)) {
			t.Errorf("value %d = %v, want %s", index, values[index], want)
		}
	}
}

// A cell must not alias the caller's row.
func TestRun_RecordsModeDoesNotAliasInputRows(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{{Name: "Cats", Kind: ColumnLabel, Field: "category"}},
	}

	row := Row{"category": {Str("Clothing")}}
	model := mustRun(t, spec, []Row{row})

	model.Root.DetailRows[0].Cells[0].Values[0] = Str("Tampered")

	if text, _ := row["category"][0].Text(); text != "Clothing" {
		t.Errorf("mutating a cell changed the caller's row: %v", row)
	}
}

// A label column on an aggregated detail row shows the group it sits under.
func TestRun_LabelColumnReadsAncestorBucket(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by", "tag"},
		Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns: []Column{
			{Name: "PaidBy", Kind: ColumnLabel, Field: "paid_by"},
			{Name: "Child", Kind: ColumnLabel, Field: "tag"},
			{Name: "Category", Kind: ColumnLabel, Field: "category"},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
	}

	rows := []Row{
		receiptRow("Dana", "Alex", "Clothing", "10.00", ""),
		receiptRow("Eve", "Kim", "Medical", "20.00", ""),
	}

	model := mustRun(t, spec, rows)

	dana := model.Root.Children[0].Children[0]
	assertRow(t, dana.DetailRows[0].Cells, map[string]string{
		"PaidBy": "Dana", "Child": "Alex", "Category": "Clothing", "Subtotal": "10.00",
	})

	eve := model.Root.Children[1].Children[0]
	assertRow(t, eve.DetailRows[0].Cells, map[string]string{
		"PaidBy": "Eve", "Child": "Kim", "Category": "Medical", "Subtotal": "20.00",
	})
}

func TestRun_SingleRecord(t *testing.T) {
	spec := workedExampleSpec()
	rows := []Row{receiptRow("Dana", "Alex", "Clothing", "50.00", "6.50")}

	model := mustRun(t, spec, rows)
	alex := model.Root.Children[0].Children[0]

	want := map[string]string{"Count": "1", "Subtotal": "50.00", "Hst": "6.50", "Total": "56.50", "AvgPerReceipt": "56.50"}
	assertRow(t, alex.DetailRows[0].Cells, want)
	assertRow(t, alex.Subtotals, want)
	assertRow(t, model.GrandTotals, want)
}

func TestRun_TogglesOffSubtotalsAndGrandTotals(t *testing.T) {
	spec := workedExampleSpec()
	spec.Subtotals = false
	spec.GrandTotals = false

	model := mustRun(t, spec, workedExampleRows())

	if model.GrandTotals != nil {
		t.Errorf("GrandTotals = %v, want nil", model.GrandTotals)
	}

	dana := model.Root.Children[0]
	if dana.Subtotals != nil {
		t.Errorf("group Subtotals = %v, want nil", dana.Subtotals)
	}
	if dana.Children[0].Subtotals != nil {
		t.Errorf("leaf Subtotals = %v, want nil", dana.Children[0].Subtotals)
	}
	// Detail rows are unaffected.
	if len(dana.Children[0].DetailRows) != 2 {
		t.Errorf("detail rows were dropped along with the subtotals")
	}
}

// The same inputs must always produce the same output, byte for byte.
func TestRun_IsDeterministic(t *testing.T) {
	spec := workedExampleSpec()
	rows := workedExampleRows()

	first := mustRun(t, spec, rows)
	second := mustRun(t, spec, rows)

	if !reflect.DeepEqual(first.Root, second.Root) {
		t.Errorf("two runs produced different trees")
	}
	if !reflect.DeepEqual(first.GrandTotals, second.GrandTotals) {
		t.Errorf("two runs produced different grand totals")
	}
}

// Aggregated output must not depend on the order rows arrived in, because it is
// grouped and sorted rather than listed.
func TestRun_AggregateOutputIsIndependentOfInputOrder(t *testing.T) {
	spec := workedExampleSpec()

	forward := workedExampleRows()
	reversed := workedExampleRows()
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}

	if !reflect.DeepEqual(mustRun(t, spec, forward).Root, mustRun(t, spec, reversed).Root) {
		t.Errorf("shuffling the input changed the report")
	}
}

// An arithmetic column may be declared before the aggregate it reads.
func TestRun_ColumnOrderIsIndependentOfDependencyOrder(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Avg", Kind: ColumnArithmetic, Expr: "Total / Count"},
			{Name: "Total", Kind: ColumnArithmetic, Expr: "Subtotal + Hst"},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Hst", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
		},
		GrandTotals: true,
	}

	rows := []Row{
		receiptRow("Dana", "Alex", "Clothing", "60.00", "6.00"),
		receiptRow("Dana", "Alex", "Clothing", "40.00", "4.00"),
	}

	model := mustRun(t, spec, rows)

	// Cells stay in declaration order even though evaluation does not.
	if model.Columns[0].Name != "Avg" {
		t.Errorf("columns were reordered: %s", model.Columns[0].Name)
	}
	assertRow(t, model.GrandTotals, map[string]string{"Total": "110.00", "Avg": "55", "Count": "2"})
}

func TestRun_Descriptors(t *testing.T) {
	model := mustRun(t, workedExampleSpec(), workedExampleRows())

	byName := make(map[string]ColumnDescriptor, len(model.Columns))
	for _, descriptor := range model.Columns {
		byName[descriptor.Name] = descriptor
	}

	if got := byName["AvgPerReceipt"].Label; got != "Avg/Receipt" {
		t.Errorf("label = %q, want Avg/Receipt", got)
	}
	if got := byName["Subtotal"].Label; got != "Subtotal" {
		t.Errorf("label defaults to the name, got %q", got)
	}

	// An arithmetic column carries a walkable tree, which is what the
	// spreadsheet renderer will translate into a live formula.
	total := byName["Total"]
	if total.Expr == nil {
		t.Errorf("arithmetic column has no parsed expression")
	}
	if total.ExprSrc != "Subtotal + Hst" {
		t.Errorf("ExprSrc = %q", total.ExprSrc)
	}
	if total.Agg != nil {
		t.Errorf("arithmetic column carries an aggregate")
	}
	if refs := columnRefs(total.Expr); len(refs) != 2 {
		t.Errorf("columnRefs(Total) = %v, want two references", refs)
	}

	// An aggregate column carries its reduction.
	subtotal := byName["Subtotal"]
	if subtotal.Agg == nil || subtotal.Agg.Func != AggSum || subtotal.Agg.Field != "amount" {
		t.Errorf("Subtotal.Agg = %+v, want SUM(amount)", subtotal.Agg)
	}
	if subtotal.Expr != nil {
		t.Errorf("aggregate column carries an expression")
	}

	// Money combined with a count is still money.
	if byName["AvgPerReceipt"].DataType != TypeCurrency {
		t.Errorf("AvgPerReceipt is not currency")
	}
	if byName["Count"].DataType != TypeNumber {
		t.Errorf("Count is not a plain number")
	}
	if byName["Category"].DataType != TypeString {
		t.Errorf("Category is not a string")
	}
}

// Money keeps every cent through grouping and rollup.
func TestRun_NoFloatDriftThroughTheTree(t *testing.T) {
	spec := ReportSpec{
		GroupBy:     []FieldKey{"tag"},
		Columns:     []Column{{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}}},
		Subtotals:   true,
		GrandTotals: true,
	}

	rows := make([]Row, 0, 300)
	for index := 0; index < 300; index++ {
		tag := "Alex"
		if index%2 == 0 {
			tag = "Sam"
		}
		rows = append(rows, Row{"tag": {Str(tag)}, "amount": {Num(dec("0.01"))}})
	}

	model := mustRun(t, spec, rows)
	assertRow(t, model.GrandTotals, map[string]string{"Subtotal": "3.00"})
	for _, child := range model.Root.Children {
		assertRow(t, child.Subtotals, map[string]string{"Subtotal": "1.50"})
	}
}

func TestRun_RejectsAnInvalidSpec(t *testing.T) {
	_, err := Run(ReportSpec{}, testCatalog(t), nil, testMeta())
	if err == nil {
		t.Fatalf("Run() with no columns did not error")
	}
}

func TestRun_MetaParamsAreCopied(t *testing.T) {
	params := map[string]string{"period": "May"}
	model, err := Run(
		ReportSpec{Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
		testCatalog(t), nil, MetaInput{GeneratedAt: testGeneratedAt, Params: params},
	)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	params["period"] = "June"
	if model.Meta.Params["period"] != "May" {
		t.Errorf("Meta.Params aliases the caller's map")
	}
}

// Grouping on a date must not split one instant across two buckets because the
// caller happened to hand it over in a different time zone.
func TestRun_DateBucketsNormalizeAcrossZones(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"date"},
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}

	utc := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	elsewhere := utc.In(time.FixedZone("plus-one", 3600))

	model := mustRun(t, spec, []Row{
		{"date": {DateVal(utc)}},
		{"date": {DateVal(elsewhere)}},
	})

	if len(model.Root.Children) != 1 {
		t.Fatalf("got %d date buckets, want 1", len(model.Root.Children))
	}
	if model.Root.Children[0].RecordCount != 2 {
		t.Errorf("RecordCount = %d, want 2", model.Root.Children[0].RecordCount)
	}
}

// false sorts before true, and both are their own bucket.
func TestRun_BooleanBucketing(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"resolved"},
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}

	model := mustRun(t, spec, []Row{
		{"resolved": {Bool(true)}},
		{"resolved": {Bool(false)}},
		{"resolved": {Bool(true)}},
	})

	if len(model.Root.Children) != 2 {
		t.Fatalf("got %d buckets, want 2", len(model.Root.Children))
	}
	if value, _ := model.Root.Children[0].Value.Boolean(); value {
		t.Errorf("true sorted before false")
	}
	if model.Root.Children[0].RecordCount != 1 {
		t.Errorf("false bucket RecordCount = %d, want 1", model.Root.Children[0].RecordCount)
	}
	if model.Root.Children[1].RecordCount != 2 {
		t.Errorf("true bucket RecordCount = %d, want 2", model.Root.Children[1].RecordCount)
	}
}

// Two amounts that differ only in scale are the same number, and grouping is
// keyed on the number rather than on how it was written.
func TestRun_NumbersDifferingOnlyInScaleShareABucket(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"category"},
		Columns: []Column{{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}}},
	}

	// bucketKey normalizes numbers even though a measure can never be grouped
	// on; this checks the normalization directly.
	if bucketKey(Num(dec("200"))) != bucketKey(Num(dec("200.00"))) {
		t.Errorf("200 and 200.00 landed in different buckets")
	}

	model := mustRun(t, spec, []Row{{"category": {Str("Clothing")}, "amount": {Num(dec("200.00"))}}})
	assertRow(t, model.Root.Children[0].DetailRows[0].Cells, map[string]string{"Total": "200"})
}

// A currency custom field is a measure, but measuring is the only thing its type
// restricts: a report may cut by it too. Buckets key on the canonical decimal, so
// 15.60 and 15.6 are one bucket, and a receipt carrying no value lands in (None).
func TestRun_GroupByACurrencyField(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"custom_1"},
		Columns: []Column{
			{Name: "Hst", Kind: ColumnLabel, Field: "custom_1"},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		},
		GrandTotals: true,
	}

	rows := []Row{
		receiptRow("Dana", "Alex", "Clothing", "100.00", "15.60"),
		receiptRow("Sam", "Alex", "Medical", "50.00", "15.6"),
		receiptRow("Dana", "Sam", "Clothing", "20.00", "2.00"),
		receiptRow("Sam", "Sam", "Medical", "10.00", ""),
	}

	model := mustRun(t, spec, rows)

	if len(model.Root.Children) != 3 {
		t.Fatalf("got %d buckets, want 3 (2.00, 15.60, (None))", len(model.Root.Children))
	}

	two, fifteen, none := model.Root.Children[0], model.Root.Children[1], model.Root.Children[2]

	if number, isNumber := two.Value.Decimal(); !isNumber || !number.Equal(dec("2.00")) {
		t.Errorf("first bucket = %v, want 2.00 (numbers sort ascending)", two.Value)
	}
	// 15.60 and 15.6 are the same amount, so they share one bucket.
	if number, _ := fifteen.Value.Decimal(); !number.Equal(dec("15.60")) || fifteen.RecordCount != 2 {
		t.Errorf("second bucket = %v with %d records, want 15.60 with 2", fifteen.Value, fifteen.RecordCount)
	}
	if !none.IsNone || none.RecordCount != 1 {
		t.Errorf("(None) bucket = %+v, want 1 record", none)
	}

	// The label column reads the bucket, so a grouped currency value is displayable.
	assertRow(t, fifteen.DetailRows[0].Cells, map[string]string{"Hst": "15.60"})
	assertRow(t, model.GrandTotals, map[string]string{"Count": "4", "Total": "180.00"})
}
