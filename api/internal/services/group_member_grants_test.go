package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"testing"
)

// setMemberGrants writes a membership's individual grants directly, bypassing the
// endpoint's ceiling validation. Used to build states the API would reject (a
// disjoint role/member pair) so the resolver's own behavior is pinned rather than
// assumed away by the write-side guard.
func setMemberGrants(t *testing.T, userId uint, groupId uint, categoryIds []uint, tagIds []uint) {
	t.Helper()

	err := repositories.NewGroupMemberRepository(nil).ReplaceMemberGrants(userId, groupId, categoryIds, tagIds)
	if err != nil {
		t.Fatalf("ReplaceMemberGrants: %v", err)
	}
}

// setRequiresIndividualGrants flips a group role's fail-closed toggles.
//
// It writes through the REPOSITORY, which deliberately bypasses the grant cache
// (only RoleService evicts), so the manual clear here is correct — it stands in
// for the eviction the service would have done. Do NOT copy this clear into a
// test that exercises the service path: see
// TestRoleUpdateEvictsStaleGrantsWithoutManualCacheClear.
func setRequiresIndividualGrants(t *testing.T, roleId uint, categories bool, tags bool) {
	t.Helper()

	err := repositories.NewRoleRepository(nil).SetGroupRoleIndividualGrantConfig(roleId, categories, tags)
	if err != nil {
		t.Fatalf("SetGroupRoleIndividualGrantConfig: %v", err)
	}
	clearGroupRoleGrantCacheAll()
}

// assertAllowedCategories asserts the resolved category set for a member.
func assertAllowedCategories(t *testing.T, userId uint, groupId uint, wantRestricted bool, want []uint) {
	t.Helper()

	allowed, unrestricted, err := NewPermissionService(nil).GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupCategoryIdsForUser: %v", err)
	}

	gotRestricted := !unrestricted
	if gotRestricted != wantRestricted {
		t.Fatalf("expected restricted=%v, got restricted=%v", wantRestricted, gotRestricted)
	}
	if !wantRestricted {
		return
	}

	got := make([]uint, 0, len(allowed))
	for id := range allowed {
		got = append(got, id)
	}
	slices.Sort(got)
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Errorf("allowed categories = %v, want %v", got, want)
	}
}

// --- Intersection matrix -----------------------------------------------------

func TestMemberGrants_RoleOnlyRestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-role-only", []uint{catA, catB}, nil)

	// No individual assignment: the member falls back to the role's full set.
	assertAllowedCategories(t, userId, groupId, true, []uint{catA, catB})
}

func TestMemberGrants_MemberOnlyRestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	makeCategory(t, "B")
	// Role grants nothing => unrestricted ceiling; the member layer alone narrows.
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-member-only", nil, nil)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)

	assertAllowedCategories(t, userId, groupId, true, []uint{catA})
}

func TestMemberGrants_BothRestrictedIntersects(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	catC := makeCategory(t, "C")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-both", []uint{catA, catB, catC}, nil)
	setMemberGrants(t, userId, groupId, []uint{catB}, nil)

	// The individual assignment narrows within the role's ceiling.
	assertAllowedCategories(t, userId, groupId, true, []uint{catB})
}

func TestMemberGrants_DisjointResolvesToNothing(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catD := makeCategory(t, "D")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-disjoint", []uint{catA}, nil)
	setMemberGrants(t, userId, groupId, []uint{catD}, nil)

	// The endpoint rejects this state, but a role narrowed AFTER the member was
	// configured can produce it. Intersection must fail closed, not fall back.
	assertAllowedCategories(t, userId, groupId, true, []uint{})
}

func TestMemberGrants_NeitherRestrictedIsUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	makeCategory(t, "A")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-open", nil, nil)

	assertAllowedCategories(t, userId, groupId, false, nil)
}

func TestMemberGrants_CategoriesAndTagsAreIndependent(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	makeCategory(t, "B")
	tagA := makeTag(t, "tag-a")

	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-independent", nil, nil)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)

	assertAllowedCategories(t, userId, groupId, true, []uint{catA})

	// Tags were never configured on either layer, so they stay unrestricted.
	_, unrestricted, err := NewPermissionService(nil).GetGroupTagIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupTagIdsForUser: %v", err)
	}
	if !unrestricted {
		t.Error("expected tags to remain unrestricted when only categories were assigned")
	}

	// Restricting tags must not disturb the category resolution.
	setMemberGrants(t, userId, groupId, []uint{catA}, []uint{tagA})
	assertAllowedCategories(t, userId, groupId, true, []uint{catA})
}

