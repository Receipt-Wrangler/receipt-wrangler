package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

// PUT /group/{groupId}/groupReceiptSettings is gated on the group's GroupUpdate permission only.
// Configuring the group's DEFAULT CUSTOM FIELDS additionally requires the app-level
// app.custom-fields.read, so a group admin cannot attach fields whose catalog they cannot read
// (mirroring enforceReceiptCustomFieldSelection on the receipt write path).

// seedDefaultCustomFieldGroup creates a group with a receipt settings row and a member who can
// update it, optionally also holding app.custom-fields.read. Returns the group id.
func seedDefaultCustomFieldGroup(t *testing.T, userId uint, canReadCustomFields bool) uint {
	t.Helper()
	db := repositories.GetDB()

	group := models.Group{Name: "default-custom-fields-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	if _, err := repositories.NewGroupReceiptSettingsRepository(nil).CreateGroupReceiptSettings(group.ID); err != nil {
		t.Fatalf("seed group receipt settings: %v", err)
	}

	grantGroupPerms(t, userId, group.ID, permissions.GroupUpdate)
	if canReadCustomFields {
		grantAppPerms(t, userId, permissions.AppCustomFieldsRead)
	} else {
		grantAppPerms(t, userId, permissions.AppAccountRead)
	}

	return group.ID
}

func seedHandlerCustomField(t *testing.T, name string) uint {
	t.Helper()
	customField := models.CustomField{Name: name, Type: models.TEXT}
	if err := repositories.GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}
	return customField.ID
}

// callUpdateGroupReceiptSettings drives the handler with the given raw JSON body, so a test can omit
// keys entirely (which is what a pointer command field's nil case actually means on the wire).
func callUpdateGroupReceiptSettings(t *testing.T, userId uint, groupId uint, body string) *httptest.ResponseRecorder {
	t.Helper()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/api", bytes.NewReader([]byte(body)))

	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("groupId", utils.UintToString(groupId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiContext))
	r = r.WithContext(context.WithValue(
		r.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}},
	))

	UpdateGroupReceiptSettings(w, r)

	return w
}

// baseSettingsBody is a valid settings payload with neither default-custom-field key present, so
// each test can splice in only the key it exercises.
const baseSettingsBody = `"quickScanPaidByEnabled":true,"quickScanPaidByRequired":true,` +
	`"quickScanStatusEnabled":true,"quickScanStatusRequired":true,` +
	`"quickScanDefaultPaidByType":"UPLOADER","quickScanDefaultStatus":"OPEN"`

func TestUpdateGroupReceiptSettingsRejectsUnknownDefaultCustomFieldId(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)
	knownField := seedHandlerCustomField(t, "Known Field")

	// Seed a valid configuration first so the failed save can be proven to change nothing.
	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"defaultCustomFieldIds":[`+utils.UintToString(knownField)+`],"applyDefaultCustomFieldsOnIngest":true}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	w = callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"hideImages":true,"defaultCustomFieldIds":[`+utils.UintToString(knownField)+`,9999]}`)
	if w.Result().StatusCode != http.StatusBadRequest {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusBadRequest)
		return
	}

	// WriteValidatorErrorResponse writes the bare errors map, not the wrapping struct.
	validationErrors := map[string]string{}
	if err := json.Unmarshal(w.Body.Bytes(), &validationErrors); err != nil {
		utils.PrintTestError(t, err, "a validator error body")
		return
	}
	if _, ok := validationErrors["defaultCustomFieldIds"]; !ok {
		utils.PrintTestError(t, validationErrors, "a defaultCustomFieldIds error")
	}

	// Nothing was written - not the rejected set, and not the unrelated scalar riding along with it.
	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(settings.DefaultCustomFieldIds) != 1 || settings.DefaultCustomFieldIds[0] != knownField {
		utils.PrintTestError(t, settings.DefaultCustomFieldIds, []uint{knownField})
	}
	if settings.HideImages {
		utils.PrintTestError(t, settings.HideImages, false)
	}
}

func TestUpdateGroupReceiptSettingsForbidsDefaultCustomFieldsWithoutCustomFieldsRead(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, false)
	knownField := seedHandlerCustomField(t, "Known Field")

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"defaultCustomFieldIds":[`+utils.UintToString(knownField)+`]}`)
	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}

	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(settings.DefaultCustomFieldIds) != 0 {
		utils.PrintTestError(t, settings.DefaultCustomFieldIds, "empty set")
	}
}

