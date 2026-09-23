package services

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
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

// --- the drill-in list -----------------------------------------------------

// pagedReceiptNames reads the receipt names off the drill-in list, in its order.
func pagedReceiptNames(t *testing.T, paged structs.PagedData) []string {
	t.Helper()
	names := make([]string, 0, len(paged.Data))
	for _, item := range paged.Data {
		receipt, ok := item.(models.Receipt)
		if !ok {
			t.Fatalf("drill-in item is %T, want models.Receipt", item)
		}
		names = append(names, receipt.Name)
	}
	return names
}

// listPeriodReceipts runs the drill-in list for a command and returns its sorted
// names and total count.
func listPeriodReceipts(t *testing.T, userId uint, command commands.ReportRequestCommand, now time.Time) ([]string, int64) {
	t.Helper()
	paged, err := NewReportService(nil).receipts(userId, command, now)
	if err != nil {
		t.Fatalf("receipts: %v", err)
	}
	names := pagedReceiptNames(t, paged)
	sort.Strings(names)
	return names, paged.TotalCount
}

// The drill-in lists exactly the receipts the report counts, on every date field.
// Both go through prepareReportFilter and the same fetch, which is what this pins.
func TestReportService_Receipts_MatchTheReportOnEachDateField(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-fields", "Household")
	june15 := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	january := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	seedReceiptOnDateField(t, "date-hit", userId, groupIds[0], commands.ReceiptFilterKeyDate, june15)
	seedReceiptOnDateField(t, "resolved-hit", userId, groupIds[0], commands.ReceiptFilterKeyResolvedDate, june15)
	seedReceiptOnDateField(t, "created-hit", userId, groupIds[0], commands.ReceiptFilterKeyCreatedAt, june15)
	seedPeriodReceipt(t, "none", userId, groupIds[0], january, &january, january)

	for _, dateField := range append([]string{""}, commands.ReceiptDateFilterKeys()...) {
		t.Run("dateField="+dateField, func(t *testing.T) {
			command := periodRecordsCommand(groupIds, commands.ReportPeriod{
				Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
				DateField: dateField,
			})
			reportNames, reportCount := buildPeriodReport(t, userId, command, june15)
			listNames, listCount := listPeriodReceipts(t, userId, command, june15)

			if len(listNames) != 1 || !reflect.DeepEqual(listNames, reportNames) {
				t.Errorf("drill-in lists %v, report covers %v; want the same single receipt", listNames, reportNames)
			}
			if listCount != int64(reportCount) {
				t.Errorf("drill-in total = %d, report count = %d", listCount, reportCount)
			}
		})
	}
}

// The period resolves on the server clock, in its time zone, for the list as for
// the report. A receipt added at 20:00 on May 31 in Los Angeles is already June 1
// in UTC: a server in Los Angeles counts it in May, and so does its drill-in,
// whatever time zone the browser asking for the list is in.
func TestReportService_Receipts_AgreeWithTheReportAcrossTimeZones(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	losAngeles, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-tz", "Household")
	// Stored in the server's own zone, as a server running there stores it.
	addedMay31 := time.Date(2026, 5, 31, 20, 0, 0, 0, losAngeles)
	seedPeriodReceipt(t, "added-may-31", userId, groupIds[0], periodFarDate, nil, addedMay31)

	now := time.Date(2026, 6, 20, 12, 0, 0, 0, losAngeles)
	tests := []struct {
		start, end string
		want       []string
	}{
		{"2026-05-01", "2026-05-31", []string{"added-may-31"}},
		{"2026-06-01", "2026-06-30", []string{}},
	}
	for _, test := range tests {
		t.Run(test.start, func(t *testing.T) {
			command := periodRecordsCommand(groupIds, commands.ReportPeriod{
				Preset: commands.ReportPeriodCustom, StartDate: test.start, EndDate: test.end,
				DateField: commands.ReceiptFilterKeyCreatedAt,
			})
			reportNames, _ := buildPeriodReport(t, userId, command, now)
			listNames, _ := listPeriodReceipts(t, userId, command, now)

			if !reflect.DeepEqual(reportNames, test.want) {
				t.Errorf("report covers %v, want %v", reportNames, test.want)
			}
			if !reflect.DeepEqual(listNames, test.want) {
				t.Errorf("drill-in lists %v, want %v", listNames, test.want)
			}
		})
	}
}

