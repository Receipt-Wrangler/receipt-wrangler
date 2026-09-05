package repositories

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

// Sorting the receipts table by a custom field. The value lives in another table,
// in a column chosen by the field's type, so each type is exercised separately.
//
// The suite runs on SQLite only (see InitTestDb), which leaves two things
// uncovered: NULL ordering, which SQLite and MySQL put first while Postgres puts
// it last, and the MySQL/Postgres form of the CURRENCY cast. The cast is ANSI, so
// it needs no per-engine test, but the ordering of value-less receipts is
// deliberately not asserted below.

// customFieldsByType indexes the fields createTestCustomFields seeds so a test
// can name the type it cares about rather than an id.
func customFieldsByType() map[models.CustomFieldType]models.CustomField {
	var customFields []models.CustomField
	GetDB().Find(&customFields)

	byType := make(map[models.CustomFieldType]models.CustomField, len(customFields))
	for _, customField := range customFields {
		byType[customField.Type] = customField
	}

	return byType
}

func createSortableReceipt(name string) models.Receipt {
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromFloat(1),
		Date:         time.Now(),
		PaidByUserID: 1,
		Status:       models.OPEN,
		GroupId:      1,
	}
	GetDB().Create(&receipt)

	return receipt
}

func sortByCustomField(customFieldId uint, sortDirection commands.SortDirection) commands.ReceiptPagedRequestCommand {
	return commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      50,
			OrderBy:       "custom_" + utils.UintToString(customFieldId),
			SortDirection: sortDirection,
		},
	}
}

// sortedNames runs the request and returns the receipt names in the order the
// query produced them, which is what every assertion below compares.
func sortedNames(t *testing.T, pagedRequest commands.ReceiptPagedRequestCommand) []string {
	repository := NewReceiptRepository(nil)

	receipts, count, err := repository.GetPagedReceiptsByGroupId(1, "1", pagedRequest, nil, nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return nil
	}

	if pagedRequest.PageSize >= int(count) && int(count) != len(receipts) {
		utils.PrintTestError(t, count, len(receipts))
	}

	names := make([]string, len(receipts))
	for i, receipt := range receipts {
		names[i] = receipt.Name
	}

	return names
}

func assertOrder(t *testing.T, got []string, want []string) {
	if len(got) != len(want) {
		utils.PrintTestError(t, got, want)
		return
	}

	for i := range want {
		if got[i] != want[i] {
			utils.PrintTestError(t, got, want)
			return
		}
	}
}

func TestShouldSortReceiptsByTextCustomField(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	textField := customFieldsByType()[models.TEXT]
	db := GetDB()

	for name, value := range map[string]string{"charlie": "c", "alpha": "a", "bravo": "b"} {
		receipt := createSortableReceipt(name)
		stringValue := value
		db.Create(&models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: textField.ID,
			StringValue:   &stringValue,
		})
	}

	assertOrder(t, sortedNames(t, sortByCustomField(textField.ID, commands.ASCENDING)),
		[]string{"alpha", "bravo", "charlie"})
	assertOrder(t, sortedNames(t, sortByCustomField(textField.ID, commands.DESCENDING)),
		[]string{"charlie", "bravo", "alpha"})
}

// The one that a naive sort gets wrong. currency_value is a text column, so
// without the cast this orders "-5", "100", "12.34", "20" - lexicographically.
func TestShouldSortReceiptsByCurrencyCustomFieldNumerically(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	currencyField := customFieldsByType()[models.CURRENCY]
	db := GetDB()

	for name, value := range map[string]string{
		"hundred": "100", "twenty": "20", "negative": "-5", "twelve": "12.34",
	} {
		receipt := createSortableReceipt(name)
		currencyValue := decimal.RequireFromString(value)
		db.Create(&models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: currencyField.ID,
			CurrencyValue: &currencyValue,
		})
	}

	assertOrder(t, sortedNames(t, sortByCustomField(currencyField.ID, commands.ASCENDING)),
		[]string{"negative", "twelve", "twenty", "hundred"})
	assertOrder(t, sortedNames(t, sortByCustomField(currencyField.ID, commands.DESCENDING)),
		[]string{"hundred", "twenty", "twelve", "negative"})
}

