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
	config "receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/modelcontextprotocol/go-sdk/auth"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpReadScope is the single read-only scope granted to MCP access tokens.
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

	// Build the JWT validator once rather than per request; it is immutable and
	// sits on the auth critical path for every MCP call.
	tokenValidator, err := services.InitTokenValidator()
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_FATAL, "Failed to initialize MCP token validator: "+err.Error())
	}

	options := &auth.RequireBearerTokenOptions{
		ResourceMetadataURL: config.GetMcpPublicUrl() + "/.well-known/oauth-protected-resource",
		Scopes:              []string{mcpReadScope},
	}

	verifier := func(ctx context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
		return verifyToken(ctx, tokenValidator, token)
	}

	return auth.RequireBearerToken(verifier, options)(streamable)
}

// verifyToken validates a bearer access token with the provided validator and
// returns its claims for downstream tool handlers. A validation failure is
// wrapped in auth.ErrInvalidToken so the middleware responds 401.
func verifyToken(ctx context.Context, tokenValidator *validator.Validator, token string) (*auth.TokenInfo, error) {
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
