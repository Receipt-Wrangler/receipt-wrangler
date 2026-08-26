package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newRefreshUser(t *testing.T) models.User {
	t.Helper()
	user := models.User{
		Username:           "refresh-user",
		Password:           "hashedpassword",
		DisplayName:        "Refresh User",
		DefaultAvatarColor: "#00FF00",
	}
	if err := repositories.GetDB().Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	return user
}

// getRefreshTokenFromRequest --------------------------------------------------

func TestGetRefreshTokenFromRequest_CookiePresent(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: "cookie-refresh-token"})
	w := httptest.NewRecorder()

	token, err := getRefreshTokenFromRequest(r, w)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if token != "cookie-refresh-token" {
		utils.PrintTestError(t, token, "cookie-refresh-token")
	}
}

func TestGetRefreshTokenFromRequest_NoCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	w := httptest.NewRecorder()

	_, err := getRefreshTokenFromRequest(r, w)
	if err == nil {
		utils.PrintTestError(t, err, "expected cookie-missing error")
	}
}

func TestGetRefreshTokenFromRequest_MobileBodyWithToken(t *testing.T) {
	body, _ := json.Marshal(commands.LogoutCommand{RefreshToken: "body-refresh-token"})
	r := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(body))
	r.Header.Set("User-Agent", "Dart/3.2 (dart:io)")
	w := httptest.NewRecorder()

	token, err := getRefreshTokenFromRequest(r, w)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if token != "body-refresh-token" {
		utils.PrintTestError(t, token, "body-refresh-token")
	}
}

func TestGetRefreshTokenFromRequest_MobileEmptyBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/refresh", strings.NewReader(""))
	r.Header.Set("User-Agent", "Dart/3.2 (dart:io)")
	w := httptest.NewRecorder()

	// Empty body -> json.Unmarshal returns "unexpected end of JSON input".
	_, err := getRefreshTokenFromRequest(r, w)
	if err == nil {
		utils.PrintTestError(t, err, "expected JSON decode error")
	}
}

// ValidateRefreshToken --------------------------------------------------------

func TestValidateRefreshToken_ValidToken(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := newRefreshUser(t)
	_, refreshToken, _, err := services.GenerateJWT(user.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: refreshToken})
	w := httptest.NewRecorder()

	nextCalled := false
	var seenRefreshTokenString interface{}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		seenRefreshTokenString = r.Context().Value("refreshTokenString")
		w.WriteHeader(http.StatusOK)
	})

	ValidateRefreshToken(nextHandler).ServeHTTP(w, r)

	if !nextCalled {
		utils.PrintTestError(t, "next not called", "next called")
	}
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}
	if seenRefreshTokenString != refreshToken {
		utils.PrintTestError(t, seenRefreshTokenString, refreshToken)
	}
}

func TestValidateRefreshToken_MissingToken(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	w := httptest.NewRecorder()

	ValidateRefreshToken(createFakeHandler()).ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusInternalServerError {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusInternalServerError)
	}
}

func TestValidateRefreshToken_InvalidToken(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: "not-a-real-jwt"})
	w := httptest.NewRecorder()

	ValidateRefreshToken(createFakeHandler()).ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusInternalServerError {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusInternalServerError)
	}
}

// RevokeRefreshToken ----------------------------------------------------------

func seedRefreshToken(t *testing.T, userId uint, used bool) string {
	t.Helper()
	raw := "raw-refresh-" + time.Now().Format(time.RFC3339Nano)
	token := models.RefreshToken{
		UserId:    userId,
		Token:     utils.Sha256Hash([]byte(raw)),
		IsUsed:    used,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := repositories.GetDB().Create(&token).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	return raw
}

func TestRevokeRefreshToken_UnusedToken(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := newRefreshUser(t)
	raw := seedRefreshToken(t, user.ID, false)

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: raw})
	w := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	RevokeRefreshToken(next).ServeHTTP(w, r)

	if !nextCalled {
		utils.PrintTestError(t, "next not called", "next called")
	}
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}

	// IsUsed should now be true.
	var dbToken models.RefreshToken
	if err := repositories.GetDB().Where("token = ?", utils.Sha256Hash([]byte(raw))).First(&dbToken).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !dbToken.IsUsed {
		utils.PrintTestError(t, dbToken.IsUsed, true)
	}
}

