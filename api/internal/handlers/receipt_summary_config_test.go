package handlers

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// The receipt summary's configuration rides the same PUT as the default custom fields, and reuses
// its helpers (seedDefaultCustomFieldGroup, callUpdateGroupReceiptSettings, baseSettingsBody) from
// group_default_custom_fields_test.go.
//
// The gate is deliberately asymmetric: the CURRENCY field selection needs app.custom-fields.read
// because it names catalog entries, but the master toggle and the status breakdown do not - gating
// those would lock an admin without that permission out of the whole feature.

func seedCurrencyHandlerCustomField(t *testing.T, name string) uint {
	t.Helper()
	customField := models.CustomField{Name: name, Type: models.CURRENCY}
	if err := repositories.GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed currency custom field: %v", err)
	}
	return customField.ID
}

func settingsValidationErrors(t *testing.T, body []byte) map[string]string {
	t.Helper()
	validationErrors := map[string]string{}
	if err := json.Unmarshal(body, &validationErrors); err != nil {
		t.Fatalf("parse validator error body: %v", err)
	}
	return validationErrors
}

func TestUpdateGroupReceiptSettingsRoundTripsSummaryConfig(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)
	hst := seedCurrencyHandlerCustomField(t, "HST")

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryEnabled":true,"receiptSummaryStatuses":["RESOLVED","OPEN"],`+
			`"receiptSummaryCustomFieldIds":[`+utils.UintToString(hst)+`]}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !settings.ReceiptSummaryEnabled {
		utils.PrintTestError(t, settings.ReceiptSummaryEnabled, true)
	}
	// Canonical enum order, not the submitted order.
	if len(settings.ReceiptSummaryStatuses) != 2 ||
		settings.ReceiptSummaryStatuses[0] != models.OPEN ||
		settings.ReceiptSummaryStatuses[1] != models.RESOLVED {
		utils.PrintTestError(t, settings.ReceiptSummaryStatuses, []models.ReceiptStatus{models.OPEN, models.RESOLVED})
	}
	if len(settings.ReceiptSummaryCustomFieldIds) != 1 || settings.ReceiptSummaryCustomFieldIds[0] != hst {
		utils.PrintTestError(t, settings.ReceiptSummaryCustomFieldIds, []uint{hst})
	}
}

func TestUpdateGroupReceiptSettingsRejectsNonCurrencySummaryCustomField(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)
	textField := seedHandlerCustomField(t, "Notes")

	// A TEXT field has no CurrencyValue, so it would total 0.00 on every row forever - which reads
	// as a group with no spend rather than as a misconfiguration.
	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryCustomFieldIds":[`+utils.UintToString(textField)+`]}`)
	if w.Result().StatusCode != http.StatusBadRequest {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusBadRequest)
		return
	}
	if _, ok := settingsValidationErrors(t, w.Body.Bytes())["receiptSummaryCustomFieldIds"]; !ok {
		utils.PrintTestError(t, w.Body.String(), "a receiptSummaryCustomFieldIds error")
	}
}

func TestUpdateGroupReceiptSettingsRejectsUnknownSummaryCustomFieldId(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryCustomFieldIds":[9999]}`)
	if w.Result().StatusCode != http.StatusBadRequest {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusBadRequest)
		return
	}
	if _, ok := settingsValidationErrors(t, w.Body.Bytes())["receiptSummaryCustomFieldIds"]; !ok {
		utils.PrintTestError(t, w.Body.String(), "a receiptSummaryCustomFieldIds error")
	}
}

func TestUpdateGroupReceiptSettingsRejectsInvalidSummaryStatus(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)

	// Without the membership check this reaches the DB layer and surfaces as a generic 500 - the
	// gap api/CLAUDE.md documents for the quick-scan default status.
	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryStatuses":["NOT_A_STATUS"]}`)
	if w.Result().StatusCode != http.StatusBadRequest {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusBadRequest)
		return
	}
	if _, ok := settingsValidationErrors(t, w.Body.Bytes())["receiptSummaryStatuses"]; !ok {
		utils.PrintTestError(t, w.Body.String(), "a receiptSummaryStatuses error")
	}
}

func TestUpdateGroupReceiptSettingsForbidsSummaryCustomFieldsWithoutPermission(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, false)
	hst := seedCurrencyHandlerCustomField(t, "HST")

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryCustomFieldIds":[`+utils.UintToString(hst)+`]}`)
	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}
}

// TestUpdateGroupReceiptSettingsAllowsSummaryToggleWithoutCustomFieldPermission is the other half of
// the asymmetry, and the one worth pinning: an admin without app.custom-fields.read must still be
// able to turn the summary on and pick which statuses break out.
func TestUpdateGroupReceiptSettingsAllowsSummaryToggleWithoutCustomFieldPermission(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, false)

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryEnabled":true,"receiptSummaryStatuses":["OPEN"]}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !settings.ReceiptSummaryEnabled {
		utils.PrintTestError(t, settings.ReceiptSummaryEnabled, true)
	}
	if len(settings.ReceiptSummaryStatuses) != 1 || settings.ReceiptSummaryStatuses[0] != models.OPEN {
		utils.PrintTestError(t, settings.ReceiptSummaryStatuses, []models.ReceiptStatus{models.OPEN})
	}
}

// TestUpdateGroupReceiptSettingsOmittingSummaryKeysLeavesConfigIntact is what the pointer command
// fields buy: the desktop omits the custom field key for an admin without the catalog permission,
// and no other client sends any of these keys at all.
func TestUpdateGroupReceiptSettingsOmittingSummaryKeysLeavesConfigIntact(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDefaultCustomFieldGroup(t, 1, true)
	hst := seedCurrencyHandlerCustomField(t, "HST")

	w := callUpdateGroupReceiptSettings(t, 1, groupId,
		`{`+baseSettingsBody+`,"receiptSummaryEnabled":true,"receiptSummaryStatuses":["OPEN"],`+
			`"receiptSummaryCustomFieldIds":[`+utils.UintToString(hst)+`]}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	// A save carrying none of the three summary keys.
	w = callUpdateGroupReceiptSettings(t, 1, groupId, `{`+baseSettingsBody+`,"hideImages":true}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !settings.ReceiptSummaryEnabled {
		utils.PrintTestError(t, settings.ReceiptSummaryEnabled, "still enabled")
	}
	if len(settings.ReceiptSummaryStatuses) != 1 {
		utils.PrintTestError(t, settings.ReceiptSummaryStatuses, "still one status")
	}
	if len(settings.ReceiptSummaryCustomFieldIds) != 1 {
		utils.PrintTestError(t, settings.ReceiptSummaryCustomFieldIds, "still one field")
	}
}
