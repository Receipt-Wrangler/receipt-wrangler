package repositories

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"testing"
)

// seedGrantGroup creates a group with two members and gives each an individual
// category grant, returning the group id, the two user ids, and the category id.
// grantFixture is the seeded world these tests act on. A struct rather than a run
// of bare uints so the tag additions stay readable at the call sites.
type grantFixture struct {
	groupId       uint
	keptUserId    uint
	removedUserId uint
	categoryId    uint
	tagId         uint
}

// seedGrantGroup seeds a group whose members are restricted on BOTH resources.
// Tags matter as much as categories here: UpdateGroup restores each restriction
// flag in its own statement, so a category-only fixture would let a regression
// that drops the tag flag pass every test in this file.
func seedGrantGroup(t *testing.T) grantFixture {
	t.Helper()
	db := GetDB()

	category := models.Category{Name: "Child A"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}

	tag := models.Tag{Name: "Respite"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("seed tag: %v", err)
	}

	keptUser := models.User{Username: "kept-user", Password: "password"}
	if err := db.Create(&keptUser).Error; err != nil {
		t.Fatalf("seed kept user: %v", err)
	}
	removedUser := models.User{Username: "removed-user", Password: "password"}
	if err := db.Create(&removedUser).Error; err != nil {
		t.Fatalf("seed removed user: %v", err)
	}

	group := models.Group{Name: "Agency"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	for _, userId := range []uint{keptUser.ID, removedUser.ID} {
		member := models.GroupMember{GroupID: group.ID, UserID: userId}
		if err := db.Create(&member).Error; err != nil {
			t.Fatalf("seed member: %v", err)
		}
		repository := NewGroupMemberRepository(nil)
		if err := repository.ReplaceMemberGrants(userId, group.ID, []uint{category.ID}, []uint{tag.ID}); err != nil {
			t.Fatalf("seed member grants: %v", err)
		}
	}

	return grantFixture{
		groupId:       group.ID,
		keptUserId:    keptUser.ID,
		removedUserId: removedUser.ID,
		categoryId:    category.ID,
		tagId:         tag.ID,
	}
}

func memberCategoryGrantIds(t *testing.T, userId uint, groupId uint) []uint {
	t.Helper()

	ids, err := NewGroupMemberRepository(nil).GetMemberCategoryGrantIds(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds: %v", err)
	}
	return ids
}

func memberTagGrantIds(t *testing.T, userId uint, groupId uint) []uint {
	t.Helper()

	ids, err := NewGroupMemberRepository(nil).GetMemberTagGrantIds(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberTagGrantIds: %v", err)
	}
	return ids
}

// memberGrantsRestricted reports both flags separately — UpdateGroup restores them
// independently, so collapsing them would hide a regression in either one.
func memberGrantsRestricted(t *testing.T, userId uint, groupId uint) (bool, bool) {
	t.Helper()

	member, err := NewGroupMemberRepository(nil).GetMemberGrantContext(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberGrantContext: %v", err)
	}
	return member.CategoryGrantsRestricted, member.TagGrantsRestricted
}

// TestUpdateGroupClearsRemovedMembersGrants exercises the REAL UpdateGroup path
// (not the cleanup helper in isolation): a roster replace that drops a member
// must take their individual grants with it, or re-adding that user later
// silently restores visibility that was deliberately revoked.
func TestUpdateGroupClearsRemovedMembersGrants(t *testing.T) {
	defer TruncateTestDb()

	fixture := seedGrantGroup(t)
	repository := NewGroupRepository(nil)

	_, err := repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: fixture.keptUserId, GroupID: fixture.groupId},
		},
	}, utils.UintToString(fixture.groupId))
	if err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}

	if got := memberCategoryGrantIds(t, fixture.removedUserId, fixture.groupId); len(got) != 0 {
		t.Errorf("removed member's category grants survived the roster replace: %v", got)
	}
	if got := memberTagGrantIds(t, fixture.removedUserId, fixture.groupId); len(got) != 0 {
		t.Errorf("removed member's tag grants survived the roster replace: %v", got)
	}

	if got := memberCategoryGrantIds(t, fixture.keptUserId, fixture.groupId); !slices.Equal(got, []uint{fixture.categoryId}) {
		t.Errorf("retained member lost their category grants: got %v, want [%d]", got, fixture.categoryId)
	}
	if got := memberTagGrantIds(t, fixture.keptUserId, fixture.groupId); !slices.Equal(got, []uint{fixture.tagId}) {
		t.Errorf("retained member lost their tag grants: got %v, want [%d]", got, fixture.tagId)
	}
}

