package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
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
