package reporting

import (
	"time"

	"github.com/expr-lang/expr/ast"
)

// ReportModel is the engine's output: a format-agnostic tree of data.
//
// It carries no layout, coordinates, or styling. There is enough structure for
// a renderer to lay out group headers, detail rows, subtotal rows and grand
// totals, and nothing about how any of them should look.
type ReportModel struct {
	Meta    Meta
	Columns []ColumnDescriptor

	// Root is a synthetic node holding the top level of the report. Its
	// Dimension is empty and its Subtotals are nil; the report-wide totals are
	// GrandTotals.
	Root GroupNode

	// GrandTotals is non-nil only when the spec asked for them.
	GrandTotals []Cell
}

// Meta is report-scope context a renderer draws around the table.
type Meta struct {
	Title       string
	GeneratedAt time.Time
	Params      map[string]string

	// Currency is the app's money-display configuration, passed through from the
	// caller for the renderers to interpret. Nil when unconfigured, in which case
	// renderers fall back to their bare numeric formatting.
	Currency  *CurrencyFormat
	NoneLabel string
}

// CurrencyFormat is the app's money-display configuration — a renderer hint the
// engine passes through untouched (it never interprets it). It mirrors the
// desktop's customCurrency pipe so a rendered report matches the rest of the UI.
type CurrencyFormat struct {
	Symbol             string // e.g. "$", "€"; "" renders no symbol
	SymbolAtEnd        bool   // false = symbol leads (START), true = symbol trails (END)
	ThousandsSeparator string // grouping separator, e.g. "," or "."
	DecimalSeparator   string // decimal point, e.g. "." or ","
	HideDecimals       bool   // drop the fractional part entirely
}

// ColumnDescriptor tells a renderer what a column is.
type ColumnDescriptor struct {
	Name  string
	Label string
	Kind  ColumnKind

	// DataType is the type of the column's cells. For an aggregate it follows
	// the measure being reduced, except COUNT which is always a plain number.
	// For arithmetic it is currency when any column it reads is currency, since
	// money combined with a count is still money.
	DataType DataType

	Format string

	// Agg is set for an aggregate column.
	Agg *Aggregate

	// Expr is set for an arithmetic column: the parsed, whitelisted syntax tree
	// of ExprSrc. It is exported so a renderer can walk it — the spreadsheet
	// renderer translates it into a live cell formula rather than emitting the
	// computed value.
	Expr    ast.Node
	ExprSrc string
}

// GroupNode is one bucket of the report tree.
//
// A node holds either Children or DetailRows, never both: the leaves of the
// grouping are where detail lives.
type GroupNode struct {
	// Dimension is the field this node's siblings were bucketed by. It is empty
	// on the synthetic root.
	Dimension FieldKey

	// Value is the bucket's value. It is null exactly when IsNone is set.
	Value  Value
	IsNone bool

	// RecordCount is how many rows were attributed to this node. A row carrying
	// two tags is attributed to both tag buckets, so counts double up on
	// purpose.
	RecordCount int

	Children   []GroupNode
	DetailRows []DetailRow

	// Subtotals is non-nil only when the spec asked for them.
	Subtotals []Cell
}

// DetailRow is one of the report's bottom rows.
type DetailRow struct {
	// Dimension and Value identify the bucket an aggregated detail row sums.
	// Both are zero in records mode, where the row is a single source record and
	// nothing keys it.
	Dimension FieldKey
	Value     Value
	IsNone    bool

	Cells []Cell
}

// Cell is one column's contribution to one row.
//
// Aggregate and arithmetic cells always hold exactly one value, which may be
// null. A label cell in records mode may hold several, because a record can
// carry several categories or tags; a renderer decides how to join them.
type Cell struct {
	Column string
	Values []Value
}

// Value returns a cell's single value, or null when it has none. Use Values
// directly to read a multi-value label cell.
func (c Cell) Value() Value {
	if len(c.Values) == 0 {
		return Null()
	}
	return c.Values[0]
}
