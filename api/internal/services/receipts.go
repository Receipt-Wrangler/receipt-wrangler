package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strconv"
	"strings"
	"time"
)

// Sentinel errors for the shared, enforced read operations. They let each ingress
// point (REST handlers, MCP tools) map a single enforcement outcome to its own
// transport without re-implementing the checks.
var (
	// ErrReceiptAccessDenied is returned when a receipt is missing, the caller
	// lacks group.receipts.read, or the receipt is hidden by paid-by visibility.
	// It is intentionally indistinct so callers don't leak a receipt's existence.
	ErrReceiptAccessDenied = errors.New("receipt access denied")
	// ErrSearchForbidden is returned when the caller lacks app.receipts.search.
	ErrSearchForbidden = errors.New("not authorized to search receipts")
)

type ReceiptService struct {
	BaseService
}

func NewReceiptService(tx *gorm.DB) ReceiptService {
	service := ReceiptService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

// GetReceiptForUser is the single, shared "read one receipt" operation used by
// both the REST handler and the MCP tool. It fetches the receipt and applies the
// full read enforcement chain — group.receipts.read, paid-by visibility, and
// category/tag grant stripping — so the two ingress points cannot drift. Missing,
// forbidden, and paid-by-hidden receipts all collapse to ErrReceiptAccessDenied so
// the caller cannot infer a receipt's existence.
func (service ReceiptService) GetReceiptForUser(userId uint, receiptId string) (models.Receipt, error) {
	receiptRepository := repositories.NewReceiptRepository(service.TX)
	permissionService := NewPermissionService(service.TX)

	// Authorize on the lightweight auth fields first (this fetch uses First, so a
	// missing row surfaces as ErrRecordNotFound). Only load the full receipt with
	// its associations once the read is allowed.
	authReceipt, err := receiptRepository.GetReceiptForAuthorization(receiptId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Receipt{}, ErrReceiptAccessDenied
		}
		return models.Receipt{}, err
	}

	hasAccess, err := permissionService.HasGroupPermissions(userId, authReceipt.GroupId, permissions.GroupReceiptsRead)
	if err != nil {
		return models.Receipt{}, err
	}
	if !hasAccess {
		return models.Receipt{}, ErrReceiptAccessDenied
	}

	visible, err := permissionService.ReceiptPaidByVisible(userId, authReceipt.GroupId, authReceipt.PaidByUserID)
	if err != nil {
		return models.Receipt{}, err
	}
	if !visible {
		return models.Receipt{}, ErrReceiptAccessDenied
	}

	receipt, err := receiptRepository.GetFullyLoadedReceiptById(receiptId)
	if err != nil {
		return models.Receipt{}, err
	}

	// Guard the window between the authorization read and this full load: if the
	// receipt's identity, group, or payer changed (or it was deleted — Find leaves
	// a zero-value row), the prior authorization no longer applies, so deny rather
	// than return data that was never authorized.
	if receipt.ID != authReceipt.ID ||
		receipt.GroupId != authReceipt.GroupId ||
		receipt.PaidByUserID != authReceipt.PaidByUserID {
		return models.Receipt{}, ErrReceiptAccessDenied
	}

	if err := permissionService.FilterReceiptCategoriesTagsForReceipt(userId, &receipt); err != nil {
		return models.Receipt{}, err
	}

	// Mask user references (created-by, charged-to) and drop non-visible comment
	// authors outside the caller's member-visible set.
	if err := permissionService.MaskReceiptForMemberVisibility(userId, &receipt); err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

// SearchReceiptsForUser is the single, shared receipt-search operation used by both
// the REST handler and the MCP tool. It enforces app.receipts.search, scopes to the
// caller's groups, applies paid-by visibility in SQL before the limit, and maps to
// SearchResult. A blank query returns no results (matching the REST search bar).
func (service ReceiptService) SearchReceiptsForUser(userId uint, query string, limit int) ([]structs.SearchResult, error) {
	permissionService := NewPermissionService(service.TX)

	hasAccess, err := permissionService.HasAppPermissions(userId, permissions.AppReceiptsSearch)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrSearchForbidden
	}

	results := make([]structs.SearchResult, 0)
	if len(strings.TrimSpace(query)) == 0 {
		return results, nil
	}

	groupMemberRepository := repositories.NewGroupMemberRepository(service.TX)
	groupIds, err := groupMemberRepository.GetGroupIdsByUserId(utils.UintToString(userId))
	if err != nil {
		return nil, err
	}

	receiptRepository := repositories.NewReceiptRepository(service.TX)
	receipts, err := receiptRepository.SearchReceiptsByGroupIds(groupIds, query, limit, permissionService.PaidByListResolver(userId))
	if err != nil {
		return nil, err
	}

	for _, receipt := range receipts {
		results = append(results, structs.SearchResult{
			ID:            receipt.ID,
			GroupID:       receipt.GroupId,
			Name:          receipt.Name,
			Date:          receipt.Date,
			Type:          "Receipt",
			Amount:        receipt.Amount,
			ReceiptStatus: receipt.Status,
			PaidByUserId:  receipt.PaidByUserID,
			CreatedAt:     receipt.CreatedAt,
		})
	}

	return results, nil
}

