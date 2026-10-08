package repositories

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"receipt-wrangler/api/internal/commands"
)

type BaseRepository struct {
	DB *gorm.DB
	TX *gorm.DB
}

func (repository BaseRepository) GetDB() *gorm.DB {
	if repository.TX != nil {
		return repository.TX
	}

	return repository.DB
}

func (repository *BaseRepository) SetTransaction(tx *gorm.DB) {
	repository.TX = tx
}

func (repository *BaseRepository) ClearTransaction() {
	repository.TX = nil
}

func (repository BaseRepository) Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageSize == -1 {
			return db
		}

		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (repository BaseRepository) Sort(db *gorm.DB, orderBy string, sortDirection commands.SortDirection) *gorm.DB {
	desc := false
	if sortDirection == commands.DESCENDING {
		desc = true
	}

	return db.Order(clause.OrderByColumn{
		Column:  clause.Column{Name: orderBy},
		Desc:    desc,
		Reorder: false,
	})
}

func (repository BaseRepository) GetCount(table string, queryWhere string, args ...interface{}) (int64, error) {
	db := repository.GetDB()
	var result int64
	err := db.Table(table).Where(queryWhere, args...).Count(&result).Error

	return result, err
}

// BuildFilterQuery translates one {operation, value} filter field into a WHERE
// clause on fieldName. Shared by the receipt filter and the system task filter
// so the two cannot drift on what an operation means.
//
// fieldName is interpolated into the clause, so it MUST be a hardcoded column
// literal supplied by the caller — never a value taken from a request.
func (repository BaseRepository) BuildFilterQuery(runningQuery *gorm.DB, value interface{}, operation commands.FilterOperation, fieldName string, isArray bool) *gorm.DB {
	if operation == commands.EQUALS && !isArray {
		return runningQuery.Where(fmt.Sprintf("%v = ?", fieldName), value)
	}

	if operation == commands.CONTAINS && !isArray {
		searchValue := value.(string)
		searchValue = "%" + searchValue + "%"
		return runningQuery.Where(fmt.Sprintf("%v LIKE ?", fieldName), searchValue)
	}

	if operation == commands.CONTAINS && isArray {
		return runningQuery.Where(fmt.Sprintf("%v IN ?", fieldName), value)
	}

	if operation == commands.GREATER_THAN && !isArray {
		return runningQuery.Where(fmt.Sprintf("%v > ?", fieldName), value)
	}

	if operation == commands.LESS_THAN && !isArray {
		return runningQuery.Where(fmt.Sprintf("%v < ?", fieldName), value)
	}

	if operation == commands.BETWEEN {
		arrayValue, ok := value.([]interface{})
		if !ok || len(arrayValue) != 2 {
			return runningQuery
		}

		return runningQuery.Where(fmt.Sprintf("%v >= ? AND %v <= ?", fieldName, fieldName), arrayValue[0], arrayValue[1])
	}

	if operation == commands.WITHIN_CURRENT_MONTH {
		// A column with a time of day is an instant, so the month is the app
		// zone's. Date columns never reach here: they go through
		// BuildDayFilterQuery, which knows which zone their days are in.
		location := repository.AppLocation()
		start, end := currentMonthBounds(time.Now(), location, location)
		return runningQuery.Where(fmt.Sprintf("%v >= ? AND %v < ?", fieldName, fieldName), start, end)
	}

	return runningQuery
}

