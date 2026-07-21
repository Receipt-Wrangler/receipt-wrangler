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

func TestShouldReturnErrorEncryptingAndEncodingWithEmptyKey(t *testing.T) {
	_, err := EncryptAndEncodeToBase64("", "superSecretData")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestDecryptDataShouldReturnErrorForShortCiphertext(t *testing.T) {
	key := "superSecureKey"
	// Non-empty but shorter than the GCM nonce size (12 bytes). This must
	// return an error rather than panic on the nonce slice.
	shortCiphertext := []byte{1, 2, 3, 4, 5}

	_, err := DecryptData(key, shortCiphertext)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestDecryptDataShouldReturnErrorForTamperedCiphertext(t *testing.T) {
	key := "superSecureKey"

	cipherText, err := EncryptData(key, []byte("superSecretData"))
	if err != nil {
		PrintTestError(t, err, nil)
		return
	}

	// Corrupt the authentication tag (last byte) so gcm.Open fails.
	cipherText[len(cipherText)-1] ^= 0xff

	_, err = DecryptData(key, cipherText)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestDecryptB64EncodedDataShouldReturnErrorForInvalidBase64(t *testing.T) {
	_, err := DecryptB64EncodedData("superSecureKey", "!!!not-base64!!!")
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}

func TestDecryptB64EncodedDataShouldReturnErrorForWrongKey(t *testing.T) {
	cipherText, err := EncryptData("superSecureKey", []byte("superSecretData"))
	if err != nil {
		PrintTestError(t, err, nil)
		return
	}

	encoded := Base64Encode(cipherText)

	_, err = DecryptB64EncodedData("wrongKey", encoded)
	if err == nil {
		PrintTestError(t, nil, "error")
	}
}
