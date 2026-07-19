package models

import (
	"net/http/httptest"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

func TestReceiptProcessingSettings_LoadDataFromRequest(t *testing.T) {
	body := `{"name": "Default", "aiType": "openAi", "model": "gpt-4"}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var settings ReceiptProcessingSettings
	err := settings.LoadDataFromRequest(w, r)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if settings.Name != "Default" {
		utils.PrintTestError(t, settings.Name, "Default")
	}
	if settings.AiType != OPEN_AI {
		utils.PrintTestError(t, settings.AiType, OPEN_AI)
	}
	if settings.Model != "gpt-4" {
		utils.PrintTestError(t, settings.Model, "gpt-4")
	}
}

func TestReceiptProcessingSettings_LoadDataFromRequest_MalformedJson(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("{"))
	w := httptest.NewRecorder()

	var settings ReceiptProcessingSettings
	err := settings.LoadDataFromRequest(w, r)
	if err == nil {
		utils.PrintTestError(t, err, "an unmarshal error")
	}
}

func TestReceiptProcessingSettings_LoadDataFromRequest_BodyReadError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", errReader{})
	w := httptest.NewRecorder()

	var settings ReceiptProcessingSettings
	err := settings.LoadDataFromRequest(w, r)
	if err == nil {
		utils.PrintTestError(t, err, "a body read error")
	}
}
