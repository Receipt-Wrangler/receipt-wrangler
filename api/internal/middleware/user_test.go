package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func runValidateUserData(command commands.SignUpCommand, roleRequired bool) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, "/api/user", nil)
	r = r.WithContext(context.WithValue(r.Context(), "user", command))
	w := httptest.NewRecorder()

	mw := ValidateUserData(roleRequired)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, r)
	return w
}

// The desktop now supplies the modern app role id instead of the legacy enum.
func TestValidateUserDataAllowsModernAppRoleId(t *testing.T) {
	defer repositories.TruncateTable(repositories.GetDB(), "users")

	appRoleId := uint(1)
	w := runValidateUserData(commands.SignUpCommand{
		Username:    "modern",
		Password:    "password",
		DisplayName: "Modern",
		AppRoleID:   &appRoleId,
	}, true)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}
}

// The legacy enum is still accepted (e.g. the e2e seed script posts userRole).
func TestValidateUserDataAllowsLegacyUserRole(t *testing.T) {
	defer repositories.TruncateTable(repositories.GetDB(), "users")

	w := runValidateUserData(commands.SignUpCommand{
		Username:    "legacy",
		Password:    "password",
		DisplayName: "Legacy",
		UserRole:    models.USER,
	}, true)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}
}

// Neither a modern app role id nor a legacy enum is a 400.
func TestValidateUserDataRejectsMissingRole(t *testing.T) {
	defer repositories.TruncateTable(repositories.GetDB(), "users")

	w := runValidateUserData(commands.SignUpCommand{
		Username:    "norole",
		Password:    "password",
		DisplayName: "No Role",
	}, true)

	if w.Result().StatusCode != http.StatusBadRequest {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusBadRequest)
	}
}
