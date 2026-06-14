package oauth

import (
	"net/http"
	"net/url"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/services"
)

// AuthorizeForm handles GET /oauth/authorize. It validates the OAuth request,
// and once the client and redirect URI are trusted, renders the login page.
// Parameter errors discovered before the redirect URI is trusted render an
// error page; errors after that point redirect back to the client per OAuth.
func AuthorizeForm(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	client, redirectUri, ok := validateClientAndRedirect(w, query.Get("client_id"), query.Get("redirect_uri"))
	if !ok {
		return
	}

	state := query.Get("state")

	if query.Get("response_type") != "code" {
		redirectWithError(w, r, redirectUri, state, "unsupported_response_type", "Only the authorization code flow is supported")
		return
	}

	if query.Get("code_challenge_method") != "S256" || len(query.Get("code_challenge")) == 0 {
		redirectWithError(w, r, redirectUri, state, "invalid_request", "PKCE with code_challenge_method=S256 is required")
		return
	}

	if resource := query.Get("resource"); len(resource) > 0 && resource != mcpResourceUrl() {
		redirectWithError(w, r, redirectUri, state, "invalid_target", "Unknown resource indicator")
		return
	}

	renderLogin(w, http.StatusOK, loginFormData{
		ClientId:      client.ClientId,
		ClientName:    client.ClientName,
		RedirectUri:   redirectUri,
		State:         state,
		Scope:         scopeOrDefault(query.Get("scope")),
		CodeChallenge: query.Get("code_challenge"),
		Resource:      query.Get("resource"),
	})
}

// Authorize handles POST /oauth/authorize: it authenticates the submitted
// credentials against the existing user store and, on success, issues a
// single-use authorization code and redirects back to the client.
func Authorize(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderError(w, http.StatusBadRequest, "Malformed request")
		return
	}

	client, redirectUri, ok := validateClientAndRedirect(w, r.PostForm.Get("client_id"), r.PostForm.Get("redirect_uri"))
	if !ok {
		return
	}

	state := r.PostForm.Get("state")
	codeChallenge := r.PostForm.Get("code_challenge")
	resource := r.PostForm.Get("resource")
	scope := scopeOrDefault(r.PostForm.Get("scope"))

	if r.PostForm.Get("code_challenge_method") == "plain" || len(codeChallenge) == 0 {
		redirectWithError(w, r, redirectUri, state, "invalid_request", "PKCE with code_challenge_method=S256 is required")
		return
	}

	formData := loginFormData{
		ClientId:      client.ClientId,
		ClientName:    client.ClientName,
		RedirectUri:   redirectUri,
		State:         state,
		Scope:         scope,
		CodeChallenge: codeChallenge,
		Resource:      resource,
	}

	user, _, err := services.LoginUser(commands.LoginCommand{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	})
	if err != nil {
		formData.ErrorMessage = "Invalid username or password."
		renderLogin(w, http.StatusUnauthorized, formData)
		return
	}

	code, err := createAuthorizationCode(client.ClientId, user.ID, redirectUri, codeChallenge, scope, resource)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to create authorization code: "+err.Error())
		redirectWithError(w, r, redirectUri, state, "server_error", "Failed to issue authorization code")
		return
	}

	redirectWithParams(w, r, redirectUri, map[string]string{
		"code":  code,
		"state": state,
	})
}

// validateClientAndRedirect loads the client and verifies the redirect URI is
// registered. On failure it writes an error page (the redirect URI is not
// trusted yet, so we must not redirect) and reports ok=false.
func validateClientAndRedirect(w http.ResponseWriter, clientId string, redirectUri string) (client clientResult, resolvedRedirect string, ok bool) {
	if len(clientId) == 0 || len(redirectUri) == 0 {
		renderError(w, http.StatusBadRequest, "Missing client_id or redirect_uri")
		return clientResult{}, "", false
	}

	dbClient, err := getClient(clientId)
	if err != nil {
		renderError(w, http.StatusBadRequest, "Unknown client")
		return clientResult{}, "", false
	}

	if !clientAllowsRedirect(dbClient, redirectUri) {
		renderError(w, http.StatusBadRequest, "redirect_uri is not registered for this client")
		return clientResult{}, "", false
	}

	return clientResult{ClientId: dbClient.ClientId, ClientName: dbClient.ClientName}, redirectUri, true
}

// clientResult is the minimal client view the authorize handlers need.
type clientResult struct {
	ClientId   string
	ClientName string
}

func scopeOrDefault(scope string) string {
	if len(scope) == 0 {
		return readScope
	}
	return scope
}

func renderLogin(w http.ResponseWriter, status int, data loginFormData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := loginTemplate.Execute(w, data); err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to render login template: "+err.Error())
	}
}

func renderError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

// redirectWithError redirects back to the client's redirect URI carrying an
// OAuth error (RFC 6749 §4.1.2.1).
func redirectWithError(w http.ResponseWriter, r *http.Request, redirectUri string, state string, oauthError string, description string) {
	redirectWithParams(w, r, redirectUri, map[string]string{
		"error":             oauthError,
		"error_description": description,
		"state":             state,
	})
}

// redirectWithParams issues a 302 to redirectUri with the given query
// parameters appended, preserving any existing query on the URI.
func redirectWithParams(w http.ResponseWriter, r *http.Request, redirectUri string, params map[string]string) {
	parsed, err := url.Parse(redirectUri)
	if err != nil {
		renderError(w, http.StatusBadRequest, "Invalid redirect_uri")
		return
	}

	query := parsed.Query()
	for key, value := range params {
		if len(value) > 0 {
			query.Set(key, value)
		}
	}
	parsed.RawQuery = query.Encode()

	http.Redirect(w, r, parsed.String(), http.StatusFound)
}
