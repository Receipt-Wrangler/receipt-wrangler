package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestGenerateHmacShouldMatchKnownVector(t *testing.T) {
	// RFC 4231, Test Case 2 (HMAC-SHA256)
	key := []byte("Jefe")
	data := []byte("what do ya want for nothing?")
	expected := "5bdcc146bf60754e6a042426089575c75a003f089d2739839dec58b964ec3843"

	mac := GenerateHmac(key, data)

	actual := hex.EncodeToString(mac)
	if actual != expected {
		PrintTestError(t, actual, expected)
	}

	if len(mac) != 32 {
		PrintTestError(t, len(mac), 32)
	}
}

func TestGenerateHmacShouldBeDeterministic(t *testing.T) {
	key := []byte("some-key")
	data := []byte("some-data")

	first := GenerateHmac(key, data)
	second := GenerateHmac(key, data)

	if !bytes.Equal(first, second) {
		t.Errorf("Expected identical HMACs for the same key/data, got %x and %x", first, second)
	}
}

func TestGenerateHmacShouldDifferWithDifferentKey(t *testing.T) {
	data := []byte("some-data")

	first := GenerateHmac([]byte("key-a"), data)
	second := GenerateHmac([]byte("key-b"), data)

	if bytes.Equal(first, second) {
		t.Errorf("Expected different HMACs for different keys, got identical output")
	}
}
