package handlers

import (
	"context"
	"fmt"
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

func claimsForUser(userId uint) *validator.ValidatedClaims {
	return &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}}
}

// seedRestrictedReceiptCreator creates a group plus a member whose group role
// grants receipts create/read/update restricted to categoryGrantIds, and returns
// the user id and group id.
func seedRestrictedReceiptCreator(t *testing.T, categoryGrantIds []uint) (uint, uint) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	group := models.Group{Name: "grant-handler-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"Restricted Creator",
		"",
		[]string{permissions.GroupReceiptsCreate, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate},
		categoryGrantIds,
		nil,
	)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: "restricted-creator", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	return user.ID, group.ID
}

func TestCreateReceiptDeniedForDisallowedCategory(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedCategory := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&allowedCategory)
	disallowedCategory := models.Category{Name: "Salary"}
	repositories.GetDB().Create(&disallowedCategory)

	userId, groupId := seedRestrictedReceiptCreator(t, []uint{allowedCategory.ID})

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","categories":[{"id":%d,"name":"Salary"}]}`,
		groupId, userId, disallowedCategory.ID,
	)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	CreateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestCreateReceiptDeniedForNewCategoryWithoutCreatePermission(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedCategory := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&allowedCategory)

	// The member's group role grants receipts.create but the user has no app
	// role, so it holds no app.categories.create — a new-by-name category is denied.
	userId, groupId := seedRestrictedReceiptCreator(t, []uint{allowedCategory.ID})

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","categories":[{"name":"Brand New"}]}`,
		groupId, userId,
	)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	CreateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestCreateReceiptAllowedForGrantedCategory(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedCategory := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&allowedCategory)

	userId, groupId := seedRestrictedReceiptCreator(t, []uint{allowedCategory.ID})

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","categories":[{"id":%d,"name":"Groceries"}]}`,
		groupId, userId, allowedCategory.ID,
	)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	CreateReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}
