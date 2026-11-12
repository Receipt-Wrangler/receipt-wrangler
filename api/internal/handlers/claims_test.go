package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

// claimRequest builds a request acting as userId, with the given chi URL params.
func claimRequest(method string, userId uint, body string, params map[string]string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, "/api/claim", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

func claimParams(claimId uint) map[string]string {
	return map[string]string{"claimId": utils.UintToString(claimId)}
}

func groupParams(groupId uint) map[string]string {
	return map[string]string{"groupId": utils.UintToString(groupId)}
}

func decodeInto(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), target); err != nil {
		t.Fatalf("decode %q: %v", w.Body.String(), err)
	}
}

// seedClaimRow writes a claim straight to the database, as if createdBy made it.
func seedClaimRow(t *testing.T, groupId uint, createdBy uint, name string) models.Claim {
	t.Helper()
	claim := models.Claim{
		BaseModel: models.BaseModel{CreatedBy: &createdBy},
		GroupId:   groupId,
		Name:      name,
		Status:    models.ClaimStatusDraft,
	}
	if err := repositories.GetDB().Create(&claim).Error; err != nil {
		t.Fatalf("seed claim: %v", err)
	}
	return claim
}

func seedReceiptWithAmount(t *testing.T, groupId uint, paidByUserId uint, amount string) uint {
	t.Helper()
	receipt := models.Receipt{
		Name:         "Claimed",
		Amount:       decimal.RequireFromString(amount),
		Date:         time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		GroupId:      groupId,
		PaidByUserID: paidByUserId,
		Status:       models.OPEN,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt.ID
}

func linkReceipts(t *testing.T, claimId uint, receiptIds ...uint) {
	t.Helper()
	if err := repositories.NewClaimRepository(nil).LinkReceipts(claimId, receiptIds); err != nil {
		t.Fatalf("link: %v", err)
	}
}

func seedLedgerEntry(t *testing.T, claimId uint, createdBy uint, entryType string, amount string) models.ClaimLedgerEntry {
	t.Helper()
	entry := models.ClaimLedgerEntry{
		BaseModel: models.BaseModel{CreatedBy: &createdBy},
		ClaimId:   claimId,
		EntryType: entryType,
		Date:      time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
	}
	if amount != "" {
		value := decimal.RequireFromString(amount)
		entry.Amount = &value
	}
	if err := repositories.GetDB().Create(&entry).Error; err != nil {
		t.Fatalf("seed entry: %v", err)
	}
	return entry
}

func getClaim(t *testing.T, userId uint, claimId uint) structs.ClaimView {
	t.Helper()
	w, r := claimRequest(http.MethodGet, userId, "", claimParams(claimId))
	GetClaimById(w, r)
	assertStatus(t, w, http.StatusOK)
	var view structs.ClaimView
	decodeInto(t, w, &view)
	return view
}

func listClaims(t *testing.T, userId uint, groupId uint) structs.ClaimListView {
	t.Helper()
	w, r := claimRequest(http.MethodGet, userId, "", groupParams(groupId))
	GetClaimsForGroup(w, r)
	assertStatus(t, w, http.StatusOK)
	var list structs.ClaimListView
	decodeInto(t, w, &list)
	return list
}

func TestClaims_CreateGetUpdateRoundTrip(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	w, r := claimRequest(http.MethodPost, 1, `{
		"name": "PCS Rome", "referenceNumber": "TA-1", "type": "Travel voucher", "status": "SUBMITTED",
		"submittedDate": "2026-09-01", "dueDate": "2026-12-01", "claimedAmount": "412.50", "notes": "n"
	}`, groupParams(1))
	CreateClaim(w, r)
	assertStatus(t, w, http.StatusOK)
	var created structs.ClaimView
	decodeInto(t, w, &created)

	if created.GroupId != 1 || created.Name != "PCS Rome" || created.Status != "SUBMITTED" ||
		created.SubmittedDate != "2026-09-01" || created.DueDate != "2026-12-01" ||
		created.ClaimedAmount != "412.50" || created.EffectiveClaimedAmount != "412.50" ||
		created.Type != "Travel voucher" || created.ReferenceNumber != "TA-1" {
		utils.PrintTestError(t, created, "the submitted fields")
	}
	if created.CreatedBy == nil || *created.CreatedBy != 1 {
		utils.PrintTestError(t, created.CreatedBy, 1)
	}
	if !created.AwaitingSettlement || created.DaysWaiting == nil {
		utils.PrintTestError(t, created.DaysWaiting, "an age for a submitted claim")
	}

	// A full replace: omitted optional fields are cleared.
	w, r = claimRequest(http.MethodPut, 1, `{"name": "PCS Rome (resubmitted)", "status": "DISPUTED"}`, claimParams(created.Id))
	UpdateClaim(w, r)
	assertStatus(t, w, http.StatusOK)

	got := getClaim(t, 1, created.Id)
	if got.Name != "PCS Rome (resubmitted)" || got.Status != "DISPUTED" {
		utils.PrintTestError(t, got, "renamed and DISPUTED")
	}
	if got.ClaimedAmount != "" || got.SubmittedDate != "" || got.ReferenceNumber != "" || got.Notes != "" {
		utils.PrintTestError(t, got, "optional fields cleared")
	}
	if got.DaysWaiting != nil {
		t.Errorf("without a submitted date there is no age")
	}

	w, r = claimRequest(http.MethodDelete, 1, "", claimParams(created.Id))
	DeleteClaim(w, r)
	assertStatus(t, w, http.StatusOK)

	w, r = claimRequest(http.MethodGet, 1, "", claimParams(created.Id))
	GetClaimById(w, r)
	assertStatus(t, w, http.StatusForbidden)
}

func TestClaims_ReadPermissionDoesNotAllowWrites(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 2, 1, permissions.GroupReceiptsRead)
	claim := seedClaimRow(t, 1, 1, "existing")

	getClaim(t, 2, claim.ID)

	writes := []func() *httptest.ResponseRecorder{
		func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPost, 2, `{"name":"x","status":"DRAFT"}`, groupParams(1))
			CreateClaim(w, r)
			return w
		},
		func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPut, 2, `{"name":"x","status":"DRAFT"}`, claimParams(claim.ID))
			UpdateClaim(w, r)
			return w
		},
		func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodDelete, 2, "", claimParams(claim.ID))
			DeleteClaim(w, r)
			return w
		},
		func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPost, 2, `{"type":"NOTE","date":"2026-10-01","note":"x"}`, claimParams(claim.ID))
			CreateClaimLedgerEntry(w, r)
			return w
		},
	}
	for i, write := range writes {
		if w := write(); w.Code != http.StatusForbidden {
			t.Errorf("write %d: got %d, want 403", i, w.Code)
		}
	}
}