func (service ReceiptService) GetReceiptByReceiptImageId(receiptImageId string) (models.Receipt, error) {
	db := service.GetDB()
	var fileData models.FileData

	err := db.Model(models.FileData{}).Where("id = ?", receiptImageId).Select("receipt_id").First(&fileData).Error
	if err != nil {
		return models.Receipt{}, err
	}

	receiptRepository := repositories.NewReceiptRepository(service.TX)
	receipt, err := receiptRepository.GetReceiptById(strconv.FormatUint(uint64(fileData.ReceiptId), 10))
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

func (service ReceiptService) DeleteReceipt(id string) error {
	db := service.GetDB()
	var receipt models.Receipt
	receiptRepository := repositories.NewReceiptRepository(service.TX)

	receipt, err := receiptRepository.GetFullyLoadedReceiptById(id)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var imagesToDelete []string
		fileRepository := repositories.NewFileRepository(tx)
		fileRepository.SetTransaction(tx)

		for _, f := range receipt.ImageFiles {
			path, _ := fileRepository.BuildFilePath(utils.UintToString(f.ReceiptId), utils.UintToString(f.ID), f.Name)
			imagesToDelete = append(imagesToDelete, path)
		}

		for _, r := range receipt.ReceiptItems {
			err = tx.Model(&r).Association("Categories").Clear()
			if err != nil {
				return err
			}

			err = tx.Model(&r).Association("Tags").Clear()
			if err != nil {
				return err
			}
		}

		err = tx.Model(&receipt).Association("ReceiptItems").Clear()
		if err != nil {
			return err
		}

		err = tx.Select(clause.Associations).Delete(&receipt).Error
		if err != nil {
			return err
		}

		for _, path := range imagesToDelete {
			utils.RemoveDataPath(path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service ReceiptService) QuickScan(
	token *structs.Claims,
	paidByUserId uint,
	groupId uint,
	status models.ReceiptStatus,
	categoryIds []uint,
	tagIds []uint,
	tempPath string,
	originalFileName string,
	asynqTaskId string,
) (models.Receipt, error) {
	db := repositories.GetDB()
	systemTaskService := NewSystemTaskService(service.TX)
	var createdReceipt models.Receipt

	fileRepository := repositories.NewFileRepository(service.TX)
	fileBytes, err := utils.ReadFile(tempPath)
	if err != nil {
		return models.Receipt{}, err
	}

	fileInfo, err := os.Stat(tempPath)
	if err != nil {
		return models.Receipt{}, err
	}

	validatedFileType, err := fileRepository.ValidateFileType(fileBytes)
	if err != nil {
		return models.Receipt{}, err
	}

	magicFillCommand := commands.MagicFillCommand{
		ImageData: fileBytes,
		Filename:  originalFileName,
	}

	receiptRepository := repositories.NewReceiptRepository(service.TX)
	receiptImageRepository := repositories.NewReceiptImageRepository(service.TX)

	groupIdString := utils.UintToString(groupId)

	now := time.Now()
	receiptCommand, receiptProcessingMetadata, magicFillErr := MagicFillFromImage(magicFillCommand, groupIdString, token.UserId)
	finishedAt := time.Now()

	quickScanSystemTasks, taskErr := systemTaskService.CreateSystemTasksFromMetadata(
		receiptProcessingMetadata,
		now,
		finishedAt,
		models.QUICK_SCAN,
		&token.UserId,
		&groupId,
		asynqTaskId, nil)
	if taskErr != nil {
		return models.Receipt{}, taskErr
	}

	if magicFillErr != nil {
		return models.Receipt{}, magicFillErr
	}

	if receiptCommand.PaidByUserID == 0 {
		receiptCommand.PaidByUserID = paidByUserId
	}

	if len(receiptCommand.Status) == 0 {
		receiptCommand.Status = models.ReceiptStatus(status)
	}

	receiptCommand.GroupId = groupId

	// Merge the user's quick-scan category/tag picks with whatever the AI auto-assigned (union,
	// deduped by id). Names are resolved from the ids so the merged selections pass receipt
	// validation, which requires a category/tag name.
	receiptCommand.Categories, err = service.mergeQuickScanCategories(receiptCommand.Categories, categoryIds)
	if err != nil {
		return models.Receipt{}, err
	}

	receiptCommand.Tags, err = service.mergeQuickScanTags(receiptCommand.Tags, tagIds)
	if err != nil {
		return models.Receipt{}, err
	}

	vErr := receiptCommand.Validate(token.UserId, true)
	if len(vErr.Errors) > 0 {
		errBytes, _ := json.Marshal(vErr.Errors)
		return models.Receipt{}, fmt.Errorf("receipt validation failed: %s", string(errBytes))
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		receiptRepository.SetTransaction(tx)
		receiptImageRepository.SetTransaction(tx)
		systemTaskService.SetTransaction(tx)
		uploadStart := time.Now()

		createdReceipt, err = receiptRepository.CreateReceipt(receiptCommand, token.UserId, false)
		_, taskErr := systemTaskService.CreateReceiptUploadedSystemTask(
			err,
			createdReceipt,
			quickScanSystemTasks,
			uploadStart,
		)
		if taskErr != nil {
			return taskErr
		}
		if err != nil {
			tx.Commit()
			return err
		}

		taskErr = systemTaskService.AssociateProcessingSystemTasksToReceipt(quickScanSystemTasks, createdReceipt.ID)
		if taskErr != nil {
			return taskErr
		}

		fileData := models.FileData{
			Name:      originalFileName,
			Size:      uint(fileInfo.Size()),
			ReceiptId: createdReceipt.ID,
			FileType:  validatedFileType,
		}
		_, err := receiptImageRepository.CreateReceiptImage(fileData, fileBytes)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Receipt{}, err
	}

	os.Remove(tempPath)
	return createdReceipt, nil
}

// mergeQuickScanCategories appends the user-selected category ids to the AI-filled categories,
// skipping any already present. Names are loaded from the ids so the result passes receipt
// validation.
func (service ReceiptService) mergeQuickScanCategories(
	existing []commands.UpsertCategoryCommand,
	categoryIds []uint,
) ([]commands.UpsertCategoryCommand, error) {
	if len(categoryIds) == 0 {
		return existing, nil
	}

	categoryRepository := repositories.NewCategoryRepository(service.TX)
	categories, err := categoryRepository.GetByIds(categoryIds)
	if err != nil {
		return nil, err
	}

	presentIds := make(map[uint]bool)
	for _, category := range existing {
		if category.Id != nil {
			presentIds[*category.Id] = true
		}
	}

	for _, category := range categories {
		if presentIds[category.ID] {
			continue
		}

		id := category.ID
		existing = append(existing, commands.UpsertCategoryCommand{
			Id:          &id,
			Name:        category.Name,
			Description: category.Description,
		})
		presentIds[id] = true
	}

	return existing, nil
}

// mergeQuickScanTags is the tag counterpart of mergeQuickScanCategories.
func (service ReceiptService) mergeQuickScanTags(
	existing []commands.UpsertTagCommand,
	tagIds []uint,
) ([]commands.UpsertTagCommand, error) {
	if len(tagIds) == 0 {
		return existing, nil
	}

	tagsRepository := repositories.NewTagsRepository(service.TX)
	tags, err := tagsRepository.GetByIds(tagIds)
	if err != nil {
		return nil, err
	}

	presentIds := make(map[uint]bool)
	for _, tag := range existing {
		if tag.Id != nil {
			presentIds[*tag.Id] = true
		}
	}

	for _, tag := range tags {
		if presentIds[tag.ID] {
			continue
		}

		id := tag.ID
		existing = append(existing, commands.UpsertTagCommand{
			Id:          &id,
			Name:        tag.Name,
			Description: tag.Description,
		})
		presentIds[id] = true
	}

	return existing, nil
}

func (service ReceiptService) DuplicateReceipt(
	userId uint,
	receiptId string,
) (models.Receipt, error) {
	db := repositories.GetDB()
	newReceipt := models.Receipt{}

	systemTaskCommand := commands.UpsertSystemTaskCommand{
		Type:                 models.RECEIPT_UPLOADED,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.RECEIPT,
		AssociatedEntityId:   0,
		StartedAt:            time.Now(),
		EndedAt:              nil,
		ResultDescription:    "",
		RanByUserId:          &userId,
		ReceiptId:            nil,
		GroupId:              nil,
	}

	receiptRepository := repositories.NewReceiptRepository(nil)
	receipt, err := receiptRepository.GetFullyLoadedReceiptById(receiptId)
	defer func() {
		systemTaskService := NewSystemTaskService(nil)
		systemTaskService.CreateSystemTaskFromError(systemTaskCommand, err)
	}()
	if err != nil {
		return models.Receipt{}, err
	}

	systemTaskCommand.GroupId = &receipt.GroupId

	// Strip categories/tags the duplicating user cannot see so they are not
	// copied onto the new receipt.
	permissionService := NewPermissionService(nil)
	err = permissionService.FilterReceiptCategoriesTagsForReceipt(userId, &receipt)
	if err != nil {
		return models.Receipt{}, err
	}

	// Mask user references outside the caller's member-visible set (e.g. an item
	// charged to a non-visible user) so they are not carried onto the copy.
	err = permissionService.MaskReceiptForMemberVisibility(userId, &receipt)
	if err != nil {
		return models.Receipt{}, err
	}

	copier.Copy(&newReceipt, receipt)

	newReceipt.ID = 0
	newReceipt.Name = newReceipt.Name + " duplicate"
	newReceipt.ImageFiles = make([]models.FileData, 0)
	newReceipt.ReceiptItems = make([]models.Item, 0)
	newReceipt.Comments = make([]models.Comment, 0)
	newReceipt.CreatedAt = time.Now()
	newReceipt.UpdatedAt = time.Now()
	newReceipt.CreatedBy = &userId

	// Remove fks from any related data
	for _, fileData := range receipt.ImageFiles {
		var newFileData models.FileData
		copier.Copy(&newFileData, fileData)

		newFileData.ID = 0
		newFileData.ReceiptId = 0
		newFileData.Receipt = models.Receipt{}
		newReceipt.ImageFiles = append(newReceipt.ImageFiles, newFileData)
	}

	// Copy items
	for _, item := range receipt.ReceiptItems {
		var newItem models.Item
		copier.Copy(&newItem, item)

		newItem.ID = 0
		newItem.ReceiptId = 0
		newItem.Receipt = models.Receipt{}
		newReceipt.ReceiptItems = append(newReceipt.ReceiptItems, newItem)
	}

	// Copy comments
	for _, comment := range receipt.Comments {
		var newComment models.Comment
		copier.Copy(&newComment, comment)

		newComment.ID = 0
		newComment.ReceiptId = 0
		newComment.Receipt = models.Receipt{}
		newReceipt.Comments = append(newReceipt.Comments, newComment)
	}

	err = db.Create(&newReceipt).Error
	if err != nil {
		return models.Receipt{}, err
	}
	systemTaskCommand.AssociatedEntityId = newReceipt.ID
	systemTaskCommand.ReceiptId = &newReceipt.ID

	resultString, err := newReceipt.ToString()
	if err != nil {
		return models.Receipt{}, err
	}

	systemTaskCommand.ResultDescription = resultString

	// Copy receipt images
	fileRepository := repositories.NewFileRepository(nil)
	for i, fileData := range newReceipt.ImageFiles {
		srcFileData := receipt.ImageFiles[i]
		srcImageBytes, err := fileRepository.GetBytesForFileData(srcFileData)
		if err != nil {
			return models.Receipt{}, err
		}

		dstPath, err := fileRepository.BuildFilePath(
			utils.UintToString(newReceipt.ID),
			utils.UintToString(fileData.ID),
			fileData.Name,
		)
		if err != nil {
			return models.Receipt{}, err
		}

		err = utils.WriteDataFile(dstPath, srcImageBytes)
		if err != nil {
			return models.Receipt{}, err
		}
	}

	return newReceipt, nil
}
