package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"strings"
	"testing"
)

const testPublicUrl = "https://receipts.example.com"
const testRedirectUri = "https://claude.ai/api/mcp/auth_callback"

func setPublicUrl(t *testing.T) {
	t.Helper()
	t.Setenv(string(constants.McpPublicUrl), testPublicUrl)
}

func decodeJson(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to decode json response %q: %v", string(body), err)
	}
	return result
}

func TestProtectedResourceMetadata(t *testing.T) {
	setPublicUrl(t)

	recorder := httptest.NewRecorder()
	ProtectedResourceMetadata(recorder, httptest.NewRequest(http.MethodGet, "/.well-known/oauth-protected-resource", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	body := decodeJson(t, recorder.Body.Bytes())
	if body["resource"] != testPublicUrl+"/mcp" {
		t.Errorf("resource = %v, want %v", body["resource"], testPublicUrl+"/mcp")
	}

	servers, ok := body["authorization_servers"].([]interface{})
	if !ok || len(servers) != 1 || servers[0] != testPublicUrl {
		t.Errorf("authorization_servers = %v, want [%v]", body["authorization_servers"], testPublicUrl)
	}
}

func TestAuthorizationServerMetadata(t *testing.T) {
	setPublicUrl(t)

	recorder := httptest.NewRecorder()
	AuthorizationServerMetadata(recorder, httptest.NewRequest(http.MethodGet, "/.well-known/oauth-authorization-server", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	body := decodeJson(t, recorder.Body.Bytes())
	expectations := map[string]string{
		"issuer":                 testPublicUrl,
		"authorization_endpoint": testPublicUrl + "/oauth/authorize",
		"token_endpoint":         testPublicUrl + "/oauth/token",
		"registration_endpoint":  testPublicUrl + "/oauth/register",
	}
	for key, want := range expectations {
		if body[key] != want {
			t.Errorf("%s = %v, want %v", key, body[key], want)
		}
	}

	methods, ok := body["code_challenge_methods_supported"].([]interface{})
	if !ok || len(methods) != 1 || methods[0] != "S256" {
		t.Errorf("code_challenge_methods_supported = %v, want [S256]", body["code_challenge_methods_supported"])
	}
}

func TestRegisterPersistsClient(t *testing.T) {
	defer repositories.TruncateTestDb()

	requestBody := `{"client_name":"Claude","redirect_uris":["` + testRedirectUri + `"]}`
	recorder := httptest.NewRecorder()
	Register(recorder, httptest.NewRequest(http.MethodPost, "/oauth/register", strings.NewReader(requestBody)))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	body := decodeJson(t, recorder.Body.Bytes())
	clientId, ok := body["client_id"].(string)
	if !ok || len(clientId) == 0 {
		t.Fatalf("expected a client_id in response, got %v", body["client_id"])
	}

	client, err := getClient(clientId)
	if err != nil {
		t.Fatalf("registered client not found in db: %v", err)
	}
	if !clientAllowsRedirect(client, testRedirectUri) {
		t.Errorf("persisted client does not allow the registered redirect uri")
	}
}

func TestRegisterRequiresRedirectUri(t *testing.T) {
	recorder := httptest.NewRecorder()
	Register(recorder, httptest.NewRequest(http.MethodPost, "/oauth/register", strings.NewReader(`{"client_name":"Claude"}`)))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing redirect_uris, got %d", recorder.Code)
	}
}

func TestRegisterRejectsInvalidRedirectUri(t *testing.T) {
	invalid := []string{
		`{"redirect_uris":["not-a-url"]}`,
		`{"redirect_uris":["ftp://example.com/cb"]}`,
		`{"redirect_uris":["https://example.com/cb#frag"]}`,
		`{"redirect_uris":["https://good.example.com/cb","javascript:alert(1)"]}`,
	}

	for _, body := range invalid {
		recorder := httptest.NewRecorder()
		Register(recorder, httptest.NewRequest(http.MethodPost, "/oauth/register", strings.NewReader(body)))
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid redirect_uri in %s, got %d", body, recorder.Code)
		}
	}
}

// postAuthorize submits the login form with the given credentials against a
// freshly registered client and returns the recorder.
func postAuthorize(t *testing.T, client models.OAuthClient, redirectUri string, username string, password string, codeChallenge string) *httptest.ResponseRecorder {
	t.Helper()

	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	form.Set("client_id", client.ClientId)
	form.Set("redirect_uri", redirectUri)
	form.Set("state", "xyz")
	form.Set("code_challenge", codeChallenge)
	form.Set("code_challenge_method", "S256")

	request := httptest.NewRequest(http.MethodPost, "/oauth/authorize", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	Authorize(recorder, request)
	return recorder
}

func TestAuthorizeIssuesCodeOnValidLogin(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createTestUserWithPassword(t, "alice", "s3cret-password")
	client, err := createClient("Claude", []string{testRedirectUri})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	recorder := postAuthorize(t, client, testRedirectUri, "alice", "s3cret-password", challengeFor("verifier-123"))

	if recorder.Code != http.StatusFound {
		t.Fatalf("expected 302 redirect, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	location, err := url.Parse(recorder.Header().Get("Location"))
	if err != nil {
		t.Fatalf("invalid redirect location: %v", err)
	}
	if location.Query().Get("state") != "xyz" {
		t.Errorf("expected state to be echoed, got %q", location.Query().Get("state"))
	}

	code := location.Query().Get("code")
	if len(code) == 0 {
		t.Fatalf("expected an authorization code in the redirect, got none")
	}

	var stored models.OAuthAuthorizationCode
	if err := repositories.GetDB().Where("code = ?", code).First(&stored).Error; err != nil {
		t.Fatalf("authorization code not persisted: %v", err)
	}
	if stored.UserId != user.ID {
		t.Errorf("code bound to user %d, want %d", stored.UserId, user.ID)
	}
}

func TestAuthorizeRejectsBadCredentials(t *testing.T) {
	defer repositories.TruncateTestDb()

	createTestUserWithPassword(t, "bob", "correct-password")
	client, _ := createClient("Claude", []string{testRedirectUri})

	recorder := postAuthorize(t, client, testRedirectUri, "bob", "wrong-password", challengeFor("verifier-123"))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for bad credentials, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "Invalid username or password") {
		t.Errorf("expected the login form to re-render with an error message")
	}
}

func TestAuthorizeRejectsNonS256Pkce(t *testing.T) {
	defer repositories.TruncateTestDb()

	createTestUserWithPassword(t, "grace", "password")
	client, _ := createClient("Claude", []string{testRedirectUri})

	form := url.Values{}
	form.Set("username", "grace")
	form.Set("password", "password")
	form.Set("client_id", client.ClientId)
	form.Set("redirect_uri", testRedirectUri)
	form.Set("state", "xyz")
	form.Set("code_challenge", "plaintext-challenge")
	form.Set("code_challenge_method", "plain")

	request := httptest.NewRequest(http.MethodPost, "/oauth/authorize", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	Authorize(recorder, request)

	if recorder.Code != http.StatusFound {
		t.Fatalf("expected a redirect with an error, got %d", recorder.Code)
	}
	location, _ := url.Parse(recorder.Header().Get("Location"))
	if location.Query().Get("error") != "invalid_request" {
		t.Errorf("expected invalid_request error for non-S256 PKCE, got %q", location.Query().Get("error"))
	}
}

func TestAuthorizeRejectsUnregisteredRedirect(t *testing.T) {
	defer repositories.TruncateTestDb()

	createTestUserWithPassword(t, "carol", "password")
	client, _ := createClient("Claude", []string{testRedirectUri})

	recorder := postAuthorize(t, client, "https://evil.example.com/callback", "carol", "password", challengeFor("verifier-123"))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unregistered redirect_uri, got %d", recorder.Code)
	}
}

// postToken submits a form-encoded token request and returns the recorder.
func postToken(t *testing.T, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	Token(recorder, request)
	return recorder
}

func TestTokenAuthorizationCodeGrant(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createTestUserWithPassword(t, "dave", "password")
	client, _ := createClient("Claude", []string{testRedirectUri})

	verifier := "the-real-code-verifier-value-1234567890"
	code, err := createAuthorizationCode(client.ClientId, user.ID, testRedirectUri, challengeFor(verifier), readScope, "")
	if err != nil {
		t.Fatalf("failed to create authorization code: %v", err)
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", client.ClientId)
	form.Set("redirect_uri", testRedirectUri)
	form.Set("code_verifier", verifier)

	recorder := postToken(t, form)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from token endpoint, got %d (%s)", recorder.Code, recorder.Body.String())
	}

	body := decodeJson(t, recorder.Body.Bytes())
	if body["token_type"] != "Bearer" {
		t.Errorf("token_type = %v, want Bearer", body["token_type"])
	}
	if access, ok := body["access_token"].(string); !ok || len(access) == 0 {
		t.Errorf("expected a non-empty access_token")
	}
	if refresh, ok := body["refresh_token"].(string); !ok || len(refresh) == 0 {
		t.Errorf("expected a non-empty refresh_token")
	}

	// A redeemed code must not be exchangeable a second time.
	replay := postToken(t, form)
	if replay.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when replaying a used code, got %d", replay.Code)
	}
}

func TestTokenMismatchedClientDoesNotBurnCode(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createTestUserWithPassword(t, "heidi", "password")
	client, _ := createClient("Claude", []string{testRedirectUri})

	verifier := "the-real-code-verifier-value-1234567890"
	code, err := createAuthorizationCode(client.ClientId, user.ID, testRedirectUri, challengeFor(verifier), readScope, "")
	if err != nil {
		t.Fatalf("failed to create authorization code: %v", err)
	}

	// An attacker redeems with the wrong client_id; this must be rejected and
	// must NOT consume the code.
	badForm := url.Values{}
	badForm.Set("grant_type", "authorization_code")
	badForm.Set("code", code)
	badForm.Set("client_id", "someone-elses-client")
	badForm.Set("redirect_uri", testRedirectUri)
	badForm.Set("code_verifier", verifier)

	if recorder := postToken(t, badForm); recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for client mismatch, got %d", recorder.Code)
	}

	// The legitimate client can still redeem the un-burned code.
	goodForm := url.Values{}
	goodForm.Set("grant_type", "authorization_code")
	goodForm.Set("code", code)
	goodForm.Set("client_id", client.ClientId)
	goodForm.Set("redirect_uri", testRedirectUri)
	goodForm.Set("code_verifier", verifier)

	if recorder := postToken(t, goodForm); recorder.Code != http.StatusOK {
		t.Fatalf("legitimate client could not redeem the code after a mismatched attempt: %d (%s)", recorder.Code, recorder.Body.String())
	}
}

func TestTokenRejectsWrongCodeVerifier(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createTestUserWithPassword(t, "erin", "password")
	client, _ := createClient("Claude", []string{testRedirectUri})

	code, _ := createAuthorizationCode(client.ClientId, user.ID, testRedirectUri, challengeFor("expected-verifier"), readScope, "")

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", client.ClientId)
	form.Set("redirect_uri", testRedirectUri)
	form.Set("code_verifier", "a-different-verifier")

	recorder := postToken(t, form)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for PKCE mismatch, got %d", recorder.Code)
	}
	if decodeJson(t, recorder.Body.Bytes())["error"] != "invalid_grant" {
		t.Errorf("expected invalid_grant error for PKCE mismatch")
	}
}

func TestTokenRefreshGrant(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createTestUserWithPassword(t, "frank", "password")

	_, refreshToken, _, err := services.GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("failed to generate initial tokens: %v", err)
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	recorder := postToken(t, form)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from refresh grant, got %d (%s)", recorder.Code, recorder.Body.String())
	}
	if access, ok := decodeJson(t, recorder.Body.Bytes())["access_token"].(string); !ok || len(access) == 0 {
		t.Errorf("expected a fresh access_token from refresh grant")
	}

	// The refresh token is single use; replay must fail.
	replay := postToken(t, form)
	if replay.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when reusing a refresh token, got %d", replay.Code)
	}
}
