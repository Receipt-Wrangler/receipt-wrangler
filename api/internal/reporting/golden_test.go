package reporting

import (
	"fmt"
	"strings"
	"testing"
)

// renderTextTable lays a report out as a fixed-width table, the way a human
// reads one. It exists so that a single assertion covers the shape of a whole
// report at once: the grouping, the ordering, the values, where subtotals sit,
// and which cells are blank.
//
// It is a test helper, not a renderer. It assumes the report's first column is
// the one naming each detail row, which is how the design's report is built.
func renderTextTable(model ReportModel, labels map[FieldKey]string) string {
	type line struct {
		cells   []string
		heading bool
	}

	var lines []line

	header := make([]string, 0, len(model.Columns))
	for _, column := range model.Columns {
		header = append(header, column.Label)
	}
	lines = append(lines, line{cells: header})

	format := func(cells []Cell, leading string, indent int) []string {
		rendered := make([]string, len(cells))
		for index, cell := range cells {
			rendered[index] = formatCell(model.Columns[index], cell, model.Meta.NoneLabel)
		}
		if leading != "" {
			rendered[0] = leading
		}
		rendered[0] = strings.Repeat("  ", indent) + rendered[0]
		return rendered
	}

	var walk func(node GroupNode, depth int)
	walk = func(node GroupNode, depth int) {
		if depth > 0 {
			name := labels[node.Dimension]
			value := node.Value.String()
			if node.IsNone {
				value = model.Meta.NoneLabel
			}
			lines = append(lines, line{
				cells:   []string{strings.Repeat("  ", depth-1) + name + ": " + value},
				heading: true,
			})
		}

		for _, row := range node.DetailRows {
			lines = append(lines, line{cells: format(row.Cells, "", depth)})
		}
		for _, child := range node.Children {
			walk(child, depth+1)
		}
		if node.Subtotals != nil {
			lines = append(lines, line{cells: format(node.Subtotals, "TOTALS", depth)})
		}
	}
	walk(model.Root, 0)

	if model.GrandTotals != nil {
		lines = append(lines, line{cells: format(model.GrandTotals, "GRAND TOTALS", 0)})
	}

	widths := make([]int, len(model.Columns))
	for _, current := range lines {
		if current.heading {
			continue
		}
		for index, cell := range current.cells {
			if len(cell) > widths[index] {
				widths[index] = len(cell)
			}
		}
	}

	var out strings.Builder
	for _, current := range lines {
		if current.heading {
			out.WriteString(current.cells[0] + "\n")
			continue
		}
		parts := make([]string, len(current.cells))
		for index, cell := range current.cells {
			if index == 0 {
				parts[index] = fmt.Sprintf("%-*s", widths[index], cell)
			} else {
				parts[index] = fmt.Sprintf("%*s", widths[index], cell)
			}
		}
		out.WriteString(strings.TrimRight(strings.Join(parts, "  "), " ") + "\n")
	}

	return out.String()
}

// formatCell renders one cell the way a report would present it: money to two
// places, counts as whole numbers, labels as text, and an undefined value as an
// empty cell.
//
// A label cell holding no values is blank — that is a subtotal row, which no
// single bucket named. A label cell holding one null value is the (None)
// bucket, which is a bucket like any other and gets the report's name for it.
// The two are different things and the model keeps them apart.
func formatCell(column ColumnDescriptor, cell Cell, noneLabel string) string {
	if len(cell.Values) == 0 {
		return ""
	}

	if column.Kind == ColumnLabel {
		parts := make([]string, 0, len(cell.Values))
		for _, value := range cell.Values {
			if value.IsNull() {
				parts = append(parts, noneLabel)
				continue
			}
			parts = append(parts, value.String())
		}
		return strings.Join(parts, ", ")
	}

	value := cell.Value()
	if value.IsNull() {
		return ""
	}

	number, isNumber := value.Decimal()
	if !isNumber {
		return value.String()
	}
	if column.DataType == TypeCurrency {
		return number.StringFixed(2)
	}
	return number.String()
}

// The golden report. Every number, every row and every position below is read
// off the design document's worked example.
//
// Money is shown to two places, which is presentation rather than data: the
// engine's Avg/Receipt for Alex is 35.933333, and it is the renderer that shows
// 35.93. The grouping labels come from the catalog, not from the engine, which
// only names the dimension each level was cut by.
func TestGolden_WorkedExampleRendersAsTheDesignDocument(t *testing.T) {
	model := mustRun(t, workedExampleSpec(), workedExampleRows())

	labels := map[FieldKey]string{"paid_by": "Paid By", "tag": "Tag"}

	want := strings.TrimPrefix(`
Category      Count  Subtotal    Hst   Total  Avg/Receipt
Paid By: Dana
  Tag: Alex
    Clothing      4    120.00  15.60  135.60        33.90
    Medical       2     80.00   0.00   80.00        40.00
    TOTALS        6    200.00  15.60  215.60        35.93
  Tag: Sam
    Clothing      3     90.00  11.70  101.70        33.90
    Mileage       1     30.00   0.00   30.00        30.00
    TOTALS        4    120.00  11.70  131.70        32.93
  TOTALS         10    320.00  27.30  347.30        34.73
GRAND TOTALS     10    320.00  27.30  347.30        34.73
`, "\n")

	if got := renderTextTable(model, labels); got != want {
		t.Errorf("report does not match the design document.\n\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// The same report with the (None) bucket in it, so the empty-value case has a
// golden shape too. A receipt with no tag and one with no category are both
// kept, and both sort last within their level.
func TestGolden_NoneBucketsRenderLast(t *testing.T) {
	spec := workedExampleSpec()
	spec.NoneLabel = "(None)"

	rows := append(workedExampleRows(),
		// Tagged, but uncategorized.
		Row{"paid_by": {Str("Dana")}, "tag": {Str("Alex")}, "amount": {Num(dec("7.00"))}},
		// Untagged and uncategorized.
		Row{"paid_by": {Str("Dana")}, "amount": {Num(dec("3.00"))}},
	)

	model := mustRun(t, spec, rows)
	labels := map[FieldKey]string{"paid_by": "Paid By", "tag": "Tag"}

	want := strings.TrimPrefix(`
Category      Count  Subtotal    Hst   Total  Avg/Receipt
Paid By: Dana
  Tag: Alex
    Clothing      4    120.00  15.60  135.60        33.90
    Medical       2     80.00   0.00   80.00        40.00
    (None)        1      7.00   0.00    7.00         7.00
    TOTALS        7    207.00  15.60  222.60        31.80
  Tag: Sam
    Clothing      3     90.00  11.70  101.70        33.90
    Mileage       1     30.00   0.00   30.00        30.00
    TOTALS        4    120.00  11.70  131.70        32.93
  Tag: (None)
    (None)        1      3.00   0.00    3.00         3.00
    TOTALS        1      3.00   0.00    3.00         3.00
  TOTALS         12    330.00  27.30  357.30        29.78
GRAND TOTALS     12    330.00  27.30  357.30        29.78
`, "\n")

	if got := renderTextTable(model, labels); got != want {
		t.Errorf("report does not match.\n\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
