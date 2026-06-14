package oauth

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/logging"
)

// registrationRequest is the subset of the RFC 7591 Dynamic Client
// Registration request we honor. Clients are treated as public (PKCE, no
// secret), so credential-related fields are ignored.
type registrationRequest struct {
	RedirectUris []string `json:"redirect_uris"`
	ClientName   string   `json:"client_name"`
}

// Register implements RFC 7591 Dynamic Client Registration. An MCP client such
// as Claude self-registers by posting its redirect URIs; we persist them and
// return a generated client_id. The client is public and authenticates the
// code exchange with PKCE rather than a secret.
func Register(w http.ResponseWriter, r *http.Request) {
	var request registrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_client_metadata", "Malformed registration request body")
		return
	}

	if len(request.RedirectUris) == 0 {
		writeOAuthError(w, http.StatusBadRequest, "invalid_redirect_uri", "At least one redirect_uri is required")
		return
	}

	client, err := createClient(request.ClientName, request.RedirectUris)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to register OAuth client: "+err.Error())
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "Failed to register client")
		return
	}

	writeJson(w, http.StatusCreated, map[string]interface{}{
		"client_id":                  client.ClientId,
		"client_name":                client.ClientName,
		"redirect_uris":              request.RedirectUris,
		"token_endpoint_auth_method": "none",
		"grant_types":                []string{"authorization_code", "refresh_token"},
		"response_types":             []string{"code"},
	})
}
