package services

import (
	"encoding/json"
	"errors"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strconv"

	"gorm.io/gorm"
)

// restrictedGroupName stands in for the name of a group the caller may not see:
// one they are not a member of, without app.groups.read.
const restrictedGroupName = "(Restricted)"

// ReceiptHistoryService builds a receipt's change history from its system tasks:
// the RECEIPT_UPLOADED row that created it and every RECEIPT_UPDATED row, failed
// ones included. The server computes the diff and the client only formats, so
// the redaction rules stay in one place and every client shows the same thing.
//
// Its callers authorize the read first: group.receipts.read and
// group.receipts.history.read in the receipt's group as it is now, plus paid-by
// visibility on the receipt.
type ReceiptHistoryService struct {
	BaseService
}

func NewReceiptHistoryService(tx *gorm.DB) ReceiptHistoryService {
	return ReceiptHistoryService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
}

// historyRecord is one system task read as a history entry, before it is diffed.
type historyRecord struct {
	task     models.SystemTask
	creation bool
	failed   bool
	actor    structs.ReceiptHistoryActor
	// before and after are the parsed snapshots. Both are nil on a failed row or
	// one whose description cannot be read; before is nil on the creation.
	before       *models.Receipt
	after        *models.Receipt
	beforeState  string
	beforeSource *receiptUpdateBeforeSource
}

// GetReceiptHistory returns the receipt's history as viewerId may see it, newest
// entry first.
func (service ReceiptHistoryService) GetReceiptHistory(viewerId uint, receiptId uint) (structs.ReceiptHistory, error) {
	records, err := service.loadRecords(viewerId, receiptId, 0)
	if err != nil {
		return structs.ReceiptHistory{}, err
	}

	context, err := service.diffContext(viewerId, records)
	if err != nil {
		return structs.ReceiptHistory{}, err
	}

	history := structs.NewReceiptHistory()
	contributors := make(map[string]bool)
	for _, record := range records {
		entry := record.toEntry(context)
		history.Entries = append(history.Entries, entry)
		contributors[receiptHistoryActorKey(entry.Actor)] = true

		if record.failed {
			continue
		}
		history.SaveCount++
		if !record.creation {
			history.ChangeCount += len(entry.Changes)
		}
	}
	history.ContributorCount = len(contributors)

	return history, nil
}

// GetReceiptHistorySnapshots returns the receipt before and after one entry,
// redacted exactly as GetReceiptHistory redacts it, and never the stored
// description. An entry that does not exist, is not this receipt's, is hidden
// from the caller or has no snapshots is ErrReceiptAccessDenied, so no answer
// tells those apart.
func (service ReceiptHistoryService) GetReceiptHistorySnapshots(
	viewerId uint,
	receiptId uint,
	systemTaskId uint,
) (structs.ReceiptHistorySnapshots, error) {
	if systemTaskId == 0 {
		return structs.ReceiptHistorySnapshots{}, ErrReceiptAccessDenied
	}

	records, err := service.loadRecords(viewerId, receiptId, systemTaskId)
	if err != nil {
		return structs.ReceiptHistorySnapshots{}, err
	}
	if len(records) == 0 || records[0].after == nil {
		return structs.ReceiptHistorySnapshots{}, ErrReceiptAccessDenied
	}
	record := records[0]

	snapshots := structs.ReceiptHistorySnapshots{}
	if snapshots.After, err = record.after.ToString(); err != nil {
		return structs.ReceiptHistorySnapshots{}, err
	}
	if record.before != nil {
		if snapshots.Before, err = record.before.ToString(); err != nil {
			return structs.ReceiptHistorySnapshots{}, err
		}
	}

	return snapshots, nil
}

// loadRecords reads the receipt's tasks into records, newest first, with every
// one the caller may not see dropped and every snapshot redacted. A non-zero
// onlyTaskId keeps that one task alone (it still reads every task, since a
// version 1 row is rebuilt from an earlier one).
func (service ReceiptHistoryService) loadRecords(viewerId uint, receiptId uint, onlyTaskId uint) ([]historyRecord, error) {
	receipt, err := repositories.NewReceiptRepository(service.TX).GetReceiptForAuthorization(utils.UintToString(receiptId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReceiptAccessDenied
		}
		return nil, err
	}

	tasks, err := service.receiptTasks(receiptId)
	if err != nil {
		return nil, err
	}
	tasks = tasksSinceCreation(tasks, receiptId)

	selected := tasks
	if onlyTaskId != 0 {
		selected = make([]models.SystemTask, 0, 1)
		for _, task := range tasks {
			if task.ID == onlyTaskId {
				selected = append(selected, task)
			}
		}
	}

	parents, err := service.parentTasks(selected)
	if err != nil {
		return nil, err
	}

	records := make([]historyRecord, 0, len(selected))
	for _, task := range selected {
		record := readHistoryRecord(task, tasks, receiptId)
		record.actor = receiptHistoryActor(record, parents)
		records = append(records, record)
	}

	return service.redact(viewerId, receipt.GroupId, records)
}

