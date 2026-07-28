package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func createTestTags() {
	db := repositories.GetDB()
	db.Create(&models.Tag{Name: "tag-a"})
	db.Create(&models.Tag{Name: "tag-b"})
}

func containsCategoryId(categories []commands.UpsertCategoryCommand, id uint) bool {
	count := 0
	for _, category := range categories {
		if category.Id != nil && *category.Id == id {
			count++
		}
	}
	return count == 1
}

func containsTagId(tags []commands.UpsertTagCommand, id uint) bool {
	count := 0
	for _, tag := range tags {
		if tag.Id != nil && *tag.Id == id {
			count++
		}
	}
	return count == 1
}

func TestResolveQuickScanCategories(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestCategories()
	service := NewReceiptService(nil)

	id1 := uint(1)

	// AI returns category 1 by id only (no name); it is resolved to its real record so the name is
	// filled in, even with no user picks. userId/groupId 0 => unrestricted (no grants apply).
	aiCategories := []commands.UpsertCategoryCommand{{Id: &id1}}
	result, err := service.resolveQuickScanCategories(aiCategories, nil, 0, 0)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 1 || !containsCategoryId(result, 1) || result[0].Name != "test" {
		utils.PrintTestError(t, result, "category 1 resolved with name")
	}

	// Union of AI (id 1) and user picks (ids 2, 3), all present exactly once, names resolved.
	result, err = service.resolveQuickScanCategories(aiCategories, []uint{2, 3}, 0, 0)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 3 || !containsCategoryId(result, 1) || !containsCategoryId(result, 2) || !containsCategoryId(result, 3) {
		utils.PrintTestError(t, result, "categories 1,2,3 exactly once")
	}
	for _, category := range result {
		if len(category.Name) == 0 {
			utils.PrintTestError(t, category, "name resolved for merged category")
		}
	}

	// Overlap between AI and user pick (both id 1) is deduped, not duplicated.
	result, err = service.resolveQuickScanCategories(aiCategories, []uint{1, 2}, 0, 0)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 2 || !containsCategoryId(result, 1) || !containsCategoryId(result, 2) {
		utils.PrintTestError(t, result, "categories 1,2 deduped")
	}

	// A hallucinated / non-existent id is dropped rather than surfaced.
	badId := uint(999)
	result, err = service.resolveQuickScanCategories([]commands.UpsertCategoryCommand{{Id: &id1}, {Id: &badId}}, nil, 0, 0)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 1 || !containsCategoryId(result, 1) {
		utils.PrintTestError(t, result, "only category 1 (999 dropped)")
	}
}

func TestResolveQuickScanTags(t *testing.T) {
	defer repositories.TruncateTestDb()
	createTestTags()
	service := NewReceiptService(nil)

	id1 := uint(1)
	aiTags := []commands.UpsertTagCommand{{Id: &id1}}

	// AI tag by id only is resolved (name filled) and unioned with the user's pick.
	result, err := service.resolveQuickScanTags(aiTags, []uint{2}, 0, 0)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 2 || !containsTagId(result, 1) || !containsTagId(result, 2) {
		utils.PrintTestError(t, result, "tags 1,2 unioned")
	}
	for _, tag := range result {
		if len(tag.Name) == 0 {
			utils.PrintTestError(t, tag, "name resolved for merged tag")
		}
	}

	// Dedup overlap.
	result, err = service.resolveQuickScanTags(aiTags, []uint{1}, 0, 0)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 1 || !containsTagId(result, 1) {
		utils.PrintTestError(t, result, "tag 1 deduped")
	}
}
