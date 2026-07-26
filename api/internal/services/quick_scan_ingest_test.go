package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

// These tests exercise the full "AI response -> created receipt" path: a controlled AI JSON body is
// fed through the real ReceiptService.QuickScan (the only place an AI response reaches CreateReceipt)
// and the persisted receipt is read back through the deep-preload GetFullyLoadedReceiptById to prove
// that every field survives - name, amount, date, paid-by, status, categories, tags, receipt items
// (amount / status / chargedToUserId / isTaxed / per-item categories+tags / linked items), comments,
// and custom-field values of all five value types. If a future prompt starts emitting these fields,
// these tests prove the plumbing already persists them; a regression that silently drops a field
// fails here.
//
// The AI is mocked at the HTTP transport layer (see ai_test.go) via a server whose body is set AFTER
// the pipeline is seeded, so fixtures can reference the seeded user's real id without guessing.

// newMutableOllamaServer serves whatever *body points at, at request time. Unlike
// newMockOllamaServerForService (fixed body) this lets a test seed the graph first, learn the seeded
// user's id, then set the AI response referencing that id - the request only fires later, inside
// QuickScan.
func newMutableOllamaServer(t *testing.T, body *string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(*body))
	}))
	t.Cleanup(server.Close)
	return server
}

// ollamaReceiptResponse wraps a receipt JSON string as the message.content of an Ollama chat
// response, JSON-escaping it so quotes/newlines in the receipt JSON can't break the envelope.
func ollamaReceiptResponse(t *testing.T, receiptJSON string) string {
	t.Helper()
	envelope := map[string]any{
		"model":      "test",
		"created_at": "2024-01-01T00:00:00Z",
		"message":    map[string]any{"role": "assistant", "content": receiptJSON},
		"done":       true,
	}
	bytes, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshal ollama envelope: %v", err)
	}
	return string(bytes)
}

// writeQuickScanTempFile writes the shared JPG fixture to a temp file and returns its path, mirroring
// what handlers.QuickScan hands ReceiptService.QuickScan.
func writeQuickScanTempFile(t *testing.T) string {
	t.Helper()
	jpg, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		t.Fatalf("read jpg fixture: %v", err)
	}
	path, err := repositories.NewFileRepository(nil).WriteTempFile(jpg)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })
	return path
}

// seededCustomFields carries the ids of the custom fields (one per value type) seeded for a test so
// the AI-response fixture can reference them.
type seededCustomFields struct {
	textId, dateId, selectId, optionId, currencyId, booleanId uint
}

func seedQuickScanCustomFields(t *testing.T) seededCustomFields {
	t.Helper()
	db := repositories.GetDB()

	text := models.CustomField{Name: "Text Field", Type: models.TEXT}
	date := models.CustomField{Name: "Date Field", Type: models.DATE}
	selectField := models.CustomField{Name: "Select Field", Type: models.SELECT}
	currency := models.CustomField{Name: "Currency Field", Type: models.CURRENCY}
	boolean := models.CustomField{Name: "Boolean Field", Type: models.BOOLEAN}
	for _, field := range []*models.CustomField{&text, &date, &selectField, &currency, &boolean} {
		if err := db.Create(field).Error; err != nil {
			t.Fatalf("seed custom field: %v", err)
		}
	}

	option := models.CustomFieldOption{Value: "Option A", CustomFieldId: selectField.ID}
	if err := db.Create(&option).Error; err != nil {
		t.Fatalf("seed custom field option: %v", err)
	}

	return seededCustomFields{
		textId:     text.ID,
		dateId:     date.ID,
		selectId:   selectField.ID,
		optionId:   option.ID,
		currencyId: currency.ID,
		booleanId:  boolean.ID,
	}
}

// runQuickScan seeds the receipt-image pipeline + categories (ids 1-3) + tags (ids 1-2), points a
// mock AI at the JSON built from the seeded user/group, and drives ReceiptService.QuickScan. The
// aiJSONFor callback runs after seeding so it can reference the seeded user's real id; paidByArg is
// the QuickScan paid-by argument (the handler-resolved default), resolved from the seeded user too.
func runQuickScan(
	t *testing.T,
	paidByArg func(user models.User) uint,
	status models.ReceiptStatus,
	categoryPicks []uint,
	tagPicks []uint,
	aiJSONFor func(user models.User, group models.Group) string,
) (models.Receipt, models.User, models.Group, error) {
	return runQuickScanWithSetup(t, paidByArg, status, categoryPicks, tagPicks, aiJSONFor, nil)
}

