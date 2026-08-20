package render

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"receipt-wrangler/api/internal/reporting"

	"github.com/xuri/excelize/v2"
)

// --- helpers --------------------------------------------------------------

// xlsxGrid renders the model, reopens the bytes, and returns the "Report"
// sheet's cells as excelize formats them for display (currency to two places,
// counts as whole numbers, blanks empty).
func xlsxGrid(t *testing.T, model reporting.ReportModel, groupBy []Dimension) [][]string {
	t.Helper()
	out, err := XLSX(model, groupBy)
	if err != nil {
		t.Fatalf("XLSX: %v", err)
	}
	file, err := excelize.OpenReader(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	defer file.Close()

	rows, err := file.GetRows(xlsxSheet)
	if err != nil {
		t.Fatalf("read rows: %v", err)
	}
	return rows
}

// assertGrid compares two grids after trimming trailing empty cells from each
// row, so it does not depend on how excelize pads short rows.
func assertGrid(t *testing.T, got, want [][]string) {
	t.Helper()
	trim := func(grid [][]string) [][]string {
		out := make([][]string, len(grid))
		for index, row := range grid {
			end := len(row)
			for end > 0 && row[end-1] == "" {
				end--
			}
			out[index] = row[:end]
		}
		return out
	}
	if g, w := trim(got), trim(want); !reflect.DeepEqual(g, w) {
		t.Errorf("grid mismatch\n--- got ---\n%v\n--- want ---\n%v", g, w)
	}
}

func cellStyle(t *testing.T, file *excelize.File, cell string) *excelize.Style {
	t.Helper()
	id, err := file.GetCellStyle(xlsxSheet, cell)
	if err != nil {
		t.Fatalf("cell style %s: %v", cell, err)
	}
	style, err := file.GetStyle(id)
	if err != nil {
		t.Fatalf("get style %s: %v", cell, err)
	}
	return style
}

func isBold(style *excelize.Style) bool {
	return style.Font != nil && style.Font.Bold
}

// oneLevelSpec is the shared fixture: group by paid_by, aggregate by category,
// with a COUNT, a SUM, and an arithmetic average.
func oneLevelSpec() reporting.ReportSpec {
	return reporting.ReportSpec{
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
}

func oneLevelRows() []reporting.Row {
	return []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("200")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Gas")}, "amount": {money("60")}},
		{"paid_by": {reporting.Str("Sam")}, "category": {reporting.Str("Food")}, "amount": {money("40")}},
	}
}

func paidByDimension() []Dimension { return []Dimension{{Key: "paid_by", Label: "Paid By"}} }

// typedCatalog is the fixture for the non-string dimensions: a boolean, a date,
// and a currency custom field, each of which a report may group by.
func typedCatalog(t *testing.T) reporting.FieldCatalog {
	t.Helper()
	return mustCatalog(t,
		reporting.FieldRef{Key: "reimbursable", Label: "Reimbursable", DataType: reporting.TypeBool},
		reporting.FieldRef{Key: "due", Label: "Due Date", DataType: reporting.TypeDate},
		reporting.FieldRef{Key: "custom_1", Label: "HST", DataType: reporting.TypeCurrency},
		reporting.FieldRef{Key: "category", Label: "Category", DataType: reporting.TypeString, Multi: true},
		reporting.FieldRef{Key: "amount", Label: "Amount", DataType: reporting.TypeCurrency},
	)
}

// --- headline: sums, aggregates, non-linear roll-up -----------------------

// The dimension shows once per group, subtotal/grand-total rows carry a Total
// marker, additive columns roll up (subtotal = sum of its details, grand total =
// sum of all), and the non-linear Avg is recomputed per level (120.00 = 360/3,
// 100.00 = 400/4) — displayed as native numbers with a currency format.
func TestXLSX_SumsAggregatesAndNonLinearRollUp(t *testing.T) {
	model := mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows())

	assertGrid(t, xlsxGrid(t, model, paidByDimension()), [][]string{
		{"Paid By", "Category", "Count", "Total", "Avg"},
		{"Dana", "Food", "2", "300.00", "150.00"},
		{"", "Gas", "1", "60.00", "60.00"},
		{"", "Total", "3", "360.00", "120.00"},
		{"Sam", "Food", "1", "40.00", "40.00"},
		{"", "Total", "1", "40.00", "40.00"},
		{"Grand Total", "", "4", "400.00", "100.00"},
	})
}

// Two grouping levels prove the staircase (a group's Total marker sits in the
// column at its depth) and dimension blanking (a value shows once per group).
func TestXLSX_NestedGroupsStaircaseAndBlanking(t *testing.T) {
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

	grid := xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows),
		[]Dimension{{Key: "paid_by", Label: "Paid By"}, {Key: "tag", Label: "Tag"}})
	assertGrid(t, grid, [][]string{
		{"Paid By", "Tag", "Category", "Total"},
		{"Dana", "Alex", "Food", "100.00"},
		{"", "", "Gas", "50.00"},
		{"", "", "Total", "150.00"},
		{"", "Bob", "Food", "30.00"},
		{"", "", "Total", "30.00"},
		{"", "Total", "", "180.00"},
		{"Sam", "Alex", "Food", "20.00"},
		{"", "", "Total", "20.00"},
		{"", "Total", "", "20.00"},
		{"Grand Total", "", "", "200.00"},
	})
}