// A non-member and a missing claim get the same answer, so the endpoint cannot
// be used to learn that a claim exists.
func TestClaims_MissingAndForbiddenReadTheSame(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 4, 2, permissions.GroupReceiptsRead)
	claim := seedClaimRow(t, 1, 1, "group one's")

	w, r := claimRequest(http.MethodGet, 4, "", claimParams(claim.ID))
	GetClaimById(w, r)
	forbidden := w

	w, r = claimRequest(http.MethodGet, 4, "", claimParams(9999))
	GetClaimById(w, r)
	missing := w

	if forbidden.Code != http.StatusForbidden || missing.Code != http.StatusForbidden {
		utils.PrintTestError(t, []int{forbidden.Code, missing.Code}, "both 403")
	}
	if forbidden.Body.String() != missing.Body.String() {
		utils.PrintTestError(t, missing.Body.String(), forbidden.Body.String())
	}
}

func TestClaims_CreateInAllGroupIsRejected(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	allGroup := models.Group{Name: "All", IsAllGroup: true}
	repositories.GetDB().Create(&allGroup)
	grantGroupPerms(t, 1, allGroup.ID, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	w, r := claimRequest(http.MethodPost, 1, `{"name":"x","status":"DRAFT"}`, groupParams(allGroup.ID))
	CreateClaim(w, r)

	assertStatus(t, w, http.StatusBadRequest)
	var count int64
	repositories.GetDB().Model(&models.Claim{}).Count(&count)
	if count != 0 {
		t.Errorf("no claim may be created")
	}
}

// The All group lists the claims of every group where the caller may read
// receipts, and nothing from a group where they may not.
func TestClaims_AllGroupListsOnlyReadableGroups(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	allGroup := models.Group{Name: "All", IsAllGroup: true}
	repositories.GetDB().Create(&allGroup)
	grantGroupPerms(t, 1, allGroup.ID, permissions.GroupReceiptsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	grantGroupPerms(t, 1, 2, permissions.GroupReceiptsCreate)

	readable := seedClaimRow(t, 1, 1, "readable")
	seedClaimRow(t, 2, 4, "unreadable")

	list := listClaims(t, 1, allGroup.ID)

	if len(list.Claims) != 1 || list.Claims[0].Id != readable.ID {
		utils.PrintTestError(t, list.Claims, "only the readable group's claim")
	}
	if list.Summary.ClaimCount != 1 {
		utils.PrintTestError(t, list.Summary.ClaimCount, 1)
	}
}

func TestClaims_FiguresAndSummary(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	claim := seedClaimRow(t, 1, 1, "PCS")
	submitted := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repositories.GetDB().Model(&models.Claim{}).Where("id = ?", claim.ID).Updates(map[string]interface{}{
		"status": models.ClaimStatusPartiallyPaid, "submitted_date": submitted,
	})
	first := seedReceiptWithAmount(t, 1, 1, "10.10")
	second := seedReceiptWithAmount(t, 1, 1, "5.25")
	linkReceipts(t, claim.ID, first, second)
	seedLedgerEntry(t, claim.ID, 1, models.ClaimLedgerEntryPayment, "10.00")
	seedLedgerEntry(t, claim.ID, 1, models.ClaimLedgerEntryRecoupment, "2.50")

	list := listClaims(t, 1, 1)
	if len(list.Claims) != 1 {
		t.Fatalf("expected one claim, got %d", len(list.Claims))
	}
	view := list.Claims[0]

	if !reflect.DeepEqual(view.ReceiptIds, []uint{first, second}) || view.ReceiptCount != 2 {
		utils.PrintTestError(t, view.ReceiptIds, []uint{first, second})
	}
	if view.ReceiptsTotal != "15.35" || view.PaidAmount != "7.50" || view.OutstandingAmount != "7.85" {
		utils.PrintTestError(t, view, "15.35 claimed, 7.50 paid, 7.85 outstanding")
	}
	if len(view.LedgerEntries) != 2 {
		utils.PrintTestError(t, len(view.LedgerEntries), 2)
	}
	if list.Summary.AwaitingCount != 1 || list.Summary.OutstandingTotal != "7.85" || list.Summary.OldestDaysWaiting == nil {
		utils.PrintTestError(t, list.Summary, "one awaiting claim, 7.85 outstanding, with an age")
	}
}

// Arrays are [] on the wire, never null: the generated Dart client has no null
// guard. Asserted on the raw body, where a null and a [] differ.
func TestClaims_EmptyCollectionsSerializeAsArrays(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	w, r := claimRequest(http.MethodGet, 1, "", groupParams(1))
	GetClaimsForGroup(w, r)
	if !strings.Contains(w.Body.String(), `"claims":[]`) {
		utils.PrintTestError(t, w.Body.String(), `"claims":[]`)
	}

	claim := seedClaimRow(t, 1, 1, "empty")
	w, r = claimRequest(http.MethodGet, 1, "", claimParams(claim.ID))
	GetClaimById(w, r)
	body := w.Body.String()
	for _, want := range []string{`"receiptIds":[]`, `"ledgerEntries":[]`} {
		if !strings.Contains(body, want) {
			utils.PrintTestError(t, body, want)
		}
	}
}

func TestClaims_AddReceiptsReportsEachRefusal(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	claim := seedClaimRow(t, 1, 1, "target")
	other := seedClaimRow(t, 1, 1, "other")
	fresh := seedReceiptWithAmount(t, 1, 1, "1.00")
	alreadyHere := seedReceiptWithAmount(t, 1, 1, "2.00")
	inOther := seedReceiptWithAmount(t, 1, 1, "3.00")
	otherGroup := seedReceiptWithAmount(t, 2, 4, "4.00")
	const missing uint = 9999
	linkReceipts(t, claim.ID, alreadyHere)
	linkReceipts(t, other.ID, inOther)

	body, _ := json.Marshal(map[string][]uint{"receiptIds": {fresh, alreadyHere, inOther, otherGroup, missing}})
	w, r := claimRequest(http.MethodPost, 1, string(body), claimParams(claim.ID))
	AddReceiptsToClaim(w, r)
	assertStatus(t, w, http.StatusOK)
	var result structs.ClaimReceiptsResult
	decodeInto(t, w, &result)

	if !reflect.DeepEqual(result.ReceiptIds, []uint{fresh, alreadyHere}) {
		utils.PrintTestError(t, result.ReceiptIds, []uint{fresh, alreadyHere})
	}
	wantFailures := []structs.ClaimReceiptsFailure{
		{ReceiptId: inOther, Message: services.ClaimReceiptInOtherClaimMessage},
		{ReceiptId: otherGroup, Message: services.ClaimReceiptNotAddableMessage},
		{ReceiptId: missing, Message: services.ClaimReceiptNotAddableMessage},
	}
	if !reflect.DeepEqual(result.Failures, wantFailures) {
		utils.PrintTestError(t, result.Failures, wantFailures)
	}

	if got := getClaim(t, 1, claim.ID).ReceiptIds; !reflect.DeepEqual(got, []uint{fresh, alreadyHere}) {
		utils.PrintTestError(t, got, []uint{fresh, alreadyHere})
	}
	if got := getClaim(t, 1, other.ID).ReceiptIds; !reflect.DeepEqual(got, []uint{inOther}) {
		t.Errorf("a receipt in another claim must not be moved")
	}
}

func TestClaims_RemoveReceipts(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	claim := seedClaimRow(t, 1, 1, "target")
	other := seedClaimRow(t, 1, 1, "other")
	linked := seedReceiptWithAmount(t, 1, 1, "1.00")
	inOther := seedReceiptWithAmount(t, 1, 1, "2.00")
	linkReceipts(t, claim.ID, linked)
	linkReceipts(t, other.ID, inOther)

	body, _ := json.Marshal(map[string][]uint{"receiptIds": {linked, inOther}})
	w, r := claimRequest(http.MethodPost, 1, string(body), claimParams(claim.ID))
	RemoveReceiptsFromClaim(w, r)
	assertStatus(t, w, http.StatusOK)
	var result structs.ClaimReceiptsResult
	decodeInto(t, w, &result)

	if !reflect.DeepEqual(result.ReceiptIds, []uint{linked}) {
		utils.PrintTestError(t, result.ReceiptIds, []uint{linked})
	}
	if !reflect.DeepEqual(result.Failures, []structs.ClaimReceiptsFailure{{ReceiptId: inOther, Message: services.ClaimReceiptNotInClaimMessage}}) {
		utils.PrintTestError(t, result.Failures, "the other claim's receipt refused")
	}
	if got := getClaim(t, 1, other.ID).ReceiptIds; !reflect.DeepEqual(got, []uint{inOther}) {
		t.Errorf("removing from one claim must not touch another")
	}
	if !receiptExists(t, linked) {
		t.Errorf("removing a receipt from a claim must not delete it")
	}
}

func TestClaims_InvalidBodiesAre400(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)
	claim := seedClaimRow(t, 1, 1, "target")

	cases := map[string]func() *httptest.ResponseRecorder{
		"unknown status": func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPost, 1, `{"name":"x","status":"LOST"}`, groupParams(1))
			CreateClaim(w, r)
			return w
		},
		"malformed json": func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPut, 1, `{"name":`, claimParams(claim.ID))
			UpdateClaim(w, r)
			return w
		},
		"non-numeric amount": func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPost, 1, `{"name":"x","status":"DRAFT","claimedAmount":"abc"}`, groupParams(1))
			CreateClaim(w, r)
			return w
		},
		"payment without amount": func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPost, 1, `{"type":"PAYMENT","date":"2026-10-01"}`, claimParams(claim.ID))
			CreateClaimLedgerEntry(w, r)
			return w
		},
		"empty receipt list": func() *httptest.ResponseRecorder {
			w, r := claimRequest(http.MethodPost, 1, `{"receiptIds":[]}`, claimParams(claim.ID))
			AddReceiptsToClaim(w, r)
			return w
		},
	}

	for name, run := range cases {
		if w := run(); w.Code != http.StatusBadRequest {
			t.Errorf("%s: got %d, want 400 (%s)", name, w.Code, w.Body.String())
		}
	}
}

