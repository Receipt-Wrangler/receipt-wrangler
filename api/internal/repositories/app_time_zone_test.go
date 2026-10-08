package repositories

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
)

func mustLoadLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	return location
}

// --- the resolver -----------------------------------------------------------

func TestGetAppLocationDefaultsToUTC(t *testing.T) {
	defer TruncateTestDb()

	if location := GetAppLocation(); location != time.UTC {
		t.Errorf("GetAppLocation() = %v, want UTC on a fresh install", location)
	}

	settings, err := NewSystemSettingsRepository(nil).GetSystemSettings()
	if err != nil {
		t.Fatal(err)
	}
	if settings.TimeZone != "UTC" {
		t.Errorf("stored TimeZone = %q, want the column default \"UTC\"", settings.TimeZone)
	}
}

func TestGetAppLocationReadsTheConfiguredZone(t *testing.T) {
	defer TruncateTestDb()

	if err := SetAppTimeZoneForTests("America/New_York"); err != nil {
		t.Fatal(err)
	}

	if location := GetAppLocation(); location.String() != "America/New_York" {
		t.Errorf("GetAppLocation() = %v, want America/New_York", location)
	}
}

// A name the runtime cannot load — stored before validation existed, or by a
// build with a newer zone database — must never take a request down.
func TestGetAppLocationFallsBackToUTCForAnInvalidStoredName(t *testing.T) {
	defer TruncateTestDb()

	if _, err := NewSystemSettingsRepository(nil).GetSystemSettings(); err != nil {
		t.Fatal(err)
	}
	if err := GetDB().Model(&models.SystemSettings{}).Where("1 = 1").Update("time_zone", "Not/AZone").Error; err != nil {
		t.Fatal(err)
	}

	if location := GetAppLocation(); location != time.UTC {
		t.Errorf("GetAppLocation() = %v, want UTC for an unknown stored name", location)
	}
}

func TestResolveAppLocationFallbacks(t *testing.T) {
	for _, name := range []string{"", "UTC", "Local", "Not/AZone"} {
		if location := ResolveAppLocation(name); location != time.UTC {
			t.Errorf("ResolveAppLocation(%q) = %v, want UTC", name, location)
		}
	}

	if location := ResolveAppLocation("Asia/Tokyo"); location.String() != "Asia/Tokyo" {
		t.Errorf("ResolveAppLocation(Asia/Tokyo) = %v", location)
	}
}

// The setting is read live: a save through the real settings update is visible
// to the very next call, with no restart and no cache.
func TestGetAppLocationFollowsASettingsUpdateImmediately(t *testing.T) {
	defer TruncateTestDb()

	for _, name := range []string{"America/New_York", "Australia/Sydney", "UTC"} {
		if err := SetAppTimeZoneForTests(name); err != nil {
			t.Fatal(err)
		}
		if location := GetAppLocation(); location.String() != name {
			t.Errorf("after saving %s, GetAppLocation() = %v", name, location)
		}
	}
}

// --- the shared day filter --------------------------------------------------

// dryRunBounds builds query through build, renders it without running it, and
// returns its time.Time bind values in order. Bounds are asserted rather than
// rows because the test DB is SQLite, which compares timestamps as text: there a
// row count alone cannot tell a widened bound from a broken one.
func dryRunBounds(t *testing.T, build func(query *gorm.DB) *gorm.DB) ([]time.Time, string) {
	t.Helper()
	query := build(GetDB().Session(&gorm.Session{DryRun: true}).Model(&models.SystemTask{}))

	var results []models.SystemTask
	statement := query.Find(&results).Statement

	bounds := []time.Time{}
	for _, variable := range statement.Vars {
		if bound, ok := variable.(time.Time); ok {
			bounds = append(bounds, bound)
		}
	}
	return bounds, statement.SQL.String()
}

func assertBounds(t *testing.T, got []time.Time, sql string, want ...time.Time) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d bounds %v, want %v -- SQL: %s", len(got), got, want, sql)
	}
	for index := range want {
		if !got[index].Equal(want[index]) {
			t.Errorf("bound %d = %v, want %v -- SQL: %s", index, got[index], want[index], sql)
		}
		// Every bound reaches the driver in UTC, so it compares correctly
		// against UTC-stored timestamps on every engine.
		if got[index].Location() != time.UTC {
			t.Errorf("bound %d is in %v, want UTC", index, got[index].Location())
		}
	}
}

