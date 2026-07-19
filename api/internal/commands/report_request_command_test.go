package commands

import (
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestReportRequestCommand_LoadDataFromRequest(t *testing.T) {
	body := `{
	  "name": "R",
	  "groupIds": ["1", "2"],
	  "period": {"preset": "this_month"},
	  "detail": {"mode": "records"},
	  "columns": [{"kind": "dimension", "name": "Name", "field": "name"}],
	  "formats": ["csv"]
	}`
	request := httptest.NewRequest("POST", "/api/report/generate", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	command := ReportRequestCommand{}
	if err := command.LoadDataFromRequest(recorder, request); err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}

	if command.Name != "R" || len(command.GroupIds) != 2 || command.Period.Preset != ReportPeriodThisMonth {
		t.Errorf("body not unmarshalled: %+v", command)
	}

	// The filter is seeded with the non-nil defaults downstream grant-narrowing
	// and query building rely on.
	if _, ok := command.Filter.PaidBy.Value.([]interface{}); !ok {
		t.Errorf("PaidBy not initialized to a slice: %#v", command.Filter.PaidBy.Value)
	}
	if command.Filter.Amount.Value != float64(0) {
		t.Errorf("Amount not initialized to 0: %#v", command.Filter.Amount.Value)
	}
	if command.Filter.Date.Value != "" {
		t.Errorf("Date not initialized to empty string: %#v", command.Filter.Date.Value)
	}
}

func TestReportRequestCommand_LoadDataFromRequest_MalformedBody(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/report/generate", strings.NewReader("{not json"))
	command := ReportRequestCommand{}
	if err := command.LoadDataFromRequest(httptest.NewRecorder(), request); err == nil {
		t.Fatal("expected an error for a malformed body, got none")
	}
}

// validReportCommand is a fully valid baseline each case mutates.
func validReportCommand() ReportRequestCommand {
	return ReportRequestCommand{
		Name:     "My Report",
		GroupIds: []string{"1"},
		Period:   ReportPeriod{Preset: ReportPeriodThisMonth},
		GroupBy:  []string{"group"},
		Detail:   ReportDetail{Mode: ReportDetailAggregate, By: "category"},
		Columns: []ReportColumn{
			{Kind: ReportColumnDimension, Name: "Category", Label: "Category", Field: "category"},
			{Kind: ReportColumnAggregate, Name: "Total", Label: "Total", AggFunc: "SUM", Measure: "amount"},
			{Kind: ReportColumnAggregate, Name: "Count", Label: "Count", AggFunc: "COUNT"},
			{Kind: ReportColumnFormula, Name: "Avg", Label: "Avg", Expr: "Total / Count"},
		},
		Formats: []string{ReportFormatCsv},
	}
}

// TestReportColumn_MarshalOmitsEmptyContextualFields guards the omitempty tags on
// ReportColumn: a dimension column must not emit `"aggFunc":""`, because the
// generated mobile dart-dio ReportColumnAggFuncEnum cannot deserialize "" and the
// whole report template would silently drop out of the mobile list.
func TestReportColumn_MarshalOmitsEmptyContextualFields(t *testing.T) {
	mustMarshal := func(c ReportColumn) string {
		b, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}
		return string(b)
	}

	// Dimension: keeps field, omits every empty contextual field.
	dim := mustMarshal(ReportColumn{Kind: ReportColumnDimension, Name: "Category", Field: "category"})
	if !strings.Contains(dim, `"field":"category"`) {
		t.Errorf("dimension should keep field, got: %s", dim)
	}
	for _, absent := range []string{"aggFunc", "label", "measure", "expr"} {
		if strings.Contains(dim, `"`+absent+`"`) {
			t.Errorf("dimension should omit %q, got: %s", absent, dim)
		}
	}

	// Aggregate: keeps aggFunc (+ measure).
	agg := mustMarshal(ReportColumn{Kind: ReportColumnAggregate, Name: "Total", AggFunc: "SUM", Measure: "amount"})
	if !strings.Contains(agg, `"aggFunc":"SUM"`) || !strings.Contains(agg, `"measure":"amount"`) {
		t.Errorf("aggregate should keep aggFunc and measure, got: %s", agg)
	}

	// COUNT: keeps aggFunc, omits the (empty) measure.
	count := mustMarshal(ReportColumn{Kind: ReportColumnAggregate, Name: "Count", AggFunc: "COUNT"})
	if !strings.Contains(count, `"aggFunc":"COUNT"`) {
		t.Errorf("COUNT should keep aggFunc, got: %s", count)
	}
	if strings.Contains(count, `"measure"`) {
		t.Errorf("COUNT should omit empty measure, got: %s", count)
	}

	// Formula: keeps expr, omits aggFunc.
	formula := mustMarshal(ReportColumn{Kind: ReportColumnFormula, Name: "Avg", Expr: "Total / Count"})
	if !strings.Contains(formula, `"expr":"Total / Count"`) {
		t.Errorf("formula should keep expr, got: %s", formula)
	}
	if strings.Contains(formula, `"aggFunc"`) {
		t.Errorf("formula should omit aggFunc, got: %s", formula)
	}
}

func TestReportRequestCommand_Validate_AcceptsValid(t *testing.T) {
	command := validReportCommand()
	if errs := command.Validate().Errors; len(errs) > 0 {
		t.Fatalf("expected a valid command, got errors %v", errs)
	}
}

func TestReportRequestCommand_Validate_AcceptsRecordsMode(t *testing.T) {
	command := validReportCommand()
	command.Detail = ReportDetail{Mode: ReportDetailRecords}
	if errs := command.Validate().Errors; len(errs) > 0 {
		t.Fatalf("expected records mode to be valid, got %v", errs)
	}
}

func TestReportRequestCommand_Validate_AcceptsCustomPeriod(t *testing.T) {
	command := validReportCommand()
	command.Period = ReportPeriod{Preset: ReportPeriodCustom, StartDate: "2026-05-01", EndDate: "2026-05-31"}
	if errs := command.Validate().Errors; len(errs) > 0 {
		t.Fatalf("expected a valid custom period, got %v", errs)
	}
}

func TestReportRequestCommand_Validate_Rejects(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*ReportRequestCommand)
		wantKey string
	}{
		{"no groups", func(c *ReportRequestCommand) { c.GroupIds = nil }, "groupIds"},
		{"empty group id", func(c *ReportRequestCommand) { c.GroupIds = []string{""} }, "groupIds"},
		{"duplicate group ids", func(c *ReportRequestCommand) { c.GroupIds = []string{"1", "1"} }, "groupIds"},
		{"missing preset", func(c *ReportRequestCommand) { c.Period.Preset = "" }, "period"},
		{"unknown preset", func(c *ReportRequestCommand) { c.Period.Preset = "someday" }, "period"},
		{"custom without dates", func(c *ReportRequestCommand) { c.Period = ReportPeriod{Preset: ReportPeriodCustom} }, "period"},
		{"custom end before start", func(c *ReportRequestCommand) {
			c.Period = ReportPeriod{Preset: ReportPeriodCustom, StartDate: "2026-05-31", EndDate: "2026-05-01"}
		}, "period"},
		{"records with by", func(c *ReportRequestCommand) { c.Detail = ReportDetail{Mode: ReportDetailRecords, By: "category"} }, "detail"},
		{"aggregate without by", func(c *ReportRequestCommand) { c.Detail = ReportDetail{Mode: ReportDetailAggregate} }, "detail"},
		{"unknown mode", func(c *ReportRequestCommand) { c.Detail = ReportDetail{Mode: "pivot"} }, "detail"},
		{"no columns", func(c *ReportRequestCommand) { c.Columns = nil }, "columns"},
		{"column missing name", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{{Kind: ReportColumnDimension, Field: "category"}}
		}, "columns"},
		{"duplicate column names", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{
				{Kind: ReportColumnDimension, Name: "Dup", Field: "category"},
				{Kind: ReportColumnDimension, Name: "Dup", Field: "group"},
			}
		}, "columns"},
		{"dimension missing field", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{{Kind: ReportColumnDimension, Name: "Category"}}
		}, "columns"},
		{"aggregate bad function", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{{Kind: ReportColumnAggregate, Name: "Total", AggFunc: "MEDIAN", Measure: "amount"}}
		}, "columns"},
		{"aggregate missing measure", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{{Kind: ReportColumnAggregate, Name: "Total", AggFunc: "SUM"}}
		}, "columns"},
		{"formula missing expr", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{{Kind: ReportColumnFormula, Name: "Avg"}}
		}, "columns"},
		{"unknown column kind", func(c *ReportRequestCommand) {
			c.Columns = []ReportColumn{{Kind: "widget", Name: "X"}}
		}, "columns"},
		{"no formats", func(c *ReportRequestCommand) { c.Formats = nil }, "formats"},
		{"unsupported format", func(c *ReportRequestCommand) { c.Formats = []string{"json"} }, "formats"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := validReportCommand()
			test.mutate(&command)
			errs := command.Validate().Errors
			if _, ok := errs[test.wantKey]; !ok {
				t.Errorf("expected an error under %q, got %v", test.wantKey, errs)
			}
		})
	}
}