// Each aggregate function renders at detail and roll-up level; the single group
// keeps its subtotal equal to the grand total.
func TestXLSX_AllAggregateFunctions(t *testing.T) {
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

	grid := xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension())
	assertGrid(t, grid, [][]string{
		{"Paid By", "Category", "Sum", "Count", "Avg", "Min", "Max"},
		{"Dana", "Food", "40.00", "2", "20.00", "10.00", "30.00"},
		{"", "Gas", "110.00", "1", "110.00", "110.00", "110.00"},
		{"", "Total", "150.00", "3", "50.00", "10.00", "110.00"},
		{"Grand Total", "", "150.00", "3", "50.00", "10.00", "110.00"},
	})
}

// Records mode emits one row per record; a multi-value label is joined with
// commas into a single cell (no CSV-style escaping needed in a spreadsheet).
func TestXLSX_RecordsModeMultiValueLabel(t *testing.T) {
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

	grid := xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension())
	assertGrid(t, grid, [][]string{
		{"Paid By", "Category", "Total"},
		{"Dana", "Food, Gas", "100.00"},
		{"", "Food", "50.00"},
		{"", "Total", "150.00"},
		{"Grand Total", "", "150.00"},
	})
}

// With no grouping the details hang off the root and there are no dimension
// columns; the grand total marker lands in the first report column.
func TestXLSX_NoGrouping(t *testing.T) {
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

	assertGrid(t, xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), nil), [][]string{
		{"Category", "Total"},
		{"Food", "125.00"},
		{"Gas", "50.00"},
		{"Grand Total", "175.00"},
	})
}

// Subtotals on, grand totals off: subtotal rows appear, no grand-total row.
func TestXLSX_SubtotalsWithoutGrandTotal(t *testing.T) {
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
		Subtotals: true,
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"paid_by": {reporting.Str("Dana")}, "category": {reporting.Str("Gas")}, "amount": {money("50")}},
		{"paid_by": {reporting.Str("Sam")}, "category": {reporting.Str("Food")}, "amount": {money("40")}},
	}

	assertGrid(t, xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension()), [][]string{
		{"Paid By", "Category", "Total"},
		{"Dana", "Food", "100.00"},
		{"", "Gas", "50.00"},
		{"", "Total", "150.00"},
		{"Sam", "Food", "40.00"},
		{"", "Total", "40.00"},
	})
}

// A row with no value for the grouping dimension is the (None) bucket; it sorts
// last and prints the report's name for it.
func TestXLSX_NoneBucket(t *testing.T) {
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

	assertGrid(t, xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension()), [][]string{
		{"Paid By", "Category", "Total"},
		{"Dana", "Food", "50.00"},
		{"", "Total", "50.00"},
		{"(None)", "Food", "100.00"},
		{"", "Total", "100.00"},
		{"Grand Total", "", "150.00"},
	})
}

// A multi-value grouping dimension fans a row into every bucket in full, so one
// receipt tagged twice contributes to both tags and doubles the grand total.
func TestXLSX_MultiValueGroupFanOut(t *testing.T) {
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

	assertGrid(t, xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), []Dimension{{Key: "tag", Label: "Tag"}}), [][]string{
		{"Tag", "Category", "Total"},
		{"Alex", "Food", "100.00"},
		{"Bob", "Food", "100.00"},
		{"Grand Total", "", "200.00"},
	})
}

// An empty report is the header plus a grand-total row (zero count, zero sum).
func TestXLSX_EmptyReport(t *testing.T) {
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

	assertGrid(t, xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), nil), paidByDimension()), [][]string{
		{"Paid By", "Category", "Count", "Total"},
		{"Grand Total", "", "0", "0.00"},
	})
}

// Unicode passes through into a cell unchanged.
func TestXLSX_Unicode(t *testing.T) {
	spec := reporting.ReportSpec{
		Detail: reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"category": {reporting.Str("日本語 🎉")}, "amount": {money("10")}},
	}

	assertGrid(t, xlsxGrid(t, mustRun(t, spec, paidByCatalog(t), rows), nil), [][]string{
		{"Category", "Total"},
		{"日本語 🎉", "10.00"},
	})
}

// Supplying fewer dimensions than the report groups by is rejected, reusing the
// shared depth guard.
func TestXLSX_GroupByDepthMismatchErrors(t *testing.T) {
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"paid_by", "tag"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"paid_by": {reporting.Str("Dana")}, "tag": {reporting.Str("Alex")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
	}

	out, err := XLSX(mustRun(t, spec, paidByCatalog(t), rows), paidByDimension())
	if err == nil {
		t.Fatalf("expected an error for a group-depth mismatch, got none")
	}
	if out != nil {
		t.Errorf("expected nil bytes on error, got %d", len(out))
	}
}