// runQuickScanWithSetup is runQuickScan with an optional postSeed hook that runs after the pipeline +
// categories/tags are seeded and before the AI body is set - so a test can, e.g., assign the seeded
// member a grant-restricted group role to exercise the out-of-grant drop.
func runQuickScanWithSetup(
	t *testing.T,
	paidByArg func(user models.User) uint,
	status models.ReceiptStatus,
	categoryPicks []uint,
	tagPicks []uint,
	aiJSONFor func(user models.User, group models.Group) string,
	postSeed func(t *testing.T, user models.User, group models.Group),
) (models.Receipt, models.User, models.Group, error) {
	t.Helper()

	var body string
	server := newMutableOllamaServer(t, &body)
	user, group, _ := seedReceiptImagePipeline(t, server.URL)
	repositories.CreateTestCategories()
	createTestTags()

	if postSeed != nil {
		postSeed(t, user, group)
	}

	body = ollamaReceiptResponse(t, aiJSONFor(user, group))
	tempPath := writeQuickScanTempFile(t)

	paidBy := uint(0)
	if paidByArg != nil {
		paidBy = paidByArg(user)
	}

	receipt, err := NewReceiptService(nil).QuickScan(
		&structs.Claims{UserId: user.ID},
		paidBy,
		group.ID,
		status,
		categoryPicks,
		tagPicks,
		tempPath,
		"test.jpg",
		"test-task",
	)
	return receipt, user, group, err
}

