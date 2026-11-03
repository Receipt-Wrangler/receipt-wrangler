package repositories

import (
	"encoding/json"
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// createFullyPopulatedReceipt creates a receipt carrying every association the
// full loader preloads: an item with categories, tags and a linked item, a
// custom field value with its definition, a comment and an image row.
func createFullyPopulatedReceipt(t *testing.T, name string, customField models.CustomField) models.Receipt {
	t.Helper()
	poNumber := "PO-" + name
	command := commands.UpsertReceiptCommand{
		Name:         name,
		Amount:       decimal.NewFromFloat(50.00),
		Date:         time.Now(),
		PaidByUserID: 1,
		Status:       models.OPEN,
		GroupId:      1,
		Categories:   []commands.UpsertCategoryCommand{{Id: uintPtr(1)}},
		Tags:         []commands.UpsertTagCommand{{Id: uintPtr(1)}},
		Items: []commands.UpsertItemCommand{{
			Name:            "Main Item",
			Amount:          decimal.NewFromFloat(40.00),
			ChargedToUserId: uintPtr(2),
			Status:          models.ITEM_OPEN,
			Categories:      []commands.UpsertCategoryCommand{{Id: uintPtr(2)}},
			Tags:            []commands.UpsertTagCommand{{Id: uintPtr(2)}},
			LinkedItems: []commands.UpsertItemCommand{{
				Name:            "Linked Item",
				Amount:          decimal.NewFromFloat(10.00),
				ChargedToUserId: uintPtr(3),
				Status:          models.ITEM_OPEN,
				Categories:      []commands.UpsertCategoryCommand{{Id: uintPtr(3)}},
			}},
		}},
		CustomFields: []commands.UpsertCustomFieldValueCommand{{CustomFieldId: customField.ID, StringValue: &poNumber}},
	}

	receipt, err := NewReceiptRepository(nil).CreateReceipt(command, 1, false)
	if err != nil {
		t.Fatalf("create receipt: %v", err)
	}

	db := GetDB()
	author := uint(1)
	db.Create(&models.Comment{Comment: "comment on " + name, ReceiptId: receipt.ID, UserId: &author})
	db.Create(&models.FileData{Name: name + ".jpg", ReceiptId: receipt.ID, FileType: "image/jpeg"})
	return receipt
}

func TestGetFullyLoadedReceiptsByIdsMatchesSingleLoader(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	customField := models.CustomField{Name: "PO Number", Type: models.TEXT}
	GetDB().Create(&customField)
	first := createFullyPopulatedReceipt(t, "First", customField)
	second := createFullyPopulatedReceipt(t, "Second", customField)

	repository := NewReceiptRepository(nil)
	batch, err := repository.GetFullyLoadedReceiptsByIds([]uint{second.ID, 999, first.ID})
	if err != nil {
		t.Fatalf("batch load: %v", err)
	}

	// Ordered by id, and an id with no receipt is simply absent.
	if len(batch) != 2 || batch[0].ID != first.ID || batch[1].ID != second.ID {
		t.Fatalf("batch = %d receipts, want receipts %d and %d in id order", len(batch), first.ID, second.ID)
	}

	for _, batched := range batch {
		single, err := repository.GetFullyLoadedReceiptById(utils.UintToString(batched.ID))
		if err != nil {
			t.Fatalf("single load: %v", err)
		}

		batchedJson, _ := batched.ToString()
		singleJson, _ := single.ToString()
		if batchedJson != singleJson {
			t.Errorf("receipt %d loads differently:\nbatch  %s\nsingle %s", batched.ID, batchedJson, singleJson)
		}

		// The depth both loaders promise: linked items filtered out of the top
		// level, nested categories/tags, and custom field definitions.
		if len(batched.ReceiptItems) != 1 || len(batched.ReceiptItems[0].LinkedItems) != 1 ||
			len(batched.ReceiptItems[0].LinkedItems[0].Categories) != 1 ||
			len(batched.ReceiptItems[0].Categories) != 1 || len(batched.ReceiptItems[0].Tags) != 1 ||
			len(batched.CustomFields) != 1 || batched.CustomFields[0].CustomField.Name != "PO Number" ||
			len(batched.Comments) != 1 || len(batched.ImageFiles) != 1 {
			t.Errorf("receipt %d was not fully loaded: %s", batched.ID, batchedJson)
		}
	}

	empty, err := repository.GetFullyLoadedReceiptsByIds(nil)
	if err != nil || empty == nil || len(empty) != 0 {
		t.Errorf("no ids = %#v, %v; want an empty slice", empty, err)
	}
}

func receiptUpdatedTasks(t *testing.T) []models.SystemTask {
	t.Helper()
	var tasks []models.SystemTask
	if err := GetDB().Where("type = ?", models.RECEIPT_UPDATED).Order("id").Find(&tasks).Error; err != nil {
		t.Fatalf("load system tasks: %v", err)
	}
	return tasks
}

func readRecordedDescription(t *testing.T, task models.SystemTask) (models.Receipt, models.Receipt, int) {
	t.Helper()
	var description struct {
		Before  string `json:"before"`
		After   string `json:"after"`
		Version int    `json:"version"`
	}
	if err := json.Unmarshal([]byte(task.ResultDescription), &description); err != nil {
		t.Fatalf("read description %q: %v", task.ResultDescription, err)
	}

	var before, after models.Receipt
	if err := json.Unmarshal([]byte(description.Before), &before); err != nil {
		t.Fatalf("read before: %v", err)
	}
	if err := json.Unmarshal([]byte(description.After), &after); err != nil {
		t.Fatalf("read after: %v", err)
	}
	return before, after, description.Version
}

func createPlainReceipt(t *testing.T, name string, groupId uint) models.Receipt {
	t.Helper()
	receipt := models.Receipt{Name: name, Amount: decimal.NewFromFloat(5), Date: time.Now(), GroupId: groupId, PaidByUserID: 1, Status: models.OPEN}
	if err := GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("create receipt: %v", err)
	}
	return receipt
}