// receiptTasks loads the receipt's creation and update rows, both statuses,
// newest first.
func (service ReceiptHistoryService) receiptTasks(receiptId uint) ([]models.SystemTask, error) {
	db := service.GetDB()
	tasks := make([]models.SystemTask, 0)

	err := db.Model(&models.SystemTask{}).
		Where("type IN ?", []models.SystemTaskType{models.RECEIPT_UPLOADED, models.RECEIPT_UPDATED}).
		Where(repositories.WhereReceiptTasks(db, []uint{receiptId})).
		Order("id DESC").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// parentTasks loads, in one query, the task each upload ran under: the QUICK_SCAN
// or EMAIL_UPLOAD that says who or what made the receipt.
func (service ReceiptHistoryService) parentTasks(tasks []models.SystemTask) (map[uint]models.SystemTask, error) {
	parentIds := make([]uint, 0)
	for _, task := range tasks {
		if task.AssociatedSystemTaskId != nil && *task.AssociatedSystemTaskId != 0 {
			parentIds = append(parentIds, *task.AssociatedSystemTaskId)
		}
	}

	parents := make(map[uint]models.SystemTask, len(parentIds))
	if len(parentIds) == 0 {
		return parents, nil
	}

	var rows []models.SystemTask
	err := service.GetDB().Model(&models.SystemTask{}).
		Select("id", "type", "ran_by_user_id").
		Where("id IN ?", parentIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		parents[row.ID] = row
	}
	return parents, nil
}

// tasksSinceCreation drops the rows of an earlier receipt that had the same id.
// A deleted receipt keeps its system tasks, and an id can be handed out again
// (SQLite reuses the highest one), so everything older than the newest successful
// creation of THIS receipt belongs to another. tasks must be newest first; with
// no such creation every row is kept.
func tasksSinceCreation(tasks []models.SystemTask, receiptId uint) []models.SystemTask {
	for i, task := range tasks {
		if task.Type != models.RECEIPT_UPLOADED || task.Status != models.SYSTEM_TASK_SUCCEEDED {
			continue
		}
		if _, ok := receiptSnapshotFromTask(task, receiptId); ok {
			return tasks[:i+1]
		}
	}
	return tasks
}

// readHistoryRecord reads one task's snapshots. tasks are all the receipt's
// tasks, newest first, which a version 1 update is rebuilt from.
func readHistoryRecord(task models.SystemTask, tasks []models.SystemTask, receiptId uint) historyRecord {
	record := historyRecord{
		task:     task,
		creation: task.Type == models.RECEIPT_UPLOADED,
		failed:   task.Status == models.SYSTEM_TASK_FAILED,
	}
	if record.failed {
		return record
	}

	if record.creation {
		record.after = parseReceiptSnapshot(task.ResultDescription)
		return record
	}

	description, ok := parseReceiptUpdateDescription(task.ResultDescription)
	if !ok {
		return record
	}

	before := description.Before
	record.beforeState = structs.ReceiptHistoryBeforeStateComplete
	if description.hasIncompleteBefore() {
		record.beforeState = structs.ReceiptHistoryBeforeStateIncomplete
		if earlier, ok := nearestEarlierReceiptCopy(tasks, receiptId, task.ID); ok {
			before = earlier.snapshot
			record.beforeState = structs.ReceiptHistoryBeforeStateRebuilt
			record.beforeSource = &earlier.source
		}
	}

	record.before = parseReceiptSnapshot(before)
	record.after = parseReceiptSnapshot(description.After)
	if record.before == nil || record.after == nil {
		// Unreadable snapshots: the entry stays, with nothing to compare.
		return historyRecord{task: task}
	}

	return record
}

func parseReceiptSnapshot(snapshot string) *models.Receipt {
	var receipt models.Receipt
	if err := json.Unmarshal([]byte(snapshot), &receipt); err != nil {
		return nil
	}
	return &receipt
}

// receiptHistoryActor says who or what made an entry. A ran-by user is the
// actor. An upload with none ran under a QUICK_SCAN (the scanning user's) or an
// EMAIL_UPLOAD; when that parent is missing, an email receipt still names its
// origin. Anything else is the system. A ran-by or parent id of 0 counts as none
// (older rows stored 0 for "nobody").
func receiptHistoryActor(record historyRecord, parents map[uint]models.SystemTask) structs.ReceiptHistoryActor {
	task := record.task
	if userId := nonZeroUint(task.RanByUserId); userId != nil {
		return structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorUser, UserId: userId}
	}

	if parentId := nonZeroUint(task.AssociatedSystemTaskId); parentId != nil {
		if parent, ok := parents[*parentId]; ok {
			switch parent.Type {
			case models.QUICK_SCAN:
				return structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorQuickScan, UserId: nonZeroUint(parent.RanByUserId)}
			case models.EMAIL_UPLOAD:
				return structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorEmail}
			}
		}
	}

	if record.creation && record.after != nil &&
		(record.after.EmailInbox != "" || record.after.CreatedByString == constants.EMAIL_INTEGRATION_CREATED_BY) {
		return structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorEmail}
	}

	return structs.ReceiptHistoryActor{Type: structs.ReceiptHistoryActorSystem}
}

