package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

// categoryWithId / tagWithId build an in-memory association row referencing an
// existing category/tag id (the grant filter only inspects ids).
func categoryWithId(id uint) models.Category {
	return models.Category{BaseModel: models.BaseModel{ID: id}}
}

func tagWithId(id uint) models.Tag {
	return models.Tag{BaseModel: models.BaseModel{ID: id}}
}

func categoryIdSet(categories []models.Category) map[uint]struct{} {
	set := make(map[uint]struct{}, len(categories))
	for _, category := range categories {
		set[category.ID] = struct{}{}
	}
	return set
}

// grantUserAppPermissions assigns userId an app role granting perms, so the
// user bypasses grant filtering for those resources.
func grantUserAppPermissions(t *testing.T, userId uint, perms []string) {
	t.Helper()
	role, err := repositories.NewRoleRepository(nil).CreateAppRole("Bypass Role", "", perms)
	if err != nil {
		t.Fatalf("create app role: %v", err)
	}
	if err := repositories.GetDB().Model(&models.User{}).Where("id = ?", userId).Update("app_role_id", role.ID).Error; err != nil {
		t.Fatalf("assign app role: %v", err)
	}
	clearRolePermissionCacheAll()
}

func TestFilterReceiptStripsDisallowedCategories(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	tag1 := makeTag(t, "Reimbursable")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-strip", []uint{cat1}, nil)

	receipts := []models.Receipt{{
		GroupId:    groupId,
		Categories: []models.Category{categoryWithId(cat1), categoryWithId(cat2)},
		Tags:       []models.Tag{tagWithId(tag1)},
	}}

	service := NewPermissionService(nil)
	if err := service.FilterReceiptCategoriesTags(userId, receipts); err != nil {
		t.Fatalf("FilterReceiptCategoriesTags: %v", err)
	}

	if len(receipts[0].Categories) != 1 || receipts[0].Categories[0].ID != cat1 {
		t.Errorf("expected only cat1 to survive, got %v", receipts[0].Categories)
	}
	// The role grants no tags, so tags are unrestricted and left untouched.
	if len(receipts[0].Tags) != 1 || receipts[0].Tags[0].ID != tag1 {
		t.Errorf("expected tags untouched (unrestricted), got %v", receipts[0].Tags)
	}
}

func TestFilterReceiptUnrestrictedNoOp(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-open", nil, nil)

	receipts := []models.Receipt{{
		GroupId:    groupId,
		Categories: []models.Category{categoryWithId(cat1), categoryWithId(cat2)},
	}}

	service := NewPermissionService(nil)
	if err := service.FilterReceiptCategoriesTags(userId, receipts); err != nil {
		t.Fatalf("FilterReceiptCategoriesTags: %v", err)
	}

	if len(receipts[0].Categories) != 2 {
		t.Errorf("expected unrestricted receipt untouched, got %v", receipts[0].Categories)
	}
}

func TestFilterReceiptAdminBypass(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	// Restricted group role (only cat1) ...
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-admin", []uint{cat1}, nil)
	// ... but the user holds app.categories.read, so grants do not apply.
	grantUserAppPermissions(t, userId, []string{permissions.AppCategoriesRead})

	receipts := []models.Receipt{{
		GroupId:    groupId,
		Categories: []models.Category{categoryWithId(cat1), categoryWithId(cat2)},
	}}

	service := NewPermissionService(nil)
	if err := service.FilterReceiptCategoriesTags(userId, receipts); err != nil {
		t.Fatalf("FilterReceiptCategoriesTags: %v", err)
	}

	if len(receipts[0].Categories) != 2 {
		t.Errorf("expected admin (app.categories.read) to bypass stripping, got %v", receipts[0].Categories)
	}
}

