package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

const (
	// bindingCookieName holds the browser-binding secret for the desktop leg.
	bindingCookieName = "oidc_session"

	// bindingCookiePath scopes the cookie to the OIDC routes: nothing else needs
	// it, and it should not ride along on every API call.
	bindingCookiePath = "/api/oidc"

	// mobileCallbackScheme is the app's private-use URL scheme, reverse-DNS of a
	// domain the project controls per RFC 8252 section 7.1. It carries a one-time,
	// PKCE-bound code -- never a token.
	mobileCallbackScheme = "io.receiptwrangler://oidc"

	// desktopCallbackPath is where the browser lands after a successful login. It
	// is a fixed, relative path: no redirect target is ever taken from the request,
	// which is what keeps this flow free of an open-redirect hole.
	desktopCallbackPath = "/auth/callback"

	// desktopLoginPath is where a failed login lands.
	desktopLoginPath = "/auth/login"

	// desktopProfilePath is where a link flow returns to.
	desktopProfilePath = "/settings/user-profile/view"

	exchangeTimeout = 15 * time.Second
)

// Error codes handed back to a client. A small fixed vocabulary -- an upstream
// identity provider's own error text is never echoed through.
const (
	errUnknownProvider = "unknown_provider"
	errInvalidRequest  = "invalid_request"
	errInvalidState    = "invalid_state"
	errNonceMismatch   = "nonce_mismatch"
	errNoIdToken       = "no_id_token"
	errProviderError   = "provider_error"
	errNoAccount       = "no_account"
	errAccountExists   = "account_exists"
	errAlreadyLinked   = "already_linked"
	errServerError     = "server_error"
)

// Login starts an OIDC login. Unauthenticated.
func Login(w http.ResponseWriter, r *http.Request) {
	startFlow(w, r, nil)
}

// LinkStart starts a "connect account" flow for the authenticated caller. It is
// the same flow as Login, differing only in that the resulting session carries
// the caller's user id -- so the callback links instead of guessing.
func LinkStart(w http.ResponseWriter, r *http.Request) {
	claims := structs.GetClaims(r)
	userId := claims.UserId

	startFlow(w, r, &userId)
}

func startFlow(w http.ResponseWriter, r *http.Request, linkUserId *uint) {
	name := chi.URLParam(r, "name")

	providerRow, err := repositories.NewOidcProviderRepository(nil).GetEnabledOidcProviderByName(name)
	if err != nil {
		redirectWithError(w, r, models.OidcClientDesktop, errUnknownProvider)
		return
	}

	clientType := models.OidcClientDesktop
	mobileChallenge := ""

	if r.URL.Query().Get("client") == models.OidcClientMobile {
		clientType = models.OidcClientMobile
		mobileChallenge = strings.TrimSpace(r.URL.Query().Get("codeChallenge"))

		// The mobile leg has no browser cookie to bind to, so the app's own PKCE
		// challenge is the only thing tying the exchange back to the app that
		// started the flow. Without it the handoff code would be bearer-only.
		if len(mobileChallenge) == 0 {
			redirectWithError(w, r, clientType, errInvalidRequest)
			return
		}
	}

	discovered, err := GetProvider(providerRow)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "OIDC discovery failed for provider "+providerRow.Name+": "+err.Error())
		redirectWithError(w, r, clientType, errProviderError)
		return
	}

	config, err := buildOauthConfig(providerRow, discovered, services.BuildOidcRedirectUri(providerRow.Name))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to build OIDC config for provider "+providerRow.Name+": "+err.Error())
		redirectWithError(w, r, clientType, errServerError)
		return
	}

	verifier := oauth2.GenerateVerifier()

	created, err := createAuthSession(newAuthSessionParams{
		ProviderId:          providerRow.ID,
		ClientType:          clientType,
		MobileCodeChallenge: mobileChallenge,
		LinkUserId:          linkUserId,
		CodeVerifier:        verifier,
	})
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to create OIDC auth session: "+err.Error())
		redirectWithError(w, r, clientType, errServerError)
		return
	}

	if clientType != models.OidcClientMobile {
		http.SetCookie(w, buildBindingCookie(created.Binding))
	}

	authUrl := config.AuthCodeURL(
		created.State,
		oidc.Nonce(created.Nonce),
		oauth2.S256ChallengeOption(verifier),
	)

	http.Redirect(w, r, authUrl, http.StatusFound)
}

