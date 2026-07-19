package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"
)

func TestFileDataView_FromFileData(t *testing.T) {
	createdBy := uint(7)
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	fileData := FileData{
		BaseModel: BaseModel{
			ID:              42,
			CreatedAt:       createdAt,
			UpdatedAt:       updatedAt,
			CreatedBy:       &createdBy,
			CreatedByString: "tester",
		},
		Name:      "receipt.jpg",
		FileType:  "image/jpeg",
		Size:      1024,
		ReceiptId: 99,
	}

	view := FileDataView{}.FromFileData(fileData)

	if view.ID != fileData.ID {
		utils.PrintTestError(t, view.ID, fileData.ID)
	}
	if !view.CreatedAt.Equal(fileData.CreatedAt) {
		utils.PrintTestError(t, view.CreatedAt, fileData.CreatedAt)
	}
	if !view.UpdatedAt.Equal(fileData.UpdatedAt) {
		utils.PrintTestError(t, view.UpdatedAt, fileData.UpdatedAt)
	}
	if view.CreatedBy != fileData.CreatedBy {
		utils.PrintTestError(t, view.CreatedBy, fileData.CreatedBy)
	}
	if view.CreatedByString != fileData.CreatedByString {
		utils.PrintTestError(t, view.CreatedByString, fileData.CreatedByString)
	}
	if view.Name != fileData.Name {
		utils.PrintTestError(t, view.Name, fileData.Name)
	}
	// FromFileData never carries the encoded image.
	if view.EncodedImage != "" {
		utils.PrintTestError(t, view.EncodedImage, "")
	}
}
