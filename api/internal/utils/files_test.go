package utils

import (
	"os"
	"testing"
)

func TestWritesFile(t *testing.T) {
	path := "test.txt"
	fileContents := "test"

	WriteFile(path, []byte(fileContents))

	_, err := os.Stat(path)

	if err != nil {
		PrintTestError(t, err, "Expected file to be written")
	}

	os.Remove(path)
}

func TestFileIsRead(t *testing.T) {
	path := "test.txt"
	fileContents := "test"

	WriteFile(path, []byte(fileContents))

	_, err := os.Stat(path)

	if err != nil {
		PrintTestError(t, err, "Expected file to be written")
	}

	contents, err := ReadFile(path)
	if err != nil {
		PrintTestError(t, err, "Expected contents to be read")
	}

	if string(contents) != "test" {
		PrintTestError(t, contents, "test")
	}

	os.Remove(path)
}

func TestShouldReturnNoErrIfDirExists(t *testing.T) {
	path := "../utils"

	err := DirectoryExists(path, false)
	if err != nil {
		PrintTestError(t, err, "Expected directory to exist")
	}
}

func TestShouldReturnErrIfDirDoesNotExists(t *testing.T) {
	path := "./fakeDir"

	err := DirectoryExists(path, false)
	if err == nil {
		PrintTestError(t, err, "Expected error to exist")
	}
}

func TestShouldCreateDirIfItDoesntExist(t *testing.T) {
	path := "./fakeDir"

	err := DirectoryExists(path, true)
	if err != nil {
		PrintTestError(t, err, "Expected no error")
	}

	err = DirectoryExists(path, false)
	if err != nil {
		PrintTestError(t, err, "Expected directory to exist")
	}

	os.Remove(path)
}

func TestShouldCreateDirectory(t *testing.T) {
	path := "./fakeDir"

	err := MakeDirectory(path)
	if err != nil {
		PrintTestError(t, err, "Expected no error")
	}

	err = DirectoryExists(path, false)
	if err != nil {
		PrintTestError(t, err, "Expected no error")
	}

	os.Remove(path)
}

func TestFileExistsShouldReturnTrueForExistingFile(t *testing.T) {
	path := "test_exists.txt"
	WriteFile(path, []byte("test"))
	defer os.Remove(path)

	exists := FileExists(path)
	if !exists {
		t.Errorf("Expected FileExists to return true for existing file")
	}
}

func TestFileExistsShouldReturnFalseForNonExistingFile(t *testing.T) {
	path := "nonexistent_file_12345.txt"

	exists := FileExists(path)
	if exists {
		t.Errorf("Expected FileExists to return false for non-existing file")
	}
}

func TestReadLastFileLineShouldReadLastLine(t *testing.T) {
	path := "test_last_line.txt"
	content := "line1\nline2\nline3\nlast line"
	WriteFile(path, []byte(content))
	defer os.Remove(path)

	lastLine, err := ReadLastFileLine(path)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if lastLine != "last line" {
		PrintTestError(t, lastLine, "last line")
	}
}

func TestReadLastFileLineShouldHandleSingleLine(t *testing.T) {
	path := "test_single_line.txt"
	content := "only one line"
	WriteFile(path, []byte(content))
	defer os.Remove(path)

	lastLine, err := ReadLastFileLine(path)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if lastLine != "only one line" {
		PrintTestError(t, lastLine, "only one line")
	}
}

func TestReadLastFileLineShouldHandleEmptyFile(t *testing.T) {
	path := "test_empty.txt"
	WriteFile(path, []byte(""))
	defer os.Remove(path)

	lastLine, err := ReadLastFileLine(path)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	if lastLine != "" {
		PrintTestError(t, lastLine, "")
	}
}

func TestReadLastFileLineShouldHandleNonExistingFileInTestEnv(t *testing.T) {
	os.Setenv("ENV", "test")
	defer os.Unsetenv("ENV")

	path := "nonexistent_file_for_lastline.txt"

	lastLine, err := ReadLastFileLine(path)
	if err != nil {
		t.Errorf("Expected no error in test env for non-existing file, got %v", err)
	}

	if lastLine != "" {
		PrintTestError(t, lastLine, "")
	}
}

