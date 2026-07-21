package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/constants"
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

func tearDownGenericHandlerTest() {
	repositories.TruncateTestDb()
	services.ClearRolePermissionCacheForTests()
}

// requestForUser builds a recorder + request carrying JWT claims for userId.
func requestForUser(userId uint) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", strings.NewReader(""))
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}})
	return w, r.WithContext(newContext)
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Result().StatusCode != want {
		utils.PrintTestError(t, w.Result().StatusCode, want)
	}
}

func okHandlerFunc(w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
}

func TestShouldSetContentTypeHeader(t *testing.T) {
	defer tearDownGenericHandlerTest()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", strings.NewReader(""))

	handler := structs.Handler{
		Writer:          w,
		Request:         r,
		ResponseType:    constants.ApplicationJson,
		HandlerFunction: okHandlerFunc,
	}

	HandleRequest(handler)

	if contentType := w.Header().Get("Content-Type"); contentType != constants.ApplicationJson {
		utils.PrintTestError(t, contentType, constants.ApplicationJson)
	}
}

func TestShouldAllowWhenNoPermissionsRequired(t *testing.T) {
	defer tearDownGenericHandlerTest()
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:          w,
		Request:         r,
		ResponseType:    constants.ApplicationJson,
		HandlerFunction: okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusOK)
}

func TestShouldRejectWhenUserLacksAppPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppUsersRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:          w,
		Request:         r,
		ResponseType:    constants.ApplicationJson,
		AppPermissions:  []string{permissions.AppUsersDelete},
		HandlerFunction: okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldAcceptWhenUserHasAppPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppUsersRead, permissions.AppUsersDelete)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:          w,
		Request:         r,
		ResponseType:    constants.ApplicationJson,
		AppPermissions:  []string{permissions.AppUsersDelete},
		HandlerFunction: okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusOK)
}

func TestShouldAcceptWhenUserHasAnyAppPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	// The caller holds one of the two any-of permissions, which satisfies the gate.
	grantAppPerms(t, 1, permissions.AppUsersRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:            w,
		Request:           r,
		ResponseType:      constants.ApplicationJson,
		AnyAppPermissions: []string{permissions.AppUsersRead, permissions.AppUsersDelete},
		HandlerFunction:   okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusOK)
}

func TestShouldRejectWhenUserHasNoneOfAnyAppPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	// The caller holds an unrelated permission but neither of the any-of set, so the
	// gate denies.
	grantAppPerms(t, 1, permissions.AppUsersCreate)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:            w,
		Request:           r,
		ResponseType:      constants.ApplicationJson,
		AnyAppPermissions: []string{permissions.AppUsersRead, permissions.AppUsersDelete},
		HandlerFunction:   okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldRejectGroupPermissionWhenNotAssignedRole(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupId:          "1",
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldRejectGroupPermissionWhenRoleLacksPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupId:          "1",
		GroupPermissions: []string{permissions.GroupReceiptsDelete},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldAcceptGroupPermissionWhenRoleGrantsPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsDelete)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupId:          "1",
		GroupPermissions: []string{permissions.GroupReceiptsDelete},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusOK)
}

func TestShouldRejectGroupPermissionWhenNoGroupProvided(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupId:          "",
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldResolveReceiptGroupAndReject(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	db := repositories.GetDB()
	db.Create(&models.Receipt{Name: "Test receipt", GroupId: 1, PaidByUserID: 1})
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		ReceiptId:        "1",
		GroupPermissions: []string{permissions.GroupReceiptsDelete},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldResolveReceiptGroupAndAccept(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	db := repositories.GetDB()
	db.Create(&models.Receipt{Name: "Test receipt", GroupId: 1, PaidByUserID: 1})
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		ReceiptId:        "1",
		GroupPermissions: []string{permissions.GroupReceiptsDelete},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusOK)
}

func TestShouldRejectMultipleReceiptsWhenOneGroupMissingPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	db := repositories.GetDB()
	db.Create(&models.Receipt{Name: "Receipt 1", GroupId: 1, PaidByUserID: 1})
	// receipt 2 belongs to group 2, where user 1 has no role
	db.Create(&models.Receipt{Name: "Receipt 2", GroupId: 2, PaidByUserID: 1})
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		ReceiptIds:       []string{"1", "2"},
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldRejectGroupIdsWhenOneMissingPermission(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	// user 1 is a member of group 1 but not group 2
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupIds:         []string{"1", "2"},
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusForbidden)
}

func TestShouldAcceptGroupViaOrAppFallback(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	// user 1 is a member of group 1 but has no group role (so the group check
	// fails), yet holds the app-scoped fallback permission.
	grantAppPerms(t, 1, permissions.AppGroupsRead)
	w, r := requestForUser(1)

	handler := structs.Handler{
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationJson,
		GroupId:          "1",
		GroupPermissions: []string{permissions.GroupView},
		OrAppPermissions: []string{permissions.AppGroupsRead},
		HandlerFunction:  okHandlerFunc,
	}

	HandleRequest(handler)

	assertStatus(t, w, http.StatusOK)
}
