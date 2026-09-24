package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// startLink drives GET /oidc/link/{name} as the given user.
func startLink(t *testing.T, providerName string, userId uint) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/api/oidc/link/"+providerName, nil)
	request = withUrlParam(request, "name", providerName)
	request = request.WithContext(context.WithValue(
		request.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}},
	))

	recorder := httptest.NewRecorder()
	LinkStart(recorder, request)

	return recorder
}

func countAuthSessions(t *testing.T) int64 {
	t.Helper()

	var count int64
	if err := repositories.GetDB().Model(&models.OidcAuthSession{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count auth sessions: %v", err)
	}

	return count
}

// TestLinkStartRequiresAccountUpdatePermission is the gate itself.
//
// Connecting a provider ADDS a way to sign into an account, so it must be at
// least as privileged as disconnecting one -- and DeleteOidcConnection already
// requires app.account.update. Without this check an administrator could build a
// role that cannot remove a sign-in method but can still add one.
func TestLinkStartRequiresAccountUpdatePermission(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	// Holds the read half but not the update half, which is exactly the role an
	// administrator would build to make an account's bindings look-but-don't-touch.
	user := createTestUserWithPermissions(t, "restricted", permissions.AppAccountRead)

	recorder := startLink(t, provider.Name, user.ID)

	if recorder.Code != http.StatusFound {
		t.Fatalf("expected a redirect, got %d", recorder.Code)
	}

	if code := oidcErrorCode(t, recorder); code != errForbidden {
		t.Errorf("expected the %q error code, got %q", errForbidden, code)
	}

	// Nothing may be started: a denied caller must not even reach the identity
	// provider, and must not leave a redeemable session row behind.
	if count := countAuthSessions(t); count != 0 {
		t.Errorf("expected no auth session to be created for a denied caller, got %d", count)
	}
}

// TestLinkStartFailureLandsOnTheProfileNotTheLoginPage pins the redirect target.
// A link caller is already signed in, so bouncing them to a login form reads as
// "you have been signed out" rather than "that provider could not be connected".
func TestLinkStartFailureLandsOnTheProfileNotTheLoginPage(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "restricted")

	recorder := startLink(t, provider.Name, user.ID)

	location := recorder.Header().Get("Location")
	if !strings.HasPrefix(location, desktopProfilePath) {
		t.Errorf("expected a redirect to %q, got %q", desktopProfilePath, location)
	}

	if strings.Contains(location, desktopLoginPath) {
		t.Errorf("a signed-in caller must not be sent to the login page, got %q", location)
	}
}

// TestLinkStartProceedsWithAccountUpdatePermission is the other half: the gate
// must not break the flow it guards.
func TestLinkStartProceedsWithAccountUpdatePermission(t *testing.T) {
	defer teardownOidcTest()
	idp, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "allowed", permissions.AppAccountUpdate)

	recorder := startLink(t, provider.Name, user.ID)

	if recorder.Code != http.StatusFound {
		t.Fatalf("expected a redirect, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	location, err := url.Parse(recorder.Header().Get("Location"))
	if err != nil {
		t.Fatalf("failed to parse the redirect: %v", err)
	}

	if !strings.HasPrefix(location.String(), idp.issuer()) {
		t.Fatalf("expected a redirect to the identity provider at %q, got %q", idp.issuer(), location)
	}

	var session models.OidcAuthSession
	err = repositories.GetDB().
		Model(&models.OidcAuthSession{}).
		Where("state_hash = ?", hashSecret(location.Query().Get("state"))).
		First(&session).Error
	if err != nil {
		t.Fatalf("expected an auth session for the minted state: %v", err)
	}

	if !session.IsLink() {
		t.Fatal("expected a link session, got a plain login session")
	}

	if *session.LinkUserId != user.ID {
		t.Errorf("expected the session to carry user %d, got %d", user.ID, *session.LinkUserId)
	}
}

// startMobileLink drives GET /oidc/link/{name}?client=mobile as the given user,
// which is how the app calls it: an ordinary authenticated API request.
func startMobileLink(t *testing.T, providerName string, userId uint) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/api/oidc/link/"+providerName+"?client=mobile", nil)
	request = withUrlParam(request, "name", providerName)
	request = request.WithContext(context.WithValue(
		request.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}},
	))

	recorder := httptest.NewRecorder()
	LinkStart(recorder, request)

	return recorder
}

