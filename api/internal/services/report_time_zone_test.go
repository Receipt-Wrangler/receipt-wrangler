package services

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting/render"
	"receipt-wrangler/api/internal/repositories"
)

func setAppTimeZone(t *testing.T, name string) {
	t.Helper()
	if err := repositories.SetAppTimeZoneForTests(name); err != nil {
		t.Fatalf("set app time zone %s: %v", name, err)
	}
}

// The reported bug. A receipt added at 9PM on Sep 30 Eastern is 01:00 UTC on
// Oct 1, so a "Last Month" report on Added At, run in October, left it out: the
// month was worked out in the server's zone. With the app zone set to New York
// it is included — on the very next report, through the real settings save, with
// no restart — and switching back to UTC moves it out again.
func TestReportService_LastMonthOnAddedAtFollowsTheAppTimeZone(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-tz-bug", "Household")
	seedPeriodReceipt(t, "late-sep-30-eastern", userId, groupIds[0], periodFarDate, nil,
		time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC))

	command := periodRecordsCommand(groupIds, commands.ReportPeriod{
		Preset:    commands.ReportPeriodLastMonth,
		DateField: commands.ReceiptFilterKeyCreatedAt,
	})
	now := time.Date(2026, 10, 8, 15, 0, 0, 0, time.UTC)

	steps := []struct {
		zone string
		want []string
	}{
		{"UTC", []string{}},
		{"America/New_York", []string{"late-sep-30-eastern"}},
		{"UTC", []string{}},
	}
	for _, step := range steps {
		setAppTimeZone(t, step.zone)

		reportNames, _ := buildPeriodReport(t, userId, command, now)
		if !reflect.DeepEqual(reportNames, step.want) {
			t.Errorf("%s: Last Month on Added At covers %v, want %v", step.zone, reportNames, step.want)
		}
		listNames, _ := listPeriodReceipts(t, userId, command, now)
		if !reflect.DeepEqual(listNames, step.want) {
			t.Errorf("%s: drill-in lists %v, want %v", step.zone, listNames, step.want)
		}
	}
}

// A receipt's Date is a calendar day stored as midnight UTC. It never shifts
// with the app zone: Sep 1 is in a September period in New York, in UTC and in
// Sydney, and in no August period.
func TestReportService_CalendarDateIsUnaffectedByTheAppTimeZone(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-tz-date", "Household")
	seedPeriodReceipt(t, "dated-sep-1", userId, groupIds[0],
		time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), nil, periodFarDate)

	period := func(start, end string) commands.ReportRequestCommand {
		return periodRecordsCommand(groupIds, commands.ReportPeriod{
			Preset: commands.ReportPeriodCustom, StartDate: start, EndDate: end,
			DateField: commands.ReceiptFilterKeyDate,
		})
	}
	now := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)

	for _, zone := range []string{"UTC", "America/New_York", "Australia/Sydney"} {
		setAppTimeZone(t, zone)

		if names, _ := buildPeriodReport(t, userId, period("2026-09-01", "2026-09-30"), now); !reflect.DeepEqual(names, []string{"dated-sep-1"}) {
			t.Errorf("%s: September covers %v, want the Sep 1 receipt", zone, names)
		}
		if names, _ := buildPeriodReport(t, userId, period("2026-08-01", "2026-08-31"), now); len(names) != 0 {
			t.Errorf("%s: August covers %v, want nothing", zone, names)
		}

		// "This month" on the receipt date is the app zone's month, bounded in
		// UTC days: mid-September everywhere here.
		thisMonth := periodRecordsCommand(groupIds, commands.ReportPeriod{
			Preset: commands.ReportPeriodThisMonth, DateField: commands.ReceiptFilterKeyDate,
		})
		if names, _ := buildPeriodReport(t, userId, thisMonth, now); !reflect.DeepEqual(names, []string{"dated-sep-1"}) {
			t.Errorf("%s: This Month on Date covers %v, want the Sep 1 receipt", zone, names)
		}
	}
}

// Presets, the period label and {{generatedAt}} are all read off the app
// zone's clock. At 02:00 UTC on Oct 1 it is still Sep 30 in New York.
func TestReportService_PresetsLabelAndGeneratedAtUseTheAppTimeZone(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-tz-label", "Household")
	now := time.Date(2026, 10, 1, 2, 30, 0, 0, time.UTC)

	tests := []struct {
		zone          string
		preset        string
		wantPeriod    string
		wantGenerated string
	}{
		{"UTC", commands.ReportPeriodThisMonth, "2026-10-01 to 2026-10-31", "Oct 1, 2026, 2:30 AM"},
		{"America/New_York", commands.ReportPeriodThisMonth, "2026-09-01 to 2026-09-30", "Sep 30, 2026, 10:30 PM"},
		{"UTC", commands.ReportPeriodLastMonth, "2026-09-01 to 2026-09-30", "Oct 1, 2026, 2:30 AM"},
		{"America/New_York", commands.ReportPeriodLastMonth, "2026-08-01 to 2026-08-31", "Sep 30, 2026, 10:30 PM"},
		{"America/New_York", commands.ReportPeriodMtd, "2026-09-01 to 2026-09-30", "Sep 30, 2026, 10:30 PM"},
		{"America/New_York", commands.ReportPeriodYtd, "2026-01-01 to 2026-09-30", "Sep 30, 2026, 10:30 PM"},
		{"Asia/Tokyo", commands.ReportPeriodQtd, "2026-10-01 to 2026-10-01", "Oct 1, 2026, 11:30 AM"},
	}
	for _, test := range tests {
		t.Run(test.zone+"/"+test.preset, func(t *testing.T) {
			setAppTimeZone(t, test.zone)

			command := periodRecordsCommand(groupIds, commands.ReportPeriod{Preset: test.preset})
			command.Document = commands.ReportDocument{Intro: "{{period}} | {{generatedAt}}"}

			build, err := NewReportService(nil).buildModel(userId, command, now, 0)
			if err != nil {
				t.Fatalf("buildModel: %v", err)
			}

			want := test.wantPeriod + " | " + test.wantGenerated
			if build.chrome.Intro != want {
				t.Errorf("intro = %q, want %q", build.chrome.Intro, want)
			}
			if build.model.Meta.Params["Period"] != test.wantPeriod {
				t.Errorf("Period param = %q, want %q", build.model.Meta.Params["Period"], test.wantPeriod)
			}
		})
	}
}

