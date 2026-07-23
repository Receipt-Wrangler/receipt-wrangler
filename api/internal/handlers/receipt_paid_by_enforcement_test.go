package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
)

// seedPaidByRestrictedMember creates a group plus a member whose group role
// grants receipts.read with the given paid-by visibility config. Returns the
// member user id and group id.
func seedPaidByRestrictedMember(t *testing.T, paidByUserGrantIds []uint, includeOwn bool) (uint, uint) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	group := models.Group{Name: "paidby-handler-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"PaidBy Restricted Reader",
		"",
		[]string{permissions.GroupReceiptsRead},
		nil,
		nil,
		paidByUserGrantIds,
		includeOwn, false,
	)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: "paidby-reader", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	return user.ID, group.ID
}

func getReceiptByIdRequest(receiptId uint, userId uint) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", utils.UintToString(receiptId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	return w, r
}

func TestGetReceiptDeniedWhenPaidByHidden(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "allowed-payer", Password: "x"}
	repositories.GetDB().Create(&payer)
	hidden := models.User{Username: "hidden-payer", Password: "x"}
	repositories.GetDB().Create(&hidden)

	// The member may see only payer's receipts (no own); a receipt paid by a
	// different user is hidden and accessing it by id must be denied.
	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	receiptId := seedReceipt(t, groupId, hidden.ID, 0)

	w, r := getReceiptByIdRequest(receiptId, userId)
	GetReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestGetReceiptAllowedWhenPaidByVisible(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "allowed-payer-2", Password: "x"}
	repositories.GetDB().Create(&payer)

	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	receiptId := seedReceipt(t, groupId, payer.ID, 0)

	w, r := getReceiptByIdRequest(receiptId, userId)
	GetReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestGetReceiptAllowedWhenPaidByUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()

	anyPayer := models.User{Username: "any-payer", Password: "x"}
	repositories.GetDB().Create(&anyPayer)

	// No paid-by config => unrestricted => visible regardless of payer.
	userId, groupId := seedPaidByRestrictedMember(t, nil, false)
	receiptId := seedReceipt(t, groupId, anyPayer.ID, 0)

	w, r := getReceiptByIdRequest(receiptId, userId)
	GetReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

// hasAccessRequest builds a GET /hasAccess request (the receipt-route guard's
// probe) for receiptId + permission, acting as userId.
func hasAccessRequest(receiptId uint, permission string, userId uint) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	url := "/hasAccess?receiptId=" + utils.UintToString(receiptId) + "&permission=" + permission
	r := httptest.NewRequest("GET", url, nil)
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

func TestHasAccessDeniedWhenPaidByHidden(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "ha-payer", Password: "x"}
	repositories.GetDB().Create(&payer)
	hidden := models.User{Username: "ha-hidden", Password: "x"}
	repositories.GetDB().Create(&hidden)

	// The member holds group.receipts.read but is paid-by-restricted to payer, so
	// the guard probe for a hidden-payer receipt must be denied (clean redirect).
	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	receiptId := seedReceipt(t, groupId, hidden.ID, 0)

	w, r := hasAccessRequest(receiptId, permissions.GroupReceiptsRead, userId)
	HasAccess(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestHasAccessAllowedWhenPaidByVisible(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "ha-payer-2", Password: "x"}
	repositories.GetDB().Create(&payer)

	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	receiptId := seedReceipt(t, groupId, payer.ID, 0)

	w, r := hasAccessRequest(receiptId, permissions.GroupReceiptsRead, userId)
	HasAccess(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestGetAmountOwedNotBlockedByPaidByVisibility(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "owed-payer", Password: "x"}
	repositories.GetDB().Create(&payer)
	hidden := models.User{Username: "owed-hidden", Password: "x"}
	repositories.GetDB().Create(&hidden)

	// The member is restricted to payer's receipts; amount-owed references a receipt
	// paid by a hidden user. Accounting must stay reachable (settlement totals are
	// the same for every member), so it is exempt from the paid-by check — only
	// group.receipts.read is required, which the member holds.
	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	receiptId := seedReceipt(t, groupId, hidden.ID, 0)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?receiptIds="+utils.UintToString(receiptId), nil)
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	GetAmountOwedForUser(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestExportByIdDeniedWhenAnyPaidByHidden(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "exp-payer", Password: "x"}
	repositories.GetDB().Create(&payer)
	hidden := models.User{Username: "exp-hidden", Password: "x"}
	repositories.GetDB().Create(&hidden)

	// A multi-id surface (export by id) must deny the whole request when ANY id is
	// hidden by the caller's paid-by filter.
	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	visibleReceipt := seedReceipt(t, groupId, payer.ID, 0)
	hiddenReceipt := seedReceipt(t, groupId, hidden.ID, 0)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?receiptIds="+utils.UintToString(visibleReceipt)+"&receiptIds="+utils.UintToString(hiddenReceipt), nil)
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	ExportReceiptsById(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}
