package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
)

type SystemTaskRepository struct {
	BaseRepository
}

func NewSystemTaskRepository(tx *gorm.DB) SystemTaskRepository {
	repository := SystemTaskRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository SystemTaskRepository) GetPagedSystemTasks(command commands.GetSystemTaskCommand) ([]models.SystemTask, int64, error) {
	db := repository.GetDB()
	var results []models.SystemTask
	var count int64

	if !isColumnNameValid(command.OrderBy) {
		return nil, 0, errors.New("invalid column name")
	}

	filteredSystemTaskTypes := []models.SystemTaskType{
		models.RECEIPT_UPLOADED,
		models.CHAT_COMPLETION,
		models.OCR_PROCESSING,
	}

	query := db.Model(&models.SystemTask{}).Where("type NOT IN ?", filteredSystemTaskTypes)

	if command.AssociatedEntityId != 0 {
		query = query.Where("associated_entity_id = ?", command.AssociatedEntityId)
	}

	if len(command.AssociatedEntityType) > 0 {
		query = query.Where("associated_entity_type = ?", command.AssociatedEntityType)
	}

	query.Count(&count)

	query = repository.Sort(query, command.OrderBy, command.SortDirection)
	query = query.Scopes(repository.Paginate(command.Page, command.PageSize))

	err := query.Preload(clause.Associations).Preload("ChildSystemTasks.ChildSystemTasks").Find(&results).Error
	if query.Error != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// ActivityVisibilityResolver reports, for a group, the ran-by user ids the caller may
// see under member isolation, or unrestricted == true (see every actor). It is the
// activity analogue of PaidByAllowedResolver: it lets the handler inject the service-layer
// per-group visible set without the repository importing services. A nil resolver skips
// isolation filtering (backward compatible).
type ActivityVisibilityResolver func(groupId uint) (visibleUserIds []uint, unrestricted bool, err error)

func (repository SystemTaskRepository) GetPagedActivities(
	command commands.PagedActivityRequestCommand,
	resolver ActivityVisibilityResolver,
) (
	[]structs.Activity,
	int64,
	error,
) {
	db := repository.GetDB()
	var results []structs.Activity
	var count int64

	if !isColumnNameValid(command.OrderBy) {
		return nil, 0, errors.New("invalid column name")
	}

	systemTaskTypesToGet := []models.SystemTaskType{
		models.QUICK_SCAN,
		models.RECEIPT_UPLOADED,
		models.RECEIPT_UPDATED,
		models.EMAIL_UPLOAD,
	}

	query := db.Model(&models.SystemTask{}).
		Omit("can_be_restarted").
		Where("type IN ?", systemTaskTypesToGet).
		Where("group_id IN ?", command.GroupIds).
		Not(db.Where("type = ? AND ran_by_user_id IS NULL", models.RECEIPT_UPLOADED))

	// Member isolation: drop activities run by a user the caller may not see in that
	// activity's group IN THE QUERY (before Count + pagination), so TotalCount and the
	// returned page both reflect only visible rows and DB-side LIMIT/OFFSET is preserved.
	query, err := repository.applyActivityVisibilityDisjunction(query, command.GroupIds, resolver)
	if err != nil {
		return nil, 0, err
	}

	query.Count(&count)

	query = repository.Sort(query, command.OrderBy, command.SortDirection)
	query = query.Scopes(repository.Paginate(command.Page, command.PageSize))

	query.Find(&results)

	return results, count, nil
}

// applyActivityVisibilityDisjunction AND-s a per-group actor-visibility disjunction onto
// the query, mirroring ReceiptRepository.ApplyPaidByDisjunction. A nil resolver adds no
// predicate. For each group the caller sees every actor (unrestricted) the clause is just
// group_id = G; otherwise it is group_id = G AND (ran_by_user_id IS NULL OR ran_by_user_id
// IN <visible ids>) — so system actions (nil ran-by) stay visible and only visible
// members' activities survive. Applied before Count so pagination stays consistent.
func (repository SystemTaskRepository) applyActivityVisibilityDisjunction(
	query *gorm.DB,
	groupIds []uint,
	resolver ActivityVisibilityResolver,
) (*gorm.DB, error) {
	if resolver == nil {
		return query, nil
	}
	if len(groupIds) == 0 {
		return query.Where("1 = 0"), nil
	}

	disjunction := repository.GetDB().Session(&gorm.Session{NewDB: true})
	for _, groupId := range groupIds {
		visibleIds, unrestricted, err := resolver(groupId)
		if err != nil {
			return nil, err
		}
		if unrestricted {
			disjunction = disjunction.Or("group_id = ?", groupId)
		} else {
			groupCondition := repository.GetDB().Session(&gorm.Session{NewDB: true}).
				Where("group_id = ?", groupId).
				Where("(ran_by_user_id IS NULL OR ran_by_user_id IN ?)", activityInValues(visibleIds))
			disjunction = disjunction.Or(groupCondition)
		}
	}

	return query.Where(disjunction), nil
}

// activityInValues guards the IN clause against an empty visible set (a restricted set
// always contains at least the caller's own id, but guard defensively): ran-by user ids
// start at 1, so 0 matches no row, yielding "see nothing" rather than a malformed IN ().
func activityInValues(visibleUserIds []uint) []uint {
	if len(visibleUserIds) == 0 {
		return []uint{0}
	}
	return visibleUserIds
}

func isColumnNameValid(columnName string) bool {
	return columnName == "type" || columnName == "status" || columnName == "associated_entity_type" || columnName == "associated_entity_id" || columnName == "started_at" || columnName == "ended_at" || columnName == "result_description" || columnName == "ran_by_user_id"
}

func (repository SystemTaskRepository) CreateSystemTask(command commands.UpsertSystemTaskCommand) (models.SystemTask, error) {
	db := repository.GetDB()

	systemTask := models.SystemTask{
		Type:                   command.Type,
		Status:                 command.Status,
		AssociatedEntityType:   command.AssociatedEntityType,
		AssociatedEntityId:     command.AssociatedEntityId,
		StartedAt:              command.StartedAt,
		EndedAt:                command.EndedAt,
		ResultDescription:      command.ResultDescription,
		RanByUserId:            command.RanByUserId,
		ReceiptId:              command.ReceiptId,
		GroupId:                command.GroupId,
		AssociatedSystemTaskId: command.AssociatedSystemTaskId,
		AsynqTaskId:            command.AsynqTaskId,
	}

	err := db.Create(&systemTask).Error
	if err != nil {
		return models.SystemTask{}, err
	}

	if command.AssociatedSystemTaskId != nil && systemTask.Status == models.SYSTEM_TASK_FAILED {
		var parentSystemTask models.SystemTask
		db.Model(&models.SystemTask{}).Where("id = ?", command.AssociatedSystemTaskId).Find(&parentSystemTask)

		if parentSystemTask.Status == models.SYSTEM_TASK_SUCCEEDED {
			db.Model(&parentSystemTask).Update("status", models.SYSTEM_TASK_FAILED)
		}

	}

	return systemTask, nil
}

func (repository SystemTaskRepository) DeleteSystemTaskByAssociatedEntityId(
	associatedEntityId string,
	emailType models.AssociatedEntityType,
) error {
	db := repository.GetDB()
	err := db.Where("associated_entity_id = ? and associated_entity_type = ?", associatedEntityId, emailType).Delete(&models.SystemTask{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository SystemTaskRepository) GetSystemTaskById(id uint) (models.SystemTask, error) {
	db := repository.GetDB()
	var systemTask models.SystemTask

	err := db.Model(&models.SystemTask{}).Where("id = ?", id).First(&systemTask).Error
	if err != nil {
		return models.SystemTask{}, err
	}

	return systemTask, nil
}

func (repository SystemTaskRepository) AssociateSystemTaskToReceipt(receiptId uint, systemTaskId uint) error {
	db := repository.GetDB()
	return db.Model(&models.SystemTask{}).Where("id = ?", systemTaskId).Update("receipt_id", receiptId).Error
}
