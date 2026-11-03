package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

func tearDownReceiptHistoryTest() {
	repositories.TruncateTestDb()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
}

// receiptHistoryRequest builds a request the way the router invokes the
// history handlers: the receipt id (and the entry's task id, for the
// snapshots) as chi URL params, and claims for userId.
func receiptHistoryRequest(userId uint, receiptId string, systemTaskId string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/receipt/"+receiptId+"/history", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", receiptId)
	if systemTaskId != "" {
		rctx.URLParams.Add("systemTaskId", systemTaskId)
	}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

// seedReceiptWithHistory creates a receipt in group 1 paid by user 1 through
// CreateReceipt, then renames it through the recorder, so it has a creation
// and an update.
func seedReceiptWithHistory(t *testing.T) models.Receipt {
	t.Helper()
	receipt, err := repositories.NewReceiptRepository(nil).CreateReceipt(commands.UpsertReceiptCommand{
		Name:         "Groceries",
		Amount:       decimal.RequireFromString("12.00"),
		Date:         time.Now(),
		GroupId:      1,
		PaidByUserID: 1,
		Status:       models.OPEN,
	}, 1, true)
	if err != nil {
		t.Fatalf("create receipt: %v", err)
	}

	err = repositories.NewReceiptChangeRecorder(1, receipt.ID).Record(func() error {
		return repositories.GetDB().Model(&models.Receipt{}).Where("id = ?", receipt.ID).Update("name", "Weekly shop").Error
	})
	if err != nil {
		t.Fatalf("rename receipt: %v", err)
	}
	return receipt
}

func assertForbiddenBody(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	assertStatus(t, w, http.StatusForbidden)
	if !strings.Contains(w.Body.String(), unauthorizedEntityMessage) {
		t.Errorf("body = %s, want the standard denial", w.Body.String())
	}
}

func TestGetReceiptHistoryAuthorization(t *testing.T) {
	tests := []struct {
		name      string
		appPerms  []string
		groupPerm []string
		receiptId func(receipt models.Receipt) string
		want      int
	}{
		{
			name:      "allowed with the app permission and group read",
			appPerms:  []string{permissions.AppReceiptHistoryRead},
			groupPerm: []string{permissions.GroupReceiptsRead},
			want:      http.StatusOK,
		},
		{
			name:      "denied without the app permission",
			appPerms:  []string{permissions.AppReceiptsSearch},
			groupPerm: []string{permissions.GroupReceiptsRead},
			want:      http.StatusForbidden,
		},
		{
			name:      "denied without group read, even for an app-wide reader",
			appPerms:  []string{permissions.AppReceiptHistoryRead, permissions.AppGroupsRead},
			groupPerm: []string{permissions.GroupReceiptsUpdate},
			want:      http.StatusForbidden,
		},
		{
			name:      "denied for a missing receipt",
			appPerms:  []string{permissions.AppReceiptHistoryRead},
			groupPerm: []string{permissions.GroupReceiptsRead},
			receiptId: func(models.Receipt) string { return "404" },
			want:      http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer tearDownReceiptHistoryTest()
			repositories.CreateTestGroupWithUsers()
			receipt := seedReceiptWithHistory(t)
			grantAppPerms(t, 2, test.appPerms...)
			grantGroupPerms(t, 2, 1, test.groupPerm...)

			receiptId := utils.UintToString(receipt.ID)
			if test.receiptId != nil {
				receiptId = test.receiptId(receipt)
			}

			for _, handler := range []struct {
				name   string
				taskId string
				call   func(w http.ResponseWriter, r *http.Request)
			}{
				{"history", "", GetReceiptHistory},
				{"snapshots", "1", GetReceiptHistorySnapshots},
			} {
				w, r := receiptHistoryRequest(2, receiptId, handler.taskId)
				handler.call(w, r)
				if test.want == http.StatusForbidden {
					assertForbiddenBody(t, w)
				} else {
					assertStatus(t, w, test.want)
				}
			}
		})
	}
}

func TestGetReceiptHistoryDeniedWhenPaidByHidesTheReceipt(t *testing.T) {
	defer tearDownReceiptHistoryTest()

	payer := models.User{Username: "history-payer", Password: "x"}
	repositories.GetDB().Create(&payer)
	userId, groupId := seedPaidByRestrictedMember(t, nil, true)
	grantAppPerms(t, userId, permissions.AppReceiptHistoryRead)

	receipt := models.Receipt{Name: "Someone else's", GroupId: groupId, PaidByUserID: payer.ID, Date: time.Now()}
	repositories.GetDB().Create(&receipt)

	w, r := receiptHistoryRequest(userId, utils.UintToString(receipt.ID), "")
	GetReceiptHistory(w, r)
	assertForbiddenBody(t, w)
}

