package services

import (
	"encoding/json"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// The diff is pure, so these tests build receipts in memory and compare the
// change rows as JSON, which is exactly what a client receives.

func diffUint(value uint) *uint {
	return &value
}

func diffString(value string) *string {
	return &value
}

func diffBool(value bool) *bool {
	return &value
}

func diffDecimal(value string) *decimal.Decimal {
	amount := decimal.RequireFromString(value)
	return &amount
}

func diffItem(id uint, name string, amount string, chargedTo *uint) models.Item {
	return models.Item{
		BaseModel:       models.BaseModel{ID: id},
		Name:            name,
		Amount:          decimal.RequireFromString(amount),
		ChargedToUserId: chargedTo,
		Status:          models.ITEM_OPEN,
	}
}

func textChange(field string, kind string, from string, to string) structs.ReceiptHistoryChange {
	change := structs.ReceiptHistoryChange{Field: field, Kind: kind}
	if kind != structs.ReceiptHistoryKindAdded {
		change.From = textValue(from)
	}
	if kind != structs.ReceiptHistoryKindRemoved {
		change.To = textValue(to)
	}
	return change
}

func typedValue(valueType string, value string) *structs.ReceiptHistoryValue {
	return &structs.ReceiptHistoryValue{Type: valueType, Value: value}
}

func assertChanges(t *testing.T, got []structs.ReceiptHistoryChange, want []structs.ReceiptHistoryChange) {
	t.Helper()
	if want == nil {
		want = []structs.ReceiptHistoryChange{}
	}
	gotJson, _ := json.MarshalIndent(got, "", "  ")
	wantJson, _ := json.MarshalIndent(want, "", "  ")
	if string(gotJson) != string(wantJson) {
		t.Errorf("changes =\n%s\nwant\n%s", gotJson, wantJson)
	}
}

func emptyDiffContext() receiptDiffContext {
	return receiptDiffContext{groupNames: map[uint]string{1: "Household", 2: "Trip"}}
}

func TestDiffReceiptSnapshotsCreationListsOnlyTheFieldsSetEmailFirst(t *testing.T) {
	created := models.Receipt{
		BaseModel:    models.BaseModel{ID: 9, CreatedBy: diffUint(4), CreatedByString: "Email Integration"},
		Name:         "Groceries",
		Amount:       decimal.RequireFromString("42.10"),
		Date:         time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC),
		GroupId:      1,
		PaidByUserID: 3,
		Status:       models.OPEN,
		Categories:   []models.Category{{BaseModel: models.BaseModel{ID: 5}, Name: "Food"}},
		EmailInbox:   "receipts@example.com",
		EmailSender:  "store@example.com",
		ReceiptItems: []models.Item{diffItem(1, "Milk", "4.20", nil), diffItem(2, "Milk", "4.20", diffUint(3))},
		Comments:     []models.Comment{{BaseModel: models.BaseModel{ID: 6}, Comment: "Paid in cash", UserId: diffUint(3)}},
		ImageFiles:   []models.FileData{{BaseModel: models.BaseModel{ID: 8}, Name: "receipt.jpg"}},
	}

	got := diffReceiptSnapshots(nil, created, emptyDiffContext())

	assertChanges(t, got, []structs.ReceiptHistoryChange{
		textChange(structs.ReceiptHistoryFieldEmailInbox, structs.ReceiptHistoryKindAdded, "", "receipts@example.com"),
		textChange(structs.ReceiptHistoryFieldEmailSender, structs.ReceiptHistoryKindAdded, "", "store@example.com"),
		textChange(structs.ReceiptHistoryFieldName, structs.ReceiptHistoryKindAdded, "", "Groceries"),
		{Field: structs.ReceiptHistoryFieldAmount, Kind: structs.ReceiptHistoryKindAdded, To: typedValue(structs.ReceiptHistoryValueCurrency, "42.1")},
		{Field: structs.ReceiptHistoryFieldDate, Kind: structs.ReceiptHistoryKindAdded, To: typedValue(structs.ReceiptHistoryValueDate, "2026-09-30")},
		textChange(structs.ReceiptHistoryFieldGroup, structs.ReceiptHistoryKindAdded, "", "Household"),
		{Field: structs.ReceiptHistoryFieldPaidBy, Kind: structs.ReceiptHistoryKindAdded, To: &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueUser, UserId: diffUint(3)}},
		{Field: structs.ReceiptHistoryFieldStatus, Kind: structs.ReceiptHistoryKindAdded, To: typedValue(structs.ReceiptHistoryValueStatus, "OPEN")},
		textChange(structs.ReceiptHistoryFieldCategories, structs.ReceiptHistoryKindAdded, "", "Food"),
		{Field: structs.ReceiptHistoryFieldItem, Kind: structs.ReceiptHistoryKindAdded, Label: "Milk", To: typedValue(structs.ReceiptHistoryValueCurrency, "4.2")},
		{Field: structs.ReceiptHistoryFieldShare, Kind: structs.ReceiptHistoryKindAdded, Label: "Milk", UserId: diffUint(3), To: typedValue(structs.ReceiptHistoryValueCurrency, "4.2")},
		{Field: structs.ReceiptHistoryFieldComment, Kind: structs.ReceiptHistoryKindAdded, UserId: diffUint(3), To: textValue("Paid in cash")},
		textChange(structs.ReceiptHistoryFieldImage, structs.ReceiptHistoryKindAdded, "", "receipt.jpg"),
	})
}

