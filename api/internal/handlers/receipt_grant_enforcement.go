package handlers

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
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

// enforceReceiptMemberVisibilitySelection checks that every user id a receipt upsert
// plants — the payer (paidByUserId) and each item/linked-item chargedToUserId — is
// within the creator's member-visible set FOR THE RECEIPT'S GROUP, so an isolated member
// cannot attach a user they may not see in that group (which would leak that user's
// presence, or hand them a receipt/charge). It mirrors enforceReceiptGrantSelection:
// returns (allowed, denyMessage, error), and the caller responds 403 with denyMessage
// when not allowed. No-op when the creator is unrestricted in that group.
func enforceReceiptMemberVisibilitySelection(userId uint, groupId uint, command commands.UpsertReceiptCommand) (bool, string, error) {
	permissionService := services.NewPermissionService(nil)

	chargedToUserIds := collectReceiptChargedToUserIds(command)

	ok, err := permissionService.ValidateReceiptUserSelection(userId, groupId, command.PaidByUserID, chargedToUserIds)
	if err != nil {
		return false, "", err
	}
	if !ok {
		return false, "You do not have access to one or more selected users", nil
	}

	return true, "", nil
}

// enforceReceiptChargedToPreservation blocks an isolated editor from updating a
// receipt whose STORED items are charged to a user they cannot see. Because the
// update replaces items wholesale and item commands carry no stable id, the hidden
// charge (masked to unset when the editor read the receipt) would be silently
// dropped on save — corrupting the split. Rather than lose the charge we reject the
// edit and tell the caller to route it through someone with full access. No-op when
// the editor is unrestricted. (Letting the edit proceed while preserving the hidden
// charge needs stable item identity — a separate, larger change.) Returns
// (allowed, denyMessage, error).
func enforceReceiptChargedToPreservation(userId uint, current models.Receipt) (bool, string, error) {
	permissionService := services.NewPermissionService(nil)

	visible, unrestricted, err := permissionService.GetVisibleUserIdsForUserInGroup(userId, current.GroupId)
	if err != nil {
		return false, "", err
	}
	if unrestricted {
		return true, "", nil
	}

	hidden := func(chargedTo *uint) bool {
		if chargedTo == nil || *chargedTo == 0 || *chargedTo == userId {
			return false
		}
		_, ok := visible[*chargedTo]
		return !ok
	}

	for _, item := range current.ReceiptItems {
		if hidden(item.ChargedToUserId) {
			return false, receiptHiddenChargeMessage, nil
		}
		for _, linkedItem := range item.LinkedItems {
			if hidden(linkedItem.ChargedToUserId) {
				return false, receiptHiddenChargeMessage, nil
			}
		}
	}

	return true, "", nil
}

const receiptHiddenChargeMessage = "This receipt includes a charge to a member you do not have access to and cannot be edited here. Ask an administrator or a member with full access to edit it."

// collectReceiptChargedToUserIds walks a receipt command's items and linked items and
// returns every charged-to user id referenced (skipping unset/zero ids).
func collectReceiptChargedToUserIds(command commands.UpsertReceiptCommand) []uint {
	var ids []uint
	for _, item := range command.Items {
		if item.ChargedToUserId != nil && *item.ChargedToUserId != 0 {
			ids = append(ids, *item.ChargedToUserId)
		}
		for _, linkedItem := range item.LinkedItems {
			if linkedItem.ChargedToUserId != nil && *linkedItem.ChargedToUserId != 0 {
				ids = append(ids, *linkedItem.ChargedToUserId)
			}
		}
	}
	return ids
}

// enforceQuickScanGrantSelection checks a quick-scan command's per-file category/tag
// picks against the caller's grants for each file's group. Quick scan creates receipts
// through the service layer (bypassing enforceReceiptGrantSelection on the receipt-upsert
// path), so this is the synchronous grant gate — the caller responds 403 with denyMessage
// when a pick is outside the caller's grants. Quick-scan pickers are id-only (no new-by-name
// category/tag), so no create-permission check is needed. It returns (allowed, denyMessage, error).
func enforceQuickScanGrantSelection(userId uint, command commands.QuickScanCommand) (bool, string, error) {
	permissionService := services.NewPermissionService(nil)

	for i := 0; i < len(command.GroupIds); i++ {
		categoryIds := command.CategoryIdsForFile(i)
		tagIds := command.TagIdsForFile(i)

		ok, err := permissionService.ValidateCategoryTagSelection(userId, command.GroupIds[i], categoryIds, tagIds)
		if err != nil {
			return false, "", err
		}
		if !ok {
			return false, "You do not have access to one or more selected categories or tags", nil
		}
	}

	return true, "", nil
}

// enforceReceiptCustomFieldSelection ensures a caller WITHOUT app.custom-fields.read
// cannot change the SET of custom fields attached to a receipt — they may edit the
// VALUES of fields already present, but may not add a new custom field id or remove an
// existing one. Read holders manage the set freely. currentCustomFieldIds is the
// receipt's existing set (nil/empty on create, so any custom field is an add). Returns
// (allowed, denyMessage, error); the caller responds 403 with denyMessage when not allowed.
func enforceReceiptCustomFieldSelection(userId uint, command commands.UpsertReceiptCommand, currentCustomFieldIds []uint) (bool, string, error) {
	permissionService := services.NewPermissionService(nil)

	canManage, err := permissionService.HasAppPermissions(userId, permissions.AppCustomFieldsRead)
	if err != nil {
		return false, "", err
	}
	if canManage {
		return true, "", nil
	}

	submitted := make(map[uint]struct{}, len(command.CustomFields))
	for _, customField := range command.CustomFields {
		submitted[customField.CustomFieldId] = struct{}{}
	}

	current := make(map[uint]struct{}, len(currentCustomFieldIds))
	for _, id := range currentCustomFieldIds {
		current[id] = struct{}{}
	}

	if !uintSetsEqual(submitted, current) {
		return false, "You do not have permission to add or remove custom fields", nil
	}

	return true, "", nil
}

// uintSetsEqual reports whether two uint sets contain exactly the same ids.
func uintSetsEqual(a map[uint]struct{}, b map[uint]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for id := range a {
		if _, ok := b[id]; !ok {
			return false
		}
	}
	return true
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
