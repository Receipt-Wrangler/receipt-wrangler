package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
)

// isoHandlerFixture is an isolated group with a supervisor (SeesAllMembers) and two
// plain isolated members, all holding read/create/update on receipts. Members A and
// B cannot see each other; both can see the supervisor.
type isoHandlerFixture struct {
	groupId      uint
	supervisorId uint
	memberAId    uint
	memberBId    uint
}

func seedIsoHandlerUser(t *testing.T, username string) uint {
	t.Helper()
	user := models.User{Username: username, Password: "password"}
	if err := repositories.GetDB().Create(&user).Error; err != nil {
		t.Fatalf("seed user %s: %v", username, err)
	}
	return user.ID
}

func seedIsolatedReceiptGroupHandler(t *testing.T, isolate bool) isoHandlerFixture {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()
	roleRepository := repositories.NewRoleRepository(nil)

	perms := []string{
		permissions.GroupReceiptsRead,
		permissions.GroupReceiptsCreate,
		permissions.GroupReceiptsUpdate,
	}

	group := models.Group{Name: "iso-handler-group", IsolateMembers: isolate}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	supRole, err := roleRepository.CreateGroupRole("Iso H Supervisor", "", perms, nil, nil, nil, false, true)
	if err != nil {
		t.Fatalf("seed supervisor role: %v", err)
	}
	memberRole, err := roleRepository.CreateGroupRole("Iso H Member", "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("seed member role: %v", err)
	}

	sup := seedIsoHandlerUser(t, "iso-h-sup")
	a := seedIsoHandlerUser(t, "iso-h-a")
	b := seedIsoHandlerUser(t, "iso-h-b")

	for _, m := range []models.GroupMember{
		{GroupID: group.ID, UserID: sup, GroupRoleID: &supRole.ID},
		{GroupID: group.ID, UserID: a, GroupRoleID: &memberRole.ID},
		{GroupID: group.ID, UserID: b, GroupRoleID: &memberRole.ID},
	} {
		if err := db.Create(&m).Error; err != nil {
			t.Fatalf("seed member: %v", err)
		}
	}

	return isoHandlerFixture{groupId: group.ID, supervisorId: sup, memberAId: a, memberBId: b}
}

func createReceiptRequest(userId uint, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

// --- Write-side guard (task 8) ---

func TestMemberIsolationCreateDeniedForNonVisiblePaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)

	// Member A tries to plant member B (whom A cannot see) as the payer.
	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`,
		fx.groupId, fx.memberBId,
	)
	w, r := createReceiptRequest(fx.memberAId, body)
	CreateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestMemberIsolationCreateDeniedForNonVisibleChargedTo(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)

	// Payer is A (visible), but an item charges the non-visible member B.
	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","receiptItems":[{"name":"I","amount":"5","status":"OPEN","chargedToUserId":%d}]}`,
		fx.groupId, fx.memberAId, fx.memberBId,
	)
	w, r := createReceiptRequest(fx.memberAId, body)
	CreateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestMemberIsolationCreateAllowedForVisibleUsers(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)

	// Payer and charged-to are both A (self, always visible).
	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","receiptItems":[{"name":"I","amount":"5","status":"OPEN","chargedToUserId":%d}]}`,
		fx.groupId, fx.memberAId, fx.memberAId,
	)
	w, r := createReceiptRequest(fx.memberAId, body)
	CreateReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestMemberIsolationUpdateDeniedForNonVisiblePaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	// A receipt A can access (paid by self), which A then tries to reassign to B.
	receiptId := seedReceipt(t, fx.groupId, fx.memberAId, 0)

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`,
		fx.groupId, fx.memberBId,
	)
	w, r := updateReceiptRequest(receiptId, fx.memberAId, body)
	UpdateReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

// --- Single receipt read (surface B via GetReceiptForUser) ---

func TestMemberIsolationGetReceiptDeniedForNonVisiblePayer(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	// A receipt paid by the non-visible member B.
	receiptId := seedReceipt(t, fx.groupId, fx.memberBId, 0)

	w, r := getReceiptByIdRequest(receiptId, fx.memberAId)
	GetReceipt(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestMemberIsolationGetReceiptAllowedForOwnReceipt(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	receiptId := seedReceipt(t, fx.groupId, fx.memberAId, 0)

	w, r := getReceiptByIdRequest(receiptId, fx.memberAId)
	GetReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

// --- Paged list (surface B + count correctness) ---

func pagedReceiptPayers(t *testing.T, w *httptest.ResponseRecorder) ([]uint, int64) {
	t.Helper()
	var pagedData structs.PagedData
	if err := json.Unmarshal(w.Body.Bytes(), &pagedData); err != nil {
		t.Fatalf("unmarshal paged data: %v", err)
	}
	payers := make([]uint, 0, len(pagedData.Data))
	for _, entry := range pagedData.Data {
		entryBytes, _ := json.Marshal(entry)
		var receipt models.Receipt
		json.Unmarshal(entryBytes, &receipt)
		payers = append(payers, receipt.PaidByUserID)
	}
	return payers, pagedData.TotalCount
}

func TestMemberIsolationPagedListHidesNonVisiblePayer(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	seedReceipt(t, fx.groupId, fx.memberAId, 0) // A's own
	seedReceipt(t, fx.groupId, fx.memberBId, 0) // B's — hidden from A

	requestBody := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{Page: 1, PageSize: 10},
	}
	body, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(string(body)))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", utils.UintToString(fx.groupId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(fx.memberAId)))

	GetPagedReceiptsForGroup(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}

	payers, count := pagedReceiptPayers(t, w)
	if count != 1 {
		t.Errorf("expected totalCount 1 (only A's receipt), got %d", count)
	}
	if len(payers) != 1 || payers[0] != fx.memberAId {
		t.Errorf("expected only A's receipt in the page, got payers %v", payers)
	}
}

func TestMemberIsolationGetReceiptsForGroupIdsHidesNonVisiblePayer(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	seedReceipt(t, fx.groupId, fx.memberAId, 0)
	seedReceipt(t, fx.groupId, fx.memberBId, 0)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?groupIds="+utils.UintToString(fx.groupId), nil)
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(fx.memberAId)))

	GetReceiptsForGroupIds(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}

	var receipts []models.Receipt
	if err := json.Unmarshal(w.Body.Bytes(), &receipts); err != nil {
		t.Fatalf("unmarshal receipts: %v", err)
	}
	if len(receipts) != 1 || receipts[0].PaidByUserID != fx.memberAId {
		t.Errorf("expected only A's receipt, got %d receipts", len(receipts))
	}
}

// --- Backward compatibility ---

func TestMemberIsolationBackwardCompatNonIsolatedGroupSeesEverything(t *testing.T) {
	defer repositories.TruncateTestDb()

	// Same setup, but the group does NOT isolate members: A is unrestricted and
	// sees B's receipt exactly as before.
	fx := seedIsolatedReceiptGroupHandler(t, false)
	receiptId := seedReceipt(t, fx.groupId, fx.memberBId, 0)

	w, r := getReceiptByIdRequest(receiptId, fx.memberAId)
	GetReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}