func TestDiffReceiptSnapshotsCreationOmitsAZeroAmountAndEmptyFields(t *testing.T) {
	got := diffReceiptSnapshots(nil, models.Receipt{Name: "Blank", Status: models.DRAFT}, emptyDiffContext())

	assertChanges(t, got, []structs.ReceiptHistoryChange{
		textChange(structs.ReceiptHistoryFieldName, structs.ReceiptHistoryKindAdded, "", "Blank"),
		{Field: structs.ReceiptHistoryFieldStatus, Kind: structs.ReceiptHistoryKindAdded, To: typedValue(structs.ReceiptHistoryValueStatus, "DRAFT")},
	})
}

func TestDiffReceiptSnapshotsScalars(t *testing.T) {
	resolvedAt := time.Date(2026, 10, 2, 15, 4, 5, 0, time.UTC)
	base := func() models.Receipt {
		return models.Receipt{
			Name:         "Dinner",
			Amount:       decimal.RequireFromString("12.50"),
			Date:         time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			GroupId:      1,
			PaidByUserID: 2,
			Status:       models.OPEN,
		}
	}

	tests := []struct {
		name   string
		change func(receipt *models.Receipt)
		want   []structs.ReceiptHistoryChange
	}{
		{
			name:   "nothing changed",
			change: func(receipt *models.Receipt) {},
		},
		{
			name: "BaseModel fields are never compared",
			change: func(receipt *models.Receipt) {
				receipt.ID = 99
				receipt.UpdatedAt = time.Now()
				receipt.CreatedBy = diffUint(7)
				receipt.CreatedByString = "someone"
			},
		},
		{
			name:   "name",
			change: func(receipt *models.Receipt) { receipt.Name = "Supper" },
			want:   []structs.ReceiptHistoryChange{textChange(structs.ReceiptHistoryFieldName, structs.ReceiptHistoryKindChanged, "Dinner", "Supper")},
		},
		{
			name:   "an amount compares by value",
			change: func(receipt *models.Receipt) { receipt.Amount = decimal.RequireFromString("12.5") },
		},
		{
			name:   "amount",
			change: func(receipt *models.Receipt) { receipt.Amount = decimal.RequireFromString("0") },
			want: []structs.ReceiptHistoryChange{{
				Field: structs.ReceiptHistoryFieldAmount, Kind: structs.ReceiptHistoryKindChanged,
				From: typedValue(structs.ReceiptHistoryValueCurrency, "12.5"), To: typedValue(structs.ReceiptHistoryValueCurrency, "0"),
			}},
		},
		{
			name: "a date compares by its UTC calendar day",
			change: func(receipt *models.Receipt) {
				// Aug 31 in New York, but Sep 1 in UTC.
				receipt.Date = time.Date(2026, 8, 31, 22, 0, 0, 0, time.FixedZone("EDT", -4*3600))
			},
		},
		{
			name:   "date",
			change: func(receipt *models.Receipt) { receipt.Date = time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC) },
			want: []structs.ReceiptHistoryChange{{
				Field: structs.ReceiptHistoryFieldDate, Kind: structs.ReceiptHistoryKindChanged,
				From: typedValue(structs.ReceiptHistoryValueDate, "2026-09-01"), To: typedValue(structs.ReceiptHistoryValueDate, "2026-09-02"),
			}},
		},
		{
			name:   "group, named only where the caller may see it",
			change: func(receipt *models.Receipt) { receipt.GroupId = 3 },
			want:   []structs.ReceiptHistoryChange{textChange(structs.ReceiptHistoryFieldGroup, structs.ReceiptHistoryKindChanged, "Household", "(Restricted)")},
		},
		{
			name:   "paid by",
			change: func(receipt *models.Receipt) { receipt.PaidByUserID = 5 },
			want: []structs.ReceiptHistoryChange{{
				Field: structs.ReceiptHistoryFieldPaidBy, Kind: structs.ReceiptHistoryKindChanged,
				From: &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueUser, UserId: diffUint(2)},
				To:   &structs.ReceiptHistoryValue{Type: structs.ReceiptHistoryValueUser, UserId: diffUint(5)},
			}},
		},
		{
			name: "status and resolved date",
			change: func(receipt *models.Receipt) {
				receipt.Status = models.RESOLVED
				receipt.ResolvedDate = &resolvedAt
			},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldStatus, Kind: structs.ReceiptHistoryKindChanged, From: typedValue(structs.ReceiptHistoryValueStatus, "OPEN"), To: typedValue(structs.ReceiptHistoryValueStatus, "RESOLVED")},
				{Field: structs.ReceiptHistoryFieldResolvedDate, Kind: structs.ReceiptHistoryKindAdded, To: typedValue(structs.ReceiptHistoryValueInstant, "2026-10-02T15:04:05Z")},
			},
		},
		{
			name: "the email origin is never compared on an update",
			change: func(receipt *models.Receipt) {
				receipt.EmailInbox = "receipts@example.com"
				receipt.EmailSender = "store@example.com"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before, after := base(), base()
			test.change(&after)
			assertChanges(t, diffReceiptSnapshots(&before, after, emptyDiffContext()), test.want)
		})
	}

	t.Run("a resolved date that is cleared is removed", func(t *testing.T) {
		before, after := base(), base()
		before.ResolvedDate = &resolvedAt
		assertChanges(t, diffReceiptSnapshots(&before, after, emptyDiffContext()), []structs.ReceiptHistoryChange{
			{Field: structs.ReceiptHistoryFieldResolvedDate, Kind: structs.ReceiptHistoryKindRemoved, From: typedValue(structs.ReceiptHistoryValueInstant, "2026-10-02T15:04:05Z")},
		})
	})
}

