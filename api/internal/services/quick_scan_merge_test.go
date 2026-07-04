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

func TestMergeQuickScanCategories(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestCategories()
	service := NewReceiptService(nil)

	id1 := uint(1)

	// No selections leaves the AI-filled categories untouched.
	existing := []commands.UpsertCategoryCommand{{Id: &id1, Name: "test"}}
	result, err := service.mergeQuickScanCategories(existing, []uint{})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 1 {
		utils.PrintTestError(t, result, "unchanged existing")
	}

	// Union of AI (id 1) and user picks (ids 2, 3), all present exactly once, names resolved.
	result, err = service.mergeQuickScanCategories(existing, []uint{2, 3})
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
	result, err = service.mergeQuickScanCategories(existing, []uint{1, 2})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 2 || !containsCategoryId(result, 1) || !containsCategoryId(result, 2) {
		utils.PrintTestError(t, result, "categories 1,2 deduped")
	}
}

func TestMergeQuickScanTags(t *testing.T) {
	defer repositories.TruncateTestDb()
	createTestTags()
	service := NewReceiptService(nil)

	id1 := uint(1)
	existing := []commands.UpsertTagCommand{{Id: &id1, Name: "tag-a"}}

	result, err := service.mergeQuickScanTags(existing, []uint{2})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 2 || !containsTagId(result, 1) || !containsTagId(result, 2) {
		utils.PrintTestError(t, result, "tags 1,2 unioned")
	}

	// Dedup overlap.
	result, err = service.mergeQuickScanTags(existing, []uint{1})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(result) != 1 || !containsTagId(result, 1) {
		utils.PrintTestError(t, result, "tag 1 deduped")
	}
}