// A saved filter's "whoever generates the report" paid-by sentinel resolves to the
// caller in the list too; sent as-is it would match no receipt at all.
func TestReportService_Receipts_ResolveTheReportGeneratorPaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-payer", "Household")
	otherPayer := makeUser(t, "rpt-drill-other-payer")
	june15 := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	seedPeriodReceipt(t, "mine", userId, groupIds[0], june15, nil, june15)
	seedPeriodReceipt(t, "theirs", otherPayer, groupIds[0], june15, nil, june15)

	command := periodRecordsCommand(groupIds, commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
	})
	command.Filter.PaidBy = commands.PagedRequestField{Operation: commands.CONTAINS, Value: []interface{}{float64(-1)}}

	names, count := listPeriodReceipts(t, userId, command, june15)
	if !reflect.DeepEqual(names, []string{"mine"}) || count != 1 {
		t.Errorf("drill-in lists %v (total %d), want only the caller's receipt", names, count)
	}
}

// Every custom field value carries its definition. Without it the value
// serializes a zero definition with "type":"", which the mobile client's closed
// CustomFieldType enum cannot read.
func TestReportService_Receipts_CarryCustomFieldDefinitions(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-custom", "Household")
	customField := models.CustomField{Name: "HST", Type: models.CURRENCY}
	if err := repositories.GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}
	hst := decimal.RequireFromString("15.60")
	receipt := models.Receipt{
		Name:         "with-hst",
		Amount:       decimal.NewFromInt(100),
		Date:         time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		PaidByUserID: userId,
		GroupId:      groupIds[0],
		Status:       models.OPEN,
		CustomFields: []models.CustomFieldValue{{CustomFieldId: customField.ID, CurrencyValue: &hst}},
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("create receipt: %v", err)
	}

	command := periodRecordsCommand(groupIds, commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
	})
	paged, err := NewReportService(nil).receipts(userId, command, time.Now())
	if err != nil {
		t.Fatalf("receipts: %v", err)
	}
	if len(paged.Data) != 1 {
		t.Fatalf("drill-in lists %d receipts, want 1", len(paged.Data))
	}
	values := paged.Data[0].(models.Receipt).CustomFields
	if len(values) != 1 || values[0].CustomField.Type != models.CURRENCY {
		t.Errorf("custom field values = %+v, want one carrying its CURRENCY definition", values)
	}
}

// Receipts from several groups come back as one list, newest receipt date first.
func TestReportService_Receipts_MergeGroupsNewestFirst(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-order", "Household", "Roommates")
	day := func(d int) time.Time { return time.Date(2026, 6, d, 0, 0, 0, 0, time.UTC) }
	seedPeriodReceipt(t, "household-10", userId, groupIds[0], day(10), nil, day(10))
	seedPeriodReceipt(t, "roommates-20", userId, groupIds[1], day(20), nil, day(20))
	seedPeriodReceipt(t, "household-5", userId, groupIds[0], day(5), nil, day(5))
	seedPeriodReceipt(t, "roommates-15", userId, groupIds[1], day(15), nil, day(15))

	command := periodRecordsCommand(groupIds, commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
	})
	paged, err := NewReportService(nil).receipts(userId, command, day(30))
	if err != nil {
		t.Fatalf("receipts: %v", err)
	}
	want := []string{"roommates-20", "roommates-15", "household-10", "household-5"}
	if got := pagedReceiptNames(t, paged); !reflect.DeepEqual(got, want) {
		t.Errorf("drill-in order = %v, want %v", got, want)
	}
}

// The list is capped, but its total still counts every receipt the report covers,
// so the drill-in's subtitle agrees with the preview's count chip.
func TestReportService_Receipts_CapTheListButCountEveryReceipt(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-cap", "Household")
	receipts := make([]models.Receipt, reportReceiptsCap+1)
	for index := range receipts {
		receipts[index] = models.Receipt{
			Name:         fmt.Sprintf("r%03d", index),
			Amount:       decimal.NewFromInt(1),
			Date:         time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
			PaidByUserID: userId,
			GroupId:      groupIds[0],
			Status:       models.OPEN,
		}
	}
	if err := repositories.GetDB().Create(&receipts).Error; err != nil {
		t.Fatalf("create receipts: %v", err)
	}

	command := periodRecordsCommand(groupIds, commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
	})
	paged, err := NewReportService(nil).receipts(userId, command, time.Now())
	if err != nil {
		t.Fatalf("receipts: %v", err)
	}
	if len(paged.Data) != reportReceiptsCap {
		t.Errorf("drill-in lists %d receipts, want the cap of %d", len(paged.Data), reportReceiptsCap)
	}
	if paged.TotalCount != int64(reportReceiptsCap+1) {
		t.Errorf("total = %d, want %d", paged.TotalCount, reportReceiptsCap+1)
	}
}

