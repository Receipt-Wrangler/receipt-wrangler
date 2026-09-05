package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting/receiptsource"
	"receipt-wrangler/api/internal/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReceiptRepository struct {
	BaseRepository
}

func NewReceiptRepository(tx *gorm.DB) ReceiptRepository {
	repository := ReceiptRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository ReceiptRepository) BeforeUpdateReceipt(currentReceipt models.Receipt, updatedReceipt models.Receipt) (err error) {
	db := repository.GetDB()
	if updatedReceipt.GroupId > 0 && currentReceipt.GroupId != updatedReceipt.GroupId && len(currentReceipt.ImageFiles) > 0 {
		var oldGroup models.Group
		var newGroup models.Group

		err = db.Table("groups").Where("id = ?", currentReceipt.GroupId).Select("id", "name").Find(&oldGroup).Error
		if err != nil {
			return err
		}

		err = db.Table("groups").Where("id = ?", updatedReceipt.GroupId).Select("id", "name").Find(&newGroup).Error
		if err != nil {
			return err
		}

		oldGroupPath, err := utils.BuildGroupPathString(utils.UintToString(oldGroup.ID), oldGroup.Name)
		if err != nil {
			return err
		}

		newGroupPath, err := utils.BuildGroupPathString(utils.UintToString(newGroup.ID), newGroup.Name)
		if err != nil {
			return err
		}

		for _, fileData := range currentReceipt.ImageFiles {
			filename := utils.BuildFileName(utils.UintToString(currentReceipt.ID), utils.UintToString(fileData.ID), fileData.Name)

			oldFilePath := filepath.Join(oldGroupPath, filename)
			newFilePathPath := filepath.Join(newGroupPath, filename)

			err := utils.RenameDataPath(oldFilePath, newFilePathPath)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func createFailedUpdateSystemTask(command commands.UpsertSystemTaskCommand, err error) {
	endedAt := time.Now()
	command.EndedAt = &endedAt
	command.Status = models.SYSTEM_TASK_FAILED
	command.ResultDescription = err.Error()

	repository := NewSystemTaskRepository(nil)
	repository.CreateSystemTask(command)
}

func (repository ReceiptRepository) UpdateReceipt(id string, command commands.UpsertReceiptCommand, userId uint) (models.Receipt, error) {
	db := repository.GetDB()

	systemTaskResultDescription := map[string]interface{}{}
	var endedAt time.Time
	stringId, err := utils.StringToUint(id)
	if err != nil {
		return models.Receipt{}, err
	}
	var currentReceipt models.Receipt
	var ranByUserId = userId

	systemTask := commands.UpsertSystemTaskCommand{
		Type:                 models.RECEIPT_UPDATED,
		AssociatedEntityType: models.RECEIPT,
		AssociatedEntityId:   stringId,
		StartedAt:            time.Now(),
		EndedAt:              &endedAt,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		RanByUserId:          &ranByUserId,
	}

	updatedReceipt, err := command.ToReceipt()
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}

	err = db.Table("receipts").Where("id = ?", id).Preload(clause.Associations).Find(&currentReceipt).Error
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}

	systemTask.GroupId = &currentReceipt.GroupId
	systemTask.ReceiptId = &currentReceipt.ID

	// NOTE: ID and field used for afterReceiptUpdated
	updatedReceipt.ID = currentReceipt.ID
	updatedReceipt.ResolvedDate = currentReceipt.ResolvedDate
	before, err := currentReceipt.ToString()
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}
	systemTaskResultDescription["before"] = before

	err = db.Transaction(func(tx *gorm.DB) error {
		repository.SetTransaction(tx)

		txErr := repository.BeforeUpdateReceipt(currentReceipt, updatedReceipt)
		if txErr != nil {
			return txErr
		}

		txErr = tx.Session(&gorm.Session{FullSaveAssociations: true}).Model(&currentReceipt).Updates(&updatedReceipt).Error
		if txErr != nil {
			return txErr
		}

		txErr = tx.Model(&currentReceipt).Association("Tags").Replace(&updatedReceipt.Tags)
		if txErr != nil {
			return txErr
		}

		txErr = tx.Model(&currentReceipt).Association("Categories").Replace(&updatedReceipt.Categories)
		if txErr != nil {
			return txErr
		}

		txErr = tx.Model(&currentReceipt).Association("ReceiptItems").Replace(&updatedReceipt.ReceiptItems)
		if txErr != nil {
			return txErr
		}

		txErr = tx.Model(&currentReceipt).Association("CustomFields").Replace(&updatedReceipt.CustomFields)
		if txErr != nil {
			return txErr
		}

		for _, item := range updatedReceipt.ReceiptItems {
			txErr = tx.Model(&item).Association("Categories").Replace(&item.Categories)
			if txErr != nil {
				return txErr
			}

			txErr = tx.Model(&item).Association("Tags").Replace(&item.Tags)
			if txErr != nil {
				return txErr
			}

			txErr = tx.Model(&item).Association("LinkedItems").Replace(&item.LinkedItems)
			if txErr != nil {
				return txErr
			}

			// Update categories and tags for linked items
			for _, linkedItem := range item.LinkedItems {
				txErr = tx.Model(&linkedItem).Association("Categories").Replace(&linkedItem.Categories)
				if txErr != nil {
					return txErr
				}

				txErr = tx.Model(&linkedItem).Association("Tags").Replace(&linkedItem.Tags)
				if txErr != nil {
					return txErr
				}
			}
		}

		err = repository.AfterReceiptUpdated(&updatedReceipt)
		if err != nil {
			return err
		}

		repository.ClearTransaction()
		return nil
	})
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}

	fullyLoadedReceipt, err := repository.GetFullyLoadedReceiptById(id)
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}

	after, err := fullyLoadedReceipt.ToString()
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}

	systemTaskResultDescription["after"] = after
	endedAt = time.Now()
	systemTask.EndedAt = &endedAt

	resultDescriptionBytes, err := json.Marshal(systemTaskResultDescription)
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}
	systemTask.ResultDescription = string(resultDescriptionBytes)

	systemTaskRepository := NewSystemTaskRepository(nil)
	_, err = systemTaskRepository.CreateSystemTask(systemTask)
	if err != nil {
		createFailedUpdateSystemTask(systemTask, err)
		return models.Receipt{}, err
	}

	return fullyLoadedReceipt, nil
}

