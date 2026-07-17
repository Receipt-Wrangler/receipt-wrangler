package commands

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// Report format identifiers. A request asks for one or more; multiple formats are
// bundled into a single zip.
const (
	ReportFormatCsv  = "csv"
	ReportFormatXlsx = "xlsx"
	ReportFormatPdf  = "pdf"
)

// Detail modes: one row per receipt, or one rolled-up row per distinct value of a
// dimension.
const (
	ReportDetailRecords   = "records"
	ReportDetailAggregate = "aggregate"
)

// Column kinds. Mirrors the reporting engine's ColumnKind so the service can map
// each one straight onto a reporting.Column.
const (
	ReportColumnDimension = "dimension"
	ReportColumnAggregate = "aggregate"
	ReportColumnFormula   = "formula"
)

// CurrentReportConfigurationVersion is the schema version stamped onto a saved
// report template's stored configuration. Bump it and write an upcaster when a
// breaking change to the ReportRequestCommand shape lands.
const CurrentReportConfigurationVersion = 1

// maxReportColumns caps how many columns one report may declare. Arithmetic
// columns may reference one another, so an unbounded column list is a way to
// chain expensive per-cell work; the engine also bounds the magnitude of any
// single computed value, but rejecting an absurd column count here turns that
// into a clean 400 instead of a silently-nulled cell. Real reports use a
// handful of columns — this is a generous ceiling, not a design target.
const maxReportColumns = 100

// Period presets. Everything but "custom" resolves to a date window from the
// server clock at generation time; "custom" uses the supplied start/end.
const (
	ReportPeriodThisMonth = "this_month"
	ReportPeriodLastMonth = "last_month"
	ReportPeriodMtd       = "mtd"
	ReportPeriodQtd       = "qtd"
	ReportPeriodYtd       = "ytd"
	ReportPeriodCustom    = "custom"
)

// reportDateLayout is the wire format for the custom period bounds.
const reportDateLayout = "2006-01-02"

var validReportFormats = map[string]bool{ReportFormatCsv: true, ReportFormatXlsx: true, ReportFormatPdf: true}
var validReportPresets = map[string]bool{
	ReportPeriodThisMonth: true, ReportPeriodLastMonth: true, ReportPeriodMtd: true,
	ReportPeriodQtd: true, ReportPeriodYtd: true, ReportPeriodCustom: true,
}
var validAggFuncs = map[string]bool{"SUM": true, "COUNT": true, "AVG": true, "MIN": true, "MAX": true}

// ReportRequestCommand is the report-builder configuration a client submits to
// generate a report. It is engine-shaped: group-by, detail, and columns carry the
// reporting engine's field keys and (for formulas) a machine column Name that
// expressions reference — the client maps its own UI vocabulary onto these before
// submitting. The deep validity of the spec (that a field key exists, has the
// right role, and that formulas form no cycle) is left to reporting.Run.
type ReportRequestCommand struct {
	Name        string                    `json:"name"`
	GroupIds    []string                  `json:"groupIds"`
	Period      ReportPeriod              `json:"period"`
	Filter      ReceiptPagedRequestFilter `json:"filter"`
	GroupBy     []string                  `json:"groupBy"`
	Detail      ReportDetail              `json:"detail"`
	Columns     []ReportColumn            `json:"columns"`
	Subtotals   bool                      `json:"subtotals"`
	GrandTotals bool                      `json:"grandTotals"`
	Document    ReportDocument            `json:"document"`
	Formats     []string                  `json:"formats"`
}