func nonZeroUint(value *uint) *uint {
	if value == nil || *value == 0 {
		return nil
	}
	copied := *value
	return &copied
}

// receiptHistoryActorKey identifies an actor for the contributor count: each
// user once, and each kind of automation once.
func receiptHistoryActorKey(actor structs.ReceiptHistoryActor) string {
	if actor.Type == structs.ReceiptHistoryActorUser && actor.UserId != nil {
		return actor.Type + ":" + strconv.FormatUint(uint64(*actor.UserId), 10)
	}
	return actor.Type
}

// redact applies what the caller may see on the receipt TODAY to every record,
// evaluating each snapshot as if it lived in the receipt's current group (a
// receipt moved between groups is judged by where it is now, the group whose
// read permission admitted the caller):
//
//   - an entry made by a user the caller cannot see there is dropped; a quick
//     scan keeps its entry but loses the scanning user;
//   - an entry any of whose snapshots was paid by someone the caller cannot see
//     there is dropped, as the receipt itself would be;
//   - every snapshot has its categories and tags stripped to the caller's grants
//     and its members masked, as GetReceiptForUser does.
func (service ReceiptHistoryService) redact(viewerId uint, groupId uint, records []historyRecord) ([]historyRecord, error) {
	permissionService := NewPermissionService(service.TX)
	visibility := permissionService.newGroupVisibilityResolver(viewerId)
	paidByAllowed, paidByUnrestricted, err := permissionService.paidByAllowedForGroup(viewerId, groupId)
	if err != nil {
		return nil, err
	}

	kept := make([]historyRecord, 0, len(records))
	for _, record := range records {
		if record.actor.UserId != nil {
			visible, err := visibility.isVisible(groupId, *record.actor.UserId)
			if err != nil {
				return nil, err
			}
			if !visible {
				if record.actor.Type == structs.ReceiptHistoryActorUser {
					continue
				}
				record.actor.UserId = nil
			}
		}

		if !paidByUnrestricted && !snapshotPayersAllowed(record, paidByAllowed) {
			continue
		}

		kept = append(kept, record)
	}

	if err := service.redactSnapshots(permissionService, viewerId, groupId, kept); err != nil {
		return nil, err
	}

	return kept, nil
}

func snapshotPayersAllowed(record historyRecord, allowed map[uint]struct{}) bool {
	for _, snapshot := range []*models.Receipt{record.before, record.after} {
		if snapshot == nil {
			continue
		}
		if _, ok := allowed[snapshot.PaidByUserID]; !ok {
			return false
		}
	}
	return true
}

// redactSnapshots strips and masks every snapshot in one batch, each placed in
// the receipt's current group for the purpose, then puts back the group it was
// recorded in, which the diff and the JSON both show.
func (service ReceiptHistoryService) redactSnapshots(
	permissionService PermissionService,
	viewerId uint,
	groupId uint,
	records []historyRecord,
) error {
	snapshots := make([]*models.Receipt, 0, len(records)*2)
	for _, record := range records {
		for _, snapshot := range []*models.Receipt{record.before, record.after} {
			if snapshot != nil {
				snapshots = append(snapshots, snapshot)
			}
		}
	}
	if len(snapshots) == 0 {
		return nil
	}

	batch := make([]models.Receipt, len(snapshots))
	for i, snapshot := range snapshots {
		batch[i] = *snapshot
		batch[i].GroupId = groupId
	}

	if err := permissionService.FilterReceiptCategoriesTags(viewerId, batch); err != nil {
		return err
	}
	if err := permissionService.MaskReceiptsForMemberVisibility(viewerId, batch); err != nil {
		return err
	}

	for i, snapshot := range snapshots {
		recordedGroupId := snapshot.GroupId
		*snapshot = batch[i]
		snapshot.GroupId = recordedGroupId
	}
	return nil
}