// TestUpdateGroupDoesNotResurrectGrantsOnRejoin is the end-to-end shape of the
// hazard: remove a member, then add them back, and confirm they return with no
// assignment rather than their revoked one.
func TestUpdateGroupDoesNotResurrectGrantsOnRejoin(t *testing.T) {
	defer TruncateTestDb()

	fixture := seedGrantGroup(t)
	repository := NewGroupRepository(nil)
	groupIdString := utils.UintToString(fixture.groupId)

	_, err := repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: fixture.keptUserId, GroupID: fixture.groupId},
		},
	}, groupIdString)
	if err != nil {
		t.Fatalf("UpdateGroup (remove): %v", err)
	}

	_, err = repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: fixture.keptUserId, GroupID: fixture.groupId},
			{UserID: fixture.removedUserId, GroupID: fixture.groupId},
		},
	}, groupIdString)
	if err != nil {
		t.Fatalf("UpdateGroup (re-add): %v", err)
	}

	if got := memberCategoryGrantIds(t, fixture.removedUserId, fixture.groupId); len(got) != 0 {
		t.Errorf("rejoined member's revoked category grants were resurrected: %v", got)
	}
	if got := memberTagGrantIds(t, fixture.removedUserId, fixture.groupId); len(got) != 0 {
		t.Errorf("rejoined member's revoked tag grants were resurrected: %v", got)
	}
	categoryRestricted, tagRestricted := memberGrantsRestricted(t, fixture.removedUserId, fixture.groupId)
	if categoryRestricted || tagRestricted {
		t.Errorf(
			"rejoined member should come back unrestricted, not carrying stale restriction flags (category=%v, tag=%v)",
			categoryRestricted, tagRestricted,
		)
	}
}

// TestUpdateGroupPreservesRetainedMemberGrantFlags pins that a plain group edit
// (a rename, roster unchanged) neither drops a member's grants nor resets the
// fail-closed flag that keeps them restricted.
func TestUpdateGroupPreservesRetainedMemberGrantFlags(t *testing.T) {
	defer TruncateTestDb()

	fixture := seedGrantGroup(t)
	repository := NewGroupRepository(nil)

	_, err := repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency Renamed",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: fixture.keptUserId, GroupID: fixture.groupId},
			{UserID: fixture.removedUserId, GroupID: fixture.groupId},
		},
	}, utils.UintToString(fixture.groupId))
	if err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}

	for _, userId := range []uint{fixture.keptUserId, fixture.removedUserId} {
		if got := memberCategoryGrantIds(t, userId, fixture.groupId); !slices.Equal(got, []uint{fixture.categoryId}) {
			t.Errorf("member %d lost category grants on a rename: got %v, want [%d]", userId, got, fixture.categoryId)
		}
		if got := memberTagGrantIds(t, userId, fixture.groupId); !slices.Equal(got, []uint{fixture.tagId}) {
			t.Errorf("member %d lost tag grants on a rename: got %v, want [%d]", userId, got, fixture.tagId)
		}

		categoryRestricted, tagRestricted := memberGrantsRestricted(t, userId, fixture.groupId)
		if !categoryRestricted {
			t.Errorf("member %d lost their CATEGORY restriction flag on a rename", userId)
		}
		if !tagRestricted {
			t.Errorf("member %d lost their TAG restriction flag on a rename", userId)
		}
	}
}