// TestMobileLinkStartReturnsTheAuthorizationUrlAsJson is the whole reason the
// mobile link needs a shape of its own.
//
// The app authenticates with a bearer token, and the external user agent RFC
// 8252 requires cannot carry one -- so a 302 into that browser arrives
// unauthenticated and is refused. The app therefore calls this itself and opens
// the URL it gets back.
func TestMobileLinkStartReturnsTheAuthorizationUrlAsJson(t *testing.T) {
	defer teardownOidcTest()
	idp, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "mobilelinker", permissions.AppAccountUpdate)

	recorder := startMobileLink(t, provider.Name, user.ID)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	var view structs.OidcLinkStartView
	if err := json.Unmarshal(recorder.Body.Bytes(), &view); err != nil {
		t.Fatalf("failed to decode the response: %v", err)
	}

	if !strings.HasPrefix(view.AuthorizationUrl, idp.issuer()) {
		t.Fatalf("expected an authorization URL at %q, got %q", idp.issuer(), view.AuthorizationUrl)
	}

	parsed, err := url.Parse(view.AuthorizationUrl)
	if err != nil {
		t.Fatalf("failed to parse the authorization URL: %v", err)
	}

	var session models.OidcAuthSession
	err = repositories.GetDB().
		Model(&models.OidcAuthSession{}).
		Where("state_hash = ?", hashSecret(parsed.Query().Get("state"))).
		First(&session).Error
	if err != nil {
		t.Fatalf("expected an auth session for the minted state: %v", err)
	}

	if !session.IsLink() || *session.LinkUserId != user.ID {
		t.Errorf("expected a link session for user %d, got %+v", user.ID, session.LinkUserId)
	}

	if session.ClientType != models.OidcClientMobile {
		t.Errorf("expected a mobile session, got %q", session.ClientType)
	}

	// This is an API call, not a navigation. A Set-Cookie here would be a session
	// the app never asked for and cannot carry.
	if len(recorder.Result().Cookies()) > 0 {
		t.Error("a mobile link start must not set cookies")
	}
}

// TestMobileLinkStartNeedsNoCodeChallenge: the challenge binds the mobile LOGIN's
// exchange code to the app that started it. A link mints no exchange code and no
// session, and the caller already proved who they are with a bearer token on the
// start request -- so demanding one would be ceremony that binds nothing.
func TestMobileLinkStartNeedsNoCodeChallenge(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "nochallenge", permissions.AppAccountUpdate)

	// startMobileLink deliberately sends no codeChallenge.
	if recorder := startMobileLink(t, provider.Name, user.ID); recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 without a code challenge, got %d (%s)", recorder.Code, recorder.Body.String())
	}
}

// TestMobileLoginStillRequiresACodeChallenge guards the other half of that split.
// Relaxing it for the link must not relax it for the login, where the challenge
// is the only thing making an intercepted handoff code worthless.
func TestMobileLoginStillRequiresACodeChallenge(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	recorder := startLogin(t, provider.Name, "client=mobile")

	if code := oidcErrorCode(t, recorder); code != errInvalidRequest {
		t.Errorf("expected %q, got %q", errInvalidRequest, code)
	}

	if count := countAuthSessions(t); count != 0 {
		t.Errorf("no session may be created without a challenge, got %d", count)
	}
}