func hasCategoryId(categories []models.Category, id uint) bool {
	for _, category := range categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

func hasTagId(tags []models.Tag, id uint) bool {
	for _, tag := range tags {
		if tag.ID == id {
			return true
		}
	}
	return false
}

func findItemByName(items []models.Item, name string) (models.Item, bool) {
	for _, item := range items {
		if item.Name == name {
			return item, true
		}
	}
	return models.Item{}, false
}

func findCustomFieldValue(values []models.CustomFieldValue, customFieldId uint) (models.CustomFieldValue, bool) {
	for _, value := range values {
		if value.CustomFieldId == customFieldId {
			return value, true
		}
	}
	return models.CustomFieldValue{}, false
}

// fullReceiptAiJSON is a maximal receipt the AI could return: every persistable field populated,
// including nested linked items, item-level categories/tags, comments, and one custom-field value per
// type. The %d placeholders are: paidByUserId, item.chargedToUserId, comment.userId, then the five
// custom-field ids (with the select option id between the select field id and the currency field id).
const fullReceiptAiJSON = `{
	"name": "Full Receipt",
	"amount": 100.00,
	"date": "2024-05-01T00:00:00Z",
	"paidByUserId": %d,
	"status": "OPEN",
	"categories": [{"id": 1, "name": "test"}, {"id": 2, "name": "test2"}],
	"tags": [{"id": 1, "name": "tag-a"}],
	"receiptItems": [
		{
			"name": "Widget",
			"amount": 40.00,
			"status": "OPEN",
			"chargedToUserId": %d,
			"isTaxed": true,
			"categories": [{"id": 1, "name": "test"}],
			"tags": [{"id": 1, "name": "tag-a"}],
			"linkedItems": [
				{"name": "Sub Widget", "amount": 10.00, "status": "OPEN", "categories": [{"id": 2, "name": "test2"}], "tags": [{"id": 2, "name": "tag-b"}]}
			]
		},
		{"name": "Gadget", "amount": 25.00, "status": "RESOLVED", "isTaxed": false}
	],
	"comments": [{"comment": "Looks good", "userId": %d}],
	"customFields": [
		{"customFieldId": %d, "stringValue": "hello world"},
		{"customFieldId": %d, "dateValue": "2024-06-15T00:00:00Z"},
		{"customFieldId": %d, "selectValue": %d},
		{"customFieldId": %d, "currencyValue": 12.34},
		{"customFieldId": %d, "booleanValue": true}
	]
}`

// TestQuickScan_PersistsEveryField is the anti-drop centerpiece: a maximal AI response is created
// through QuickScan and every scalar and nested association is asserted on the read-back receipt.
func TestQuickScan_PersistsEveryField(t *testing.T) {
	defer repositories.TruncateTestDb()

	cf := seedQuickScanCustomFields(t)

	created, user, group, err := runQuickScan(
		t,
		func(u models.User) uint { return 0 }, // AI provides paid-by; no backfill
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			return fmt.Sprintf(
				fullReceiptAiJSON,
				u.ID, u.ID, u.ID,
				cf.textId, cf.dateId, cf.selectId, cf.optionId, cf.currencyId, cf.booleanId,
			)
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error creating receipt from AI response")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error reading receipt back")
		return
	}

	// Scalars.
	if receipt.Name != "Full Receipt" {
		utils.PrintTestError(t, receipt.Name, "Full Receipt")
	}
	if !receipt.Amount.Equal(decimal.NewFromInt(100)) {
		utils.PrintTestError(t, receipt.Amount.String(), "100")
	}
	if receipt.Date.UTC().Format("2006-01-02") != "2024-05-01" {
		utils.PrintTestError(t, receipt.Date.UTC().Format("2006-01-02"), "2024-05-01")
	}
	if receipt.PaidByUserID != user.ID {
		utils.PrintTestError(t, receipt.PaidByUserID, user.ID)
	}
	if receipt.Status != models.OPEN {
		utils.PrintTestError(t, receipt.Status, models.OPEN)
	}
	if receipt.GroupId != group.ID {
		utils.PrintTestError(t, receipt.GroupId, group.ID)
	}

	// Receipt-level categories & tags.
	if len(receipt.Categories) != 2 || !hasCategoryId(receipt.Categories, 1) || !hasCategoryId(receipt.Categories, 2) {
		utils.PrintTestError(t, receipt.Categories, "categories 1 and 2")
	}
	if len(receipt.Tags) != 1 || !hasTagId(receipt.Tags, 1) {
		utils.PrintTestError(t, receipt.Tags, "tag 1")
	}

	// Items (linked items are filtered out of the top level, so 2 parents remain).
	if len(receipt.ReceiptItems) != 2 {
		utils.PrintTestError(t, len(receipt.ReceiptItems), 2)
		return
	}

	widget, ok := findItemByName(receipt.ReceiptItems, "Widget")
	if !ok {
		utils.PrintTestError(t, receipt.ReceiptItems, "an item named Widget")
		return
	}
	if !widget.Amount.Equal(decimal.NewFromInt(40)) {
		utils.PrintTestError(t, widget.Amount.String(), "40")
	}
	if widget.Status != models.ITEM_OPEN {
		utils.PrintTestError(t, widget.Status, models.ITEM_OPEN)
	}
	if widget.ChargedToUserId == nil || *widget.ChargedToUserId != user.ID {
		utils.PrintTestError(t, widget.ChargedToUserId, user.ID)
	}
	if !widget.IsTaxed {
		utils.PrintTestError(t, widget.IsTaxed, true)
	}
	if len(widget.Categories) != 1 || !hasCategoryId(widget.Categories, 1) {
		utils.PrintTestError(t, widget.Categories, "item category 1")
	}
	if len(widget.Tags) != 1 || !hasTagId(widget.Tags, 1) {
		utils.PrintTestError(t, widget.Tags, "item tag 1")
	}

	// Linked item under Widget, with its own categories/tags.
	if len(widget.LinkedItems) != 1 {
		utils.PrintTestError(t, len(widget.LinkedItems), 1)
		return
	}
	sub := widget.LinkedItems[0]
	if sub.Name != "Sub Widget" {
		utils.PrintTestError(t, sub.Name, "Sub Widget")
	}
	if !sub.Amount.Equal(decimal.NewFromInt(10)) {
		utils.PrintTestError(t, sub.Amount.String(), "10")
	}
	if len(sub.Categories) != 1 || !hasCategoryId(sub.Categories, 2) {
		utils.PrintTestError(t, sub.Categories, "linked item category 2")
	}
	if len(sub.Tags) != 1 || !hasTagId(sub.Tags, 2) {
		utils.PrintTestError(t, sub.Tags, "linked item tag 2")
	}

	gadget, ok := findItemByName(receipt.ReceiptItems, "Gadget")
	if !ok {
		utils.PrintTestError(t, receipt.ReceiptItems, "an item named Gadget")
		return
	}
	if !gadget.Amount.Equal(decimal.NewFromInt(25)) {
		utils.PrintTestError(t, gadget.Amount.String(), "25")
	}
	if gadget.Status != models.ITEM_RESOLVED {
		utils.PrintTestError(t, gadget.Status, models.ITEM_RESOLVED)
	}
	if gadget.IsTaxed {
		utils.PrintTestError(t, gadget.IsTaxed, false)
	}
	if gadget.ChargedToUserId != nil {
		utils.PrintTestError(t, gadget.ChargedToUserId, "nil charged-to")
	}

	// Comment.
	if len(receipt.Comments) != 1 {
		utils.PrintTestError(t, len(receipt.Comments), 1)
		return
	}
	if receipt.Comments[0].Comment != "Looks good" {
		utils.PrintTestError(t, receipt.Comments[0].Comment, "Looks good")
	}
	if receipt.Comments[0].UserId == nil || *receipt.Comments[0].UserId != user.ID {
		utils.PrintTestError(t, receipt.Comments[0].UserId, user.ID)
	}

	// Custom-field values: each must land in the right column with the right field id, and nothing
	// else. Even though Validate ignores custom fields, they must persist unchanged.
	if len(receipt.CustomFields) != 5 {
		utils.PrintTestError(t, len(receipt.CustomFields), 5)
		return
	}

	if text, ok := findCustomFieldValue(receipt.CustomFields, cf.textId); !ok {
		utils.PrintTestError(t, receipt.CustomFields, "a value for the text field")
	} else if text.StringValue == nil || *text.StringValue != "hello world" {
		utils.PrintTestError(t, text.StringValue, "hello world")
	} else if text.CustomField.Type != models.TEXT {
		utils.PrintTestError(t, text.CustomField.Type, models.TEXT)
	}

	if date, ok := findCustomFieldValue(receipt.CustomFields, cf.dateId); !ok {
		utils.PrintTestError(t, receipt.CustomFields, "a value for the date field")
	} else if date.DateValue == nil || date.DateValue.UTC().Format("2006-01-02") != "2024-06-15" {
		utils.PrintTestError(t, date.DateValue, "2024-06-15")
	}

	if selectValue, ok := findCustomFieldValue(receipt.CustomFields, cf.selectId); !ok {
		utils.PrintTestError(t, receipt.CustomFields, "a value for the select field")
	} else if selectValue.SelectValue == nil || *selectValue.SelectValue != cf.optionId {
		utils.PrintTestError(t, selectValue.SelectValue, cf.optionId)
	} else if len(selectValue.CustomField.Options) != 1 {
		utils.PrintTestError(t, selectValue.CustomField.Options, "one preloaded option")
	}

	if currency, ok := findCustomFieldValue(receipt.CustomFields, cf.currencyId); !ok {
		utils.PrintTestError(t, receipt.CustomFields, "a value for the currency field")
	} else if currency.CurrencyValue == nil || !currency.CurrencyValue.Equal(decimal.RequireFromString("12.34")) {
		utils.PrintTestError(t, currency.CurrencyValue, "12.34")
	}

	if boolean, ok := findCustomFieldValue(receipt.CustomFields, cf.booleanId); !ok {
		utils.PrintTestError(t, receipt.CustomFields, "a value for the boolean field")
	} else if boolean.BooleanValue == nil || *boolean.BooleanValue != true {
		utils.PrintTestError(t, boolean.BooleanValue, true)
	}
}

