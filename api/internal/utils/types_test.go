package utils

import "testing"

func TestUintToStringShouldFormat(t *testing.T) {
	if UintToString(123) != "123" {
		PrintTestError(t, UintToString(123), "123")
	}

	if UintToString(0) != "0" {
		PrintTestError(t, UintToString(0), "0")
	}
}

func TestStringToUintShouldParse(t *testing.T) {
	result, err := StringToUint("123")
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if result != 123 {
		PrintTestError(t, result, 123)
	}
}

func TestStringToUintShouldTrimWhitespace(t *testing.T) {
	result, err := StringToUint("  456  ")
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if result != 456 {
		PrintTestError(t, result, 456)
	}
}

func TestStringToUintShouldErrorOnNonNumeric(t *testing.T) {
	_, err := StringToUint("notanumber")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestStringToUint64ShouldParse(t *testing.T) {
	result, err := StringToUint64("789")
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if result != 789 {
		PrintTestError(t, result, 789)
	}
}

func TestStringToUint64ShouldErrorOnNonNumeric(t *testing.T) {
	_, err := StringToUint64("abc")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestStringToIntShouldParse(t *testing.T) {
	result, err := StringToInt("-42")
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if result != -42 {
		PrintTestError(t, result, -42)
	}
}

func TestStringToIntShouldErrorOnNonNumeric(t *testing.T) {
	_, err := StringToInt("xyz")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}