func TestDiffReceiptSnapshotsCategoriesAndTagsCompareById(t *testing.T) {
	food := models.Category{BaseModel: models.BaseModel{ID: 1}, Name: "Food"}
	travel := models.Category{BaseModel: models.BaseModel{ID: 2}, Name: "Travel"}
	work := models.Tag{BaseModel: models.BaseModel{ID: 1}, Name: "Work"}

	before := models.Receipt{Categories: []models.Category{food, travel}, Tags: []models.Tag{work}}
	// Renaming a category keeps its id, so it is not a change.
	renamedFood := food
	renamedFood.Name = "Groceries"
	after := models.Receipt{
		Categories: []models.Category{renamedFood, {BaseModel: models.BaseModel{ID: 3}, Name: "Home"}},
		Tags:       []models.Tag{},
	}

	assertChanges(t, diffReceiptSnapshots(&before, after, emptyDiffContext()), []structs.ReceiptHistoryChange{
		textChange(structs.ReceiptHistoryFieldCategories, structs.ReceiptHistoryKindAdded, "", "Home"),
		textChange(structs.ReceiptHistoryFieldCategories, structs.ReceiptHistoryKindRemoved, "Travel", ""),
		textChange(structs.ReceiptHistoryFieldTags, structs.ReceiptHistoryKindRemoved, "Work", ""),
	})
}

