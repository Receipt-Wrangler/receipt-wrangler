package handlers

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
)

// enforceReceiptGrantSelection checks a receipt upsert command against the
// caller's category/tag grants for the group: every EXISTING category/tag id
// (receipt-level and item-level) must be within the allowed set, and any NEW
// (by-name) category/tag requires the matching app create permission. It returns
// (allowed, denyMessage, error); when allowed is false the caller should respond
// 403 with denyMessage.
func enforceReceiptGrantSelection(userId uint, groupId uint, command commands.UpsertReceiptCommand) (bool, string, error) {
	permissionService := services.NewPermissionService(nil)

	categoryIds, tagIds, hasNewCategory, hasNewTag := collectReceiptSelection(command)

	if hasNewCategory {
		ok, err := permissionService.HasAppPermissions(userId, permissions.AppCategoriesCreate)
		if err != nil {
			return false, "", err
		}
		if !ok {
			return false, "You do not have permission to create new categories", nil
		}
	}

	if hasNewTag {
		ok, err := permissionService.HasAppPermissions(userId, permissions.AppTagsCreate)
		if err != nil {
			return false, "", err
		}
		if !ok {
			return false, "You do not have permission to create new tags", nil
		}
	}

	ok, err := permissionService.ValidateCategoryTagSelection(userId, groupId, categoryIds, tagIds)
	if err != nil {
		return false, "", err
	}
	if !ok {
		return false, "You do not have access to one or more selected categories or tags", nil
	}

	return true, "", nil
}

// collectReceiptSelection walks a receipt command (receipt-level, item-level, and
// linked-item-level) and returns the existing category/tag ids referenced plus
// whether any new-by-name category/tag is present.
func collectReceiptSelection(command commands.UpsertReceiptCommand) (categoryIds []uint, tagIds []uint, hasNewCategory bool, hasNewTag bool) {
	addCategories := func(categories []commands.UpsertCategoryCommand) {
		for _, category := range categories {
			if category.Id != nil && *category.Id != 0 {
				categoryIds = append(categoryIds, *category.Id)
			} else {
				hasNewCategory = true
			}
		}
	}

	addTags := func(tags []commands.UpsertTagCommand) {
		for _, tag := range tags {
			if tag.Id != nil && *tag.Id != 0 {
				tagIds = append(tagIds, *tag.Id)
			} else {
				hasNewTag = true
			}
		}
	}

	addCategories(command.Categories)
	addTags(command.Tags)

	for _, item := range command.Items {
		addCategories(item.Categories)
		addTags(item.Tags)
		for _, linkedItem := range item.LinkedItems {
			addCategories(linkedItem.Categories)
			addTags(linkedItem.Tags)
		}
	}

	return categoryIds, tagIds, hasNewCategory, hasNewTag
}