func TestReceiptChangeRecorderRecordsAVersionTwoRowPerReceipt(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	first := createPlainReceipt(t, "First", 1)
	second := createPlainReceipt(t, "Second", 2)

	// A duplicate id is recorded once, and an id with no receipt not at all.
	recorder := NewReceiptChangeRecorder(3, first.ID, second.ID, first.ID, 999)
	if err := recorder.CaptureBefore(); err != nil {
		t.Fatalf("capture: %v", err)
	}

	author := uint(3)
	db := GetDB()
	db.Create(&models.Comment{Comment: "bulk note", ReceiptId: first.ID, UserId: &author})
	db.Create(&models.Comment{Comment: "bulk note", ReceiptId: second.ID, UserId: &author})

	after, err := recorder.RecordSuccess()
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if len(after) != 2 || len(after[first.ID].Comments) != 1 || len(after[second.ID].Comments) != 1 {
		t.Errorf("RecordSuccess returned %d receipts, want both reloaded with their comment", len(after))
	}

	tasks := receiptUpdatedTasks(t)
	if len(tasks) != 2 {
		t.Fatalf("recorded %d rows, want 2", len(tasks))
	}

	for i, receipt := range []models.Receipt{first, second} {
		task := tasks[i]
		if task.Status != models.SYSTEM_TASK_SUCCEEDED ||
			task.AssociatedEntityType != models.RECEIPT || task.AssociatedEntityId != receipt.ID ||
			task.ReceiptId == nil || *task.ReceiptId != receipt.ID ||
			task.GroupId == nil || *task.GroupId != receipt.GroupId ||
			task.RanByUserId == nil || *task.RanByUserId != 3 ||
			task.EndedAt == nil {
			t.Errorf("row %d = %+v, want a succeeded row for receipt %d in group %d run by 3", i, task, receipt.ID, receipt.GroupId)
		}

		before, afterSnapshot, version := readRecordedDescription(t, task)
		if version != ReceiptUpdateDescriptionVersion {
			t.Errorf("version = %d, want %d", version, ReceiptUpdateDescriptionVersion)
		}
		if len(before.Comments) != 0 || len(afterSnapshot.Comments) != 1 || afterSnapshot.Comments[0].Comment != "bulk note" {
			t.Errorf("receipt %d: before has %d comments, after %d; want 0 and the bulk note", receipt.ID, len(before.Comments), len(afterSnapshot.Comments))
		}
	}
}

