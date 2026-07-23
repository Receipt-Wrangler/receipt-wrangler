package handlers

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
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
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

func claimsForUser(userId uint) *validator.ValidatedClaims {
	return &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}}
}

// seedRestrictedReceiptCreator creates a group plus a member whose group role
// grants receipts create/read/update restricted to categoryGrantIds (tags
// unrestricted), and returns the user id and group id.
func seedRestrictedReceiptCreator(t *testing.T, categoryGrantIds []uint) (uint, uint) {
	return seedRestrictedReceiptCreatorWithGrants(t, categoryGrantIds, nil)
}

// seedRestrictedReceiptCreatorWithGrants is seedRestrictedReceiptCreator with an
// explicit tag grant set, so tests can restrict tags as well as categories.
func seedRestrictedReceiptCreatorWithGrants(t *testing.T, categoryGrantIds []uint, tagGrantIds []uint) (uint, uint) {
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
		tagGrantIds,
		nil,
		false, false,
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

// seedRestrictedQuickScanner creates a group (with default receipt settings) plus a member whose
// group role grants only group.receipts.quick-scan, restricted to categoryGrantIds. It returns the
// user id and group id. Receipt settings are seeded because resolveQuickScanFields reads them with
// .First() and would otherwise 500 on a missing row.
func seedRestrictedQuickScanner(t *testing.T, categoryGrantIds []uint) (uint, uint) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	group := models.Group{Name: "quick-scan-grant-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	if _, err := repositories.NewGroupReceiptSettingsRepository(nil).CreateGroupReceiptSettings(group.ID); err != nil {
		t.Fatalf("seed group receipt settings: %v", err)
	}

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"Restricted Quick Scanner",
		"",
		[]string{permissions.GroupReceiptsQuickScan},
		categoryGrantIds,
		nil,
		nil,
		false, false,
	)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: "restricted-quick-scanner", Password: "password"}
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

func TestEnforceQuickScanGrantSelection(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedCategory := models.Category{Name: "Groceries"}
	if err := repositories.GetDB().Create(&allowedCategory).Error; err != nil {
		t.Fatalf("seed allowed category: %v", err)
	}
	disallowedCategory := models.Category{Name: "Salary"}
	if err := repositories.GetDB().Create(&disallowedCategory).Error; err != nil {
		t.Fatalf("seed disallowed category: %v", err)
	}
	allowedTag := models.Tag{Name: "Reimbursable"}
	if err := repositories.GetDB().Create(&allowedTag).Error; err != nil {
		t.Fatalf("seed allowed tag: %v", err)
	}
	disallowedTag := models.Tag{Name: "Personal"}
	if err := repositories.GetDB().Create(&disallowedTag).Error; err != nil {
		t.Fatalf("seed disallowed tag: %v", err)
	}

	userId, groupId := seedRestrictedReceiptCreatorWithGrants(t, []uint{allowedCategory.ID}, []uint{allowedTag.ID})

	// A category pick outside the caller's grants is denied.
	deniedCommand := commands.QuickScanCommand{
		GroupIds:    []uint{groupId},
		CategoryIds: [][]uint{{disallowedCategory.ID}},
	}
	allowed, denyMessage, err := enforceQuickScanGrantSelection(userId, deniedCommand)
	if err != nil {
		t.Fatalf("enforce category denied: %v", err)
	}
	if allowed {
		t.Error("expected out-of-grant category pick to be denied")
	}
	if len(denyMessage) == 0 {
		t.Error("expected a deny message when a category pick is rejected")
	}

	// A tag pick outside the caller's grants is denied.
	tagDeniedCommand := commands.QuickScanCommand{
		GroupIds: []uint{groupId},
		TagIds:   [][]uint{{disallowedTag.ID}},
	}
	allowed, denyMessage, err = enforceQuickScanGrantSelection(userId, tagDeniedCommand)
	if err != nil {
		t.Fatalf("enforce tag denied: %v", err)
	}
	if allowed {
		t.Error("expected out-of-grant tag pick to be denied")
	}
	if len(denyMessage) == 0 {
		t.Error("expected a deny message when a tag pick is rejected")
	}

	// Picks within the caller's grants (category and tag) are allowed.
	allowedCommand := commands.QuickScanCommand{
		GroupIds:    []uint{groupId},
		CategoryIds: [][]uint{{allowedCategory.ID}},
		TagIds:      [][]uint{{allowedTag.ID}},
	}
	allowed, _, err = enforceQuickScanGrantSelection(userId, allowedCommand)
	if err != nil {
		t.Fatalf("enforce allowed: %v", err)
	}
	if !allowed {
		t.Error("expected in-grant category and tag picks to be allowed")
	}

	// The loop rejects the whole request when any file carries an out-of-grant pick.
	mixedCommand := commands.QuickScanCommand{
		GroupIds:    []uint{groupId, groupId},
		CategoryIds: [][]uint{{allowedCategory.ID}, {disallowedCategory.ID}},
	}
	allowed, _, err = enforceQuickScanGrantSelection(userId, mixedCommand)
	if err != nil {
		t.Fatalf("enforce mixed: %v", err)
	}
	if allowed {
		t.Error("expected a request with any out-of-grant file to be denied")
	}
}

