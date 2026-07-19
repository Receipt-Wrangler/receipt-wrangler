package utils

import "testing"

func TestSha256Hash128Bit(t *testing.T) {
	value := "superSecretData"
	expected := "4M2yAEADbol2mSGOXAMLNA=="

	hashedValue := Sha256Hash128Bit(value)
	hashedString := Base64Encode(hashedValue)

	if hashedString != expected {
		PrintTestError(t, hashedString, expected)
	}

	if len(hashedValue) != 16 {
		PrintTestError(t, len(hashedValue), 16)
	}
}

func TestSha256HashShouldMatchKnownVector(t *testing.T) {
	// SHA-256 of "abc"
	expected := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"

	actual := Sha256Hash([]byte("abc"))
	if actual != expected {
		PrintTestError(t, actual, expected)
	}
}

func TestSha256HashShouldHashEmptyInput(t *testing.T) {
	// SHA-256 of an empty input
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	actual := Sha256Hash([]byte(""))
	if actual != expected {
		PrintTestError(t, actual, expected)
	}
}
