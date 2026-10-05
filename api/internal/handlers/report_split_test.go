package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
)

// splitReportBody aggregates group 1's receipts by category, splitting each
// receipt equally across its categories.
const splitReportBody = `{
  "name": "Split Report",
  "groupIds": ["1"],
  "period": {"preset": "custom", "startDate": "2026-05-01", "endDate": "2026-05-31"},
  "detail": {"mode": "aggregate", "by": "category"},
  "columns": [
    {"kind": "dimension", "name": "Category", "label": "Category", "field": "category"},
    {"kind": "aggregate", "name": "Total", "label": "Total", "aggFunc": "SUM", "measure": "amount"}
  ],
  "grandTotals": true,
  "splitCategoriesEqually": true,
  "splitExcludedFields": ["custom_9"],
  "formats": ["csv"]
}`

// seedSplitReceipt stores a 100.00 receipt in group 1 carrying two categories,
// so a split report shows 50.00 per category and an unsplit one 100.00.
func seedSplitReceipt(t *testing.T) {
	t.Helper()
	db := repositories.GetDB()
	categories := []models.Category{{Name: "Food"}, {Name: "Fuel"}}
	for index := range categories {
		if err := db.Create(&categories[index]).Error; err != nil {
			t.Fatalf("seed category: %v", err)
		}
	}
	receipt := models.Receipt{
		Name:         "shared",
		Amount:       decimal.NewFromInt(100),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		GroupId:      1,
		PaidByUserID: 1,
		Status:       models.OPEN,
		Categories:   categories,
	}
	if err := db.Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
}

func TestPreviewReport_SplitsAcrossCategories(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedSplitReceipt(t)

	w, r := generateReportRequest(1, splitReportBody)
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusOK)
	preview := decodePreview(t, w.Body.Bytes())
	if preview.ReceiptCount != 1 {
		t.Errorf("receiptCount = %d, want 1 receipt, not one per category", preview.ReceiptCount)
	}
	if !strings.Contains(preview.Html, "50.00") {
		t.Errorf("expected 50.00 per category, got:\n%s", preview.Html)
	}
	if strings.Contains(preview.Html, "200.00") {
		t.Errorf("grand total still double counts (200.00):\n%s", preview.Html)
	}

	// The same body without the split attributes the receipt to both in full.
	w, r = generateReportRequest(1, strings.Replace(splitReportBody, `"splitCategoriesEqually": true`, `"splitCategoriesEqually": false`, 1))
	PreviewReport(w, r)
	assertStatus(t, w, http.StatusOK)
	if html := decodePreview(t, w.Body.Bytes()).Html; !strings.Contains(html, "200.00") {
		t.Errorf("unsplit grand total should be 200.00, got:\n%s", html)
	}
}

func TestPreviewReport_RejectsMalformedSplitExcludedField(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	body := strings.Replace(splitReportBody, `["custom_9"]`, `["amount"]`, 1)
	w, r := generateReportRequest(1, body)
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
	if !strings.Contains(w.Body.String(), "splitExcludedFields") {
		t.Errorf("expected the error under splitExcludedFields, got %s", w.Body.String())
	}
}

// A template stores the split options and hands them back unchanged.
func TestReportTemplate_SplitOptionsRoundTrip(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsCreate, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, splitReportBody)
	CreateReportTemplate(w, r)
	assertStatus(t, w, http.StatusOK)

	var created models.ReportTemplate
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created template: %v", err)
	}

	w, r = reportTemplateIdRequest("GET", 1, fmt.Sprint(created.ID))
	GetReportTemplate(w, r)
	assertStatus(t, w, http.StatusOK)

	var fetched models.ReportTemplate
	if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decode fetched template: %v", err)
	}
	var configuration commands.ReportRequestCommand
	if err := json.Unmarshal(fetched.Configuration, &configuration); err != nil {
		t.Fatalf("decode configuration: %v", err)
	}
	if !configuration.SplitCategoriesEqually || configuration.SplitTagsEqually ||
		len(configuration.SplitExcludedFields) != 1 || configuration.SplitExcludedFields[0] != "custom_9" {
		t.Errorf("split options did not round-trip: %+v", configuration)
	}
}

// A template that never used the split stores none of its keys.
func TestCreateReportTemplate_StoredConfigOmitsUnsetSplitOptions(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsCreate)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, recordsReportBody)
	CreateReportTemplate(w, r)
	assertStatus(t, w, http.StatusOK)

	var created models.ReportTemplate
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created template: %v", err)
	}
	var stored models.ReportTemplate
	if err := repositories.GetDB().First(&stored, created.ID).Error; err != nil {
		t.Fatalf("load stored template: %v", err)
	}
	for _, key := range []string{"splitCategoriesEqually", "splitTagsEqually", "splitExcludedFields"} {
		if strings.Contains(string(stored.Configuration), key) {
			t.Errorf("stored configuration carries unset %s: %s", key, stored.Configuration)
		}
	}
}

func seedSplitTemplate(t *testing.T) models.ReportTemplate {
	t.Helper()
	var command commands.ReportRequestCommand
	if err := json.Unmarshal([]byte(splitReportBody), &command); err != nil {
		t.Fatalf("decode split body: %v", err)
	}
	template, err := repositories.NewReportTemplateRepository(nil).CreateReportTemplate(command, 1)
	if err != nil {
		t.Fatalf("seed report template: %v", err)
	}
	return template
}

func TestGenerateReportFromTemplate_AppliesSplit(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsGenerate)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedSplitReceipt(t)
	template := seedSplitTemplate(t)

	w, r := reportTemplateIdRequest("POST", 1, fmt.Sprint(template.ID))
	GenerateReportFromTemplate(w, r)

	assertStatus(t, w, http.StatusOK)
	content := w.Body.String()
	if !strings.Contains(content, "50.00") || strings.Contains(content, "200.00") {
		t.Errorf("expected split 50.00 buckets and a 100.00 grand total, got:\n%s", content)
	}
}

func TestRenderReportTemplate_AppliesSplit(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
	seedSplitReceipt(t)
	template := seedSplitTemplate(t)

	w, r := reportTemplateIdRequest("POST", 1, fmt.Sprint(template.ID))
	RenderReportTemplate(w, r)

	assertStatus(t, w, http.StatusOK)
	preview := decodePreview(t, w.Body.Bytes())
	if !strings.Contains(preview.Html, "50.00") || strings.Contains(preview.Html, "200.00") {
		t.Errorf("expected split 50.00 buckets and no 200.00 double count, got:\n%s", preview.Html)
	}
}
