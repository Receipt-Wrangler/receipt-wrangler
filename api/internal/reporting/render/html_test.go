package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/reporting"

	"golang.org/x/net/html"
)

// --- helpers --------------------------------------------------------------

func parseHTML(t *testing.T, out []byte) *html.Node {
	t.Helper()
	doc, err := html.Parse(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("parse html: %v", err)
	}
	return doc
}

// findFirst returns the first element of the given tag in document order.
func findFirst(node *html.Node, tag string) *html.Node {
	if node.Type == html.ElementNode && node.Data == tag {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findFirst(child, tag); found != nil {
			return found
		}
	}
	return nil
}

// elementChildren returns the direct child elements of the given tag.
func elementChildren(node *html.Node, tag string) []*html.Node {
	var out []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == tag {
			out = append(out, child)
		}
	}
	return out
}

func textOf(node *html.Node) string {
	var builder strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			builder.WriteString(n.Data)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return builder.String()
}

func attrOf(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func cellTexts(row *html.Node, cellTag string) []string {
	cells := elementChildren(row, cellTag)
	texts := make([]string, len(cells))
	for index, cell := range cells {
		texts[index] = textOf(cell)
	}
	return texts
}

// htmlGrid renders the model and returns the table as a grid whose first row is
// the header, so it can be asserted against the same golden grids as the XLSX
// renderer (via assertGrid, which trims trailing blanks).
func htmlGrid(t *testing.T, model reporting.ReportModel, groupBy []Dimension) [][]string {
	t.Helper()
	out, err := HTML(model, groupBy)
	if err != nil {
		t.Fatalf("HTML: %v", err)
	}
	table := findFirst(parseHTML(t, out), "table")
	if table == nil {
		t.Fatal("no <table> in rendered document")
	}

	grid := [][]string{cellTexts(findFirst(findFirst(table, "thead"), "tr"), "th")}
	for _, row := range elementChildren(findFirst(table, "tbody"), "tr") {
		grid = append(grid, cellTexts(row, "td"))
	}
	return grid
}

// bodyRows renders the model and returns the <tbody> <tr> nodes, for asserting
// per-row classes and per-cell attributes the grid text drops.
func bodyRows(t *testing.T, model reporting.ReportModel, groupBy []Dimension) []*html.Node {
	t.Helper()
	out, err := HTML(model, groupBy)
	if err != nil {
		t.Fatalf("HTML: %v", err)
	}
	return elementChildren(findFirst(findFirst(parseHTML(t, out), "table"), "tbody"), "tr")
}

// --- faithful parity: the same grids the XLSX renderer produces -----------

// The dimension shows once per group, subtotal/grand-total rows carry a Total
// marker, additive columns roll up, and the non-linear Avg is recomputed per
// level — the same faithful grid the XLSX renderer emits.
func TestHTML_SumsAggregatesAndNonLinearRollUp(t *testing.T) {
	model := mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows())

	assertGrid(t, htmlGrid(t, model, paidByDimension()), [][]string{
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
func TestHTML_NestedGroupsStaircaseAndBlanking(t *testing.T) {
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

	grid := htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows),
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

// Each aggregate function renders at detail and roll-up level.
func TestHTML_AllAggregateFunctions(t *testing.T) {
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

	grid := htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension())
	assertGrid(t, grid, [][]string{
		{"Paid By", "Category", "Sum", "Count", "Avg", "Min", "Max"},
		{"Dana", "Food", "40.00", "2", "20.00", "10.00", "30.00"},
		{"", "Gas", "110.00", "1", "110.00", "110.00", "110.00"},
		{"", "Total", "150.00", "3", "50.00", "10.00", "110.00"},
		{"Grand Total", "", "150.00", "3", "50.00", "10.00", "110.00"},
	})
}

// Records mode emits one row per record; a multi-value label is joined with
// commas into a single cell.
func TestHTML_RecordsModeMultiValueLabel(t *testing.T) {
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

	grid := htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension())
	assertGrid(t, grid, [][]string{
		{"Paid By", "Category", "Total"},
		{"Dana", "Food, Gas", "100.00"},
		{"", "Food", "50.00"},
		{"", "Total", "150.00"},
		{"Grand Total", "", "150.00"},
	})
}

// With no grouping there are no dimension columns and the grand-total marker
// lands in the first report column.
func TestHTML_NoGrouping(t *testing.T) {
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

	assertGrid(t, htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), nil), [][]string{
		{"Category", "Total"},
		{"Food", "125.00"},
		{"Gas", "50.00"},
		{"Grand Total", "175.00"},
	})
}

