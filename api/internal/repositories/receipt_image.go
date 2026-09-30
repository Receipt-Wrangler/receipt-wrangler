package repositories

import (
	"path/filepath"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReceiptImageRepository struct {
	BaseRepository
}

// writeReceiptImageFile is the write CreateReceiptImage makes. It is a variable
// only so a test can make it fail partway through a file (see
// SetReceiptImageWriterForTests); nothing else should reassign it.
var writeReceiptImageFile = utils.WriteDataFile

// SetReceiptImageWriterForTests swaps the file write CreateReceiptImage makes and
// returns a function restoring the real one. Tests only: it exists so a caller's
// clean-up of a partially written image can be proven, which a real filesystem
// will not produce on demand.
func SetReceiptImageWriterForTests(write func(path string, data []byte) error) (restore func()) {
	previous := writeReceiptImageFile
	writeReceiptImageFile = write
	return func() { writeReceiptImageFile = previous }
}

func NewReceiptImageRepository(tx *gorm.DB) ReceiptImageRepository {
	repository := ReceiptImageRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

// TODO: Move to service
func (repository ReceiptImageRepository) CreateReceiptImage(fileData models.FileData, fileBytes []byte) (models.FileData, error) {
	fileRepository := NewFileRepository(repository.TX)
	db := repository.GetDB()

	// TODO: refactor to use command
	validatedFileType, err := fileRepository.ValidateFileType(fileBytes)
	if err != nil {
		return models.FileData{}, err
	}

	fileData.FileType = validatedFileType

	// Ensure the data directory exists
	dataDir, err := utils.GetDataDir()
	if err != nil {
		return models.FileData{}, err
	}
	err = utils.EnsureDataDirectory(dataDir)
	if err != nil {
		return models.FileData{}, err
	}

	// Get initial group directory to see if it exists
	filePath, err := fileRepository.BuildFilePath(utils.UintToString(fileData.ReceiptId), "", fileData.Name)
	if err != nil {
		return models.FileData{}, err
	}
	groupDir, _ := filepath.Split(filePath)

	err = db.Model(models.FileData{}).Create(&fileData).Error
	if err != nil {
		utils.RemoveDataPath(filePath)
		return models.FileData{}, err
	}

	// Check if group's path exists
	err = utils.EnsureDataDirectory(groupDir)
	if err != nil {
		return models.FileData{}, err
	}

	// Rebuild file path with correct file id
	filePath, err = fileRepository.BuildFilePath(utils.UintToString(fileData.ReceiptId), utils.UintToString(fileData.ID), fileData.Name)
	if err != nil {
		return models.FileData{}, err
	}

	// The row (and its id) is returned alongside a write error so a caller running
	// this inside its own transaction can locate, and remove, a partially written
	// file once that transaction rolls back.
	err = writeReceiptImageFile(filePath, fileBytes)
	if err != nil {
		return fileData, err
	}

	return fileData, nil
}

func (repository ReceiptImageRepository) GetReceiptImageById(receiptImageId uint) (models.FileData, error) {
	db := repository.GetDB()
	var result models.FileData

	err := db.Model(models.FileData{}).Where("id = ?", receiptImageId).Preload(clause.Associations).Find(&result).Error
	if err != nil {
		return models.FileData{}, err
	}

	return result, nil
}

func (repository ReceiptImageRepository) GetReceiptImagesByIdArray(receiptImageIds []uint) ([]models.FileData, error) {
	db := repository.GetDB()
	var result = make([]models.FileData, 0)

	err := db.Model(models.FileData{}).Where("id IN ?", receiptImageIds).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
