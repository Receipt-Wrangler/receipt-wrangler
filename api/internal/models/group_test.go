package models

import (
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/utils"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newInMemoryDb opens a fresh, single-connection in-memory SQLite database. A
// single connection keeps every query pinned to the same in-memory instance so
// tables created here stay visible for the rest of the test.
func newInMemoryDb(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	return db
}

// newGroupsTableDb returns an in-memory database with a minimal groups table.
// Group.BeforeUpdate only selects id and name, so no other columns are needed.
func newGroupsTableDb(t *testing.T) *gorm.DB {
	db := newInMemoryDb(t)
	if err := db.Exec("CREATE TABLE groups (id integer primary key, name text)").Error; err != nil {
		t.Fatalf("failed to create groups table: %v", err)
	}
	return db
}

func TestGroup_BeforeUpdate_NoOpWhenIdZero(t *testing.T) {
	// With a zero ID the hook returns immediately and never touches the tx.
	group := Group{Name: "Test"}
	err := group.BeforeUpdate(nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
}

func TestGroup_BeforeUpdate_NameUnchanged(t *testing.T) {
	db := newGroupsTableDb(t)
	if err := db.Exec("INSERT INTO groups (id, name) VALUES (1, 'Test')").Error; err != nil {
		t.Fatalf("failed to seed group: %v", err)
	}

	group := Group{BaseModel: BaseModel{ID: 1}, Name: "Test"}
	err := group.BeforeUpdate(db)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
}

func TestGroup_BeforeUpdate_NameChanged(t *testing.T) {
	db := newGroupsTableDb(t)
	if err := db.Exec("INSERT INTO groups (id, name) VALUES (1, 'Old')").Error; err != nil {
		t.Fatalf("failed to seed group: %v", err)
	}

	// The old directory does not exist, so the internal utils.RenameDataPath fails
	// silently (its error is deliberately ignored) and the hook still returns nil.
	group := Group{BaseModel: BaseModel{ID: 1}, Name: "New"}
	err := group.BeforeUpdate(db)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
}

func TestGroup_BeforeUpdate_DbError(t *testing.T) {
	// No groups table exists, so the lookup query errors and the hook propagates it.
	db := newInMemoryDb(t)

	group := Group{BaseModel: BaseModel{ID: 1}, Name: "New"}
	err := group.BeforeUpdate(db)
	if err == nil {
		utils.PrintTestError(t, err, "a db error")
	}
}

func TestGroup_AfterDelete_NoOpWhenIdZero(t *testing.T) {
	// With a zero ID the hook returns immediately without touching the filesystem.
	group := Group{Name: "Test"}
	err := group.AfterDelete(nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
}

func TestGroup_AfterDelete_RemovesGroupPath(t *testing.T) {
	// The group path does not exist, so utils.RemoveAllInDataDir is a no-op that
	// returns nil.
	group := Group{BaseModel: BaseModel{ID: 1}, Name: "Test"}
	err := group.AfterDelete(nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
}

// AfterDelete runs a recursive remove on the group's name-derived storage path. A
// group whose name resolves outside the data directory (CWE-22) must be refused so
// the delete hook can never recursively remove an arbitrary directory.
func TestGroup_AfterDelete_RejectsPathTraversalName(t *testing.T) {
	// The sentinel is placed at exactly the location the crafted name below
	// resolves to (10 "../" climb out of <cwd>/data to the filesystem root, then
	// into /tmp), so the "not deleted" assertion is meaningful: without the guard
	// AfterDelete's recursive remove would wipe this directory.
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
