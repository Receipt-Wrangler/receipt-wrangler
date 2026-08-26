package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// GetRefreshTokenExpiryDate returns the expiry for a refresh token minted now.
// The lifetime is configurable per install, so callers resolve it (see
// services.GetRefreshTokenLifetime) and pass it in — this package deliberately
// imports nothing from internal/, so it cannot read System Settings itself.
func GetRefreshTokenExpiryDate(lifetime time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(lifetime))
}

// GetAccessTokenExpiryDate returns the expiry for an access token minted now.
// This one is intentionally NOT configurable: the access token is stateless and
// cannot be revoked, and both clients proactively refresh on a 15-minute timer
// sized against this 20-minute window.
func GetAccessTokenExpiryDate() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(20 * time.Minute))
}