func TestShouldSortReceiptsByDateCustomField(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	dateField := customFieldsByType()[models.DATE]
	db := GetDB()

	for name, day := range map[string]int{"third": 3, "first": 1, "second": 2} {
		receipt := createSortableReceipt(name)
		dateValue := time.Date(2026, time.March, day, 0, 0, 0, 0, time.UTC)
		db.Create(&models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: dateField.ID,
			DateValue:     &dateValue,
		})
	}

	assertOrder(t, sortedNames(t, sortByCustomField(dateField.ID, commands.ASCENDING)),
		[]string{"first", "second", "third"})
}

// A select stores an option id, but readers see the option's text, so that is
// what the column sorts by. The ids here are deliberately the reverse of the
// alphabetical order they must come back in.
func TestShouldSortReceiptsBySelectCustomFieldOptionText(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	selectField := customFieldsByType()[models.SELECT]
	db := GetDB()

	var options []models.CustomFieldOption
	db.Where("custom_field_id = ?", selectField.ID).Order("id").Find(&options)
	if len(options) != 2 {
		utils.PrintTestError(t, len(options), 2)
		return
	}

	db.Model(&models.CustomFieldOption{}).Where("id = ?", options[0].ID).Update("value", "Zulu")
	db.Model(&models.CustomFieldOption{}).Where("id = ?", options[1].ID).Update("value", "Alpha")

	for name, optionIndex := range map[string]int{"zulu-receipt": 0, "alpha-receipt": 1} {
		receipt := createSortableReceipt(name)
		selectValue := options[optionIndex].ID
		db.Create(&models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: selectField.ID,
			SelectValue:   &selectValue,
		})
	}

	assertOrder(t, sortedNames(t, sortByCustomField(selectField.ID, commands.ASCENDING)),
		[]string{"alpha-receipt", "zulu-receipt"})
}

func TestShouldSortReceiptsByBooleanCustomField(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	db := GetDB()
	booleanField := models.CustomField{Name: "Reimbursed", Type: models.BOOLEAN}
	db.Create(&booleanField)

	for name, value := range map[string]bool{"yes": true, "no": false} {
		receipt := createSortableReceipt(name)
		booleanValue := value
		db.Create(&models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: booleanField.ID,
			BooleanValue:  &booleanValue,
		})
	}

	assertOrder(t, sortedNames(t, sortByCustomField(booleanField.ID, commands.ASCENDING)),
		[]string{"no", "yes"})
	assertOrder(t, sortedNames(t, sortByCustomField(booleanField.ID, commands.DESCENDING)),
		[]string{"yes", "no"})
}

// A receipt without a value for the field must still be listed, and the total
// count must still match the rows returned - the regression a join would cause.
func TestShouldKeepReceiptsWithoutACustomFieldValueWhenSorting(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	textField := customFieldsByType()[models.TEXT]
	db := GetDB()

	withValue := createSortableReceipt("has-value")
	createSortableReceipt("no-value")

	stringValue := "a"
	db.Create(&models.CustomFieldValue{
		ReceiptId:     withValue.ID,
		CustomFieldId: textField.ID,
		StringValue:   &stringValue,
	})

	names := sortedNames(t, sortByCustomField(textField.ID, commands.ASCENDING))
	if len(names) != 2 {
		utils.PrintTestError(t, names, 2)
	}
}

// Nothing stops a receipt holding several values for one field. The lowest id
// among the values that actually resolve wins - the same rule the reporting
// engine reads by - so an empty low-id row must not hide a real one, and the
// receipt must appear exactly once.
func TestShouldResolveDuplicateCustomFieldValuesByLowestNonNullId(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	textField := customFieldsByType()[models.TEXT]
	db := GetDB()

	duplicated := createSortableReceipt("duplicated")
	other := createSortableReceipt("other")

	// Lowest id, but empty: it must not win.
	db.Create(&models.CustomFieldValue{ReceiptId: duplicated.ID, CustomFieldId: textField.ID})

	winning := "a"
	losing := "z"
	db.Create(&models.CustomFieldValue{
		ReceiptId: duplicated.ID, CustomFieldId: textField.ID, StringValue: &winning,
	})
	db.Create(&models.CustomFieldValue{
		ReceiptId: duplicated.ID, CustomFieldId: textField.ID, StringValue: &losing,
	})

	middle := "m"
	db.Create(&models.CustomFieldValue{
		ReceiptId: other.ID, CustomFieldId: textField.ID, StringValue: &middle,
	})

	// "a" wins over "z", so the duplicated receipt sorts first - and appears once.
	assertOrder(t, sortedNames(t, sortByCustomField(textField.ID, commands.ASCENDING)),
		[]string{"duplicated", "other"})
}

