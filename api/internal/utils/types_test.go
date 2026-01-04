package utils

import "testing"

func TestUintToStringShouldConvertUintToString(t *testing.T) {
	testCases := []struct {
		input    uint
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{123, "123"},
		{4294967295, "4294967295"}, // Max uint32
	}

	for _, tc := range testCases {
		result := UintToString(tc.input)
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestStringToUintShouldConvertStringToUint(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint
	}{
		{"0", 0},
		{"1", 1},
		{"123", 123},
		{"  456  ", 456}, // Test trimming
	}

	for _, tc := range testCases {
		result, err := StringToUint(tc.input)
		if err != nil {
			PrintTestError(t, err, nil)
		}
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestStringToUintShouldReturnErrorForInvalidInput(t *testing.T) {
	testCases := []string{
		"abc",
		"12.34",
		"-1",
		"",
		"not a number",
	}

	for _, tc := range testCases {
		_, err := StringToUint(tc)
		if err == nil {
			t.Errorf("Expected error for input %q, got nil", tc)
		}
	}
}

func TestStringToUint64ShouldConvertStringToUint64(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint64
	}{
		{"0", 0},
		{"1", 1},
		{"123", 123},
		{"4294967295", 4294967295}, // Max uint32
	}

	for _, tc := range testCases {
		result, err := StringToUint64(tc.input)
		if err != nil {
			PrintTestError(t, err, nil)
		}
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestStringToUint64ShouldReturnErrorForInvalidInput(t *testing.T) {
	testCases := []string{
		"abc",
		"12.34",
		"-1",
		"",
		"not a number",
	}

	for _, tc := range testCases {
		_, err := StringToUint64(tc)
		if err == nil {
			t.Errorf("Expected error for input %q, got nil", tc)
		}
	}
}

func TestStringToIntShouldConvertStringToInt(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"1", 1},
		{"123", 123},
		{"-456", -456},
		{"2147483647", 2147483647}, // Max int32
	}

	for _, tc := range testCases {
		result, err := StringToInt(tc.input)
		if err != nil {
			PrintTestError(t, err, nil)
		}
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestStringToIntShouldReturnErrorForInvalidInput(t *testing.T) {
	testCases := []string{
		"abc",
		"12.34",
		"",
		"not a number",
	}

	for _, tc := range testCases {
		_, err := StringToInt(tc)
		if err == nil {
			t.Errorf("Expected error for input %q, got nil", tc)
		}
	}
}