func TestFilterReceiptBatchPerGroup(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	catA := makeCategory(t, "A")
	catB := makeCategory(t, "B")

	// Same user restricted differently in two groups.
	userId, groupA, _ := seedMemberWithGroupRoleGrants(t, "u-batch", []uint{catA}, nil)

	db := repositories.GetDB()
	groupB := models.Group{Name: "group-b"}
	db.Create(&groupB)
	roleB, err := repositories.NewRoleRepository(nil).CreateGroupRole("Role B", "", []string{permissions.GroupReceiptsRead}, []uint{catB}, nil)
	if err != nil {
		t.Fatalf("create role B: %v", err)
	}
	memberB := models.GroupMember{GroupID: groupB.ID, UserID: userId, GroupRoleID: &roleB.ID}
	db.Create(&memberB)

	receipts := []models.Receipt{
		{GroupId: groupA, Categories: []models.Category{categoryWithId(catA), categoryWithId(catB)}},
		{GroupId: groupB.ID, Categories: []models.Category{categoryWithId(catA), categoryWithId(catB)}},
	}

	service := NewPermissionService(nil)
	if err := service.FilterReceiptCategoriesTags(userId, receipts); err != nil {
		t.Fatalf("FilterReceiptCategoriesTags: %v", err)
	}

	if _, ok := categoryIdSet(receipts[0].Categories)[catA]; !ok || len(receipts[0].Categories) != 1 {
		t.Errorf("group A receipt should keep only catA, got %v", receipts[0].Categories)
	}
	if _, ok := categoryIdSet(receipts[1].Categories)[catB]; !ok || len(receipts[1].Categories) != 1 {
		t.Errorf("group B receipt should keep only catB, got %v", receipts[1].Categories)
	}
}

func TestFilterReceiptForReceiptSingle(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-single", []uint{cat1}, nil)

	receipt := models.Receipt{
		GroupId:    groupId,
		Categories: []models.Category{categoryWithId(cat1), categoryWithId(cat2)},
	}

	service := NewPermissionService(nil)
	if err := service.FilterReceiptCategoriesTagsForReceipt(userId, &receipt); err != nil {
		t.Fatalf("FilterReceiptCategoriesTagsForReceipt: %v", err)
	}

	if len(receipt.Categories) != 1 || receipt.Categories[0].ID != cat1 {
		t.Errorf("expected only cat1 to survive, got %v", receipt.Categories)
	}
}

func TestValidateCategoryTagSelection(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	cat1 := makeCategory(t, "Groceries")
	cat2 := makeCategory(t, "Utilities")
	tag1 := makeTag(t, "Reimbursable")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-validate", []uint{cat1}, nil)

	service := NewPermissionService(nil)

	// In-set category passes.
	ok, err := service.ValidateCategoryTagSelection(userId, groupId, []uint{cat1}, nil)
	if err != nil {
		t.Fatalf("validate in-set: %v", err)
	}
	if !ok {
		t.Error("expected in-set category selection to be allowed")
	}

	// Out-of-set category fails.
	ok, err = service.ValidateCategoryTagSelection(userId, groupId, []uint{cat2}, nil)
	if err != nil {
		t.Fatalf("validate out-of-set: %v", err)
	}
	if ok {
		t.Error("expected out-of-set category selection to be rejected")
	}

	// Tags are unrestricted for this role, so any tag passes.
	ok, err = service.ValidateCategoryTagSelection(userId, groupId, []uint{cat1}, []uint{tag1})
	if err != nil {
		t.Fatalf("validate tag: %v", err)
	}
	if !ok {
		t.Error("expected tag selection allowed when tags unrestricted")
	}
}

func TestValidateCategoryTagSelectionUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-validate-open", nil, nil)

	service := NewPermissionService(nil)
	// Unrestricted role: any id passes (here an id that does not even exist).
	ok, err := service.ValidateCategoryTagSelection(userId, groupId, []uint{424242}, nil)
	if err != nil {
		t.Fatalf("validate unrestricted: %v", err)
	}
	if !ok {
		t.Error("expected unrestricted role to allow any category selection")
	}
}