func TestGetReceiptHistoryResponseShape(t *testing.T) {
	defer tearDownReceiptHistoryTest()
	repositories.CreateTestGroupWithUsers()
	receipt := seedReceiptWithHistory(t)
	grantAppPerms(t, 2, permissions.AppReceiptHistoryRead)
	grantGroupPerms(t, 2, 1, permissions.GroupReceiptsRead)

	w, r := receiptHistoryRequest(2, utils.UintToString(receipt.ID), "")
	GetReceiptHistory(w, r)
	assertStatus(t, w, http.StatusOK)

	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("content type = %q", contentType)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("read body %s: %v", w.Body.String(), err)
	}
	for _, key := range []string{"entries", "changeCount", "saveCount", "contributorCount"} {
		if _, ok := body[key]; !ok {
			t.Errorf("response is missing %q: %s", key, w.Body.String())
		}
	}

	entries, _ := body["entries"].([]interface{})
	if len(entries) != 2 {
		t.Fatalf("entries = %v, want the rename and the creation", body["entries"])
	}
	update := entries[0].(map[string]interface{})
	if update["action"] != "UPDATED" || update["hasSnapshots"] != true || update["beforeState"] != "COMPLETE" {
		t.Errorf("update entry = %v", update)
	}
	changes, _ := update["changes"].([]interface{})
	if len(changes) != 1 {
		t.Fatalf("changes = %v, want the rename", update["changes"])
	}
	change := changes[0].(map[string]interface{})
	if change["field"] != "NAME" || change["kind"] != "CHANGED" ||
		change["from"].(map[string]interface{})["value"] != "Groceries" || change["to"].(map[string]interface{})["value"] != "Weekly shop" {
		t.Errorf("change = %v", change)
	}
	if actor := update["actor"].(map[string]interface{}); actor["type"] != "USER" || actor["userId"] != float64(1) {
		t.Errorf("actor = %v", actor)
	}
}

func TestGetReceiptHistorySnapshotsEndpoint(t *testing.T) {
	defer tearDownReceiptHistoryTest()
	repositories.CreateTestGroupWithUsers()
	receipt := seedReceiptWithHistory(t)
	grantAppPerms(t, 2, permissions.AppReceiptHistoryRead)
	grantGroupPerms(t, 2, 1, permissions.GroupReceiptsRead)

	var update models.SystemTask
	repositories.GetDB().Where("type = ? AND receipt_id = ?", models.RECEIPT_UPDATED, receipt.ID).First(&update)

	w, r := receiptHistoryRequest(2, utils.UintToString(receipt.ID), utils.UintToString(update.ID))
	GetReceiptHistorySnapshots(w, r)
	assertStatus(t, w, http.StatusOK)

	var snapshots struct {
		Before string `json:"before"`
		After  string `json:"after"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &snapshots); err != nil {
		t.Fatalf("read body: %v", err)
	}
	var before, after models.Receipt
	if json.Unmarshal([]byte(snapshots.Before), &before) != nil || json.Unmarshal([]byte(snapshots.After), &after) != nil ||
		before.Name != "Groceries" || after.Name != "Weekly shop" {
		t.Errorf("snapshots = %+v, want Groceries -> Weekly shop", snapshots)
	}

	for _, taskId := range []string{"9999", "not-a-number"} {
		w, r = receiptHistoryRequest(2, utils.UintToString(receipt.ID), taskId)
		GetReceiptHistorySnapshots(w, r)
		assertForbiddenBody(t, w)
	}
}

// recordedReceiptChanges returns the RECEIPT_UPDATED rows recorded for a
// receipt, oldest first. Comments, images and the bulk status update record one
// each through the ReceiptChangeRecorder; these helpers let their own handler
// tests check it.
func recordedReceiptChanges(t *testing.T, receiptId uint) []models.SystemTask {
	t.Helper()
	var tasks []models.SystemTask
	err := repositories.GetDB().
		Where("type = ? AND associated_entity_type = ? AND associated_entity_id = ?", models.RECEIPT_UPDATED, models.RECEIPT, receiptId).
		Order("id").
		Find(&tasks).Error
	if err != nil {
		t.Fatalf("load recorded changes: %v", err)
	}
	return tasks
}

// assertOneRecordedChange checks that exactly one successful change was
// recorded for the receipt, run by userId, and returns its snapshots.
func assertOneRecordedChange(t *testing.T, receiptId uint, userId uint) (before models.Receipt, after models.Receipt) {
	t.Helper()
	tasks := recordedReceiptChanges(t, receiptId)
	if len(tasks) != 1 {
		t.Fatalf("recorded %d changes for receipt %d, want 1", len(tasks), receiptId)
	}
	task := tasks[0]
	if task.Status != models.SYSTEM_TASK_SUCCEEDED || task.RanByUserId == nil || *task.RanByUserId != userId ||
		task.ReceiptId == nil || *task.ReceiptId != receiptId || task.GroupId == nil {
		t.Fatalf("recorded change = %+v, want a succeeded row for receipt %d run by %d", task, receiptId, userId)
	}

	var description struct {
		Before  string `json:"before"`
		After   string `json:"after"`
		Version int    `json:"version"`
	}
	if err := json.Unmarshal([]byte(task.ResultDescription), &description); err != nil {
		t.Fatalf("read description: %v", err)
	}
	if description.Version != repositories.ReceiptUpdateDescriptionVersion {
		t.Errorf("version = %d, want %d", description.Version, repositories.ReceiptUpdateDescriptionVersion)
	}
	if err := json.Unmarshal([]byte(description.Before), &before); err != nil {
		t.Fatalf("read before: %v", err)
	}
	if err := json.Unmarshal([]byte(description.After), &after); err != nil {
		t.Fatalf("read after: %v", err)
	}
	return before, after
}
