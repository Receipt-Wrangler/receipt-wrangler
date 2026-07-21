package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetBodyDataGetsData(t *testing.T) {
	testString := "my test string wowzer"
	reader := strings.NewReader(testString)
	r := httptest.NewRequest(http.MethodGet, "/api", reader)
	w := httptest.NewRecorder()

	bytes, err := GetBodyData(w, r)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if string(bytes) != testString {
		PrintTestError(t, string(bytes), testString)
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

func TestMarshalResponseDataShouldMarshal(t *testing.T) {
	data := map[string]string{"hello": "world"}

	bytes, err := MarshalResponseData(data)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	expected := `{"hello":"world"}`
	if string(bytes) != expected {
		PrintTestError(t, string(bytes), expected)
	}
}

func TestMarshalResponseDataShouldErrorOnUnmarshalableValue(t *testing.T) {
	// A channel cannot be marshaled to JSON.
	_, err := MarshalResponseData(make(chan int))
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestSetJSONResponseHeadersShouldSetContentType(t *testing.T) {
	w := httptest.NewRecorder()

	SetJSONResponseHeaders(w)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		PrintTestError(t, contentType, "application/json")
	}
}

func TestIsMobileAppShouldReturnTrueForDartUserAgent(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	r.Header.Set("User-Agent", "MyApp/1.0 (dart:io)")

	if !IsMobileApp(r) {
		t.Errorf("Expected IsMobileApp to return true for a dart:io user agent")
	}
}

func TestIsMobileAppShouldReturnFalseForBrowserUserAgent(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	r.Header.Set("User-Agent", "Mozilla/5.0")

	if IsMobileApp(r) {
		t.Errorf("Expected IsMobileApp to return false for a browser user agent")
	}
}

// TODO: move test
// func TestWriteValidatorErrorResponseWritesResponse(t *testing.T) {
// 	var errBytes = make([]byte, 100)
// 	var bodyVErr structs.ValidatorError
// 	vErr := structs.ValidatorError{
// 		Errors: make(map[string]string),
// 	}
// 	nameErr := "error"
// 	amountErr := "amount cannot be empty"

// 	vErr.Errors["name"] = nameErr
// 	vErr.Errors["amount"] = amountErr

// 	w := httptest.NewRecorder()

// 	WriteValidatorErrorResponse(w, vErr, 400)

// 	if w.Result().StatusCode != 400 {
// 		PrintTestError(t, w.Result().StatusCode, 400)
// 	}

// 	w.Body.Read(errBytes)
// 	json.Unmarshal(errBytes[0:50], &bodyVErr)

// 	if reflect.DeepEqual(vErr, bodyVErr) {
// 		PrintTestError(t, vErr, bodyVErr)
// 	}
// }
