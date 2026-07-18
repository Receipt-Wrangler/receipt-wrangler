package services

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
)

// seedRenderableTemplate creates a report template whose stored configuration is a
// real ReportRequestCommand (so RenderTemplateForUser can unmarshal and render it)
// and indexes it under groupIds.
func seedRenderableTemplate(t *testing.T, name string, groupIds []uint, command commands.ReportRequestCommand) uint {
	t.Helper()
	db := repositories.GetDB()

	config, err := json.Marshal(command)
	if err != nil {
		t.Fatalf("marshal template config: %v", err)
	}
	template := models.ReportTemplate{Name: name, Configuration: config, ConfigurationVersion: 1}
	if err := db.Create(&template).Error; err != nil {
		t.Fatalf("seed template: %v", err)
	}
	for _, groupId := range groupIds {
		if err := db.Create(&models.ReportTemplateGroup{ReportTemplateID: template.ID, GroupID: groupId}).Error; err != nil {
			t.Fatalf("seed template group: %v", err)
		}
	}
	return template.ID
}

// A user who can read AND generate the template gets the rendered report HTML, the
// true receipt count, and an AllowedActions list containing both actions (so the
// widget shows its download button).
func TestReportService_RenderTemplateForUser_ReadAndGenerate(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	category := loadCategory(t, makeCategory(t, "Groceries"))
	userId := seedAppUser(t, "widget-rw", []string{permissions.AppReportsRead, permissions.AppReportsGenerate})
	groupId, _ := joinGroup(t, userId, "Household", []string{permissions.GroupReportsRead, permissions.GroupReceiptsRead})
	createReportReceipt(t, "h1", userId, groupId, []models.Category{category})
	createReportReceipt(t, "h2", userId, groupId, []models.Category{category})
	templateId := seedRenderableTemplate(t, "Widget", []uint{groupId}, aggregateReportCommand("Widget", []uint{groupId}, nil))

	preview, err := NewReportService(nil).RenderTemplateForUser(userId, utils.UintToString(templateId))
	if err != nil {
		t.Fatalf("RenderTemplateForUser: %v", err)
	}

	if preview.ReceiptCount != 2 {
		t.Errorf("receipt count = %d, want 2", preview.ReceiptCount)
	}
	if !strings.Contains(preview.Html, "<h1>Widget</h1>") || !strings.Contains(preview.Html, "Groceries") {
		t.Errorf("expected rendered report HTML with heading + category, got:\n%s", preview.Html)
	}
	if !slices.Contains(preview.AllowedActions, "read") || !slices.Contains(preview.AllowedActions, "generate") {
		t.Errorf("allowedActions = %v, want read + generate", preview.AllowedActions)
	}
}

// A user who can read but not generate still sees the rendered report, but
// AllowedActions omits "generate" so the widget hides its download button.
func TestReportService_RenderTemplateForUser_ReadOnlyOmitsGenerate(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	category := loadCategory(t, makeCategory(t, "Groceries"))
	userId := seedAppUser(t, "widget-ro", []string{permissions.AppReportsRead})
	groupId, _ := joinGroup(t, userId, "Household", []string{permissions.GroupReportsRead, permissions.GroupReceiptsRead})
	createReportReceipt(t, "h1", userId, groupId, []models.Category{category})
	templateId := seedRenderableTemplate(t, "Widget", []uint{groupId}, aggregateReportCommand("Widget", []uint{groupId}, nil))

	preview, err := NewReportService(nil).RenderTemplateForUser(userId, utils.UintToString(templateId))
	if err != nil {
		t.Fatalf("RenderTemplateForUser: %v", err)
	}

	if !strings.Contains(preview.Html, "<h1>Widget</h1>") {
		t.Errorf("expected rendered report HTML, got:\n%s", preview.Html)
	}
	if !slices.Contains(preview.AllowedActions, "read") {
		t.Errorf("allowedActions = %v, want read present", preview.AllowedActions)
	}
	if slices.Contains(preview.AllowedActions, "generate") {
		t.Errorf("allowedActions = %v, must omit generate (no generate permission)", preview.AllowedActions)
	}
}