// TestMobileLinkStartDeniedAnswersWithJson: a denial has to reach the app in the
// same shape as a success, or it surfaces as an unparseable body rather than a
// message.
func TestMobileLinkStartDeniedAnswersWithJson(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "mobilerestricted", permissions.AppAccountRead)

	recorder := startMobileLink(t, provider.Name, user.ID)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}

	if location := recorder.Header().Get("Location"); len(location) > 0 {
		t.Errorf("an API call must not be answered with a redirect, got %q", location)
	}

	var body structs.OidcFlowError
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode the error: %v", err)
	}

	if body.ErrorCode != errForbidden {
		t.Errorf("expected %q, got %q", errForbidden, body.ErrorCode)
	}

	if count := countAuthSessions(t); count != 0 {
		t.Errorf("expected no auth session for a denied caller, got %d", count)
	}
}

// TestMobileLinkStartUnknownProviderAnswersWithJson pins the reordering that made
// the client type known before the provider lookup. Without it this failure took
// the desktop redirect branch and the app received HTML it could not read.
func TestMobileLinkStartUnknownProviderAnswersWithJson(t *testing.T) {
	defer teardownOidcTest()
	setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "unknownprovider", permissions.AppAccountUpdate)

	recorder := startMobileLink(t, "no-such-provider", user.ID)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	var body structs.OidcFlowError
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode the error: %v", err)
	}

	if body.ErrorCode != errUnknownProvider {
		t.Errorf("expected %q, got %q", errUnknownProvider, body.ErrorCode)
	}
}

// TestMobileLinkCallbackReturnsToTheAppScheme is the far end of the flow: the
// callback lands in the external browser, and the only way back into the app is
// its private-use scheme. Redirecting to the desktop profile path -- as this did
// before -- strands the browser on a page the app never sees.
func TestMobileLinkCallbackReturnsToTheAppScheme(t *testing.T) {
	defer teardownOidcTest()
	idp, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUserWithPermissions(t, "mobilecallback", permissions.AppAccountUpdate)

	recorder := startMobileLink(t, provider.Name, user.ID)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected the link to start, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	var view structs.OidcLinkStartView
	if err := json.Unmarshal(recorder.Body.Bytes(), &view); err != nil {
		t.Fatalf("failed to decode the response: %v", err)
	}

	authUrl, err := url.Parse(view.AuthorizationUrl)
	if err != nil {
		t.Fatalf("failed to parse the authorization URL: %v", err)
	}

	idp.setClaims(claimsFor(idp, "mobile-link-subject", authUrl.Query().Get("nonce"), "whoever"))

	// No binding cookie: the flow ran in an external browser that never had one.
	callback := runCallback(t, provider.Name, url.Values{
		"code":  {"abc"},
		"state": {authUrl.Query().Get("state")},
	}, nil)

	location := callback.Header().Get("Location")
	if !strings.HasPrefix(location, mobileCallbackScheme) {
		t.Fatalf("expected a return to %q, got %q", mobileCallbackScheme, location)
	}

	parsed, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse the app redirect: %v", err)
	}

	if parsed.Query().Get("linked") != provider.Name {
		t.Errorf("expected linked=%q, got %q (error %q)", provider.Name, parsed.Query().Get("linked"), parsed.Query().Get("error"))
	}

	// A link mints no session, so nothing credential-like may ride back.
	if len(callback.Result().Cookies()) > 0 {
		t.Error("a link must not set cookies")
	}

	for _, key := range []string{"code", "jwt", "token", "access_token", "refresh_token", "refreshToken"} {
		if len(parsed.Query().Get(key)) > 0 {
			t.Errorf("the app redirect must not carry %s", key)
		}
	}

	identity, err := repositories.NewOidcIdentityRepository(nil).GetIdentityBySubject(provider.ID, "mobile-link-subject")
	if err != nil {
		t.Fatalf("expected the identity to be linked: %v", err)
	}

	if identity.UserId != user.ID {
		t.Errorf("expected the link to land on user %d, got %d", user.ID, identity.UserId)
	}

	if identity.ProvisionedUser {
		t.Error("a link must not mark the account as provisioned -- it has its own password")
	}
}
