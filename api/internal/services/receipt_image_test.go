package services

import (
	"net/http"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"runtime"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

// testApiRoot returns the absolute path of the api/ directory. See the
// identical helper in repositories/files_test.go for rationale.
func testApiRoot() string {
	_, file, _, _ := runtime.Caller(0)
	// file = .../api/internal/services/receipt_image_test.go
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// jpgToPdf converts a JPG blob to a one-page PDF using ImageMagick. Used by
// the PDF-branch receipt-image test.
func jpgToPdf(t *testing.T, jpgBytes []byte) []byte {
	t.Helper()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImageBlob(jpgBytes); err != nil {
		t.Fatalf("ImageMagick ReadImageBlob: %v", err)
	}
	if err := mw.SetImageFormat("pdf"); err != nil {
		t.Fatalf("ImageMagick SetImageFormat pdf: %v", err)
	}
	out, err := mw.GetImageBlob()
	if err != nil {
		t.Fatalf("ImageMagick GetImageBlob: %v", err)
	}
	return out
}

// seedReceiptImagePipeline builds the full receipt-image pipeline graph:
//
//   - A User + a Group (non-All), with the user as a member.
//   - A Prompt + ReceiptProcessingSettings (vision model) pointed at the
//     given mock server URL.
//   - A SystemSettings row linking to those processing settings.
//   - A Receipt + FileData for that receipt, with the JPG fixture written
//     to the path BuildFilePath will construct.
//
// Returns the user, the group, and the FileData (for
// GetReceiptFromReceiptImageId / ReadReceiptImage invocations).
func seedReceiptImagePipeline(t *testing.T, url string) (models.User, models.Group, models.FileData) {
	t.Helper()
	t.Setenv("BASE_PATH", testApiRoot())
	db := repositories.GetDB()

	user := models.User{Username: "ri-user", Password: "p", DisplayName: "x"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	group := models.Group{Name: "ri-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	member := models.GroupMember{GroupID: group.ID, UserID: user.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}
	groupSettings := models.GroupSettings{GroupId: group.ID}
	if err := db.Create(&groupSettings).Error; err != nil {
		t.Fatalf("seed group settings: %v", err)
	}

	prompt := models.Prompt{Name: "ri-prompt", Prompt: "Extract: @ocrText"}
	if err := db.Create(&prompt).Error; err != nil {
		t.Fatalf("seed prompt: %v", err)
	}

	settings := models.ReceiptProcessingSettings{
		Name:          "ri-rps",
		AiType:        models.OLLAMA,
		Url:           url,
		Model:         "m",
		IsVisionModel: true,
		PromptId:      prompt.ID,
	}
	if err := db.Create(&settings).Error; err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	sysSettings := models.SystemSettings{
		BaseModel:                   models.BaseModel{ID: 1},
		ReceiptProcessingSettingsId: &settings.ID,
	}
	if err := db.Create(&sysSettings).Error; err != nil {
		t.Fatalf("seed system settings: %v", err)
	}

	receipt := models.Receipt{Name: "r", GroupId: group.ID, PaidByUserID: user.ID}
	if err := db.Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}

	fileData := models.FileData{Name: "img.jpg", ReceiptId: receipt.ID, FileType: "image/jpeg"}
	if err := db.Create(&fileData).Error; err != nil {
		t.Fatalf("seed file data: %v", err)
	}

	// Write the jpg to the path BuildFilePath will resolve.
	fileRepo := repositories.NewFileRepository(nil)
	path, err := fileRepo.BuildFilePath(utils.UintToString(receipt.ID), utils.UintToString(fileData.ID), fileData.Name)
	if err != nil {
		t.Fatalf("BuildFilePath: %v", err)
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	jpg, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(path, jpg, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })

	return user, group, fileData
}

// ---------- ReadReceiptImage ----------

func TestReadReceiptImage_HappyPath(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("FromImage", "10"))
	_, _, fileData := seedReceiptImagePipeline(t, server.URL)

	receipt, _, err := ReadReceiptImage(utils.UintToString(fileData.ID))
	if err != nil {
		t.Fatalf("ReadReceiptImage: %v", err)
	}
	if receipt.Name != "FromImage" {
		t.Errorf("expected 'FromImage', got %q", receipt.Name)
	}
}

func TestReadReceiptImage_UnknownId(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("x", "1"))
	seedReceiptImagePipeline(t, server.URL)

	_, _, err := ReadReceiptImage("999999")
	if err == nil {
		t.Fatal("expected not-found error")
	}
}

// PDF image path → ConvertPdfToJpg runs, temp file written + removed.
func TestReadReceiptImage_PdfRoutesThroughConversion(t *testing.T) {
	defer repositories.TruncateTestDb()
	t.Setenv("BASE_PATH", testApiRoot())

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("PdfReceipt", "99"))
	_, _, fileData := seedReceiptImagePipeline(t, server.URL)

	// Replace the on-disk image with a PDF and update the FileData row to
	// match.
	db := repositories.GetDB()
	fileRepo := repositories.NewFileRepository(nil)
	path, err := fileRepo.BuildFilePath(utils.UintToString(fileData.ReceiptId), utils.UintToString(fileData.ID), fileData.Name)
	if err != nil {
		t.Fatalf("BuildFilePath: %v", err)
	}

	// Read original JPG, convert to PDF, overwrite target.
	jpg, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seeded jpg: %v", err)
	}
	// Reuse the files_test helper via a direct ImageMagick conversion.
	fileRepoImpl := repositories.NewFileRepository(nil)
	_ = fileRepoImpl // reserved for readability

	pdfBytes := jpgToPdf(t, jpg)
	if err := os.WriteFile(path, pdfBytes, 0o644); err != nil {
		t.Fatalf("overwrite with pdf: %v", err)
	}
	// Update file type so ReadReceiptImage takes the PDF branch.
	if err := db.Model(&models.FileData{}).Where("id = ?", fileData.ID).Update("file_type", constants.ApplicationPdf).Error; err != nil {
		t.Fatalf("update file type: %v", err)
	}

	receipt, _, err := ReadReceiptImage(utils.UintToString(fileData.ID))
	if err != nil {
		t.Fatalf("ReadReceiptImage: %v", err)
	}
	if receipt.Name != "PdfReceipt" {
		t.Errorf("expected 'PdfReceipt', got %q", receipt.Name)
	}
}

