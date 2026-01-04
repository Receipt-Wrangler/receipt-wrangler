package utils

import (
	"bytes"
	"testing"
)

func TestBase64EncodeShouldEncodeBytes(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8="},
		{[]byte("world"), "d29ybGQ="},
		{[]byte(""), ""},
		{[]byte("a"), "YQ=="},
		{[]byte("ab"), "YWI="},
		{[]byte("abc"), "YWJj"},
	}

	for _, tc := range testCases {
		result := Base64Encode(tc.input)
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestBase64DecodeShouldDecodeString(t *testing.T) {
	testCases := []struct {
		input    string
		expected []byte
	}{
		{"aGVsbG8=", []byte("hello")},
		{"d29ybGQ=", []byte("world")},
		{"", []byte{}},
		{"YQ==", []byte("a")},
		{"YWI=", []byte("ab")},
		{"YWJj", []byte("abc")},
	}

	for _, tc := range testCases {
		result, err := Base64Decode(tc.input)
		if err != nil {
			PrintTestError(t, err, nil)
		}
		if !bytes.Equal(result, tc.expected) {
			PrintTestError(t, string(result), string(tc.expected))
		}
	}
}

func TestBase64DecodeShouldReturnErrorForInvalidInput(t *testing.T) {
	testCases := []string{
		"not-valid-base64!!!",
		"====",
		"YQ===", // Invalid padding
	}

	for _, tc := range testCases {
		_, err := Base64Decode(tc)
		if err == nil {
			t.Errorf("Expected error for invalid base64 input %q, got nil", tc)
		}
	}
}

func TestBase64URLEncodeShouldEncodeBytes(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8="},
		{[]byte("test+data/here"), "dGVzdCtkYXRhL2hlcmU="},
		{[]byte(""), ""},
	}

	for _, tc := range testCases {
		result := Base64URLEncode(tc.input)
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestBase64URLDecodeShouldDecodeString(t *testing.T) {
	testCases := []struct {
		input    string
		expected []byte
	}{
		{"aGVsbG8=", []byte("hello")},
		{"dGVzdCtkYXRhL2hlcmU=", []byte("test+data/here")},
		{"", []byte{}},
	}

	for _, tc := range testCases {
		result, err := Base64URLDecode(tc.input)
		if err != nil {
			PrintTestError(t, err, nil)
		}
		if !bytes.Equal(result, tc.expected) {
			PrintTestError(t, string(result), string(tc.expected))
		}
	}
}

func TestBase64URLDecodeShouldReturnErrorForInvalidInput(t *testing.T) {
	testCases := []string{
		"not-valid!!!",
		"====",
	}

	for _, tc := range testCases {
		_, err := Base64URLDecode(tc)
		if err == nil {
			t.Errorf("Expected error for invalid URL base64 input %q, got nil", tc)
		}
	}
}

func TestBase64EncodeStringShouldEncodeString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello", "aGVsbG8="},
		{"world", "d29ybGQ="},
		{"", ""},
		{"test string", "dGVzdCBzdHJpbmc="},
	}

	for _, tc := range testCases {
		result := Base64EncodeString(tc.input)
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestBuildDataURIShouldCreateValidDataURI(t *testing.T) {
	testCases := []struct {
		mimeType string
		data     []byte
		expected string
	}{
		{"image/png", []byte("PNG"), "data:image/png;base64,UE5H"},
		{"application/json", []byte(`{"key":"value"}`), "data:application/json;base64,eyJrZXkiOiJ2YWx1ZSJ9"},
		{"text/plain", []byte("Hello World"), "data:text/plain;base64,SGVsbG8gV29ybGQ="},
		{"image/jpeg", []byte{}, "data:image/jpeg;base64,"},
	}

	for _, tc := range testCases {
		result := BuildDataURI(tc.mimeType, tc.data)
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}