// TestQuickScan_BackfillsPaidByAndStatusWhenAiOmitsThem proves QuickScan fills the paid-by/status
// from its arguments (the group's resolved defaults) when the AI omits them.
func TestQuickScan_BackfillsPaidByAndStatusWhenAiOmitsThem(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, user, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID }, // handler-resolved default paid-by
		models.NEEDS_ATTENTION,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			return `{"name": "No Payer", "amount": 12.50, "date": "2024-03-03T00:00:00Z"}`
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if receipt.PaidByUserID != user.ID {
		utils.PrintTestError(t, receipt.PaidByUserID, user.ID)
	}
	if receipt.Status != models.NEEDS_ATTENTION {
		utils.PrintTestError(t, receipt.Status, models.NEEDS_ATTENTION)
	}
}

// TestQuickScan_KeepsAiPaidByAndStatusWhenPresent proves an AI-supplied paid-by/status is not
// overwritten by the QuickScan default arguments.
func TestQuickScan_KeepsAiPaidByAndStatusWhenPresent(t *testing.T) {
	defer repositories.TruncateTestDb()

	var aiPayerId uint
	created, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID }, // default paid-by (should be ignored)
		models.OPEN,                              // default status (should be ignored)
		nil,
		nil,
		func(u models.User, g models.Group) string {
			// A second, distinct user for the AI paid-by so the assertion fails if
			// QuickScan drops the AI value and falls back to the default arg (the
			// seeded user). Same-user ids would pass either way.
			aiPayer := models.User{Username: "ai-payer", Password: "p", DisplayName: "x"}
			if err := repositories.GetDB().Create(&aiPayer).Error; err != nil {
				t.Fatalf("seed ai payer: %v", err)
			}
			aiPayerId = aiPayer.ID
			return fmt.Sprintf(
				`{"name": "Has Payer", "amount": 5.00, "date": "2024-02-02T00:00:00Z", "paidByUserId": %d, "status": "RESOLVED"}`,
				aiPayer.ID,
			)
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if receipt.PaidByUserID != aiPayerId {
		utils.PrintTestError(t, receipt.PaidByUserID, aiPayerId)
	}
	if receipt.Status != models.RESOLVED {
		utils.PrintTestError(t, receipt.Status, models.RESOLVED)
	}
}

