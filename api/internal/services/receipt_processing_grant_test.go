package services

import (
	"encoding/json"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

func TestGetCategoriesStringRestrictedToUserGrants(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	grantedCategory := makeCategory(t, "Groceries")
	makeCategory(t, "Salary")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-ai-cat", []uint{grantedCategory}, nil)

	service := ReceiptProcessingService{
		Group:  models.Group{BaseModel: models.BaseModel{ID: groupId}},
		UserId: userId,
	}

	result, err := service.getCategoriesString()
	if err != nil {
		t.Fatalf("getCategoriesString: %v", err)
	}

	var categories []models.Category
	if err := json.Unmarshal([]byte(result), &categories); err != nil {
		t.Fatalf("unmarshal categories: %v", err)
	}
	if len(categories) != 1 || categories[0].ID != grantedCategory {
		t.Errorf("expected the AI prompt to see only the granted category, got %v", categories)
	}
}

func TestGetCategoriesStringUnrestrictedForSystemProcessing(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	makeCategory(t, "Groceries")
	makeCategory(t, "Salary")

	// UserId 0 = system-initiated (e.g. email polling): no restriction.
	service := ReceiptProcessingService{}

	result, err := service.getCategoriesString()
	if err != nil {
		t.Fatalf("getCategoriesString: %v", err)
	}

	var categories []models.Category
	if err := json.Unmarshal([]byte(result), &categories); err != nil {
		t.Fatalf("unmarshal categories: %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("expected all categories for system processing, got %d", len(categories))
	}
}

func TestGetTagsStringRestrictedToUserGrants(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	grantedTag := makeTag(t, "Reimbursable")
	makeTag(t, "Personal")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "u-ai-tag", nil, []uint{grantedTag})

	service := ReceiptProcessingService{
		Group:  models.Group{BaseModel: models.BaseModel{ID: groupId}},
		UserId: userId,
	}

	result, err := service.getTagsString()
	if err != nil {
		t.Fatalf("getTagsString: %v", err)
	}

	var tags []models.Tag
	if err := json.Unmarshal([]byte(result), &tags); err != nil {
		t.Fatalf("unmarshal tags: %v", err)
	}
	if len(tags) != 1 || tags[0].ID != grantedTag {
		t.Errorf("expected the AI prompt to see only the granted tag, got %v", tags)
	}
}
