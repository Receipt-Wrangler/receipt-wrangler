package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
)

func TestFormatEmailContent_BothPresent(t *testing.T) {
	result := FormatEmailContent("OCR extracted text", true, "Email body content")
	expected := "Image Data:\nOCR extracted text\n\nEmail Body:\nEmail body content"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatEmailContent_ImageDataOnly(t *testing.T) {
	result := FormatEmailContent("OCR extracted text", true, "")
	expected := "Image Data:\nOCR extracted text\n\nEmail Body:\nNo email body found"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatEmailContent_EmailBodyOnly(t *testing.T) {
	result := FormatEmailContent("", false, "Receipt from Amazon: $45.00")
	expected := "Image Data:\nNo attachments found\n\nEmail Body:\nReceipt from Amazon: $45.00"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatEmailContent_NeitherPresent(t *testing.T) {
	result := FormatEmailContent("", false, "")
	expected := "Image Data:\nNo attachments found\n\nEmail Body:\nNo email body found"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatEmailContent_HasImageButEmptyOcrText(t *testing.T) {
	// When hasImage is true but ocrText is empty (e.g., vision model path),
	// we should NOT show "No attachments found" since there IS an image.
	result := FormatEmailContent("", true, "Some body")
	expected := "Image Data:\n\n\nEmail Body:\nSome body"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestCombineOcrResults_Empty(t *testing.T) {
	text, cmd, err := combineOcrResults(nil)
	if err != nil {
		t.Errorf("expected no error for empty input, got: %v", err)
	}
	if text != "" {
		t.Errorf("expected empty text, got: %q", text)
	}
	if cmd.Status != "" {
		t.Errorf("expected zero-value command, got status %q", cmd.Status)
	}
}

func TestCombineOcrResults_SingleSuccess(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Second)
	in := []ocrImageResult{
		{
			Text: "receipt total $5",
			Command: commands.UpsertSystemTaskCommand{
				Status:            models.SYSTEM_TASK_SUCCEEDED,
				StartedAt:         start,
				EndedAt:           &end,
				ResultDescription: "receipt total $5",
			},
		},
	}

	text, cmd, err := combineOcrResults(in)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if text != "receipt total $5" {
		t.Errorf("expected single-image text passthrough, got: %q", text)
	}
	if cmd.Status != models.SYSTEM_TASK_SUCCEEDED {
		t.Errorf("expected SUCCEEDED, got: %s", cmd.Status)
	}
	// Single-image success: ResultDescription gets overwritten with the
	// combined text, which equals the original text — idempotent.
	if cmd.ResultDescription != "receipt total $5" {
		t.Errorf("expected ResultDescription %q, got: %q", "receipt total $5", cmd.ResultDescription)
	}
}

func TestCombineOcrResults_SingleFailurePreservesErrorDescription(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Millisecond)
	failErr := errors.New("ocr blew up")
	in := []ocrImageResult{
		{
			Text: "",
			Err:  failErr,
			Command: commands.UpsertSystemTaskCommand{
				Status:            models.SYSTEM_TASK_FAILED,
				StartedAt:         start,
				EndedAt:           &end,
				ResultDescription: "ocr blew up",
			},
		},
	}

	text, cmd, err := combineOcrResults(in)

	if !errors.Is(err, failErr) {
		t.Errorf("expected the original failure error, got: %v", err)
	}
	if text != "" {
		t.Errorf("expected empty combined text on sole failure, got: %q", text)
	}
	if cmd.Status != models.SYSTEM_TASK_FAILED {
		t.Errorf("expected FAILED status, got: %s", cmd.Status)
	}
	// When every image failed there's no combined text to write — the
	// error-message description from OcrService must survive.
	if cmd.ResultDescription != "ocr blew up" {
		t.Errorf("expected failure description preserved, got: %q", cmd.ResultDescription)
	}
}

func TestCombineOcrResults_MultiSuccessJoinsWithSeparator(t *testing.T) {
	start := time.Now()
	mid := start.Add(time.Second)
	end := start.Add(2 * time.Second)
	in := []ocrImageResult{
		{
			Text: "first image",
			Command: commands.UpsertSystemTaskCommand{
				Status:            models.SYSTEM_TASK_SUCCEEDED,
				StartedAt:         start,
				EndedAt:           &mid,
				ResultDescription: "first image",
			},
		},
		{
			Text: "second image",
			Command: commands.UpsertSystemTaskCommand{
				Status:            models.SYSTEM_TASK_SUCCEEDED,
				StartedAt:         mid,
				EndedAt:           &end,
				ResultDescription: "second image",
			},
		},
	}

	text, cmd, err := combineOcrResults(in)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	expectedText := "first image" + ocrImageSeparator + "second image"
	if text != expectedText {
		t.Errorf("expected text %q, got %q", expectedText, text)
	}
	if cmd.Status != models.SYSTEM_TASK_SUCCEEDED {
		t.Errorf("expected SUCCEEDED, got: %s", cmd.Status)
	}
	if cmd.StartedAt != start {
		t.Errorf("expected aggregated start to be the first start, got: %v", cmd.StartedAt)
	}
	if cmd.EndedAt == nil || *cmd.EndedAt != end {
		t.Errorf("expected aggregated end to be the last end (%v), got: %v", end, cmd.EndedAt)
	}
	if cmd.ResultDescription != expectedText {
		t.Errorf("multi-image run should overwrite ResultDescription with combined text, got: %q", cmd.ResultDescription)
	}
}

func TestCombineOcrResults_MultiSecondFails(t *testing.T) {
	start := time.Now()
	mid := start.Add(time.Second)
	end := start.Add(2 * time.Second)
	failErr := errors.New("second image failed")
	in := []ocrImageResult{
		{
			Text: "first ok",
			Command: commands.UpsertSystemTaskCommand{
				Status:            models.SYSTEM_TASK_SUCCEEDED,
				StartedAt:         start,
				EndedAt:           &mid,
				ResultDescription: "first ok",
			},
		},
		{
			Text: "",
			Err:  failErr,
			Command: commands.UpsertSystemTaskCommand{
				Status:    models.SYSTEM_TASK_FAILED,
				StartedAt: mid,
				EndedAt:   &end,
			},
		},
	}

	text, cmd, err := combineOcrResults(in)

	if !errors.Is(err, failErr) {
		t.Errorf("expected the second image's error, got: %v", err)
	}
	// Aggregated text should still contain the first image's content; failed images contribute nothing.
	if text != "first ok" {
		t.Errorf("expected only the successful image's text, got: %q", text)
	}
	if cmd.Status != models.SYSTEM_TASK_FAILED {
		t.Errorf("expected FAILED when any image fails, got: %s", cmd.Status)
	}
	if cmd.EndedAt == nil || *cmd.EndedAt != end {
		t.Errorf("expected aggregated end to be the failing image's end, got: %v", cmd.EndedAt)
	}
}

func TestEncodeImageForAi_UnsupportedAiTypeReturnsError(t *testing.T) {
	service := ReceiptProcessingService{}
	settings := models.ReceiptProcessingSettings{
		AiType: models.AiClientType("totally-bogus-provider"),
	}

	encoded, err := service.encodeImageForAi("/does/not/matter", settings)
	if err == nil {
		t.Fatal("expected error for unsupported AI type, got nil")
	}
	if encoded != "" {
		t.Errorf("expected empty encoded string on error, got: %q", encoded)
	}
	if !strings.Contains(err.Error(), "unsupported AI type") {
		t.Errorf("expected error to mention 'unsupported AI type', got: %v", err)
	}
}

// ---------- cleanResponse ----------

func TestCleanResponse_StripsJsonMarkers(t *testing.T) {
	service := ReceiptProcessingService{}
	got := service.cleanResponse("```json\n{\"a\":1}\n```")
	if strings.Contains(got, "```") {
		t.Errorf("expected backticks stripped, got: %q", got)
	}
	if !strings.Contains(got, `{"a":1}`) {
		t.Errorf("expected payload preserved, got: %q", got)
	}
}

func TestCleanResponse_NoMarkersPassThrough(t *testing.T) {
	service := ReceiptProcessingService{}
	got := service.cleanResponse(`{"a":1}`)
	if got != `{"a":1}` {
		t.Errorf("expected %q, got %q", `{"a":1}`, got)
	}
}

func TestCleanResponse_StripsTrailingCommas(t *testing.T) {
	service := ReceiptProcessingService{}
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"object", `{"a":1,}`, `{"a":1}`},
		{"array", `[1,2,3,]`, `[1,2,3]`},
		{"nested", `{"x":[1,2,],"y":{"z":1,},}`, `{"x":[1,2],"y":{"z":1}}`},
		{"whitespace before close", "{\"a\":1,\n}", "{\"a\":1\n}"},
		{"valid json untouched", `{"a":1,"b":2}`, `{"a":1,"b":2}`},
		{"comma inside string preserved", `{"name":"Beans, baked","amount":1,}`, `{"name":"Beans, baked","amount":1}`},
		// String-aware safety: a ",}" or ",]" sequence INSIDE a quoted value must
		// not be touched; only the genuine trailing comma after the value is removed.
		{"comma-brace inside string preserved", `{"note":"foo,}","amt":1,}`, `{"note":"foo,}","amt":1}`},
		{"comma-bracket inside string preserved", `{"note":"foo,]","amt":1,}`, `{"note":"foo,]","amt":1}`},
		{"escaped quote inside string", `{"note":"he said \"hi,\"","amt":1,}`, `{"note":"he said \"hi,\"","amt":1}`},
		// cleanResponse chains fence-stripping then comma-stripping; cover the
		// interaction directly so a regression in ordering is caught here.
		{"fenced json with trailing comma", "```json\n{\"a\":1,}\n```", "\n{\"a\":1}\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.cleanResponse(tt.in)
			if got != tt.want {
				t.Errorf("cleanResponse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestCleanResponse_TrailingCommaOutputIsParseable guards the actual purpose
// of the fix: AI output with trailing commas must unmarshal cleanly into the
// real production type used at the parse site (commands.UpsertReceiptCommand),
// including a populated items array. Using the production type rather than a
// hand-picked struct catches drift in the actual integration contract.
func TestCleanResponse_TrailingCommaOutputIsParseable(t *testing.T) {
	service := ReceiptProcessingService{}
	// Note: UpsertReceiptCommand maps items via the "receiptItems" JSON key.
	raw := `{
  "name": "BILLA",
  "amount": 2306.86,
  "date": "2026-06-03T00:00:00Z",
  "receiptItems": [
    { "name": "COTTAGE FIT", "amount": 197.4, },
    { "name": "Beans, baked", "amount": 51.8, },
  ],
  "categories": [],
  "tags": [],
}`
	cleaned := service.cleanResponse(raw)

	var receipt commands.UpsertReceiptCommand
	if err := json.Unmarshal([]byte(cleaned), &receipt); err != nil {
		t.Fatalf("expected cleaned response to parse, got error: %v\ncleaned: %s", err, cleaned)
	}
	if receipt.Name != "BILLA" {
		t.Errorf("expected name BILLA, got %q", receipt.Name)
	}
	if len(receipt.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(receipt.Items))
	}
	if receipt.Items[1].Name != "Beans, baked" {
		t.Errorf("expected comma-containing item name preserved, got %q", receipt.Items[1].Name)
	}
}

// ---------- Fixture helpers for the DB/orchestration tests ----------

// seedReceiptProcessingFixtures creates the minimum viable graph to exercise
// processing: a Prompt, a ReceiptProcessingSettings pointed at the given
// mock-server URL with AiType=OLLAMA, and a matching
// SystemReceiptProcessingSettings so NewSystemReceiptProcessingService works
// when invoked with an empty groupId.
func seedReceiptProcessingFixtures(t *testing.T, url string) models.ReceiptProcessingSettings {
	t.Helper()
	db := repositories.GetDB()

	prompt := models.Prompt{Name: "p", Prompt: "Receipt: @ocrText"}
	if err := db.Create(&prompt).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}

	settings := models.ReceiptProcessingSettings{
		Name:          "test-rps",
		AiType:        models.OLLAMA,
		Url:           url,
		Model:         "test-model",
		IsVisionModel: false,
		PromptId:      prompt.ID,
	}
	if err := db.Create(&settings).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}

	// Link as the active system settings so
	// NewSystemReceiptProcessingService picks it up without a groupId.
	sysSettings := models.SystemSettings{
		BaseModel:                   models.BaseModel{ID: 1},
		ReceiptProcessingSettingsId: &settings.ID,
	}
	if err := db.Create(&sysSettings).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}

	return settings
}

func ollamaReceiptJson(name, amount string) string {
	content := fmt.Sprintf(`{"name":"%s","amount":"%s"}`, name, amount)
	contentEscaped, _ := json.Marshal(content)
	return `{"model":"test","created_at":"2024-01-01T00:00:00Z","message":{"role":"assistant","content":` + string(contentEscaped) + `},"done":true}`
}

func setupImagePathFixture(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	path := filepath.Join(tmp, "img.jpg")

	jpg, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(path, jpg, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// ---------- Constructors ----------

func TestNewReceiptProcessingService_InvalidId(t *testing.T) {
	defer repositories.TruncateTestDb()

	_, err := NewReceiptProcessingService(nil, "999999", "")
	if err == nil {
		t.Fatal("expected not-found error for invalid settings id")
	}
}

func TestNewReceiptProcessingService_HappyPath(t *testing.T) {
	defer repositories.TruncateTestDb()

	settings := seedReceiptProcessingFixtures(t, "http://example.invalid")
	service, err := NewReceiptProcessingService(nil, utils.UintToString(settings.ID), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service.ReceiptProcessingSettings.ID != settings.ID {
		t.Errorf("expected settings ID %d, got %d", settings.ID, service.ReceiptProcessingSettings.ID)
	}
}

func TestNewReceiptProcessingService_WithFallback(t *testing.T) {
	defer repositories.TruncateTestDb()

	primary := seedReceiptProcessingFixtures(t, "http://primary.invalid")
	db := repositories.GetDB()
	fallbackPrompt := models.Prompt{Name: "p-fallback", Prompt: "fallback"}
	if err := db.Create(&fallbackPrompt).Error; err != nil {
		t.Fatalf("seed fallback prompt: %v", err)
	}
	fallback := models.ReceiptProcessingSettings{
		Name:     "fallback-rps",
		AiType:   models.OLLAMA,
		Url:      "http://fallback.invalid",
		Model:    "m",
		PromptId: fallbackPrompt.ID,
	}
	if err := db.Create(&fallback).Error; err != nil {
		t.Fatalf("seed fallback: %v", err)
	}

	service, err := NewReceiptProcessingService(nil, utils.UintToString(primary.ID), utils.UintToString(fallback.ID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service.FallbackReceiptProcessingSettings.ID != fallback.ID {
		t.Errorf("expected fallback ID %d, got %d", fallback.ID, service.FallbackReceiptProcessingSettings.ID)
	}
}

func TestNewSystemReceiptProcessingService_NoGroupId(t *testing.T) {
	defer repositories.TruncateTestDb()

	seedReceiptProcessingFixtures(t, "http://sys.invalid")
	service, err := NewSystemReceiptProcessingService(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service.ReceiptProcessingSettings.ID == 0 {
		t.Error("expected loaded ReceiptProcessingSettings to have nonzero ID")
	}
}

// ---------- buildPrompt ----------

func TestBuildPrompt_InterpolatesVariables(t *testing.T) {
	defer repositories.TruncateTestDb()

	db := repositories.GetDB()
	prompt := models.Prompt{
		Name:   "bp",
		Prompt: "categories=@categories tags=@tags ocr=@ocrText year=@currentYear",
	}
	if err := db.Create(&prompt).Error; err != nil {
		t.Fatalf("seed prompt: %v", err)
	}

	service := ReceiptProcessingService{
		ReceiptProcessingSettings: models.ReceiptProcessingSettings{PromptId: prompt.ID},
	}
	realPrompt, cmd, err := service.buildPrompt(service.ReceiptProcessingSettings, "some ocr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(realPrompt, "ocr=some ocr") {
		t.Errorf("expected ocr substituted, got: %q", realPrompt)
	}
	if !strings.Contains(realPrompt, "year=") {
		t.Errorf("expected current-year substitution, got: %q", realPrompt)
	}
	if cmd.Status != models.SYSTEM_TASK_SUCCEEDED {
		t.Errorf("expected SUCCEEDED, got: %s", cmd.Status)
	}
}

func TestBuildPrompt_InvalidPromptId(t *testing.T) {
	defer repositories.TruncateTestDb()

	service := ReceiptProcessingService{
		ReceiptProcessingSettings: models.ReceiptProcessingSettings{PromptId: 999999},
	}
	_, cmd, err := service.buildPrompt(service.ReceiptProcessingSettings, "")
	if err == nil {
		t.Fatal("expected error for missing prompt")
	}
	if cmd.Status != models.SYSTEM_TASK_FAILED {
		t.Errorf("expected FAILED, got: %s", cmd.Status)
	}
}

// ---------- Image encoding helpers ----------

func TestEncodeImageForAi_OpenAi(t *testing.T) {
	path := setupImagePathFixture(t)
	service := ReceiptProcessingService{}

	got, err := service.encodeImageForAi(path, models.ReceiptProcessingSettings{AiType: models.OPEN_AI_NEW})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(got, "data:image/") {
		t.Errorf("expected data URI prefix, got: %q", got[:min(25, len(got))])
	}
}

func TestEncodeImageForAi_Ollama(t *testing.T) {
	path := setupImagePathFixture(t)
	service := ReceiptProcessingService{}

	got, err := service.encodeImageForAi(path, models.ReceiptProcessingSettings{AiType: models.OLLAMA})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty base64 string")
	}
	if strings.HasPrefix(got, "data:") {
		t.Errorf("ollama path should return raw base64, got data URI: %q", got[:min(25, len(got))])
	}
}

func TestEncodeImageForAi_Gemini(t *testing.T) {
	path := setupImagePathFixture(t)
	service := ReceiptProcessingService{}

	got, err := service.encodeImageForAi(path, models.ReceiptProcessingSettings{AiType: models.GEMINI_NEW})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty base64 string")
	}
}

func TestGetOpenAiBase64Image_ReadFileError(t *testing.T) {
	service := ReceiptProcessingService{}
	_, err := service.getOpenAiBase64Image("/definitely/not/a/real/path/nope.jpg")
	if err == nil {
		t.Fatal("expected read-file error")
	}
}

// ---------- Full orchestration: text-only happy path ----------

func TestReadReceiptText_HappyPath(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("Coffee Shop", "5.25"))
	settings := seedReceiptProcessingFixtures(t, server.URL)

	service, err := NewReceiptProcessingService(nil, utils.UintToString(settings.ID), "")
	if err != nil {
		t.Fatalf("constructor: %v", err)
	}

	receipt, metadata, err := service.ReadReceiptText("email body goes here")
	if err != nil {
		t.Fatalf("ReadReceiptText: %v", err)
	}
	if receipt.Name != "Coffee Shop" {
		t.Errorf("expected name 'Coffee Shop', got: %q", receipt.Name)
	}
	if !metadata.DidReceiptProcessingSettingsSucceed {
		t.Error("expected DidReceiptProcessingSettingsSucceed=true")
	}
	if metadata.ReceiptProcessingSettingsIdRan != settings.ID {
		t.Errorf("expected settings id %d, got %d", settings.ID, metadata.ReceiptProcessingSettingsIdRan)
	}
}

func TestReadReceiptImagesWithEmailBody_HappyPath_Vision(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("Grocery", "42"))
	settings := seedReceiptProcessingFixtures(t, server.URL)
	// Switch to a vision model so we skip the OCR path and base64-encode images.
	repositories.GetDB().Model(&models.ReceiptProcessingSettings{}).Where("id = ?", settings.ID).Update("is_vision_model", true)

	service, err := NewReceiptProcessingService(nil, utils.UintToString(settings.ID), "")
	if err != nil {
		t.Fatalf("constructor: %v", err)
	}
	// Reload settings with the updated IsVisionModel flag so processImages sees it.
	repositories.GetDB().Where("id = ?", settings.ID).First(&service.ReceiptProcessingSettings)

	imgPath := setupImagePathFixture(t)
	receipt, metadata, err := service.ReadReceiptImagesWithEmailBody([]string{imgPath}, "body", false)
	if err != nil {
		t.Fatalf("ReadReceiptImagesWithEmailBody: %v", err)
	}
	if receipt.Name != "Grocery" {
		t.Errorf("expected name 'Grocery', got: %q", receipt.Name)
	}
	if !metadata.DidReceiptProcessingSettingsSucceed {
		t.Error("expected DidReceiptProcessingSettingsSucceed=true")
	}
}

func TestReadReceipt_FallbackUsedWhenPrimaryFails(t *testing.T) {
	defer repositories.TruncateTestDb()

	primaryServer, _ := newMockOllamaServerForService(t, http.StatusOK, "")
	primaryUrl := primaryServer.URL
	primaryServer.Close()

	fallbackServer, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("From Fallback", "7"))

	primary := seedReceiptProcessingFixtures(t, primaryUrl)
	db := repositories.GetDB()
	fallbackPrompt := models.Prompt{Name: "fallback-p", Prompt: "fallback prompt"}
	if err := db.Create(&fallbackPrompt).Error; err != nil {
		t.Fatalf("seed fallback prompt: %v", err)
	}
	fallback := models.ReceiptProcessingSettings{
		Name:     "fallback-rps",
		AiType:   models.OLLAMA,
		Url:      fallbackServer.URL,
		Model:    "m",
		PromptId: fallbackPrompt.ID,
	}
	if err := db.Create(&fallback).Error; err != nil {
		t.Fatalf("seed fallback: %v", err)
	}

	service, err := NewReceiptProcessingService(nil, utils.UintToString(primary.ID), utils.UintToString(fallback.ID))
	if err != nil {
		t.Fatalf("constructor: %v", err)
	}

	receipt, metadata, err := service.ReadReceiptText("body")
	if err != nil {
		t.Fatalf("ReadReceiptText: %v", err)
	}
	if receipt.Name != "From Fallback" {
		t.Errorf("expected 'From Fallback', got: %q", receipt.Name)
	}
	if metadata.DidReceiptProcessingSettingsSucceed {
		t.Error("expected primary to NOT succeed")
	}
	if !metadata.DidFallbackReceiptProcessingSettingsSucceed {
		t.Error("expected fallback to succeed")
	}
	if metadata.FallbackReceiptProcessingSettingsIdRan != fallback.ID {
		t.Errorf("expected fallback id %d, got %d", fallback.ID, metadata.FallbackReceiptProcessingSettingsIdRan)
	}
}

// ReadReceiptImage and ReadReceiptImageWithEmailBody are thin one-line
// wrappers around readReceipt; cover them so the wrappers aren't 0%.
func TestReadReceiptImage_SingleImageWrapper(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("Wrapper", "1"))
	settings := seedReceiptProcessingFixtures(t, server.URL)
	repositories.GetDB().Model(&models.ReceiptProcessingSettings{}).Where("id = ?", settings.ID).Update("is_vision_model", true)

	service, err := NewReceiptProcessingService(nil, utils.UintToString(settings.ID), "")
	if err != nil {
		t.Fatalf("constructor: %v", err)
	}
	repositories.GetDB().Where("id = ?", settings.ID).First(&service.ReceiptProcessingSettings)

	imgPath := setupImagePathFixture(t)
	receipt, _, err := service.ReadReceiptImage(imgPath)
	if err != nil {
		t.Fatalf("ReadReceiptImage: %v", err)
	}
	if receipt.Name != "Wrapper" {
		t.Errorf("expected 'Wrapper', got %q", receipt.Name)
	}
}

func TestReadReceiptImageWithEmailBody_Wrapper(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("WrapperBody", "2"))
	settings := seedReceiptProcessingFixtures(t, server.URL)
	repositories.GetDB().Model(&models.ReceiptProcessingSettings{}).Where("id = ?", settings.ID).Update("is_vision_model", true)

	service, err := NewReceiptProcessingService(nil, utils.UintToString(settings.ID), "")
	if err != nil {
		t.Fatalf("constructor: %v", err)
	}
	repositories.GetDB().Where("id = ?", settings.ID).First(&service.ReceiptProcessingSettings)

	imgPath := setupImagePathFixture(t)
	receipt, _, err := service.ReadReceiptImageWithEmailBody(imgPath, "email body")
	if err != nil {
		t.Fatalf("ReadReceiptImageWithEmailBody: %v", err)
	}
	if receipt.Name != "WrapperBody" {
		t.Errorf("expected 'WrapperBody', got %q", receipt.Name)
	}
}

func TestReadReceipt_AiReturnsInvalidJson(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK,
		`{"model":"test","created_at":"2024-01-01T00:00:00Z","message":{"role":"assistant","content":"not json"},"done":true}`)
	settings := seedReceiptProcessingFixtures(t, server.URL)

	service, err := NewReceiptProcessingService(nil, utils.UintToString(settings.ID), "")
	if err != nil {
		t.Fatalf("constructor: %v", err)
	}

	_, metadata, err := service.ReadReceiptText("body")
	if err == nil {
		t.Fatal("expected JSON unmarshal error")
	}
	if metadata.DidReceiptProcessingSettingsSucceed {
		t.Error("expected primary to NOT succeed on invalid JSON")
	}
}
