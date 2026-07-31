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
func seedGrantGroup(t *testing.T) (uint, uint, uint, uint) {
	t.Helper()
	db := GetDB()

	category := models.Category{Name: "Child A"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("seed category: %v", err)
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
		if err := repository.ReplaceMemberGrants(userId, group.ID, []uint{category.ID}, nil); err != nil {
			t.Fatalf("seed member grants: %v", err)
		}
	}

	return group.ID, keptUser.ID, removedUser.ID, category.ID
}

func memberCategoryGrantIds(t *testing.T, userId uint, groupId uint) []uint {
	t.Helper()

	ids, err := NewGroupMemberRepository(nil).GetMemberCategoryGrantIds(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds: %v", err)
	}
	return ids
}

func memberGrantsRestricted(t *testing.T, userId uint, groupId uint) bool {
	t.Helper()

	member, err := NewGroupMemberRepository(nil).GetMemberGrantContext(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberGrantContext: %v", err)
	}
	return member.CategoryGrantsRestricted
}

// TestUpdateGroupClearsRemovedMembersGrants exercises the REAL UpdateGroup path
// (not the cleanup helper in isolation): a roster replace that drops a member
// must take their individual grants with it, or re-adding that user later
// silently restores visibility that was deliberately revoked.
func TestUpdateGroupClearsRemovedMembersGrants(t *testing.T) {
	defer TruncateTestDb()

	groupId, keptUserId, removedUserId, categoryId := seedGrantGroup(t)
	repository := NewGroupRepository(nil)

	_, err := repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: keptUserId, GroupID: groupId},
		},
	}, utils.UintToString(groupId))
	if err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}

	if got := memberCategoryGrantIds(t, removedUserId, groupId); len(got) != 0 {
		t.Errorf("removed member's grants survived the roster replace: %v", got)
	}

	if got := memberCategoryGrantIds(t, keptUserId, groupId); !slices.Equal(got, []uint{categoryId}) {
		t.Errorf("retained member lost their grants: got %v, want [%d]", got, categoryId)
	}
}

// TestUpdateGroupDoesNotResurrectGrantsOnRejoin is the end-to-end shape of the
// hazard: remove a member, then add them back, and confirm they return with no
// assignment rather than their revoked one.
func TestUpdateGroupDoesNotResurrectGrantsOnRejoin(t *testing.T) {
	defer TruncateTestDb()

	groupId, keptUserId, removedUserId, _ := seedGrantGroup(t)
	repository := NewGroupRepository(nil)
	groupIdString := utils.UintToString(groupId)

	_, err := repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: keptUserId, GroupID: groupId},
		},
	}, groupIdString)
	if err != nil {
		t.Fatalf("UpdateGroup (remove): %v", err)
	}

	_, err = repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: keptUserId, GroupID: groupId},
			{UserID: removedUserId, GroupID: groupId},
		},
	}, groupIdString)
	if err != nil {
		t.Fatalf("UpdateGroup (re-add): %v", err)
	}

	if got := memberCategoryGrantIds(t, removedUserId, groupId); len(got) != 0 {
		t.Errorf("rejoined member's revoked grants were resurrected: %v", got)
	}
	if memberGrantsRestricted(t, removedUserId, groupId) {
		t.Error("rejoined member should come back unrestricted, not carrying a stale restriction flag")
	}
}

// TestUpdateGroupPreservesRetainedMemberGrantFlags pins that a plain group edit
// (a rename, roster unchanged) neither drops a member's grants nor resets the
// fail-closed flag that keeps them restricted.
func TestUpdateGroupPreservesRetainedMemberGrantFlags(t *testing.T) {
	defer TruncateTestDb()

	groupId, keptUserId, removedUserId, categoryId := seedGrantGroup(t)
	repository := NewGroupRepository(nil)

	_, err := repository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "Agency Renamed",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: keptUserId, GroupID: groupId},
			{UserID: removedUserId, GroupID: groupId},
		},
	}, utils.UintToString(groupId))
	if err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}

	for _, userId := range []uint{keptUserId, removedUserId} {
		if got := memberCategoryGrantIds(t, userId, groupId); !slices.Equal(got, []uint{categoryId}) {
			t.Errorf("member %d lost grants on a rename: got %v, want [%d]", userId, got, categoryId)
		}
		if !memberGrantsRestricted(t, userId, groupId) {
			t.Errorf("member %d lost their restriction flag on a rename", userId)
		}
	}
}