// Subtotals on, grand totals off: subtotal rows appear, no grand-total row.
func TestHTML_SubtotalsWithoutGrandTotal(t *testing.T) {
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

	assertGrid(t, htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension()), [][]string{
		{"Paid By", "Category", "Total"},
		{"Dana", "Food", "100.00"},
		{"", "Gas", "50.00"},
		{"", "Total", "150.00"},
		{"Sam", "Food", "40.00"},
		{"", "Total", "40.00"},
	})
}

// A row with no value for the grouping dimension is the (None) bucket; it prints
// the report's name for it.
func TestHTML_NoneBucket(t *testing.T) {
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

	assertGrid(t, htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), paidByDimension()), [][]string{
		{"Paid By", "Category", "Total"},
		{"Dana", "Food", "50.00"},
		{"", "Total", "50.00"},
		{"(None)", "Food", "100.00"},
		{"", "Total", "100.00"},
		{"Grand Total", "", "150.00"},
	})
}

// A multi-value grouping dimension fans a row into every bucket in full.
func TestHTML_MultiValueGroupFanOut(t *testing.T) {
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

	assertGrid(t, htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), []Dimension{{Key: "tag", Label: "Tag"}}), [][]string{
		{"Tag", "Category", "Total"},
		{"Alex", "Food", "100.00"},
		{"Bob", "Food", "100.00"},
		{"Grand Total", "", "200.00"},
	})
}

// An empty report is the header plus a grand-total row (zero count, zero sum).
func TestHTML_EmptyReport(t *testing.T) {
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

	assertGrid(t, htmlGrid(t, mustRun(t, spec, paidByCatalog(t), nil), paidByDimension()), [][]string{
		{"Paid By", "Category", "Count", "Total"},
		{"Grand Total", "", "0", "0.00"},
	})
}

// Unicode passes through into a cell unchanged.
func TestHTML_Unicode(t *testing.T) {
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

	assertGrid(t, htmlGrid(t, mustRun(t, spec, paidByCatalog(t), rows), nil), [][]string{
		{"Category", "Total"},
		{"日本語 🎉", "10.00"},
	})
}

// Supplying fewer dimensions than the report groups by is rejected, reusing the
// shared depth guard, and yields no bytes.
func TestHTML_GroupByDepthMismatchErrors(t *testing.T) {
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

	out, err := HTML(mustRun(t, spec, paidByCatalog(t), rows), paidByDimension())
	if err == nil {
		t.Fatalf("expected an error for a group-depth mismatch, got none")
	}
	if out != nil {
		t.Errorf("expected nil bytes on error, got %d", len(out))
	}
}

// --- HTML-specific structure: row classes and numeric alignment -----------

// Detail rows are unclassed; subtotal and grand-total rows carry their CSS class
// so the stylesheet can emphasize them.
func TestHTML_RowClassesMarkRollUps(t *testing.T) {
	rows := bodyRows(t, mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows()), paidByDimension())

	want := []string{"", "", "subtotal", "", "subtotal", "grand-total"}
	if len(rows) != len(want) {
		t.Fatalf("row count = %d, want %d", len(rows), len(want))
	}
	for index, class := range want {
		if got := attrOf(rows[index], "class"); got != class {
			t.Errorf("row %d class = %q, want %q", index, got, class)
		}
	}
}

