package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

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

func TestGenerateReport_MalformedBodyIsBadRequest(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	// Invalid JSON is a client payload error, not a server failure.
	w, r := generateReportRequest(1, "{not valid json")
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

// errReader fails on read, standing in for a genuine request-body I/O failure
// (distinct from a malformed but readable body).
type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("body read failed") }

func TestGenerateReport_BodyReadErrorIsServerError(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/report/generate", errReader{})
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}}))
	GenerateReport(w, r)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestPreviewReport_ReturnsHtmlWhenAuthorized(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedReportReceipt(1, 1)

	w, r := generateReportRequest(1, recordsReportBody)
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusOK)
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, constants.ApplicationJson) {
		t.Errorf("Content-Type = %q, want %q", got, constants.ApplicationJson)
	}

	var preview services.ReportPreview
	if err := json.Unmarshal(w.Body.Bytes(), &preview); err != nil {
		t.Fatalf("decode preview body: %v", err)
	}
	if !strings.Contains(preview.Html, "<") {
		t.Errorf("expected HTML in the preview body, got %q", preview.Html)
	}
	if preview.ReceiptCount != 1 {
		t.Errorf("receiptCount = %d, want 1", preview.ReceiptCount)
	}
}

func TestPreviewReport_ForbidsWhenMemberLacksPermission(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers() // member of group 1 with no granting role

	w, r := generateReportRequest(1, recordsReportBody)
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusForbidden)
}

func TestPreviewReport_MapsInvalidSpecToBadRequest(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedReportReceipt(1, 1)

	body := strings.Replace(recordsReportBody, `"detail": {"mode": "records"},`, `"groupBy": ["amount"],
  "detail": {"mode": "records"},`, 1)
	w, r := generateReportRequest(1, body)
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestPreviewReport_MalformedBodyIsBadRequest(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, "{not valid json")
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
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

	// The streamed body must be a valid zip carrying both non-empty entries.
	reader, err := zip.NewReader(bytes.NewReader(w.Body.Bytes()), int64(w.Body.Len()))
	if err != nil {
		t.Fatalf("open zip body: %v", err)
	}
	sizes := make(map[string]uint64, len(reader.File))
	for _, file := range reader.File {
		sizes[file.Name] = file.UncompressedSize64
	}
	for _, name := range []string{"HTTP_Report.csv", "HTTP_Report.xlsx"} {
		if size, ok := sizes[name]; !ok || size == 0 {
			t.Errorf("zip entry %q missing or empty (entries: %v)", name, sizes)
		}
	}
}

func TestCreateReportTemplate_SavesWhenAuthorized(t *testing.T) {
	defer tearDownReportTest()
	grantAppPerms(t, 1, permissions.AppReportsCreate)

	w, r := generateReportRequest(1, recordsReportBody)
	CreateReportTemplate(w, r)

	assertStatus(t, w, http.StatusOK)

	var template models.ReportTemplate
	if err := json.Unmarshal(w.Body.Bytes(), &template); err != nil {
		t.Fatalf("decode template body: %v", err)
	}
	if template.ID == 0 {
		t.Error("expected the saved template to carry an id")
	}
	if template.Name != "HTTP Report" {
		t.Errorf("template name = %q, want %q", template.Name, "HTTP Report")
	}
	if template.CreatedBy == nil || *template.CreatedBy != 1 {
		t.Errorf("template createdBy = %v, want 1", template.CreatedBy)
	}
	if template.ConfigurationVersion != 1 {
		t.Errorf("template configurationVersion = %d, want 1", template.ConfigurationVersion)
	}
}

func TestCreateReportTemplate_ForbidsWithoutCreatePermission(t *testing.T) {
	defer tearDownReportTest()
	// The caller can access reports but cannot create templates — read must not
	// imply create.
	grantAppPerms(t, 1, permissions.AppReportsRead)

	w, r := generateReportRequest(1, recordsReportBody)
	CreateReportTemplate(w, r)

	assertStatus(t, w, http.StatusForbidden)
}

func TestCreateReportTemplate_RejectsInvalidCommand(t *testing.T) {
	defer tearDownReportTest()
	grantAppPerms(t, 1, permissions.AppReportsCreate)

	// No formats — the shared loadReportCommand validator rejects it, so a template
	// can only ever store a complete, buildable configuration.
	body := strings.Replace(recordsReportBody, `"formats": ["csv"]`, `"formats": []`, 1)
	w, r := generateReportRequest(1, body)
	CreateReportTemplate(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

// deleteReportTemplateRequest builds a DELETE carrying JWT claims for userId and a
// chi route context supplying the {id} URL param, mirroring how the router invokes
// DeleteReportTemplate.
func deleteReportTemplateRequest(userId uint, id string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/report/template/"+id, nil)
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}}))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chi.RouteContext(r.Context()).URLParams.Add("id", id)
	return w, r
}

func TestDeleteReportTemplate_DeletesWhenAuthorized(t *testing.T) {
	defer tearDownReportTest()
	grantAppPerms(t, 1, permissions.AppReportsCreate, permissions.AppReportsDelete)

	// Save a template through the create handler, then delete it.
	cw, cr := generateReportRequest(1, recordsReportBody)
	CreateReportTemplate(cw, cr)
	assertStatus(t, cw, http.StatusOK)

	var created models.ReportTemplate
	if err := json.Unmarshal(cw.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created template: %v", err)
	}

	dw, dr := deleteReportTemplateRequest(1, fmt.Sprint(created.ID))
	DeleteReportTemplate(dw, dr)
	assertStatus(t, dw, http.StatusOK)

	// The row is gone.
	var count int64
	repositories.GetDB().Model(&models.ReportTemplate{}).Where("id = ?", created.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected the template to be deleted, but %d row(s) remain", count)
	}
}

func TestDeleteReportTemplate_ForbidsWithoutDeletePermission(t *testing.T) {
	defer tearDownReportTest()
	// Read access does not imply delete.
	grantAppPerms(t, 1, permissions.AppReportsRead)

	w, r := deleteReportTemplateRequest(1, "1")
	DeleteReportTemplate(w, r)

	assertStatus(t, w, http.StatusForbidden)
}