// With several groups, each loads only its newest receipts, yet the list is still
// exactly the newest reportReceiptsCap across all of them and the total counts
// every receipt in every group.
func TestReportService_Receipts_CapAcrossGroupsKeepsTheNewestAndTheFullTotal(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-drill-cap-groups", "Household", "Roommates")
	type seeded struct {
		name string
		date time.Time
	}
	var all []seeded
	seedGroup := func(groupId uint, prefix string, count int, offset int) {
		receipts := make([]models.Receipt, count)
		for index := range receipts {
			// Interleave the two groups' dates, so the newest overall draws from both.
			date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC).Add(time.Duration(index*2+offset) * time.Minute)
			name := fmt.Sprintf("%s-%03d", prefix, index)
			receipts[index] = models.Receipt{
				Name:         name,
				Amount:       decimal.NewFromInt(1),
				Date:         date,
				PaidByUserID: userId,
				GroupId:      groupId,
				Status:       models.OPEN,
			}
			all = append(all, seeded{name, date})
		}
		if err := repositories.GetDB().Create(&receipts).Error; err != nil {
			t.Fatalf("create receipts: %v", err)
		}
	}
	seedGroup(groupIds[0], "household", 150, 0)
	seedGroup(groupIds[1], "roommates", 120, 1)

	sort.Slice(all, func(i, j int) bool { return all[i].date.After(all[j].date) })
	want := make([]string, reportReceiptsCap)
	for index := range want {
		want[index] = all[index].name
	}

	command := periodRecordsCommand(groupIds, commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-06-01", EndDate: "2026-06-30",
	})
	paged, err := NewReportService(nil).receipts(userId, command, time.Now())
	if err != nil {
		t.Fatalf("receipts: %v", err)
	}
	if paged.TotalCount != 270 {
		t.Errorf("total = %d, want 270", paged.TotalCount)
	}
	if got := pagedReceiptNames(t, paged); !reflect.DeepEqual(got, want) {
		t.Errorf("drill-in lists %d receipts, want the %d newest across both groups in order", len(got), len(want))
	}
}

// A limited fetch loads at most the limit, newest first, while its count is every
// receipt the caller may see: the count is taken after paid-by visibility, so a
// hidden payer's receipts are in neither.
func TestReportDataService_Receipts_LimitsTheFetchButCountsWhatTheCallerMaySee(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	allowedPayer := makeUser(t, "rpt-drill-allowed-payer")
	hiddenPayer := makeUser(t, "rpt-drill-hidden-payer")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "rpt-drill-reviewer", []uint{allowedPayer}, false)

	day := func(d int) time.Time { return time.Date(2026, 6, d, 0, 0, 0, 0, time.UTC) }
	seedPeriodReceipt(t, "visible-1", allowedPayer, groupId, day(1), nil, day(1))
	seedPeriodReceipt(t, "visible-2", allowedPayer, groupId, day(2), nil, day(2))
	seedPeriodReceipt(t, "visible-3", allowedPayer, groupId, day(3), nil, day(3))
	seedPeriodReceipt(t, "hidden-4", hiddenPayer, groupId, day(4), nil, day(4))
	seedPeriodReceipt(t, "hidden-5", hiddenPayer, groupId, day(5), nil, day(5))

	receipts, count, err := NewReportDataService(nil).Receipts(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{}, 2)
	if err != nil {
		t.Fatalf("Receipts: %v", err)
	}
	names := make([]string, len(receipts))
	for index, receipt := range receipts {
		names[index] = receipt.Name
	}
	if want := []string{"visible-3", "visible-2"}; !reflect.DeepEqual(names, want) {
		t.Errorf("fetched %v, want the newest two visible receipts %v", names, want)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3 (every visible receipt, none hidden)", count)
	}
}