// TestQuickScan_MergesUserCategoryAndTagPicksWithAi proves the user's per-file category/tag picks are
// unioned (deduped by id) with the AI-assigned ones, with names resolved so validation passes.
func TestQuickScan_MergesUserCategoryAndTagPicksWithAi(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		[]uint{2, 3}, // user category picks
		[]uint{2},    // user tag pick
		func(u models.User, g models.Group) string {
			// AI assigns category 1 and tag 1 (with names, as required by validation).
			return `{"name": "Merge", "amount": 9.99, "date": "2024-01-01T00:00:00Z", "categories": [{"id": 1, "name": "test"}], "tags": [{"id": 1, "name": "tag-a"}]}`
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(receipt.Categories) != 3 ||
		!hasCategoryId(receipt.Categories, 1) ||
		!hasCategoryId(receipt.Categories, 2) ||
		!hasCategoryId(receipt.Categories, 3) {
		utils.PrintTestError(t, receipt.Categories, "categories 1,2,3 unioned")
	}
	if len(receipt.Tags) != 2 || !hasTagId(receipt.Tags, 1) || !hasTagId(receipt.Tags, 2) {
		utils.PrintTestError(t, receipt.Tags, "tags 1,2 unioned")
	}
}

// TestQuickScan_RefundNegativeAmountsPersist proves a refund (negative receipt and item amounts)
// round-trips - the item's abs amount is still <= the receipt's abs amount, so validation passes.
func TestQuickScan_RefundNegativeAmountsPersist(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			return `{"name": "Refund", "amount": -30.00, "date": "2024-04-04T00:00:00Z", "receiptItems": [{"name": "Returned", "amount": -30.00, "status": "OPEN"}]}`
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !receipt.Amount.Equal(decimal.NewFromInt(-30)) {
		utils.PrintTestError(t, receipt.Amount.String(), "-30")
	}
	if len(receipt.ReceiptItems) != 1 {
		utils.PrintTestError(t, len(receipt.ReceiptItems), 1)
		return
	}
	if !receipt.ReceiptItems[0].Amount.Equal(decimal.NewFromInt(-30)) {
		utils.PrintTestError(t, receipt.ReceiptItems[0].Amount.String(), "-30")
	}
}

// TestQuickScan_MinimalReceiptCreatesNoSpuriousAssociations proves a bare AI response (only
// name/amount/date) persists with the backfilled paid-by/status and no invented associations.
func TestQuickScan_MinimalReceiptCreatesNoSpuriousAssociations(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, user, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			return `{"name": "Bare", "amount": 1.00, "date": "2024-07-07T00:00:00Z"}`
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if receipt.Name != "Bare" {
		utils.PrintTestError(t, receipt.Name, "Bare")
	}
	if receipt.PaidByUserID != user.ID {
		utils.PrintTestError(t, receipt.PaidByUserID, user.ID)
	}
	if receipt.Status != models.OPEN {
		utils.PrintTestError(t, receipt.Status, models.OPEN)
	}
	if len(receipt.ReceiptItems) != 0 {
		utils.PrintTestError(t, len(receipt.ReceiptItems), 0)
	}
	if len(receipt.Categories) != 0 {
		utils.PrintTestError(t, len(receipt.Categories), 0)
	}
	if len(receipt.Tags) != 0 {
		utils.PrintTestError(t, len(receipt.Tags), 0)
	}
	if len(receipt.Comments) != 0 {
		utils.PrintTestError(t, len(receipt.Comments), 0)
	}
	if len(receipt.CustomFields) != 0 {
		utils.PrintTestError(t, len(receipt.CustomFields), 0)
	}
}

// TestQuickScan_MissingItemStatusFailsAndPersistsNothing proves receipt validation runs on the
// AI-built command: an item missing its (required) status fails the whole quick scan, and no receipt
// row is written (all-or-nothing).
func TestQuickScan_MissingItemStatusFailsAndPersistsNothing(t *testing.T) {
	defer repositories.TruncateTestDb()

	_, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			// Item has no status -> UpsertItemCommand.Validate fails.
			return `{"name": "Bad Item", "amount": 20.00, "date": "2024-08-08T00:00:00Z", "receiptItems": [{"name": "No Status", "amount": 5.00}]}`
		},
	)
	if err == nil {
		utils.PrintTestError(t, nil, "a validation error")
	}

	// Only receipts created through CreateReceipt carry a created_by; the pipeline's seeded baseline
	// receipt leaves it NULL. Zero here means the failed scan wrote nothing.
	if count := quickScanCreatedReceiptCount(); count != 0 {
		utils.PrintTestError(t, count, int64(0))
	}

	// The pre-transaction failure must be recorded (not just logged): the QUICK_SCAN parent flips to
	// FAILED and a FAILED RECEIPT_UPLOADED task carries the error.
	if parent, ok := quickScanSystemTaskByType(models.QUICK_SCAN); !ok || parent.Status != models.SYSTEM_TASK_FAILED {
		utils.PrintTestError(t, parent.Status, models.SYSTEM_TASK_FAILED)
	}
	if uploaded, ok := quickScanSystemTaskByType(models.RECEIPT_UPLOADED); !ok || uploaded.Status != models.SYSTEM_TASK_FAILED {
		utils.PrintTestError(t, "missing FAILED RECEIPT_UPLOADED task", "a FAILED RECEIPT_UPLOADED task")
	}
}