// TestQuickScanHandlerDeniesOutOfGrantCategory drives the real QuickScan handler to prove the grant
// gate is actually wired in (returns 403), not just unit-tested in isolation. The deny returns before
// the enqueue loop, so no Redis/AI/file processing is exercised.
func TestQuickScanHandlerDeniesOutOfGrantCategory(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedCategory := models.Category{Name: "Groceries"}
	if err := repositories.GetDB().Create(&allowedCategory).Error; err != nil {
		t.Fatalf("seed allowed category: %v", err)
	}
	disallowedCategory := models.Category{Name: "Salary"}
	if err := repositories.GetDB().Create(&disallowedCategory).Error; err != nil {
		t.Fatalf("seed disallowed category: %v", err)
	}

	userId, groupId := seedRestrictedQuickScanner(t, []uint{allowedCategory.ID})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "receipt.png")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	part.Write([]byte("not-a-real-image"))

	// One aligned entry per file. Paid-by is non-zero and status non-empty so the required-field
	// config check passes and control reaches the grant check; the picked category is out of grant.
	writer.WriteField("groupIds", utils.UintToString(groupId))
	writer.WriteField("paidByUserIds", utils.UintToString(userId))
	writer.WriteField("statuses", string(models.OPEN))
	writer.WriteField("categoryIds", utils.UintToString(disallowedCategory.ID))

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	QuickScan(w, r)

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

// seedCustomField creates a custom field and returns it.
func seedCustomField(t *testing.T) models.CustomField {
	t.Helper()
	field := models.CustomField{Name: "Field", Type: models.TEXT}
	if err := repositories.GetDB().Create(&field).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}
	return field
}

// seedReceipt creates a persisted receipt in groupId, optionally attaching a
// custom field value for customFieldId (0 = none), and returns the receipt id.
func seedReceipt(t *testing.T, groupId uint, paidByUserId uint, customFieldId uint) uint {
	t.Helper()
	receipt := models.Receipt{
		Name:         "Existing",
		Amount:       decimal.NewFromInt(5),
		Date:         time.Now(),
		GroupId:      groupId,
		PaidByUserID: paidByUserId,
		Status:       models.OPEN,
	}
	if customFieldId != 0 {
		receipt.CustomFields = []models.CustomFieldValue{{CustomFieldId: customFieldId}}
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt.ID
}

// updateReceiptRequest builds an UpdateReceipt request for receiptId acting as
// userId, with the given JSON body.
func updateReceiptRequest(receiptId uint, userId uint, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/api", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", utils.UintToString(receiptId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	return w, r
}

func TestCreateReceiptDeniedForCustomFieldWithoutAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	field := seedCustomField(t)

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[{"customFieldId":%d}]}`,
		groupId, userId, field.ID,
	)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	CreateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestCreateReceiptAllowedForCustomFieldWithAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	grantAppPerms(t, userId, permissions.AppCustomFieldsRead)
	field := seedCustomField(t)

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[{"customFieldId":%d}]}`,
		groupId, userId, field.ID,
	)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	CreateReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestUpdateReceiptDeniedWhenAddingCustomFieldWithoutAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	field := seedCustomField(t)
	receiptId := seedReceipt(t, groupId, userId, 0)

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[{"customFieldId":%d}]}`,
		groupId, userId, field.ID,
	)
	w, r := updateReceiptRequest(receiptId, userId, body)

	UpdateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestUpdateReceiptDeniedWhenRemovingCustomFieldWithoutAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	field := seedCustomField(t)
	receiptId := seedReceipt(t, groupId, userId, field.ID)

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[]}`,
		groupId, userId,
	)
	w, r := updateReceiptRequest(receiptId, userId, body)

	UpdateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestUpdateReceiptAllowedWhenEditingCustomFieldValueWithoutAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	field := seedCustomField(t)
	receiptId := seedReceipt(t, groupId, userId, field.ID)

	// Same custom field set, only the value changes — allowed without read.
	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[{"customFieldId":%d,"stringValue":"edited"}]}`,
		groupId, userId, field.ID,
	)
	w, r := updateReceiptRequest(receiptId, userId, body)

	UpdateReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestUpdateReceiptAllowedWhenChangingCustomFieldsWithAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	grantAppPerms(t, userId, permissions.AppCustomFieldsRead)
	field := seedCustomField(t)
	receiptId := seedReceipt(t, groupId, userId, 0)

	// Read holder may add a custom field freely.
	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[{"customFieldId":%d}]}`,
		groupId, userId, field.ID,
	)
	w, r := updateReceiptRequest(receiptId, userId, body)

	UpdateReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}
