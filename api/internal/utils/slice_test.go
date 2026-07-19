package utils

import "testing"

func TestContainsShouldFindExistingValue(t *testing.T) {
	slice := []interface{}{"a", "b", "c"}

	if !Contains(slice, "b") {
		t.Errorf("Expected Contains to find 'b' in the slice")
	}
}

func TestContainsShouldNotFindMissingValue(t *testing.T) {
	slice := []interface{}{"a", "b", "c"}

	if Contains(slice, "z") {
		t.Errorf("Expected Contains to not find 'z' in the slice")
	}
}

func TestContainsShouldReturnFalseForEmptySlice(t *testing.T) {
	if Contains([]interface{}{}, "a") {
		t.Errorf("Expected Contains to return false for an empty slice")
	}
}

func TestContainsShouldMatchIntValues(t *testing.T) {
	slice := []interface{}{1, 2, 3}

	if !Contains(slice, 2) {
		t.Errorf("Expected Contains to find int 2 in the slice")
	}
}
