package render

import (
	"strconv"
	"strings"

	"receipt-wrangler/api/internal/reporting"

	"github.com/xuri/excelize/v2"
)

const (
	xlsxSheet          = "Report"
	xlsxTotalLabel     = "Total"
	xlsxGrandLabel     = "Grand Total"
	defaultCurrencyFmt = "#,##0.00"
)

// XLSX renders a ReportModel as a faithful, grouped spreadsheet: the group-by
// dimensions are leading columns (each value shown once per group), the report
// columns follow, and subtotal/grand-total rows carry a Total marker in the
// column at the group's depth. Numbers are written as native, typed cells with a
// number format — not strings — so the workbook is analyzable, and the header,
// subtotal, and grand-total rows are bold.
//
// It writes the engine-computed values statically; it does not (yet) translate
// arithmetic columns into live cell formulas. Like every renderer it is a pure
// consumer of the model — it fetches and computes nothing, and reaches back into
// no upstream code.
//
// groupBy supplies the dimension order and header labels (the model carries
// dimension keys but not labels), and must match the report's grouping depth —
// see validateGroupByDepth.
func XLSX(model reporting.ReportModel, groupBy []Dimension) ([]byte, error) {
	if err := validateGroupByDepth(model, groupBy); err != nil {
		return nil, err
	}

	writer := &xlsxWriter{
		file:    excelize.NewFile(),
		model:   model,
		groupBy: groupBy,
		numDims: len(groupBy),
		styles:  map[string]int{},
		row:     1,
	}
	defer writer.file.Close()
	if err := writer.file.SetSheetName(writer.file.GetSheetName(0), xlsxSheet); err != nil {
		return nil, err
	}

	if err := writer.writeHeader(); err != nil {
		return nil, err
	}
	if err := writer.walk(model.Root, 0, nil); err != nil {
		return nil, err
	}
	if model.GrandTotals != nil {
		if err := writer.emitTotal(0, xlsxGrandLabel, model.GrandTotals); err != nil {
			return nil, err
		}
	}

	buffer, err := writer.file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// xlsxWriter carries the workbook and the walk's running state: the next row to
// write, the previous detail row's group path (for dimension blanking), and a
// cache of registered styles keyed by (number format, bold).
type xlsxWriter struct {
	file    *excelize.File
	model   reporting.ReportModel
	groupBy []Dimension
	numDims int

	styles   map[string]int
	row      int
	prevPath []reporting.Value
}

func (w *xlsxWriter) writeHeader() error {
	column := 1
	for _, dimension := range w.groupBy {
		if err := w.setString(column, w.row, dimensionHeading(dimension), true); err != nil {
			return err
		}
		column++
	}
	for _, descriptor := range w.model.Columns {
		if err := w.setString(column, w.row, columnHeading(descriptor), true); err != nil {
			return err
		}
		column++
	}
	w.row++
	return nil
}

// walk emits a node's descendants then its own subtotal, so a subtotal always
// follows the rows it sums. level is the node's depth (root is 0), which places
// the subtotal marker's column.
func (w *xlsxWriter) walk(node reporting.GroupNode, level int, path []reporting.Value) error {
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			childPath := make([]reporting.Value, len(path), len(path)+1)
			copy(childPath, path)
			childPath = append(childPath, child.Value)
			if err := w.walk(child, level+1, childPath); err != nil {
				return err
			}
		}
	} else {
		for _, detail := range node.DetailRows {
			if err := w.emitDetail(path, detail.Cells); err != nil {
				return err
			}
		}
	}

	if node.Subtotals != nil {
		if err := w.emitTotal(level, xlsxTotalLabel, node.Subtotals); err != nil {
			return err
		}
	}
	return nil
}

// emitDetail writes one detail row. A dimension value is printed only from the
// first level that differs from the previous detail row's path; the unchanged
// leading levels stay blank, so a value shows once per group.
func (w *xlsxWriter) emitDetail(path []reporting.Value, cells []reporting.Cell) error {
	firstDiff := w.firstDiff(path)
	for index := firstDiff; index < len(path); index++ {
		if err := w.setString(index+1, w.row, formatDimension(path[index], w.model.Meta.NoneLabel), false); err != nil {
			return err
		}
	}

	for offset, descriptor := range w.model.Columns {
		if err := w.writeReportCell(w.numDims+offset+1, w.row, descriptor, cells[offset], false); err != nil {
			return err
		}
	}

	w.prevPath = append(w.prevPath[:0], path...)
	w.row++
	return nil
}