func TestClaims_LedgerEntryLifecycle(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)
	claim := seedClaimRow(t, 1, 1, "target")
	otherClaim := seedClaimRow(t, 1, 1, "other")

	w, r := claimRequest(http.MethodPost, 1, `{"type":"PAYMENT","date":"2026-10-01","amount":"100","note":"EFT"}`, claimParams(claim.ID))
	CreateClaimLedgerEntry(w, r)
	assertStatus(t, w, http.StatusOK)
	var entry structs.ClaimLedgerEntryView
	decodeInto(t, w, &entry)
	if entry.Type != "PAYMENT" || entry.Date != "2026-10-01" || entry.Amount != "100.00" || entry.ClaimId != claim.ID {
		utils.PrintTestError(t, entry, "the payment as sent")
	}

	entryParams := map[string]string{"claimId": utils.UintToString(claim.ID), "ledgerEntryId": utils.UintToString(entry.Id)}
	w, r = claimRequest(http.MethodPut, 1, `{"type":"RECOUPMENT","date":"2026-10-02","amount":"25.5"}`, entryParams)
	UpdateClaimLedgerEntry(w, r)
	assertStatus(t, w, http.StatusOK)
	var updated structs.ClaimLedgerEntryView
	decodeInto(t, w, &updated)
	if updated.Type != "RECOUPMENT" || updated.Amount != "25.50" || updated.Note != "" {
		utils.PrintTestError(t, updated, "the recoupment with its note cleared")
	}

	// The entry addressed through another claim is a 403, not an edit.
	wrongParams := map[string]string{"claimId": utils.UintToString(otherClaim.ID), "ledgerEntryId": utils.UintToString(entry.Id)}
	w, r = claimRequest(http.MethodDelete, 1, "", wrongParams)
	DeleteClaimLedgerEntry(w, r)
	assertStatus(t, w, http.StatusForbidden)

	w, r = claimRequest(http.MethodDelete, 1, "", entryParams)
	DeleteClaimLedgerEntry(w, r)
	assertStatus(t, w, http.StatusOK)
	if got := getClaim(t, 1, claim.ID).LedgerEntries; len(got) != 0 {
		utils.PrintTestError(t, got, "no entries")
	}
}

