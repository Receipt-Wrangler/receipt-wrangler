package services

import (
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/repositories"
)

// These tests run the report pipeline against real receipts to prove a period
// covers exactly the receipt date its dateField names. Every timestamp is UTC and
// every buildModel call gets a UTC now: the period's bounds are built in now's
// location, and SQLite compares the stored timestamps as text, so mixed offsets
// would make the edge cases pass or fail for the wrong reason.

var periodFarDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

// seedPeriodReceipt creates a receipt carrying the three dates a report period can
// cover. created_at is stamped after the insert, so the test owns it rather than
// GORM's create-time default.
func seedPeriodReceipt(
	t *testing.T,
	name string,
	userId uint,
	groupId uint,
	date time.Time,
	resolved *time.Time,
	created time.Time,
) {
	t.Helper()
	status := models.OPEN
	if resolved != nil {
		status = models.RESOLVED
	}
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(10),
		Date:         date,
		ResolvedDate: resolved,
		PaidByUserID: userId,
		GroupId:      groupId,
		Status:       status,
	}
	db := repositories.GetDB()
	if err := db.Create(&receipt).Error; err != nil {
		t.Fatalf("create receipt %q: %v", name, err)
	}
	if err := db.Model(&receipt).UpdateColumn("created_at", created).Error; err != nil {
		t.Fatalf("stamp created_at on %q: %v", name, err)
	}
}

// seedReceiptOnDateField creates a receipt whose dateField holds at, with its other
// dates far outside any window under test (and no resolved date at all).
func seedReceiptOnDateField(t *testing.T, name string, userId uint, groupId uint, dateField string, at time.Time) {
	t.Helper()
	date, created := periodFarDate, periodFarDate
	var resolved *time.Time
	switch dateField {
	case commands.ReceiptFilterKeyDate:
		date = at
	case commands.ReceiptFilterKeyResolvedDate:
		resolved = &at
	case commands.ReceiptFilterKeyCreatedAt:
		created = at
	default:
		t.Fatalf("unknown date field %q", dateField)
	}
	seedPeriodReceipt(t, name, userId, groupId, date, resolved, created)
}

// periodRecordsCommand is a records-mode report listing each covered receipt's name.
func periodRecordsCommand(groupIds []uint, period commands.ReportPeriod) commands.ReportRequestCommand {
	ids := make([]string, len(groupIds))
	for index, groupId := range groupIds {
		ids[index] = groupIdString(groupId)
	}
	return commands.ReportRequestCommand{
		Name:     "Period",
		GroupIds: ids,
		Period:   period,
		Detail:   commands.ReportDetail{Mode: commands.ReportDetailRecords},
		Columns: []commands.ReportColumn{
			{Kind: commands.ReportColumnDimension, Name: "Name", Label: "Name", Field: "name"},
		},
		Formats: []string{commands.ReportFormatCsv},
	}
}

