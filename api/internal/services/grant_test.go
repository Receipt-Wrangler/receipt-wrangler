package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"slices"
	"testing"
)

// seedMemberWithGroupRoleGrants creates a group, a group role granting
// receipts.read plus the given category/tag grants, and a member assigned to it.
// Returns the user id, group id, and role id.
func seedMemberWithGroupRoleGrants(t *testing.T, username string, categoryGrantIds []uint, tagGrantIds []uint) (uint, uint, uint) {
	t.Helper()
	db := repositories.GetDB()

	group := models.Group{Name: "grant-group-" + username}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole("Grant Role "+username, "", []string{permissions.GroupReceiptsRead}, categoryGrantIds, tagGrantIds)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: username, Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}

	return user.ID, group.ID, role.ID
}

func makeCategory(t *testing.T, name string) uint {
	t.Helper()
	category := models.Category{Name: name}
	if err := repositories.GetDB().Create(&category).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}
	return category.ID
}

func makeTag(t *testing.T, name string) uint {
	t.Helper()
	tag := models.Tag{Name: name}
	if err := repositories.GetDB().Create(&tag).Error; err != nil {
		t.Fatalf("seed tag: %v", err)
	}
	return tag.ID
}

func TestGetGroupCategoryIdsUnrestrictedWhenNoGrants(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-open", nil, nil)
	service := NewPermissionService(nil)

	allowed, unrestricted, err := service.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupCategoryIdsForUser: %v", err)
	}
	if !unrestricted {
		t.Error("expected unrestricted when role has no category grants")
	}
	if allowed != nil {
		t.Errorf("expected nil allowed set when unrestricted, got %v", allowed)
	}
}

func TestGetGroupCategoryIdsRestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-restricted", []uint{cat1, cat2}, nil)
	service := NewPermissionService(nil)

	allowed, unrestricted, err := service.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupCategoryIdsForUser: %v", err)
	}
	if unrestricted {
		t.Error("expected restricted when role has category grants")
	}
	if len(allowed) != 2 {
		t.Fatalf("expected 2 allowed categories, got %d", len(allowed))
	}
	if _, ok := allowed[cat1]; !ok {
		t.Error("expected cat1 in allowed set")
	}
	if _, ok := allowed[cat2]; !ok {
		t.Error("expected cat2 in allowed set")
	}
}

func TestGetGroupCategoryIdsNonMemberUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	cat1 := makeCategory(t, "Groceries")
	_, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-member", []uint{cat1}, nil)

	// A user who is NOT a member of the group resolves to unrestricted: grants
	// only narrow access within an already-permitted group, never grant it.
	outsider := models.User{Username: "outsider", Password: "password"}
	if err := repositories.GetDB().Create(&outsider).Error; err != nil {
		t.Fatalf("seed outsider: %v", err)
	}

	service := NewPermissionService(nil)
	_, unrestricted, err := service.GetGroupCategoryIdsForUser(outsider.ID, groupId)
	if err != nil {
		t.Fatalf("GetGroupCategoryIdsForUser: %v", err)
	}
	if !unrestricted {
		t.Error("expected a non-member to resolve as unrestricted")
	}
}

func TestCategoryRestrictedLeavesTagsUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	cat1 := makeCategory(t, "Groceries")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-cat-only", []uint{cat1}, nil)
	service := NewPermissionService(nil)

	_, catUnrestricted, err := service.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("category resolve: %v", err)
	}
	if catUnrestricted {
		t.Error("expected categories restricted")
	}

	_, tagUnrestricted, err := service.GetGroupTagIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("tag resolve: %v", err)
	}
	if !tagUnrestricted {
		t.Error("expected tags unrestricted when role grants no tags")
	}
}

func TestGetVisibleCategoriesForUserFilters(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	cat1 := makeCategory(t, "Groceries")
	makeCategory(t, "Utilities")
	cat3 := makeCategory(t, "Travel")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-vis", []uint{cat1, cat3}, nil)
	service := NewPermissionService(nil)

	allCategories, err := repositories.NewCategoryRepository(nil).GetAllCategories("*")
	if err != nil {
		t.Fatalf("get all categories: %v", err)
	}

	visible, err := service.GetVisibleCategoriesForUser(userId, groupId, allCategories)
	if err != nil {
		t.Fatalf("GetVisibleCategoriesForUser: %v", err)
	}
	if len(visible) != 2 {
		t.Fatalf("expected 2 visible categories, got %d", len(visible))
	}

	gotIds := []uint{visible[0].ID, visible[1].ID}
	slices.Sort(gotIds)
	want := []uint{cat1, cat3}
	slices.Sort(want)
	if !slices.Equal(gotIds, want) {
		t.Errorf("visible ids = %v, want %v", gotIds, want)
	}
}