// A member whose role hides some payers sees only the receipts they may see, both
// in the ids and in every figure, and can neither add nor remove a hidden one.
func TestClaims_PaidByHiddenReceiptsAreLeftOut(t *testing.T) {
	defer tearDownGenericHandlerTest()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	repositories.CreateTestGroupWithUsers()
	db := repositories.GetDB()

	// User 2 sees only receipts paid by user 1.
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"Claims PaidBy Restricted", "",
		[]string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate},
		nil, nil, []uint{1}, false, false,
	)
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	db.Model(&models.GroupMember{}).Where("user_id = ? AND group_id = ?", 2, 1).Update("group_role_id", role.ID)

	claim := seedClaimRow(t, 1, 1, "shared")
	visible := seedReceiptWithAmount(t, 1, 1, "10.00")
	hidden := seedReceiptWithAmount(t, 1, 3, "90.00")
	linkReceipts(t, claim.ID, visible, hidden)

	view := getClaim(t, 2, claim.ID)
	if !reflect.DeepEqual(view.ReceiptIds, []uint{visible}) || view.ReceiptsTotal != "10.00" || view.EffectiveClaimedAmount != "10.00" {
		utils.PrintTestError(t, view, "only the visible receipt counted")
	}

	body, _ := json.Marshal(map[string][]uint{"receiptIds": {hidden}})
	w, r := claimRequest(http.MethodPost, 2, string(body), claimParams(claim.ID))
	RemoveReceiptsFromClaim(w, r)
	var result structs.ClaimReceiptsResult
	decodeInto(t, w, &result)
	if len(result.ReceiptIds) != 0 || len(result.Failures) != 1 || result.Failures[0].Message != services.ClaimReceiptNotInClaimMessage {
		utils.PrintTestError(t, result, "the hidden receipt reads as not in the claim")
	}

	otherHidden := seedReceiptWithAmount(t, 1, 3, "5.00")
	body, _ = json.Marshal(map[string][]uint{"receiptIds": {otherHidden}})
	w, r = claimRequest(http.MethodPost, 2, string(body), claimParams(claim.ID))
	AddReceiptsToClaim(w, r)
	decodeInto(t, w, &result)
	if len(result.Failures) != 1 || result.Failures[0].Message != services.ClaimReceiptNotAddableMessage {
		utils.PrintTestError(t, result, "a hidden receipt cannot be added")
	}

	// The unrestricted creator still sees both.
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	if got := getClaim(t, 1, claim.ID).ReceiptIds; !reflect.DeepEqual(got, []uint{visible, hidden}) {
		utils.PrintTestError(t, got, []uint{visible, hidden})
	}
}

