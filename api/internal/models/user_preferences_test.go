package models

import (
	"net/http/httptest"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

func TestUserPrefernces_LoadDataFromRequest(t *testing.T) {
	body := `{"userId": 5, "quickScanDefaultStatus": "OPEN", "closeChipSelectOnSelect": true}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var preferences UserPrefernces
	err := preferences.LoadDataFromRequest(w, r)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if preferences.UserId != 5 {
		utils.PrintTestError(t, preferences.UserId, uint(5))
	}
	if preferences.QuickScanDefaultStatus != OPEN {
		utils.PrintTestError(t, preferences.QuickScanDefaultStatus, OPEN)
	}
	if !preferences.CloseChipSelectOnSelect {
		utils.PrintTestError(t, preferences.CloseChipSelectOnSelect, true)
	}
}

// A body that omits closeChipSelectOnSelect decodes it as false, which is the
// default behavior (the option list stays open). Pinned because the repository
// writes the field unconditionally, so "omitted" and "explicitly false" are the
// same request as far as the stored value is concerned.
func TestUserPrefernces_LoadDataFromRequestDefaultsCloseChipSelectOnSelect(t *testing.T) {
	body := `{"userId": 5, "quickScanDefaultStatus": "OPEN"}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var preferences UserPrefernces
	err := preferences.LoadDataFromRequest(w, r)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if preferences.CloseChipSelectOnSelect {
		utils.PrintTestError(t, preferences.CloseChipSelectOnSelect, false)
	}
}

// showLargeImagePreviews was removed from this model, but already-released
// mobile builds still send it. LoadDataFromRequest uses encoding/json's default
// decoder and the API sets DisallowUnknownFields nowhere, so the unknown key is
// ignored rather than rejected and the fields that remain still decode. Keep
// the removed field in this payload - it is the whole point of the test.
func TestUserPrefernces_LoadDataFromRequestIgnoresRemovedFields(t *testing.T) {
	body := `{"userId": 5, "showLargeImagePreviews": true, "quickScanDefaultStatus": "OPEN"}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var preferences UserPrefernces
	err := preferences.LoadDataFromRequest(w, r)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if preferences.UserId != 5 {
		utils.PrintTestError(t, preferences.UserId, uint(5))
	}
	if preferences.QuickScanDefaultStatus != OPEN {
		utils.PrintTestError(t, preferences.QuickScanDefaultStatus, OPEN)
	}
}

func TestUserPrefernces_LoadDataFromRequest_MalformedJson(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("{"))
	w := httptest.NewRecorder()

	var preferences UserPrefernces
	err := preferences.LoadDataFromRequest(w, r)
	if err == nil {
		utils.PrintTestError(t, err, "an unmarshal error")
	}
}

func TestUserPrefernces_LoadDataFromRequest_BodyReadError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", errReader{})
	w := httptest.NewRecorder()

	var preferences UserPrefernces
	err := preferences.LoadDataFromRequest(w, r)
	if err == nil {
		utils.PrintTestError(t, err, "a body read error")
	}
}
