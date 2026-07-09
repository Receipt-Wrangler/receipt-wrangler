package reporting

import "time"

// defaultNoneLabel is what a renderer draws for the bucket holding rows with no
// value for a dimension.
const defaultNoneLabel = "(None)"

// defaultDivisionScale bounds division to a scale far finer than any currency
// needs, so a quotient carries enough precision for a renderer to round it for
// display.
const defaultDivisionScale = 6

// DetailMode is the shape of the innermost rows of the report.
type DetailMode uint8

const (
	// DetailRecords emits one row per source record.
	DetailRecords DetailMode = iota

	// DetailAggregate emits one summed row per distinct value of a dimension.
	DetailAggregate
)

func (m DetailMode) String() string {
	switch m {
	case DetailRecords:
		return "records"
	case DetailAggregate:
		return "aggregate"
	}
	return "unknown"
}

// DetailSpec chooses what the bottom rows of the report are.
type DetailSpec struct {
	Mode DetailMode

	// By is the dimension a DetailAggregate keys its rows on. It must be empty
	// for DetailRecords.
	By FieldKey
}

// EngineConfig tunes numeric behavior.
type EngineConfig struct {
	// DivisionScale is the number of decimal places division and AVG round to.
	//
	// It exists because shopspring's Div consults decimal.DivisionPrecision, a
	// mutable process-wide global. Reading it would make a report's output
	// depend on whatever else the process had done, so the engine never touches
	// it and rounds explicitly instead.
	//
	// Zero means unset, and becomes the default. A column that should show whole
	// numbers asks for them with ROUND(x, 0) rather than by dividing coarsely.
	DivisionScale int32
}

// DefaultConfig returns the configuration a caller gets when it leaves
// EngineConfig zero.
func DefaultConfig() EngineConfig {
	return EngineConfig{DivisionScale: defaultDivisionScale}
}

// ReportSpec is a runnable report definition: what to group by, what the bottom
// rows are, and what the columns compute.
//
// It says nothing about how the result looks. Slots, branding, page size and
// currency symbols belong to whatever renders the ReportModel.
type ReportSpec struct {
	Title string

	// GroupBy is the ordered list of dimensions to nest the report by. It may
	// be empty, in which case the detail rows hang directly off the root.
	GroupBy []FieldKey

	Detail  DetailSpec
	Columns []Column

	Subtotals   bool
	GrandTotals bool

	// NoneLabel names the bucket holding rows with no value for a dimension. It
	// defaults to "(None)" and is carried to the model for renderers; the data
	// itself marks those buckets with IsNone.
	NoneLabel string

	Config EngineConfig
}

// withDefaults fills in the values a caller may leave zero. It does not
// validate; Validate does that.
func (s ReportSpec) withDefaults() ReportSpec {
	if len(s.NoneLabel) == 0 {
		s.NoneLabel = defaultNoneLabel
	}
	if s.Config.DivisionScale == 0 {
		s.Config.DivisionScale = defaultDivisionScale
	}
	return s
}

// MetaInput is the report-scope context a caller supplies at run time.
type MetaInput struct {
	// GeneratedAt is an input rather than something the engine reads from the
	// clock, so that the same inputs always produce the same output. A golden
	// test could not exist otherwise.
	GeneratedAt time.Time

	// Params are the resolved run-time parameter values, such as the period a
	// report covers. The engine carries them through to the model for a
	// renderer to substitute into its copy.
	Params map[string]string

	// CurrencyFormat is an opaque default presentation hint, overridable per
	// column. The engine does not interpret it.
	CurrencyFormat string
}