// In a group with isolated members, a plain member does not see a claim or a
// ledger entry recorded by another plain member.
func TestClaims_IsolatedGroupHidesOtherMembersRecords(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	repositories.GetDB().Model(&models.Group{}).Where("id = ?", 1).Update("isolate_members", true)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)
	grantGroupPerms(t, 2, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	mine := seedClaimRow(t, 1, 1, "mine")
	theirs := seedClaimRow(t, 1, 2, "theirs")
	seedLedgerEntry(t, mine.ID, 1, models.ClaimLedgerEntryPayment, "10.00")
	seedLedgerEntry(t, mine.ID, 2, models.ClaimLedgerEntryPayment, "500.00")

	list := listClaims(t, 1, 1)
	if len(list.Claims) != 1 || list.Claims[0].Id != mine.ID {
		utils.PrintTestError(t, list.Claims, "only the caller's own claim")
	}
	if len(list.Claims) == 1 && (len(list.Claims[0].LedgerEntries) != 1 || list.Claims[0].PaidAmount != "10.00") {
		utils.PrintTestError(t, list.Claims[0], "only the caller's own entry, and only it counted")
	}

	w, r := claimRequest(http.MethodGet, 1, "", claimParams(theirs.ID))
	GetClaimById(w, r)
	assertStatus(t, w, http.StatusForbidden)
}

