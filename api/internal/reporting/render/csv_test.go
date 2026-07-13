package render

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/reporting"

	"github.com/shopspring/decimal"
)

// --- fixtures -------------------------------------------------------------

func num(v int64) reporting.Value { return reporting.Num(decimal.NewFromInt(v)) }

func money(s string) reporting.Value { return reporting.Num(decimal.RequireFromString(s)) }

func mustCatalog(t *testing.T, fields ...reporting.FieldRef) reporting.FieldCatalog {
	t.Helper()
	catalog, err := reporting.NewFieldCatalog(fields...)
	if err != nil {
		t.Fatalf("build catalog: %v", err)
	}
	return catalog
}

func mustRun(t *testing.T, spec reporting.ReportSpec, catalog reporting.FieldCatalog, rows []reporting.Row) reporting.ReportModel {
	t.Helper()
	model, err := reporting.Run(spec, catalog, rows, reporting.MetaInput{})
	if err != nil {
		t.Fatalf("run report: %v", err)
	}
	return model
}

// crlf joins the given lines with the writer's RFC 4180 line ending and a
// trailing terminator, which is what encoding/csv emits with UseCRLF set.
func crlf(lines ...string) string {
	return strings.Join(lines, "\r\n") + "\r\n"
}

func mustCSV(t *testing.T, model reporting.ReportModel, groupBy []Dimension) string {
	t.Helper()
	out, err := CSV(model, groupBy)
	if err != nil {
		t.Fatalf("CSV: %v", err)
	}
	return string(out)
}

// paidByCatalog is the common set of fields the grouped tests reference.
func paidByCatalog(t *testing.T) reporting.FieldCatalog {
	t.Helper()
	return mustCatalog(t,
		reporting.FieldRef{Key: "paid_by", Label: "Paid By", DataType: reporting.TypeString},
		reporting.FieldRef{Key: "tag", Label: "Tag", DataType: reporting.TypeString, Multi: true},
		reporting.FieldRef{Key: "category", Label: "Category", DataType: reporting.TypeString, Multi: true},
		reporting.FieldRef{Key: "amount", Label: "Amount", DataType: reporting.TypeCurrency},
	)
}

// --- headline: sums, aggregates, and non-linear roll-up -------------------

