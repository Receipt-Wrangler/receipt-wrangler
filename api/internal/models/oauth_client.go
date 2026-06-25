package models

import "time"

// OAuthClient is an OAuth 2.1 client dynamically registered (RFC 7591) by an
// MCP client such as Claude. Clients are public (no secret) and authenticate
// the authorization-code exchange with PKCE. RedirectUris is a JSON-encoded
// list of the exact redirect URIs the client registered; the authorize and
// token endpoints require an exact match against this list.
type OAuthClient struct {
	ClientId     string    `gorm:"primarykey" json:"clientId"`
	ClientName   string    `json:"clientName"`
	RedirectUris string    `json:"redirectUris"`
	CreatedAt    time.Time `json:"createdAt"`
}
