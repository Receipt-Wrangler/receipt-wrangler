package services

import (
	"encoding/json"
	"fmt"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
	"time"
)

// The receipt copies are minimal JSON objects: the upcast only needs a copy to
// parse and to name the right receipt.
func receiptJson(id uint, name string) string {
	return fmt.Sprintf(`{"id":%d,"name":%q}`, id, name)
}

// updateDescription builds a stored RECEIPT_UPDATED description; version 0
// leaves the key out, as rows written before the marker did.
func updateDescription(before string, after string, version int) string {
	bytes, _ := json.Marshal(receiptUpdateDescription{Before: before, After: after, Version: version})
	return string(bytes)
}

type historyTaskOptions struct {
	taskType    models.SystemTaskType
	status      models.SystemTaskStatus
	receiptId   uint
	description string
	// Leaves receipt_id empty and names the receipt only as the associated
	// entity, as rows from before that column did.
	associatedEntityOnly bool
}

func createHistoryTask(t *testing.T, options historyTaskOptions) models.SystemTask {
	t.Helper()
	status := options.status
	if status == "" {
		status = models.SYSTEM_TASK_SUCCEEDED
	}
	endedAt := time.Now()
	task := models.SystemTask{
		Type:                 options.taskType,
		Status:               status,
		AssociatedEntityType: models.RECEIPT,
		AssociatedEntityId:   options.receiptId,
		StartedAt:            time.Now(),
		EndedAt:              &endedAt,
		ResultDescription:    options.description,
	}
	if !options.associatedEntityOnly {
		receiptId := options.receiptId
		task.ReceiptId = &receiptId
	}
	if err := repositories.GetDB().Create(&task).Error; err != nil {
		t.Fatalf("creating system task: %v", err)
	}
	return task
}

func upcast(t *testing.T, tasks ...models.SystemTask) []models.SystemTask {
	t.Helper()
	if err := NewSystemTaskService(nil).UpcastReceiptUpdateDescriptions(tasks); err != nil {
		t.Fatalf("upcasting: %v", err)
	}
	return tasks
}

func readDescription(t *testing.T, task models.SystemTask) receiptUpdateDescription {
	t.Helper()
	var description receiptUpdateDescription
	if err := json.Unmarshal([]byte(task.ResultDescription), &description); err != nil {
		t.Fatalf("reading description %q: %v", task.ResultDescription, err)
	}
	return description
}

func assertRebuiltFrom(t *testing.T, task models.SystemTask, before string, source models.SystemTask) {
	t.Helper()
	description := readDescription(t, task)
	if description.Before != before {
		t.Errorf("before = %s, want %s", description.Before, before)
	}
	if description.Version != 1 {
		t.Errorf("version = %d, want 1", description.Version)
	}
	if description.BeforeSource == nil ||
		description.BeforeSource.SystemTaskId != source.ID ||
		description.BeforeSource.Type != source.Type {
		t.Errorf("beforeSource = %+v, want task %d (%s)", description.BeforeSource, source.ID, source.Type)
	}
}

func TestUpcastReceiptUpdateDescriptions(t *testing.T) {
	const shallow = `{"id":7,"name":"incomplete"}`

	t.Run("an update compares against the previous update's after", func(t *testing.T) {
		defer repositories.TruncateTestDb()

		created := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: receiptJson(7, "Created")})
		first := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 7, description: updateDescription(shallow, receiptJson(7, "One"), 0)})
		second := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 7, description: updateDescription(shallow, receiptJson(7, "Two"), 0)})

		tasks := upcast(t, second, first)

		assertRebuiltFrom(t, tasks[0], receiptJson(7, "One"), first)
		// The first update has no earlier update, so it uses the creation copy,
		// and a later update is never used.
		assertRebuiltFrom(t, tasks[1], receiptJson(7, "Created"), created)
		if after := readDescription(t, tasks[0]).After; after != receiptJson(7, "Two") {
			t.Errorf("after = %s, want the stored after", after)
		}
	})

	t.Run("skips failed and unusable copies for an older usable one", func(t *testing.T) {
		defer repositories.TruncateTestDb()

		created := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: receiptJson(7, "Created")})
		createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, status: models.SYSTEM_TASK_FAILED, receiptId: 7, description: "record not found"})
		createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: "not json"})
		createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: receiptJson(8, "Someone else's copy")})
		createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 8, description: updateDescription(shallow, receiptJson(8, "Other receipt"), 2)})
		update := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 7, description: updateDescription(shallow, receiptJson(7, "One"), 0)})

		assertRebuiltFrom(t, upcast(t, update)[0], receiptJson(7, "Created"), created)
	})

	t.Run("matches rows that name the receipt only as their associated entity", func(t *testing.T) {
		defer repositories.TruncateTestDb()

		created := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: receiptJson(7, "Created"), associatedEntityOnly: true})
		update := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 7, description: updateDescription(shallow, receiptJson(7, "One"), 0), associatedEntityOnly: true})

		assertRebuiltFrom(t, upcast(t, update)[0], receiptJson(7, "Created"), created)
	})

	t.Run("leaves version 2 rows, error text and rows with no earlier copy as stored", func(t *testing.T) {
		defer repositories.TruncateTestDb()

		createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: receiptJson(7, "Created")})
		versionTwo := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 7, description: updateDescription(receiptJson(7, "Created"), receiptJson(7, "One"), 2)})
		failed := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, status: models.SYSTEM_TASK_FAILED, receiptId: 7, description: "record not found"})
		noHistory := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 9, description: updateDescription(`{"id":9}`, receiptJson(9, "One"), 0)})

		tasks := upcast(t, versionTwo, failed, noHistory)

		for i, original := range []models.SystemTask{versionTwo, failed, noHistory} {
			if tasks[i].ResultDescription != original.ResultDescription {
				t.Errorf("task %d changed: %s", original.ID, tasks[i].ResultDescription)
			}
		}
	})

	t.Run("never writes the rebuilt description back", func(t *testing.T) {
		defer repositories.TruncateTestDb()

		createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPLOADED, receiptId: 7, description: receiptJson(7, "Created")})
		update := createHistoryTask(t, historyTaskOptions{taskType: models.RECEIPT_UPDATED, receiptId: 7, description: updateDescription(shallow, receiptJson(7, "One"), 0)})
		upcast(t, update)

		var stored models.SystemTask
		repositories.GetDB().First(&stored, update.ID)
		if stored.ResultDescription != update.ResultDescription {
			t.Errorf("stored description changed: %s", stored.ResultDescription)
		}
	})
}
