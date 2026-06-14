package oauth

import (
	"encoding/json"
	"errors"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"time"

	"gorm.io/gorm"
)

// createClient persists a dynamically registered OAuth client and returns it.
func createClient(clientName string, redirectUris []string) (models.OAuthClient, error) {
	clientId, err := randomToken(24)
	if err != nil {
		return models.OAuthClient{}, err
	}

	encodedUris, err := json.Marshal(redirectUris)
	if err != nil {
		return models.OAuthClient{}, err
	}

	client := models.OAuthClient{
		ClientId:     clientId,
		ClientName:   clientName,
		RedirectUris: string(encodedUris),
	}

	if err := repositories.GetDB().Create(&client).Error; err != nil {
		return models.OAuthClient{}, err
	}

	return client, nil
}

// getClient loads a registered client by id. A missing client is reported as
// gorm.ErrRecordNotFound so callers can distinguish "unknown client" from a
// database failure.
func getClient(clientId string) (models.OAuthClient, error) {
	var client models.OAuthClient
	err := repositories.GetDB().Where("client_id = ?", clientId).First(&client).Error
	return client, err
}

// clientAllowsRedirect reports whether redirectUri exactly matches one of the
// client's registered redirect URIs. Exact match (no prefix/substring) is
// required by OAuth 2.1 to prevent open-redirect attacks.
func clientAllowsRedirect(client models.OAuthClient, redirectUri string) bool {
	var registered []string
	if err := json.Unmarshal([]byte(client.RedirectUris), &registered); err != nil {
		return false
	}

	for _, uri := range registered {
		if uri == redirectUri {
			return true
		}
	}

	return false
}

// createAuthorizationCode issues and persists a single-use authorization code
// bound to the authenticated user, client, redirect URI, PKCE challenge, and
// requested resource.
func createAuthorizationCode(
	clientId string,
	userId uint,
	redirectUri string,
	codeChallenge string,
	scope string,
	resource string,
) (string, error) {
	code, err := randomToken(32)
	if err != nil {
		return "", err
	}

	authCode := models.OAuthAuthorizationCode{
		Code:          code,
		ClientId:      clientId,
		UserId:        userId,
		RedirectUri:   redirectUri,
		CodeChallenge: codeChallenge,
		Scope:         scope,
		Resource:      resource,
		ExpiresAt:     time.Now().Add(authCodeTTLSeconds * time.Second),
	}

	if err := repositories.GetDB().Create(&authCode).Error; err != nil {
		return "", err
	}

	return code, nil
}

// consumeAuthorizationCode atomically loads and marks an authorization code as
// used. It returns an error if the code is unknown, already used, or expired,
// so a redeemed or replayed code can never be exchanged twice.
func consumeAuthorizationCode(code string) (models.OAuthAuthorizationCode, error) {
	var authCode models.OAuthAuthorizationCode

	err := repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("code = ?", code).First(&authCode).Error; err != nil {
			return err
		}

		if authCode.Used {
			return errors.New("authorization code already used")
		}

		if time.Now().After(authCode.ExpiresAt) {
			return errors.New("authorization code expired")
		}

		return tx.Model(&authCode).Update("used", true).Error
	})
	if err != nil {
		return models.OAuthAuthorizationCode{}, err
	}

	return authCode, nil
}
