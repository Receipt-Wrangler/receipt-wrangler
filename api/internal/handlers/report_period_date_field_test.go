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
	"receipt-wrangler/api/internal/services"
)

// seedDatedReportReceipt creates a receipt in group 1 dated in May 2026, resolved
// in July and added in August, so each date field lands in its own month and a
// period picks it up on exactly one of them.
func seedDatedReportReceipt(t *testing.T, name string) {
	t.Helper()
	resolved := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(100),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		ResolvedDate: &resolved,
		GroupId:      1,
		PaidByUserID: 1,
		Status:       models.RESOLVED,
	}
	db := repositories.GetDB()
	if err := db.Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	added := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	if err := db.Model(&receipt).UpdateColumn("created_at", added).Error; err != nil {
		t.Fatalf("stamp created_at: %v", err)
	}
}

// periodReportBody is recordsReportBody over one calendar month of 2026, with the
// period's date field set when dateField is non-empty.
func periodReportBody(month int, dateField string) string {
	lastDay := time.Date(2026, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	period := fmt.Sprintf(`"period": {"preset": "custom", "startDate": "2026-%02d-01", "endDate": "2026-%02d-%02d"`, month, month, lastDay)
	if dateField != "" {
		period += fmt.Sprintf(`, "dateField": %q`, dateField)
	}
	period += "}"
	return strings.Replace(recordsReportBody,
		`"period": {"preset": "custom", "startDate": "2026-05-01", "endDate": "2026-05-31"}`, period, 1)
}

// seedPeriodTemplate stores a template over group 1 whose custom period covers
// one calendar month of 2026 on dateField (empty for a template saved before the
// field existed). It goes through the repository, not the handler, so it can also
// store a value the handler's validation would reject.
func seedPeriodTemplate(t *testing.T, dateField string, month int) models.ReportTemplate {
	t.Helper()
	lastDay := time.Date(2026, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	command := commands.ReportRequestCommand{
		Name:     "Period Template",
		GroupIds: []string{"1"},
		Period: commands.ReportPeriod{
			Preset:    commands.ReportPeriodCustom,
			StartDate: fmt.Sprintf("2026-%02d-01", month),
			EndDate:   fmt.Sprintf("2026-%02d-%02d", month, lastDay),
			DateField: dateField,
		},
		Detail:  commands.ReportDetail{Mode: commands.ReportDetailRecords},
		Columns: []commands.ReportColumn{{Kind: "dimension", Name: "Name", Label: "Name", Field: "name"}},
		Formats: []string{commands.ReportFormatCsv},
	}
	template, err := repositories.NewReportTemplateRepository(nil).CreateReportTemplate(command, 1)
	if err != nil {
		t.Fatalf("seed report template: %v", err)
	}
	return template
}

func decodePreview(t *testing.T, body []byte) services.ReportPreview {
	t.Helper()
	var preview services.ReportPreview
	if err := json.Unmarshal(body, &preview); err != nil {
		t.Fatalf("decode preview body: %v", err)
	}
	return preview
}

// Over HTTP, a period covers the receipt only in the month its chosen date field
// falls in: May on the receipt date (also the default), July on the resolved date,
// August on the date it was added.
func TestPreviewReport_HonorsPeriodDateField(t *testing.T) {
	tests := []struct {
		dateField string
		month     int
		want      int
	}{
		{"", 5, 1},
		{commands.ReceiptFilterKeyDate, 5, 1},
		{commands.ReceiptFilterKeyDate, 8, 0},
		{commands.ReceiptFilterKeyResolvedDate, 7, 1},
		{commands.ReceiptFilterKeyResolvedDate, 5, 0},
		{commands.ReceiptFilterKeyCreatedAt, 8, 1},
		{commands.ReceiptFilterKeyCreatedAt, 5, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("dateField=%s/month=%d", test.dateField, test.month), func(t *testing.T) {
			defer tearDownReportTest()
			repositories.CreateTestGroupWithUsers()
			grantAppPerms(t, 1, permissions.AppReportsRead)
			grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
			seedDatedReportReceipt(t, "dated")

			w, r := generateReportRequest(1, periodReportBody(test.month, test.dateField))
			PreviewReport(w, r)

			assertStatus(t, w, http.StatusOK)
			if got := decodePreview(t, w.Body.Bytes()).ReceiptCount; got != test.want {
				t.Errorf("receiptCount = %d, want %d", got, test.want)
			}
		})
	}
}

func TestPreviewReport_RejectsUnknownPeriodDateField(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, periodReportBody(5, "created_at"))
	PreviewReport(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}

// A template can only store a date field the report can run on.
func TestCreateReportTemplate_RejectsUnknownPeriodDateField(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsCreate)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, periodReportBody(5, "paidAt"))
	CreateReportTemplate(w, r)

	assertStatus(t, w, http.StatusBadRequest)
	var count int64
	repositories.GetDB().Model(&models.ReportTemplate{}).Count(&count)
	if count != 0 {
		t.Errorf("report_templates rows = %d, want 0", count)
	}
}