func TestReceiptChangeRecorderRecordFailure(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	receipt := createPlainReceipt(t, "Receipt", 1)
	changeErr := errors.New("the change failed")

	t.Run("before the capture, a row names the receipt only as its entity", func(t *testing.T) {
		NewReceiptChangeRecorder(2, receipt.ID).RecordFailure(changeErr)

		tasks := receiptUpdatedTasks(t)
		if len(tasks) != 1 {
			t.Fatalf("recorded %d rows, want 1", len(tasks))
		}
		task := tasks[0]
		if task.Status != models.SYSTEM_TASK_FAILED || task.ResultDescription != changeErr.Error() ||
			task.AssociatedEntityId != receipt.ID || task.ReceiptId != nil || task.GroupId != nil {
			t.Errorf("row = %+v, want a failed row naming receipt %d only as its entity", task, receipt.ID)
		}
		GetDB().Where("1 = 1").Delete(&models.SystemTask{})
	})

	t.Run("after the capture, a row names the receipt and its group", func(t *testing.T) {
		recorder := NewReceiptChangeRecorder(2, receipt.ID)
		if err := recorder.CaptureBefore(); err != nil {
			t.Fatalf("capture: %v", err)
		}
		recorder.RecordFailure(changeErr)

		tasks := receiptUpdatedTasks(t)
		if len(tasks) != 1 {
			t.Fatalf("recorded %d rows, want 1", len(tasks))
		}
		task := tasks[0]
		if task.Status != models.SYSTEM_TASK_FAILED || task.ResultDescription != changeErr.Error() ||
			task.ReceiptId == nil || *task.ReceiptId != receipt.ID || task.GroupId == nil || *task.GroupId != 1 ||
			task.RanByUserId == nil || *task.RanByUserId != 2 {
			t.Errorf("row = %+v, want a failed row for receipt %d in group 1 run by 2", task, receipt.ID)
		}
		GetDB().Where("1 = 1").Delete(&models.SystemTask{})
	})
}

func TestReceiptChangeRecorderRecord(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	t.Run("a failed change is recorded as failed and its error returned", func(t *testing.T) {
		defer GetDB().Where("1 = 1").Delete(&models.SystemTask{})
		receipt := createPlainReceipt(t, "Receipt", 1)
		changeErr := errors.New("the change failed")

		err := NewReceiptChangeRecorder(1, receipt.ID).Record(func() error { return changeErr })
		if !errors.Is(err, changeErr) {
			t.Errorf("Record = %v, want the change's error", err)
		}

		tasks := receiptUpdatedTasks(t)
		if len(tasks) != 1 || tasks[0].Status != models.SYSTEM_TASK_FAILED || tasks[0].ResultDescription != changeErr.Error() {
			t.Errorf("rows = %+v, want one failed row", tasks)
		}
	})

	t.Run("a change that cannot be recorded afterwards still succeeds, and records nothing", func(t *testing.T) {
		defer GetDB().Where("1 = 1").Delete(&models.SystemTask{})
		receipt := createPlainReceipt(t, "Doomed", 1)

		// Deleting the receipt makes the "after" load find nothing.
		err := NewReceiptChangeRecorder(1, receipt.ID).Record(func() error {
			return GetDB().Delete(&models.Receipt{}, receipt.ID).Error
		})
		if err != nil {
			t.Errorf("Record = %v, want nil once the change committed", err)
		}

		if tasks := receiptUpdatedTasks(t); len(tasks) != 0 {
			t.Errorf("rows = %+v, want none: the change committed, so it is not a failure", tasks)
		}
	})

	t.Run("a successful change is recorded", func(t *testing.T) {
		defer GetDB().Where("1 = 1").Delete(&models.SystemTask{})
		receipt := createPlainReceipt(t, "Before", 1)

		err := NewReceiptChangeRecorder(1, receipt.ID).Record(func() error {
			return GetDB().Model(&models.Receipt{}).Where("id = ?", receipt.ID).Update("name", "After").Error
		})
		if err != nil {
			t.Fatalf("Record = %v", err)
		}

		tasks := receiptUpdatedTasks(t)
		if len(tasks) != 1 || tasks[0].Status != models.SYSTEM_TASK_SUCCEEDED {
			t.Fatalf("rows = %+v, want one succeeded row", tasks)
		}
		before, after, _ := readRecordedDescription(t, tasks[0])
		if before.Name != "Before" || after.Name != "After" {
			t.Errorf("recorded %q -> %q, want Before -> After", before.Name, after.Name)
		}
	})

	t.Run("recording success without a capture is an error", func(t *testing.T) {
		if _, err := NewReceiptChangeRecorder(1, 1).RecordSuccess(); err == nil {
			t.Error("RecordSuccess without CaptureBefore = nil, want an error")
		}
	})
}
