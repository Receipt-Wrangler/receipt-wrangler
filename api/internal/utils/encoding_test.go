package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestBase64URLEncodeDecodeRoundTrip(t *testing.T) {
	data := []byte{0xff, 0xfe, 0xfd, 0x00, 0x10}

	encoded := Base64URLEncode(data)

	// The URL-safe alphabet must not contain '+' or '/'
	if strings.ContainsAny(encoded, "+/") {
		t.Errorf("Expected URL-safe encoding without '+' or '/', got %s", encoded)
	}

	decoded, err := Base64URLDecode(encoded)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if !bytes.Equal(decoded, data) {
		t.Errorf("Expected round-trip to return the original bytes, got %v", decoded)
	}
}

func TestBase64URLDecodeShouldErrorOnInvalidInput(t *testing.T) {
	_, err := Base64URLDecode("not valid base64!!!")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestBase64DecodeShouldErrorOnInvalidInput(t *testing.T) {
	_, err := Base64Decode("!!!not-valid!!!")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestBase64EncodeStringShouldEncode(t *testing.T) {
	expected := "c3VwZXJTZWNyZXREYXRh"

	actual := Base64EncodeString("superSecretData")
	if actual != expected {
		PrintTestError(t, actual, expected)
	}
}

func TestBuildDataURIShouldBuildURI(t *testing.T) {
	data := []byte("hello")
	expected := "data:text/plain;base64,aGVsbG8="

	actual := BuildDataURI("text/plain", data)
	if actual != expected {
		PrintTestError(t, actual, expected)
	}
}