// ---------- ReadReceiptImageFromFileOnly / WithEmailBody / ImagesWithEmailBody ----------

func TestReadReceiptImageFromFileOnly(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("FileOnly", "3"))
	_, group, _ := seedReceiptImagePipeline(t, server.URL)
	imgPath := setupImagePathFixture(t)

	receipt, _, err := ReadReceiptImageFromFileOnly(imgPath, utils.UintToString(group.ID))
	if err != nil {
		t.Fatalf("ReadReceiptImageFromFileOnly: %v", err)
	}
	if receipt.Name != "FileOnly" {
		t.Errorf("expected 'FileOnly', got %q", receipt.Name)
	}
}

func TestReadReceiptImageWithEmailBody(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("WithBody", "4"))
	_, group, _ := seedReceiptImagePipeline(t, server.URL)
	imgPath := setupImagePathFixture(t)

	receipt, _, err := ReadReceiptImageWithEmailBody(imgPath, "email body", utils.UintToString(group.ID))
	if err != nil {
		t.Fatalf("ReadReceiptImageWithEmailBody: %v", err)
	}
	if receipt.Name != "WithBody" {
		t.Errorf("expected 'WithBody', got %q", receipt.Name)
	}
}

func TestReadReceiptImagesWithEmailBody_Multi(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("MultiBody", "5"))
	_, group, _ := seedReceiptImagePipeline(t, server.URL)
	imgPath := setupImagePathFixture(t)

	receipt, _, err := ReadReceiptImagesWithEmailBody([]string{imgPath, imgPath}, "body", true, utils.UintToString(group.ID))
	if err != nil {
		t.Fatalf("ReadReceiptImagesWithEmailBody: %v", err)
	}
	if receipt.Name != "MultiBody" {
		t.Errorf("expected 'MultiBody', got %q", receipt.Name)
	}
}

