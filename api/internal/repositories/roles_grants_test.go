package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"testing"
)

// makeTestCategory inserts a category and returns its id.
func makeTestCategory(t *testing.T, name string) uint {
	category := models.Category{Name: name}
	if err := GetDB().Create(&category).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	return category.ID
}

// makeTestTag inserts a tag and returns its id.
func makeTestTag(t *testing.T, name string) uint {
	tag := models.Tag{Name: name}
	if err := GetDB().Create(&tag).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	return tag.ID
}

func TestCreateGroupRolePersistsCategoryAndTagGrants(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	cat1 := makeTestCategory(t, "Groceries")
	cat2 := makeTestCategory(t, "Utilities")
	tag1 := makeTestTag(t, "Reimbursable")

	role, err := repository.CreateGroupRole(
		"Restricted Role",
		"",
		[]string{permissions.GroupReceiptsRead},
		[]uint{cat1, cat2},
		[]uint{tag1},
		nil,
		false,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(role.CategoryGrants) != 2 {
		utils.PrintTestError(t, len(role.CategoryGrants), 2)
	}
	if len(role.TagGrants) != 1 {
		utils.PrintTestError(t, len(role.TagGrants), 1)
	}

	categoryIds, err := repository.GetGroupRoleCategoryIds(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	slices.Sort(categoryIds)
	expected := []uint{cat1, cat2}
	slices.Sort(expected)
	if !slices.Equal(categoryIds, expected) {
		utils.PrintTestError(t, categoryIds, expected)
	}

	tagIds, err := repository.GetGroupRoleTagIds(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Equal(tagIds, []uint{tag1}) {
		utils.PrintTestError(t, tagIds, []uint{tag1})
	}
}

func TestCreateGroupRoleWithNoGrantsIsUnrestricted(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	role, err := repository.CreateGroupRole("Open Role", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	categoryIds, err := repository.GetGroupRoleCategoryIds(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(categoryIds) != 0 {
		utils.PrintTestError(t, len(categoryIds), 0)
	}
}

func TestUpdateGroupRoleReplacesGrants(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	cat1 := makeTestCategory(t, "Groceries")
	cat2 := makeTestCategory(t, "Utilities")
	cat3 := makeTestCategory(t, "Travel")
	tag1 := makeTestTag(t, "Reimbursable")

	created, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReceiptsRead}, []uint{cat1, cat2}, []uint{tag1}, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Replace the grant sets entirely with cat3 and no tags.
	updated, err := repository.UpdateGroupRole(created.ID, "Role", "", []string{permissions.GroupReceiptsRead}, []uint{cat3}, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(updated.CategoryGrants) != 1 || updated.CategoryGrants[0].CategoryID != cat3 {
		utils.PrintTestError(t, updated.CategoryGrants, []uint{cat3})
	}
	if len(updated.TagGrants) != 0 {
		utils.PrintTestError(t, len(updated.TagGrants), 0)
	}

	categoryIds, err := repository.GetGroupRoleCategoryIds(created.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Equal(categoryIds, []uint{cat3}) {
		utils.PrintTestError(t, categoryIds, []uint{cat3})
	}
}

func TestDeleteGroupRoleCascadesGrants(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	cat1 := makeTestCategory(t, "Groceries")
	tag1 := makeTestTag(t, "Reimbursable")

	created, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReceiptsRead}, []uint{cat1}, []uint{tag1}, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := repository.DeleteGroupRole(created.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var categoryGrantCount int64
	GetDB().Model(&models.GroupRoleCategoryGrant{}).Where("group_role_id = ?", created.ID).Count(&categoryGrantCount)
	if categoryGrantCount != 0 {
		utils.PrintTestError(t, categoryGrantCount, 0)
	}

	var tagGrantCount int64
	GetDB().Model(&models.GroupRoleTagGrant{}).Where("group_role_id = ?", created.ID).Count(&tagGrantCount)
	if tagGrantCount != 0 {
		utils.PrintTestError(t, tagGrantCount, 0)
	}
}

func TestCategoryAndTagCountByIds(t *testing.T) {
	defer TruncateTestDb()

	cat1 := makeTestCategory(t, "Groceries")
	cat2 := makeTestCategory(t, "Utilities")
	tag1 := makeTestTag(t, "Reimbursable")

	categoryCount, err := NewCategoryRepository(nil).CountByIds([]uint{cat1, cat2})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if categoryCount != 2 {
		utils.PrintTestError(t, categoryCount, 2)
	}

	// A non-existent id should not be counted.
	categoryCount, err = NewCategoryRepository(nil).CountByIds([]uint{cat1, 999999})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if categoryCount != 1 {
		utils.PrintTestError(t, categoryCount, 1)
	}

	tagCount, err := NewTagsRepository(nil).CountByIds([]uint{tag1})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if tagCount != 1 {
		utils.PrintTestError(t, tagCount, 1)
	}
}

// makeTestUser inserts a user and returns its id.
func makeTestUser(t *testing.T, username string) uint {
	user := models.User{Username: username, Password: "x"}
	if err := GetDB().Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	return user.ID
}

func TestCreateGroupRolePersistsPaidByGrants(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	payer1 := makeTestUser(t, "payer1")
	payer2 := makeTestUser(t, "payer2")

	role, err := repository.CreateGroupRole(
		"Paid-By Role",
		"",
		[]string{permissions.GroupReceiptsRead},
		nil,
		nil,
		[]uint{payer1, payer2},
		true,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(role.PaidByUserGrants) != 2 {
		utils.PrintTestError(t, len(role.PaidByUserGrants), 2)
	}
	if !role.IncludeOwnPaidReceipts {
		utils.PrintTestError(t, role.IncludeOwnPaidReceipts, true)
	}

	userIds, err := repository.GetGroupRolePaidByUserIds(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	slices.Sort(userIds)
	expected := []uint{payer1, payer2}
	slices.Sort(expected)
	if !slices.Equal(userIds, expected) {
		utils.PrintTestError(t, userIds, expected)
	}

	includeOwn, err := repository.GetGroupRoleIncludeOwnPaidReceipts(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !includeOwn {
		utils.PrintTestError(t, includeOwn, true)
	}
}

func TestUpdateGroupRolePersistsIncludeOwnToggleOff(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	payer := makeTestUser(t, "toggle-payer")

	created, err := repository.CreateGroupRole("Toggle Role", "", []string{permissions.GroupReceiptsRead}, nil, nil, []uint{payer}, true)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Toggle include-own OFF and drop the user grant. The map-based update must
	// persist the false value (a struct update would skip the zero-value bool).
	updated, err := repository.UpdateGroupRole(created.ID, "Toggle Role", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if updated.IncludeOwnPaidReceipts {
		utils.PrintTestError(t, updated.IncludeOwnPaidReceipts, false)
	}
	if len(updated.PaidByUserGrants) != 0 {
		utils.PrintTestError(t, len(updated.PaidByUserGrants), 0)
	}

	includeOwn, err := repository.GetGroupRoleIncludeOwnPaidReceipts(created.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if includeOwn {
		utils.PrintTestError(t, includeOwn, false)
	}
}

func TestDeleteGroupRoleCascadesPaidByGrants(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	payer := makeTestUser(t, "cascade-payer")

	created, err := repository.CreateGroupRole("Cascade Role", "", []string{permissions.GroupReceiptsRead}, nil, nil, []uint{payer}, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := repository.DeleteGroupRole(created.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var paidByGrantCount int64
	GetDB().Model(&models.GroupRolePaidByUserGrant{}).Where("group_role_id = ?", created.ID).Count(&paidByGrantCount)
	if paidByGrantCount != 0 {
		utils.PrintTestError(t, paidByGrantCount, 0)
	}
}

func TestGroupRolePaidByVisibilityRestrictedPersists(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)
	payer := makeTestUser(t, "restricted-flag-payer")

	assertRestricted := func(roleId uint, want bool, desc string) {
		got, err := repository.GetGroupRolePaidByVisibilityRestricted(roleId)
		if err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
		if got != want {
			utils.PrintTestError(t, got, desc)
		}
	}

	// Specific-user grant (no include-own) is restricted.
	usersOnly, err := repository.CreateGroupRole("Users Only", "", []string{permissions.GroupReceiptsRead}, nil, nil, []uint{payer}, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	assertRestricted(usersOnly.ID, true, "users-only role should be restricted")

	// Include-own only is restricted.
	ownOnly, err := repository.CreateGroupRole("Own Only", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, true)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	assertRestricted(ownOnly.ID, true, "include-own role should be restricted")

	// Empty paid-by config is NOT restricted (unrestricted = see everyone).
	open, err := repository.CreateGroupRole("Open", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	assertRestricted(open.ID, false, "empty paid-by config should not be restricted")

	// Updating a restricted role to an empty config clears the flag (map update
	// persists the false value).
	_, err = repository.UpdateGroupRole(usersOnly.ID, "Users Only", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	assertRestricted(usersOnly.ID, false, "clearing the paid-by config should clear restricted")
}

func TestUserCountByIds(t *testing.T) {
	defer TruncateTestDb()

	user1 := makeTestUser(t, "count-user1")
	user2 := makeTestUser(t, "count-user2")

	count, err := NewUserRepository(nil).CountByIds([]uint{user1, user2})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 2 {
		utils.PrintTestError(t, count, 2)
	}

	count, err = NewUserRepository(nil).CountByIds([]uint{user1, 999999})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 1 {
		utils.PrintTestError(t, count, 1)
	}
}
