package models

import "time"

// OAuthAuthorizationCode is a single-use OAuth 2.1 authorization code issued by
// the authorize endpoint after a user logs in, and redeemed once at the token
// endpoint. It is bound to the client, the user, the exact redirect URI, the
// PKCE code challenge (S256), and the requested resource. Codes are short
// lived (see ExpiresAt) and marked Used on redemption to prevent replay.
type OAuthAuthorizationCode struct {
	Code          string    `gorm:"primarykey" json:"code"`
	ClientId      string    `json:"clientId"`
	UserId        uint      `json:"userId"`
	RedirectUri   string    `json:"redirectUri"`
	CodeChallenge string    `json:"codeChallenge"`
	Scope         string    `json:"scope"`
	Resource      string    `json:"resource"`
	ExpiresAt     time.Time `json:"expiresAt"`
	Used          bool      `json:"used"`
	CreatedAt     time.Time `json:"createdAt"`
}
