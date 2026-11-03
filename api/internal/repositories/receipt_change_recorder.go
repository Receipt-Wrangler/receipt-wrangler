package repositories

import (
	"encoding/json"
	"fmt"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"slices"
	"time"
)

// ReceiptChangeRecorder records a change to one or more receipts as RECEIPT_UPDATED
// system tasks, one per receipt, in the version 2 format
// (ReceiptUpdateDescriptionVersion): {"before", "after", "version"}, each side the
// receipt fully loaded (GetFullyLoadedReceiptsByIds) and serialized.
//
// It is the one place that format is written. UpdateReceipt uses it, and so does
// every other change that used to write no row at all — comments, images and the
// bulk status update — so the receipt history and the System Tasks diff see them.
//
//	recorder := NewReceiptChangeRecorder(userId, receiptId)
//	CaptureBefore() → make the change → RecordSuccess(), or RecordFailure(err)
//
// Record wraps that sequence for a change that has no use for the reloaded
// receipts.
type ReceiptChangeRecorder struct {
	ranByUserId uint
	receiptIds  []uint
	startedAt   time.Time
	// captured holds each receipt that existed when CaptureBefore ran, by id.
	// It is nil until CaptureBefore succeeds.
	captured map[uint]capturedReceipt
}

type capturedReceipt struct {
	groupId uint
	before  string
}

func NewReceiptChangeRecorder(ranByUserId uint, receiptIds ...uint) *ReceiptChangeRecorder {
	uniqueIds := make([]uint, 0, len(receiptIds))
	for _, id := range receiptIds {
		if !slices.Contains(uniqueIds, id) {
			uniqueIds = append(uniqueIds, id)
		}
	}

	return &ReceiptChangeRecorder{
		ranByUserId: ranByUserId,
		receiptIds:  uniqueIds,
		startedAt:   time.Now(),
	}
}

// CaptureBefore loads and serializes every receipt as it is before the change.
// A receipt that does not exist is skipped: RecordSuccess records nothing for it.
func (recorder *ReceiptChangeRecorder) CaptureBefore() error {
	receipts, err := NewReceiptRepository(nil).GetFullyLoadedReceiptsByIds(recorder.receiptIds)
	if err != nil {
		return err
	}

	captured := make(map[uint]capturedReceipt, len(receipts))
	for i := range receipts {
		before, err := receipts[i].ToString()
		if err != nil {
			return err
		}
		captured[receipts[i].ID] = capturedReceipt{groupId: receipts[i].GroupId, before: before}
	}

	recorder.captured = captured
	return nil
}

// RecordSuccess loads every captured receipt again, now that the change is
// committed, and records its before and after. It returns the reloaded receipts by
// id. It records nothing when it fails — a receipt that could not be reloaded,
// say — and leaves it to the caller whether that fails the change.
func (recorder *ReceiptChangeRecorder) RecordSuccess() (map[uint]models.Receipt, error) {
	if recorder.captured == nil {
		return nil, fmt.Errorf("the receipts were not captured before the change")
	}

	ids := recorder.capturedIds()
	receipts, err := NewReceiptRepository(nil).GetFullyLoadedReceiptsByIds(ids)
	if err != nil {
		return nil, err
	}

	after := make(map[uint]models.Receipt, len(receipts))
	for _, receipt := range receipts {
		after[receipt.ID] = receipt
	}

	// Every row is built before any is written, so a receipt that cannot be
	// reloaded or serialized leaves no partial record behind.
	systemTasks := make([]commands.UpsertSystemTaskCommand, 0, len(ids))
	for _, id := range ids {
		receipt, ok := after[id]
		if !ok {
			return nil, fmt.Errorf("receipt %d could not be reloaded", id)
		}

		afterString, err := receipt.ToString()
		if err != nil {
			return nil, err
		}

		description, err := json.Marshal(map[string]interface{}{
			"before":  recorder.captured[id].before,
			"after":   afterString,
			"version": ReceiptUpdateDescriptionVersion,
		})
		if err != nil {
			return nil, err
		}

		endedAt := time.Now()
		systemTask := recorder.systemTask(id)
		systemTask.EndedAt = &endedAt
		systemTask.ResultDescription = string(description)
		systemTasks = append(systemTasks, systemTask)
	}

	systemTaskRepository := NewSystemTaskRepository(nil)
	for _, systemTask := range systemTasks {
		if _, err := systemTaskRepository.CreateSystemTask(systemTask); err != nil {
			return nil, err
		}
	}

	return after, nil
}

// RecordFailure records a FAILED row, carrying err's text, for every receipt the
// recorder was made for. Once CaptureBefore has run, a row also names the
// receipt's group, so it shows in that group's activity.
func (recorder *ReceiptChangeRecorder) RecordFailure(err error) {
	for _, id := range recorder.receiptIds {
		createFailedUpdateSystemTask(recorder.systemTask(id), err)
	}
}

// Record captures the receipts, runs change, and records the result. A change that
// fails, or cannot be captured beforehand, is recorded as failed and its error
// returned. Once change has succeeded, a failure to record it is only logged: the
// change is committed, and failing the request would invite a retry that repeats
// it (a second comment, a second image).
func (recorder *ReceiptChangeRecorder) Record(change func() error) error {
	if err := recorder.CaptureBefore(); err != nil {
		recorder.RecordFailure(err)
		return err
	}

	if err := change(); err != nil {
		recorder.RecordFailure(err)
		return err
	}

	if _, err := recorder.RecordSuccess(); err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, fmt.Sprintf("Failed to record a change to receipts %v: %s", recorder.receiptIds, err.Error()))
	}

	return nil
}

// capturedIds lists the captured receipts in the order they were asked for.
func (recorder *ReceiptChangeRecorder) capturedIds() []uint {
	ids := make([]uint, 0, len(recorder.captured))
	for _, id := range recorder.receiptIds {
		if _, ok := recorder.captured[id]; ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func (recorder *ReceiptChangeRecorder) systemTask(receiptId uint) commands.UpsertSystemTaskCommand {
	ranByUserId := recorder.ranByUserId
	systemTask := commands.UpsertSystemTaskCommand{
		Type:                 models.RECEIPT_UPDATED,
		AssociatedEntityType: models.RECEIPT,
		AssociatedEntityId:   receiptId,
		StartedAt:            recorder.startedAt,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		RanByUserId:          &ranByUserId,
	}

	if captured, ok := recorder.captured[receiptId]; ok {
		id := receiptId
		groupId := captured.groupId
		systemTask.ReceiptId = &id
		systemTask.GroupId = &groupId
	}

	return systemTask
}
