package handlers

import (
	"context"
	"fmt"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
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

func TestUintSetsEqual(t *testing.T) {
	cases := []struct {
		name string
		a    []uint
		b    []uint
		want bool
	}{
		{"both empty", nil, nil, true},
		{"same single", []uint{1}, []uint{1}, true},
		{"same multi unordered", []uint{1, 2, 3}, []uint{3, 1, 2}, true},
		{"added", []uint{1, 2}, []uint{1}, false},
		{"removed", []uint{1}, []uint{1, 2}, false},
		{"different", []uint{1}, []uint{2}, false},
	}

	toSet := func(ids []uint) map[uint]struct{} {
		set := make(map[uint]struct{}, len(ids))
		for _, id := range ids {
			set[id] = struct{}{}
		}
		return set
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := uintSetsEqual(toSet(c.a), toSet(c.b)); got != c.want {
				utils.PrintTestError(t, got, c.want)
			}
		})
	}
}

func customFieldSelectionCommand(ids ...uint) commands.UpsertReceiptCommand {
	values := make([]commands.UpsertCustomFieldValueCommand, 0, len(ids))
	for _, id := range ids {
		values = append(values, commands.UpsertCustomFieldValueCommand{CustomFieldId: id})
	}
	return commands.UpsertReceiptCommand{CustomFields: values}
}

func TestEnforceReceiptCustomFieldSelectionWithoutAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	// The seeded creator has only group permissions, so no app.custom-fields.read.
	userId, _ := seedRestrictedReceiptCreator(t, nil)

	cases := []struct {
		name      string
		submitted []uint
		current   []uint
		wantAllow bool
	}{
		{"add a field is denied", []uint{1}, nil, false},
		{"remove a field is denied", []uint{1}, []uint{1, 2}, false},
		{"same set (value edit) is allowed", []uint{1, 2}, []uint{2, 1}, true},
		{"no custom fields is allowed", nil, nil, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			allowed, _, err := enforceReceiptCustomFieldSelection(userId, customFieldSelectionCommand(c.submitted...), c.current)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if allowed != c.wantAllow {
				utils.PrintTestError(t, allowed, c.wantAllow)
			}
		})
	}
}

func TestEnforceReceiptCustomFieldSelectionWithReadPermission(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, _ := seedRestrictedReceiptCreator(t, nil)
	grantAppPerms(t, userId, permissions.AppCustomFieldsRead)

	// A read holder may change the attached set freely (here, an add).
	allowed, _, err := enforceReceiptCustomFieldSelection(userId, customFieldSelectionCommand(1), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		utils.PrintTestError(t, allowed, true)
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
