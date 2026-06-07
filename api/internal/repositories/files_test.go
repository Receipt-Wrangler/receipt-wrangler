package repositories

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

// testBasePath returns the absolute path of the api/ directory so tests
// can locate the shared testing/ fixtures regardless of CWD or where the
// repo is checked out (dev container vs. CI vs. local clone).
func testBasePath() string {
	_, file, _, _ := runtime.Caller(0)
	// file = .../api/internal/repositories/files_test.go
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// readTestJpgBytes is a local helper so tests don't depend on basePath
// being set at package init.
func readTestJpgBytes(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join(testBasePath(), "testing", "test.jpg")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to read test fixture %s: %v", path, err)
	}
	return b
}

// makePdfFromJpg converts a JPG byte blob into a minimal PDF using
// ImageMagick. Returns the PDF bytes for use as a ConvertPdfToJpg fixture.
func makePdfFromJpg(t *testing.T, jpgBytes []byte) []byte {
	t.Helper()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(jpgBytes); err != nil {
		t.Fatalf("ImageMagick ReadImageBlob failed: %v", err)
	}
	if err := mw.SetImageFormat("pdf"); err != nil {
		t.Fatalf("ImageMagick SetImageFormat pdf failed: %v", err)
	}
	pdf, err := mw.GetImageBlob()
	if err != nil {
		t.Fatalf("ImageMagick GetImageBlob failed: %v", err)
	}
	return pdf
}

func TestShouldZipMultipleFilesSuccessfully(t *testing.T) {
	repository := NewFileRepository(nil)

	filenames := []string{"file1.txt", "file2.txt", "file3.txt"}
	fileContents := [][]byte{
		[]byte("Content of file 1"),
		[]byte("Content of file 2"),
		[]byte("Content of file 3"),
	}

	zipData, err := repository.ZipFiles(filenames, fileContents)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	// Verify zip contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(zipReader.File) != len(filenames) {
		utils.PrintTestError(t, len(zipReader.File), len(filenames))
		return
	}

	// Check each file in the zip
	for i, zipFile := range zipReader.File {
		if zipFile.Name != filenames[i] {
			utils.PrintTestError(t, zipFile.Name, filenames[i])
		}

		rc, err := zipFile.Open()
		if err != nil {
			utils.PrintTestError(t, err, nil)
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			utils.PrintTestError(t, err, nil)
			continue
		}

		if string(content) != string(fileContents[i]) {
			utils.PrintTestError(t, string(content), string(fileContents[i]))
		}
	}
}

func TestShouldReturnErrorWhenFilenamesAndContentsDontMatch(t *testing.T) {
	repository := NewFileRepository(nil)

	filenames := []string{"file1.txt", "file2.txt"}
	fileContents := [][]byte{[]byte("Content of file 1")}

	_, err := repository.ZipFiles(filenames, fileContents)

	expectedError := "number of filenames does not match number of file contents"
	if err == nil || err.Error() != expectedError {
		utils.PrintTestError(t, err, expectedError)
	}
}

func TestShouldReturnErrorWhenNoFilesAreProvided(t *testing.T) {
	repository := NewFileRepository(nil)

	filenames := []string{}
	fileContents := [][]byte{}

	_, err := repository.ZipFiles(filenames, fileContents)

	expectedError := "no files to zip"
	if err == nil || err.Error() != expectedError {
		utils.PrintTestError(t, err, expectedError)
	}
}

func TestShouldHandleEmptyFileContent(t *testing.T) {
	repository := NewFileRepository(nil)

	filenames := []string{"empty.txt"}
	fileContents := [][]byte{[]byte("")}

	zipData, err := repository.ZipFiles(filenames, fileContents)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Verify zip contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(zipReader.File) != 1 {
		utils.PrintTestError(t, len(zipReader.File), 1)
		return
	}

	rc, err := zipReader.File[0].Open()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	content, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(content) != 0 {
		utils.PrintTestError(t, len(content), 0)
	}
}