// A user without view access to the template gets the restricted notice at a normal
// success (empty count + actions), not the report and not an error — the widget just
// renders whatever HTML it receives.
func TestReportService_RenderTemplateForUser_DeniedReadIsRestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "widget-denied", []string{permissions.AppAccountRead})
	groupId, _ := joinGroup(t, userId, "Household", []string{permissions.GroupReportsRead, permissions.GroupReceiptsRead})
	createReportReceipt(t, "h1", userId, groupId, nil)
	templateId := seedRenderableTemplate(t, "Widget", []uint{groupId}, aggregateReportCommand("Widget", []uint{groupId}, nil))

	preview, err := NewReportService(nil).RenderTemplateForUser(userId, utils.UintToString(templateId))
	if err != nil {
		t.Fatalf("RenderTemplateForUser: %v", err)
	}

	if preview.Html != restrictedReportHTML {
		t.Errorf("expected the restricted notice, got:\n%s", preview.Html)
	}
	if preview.ReceiptCount != 0 || len(preview.AllowedActions) != 0 {
		t.Errorf("restricted response must have count 0 and no actions, got count=%d actions=%v", preview.ReceiptCount, preview.AllowedActions)
	}
}

// A deleted (or never-existent) template resolves to the restricted notice, not an
// error, so a widget pinned to a removed template degrades gracefully.
func TestReportService_RenderTemplateForUser_MissingTemplateIsRestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "widget-missing", []string{permissions.AppReportsRead})

	preview, err := NewReportService(nil).RenderTemplateForUser(userId, "999999")
	if err != nil {
		t.Fatalf("RenderTemplateForUser: %v", err)
	}

	if preview.Html != restrictedReportHTML {
		t.Errorf("expected the restricted notice for a missing template, got:\n%s", preview.Html)
	}
}

// The widget renders the FULL dataset, not the row-capped builder-preview sample:
// with more than reportPreviewRowCap receipts, the grand total sums every receipt
// (would read short if this path reused the preview cap).
func TestReportService_RenderTemplateForUser_RendersFullDataset(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "widget-full", []string{permissions.AppReportsRead})
	groupId, _ := joinGroup(t, userId, "Household", []string{permissions.GroupReportsRead, permissions.GroupReceiptsRead})

	const receiptCount = reportPreviewRowCap + 1 // one past the preview cap
	receipts := make([]models.Receipt, receiptCount)
	for i := range receipts {
		receipts[i] = models.Receipt{
			Name:         "r" + strconv.Itoa(i),
			Amount:       decimal.NewFromInt(100),
			Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
			PaidByUserID: userId,
			GroupId:      groupId,
			Status:       models.OPEN,
		}
	}
	if err := repositories.GetDB().CreateInBatches(receipts, 200).Error; err != nil {
		t.Fatalf("seed receipts: %v", err)
	}
	templateId := seedRenderableTemplate(t, "Full", []uint{groupId}, aggregateReportCommand("Full", []uint{groupId}, nil))

	preview, err := NewReportService(nil).RenderTemplateForUser(userId, utils.UintToString(templateId))
	if err != nil {
		t.Fatalf("RenderTemplateForUser: %v", err)
	}

	if preview.ReceiptCount != receiptCount {
		t.Errorf("receipt count = %d, want %d", preview.ReceiptCount, receiptCount)
	}
	// 1001 receipts * 100 = 100,100 (the full total). A capped render would sum only
	// the first 1000 → 100,000, which does not contain "100,100".
	if !strings.Contains(preview.Html, "100,100") {
		t.Errorf("widget render must include the full dataset (grand total 100,100), got:\n%s", preview.Html)
	}
}
