// Package render turns the reporting engine's format-agnostic ReportModel into a
// concrete output format. It is a pure downstream consumer: it reads a
// ReportModel and never fetches, computes, or reaches back into the engine.
package render

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"receipt-wrangler/api/internal/reporting"
)

// Dimension is one group-by level: the field it cuts on and the header a
// renderer prints for it.
//
// The ReportModel carries a dimension key on each group node but not its label,
// and a flat table needs the levels up front and in order — an empty report
// still needs its dimension columns. So a caller supplies them, in the same
// order as the spec's GroupBy. A future orchestrator builds this from the spec
// and the field catalog.
type Dimension struct {
	Key   reporting.FieldKey
	Label string
}

const (
	rowTypeHeader = "Row Type"
	detailRow     = "Detail"
	subtotalRow   = "Subtotal"
	grandTotalRow = "Grand Total"
)

// CSV renders a ReportModel as a flat, data-only CSV table: the group-by
// dimensions become leading columns, each detail leaf is one row, and a leading
// "Row Type" column marks whether a row is a Detail, a Subtotal, or the Grand
// Total. Subtotal and grand-total rows are emitted only when the model carries
// them (the spec's Subtotals/GrandTotals toggles decide that upstream); the
// renderer never drops what it is given.
//
// The "Row Type" column is what keeps the file safe to aggregate: a consumer
// filters Row Type=Detail before summing an additive column — otherwise the
// subtotal and grand-total rows double-count it — and reads the roll-up rows
// for a non-additive column such as an average, which the engine recomputed at
// each level rather than summing.
//
// It renders no document chrome (title, intro, branding); that is the richer
// formats' job. This is the minimal, machine-friendly export.
func CSV(model reporting.ReportModel, groupBy []Dimension) ([]byte, error) {
	if err := validateGroupByDepth(model, groupBy); err != nil {
		return nil, err
	}

	records := [][]string{header(model, groupBy)}
	records = appendGroup(records, model, groupBy, model.Root, nil)

	if model.GrandTotals != nil {
		records = append(records, buildRecord(model, groupBy, grandTotalRow, nil, model.GrandTotals))
	}

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.UseCRLF = true // RFC 4180 line endings, which Excel prefers.
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// header is: Row Type, then one column per group-by dimension, then one per
// report column.
func header(model reporting.ReportModel, groupBy []Dimension) []string {
	head := make([]string, 0, 1+len(groupBy)+len(model.Columns))
	head = append(head, rowTypeHeader)
	for _, dimension := range groupBy {
		head = append(head, dimensionHeading(dimension))
	}
	for _, column := range model.Columns {
		head = append(head, columnHeading(column))
	}
	return head
}

// dimensionHeading falls back to the raw field key when a caller supplies no
// label, so a column is never blank-headed.
func dimensionHeading(dimension Dimension) string {
	if len(dimension.Label) > 0 {
		return SanitizeCSVField(dimension.Label)
	}
	return SanitizeCSVField(string(dimension.Key))
}

// validateGroupByDepth rejects a groupBy that does not match how deep the report
// actually groups. groupBy is a contract between the caller and the model's
// shape: too few dimensions would drop deeper ancestor values, so distinct
// buckets would print identical leading columns; too many would pad every row
// with blanks. An empty grouped report has a childless root and no discoverable
// depth, so it is left to render — its dimension columns come from groupBy. Both
// renderers share this check.
func validateGroupByDepth(model reporting.ReportModel, groupBy []Dimension) error {
	if depth := groupDepth(model.Root); depth > 0 && depth != len(groupBy) {
		return fmt.Errorf("render: report groups %d level(s) deep, but %d dimension(s) were supplied", depth, len(groupBy))
	}
	return nil
}

// groupDepth is how many grouping levels the tree nests. Every branch of a tree
// for a fixed spec has the same depth, so walking the first child suffices. It is
// 0 for an ungrouped report and for an empty one (a childless root), which is why
// the caller's depth check only fires when the tree actually groups.
func groupDepth(node reporting.GroupNode) int {
	depth := 0
	for len(node.Children) > 0 {
		depth++
		node = node.Children[0]
	}
	return depth
}

// appendGroup walks the tree in the order a reader expects: a node's children
// (each of which emits its own detail and subtotal rows) come before the node's
// own subtotal, so a subtotal always follows the rows it sums.
//
// path carries the bucket value of each ancestor level, which fills the leading
// dimension columns.
func appendGroup(records [][]string, model reporting.ReportModel, groupBy []Dimension, node reporting.GroupNode, path []reporting.Value) [][]string {
	// A node holds either Children or DetailRows, never both; an empty root has
	// neither, and falls through the detail loop harmlessly.
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			childPath := make([]reporting.Value, len(path), len(path)+1)
			copy(childPath, path)
			childPath = append(childPath, child.Value)
			records = appendGroup(records, model, groupBy, child, childPath)
		}
	} else {
		for _, detail := range node.DetailRows {
			records = append(records, buildRecord(model, groupBy, detailRow, path, detail.Cells))
		}
	}

	if node.Subtotals != nil {
		records = append(records, buildRecord(model, groupBy, subtotalRow, path, node.Subtotals))
	}
	return records
}