// --- native typing, number formats, styling -------------------------------

// Numbers are written as native numeric cells (not text) with a currency number
// format, and the header, subtotal, and grand-total rows are bold while detail
// rows are not.
func TestXLSX_NativeTypesFormatsAndBold(t *testing.T) {
	out, err := XLSX(mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows()), paidByDimension())
	if err != nil {
		t.Fatalf("XLSX: %v", err)
	}
	file, err := excelize.OpenReader(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	defer file.Close()

	if name := file.GetSheetName(0); name != xlsxSheet {
		t.Errorf("sheet name = %q, want %q", name, xlsxSheet)
	}

	// D2 is the Total of the first detail row. It is a native number, not text:
	// its raw stored value is unformatted while its displayed value has the
	// currency format applied — a text cell would show the same string for both.
	// (excelize omits the type attribute for numbers, so GetCellType is not the
	// tell here.)
	raw, err := file.GetCellValue(xlsxSheet, "D2", excelize.Options{RawCellValue: true})
	if err != nil {
		t.Fatalf("raw D2: %v", err)
	}
	shown, err := file.GetCellValue(xlsxSheet, "D2")
	if err != nil {
		t.Fatalf("shown D2: %v", err)
	}
	if shown != "300.00" {
		t.Errorf("D2 displayed = %q, want 300.00", shown)
	}
	if raw == shown {
		t.Errorf("D2 stored as text, not a native number (raw == shown == %q)", raw)
	}
	if fmtStr := cellStyle(t, file, "D2").CustomNumFmt; fmtStr == nil || *fmtStr != defaultCurrencyFmt {
		t.Errorf("D2 number format = %v, want %q", fmtStr, defaultCurrencyFmt)
	}

	// Header (A1), a subtotal number (D4), and the grand-total label (A7) are
	// bold; a detail cell (D2) is not.
	if !isBold(cellStyle(t, file, "A1")) {
		t.Error("header A1 should be bold")
	}
	if !isBold(cellStyle(t, file, "D4")) {
		t.Error("subtotal D4 should be bold")
	}
	if !isBold(cellStyle(t, file, "A7")) {
		t.Error("grand-total A7 should be bold")
	}
	if isBold(cellStyle(t, file, "D2")) {
		t.Error("detail D2 should not be bold")
	}
}

// A dimension and a label column are presented by their declared type — a
// boolean as Yes/No, a date as its calendar day, and money per the report's
// currency configuration — rather than by the engine's raw rendering. Label
// cells stay strings even when they hold money: they name a bucket, they do not
// measure one.
func TestXLSX_TypedDimensionsAndLabels(t *testing.T) {
	due := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"reimbursable"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "due"},
		Columns: []reporting.Column{
			{Name: "Due", Label: "Due", Kind: reporting.ColumnLabel, Field: "due"},
			{Name: "Reimbursable", Label: "Reimbursable", Kind: reporting.ColumnLabel, Field: "reimbursable"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(custom_1)"},
		},
	}
	rows := []reporting.Row{
		{"reimbursable": {reporting.Bool(true)}, "due": {reporting.DateVal(due)}, "custom_1": {money("1500.5")}},
		{"reimbursable": {reporting.Bool(false)}, "due": {reporting.DateVal(due)}, "custom_1": {money("2")}},
	}

	model := mustRun(t, spec, typedCatalog(t), rows)
	model.Meta.Currency = &reporting.CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: "."}

	grid := xlsxGrid(t, model, []Dimension{{Key: "reimbursable", Label: "Reimbursable", DataType: reporting.TypeBool}})
	assertGrid(t, grid, [][]string{
		{"Reimbursable", "Due", "Reimbursable", "Total"},
		{"No", "2026-06-01", "No", "$2.00"},
		{"Yes", "2026-06-01", "Yes", "$1,500.50"},
	})
}

// Grouping by a currency field renders its buckets as money, and a receipt with
// no value for it falls into the (None) bucket like any other dimension.
func TestXLSX_CurrencyGroupDimension(t *testing.T) {
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"custom_1"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"custom_1": {money("15.60")}, "category": {reporting.Str("Food")}, "amount": {money("100")}},
		{"custom_1": {money("15.6")}, "category": {reporting.Str("Food")}, "amount": {money("50")}},
		{"category": {reporting.Str("Gas")}, "amount": {money("10")}},
	}

	model := mustRun(t, spec, typedCatalog(t), rows)
	model.Meta.Currency = &reporting.CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: "."}

	grid := xlsxGrid(t, model, []Dimension{{Key: "custom_1", Label: "HST", DataType: reporting.TypeCurrency}})
	assertGrid(t, grid, [][]string{
		{"HST", "Category", "Total"},
		// 15.60 and 15.6 are one bucket, so the two receipts sum together.
		{"$15.60", "Food", "$150.00"},
		{"(None)", "Gas", "$10.00"},
	})
}
