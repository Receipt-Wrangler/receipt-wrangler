package commands

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
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

// Every receipt date key is a valid period date field on both a preset and a
// custom period, and so is an empty one (a template saved before the field
// existed).
func TestReportRequestCommand_Validate_AcceptsEveryPeriodDateField(t *testing.T) {
	periods := map[string]ReportPeriod{
		"preset": {Preset: ReportPeriodThisMonth},
		"custom": {Preset: ReportPeriodCustom, StartDate: "2026-05-01", EndDate: "2026-05-31"},
	}
	for periodName, period := range periods {
		for _, dateField := range append([]string{""}, ReceiptDateFilterKeys()...) {
			t.Run(periodName+"/"+dateField, func(t *testing.T) {
				command := validReportCommand()
				command.Period = period
				command.Period.DateField = dateField
				if errs := command.Validate().Errors; len(errs) > 0 {
					t.Errorf("expected date field %q to be valid, got %v", dateField, errs)
				}
			})
		}
	}
}

// A missing preset still reports the preset, not the date field, so the user is
// told about the thing actually missing.
func TestReportRequestCommand_Validate_PresetErrorWinsOverDateField(t *testing.T) {
	command := validReportCommand()
	command.Period = ReportPeriod{DateField: "paidAt"}
	if got := command.Validate().Errors["period"]; got != "A reporting period is required" {
		t.Errorf("period error = %q, want the missing-preset message", got)
	}
}

func TestReportRequestCommand_LoadDataFromRequest_ReadsPeriodDateField(t *testing.T) {
	for _, dateField := range ReceiptDateFilterKeys() {
		t.Run(dateField, func(t *testing.T) {
			body := `{"groupIds": ["1"], "period": {"preset": "this_month", "dateField": "` + dateField + `"}}`
			request := httptest.NewRequest("POST", "/api/report/generate", strings.NewReader(body))

			command := ReportRequestCommand{}
			if err := command.LoadDataFromRequest(httptest.NewRecorder(), request); err != nil {
				t.Fatalf("LoadDataFromRequest: %v", err)
			}
			if command.Period.DateField != dateField {
				t.Errorf("period date field = %q, want %q", command.Period.DateField, dateField)
			}
		})
	}
}

func TestReportPeriod_DateFilterKey(t *testing.T) {
	if got := (ReportPeriod{}).DateFilterKey(); got != ReceiptFilterKeyDate {
		t.Errorf("empty date field = %q, want the receipt date %q", got, ReceiptFilterKeyDate)
	}
	for _, dateField := range ReceiptDateFilterKeys() {
		if got := (ReportPeriod{DateField: dateField}).DateFilterKey(); got != dateField {
			t.Errorf("date field %q = %q, want it passed through", dateField, got)
		}
	}
}

// TestReportPeriod_MarshalOmitsEmptyDateField guards the omitempty tag: templates
// are stored with json.Marshal, and an empty `"dateField":""` in that blob is a
// value no client ever sent, where an absent key already means the receipt date.
func TestReportPeriod_MarshalOmitsEmptyDateField(t *testing.T) {
	empty, err := json.Marshal(ReportPeriod{Preset: ReportPeriodThisMonth})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if strings.Contains(string(empty), "dateField") {
		t.Errorf("an empty date field should be omitted, got: %s", empty)
	}

	set, err := json.Marshal(ReportPeriod{Preset: ReportPeriodThisMonth, DateField: ReceiptFilterKeyCreatedAt})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if !strings.Contains(string(set), `"dateField":"createdAt"`) {
		t.Errorf("a set date field should be kept, got: %s", set)
	}
}

// TestReportPeriodDateFieldIsAnOpenStringOnTheContract pins the swagger shape of
// ReportPeriod.dateField. ReportPeriod rides inside ReportTemplate.configuration, a
// response the mobile client deserializes; were dateField a closed enum, the
// generated dart-dio EnumClass would throw on any date key added later and fail
// the whole template payload on every already-released build. So it must stay a
// plain optional string, and its description must name every accepted key.
func TestReportPeriodDateFieldIsAnOpenStringOnTheContract(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "swagger.yml"))
	if err != nil {
		t.Fatalf("read swagger.yml: %v", err)
	}

	var doc struct {
		Components struct {
			Schemas map[string]struct {
				Required   []string `yaml:"required"`
				Properties map[string]struct {
					Type        string   `yaml:"type"`
					Enum        []string `yaml:"enum"`
					Ref         string   `yaml:"$ref"`
					Description string   `yaml:"description"`
				} `yaml:"properties"`
			} `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse swagger.yml: %v", err)
	}

	period, ok := doc.Components.Schemas["ReportPeriod"]
	if !ok {
		t.Fatal("swagger.yml is missing the ReportPeriod schema")
	}
	dateField, ok := period.Properties["dateField"]
	if !ok {
		t.Fatal("ReportPeriod is missing the dateField property")
	}
	if dateField.Type != "string" || len(dateField.Enum) > 0 || dateField.Ref != "" {
		t.Errorf("dateField must be a plain string, got type=%q enum=%v $ref=%q", dateField.Type, dateField.Enum, dateField.Ref)
	}
	for _, required := range period.Required {
		if required == "dateField" {
			t.Error("dateField must stay optional; templates saved before it existed omit it")
		}
	}
	for _, key := range ReceiptDateFilterKeys() {
		if !strings.Contains(dateField.Description, key) {
			t.Errorf("dateField description does not name the accepted key %q", key)
		}
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
		// The period's date field is a filter JSON key, never a column name, a
		// differently cased key, or a filter key that isn't a date.
		{"date field as column name created_at", func(c *ReportRequestCommand) { c.Period.DateField = "created_at" }, "period"},
		{"date field as column name resolved_date", func(c *ReportRequestCommand) { c.Period.DateField = "resolved_date" }, "period"},
		{"date field wrong case Date", func(c *ReportRequestCommand) { c.Period.DateField = "Date" }, "period"},
		{"date field wrong case CreatedAt", func(c *ReportRequestCommand) { c.Period.DateField = "CreatedAt" }, "period"},
		{"date field with whitespace", func(c *ReportRequestCommand) { c.Period.DateField = " date" }, "period"},
		{"date field non-date filter key amount", func(c *ReportRequestCommand) { c.Period.DateField = "amount" }, "period"},
		{"date field non-date filter key name", func(c *ReportRequestCommand) { c.Period.DateField = "name" }, "period"},
		{"date field unknown", func(c *ReportRequestCommand) { c.Period.DateField = "paidAt" }, "period"},
		{"date field unknown on a valid custom range", func(c *ReportRequestCommand) {
			c.Period = ReportPeriod{Preset: ReportPeriodCustom, StartDate: "2026-05-01", EndDate: "2026-05-31", DateField: "paidAt"}
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
