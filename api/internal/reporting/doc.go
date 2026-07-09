// Package reporting is a pure report engine: it turns a report definition plus
// already-fetched data into a format-agnostic ReportModel.
//
// The engine fetches nothing, renders nothing, and mutates nothing. It imports
// no database, HTTP, or model packages — the only third-party dependencies are
// shopspring/decimal and expr-lang/expr. Producers map their own domain objects
// into Rows before calling Run; see internal/reporting/receiptsource for the
// receipt mapping. Consumers (CSV/XLSX/PDF renderers, a dashboard widget) read
// the ReportModel. Because Run is pure, the same inputs always yield the same
// output.
//
// The pipeline is:
//
//	rows -> group -> aggregate -> ReportModel
//
// Four invariants hold everywhere in this package, and each is covered by tests:
//
// Money is decimal, never float. Division goes through DivRound at a configured
// scale, and a zero divisor yields a null cell rather than a panic. The
// package-level decimal.DivisionPrecision global is never read or written.
//
// Aggregate columns roll up by merging their children's accumulators, never by
// combining finalized values. This is what makes AVG correct at every level: a
// subtotal's average is sum(all descendants) / count(all descendants), not the
// average of its children's averages. Arithmetic columns are recomputed from
// the other columns on the same row at every level (detail, subtotal, grand
// total), which keeps non-linear formulas such as ratios and averages correct
// where summing the computed column would not be.
//
// Multi-value dimensions fan out. A row carrying two tags is attributed in full
// to both tag buckets, so it double-counts, and that double-count propagates to
// the grand total. This is deliberate and matches the dashboard pie chart's
// attribution (see services/pie_chart.go). An empty dimension value becomes an
// explicit (None) bucket and is never silently dropped.
//
// Output is deterministic. Buckets sort by their typed value with (None) last,
// records preserve input order, and no output is ever produced by ranging over
// a Go map.
//
// The engine emits raw typed values. Currency symbols, separators, and decimal
// places are presentation, and belong to renderers.
package reporting