// A one-level grouped report with a COUNT, a SUM, and an arithmetic average
// proves the whole roll-up: additive columns (Count, Total) sum up — a subtotal
// is the sum of its detail rows and the grand total the sum of them all — while
// the non-linear Avg is recomputed at every level (120.00 = 360/3 and 100.00 =
// 400/4), which is neither the sum nor the average of the rows above it.
func TestCSV_SumsAggregatesAndNonLinearRollUp(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Count", Label: "Count", Kind: reporting.ColumnAggregate, AggSrc: "COUNT()"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
			{Name: "Avg", Label: "Avg", Kind: reporting.ColumnArithmetic, Expr: "Total / Count"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("200")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Gas")}, "amount": {money("60")}},
		{"paid_by": {reporting.Str("Sam")}, "category": {reporting.Str("Food")}, "amount": {money("40")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), []Dimension{{Key: "paid_by", Label: "Paid By"}})
	want := crlf(
		"Row Type,Paid By,Category,Count,Total,Avg",
		"Detail,Dana,Food,2,300.00,150.00",
		"Detail,Dana,Gas,1,60.00,60.00",
		"Subtotal,Dana,,3,360.00,120.00",
		"Detail,Sam,Food,1,40.00,40.00",
		"Subtotal,Sam,,1,40.00,40.00",
		"Grand Total,,,4,400.00,100.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// Two grouping levels prove nested subtotals: each level emits its own subtotal
// row after the rows it sums, the leading dimension columns fill to that level
// and blank below it, and the ordering is deterministic (values ascending).
func TestCSV_NestedGroupsSubtotalsAndPadding(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by", "tag"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "tag": {reporting.Str("Alex")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"paid_by": {reporting.Str("Dana")}, "tag": {reporting.Str("Alex")}, "category": {reporting.Str("Gas")}, "amount": {money("50")}},
		{"paid_by": {reporting.Str("Dana")}, "tag": {reporting.Str("Bob")}, "category": {reporting.Str("Food")}, "amount": {money("30")}},
		{"paid_by": {reporting.Str("Sam")}, "tag": {reporting.Str("Alex")}, "category": {reporting.Str("Food")}, "amount": {money("20")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows),
		[]Dimension{{Key: "paid_by", Label: "Paid By"}, {Key: "tag", Label: "Tag"}})
	want := crlf(
		"Row Type,Paid By,Tag,Category,Total",
		"Detail,Dana,Alex,Food,100.00",
		"Detail,Dana,Alex,Gas,50.00",
		"Subtotal,Dana,Alex,,150.00",
		"Detail,Dana,Bob,Food,30.00",
		"Subtotal,Dana,Bob,,30.00",
		"Subtotal,Dana,,,180.00",
		"Detail,Sam,Alex,Food,20.00",
		"Subtotal,Sam,Alex,,20.00",
		"Subtotal,Sam,,,20.00",
		"Grand Total,,,,200.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// Every aggregate function renders at detail, subtotal, and grand-total level.
// The single group keeps its subtotal equal to the grand total, so one set of
// expected values covers both roll-up levels.
func TestCSV_AllAggregateFunctions(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Sum", Label: "Sum", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
			{Name: "Count", Label: "Count", Kind: reporting.ColumnAggregate, AggSrc: "COUNT()"},
			{Name: "Avg", Label: "Avg", Kind: reporting.ColumnAggregate, AggSrc: "AVG(amount)"},
			{Name: "Min", Label: "Min", Kind: reporting.ColumnAggregate, AggSrc: "MIN(amount)"},
			{Name: "Max", Label: "Max", Kind: reporting.ColumnAggregate, AggSrc: "MAX(amount)"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("10")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("30")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Gas")}, "amount": {money("110")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), []Dimension{{Key: "paid_by", Label: "Paid By"}})
	// SUM/MIN/MAX/AVG follow their measure (currency, 2dp); COUNT is a plain number.
	want := crlf(
		"Row Type,Paid By,Category,Sum,Count,Avg,Min,Max",
		"Detail,Dana,Food,40.00,2,20.00,10.00,30.00",
		"Detail,Dana,Gas,110.00,1,110.00,110.00,110.00",
		"Subtotal,Dana,,150.00,3,50.00,10.00,110.00",
		"Grand Total,,,150.00,3,50.00,10.00,110.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// --- modes and shapes -----------------------------------------------------

// Records mode emits one row per source record in input order, and a label cell
// carrying several values (a receipt with two categories) is joined with commas
// — which forces the field to be quoted.
func TestCSV_RecordsModeMultiValueLabel(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailRecords},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food"), reporting.Str("Gas")}, "amount": {money("100")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("50")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), []Dimension{{Key: "paid_by", Label: "Paid By"}})
	want := crlf(
		"Row Type,Paid By,Category,Total",
		`Detail,Dana,"Food, Gas",100.00`,
		"Detail,Dana,Food,50.00",
		"Subtotal,Dana,,150.00",
		"Grand Total,,,150.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// With no grouping the detail rows hang off the root and there are no leading
// dimension columns.
func TestCSV_NoGrouping(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		Detail: reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"category": {reporting.Str("Gas")}, "amount": {money("50")}},
		{"category": {reporting.Str("Food")}, "amount": {money("25")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), nil)
	want := crlf(
		"Row Type,Category,Total",
		"Detail,Food,125.00",
		"Detail,Gas,50.00",
		"Grand Total,,175.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// With subtotals and grand totals off, every row is a Detail — but the Row Type
// column is still present, so the schema is stable regardless of the toggles.
func TestCSV_NoTotalsStillHasRowTypeColumn(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		Detail: reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"category": {reporting.Str("Gas")}, "amount": {money("50")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), nil)
	want := crlf(
		"Row Type,Category,Total",
		"Detail,Food,100.00",
		"Detail,Gas,50.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// A row with no value for the grouping dimension falls into the (None) bucket,
// which sorts last and prints the report's name for it in the dimension column.
func TestCSV_NoneBucket(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("50")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), []Dimension{{Key: "paid_by", Label: "Paid By"}})
	want := crlf(
		"Row Type,Paid By,Category,Total",
		"Detail,Dana,Food,50.00",
		"Subtotal,Dana,,50.00",
		"Detail,(None),Food,100.00",
		"Subtotal,(None),,100.00",
		"Grand Total,,,150.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// A multi-value grouping dimension fans a row out into every bucket, in full, so
// one $100 receipt tagged twice contributes $100 to each tag and $200 to the
// grand total — the intended double count, surfaced honestly in the CSV.
func TestCSV_MultiValueGroupFanOut(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"tag"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		GrandTotals: true,
	}
	rows := []reporting.Row{
		{"tag": {reporting.Str("Alex"), reporting.Str("Bob")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), []Dimension{{Key: "tag", Label: "Tag"}})
	want := crlf(
		"Row Type,Tag,Category,Total",
		"Detail,Alex,Food,100.00",
		"Detail,Bob,Food,100.00",
		"Grand Total,,,200.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// An empty report is just the header, plus a grand-total row of blanks (and a
// zero count) when grand totals are on.
func TestCSV_EmptyReportHeaderAndGrandTotalOnly(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Count", Label: "Count", Kind: reporting.ColumnAggregate, AggSrc: "COUNT()"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		GrandTotals: true,
	}

	// SUM over no rows is zero (not null), and COUNT is zero, so the grand total
	// row is fully populated even with an empty report.
	got := mustCSV(t, mustRun(t, spec, catalog, nil), []Dimension{{Key: "paid_by", Label: "Paid By"}})
	want := crlf(
		"Row Type,Paid By,Category,Count,Total",
		"Grand Total,,,0,0.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// A dimension supplied without a label falls back to its raw field key, so a
// column is never blank-headed.
func TestCSV_MissingDimensionLabelFallsBackToKey(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows), []Dimension{{Key: "paid_by"}})
	if first, _, _ := strings.Cut(got, "\r\n"); first != "Row Type,paid_by,Category,Total" {
		t.Errorf("header = %q, want the raw key as the dimension heading", first)
	}
}

// --- escaping and unicode -------------------------------------------------

// Values carrying commas, quotes, and newlines survive a round trip through a
// CSV reader unchanged, and unicode is preserved — the writer quotes and escapes
// them, so the file is safe to re-parse.
func TestCSV_EscapingAndUnicodeRoundTrip(t *testing.T) {
	catalog := paidByCatalog(t)
	spec := reporting.ReportSpec{
		Detail: reporting.DetailSpec{Mode: reporting.DetailRecords},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"category": {reporting.Str("Groceries, Inc")}, "amount": {money("10")}},
		{"category": {reporting.Str(`Say "hi"`)}, "amount": {money("20")}},
		{"category": {reporting.Str("line1\nline2")}, "amount": {money("30")}},
		{"category": {reporting.Str("日本語 🎉")}, "amount": {money("40")}},
	}

	out, err := CSV(mustRun(t, spec, catalog, rows), nil)
	if err != nil {
		t.Fatalf("CSV: %v", err)
	}

	got, err := csv.NewReader(strings.NewReader(string(out))).ReadAll()
	if err != nil {
		t.Fatalf("re-parse CSV: %v", err)
	}
	want := [][]string{
		{"Row Type", "Category", "Total"},
		{"Detail", "Groceries, Inc", "10.00"},
		{"Detail", `Say "hi"`, "20.00"},
		{"Detail", "line1\nline2", "30.00"},
		{"Detail", "日本語 🎉", "40.00"},
	}
	if len(got) != len(want) {
		t.Fatalf("row count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if strings.Join(got[i], "\x00") != strings.Join(want[i], "\x00") {
			t.Errorf("row %d = %q, want %q", i, got[i], want[i])
		}
	}
}

// --- cell formatting ------------------------------------------------------

func TestFormatCell_ByTypeAndKind(t *testing.T) {
	currency := reporting.ColumnDescriptor{Kind: reporting.ColumnAggregate, DataType: reporting.TypeCurrency}
	number := reporting.ColumnDescriptor{Kind: reporting.ColumnAggregate, DataType: reporting.TypeNumber}
	label := reporting.ColumnDescriptor{Kind: reporting.ColumnLabel, DataType: reporting.TypeString}

	date := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name   string
		column reporting.ColumnDescriptor
		cell   reporting.Cell
		want   string
	}{
		{"currency two places", currency, reporting.Cell{Values: []reporting.Value{money("1.5")}}, "1.50"},
		{"currency rounds", currency, reporting.Cell{Values: []reporting.Value{money("1234.5")}}, "1234.50"},
		{"plain number full precision", number, reporting.Cell{Values: []reporting.Value{money("1.5")}}, "1.5"},
		{"plain integer", number, reporting.Cell{Values: []reporting.Value{num(2)}}, "2"},
		{"string label", label, reporting.Cell{Values: []reporting.Value{reporting.Str("Food")}}, "Food"},
		{"date label utc", label, reporting.Cell{Values: []reporting.Value{reporting.DateVal(date)}}, "2026-05-15T00:00:00Z"},
		{"bool label", label, reporting.Cell{Values: []reporting.Value{reporting.Bool(true)}}, "true"},
		{"null measure is blank", currency, reporting.Cell{Values: []reporting.Value{reporting.Null()}}, ""},
		{"null label is none", label, reporting.Cell{Values: []reporting.Value{reporting.Null()}}, "(None)"},
		{"no values is blank", label, reporting.Cell{}, ""},
		{"multi-value label joined", label, reporting.Cell{Values: []reporting.Value{reporting.Str("A"), reporting.Str("B")}}, "A, B"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := formatCell(testCase.column, testCase.cell, "(None)"); got != testCase.want {
				t.Errorf("formatCell = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestFormatDimension_NullIsNoneLabel(t *testing.T) {
	if got := formatDimension(reporting.Null(), "(None)"); got != "(None)" {
		t.Errorf("null dimension = %q, want (None)", got)
	}
	if got := formatDimension(reporting.Str("Dana"), "(None)"); got != "Dana" {
		t.Errorf("string dimension = %q, want Dana", got)
	}
}
