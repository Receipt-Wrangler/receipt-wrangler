package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"reflect"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
)

func bulkDeleteReceiptsRequest(userId uint, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/receipt/bulkDelete", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

func bulkDeleteBody(t *testing.T, receiptIds ...uint) string {
	t.Helper()
	bytes, err := json.Marshal(map[string][]uint{"receiptIds": receiptIds})
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	return string(bytes)
}

func runBulkDelete(t *testing.T, userId uint, receiptIds ...uint) structs.BulkReceiptDeleteResult {
	t.Helper()
	w, r := bulkDeleteReceiptsRequest(userId, bulkDeleteBody(t, receiptIds...))
	BulkDeleteReceipts(w, r)
	assertStatus(t, w, http.StatusOK)

	var result structs.BulkReceiptDeleteResult
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode result %q: %v", w.Body.String(), err)
	}
	return result
}

func receiptExists(t *testing.T, receiptId uint) bool {
	t.Helper()
	var count int64
	repositories.GetDB().Model(&models.Receipt{}).Where("id = ?", receiptId).Count(&count)
	return count > 0
}

func deniedFailure(receiptId uint) structs.BulkReceiptDeleteFailure {
	return structs.BulkReceiptDeleteFailure{ReceiptId: receiptId, Message: services.BulkReceiptDeleteDeniedMessage}
}

func TestBulkDeleteReceipts_DeletesEveryPermittedReceipt(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)
	first := seedReceipt(t, 1, 1, 0)
	second := seedReceipt(t, 1, 2, 0)

	result := runBulkDelete(t, 1, first, second)

	if !reflect.DeepEqual(result.DeletedReceiptIds, []uint{first, second}) {
		utils.PrintTestError(t, result.DeletedReceiptIds, []uint{first, second})
	}
	if len(result.Failures) != 0 {
		utils.PrintTestError(t, result.Failures, "no failures")
	}
	if receiptExists(t, first) || receiptExists(t, second) {
		t.Errorf("expected both receipts to be deleted")
	}
}

// The case the endpoint exists for: a selection spanning groups (the All view)
// where the caller may delete in only some of them. The permitted receipt is
// deleted, every other one is reported in request order, and a receipt that does
// not exist reads exactly like one the caller may not delete.
func TestBulkDeleteReceipts_ReportsEachReceiptItCannotDelete(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)
	grantGroupPerms(t, 1, 2, permissions.GroupReceiptsRead)
	readOnly := seedReceipt(t, 2, 1, 0)
	deletable := seedReceipt(t, 1, 1, 0)
	const missing uint = 9999

	result := runBulkDelete(t, 1, readOnly, deletable, missing)

	if !reflect.DeepEqual(result.DeletedReceiptIds, []uint{deletable}) {
		utils.PrintTestError(t, result.DeletedReceiptIds, []uint{deletable})
	}
	wantFailures := []structs.BulkReceiptDeleteFailure{deniedFailure(readOnly), deniedFailure(missing)}
	if !reflect.DeepEqual(result.Failures, wantFailures) {
		utils.PrintTestError(t, result.Failures, wantFailures)
	}
	if !receiptExists(t, readOnly) {
		t.Errorf("receipt in the read-only group must survive")
	}
	if receiptExists(t, deletable) {
		t.Errorf("permitted receipt must be deleted")
	}
}

func TestBulkDeleteReceipts_NothingDeletableStillAnswersOk(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	// User 4 belongs only to group 2 and holds no role there.
	receiptId := seedReceipt(t, 1, 1, 0)

	result := runBulkDelete(t, 4, receiptId)

	if len(result.DeletedReceiptIds) != 0 {
		utils.PrintTestError(t, result.DeletedReceiptIds, "nothing deleted")
	}
	if !reflect.DeepEqual(result.Failures, []structs.BulkReceiptDeleteFailure{deniedFailure(receiptId)}) {
		utils.PrintTestError(t, result.Failures, deniedFailure(receiptId))
	}
	if !receiptExists(t, receiptId) {
		t.Errorf("receipt must survive")
	}
}

// Both lists must be [] rather than null: the generated Dart client has no null guard.
// Asserted on the raw body, since a decoded null and a decoded [] look the same.
func TestBulkDeleteReceipts_EmptyListsSerializeAsArrays(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)
	receiptId := seedReceipt(t, 1, 1, 0)
	const missing uint = 9999

	cases := map[string]struct {
		receiptId uint
		want      string
	}{
		"nothing refused": {receiptId: receiptId, want: `"failures":[]`},
		"nothing deleted": {receiptId: missing, want: `"deletedReceiptIds":[]`},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			w, r := bulkDeleteReceiptsRequest(1, bulkDeleteBody(t, testCase.receiptId))
			BulkDeleteReceipts(w, r)

			assertStatus(t, w, http.StatusOK)
			if !strings.Contains(w.Body.String(), testCase.want) {
				utils.PrintTestError(t, w.Body.String(), testCase.want)
			}
		})
	}
}

