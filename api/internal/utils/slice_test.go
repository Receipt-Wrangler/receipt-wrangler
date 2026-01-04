package utils

import "testing"

func TestContainsShouldReturnTrueWhenElementExists(t *testing.T) {
	slice := []interface{}{"a", "b", "c"}

	if !Contains(slice, "a") {
		t.Errorf("Expected slice to contain 'a'")
	}

	if !Contains(slice, "b") {
		t.Errorf("Expected slice to contain 'b'")
	}

	if !Contains(slice, "c") {
		t.Errorf("Expected slice to contain 'c'")
	}
}

func TestContainsShouldReturnFalseWhenElementDoesNotExist(t *testing.T) {
	slice := []interface{}{"a", "b", "c"}

	if Contains(slice, "d") {
		t.Errorf("Expected slice to not contain 'd'")
	}

	if Contains(slice, "z") {
		t.Errorf("Expected slice to not contain 'z'")
	}
}

func TestContainsShouldWorkWithNumbers(t *testing.T) {
	slice := []interface{}{1, 2, 3, 4, 5}

	if !Contains(slice, 1) {
		t.Errorf("Expected slice to contain 1")
	}

	if !Contains(slice, 5) {
		t.Errorf("Expected slice to contain 5")
	}

	if Contains(slice, 6) {
		t.Errorf("Expected slice to not contain 6")
	}
}

func TestContainsShouldWorkWithEmptySlice(t *testing.T) {
	slice := []interface{}{}

	if Contains(slice, "anything") {
		t.Errorf("Expected empty slice to not contain any element")
	}
}

func TestContainsShouldWorkWithNilValues(t *testing.T) {
	slice := []interface{}{nil, "a", "b"}

	if !Contains(slice, nil) {
		t.Errorf("Expected slice to contain nil")
	}
}

func TestContainsShouldWorkWithMixedTypes(t *testing.T) {
	slice := []interface{}{"string", 123, true, 3.14}

	if !Contains(slice, "string") {
		t.Errorf("Expected slice to contain 'string'")
	}

	if !Contains(slice, 123) {
		t.Errorf("Expected slice to contain 123")
	}

	if !Contains(slice, true) {
		t.Errorf("Expected slice to contain true")
	}

	if !Contains(slice, 3.14) {
		t.Errorf("Expected slice to contain 3.14")
	}

	if Contains(slice, false) {
		t.Errorf("Expected slice to not contain false")
	}
}