// Callback is where the identity provider returns the user. Every check below is
// load-bearing; they run in this order and fail closed.
func Callback(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	providerRow, err := repositories.NewOidcProviderRepository(nil).GetEnabledOidcProviderByName(name)
	if err != nil {
		redirectWithError(w, r, models.OidcClientDesktop, errUnknownProvider)
		return
	}

	// 1. The identity provider itself refused. Nothing is consumed; the pending
	// session simply expires. Never echo the upstream text.
	if len(r.URL.Query().Get("error")) > 0 {
		clearBindingCookie(w)
		redirectWithError(w, r, models.OidcClientDesktop, errProviderError)
		return
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if len(state) == 0 || len(code) == 0 {
		clearBindingCookie(w)
		redirectWithError(w, r, models.OidcClientDesktop, errInvalidRequest)
		return
	}

	// 2. Claim the session ATOMICALLY, before any network call. A replayed state
	// therefore cannot even reach the identity provider, let alone mint a second
	// session. Unknown, expired and already-used states are indistinguishable.
	session, claimed, err := consumeAuthSession(state)
	clearBindingCookie(w)

	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to consume OIDC auth session: "+err.Error())
		redirectWithError(w, r, models.OidcClientDesktop, errServerError)
		return
	}

	if !claimed {
		redirectWithError(w, r, models.OidcClientDesktop, errInvalidState)
		return
	}

	clientType := session.ClientType

	// 3. The state must belong to THIS provider. Otherwise a state minted for a
	// permissive provider could be redeemed at a stricter one's callback.
	if session.OidcProviderId != providerRow.ID {
		redirectWithError(w, r, clientType, errInvalidState)
		return
	}

	// 4. Browser binding -- the login-CSRF defense. Without it an attacker can
	// start a flow, harvest the state and code, and plant them in a victim's
	// browser, silently signing the victim into the ATTACKER's account. The
	// attacker cannot produce the victim's cookie.
	if len(session.BindingHash) > 0 && !bindingMatches(r, session.BindingHash) {
		redirectWithError(w, r, clientType, errInvalidState)
		return
	}

	verifier, err := decryptVerifier(session)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to decrypt OIDC code verifier: "+err.Error())
		redirectWithError(w, r, clientType, errServerError)
		return
	}

	discovered, err := GetProvider(providerRow)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "OIDC discovery failed on callback: "+err.Error())
		redirectWithError(w, r, clientType, errProviderError)
		return
	}

	config, err := buildOauthConfig(providerRow, discovered, services.BuildOidcRedirectUri(providerRow.Name))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to build OIDC config on callback: "+err.Error())
		redirectWithError(w, r, clientType, errServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), exchangeTimeout)
	defer cancel()
	ctx = oidc.ClientContext(ctx, &http.Client{Timeout: exchangeTimeout})

	// 5. Exchange the code, proving possession of our PKCE verifier.
	token, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "OIDC code exchange failed: "+err.Error())
		redirectWithError(w, r, clientType, errProviderError)
		return
	}

	// 6. No ID token means this was an OAuth response, not an OIDC one -- there is
	// no verified identity to act on. Hard fail rather than falling back to
	// userinfo, which is not signed.
	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok || len(rawIdToken) == 0 {
		redirectWithError(w, r, clientType, errNoIdToken)
		return
	}

	// 7. Verify signature (via JWKS), issuer, audience and expiry.
	idToken, err := discovered.Verifier(&oidc.Config{ClientID: providerRow.ClientId}).Verify(ctx, rawIdToken)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "OIDC ID token verification failed: "+err.Error())
		redirectWithError(w, r, clientType, errProviderError)
		return
	}

	// 8. Verify the nonce OURSELVES. go-oidc documents that Verify does NOT do
	// nonce validation -- "that is the caller's responsibility" -- and skipping it
	// re-opens ID token replay. The empty check matters independently: an empty
	// nonce hashes to a fixed value, and must never be allowed to match.
	if len(idToken.Nonce) == 0 || !utils.SecureCompare(hashSecret(idToken.Nonce), session.NonceHash) {
		redirectWithError(w, r, clientType, errNonceMismatch)
		return
	}

	claims, err := extractClaims(ctx, discovered, config, token, idToken)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to read OIDC claims: "+err.Error())
		redirectWithError(w, r, clientType, errProviderError)
		return
	}

	if len(claims.Subject) == 0 {
		redirectWithError(w, r, clientType, errProviderError)
		return
	}

	if session.IsLink() {
		finishLink(w, r, providerRow, claims, session)
		return
	}

	finishLogin(w, r, providerRow, claims, session)
}

// extractClaims decodes the ID token's claims, falling back to the userinfo
// endpoint only to fill in an email the ID token omitted. Identity itself always
// comes from the signed ID token.
func extractClaims(
	ctx context.Context,
	discovered *oidc.Provider,
	config *oauth2.Config,
	token *oauth2.Token,
	idToken *oidc.IDToken,
) (idTokenClaims, error) {
	var raw json.RawMessage
	if err := idToken.Claims(&raw); err != nil {
		return idTokenClaims{}, err
	}

	claims, err := decodeIdTokenClaims(raw)
	if err != nil {
		return idTokenClaims{}, err
	}

	if len(claims.Email) > 0 || len(discovered.UserInfoEndpoint()) == 0 {
		return claims, nil
	}

	// Best effort: several providers (Twitch notably) will not put email in the ID
	// token. Email is display-only here, so a failure is logged and ignored rather
	// than failing an otherwise valid login.
	info, err := discovered.UserInfo(ctx, config.TokenSource(ctx, token))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "OIDC userinfo lookup failed: "+err.Error())
		return claims, nil
	}

	if info.Subject == claims.Subject {
		claims.Email = info.Email
		claims.EmailVerified = tolerantBool(info.EmailVerified)
	}

	return claims, nil
}

