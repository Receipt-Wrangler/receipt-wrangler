package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
)

func uploadImageRequest(t *testing.T, userId uint, receiptId uint, name string, fileBytes []byte) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("receiptId", utils.UintToString(receiptId))
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatalf("create file part: %v", err)
	}
	part.Write(fileBytes)
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/receiptImage/", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

func TestUploadReceiptImageRecordsAReceiptChange(t *testing.T) {
	defer repositories.TruncateTestDb()
	defer services.ClearRolePermissionCacheForTests()
	userId, groupId := seedGroupWithRole(t, "image-upload-recorded", []string{permissions.GroupReceiptsUpdate})
	cleanUpGroupDir(t, groupId)
	receiptId := seedReceipt(t, groupId, userId, 0)

	w, r := uploadImageRequest(t, userId, receiptId, "receipt.jpg", readHandlerTestJpg(t))
	UploadReceiptImage(w, r)
	assertStatus(t, w, http.StatusOK)

	before, after := assertOneRecordedChange(t, receiptId, userId)
	if len(before.ImageFiles) != 0 || len(after.ImageFiles) != 1 || after.ImageFiles[0].Name != "receipt.jpg" {
		t.Errorf("recorded %d -> %d images, want the upload in the after alone", len(before.ImageFiles), len(after.ImageFiles))
	}
}

func TestRemoveReceiptImageRecordsAReceiptChangeAndRemovesTheFile(t *testing.T) {
	defer repositories.TruncateTestDb()
	defer services.ClearRolePermissionCacheForTests()
	userId, groupId := seedGroupWithRole(t, "image-remove-recorded", []string{permissions.GroupReceiptsUpdate})
	cleanUpGroupDir(t, groupId)
	receiptId := seedReceipt(t, groupId, userId, 0)
	image := addStoredImage(t, receiptId, "receipt.jpg")

	path, err := repositories.NewFileRepository(nil).BuildFilePath(utils.UintToString(receiptId), utils.UintToString(image.ID), image.Name)
	if err != nil {
		t.Fatalf("BuildFilePath: %v", err)
	}
	if !utils.FileExists(path) {
		t.Fatalf("the image file %s was not written", path)
	}

	w, r := removeImageRequest(userId, image.ID)
	RemoveReceiptImage(w, r)
	assertStatus(t, w, http.StatusOK)

	before, after := assertOneRecordedChange(t, receiptId, userId)
	if len(before.ImageFiles) != 1 || len(after.ImageFiles) != 0 {
		t.Errorf("recorded %d -> %d images, want 1 -> 0", len(before.ImageFiles), len(after.ImageFiles))
	}
	if utils.FileExists(path) {
		t.Errorf("the image file %s is still on disk", path)
	}
	if n := countRows(t, &models.FileData{}); n != 0 {
		t.Errorf("image rows = %d, want 0", n)
	}
}

// The role-requirement refusal comes before the change, so nothing is recorded.
func TestRemoveReceiptImageRefusedByRoleRequirementRecordsNothing(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "image-refused", false, true)
	receiptId := seedReceipt(t, groupId, userId, 0)
	only := addStoredImage(t, receiptId, "only.jpg")

	w, r := removeImageRequest(userId, only.ID)
	RemoveReceiptImage(w, r)
	assertStatus(t, w, http.StatusBadRequest)

	if tasks := recordedReceiptChanges(t, receiptId); len(tasks) != 0 {
		t.Errorf("recorded %+v, want nothing for a refused delete", tasks)
	}
}
