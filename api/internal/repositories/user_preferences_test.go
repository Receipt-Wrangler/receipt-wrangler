package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// UpdateUserPreferences copies the request onto the stored row one field at a
// time, so a field missing from that block silently never persists - no compile
// error, no failing test anywhere else. These round-trips are what catch it.
func TestUpdateUserPreferencesRoundTripsCloseChipSelectOnSelect(t *testing.T) {
	defer TruncateTestDb()
	CreateTestUser()

	repository := NewUserPreferencesRepository(nil)

	created, err := repository.GetUserPreferencesOrCreate(1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if created.CloseChipSelectOnSelect {
		utils.PrintTestError(t, created.CloseChipSelectOnSelect, false)
	}

	updated, err := repository.UpdateUserPreferences(1, models.UserPrefernces{
		CloseChipSelectOnSelect: true,
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !updated.CloseChipSelectOnSelect {
		utils.PrintTestError(t, updated.CloseChipSelectOnSelect, true)
	}

	read, err := repository.GetUserPreferencesOrCreate(1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !read.CloseChipSelectOnSelect {
		utils.PrintTestError(t, read.CloseChipSelectOnSelect, true)
	}
}

// GORM's struct-form Updates skips zero values, so a boolean that is turned
// back off is the case that regresses. The repository writes with Select("*"),
// which is what makes the false stick.
func TestUpdateUserPreferencesTurnsCloseChipSelectOnSelectBackOff(t *testing.T) {
	defer TruncateTestDb()
	CreateTestUser()

	repository := NewUserPreferencesRepository(nil)

	_, err := repository.GetUserPreferencesOrCreate(1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	_, err = repository.UpdateUserPreferences(1, models.UserPrefernces{
		CloseChipSelectOnSelect: true,
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	_, err = repository.UpdateUserPreferences(1, models.UserPrefernces{
		CloseChipSelectOnSelect: false,
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	read, err := repository.GetUserPreferencesOrCreate(1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if read.CloseChipSelectOnSelect {
		utils.PrintTestError(t, read.CloseChipSelectOnSelect, false)
	}
}
