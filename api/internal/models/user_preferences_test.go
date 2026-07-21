package models

import (
	"net/http/httptest"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

func TestUserPrefernces_LoadDataFromRequest(t *testing.T) {
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
	if !preferences.ShowLargeImagePreviews {
		utils.PrintTestError(t, preferences.ShowLargeImagePreviews, true)
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