func TestReportRequestCommand_Validate_CountNeedsNoMeasure(t *testing.T) {
	command := validReportCommand()
	command.Columns = []ReportColumn{{Kind: ReportColumnAggregate, Name: "Count", AggFunc: "COUNT"}}
	if _, ok := command.Validate().Errors["columns"]; ok {
		t.Error("COUNT should not require a measure")
	}
}

// An unbounded column list is a way to chain expensive per-cell work (arithmetic
// columns reference one another), so the count is capped at the boundary.
func TestReportRequestCommand_Validate_BoundsColumnCount(t *testing.T) {
	dimensionColumns := func(n int) []ReportColumn {
		columns := make([]ReportColumn, n)
		for i := range columns {
			columns[i] = ReportColumn{Kind: ReportColumnDimension, Name: "c" + strconv.Itoa(i), Field: "category"}
		}
		return columns
	}

	t.Run("a report at the column ceiling is accepted", func(t *testing.T) {
		command := validReportCommand()
		command.Detail = ReportDetail{Mode: ReportDetailRecords}
		command.Columns = dimensionColumns(maxReportColumns)
		if _, ok := command.Validate().Errors["columns"]; ok {
			t.Errorf("exactly %d columns should be accepted", maxReportColumns)
		}
	})

	t.Run("one column past the ceiling is rejected", func(t *testing.T) {
		command := validReportCommand()
		command.Detail = ReportDetail{Mode: ReportDetailRecords}
		command.Columns = dimensionColumns(maxReportColumns + 1)
		if _, ok := command.Validate().Errors["columns"]; !ok {
			t.Errorf("more than %d columns should be rejected", maxReportColumns)
		}
	})
}