func TestReadReceiptFromTextOnly(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("TextOnly", "6"))
	_, group, _ := seedReceiptImagePipeline(t, server.URL)

	receipt, _, err := ReadReceiptFromTextOnly("the full email body", utils.UintToString(group.ID))
	if err != nil {
		t.Fatalf("ReadReceiptFromTextOnly: %v", err)
	}
	if receipt.Name != "TextOnly" {
		t.Errorf("expected 'TextOnly', got %q", receipt.Name)
	}
}

func TestReadReceiptFromTextOnly_MissingSettings(t *testing.T) {
	defer repositories.TruncateTestDb()

	// No seeds → NewSystemReceiptProcessingService returns an error.
	_, _, err := ReadReceiptFromTextOnly("body", "")
	if err == nil {
		t.Fatal("expected error when system settings missing")
	}
}

// ---------- MagicFillFromImage ----------

func TestMagicFillFromImage_HappyPath(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("Magic", "7"))
	_, group, _ := seedReceiptImagePipeline(t, server.URL)

	jpg, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	cmd := commands.MagicFillCommand{ImageData: jpg}
	receipt, _, err := MagicFillFromImage(cmd, utils.UintToString(group.ID), 0)
	if err != nil {
		t.Fatalf("MagicFillFromImage: %v", err)
	}
	if receipt.Name != "Magic" {
		t.Errorf("expected 'Magic', got %q", receipt.Name)
	}
}

func TestMagicFillFromImage_InvalidImageData(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("x", "1"))
	_, group, _ := seedReceiptImagePipeline(t, server.URL)

	cmd := commands.MagicFillCommand{ImageData: []byte("not an image")}
	_, _, err := MagicFillFromImage(cmd, utils.UintToString(group.ID), 0)
	if err == nil {
		t.Fatal("expected invalid file type error")
	}
}

// ---------- GetReceiptImagesForGroup / GetReceiptFromReceiptImageId ----------

func TestGetReceiptImagesForGroup_FiltersByGroup(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("x", "1"))
	user, group, _ := seedReceiptImagePipeline(t, server.URL)

	// Add a second group + receipt + FileData that should NOT appear.
	db := repositories.GetDB()
	otherGroup := models.Group{Name: "other-group"}
	if err := db.Create(&otherGroup).Error; err != nil {
		t.Fatalf("seed other group: %v", err)
	}
	otherReceipt := models.Receipt{Name: "o", GroupId: otherGroup.ID, PaidByUserID: user.ID}
	if err := db.Create(&otherReceipt).Error; err != nil {
		t.Fatalf("seed other receipt: %v", err)
	}
	if err := db.Create(&models.FileData{Name: "other.jpg", ReceiptId: otherReceipt.ID, FileType: "image/jpeg"}).Error; err != nil {
		t.Fatalf("seed other FileData: %v", err)
	}

	got, err := GetReceiptImagesForGroup(utils.UintToString(group.ID), utils.UintToString(user.ID))
	if err != nil {
		t.Fatalf("GetReceiptImagesForGroup: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 FileData for the group, got %d", len(got))
	}
	if len(got) > 0 && got[0].Name != "img.jpg" {
		t.Errorf("expected img.jpg, got %q", got[0].Name)
	}
}

func TestGetReceiptImagesForGroup_UnknownGroup(t *testing.T) {
	defer repositories.TruncateTestDb()

	_, err := GetReceiptImagesForGroup("999999", "1")
	if err == nil {
		t.Fatal("expected group-not-found error")
	}
}

func TestGetReceiptFromReceiptImageId_Found(t *testing.T) {
	defer repositories.TruncateTestDb()

	server, _ := newMockOllamaServerForService(t, http.StatusOK, ollamaReceiptJson("x", "1"))
	_, _, fileData := seedReceiptImagePipeline(t, server.URL)

	got, err := GetReceiptFromReceiptImageId(utils.UintToString(fileData.ID))
	if err != nil {
		t.Fatalf("GetReceiptFromReceiptImageId: %v", err)
	}
	if got.ID != fileData.ReceiptId {
		t.Errorf("expected receipt id %d, got %d", fileData.ReceiptId, got.ID)
	}
}

func TestGetReceiptFromReceiptImageId_MissingImage(t *testing.T) {
	defer repositories.TruncateTestDb()

	_, err := GetReceiptFromReceiptImageId("999999")
	if err == nil {
		t.Fatal("expected FileData-not-found error")
	}
}