// TestUpdateGroupReceiptSettingsForbidsIngestToggleWithoutCustomFieldsRead pins that the
// app.custom-fields.read gate covers the ingest toggle, not just the id list. The toggle decides
// whether the group's default fields ride every quick-scan and email receipt, so a caller who
// cannot read the catalog must not be able to flip it -- gating only defaultCustomFieldIds would
// let them change what those (to them invisible) fields do by sending this key alone.
func TestUpdateGroupReceiptSettingsForbidsIngestToggleWithoutCustomFieldsRead(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, false)

	// Only the toggle - no defaultCustomFieldIds. This is the request that used to slip through.
	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"applyDefaultCustomFieldsOnIngest":true}`)
	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}

	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if settings.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, settings.ApplyDefaultCustomFieldsOnIngest, false)
	}
}

// TestUpdateGroupReceiptSettingsOmittingDefaultCustomFieldKeysLeavesConfigUntouched is the reason
// both command fields are pointers: the desktop hides this whole section from an admin without
// app.custom-fields.read, so its payload omits both keys. A non-pointer bool would unmarshal as
// false and wipe the stored toggle.
func TestUpdateGroupReceiptSettingsOmittingDefaultCustomFieldKeysLeavesConfigUntouched(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)
	knownField := seedHandlerCustomField(t, "Known Field")

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"defaultCustomFieldIds":[`+utils.UintToString(knownField)+`],"applyDefaultCustomFieldsOnIngest":true}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	// A save that omits both keys - exactly what a permission-gated client sends.
	w = callUpdateGroupReceiptSettings(t, 1, groupId, `{`+baseSettingsBody+`,"hideImages":true}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	var responseSettings models.GroupReceiptSettings
	if err := json.Unmarshal(w.Body.Bytes(), &responseSettings); err != nil {
		utils.PrintTestError(t, err, "a group receipt settings body")
		return
	}
	if len(responseSettings.DefaultCustomFieldIds) != 1 || responseSettings.DefaultCustomFieldIds[0] != knownField {
		utils.PrintTestError(t, responseSettings.DefaultCustomFieldIds, []uint{knownField})
	}
	if !responseSettings.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, responseSettings.ApplyDefaultCustomFieldsOnIngest, true)
	}

	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(settings.DefaultCustomFieldIds) != 1 || settings.DefaultCustomFieldIds[0] != knownField {
		utils.PrintTestError(t, settings.DefaultCustomFieldIds, []uint{knownField})
	}
	if !settings.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, settings.ApplyDefaultCustomFieldsOnIngest, true)
	}
	if !settings.HideImages {
		utils.PrintTestError(t, settings.HideImages, true)
	}
}

// TestUpdateGroupReceiptSettingsSerializesEmptyDefaultCustomFieldsAsArray guards the wire contract on
// the PUT response itself: a null here would fail the whole payload on released Dart clients.
func TestUpdateGroupReceiptSettingsSerializesEmptyDefaultCustomFieldsAsArray(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)

	w := callUpdateGroupReceiptSettings(t, 1, groupId, `{`+baseSettingsBody+`}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		utils.PrintTestError(t, err, "a group receipt settings body")
		return
	}
	ids, ok := raw["defaultCustomFieldIds"]
	if !ok || ids == nil {
		utils.PrintTestError(t, raw["defaultCustomFieldIds"], "[] rather than null/absent")
		return
	}
	if _, ok := ids.([]any); !ok {
		utils.PrintTestError(t, ids, "a JSON array")
	}
}

// Compile-time reminder that the command's two new fields are pointers; a plain bool/slice here
// would not build.
var _ = commands.UpdateGroupReceiptSettingsCommand{
	DefaultCustomFieldIds:            &[]uint{},
	ApplyDefaultCustomFieldsOnIngest: new(bool),
}
