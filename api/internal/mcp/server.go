// Package mcp exposes Receipt Wrangler as a remote Model Context Protocol
// server over Streamable HTTP, so MCP clients such as Claude can read a user's
// receipts, groups, categories, tags, and dashboards.
//
// Requests are authenticated with the same HS512 access JWTs the REST API
// issues; the OAuth 2.1 handshake that mints them lives in the sibling
// internal/oauth package. Bearer-token verification, the 401 +
// WWW-Authenticate challenge, and propagation of the authenticated user into
// tool handlers are provided by the official go-sdk auth middleware.
package mcp

import (
	"context"
	"fmt"
	"net/http"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/modelcontextprotocol/go-sdk/auth"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpReadScope is the single scope advertised on MCP access tokens.
//
// NOTE: this scope is currently DEAD as an authorization control. The bearer
// middleware requires its presence, but every issued token carries it, so it
// gates nothing. Read-only behavior is guaranteed structurally by the fact
// that registerTools only registers read tools (no write/delete tools), not by
// this scope. The moment a write or delete tool is added, that structural
// guarantee is gone and real per-tool scope enforcement must be introduced.
// See the matching note on readScope in internal/oauth/oauth.go.
const mcpReadScope = "mcp:read"

// claimsKey is the auth.TokenInfo.Extra key under which the verified
// *structs.Claims is stashed for tool handlers to retrieve.
const claimsKey = "claims"

// NewHandler builds the Streamable HTTP MCP handler wrapped in bearer-token
// authentication. On a missing or invalid token it returns 401 with a
// WWW-Authenticate header pointing at the protected-resource metadata, which
// is how an MCP client discovers it must run the OAuth flow.
func NewHandler() http.Handler {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    "receipt-wrangler",
		Title:   "Receipt Wrangler",
		Version: "1.0.0",
	}, nil)

	registerTools(server)

	streamable := mcpsdk.NewStreamableHTTPHandler(func(*http.Request) *mcpsdk.Server {
		return server
	}, nil)

	// The public URL — and therefore the token audience and the resource
	// metadata URL — is a live System Setting, so it is read per request rather
	// than captured once at startup. This keeps the bearer challenge and the
	// audience the validator enforces in sync with the current configuration.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicUrl := services.GetMcpPublicUrl()
		resourceUrl := publicUrl + "/mcp"

		options := &auth.RequireBearerTokenOptions{
			ResourceMetadataURL: publicUrl + "/.well-known/oauth-protected-resource",
			Scopes:              []string{mcpReadScope},
		}

		verifier := func(ctx context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
			return verifyToken(ctx, resourceUrl, token)
		}

		auth.RequireBearerToken(verifier, options)(streamable).ServeHTTP(w, r)
	})
}

// verifyToken validates a bearer access token against the MCP audience and
// returns its claims for downstream tool handlers. A validation failure is
// wrapped in auth.ErrInvalidToken so the middleware responds 401. The validator
// is built per call because the MCP audience is derived from a live setting.
func verifyToken(ctx context.Context, audience string, token string) (*auth.TokenInfo, error) {
	tokenValidator, err := services.InitMcpTokenValidator(audience)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}

	validated, err := tokenValidator.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}

	claims, ok := validated.(*validator.ValidatedClaims).CustomClaims.(*structs.Claims)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected claims type", auth.ErrInvalidToken)
	}

	info := &auth.TokenInfo{
		UserID: utils.UintToString(claims.UserId),
		Scopes: []string{mcpReadScope},
		Extra:  map[string]any{claimsKey: claims},
	}
	if claims.ExpiresAt != nil {
		info.Expiration = claims.ExpiresAt.Time
	}

	return info, nil
}
