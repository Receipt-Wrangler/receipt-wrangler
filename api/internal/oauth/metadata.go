package oauth

import (
	"net/http"
	"receipt-wrangler/api/internal/services"
)

// issuer returns the OAuth issuer / public origin (no trailing slash), read
// live from System Settings so it tracks runtime configuration changes.
func issuer() string {
	return services.GetMcpPublicUrl()
}

// mcpResourceUrl is the protected resource identifier advertised to clients,
// echoed back as the RFC 8707 resource indicator, and bound as the audience of
// MCP tokens.
func mcpResourceUrl() string {
	return services.GetMcpResourceUrl()
}

// ProtectedResourceMetadata serves RFC 9728 OAuth 2.0 Protected Resource
// Metadata. It tells the MCP client which authorization server protects the
// /mcp resource so the client can discover the auth endpoints.
func ProtectedResourceMetadata(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, map[string]interface{}{
		"resource":                 mcpResourceUrl(),
		"authorization_servers":    []string{issuer()},
		"scopes_supported":         []string{readScope},
		"bearer_methods_supported": []string{"header"},
	})
}

// AuthorizationServerMetadata serves RFC 8414 OAuth 2.0 Authorization Server
// Metadata describing the endpoints and capabilities of this self-hosted
// authorization server.
func AuthorizationServerMetadata(w http.ResponseWriter, r *http.Request) {
	base := issuer()

	writeJson(w, http.StatusOK, map[string]interface{}{
		"issuer":                                base,
		"authorization_endpoint":                base + "/oauth/authorize",
		"token_endpoint":                        base + "/oauth/token",
		"registration_endpoint":                 base + "/oauth/register",
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token"},
		"code_challenge_methods_supported":      []string{"S256"},
		"token_endpoint_auth_methods_supported": []string{"none"},
		"scopes_supported":                      []string{readScope},
	})
}
