package utils

import (
	"strings"
	"testing"
)

func TestGetRandomStringShouldReturnDecodableStringOfRequestedByteLength(t *testing.T) {
	length := 32

	random, err := GetRandomString(length)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	decoded, err := Base64URLDecode(random)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if len(decoded) != length {
		PrintTestError(t, len(decoded), length)
	}
}

func TestGetRandomStringShouldReturnUniqueValues(t *testing.T) {
	first, err := GetRandomString(16)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	second, err := GetRandomString(16)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if first == second {
		t.Errorf("Expected two random strings to differ, got identical values")
	}
}

func TestRemoveJsonFormatShouldStripFences(t *testing.T) {
	input := "```json\n{\"a\":1}\n```"
	expected := "\n{\"a\":1}\n"

	actual := RemoveJsonFormat(input)
	if actual != expected {
		PrintTestError(t, actual, expected)
	}

	if strings.Contains(actual, "```") {
		t.Errorf("Expected all code fences to be removed, got %q", actual)
	}
}

func TestRemoveJsonFormatShouldLeavePlainStringUnchanged(t *testing.T) {
	input := "just a plain string"

	actual := RemoveJsonFormat(input)
	if actual != input {
		PrintTestError(t, actual, input)
	}
}