func TestMemberGrants_NonMemberIsUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	_, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-owner", nil, nil)

	outsider := models.User{Username: "m-outsider", Password: "password"}
	if err := repositories.GetDB().Create(&outsider).Error; err != nil {
		t.Fatalf("seed outsider: %v", err)
	}

	assertAllowedCategories(t, outsider.ID, groupId, false, nil)
}

// --- Fail-closed toggle ------------------------------------------------------

func TestMemberGrants_RequiresIndividualHidesUnassignedMember(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-requires", []uint{catA, catB}, nil)
	setRequiresIndividualGrants(t, roleId, true, false)

	// Without the toggle this member would see the role's full set; with it, an
	// unassigned member sees nothing.
	assertAllowedCategories(t, userId, groupId, true, []uint{})
}

func TestMemberGrants_RequiresIndividualHonorsAssignment(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-requires-set", []uint{catA, catB}, nil)
	setRequiresIndividualGrants(t, roleId, true, false)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)

	assertAllowedCategories(t, userId, groupId, true, []uint{catA})
}

func TestMemberGrants_RequiresIndividualOffFallsBackToRole(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-requires-off", []uint{catA}, nil)
	setRequiresIndividualGrants(t, roleId, false, false)

	assertAllowedCategories(t, userId, groupId, true, []uint{catA})
}

func TestMemberGrants_RequiresIndividualTagsIndependentOfCategories(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-requires-tags", []uint{catA}, nil)
	setRequiresIndividualGrants(t, roleId, false, true)

	// Categories fall back to the role; tags fail closed.
	assertAllowedCategories(t, userId, groupId, true, []uint{catA})

	allowedTags, unrestricted, err := NewPermissionService(nil).GetGroupTagIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupTagIdsForUser: %v", err)
	}
	if unrestricted {
		t.Fatal("expected tags restricted when the role requires individual tag grants")
	}
	if len(allowedTags) != 0 {
		t.Errorf("expected no allowed tags, got %v", allowedTags)
	}
}

// --- Fail-closed restriction flag --------------------------------------------

func TestMemberGrants_StaysRestrictedAfterGrantedCategoryDeleted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	makeCategory(t, "B")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-cascade", nil, nil)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)

	// Delete the member's only granted category, then clear its grant rows the way
	// an ON DELETE CASCADE would.
	db := repositories.GetDB()
	if err := db.Delete(&models.Category{}, catA).Error; err != nil {
		t.Fatalf("delete category: %v", err)
	}
	if err := db.Where("category_id = ?", catA).Delete(&models.GroupMemberCategoryGrant{}).Error; err != nil {
		t.Fatalf("cascade grant rows: %v", err)
	}

	// The membership must stay restricted (seeing nothing), NOT widen back to
	// see-all just because its rows are gone.
	assertAllowedCategories(t, userId, groupId, true, []uint{})
}

// --- Lifecycle: no resurrection ----------------------------------------------

func TestMemberGrants_RemovedMemberGrantsDoNotResurrectOnRejoin(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-rejoin", nil, nil)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)

	db := repositories.GetDB()

	// Remove the membership the way the teardown paths do (raw delete), plus the
	// grant cleanup those paths now perform.
	if err := db.Delete(&models.GroupMember{}, "user_id = ? AND group_id = ?", userId, groupId).Error; err != nil {
		t.Fatalf("remove member: %v", err)
	}
	if err := repositories.DeleteMemberGrants(db, userId, groupId); err != nil {
		t.Fatalf("DeleteMemberGrants: %v", err)
	}

	// Re-add the same user to the same group.
	rejoined := models.GroupMember{GroupID: groupId, UserID: userId, GroupRoleID: &roleId}
	if err := db.Create(&rejoined).Error; err != nil {
		t.Fatalf("re-add member: %v", err)
	}

	// They must come back with NO assignment, not their revoked one.
	assertAllowedCategories(t, userId, groupId, false, nil)

	categoryIds, err := repositories.NewGroupMemberRepository(nil).GetMemberCategoryGrantIds(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds: %v", err)
	}
	if len(categoryIds) != 0 {
		t.Errorf("expected no resurrected grants, got %v (catB=%d unused)", categoryIds, catB)
	}
}

