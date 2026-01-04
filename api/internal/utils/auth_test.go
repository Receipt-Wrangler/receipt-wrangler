package utils

import (
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordShouldHashPassword(t *testing.T) {
	password := "testPassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if len(hashedPassword) == 0 {
		PrintTestError(t, len(hashedPassword), "> 0")
	}

	// Verify the hash can be compared with the original password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		PrintTestError(t, err, nil)
	}
}

func TestHashPasswordShouldHashEmptyPassword(t *testing.T) {
	password := ""

	hashedPassword, err := HashPassword(password)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if len(hashedPassword) == 0 {
		PrintTestError(t, len(hashedPassword), "> 0")
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		PrintTestError(t, err, nil)
	}
}

func TestHashPasswordShouldReturnErrorForLongPassword(t *testing.T) {
	// bcrypt returns an error for passwords longer than 72 bytes
	password := "thisIsAVeryLongPasswordThatExceeds72BytesAndShouldStillWorkWithBcryptHashing123456789"

	_, err := HashPassword(password)
	if err == nil {
		t.Errorf("Expected error for password > 72 bytes, got nil")
	}
}

func TestHashPasswordShouldProduceDifferentHashesForSamePassword(t *testing.T) {
	password := "testPassword123"

	hash1, err := HashPassword(password)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	// bcrypt uses random salt, so hashes should be different
	if string(hash1) == string(hash2) {
		t.Errorf("Expected different hashes for same password, but got identical hashes")
	}
}

func TestGetRefreshTokenExpiryDateShouldReturn24HoursFromNow(t *testing.T) {
	before := time.Now().Truncate(time.Second)
	expiryDate := GetRefreshTokenExpiryDate()
	after := time.Now().Add(time.Second).Truncate(time.Second)

	if expiryDate == nil {
		t.Errorf("Expected non-nil expiry date, got nil")
		return
	}

	expected24HoursFromBefore := before.Add(24 * time.Hour)
	expected24HoursFromAfter := after.Add(24 * time.Hour)

	// The expiry date should be between 24 hours from before and 24 hours from after
	// JWT NumericDate has second precision, so we truncate comparisons to seconds
	if expiryDate.Time.Before(expected24HoursFromBefore) || expiryDate.Time.After(expected24HoursFromAfter) {
		t.Errorf("Expected expiry date to be approximately 24 hours from now, got %v", expiryDate.Time)
	}
}

func TestGetAccessTokenExpiryDateShouldReturn20MinutesFromNow(t *testing.T) {
	before := time.Now().Truncate(time.Second)
	expiryDate := GetAccessTokenExpiryDate()
	after := time.Now().Add(time.Second).Truncate(time.Second)

	if expiryDate == nil {
		t.Errorf("Expected non-nil expiry date, got nil")
		return
	}

	expected20MinutesFromBefore := before.Add(20 * time.Minute)
	expected20MinutesFromAfter := after.Add(20 * time.Minute)

	// The expiry date should be between 20 minutes from before and 20 minutes from after
	// JWT NumericDate has second precision, so we truncate comparisons to seconds
	if expiryDate.Time.Before(expected20MinutesFromBefore) || expiryDate.Time.After(expected20MinutesFromAfter) {
		t.Errorf("Expected expiry date to be approximately 20 minutes from now, got %v", expiryDate.Time)
	}
}