// buildRecord assembles one CSV record: the row-type marker, the leading
// dimension columns (the ancestor path, padded with blanks for the levels below
// this row), then one formatted cell per report column.
func buildRecord(model reporting.ReportModel, groupBy []Dimension, rowType string, path []reporting.Value, cells []reporting.Cell) []string {
	record := make([]string, 0, 1+len(groupBy)+len(model.Columns))
	record = append(record, rowType)

	for index := range groupBy {
		if index < len(path) {
			record = append(record, formatDimension(path[index], model.Meta.NoneLabel))
		} else {
			record = append(record, "")
		}
	}

	for index, column := range model.Columns {
		record = append(record, formatCell(column, cells[index], model.Meta.NoneLabel, model.Meta.Currency))
	}
	return record
}

func columnHeading(column reporting.ColumnDescriptor) string {
	if len(column.Label) > 0 {
		return SanitizeCSVField(column.Label)
	}
	return SanitizeCSVField(column.Name)
}

// formatDimension renders a group bucket's value for a leading column. The
// (None) bucket — a null value — gets the report's name for it.
func formatDimension(value reporting.Value, noneLabel string) string {
	if value.IsNull() {
		return SanitizeCSVField(noneLabel)
	}
	return SanitizeCSVField(value.String())
}

// formatCell renders one report cell the way this format presents it: money per
// the report's currency configuration (a bare two places when none is set),
// other numbers at full precision, an undefined value as an empty cell, and a
// multi-value label (a record carrying several categories/tags) joined with
// commas.
//
// A cell with no values is a label column on a subtotal row — no single bucket
// named it — and stays blank. A label cell holding one null value is the (None)
// bucket, a bucket like any other, and gets the report's name for it.
func formatCell(column reporting.ColumnDescriptor, cell reporting.Cell, noneLabel string, currency *reporting.CurrencyFormat) string {
	if len(cell.Values) == 0 {
		return ""
	}

	if column.Kind == reporting.ColumnLabel {
		parts := make([]string, 0, len(cell.Values))
		for _, value := range cell.Values {
			if value.IsNull() {
				parts = append(parts, SanitizeCSVField(noneLabel))
				continue
			}
			parts = append(parts, SanitizeCSVField(value.String()))
		}
		return strings.Join(parts, ", ")
	}

	value := cell.Value()
	if value.IsNull() {
		return ""
	}

	number, isNumber := value.Decimal()
	if !isNumber {
		return SanitizeCSVField(value.String())
	}
	if column.DataType == reporting.TypeCurrency {
		if currency != nil {
			return formatCurrency(number, *currency)
		}
		return number.StringFixed(2)
	}
	return number.String()
}

// SanitizeCSVField neutralizes spreadsheet formula/DDE injection: a text cell whose
// first byte is a character Excel/Sheets may treat as a formula lead (= + - @) or a
// control character (tab, CR) is prefixed with a single quote so the app renders it as
// literal text. A value that parses as a number (e.g. a negative amount) is left intact
// so the CSV stays machine-readable.
//
// It is exported so the receipt CSV export (services.ReceiptCsvService) can apply the
// same guard to its user-controlled text columns without duplicating the logic.
func SanitizeCSVField(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		if _, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
			return s
		}
		return "'" + s
	}
	return s
}