func TestMemberGrants_OrphanCleanupKeepsRetainedMembers(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	keptUserId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-kept", nil, nil)

	db := repositories.GetDB()
	removed := models.User{Username: "m-removed", Password: "password"}
	if err := db.Create(&removed).Error; err != nil {
		t.Fatalf("seed removed user: %v", err)
	}
	removedMember := models.GroupMember{GroupID: groupId, UserID: removed.ID, GroupRoleID: &roleId}
	if err := db.Create(&removedMember).Error; err != nil {
		t.Fatalf("seed removed member: %v", err)
	}

	setMemberGrants(t, keptUserId, groupId, []uint{catA}, nil)
	setMemberGrants(t, removed.ID, groupId, []uint{catA}, nil)

	// Drop one membership, then run the cleanup UpdateGroup performs after its
	// wholesale roster replace.
	if err := db.Delete(&models.GroupMember{}, "user_id = ? AND group_id = ?", removed.ID, groupId).Error; err != nil {
		t.Fatalf("remove member: %v", err)
	}
	if err := repositories.DeleteOrphanedMemberGrants(db, groupId); err != nil {
		t.Fatalf("DeleteOrphanedMemberGrants: %v", err)
	}

	groupMemberRepository := repositories.NewGroupMemberRepository(nil)

	keptIds, err := groupMemberRepository.GetMemberCategoryGrantIds(keptUserId, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds(kept): %v", err)
	}
	if !slices.Equal(keptIds, []uint{catA}) {
		t.Errorf("retained member lost grants: got %v, want [%d]", keptIds, catA)
	}

	orphanIds, err := groupMemberRepository.GetMemberCategoryGrantIds(removed.ID, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds(removed): %v", err)
	}
	if len(orphanIds) != 0 {
		t.Errorf("removed member's grants were not cleaned up: %v", orphanIds)
	}
}

// --- Write-side ceiling validation -------------------------------------------

func TestUpdateMemberGrants_RejectsIdsOutsideRoleCeiling(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catD := makeCategory(t, "D")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-ceiling", []uint{catA}, nil)

	err := NewGroupService(nil).UpdateMemberGrants(groupId, userId, commands.UpdateGroupMemberGrantsCommand{
		CategoryIds: []uint{catD},
	})

	var violation *GrantCeilingViolation
	if !errors.As(err, &violation) {
		t.Fatalf("expected GrantCeilingViolation, got %v", err)
	}
	if !slices.Equal(violation.CategoryIds, []uint{catD}) {
		t.Errorf("violation should name the offending id: got %v, want [%d]", violation.CategoryIds, catD)
	}

	// Nothing may have been written.
	categoryIds, err := repositories.NewGroupMemberRepository(nil).GetMemberCategoryGrantIds(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds: %v", err)
	}
	if len(categoryIds) != 0 {
		t.Errorf("rejected write must not persist, got %v", categoryIds)
	}
}

func TestUpdateMemberGrants_AcceptsIdsWithinRoleCeiling(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-within", []uint{catA, catB}, nil)

	err := NewGroupService(nil).UpdateMemberGrants(groupId, userId, commands.UpdateGroupMemberGrantsCommand{
		CategoryIds: []uint{catB},
	})
	if err != nil {
		t.Fatalf("UpdateMemberGrants: %v", err)
	}

	assertAllowedCategories(t, userId, groupId, true, []uint{catB})
}

func TestUpdateMemberGrants_UnrestrictedRoleImposesNoCeiling(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-no-ceiling", nil, nil)

	err := NewGroupService(nil).UpdateMemberGrants(groupId, userId, commands.UpdateGroupMemberGrantsCommand{
		CategoryIds: []uint{catA},
	})
	if err != nil {
		t.Fatalf("UpdateMemberGrants: %v", err)
	}

	assertAllowedCategories(t, userId, groupId, true, []uint{catA})
}

func TestUpdateMemberGrants_RejectsNonExistentCategory(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-missing-cat", nil, nil)

	err := NewGroupService(nil).UpdateMemberGrants(groupId, userId, commands.UpdateGroupMemberGrantsCommand{
		CategoryIds: []uint{99999},
	})
	if !errors.Is(err, ErrInvalidGrant) {
		t.Fatalf("expected ErrInvalidGrant, got %v", err)
	}
}

func TestUpdateMemberGrants_RejectsNonMember(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	_, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-nonmember-owner", nil, nil)

	outsider := models.User{Username: "m-nonmember", Password: "password"}
	if err := repositories.GetDB().Create(&outsider).Error; err != nil {
		t.Fatalf("seed outsider: %v", err)
	}

	err := NewGroupService(nil).UpdateMemberGrants(groupId, outsider.ID, commands.UpdateGroupMemberGrantsCommand{})
	if !errors.Is(err, ErrMemberNotInGroup) {
		t.Fatalf("expected ErrMemberNotInGroup, got %v", err)
	}
}