func TestBuildDayFilterQueryBoundsPerOperationAndZone(t *testing.T) {
	newYork := mustLoadLocation(t, "America/New_York")
	sydney := mustLoadLocation(t, "Australia/Sydney")
	repository := BaseRepository{}

	zones := map[string]*time.Location{"UTC": time.UTC, "America/New_York": newYork, "Australia/Sydney": sydney}

	for zoneName, zone := range zones {
		day := time.Date(2026, 3, 11, 0, 0, 0, 0, zone)
		nextDay := time.Date(2026, 3, 12, 0, 0, 0, 0, zone)
		previousDay := time.Date(2026, 3, 10, 0, 0, 0, 0, zone)

		tests := map[string]struct {
			operation commands.FilterOperation
			value     interface{}
			want      []time.Time
		}{
			"equals covers the whole day":             {commands.EQUALS, "2026-03-11", []time.Time{day, nextDay}},
			"greater than starts at the next day":     {commands.GREATER_THAN, "2026-03-11", []time.Time{nextDay}},
			"less than stops at the start of the day": {commands.LESS_THAN, "2026-03-11", []time.Time{day}},
			"between runs to the end of the last day": {
				commands.BETWEEN, []interface{}{"2026-03-10", "2026-03-11"}, []time.Time{previousDay, nextDay},
			},
		}

		for name, test := range tests {
			t.Run(zoneName+"/"+name, func(t *testing.T) {
				bounds, sql := dryRunBounds(t, func(query *gorm.DB) *gorm.DB {
					return repository.BuildDayFilterQuery(query, test.value, test.operation, "started_at", zone, zone)
				})
				assertBounds(t, bounds, sql, test.want...)
			})
		}
	}
}

// March 8 2026 is the day New York springs forward, so it is 23 hours long. A
// fixed 24h step would end it at 01:00 on the 9th.
func TestBuildDayFilterQueryFollowsDaylightSavingDayLengths(t *testing.T) {
	newYork := mustLoadLocation(t, "America/New_York")

	bounds, sql := dryRunBounds(t, func(query *gorm.DB) *gorm.DB {
		return BaseRepository{}.BuildDayFilterQuery(query, "2026-03-08", commands.EQUALS, "created_at", newYork, newYork)
	})
	assertBounds(t, bounds, sql,
		time.Date(2026, 3, 8, 5, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 9, 4, 0, 0, 0, time.UTC),
	)
}

// An older mobile build sends the wall time the user picked with a literal Z
// ("2026-09-01T00:00:00Z" for Sep 1 in New York), and a browser's local midnight
// arrives with its own offset. Read as written, every form names the day the
// user picked, so each gives exactly the bare day's range — in every zone.
func TestBuildDayFilterQueryReadsEveryWireFormAsTheDayWritten(t *testing.T) {
	newYork := mustLoadLocation(t, "America/New_York")
	sydney := mustLoadLocation(t, "Australia/Sydney")

	values := map[string]interface{}{
		"bare day":                  "2026-09-01",
		"old mobile fake Z":         "2026-09-01T00:00:00Z",
		"old mobile fake Z, millis": "2026-09-01T00:00:00.000Z",
		"browser local midnight":    "2026-09-01T00:00:00-04:00",
		"time.Time in its own zone": time.Date(2026, 9, 1, 23, 30, 0, 0, sydney),
	}

	for _, zone := range []*time.Location{time.UTC, newYork, sydney} {
		want := []time.Time{
			time.Date(2026, 9, 1, 0, 0, 0, 0, zone),
			time.Date(2026, 9, 2, 0, 0, 0, 0, zone),
		}
		for name, value := range values {
			t.Run(zone.String()+"/"+name, func(t *testing.T) {
				bounds, sql := dryRunBounds(t, func(query *gorm.DB) *gorm.DB {
					return BaseRepository{}.BuildDayFilterQuery(query, value, commands.EQUALS, "created_at", zone, zone)
				})
				assertBounds(t, bounds, sql, want...)
			})
		}
	}
}

// The day a value names never depends on the server process's zone: this is
// what the setting replaces. Before it, a UTC-4 browser's 2026-09-22T04:00:00Z
// became the 21st under an API in America/Los_Angeles.
func TestStartOfCalendarDayIgnoresTheProcessZone(t *testing.T) {
	for _, zoneName := range []string{"UTC", "America/Los_Angeles", "Australia/Sydney"} {
		t.Run(zoneName, func(t *testing.T) {
			original := time.Local
			time.Local = mustLoadLocation(t, zoneName)
			defer func() { time.Local = original }()

			for _, value := range []interface{}{"2026-09-22", "2026-09-22T04:00:00Z"} {
				day, ok := StartOfCalendarDay(value, time.UTC)
				if !ok {
					t.Fatalf("%v did not parse", value)
				}
				if !day.Equal(time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)) {
					t.Errorf("%v -> %v, want 2026-09-22 00:00 UTC", value, day)
				}
			}
		})
	}
}