func TestRevokeRefreshToken_AlreadyUsedToken(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := newRefreshUser(t)
	raw := seedRefreshToken(t, user.ID, true)

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: raw})
	w := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	RevokeRefreshToken(next).ServeHTTP(w, r)

	if nextCalled {
		utils.PrintTestError(t, "next called on already-used token", "next not called")
	}
	if w.Result().StatusCode != http.StatusInternalServerError {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusInternalServerError)
	}

	// Verify the handler set clearing cookies.
	foundJwt, foundRefresh := false, false
	for _, c := range w.Result().Cookies() {
		if c.Name == constants.JwtKey && c.MaxAge == -1 {
			foundJwt = true
		}
		if c.Name == constants.RefreshTokenKey && c.MaxAge == -1 {
			foundRefresh = true
		}
	}
	if !foundJwt {
		utils.PrintTestError(t, "jwt clearing cookie missing", "jwt cookie with MaxAge=-1")
	}
	if !foundRefresh {
		utils.PrintTestError(t, "refresh clearing cookie missing", "refresh cookie with MaxAge=-1")
	}
}

// Covers the fallback path where the context does NOT carry
// "refreshTokenString" AND the request has no cookie — we expect a 500.
func TestRevokeRefreshToken_NoContextNoCookie(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	w := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	RevokeRefreshToken(next).ServeHTTP(w, r)

	if nextCalled {
		utils.PrintTestError(t, "next called with no refresh token", "next not called")
	}
	if w.Result().StatusCode != http.StatusInternalServerError {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusInternalServerError)
	}
}

// Token not found: the middleware must treat an unknown token as an auth
// failure, clear the auth cookies, and return 500 — the same shape as the
// already-used branch.
func TestRevokeRefreshToken_TokenNotFound(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: "never-stored"})
	w := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	RevokeRefreshToken(next).ServeHTTP(w, r)

	if nextCalled {
		utils.PrintTestError(t, "next called for unknown refresh token", "next not called")
	}
	if w.Result().StatusCode != http.StatusInternalServerError {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusInternalServerError)
	}

	foundJwt, foundRefresh := false, false
	for _, c := range w.Result().Cookies() {
		if c.Name == constants.JwtKey && c.MaxAge == -1 {
			foundJwt = true
		}
		if c.Name == constants.RefreshTokenKey && c.MaxAge == -1 {
			foundRefresh = true
		}
	}
	if !foundJwt {
		utils.PrintTestError(t, "jwt clearing cookie missing", "jwt cookie with MaxAge=-1")
	}
	if !foundRefresh {
		utils.PrintTestError(t, "refresh clearing cookie missing", "refresh cookie with MaxAge=-1")
	}
}

// Covers the branch where upstream middleware has already placed
// "refreshTokenString" into the context — RevokeRefreshToken should consume
// that value instead of re-reading the cookie.
func TestRevokeRefreshToken_UsesContextValue(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := newRefreshUser(t)
	raw := seedRefreshToken(t, user.ID, false)

	// Note: no cookie on the request — the context value is the only source.
	r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
	r = r.WithContext(context.WithValue(r.Context(), "refreshTokenString", raw))
	w := httptest.NewRecorder()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	RevokeRefreshToken(next).ServeHTTP(w, r)

	if !nextCalled {
		utils.PrintTestError(t, "next not called", "next called")
	}
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}
}

// Refresh tokens are single-use. Consumption used to be a read followed by a
// separate update, so two concurrent refreshes of the same token could both see
// is_used = false and rotate — defeating replay detection, and logging the loser
// out of every tab once the rotation it lost landed. Exactly one racer must win.
func TestRevokeRefreshToken_ConcurrentRefreshesOnlyOneSucceeds(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := newRefreshUser(t)
	raw := seedRefreshToken(t, user.ID, false)

	const racers = 8
	var passed atomic.Int32
	var start sync.WaitGroup
	var done sync.WaitGroup
	start.Add(1)

	for i := 0; i < racers; i++ {
		done.Add(1)
		go func() {
			defer done.Done()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				passed.Add(1)
				w.WriteHeader(http.StatusOK)
			})

			r := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
			r.AddCookie(&http.Cookie{Name: constants.RefreshTokenKey, Value: raw})
			w := httptest.NewRecorder()

			start.Wait()
			RevokeRefreshToken(next).ServeHTTP(w, r)
		}()
	}

	start.Done()
	done.Wait()

	if got := passed.Load(); got != 1 {
		utils.PrintTestError(t, got, "exactly 1 racer to pass through")
	}

	var dbToken models.RefreshToken
	if err := repositories.GetDB().Where("token = ?", utils.Sha256Hash([]byte(raw))).First(&dbToken).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !dbToken.IsUsed {
		utils.PrintTestError(t, dbToken.IsUsed, true)
	}
}