// emitTotal writes a subtotal or grand-total row. The label lands in the column
// at the group's depth (the staircase); every other dimension cell is blank, and
// the measure columns carry the roll-up numbers. The whole row is bold.
func (w *xlsxWriter) emitTotal(level int, label string, cells []reporting.Cell) error {
	if level < w.numDims {
		if err := w.setString(level+1, w.row, label, true); err != nil {
			return err
		}
	}

	for offset, descriptor := range w.model.Columns {
		column := w.numDims + offset
		if column == level {
			if err := w.setString(column+1, w.row, label, true); err != nil {
				return err
			}
			continue
		}
		if err := w.writeReportCell(column+1, w.row, descriptor, cells[offset], true); err != nil {
			return err
		}
	}

	w.row++
	return nil
}

// firstDiff is the first group level whose value differs from the previous detail
// row's, or the path length when the row repeats the previous group entirely.
func (w *xlsxWriter) firstDiff(path []reporting.Value) int {
	for index := range path {
		if index >= len(w.prevPath) || !path[index].Equal(w.prevPath[index]) {
			return index
		}
	}
	return len(path)
}

// writeReportCell writes one report column's cell. A label column is joined text;
// a numeric column is a native number with the column's format; a null is left
// blank. bold styles the whole cell (subtotal/grand-total rows).
func (w *xlsxWriter) writeReportCell(column, row int, descriptor reporting.ColumnDescriptor, cell reporting.Cell, bold bool) error {
	if descriptor.Kind == reporting.ColumnLabel {
		return w.setString(column, row, joinLabel(cell, w.model.Meta.NoneLabel), bold)
	}

	value := cell.Value()
	if value.IsNull() {
		return nil
	}
	if number, ok := value.Decimal(); ok {
		return w.setNumber(column, row, number.InexactFloat64(), w.numberFormat(descriptor), bold)
	}
	// A non-number, non-null measure cell cannot arise through Run; render its
	// text defensively, matching the CSV renderer.
	return w.setString(column, row, value.String(), bold)
}

// numberFormat resolves a numeric column's Excel format: the column's own format
// override, else a currency default (the report's currency format, or a plain
// two-place format), else General for other numbers.
func (w *xlsxWriter) numberFormat(descriptor reporting.ColumnDescriptor) string {
	if descriptor.Format != "" {
		return descriptor.Format
	}
	if descriptor.DataType == reporting.TypeCurrency {
		if w.model.Meta.CurrencyFormat != "" {
			return w.model.Meta.CurrencyFormat
		}
		return defaultCurrencyFmt
	}
	return ""
}

func (w *xlsxWriter) setString(column, row int, text string, bold bool) error {
	if text == "" {
		return nil
	}
	cell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		return err
	}
	if err := w.file.SetCellValue(xlsxSheet, cell, text); err != nil {
		return err
	}
	return w.applyStyle(cell, "", bold)
}

func (w *xlsxWriter) setNumber(column, row int, value float64, numberFormat string, bold bool) error {
	cell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		return err
	}
	if err := w.file.SetCellValue(xlsxSheet, cell, value); err != nil {
		return err
	}
	return w.applyStyle(cell, numberFormat, bold)
}

func (w *xlsxWriter) applyStyle(cell, numberFormat string, bold bool) error {
	id, err := w.styleID(numberFormat, bold)
	if err != nil || id == 0 {
		return err
	}
	return w.file.SetCellStyle(xlsxSheet, cell, cell, id)
}

// styleID registers (and caches) a style for a number-format/bold combination.
// The default cell — no format, not bold — needs no style, so it returns 0.
func (w *xlsxWriter) styleID(numberFormat string, bold bool) (int, error) {
	if numberFormat == "" && !bold {
		return 0, nil
	}
	key := numberFormat + "|" + strconv.FormatBool(bold)
	if id, ok := w.styles[key]; ok {
		return id, nil
	}

	style := &excelize.Style{}
	if bold {
		style.Font = &excelize.Font{Bold: true}
	}
	if numberFormat != "" {
		format := numberFormat
		style.CustomNumFmt = &format
	}

	id, err := w.file.NewStyle(style)
	if err != nil {
		return 0, err
	}
	w.styles[key] = id
	return id, nil
}

// joinLabel renders a label cell: its values joined with commas, a null as the
// report's (None) name, and no values (a subtotal's label column) as empty.
func joinLabel(cell reporting.Cell, noneLabel string) string {
	if len(cell.Values) == 0 {
		return ""
	}
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