func TestDiffReceiptSnapshotsCustomFields(t *testing.T) {
	poNumber := models.CustomField{BaseModel: models.BaseModel{ID: 1}, Name: "PO Number", Type: models.TEXT}
	project := models.CustomField{
		BaseModel: models.BaseModel{ID: 2}, Name: "Project", Type: models.SELECT,
		Options: []models.CustomFieldOption{{BaseModel: models.BaseModel{ID: 10}, Value: "Alpha"}, {BaseModel: models.BaseModel{ID: 11}, Value: "Beta"}},
	}
	dueDate := models.CustomField{BaseModel: models.BaseModel{ID: 3}, Name: "Due", Type: models.DATE}
	tip := models.CustomField{BaseModel: models.BaseModel{ID: 4}, Name: "Tip", Type: models.CURRENCY}
	billable := models.CustomField{BaseModel: models.BaseModel{ID: 5}, Name: "Billable", Type: models.BOOLEAN}

	value := func(definition models.CustomField, set func(value *models.CustomFieldValue)) models.CustomFieldValue {
		customFieldValue := models.CustomFieldValue{CustomFieldId: definition.ID, CustomField: definition}
		if set != nil {
			set(&customFieldValue)
		}
		return customFieldValue
	}
	due := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		before  []models.CustomFieldValue
		after   []models.CustomFieldValue
		context map[uint]models.CustomField
		want    []structs.ReceiptHistoryChange
	}{
		{
			name:   "an absent value equals an empty one",
			before: nil,
			after: []models.CustomFieldValue{
				value(poNumber, nil),
				value(poNumber, func(v *models.CustomFieldValue) { v.StringValue = diffString("") }),
				value(project, nil),
				value(dueDate, nil),
				value(tip, nil),
				value(billable, func(v *models.CustomFieldValue) { v.BooleanValue = diffBool(false) }),
			},
		},
		{
			name:   "text",
			before: []models.CustomFieldValue{value(poNumber, func(v *models.CustomFieldValue) { v.StringValue = diffString("A-1") })},
			after:  []models.CustomFieldValue{value(poNumber, func(v *models.CustomFieldValue) { v.StringValue = diffString("A-2") })},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindChanged, Label: "PO Number", From: textValue("A-1"), To: textValue("A-2")},
			},
		},
		{
			name:   "a select compares by option id and shows its text",
			before: []models.CustomFieldValue{value(project, func(v *models.CustomFieldValue) { v.SelectValue = diffUint(10) })},
			after:  []models.CustomFieldValue{value(project, func(v *models.CustomFieldValue) { v.SelectValue = diffUint(11) })},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindChanged, Label: "Project", From: textValue("Alpha"), To: textValue("Beta")},
			},
		},
		{
			name:   "date, currency and boolean",
			before: []models.CustomFieldValue{value(billable, func(v *models.CustomFieldValue) { v.BooleanValue = diffBool(false) })},
			after: []models.CustomFieldValue{
				value(dueDate, func(v *models.CustomFieldValue) { v.DateValue = &due }),
				value(tip, func(v *models.CustomFieldValue) { v.CurrencyValue = diffDecimal("2.50") }),
				value(billable, func(v *models.CustomFieldValue) { v.BooleanValue = diffBool(true) }),
			},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindAdded, Label: "Due", To: typedValue(structs.ReceiptHistoryValueDate, "2026-11-01")},
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindAdded, Label: "Tip", To: typedValue(structs.ReceiptHistoryValueCurrency, "2.5")},
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindChanged, Label: "Billable", From: typedValue(structs.ReceiptHistoryValueBoolean, "false"), To: typedValue(structs.ReceiptHistoryValueBoolean, "true")},
			},
		},
		{
			name:   "a cleared value is removed",
			before: []models.CustomFieldValue{value(billable, func(v *models.CustomFieldValue) { v.BooleanValue = diffBool(true) })},
			after:  []models.CustomFieldValue{},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindRemoved, Label: "Billable", From: typedValue(structs.ReceiptHistoryValueBoolean, "true")},
			},
		},
		{
			name: "a value stored without its definition is read against the stored one",
			before: []models.CustomFieldValue{
				{CustomFieldId: project.ID, SelectValue: diffUint(10)},
			},
			after:   []models.CustomFieldValue{value(project, func(v *models.CustomFieldValue) { v.SelectValue = diffUint(10) })},
			context: map[uint]models.CustomField{project.ID: project},
		},
		{
			name:   "the snapshot's own option text wins over the stored one",
			before: []models.CustomFieldValue{{CustomFieldId: project.ID, SelectValue: diffUint(10)}},
			after: []models.CustomFieldValue{value(models.CustomField{
				BaseModel: project.BaseModel, Name: "Project", Type: models.SELECT,
				Options: []models.CustomFieldOption{{BaseModel: models.BaseModel{ID: 11}, Value: "Beta (old name)"}},
			}, func(v *models.CustomFieldValue) { v.SelectValue = diffUint(11) })},
			context: map[uint]models.CustomField{project.ID: project},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldCustomField, Kind: structs.ReceiptHistoryKindChanged, Label: "Project", From: textValue("Alpha"), To: textValue("Beta (old name)")},
			},
		},
		{
			name:   "a value with no definition anywhere is skipped",
			before: []models.CustomFieldValue{{CustomFieldId: 77, StringValue: diffString("x")}},
			after:  []models.CustomFieldValue{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context := emptyDiffContext()
			context.customFields = test.context
			before := models.Receipt{CustomFields: test.before}
			after := models.Receipt{CustomFields: test.after}
			assertChanges(t, diffReceiptSnapshots(&before, after, context), test.want)
		})
	}
}