// Report columns are right-aligned via the "num" class, except label columns;
// dimension columns are never numeric.
func TestHTML_NumericCellsAreRightAligned(t *testing.T) {
	rows := bodyRows(t, mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows()), paidByDimension())

	// First detail row: Paid By, Category (both label), then Count, Total, Avg.
	cells := elementChildren(rows[0], "td")
	wantNum := []bool{false, false, true, true, true}
	if len(cells) != len(wantNum) {
		t.Fatalf("cell count = %d, want %d", len(cells), len(wantNum))
	}
	for index, num := range wantNum {
		isNum := attrOf(cells[index], "class") == "num"
		if isNum != num {
			t.Errorf("cell %d numeric = %v, want %v", index, isNum, num)
		}
	}
}

// --- document chrome ------------------------------------------------------

// A title, resolved parameters, and a generated-at timestamp render as the
// document's heading, a sorted preamble, and a footer.
func TestHTML_DocumentChromeRendersFromMeta(t *testing.T) {
	spec := oneLevelSpec()
	spec.Title = "Quarterly Expenses"
	meta := reporting.MetaInput{
		GeneratedAt: time.Date(2026, 7, 13, 9, 30, 0, 0, time.UTC),
		Params:      map[string]string{"Period": "2026-Q2", "Group": "Household"},
	}
	model, err := reporting.Run(spec, paidByCatalog(t), oneLevelRows(), meta)
	if err != nil {
		t.Fatalf("run report: %v", err)
	}

	out, err := HTML(model, paidByDimension())
	if err != nil {
		t.Fatalf("HTML: %v", err)
	}
	rendered := string(out)

	if !strings.Contains(rendered, "<h1>Quarterly Expenses</h1>") {
		t.Errorf("title heading missing:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Generated 2026-07-13 09:30:00 UTC") {
		t.Errorf("generated-at footer missing:\n%s", rendered)
	}
	for _, want := range []string{
		`<span class="key">Group</span>: Household`,
		`<span class="key">Period</span>: 2026-Q2`,
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("preamble missing %q", want)
		}
	}
	// Params render alphabetically: Group before Period.
	if strings.Index(rendered, ">Group<") > strings.Index(rendered, ">Period<") {
		t.Error("params not sorted: Group should precede Period")
	}
}

// With no title, no params, and a zero timestamp the document omits its chrome
// entirely rather than rendering empty elements.
func TestHTML_DocumentChromeOmittedWhenEmpty(t *testing.T) {
	out, err := HTML(mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows()), paidByDimension())
	if err != nil {
		t.Fatalf("HTML: %v", err)
	}
	rendered := string(out)

	for _, absent := range []string{"<h1>", `class="preamble"`, "<footer>"} {
		if strings.Contains(rendered, absent) {
			t.Errorf("expected %q to be omitted:\n%s", absent, rendered)
		}
	}
}

// Cell text is HTML-escaped, so a value containing markup cannot break out of
// its cell.
func TestHTML_EscapesCellText(t *testing.T) {
	spec := reporting.ReportSpec{
		Detail: reporting.DetailSpec{Mode: reporting.DetailAggregate, By: "category"},
		Columns: []reporting.Column{
			{Name: "Category", Label: "Category", Kind: reporting.ColumnLabel, Field: "category"},
			{Name: "Total", Label: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{"category": {reporting.Str("A & <b>")}, "amount": {money("10")}},
	}

	out, err := HTML(mustRun(t, spec, paidByCatalog(t), rows), nil)
	if err != nil {
		t.Fatalf("HTML: %v", err)
	}
	rendered := string(out)

	if !strings.Contains(rendered, "A &amp; &lt;b&gt;") {
		t.Errorf("cell text not escaped:\n%s", rendered)
	}
	if strings.Contains(rendered, "<b>") {
		t.Errorf("unescaped markup leaked into output:\n%s", rendered)
	}
}
