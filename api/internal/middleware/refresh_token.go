package middleware

import (
	"context"
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
)

func ValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenValidator, err := services.InitTokenValidator()
		errMessage := "Error refreshing token"

		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_FATAL, err.Error())
			return
		}

		refreshTokenString, err := getRefreshTokenFromRequest(r, w)
		if err != nil {
			utils.WriteCustomErrorResponse(w, errMessage, http.StatusInternalServerError)
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Refresh token not found")
			return
		}

		refreshToken, err := tokenValidator.ValidateToken(context.TODO(), refreshTokenString)
		if err != nil {
			utils.WriteCustomErrorResponse(w, errMessage, http.StatusInternalServerError)
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "refreshToken", refreshToken)
		ctx = context.WithValue(ctx, "refreshTokenString", refreshTokenString)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RevokeRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := repositories.GetDB()
		err := error(nil)
		errMessage := "Error refreshing token"

		refreshTokenString := r.Context().Value("refreshTokenString")
		if refreshTokenString == nil {
			refreshTokenString, err = getRefreshTokenFromRequest(r, w)
			if err != nil {
				utils.WriteCustomErrorResponse(w, errMessage, http.StatusInternalServerError)
				logging.LogStd(logging.LOG_LEVEL_ERROR, "Refresh token not found")
				return
			}
		}

		hashTokenString := utils.Sha256Hash([]byte(refreshTokenString.(string)))

		// Atomically mark the stored token used via the WHERE clause so two
		// concurrent refreshes of the same token cannot both succeed. A read
		// followed by a separate update lets both racers observe is_used = false
		// and rotate, which defeats replay detection; it also logs the loser out
		// of every tab once the rotation it lost lands. Mirrors the same fix in
		// oauth.rotateRefreshToken. A zero row count means the token is unknown
		// or already used — indistinguishable on purpose, and both already
		// produced this identical response.
		result := db.Model(&models.RefreshToken{}).
			Where("token = ? AND is_used = ?", hashTokenString, false).
			Update("is_used", true)
		if result.Error != nil {
			utils.WriteCustomErrorResponse(w, errMessage, http.StatusInternalServerError)
			logging.LogStd(logging.LOG_LEVEL_ERROR, result.Error.Error())
			return
		}

		if result.RowsAffected == 0 {
			emptyAccessTokenCookie := services.GetEmptyAccessTokenCookie()
			emptyRefreshTokenCookie := services.GetEmptyRefreshTokenCookie()

			http.SetCookie(w, &emptyAccessTokenCookie)
			http.SetCookie(w, &emptyRefreshTokenCookie)

			utils.WriteCustomErrorResponse(w, errMessage, http.StatusInternalServerError)
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Refresh token is invalid or has already been used.")

			return
		}

		next.ServeHTTP(w, r)
	})
}

func getRefreshTokenFromRequest(r *http.Request, w http.ResponseWriter) (string, error) {
	if utils.IsMobileApp(r) {
		var command commands.LogoutCommand
		err := command.LoadDataFromRequest(w, r)
		if err != nil {
			return "", err
		}

		return command.RefreshToken, nil
	} else {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			return "", err
		}

		return cookie.Value, nil
	}
}
