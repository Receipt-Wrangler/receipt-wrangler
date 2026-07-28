package repositories

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"testing"
	"time"
)

// GetPagedActivities applies the actor-visibility disjunction in SQL BEFORE COUNT/LIMIT:
// a hidden actor's activities affect neither the returned rows nor the total, even when
// they are newer and would otherwise fill the first page.
func TestGetPagedActivitiesFiltersByVisibilityBeforeCountAndLimit(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	group := models.Group{Name: "act-vis-grp"}
	db.Create(&group)
	visible := models.User{Username: "act-visible", Password: "x"}
	hidden := models.User{Username: "act-hidden", Password: "x"}
	db.Create(&visible)
	db.Create(&hidden)

	gid := group.ID
	// Hidden activities are NEWER (sort first by started_at desc); the visible one is the
	// oldest, so a limit applied before filtering would drop it.
	for i := 0; i < 5; i++ {
		hid := hidden.ID
		db.Create(&models.SystemTask{
			Type:                 models.QUICK_SCAN,
			Status:               models.SYSTEM_TASK_SUCCEEDED,
			AssociatedEntityType: models.NOOP_ENTITY_TYPE,
			GroupId:              &gid,
			RanByUserId:          &hid,
			StartedAt:            time.Now().Add(time.Duration(i+1) * time.Hour),
		})
	}
	vid := visible.ID
	db.Create(&models.SystemTask{
		Type:                 models.QUICK_SCAN,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.NOOP_ENTITY_TYPE,
		GroupId:              &gid,
		RanByUserId:          &vid,
		StartedAt:            time.Now(),
	})

	// The resolver restricts the group to the visible actor only.
	resolver := func(groupId uint) ([]uint, bool, error) {
		return []uint{visible.ID}, false, nil
	}

	command := commands.PagedActivityRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      3,
			OrderBy:       "started_at",
			SortDirection: commands.DESCENDING,
		},
		GroupIds: []uint{group.ID},
	}

	activities, count, err := NewSystemTaskRepository(nil).GetPagedActivities(command, resolver)
	if err != nil {
		t.Fatalf("GetPagedActivities: %v", err)
	}

	// Only the visible actor's one activity is counted and returned, even though the
	// hidden ones are newer and a pre-filter limit of 3 would have returned them.
	if count != 1 {
		t.Errorf("TotalCount = %d, want 1 (filtered before count)", count)
	}
	if len(activities) != 1 {
		t.Fatalf("returned %d rows, want 1", len(activities))
	}
	if activities[0].RanByUserId == nil || *activities[0].RanByUserId != visible.ID {
		t.Errorf("returned activity should be the visible actor's, got %+v", activities[0])
	}
}

// A multi-group request combining an unrestricted (open) group and a restricted (isolated)
// group exercises the `unrestricted`-per-group OR branch: the open group's hidden-actor
// activity is kept while the isolated group's hidden-actor activity is dropped.
func TestGetPagedActivitiesMixedVisibilityAcrossGroups(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	openGroup := models.Group{Name: "act-open-grp"}
	isolatedGroup := models.Group{Name: "act-isolated-grp"}
	db.Create(&openGroup)
	db.Create(&isolatedGroup)
	visible := models.User{Username: "act-mix-visible", Password: "x"}
	hidden := models.User{Username: "act-mix-hidden", Password: "x"}
	db.Create(&visible)
	db.Create(&hidden)

	oid, iid, hid, vid := openGroup.ID, isolatedGroup.ID, hidden.ID, visible.ID

	newTask := func(groupId, ranBy uint) *models.SystemTask {
		gid, uid := groupId, ranBy
		return &models.SystemTask{
			Type:                 models.QUICK_SCAN,
			Status:               models.SYSTEM_TASK_SUCCEEDED,
			AssociatedEntityType: models.NOOP_ENTITY_TYPE,
			GroupId:              &gid,
			RanByUserId:          &uid,
			StartedAt:            time.Now(),
		}
	}
	// Open group: the hidden user's activity should still show (unrestricted).
	db.Create(newTask(oid, hid))
	// Isolated group: the hidden user's activity is dropped; the visible user's is kept.
	db.Create(newTask(iid, hid))
	db.Create(newTask(iid, vid))

	resolver := func(groupId uint) ([]uint, bool, error) {
		if groupId == oid {
			return nil, true, nil // unrestricted for the open group
		}
		return []uint{vid}, false, nil // restricted to the visible user in the isolated group
	}

	command := commands.PagedActivityRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      10,
			OrderBy:       "started_at",
			SortDirection: commands.DESCENDING,
		},
		GroupIds: []uint{oid, iid},
	}

	activities, count, err := NewSystemTaskRepository(nil).GetPagedActivities(command, resolver)
	if err != nil {
		t.Fatalf("GetPagedActivities: %v", err)
	}

	// Open-group hidden actor (kept) + isolated-group visible actor (kept) = 2; the
	// isolated-group hidden actor is dropped.
	if count != 2 {
		t.Errorf("TotalCount = %d, want 2", count)
	}
	if len(activities) != 2 {
		t.Fatalf("returned %d rows, want 2", len(activities))
	}
	for _, a := range activities {
		if a.GroupId != nil && *a.GroupId == iid && a.RanByUserId != nil && *a.RanByUserId == hid {
			t.Errorf("hidden actor's isolated-group activity should be dropped")
		}
	}
}

// A nil resolver adds no predicate (backward compatible): every matching activity is
// returned and counted.
func TestGetPagedActivitiesNilResolverUnfiltered(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	group := models.Group{Name: "act-nil-grp"}
	db.Create(&group)
	user := models.User{Username: "act-nil-user", Password: "x"}
	db.Create(&user)

	gid := group.ID
	uid := user.ID
	for i := 0; i < 2; i++ {
		db.Create(&models.SystemTask{
			Type:                 models.QUICK_SCAN,
			Status:               models.SYSTEM_TASK_SUCCEEDED,
			AssociatedEntityType: models.NOOP_ENTITY_TYPE,
			GroupId:              &gid,
			RanByUserId:          &uid,
			StartedAt:            time.Now(),
		})
	}

	command := commands.PagedActivityRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      10,
			OrderBy:       "started_at",
			SortDirection: commands.DESCENDING,
		},
		GroupIds: []uint{group.ID},
	}

	activities, count, err := NewSystemTaskRepository(nil).GetPagedActivities(command, nil)
	if err != nil {
		t.Fatalf("GetPagedActivities: %v", err)
	}
	if count != 2 || len(activities) != 2 {
		t.Errorf("nil resolver should return all rows unfiltered, got count=%d rows=%d", count, len(activities))
	}
}
