package render

import (
	"strconv"
	"strings"

	"receipt-wrangler/api/internal/reporting"

	"github.com/xuri/excelize/v2"
)

const (
	xlsxSheet          = "Report"
	defaultCurrencyFmt = "#,##0.00"
)

// XLSX renders a ReportModel as a faithful, grouped spreadsheet: the group-by
// dimensions are leading columns (each value shown once per group), the report
// columns follow, and subtotal/grand-total rows carry a marker in the column at
// the group's depth. Numbers are written as native, typed cells with a number
// format — not strings — so the workbook is analyzable, and the header,
// subtotal, and grand-total rows are bold.
//
// It writes the engine-computed values statically; it does not (yet) translate
// arithmetic columns into live cell formulas. Like every renderer it is a pure
// consumer of the model. It shares the faithful walk (faithfulWalk) with the
// HTML/PDF renderer; only the per-cell writing below is spreadsheet-specific.
//
// groupBy supplies the dimension order and header labels (the model carries
// dimension keys but not labels), and must match the report's grouping depth —
// see validateGroupByDepth.
func XLSX(model reporting.ReportModel, groupBy []Dimension) ([]byte, error) {
	file := excelize.NewFile()
	defer file.Close()
	if err := file.SetSheetName(file.GetSheetName(0), xlsxSheet); err != nil {
		return nil, err
	}

	sink := &xlsxSink{
		file:   file,
		model:  model,
		styles: map[string]int{},
		row:    1,
	}
	if err := faithfulWalk(model, groupBy, sink); err != nil {
		return nil, err
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// xlsxSink is the faithfulSink that writes the walk into an excelize worksheet.
// It tracks the next row to write and caches registered styles keyed by (number
// format, bold).
type xlsxSink struct {
	file  *excelize.File
	model reporting.ReportModel

	styles map[string]int
	row    int
}

func (s *xlsxSink) writeHeader(headings []string) error {
	for index, heading := range headings {
		if err := s.setString(index+1, s.row, heading, true); err != nil {
			return err
		}
	}
	s.row++
	return nil
}

func (s *xlsxSink) writeRow(kind rowKind, cells []faithfulCell) error {
	bold := kind != detailKind
	for index, cell := range cells {
		column := index + 1
		switch cell.kind {
		case blankCell:
			// An empty grid position writes nothing.
		case textCell:
			if err := s.setString(column, s.row, cell.text, bold); err != nil {
				return err
			}
		case reportCell:
			if err := s.writeReportCell(column, s.row, cell.descriptor, cell.cell, bold); err != nil {
				return err
			}
		}
	}
	s.row++
	return nil
}

// writeReportCell writes one report column's cell. A label column is joined
// text; a numeric column is a native number with the column's format; a null is
// left blank. bold styles the whole cell (subtotal/grand-total rows).
func (s *xlsxSink) writeReportCell(column, row int, descriptor reporting.ColumnDescriptor, cell reporting.Cell, bold bool) error {
	if descriptor.Kind == reporting.ColumnLabel {
		return s.setString(column, row, joinLabel(cell, descriptor, s.model.Meta.Currency, s.model.Meta.NoneLabel), bold)
	}

	value := cell.Value()
	if value.IsNull() {
		return nil
	}
	if number, ok := value.Decimal(); ok {
		return s.setNumber(column, row, number.InexactFloat64(), s.numberFormat(descriptor), bold)
	}
	// A non-number, non-null measure cell cannot arise through Run; render its
	// text defensively, matching the CSV renderer.
	return s.setString(column, row, value.String(), bold)
}

// numberFormat resolves a numeric column's Excel format: the column's own format
// override, else a currency default (the report's currency format, or a plain
// two-place format), else General for other numbers.
func (s *xlsxSink) numberFormat(descriptor reporting.ColumnDescriptor) string {
	if descriptor.Format != "" {
		return descriptor.Format
	}
	if descriptor.DataType == reporting.TypeCurrency {
		if s.model.Meta.Currency != nil {
			return excelCurrencyFormat(*s.model.Meta.Currency)
		}
		return defaultCurrencyFmt
	}
	return ""
}

func (s *xlsxSink) setString(column, row int, text string, bold bool) error {
	if text == "" {
		return nil
	}
	cell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		return err
	}
	if err := s.file.SetCellValue(xlsxSheet, cell, text); err != nil {
		return err
	}
	return s.applyStyle(cell, "", bold)
}

func (s *xlsxSink) setNumber(column, row int, value float64, numberFormat string, bold bool) error {
	cell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		return err
	}
	if err := s.file.SetCellValue(xlsxSheet, cell, value); err != nil {
		return err
	}
	return s.applyStyle(cell, numberFormat, bold)
}

func (s *xlsxSink) applyStyle(cell, numberFormat string, bold bool) error {
	id, err := s.styleID(numberFormat, bold)
	if err != nil || id == 0 {
		return err
	}
	return s.file.SetCellStyle(xlsxSheet, cell, cell, id)
}

// styleID registers (and caches) a style for a number-format/bold combination.
// The default cell — no format, not bold — needs no style, so it returns 0.
func (s *xlsxSink) styleID(numberFormat string, bold bool) (int, error) {
	if numberFormat == "" && !bold {
		return 0, nil
	}
	key := numberFormat + "|" + strconv.FormatBool(bold)
	if id, ok := s.styles[key]; ok {
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

	id, err := s.file.NewStyle(style)
	if err != nil {
		return 0, err
	}
	s.styles[key] = id
	return id, nil
}

// joinLabel renders a label cell: its values formatted per the column's declared
// type and joined with commas, a null as the report's (None) name, and no values
// (a subtotal's label column) as empty.
//
// A label column is always a string cell, even when it displays money — its
// values name a bucket rather than measure one, and a comma-joined multi-value
// label has no numeric reading at all. Aggregate columns are what carry native,
// analyzable numbers into the workbook.
func joinLabel(cell reporting.Cell, descriptor reporting.ColumnDescriptor, currency *reporting.CurrencyFormat, noneLabel string) string {
	if len(cell.Values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(cell.Values))
	for _, value := range cell.Values {
		parts = append(parts, formatLabelValue(value, descriptor.DataType, currency, noneLabel))
	}
	return strings.Join(parts, ", ")
}