func TestGetVisibleCategoriesUnrestrictedPassThrough(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	makeCategory(t, "Groceries")
	makeCategory(t, "Utilities")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-open2", nil, nil)
	service := NewPermissionService(nil)

	allCategories, err := repositories.NewCategoryRepository(nil).GetAllCategories("*")
	if err != nil {
		t.Fatalf("get all categories: %v", err)
	}

	visible, err := service.GetVisibleCategoriesForUser(userId, groupId, allCategories)
	if err != nil {
		t.Fatalf("GetVisibleCategoriesForUser: %v", err)
	}
	if len(visible) != len(allCategories) {
		t.Errorf("expected pass-through of all %d categories, got %d", len(allCategories), len(visible))
	}
}

func TestGetVisibleTagsForUserFilters(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	tag1 := makeTag(t, "Reimbursable")
	makeTag(t, "Personal")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-tagvis", nil, []uint{tag1})
	service := NewPermissionService(nil)

	allTags, err := repositories.NewTagsRepository(nil).GetAllTags("*")
	if err != nil {
		t.Fatalf("get all tags: %v", err)
	}

	visible, err := service.GetVisibleTagsForUser(userId, groupId, allTags)
	if err != nil {
		t.Fatalf("GetVisibleTagsForUser: %v", err)
	}
	if len(visible) != 1 || visible[0].ID != tag1 {
		t.Errorf("expected only tag1 visible, got %v", visible)
	}
}

func TestGrantCacheInvalidatedOnRoleUpdate(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "u-cache", []uint{cat1}, nil)

	permissionService := NewPermissionService(nil)

	// Populate the cache with the initial grant set.
	allowed, _, err := permissionService.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("initial resolve: %v", err)
	}
	if _, ok := allowed[cat1]; !ok || len(allowed) != 1 {
		t.Fatalf("expected initial allowed {cat1}, got %v", allowed)
	}

	// Replace the grants via the service — should evict the grant cache.
	roleService := NewRoleService(nil)
	_, err = roleService.UpdateRole(roleId, commands.UpsertRoleCommand{
		Name:           "Grant Role u-cache",
		Scope:          permissions.ScopeGroup,
		Permissions:    []string{permissions.GroupReceiptsRead},
		CategoryGrants: []uint{cat2},
	})
	if err != nil {
		t.Fatalf("update role: %v", err)
	}

	allowed, _, err = permissionService.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("re-resolve: %v", err)
	}
	if _, ok := allowed[cat2]; !ok || len(allowed) != 1 {
		t.Errorf("expected allowed {cat2} after update, got %v", allowed)
	}
}

func TestGrantCacheReturnsCachedUntilCleared(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	userId, groupId, roleId := seedMemberWithGroupRoleGrants(t, "u-cachehit", []uint{cat1}, nil)
	service := NewPermissionService(nil)

	// Populate the cache.
	if _, _, err := service.GetGroupCategoryIdsForUser(userId, groupId); err != nil {
		t.Fatalf("populate: %v", err)
	}

	// Mutate grant rows directly, bypassing the service so the cache is NOT
	// invalidated.
	db := repositories.GetDB()
	db.Where("group_role_id = ?", roleId).Delete(&models.GroupRoleCategoryGrant{})
	db.Omit("Category").Create(&models.GroupRoleCategoryGrant{GroupRoleID: roleId, CategoryID: cat2})

	allowed, _, err := service.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("cached resolve: %v", err)
	}
	if _, ok := allowed[cat1]; !ok || len(allowed) != 1 {
		t.Errorf("expected cached {cat1}, got %v", allowed)
	}

	// After an explicit eviction, the fresh DB value {cat2} is read.
	clearGroupRoleGrantCache(roleId)
	allowed, _, err = service.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("fresh resolve: %v", err)
	}
	if _, ok := allowed[cat2]; !ok || len(allowed) != 1 {
		t.Errorf("expected fresh {cat2}, got %v", allowed)
	}
}