// A caller's membership in the synthetic All group carries an unrestricted role,
// but no receipt lives there: each receipt is judged in its own group.
func TestBulkDeleteReceipts_AllGroupMembershipGrantsNothing(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	allGroup := models.Group{Name: "All", IsAllGroup: true}
	if err := repositories.GetDB().Create(&allGroup).Error; err != nil {
		t.Fatalf("seed all group: %v", err)
	}
	grantGroupPerms(t, 1, allGroup.ID, permissions.GroupReceiptsRead, permissions.GroupReceiptsDelete)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	receiptId := seedReceipt(t, 1, 1, 0)

	result := runBulkDelete(t, 1, receiptId)

	if !reflect.DeepEqual(result.Failures, []structs.BulkReceiptDeleteFailure{deniedFailure(receiptId)}) {
		utils.PrintTestError(t, result.Failures, deniedFailure(receiptId))
	}
	if !receiptExists(t, receiptId) {
		t.Errorf("receipt must survive")
	}
}

// The single delete refuses a receipt hidden by the caller's paid-by visibility
// (HandleRequest's enforcePaidByVisibility); the bulk delete must too.
func TestBulkDeleteReceipts_RefusesReceiptHiddenByPaidBy(t *testing.T) {
	defer tearDownGenericHandlerTest()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	payer := models.User{Username: "visible-payer", Password: "x"}
	db.Create(&payer)
	hidden := models.User{Username: "hidden-payer", Password: "x"}
	db.Create(&hidden)
	member := models.User{Username: "paidby-deleter", Password: "x"}
	db.Create(&member)
	group := models.Group{Name: "paidby-bulk-delete"}
	db.Create(&group)

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"PaidBy Restricted Deleter",
		"",
		[]string{permissions.GroupReceiptsRead, permissions.GroupReceiptsDelete},
		nil,
		nil,
		[]uint{payer.ID},
		false, false,
	)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}
	db.Create(&models.GroupMember{GroupID: group.ID, UserID: member.ID, GroupRoleID: &role.ID})

	visibleReceipt := seedReceipt(t, group.ID, payer.ID, 0)
	hiddenReceipt := seedReceipt(t, group.ID, hidden.ID, 0)

	result := runBulkDelete(t, member.ID, visibleReceipt, hiddenReceipt)

	if !reflect.DeepEqual(result.DeletedReceiptIds, []uint{visibleReceipt}) {
		utils.PrintTestError(t, result.DeletedReceiptIds, []uint{visibleReceipt})
	}
	if !reflect.DeepEqual(result.Failures, []structs.BulkReceiptDeleteFailure{deniedFailure(hiddenReceipt)}) {
		utils.PrintTestError(t, result.Failures, deniedFailure(hiddenReceipt))
	}
	if !receiptExists(t, hiddenReceipt) {
		t.Errorf("hidden receipt must survive")
	}
}

func TestBulkDeleteReceipts_RejectsInvalidLists(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)
	receiptId := seedReceipt(t, 1, 1, 0)

	cases := map[string]string{
		"empty":     `{"receiptIds":[]}`,
		"missing":   `{}`,
		"duplicate": bulkDeleteBody(t, receiptId, receiptId),
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			w, r := bulkDeleteReceiptsRequest(1, body)
			BulkDeleteReceipts(w, r)

			assertStatus(t, w, http.StatusBadRequest)
			var errs map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &errs); err != nil || errs["receiptIds"] == "" {
				utils.PrintTestError(t, w.Body.String(), "a receiptIds validation error")
			}
		})
	}

	if !receiptExists(t, receiptId) {
		t.Errorf("an invalid request must delete nothing")
	}
}

func TestBulkDeleteReceipts_RejectsMalformedBody(t *testing.T) {
	defer tearDownGenericHandlerTest()

	w, r := bulkDeleteReceiptsRequest(1, `{"receiptIds":["1"]}`)
	BulkDeleteReceipts(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

// Run with foreign keys OFF, so the comment and custom-field value can only be
// gone because the delete removed them, not because a cascade did.
func TestBulkDeleteReceipts_RemovesTheReceiptsOwnRows(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsDelete)

	customField := models.CustomField{Name: "Note", Type: models.TEXT}
	if err := repositories.GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}
	receiptId := seedReceipt(t, 1, 1, customField.ID)
	userId := uint(1)
	repositories.GetDB().Create(&models.Comment{Comment: "keep me?", ReceiptId: receiptId, UserId: &userId})

	restore, err := repositories.UseTestDbWithoutForeignKeysForTests()
	if err != nil {
		t.Fatalf("open db without foreign keys: %v", err)
	}
	defer restore()

	result := runBulkDelete(t, 1, receiptId)
	if !reflect.DeepEqual(result.DeletedReceiptIds, []uint{receiptId}) {
		utils.PrintTestError(t, result.DeletedReceiptIds, []uint{receiptId})
	}

	db := repositories.GetDB()
	for name, model := range map[string]interface{}{
		"comments":            &models.Comment{},
		"custom field values": &models.CustomFieldValue{},
	} {
		var count int64
		db.Model(model).Where("receipt_id = ?", receiptId).Count(&count)
		if count != 0 {
			t.Errorf("%s left behind: %d", name, count)
		}
	}
}
