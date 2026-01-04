package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TODO: Fix
func TestGetBodyDataGetsData(t *testing.T) {
	var unmarshalResult any
	testString := "my test string wowzer"
	reader := strings.NewReader(testString)
	r := httptest.NewRequest(http.MethodGet, "/api", reader)
	w := httptest.NewRecorder()
	bytes, _ := GetBodyData(w, r)

	json.Unmarshal(bytes, &unmarshalResult)

	if testString != unmarshalResult {
		// repositories.PrintTestError(t, unmarshalResult, testString)
	}
}

func TestWriteErrorResponseWritesReponse(t *testing.T) {
	var errBytes = make([]byte, 100)
	var errMap map[string]string

	w := httptest.NewRecorder()
	err := fmt.Errorf("Test error")

	WriteErrorResponse(w, err, 500)

	if w.Result().StatusCode != 500 {
		PrintTestError(t, w.Result().StatusCode, 500)
	}

	w.Body.Read(errBytes)
	json.Unmarshal(errBytes[0:25], &errMap)

	if errMap[errKey] != "Test error" {
		PrintTestError(t, errMap[errKey], "Test error")
	}
}

func TestWriteCustomErrorResponseWritesResponse(t *testing.T) {
	var errBytes = make([]byte, 100)
	var errMap map[string]string

	customMsg := "Hello world"

	w := httptest.NewRecorder()

	WriteCustomErrorResponse(w, customMsg, 200)

	if w.Result().StatusCode != 200 {
		PrintTestError(t, w.Result().StatusCode, 200)
	}

	w.Body.Read(errBytes)
	json.Unmarshal(errBytes[0:26], &errMap)

	if errMap[errKey] != customMsg {
		PrintTestError(t, errMap[errKey], customMsg)
	}
}

func TestMarshalResponseDataShouldMarshalStruct(t *testing.T) {
	data := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		Name:  "test",
		Value: 123,
	}

	result, err := MarshalResponseData(data)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	expected := `{"name":"test","value":123}`
	if string(result) != expected {
		PrintTestError(t, string(result), expected)
	}
}

func TestMarshalResponseDataShouldMarshalMap(t *testing.T) {
	data := map[string]interface{}{
		"key": "value",
		"num": 42,
	}

	result, err := MarshalResponseData(data)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	// Unmarshal and verify
	var resultMap map[string]interface{}
	json.Unmarshal(result, &resultMap)

	if resultMap["key"] != "value" {
		PrintTestError(t, resultMap["key"], "value")
	}
}

func TestMarshalResponseDataShouldMarshalSlice(t *testing.T) {
	data := []string{"a", "b", "c"}

	result, err := MarshalResponseData(data)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	expected := `["a","b","c"]`
	if string(result) != expected {
		PrintTestError(t, string(result), expected)
	}
}

func TestMarshalResponseDataShouldReturnErrorForInvalidData(t *testing.T) {
	// Channels cannot be marshaled to JSON
	data := make(chan int)

	_, err := MarshalResponseData(data)
	if err == nil {
		t.Errorf("Expected error for unmarshallable data, got nil")
	}
}

func TestSetJSONResponseHeadersShouldSetContentType(t *testing.T) {
	w := httptest.NewRecorder()

	SetJSONResponseHeaders(w)

	contentType := w.Header().Get("Content-Type")
	expected := "application/json"

	if contentType != expected {
		PrintTestError(t, contentType, expected)
	}
}

func TestIsMobileAppShouldReturnTrueForDartUserAgent(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	r.Header.Set("User-Agent", "Dart/2.18 (dart:io)")

	result := IsMobileApp(r)

	if !result {
		t.Errorf("Expected IsMobileApp to return true for dart:io user agent")
	}
}

func TestIsMobileAppShouldReturnFalseForBrowserUserAgent(t *testing.T) {
	testCases := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"Chrome/91.0.4472.124 Safari/537.36",
		"",
	}

	for _, userAgent := range testCases {
		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		r.Header.Set("User-Agent", userAgent)

		result := IsMobileApp(r)

		if result {
			t.Errorf("Expected IsMobileApp to return false for user agent %q", userAgent)
		}
	}
}

func TestIsMobileAppShouldReturnTrueForFlutterUserAgent(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	r.Header.Set("User-Agent", "Flutter/3.0 (dart:io) package:http/http.dart")

	result := IsMobileApp(r)

	if !result {
		t.Errorf("Expected IsMobileApp to return true for Flutter/dart:io user agent")
	}
}