func TestDiffReceiptSnapshotsItemsAndShares(t *testing.T) {
	alice := diffUint(2)
	food := models.Category{BaseModel: models.BaseModel{ID: 1}, Name: "Food"}
	shared := models.Tag{BaseModel: models.BaseModel{ID: 1}, Name: "Shared"}

	tests := []struct {
		name   string
		before []models.Item
		after  []models.Item
		want   []structs.ReceiptHistoryChange
	}{
		{
			name:   "items recreated with new ids are not a change",
			before: []models.Item{diffItem(1, "Milk", "4.20", nil), diffItem(2, "Bread", "3.00", alice)},
			after:  []models.Item{diffItem(7, "Bread ", "3", alice), diffItem(8, "Milk", "4.2", nil)},
		},
		{
			name:   "duplicate names match equal amounts first",
			before: []models.Item{diffItem(1, "Milk", "4.20", nil), diffItem(2, "Milk", "2.00", nil)},
			after:  []models.Item{diffItem(3, "Milk", "2.00", nil), diffItem(4, "Milk", "5.00", nil), diffItem(5, "Milk", "1.00", nil)},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldItem, Kind: structs.ReceiptHistoryKindChanged, Label: "Milk", From: typedValue(structs.ReceiptHistoryValueCurrency, "4.2"), To: typedValue(structs.ReceiptHistoryValueCurrency, "5")},
				{Field: structs.ReceiptHistoryFieldItem, Kind: structs.ReceiptHistoryKindAdded, Label: "Milk", To: typedValue(structs.ReceiptHistoryValueCurrency, "1")},
			},
		},
		{
			name:   "added and removed items and shares",
			before: []models.Item{diffItem(1, "Eggs", "3.00", nil)},
			after:  []models.Item{diffItem(2, "Eggs", "3.00", alice)},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldItem, Kind: structs.ReceiptHistoryKindRemoved, Label: "Eggs", From: typedValue(structs.ReceiptHistoryValueCurrency, "3")},
				{Field: structs.ReceiptHistoryFieldShare, Kind: structs.ReceiptHistoryKindAdded, Label: "Eggs", UserId: alice, To: typedValue(structs.ReceiptHistoryValueCurrency, "3")},
			},
		},
		{
			name: "a share's status is reported, an item's is not",
			before: []models.Item{
				diffItem(1, "Milk", "4.20", nil),
				diffItem(2, "Bread", "3.00", alice),
			},
			after: func() []models.Item {
				item := diffItem(3, "Milk", "4.20", nil)
				item.Status = models.ITEM_RESOLVED
				share := diffItem(4, "Bread", "3.00", alice)
				share.Status = models.ITEM_RESOLVED
				return []models.Item{item, share}
			}(),
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldShare, Kind: structs.ReceiptHistoryKindChanged, Label: "Bread", UserId: alice, From: typedValue(structs.ReceiptHistoryValueStatus, "OPEN"), To: typedValue(structs.ReceiptHistoryValueStatus, "RESOLVED")},
			},
		},
		{
			name: "item categories and tags",
			before: []models.Item{func() models.Item {
				item := diffItem(1, "Milk", "4.20", nil)
				item.Tags = []models.Tag{shared}
				return item
			}()},
			after: []models.Item{func() models.Item {
				item := diffItem(2, "Milk", "4.20", nil)
				item.Categories = []models.Category{food}
				return item
			}()},
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldItemCategories, Kind: structs.ReceiptHistoryKindAdded, Label: "Milk", To: textValue("Food")},
				{Field: structs.ReceiptHistoryFieldItemTags, Kind: structs.ReceiptHistoryKindRemoved, Label: "Milk", From: textValue("Shared")},
			},
		},
		{
			name: "linked items are compared once each",
			before: func() []models.Item {
				linked := diffItem(2, "Tip", "1.00", alice)
				parent := diffItem(1, "Meal", "10.00", nil)
				parent.LinkedItems = []models.Item{linked}
				// A version 1 snapshot also lists the linked item at the top level.
				return []models.Item{parent, linked}
			}(),
			after: func() []models.Item {
				parent := diffItem(3, "Meal", "10.00", nil)
				parent.LinkedItems = []models.Item{diffItem(4, "Tip", "2.00", alice)}
				return []models.Item{parent}
			}(),
			want: []structs.ReceiptHistoryChange{
				{Field: structs.ReceiptHistoryFieldShare, Kind: structs.ReceiptHistoryKindChanged, Label: "Tip", UserId: alice, From: typedValue(structs.ReceiptHistoryValueCurrency, "1"), To: typedValue(structs.ReceiptHistoryValueCurrency, "2")},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := models.Receipt{ReceiptItems: test.before}
			after := models.Receipt{ReceiptItems: test.after}
			assertChanges(t, diffReceiptSnapshots(&before, after, emptyDiffContext()), test.want)
		})
	}
}

