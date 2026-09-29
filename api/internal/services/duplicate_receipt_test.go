package services

import (
	"bytes"
	"os"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// seedDuplicateSource creates a user, a group they belong to, and a receipt in it
// with one FileData per image name. It writes no image files; each test decides
// which ones are on disk. The group's data directory is removed afterwards.
func seedDuplicateSource(t *testing.T, imageNames ...string) (models.User, models.Receipt, []models.FileData, string) {
	t.Helper()
	db := repositories.GetDB()

	user := models.User{Username: "dup-user", Password: "p", DisplayName: "Dup"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	group := models.Group{Name: "dup-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: user.ID}).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}

	receipt := models.Receipt{Name: "Source", GroupId: group.ID, PaidByUserID: user.ID}
	if err := db.Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}

	images := make([]models.FileData, 0, len(imageNames))
	for _, name := range imageNames {
		image := models.FileData{Name: name, ReceiptId: receipt.ID, FileType: "application/pdf"}
		if err := db.Create(&image).Error; err != nil {
			t.Fatalf("seed file data: %v", err)
		}
		images = append(images, image)
	}

	groupPath, err := repositories.NewFileRepository(nil).BuildGroupPath(group.ID, "")
	if err != nil {
		t.Fatalf("BuildGroupPath: %v", err)
	}
	if err := os.MkdirAll(groupPath, 0o755); err != nil {
		t.Fatalf("create group dir: %v", err)
	}
	t.Cleanup(func() { utils.RemoveAllInDataDir(groupPath) })

	return user, receipt, images, groupPath
}

func imagePath(t *testing.T, image models.FileData) string {
	t.Helper()
	path, err := repositories.NewFileRepository(nil).BuildFilePath(
		utils.UintToString(image.ReceiptId),
		utils.UintToString(image.ID),
		image.Name,
	)
	if err != nil {
		t.Fatalf("BuildFilePath: %v", err)
	}
	return path
}

func receiptUploadedTasks(t *testing.T) []models.SystemTask {
	t.Helper()
	var tasks []models.SystemTask
	if err := repositories.GetDB().Where("type = ?", models.RECEIPT_UPLOADED).Find(&tasks).Error; err != nil {
		t.Fatalf("load tasks: %v", err)
	}
	return tasks
}

// The copy is the stored bytes, not the display conversion: a PDF must stay the
// PDF its name and file type say it is. These bytes cannot be rasterized, so a
// copy that converts fails here rather than passing by accident.
func TestDuplicateReceipt_CopiesImageBytesUnconverted(t *testing.T) {
	defer repositories.TruncateTestDb()

	user, receipt, images, _ := seedDuplicateSource(t, "scan.pdf")
	source := []byte("%PDF-1.4 not really a pdf")
	if err := os.WriteFile(imagePath(t, images[0]), source, 0o644); err != nil {
		t.Fatalf("write source image: %v", err)
	}

	duplicate, err := NewReceiptService(nil).DuplicateReceipt(user.ID, utils.UintToString(receipt.ID))
	if err != nil {
		t.Fatalf("DuplicateReceipt: %v", err)
	}

	if len(duplicate.ImageFiles) != 1 {
		t.Fatalf("expected 1 image on the duplicate, got %d", len(duplicate.ImageFiles))
	}
	copied, err := os.ReadFile(imagePath(t, duplicate.ImageFiles[0]))
	if err != nil {
		t.Fatalf("read copied image: %v", err)
	}
	if !bytes.Equal(copied, source) {
		t.Errorf("copied image = %q, want the source bytes %q", copied, source)
	}

	tasks := receiptUploadedTasks(t)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 RECEIPT_UPLOADED task, got %d", len(tasks))
	}
	if tasks[0].Status != models.SYSTEM_TASK_SUCCEEDED {
		t.Errorf("task status = %v, want %v", tasks[0].Status, models.SYSTEM_TASK_SUCCEEDED)
	}
	if tasks[0].ReceiptId == nil || *tasks[0].ReceiptId != duplicate.ID {
		t.Errorf("task receipt id = %v, want %d", tasks[0].ReceiptId, duplicate.ID)
	}
}

// A duplicate whose second image cannot be copied leaves nothing behind — no
// receipt, no FileData, not even the first image it did copy — and its task says
// it failed.
func TestDuplicateReceipt_RollsBackAndRecordsFailureWhenAnImageIsMissing(t *testing.T) {
	defer repositories.TruncateTestDb()

	user, receipt, images, groupPath := seedDuplicateSource(t, "first.pdf", "second.pdf")
	if err := os.WriteFile(imagePath(t, images[0]), []byte("first image"), 0o644); err != nil {
		t.Fatalf("write source image: %v", err)
	}

	_, duplicateErr := NewReceiptService(nil).DuplicateReceipt(user.ID, utils.UintToString(receipt.ID))
	if duplicateErr == nil {
		t.Fatal("expected DuplicateReceipt to fail on the missing image")
	}

	db := repositories.GetDB()
	var receiptCount, fileDataCount int64
	db.Model(&models.Receipt{}).Count(&receiptCount)
	db.Model(&models.FileData{}).Count(&fileDataCount)
	if receiptCount != 1 || fileDataCount != 2 {
		t.Errorf("expected only the source's receipt and 2 images, got %d receipts and %d images", receiptCount, fileDataCount)
	}

	entries, err := os.ReadDir(groupPath)
	if err != nil {
		t.Fatalf("read group dir: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected only the source image on disk, got %d files", len(entries))
	}

	tasks := receiptUploadedTasks(t)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 RECEIPT_UPLOADED task, got %d", len(tasks))
	}
	task := tasks[0]
	if task.Status != models.SYSTEM_TASK_FAILED {
		t.Errorf("task status = %v, want %v", task.Status, models.SYSTEM_TASK_FAILED)
	}
	if task.ResultDescription != duplicateErr.Error() {
		t.Errorf("task description = %q, want the error %q", task.ResultDescription, duplicateErr.Error())
	}
	if task.ReceiptId != nil {
		t.Errorf("task receipt id = %d, want none for a rolled-back duplicate", *task.ReceiptId)
	}
	if task.RanByUserId == nil || *task.RanByUserId != user.ID {
		t.Errorf("task ran by = %v, want %d", task.RanByUserId, user.ID)
	}
}
