package services

import (
	"receipt-wrangler/api/internal/commands"
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

// SubstituteRestrictedCategoriesTags rewrites every receipt's Categories/Tags in
// place, replacing the entries the user may not see with a single (Restricted)
// marker rather than dropping them (as FilterReceiptCategoriesTags does).
// Aggregation surfaces — reports, charts — use this so a hidden category is
// attributed to its own (Restricted) bucket instead of vanishing from the totals
// or collapsing into (None). Like the strip variant, the allowed sets resolve
// once per group, and a receipt whose group is unrestricted (or whose viewer
// bypasses grants) is left untouched.
//
// It operates on receipt-level Categories/Tags only, which is all an aggregation
// reads; it does not recurse into receipt items.
func (service PermissionService) SubstituteRestrictedCategoriesTags(userId uint, receipts []models.Receipt) error {
	if len(receipts) == 0 {
		return nil
	}

	filter, err := service.newReceiptGrantFilter(userId)
	if err != nil {
		return err
	}

	for i := range receipts {
		if err := filter.substitute(&receipts[i]); err != nil {
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

// MergeHiddenReceiptCategoriesTags appends to an update command the
// receipt-level categories/tags that exist on the current receipt but are
// outside the user's grants (and so were hidden from them on read). Without this
// a restricted user's full-replace update would silently delete categories/tags
// they could not see. No-op when the user is unrestricted or bypasses grants.
//
// NOTE: this only merges RECEIPT-level associations. Receipt items have no stable
// id across an update (UpdateReceipt recreates them), so hidden item-level
// categories/tags cannot be matched back and are not preserved.
func (service PermissionService) MergeHiddenReceiptCategoriesTags(
	userId uint,
	groupId uint,
	currentCategories []models.Category,
	currentTags []models.Tag,
	command *commands.UpsertReceiptCommand,
) error {
	bypassCategories, err := service.userBypassesGrants(userId, permissions.AppCategoriesRead)
	if err != nil {
		return err
	}
	if !bypassCategories {
		allowed, unrestricted, err := service.GetGroupCategoryIdsForUser(userId, groupId)
		if err != nil {
			return err
		}
		if !unrestricted {
			for i := range currentCategories {
				category := currentCategories[i]
				if _, ok := allowed[category.ID]; !ok {
					categoryId := category.ID
					command.Categories = append(command.Categories, commands.UpsertCategoryCommand{
						Id:          &categoryId,
						Name:        category.Name,
						Description: category.Description,
					})
				}
			}
		}
	}

	bypassTags, err := service.userBypassesGrants(userId, permissions.AppTagsRead)
	if err != nil {
		return err
	}
	if !bypassTags {
		allowed, unrestricted, err := service.GetGroupTagIdsForUser(userId, groupId)
		if err != nil {
			return err
		}
		if !unrestricted {
			for i := range currentTags {
				tag := currentTags[i]
				if _, ok := allowed[tag.ID]; !ok {
					tagId := tag.ID
					command.Tags = append(command.Tags, commands.UpsertTagCommand{
						Id:          &tagId,
						Name:        tag.Name,
						Description: tag.Description,
					})
				}
			}
		}
	}

	return nil
}

// IntersectReceiptFilterWithGrants narrows a paged-request category/tag filter to
// the ids the user may see, so a restricted user cannot probe receipt existence
// by filtering on a category/tag they have no access to. When a filter is made
// up entirely of disallowed ids it is forced to a no-match sentinel rather than
// dropped (dropping would widen the result set). No-op for unrestricted/bypass.
func (service PermissionService) IntersectReceiptFilterWithGrants(userId uint, groupId uint, filter *commands.ReceiptPagedRequestFilter) error {
	bypassCategories, err := service.userBypassesGrants(userId, permissions.AppCategoriesRead)
	if err != nil {
		return err
	}
	if !bypassCategories {
		allowed, unrestricted, err := service.GetGroupCategoryIdsForUser(userId, groupId)
		if err != nil {
			return err
		}
		if !unrestricted {
			intersectFilterFieldWithGrants(&filter.Categories, allowed)
		}
	}

	bypassTags, err := service.userBypassesGrants(userId, permissions.AppTagsRead)
	if err != nil {
		return err
	}
	if !bypassTags {
		allowed, unrestricted, err := service.GetGroupTagIdsForUser(userId, groupId)
		if err != nil {
			return err
		}
		if !unrestricted {
			intersectFilterFieldWithGrants(&filter.Tags, allowed)
		}
	}

	return nil
}

func intersectFilterFieldWithGrants(field *commands.PagedRequestField, allowed map[uint]struct{}) {
	if field.Value == nil {
		return
	}
	values, ok := field.Value.([]interface{})
	if !ok || len(values) == 0 {
		return
	}

	filtered := make([]interface{}, 0, len(values))
	for _, value := range values {
		id, ok := filterValueToUint(value)
		if !ok {
			continue
		}
		if _, allowedOk := allowed[id]; allowedOk {
			filtered = append(filtered, value)
		}
	}

	if len(filtered) == 0 {
		// The user filtered entirely by ids they cannot see; force a no-match so
		// receipt existence is not revealed. Category/tag ids start at 1, so 0
		// matches nothing.
		filtered = []interface{}{float64(0)}
	}

	field.Value = filtered
}

// filterValueToUint coerces a JSON-decoded filter id (typically float64) to uint.
func filterValueToUint(value interface{}) (uint, bool) {
	switch typed := value.(type) {
	case float64:
		return uint(typed), true
	case int:
		return uint(typed), true
	case uint:
		return typed, true
	default:
		return 0, false
	}
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

	filter.stripCategoriesTags(&receipt.Categories, &receipt.Tags, grants)

	// Receipt items (and their linked items) carry their own categories/tags —
	// strip them with the same rule so item-level associations don't leak.
	for i := range receipt.ReceiptItems {
		item := &receipt.ReceiptItems[i]
		filter.stripCategoriesTags(&item.Categories, &item.Tags, grants)
		for j := range item.LinkedItems {
			linkedItem := &item.LinkedItems[j]
			filter.stripCategoriesTags(&linkedItem.Categories, &linkedItem.Tags, grants)
		}
	}
	return nil
}

// stripCategoriesTags rewrites a categories/tags pair in place to the allowed
// subset, honoring the per-resource bypass and unrestricted flags.
func (filter *receiptGrantFilter) stripCategoriesTags(categories *[]models.Category, tags *[]models.Tag, grants receiptGroupGrants) {
	if !filter.bypassCategories && !grants.categoryUnrestricted {
		*categories = filterCategoriesBySet(*categories, grants.categoryAllowed)
	}
	if !filter.bypassTags && !grants.tagUnrestricted {
		*tags = filterTagsBySet(*tags, grants.tagAllowed)
	}
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

// restrictedCategoryTagName is the marker a hidden category or tag is replaced
// with on an aggregation surface. It forms its own bucket, kept distinct from
// (None) — a receipt that genuinely carries no category or tag.
const restrictedCategoryTagName = "(Restricted)"

// substitute is the substitution counterpart of apply: it keeps every receipt
// but replaces the categories/tags the user may not see with a (Restricted)
// marker. It stays at receipt level because that is all an aggregation reads.
func (filter *receiptGrantFilter) substitute(receipt *models.Receipt) error {
	if filter.bypassCategories && filter.bypassTags {
		return nil
	}

	grants, err := filter.grantsForGroup(receipt.GroupId)
	if err != nil {
		return err
	}

	if !filter.bypassCategories && !grants.categoryUnrestricted {
		receipt.Categories = substituteCategoriesBySet(receipt.Categories, grants.categoryAllowed)
	}
	if !filter.bypassTags && !grants.tagUnrestricted {
		receipt.Tags = substituteTagsBySet(receipt.Tags, grants.tagAllowed)
	}
	return nil
}

// substituteCategoriesBySet keeps the allowed categories and, if any category was
// disallowed, appends exactly one (Restricted) marker. Several hidden categories
// collapse to that one marker, which is enough: the engine attributes a receipt
// to a multi-valued bucket once no matter how many times it carries the value. A
// receipt with no disallowed category gets no marker.
func substituteCategoriesBySet(categories []models.Category, allowed map[uint]struct{}) []models.Category {
	result := make([]models.Category, 0, len(categories))
	hidden := false
	for _, category := range categories {
		if _, ok := allowed[category.ID]; ok {
			result = append(result, category)
			continue
		}
		hidden = true
	}
	if hidden {
		result = append(result, models.Category{Name: restrictedCategoryTagName})
	}
	return result
}

func substituteTagsBySet(tags []models.Tag, allowed map[uint]struct{}) []models.Tag {
	result := make([]models.Tag, 0, len(tags))
	hidden := false
	for _, tag := range tags {
		if _, ok := allowed[tag.ID]; ok {
			result = append(result, tag)
			continue
		}
		hidden = true
	}
	if hidden {
		result = append(result, models.Tag{Name: restrictedCategoryTagName})
	}
	return result
}
