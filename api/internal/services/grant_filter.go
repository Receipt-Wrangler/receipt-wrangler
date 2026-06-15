package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
)

// This file holds the shared category/tag grant enforcement used by every read
// surface (list/search, single receipt, widgets, export, duplicate) and the
// write surface (receipt create/update), so the rules live in exactly one place.

// userBypassesGrants reports whether a user holds the app-level read permission
// for a resource (app.categories.read / app.tags.read). Such users (admins,
// category/tag managers) can already see the entire global pool, so grant
// restrictions are not applied to them — keeping their view consistent with the
// global lists they are entitled to.
func (service PermissionService) userBypassesGrants(userId uint, appReadPermission string) (bool, error) {
	return service.HasAppPermissions(userId, appReadPermission)
}

// FilterReceiptCategoriesTags strips every receipt's Categories/Tags down to the
// set the user may see, in place. The allowed sets are resolved once per group
// (cache-backed), so a page of receipts spanning N groups does at most N
// resolutions. Receipts whose group is unrestricted (or where the user bypasses
// grants) are left untouched.
func (service PermissionService) FilterReceiptCategoriesTags(userId uint, receipts []models.Receipt) error {
	if len(receipts) == 0 {
		return nil
	}

	filter, err := service.newReceiptGrantFilter(userId)
	if err != nil {
		return err
	}

	for i := range receipts {
		if err := filter.apply(&receipts[i]); err != nil {
			return err
		}
	}
	return nil
}

// FilterReceiptCategoriesTagsForReceipt is the single-receipt variant of
// FilterReceiptCategoriesTags (single GetReceipt, duplicate source, etc.).
func (service PermissionService) FilterReceiptCategoriesTagsForReceipt(userId uint, receipt *models.Receipt) error {
	if receipt == nil {
		return nil
	}

	filter, err := service.newReceiptGrantFilter(userId)
	if err != nil {
		return err
	}

	return filter.apply(receipt)
}

// ValidateCategoryTagSelection reports whether every category/tag id is within
// the user's allowed set for the group. Unrestricted resources, and users who
// bypass grants, always pass. Used on receipt create/update for ids that
// reference EXISTING categories/tags; new-by-name selections are not id-checked
// here — they are gated by the create permission at the call site.
func (service PermissionService) ValidateCategoryTagSelection(userId uint, groupId uint, categoryIds []uint, tagIds []uint) (bool, error) {
	if len(categoryIds) > 0 {
		ok, err := service.selectionWithinGrants(userId, groupId, permissions.AppCategoriesRead, categoryIds, service.GetGroupCategoryIdsForUser)
		if err != nil || !ok {
			return false, err
		}
	}

	if len(tagIds) > 0 {
		ok, err := service.selectionWithinGrants(userId, groupId, permissions.AppTagsRead, tagIds, service.GetGroupTagIdsForUser)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

// selectionWithinGrants checks that every id is allowed for the user in the
// group, short-circuiting when the user bypasses grants or the resource is
// unrestricted. resolve is the per-resource id resolver
// (GetGroupCategoryIdsForUser / GetGroupTagIdsForUser).
func (service PermissionService) selectionWithinGrants(
	userId uint,
	groupId uint,
	appReadPermission string,
	ids []uint,
	resolve func(uint, uint) (map[uint]struct{}, bool, error),
) (bool, error) {
	bypass, err := service.userBypassesGrants(userId, appReadPermission)
	if err != nil {
		return false, err
	}
	if bypass {
		return true, nil
	}

	allowed, unrestricted, err := resolve(userId, groupId)
	if err != nil {
		return false, err
	}
	if unrestricted {
		return true, nil
	}

	for _, id := range ids {
		if _, ok := allowed[id]; !ok {
			return false, nil
		}
	}
	return true, nil
}

// receiptGrantFilter caches a user's bypass flags and per-group allowed sets for
// the lifetime of one filter pass, so a batch of receipts resolves each group's
// grants at most once.
type receiptGrantFilter struct {
	service          PermissionService
	userId           uint
	bypassCategories bool
	bypassTags       bool
	groupCache       map[uint]receiptGroupGrants
}

type receiptGroupGrants struct {
	categoryAllowed      map[uint]struct{}
	categoryUnrestricted bool
	tagAllowed           map[uint]struct{}
	tagUnrestricted      bool
}

func (service PermissionService) newReceiptGrantFilter(userId uint) (*receiptGrantFilter, error) {
	bypassCategories, err := service.userBypassesGrants(userId, permissions.AppCategoriesRead)
	if err != nil {
		return nil, err
	}

	bypassTags, err := service.userBypassesGrants(userId, permissions.AppTagsRead)
	if err != nil {
		return nil, err
	}

	return &receiptGrantFilter{
		service:          service,
		userId:           userId,
		bypassCategories: bypassCategories,
		bypassTags:       bypassTags,
		groupCache:       map[uint]receiptGroupGrants{},
	}, nil
}

func (filter *receiptGrantFilter) apply(receipt *models.Receipt) error {
	if filter.bypassCategories && filter.bypassTags {
		return nil
	}

	grants, err := filter.grantsForGroup(receipt.GroupId)
	if err != nil {
		return err
	}

	if !filter.bypassCategories && !grants.categoryUnrestricted {
		receipt.Categories = filterCategoriesBySet(receipt.Categories, grants.categoryAllowed)
	}
	if !filter.bypassTags && !grants.tagUnrestricted {
		receipt.Tags = filterTagsBySet(receipt.Tags, grants.tagAllowed)
	}
	return nil
}

func (filter *receiptGrantFilter) grantsForGroup(groupId uint) (receiptGroupGrants, error) {
	if cached, ok := filter.groupCache[groupId]; ok {
		return cached, nil
	}

	categoryAllowed, categoryUnrestricted, err := filter.service.GetGroupCategoryIdsForUser(filter.userId, groupId)
	if err != nil {
		return receiptGroupGrants{}, err
	}

	tagAllowed, tagUnrestricted, err := filter.service.GetGroupTagIdsForUser(filter.userId, groupId)
	if err != nil {
		return receiptGroupGrants{}, err
	}

	grants := receiptGroupGrants{
		categoryAllowed:      categoryAllowed,
		categoryUnrestricted: categoryUnrestricted,
		tagAllowed:           tagAllowed,
		tagUnrestricted:      tagUnrestricted,
	}
	filter.groupCache[groupId] = grants
	return grants, nil
}

func filterCategoriesBySet(categories []models.Category, allowed map[uint]struct{}) []models.Category {
	filtered := make([]models.Category, 0, len(categories))
	for _, category := range categories {
		if _, ok := allowed[category.ID]; ok {
			filtered = append(filtered, category)
		}
	}
	return filtered
}

func filterTagsBySet(tags []models.Tag, allowed map[uint]struct{}) []models.Tag {
	filtered := make([]models.Tag, 0, len(tags))
	for _, tag := range tags {
		if _, ok := allowed[tag.ID]; ok {
			filtered = append(filtered, tag)
		}
	}
	return filtered
}