// reportRecordNames reads the sorted receipt names off an ungrouped records report.
func reportRecordNames(model reporting.ReportModel) []string {
	names := []string{}
	for _, row := range model.Root.DetailRows {
		for _, cell := range row.Cells {
			if cell.Column != "Name" {
				continue
			}
			if name, ok := cell.Value().Text(); ok {
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)
	return names
}

// buildPeriodReport runs buildModel and returns the covered receipt names plus the
// report's receipt count.
func buildPeriodReport(t *testing.T, userId uint, command commands.ReportRequestCommand, now time.Time) ([]string, int) {
	t.Helper()
	build, err := NewReportService(nil).buildModel(userId, command, now, 0)
	if err != nil {
		t.Fatalf("buildModel: %v", err)
	}
	return reportRecordNames(build.model), build.receiptCount
}

// Each date field selects the one receipt whose date of that kind falls in the
// window, and nothing else — an empty field being the receipt date, as it is for a
// template saved before the field existed.
func TestReportService_BuildModel_EachDateFieldSelectsItsReceipt(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-period-fields", "Household")
	june15 := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	january := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	seedReceiptOnDateField(t, "date-hit", userId, groupIds[0], commands.ReceiptFilterKeyDate, june15)
	seedReceiptOnDateField(t, "resolved-hit", userId, groupIds[0], commands.ReceiptFilterKeyResolvedDate, june15)
	seedReceiptOnDateField(t, "created-hit", userId, groupIds[0], commands.ReceiptFilterKeyCreatedAt, june15)
	// Every date outside the window, including a resolved date that is set.
	seedPeriodReceipt(t, "none", userId, groupIds[0], january, &january, january)

	tests := []struct {
		dateField string
		want      []string
	}{
		{"", []string{"date-hit"}},
		{commands.ReceiptFilterKeyDate, []string{"date-hit"}},
		// Receipts with no resolved date never match a resolved-date period.
		{commands.ReceiptFilterKeyResolvedDate, []string{"resolved-hit"}},
		{commands.ReceiptFilterKeyCreatedAt, []string{"created-hit"}},
	}
	for _, test := range tests {
		t.Run("dateField="+test.dateField, func(t *testing.T) {
			period := commands.ReportPeriod{
				Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
				DateField: test.dateField,
			}
			names, count := buildPeriodReport(t, userId, periodRecordsCommand(groupIds, period), june15)

			if !reflect.DeepEqual(names, test.want) {
				t.Errorf("covered receipts = %v, want %v", names, test.want)
			}
			if count != len(test.want) {
				t.Errorf("receipt count = %d, want %d", count, len(test.want))
			}
		})
	}
}

// The window is inclusive at both ends on every date field: the first and last
// instants of the period are in, the instants either side of it are out.
func TestReportService_BuildModel_PeriodBoundsAreInclusivePerField(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	for _, dateField := range commands.ReceiptDateFilterKeys() {
		t.Run(dateField, func(t *testing.T) {
			userId, groupIds := seedReportUserInGroups(t, "rpt-period-bounds-"+dateField, "Bounds "+dateField)
			seeds := map[string]time.Time{
				"before":     time.Date(2026, 5, 31, 23, 59, 59, 0, time.UTC),
				"start-edge": time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
				"end-edge":   time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
				"after":      time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			}
			for name, at := range seeds {
				seedReceiptOnDateField(t, name, userId, groupIds[0], dateField, at)
			}

			period := commands.ReportPeriod{
				Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
				DateField: dateField,
			}
			now := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
			names, _ := buildPeriodReport(t, userId, periodRecordsCommand(groupIds, period), now)

			if want := []string{"end-edge", "start-edge"}; !reflect.DeepEqual(names, want) {
				t.Errorf("covered receipts = %v, want %v", names, want)
			}
		})
	}
}

// A preset period resolves from now and lands on the chosen field too: month to
// date covers the start of the month through the end of today, whichever date the
// period is on.
func TestReportService_BuildModel_PresetPeriodPerField(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	now := time.Date(2026, 6, 20, 15, 0, 0, 0, time.UTC)
	for _, dateField := range commands.ReceiptDateFilterKeys() {
		t.Run(dateField, func(t *testing.T) {
			userId, groupIds := seedReportUserInGroups(t, "rpt-period-preset-"+dateField, "Preset "+dateField)
			seedReceiptOnDateField(t, "this-month", userId, groupIds[0], dateField, time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC))
			seedReceiptOnDateField(t, "later-today", userId, groupIds[0], dateField, time.Date(2026, 6, 20, 23, 0, 0, 0, time.UTC))
			seedReceiptOnDateField(t, "tomorrow", userId, groupIds[0], dateField, time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC))
			seedReceiptOnDateField(t, "last-month", userId, groupIds[0], dateField, time.Date(2026, 5, 31, 23, 59, 59, 0, time.UTC))

			period := commands.ReportPeriod{Preset: commands.ReportPeriodMtd, DateField: dateField}
			names, _ := buildPeriodReport(t, userId, periodRecordsCommand(groupIds, period), now)

			if want := []string{"later-today", "this-month"}; !reflect.DeepEqual(names, want) {
				t.Errorf("covered receipts = %v, want %v", names, want)
			}
		})
	}
}

// The period only overwrites the slot it covers. A Date filter set in the builder
// still narrows a period on another field (both must hold), and is replaced — as
// it always has been — when the period is on the receipt date itself.
func TestReportService_BuildModel_DateFilterSurvivesANonDatePeriod(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-period-filter", "Household")
	january := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	june := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	seedPeriodReceipt(t, "jan-receipt-june-added", userId, groupIds[0], january, nil, june)
	seedPeriodReceipt(t, "june-receipt-june-added", userId, groupIds[0], june, nil, june)
	seedPeriodReceipt(t, "june-receipt-jan-added", userId, groupIds[0], june, nil, january)

	januaryFilter := commands.PagedRequestField{
		Operation: commands.BETWEEN,
		Value: []interface{}{
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
		},
	}
	june2026 := func(dateField string) commands.ReportPeriod {
		return commands.ReportPeriod{
			Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
			DateField: dateField,
		}
	}

	t.Run("period on created at keeps the Date filter", func(t *testing.T) {
		command := periodRecordsCommand(groupIds, june2026(commands.ReceiptFilterKeyCreatedAt))
		command.Filter.Date = januaryFilter
		names, _ := buildPeriodReport(t, userId, command, june)

		if want := []string{"jan-receipt-june-added"}; !reflect.DeepEqual(names, want) {
			t.Errorf("covered receipts = %v, want %v", names, want)
		}
	})

	t.Run("period on the receipt date replaces the Date filter", func(t *testing.T) {
		command := periodRecordsCommand(groupIds, june2026(commands.ReceiptFilterKeyDate))
		command.Filter.Date = januaryFilter
		names, _ := buildPeriodReport(t, userId, command, june)

		if want := []string{"june-receipt-jan-added", "june-receipt-june-added"}; !reflect.DeepEqual(names, want) {
			t.Errorf("covered receipts = %v, want %v", names, want)
		}
	})
}

// The public Preview entry point honours the date field end to end.
func TestReportService_Preview_HonorsPeriodDateField(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-period-preview", "Household")
	june15 := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	seedReceiptOnDateField(t, "date-hit", userId, groupIds[0], commands.ReceiptFilterKeyDate, june15)
	seedReceiptOnDateField(t, "created-hit", userId, groupIds[0], commands.ReceiptFilterKeyCreatedAt, june15)

	period := commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
		DateField: commands.ReceiptFilterKeyCreatedAt,
	}
	preview, err := NewReportService(nil).Preview(userId, periodRecordsCommand(groupIds, period))
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}

	if preview.ReceiptCount != 1 {
		t.Errorf("receipt count = %d, want 1", preview.ReceiptCount)
	}
	if !strings.Contains(preview.Html, "created-hit") || strings.Contains(preview.Html, "date-hit") {
		t.Errorf("preview should list only created-hit:\n%s", preview.Html)
	}
}
