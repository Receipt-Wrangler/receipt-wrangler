package render

import (
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/reporting"
)

func loadLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	return location
}

// A date names its calendar day in the zone its field declares. 01:00 UTC on
// Oct 1 is Sep 30 in New York, so an Added At bucket reads 2026-09-30 there;
// with no zone declared it stays the UTC day.
func TestFormatLabelValue_DateUsesTheFieldsLocation(t *testing.T) {
	moment := reporting.DateVal(time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC))
	newYork := loadLocation(t, "America/New_York")

	if got := formatLabelValue(moment, reporting.TypeDate, nil, nil, "(None)"); got != "2026-10-01" {
		t.Errorf("no location = %q, want the UTC day 2026-10-01", got)
	}
	if got := formatLabelValue(moment, reporting.TypeDate, time.UTC, nil, "(None)"); got != "2026-10-01" {
		t.Errorf("UTC = %q, want 2026-10-01", got)
	}
	if got := formatLabelValue(moment, reporting.TypeDate, newYork, nil, "(None)"); got != "2026-09-30" {
		t.Errorf("New York = %q, want 2026-09-30", got)
	}
}

// Both halves of a rendered report — the leading dimension column and a label
// column over the same kind of field — name the day in the field's zone, while
// a calendar-day field (Location UTC) is unmoved. Bucket identity is unchanged:
// the two instants stay two buckets.
func TestCSV_DateDimensionsAndLabelsUseTheFieldsLocation(t *testing.T) {
	newYork := loadLocation(t, "America/New_York")
	catalog := mustCatalog(t,
		reporting.FieldRef{Key: "created_at", Label: "Added At", DataType: reporting.TypeDate, Location: newYork},
		reporting.FieldRef{Key: "date", Label: "Date", DataType: reporting.TypeDate, Location: time.UTC},
		reporting.FieldRef{Key: "amount", Label: "Amount", DataType: reporting.TypeCurrency},
	)
	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{"created_at"},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailRecords},
		Columns: []reporting.Column{
			{Name: "Added", Label: "Added", Kind: reporting.ColumnLabel, Field: "created_at"},
			{Name: "Date", Label: "Date", Kind: reporting.ColumnLabel, Field: "date"},
			{Name: "Amount", Label: "Amount", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
		},
	}
	rows := []reporting.Row{
		{
			"created_at": {reporting.DateVal(time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC))},
			"date":       {reporting.DateVal(time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC))},
			"amount":     {money("5")},
		},
		{
			"created_at": {reporting.DateVal(time.Date(2026, 10, 1, 3, 0, 0, 0, time.UTC))},
			"date":       {reporting.DateVal(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC))},
			"amount":     {money("7")},
		},
	}

	got := mustCSV(t, mustRun(t, spec, catalog, rows),
		[]Dimension{{Key: "created_at", Label: "Added At", DataType: reporting.TypeDate, Location: newYork}})
	want := crlf(
		"Row Type,Added At,Added,Date,Amount",
		"Detail,2026-09-30,2026-09-30,2026-09-30,5.00",
		"Detail,2026-09-30,2026-09-30,2026-10-01,7.00",
	)
	if got != want {
		t.Errorf("csv mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// The automatic footer prints GeneratedAt in the app zone and names it, so a
// reader is never left guessing which clock it was read off.
func TestHTML_FooterPrintsGeneratedAtInTheAppZone(t *testing.T) {
	generatedAt := time.Date(2026, 7, 13, 9, 30, 0, 0, time.UTC)

	cases := []struct {
		location *time.Location
		want     string
	}{
		{nil, "Generated 2026-07-13 09:30:00 UTC"},
		{loadLocation(t, "America/New_York"), "Generated 2026-07-13 05:30:00 EDT"},
		{loadLocation(t, "Asia/Tokyo"), "Generated 2026-07-13 18:30:00 JST"},
	}
	for _, test := range cases {
		model, err := reporting.Run(oneLevelSpec(), paidByCatalog(t), oneLevelRows(),
			reporting.MetaInput{GeneratedAt: generatedAt, Location: test.location})
		if err != nil {
			t.Fatalf("run report: %v", err)
		}
		out, err := HTML(model, paidByDimension(), DocumentChrome{})
		if err != nil {
			t.Fatalf("HTML: %v", err)
		}
		if !strings.Contains(string(out), test.want) {
			t.Errorf("footer missing %q:\n%s", test.want, out)
		}
	}
}
