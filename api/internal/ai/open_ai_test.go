package ai

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

type capturedOpenAiRequest struct {
	Method      string
	Path        string
	RawQuery    string
	ContentType string
	Header      http.Header
	Body        map[string]interface{}
}

func newMockOpenAiServer(t *testing.T, statusCode int, responseBody string) (*httptest.Server, *capturedOpenAiRequest) {
	t.Helper()
	captured := &capturedOpenAiRequest{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.Method = r.Method
		captured.Path = r.URL.Path
		captured.RawQuery = r.URL.RawQuery
		captured.ContentType = r.Header.Get("Content-Type")
		captured.Header = r.Header.Clone()

		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil && len(bodyBytes) > 0 {
			_ = json.Unmarshal(bodyBytes, &captured.Body)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	t.Cleanup(server.Close)
	return server, captured
}

func openAiSuccessBody(content string) string {
	resp := map[string]interface{}{
		"id":      "chatcmpl-test",
		"object":  "chat.completion",
		"created": 1234567890,
		"model":   "gpt-3.5-turbo",
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": content,
				},
				"finish_reason": "stop",
			},
		},
	}
	raw, _ := json.Marshal(resp)
	return string(raw)
}

func newOpenAiClient(url, model, apiKey string, enforceJSON bool, messages []structs.AiClientMessage) *OpenAiClient {
	return NewOpenAiClient(
		structs.AiChatCompletionOptions{
			Messages:   messages,
			DecryptKey: false,
		},
		models.ReceiptProcessingSettings{
			Url:                       url,
			Model:                     model,
			Key:                       apiKey,
			EnforceJsonResponseFormat: enforceJSON,
		},
	)
}

func TestOpenAiGetChatCompletion_HappyPath_TextOnly(t *testing.T) {
	server, captured := newMockOpenAiServer(t, http.StatusOK, openAiSuccessBody("hello from openai"))

	client := newOpenAiClient(server.URL, "gpt-4o-mini", "test-key", false, []structs.AiClientMessage{
		{Role: "user", Content: "hi"},
	})
	result, err := client.GetChatCompletion()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if result.Response != "hello from openai" {
		utils.PrintTestError(t, result.Response, "hello from openai")
	}
	if result.RawResponse == "" {
		utils.PrintTestError(t, "empty RawResponse", "non-empty")
	}

	if captured.Method != http.MethodPost {
		utils.PrintTestError(t, captured.Method, http.MethodPost)
	}
	if !strings.Contains(captured.Path, "chat/completions") {
		utils.PrintTestError(t, captured.Path, "path containing chat/completions")
	}
	if captured.Body["model"] != "gpt-4o-mini" {
		utils.PrintTestError(t, captured.Body["model"], "gpt-4o-mini")
	}
	if _, hasRespFmt := captured.Body["response_format"]; hasRespFmt {
		utils.PrintTestError(t, "response_format present for non-JSON mode", "absent")
	}
}

func TestOpenAiGetChatCompletion_EnforceJsonResponseFormat(t *testing.T) {
	server, captured := newMockOpenAiServer(t, http.StatusOK, openAiSuccessBody("{}"))

	client := newOpenAiClient(server.URL, "gpt-4o-mini", "test-key", true, []structs.AiClientMessage{
		{Role: "user", Content: "hi"},
	})
	_, err := client.GetChatCompletion()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	respFmt, ok := captured.Body["response_format"].(map[string]interface{})
	if !ok {
		utils.PrintTestError(t, captured.Body["response_format"], "json_object response_format object")
		return
	}
	if respFmt["type"] != "json_object" {
		utils.PrintTestError(t, respFmt["type"], "json_object")
	}
}

func TestOpenAiGetChatCompletion_DefaultModelWhenEmpty(t *testing.T) {
	server, captured := newMockOpenAiServer(t, http.StatusOK, openAiSuccessBody("ok"))

	// Empty model → default gpt-3.5-turbo.
	client := newOpenAiClient(server.URL, "", "test-key", false, []structs.AiClientMessage{
		{Role: "user", Content: "hi"},
	})
	_, err := client.GetChatCompletion()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if captured.Body["model"] != "gpt-3.5-turbo" {
		utils.PrintTestError(t, captured.Body["model"], "gpt-3.5-turbo")
	}
}

func TestOpenAiGetChatCompletion_VisionMessageIncludesImagePart(t *testing.T) {
	server, captured := newMockOpenAiServer(t, http.StatusOK, openAiSuccessBody("looks like a receipt"))

	client := newOpenAiClient(server.URL, "gpt-4o-mini", "test-key", false, []structs.AiClientMessage{
		{
			Role:    "user",
			Content: "What's in this image?",
			Images:  []string{"data:image/png;base64,iVBORw0K"},
		},
	})
	_, err := client.GetChatCompletion()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	// Messages is an array; each element has a content array with text + image_url parts.
	messages, ok := captured.Body["messages"].([]interface{})
	if !ok || len(messages) == 0 {
		utils.PrintTestError(t, captured.Body["messages"], "non-empty array")
		return
	}
	first, _ := messages[0].(map[string]interface{})
	parts, ok := first["content"].([]interface{})
	if !ok || len(parts) != 2 {
		utils.PrintTestError(t, first["content"], "two-part content array")
		return
	}

	textPart, _ := parts[0].(map[string]interface{})
	imagePart, _ := parts[1].(map[string]interface{})
	if textPart["type"] != "text" {
		utils.PrintTestError(t, textPart["type"], "text")
	}
	if imagePart["type"] != "image_url" {
		utils.PrintTestError(t, imagePart["type"], "image_url")
	}
	imageUrl, _ := imagePart["image_url"].(map[string]interface{})
	if imageUrl["url"] != "data:image/png;base64,iVBORw0K" {
		utils.PrintTestError(t, imageUrl["url"], "data:image/png;base64,iVBORw0K")
	}
}