func TestDiffReceiptSnapshotsCommentsAndImagesCompareById(t *testing.T) {
	bob := diffUint(3)
	before := models.Receipt{
		Comments: []models.Comment{
			{BaseModel: models.BaseModel{ID: 1}, Comment: "kept", UserId: bob},
			{BaseModel: models.BaseModel{ID: 2}, Comment: "typo", UserId: bob},
			{BaseModel: models.BaseModel{ID: 3}, Comment: "gone", UserId: bob},
		},
		ImageFiles: []models.FileData{{BaseModel: models.BaseModel{ID: 1}, Name: "old.jpg"}},
	}
	after := models.Receipt{
		Comments: []models.Comment{
			{BaseModel: models.BaseModel{ID: 1}, Comment: "kept", UserId: bob},
			{BaseModel: models.BaseModel{ID: 2}, Comment: "fixed", UserId: bob, Replies: []models.Comment{
				{BaseModel: models.BaseModel{ID: 4}, Comment: "a reply", UserId: diffUint(2)},
			}},
		},
		ImageFiles: []models.FileData{{BaseModel: models.BaseModel{ID: 2}, Name: "new.jpg"}},
	}

	assertChanges(t, diffReceiptSnapshots(&before, after, emptyDiffContext()), []structs.ReceiptHistoryChange{
		{Field: structs.ReceiptHistoryFieldComment, Kind: structs.ReceiptHistoryKindChanged, UserId: bob, From: textValue("typo"), To: textValue("fixed")},
		{Field: structs.ReceiptHistoryFieldComment, Kind: structs.ReceiptHistoryKindRemoved, UserId: bob, From: textValue("gone")},
		{Field: structs.ReceiptHistoryFieldComment, Kind: structs.ReceiptHistoryKindAdded, UserId: diffUint(2), To: textValue("a reply")},
		textChange(structs.ReceiptHistoryFieldImage, structs.ReceiptHistoryKindRemoved, "old.jpg", ""),
		textChange(structs.ReceiptHistoryFieldImage, structs.ReceiptHistoryKindAdded, "", "new.jpg"),
	})
}

