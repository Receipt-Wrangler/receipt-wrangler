package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

func callGetReceiptSummary(t *testing.T, userId uint, groupId uint, body string) *httptest.ResponseRecorder {
	t.Helper()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api", bytes.NewReader([]byte(body)))

	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("groupId", utils.UintToString(groupId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiContext))
	r = r.WithContext(context.WithValue(
		r.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}},
	))

	GetReceiptSummaryForGroup(w, r)

	return w
}

func enableSummaryForGroup(t *testing.T, groupId uint, statuses []models.ReceiptStatus) {
	t.Helper()

	repository := repositories.NewGroupReceiptSettingsRepository(nil)
	if _, err := repository.GetGroupReceiptSettingsByGroupId(groupId); err != nil {
		if _, err := repository.CreateGroupReceiptSettings(groupId); err != nil {
			t.Fatalf("create settings: %v", err)
		}
	}

	if err := repositories.GetDB().Model(&models.GroupReceiptSettings{}).
		Where("group_id = ?", groupId).
		Update("receipt_summary_enabled", true).Error; err != nil {
		t.Fatalf("enable summary: %v", err)
	}

	for _, status := range statuses {
		row := models.GroupReceiptSettingsSummaryStatus{GroupId: groupId, Status: status}
		if err := repositories.GetDB().Create(&row).Error; err != nil {
			t.Fatalf("seed summary status: %v", err)
		}
	}
}

func seedSummaryReceiptForGroup(t *testing.T, groupId uint, paidByUserId uint, amount string) {
	t.Helper()

	receipt := models.Receipt{
		Name:         "summary receipt",
		Amount:       decimal.RequireFromString(amount),
		Date:         time.Now(),
		PaidByUserID: paidByUserId,
		GroupId:      groupId,
		Status:       models.OPEN,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
}

func decodeSummary(t *testing.T, w *httptest.ResponseRecorder) structs.ReceiptSummary {
	t.Helper()
	summary := structs.ReceiptSummary{}
	if err := json.Unmarshal(w.Body.Bytes(), &summary); err != nil {
		t.Fatalf("decode summary: %v (body %s)", err, w.Body.String())
	}
	return summary
}

func TestGetReceiptSummaryRequiresReceiptsRead(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestGroupWithUsers()
	enableSummaryForGroup(t, 1, []models.ReceiptStatus{models.OPEN})

	// User 1 is a member of group 1 but holds no receipts.read.
	w := callGetReceiptSummary(t, 1, 1, `{"filter":{}}`)
	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}
}

func TestGetReceiptSummaryReturnsTotals(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	enableSummaryForGroup(t, 1, []models.ReceiptStatus{models.OPEN})

	seedSummaryReceiptForGroup(t, 1, 1, "10.50")
	seedSummaryReceiptForGroup(t, 1, 1, "4.25")

	w := callGetReceiptSummary(t, 1, 1, `{"filter":{}}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	summary := decodeSummary(t, w)
	if !summary.Enabled || summary.Overall.ReceiptCount != 2 {
		utils.PrintTestError(t, summary.Overall, "enabled, 2 receipts")
	}
	if !summary.Overall.Total.Equal(decimal.RequireFromString("14.75")) {
		utils.PrintTestError(t, summary.Overall.Total.String(), "14.75")
	}
}

// TestGetReceiptSummarySerializesMoneyAsString pins the wire contract: swagger declares these as
// strings, matching Receipt.amount, and the generated clients parse them as such.
func TestGetReceiptSummarySerializesMoneyAsString(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	enableSummaryForGroup(t, 1, nil)
	seedSummaryReceiptForGroup(t, 1, 1, "10.50")

	w := callGetReceiptSummary(t, 1, 1, `{"filter":{}}`)

	raw := map[string]interface{}{}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		utils.PrintTestError(t, err, "a JSON object")
		return
	}

	overall, _ := raw["overall"].(map[string]interface{})
	if _, isString := overall["total"].(string); !isString {
		utils.PrintTestError(t, overall["total"], "total as a JSON string")
	}

	// And the array-valued keys must be [] rather than null, or a released Dart client fails the
	// whole payload.
	if statuses, ok := raw["statuses"].([]interface{}); !ok || statuses == nil {
		utils.PrintTestError(t, raw["statuses"], "statuses as []")
	}
	if totals, ok := overall["customFieldTotals"].([]interface{}); !ok || totals == nil {
		utils.PrintTestError(t, overall["customFieldTotals"], "customFieldTotals as []")
	}
}

// TestGetReceiptSummaryExcludesHiddenPayersReceipts is the regression that matters most here: the
// totals must describe exactly the rows the viewer can see. A receipt hidden by paid-by visibility
// leaking into a SUM would disclose spend the viewer is not allowed to know about.
func TestGetReceiptSummaryExcludesHiddenPayersReceipts(t *testing.T) {
	defer repositories.TruncateTestDb()

	payer := models.User{Username: "summary-allowed-payer", Password: "x"}
	repositories.GetDB().Create(&payer)
	hidden := models.User{Username: "summary-hidden-payer", Password: "x"}
	repositories.GetDB().Create(&hidden)

	userId, groupId := seedPaidByRestrictedMember(t, []uint{payer.ID}, false)
	enableSummaryForGroup(t, groupId, []models.ReceiptStatus{models.OPEN})

	seedSummaryReceiptForGroup(t, groupId, payer.ID, "10.00")
	seedSummaryReceiptForGroup(t, groupId, hidden.ID, "999.00")

	w := callGetReceiptSummary(t, userId, groupId, `{"filter":{}}`)
	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	summary := decodeSummary(t, w)
	if summary.Overall.ReceiptCount != 1 {
		utils.PrintTestError(t, summary.Overall.ReceiptCount, 1)
	}
	if !summary.Overall.Total.Equal(decimal.RequireFromString("10.00")) {
		utils.PrintTestError(t, summary.Overall.Total.String(), "10.00 (the hidden payer excluded)")
	}
}

// TestGetReceiptSummaryForbidsUnreadableConfigurationGroup: the All-group chip row lets the client
// name which group's configuration to apply, so that id must be authorized or it becomes a way to
// enumerate another group's configured custom field names.
func TestGetReceiptSummaryForbidsUnreadableConfigurationGroup(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	// Group 2 exists; user 1 is not a member and holds nothing on it.
	enableSummaryForGroup(t, 2, []models.ReceiptStatus{models.OPEN})

	w := callGetReceiptSummary(t, 1, 1, `{"configurationGroupId":2,"filter":{}}`)
	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}
}