// A malformed body is a no-op, never a panic.
func TestBuildDayFilterQueryIgnoresWrongTypedValues(t *testing.T) {
	cases := map[string]struct {
		operation commands.FilterOperation
		value     interface{}
	}{
		"number":             {commands.EQUALS, float64(20260901)},
		"unparseable string": {commands.EQUALS, "not-a-date"},
		"empty string":       {commands.LESS_THAN, ""},
		"nil":                {commands.GREATER_THAN, nil},
		"between scalar":     {commands.BETWEEN, "2026-09-01"},
		"between one bound":  {commands.BETWEEN, []interface{}{"2026-09-01"}},
		"between bad bound":  {commands.BETWEEN, []interface{}{"2026-09-01", 7}},
		"no operation":       {"", "2026-09-01"},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			bounds, sql := dryRunBounds(t, func(query *gorm.DB) *gorm.DB {
				return BaseRepository{}.BuildDayFilterQuery(query, test.value, test.operation, "created_at", time.UTC, time.UTC)
			})
			if len(bounds) != 0 {
				t.Errorf("expected no predicate, got bounds %v -- SQL: %s", bounds, sql)
			}
		})
	}
}

// WITHIN_CURRENT_MONTH is month start through the end of today, where today is
// the app zone's. At 01:00 UTC on Oct 1 it is still Sep 30 in New York, so the
// month is September there and October in UTC.
func TestCurrentMonthBoundsUseTheAppZonesToday(t *testing.T) {
	newYork := mustLoadLocation(t, "America/New_York")
	now := time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC)

	tests := []struct {
		name               string
		appZone, colZone   *time.Location
		wantStart, wantEnd time.Time
	}{
		{"instant column, UTC app", time.UTC, time.UTC,
			time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)},
		{"instant column, NY app", newYork, newYork,
			time.Date(2026, 9, 1, 4, 0, 0, 0, time.UTC), time.Date(2026, 10, 1, 4, 0, 0, 0, time.UTC)},
		// The receipt Date is a UTC calendar day: NY's September, bounded in
		// UTC days.
		{"calendar column, NY app", newYork, time.UTC,
			time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start, end := currentMonthBounds(now, test.appZone, test.colZone)
			if !start.Equal(test.wantStart) || !end.Equal(test.wantEnd) {
				t.Errorf("bounds = %v..%v, want %v..%v", start, end, test.wantStart, test.wantEnd)
			}
		})
	}
}

// --- the receipt filter -------------------------------------------------------

func receiptFilterBounds(t *testing.T, filter commands.ReceiptPagedRequestFilter) ([]time.Time, string) {
	t.Helper()
	query, err := NewReceiptRepository(nil).BuildGormFilterQuery(commands.ReceiptPagedRequestCommand{Filter: filter})
	if err != nil {
		t.Fatal(err)
	}

	var receipts []models.Receipt
	statement := query.Session(&gorm.Session{DryRun: true}).Find(&receipts).Statement
	bounds := []time.Time{}
	for _, variable := range statement.Vars {
		if bound, ok := variable.(time.Time); ok {
			bounds = append(bounds, bound)
		}
	}
	return bounds, statement.SQL.String()
}

// date is a calendar day (UTC days); resolved_date and created_at are instants
// (app-zone days).
func TestReceiptFilterReadsEachDateColumnInItsOwnZone(t *testing.T) {
	defer TruncateTestDb()

	if err := SetAppTimeZoneForTests("America/New_York"); err != nil {
		t.Fatal(err)
	}

	sept := commands.PagedRequestField{Operation: commands.EQUALS, Value: "2026-09-30"}
	newYorkDay := []time.Time{
		time.Date(2026, 9, 30, 4, 0, 0, 0, time.UTC),
		time.Date(2026, 10, 1, 4, 0, 0, 0, time.UTC),
	}

	bounds, sql := receiptFilterBounds(t, commands.ReceiptPagedRequestFilter{Date: sept})
	assertBounds(t, bounds, sql, time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC), time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC))

	bounds, sql = receiptFilterBounds(t, commands.ReceiptPagedRequestFilter{CreatedAt: sept})
	assertBounds(t, bounds, sql, newYorkDay...)

	bounds, sql = receiptFilterBounds(t, commands.ReceiptPagedRequestFilter{ResolvedDate: sept})
	assertBounds(t, bounds, sql, newYorkDay...)
}

// The receipt builder used bare .(string) / .([]interface{}) assertions and
// panicked on a malformed body.
func TestReceiptFilterIgnoresWrongTypedDateValues(t *testing.T) {
	defer TruncateTestDb()

	filter := commands.ReceiptPagedRequestFilter{
		Date:         commands.PagedRequestField{Operation: commands.EQUALS, Value: float64(7)},
		ResolvedDate: commands.PagedRequestField{Operation: commands.BETWEEN, Value: "2026-09-01"},
		CreatedAt:    commands.PagedRequestField{Operation: commands.LESS_THAN, Value: []interface{}{"x"}},
	}
	bounds, sql := receiptFilterBounds(t, filter)
	if len(bounds) != 0 {
		t.Errorf("expected no date predicate, got %v -- SQL: %s", bounds, sql)
	}
}

