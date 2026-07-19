package models

import (
	"os"
	"path/filepath"
	"testing"
)

// AfterDelete runs os.RemoveAll on the group's name-derived storage path. A group
// whose name resolves outside the data directory (CWE-22) must be refused so the
// delete hook can never recursively remove an arbitrary directory.
func TestGroup_AfterDelete_RejectsPathTraversalName(t *testing.T) {
	// The sentinel is placed at exactly the location the crafted name below
	// resolves to (10 "../" climb out of <cwd>/data to the filesystem root, then
	// into /tmp), so the "not deleted" assertion is meaningful: without the guard
	// AfterDelete's os.RemoveAll would wipe this directory.
	sentinel := "/tmp/rw_model_afterdelete_sentinel"
	if err := os.MkdirAll(sentinel, 0o755); err != nil {
		t.Fatalf("setup sentinel: %v", err)
	}
	defer os.RemoveAll(sentinel)
	if err := os.WriteFile(filepath.Join(sentinel, "marker.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatalf("setup marker: %v", err)
	}

	group := Group{
		BaseModel: BaseModel{ID: 999999},
		Name:      "../../../../../../../../../../tmp/rw_model_afterdelete_sentinel",
	}

	if err := group.AfterDelete(nil); err == nil {
		t.Fatalf("expected AfterDelete to reject a path-traversal group name")
	}

	if _, err := os.Stat(sentinel); os.IsNotExist(err) {
		t.Fatalf("sentinel dir was recursively deleted — traversal guard failed")
	}
}
