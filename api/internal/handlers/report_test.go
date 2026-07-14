package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/shopspring/decimal"
)

func tearDownReportTest() {
	repositories.TruncateTestDb()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
}

// generateReportRequest builds a POST carrying the JSON command body and JWT
// claims for userId, mirroring how the router invokes GenerateReport.
func generateReportRequest(userId uint, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/report/generate", strings.NewReader(body))
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}})
	return w, r.WithContext(newContext)
}

func seedReportReceipt(groupId uint, userId uint) {
	repositories.GetDB().Create(&models.Receipt{
		Name:         "r1",
		Amount:       decimal.NewFromInt(100),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		GroupId:      groupId,
		PaidByUserID: userId,
		Status:       models.OPEN,
	})
}

// recordsReportBody is a minimal valid single-group records-mode CSV request.
const recordsReportBody = `{
  "name": "HTTP Report",
  "groupIds": ["1"],
  "period": {"preset": "custom", "startDate": "2026-05-01", "endDate": "2026-05-31"},
  "detail": {"mode": "records"},
  "columns": [
    {"kind": "dimension", "name": "Name", "label": "Name", "field": "name"},
    {"kind": "aggregate", "name": "Amount", "label": "Amount", "aggFunc": "SUM", "measure": "amount"}
  ],
  "formats": ["csv"]
}`

func TestGenerateReport_StreamsFileWhenAuthorized(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedReportReceipt(1, 1)

	w, r := generateReportRequest(1, recordsReportBody)
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusOK)
	if got := w.Header().Get("Content-Type"); got != constants.TextCsv {
		t.Errorf("Content-Type = %q, want %q", got, constants.TextCsv)
	}
	if got := w.Header().Get("Content-Disposition"); !strings.Contains(got, `filename="HTTP_Report.csv"`) {
		t.Errorf("Content-Disposition = %q, want it to name HTTP_Report.csv", got)
	}
	if w.Body.Len() == 0 {
		t.Error("expected a non-empty report body")
	}
}

func TestGenerateReport_ForbidsWhenMemberLacksPermission(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers() // user 1 is a member of group 1 with no granting role

	w, r := generateReportRequest(1, recordsReportBody)
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusForbidden)
}

func TestGenerateReport_ForbidsWhenNotAMemberOfARequestedGroup(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead) // reports.read in group 1 only

	// user 1 is not a member of group 2, so the report over it must be denied.
	body := strings.Replace(recordsReportBody, `"groupIds": ["1"]`, `"groupIds": ["2"]`, 1)
	w, r := generateReportRequest(1, body)
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusForbidden)
}

func TestGenerateReport_RejectsInvalidCommand(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	// No formats — the command validator rejects it before the permission gate.
	body := strings.Replace(recordsReportBody, `"formats": ["csv"]`, `"formats": []`, 1)
	w, r := generateReportRequest(1, body)
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestGenerateReport_MapsInvalidSpecToBadRequest(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedReportReceipt(1, 1)

	// The command validates, but grouping by a measure is an invalid spec the
	// engine rejects — the handler must map that ReportSpecError to a 400, not a 500.
	body := strings.Replace(recordsReportBody, `"detail": {"mode": "records"},`, `"groupBy": ["amount"],
  "detail": {"mode": "records"},`, 1)
	w, r := generateReportRequest(1, body)
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestGenerateReport_RejectsMalformedBody(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, "{not valid json")
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestGenerateReport_StreamsZipForMultipleFormats(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedReportReceipt(1, 1)

	body := strings.Replace(recordsReportBody, `"formats": ["csv"]`, `"formats": ["csv", "xlsx"]`, 1)
	w, r := generateReportRequest(1, body)
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusOK)
	if got := w.Header().Get("Content-Type"); got != constants.ApplicationZip {
		t.Errorf("Content-Type = %q, want %q", got, constants.ApplicationZip)
	}
	if got := w.Header().Get("Content-Disposition"); !strings.Contains(got, `filename="HTTP_Report.zip"`) {
		t.Errorf("Content-Disposition = %q, want it to name HTTP_Report.zip", got)
	}
}