func TestUpdateMemberGrants_EmptySelectionClearsRestriction(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-clear", []uint{catA, catB}, nil)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)
	assertAllowedCategories(t, userId, groupId, true, []uint{catA})

	err := NewGroupService(nil).UpdateMemberGrants(groupId, userId, commands.UpdateGroupMemberGrantsCommand{})
	if err != nil {
		t.Fatalf("UpdateMemberGrants: %v", err)
	}

	// Clearing the individual assignment hands the member back to their role.
	assertAllowedCategories(t, userId, groupId, true, []uint{catA, catB})
}

// --- Cache eviction on the real update path ---------------------------------

// TestRoleUpdateEvictsStaleGrantsWithoutManualCacheClear pins that narrowing a
// role through RoleService takes effect immediately.
//
// The grant cache is keyed by group-role id and is only evicted by
// RoleService.UpdateRole / DeleteRole. Every other test in this file writes
// grants through the repository and clears the cache by hand — which means that
// if the service's eviction regressed, they would all still pass while
// production kept serving the revoked categories.
//
// So this test MUST NOT call clearGroupRoleGrantCacheAll after the update: the
// whole point is that the service does it. It resolves the member's set first to
// force the pre-update grants into the cache, so a missing eviction actually
// shows up as a stale read rather than a cold miss.
func TestRoleUpdateEvictsStaleGrantsWithoutManualCacheClear(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-evict", []uint{catA, catB}, nil)

	// Prime the cache with the pre-update grants.
	assertAllowedCategories(t, userId, groupId, true, []uint{catA, catB})

	// Narrow the role through the production path.
	_, err := NewRoleService(nil).UpdateRole(roleId, commands.UpsertRoleCommand{
		Name:           "Grant Role m-evict",
		Description:    "",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{catA},
	})
	if err != nil {
		t.Fatalf("UpdateRole: %v", err)
	}

	// No manual cache clear here, deliberately.
	assertAllowedCategories(t, userId, groupId, true, []uint{catA})
}

// TestRoleUpdateEvictsRequiresIndividualToggle covers the same eviction for the
// fail-closed toggle, which is written by a SEPARATE repository call inside
// RoleService.UpdateRole's transaction and so could plausibly be missed by an
// eviction that only considered the grant rows.
func TestRoleUpdateEvictsRequiresIndividualToggle(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "m-evict-toggle", []uint{catA}, nil)

	// Prime: the unassigned member currently inherits the role's set.
	assertAllowedCategories(t, userId, groupId, true, []uint{catA})

	_, err := NewRoleService(nil).UpdateRole(roleId, commands.UpsertRoleCommand{
		Name:                             "Grant Role m-evict-toggle",
		Description:                      "",
		Scope:                            permissions.ScopeGroup,
		Permissions:                      []string{permissions.GroupReceiptsRead},
		CategoryGrants:                   []uint{catA},
		RequiresIndividualCategoryGrants: true,
	})
	if err != nil {
		t.Fatalf("UpdateRole: %v", err)
	}

	// Again no manual clear: the member must now fail closed immediately.
	assertAllowedCategories(t, userId, groupId, true, []uint{})
}

// --- Delete-path cleanup ----------------------------------------------------

// TestDeleteGroupRemovesMemberGrants pins the group-teardown half of the
// resurrection hazard. GroupService.DeleteGroup removes memberships with a raw
// delete that cascades nothing, so the grant rows must be cleared explicitly.
func TestDeleteGroupRemovesMemberGrants(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	catA := makeCategory(t, "A")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "m-group-delete", nil, nil)
	setMemberGrants(t, userId, groupId, []uint{catA}, nil)

	before, err := repositories.NewGroupMemberRepository(nil).GetMemberCategoryGrantIds(userId, groupId)
	if err != nil {
		t.Fatalf("GetMemberCategoryGrantIds: %v", err)
	}
	if len(before) != 1 {
		t.Fatalf("expected a seeded grant before delete, got %v", before)
	}

	if err := NewGroupService(nil).DeleteGroup(utils.UintToString(groupId), true); err != nil {
		t.Fatalf("DeleteGroup: %v", err)
	}

	var remaining int64
	err = repositories.GetDB().Model(&models.GroupMemberCategoryGrant{}).
		Where("group_id = ?", groupId).Count(&remaining).Error
	if err != nil {
		t.Fatalf("count grants: %v", err)
	}
	if remaining != 0 {
		t.Errorf("deleting the group left %d member grant rows behind", remaining)
	}
}
