package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func adminContext() *validator.ValidatedClaims {
	return &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2, UserRole: models.ADMIN}}
}

func userContext() *validator.ValidatedClaims {
	return &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2, UserRole: models.USER}}
}

func TestShouldCreateAppRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	role := structs.RoleView{}

	reader := strings.NewReader(`{"name": "App Role", "description": "desc", "scope": "APP", "permissions": ["app.users.create", "app.users.read"]}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, adminContext())
	r = r.WithContext(newContext)

	CreateRole(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	err := json.Unmarshal(w.Body.Bytes(), &role)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if role.Scope != permissions.ScopeApp {
		utils.PrintTestError(t, role.Scope, permissions.ScopeApp)
	}

	if len(role.Permissions) != 2 {
		utils.PrintTestError(t, len(role.Permissions), 2)
	}
}

func TestShouldCreateGroupRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	role := structs.RoleView{}

	reader := strings.NewReader(`{"name": "Group Role", "description": "desc", "scope": "GROUP", "permissions": ["group.receipts.create"]}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, adminContext())
	r = r.WithContext(newContext)

	CreateRole(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	err := json.Unmarshal(w.Body.Bytes(), &role)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if role.Scope != permissions.ScopeGroup {
		utils.PrintTestError(t, role.Scope, permissions.ScopeGroup)
	}
}

func TestShouldNotCreateRoleDueToValidation(t *testing.T) {
	defer repositories.TruncateTestDb()

	reader := strings.NewReader(`{"name": "", "scope": "APP", "permissions": []}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, adminContext())
	r = r.WithContext(newContext)

	CreateRole(w, r)

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestShouldNotCreateRoleDueToRole(t *testing.T) {
	defer repositories.TruncateTestDb()

	reader := strings.NewReader(`{"name": "App Role", "scope": "APP", "permissions": ["app.users.create"]}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, userContext())
	r = r.WithContext(newContext)

	CreateRole(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestShouldGetRoles(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()
	roles := make([]structs.RoleView, 0)

	reader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, adminContext())
	r = r.WithContext(newContext)

	GetRoles(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	err := json.Unmarshal(w.Body.Bytes(), &roles)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roles) != 2 {
		utils.PrintTestError(t, len(roles), 2)
	}
}