// Ties are broken by receipt id, so paging through equal values never repeats or
// skips a row. Without the tiebreaker the two pages below can return the same
// receipt twice.
func TestShouldPageStablyThroughEqualCustomFieldValues(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	textField := customFieldsByType()[models.TEXT]
	db := GetDB()

	for _, name := range []string{"one", "two", "three", "four"} {
		receipt := createSortableReceipt(name)
		stringValue := "same"
		db.Create(&models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: textField.ID,
			StringValue:   &stringValue,
		})
	}

	seen := make(map[string]bool)
	for page := 1; page <= 2; page++ {
		pagedRequest := sortByCustomField(textField.ID, commands.ASCENDING)
		pagedRequest.Page = page
		pagedRequest.PageSize = 2

		for _, name := range sortedNames(t, pagedRequest) {
			if seen[name] {
				utils.PrintTestError(t, "receipt "+name+" returned on two pages", "each receipt once")
			}
			seen[name] = true
		}
	}

	if len(seen) != 4 {
		utils.PrintTestError(t, len(seen), 4)
	}
}

// A client persists its sort, so a custom field deleted afterwards must not make
// every subsequent list load fail. It falls back to the default ordering.
func TestShouldFallBackToDefaultOrderForUnknownCustomField(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestReceipts()

	repository := NewReceiptRepository(nil)
	receipts, count, err := repository.GetPagedReceiptsByGroupId(
		1, "1", sortByCustomField(999999, commands.DESCENDING), nil, nil,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if count != 2 || len(receipts) != 2 {
		utils.PrintTestError(t, count, 2)
	}
}

// A malformed orderBy is still rejected outright - that is the guard keeping the
// sort column out of the SQL it is concatenated into.
func TestShouldRejectMalformedCustomFieldOrderBy(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestReceipts()

	repository := NewReceiptRepository(nil)

	for _, orderBy := range []string{"custom_abc", "custom_", "custom_1; DROP TABLE receipts", "custom_1_month"} {
		pagedRequest := commands.ReceiptPagedRequestCommand{
			PagedRequestCommand: commands.PagedRequestCommand{
				Page: 1, PageSize: 50, OrderBy: orderBy, SortDirection: commands.ASCENDING,
			},
		}

		_, _, err := repository.GetPagedReceiptsByGroupId(1, "1", pagedRequest, nil, nil)
		if err == nil {
			utils.PrintTestError(t, "no error for orderBy "+orderBy, "untrusted value error")
		}
	}
}

// The list response is deserialized by the mobile client into a typed Receipt,
// and its generated CustomFieldType enum has no empty member. Loading a value
// without its definition would serialize "type":"" and collapse every row, so the
// preload set must always carry the definitions.
func TestPagedReceiptCustomFieldsAlwaysCarryTheirDefinition(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()
	createTestCustomFields()

	textField := customFieldsByType()[models.TEXT]
	receipt := createSortableReceipt("with-custom-field")

	stringValue := "a"
	GetDB().Create(&models.CustomFieldValue{
		ReceiptId:     receipt.ID,
		CustomFieldId: textField.ID,
		StringValue:   &stringValue,
	})

	repository := NewReceiptRepository(nil)
	receipts, _, err := repository.GetPagedReceiptsByGroupId(
		1, "1", pagedRequestAllReceipts(), constants.CUSTOM_FIELD_ASSOCIATIONS, nil,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	bytes, err := json.Marshal(receipts)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if strings.Contains(string(bytes), `"type":""`) {
		utils.PrintTestError(t, string(bytes), `no empty "type" in the response`)
	}

	for _, marshalled := range receipts {
		for _, customFieldValue := range marshalled.CustomFields {
			if customFieldValue.CustomField.Type != models.TEXT {
				utils.PrintTestError(t, customFieldValue.CustomField.Type, models.TEXT)
			}
		}
	}
}
