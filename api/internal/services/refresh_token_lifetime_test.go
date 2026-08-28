package services

import (
	"context"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"testing"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// setRefreshTokenLifetimeSettings writes the two configurable lifetimes straight
// to their columns, mirroring setLoginQrSettings in system_settings_test.go.
// Raw column names are used deliberately so a value the command validator would
// reject can still be stored — that is exactly what the read-side clamp exists
// to survive.
func setRefreshTokenLifetimeSettings(t *testing.T, appHours int, mcpHours int) {
	t.Helper()

	settings, err := repositories.NewSystemSettingsRepository(nil).GetSystemSettings()
	if err != nil {
		t.Fatalf("failed to load system settings: %v", err)
	}

	err = repositories.GetDB().
		Model(&models.SystemSettings{}).
		Where("id = ?", settings.ID).
		Updates(map[string]interface{}{
			"refresh_token_valid_for_hours":     appHours,
			"mcp_refresh_token_valid_for_hours": mcpHours,
		}).Error
	if err != nil {
		t.Fatalf("failed to set refresh token lifetime settings: %v", err)
	}
}

func TestClampRefreshTokenLifetime(t *testing.T) {
	tests := []struct {
		name     string
		hours    int
		expected time.Duration
	}{
		{"unset falls back to the default", 0, 24 * time.Hour},
		{"negative falls back to the default", -5, 24 * time.Hour},
		{"above the maximum falls back to the default", 721, 24 * time.Hour},
		{"absurd value falls back to the default", 1000000, 24 * time.Hour},
		{"the minimum is honored", 1, time.Hour},
		{"an in-range value is honored", 168, 168 * time.Hour},
		{"the maximum is honored", 720, 720 * time.Hour},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := clampRefreshTokenLifetime(test.hours); got != test.expected {
				t.Errorf("clampRefreshTokenLifetime(%d) = %v, expected %v", test.hours, got, test.expected)
			}
		})
	}
}

func TestGetRefreshTokenLifetimeReadsTheConfiguredValue(t *testing.T) {
	defer repositories.TruncateTestDb()

	setRefreshTokenLifetimeSettings(t, 168, 24)

	if got := GetRefreshTokenLifetime(); got != 168*time.Hour {
		t.Errorf("GetRefreshTokenLifetime() = %v, expected %v", got, 168*time.Hour)
	}
}

// The whole point of the MCP setting being separate is that a long app window
// must not extend third-party connector tokens. Assert both directions so a
// future refactor cannot quietly collapse them onto one field.
func TestRefreshTokenLifetimesAreIndependent(t *testing.T) {
	defer repositories.TruncateTestDb()

	setRefreshTokenLifetimeSettings(t, 720, 2)

	if got := GetRefreshTokenLifetime(); got != 720*time.Hour {
		t.Errorf("GetRefreshTokenLifetime() = %v, expected %v", got, 720*time.Hour)
	}

	if got := GetMcpRefreshTokenLifetime(); got != 2*time.Hour {
		t.Errorf("GetMcpRefreshTokenLifetime() = %v, expected %v", got, 2*time.Hour)
	}

	setRefreshTokenLifetimeSettings(t, 3, 500)

	if got := GetRefreshTokenLifetime(); got != 3*time.Hour {
		t.Errorf("GetRefreshTokenLifetime() = %v, expected %v", got, 3*time.Hour)
	}

	if got := GetMcpRefreshTokenLifetime(); got != 500*time.Hour {
		t.Errorf("GetMcpRefreshTokenLifetime() = %v, expected %v", got, 500*time.Hour)
	}
}

func TestGenerateJWTStampsTheConfiguredRefreshLifetime(t *testing.T) {
	defer repositories.TruncateTestDb()

	setRefreshTokenLifetimeSettings(t, 168, 24)

	db := repositories.GetDB()
	user := models.User{Username: "refresh-lifetime-user", Password: "password", DisplayName: "Refresh Lifetime User"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	before := time.Now()
	_, refreshToken, accessTokenClaims, err := GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("failed to generate jwt: %v", err)
	}
	after := time.Now()

	// The access token is deliberately NOT configurable — both clients size
	// their 15-minute refresh timer against this fixed 20-minute window.
	assertWithinWindow(t, "access token claim", accessTokenClaims.ExpiresAt.Time, before, after, 20*time.Minute)

	restValidator, err := InitTokenValidator()
	if err != nil {
		t.Fatalf("failed to build validator: %v", err)
	}
	refreshClaims := validateClaims(t, restValidator, refreshToken)
	assertWithinWindow(t, "refresh token claim", refreshClaims.ExpiresAt.Time, before, after, 168*time.Hour)

	// The persisted row must agree with the signed claim, or the cleanup job and
	// the token itself would disagree about when the session ends.
	var stored models.RefreshToken
	if err := db.Where("user_id = ?", user.ID).First(&stored).Error; err != nil {
		t.Fatalf("failed to load stored refresh token: %v", err)
	}
	assertWithinWindow(t, "stored refresh token row", stored.ExpiresAt, before, after, 168*time.Hour)
}

func TestGenerateMcpJWTUsesTheMcpRefreshLifetime(t *testing.T) {
	defer repositories.TruncateTestDb()

	// App window far longer than the MCP one, so a token minted on the app
	// setting would be unmistakable.
	setRefreshTokenLifetimeSettings(t, 720, 6)

	db := repositories.GetDB()
	user := models.User{Username: "mcp-refresh-lifetime-user", Password: "password", DisplayName: "Mcp Refresh Lifetime User"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	audience := "https://receipts.example.com/mcp"

	before := time.Now()
	_, refreshToken, _, err := GenerateMcpJWT(user.ID, audience)
	if err != nil {
		t.Fatalf("failed to generate mcp jwt: %v", err)
	}
	after := time.Now()

	mcpValidator, err := InitMcpTokenValidator(audience)
	if err != nil {
		t.Fatalf("failed to build mcp validator: %v", err)
	}
	refreshClaims := validateClaims(t, mcpValidator, refreshToken)
	assertWithinWindow(t, "mcp refresh token claim", refreshClaims.ExpiresAt.Time, before, after, 6*time.Hour)
}

// assertWithinWindow checks that expiry lands inside [before+lifetime,
// after+lifetime], allowing a second of slack for the JWT NumericDate's
// second-level precision.
func assertWithinWindow(t *testing.T, what string, expiry time.Time, before time.Time, after time.Time, lifetime time.Duration) {
	t.Helper()

	earliest := before.Add(lifetime).Add(-time.Second)
	latest := after.Add(lifetime).Add(time.Second)

	if expiry.Before(earliest) || expiry.After(latest) {
		t.Errorf("%s expiry %v is not approximately %v from now (expected between %v and %v)", what, expiry, lifetime, earliest, latest)
	}
}

// validateClaims verifies a token against the given validator and unwraps the
// custom claims, failing the test on any error.
func validateClaims(t *testing.T, v *validator.Validator, token string) *structs.Claims {
	t.Helper()

	rawClaims, err := v.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	return rawClaims.(*validator.ValidatedClaims).CustomClaims.(*structs.Claims)
}
