package reporting

// ColumnKind decides how a column is computed, and therefore how it rolls up.
// The kind is declared, not inferred from whether summing happens to work.
type ColumnKind uint8

const (
	// ColumnLabel displays a field's value. It cuts rather than measures, so it
	// never rolls up: a label cell on a subtotal row is empty.
	ColumnLabel ColumnKind = iota

	// ColumnAggregate reduces rows vertically, and rolls up by merging its
	// children's accumulators.
	ColumnAggregate

	// ColumnArithmetic combines other columns horizontally, and is recomputed
	// from the other columns on the same row at every level. It is never summed.
	//
	// For an additive formula such as Total = Subtotal + Hst the two give the
	// same answer, so a subtotal matches what a spreadsheet's =SUM would show.
	// For a non-linear one such as Avg = Total / Count, summing or averaging the
	// computed column is nonsense while recomputing stays correct. One rule
	// covers both.
	ColumnArithmetic
)

func (k ColumnKind) String() string {
	switch k {
	case ColumnLabel:
		return "label"
	case ColumnAggregate:
		return "aggregate"
	case ColumnArithmetic:
		return "arithmetic"
	}
	return "unknown"
}

// Column is one column of the report.
type Column struct {
	// Name identifies the column to arithmetic formulas. It must be a plain
	// identifier, unique within the spec, and not a reserved word.
	//
	// It is deliberately separate from Label so that renaming what a reader sees
	// cannot break a formula that references it.
	Name string

	// Label is the heading a renderer draws. It is free text, and defaults to
	// Name when empty.
	Label string

	Kind ColumnKind

	// Field is the field a ColumnLabel displays.
	Field FieldKey

	// Agg is the reduction a ColumnAggregate applies. It is the canonical form;
	// AggSrc is parsed into it when set.
	Agg Aggregate

	// AggSrc optionally carries a ColumnAggregate in the source form a persisted
	// template stores, such as "SUM(amount)". When present it overrides Agg.
	AggSrc string

	// Expr is a ColumnArithmetic's expression over other column names, such as
	// "Subtotal + Hst".
	Expr string

	// Format is an opaque per-column presentation override, such as
	// "$ #,##0.00". The engine carries it through without interpreting it;
	// renderers decide what it means.
	Format string
}

// heading returns the text a renderer should draw for the column.
func (c Column) heading() string {
	if len(c.Label) > 0 {
		return c.Label
	}
	return c.Name
}