func finishLogin(
	w http.ResponseWriter,
	r *http.Request,
	providerRow models.OidcProvider,
	claims idTokenClaims,
	session models.OidcAuthSession,
) {
	user, err := resolveUser(providerRow, claims)
	if err != nil {
		redirectWithError(w, r, session.ClientType, resolutionErrorCode(err))
		return
	}

	issueSession(w, r, user.ID, session)
}

func finishLink(
	w http.ResponseWriter,
	r *http.Request,
	providerRow models.OidcProvider,
	claims idTokenClaims,
	session models.OidcAuthSession,
) {
	err := resolveLink(providerRow, claims, *session.LinkUserId)
	if err != nil {
		http.Redirect(w, r, desktopProfilePath+"?tab=user-profile&oidcError="+resolutionErrorCode(err), http.StatusFound)
		return
	}

	http.Redirect(w, r, desktopProfilePath+"?tab=user-profile&oidcLinked="+url.QueryEscape(providerRow.Name), http.StatusFound)
}

// issueSession mints Receipt Wrangler tokens and hands them to the client in the
// way that client can actually receive them.
func issueSession(w http.ResponseWriter, r *http.Request, userId uint, session models.OidcAuthSession) {
	if session.ClientType == models.OidcClientMobile {
		// Deliberately NOT the tokens themselves: a private-use scheme is
		// unverifiable on Android, so any installed app can read this redirect. The
		// code is single-use, short-lived and PKCE-bound, so intercepting it is
		// useless without the verifier the app kept.
		code, err := createExchangeCode(userId, session.MobileCodeChallenge)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to create OIDC exchange code: "+err.Error())
			redirectWithError(w, r, session.ClientType, errServerError)
			return
		}

		http.Redirect(w, r, mobileCallbackScheme+"?code="+url.QueryEscape(code), http.StatusFound)
		return
	}

	jwt, refreshToken, _, err := services.GenerateJWT(userId)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to mint tokens after OIDC login: "+err.Error())
		redirectWithError(w, r, session.ClientType, errServerError)
		return
	}

	_, err = repositories.NewUserRepository(nil).UpdateUserLastLoginDate(userId)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Failed to stamp last login date after OIDC login: "+err.Error())
	}

	accessTokenCookie, refreshTokenCookie := services.BuildTokenCookies(jwt, refreshToken)
	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	http.Redirect(w, r, desktopCallbackPath, http.StatusFound)
}

func resolutionErrorCode(err error) string {
	switch {
	case errors.Is(err, ErrNoAccount), errors.Is(err, ErrUserIsDummy):
		return errNoAccount
	case errors.Is(err, ErrAccountExists):
		return errAccountExists
	case errors.Is(err, ErrAlreadyLinked), errors.Is(err, ErrIdentityLinkedElsewhere):
		return errAlreadyLinked
	default:
		logging.LogStd(logging.LOG_LEVEL_ERROR, "OIDC identity resolution failed: "+err.Error())
		return errServerError
	}
}

func redirectWithError(w http.ResponseWriter, r *http.Request, clientType string, code string) {
	if clientType == models.OidcClientMobile {
		http.Redirect(w, r, mobileCallbackScheme+"?error="+url.QueryEscape(code), http.StatusFound)
		return
	}

	http.Redirect(w, r, desktopLoginPath+"?oidcError="+url.QueryEscape(code), http.StatusFound)
}

func buildBindingCookie(value string) *http.Cookie {
	secure := false
	if env.GetDeployEnv() == "dev" {
		secure = true
	}

	return &http.Cookie{
		Name:     bindingCookieName,
		Value:    value,
		HttpOnly: true,
		Path:     bindingCookiePath,
		MaxAge:   int(authSessionTTL.Seconds()),
		Secure:   secure,
		// Lax, NOT Strict. BuildTokenCookies uses Strict in production, but the
		// callback arrives as a cross-site top-level GET from the identity provider:
		// Strict would drop this cookie on exactly the request that needs it, and the
		// flow would fail only in production.
		SameSite: http.SameSiteLaxMode,
	}
}

func clearBindingCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     bindingCookieName,
		Value:    "",
		HttpOnly: true,
		Path:     bindingCookiePath,
		MaxAge:   -1,
	})
}

func bindingMatches(r *http.Request, bindingHash string) bool {
	cookie, err := r.Cookie(bindingCookieName)
	if err != nil || len(cookie.Value) == 0 {
		return false
	}

	return utils.SecureCompare(hashSecret(cookie.Value), bindingHash)
}