func TestOpenAiGetChatCompletion_HttpError500(t *testing.T) {
	server, _ := newMockOpenAiServer(t, http.StatusInternalServerError, `{"error":{"message":"server error"}}`)

	client := newOpenAiClient(server.URL, "gpt-4o-mini", "test-key", false, []structs.AiClientMessage{
		{Role: "user", Content: "hi"},
	})
	result, err := client.GetChatCompletion()
	if err == nil {
		utils.PrintTestError(t, err, "expected error from 500 response")
	}
	// RawResponse should carry a marshaled (possibly empty) response struct.
	if result.RawResponse == "" {
		utils.PrintTestError(t, "empty RawResponse on error", "non-empty marshaled response")
	}
}

func TestOpenAiGetChatCompletion_DecryptError(t *testing.T) {
	server, _ := newMockOpenAiServer(t, http.StatusOK, openAiSuccessBody("ok"))

	// DecryptKey=true with an invalid base64 key → getKey returns an error
	// BEFORE any HTTP is attempted.
	t.Setenv("ENCRYPTION_KEY", "whatever")

	client := NewOpenAiClient(
		structs.AiChatCompletionOptions{
			Messages:   []structs.AiClientMessage{{Role: "user", Content: "hi"}},
			DecryptKey: true,
		},
		models.ReceiptProcessingSettings{
			Url: server.URL,
			Key: "@@@", // not valid base64
		},
	)
	_, err := client.GetChatCompletion()
	if err == nil {
		utils.PrintTestError(t, err, "expected decrypt error")
	}
}

// Empty choices[] response: the client must return a clean error and keep
// RawResponse populated for debuggability, rather than panicking on the
// pre-fix `resp.Choices[0]` read (BUG-4).
func TestOpenAiGetChatCompletion_EmptyChoicesReturnsError(t *testing.T) {
	server, _ := newMockOpenAiServer(t, http.StatusOK, `{"id":"x","object":"chat.completion","created":1,"model":"gpt-3.5-turbo","choices":[]}`)

	client := newOpenAiClient(server.URL, "gpt-3.5-turbo", "test-key", false, []structs.AiClientMessage{
		{Role: "user", Content: "hi"},
	})
	result, err := client.GetChatCompletion()
	if err == nil {
		utils.PrintTestError(t, err, "expected error for empty choices")
	}
	if result.RawResponse == "" {
		utils.PrintTestError(t, "empty RawResponse", "raw server body preserved for logging")
	}
	if result.Response != "" {
		utils.PrintTestError(t, result.Response, "")
	}
}

// An Azure URL must be used verbatim, exactly like any other OpenAI-compatible
// endpoint. The client used to sniff the URL for "azure" and switch the SDK into
// Azure mode, which rewrote a correct Foundry base URL
// (https://<resource>.services.ai.azure.com/openai/v1) into
// .../openai/v1/openai/deployments/<model>/chat/completions?api-version=2023-05-15
// and swapped Bearer auth for the api-key header — a 404 from Azure.
func TestOpenAiGetChatCompletion_AzureUrlIsUsedVerbatim(t *testing.T) {
	server, captured := newMockOpenAiServer(t, http.StatusOK, openAiSuccessBody("azure-ok"))

	// The mock's host has no "azure" in it, so put it in the path to reproduce a
	// URL the old heuristic would have matched.
	url := server.URL + "/my-resource.services.ai.azure.com/openai/v1"

	// A dotted deployment name: Azure mode also ran the model through
	// AzureModelMapperFunc, which stripped "." and ":".
	client := newOpenAiClient(url, "gpt-4.1", "test-key", false, []structs.AiClientMessage{
		{Role: "user", Content: "hi"},
	})
	result, err := client.GetChatCompletion()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if result.Response != "azure-ok" {
		utils.PrintTestError(t, result.Response, "azure-ok")
	}

	expectedPath := "/my-resource.services.ai.azure.com/openai/v1/chat/completions"
	if captured.Path != expectedPath {
		utils.PrintTestError(t, captured.Path, expectedPath)
	}
	// Azure mode appended ?api-version=<hardcoded>; the plain path adds no query.
	if captured.RawQuery != "" {
		utils.PrintTestError(t, captured.RawQuery, "no query string")
	}
	if captured.Header.Get("Authorization") != "Bearer test-key" {
		utils.PrintTestError(t, captured.Header.Get("Authorization"), "Bearer test-key")
	}
	if captured.Header.Get("api-key") != "" {
		utils.PrintTestError(t, captured.Header.Get("api-key"), "no api-key header")
	}
	if captured.Body["model"] != "gpt-4.1" {
		utils.PrintTestError(t, captured.Body["model"], "gpt-4.1")
	}
}