// TODO: Delete categories/tags here associated with items before deleting the items mkay
func (repository ReceiptRepository) AfterReceiptUpdated(updatedReceipt *models.Receipt) error {
	db := repository.GetDB()

	// TODO: Move this  to a scheduled job
	// Clean up junction tables for orphaned items
	orphanedItemsSubquery := db.Table("items").Select("id").Where("receipt_id IS NULL")

	// Clean up item_linked_items junction table - remove associations where either side is orphaned
	err := db.Table("item_linked_items").Where("item_id IN (?) OR linked_item_id IN (?)",
		orphanedItemsSubquery,
		orphanedItemsSubquery,
	).Delete(&struct{}{}).Error
	if err != nil {
		return err
	}

	// Clean up item_categories junction table
	err = db.Table("item_categories").Where("item_id IN (?)",
		orphanedItemsSubquery,
	).Delete(&struct{}{}).Error
	if err != nil {
		return err
	}

	// Clean up item_tags junction table
	err = db.Table("item_tags").Where("item_id IN (?)",
		orphanedItemsSubquery,
	).Delete(&struct{}{}).Error
	if err != nil {
		return err
	}

	// TODO: Move this  to a scheduled job
	// Delete the orphaned items themselves
	err = db.Where("receipt_id IS NULL").Delete(&models.Item{}).Error
	if err != nil {
		return err
	}

	if updatedReceipt.ID > 0 && updatedReceipt.Status == models.RESOLVED && updatedReceipt.ResolvedDate == nil {
		now := time.Now().UTC()
		err = db.Table("receipts").Where("id = ?", updatedReceipt.ID).Update("resolved_date", now).Error
	} else if updatedReceipt.ID > 0 && updatedReceipt.Status != models.RESOLVED && updatedReceipt.ResolvedDate != nil {
		err = db.Table("receipts").Where("id = ?", updatedReceipt.ID).Update("resolved_date", nil).Error
	}
	if err != nil {
		return err
	}

	if updatedReceipt.Status == models.RESOLVED && updatedReceipt.ID > 0 {
		err := repository.UpdateItemsToStatus(updatedReceipt, models.ITEM_RESOLVED)
		if err != nil {
			return err
		}
	}

	if updatedReceipt.Status == models.DRAFT && updatedReceipt.ID > 0 {
		err := repository.UpdateItemsToStatus(updatedReceipt, models.ITEM_DRAFT)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository ReceiptRepository) UpdateItemsToStatus(receipt *models.Receipt, status models.ItemStatus) error {
	db := repository.GetDB()
	var items []models.Item
	var itemIdsToUpdate []uint

	err := db.Table("items").Where("receipt_id = ?", receipt.ID).Find(&items).Error
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Status != status {
			itemIdsToUpdate = append(itemIdsToUpdate, item.ID)
		}
	}

	if len(itemIdsToUpdate) > 0 {
		err := db.Table("items").Where("id IN ?", itemIdsToUpdate).UpdateColumn("status", status).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository ReceiptRepository) CreateReceipt(
	command commands.UpsertReceiptCommand,
	createdByUserID uint,
	createSystemTask bool,
) (models.Receipt, error) {
	db := repository.GetDB()
	notificationRepository := NewNotificationRepository(nil)
	receipt, err := command.ToReceipt()
	if err != nil {
		return models.Receipt{}, err
	}

	if receipt.GroupId > 0 {
		receipt.CreatedBy = &createdByUserID
	}

	systemTask := commands.UpsertSystemTaskCommand{
		Type:                 models.RECEIPT_UPLOADED,
		AssociatedEntityType: models.RECEIPT,
		StartedAt:            time.Now(),
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		RanByUserId:          &createdByUserID,
	}

	// Extract linked items before creating receipt
	type LinkedItemData struct {
		ParentItemIndex int
		LinkedItem      models.Item
	}
	var linkedItemsData []LinkedItemData

	for i := range receipt.ReceiptItems {
		if len(receipt.ReceiptItems[i].LinkedItems) > 0 {
			for _, linkedItem := range receipt.ReceiptItems[i].LinkedItems {
				linkedItemsData = append(linkedItemsData, LinkedItemData{
					ParentItemIndex: i,
					LinkedItem:      linkedItem,
				})
			}
			// Clear linked items from the receipt item for initial creation
			receipt.ReceiptItems[i].LinkedItems = []models.Item{}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		repository.SetTransaction(tx)
		notificationRepository.SetTransaction(tx)

		// First nested transaction: Create receipt without linked items
		err := tx.Transaction(func(tx2 *gorm.DB) error {
			err := tx2.Model(models.Receipt{}).Select("*").Create(&receipt).Error
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		// Second nested transaction: Handle linked items
		if len(linkedItemsData) > 0 {
			err = tx.Transaction(func(tx3 *gorm.DB) error {
				for _, linkedData := range linkedItemsData {
					// Set the receipt ID for the linked item
					linkedData.LinkedItem.ReceiptId = receipt.ID

					// Create the linked item
					err := tx3.Model(models.Item{}).Create(&linkedData.LinkedItem).Error
					if err != nil {
						return err
					}

					// Handle linked item's categories
					if len(linkedData.LinkedItem.Categories) > 0 {
						err = tx3.Model(&linkedData.LinkedItem).Association("Categories").Replace(&linkedData.LinkedItem.Categories)
						if err != nil {
							return err
						}
					}

					// Handle linked item's tags
					if len(linkedData.LinkedItem.Tags) > 0 {
						err = tx3.Model(&linkedData.LinkedItem).Association("Tags").Replace(&linkedData.LinkedItem.Tags)
						if err != nil {
							return err
						}
					}

					// Update the parent item's LinkedItems association
					parentItem := &receipt.ReceiptItems[linkedData.ParentItemIndex]
					err = tx3.Model(parentItem).Association("LinkedItems").Append(&linkedData.LinkedItem)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		var userIdsToOmit []interface{} = make([]interface{}, 1)
		userIdsToOmit = append(userIdsToOmit, *receipt.CreatedBy)

		notificationBody := fmt.Sprintf("The receipt: %s has been uploaded to the group %s. Check it out! %s", receipt.Name, BuildParamaterisedString("groupId", receipt.GroupId, "name", "string"), BuildParamaterisedString("receiptId", receipt.ID, "", "link"))
		err = notificationRepository.SendNotificationToGroup(receipt.GroupId, "Receipt Uploaded", notificationBody, models.NOTIFICATION_TYPE_NORMAL, userIdsToOmit)
		if err != nil {
			return err
		}

		err = repository.AfterReceiptUpdated(&receipt)
		if err != nil {
			return err
		}

		// Restore the DB/transaction this method was invoked with (captured as
		// `db` above) so the read-back below runs on the same connection that
		// created the receipt. Clearing to nil would fall back to a fresh pooled
		// connection which, when CreateReceipt runs inside an outer transaction
		// (e.g. QuickScan, email attachment ingest), cannot see the still
		// uncommitted receipt row and would return a zero-ID receipt.
		repository.SetTransaction(db)
		notificationRepository.ClearTransaction()
		return nil
	})
	if err != nil {
		if !createSystemTask {
			createFailedUpdateSystemTask(systemTask, err)
		}
		return models.Receipt{}, err
	}

	fullyLoadedReceipt, err := repository.GetFullyLoadedReceiptById(utils.UintToString(receipt.ID))
	if err != nil {
		if !createSystemTask {
			createFailedUpdateSystemTask(systemTask, err)
		}
		return models.Receipt{}, err
	}

	// GetFullyLoadedReceiptById uses Find, which returns an empty receipt (ID 0)
	// with no error when the row is not visible. Guard against that so callers get
	// a clear failure here instead of a downstream foreign key violation.
	if fullyLoadedReceipt.ID == 0 {
		err = fmt.Errorf("created receipt %s could not be reloaded", utils.UintToString(receipt.ID))
		if !createSystemTask {
			createFailedUpdateSystemTask(systemTask, err)
		}
		return models.Receipt{}, err
	}

	if createSystemTask {
		endedAt := time.Now()
		systemTask.EndedAt = &endedAt
		systemTask.AssociatedEntityId = fullyLoadedReceipt.ID
		newReceiptString, err := fullyLoadedReceipt.ToString()
		if err != nil {
			return models.Receipt{}, err
		}

		systemTask.ReceiptId = &fullyLoadedReceipt.ID
		systemTask.GroupId = &fullyLoadedReceipt.GroupId
		systemTask.ResultDescription = newReceiptString

		systemTaskRepository := NewSystemTaskRepository(nil)
		_, err = systemTaskRepository.CreateSystemTask(systemTask)
		if err != nil {
			return models.Receipt{}, err
		}
	}

	return fullyLoadedReceipt, nil
}

func (repository ReceiptRepository) GetReceiptById(receiptId string) (models.Receipt, error) {
	db := GetDB()
	var receipt models.Receipt
	err := db.Model(models.Receipt{}).Where("id = ?", receiptId).First(&receipt).Debug().Error
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

// PaidByAllowedResolver returns the paid_by_user_id values a user may see in a
// group, and whether they are unrestricted (see every payer). It lets the receipt
// repository apply the role-based "paid by" visibility filter without importing
// the service layer that resolves grants. Pass a nil resolver to
// GetPagedReceiptsByGroupId to skip paid-by filtering entirely (internal/system
// callers).
type PaidByAllowedResolver func(groupId uint) (allowedUserIds []uint, unrestricted bool, err error)

func (repository ReceiptRepository) GetPagedReceiptsByGroupId(
	userId uint,
	groupId string,
	pagedRequest commands.ReceiptPagedRequestCommand,
	associations []string,
	paidByResolver PaidByAllowedResolver,
) ([]models.Receipt, int64, error) {
	var receipts []models.Receipt
	var count int64

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return nil, 0, err
	}
	groupRepository := NewGroupRepository(nil)
	isAllGroup, err := groupRepository.IsAllGroup(uintGroupId)
	if err != nil {
		return nil, 0, err
	}

	// Apply filter
	query, err := repository.BuildGormFilterQuery(pagedRequest)
	if err != nil {
		return nil, 0, err
	}

	// Filter receipts by group
	var memberGroupIds []uint
	if isAllGroup {
		groupMemberRepository := NewGroupMemberRepository(nil)
		memberGroupIds, err = groupMemberRepository.GetGroupIdsByUserId(utils.UintToString(userId))
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("group_id IN ?", memberGroupIds)
	} else {
		query = query.Where("group_id = ?", groupId)
	}

	// Apply role-based "paid by" visibility, AND-ed with the group scope above and
	// BEFORE the count below so totalCount matches the rows actually returned (a
	// post-fetch filter would corrupt pagination).
	if paidByResolver != nil {
		query, err = repository.applyPaidByVisibility(query, uintGroupId, isAllGroup, memberGroupIds, paidByResolver)
		if err != nil {
			return nil, 0, err
		}
	}

	// Set order by
	if len(pagedRequest.OrderBy) == 0 {
		pagedRequest.OrderBy = constants.DEFAULT_RECEIPT_ORDER_BY
	}

	if !commands.IsValidSortDirection(pagedRequest.SortDirection) {
		return nil, 0, errors.New("untrusted value " + pagedRequest.OrderBy + " " + string(pagedRequest.SortDirection))
	}

	if customFieldId, isCustomField := receiptsource.ParseCustomFieldKey(pagedRequest.OrderBy); isCustomField {
		query, err = repository.orderByCustomField(query, customFieldId, pagedRequest.SortDirection)
		if err != nil {
			return nil, 0, err
		}
	} else if repository.isTrustedValue(pagedRequest) {
		orderBy := pagedRequest.OrderBy
		query = query.Order(orderBy + " " + string(pagedRequest.SortDirection))
	} else {
		return nil, 0, errors.New("untrusted value " + pagedRequest.OrderBy + " " + string(pagedRequest.SortDirection))
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Preload("Categories").Preload("Tags")

	if pagedRequest.PageSize > 0 && pagedRequest.Page > 0 {
		query = query.Scopes(repository.Paginate(pagedRequest.Page, pagedRequest.PageSize))
	}

	if associations != nil && len(associations) > 0 {
		for _, association := range associations {
			query = query.Preload(association)
		}
	}

	// Run Query
	err = query.Find(&receipts).Error
	if err != nil {
		return nil, 0, err
	}

	return receipts, count, nil
}

// applyPaidByVisibility narrows query to the receipts a user may see by the
// "paid by" user, using resolver to look up each group's allowed set. For a
// single group it adds a simple paid_by_user_id IN (...) constraint; for the
// all-group view it builds a per-group disjunction so each group applies its own
// allowed set (a member may hold a different group role per group). The whole
// disjunction is AND-ed onto query, preserving the existing group scope and
// keeping the row count correct.
func (repository ReceiptRepository) applyPaidByVisibility(
	query *gorm.DB,
	uintGroupId uint,
	isAllGroup bool,
	memberGroupIds []uint,
	resolver PaidByAllowedResolver,
) (*gorm.DB, error) {
	if !isAllGroup {
		allowed, unrestricted, err := resolver(uintGroupId)
		if err != nil {
			return nil, err
		}
		if unrestricted {
			return query, nil
		}
		return query.Where("paid_by_user_id IN ?", paidByInValues(allowed)), nil
	}

	return repository.ApplyPaidByDisjunction(query, memberGroupIds, resolver)
}

// ApplyPaidByDisjunction AND-s a per-group "paid by" visibility disjunction onto
// query across memberGroupIds: each group contributes either `group_id = G`
// (unrestricted) or `(group_id = G AND paid_by_user_id IN (allowed))`, OR-ed
// together. It is shared by the all-group paged read and by search — both scope
// receipts to a member's groups and must apply each group's own paid-by role in
// SQL, BEFORE any LIMIT, so visible rows are not lost to a pre-filter row cap.
func (repository ReceiptRepository) ApplyPaidByDisjunction(
	query *gorm.DB,
	memberGroupIds []uint,
	resolver PaidByAllowedResolver,
) (*gorm.DB, error) {
	// memberGroupIds is the caller's member groups. With none, there is nothing to
	// see — fail closed explicitly rather than leave the (empty) disjunction as a
	// silent no-op that adds no predicate, mirroring the single-group IN (0) guard.
	if len(memberGroupIds) == 0 {
		return query.Where("1 = 0"), nil
	}

	disjunction := repository.GetDB().Session(&gorm.Session{NewDB: true})
	for _, groupId := range memberGroupIds {
		allowed, unrestricted, err := resolver(groupId)
		if err != nil {
			return nil, err
		}
		if unrestricted {
			disjunction = disjunction.Or("group_id = ?", groupId)
		} else {
			groupCondition := repository.GetDB().Session(&gorm.Session{NewDB: true}).
				Where("group_id = ?", groupId).
				Where("paid_by_user_id IN ?", paidByInValues(allowed))
			disjunction = disjunction.Or(groupCondition)
		}
	}

	return query.Where(disjunction), nil
}

// paidByInValues guards the IN clause against an empty restricted set: paid-by
// user ids start at 1, so 0 matches no receipt, yielding "see nothing" rather
// than a malformed IN (). A restricted role normally always has at least one id
// (a grant or the resolved self id), but a role whose only granted user was
// deleted lands here.
func paidByInValues(allowedUserIds []uint) []uint {
	if len(allowedUserIds) == 0 {
		return []uint{0}
	}
	return allowedUserIds
}

func (repository ReceiptRepository) BuildGormFilterQuery(pagedRequest commands.ReceiptPagedRequestCommand) (*gorm.DB, error) {
	query := repository.GetDB().Model(models.Receipt{})
	// Name
	if pagedRequest.Filter.Name.Value != nil {
		name := pagedRequest.Filter.Name.Value.(string)
		if len(name) > 0 {
			query = repository.buildFilterQuery(query, name, pagedRequest.Filter.Name.Operation, "name", false)
		}
	}

	// Date
	if pagedRequest.Filter.Date.Value != nil {
		var date interface{}
		isBetweenOperation := pagedRequest.Filter.Date.Operation == commands.BETWEEN
		if isBetweenOperation {
			date = pagedRequest.Filter.Date.Value.([]interface{})
		} else {
			date = pagedRequest.Filter.Date.Value.(string)
		}

		query = repository.buildFilterQuery(query, date, pagedRequest.Filter.Date.Operation, "date", isBetweenOperation)
	}

	// Paid By
	if pagedRequest.Filter.PaidBy.Value != nil {
		paidBy := pagedRequest.Filter.PaidBy.Value.([]interface{})
		if len(paidBy) > 0 {
			query = repository.buildFilterQuery(query, paidBy, pagedRequest.Filter.PaidBy.Operation, "paid_by_user_id", true)
		}
	}

	// Categories
	if pagedRequest.Filter.Categories.Value != nil {
		categories := pagedRequest.Filter.Categories.Value.([]interface{})
		if len(categories) > 0 {
			if pagedRequest.Filter.Categories.Operation == commands.CONTAINS {
				query = query.Where("id IN (?)", db.Table("receipt_categories").Select("receipt_id").Where("category_id IN ?", categories))
			}
		}

	}

	// Tags
	if pagedRequest.Filter.Tags.Value != nil {
		tags := pagedRequest.Filter.Tags.Value.([]interface{})
		if len(tags) > 0 {
			if pagedRequest.Filter.Tags.Operation == commands.CONTAINS {
				query = query.Where("id IN (?)", db.Table("receipt_tags").Select("receipt_id").Where("tag_id IN ?", tags))
			}
		}
	}

	// Amount
	if pagedRequest.Filter.Amount.Value != nil {
		var amount interface{}
		if pagedRequest.Filter.Amount.Operation == commands.BETWEEN {
			amount = pagedRequest.Filter.Amount.Value.([]interface{})
		} else {
			amount = pagedRequest.Filter.Amount.Value.(float64)
		}
		query = repository.buildFilterQuery(
			query,
			amount,
			pagedRequest.Filter.Amount.Operation,
			"amount", pagedRequest.Filter.Amount.Operation == commands.BETWEEN,
		)
	}

	// Status
	if pagedRequest.Filter.Status.Value != nil {
		status := pagedRequest.Filter.Status.Value.([]interface{})
		if len(status) > 0 {
			query = repository.buildFilterQuery(query, status, pagedRequest.Filter.Status.Operation, "status", true)
		}
	}

	// Group
	if pagedRequest.Filter.Group.Value != nil {
		groups := pagedRequest.Filter.Group.Value.([]interface{})
		if len(groups) > 0 {
			query = repository.buildFilterQuery(query, groups, pagedRequest.Filter.Group.Operation, "group_id", true)
		}
	}

	// Resolved Date
	if pagedRequest.Filter.ResolvedDate.Value != nil {
		var resolvedDate interface{}
		isBetweenOperation := pagedRequest.Filter.ResolvedDate.Operation == commands.BETWEEN
		if isBetweenOperation {
			resolvedDate = pagedRequest.Filter.ResolvedDate.Value.(interface{})
		} else {
			resolvedDate = pagedRequest.Filter.ResolvedDate.Value.(string)
		}

		query = repository.buildFilterQuery(
			query,
			resolvedDate,
			pagedRequest.Filter.ResolvedDate.Operation,
			"resolved_date",
			isBetweenOperation,
		)
	}

	// Added At
	if pagedRequest.Filter.CreatedAt.Value != nil {
		var addedAt interface{}
		isBetweenOperation := pagedRequest.Filter.CreatedAt.Operation == commands.BETWEEN
		if isBetweenOperation {
			addedAt = pagedRequest.Filter.CreatedAt.Value.([]interface{})
		} else {
			addedAt = pagedRequest.Filter.CreatedAt.Value.(string)
		}

		query = repository.buildFilterQuery(
			query,
			addedAt,
			pagedRequest.Filter.CreatedAt.Operation,
			"created_at",
			isBetweenOperation,
		)
	}

	return query, nil
}

func (repository ReceiptRepository) buildFilterQuery(runningQuery *gorm.DB, value interface{}, operation commands.FilterOperation, fieldName string, isArray bool) *gorm.DB {
	if operation == commands.EQUALS && !isArray {
		return runningQuery.Where(fmt.Sprintf("%v = ?", fieldName), value)
	}

	if operation == commands.CONTAINS && !isArray {
		searchValue := value.(string)
		searchValue = "%" + searchValue + "%"
		return runningQuery.Where(fmt.Sprintf("%v LIKE ?", fieldName), searchValue)
	}

	if operation == commands.CONTAINS && isArray {
		return runningQuery.Where(fmt.Sprintf("%v IN ?", fieldName), value)
	}

	if operation == commands.GREATER_THAN && !isArray {
		return runningQuery.Where(fmt.Sprintf("%v > ?", fieldName), value)
	}

	if operation == commands.LESS_THAN && !isArray {
		return runningQuery.Where(fmt.Sprintf("%v < ?", fieldName), value)
	}

	if operation == commands.BETWEEN {
		arrayValue := value.([]interface{})
		if len(arrayValue) != 2 {
			return runningQuery
		}

		return runningQuery.Where(fmt.Sprintf("%v >= ? AND %v <= ?", fieldName, fieldName), arrayValue[0], arrayValue[1])
	}

	if operation == commands.WITHIN_CURRENT_MONTH {
		now := time.Now()
		beginningOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

		return runningQuery.Where(fmt.Sprintf("%v >= ? AND %v <= ?", fieldName, fieldName), beginningOfMonth, endOfToday)
	}

	return runningQuery
}

// customFieldSortExpressions returns, for a custom field's type, the expression a
// value row is sorted by and the column that must be non-null for the row to
// count. They differ only for CURRENCY.
//
// CurrencyValue is a *decimal.Decimal with no gorm type tag, and gorm types an
// untagged driver.Valuer by what Value() returns - a string for decimal - so
// currency_value is a text column and a bare sort on it is lexicographic
// ("100" < "20"). CAST(x AS DECIMAL(20,6)) fixes that and is portable: MySQL and
// Postgres take DECIMAL(p,s) directly, and SQLite gives that type name NUMERIC
// affinity. Retyping the column would remove the cast, but that is a migration
// over existing data and belongs in its own change (see api/CLAUDE.md).
func customFieldSortExpressions(customFieldType models.CustomFieldType) (sortExpression string, notNullColumn string, ok bool) {
	// Every column is table-qualified: a SELECT joins custom_field_options, which
	// carries a custom_field_id of its own.
	switch customFieldType {
	case models.TEXT:
		return "custom_field_values.string_value", "custom_field_values.string_value", true
	case models.DATE:
		return "custom_field_values.date_value", "custom_field_values.date_value", true
	case models.BOOLEAN:
		return "custom_field_values.boolean_value", "custom_field_values.boolean_value", true
	case models.CURRENCY:
		return "CAST(custom_field_values.currency_value AS DECIMAL(20,6))", "custom_field_values.currency_value", true
	case models.SELECT:
		// A select stores an option id; readers see the option's text.
		return "custom_field_options.value", "custom_field_values.select_value", true
	}

	return "", "", false
}

// orderByCustomField orders query by a receipt's value for one custom field.
//
// The value is read with a correlated subquery rather than a join: nothing stops
// a receipt holding several values for one field (custom_field_values carries no
// unique index on receipt_id + custom_field_id), and a join would multiply the
// receipt rows, corrupting both the total count and pagination.
//
// Which of several values wins matches the reporting engine
// (receiptsource.addCustomFields): the lowest id among the values that actually
// resolve, so an empty low-id row cannot hide a real one. That is what the
// IS NOT NULL clause is for, and why a SELECT joins its options rather than
// left-joining them - an option id that no longer resolves is skipped, not
// preferred.
//
// A field that no longer exists sorts by the default column instead of erroring:
// clients persist their sort, and a deleted custom field must not make every
// subsequent list load fail.
func (repository ReceiptRepository) orderByCustomField(
	query *gorm.DB,
	customFieldId uint,
	sortDirection commands.SortDirection,
) (*gorm.DB, error) {
	customFieldRepository := NewCustomFieldRepository(nil)
	customFields, err := customFieldRepository.GetCustomFieldsByIds([]uint{customFieldId})
	if err != nil {
		return nil, err
	}

	if len(customFields) == 0 {
		return query.Order(constants.DEFAULT_RECEIPT_ORDER_BY + " " + string(sortDirection)), nil
	}

	sortExpression, notNullColumn, ok := customFieldSortExpressions(customFields[0].Type)
	if !ok {
		return query.Order(constants.DEFAULT_RECEIPT_ORDER_BY + " " + string(sortDirection)), nil
	}

	valueQuery := repository.GetDB().
		Model(&models.CustomFieldValue{}).
		Select(sortExpression).
		Where("custom_field_values.receipt_id = receipts.id").
		Where("custom_field_values.custom_field_id = ?", customFieldId).
		Where(notNullColumn + " IS NOT NULL").
		Order("custom_field_values.id").
		Limit(1)

	if customFields[0].Type == models.SELECT {
		valueQuery = valueQuery.Joins(
			"JOIN custom_field_options ON custom_field_options.id = custom_field_values.select_value",
		)
	}

	// The direction and the tiebreaker have to live in this one expression.
	// clause.OrderBy builds its Expression *instead of* its Columns and Desc, and
	// a second Order() call would silently replace this clause rather than append
	// to it. The tiebreaker is not cosmetic: a boolean or select field has a
	// handful of distinct values, and without a unique last term LIMIT/OFFSET
	// paging repeats and skips rows between pages.
	return query.Order(clause.OrderBy{Expression: clause.Expr{
		SQL:  "(?) " + string(sortDirection) + ", receipts.id DESC",
		Vars: []any{valueQuery},
	}}), nil
}

func (repository ReceiptRepository) isTrustedValue(pagedRequest commands.ReceiptPagedRequestCommand) bool {
	orderByTrusted := []interface{}{"date", "name", "paid_by_user_id", "amount", "categories", "tags", "status", "resolved_date", "created_at"}
	directionTrusted := commands.GetValidSortDirections()

	isOrderByTrusted := utils.Contains(orderByTrusted, pagedRequest.OrderBy)
	isDirectionTrusted := utils.Contains(directionTrusted, pagedRequest.SortDirection)

	return isOrderByTrusted && isDirectionTrusted
}

func (repository ReceiptRepository) GetReceiptGroupIdByReceiptId(id string) (uint, error) {
	db := repository.GetDB()
	var receipt models.Receipt

	err := db.Model(models.Receipt{}).Where("id = ?", id).Select("group_id").Find(&receipt).Error
	if err != nil {
		return 0, err
	}

	return receipt.GroupId, nil
}

// GetReceiptForAuthorization loads only the fields needed to authorize a receipt
// read (id, group_id, paid_by_user_id). It uses First, so a missing row returns
// gorm.ErrRecordNotFound — letting callers authorize (and detect not-found)
// before paying to preload the full receipt's associations.
func (repository ReceiptRepository) GetReceiptForAuthorization(id string) (models.Receipt, error) {
	db := repository.GetDB()
	var receipt models.Receipt

	err := db.Model(models.Receipt{}).
		Where("id = ?", id).
		Select("id", "group_id", "paid_by_user_id").
		First(&receipt).Error
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

func (repository ReceiptRepository) FilterLinkedItemsFromReceiptItems(receipt *models.Receipt) {
	if len(receipt.ReceiptItems) == 0 {
		return
	}

	// Collect all linked item IDs
	linkedItemIds := make(map[uint]bool)
	for _, item := range receipt.ReceiptItems {
		for _, linkedItem := range item.LinkedItems {
			linkedItemIds[linkedItem.ID] = true
		}
	}

	// Filter out linked items from ReceiptItems
	var filteredItems []models.Item
	for _, item := range receipt.ReceiptItems {
		if !linkedItemIds[item.ID] {
			filteredItems = append(filteredItems, item)
		}
	}

	receipt.ReceiptItems = filteredItems
}

func (repository ReceiptRepository) GetFullyLoadedReceiptById(id string) (models.Receipt, error) {
	db := repository.GetDB()
	var receipt models.Receipt

	query := db.Model(models.Receipt{}).Where("id = ?", id).Preload(clause.Associations)

	for _, association := range constants.FULL_RECEIPT_ASSOCIATIONS {
		query = query.Preload(association)
	}

	err := query.Find(&receipt).Error
	if err != nil {
		return models.Receipt{}, err
	}

	repository.FilterLinkedItemsFromReceiptItems(&receipt)

	return receipt, nil
}

func (repository ReceiptRepository) GetReceiptsByGroupIds(groupIds []string, querySelect string, queryPreload string) ([]models.Receipt, error) {
	db := repository.GetDB()
	var receipts []models.Receipt

	query := db.Model(models.Receipt{}).Where("group_id IN ?", groupIds).Select(querySelect)
	if len(queryPreload) > 0 {
		query = query.Preload(queryPreload)
	}

	err := query.Find(&receipts).Error
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

// SearchReceiptsByGroupIds returns receipts within the given groups whose name
// matches nameQuery (a substring match; an empty nameQuery matches all),
// ordered by most recent date first and capped at limit. Scoping to the
// caller's group ids is the caller's responsibility. The paidByResolver applies
// the caller's paid-by visibility in SQL BEFORE the limit, so hidden receipts
// can't push visible matches out of the capped result set.
func (repository ReceiptRepository) SearchReceiptsByGroupIds(groupIds []uint, nameQuery string, limit int, paidByResolver PaidByAllowedResolver) ([]models.Receipt, error) {
	db := repository.GetDB()
	var receipts []models.Receipt

	query := db.Model(models.Receipt{}).Where("group_id IN ?", groupIds)
	if len(nameQuery) > 0 {
		query = query.Where("name LIKE ?", "%"+nameQuery+"%")
	}

	query, err := repository.ApplyPaidByDisjunction(query, groupIds, paidByResolver)
	if err != nil {
		return nil, err
	}

	err = query.Order("date desc").Limit(limit).Find(&receipts).Error
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func (repository ReceiptRepository) GetReceiptsByIds(ids []string, associations []string) ([]models.Receipt, error) {
	query := repository.GetDB().Model(models.Receipt{}).Where("id IN ?", ids)

	hasLinkedItems := false
	if associations != nil {
		for _, association := range associations {
			query = query.Preload(association)
			if association == "ReceiptItems.LinkedItems" {
				hasLinkedItems = true
			}
		}
	}

	var receipts []models.Receipt
	err := query.Find(&receipts).Error
	if err != nil {
		return nil, err
	}

	// Filter linked items if they were loaded
	if hasLinkedItems {
		for i := range receipts {
			repository.FilterLinkedItemsFromReceiptItems(&receipts[i])
		}
	}

	return receipts, nil
}
