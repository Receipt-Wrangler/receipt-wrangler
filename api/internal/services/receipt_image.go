package services

import (
	"os"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
)

func ReadReceiptImage(receiptImageId string) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	var result commands.UpsertReceiptCommand
	var pathToReadFrom string
	receiptService := NewReceiptService(nil)

	receipt, err := receiptService.GetReceiptByReceiptImageId(receiptImageId)
	if err != nil {
		return result, commands.ReceiptProcessingMetadata{}, err
	}

	groupIdString := utils.UintToString(receipt.GroupId)

	systemReceiptProcessingService, err := NewSystemReceiptProcessingService(nil, groupIdString)
	if err != nil {
		return result, commands.ReceiptProcessingMetadata{}, err
	}

	receiptImageUint, err := utils.StringToUint(receiptImageId)
	if err != nil {
		return result, commands.ReceiptProcessingMetadata{}, err
	}

	receiptImageRepository := repositories.NewReceiptImageRepository(nil)
	receiptImage, err := receiptImageRepository.GetReceiptImageById(receiptImageUint)
	if err != nil {
		return result, commands.ReceiptProcessingMetadata{}, err
	}
	fileRepository := repositories.NewFileRepository(nil)

	receiptImagePath, err := fileRepository.BuildFilePath(utils.UintToString(receiptImage.ReceiptId), receiptImageId, receiptImage.Name)
	if err != nil {
		return result, commands.ReceiptProcessingMetadata{}, err
	}

	receiptImageBytes, err := utils.ReadDataFile(receiptImagePath)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}

	// TODO: make generic
	if receiptImage.FileType == constants.ApplicationPdf {
		bytes, err := fileRepository.ConvertPdfToJpg(receiptImageBytes)
		if err != nil {
			return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
		}

		pathToReadFrom, err = fileRepository.WriteTempFile(bytes)
		if err != nil {
			return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
		}

		defer os.Remove(pathToReadFrom)
	} else {
		pathToReadFrom = receiptImagePath
	}

	return systemReceiptProcessingService.ReadReceiptImage(pathToReadFrom)
}

func ReadReceiptImageFromFileOnly(path string, groupId string) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	receiptProcessingService, err := NewSystemReceiptProcessingService(nil, groupId)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}

	return receiptProcessingService.ReadReceiptImage(path)
}

func ReadReceiptImageWithEmailBody(path string, emailBody string, groupId string) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	receiptProcessingService, err := NewSystemReceiptProcessingService(nil, groupId)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}

	return receiptProcessingService.ReadReceiptImageWithEmailBody(path, emailBody)
}

func ReadReceiptImagesWithEmailBody(paths []string, emailBody string, bodySentAsImage bool, groupId string) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	receiptProcessingService, err := NewSystemReceiptProcessingService(nil, groupId)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}

	return receiptProcessingService.ReadReceiptImagesWithEmailBody(paths, emailBody, bodySentAsImage)
}

func ReadReceiptFromTextOnly(bodyText string, groupId string) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	receiptProcessingService, err := NewSystemReceiptProcessingService(nil, groupId)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}

	return receiptProcessingService.ReadReceiptText(bodyText)
}

func MagicFillFromImage(command commands.MagicFillCommand, groupId string, userId uint) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	fileRepository := repositories.NewFileRepository(nil)
	receiptProcessingService, err := NewSystemReceiptProcessingService(nil, groupId)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}
	// Restrict the AI prompt's candidate categories/tags to this user's grants
	// (0 when there is no triggering user, e.g. system processing).
	receiptProcessingService.UserId = userId

	bytes, err := fileRepository.GetBytesFromImageBytes(command.ImageData)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}

	filePath, err := fileRepository.WriteTempFile(bytes)
	if err != nil {
		return commands.UpsertReceiptCommand{}, commands.ReceiptProcessingMetadata{}, err
	}
	defer os.Remove(filePath)

	return receiptProcessingService.ReadReceiptImage(filePath)
}

func GetReceiptImagesForGroup(groupId string, userId string) ([]models.FileData, error) {
	db := repositories.GetDB()
	groupRepository := repositories.NewGroupRepository(nil)
	groupService := NewGroupService(nil)
	groupIds := make([]uint, 0)

	group, err := groupRepository.GetGroupById(groupId, false, true, false)
	if err != nil {
		return nil, err
	}

	if group.IsAllGroup {
		groups, err := groupService.GetGroupsForUser(userId)
		if err != nil {
			return nil, err
		}

		for _, group := range groups {
			groupIds = append(groupIds, group.ID)
		}
	} else {
		uintGroupId, err := utils.StringToUint(groupId)
		if err != nil {
			return nil, err
		}

		groupIds = append(groupIds, uintGroupId)
	}

	fileDataResults := make([]models.FileData, 0)
	err = db.Table("receipts").Select("receipts.id, receipts.group_id, file_data.*").Joins("inner join file_data on file_data.receipt_id=receipts.id").Where("receipts.group_id IN ?", groupIds).Scan(&fileDataResults).Error
	if err != nil {
		return nil, err
	}

	return fileDataResults, nil
}

func GetReceiptFromReceiptImageId(receiptImageId string) (models.Receipt, error) {
	db := repositories.GetDB()
	var receipt models.Receipt
	var fileData models.FileData

	err := db.Model(models.FileData{}).Where("id = ?", receiptImageId).Select("receipt_id").First(&fileData).Error
	if err != nil {
		return models.Receipt{}, err
	}

	receiptIdString := utils.UintToString(fileData.ReceiptId)

	receiptRepository := repositories.NewReceiptRepository(nil)
	receipt, err = receiptRepository.GetReceiptById(receiptIdString)
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}
