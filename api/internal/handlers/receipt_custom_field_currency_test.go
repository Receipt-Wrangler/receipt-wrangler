package handlers

import (
	"fmt"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// The desktop's masked currency input sends "" for a CURRENCY custom field the
// user emptied. That used to fail the whole save with a decimal decode error; it
// must now save, and store the field as blank (nil), like an untouched one.
func TestUpdateReceipt_BlankCurrencyCustomFieldSavesAsNoValue(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	db := repositories.GetDB()

	field := models.CustomField{Name: "Tax", Type: models.CURRENCY}
	if err := db.Create(&field).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}

	stored := decimal.RequireFromString("12.34")
	receipt := models.Receipt{
		Name:         "Existing",
		Amount:       decimal.NewFromInt(5),
		Date:         time.Now(),
		GroupId:      groupId,
		PaidByUserID: userId,
		Status:       models.OPEN,
		CustomFields: []models.CustomFieldValue{{CustomFieldId: field.ID, CurrencyValue: &stored}},
	}
	if err := db.Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}

	body := fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2024-01-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN","customFields":[{"customFieldId":%d,"currencyValue":""}]}`,
		groupId, userId, field.ID,
	)
	w, r := updateReceiptRequest(receipt.ID, userId, body)

	UpdateReceipt(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		t.FailNow()
	}

	var values []models.CustomFieldValue
	if err := db.Where("receipt_id = ?", receipt.ID).Find(&values).Error; err != nil {
		t.Fatalf("load custom field values: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("receipt has %d custom field values, want 1", len(values))
	}
	if values[0].CustomFieldId != field.ID {
		t.Errorf("CustomFieldId = %d, want %d", values[0].CustomFieldId, field.ID)
	}
	if values[0].CurrencyValue != nil {
		t.Errorf("CurrencyValue = %s, want nil", values[0].CurrencyValue)
	}
}
