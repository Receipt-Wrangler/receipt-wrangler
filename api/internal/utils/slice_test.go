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

func TestContainsShouldHandleUncomparableTypesWithoutPanic(t *testing.T) {
	// Slices are not comparable with ==, which previously panicked. Contains
	// must now compare by value without panicking.
	slice := []interface{}{[]int{1, 2}}

	if !Contains(slice, []int{1, 2}) {
		t.Errorf("Expected Contains to find the matching slice value")
	}

	if Contains(slice, []int{9}) {
		t.Errorf("Expected Contains to not find a non-matching slice value")
	}
}