// ReportPeriod is the reporting window. Preset is one of the ReportPeriod*
// constants; StartDate/EndDate (YYYY-MM-DD) are only read when Preset is custom.
type ReportPeriod struct {
	Preset    string `json:"preset"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// ReportDetail selects the bottom-row mode. By is the dimension an aggregate keys
// its leaf rows on, and must be empty in records mode.
type ReportDetail struct {
	Mode string `json:"mode"`
	By   string `json:"by"`
}

// ReportColumn is one output column. Name is a plain identifier that formula
// columns reference; Label is the heading. Field is set for dimensions; AggFunc
// (+ Measure, except COUNT) for aggregates; Expr for formulas.
type ReportColumn struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	Field   string `json:"field"`
	AggFunc string `json:"aggFunc"`
	Measure string `json:"measure"`
	Expr    string `json:"expr"`
}

// ReportDocument is the authored document copy. Any of the four {{variables}} the
// service knows about are resolved at generation time.
type ReportDocument struct {
	Title  string `json:"title"`
	Intro  string `json:"intro"`
	Footer string `json:"footer"`
}

func (command *ReportRequestCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &command); err != nil {
		return err
	}

	// The reporting data source reuses the receipt filter, whose grant-narrowing
	// and query code expect non-nil values — mirror the initialization the other
	// receipt-filter commands do.
	initReceiptFilterValues(&command.Filter)
	return nil
}

func (command *ReportRequestCommand) Validate() structs.ValidatorError {
	errorMap := make(map[string]string)

	if len(command.GroupIds) == 0 {
		errorMap["groupIds"] = "At least one group is required"
	} else {
		seenGroups := make(map[string]bool, len(command.GroupIds))
		for _, groupId := range command.GroupIds {
			if groupId == "" {
				errorMap["groupIds"] = "Group ids must not be empty"
				break
			}
			// A duplicate id would load and sum that group's receipts twice,
			// silently inflating the report totals.
			if seenGroups[groupId] {
				errorMap["groupIds"] = "Group ids must be unique: " + groupId
				break
			}
			seenGroups[groupId] = true
		}
	}

	command.validatePeriod(errorMap)
	command.validateDetail(errorMap)
	command.validateColumns(errorMap)
	command.validateFormats(errorMap)

	return structs.ValidatorError{Errors: errorMap}
}

func (command *ReportRequestCommand) validatePeriod(errorMap map[string]string) {
	if command.Period.Preset == "" {
		errorMap["period"] = "A reporting period is required"
		return
	}
	if !validReportPresets[command.Period.Preset] {
		errorMap["period"] = "Invalid reporting period"
		return
	}
	if command.Period.Preset == ReportPeriodCustom {
		start, startErr := time.Parse(reportDateLayout, command.Period.StartDate)
		end, endErr := time.Parse(reportDateLayout, command.Period.EndDate)
		if startErr != nil || endErr != nil {
			errorMap["period"] = "A custom period needs a valid start and end date"
			return
		}
		if end.Before(start) {
			errorMap["period"] = "The period end must not be before its start"
		}
	}
}

func (command *ReportRequestCommand) validateDetail(errorMap map[string]string) {
	switch command.Detail.Mode {
	case ReportDetailRecords:
		if command.Detail.By != "" {
			errorMap["detail"] = "Records mode does not take an aggregate dimension"
		}
	case ReportDetailAggregate:
		if command.Detail.By == "" {
			errorMap["detail"] = "Aggregate mode needs a dimension to aggregate by"
		}
	default:
		errorMap["detail"] = "Detail mode must be records or aggregate"
	}
}

func (command *ReportRequestCommand) validateColumns(errorMap map[string]string) {
	if len(command.Columns) == 0 {
		errorMap["columns"] = "At least one column is required"
		return
	}

	if len(command.Columns) > maxReportColumns {
		errorMap["columns"] = "A report may have at most " + strconv.Itoa(maxReportColumns) + " columns"
		return
	}

	seenNames := make(map[string]bool, len(command.Columns))
	for _, column := range command.Columns {
		if column.Name == "" {
			errorMap["columns"] = "Every column needs a name"
			continue
		}
		if seenNames[column.Name] {
			errorMap["columns"] = "Column names must be unique: " + column.Name
			continue
		}
		seenNames[column.Name] = true

		switch column.Kind {
		case ReportColumnDimension:
			if column.Field == "" {
				errorMap["columns"] = "Dimension column " + column.Name + " needs a field"
			}
		case ReportColumnAggregate:
			if !validAggFuncs[column.AggFunc] {
				errorMap["columns"] = "Aggregate column " + column.Name + " needs a valid function"
			} else if column.AggFunc != "COUNT" && column.Measure == "" {
				errorMap["columns"] = "Aggregate column " + column.Name + " needs a measure"
			}
		case ReportColumnFormula:
			if column.Expr == "" {
				errorMap["columns"] = "Formula column " + column.Name + " needs an expression"
			}
		default:
			errorMap["columns"] = "Column " + column.Name + " has an unknown kind"
		}
	}
}

func (command *ReportRequestCommand) validateFormats(errorMap map[string]string) {
	if len(command.Formats) == 0 {
		errorMap["formats"] = "At least one output format is required"
		return
	}
	for _, format := range command.Formats {
		if !validReportFormats[format] {
			errorMap["formats"] = "Unsupported format: " + format
			return
		}
	}
}
