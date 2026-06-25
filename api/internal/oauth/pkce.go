package oauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

// verifyPkceS256 reports whether codeVerifier matches the stored S256
// codeChallenge per RFC 7636: BASE64URL-WITHOUT-PADDING(SHA256(verifier)).
// Only the S256 method is supported; "plain" is intentionally rejected by
// callers before reaching here. The comparison is constant time.
func verifyPkceS256(codeVerifier string, codeChallenge string) bool {
	if len(codeVerifier) == 0 || len(codeChallenge) == 0 {
		return false
	}

	sum := sha256.Sum256([]byte(codeVerifier))
	computed := base64.RawURLEncoding.EncodeToString(sum[:])

	return subtle.ConstantTimeCompare([]byte(computed), []byte(codeChallenge)) == 1
}
