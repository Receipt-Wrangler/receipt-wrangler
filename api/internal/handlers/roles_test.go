package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strconv"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

// Role-management handlers are now permission-gated (app.roles.*). adminContext
// returns claims for user 2 and ensures that user actually holds the role
// permissions in the database (the JWT role itself is no longer trusted).
// userContext returns claims for a different user that holds no role, so
// "due to role" cases are rejected with 403.
func adminContext() *validator.ValidatedClaims {
	ensureRoleAdmin(2)
	return &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2, UserRole: models.ADMIN}}
}

func userContext() *validator.ValidatedClaims {
	return &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 3, UserRole: models.USER}}
}

// ensureRoleAdmin idempotently gives userId an app role granting every app.roles
// permission, so role-management handlers authorize the user. It creates a single
// reusable "Role Admin" role within the test's (truncated) database.
func ensureRoleAdmin(userId uint) {
	services.ClearRolePermissionCacheForTests()
	db := repositories.GetDB()

	var role models.AppRole
	if err := db.Where("name = ?", "Role Admin").First(&role).Error; err != nil {
		roleRepository := repositories.NewRoleRepository(nil)
		created, cErr := roleRepository.CreateAppRole("Role Admin", "", []string{
			permissions.AppRolesCreate,
			permissions.AppRolesRead,
			permissions.AppRolesUpdate,
			permissions.AppRolesDelete,
		})
		if cErr != nil {
			panic(cErr)
		}
		role = created
	}

	var count int64
	db.Model(&models.User{}).Where("id = ?", userId).Count(&count)
	if count == 0 {
		db.Create(&models.User{BaseModel: models.BaseModel{ID: userId}, Username: "role-admin", Password: "password", AppRoleID: &role.ID})
	} else {
		db.Model(&models.User{}).Where("id = ?", userId).Update("app_role_id", role.ID)
	}
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

	// CreateTestRoles seeds two roles; adminContext adds the "Role Admin" role
	// that authorizes this request, so GetRoles returns three.
	if len(roles) != 3 {
		utils.PrintTestError(t, len(roles), 3)
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
	w := updateRoleRequest(strconv.FormatUint(uint64(created.ID), 10), body, adminContext())

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
	w := updateRoleRequest(strconv.FormatUint(uint64(systemRole.ID), 10), body, adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestShouldReturnBadRequestForInvalidRoleId(t *testing.T) {
	defer repositories.TruncateTestDb()

	body := `{"name": "Role", "scope": "APP", "permissions": ["app.users.read"]}`
	w := updateRoleRequest("abc", body, adminContext())

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

func deleteRoleRequest(roleId string, scope string, claims *validator.ValidatedClaims) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api?scope="+scope, nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("roleId", roleId)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claims))

	DeleteRole(w, r)
	return w
}

func TestShouldDeleteAppRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()

	w := deleteRoleRequest("1", "APP", adminContext())

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	roleRepository := repositories.NewRoleRepository(nil)
	if _, err := roleRepository.GetAppRoleById(1); err == nil {
		utils.PrintTestError(t, "app role still exists", "app role should be deleted")
	}
}

func TestShouldDeleteGroupRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()

	w := deleteRoleRequest("1", "GROUP", adminContext())

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	roleRepository := repositories.NewRoleRepository(nil)
	if _, err := roleRepository.GetGroupRoleById(1); err == nil {
		utils.PrintTestError(t, "group role still exists", "group role should be deleted")
	}
}

func TestShouldNotDeleteRoleDueToRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()

	w := deleteRoleRequest("1", "APP", userContext())

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestShouldReturnBadRequestForInvalidDeleteRoleId(t *testing.T) {
	defer repositories.TruncateTestDb()

	w := deleteRoleRequest("abc", "APP", adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestShouldReturnBadRequestForInvalidDeleteScope(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestRoles()

	for _, scope := range []string{"", "BOGUS"} {
		w := deleteRoleRequest("1", scope, adminContext())
		if w.Result().StatusCode != 400 {
			utils.PrintTestError(t, w.Result().StatusCode, 400)
		}
	}
}

func TestShouldReturnNotFoundForMissingDeleteRole(t *testing.T) {
	defer repositories.TruncateTestDb()

	w := deleteRoleRequest("999", "APP", adminContext())

	if w.Result().StatusCode != 404 {
		utils.PrintTestError(t, w.Result().StatusCode, 404)
	}
}

func TestShouldNotDeleteRoleDueToScopeMismatch(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)
	created, err := roleRepository.CreateAppRole("App Role", "desc", []string{permissions.AppUsersCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Deleting an app role id under the GROUP scope is a type mismatch, not a delete.
	w := deleteRoleRequest(strconv.FormatUint(uint64(created.ID), 10), "GROUP", adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
		return
	}

	if _, err := roleRepository.GetAppRoleById(created.ID); err != nil {
		utils.PrintTestError(t, err, "app role should be untouched")
	}
}

func TestShouldNotDeleteSystemRole(t *testing.T) {
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

	w := deleteRoleRequest(strconv.FormatUint(uint64(systemRole.ID), 10), "APP", adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
		return
	}

	roleRepository := repositories.NewRoleRepository(nil)
	if _, err := roleRepository.GetAppRoleById(systemRole.ID); err != nil {
		utils.PrintTestError(t, err, "system role should be untouched")
	}
}

func TestShouldNotDeleteAssignedAppRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)
	created, err := roleRepository.CreateAppRole("App Role", "desc", []string{permissions.AppUsersCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	db := repositories.GetDB()
	user := models.User{Username: "assigned-user", Password: "password", AppRoleID: &created.ID}
	if err := db.Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	w := deleteRoleRequest(strconv.FormatUint(uint64(created.ID), 10), "APP", adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
		return
	}

	if _, err := roleRepository.GetAppRoleById(created.ID); err != nil {
		utils.PrintTestError(t, err, "assigned app role should be untouched")
	}
}

func TestShouldNotDeleteAssignedGroupRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	roleRepository := repositories.NewRoleRepository(nil)
	created, err := roleRepository.CreateGroupRole("Group Role", "desc", []string{permissions.GroupReceiptsCreate})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	db := repositories.GetDB()
	user := models.User{Username: "member-user", Password: "password"}
	if err := db.Table("users").Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	group := models.Group{Name: "Test Group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	member := models.GroupMember{
		UserID:      user.ID,
		GroupID:     group.ID,
		GroupRole:   models.OWNER,
		GroupRoleID: &created.ID,
	}
	if err := db.Create(&member).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	w := deleteRoleRequest(strconv.FormatUint(uint64(created.ID), 10), "GROUP", adminContext())

	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
		return
	}

	if _, err := roleRepository.GetGroupRoleById(created.ID); err != nil {
		utils.PrintTestError(t, err, "assigned group role should be untouched")
	}
}