// The derived Added At month buckets in the app zone; the receipt Date's month
// does not move; the rendered footer names the zone.
func TestReportService_BucketsAndFooterUseTheAppTimeZone(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-tz-bucket", "Household")
	seedPeriodReceipt(t, "late-sep-30-eastern", userId, groupIds[0],
		time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC), nil, time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC))

	command := func(groupBy string) commands.ReportRequestCommand {
		ids := []string{groupIdString(groupIds[0])}
		return commands.ReportRequestCommand{
			Name:     "Buckets",
			GroupIds: ids,
			Period: commands.ReportPeriod{
				Preset: commands.ReportPeriodCustom, StartDate: "2020-01-01", EndDate: "2027-12-31",
				DateField: commands.ReceiptFilterKeyDate,
			},
			GroupBy: []string{groupBy},
			Detail:  commands.ReportDetail{Mode: commands.ReportDetailRecords},
			Columns: []commands.ReportColumn{
				{Kind: commands.ReportColumnDimension, Name: "Name", Label: "Name", Field: "name"},
			},
			Formats: []string{commands.ReportFormatCsv},
		}
	}
	now := time.Date(2026, 10, 8, 15, 0, 0, 0, time.UTC)

	bucket := func(groupBy string) string {
		build, err := NewReportService(nil).buildModel(userId, command(groupBy), now, 0)
		if err != nil {
			t.Fatalf("buildModel: %v", err)
		}
		if len(build.model.Root.Children) != 1 {
			t.Fatalf("%s: want one bucket, got %d", groupBy, len(build.model.Root.Children))
		}
		text, _ := build.model.Root.Children[0].Value.Text()
		return text
	}

	setAppTimeZone(t, "UTC")
	if got := bucket("created_at_month"); got != "2026-10" {
		t.Errorf("UTC created_at_month = %q, want 2026-10", got)
	}
	if got := bucket("date_month"); got != "2026-09" {
		t.Errorf("UTC date_month = %q, want 2026-09", got)
	}

	setAppTimeZone(t, "America/New_York")
	if got := bucket("created_at_month"); got != "2026-09" {
		t.Errorf("New York created_at_month = %q, want 2026-09", got)
	}
	if got := bucket("date_month"); got != "2026-09" {
		t.Errorf("New York date_month = %q, want 2026-09", got)
	}

	// Grouping by the raw Added At labels the bucket with its New York day, and
	// the automatic footer prints in New York time with the zone named.
	build, err := NewReportService(nil).buildModel(userId, command("created_at"), now, 0)
	if err != nil {
		t.Fatalf("buildModel: %v", err)
	}
	html, err := render.HTML(build.model, build.dimensions, render.DocumentChrome{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	rendered := string(html)
	if !strings.Contains(rendered, "2026-09-30") || strings.Contains(rendered, "2026-10-01") {
		t.Errorf("raw Added At bucket should read 2026-09-30 in New York:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Generated 2026-10-08 11:00:00 EDT") {
		t.Errorf("footer should print New York time with its zone:\n%s", rendered)
	}
}

// AppData carries the resolved zone for the clients: the configured name, and
// "UTC" — never "" — when nothing usable is stored.
func TestGetAppData_CarriesTheAppTimeZone(t *testing.T) {
	defer repositories.TruncateTestDb()

	user, err := repositories.NewUserRepository(nil).CreateUser(commands.SignUpCommand{
		Username: "appdata-tz-user", Password: "Password", DisplayName: "TZ",
	})
	if err != nil {
		t.Fatal(err)
	}

	appData, err := GetAppData(user.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if appData.TimeZone != "UTC" {
		t.Errorf("default AppData.TimeZone = %q, want UTC", appData.TimeZone)
	}

	setAppTimeZone(t, "America/New_York")
	appData, err = GetAppData(user.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if appData.TimeZone != "America/New_York" {
		t.Errorf("AppData.TimeZone = %q, want America/New_York", appData.TimeZone)
	}

	if err := repositories.GetDB().Model(&models.SystemSettings{}).Where("1 = 1").Update("time_zone", "").Error; err != nil {
		t.Fatal(err)
	}
	appData, err = GetAppData(user.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if appData.TimeZone != "UTC" {
		t.Errorf("blank stored zone: AppData.TimeZone = %q, want UTC", appData.TimeZone)
	}
}
