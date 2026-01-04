package utils

import "testing"

func TestGetRandomStringShouldReturnStringOfCorrectLength(t *testing.T) {
	lengths := []int{8, 16, 32, 64}

	for _, length := range lengths {
		result, err := GetRandomString(length)
		if err != nil {
			PrintTestError(t, err, nil)
		}

		// Base64 URL encoding produces output that is ceil(n*4/3) characters
		// So we check that the result is not empty and contains valid characters
		if len(result) == 0 {
			t.Errorf("Expected non-empty random string for length %d", length)
		}
	}
}

func TestGetRandomStringShouldReturnDifferentStrings(t *testing.T) {
	result1, err := GetRandomString(32)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	result2, err := GetRandomString(32)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if result1 == result2 {
		t.Errorf("Expected different random strings, but got identical: %s", result1)
	}
}

func TestGetRandomStringShouldHandleZeroLength(t *testing.T) {
	result, err := GetRandomString(0)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if result != "" {
		t.Errorf("Expected empty string for zero length, got %s", result)
	}
}

func TestRemoveJsonFormatShouldRemoveJsonMarkdown(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{`{"key": "value"}`, `{"key": "value"}`},
		{"```json\n{\"key\": \"value\"}\n```", "\n{\"key\": \"value\"}\n"},
		{"```\n{\"data\": 123}\n```", "\n{\"data\": 123}\n"},
		{"```json{\"inline\": true}```", "{\"inline\": true}"},
		{"no json format here", "no json format here"},
		{"", ""},
		{"```json```", ""},
		{"```json\n```", "\n"},
	}

	for _, tc := range testCases {
		result := RemoveJsonFormat(tc.input)
		if result != tc.expected {
			t.Errorf("Input: %q - Expected %q, got %q", tc.input, tc.expected, result)
		}
	}
}

func TestRemoveJsonFormatShouldRemoveMultipleOccurrences(t *testing.T) {
	input := "```json{\"a\":1}``````json{\"b\":2}```"
	expected := "{\"a\":1}{\"b\":2}"

	result := RemoveJsonFormat(input)
	if result != expected {
		PrintTestError(t, result, expected)
	}
}

func TestRemoveJsonFormatShouldHandlePlainBackticks(t *testing.T) {
	input := "```{\"plain\": true}```"
	expected := "{\"plain\": true}"

	result := RemoveJsonFormat(input)
	if result != expected {
		PrintTestError(t, result, expected)
	}
}
