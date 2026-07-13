package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/reporting/receiptsource"
	"receipt-wrangler/api/internal/repositories"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func loadCategory(t *testing.T, id uint) models.Category {
	t.Helper()
	var category models.Category
	if err := repositories.GetDB().First(&category, id).Error; err != nil {
		t.Fatalf("load category %d: %v", id, err)
	}
	return category
}

// createReportReceipt inserts a receipt in a group. Categories must be loaded
// rows (not bare ids), so the association write does not blank their names.
func createReportReceipt(t *testing.T, name string, paidByUserId uint, groupId uint, categories []models.Category) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(100),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		PaidByUserID: paidByUserId,
		GroupId:      groupId,
		Status:       models.OPEN,
		Categories:   categories,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("create receipt %q: %v", name, err)
	}
	return receipt
}

func groupIdString(groupId uint) string {
	return strconv.FormatUint(uint64(groupId), 10)
}

// Paid-by visibility hides a whole receipt: one paid by an allowed payer is
// returned, one paid by a disallowed payer never reaches the engine.
func TestReportDataService_HidesReceiptsByPaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	allowedPayer := makeUser(t, "rpt-allowed-payer")
	hiddenPayer := makeUser(t, "rpt-hidden-payer")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "rpt-reviewer", []uint{allowedPayer}, false)

	visible := createReportReceipt(t, "visible", allowedPayer, groupId, nil)
	createReportReceipt(t, "hidden", hiddenPayer, groupId, nil)

	_, rows, err := NewReportDataService(nil).Rows(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{})
	if err != nil {
		t.Fatalf("Rows: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected only the allowed-payer receipt, got %d rows", len(rows))
	}
	id, _ := rows[0].Measure(receiptsource.KeyReceiptID).Decimal()
	if id.IntPart() != int64(visible.ID) {
		t.Errorf("row is receipt %s, want %d", id, visible.ID)
	}
}

// A category the caller may not see becomes (Restricted) in the row rather than
// being stripped (which would silently drop it from the totals).
func TestReportDataService_SubstitutesRestrictedCategory(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	allowedCatId := makeCategory(t, "Groceries")
	hiddenCatId := makeCategory(t, "Salary")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "rpt-cat", []uint{allowedCatId}, nil)

	createReportReceipt(t, "mixed", userId, groupId, []models.Category{loadCategory(t, allowedCatId), loadCategory(t, hiddenCatId)})

	_, rows, err := NewReportDataService(nil).Rows(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{})
	if err != nil {
		t.Fatalf("Rows: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	got := make([]string, 0)
	for _, value := range rows[0].Get(receiptsource.KeyCategory) {
		text, _ := value.Text()
		got = append(got, text)
	}
	if !slices.Contains(got, "Groceries") || !slices.Contains(got, "(Restricted)") {
		t.Errorf("category values = %v, want Groceries and (Restricted)", got)
	}
	if slices.Contains(got, "Salary") {
		t.Errorf("hidden category leaked into the row: %v", got)
	}
}

// The catalog carries the built-ins, the derived date fields, and one field per
// custom field; and the rows resolve a receipt's custom field value.
func TestReportDataService_ResolvesCustomFieldsInCatalogAndRows(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "rpt-custom", nil, nil)

	customField := models.CustomField{Name: "HST", Type: models.CURRENCY}
	if err := repositories.GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}

	hst := decimal.RequireFromString("15.60")
	receipt := models.Receipt{
		Name:         "with-hst",
		Amount:       decimal.NewFromInt(100),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		PaidByUserID: userId,
		GroupId:      groupId,
		Status:       models.OPEN,
		CustomFields: []models.CustomFieldValue{{CustomFieldId: customField.ID, CurrencyValue: &hst}},
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("create receipt: %v", err)
	}

	customKey := receiptsource.CustomFieldKey(customField.ID)
	catalog, rows, err := NewReportDataService(nil).Rows(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{})
	if err != nil {
		t.Fatalf("Rows: %v", err)
	}

	for _, key := range []reporting.FieldKey{receiptsource.KeyAmount, receiptsource.KeyDateMonth, customKey} {
		if _, ok := catalog.Get(key); !ok {
			t.Errorf("catalog missing %s", key)
		}
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	value, isNumber := rows[0].Measure(customKey).Decimal()
	if !isNumber || !value.Equal(hst) {
		t.Errorf("custom field value = %v, want %v", rows[0].Measure(customKey), hst)
	}
}