// A saved date field comes back on the template, and a template saved without one
// carries no dateField key at all — the shape every client already reads as the
// receipt date.
func TestGetReportTemplate_ReturnsPeriodDateField(t *testing.T) {
	tests := []struct {
		dateField string
		wantKey   bool
	}{
		{commands.ReceiptFilterKeyCreatedAt, true},
		{"", false},
	}
	for _, test := range tests {
		t.Run("dateField="+test.dateField, func(t *testing.T) {
			defer tearDownReportTest()
			repositories.CreateTestGroupWithUsers()
			grantAppPerms(t, 1, permissions.AppReportsRead)
			grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
			seeded := seedPeriodTemplate(t, test.dateField, 8)

			w, r := reportTemplateIdRequest("GET", 1, fmt.Sprint(seeded.ID))
			GetReportTemplate(w, r)

			assertStatus(t, w, http.StatusOK)
			var body struct {
				Configuration struct {
					Period map[string]interface{} `json:"period"`
				} `json:"configuration"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode template body: %v", err)
			}
			got, hasKey := body.Configuration.Period["dateField"]
			if hasKey != test.wantKey {
				t.Fatalf("period = %v, want dateField present = %v", body.Configuration.Period, test.wantKey)
			}
			if test.wantKey && got != test.dateField {
				t.Errorf("dateField = %v, want %q", got, test.dateField)
			}
		})
	}
}

// The dashboard widget renders a stored configuration without re-validating it.
// A stored date field is honoured there; a template saved before the field existed
// stays on the receipt date, and so does a stored value the report cannot run on.
// The receipt's three dates fall in three different months, so each case can only
// match through the date it names: May is the receipt date, July the resolved
// date, August the date it was added.
func TestRenderReportTemplate_HonorsPeriodDateField(t *testing.T) {
	tests := []struct {
		dateField string
		month     int
		want      int
	}{
		{commands.ReceiptFilterKeyCreatedAt, 8, 1},
		{"", 5, 1},
		{"", 8, 0},
		{"paidAt", 5, 1},
		{"paidAt", 8, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("dateField=%s/month=%d", test.dateField, test.month), func(t *testing.T) {
			defer tearDownReportTest()
			repositories.CreateTestGroupWithUsers()
			grantAppPerms(t, 1, permissions.AppReportsRead)
			grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
			seedDatedReportReceipt(t, "dated")
			seeded := seedPeriodTemplate(t, test.dateField, test.month)

			w, r := reportTemplateIdRequest("POST", 1, fmt.Sprint(seeded.ID))
			RenderReportTemplate(w, r)

			assertStatus(t, w, http.StatusOK)
			if got := decodePreview(t, w.Body.Bytes()).ReceiptCount; got != test.want {
				t.Errorf("receiptCount = %d, want %d", got, test.want)
			}
		})
	}
}

func TestGenerateReportFromTemplate_HonorsPeriodDateField(t *testing.T) {
	tests := []struct {
		dateField    string
		month        int
		wantIncluded bool
	}{
		{commands.ReceiptFilterKeyCreatedAt, 8, true},
		{"", 5, true},
		{"", 8, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("dateField=%s/month=%d", test.dateField, test.month), func(t *testing.T) {
			defer tearDownReportTest()
			repositories.CreateTestGroupWithUsers()
			grantAppPerms(t, 1, permissions.AppReportsGenerate)
			grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
			seedDatedReportReceipt(t, "added-in-august")
			seeded := seedPeriodTemplate(t, test.dateField, test.month)

			w, r := reportTemplateIdRequest("POST", 1, fmt.Sprint(seeded.ID))
			GenerateReportFromTemplate(w, r)

			assertStatus(t, w, http.StatusOK)
			if got := strings.Contains(w.Body.String(), "added-in-august"); got != test.wantIncluded {
				t.Errorf("csv includes the receipt = %v, want %v:\n%s", got, test.wantIncluded, w.Body.String())
			}
		})
	}
}

// --- the drill-in list -----------------------------------------------------

func decodeReportReceipts(t *testing.T, body []byte) ([]models.Receipt, int64) {
	t.Helper()
	var paged struct {
		Data       []models.Receipt `json:"data"`
		TotalCount int64            `json:"totalCount"`
	}
	if err := json.Unmarshal(body, &paged); err != nil {
		t.Fatalf("decode report receipts body: %v", err)
	}
	return paged.Data, paged.TotalCount
}

// The drill-in lists the receipts the preview counts, through the same date field.
func TestGetReportReceipts_ListsTheReceiptsOnTheChosenDateField(t *testing.T) {
	tests := []struct {
		dateField string
		month     int
		want      int
	}{
		{"", 5, 1},
		{commands.ReceiptFilterKeyDate, 5, 1},
		{commands.ReceiptFilterKeyDate, 8, 0},
		{commands.ReceiptFilterKeyResolvedDate, 7, 1},
		{commands.ReceiptFilterKeyResolvedDate, 5, 0},
		{commands.ReceiptFilterKeyCreatedAt, 8, 1},
		{commands.ReceiptFilterKeyCreatedAt, 5, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("dateField=%s/month=%d", test.dateField, test.month), func(t *testing.T) {
			defer tearDownReportTest()
			repositories.CreateTestGroupWithUsers()
			grantAppPerms(t, 1, permissions.AppReportsRead)
			grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)
			seedDatedReportReceipt(t, "dated")

			w, r := generateReportRequest(1, periodReportBody(test.month, test.dateField))
			GetReportReceipts(w, r)

			assertStatus(t, w, http.StatusOK)
			receipts, total := decodeReportReceipts(t, w.Body.Bytes())
			if len(receipts) != test.want || total != int64(test.want) {
				t.Fatalf("listed %d receipts (total %d), want %d", len(receipts), total, test.want)
			}
			if test.want == 1 && receipts[0].Name != "dated" {
				t.Errorf("listed %q, want the seeded receipt", receipts[0].Name)
			}
		})
	}
}

// Gated exactly like the preview whose count it lists: the app-level report
// permission, and report access in every covered group.
func TestGetReportReceipts_ForbidsWithoutAppReportsPermission(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, recordsReportBody)
	GetReportReceipts(w, r)

	assertStatus(t, w, http.StatusForbidden)
}

func TestGetReportReceipts_ForbidsWithoutGroupReportsPermission(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	w, r := generateReportRequest(1, recordsReportBody)
	GetReportReceipts(w, r)

	assertStatus(t, w, http.StatusForbidden)
}

func TestGetReportReceipts_RejectsInvalidCommand(t *testing.T) {
	defer tearDownReportTest()
	repositories.CreateTestGroupWithUsers()
	grantAppPerms(t, 1, permissions.AppReportsRead)
	grantGroupPerms(t, 1, 1, permissions.GroupReportsRead)

	w, r := generateReportRequest(1, periodReportBody(5, "created_at"))
	GetReportReceipts(w, r)

	assertStatus(t, w, http.StatusBadRequest)
}
