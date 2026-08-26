package utils

import (
	"testing"
	"time"
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
	err = VerifyPassword(string(hashedPassword), password)
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

	err = VerifyPassword(string(hashedPassword), password)
	if err != nil {
		PrintTestError(t, err, nil)
	}
}

func TestVerifyPasswordShouldSucceedWithCorrectPassword(t *testing.T) {
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = VerifyPassword(string(hashedPassword), password)
	if err != nil {
		t.Errorf("VerifyPassword should succeed with correct password, got: %v", err)
	}
}

func TestVerifyPasswordShouldFailWithWrongPassword(t *testing.T) {
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = VerifyPassword(string(hashedPassword), "wrongPassword")
	if err == nil {
		t.Errorf("VerifyPassword should fail with wrong password")
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

func TestGetRefreshTokenExpiryDateShouldHonorTheGivenLifetime(t *testing.T) {
	// The lifetime is configurable per install (System Settings), so this helper
	// must apply whatever it is handed rather than a baked-in constant.
	lifetimes := map[string]time.Duration{
		"one hour":             time.Hour,
		"the default 24 hours": 24 * time.Hour,
		"the 30 day maximum":   720 * time.Hour,
	}

	for name, lifetime := range lifetimes {
		t.Run(name, func(t *testing.T) {
			before := time.Now().Truncate(time.Second)
			expiryDate := GetRefreshTokenExpiryDate(lifetime)
			after := time.Now().Add(time.Second).Truncate(time.Second)

			if expiryDate == nil {
				t.Errorf("Expected non-nil expiry date, got nil")
				return
			}

			// The expiry date should be between <lifetime> from before and <lifetime>
			// from after. JWT NumericDate has second precision, so we truncate
			// comparisons to seconds.
			if expiryDate.Time.Before(before.Add(lifetime)) || expiryDate.Time.After(after.Add(lifetime)) {
				t.Errorf("Expected expiry date to be approximately %v from now, got %v", lifetime, expiryDate.Time)
			}
		})
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