// The reported bug, at the receipts table: added at 9PM Sep 30 in New York is
// 01:00 UTC on Oct 1. A September "Added At" filter finds it only once the app
// zone is New York.
func TestReceiptFilterAddedAtFollowsTheAppZone(t *testing.T) {
	defer TruncateTestDb()
	setupReceiptTest()
	db := GetDB()

	late := models.Receipt{
		Name: "late Sep 30 Eastern", Amount: decimal.NewFromInt(1), Date: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC),
		PaidByUserID: 1, Status: models.OPEN, GroupId: 1,
	}
	if err := db.Create(&late).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&late).UpdateColumn("created_at", time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatal(err)
	}

	september := commands.ReceiptPagedRequestFilter{
		CreatedAt: commands.PagedRequestField{Operation: commands.BETWEEN, Value: []interface{}{"2026-09-01", "2026-09-30"}},
	}
	count := func() int64 {
		query, err := NewReceiptRepository(nil).BuildGormFilterQuery(commands.ReceiptPagedRequestCommand{Filter: september})
		if err != nil {
			t.Fatal(err)
		}
		var total int64
		if err := query.Count(&total).Error; err != nil {
			t.Fatal(err)
		}
		return total
	}

	if got := count(); got != 0 {
		t.Errorf("UTC: September matched %d, want 0 (it was Oct 1 in UTC)", got)
	}

	if err := SetAppTimeZoneForTests("America/New_York"); err != nil {
		t.Fatal(err)
	}
	if got := count(); got != 1 {
		t.Errorf("New York: September matched %d, want 1", got)
	}
}

// --- calendar dates on write ----------------------------------------------------

func TestCreateAndUpdateReceiptStoreCalendarDatesAsMidnightUTC(t *testing.T) {
	defer TruncateTestDb()
	setupReceiptTest()
	db := GetDB()

	dateField := models.CustomField{Name: "Due", Type: models.DATE}
	if err := db.Create(&dateField).Error; err != nil {
		t.Fatal(err)
	}

	newYork := mustLoadLocation(t, "America/New_York")
	tokyo := mustLoadLocation(t, "Asia/Tokyo")
	september1 := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	october5 := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)

	command := func(date time.Time, due time.Time) commands.UpsertReceiptCommand {
		return commands.UpsertReceiptCommand{
			Name: "calendar", Amount: decimal.NewFromInt(5), Date: date,
			PaidByUserID: 1, Status: models.OPEN, GroupId: 1,
			CustomFields: []commands.UpsertCustomFieldValueCommand{
				{CustomFieldId: dateField.ID, DateValue: &due},
			},
		}
	}

	assertStored := func(t *testing.T, id uint, wantDate, wantDue time.Time) {
		t.Helper()
		stored, err := NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(id))
		if err != nil {
			t.Fatal(err)
		}
		if !stored.Date.Equal(wantDate) {
			t.Errorf("Date = %v, want %v", stored.Date, wantDate)
		}
		if len(stored.CustomFields) != 1 || stored.CustomFields[0].DateValue == nil ||
			!stored.CustomFields[0].DateValue.Equal(wantDue) {
			t.Errorf("custom DateValue = %+v, want %v", stored.CustomFields, wantDue)
		}
	}

	inputs := map[string]time.Time{
		// The desktop datepicker's local midnight, with its offset.
		"browser local midnight": time.Date(2026, 9, 1, 0, 0, 0, 0, newYork),
		// What an AI returns.
		"midnight UTC": september1,
		// An older mobile build: local wall time sent with a literal Z.
		"mobile fake Z": time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		// East of Greenwich a local midnight is the previous day in UTC; as
		// written it is still Sep 1.
		"east local midnight": time.Date(2026, 9, 1, 0, 0, 0, 0, tokyo),
		// A time of day is dropped, never rounded into the next day.
		"late in the day": time.Date(2026, 9, 1, 23, 59, 0, 0, newYork),
	}

	for name, input := range inputs {
		t.Run(name, func(t *testing.T) {
			created, err := NewReceiptRepository(nil).CreateReceipt(command(input, input), 1, false)
			if err != nil {
				t.Fatal(err)
			}
			assertStored(t, created.ID, september1, september1)

			updated := time.Date(2026, 10, 5, 0, 0, 0, 0, newYork)
			if _, err := NewReceiptRepository(nil).UpdateReceipt(utils.UintToString(created.ID), command(updated, updated), 1); err != nil {
				t.Fatal(err)
			}
			assertStored(t, created.ID, october5, october5)
		})
	}
}
