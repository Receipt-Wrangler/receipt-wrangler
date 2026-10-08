package utils

import "time"

// CalendarDayUTC returns midnight UTC of the calendar day value names in its
// own location — the day as written, never converted. It is the canonical
// stored form of a calendar-day value (a receipt's Date, a DATE custom field):
// the instant is fixed, so the day reads the same in every zone as long as it
// is read in UTC.
func CalendarDayUTC(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
