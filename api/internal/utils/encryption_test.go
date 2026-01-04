package utils

import (
	"testing"
)

func TestShouldEncryptStringWithAES128(t *testing.T) {
	key := "superSecureKey"
	value := []byte("superSecretData")

	cipherText, err := EncryptData(key, value)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	encodedCipherText := Base64Encode(cipherText)

	if len(encodedCipherText) != 60 {
		PrintTestError(t, len(encodedCipherText), 60)
	}
}

func TestShouldEncryptStringWithAES128InOneCall(t *testing.T) {
	key := "superSecureKey"
	value := "superSecretData"

	encodedCipherText, err := EncryptAndEncodeToBase64(key, value)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if len(encodedCipherText) != 60 {
		PrintTestError(t, len(encodedCipherText), 60)
	}
}

func TestShouldReturnErrorEncryptingWithEmptyKey(t *testing.T) {
	key := ""
	value := []byte("superSecretData")

	_, err := EncryptData(key, value)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestShouldReturnErrorEncryptingWithEmptyValue(t *testing.T) {
	key := "superSecureKey"
	value := []byte("")

	_, err := EncryptData(key, value)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestShouldDecryptStringWithAES128(t *testing.T) {
	key := "superSecureKey"
	value := []byte("superSecretData")

	cipherText, err := EncryptData(key, value)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	encodedCipherText := Base64Encode(cipherText)

	if len(encodedCipherText) != 60 {
		PrintTestError(t, len(encodedCipherText), 60)
	}

	clearText, err := DecryptData(key, cipherText)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if clearText != "superSecretData" {
		PrintTestError(t, clearText, "superSecretData")
	}
}

func TestShouldDecryptB64EncodedData(t *testing.T) {
	expected := "superSecretData"

	key := "superSecureKey"
	value := []byte(expected)

	encryptedData, err := EncryptData(key, value)
	if err != nil {
		PrintTestError(t, err, nil)
		return
	}

	encodedCipherText := Base64Encode(encryptedData)
	if len(encodedCipherText) != 60 {
		PrintTestError(t, len(encodedCipherText), 60)
		return
	}

	cleartext, err := DecryptB64EncodedData(key, encodedCipherText)
	if err != nil {
		PrintTestError(t, err, nil)
		return
	}

	if cleartext != expected {
		PrintTestError(t, cleartext, expected)
		return
	}
}

func TestShouldReturnErrorDecryptingWithEmptyKey(t *testing.T) {
	key := ""
	value := []byte("superSecretData")

	_, err := DecryptData(key, value)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestShouldReturnErrorDecryptingWithEmptyValue(t *testing.T) {
	key := "superSecureKey"
	value := []byte("")

	_, err := DecryptData(key, value)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestShouldEncodeValueToBase64(t *testing.T) {
	value := []byte("superSecretData")
	expected := "c3VwZXJTZWNyZXREYXRh"

	encodedValue := Base64Encode(value)

	if encodedValue != expected {
		PrintTestError(t, encodedValue, expected)
	}
}

func TestDecryptB64EncodedDataShouldReturnErrorForInvalidBase64(t *testing.T) {
	key := "superSecureKey"
	invalidBase64 := "not-valid-base64!!!"

	_, err := DecryptB64EncodedData(key, invalidBase64)
	if err == nil {
		t.Errorf("Expected error for invalid base64 input, got nil")
	}
}

func TestDecryptB64EncodedDataShouldReturnErrorForInvalidCipherText(t *testing.T) {
	key := "superSecureKey"
	// Valid base64 but invalid cipher text
	validBase64InvalidCipher := Base64Encode([]byte("this is not encrypted data but is longer than nonce"))

	_, err := DecryptB64EncodedData(key, validBase64InvalidCipher)
	if err == nil {
		t.Errorf("Expected error for invalid cipher text, got nil")
	}
}

func TestDecryptDataShouldReturnErrorForCorruptCipherText(t *testing.T) {
	key := "superSecureKey"
	// Create corrupt cipher text that's long enough to have nonce + ciphertext
	corruptCipher := make([]byte, 50)
	for i := range corruptCipher {
		corruptCipher[i] = byte(i)
	}

	_, err := DecryptData(key, corruptCipher)
	if err == nil {
		t.Errorf("Expected error for corrupt cipher text, got nil")
	}
}

func TestEncryptAndEncodeToBase64ShouldReturnErrorForEmptyKey(t *testing.T) {
	key := ""
	value := "someValue"

	_, err := EncryptAndEncodeToBase64(key, value)
	if err == nil {
		t.Errorf("Expected error for empty key, got nil")
	}
}
