package repositories

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/logging"
)

// loadedLocations caches parsed zones by name. It caches the zone DATABASE
// lookup only, never the setting: every GetAppLocation call still reads the
// stored name, so changing the setting takes effect on the very next call.
var loadedLocations sync.Map // string -> *time.Location

// reportedInvalidZones remembers which bad names have already been logged, so
// a stored name the runtime cannot load is reported once rather than on every
// request that resolves the zone.
var reportedInvalidZones sync.Map // string -> struct{}

// GetAppLocation returns the app time zone configured in System Settings — the
// zone every calendar boundary is worked out in: which day an instant such as
// created_at falls on, "today", "this month", report periods and
// {{generatedAt}}. A receipt's own Date is a calendar day and is read in UTC
// instead; see BuildDayFilterQuery.
//
// It is read live with no cache, like services.IsMcpEnabled, so a change made in
// System Settings applies to the next request without a restart. An unreadable
// settings row, an empty name or a name the runtime cannot load all resolve to
// UTC, which is also the setting's default.
func GetAppLocation() *time.Location {
	return appLocation(nil)
}

// AppLocation is GetAppLocation read through the repository's transaction, so a
// filter built inside a transaction never opens a second connection for it.
func (repository BaseRepository) AppLocation() *time.Location {
	return appLocation(repository.TX)
}

func appLocation(tx *gorm.DB) *time.Location {
	settings, err := NewSystemSettingsRepository(tx).GetSystemSettings()
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Could not read the app time zone, using UTC: "+err.Error())
		return time.UTC
	}

	return ResolveAppLocation(settings.TimeZone)
}

// ResolveAppLocation turns a stored time zone name into a location, falling back
// to UTC for an empty or unknown name. "Local" is refused as well: it would make
// the answer depend on the server process's zone, which is exactly what the
// setting replaces.
func ResolveAppLocation(name string) *time.Location {
	if name == "" || name == "UTC" {
		return time.UTC
	}

	if cached, ok := loadedLocations.Load(name); ok {
		return cached.(*time.Location)
	}

	if name != "Local" {
		if location, err := time.LoadLocation(name); err == nil {
			loadedLocations.Store(name, location)
			return location
		}
	}

	if _, alreadyReported := reportedInvalidZones.LoadOrStore(name, struct{}{}); !alreadyReported {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Unknown app time zone "+name+", using UTC")
	}

	return time.UTC
}
