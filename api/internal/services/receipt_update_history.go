package services

import (
	"encoding/json"
	"receipt-wrangler/api/internal/models"
	"sort"
	"time"
)

// receiptUpdateDescription is a RECEIPT_UPDATED system task's description:
// Before and After are each a receipt serialized as a JSON string. Version is
// absent (0 here) on rows written before it existed, which mean version 1; see
// repositories.ReceiptUpdateDescriptionVersion.
type receiptUpdateDescription struct {
	Before       string                     `json:"before"`
	After        string                     `json:"after"`
	Version      int                        `json:"version,omitempty"`
	BeforeSource *receiptUpdateBeforeSource `json:"beforeSource,omitempty"`
}

// receiptUpdateBeforeSource names the recorded copy a rebuilt "before" came
// from, so the desktop can say what it is comparing against.
type receiptUpdateBeforeSource struct {
	Type         models.SystemTaskType `json:"type"`
	SystemTaskId uint                  `json:"systemTaskId"`
	RecordedAt   time.Time             `json:"recordedAt"`
}

// UpcastReceiptUpdateDescriptions rebuilds the comparison of every version 1
// RECEIPT_UPDATED row in tasks, in place, for the response only — nothing is
// written back.
//
// A version 1 row's own "before" was loaded one level deep, so diffing it
// reports changes that never happened. Its "after", like every later row's
// and like the copy stored when the receipt was created (RECEIPT_UPLOADED),
// is complete, and the receipt just before an update is the receipt just
// after the previous recorded change. So "before" is replaced with the
// nearest earlier complete copy of the same receipt, and BeforeSource says
// which one. The row keeps version 1: a change that writes no update row (a
// bulk status change, say) now falls inside this comparison, which a version
// 2 row's own "before" would not have.
//
// A row with no earlier usable copy is left as stored.
func (service SystemTaskService) UpcastReceiptUpdateDescriptions(tasks []models.SystemTask) error {
	type versionOneRow struct {
		index       int
		receiptId   uint
		description receiptUpdateDescription
	}

	rows := make([]versionOneRow, 0)
	receiptIds := make([]uint, 0)
	var highestId uint
	for i, task := range tasks {
		if task.Type != models.RECEIPT_UPDATED || task.Status != models.SYSTEM_TASK_SUCCEEDED {
			continue
		}
		description, ok := parseReceiptUpdateDescription(task.ResultDescription)
		if !ok || description.Version != 0 {
			continue
		}
		receiptId := systemTaskReceiptId(task)
		if receiptId == 0 {
			continue
		}

		rows = append(rows, versionOneRow{index: i, receiptId: receiptId, description: description})
		receiptIds = append(receiptIds, receiptId)
		if task.ID > highestId {
			highestId = task.ID
		}
	}
	if len(rows) == 0 {
		return nil
	}

	candidatesByReceipt, err := service.getReceiptSnapshotCandidates(receiptIds, highestId)
	if err != nil {
		return err
	}

	for _, row := range rows {
		task := &tasks[row.index]
		for _, candidate := range candidatesByReceipt[row.receiptId] {
			if candidate.ID >= task.ID {
				continue
			}
			snapshot, ok := receiptSnapshotFromTask(candidate, row.receiptId)
			if !ok {
				continue
			}

			recordedAt := candidate.StartedAt
			if candidate.EndedAt != nil {
				recordedAt = *candidate.EndedAt
			}
			rebuilt, err := json.Marshal(receiptUpdateDescription{
				Before:  snapshot,
				After:   row.description.After,
				Version: 1,
				BeforeSource: &receiptUpdateBeforeSource{
					Type:         candidate.Type,
					SystemTaskId: candidate.ID,
					RecordedAt:   recordedAt,
				},
			})
			if err != nil {
				return err
			}
			task.ResultDescription = string(rebuilt)
			break
		}
	}

	return nil
}

// getReceiptSnapshotCandidates loads every successful task that may hold a
// complete copy of one of the receipts, below maxId, grouped by receipt and
// newest first. One query covers the whole page.
func (service SystemTaskService) getReceiptSnapshotCandidates(
	receiptIds []uint,
	maxId uint,
) (map[uint][]models.SystemTask, error) {
	db := service.GetDB()
	var candidates []models.SystemTask

	// Rows that predate the receipt_id column carry the receipt only as their
	// associated entity; Quick Scan and email uploads carry it only in
	// receipt_id (their associated entity is the processing settings).
	receiptMatch := db.Where("receipt_id IN ?", receiptIds).
		Or("associated_entity_type = ? AND associated_entity_id IN ?", models.RECEIPT, receiptIds)

	err := db.Model(&models.SystemTask{}).
		Select("id", "type", "status", "receipt_id", "associated_entity_type", "associated_entity_id",
			"result_description", "started_at", "ended_at").
		Where("type IN ?", []models.SystemTaskType{models.RECEIPT_UPDATED, models.RECEIPT_UPLOADED}).
		Where("status = ?", models.SYSTEM_TASK_SUCCEEDED).
		Where("id < ?", maxId).
		Where(receiptMatch).
		Find(&candidates).Error
	if err != nil {
		return nil, err
	}

	byReceipt := make(map[uint][]models.SystemTask)
	for _, candidate := range candidates {
		receiptId := systemTaskReceiptId(candidate)
		byReceipt[receiptId] = append(byReceipt[receiptId], candidate)
	}
	for receiptId := range byReceipt {
		sort.Slice(byReceipt[receiptId], func(i, j int) bool {
			return byReceipt[receiptId][i].ID > byReceipt[receiptId][j].ID
		})
	}

	return byReceipt, nil
}

// receiptSnapshotFromTask returns the complete receipt copy a task recorded:
// an update's "after", or the receipt an upload created. It is used only if it
// parses as that same receipt.
func receiptSnapshotFromTask(task models.SystemTask, receiptId uint) (string, bool) {
	snapshot := task.ResultDescription
	if task.Type == models.RECEIPT_UPDATED {
		description, ok := parseReceiptUpdateDescription(task.ResultDescription)
		if !ok {
			return "", false
		}
		snapshot = description.After
	}

	var receipt struct {
		ID uint `json:"id"`
	}
	if err := json.Unmarshal([]byte(snapshot), &receipt); err != nil || receipt.ID != receiptId {
		return "", false
	}

	return snapshot, true
}

func parseReceiptUpdateDescription(value string) (receiptUpdateDescription, bool) {
	var description receiptUpdateDescription
	if err := json.Unmarshal([]byte(value), &description); err != nil {
		return receiptUpdateDescription{}, false
	}

	return description, description.Before != "" && description.After != ""
}

func systemTaskReceiptId(task models.SystemTask) uint {
	if task.ReceiptId != nil && *task.ReceiptId != 0 {
		return *task.ReceiptId
	}
	if task.AssociatedEntityType == models.RECEIPT {
		return task.AssociatedEntityId
	}

	return 0
}