func TestShouldHandleLargeFileContent(t *testing.T) {
	repository := NewFileRepository(nil)

	// Create a 100KB file
	largeContent := bytes.Repeat([]byte("A"), 100*1024)

	filenames := []string{"large.txt"}
	fileContents := [][]byte{largeContent}

	zipData, err := repository.ZipFiles(filenames, fileContents)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Verify zip contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	rc, err := zipReader.File[0].Open()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	content, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(content) != len(largeContent) {
		utils.PrintTestError(t, len(content), len(largeContent))
	}
}

func TestShouldHandleSpecialCharactersInFilenames(t *testing.T) {
	repository := NewFileRepository(nil)

	filenames := []string{"special!@#$%^&*.txt", "path/with/slashes.txt", "空白.txt"}
	fileContents := [][]byte{
		[]byte("Special content"),
		[]byte("Path content"),
		[]byte("Unicode content"),
	}

	zipData, err := repository.ZipFiles(filenames, fileContents)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Verify zip contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Check each filename exists in the zip
	for i, expectedName := range filenames {
		found := false
		for _, file := range zipReader.File {
			if file.Name == expectedName {
				found = true

				rc, err := file.Open()
				if err != nil {
					utils.PrintTestError(t, err, nil)
					break
				}

				content, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					utils.PrintTestError(t, err, nil)
					break
				}

				if string(content) != string(fileContents[i]) {
					utils.PrintTestError(t, string(content), string(fileContents[i]))
				}

				break
			}
		}

		if !found {
			utils.PrintTestError(t, "File not found", expectedName)
		}
	}
}

// ---------- Pure-function tests (no filesystem/DB) ----------

func TestValidateFileType_AcceptsJpg(t *testing.T) {
	repository := NewFileRepository(nil)
	jpg := readTestJpgBytes(t)

	got, err := repository.ValidateFileType(jpg)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.HasPrefix(got, "image/") {
		utils.PrintTestError(t, got, "image/*")
	}
}

func TestValidateFileType_RejectsText(t *testing.T) {
	repository := NewFileRepository(nil)

	_, err := repository.ValidateFileType([]byte("just plain text"))
	if err == nil {
		utils.PrintTestError(t, err, "expected invalid file type error")
	}
}

func TestValidateJsonFileType_AcceptsJson(t *testing.T) {
	repository := NewFileRepository(nil)

	got, err := repository.ValidateJsonFileType([]byte(`{"a":1}`))
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.Contains(got, "json") {
		utils.PrintTestError(t, got, "*json*")
	}
}

func TestValidateJsonFileType_RejectsImage(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.ValidateJsonFileType(readTestJpgBytes(t))
	if err == nil {
		utils.PrintTestError(t, err, "expected invalid file type error for image")
	}
}

func TestIsImage_TrueForJpg(t *testing.T) {
	repository := NewFileRepository(nil)

	got, err := repository.IsImage(readTestJpgBytes(t))
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !got {
		utils.PrintTestError(t, got, true)
	}
}

func TestIsImage_ErrorOnInvalidType(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.IsImage([]byte("plain text"))
	if err == nil {
		utils.PrintTestError(t, err, "expected validate error")
	}
}

func TestIsPdf_FalseForJpg(t *testing.T) {
	repository := NewFileRepository(nil)
	got, err := repository.IsPdf(readTestJpgBytes(t))
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if got {
		utils.PrintTestError(t, got, false)
	}
}