// BuildDayFilterQuery compares a date or timestamp column against whole
// calendar days, so EQUALS means "that day" and BETWEEN runs to the end of its
// last day. Shared by the receipt filter (date, resolved_date, created_at) and
// the system task filter (started_at, ended_at).
//
// Two zones are involved and they are not interchangeable:
//   - columnLocation is the zone the column's days are counted in. A receipt's
//     Date is a calendar day stored as midnight UTC, so it is read in UTC; an
//     instant such as created_at is read in the app zone.
//   - appLocation is the zone "today" is in, which only WITHIN_CURRENT_MONTH
//     reads. For the date column the month is the app zone's current month,
//     bounded in UTC days.
//
// The value is a calendar day, and the day it names never depends on any zone:
//   - a bare "yyyy-MM-dd" (what clients send);
//   - an RFC 3339 string, taken as the day AS WRITTEN in its own offset, never
//     converted. An older mobile build sends local wall time with a literal Z
//     ("2026-09-01T00:00:00Z" for a user who picked Sep 1 in New York), and
//     this is what lands it on the day the user picked;
//   - a time.Time, taken as the day in its own location.
//
// Every bound is converted to UTC before it reaches the driver, so it compares
// correctly against UTC-stored timestamps on every engine (SQLite compares
// them as text).
//
// Every unwrap is comma-ok: a wrong-typed or unparseable value adds no
// predicate rather than panicking the handler. column is interpolated into the
// clause, so it MUST be a hardcoded column literal, never a request value.
func (repository BaseRepository) BuildDayFilterQuery(
	query *gorm.DB,
	value interface{},
	operation commands.FilterOperation,
	column string,
	columnLocation *time.Location,
	appLocation *time.Location,
) *gorm.DB {
	if operation == commands.WITHIN_CURRENT_MONTH {
		start, end := currentMonthBounds(time.Now(), appLocation, columnLocation)
		return query.Where(column+" >= ? AND "+column+" < ?", start, end)
	}

	if operation == commands.BETWEEN {
		bounds, ok := value.([]interface{})
		if !ok || len(bounds) != 2 {
			return query
		}

		start, startOk := StartOfCalendarDay(bounds[0], columnLocation)
		end, endOk := StartOfCalendarDay(bounds[1], columnLocation)
		if !startOk || !endOk {
			return query
		}

		return query.Where(column+" >= ? AND "+column+" < ?", start.UTC(), end.AddDate(0, 0, 1).UTC())
	}

	day, ok := StartOfCalendarDay(value, columnLocation)
	if !ok {
		return query
	}

	switch operation {
	case commands.EQUALS:
		return query.Where(column+" >= ? AND "+column+" < ?", day.UTC(), day.AddDate(0, 0, 1).UTC())
	case commands.GREATER_THAN:
		return query.Where(column+" >= ?", day.AddDate(0, 0, 1).UTC())
	case commands.LESS_THAN:
		return query.Where(column+" < ?", day.UTC())
	default:
		return query
	}
}

// StartOfCalendarDay reads a filter value as a calendar day (see
// BuildDayFilterQuery for the accepted forms) and returns the midnight that
// begins that day in location.
func StartOfCalendarDay(value interface{}, location *time.Location) (time.Time, bool) {
	year, month, day, ok := CalendarDayOf(value)
	if !ok {
		return time.Time{}, false
	}

	return time.Date(year, month, day, 0, 0, 0, 0, location), true
}

// CalendarDayOf extracts the calendar day a value names, without converting it
// into any zone: a bare yyyy-MM-dd, an RFC 3339 string read as written, or a
// time.Time read in its own location.
func CalendarDayOf(value interface{}) (int, time.Month, int, bool) {
	var parsed time.Time

	switch typed := value.(type) {
	case time.Time:
		if typed.IsZero() {
			return 0, 0, 0, false
		}
		parsed = typed
	case string:
		if len(typed) == 0 {
			return 0, 0, 0, false
		}

		var err error
		parsed, err = time.Parse(time.DateOnly, typed)
		if err != nil {
			// time.Parse keeps the written offset as the result's location, so
			// Date() below is the day as written, not the day in UTC.
			parsed, err = time.Parse(time.RFC3339, typed)
			if err != nil {
				return 0, 0, 0, false
			}
		}
	default:
		return 0, 0, 0, false
	}

	year, month, day := parsed.Date()
	return year, month, day, true
}

// currentMonthBounds returns [first of this month, start of tomorrow) — month
// start through the end of today — where "today" is now's day in appLocation
// and the bounds are those calendar days' midnights in columnLocation,
// converted to UTC for the driver.
func currentMonthBounds(now time.Time, appLocation *time.Location, columnLocation *time.Location) (time.Time, time.Time) {
	today := now.In(appLocation)
	start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, columnLocation)
	end := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, columnLocation)

	return start.UTC(), end.UTC()
}