func TestDeleteReceipt_RefusedWhileInAClaim(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate, permissions.GroupReceiptsDelete)
	claim := seedClaimRow(t, 1, 1, "Rome PCS")
	receiptId := seedReceiptWithAmount(t, 1, 1, "1.00")
	linkReceipts(t, claim.ID, receiptId)

	w, r := deleteReceiptRequest(1, receiptId)
	DeleteReceipt(w, r)

	assertStatus(t, w, http.StatusBadRequest)
	if !strings.Contains(w.Body.String(), "Rome PCS") {
		utils.PrintTestError(t, w.Body.String(), "a message naming the claim")
	}
	if !receiptExists(t, receiptId) {
		t.Fatalf("the receipt must survive")
	}

	body, _ := json.Marshal(map[string][]uint{"receiptIds": {receiptId}})
	w, r = claimRequest(http.MethodPost, 1, string(body), claimParams(claim.ID))
	RemoveReceiptsFromClaim(w, r)
	assertStatus(t, w, http.StatusOK)

	w, r = deleteReceiptRequest(1, receiptId)
	DeleteReceipt(w, r)
	assertStatus(t, w, http.StatusOK)
	if receiptExists(t, receiptId) {
		t.Errorf("once out of the claim the receipt deletes normally")
	}
}

func TestBulkDeleteReceipts_ReceiptInAClaimIsAFailure(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)
	claim := seedClaimRow(t, 1, 1, "Rome PCS")
	claimed := seedReceipt(t, 1, 1, 0)
	free := seedReceipt(t, 1, 1, 0)
	linkReceipts(t, claim.ID, claimed)

	result := runBulkDelete(t, 1, claimed, free)

	if !reflect.DeepEqual(result.DeletedReceiptIds, []uint{free}) {
		utils.PrintTestError(t, result.DeletedReceiptIds, []uint{free})
	}
	if len(result.Failures) != 1 || result.Failures[0].ReceiptId != claimed || !strings.Contains(result.Failures[0].Message, "Rome PCS") {
		utils.PrintTestError(t, result.Failures, "the claimed receipt refused, naming its claim")
	}
	if !receiptExists(t, claimed) {
		t.Errorf("the claimed receipt must survive")
	}
}

func deleteReceiptRequest(userId uint, receiptId uint) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/api/receipt", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", utils.UintToString(receiptId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}