func TestDiffReceiptSnapshotsListsFieldsInAFixedOrder(t *testing.T) {
	before := models.Receipt{Name: "A", Status: models.OPEN, ImageFiles: []models.FileData{{BaseModel: models.BaseModel{ID: 1}, Name: "a.jpg"}}}
	after := models.Receipt{
		Name:         "B",
		Status:       models.RESOLVED,
		Tags:         []models.Tag{{BaseModel: models.BaseModel{ID: 1}, Name: "T"}},
		ReceiptItems: []models.Item{diffItem(1, "I", "1", nil)},
		Comments:     []models.Comment{{BaseModel: models.BaseModel{ID: 1}, Comment: "c"}},
	}

	fields := make([]string, 0)
	for _, change := range diffReceiptSnapshots(&before, after, emptyDiffContext()) {
		fields = append(fields, change.Field)
	}

	want := []string{
		structs.ReceiptHistoryFieldName,
		structs.ReceiptHistoryFieldStatus,
		structs.ReceiptHistoryFieldTags,
		structs.ReceiptHistoryFieldItem,
		structs.ReceiptHistoryFieldComment,
		structs.ReceiptHistoryFieldImage,
	}
	if !reflect.DeepEqual(fields, want) {
		t.Errorf("fields = %v, want %v", fields, want)
	}
}

func TestReceiptHistoryAction(t *testing.T) {
	status := func(to string) structs.ReceiptHistoryChange {
		return structs.ReceiptHistoryChange{Field: structs.ReceiptHistoryFieldStatus, Kind: structs.ReceiptHistoryKindChanged, To: typedValue(structs.ReceiptHistoryValueStatus, to)}
	}
	field := func(name string) structs.ReceiptHistoryChange {
		return structs.ReceiptHistoryChange{Field: name, Kind: structs.ReceiptHistoryKindAdded, To: textValue("x")}
	}
	shareStatus := structs.ReceiptHistoryChange{Field: structs.ReceiptHistoryFieldShare, Kind: structs.ReceiptHistoryKindChanged, To: typedValue(structs.ReceiptHistoryValueStatus, "RESOLVED")}
	shareAmount := structs.ReceiptHistoryChange{Field: structs.ReceiptHistoryFieldShare, Kind: structs.ReceiptHistoryKindChanged, To: typedValue(structs.ReceiptHistoryValueCurrency, "1")}

	tests := []struct {
		name    string
		changes []structs.ReceiptHistoryChange
		want    string
	}{
		{"no visible change", nil, structs.ReceiptHistoryActionUpdated},
		{"resolved", []structs.ReceiptHistoryChange{status("RESOLVED"), field(structs.ReceiptHistoryFieldResolvedDate), shareStatus, field(structs.ReceiptHistoryFieldComment)}, structs.ReceiptHistoryActionResolved},
		{"another status", []structs.ReceiptHistoryChange{status("DECLINED")}, structs.ReceiptHistoryActionStatusChanged},
		{"a status with an image is an update", []structs.ReceiptHistoryChange{status("RESOLVED"), field(structs.ReceiptHistoryFieldImage)}, structs.ReceiptHistoryActionUpdated},
		{"a status with a share amount is an update", []structs.ReceiptHistoryChange{status("RESOLVED"), shareAmount}, structs.ReceiptHistoryActionUpdated},
		{"a share status alone is an update", []structs.ReceiptHistoryChange{shareStatus}, structs.ReceiptHistoryActionUpdated},
		{"a resolved date alone is an update", []structs.ReceiptHistoryChange{field(structs.ReceiptHistoryFieldResolvedDate)}, structs.ReceiptHistoryActionUpdated},
		{"comments", []structs.ReceiptHistoryChange{field(structs.ReceiptHistoryFieldComment)}, structs.ReceiptHistoryActionImagesCommentsUpdated},
		{"comments and images", []structs.ReceiptHistoryChange{field(structs.ReceiptHistoryFieldComment), field(structs.ReceiptHistoryFieldImage)}, structs.ReceiptHistoryActionImagesCommentsUpdated},
		{"anything else", []structs.ReceiptHistoryChange{field(structs.ReceiptHistoryFieldName), field(structs.ReceiptHistoryFieldComment)}, structs.ReceiptHistoryActionUpdated},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := receiptHistoryAction(test.changes); got != test.want {
				t.Errorf("action = %s, want %s", got, test.want)
			}
		})
	}
}