func TestBuildGroupPathStringShouldBuildCorrectPath(t *testing.T) {
	groupId := "123"
	groupName := "TestGroup"

	path, err := BuildGroupPathString(groupId, groupName)
	if err != nil {
		PrintTestError(t, err, nil)
	}

	// Path should end with data/123-TestGroup
	expectedSuffix := "data/123-TestGroup"
	if len(path) < len(expectedSuffix) {
		t.Errorf("Path too short: %s", path)
	}

	// Check that path ends with expected suffix (using filepath separators)
	if path[len(path)-len(expectedSuffix):] != expectedSuffix &&
	   path[len(path)-len(expectedSuffix):] != "data\\123-TestGroup" {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestBuildFileNameShouldBuildCorrectFileName(t *testing.T) {
	rid := "receipt123"
	fid := "file456"
	fname := "document.pdf"

	result := BuildFileName(rid, fid, fname)
	expected := "receipt123-file456-document.pdf"

	if result != expected {
		PrintTestError(t, result, expected)
	}
}

func TestBuildFileNameShouldHandleEmptyParts(t *testing.T) {
	result := BuildFileName("", "", "")
	expected := "--"

	if result != expected {
		PrintTestError(t, result, expected)
	}
}

func TestGetMimeTypeShouldDetectPNG(t *testing.T) {
	// PNG magic bytes
	pngBytes := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	mimeType := GetMimeType(pngBytes)

	if mimeType.String() != "image/png" {
		PrintTestError(t, mimeType.String(), "image/png")
	}
}

func TestGetMimeTypeShouldDetectJPEG(t *testing.T) {
	// JPEG magic bytes
	jpegBytes := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}

	mimeType := GetMimeType(jpegBytes)

	if mimeType.String() != "image/jpeg" {
		PrintTestError(t, mimeType.String(), "image/jpeg")
	}
}

func TestGetMimeTypeShouldDetectPDF(t *testing.T) {
	// PDF magic bytes
	pdfBytes := []byte{0x25, 0x50, 0x44, 0x46, 0x2D} // %PDF-

	mimeType := GetMimeType(pdfBytes)

	if mimeType.String() != "application/pdf" {
		PrintTestError(t, mimeType.String(), "application/pdf")
	}
}

func TestGetMimeTypeShouldHandleUnknownBytes(t *testing.T) {
	// Random bytes
	randomBytes := []byte{0x00, 0x01, 0x02, 0x03}

	mimeType := GetMimeType(randomBytes)

	// Should return a valid MIME type (usually application/octet-stream or text/plain)
	if mimeType == nil {
		t.Errorf("Expected GetMimeType to return a MIME type, got nil")
	}
}

func TestReadFileShouldReturnNilForNonExistingFile(t *testing.T) {
	path := "nonexistent_file_for_read.txt"

	content, err := ReadFile(path)
	// Based on the implementation, it returns nil, nil for non-existing files
	if err != nil {
		t.Errorf("Expected nil error for non-existing file, got %v", err)
	}
	if content != nil {
		t.Errorf("Expected nil content for non-existing file, got %v", content)
	}
}

func TestWriteFileShouldReturnErrorForInvalidPath(t *testing.T) {
	// Try to write to an invalid path (directory that doesn't exist)
	path := "/nonexistent_directory_12345/test.txt"

	err := WriteFile(path, []byte("test"))
	if err == nil {
		t.Errorf("Expected error for invalid path, got nil")
		// Clean up in case it somehow succeeded
		os.Remove(path)
	}
}

func TestMakeDirectoryShouldReturnErrorForExistingDirectory(t *testing.T) {
	path := "./existingDir"
	MakeDirectory(path)
	defer os.Remove(path)

	// Try to create it again
	err := MakeDirectory(path)
	if err == nil {
		t.Errorf("Expected error when creating existing directory, got nil")
	}
}

func TestReadLastFileLineShouldReturnErrorForNonExistingFileOutsideTestEnv(t *testing.T) {
	// Ensure ENV is not set to "test"
	originalEnv := os.Getenv("ENV")
	os.Unsetenv("ENV")
	defer func() {
		if originalEnv != "" {
			os.Setenv("ENV", originalEnv)
		}
	}()

	path := "nonexistent_file_for_lastline_error.txt"

	_, err := ReadLastFileLine(path)
	if err == nil {
		t.Errorf("Expected error for non-existing file outside test env, got nil")
	}
}

func TestDirectoryExistsShouldReturnErrorWhenCreateFails(t *testing.T) {
	// Try to create a directory in a path that doesn't exist
	path := "/nonexistent_parent_dir_12345/newdir"

	err := DirectoryExists(path, true)
	if err == nil {
		t.Errorf("Expected error when creating directory in non-existent parent, got nil")
		// Clean up in case it somehow succeeded
		os.RemoveAll(path)
	}
}