// TestQuickScan_InvalidAiJsonReturnsError proves a non-JSON AI response surfaces as an error from
// QuickScan and persists nothing.
func TestQuickScan_InvalidAiJsonReturnsError(t *testing.T) {
	defer repositories.TruncateTestDb()

	_, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			return `this is not json`
		},
	)
	if err == nil {
		utils.PrintTestError(t, nil, "a JSON parse error")
	}

	if count := quickScanCreatedReceiptCount(); count != 0 {
		utils.PrintTestError(t, count, int64(0))
	}
}

// TestQuickScan_ResolvesIdOnlyAiCategory covers the quick-scan bug: the default prompt tells the AI to
// return categories/tags by id only (no name), which previously failed receipt validation ("name is
// required") and silently dropped the whole receipt. QuickScan now resolves each id to its real record
// so the receipt is created with the name filled in.
func TestQuickScan_ResolvesIdOnlyAiCategory(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			// Category/tag by id only - the exact shape the default prompt produces.
			return `{"name": "Walmart", "amount": 98.21, "date": "2024-01-01T00:00:00Z", "categories": [{"id": 1}], "tags": [{"id": 1}]}`
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(receipt.Categories) != 1 || !hasCategoryId(receipt.Categories, 1) {
		utils.PrintTestError(t, receipt.Categories, "category 1 resolved")
		return
	}
	if receipt.Categories[0].Name != "test" {
		utils.PrintTestError(t, receipt.Categories[0].Name, "test")
	}
	if len(receipt.Tags) != 1 || !hasTagId(receipt.Tags, 1) {
		utils.PrintTestError(t, receipt.Tags, "tag 1 resolved")
		return
	}
	if receipt.Tags[0].Name != "tag-a" {
		utils.PrintTestError(t, receipt.Tags[0].Name, "tag-a")
	}
}

// TestQuickScan_DropsUnresolvableAiCategory proves a hallucinated AI category id (no matching row) is
// dropped rather than failing the whole scan.
func TestQuickScan_DropsUnresolvableAiCategory(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			return `{"name": "Junk", "amount": 5.00, "date": "2024-01-01T00:00:00Z", "categories": [{"id": 1}, {"id": 999}]}`
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(receipt.Categories) != 1 || !hasCategoryId(receipt.Categories, 1) {
		utils.PrintTestError(t, receipt.Categories, "only category 1 (999 dropped)")
	}
}