// diffContext resolves what the diffs show but the snapshots may lack: the name
// of every group a snapshot sits in, and the stored definition of every custom
// field a snapshot carries without one.
func (service ReceiptHistoryService) diffContext(viewerId uint, records []historyRecord) (receiptDiffContext, error) {
	context := receiptDiffContext{
		groupNames:   make(map[uint]string),
		customFields: make(map[uint]models.CustomField),
	}

	groupIds := make([]uint, 0)
	customFieldIds := make([]uint, 0)
	for _, record := range records {
		for _, snapshot := range []*models.Receipt{record.before, record.after} {
			if snapshot == nil {
				continue
			}
			if snapshot.GroupId != 0 {
				groupIds = append(groupIds, snapshot.GroupId)
			}
			for _, value := range snapshot.CustomFields {
				if value.CustomField.Type == "" {
					customFieldIds = append(customFieldIds, value.CustomFieldId)
				}
			}
		}
	}

	groupNames, err := service.visibleGroupNames(viewerId, groupIds)
	if err != nil {
		return receiptDiffContext{}, err
	}
	context.groupNames = groupNames

	if len(customFieldIds) > 0 {
		var customFields []models.CustomField
		err := service.GetDB().Model(&models.CustomField{}).
			Preload("Options").
			Where("id IN ?", customFieldIds).
			Find(&customFields).Error
		if err != nil {
			return receiptDiffContext{}, err
		}
		for _, customField := range customFields {
			context.customFields[customField.ID] = customField
		}
	}

	return context, nil
}

// visibleGroupNames names each group the caller is a member of, or may read
// through app.groups.read, and leaves every other group "(Restricted)".
func (service ReceiptHistoryService) visibleGroupNames(viewerId uint, groupIds []uint) (map[uint]string, error) {
	names := make(map[uint]string)
	if len(groupIds) == 0 {
		return names, nil
	}

	readsAllGroups, err := NewPermissionService(service.TX).HasAppPermissions(viewerId, permissions.AppGroupsRead)
	if err != nil {
		return nil, err
	}
	memberGroupIds, err := repositories.NewGroupMemberRepository(service.TX).GetGroupIdsByUserId(utils.UintToString(viewerId))
	if err != nil {
		return nil, err
	}
	memberOf := make(map[uint]bool, len(memberGroupIds))
	for _, id := range memberGroupIds {
		memberOf[id] = true
	}

	var groups []models.Group
	err = service.GetDB().Model(&models.Group{}).Select("id", "name").Where("id IN ?", groupIds).Find(&groups).Error
	if err != nil {
		return nil, err
	}

	for _, id := range groupIds {
		names[id] = restrictedGroupName
	}
	for _, group := range groups {
		if readsAllGroups || memberOf[group.ID] {
			names[group.ID] = group.Name
		}
	}
	return names, nil
}

// toEntry diffs a record into the entry the client shows.
func (record historyRecord) toEntry(context receiptDiffContext) structs.ReceiptHistoryEntry {
	entry := structs.ReceiptHistoryEntry{
		SystemTaskId: record.task.ID,
		RecordedAt:   systemTaskRecordedAt(record.task),
		Actor:        record.actor,
		Changes:      make([]structs.ReceiptHistoryChange, 0),
	}

	switch {
	case record.failed:
		entry.Action = structs.ReceiptHistoryActionFailed
		entry.ErrorMessage = record.task.ResultDescription
		return entry
	case record.creation:
		entry.Action = structs.ReceiptHistoryActionCreated
		if record.after != nil {
			entry.HasSnapshots = true
			entry.Changes = diffReceiptSnapshots(nil, *record.after, context)
		}
		return entry
	}

	entry.BeforeState = record.beforeState
	if record.beforeSource != nil {
		recordedAt := record.beforeSource.RecordedAt
		entry.BeforeSourceType = string(record.beforeSource.Type)
		entry.BeforeSourceRecordedAt = &recordedAt
	}
	if record.before != nil && record.after != nil {
		entry.HasSnapshots = true
		entry.Changes = diffReceiptSnapshots(record.before, *record.after, context)
	}
	entry.Action = receiptHistoryAction(entry.Changes)

	return entry
}
