package render

import "receipt-wrangler/api/internal/reporting"

// The faithful renderers — XLSX and the HTML/PDF document — lay a report out the
// way it looks on screen: the group-by dimensions are leading columns with each
// value shown once per group, and a subtotal/grand-total row carries its marker
// in the column at the group's depth (the "staircase"). This file holds that
// shared traversal. A format supplies a faithfulSink that decides how a cell
// looks; the walk decides where every cell goes. The CSV renderer is
// deliberately not built on this — it is flat, not faithful, and keeps its own
// walk.

// totalMarker and grandTotalMarker are the labels a roll-up row shows in its
// staircase column.
const (
	totalMarker      = "Total"
	grandTotalMarker = "Grand Total"
)

// rowKind is the kind of data row the walk emits. Emphasis (bold, shading) is
// uniform per kind — detail rows are plain, subtotal and grand-total rows are
// emphasized — so a sink derives it from the kind rather than from each cell.
type rowKind int

const (
	detailKind rowKind = iota
	subtotalKind
	grandTotalKind
)

// cellKind tags one positioned cell of a faithful row.
type cellKind int

const (
	// blankCell is an empty grid position: a blanked repeated dimension, or a
	// dimension column on a roll-up row that names no bucket.
	blankCell cellKind = iota
	// textCell is a dimension value or a Total/Grand Total marker.
	textCell
	// reportCell is a report column's cell; it carries its descriptor and value
	// so each sink formats it its own way — a native number in XLSX, a formatted
	// string in HTML.
	reportCell
)

// faithfulCell is one cell of a faithful row, already placed at its grid column.
type faithfulCell struct {
	kind       cellKind
	text       string
	descriptor reporting.ColumnDescriptor
	cell       reporting.Cell
}

// faithfulSink consumes the shared walk. The walk hands it fully positioned rows
// (blanks included), so a sink only decides how a cell looks, never where it
// goes.
type faithfulSink interface {
	writeHeader(headings []string) error
	writeRow(kind rowKind, cells []faithfulCell) error
}

// faithfulWalk drives sink over the report the way a reader expects it: a header
// row, then each group's detail rows with its subtotal beneath the rows it sums,
// then the grand total last. It runs the group-by depth guard first, so a sink
// never has to.
func faithfulWalk(model reporting.ReportModel, groupBy []Dimension, sink faithfulSink) error {
	if err := validateGroupByDepth(model, groupBy); err != nil {
		return err
	}

	headings := make([]string, 0, len(groupBy)+len(model.Columns))
	for _, dimension := range groupBy {
		headings = append(headings, dimensionHeading(dimension))
	}
	for _, descriptor := range model.Columns {
		headings = append(headings, columnHeading(descriptor))
	}
	if err := sink.writeHeader(headings); err != nil {
		return err
	}

	walker := &faithfulWalker{
		model:     model,
		numDims:   len(groupBy),
		totalCols: len(groupBy) + len(model.Columns),
		sink:      sink,
	}
	if err := walker.walk(model.Root, 0, nil); err != nil {
		return err
	}
	if model.GrandTotals != nil {
		return walker.emitTotal(0, grandTotalKind, model.GrandTotals)
	}
	return nil
}

// faithfulWalker carries the walk's running state: the sink, the grid width, and
// the previous detail row's group path for dimension blanking.
type faithfulWalker struct {
	model     reporting.ReportModel
	numDims   int
	totalCols int
	sink      faithfulSink

	prevPath []reporting.Value
}

// walk emits a node's descendants then its own subtotal, so a subtotal always
// follows the rows it sums. level is the node's depth (root is 0), which places
// its subtotal marker's column.
func (w *faithfulWalker) walk(node reporting.GroupNode, level int, path []reporting.Value) error {
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
		if err := w.emitTotal(level, subtotalKind, node.Subtotals); err != nil {
			return err
		}
	}
	return nil
}

// emitDetail lays out one detail row. A dimension value is placed only from the
// first level that differs from the previous detail row's path; the unchanged
// leading levels stay blank, so a value shows once per group.
func (w *faithfulWalker) emitDetail(path []reporting.Value, cells []reporting.Cell) error {
	row := make([]faithfulCell, w.totalCols)
	for index := w.firstDiff(path); index < len(path); index++ {
		row[index] = faithfulCell{kind: textCell, text: formatDimension(path[index], w.model.Meta.NoneLabel)}
	}
	for offset, descriptor := range w.model.Columns {
		row[w.numDims+offset] = faithfulCell{kind: reportCell, descriptor: descriptor, cell: cells[offset]}
	}

	w.prevPath = append(w.prevPath[:0], path...)
	return w.sink.writeRow(detailKind, row)
}

// emitTotal lays out a subtotal or grand-total row: the marker lands in the
// column at the group's depth (the staircase) and the report columns carry the
// roll-up values. kind picks the marker text and tells the sink how to
// emphasize the row.
func (w *faithfulWalker) emitTotal(level int, kind rowKind, cells []reporting.Cell) error {
	label := totalMarker
	if kind == grandTotalKind {
		label = grandTotalMarker
	}

	row := make([]faithfulCell, w.totalCols)
	if level < w.numDims {
		row[level] = faithfulCell{kind: textCell, text: label}
	}
	for offset, descriptor := range w.model.Columns {
		column := w.numDims + offset
		if column == level {
			row[column] = faithfulCell{kind: textCell, text: label}
			continue
		}
		row[column] = faithfulCell{kind: reportCell, descriptor: descriptor, cell: cells[offset]}
	}
	return w.sink.writeRow(kind, row)
}

// firstDiff is the first group level whose value differs from the previous
// detail row's, or the path length when the row repeats the previous group
// entirely.
func (w *faithfulWalker) firstDiff(path []reporting.Value) int {
	for index := range path {
		if index >= len(w.prevPath) || !path[index].Equal(w.prevPath[index]) {
			return index
		}
	}
	return len(path)
}
