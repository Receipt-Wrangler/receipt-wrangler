package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestGenerateHmacShouldProduceConsistentOutput(t *testing.T) {
	key := []byte("secret-key")
	data := []byte("message to authenticate")

	result1 := GenerateHmac(key, data)
	result2 := GenerateHmac(key, data)

	if !bytes.Equal(result1, result2) {
		t.Errorf("Expected HMAC to be consistent for same key and data")
	}
}

func TestGenerateHmacShouldProduceDifferentOutputForDifferentData(t *testing.T) {
	key := []byte("secret-key")
	data1 := []byte("message 1")
	data2 := []byte("message 2")

	result1 := GenerateHmac(key, data1)
	result2 := GenerateHmac(key, data2)

	if bytes.Equal(result1, result2) {
		t.Errorf("Expected different HMAC for different data")
	}
}

func TestGenerateHmacShouldProduceDifferentOutputForDifferentKeys(t *testing.T) {
	key1 := []byte("key-1")
	key2 := []byte("key-2")
	data := []byte("same message")

	result1 := GenerateHmac(key1, data)
	result2 := GenerateHmac(key2, data)

	if bytes.Equal(result1, result2) {
		t.Errorf("Expected different HMAC for different keys")
	}
}

func TestGenerateHmacShouldReturn32Bytes(t *testing.T) {
	key := []byte("secret-key")
	data := []byte("test data")

	result := GenerateHmac(key, data)

	// SHA256 produces 32 bytes
	if len(result) != 32 {
		t.Errorf("Expected HMAC length of 32 bytes, got %d", len(result))
	}
}

func TestGenerateHmacShouldProduceKnownOutput(t *testing.T) {
	// Test against a known HMAC-SHA256 output
	key := []byte("key")
	data := []byte("The quick brown fox jumps over the lazy dog")

	result := GenerateHmac(key, data)
	resultHex := hex.EncodeToString(result)

	// Known HMAC-SHA256 output for this key/data combination
	expected := "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8"

	if resultHex != expected {
		t.Errorf("Expected HMAC %s, got %s", expected, resultHex)
	}
}

func TestGenerateHmacShouldHandleEmptyData(t *testing.T) {
	key := []byte("secret-key")
	data := []byte("")

	result := GenerateHmac(key, data)

	if len(result) != 32 {
		t.Errorf("Expected HMAC length of 32 bytes for empty data, got %d", len(result))
	}
}

func TestGenerateHmacShouldHandleEmptyKey(t *testing.T) {
	key := []byte("")
	data := []byte("test data")

	result := GenerateHmac(key, data)

	if len(result) != 32 {
		t.Errorf("Expected HMAC length of 32 bytes for empty key, got %d", len(result))
	}
}