// TestQuickScan_DropsOutOfGrantAiCategory proves an AI-assigned category the triggering user isn't
// allowed to see (their group role grants only a subset) is dropped, matching the write-side grant
// enforcement.
func TestQuickScan_DropsOutOfGrantAiCategory(t *testing.T) {
	defer repositories.TruncateTestDb()

	created, _, _, err := runQuickScanWithSetup(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			// AI assigns categories 1 and 2; the user's role grants only category 1.
			return `{"name": "Restricted", "amount": 7.00, "date": "2024-01-01T00:00:00Z", "categories": [{"id": 1}, {"id": 2}]}`
		},
		func(t *testing.T, user models.User, group models.Group) {
			clearGroupRoleGrantCacheAll()
			clearRolePermissionCacheAll()
			role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
				"restricted-cat", "", []string{permissions.GroupReceiptsRead}, []uint{1}, nil, nil, false)
			if err != nil {
				t.Fatalf("create restricted role: %v", err)
			}
			if err := repositories.GetDB().Model(&models.GroupMember{}).
				Where("group_id = ? AND user_id = ?", group.ID, user.ID).
				Update("group_role_id", role.ID).Error; err != nil {
				t.Fatalf("assign role to member: %v", err)
			}
		},
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	receipt, err := repositories.NewReceiptRepository(nil).GetFullyLoadedReceiptById(utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if len(receipt.Categories) != 1 || !hasCategoryId(receipt.Categories, 1) {
		utils.PrintTestError(t, receipt.Categories, "only category 1 (2 out-of-grant dropped)")
	}
}

// TestQuickScan_ValidationFailureRecordsFailedSystemTask covers the "missing system task" gap: a
// failure AFTER AI processing succeeds but BEFORE the receipt is created (here receipt validation) used
// to vanish into the log, leaving the AI tasks marked SUCCEEDED and no record of why no receipt
// appeared. The AI JSON parses (so the QUICK_SCAN task is SUCCEEDED), then fails validation (no name);
// QuickScan now records a FAILED RECEIPT_UPLOADED task and flips the QUICK_SCAN parent to FAILED.
func TestQuickScan_ValidationFailureRecordsFailedSystemTask(t *testing.T) {
	defer repositories.TruncateTestDb()

	_, _, _, err := runQuickScan(
		t,
		func(u models.User) uint { return u.ID },
		models.OPEN,
		nil,
		nil,
		func(u models.User, g models.Group) string {
			// Parses fine (QUICK_SCAN succeeds) but has no name -> validation fails after AI processing.
			return `{"amount": 12.34, "date": "2024-01-01T00:00:00Z"}`
		},
	)
	if err == nil {
		utils.PrintTestError(t, nil, "a validation error")
	}
	if count := quickScanCreatedReceiptCount(); count != 0 {
		utils.PrintTestError(t, count, int64(0))
	}

	// The QUICK_SCAN parent must be flipped to FAILED - this is what the activity feed shows.
	parent, ok := quickScanSystemTaskByType(models.QUICK_SCAN)
	if !ok {
		utils.PrintTestError(t, "no QUICK_SCAN task", "a QUICK_SCAN task")
		return
	}
	if parent.Status != models.SYSTEM_TASK_FAILED {
		utils.PrintTestError(t, parent.Status, models.SYSTEM_TASK_FAILED)
	}

	// A FAILED RECEIPT_UPLOADED task carrying the validation error must be recorded, attributed to the
	// real group.
	uploaded, ok := quickScanSystemTaskByType(models.RECEIPT_UPLOADED)
	if !ok {
		utils.PrintTestError(t, "no RECEIPT_UPLOADED task", "a FAILED RECEIPT_UPLOADED task")
		return
	}
	if uploaded.Status != models.SYSTEM_TASK_FAILED {
		utils.PrintTestError(t, uploaded.Status, models.SYSTEM_TASK_FAILED)
	}
	if !strings.Contains(uploaded.ResultDescription, "validation failed") {
		utils.PrintTestError(t, uploaded.ResultDescription, "contains 'validation failed'")
	}
	if uploaded.GroupId == nil || *uploaded.GroupId == 0 {
		utils.PrintTestError(t, uploaded.GroupId, "the real group id")
	}
}

// TestCombineEarlyFailureErrors pins the pre-transaction failure-wrapping used by
// recordEarlyQuickScanFailure: when system-task recording itself fails, both the original failure and
// the recording error must stay reachable via errors.Is (the recording error is wrapped, not just
// string-interpolated), with the original failure leading.
func TestCombineEarlyFailureErrors(t *testing.T) {
	failureErr := errors.New("receipt validation failed")
	taskErr := errors.New("system task write failed")

	// No recording error: the original failure passes through unchanged.
	if got := combineEarlyFailureErrors(failureErr, nil); got != failureErr {
		utils.PrintTestError(t, got, failureErr)
	}

	// Both errors stay reachable via errors.Is, with failureErr's message leading.
	combined := combineEarlyFailureErrors(failureErr, taskErr)
	if !errors.Is(combined, failureErr) {
		utils.PrintTestError(t, combined, "errors.Is finds failureErr")
	}
	if !errors.Is(combined, taskErr) {
		utils.PrintTestError(t, combined, "errors.Is finds taskErr")
	}
	if !strings.Contains(combined.Error(), "receipt validation failed") {
		utils.PrintTestError(t, combined.Error(), "message leads with the original failure")
	}
}

