package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// putSystemSettings sends a valid settings body through the real handler, with
// extra merged in (a key set to nil is left out entirely).
func putSystemSettings(t *testing.T, extra map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()

	queueConfigs := make([]map[string]interface{}, 0)
	for _, config := range models.GetAllDefaultQueueConfigurations() {
		queueConfigs = append(queueConfigs, map[string]interface{}{"name": config.Name, "priority": 1})
	}
	body := map[string]interface{}{
		"currencyDisplay":              "$",
		"currencySymbolPosition":       models.START,
		"currencyThousandthsSeparator": models.COMMA,
		"currencyDecimalSeparator":     models.DOT,
		"currencyHideDecimalPlaces":    false,
		"taskConcurrency":              1,
		"emailPollingInterval":         60,
		"taskQueueConfigurations":      queueConfigs,
	}
	for key, value := range extra {
		body[key] = value
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("PUT", "/api", strings.NewReader(string(bodyBytes)))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}}))
	w := httptest.NewRecorder()
	UpdateSystemSettings(w, r)
	return w
}

func storedTimeZone(t *testing.T) string {
	t.Helper()
	settings, err := repositories.NewSystemSettingsRepository(nil).GetSystemSettings()
	if err != nil {
		t.Fatal(err)
	}
	return settings.TimeZone
}

func TestUpdateSystemSettingsTimeZone(t *testing.T) {
	defer tearDownSystemSettingsTest()

	db := repositories.GetDB()
	db.Create(&models.SystemSettings{})
	grantAllAppPerms(t, 1)

	if got := storedTimeZone(t); got != "UTC" {
		t.Fatalf("fresh install time zone = %q, want UTC", got)
	}

	// An explicit zone persists, and the very next resolve sees it.
	w := putSystemSettings(t, map[string]interface{}{"timeZone": "America/New_York"})
	if w.Code != http.StatusOK {
		t.Fatalf("explicit zone: status %d, body %s", w.Code, w.Body.String())
	}
	if got := storedTimeZone(t); got != "America/New_York" {
		t.Errorf("stored = %q, want America/New_York", got)
	}
	if got := repositories.GetAppLocation().String(); got != "America/New_York" {
		t.Errorf("GetAppLocation() = %q right after the save, want America/New_York", got)
	}

	// A body without the key (an older client) leaves it alone, and the
	// response echoes the stored value rather than "".
	w = putSystemSettings(t, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("omitted zone: status %d, body %s", w.Code, w.Body.String())
	}
	if got := storedTimeZone(t); got != "America/New_York" {
		t.Errorf("omitted key reset the stored zone to %q", got)
	}
	var echoed models.SystemSettings
	if err := json.Unmarshal(w.Body.Bytes(), &echoed); err != nil {
		t.Fatal(err)
	}
	if echoed.TimeZone != "America/New_York" {
		t.Errorf("response timeZone = %q, want the stored America/New_York", echoed.TimeZone)
	}

	// An unknown name is a field-level 400 and changes nothing.
	w = putSystemSettings(t, map[string]interface{}{"timeZone": "Mars/Olympus_Mons"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid zone: status %d, want 400", w.Code)
	}
	if !strings.Contains(w.Body.String(), "timeZone") {
		t.Errorf("400 body should name timeZone: %s", w.Body.String())
	}
	if got := storedTimeZone(t); got != "America/New_York" {
		t.Errorf("a rejected save changed the stored zone to %q", got)
	}
}
