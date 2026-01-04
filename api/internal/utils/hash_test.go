package utils

import "testing"

func TestSha256HashShouldHashData(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{[]byte("world"), "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
		{[]byte(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte("test"), "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
	}

	for _, tc := range testCases {
		result := Sha256Hash(tc.input)
		if result != tc.expected {
			PrintTestError(t, result, tc.expected)
		}
	}
}

func TestSha256HashShouldProduceConsistentOutput(t *testing.T) {
	data := []byte("consistent data")

	result1 := Sha256Hash(data)
	result2 := Sha256Hash(data)

	if result1 != result2 {
		t.Errorf("Expected consistent hash output, got different results")
	}
}

func TestSha256HashShouldReturn64CharHexString(t *testing.T) {
	data := []byte("test data")

	result := Sha256Hash(data)

	// SHA256 produces 32 bytes, which is 64 hex characters
	if len(result) != 64 {
		t.Errorf("Expected hash length of 64 characters, got %d", len(result))
	}
}

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