// quickScanCreatedReceiptCount counts receipts written through CreateReceipt (which stamps
// created_by), excluding the pipeline's seeded baseline receipt (created_by NULL). It is the "nothing
// was persisted" probe for the failure-path tests.
func quickScanCreatedReceiptCount() int64 {
	var count int64
	repositories.GetDB().Model(&models.Receipt{}).Where("created_by IS NOT NULL").Count(&count)
	return count
}

// quickScanSystemTaskByType fetches the single system task of the given type written by a runQuickScan
// call (all carry asynqTaskId "test-task"). Returns false when none exists.
func quickScanSystemTaskByType(taskType models.SystemTaskType) (models.SystemTask, bool) {
	var task models.SystemTask
	err := repositories.GetDB().
		Where("asynq_task_id = ? AND type = ?", "test-task", taskType).
		First(&task).Error
	if err != nil {
		return models.SystemTask{}, false
	}
	return task, true
}

// TestMagicFill_ParsesAllReceiptFields isolates the deserialization layer: a maximal AI response is
// parsed by MagicFillFromImage and the returned UpsertReceiptCommand is asserted to carry every
// field, so a failure localizes to parse (this test) vs. persist (the QuickScan tests). It also
// covers cleanResponse stripping a ```json code fence (the common LLM wrapping).
func TestMagicFill_ParsesAllReceiptFields(t *testing.T) {
	defer repositories.TruncateTestDb()

	cf := seedQuickScanCustomFields(t)

	var body string
	server := newMutableOllamaServer(t, &body)
	user, group, _ := seedReceiptImagePipeline(t, server.URL)

	receiptJSON := "```json\n" + fmt.Sprintf(
		fullReceiptAiJSON,
		user.ID, user.ID, user.ID,
		cf.textId, cf.dateId, cf.selectId, cf.optionId, cf.currencyId, cf.booleanId,
	) + "\n```" // code fence exercises cleanResponse
	body = ollamaReceiptResponse(t, receiptJSON)

	jpg, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	command, _, err := MagicFillFromImage(
		commands.MagicFillCommand{ImageData: jpg},
		utils.UintToString(group.ID),
		user.ID,
	)
	if err != nil {
		utils.PrintTestError(t, err, "no error parsing AI response")
		return
	}

	if command.Name != "Full Receipt" {
		utils.PrintTestError(t, command.Name, "Full Receipt")
	}
	if !command.Amount.Equal(decimal.NewFromInt(100)) {
		utils.PrintTestError(t, command.Amount.String(), "100")
	}
	if len(command.Categories) != 2 {
		utils.PrintTestError(t, len(command.Categories), 2)
	}
	if len(command.Tags) != 1 {
		utils.PrintTestError(t, len(command.Tags), 1)
	}
	if len(command.Items) != 2 {
		utils.PrintTestError(t, len(command.Items), 2)
		return
	}
	if len(command.Comments) != 1 {
		utils.PrintTestError(t, len(command.Comments), 1)
	}
	if len(command.CustomFields) != 5 {
		utils.PrintTestError(t, len(command.CustomFields), 5)
	}

	// Nested fields on the first item survive the parse.
	widget, ok := findUpsertItemByName(command.Items, "Widget")
	if !ok {
		utils.PrintTestError(t, command.Items, "an item named Widget")
		return
	}
	if widget.ChargedToUserId == nil || *widget.ChargedToUserId != user.ID {
		utils.PrintTestError(t, widget.ChargedToUserId, user.ID)
	}
	if !widget.IsTaxed {
		utils.PrintTestError(t, widget.IsTaxed, true)
	}
	if len(widget.LinkedItems) != 1 {
		utils.PrintTestError(t, len(widget.LinkedItems), 1)
	}
}

func findUpsertItemByName(items []commands.UpsertItemCommand, name string) (commands.UpsertItemCommand, bool) {
	for _, item := range items {
		if item.Name == name {
			return item, true
		}
	}
	return commands.UpsertItemCommand{}, false
}