func TestIsPdf_TrueForPdf(t *testing.T) {
	repository := NewFileRepository(nil)
	pdf := makePdfFromJpg(t, readTestJpgBytes(t))

	got, err := repository.IsPdf(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !got {
		utils.PrintTestError(t, got, true)
	}
}

func TestIsPdf_ErrorOnInvalidType(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.IsPdf([]byte("plain text"))
	if err == nil {
		utils.PrintTestError(t, err, "expected validate error")
	}
}

func TestGetFileType_PdfMappedToImageJpeg(t *testing.T) {
	repository := NewFileRepository(nil)
	pdf := makePdfFromJpg(t, readTestJpgBytes(t))

	got, err := repository.GetFileType(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if got != "image/jpeg" {
		utils.PrintTestError(t, got, "image/jpeg")
	}
}

func TestGetFileType_JpgPassThrough(t *testing.T) {
	repository := NewFileRepository(nil)
	got, err := repository.GetFileType(readTestJpgBytes(t))
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.HasPrefix(got, "image/") {
		utils.PrintTestError(t, got, "image/*")
	}
}

func TestGetFileType_ErrorOnInvalidType(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.GetFileType([]byte("plain text"))
	if err == nil {
		utils.PrintTestError(t, err, "expected validate error")
	}
}

func TestBuildEncodedImageString_ProducesDataUri(t *testing.T) {
	repository := NewFileRepository(nil)
	got, err := repository.BuildEncodedImageString(readTestJpgBytes(t))
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.HasPrefix(got, "data:image/") {
		utils.PrintTestError(t, got[:25], "data:image/... prefix")
	}
	if !strings.Contains(got, "base64,") {
		utils.PrintTestError(t, got, "contains base64,")
	}
}

func TestBuildEncodedImageString_ErrorOnInvalidType(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.BuildEncodedImageString([]byte("plain text"))
	if err == nil {
		utils.PrintTestError(t, err, "expected error")
	}
}

// ---------- Temp-path/filesystem tests ----------

func TestGetTempDirectoryPath_UsesBasePath(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	got := repository.GetTempDirectoryPath()
	want := filepath.Join(testBasePath(), "temp")
	if got != want {
		utils.PrintTestError(t, got, want)
	}
}

func TestBuildTempFilePath_FormatAndExtension(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	got, err := repository.BuildTempFilePath("jpg")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.HasPrefix(got, filepath.Join(testBasePath(), "temp")+"/") {
		utils.PrintTestError(t, got, "prefixed with temp dir")
	}
	if !strings.HasSuffix(got, ".jpg") {
		utils.PrintTestError(t, got, "*.jpg")
	}
}

func TestGetTestJpgBytes_ReadsFixture(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	got, err := repository.GetTestJpgBytes()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if len(got) == 0 {
		utils.PrintTestError(t, len(got), ">0")
	}
}

func TestWriteTempFile_WritesJpgBytes(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	jpg := readTestJpgBytes(t)

	path, err := repository.WriteTempFile(jpg)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	defer os.Remove(path)

	if !strings.HasSuffix(path, ".jpeg") {
		utils.PrintTestError(t, path, "*.jpeg suffix")
	}
	info, err := os.Stat(path)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if info.Size() == 0 {
		utils.PrintTestError(t, info.Size(), ">0")
	}
}

func TestWriteTempFile_RejectsInvalidFileType(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	_, err := repository.WriteTempFile([]byte("not an image"))
	if err == nil {
		utils.PrintTestError(t, err, "expected invalid-file-type error")
	}
}

// ---------- ImageMagick conversion tests ----------

func TestConvertPdfToJpg_HappyPath(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	pdf := makePdfFromJpg(t, readTestJpgBytes(t))

	out, err := repository.ConvertPdfToJpg(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if len(out) == 0 {
		utils.PrintTestError(t, len(out), ">0")
	}
	isImage, err := repository.IsImage(out)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !isImage {
		utils.PrintTestError(t, isImage, true)
	}
}

func TestConvertPdfToJpg_InvalidBytes(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	_, err := repository.ConvertPdfToJpg([]byte("not a pdf"))
	if err == nil {
		utils.PrintTestError(t, err, "expected ImageMagick error")
	}
}

func TestPdfRasterizationDpi_DefaultsTo300(t *testing.T) {
	t.Setenv("PDF_DPI", "")
	if got := pdfRasterizationDpi(); got != 300 {
		utils.PrintTestError(t, got, 300)
	}
}

func TestPdfRasterizationDpi_HonorsEnv(t *testing.T) {
	t.Setenv("PDF_DPI", "150")
	if got := pdfRasterizationDpi(); got != 150 {
		utils.PrintTestError(t, got, 150)
	}
}

func TestPdfRasterizationDpi_FallsBackOnInvalid(t *testing.T) {
	for _, bad := range []string{"abc", "0", "-100"} {
		t.Setenv("PDF_DPI", bad)
		if got := pdfRasterizationDpi(); got != 300 {
			utils.PrintTestError(t, got, 300)
		}
	}
}

// Higher DPI must rasterize the same PDF into a physically larger image than
// the previous hard-coded 72-DPI default, proving the resolution is applied
// before ReadImageBlob. We compare encoded JPEG byte length as a proxy for
// pixel dimensions; a 300-DPI raster of the same page is reliably larger.
func TestConvertPdfToJpg_HigherDpiProducesLargerImage(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())
	repository := NewFileRepository(nil)
	pdf := makePdfFromJpg(t, readTestJpgBytes(t))

	t.Setenv("PDF_DPI", "72")
	low, err := repository.ConvertPdfToJpg(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	t.Setenv("PDF_DPI", "300")
	high, err := repository.ConvertPdfToJpg(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if !(len(high) > len(low)) {
		utils.PrintTestError(t, len(high), "greater than "+utils.UintToString(uint(len(low))))
	}
}

func TestConvertHeicToJpg_InvalidBytes(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.ConvertHeicToJpg([]byte("not heic"))
	if err == nil {
		utils.PrintTestError(t, err, "expected ImageMagick error")
	}
}

func TestGetBytesFromImageBytes_JpgPassThrough(t *testing.T) {
	repository := NewFileRepository(nil)
	jpg := readTestJpgBytes(t)

	got, err := repository.GetBytesFromImageBytes(jpg)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if len(got) != len(jpg) {
		utils.PrintTestError(t, len(got), len(jpg))
	}
}

func TestGetBytesFromImageBytes_PdfRoutesThroughConversion(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	pdf := makePdfFromJpg(t, readTestJpgBytes(t))

	got, err := repository.GetBytesFromImageBytes(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if len(got) == 0 {
		utils.PrintTestError(t, len(got), ">0")
	}
	// Result should be an image, not a PDF any more.
	isImage, _ := repository.IsImage(got)
	if !isImage {
		utils.PrintTestError(t, isImage, true)
	}
}

func TestGetBytesFromImageBytes_InvalidTypeErrors(t *testing.T) {
	repository := NewFileRepository(nil)
	_, err := repository.GetBytesFromImageBytes([]byte("nope"))
	if err == nil {
		utils.PrintTestError(t, err, "expected invalid-file-type error")
	}
}

// ---------- DB-backed path helpers ----------

func TestBuildGroupPath_WithAlternateName(t *testing.T) {
	defer TruncateTestDb()

	repository := NewFileRepository(nil)
	got, err := repository.BuildGroupPath(42, "alt-name")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.Contains(got, "42") || !strings.Contains(got, "alt-name") {
		utils.PrintTestError(t, got, "contains 42 and alt-name")
	}
}

func TestBuildGroupPath_LooksUpGroupByName(t *testing.T) {
	defer TruncateTestDb()

	db := GetDB()
	group := models.Group{Name: "lookup-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}

	repository := NewFileRepository(nil)
	got, err := repository.BuildGroupPath(group.ID, "")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.Contains(got, "lookup-group") {
		utils.PrintTestError(t, got, "contains lookup-group")
	}
}

func TestBuildFilePath_ComposesGroupAndFileName(t *testing.T) {
	defer TruncateTestDb()

	db := GetDB()
	user := models.User{Username: "file-path-user", Password: "p", DisplayName: "x"}
	if err := db.Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	group := models.Group{Name: "files-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	receipt := models.Receipt{Name: "r", GroupId: group.ID, PaidByUserID: user.ID}
	if err := db.Create(&receipt).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}

	repository := NewFileRepository(nil)
	got, err := repository.BuildFilePath(utils.UintToString(receipt.ID), "99", "photo.jpg")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if !strings.Contains(got, "files-group") {
		utils.PrintTestError(t, got, "contains files-group")
	}
	if !strings.HasSuffix(got, ".jpg") {
		utils.PrintTestError(t, got, "*.jpg")
	}
}

// ---------- Zip from temp files ----------

func TestCreateZipFromTempFiles_ArchivesSeededFiles(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())

	repository := NewFileRepository(nil)
	tempDir := repository.GetTempDirectoryPath()
	_ = os.MkdirAll(tempDir, 0o755)

	f1 := "zip-src-a.txt"
	f2 := "zip-src-b.txt"
	if err := os.WriteFile(filepath.Join(tempDir, f1), []byte("alpha"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, f2), []byte("beta"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	defer os.Remove(filepath.Join(tempDir, f1))
	defer os.Remove(filepath.Join(tempDir, f2))

	zipName := "zipped-output.zip"
	zipPath, err := repository.CreateZipFromTempFiles(zipName, []string{f1, f2})
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	defer os.Remove(zipPath)

	info, err := os.Stat(zipPath)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if info.Size() == 0 {
		utils.PrintTestError(t, info.Size(), ">0")
	}

	zipBytes, _ := os.ReadFile(zipPath)
	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(reader.File) != 2 {
		utils.PrintTestError(t, len(reader.File), 2)
	}
}

func TestGetBytesForFileData_ReadsAndReturnsBytes(t *testing.T) {
	t.Setenv("BASE_PATH", testBasePath())
	defer TruncateTestDb()

	db := GetDB()
	user := models.User{Username: "fd-user", Password: "p", DisplayName: "x"}
	if err := db.Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	group := models.Group{Name: "bytes-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	receipt := models.Receipt{Name: "r", GroupId: group.ID, PaidByUserID: user.ID}
	if err := db.Create(&receipt).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	fileData := models.FileData{Name: "img.jpg", ReceiptId: receipt.ID, FileType: "image/jpeg"}
	if err := db.Create(&fileData).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}

	// Write the JPG fixture to the exact path BuildFilePath will construct.
	repository := NewFileRepository(nil)
	targetPath, err := repository.BuildFilePath(utils.UintToString(receipt.ID), utils.UintToString(fileData.ID), fileData.Name)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	_ = os.MkdirAll(filepath.Dir(targetPath), 0o755)
	jpg := readTestJpgBytes(t)
	if err := os.WriteFile(targetPath, jpg, 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}
	defer os.Remove(targetPath)

	got, err := repository.GetBytesForFileData(fileData)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if len(got) == 0 {
		utils.PrintTestError(t, len(got), ">0")
	}
}

// Sanity: a PDF fixture really is recognised as ApplicationPdf.
func TestValidateFileType_AcceptsPdf(t *testing.T) {
	repository := NewFileRepository(nil)
	pdf := makePdfFromJpg(t, readTestJpgBytes(t))

	got, err := repository.ValidateFileType(pdf)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if got != constants.ApplicationPdf {
		utils.PrintTestError(t, got, constants.ApplicationPdf)
	}
}

// Example of a table test in the style of the existing tests
func TestShouldValidateZipFilesInput(t *testing.T) {
	repository := NewFileRepository(nil)

	tests := map[string]struct {
		filenames    []string
		fileContents [][]byte
		expectErr    bool
		expectedMsg  string
	}{
		"mismatched counts": {
			filenames:    []string{"file1.txt", "file2.txt"},
			fileContents: [][]byte{[]byte("Content")},
			expectErr:    true,
			expectedMsg:  "number of filenames does not match number of file contents",
		},
		"no files": {
			filenames:    []string{},
			fileContents: [][]byte{},
			expectErr:    true,
			expectedMsg:  "no files to zip",
		},
		"valid input": {
			filenames:    []string{"file.txt"},
			fileContents: [][]byte{[]byte("Content")},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		zipData, err := repository.ZipFiles(test.filenames, test.fileContents)

		if test.expectErr {
			if err == nil || err.Error() != test.expectedMsg {
				utils.PrintTestError(t, err, test.expectedMsg)
			}
		} else {
			if err != nil {
				utils.PrintTestError(t, err, nil)
			}

			// Basic check that the zip was created
			if zipData == nil || len(zipData) == 0 {
				utils.PrintTestError(t, "empty data", "non-empty zip data")
			}
		}
	}
}
