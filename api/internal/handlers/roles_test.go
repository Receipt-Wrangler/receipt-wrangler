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
	"github.com/go-chi/chi/v5"
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

func updateRoleRequest(roleId string, body string, claims *validator.ValidatedClaims) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/api", strings.NewReader(body))

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("roleId", roleId)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claims))

	UpdateRole(w, r)
	return w
}

func TestShouldUpdateAppRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()
	role := structs.RoleView{}

	body := `{"name": "Updated App Role", "description": "updated", "scope": "APP", "permissions": ["app.users.read"]}`
	w := updateRoleRequest("1", body, adminContext())

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	err := json.Unmarshal(w.Body.Bytes(), &role)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if role.Name != "Updated App Role" {
		utils.PrintTestError(t, role.Name, "Updated App Role")
	}

	if role.Scope != permissions.ScopeApp {
		utils.PrintTestError(t, role.Scope, permissions.ScopeApp)
	}

	if len(role.Permissions) != 1 || role.Permissions[0] != permissions.AppUsersRead {
		utils.PrintTestError(t, role.Permissions, []string{permissions.AppUsersRead})
	}
}

func TestShouldUpdateGroupRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()
	role := structs.RoleView{}

	body := `{"name": "Updated Group Role", "description": "updated", "scope": "GROUP", "permissions": ["group.receipts.read"]}`
	w := updateRoleRequest("1", body, adminContext())

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	err := json.Unmarshal(w.Body.Bytes(), &role)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if role.Name != "Updated Group Role" {
		utils.PrintTestError(t, role.Name, "Updated Group Role")
	}

	if role.Scope != permissions.ScopeGroup {
		utils.PrintTestError(t, role.Scope, permissions.ScopeGroup)
	}
}

func TestShouldNotUpdateRoleDueToValidation(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()

	body := `{"name": "", "scope": "APP", "permissions": []}`
	w := updateRoleRequest("1", body, adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestShouldNotUpdateRoleDueToRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()

	body := `{"name": "Updated App Role", "scope": "APP", "permissions": ["app.users.read"]}`
	w := updateRoleRequest("1", body, userContext())

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestShouldNotChangeRoleType(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)
	created, err := roleRepository.CreateAppRole("App Role", "desc", []string{permissions.AppUsersCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	body := `{"name": "App Role", "scope": "GROUP", "permissions": ["group.receipts.create"]}`
	w := updateRoleRequest("1", body, adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
		return
	}

	// The app role must be untouched.
	unchanged, err := roleRepository.GetAppRoleById(created.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if unchanged.Name != "App Role" {
		utils.PrintTestError(t, unchanged.Name, "App Role")
	}
}

func TestShouldNotUpdateSystemRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()
	systemRole := models.AppRole{
		Name:        "System Role",
		Description: "system role",
		IsSystem:    true,
		Permissions: []models.AppRolePermission{
			{Permission: permissions.AppUsersCreate},
		},
	}
	if err := db.Create(&systemRole).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	body := `{"name": "Renamed System Role", "scope": "APP", "permissions": ["app.users.read"]}`
	w := updateRoleRequest("1", body, adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestShouldReturnNotFoundForMissingRole(t *testing.T) {
	defer repositories.TruncateTestDb()

	body := `{"name": "Missing Role", "scope": "APP", "permissions": ["app.users.read"]}`
	w := updateRoleRequest("999", body, adminContext())

	if w.Result().StatusCode != 404 {
		utils.PrintTestError(t, w.Result().StatusCode, 404)
	}
}
