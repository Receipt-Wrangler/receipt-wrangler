package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	config "receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

const testIssuer = "https://receiptWrangler.io"
const testAudience = "https://receiptWrangler.io"

// signClaims signs a Claims value with the server secret so it passes the token
// validator (same issuer/audience/algorithm the middleware requires).
func signClaims(t *testing.T, claims structs.Claims) string {
	t.Helper()
	claims.Issuer = testIssuer
	claims.Audience = []string{testAudience}
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(10 * time.Minute))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte(config.GetSecretKey()))
	if err != nil {
		t.Fatalf("sign claims: %v", err)
	}
	return signed
}

func callAuth(t *testing.T, bearer string) int {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	r.Header.Set("Authorization", "Bearer "+bearer)
	w := httptest.NewRecorder()
	UnifiedAuthMiddleware(createFakeHandler()).ServeHTTP(w, r)
	return w.Result().StatusCode
}

// The exploit: a refresh token must NOT be accepted as an API access token.
// A same-session access token must still be accepted, and a legacy typeless
// access token (no tokenType, no jti — issued before this fix) must keep working
// so no user is forced to re-login.
func TestUnifiedAuthMiddleware_RejectsRefreshTokenAsAccess(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := createTestUser()

	accessToken, refreshToken, _, err := services.GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("generate tokens: %v", err)
	}

	// Refresh token as a Bearer access credential -> rejected.
	if status := callAuth(t, refreshToken); status != http.StatusForbidden {
		utils.PrintTestError(t, status, http.StatusForbidden)
	}

	// The real access token -> accepted.
	if status := callAuth(t, accessToken); status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Legacy access token (no tokenType, no jti) -> still accepted (no forced re-login).
	legacyAccess := signClaims(t, structs.Claims{UserId: user.ID, Username: user.Username, Displayname: user.DisplayName})
	if status := callAuth(t, legacyAccess); status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Legacy refresh token (no tokenType but WITH a jti) -> rejected via the jti discriminator.
	legacyRefresh := signClaims(t, structs.Claims{
		UserId:           user.ID,
		Username:         user.Username,
		Displayname:      user.DisplayName,
		RegisteredClaims: jwt.RegisteredClaims{ID: "legacy-jti"},
	})
	if status := callAuth(t, legacyRefresh); status != http.StatusForbidden {
		utils.PrintTestError(t, status, http.StatusForbidden)
	}
}
