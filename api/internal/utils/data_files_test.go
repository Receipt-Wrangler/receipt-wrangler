package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSafePathComponent(t *testing.T) {
	tests := map[string]bool{
		"My Receipts":         true,
		"Reporting Load Test": true,
		"group-123":           true,
		"Mom & Dad's":         true,
		"":                    false,
		".":                   false,
		"..":                  false,
		"../etc":              false,
		"a/b":                 false,
		"a\\b":                false,
		"/etc/passwd":         false,
		"foo/../bar":          false,
	}

	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			if got := IsSafePathComponent(name); got != want {
				PrintTestError(t, got, want)
			}
		})
	}

	if IsSafePathComponent("foo\x00bar") {
		PrintTestError(t, true, false)
	}
}

func TestBuildGroupPathString_ValidNameStaysInDataDir(t *testing.T) {
	dataDir, err := GetDataDir()
	if err != nil {
		t.Fatalf("GetDataDir: %v", err)
	}

	got, err := BuildGroupPathString("42", "My Receipts")
	if err != nil {
		t.Fatalf("unexpected error for a valid name: %v", err)
	}

	want := filepath.Join(dataDir, "42-My Receipts")
	if got != want {
		PrintTestError(t, got, want)
	}
}

func TestBuildGroupPathString_RejectsTraversal(t *testing.T) {
	// Enough ".." to climb out of the data dir regardless of the test cwd depth.
	names := []string{
		"../../../../../../../../../../tmp/x",
		"../../../../../../../../../../etc",
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			got, err := BuildGroupPathString("42", name)
			if err == nil {
				PrintTestError(t, got, "error: path escapes data directory")
			}
		})
	}
}

// The centralized data-file helpers must refuse any path outside the data dir
// and leave an out-of-tree target untouched.
func TestDataFileHelpers_RejectPathsOutsideDataDir(t *testing.T) {
	sentinel := filepath.Join(os.TempDir(), "rw_utils_out_sentinel")
	if err := os.MkdirAll(sentinel, 0o755); err != nil {
		t.Fatalf("setup sentinel: %v", err)
	}
	defer os.RemoveAll(sentinel)
	marker := filepath.Join(sentinel, "marker.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatalf("setup marker: %v", err)
	}

	assertRejected := func(name string, err error) {
		if err == nil {
			PrintTestError(t, name+" accepted an out-of-tree path", "error")
		}
	}

	assertRejected("MakeDataDirectory", MakeDataDirectory(filepath.Join(sentinel, "child")))
	assertRejected("RemoveDataPath", RemoveDataPath(marker))
	assertRejected("RemoveAllInDataDir", RemoveAllInDataDir(sentinel))
	assertRejected("RenameDataPath(dst-outside)", RenameDataPath(marker, filepath.Join(os.TempDir(), "rw_utils_out_moved")))
	assertRejected("WriteDataFile", WriteDataFile(filepath.Join(sentinel, "w.txt"), []byte("x")))
	assertRejected("EnsureDataDirectory", EnsureDataDirectory(filepath.Join(sentinel, "d")))
	_, readErr := ReadDataFile(marker)
	assertRejected("ReadDataFile", readErr)

	if _, err := os.Stat(sentinel); os.IsNotExist(err) {
		t.Fatalf("out-of-tree sentinel was removed — containment guard failed")
	}
	if _, err := os.Stat(marker); os.IsNotExist(err) {
		t.Fatalf("out-of-tree marker was removed — containment guard failed")
	}
}

func TestDataFileHelpers_AllowPathsInsideDataDir(t *testing.T) {
	dataDir, err := GetDataDir()
	if err != nil {
		t.Fatalf("GetDataDir: %v", err)
	}
	createdDataDir := !FileExists(dataDir)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatalf("setup data dir: %v", err)
	}
	if createdDataDir {
		defer os.RemoveAll(dataDir)
	}

	dir := filepath.Join(dataDir, "rw_utils_in_dir")
	defer os.RemoveAll(dir)
	if err := MakeDataDirectory(dir); err != nil {
		t.Fatalf("MakeDataDirectory in data dir failed: %v", err)
	}
	if !FileExists(dir) {
		t.Fatalf("expected directory to be created")
	}

	renamed := filepath.Join(dataDir, "rw_utils_in_dir_renamed")
	defer os.RemoveAll(renamed)
	if err := RenameDataPath(dir, renamed); err != nil {
		t.Fatalf("RenameDataPath in data dir failed: %v", err)
	}

	if err := RemoveAllInDataDir(renamed); err != nil {
		t.Fatalf("RemoveAllInDataDir in data dir failed: %v", err)
	}
	if FileExists(renamed) {
		t.Fatalf("expected directory to be removed")
	}

	// WriteDataFile + ReadDataFile round-trip inside the data dir.
	file := filepath.Join(dataDir, "rw_utils_in_file.txt")
	defer os.RemoveAll(file)
	if err := WriteDataFile(file, []byte("hello")); err != nil {
		t.Fatalf("WriteDataFile in data dir failed: %v", err)
	}
	got, err := ReadDataFile(file)
	if err != nil {
		t.Fatalf("ReadDataFile in data dir failed: %v", err)
	}
	if string(got) != "hello" {
		PrintTestError(t, string(got), "hello")
	}

	// EnsureDataDirectory tolerates an already-existing directory.
	if err := EnsureDataDirectory(dataDir); err != nil {
		t.Fatalf("EnsureDataDirectory on an existing dir failed: %v", err)
	}
}

// ReadDataFile must surface a real read error (unlike utils.ReadFile, which
// swallows errors as nil, nil).
func TestReadDataFile_PropagatesReadError(t *testing.T) {
	dataDir, err := GetDataDir()
	if err != nil {
		t.Fatalf("GetDataDir: %v", err)
	}
	createdDataDir := !FileExists(dataDir)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatalf("setup data dir: %v", err)
	}
	if createdDataDir {
		defer os.RemoveAll(dataDir)
	}

	if _, err := ReadDataFile(filepath.Join(dataDir, "definitely-missing.txt")); err == nil {
		t.Fatalf("expected ReadDataFile to propagate a read error for a missing file")
	}
}
