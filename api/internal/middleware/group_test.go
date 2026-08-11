package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func groupSetup() (http.HandlerFunc, *http.Request, *httptest.ResponseRecorder) {
	createUserAndGroup()
	fakeHandler, r, w := createFakeGroupHandler()

	return fakeHandler, r, w
}

func createUserAndGroup() {
	user := models.User{
		Username:    "test",
		Password:    "Password",
		DisplayName: "test",
	}
	db := repositories.GetDB()
	db.Create(&user)

	groupMembers := make([]models.GroupMember, 1)
	groupMembers = append(groupMembers, models.GroupMember{UserID: user.ID})

	group := models.Group{
		Name:         "Test",
		GroupMembers: groupMembers,
	}

	db.Create(&group)
}

func createFakeGroupHandler() (http.HandlerFunc, *http.Request, *httptest.ResponseRecorder) {
	reader := strings.NewReader("")
	r := httptest.NewRequest(http.MethodGet, "/api/1", reader)
	w := httptest.NewRecorder()

	var vClaims validator.ValidatedClaims
	vClaims.CustomClaims = &structs.Claims{UserId: 1}

	ctx := context.WithValue(r.Context(), "groupId", "1")
	ctx = context.WithValue(ctx, jwtmiddleware.ContextKey{}, &vClaims)
	r = r.WithContext(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), r, w
}

func teardownGroupTest() {
	db := repositories.GetDB()
	repositories.TruncateTable(db, "group_members")
	repositories.TruncateTable(db, "groups")
	repositories.TruncateTable(db, "users")
	repositories.TruncateTable(db, "app_role_permissions")
	repositories.TruncateTable(db, "app_roles")
	services.ClearRolePermissionCacheForTests()
}

// grantAppPermsToUser gives userId an app role granting exactly perms. Role ids
// are reused across truncations, so the permission cache is cleared alongside.
func grantAppPermsToUser(t *testing.T, userId uint, perms ...string) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()

	role, err := repositories.NewRoleRepository(nil).CreateAppRole("Test App Role", "", perms, false)
	if err != nil {
		t.Fatalf("create app role: %v", err)
	}

	err = repositories.GetDB().Model(&models.User{}).Where("id = ?", userId).Update("app_role_id", role.ID).Error
	if err != nil {
		t.Fatalf("assign app role: %v", err)
	}
}

func TestCanDeleteGroupShouldReject1(t *testing.T) {
	defer teardownGroupTest()
	fakeHandler, r, w := groupSetup()
	handler := CanDeleteGroup(fakeHandler)
	handler.ServeHTTP(w, r)

	if w.Result().StatusCode != 500 {
		utils.PrintTestError(t, w.Result().StatusCode, 500)
	}
}

// The "stay in at least one group" rule is self-protection for an ordinary user
// deleting their own group. An administrator holding app.groups.delete cleans up
// groups they are not a member of, so the rule must not apply to them — the
// same setup that rejects in TestCanDeleteGroupShouldReject1 must pass here.
func TestCanDeleteGroupShouldAllowAppGroupsDeleteHolder(t *testing.T) {
	defer teardownGroupTest()
	_, r, w := groupSetup()
	grantAppPermsToUser(t, 1, permissions.AppGroupsDelete)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	CanDeleteGroup(next).ServeHTTP(w, r)

	if !nextCalled {
		utils.PrintTestError(t, "next handler not called", "next handler called")
	}
	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestCanDeleteGroupShouldReject2(t *testing.T) {
	defer teardownGroupTest()
	fakeHandler, _, w := groupSetup()
	groupMembers := make([]models.GroupMember, 1)
	groupMembers = append(groupMembers, models.GroupMember{UserID: 1})

	group := models.Group{
		Name:         "Another group",
		GroupMembers: groupMembers,
	}

	repositories.GetDB().Create(&group)

	CanDeleteGroup(fakeHandler)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}
