package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

// The group serialization boundaries must all hydrate GroupReceiptSettings.DefaultCustomFieldIds.
// It is `gorm:"-"`, so nothing loads it implicitly, and it carries no `omitempty` — an unhydrated
// slice therefore serializes as `null` rather than being absent, which swagger does not allow (the
// property is documented as always present) and which the generated Dart deserializer has no null
// guard for. GetGroupsForUser, GetGroupById and UpdateGroupReceiptSettings already hydrate;
// GetPagedGroups and CreateGroup are covered here.

// rawDefaultCustomFieldIds pulls groupReceiptSettings.defaultCustomFieldIds out of a serialized
// group as its RAW json value.
//
// Reading it raw is the whole point: `[]uint(nil)` and `[]uint{}` are indistinguishable once
// unmarshalled back into the model, so a struct-level assertion passes whether or not the field was
// hydrated. Decoded into `any`, encoding/json turns `null` into nil and `[]` into `[]any{}` — the
// one place the difference survives.
func rawDefaultCustomFieldIds(t *testing.T, group map[string]any) any {
	t.Helper()

	settings, ok := group["groupReceiptSettings"].(map[string]any)
	if !ok {
		t.Fatalf("serialized group carries no groupReceiptSettings object: %v", group)
	}

	value, present := settings["defaultCustomFieldIds"]
	if !present {
		t.Fatalf("groupReceiptSettings omits defaultCustomFieldIds entirely: %v", settings)
	}

	return value
}

// expectDefaultCustomFieldIds asserts the raw value is a json ARRAY holding exactly wantIds. A nil
// (i.e. `null` on the wire) fails here, which is the regression being guarded.
func expectDefaultCustomFieldIds(t *testing.T, group map[string]any, wantIds []uint) {
	t.Helper()

	raw := rawDefaultCustomFieldIds(t, group)
	ids, ok := raw.([]any)
	if !ok {
		t.Errorf("defaultCustomFieldIds should serialize as an array, got %#v", raw)
		return
	}

	if len(ids) != len(wantIds) {
		t.Errorf("Expected %v, but got %v", wantIds, ids)
		return
	}

	for i, want := range wantIds {
		got, isNumber := ids[i].(float64)
		if !isNumber || uint(got) != want {
			t.Errorf("Expected %v, but got %v", wantIds, ids)
			return
		}
	}
}

func TestGetPagedGroupsHydratesDefaultCustomFieldIds(t *testing.T) {
	defer tearDownGroupTests()
	db := repositories.GetDB()

	// Two groups so both branches are covered in one response: one with a configured set, one
	// with none — the second is the case that used to serialize as null.
	configured := models.Group{Name: "aaa-paged-with-defaults"}
	if err := db.Create(&configured).Error; err != nil {
		t.Fatalf("seed configured group: %v", err)
	}
	bare := models.Group{Name: "bbb-paged-without-defaults"}
	if err := db.Create(&bare).Error; err != nil {
		t.Fatalf("seed bare group: %v", err)
	}

	settingsRepository := repositories.NewGroupReceiptSettingsRepository(nil)
	for _, groupId := range []uint{configured.ID, bare.ID} {
		if _, err := settingsRepository.CreateGroupReceiptSettings(groupId); err != nil {
			t.Fatalf("seed group receipt settings: %v", err)
		}
	}

	customFieldId := seedHandlerCustomField(t, "Paged Default Field")
	if err := db.Create(&models.GroupReceiptSettingsCustomField{
		GroupId:       configured.ID,
		CustomFieldId: customFieldId,
	}).Error; err != nil {
		t.Fatalf("seed default custom field: %v", err)
	}

	grantAppPerms(t, 1, permissions.AppAccountRead)
	grantGroupPerms(t, 1, configured.ID, permissions.GroupView)
	grantGroupPerms(t, 1, bare.ID, permissions.GroupView)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api", bytes.NewReader([]byte(
		`{"page":1,"pageSize":10,"orderBy":"name","sortDirection":"asc","filter":{"associatedGroup":"MINE"}}`,
	)))
	r = r.WithContext(context.WithValue(
		r.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}},
	))

	GetPagedGroups(w, r)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	var paged struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &paged); err != nil {
		t.Fatalf("unmarshal paged response: %v", err)
	}
	if len(paged.Data) != 2 {
		utils.PrintTestError(t, len(paged.Data), 2)
		return
	}

	// Ordered by name ascending, so the configured group comes first.
	expectDefaultCustomFieldIds(t, paged.Data[0], []uint{customFieldId})
	expectDefaultCustomFieldIds(t, paged.Data[1], []uint{})
}

func TestCreateGroupHydratesDefaultCustomFieldIds(t *testing.T) {
	defer tearDownGroupTests()

	// CreateGroup makes the group's on-disk directory after marshalling the response, and the
	// data root does not exist in the test package. Without it the handler 500s before writing
	// the body this test is about. Removed again so the package leaves no untracked directory
	// behind (api/.gitignore only covers api/data).
	if err := os.MkdirAll("data", 0o755); err != nil {
		t.Fatalf("create test data directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll("data") })

	grantAppPerms(t, 1, permissions.AppGroupsCreate)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api", bytes.NewReader([]byte(
		`{"name":"created-group","status":"ACTIVE","groupMembers":[]}`,
	)))
	r = r.WithContext(context.WithValue(
		r.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}},
	))

	CreateGroup(w, r)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Body.String(), http.StatusOK)
		return
	}

	var group map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &group); err != nil {
		t.Fatalf("unmarshal created group: %v", err)
	}

	// A brand-new group has no defaults, so the assertion is precisely that the empty case is an
	// empty ARRAY rather than null.
	expectDefaultCustomFieldIds(t, group, []uint{})
}
