package services

import (
	"errors"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"

	"gorm.io/gorm"
)

// GetGroupCategoryIdsForUser returns the set of category ids a user may use in a
// group, plus an "unrestricted" flag. When unrestricted is true the returned set
// is nil and the caller must treat every category as allowed (the empty-grants =
// see-all rule). A non-member, or a member whose group role has no category
// grants, is unrestricted — grants only NARROW access within an already-permitted
// group; they never grant access (the handler permission gate is the access
// control). The returned map is read-only shared cache state.
func (service PermissionService) GetGroupCategoryIdsForUser(userId uint, groupId uint) (map[uint]struct{}, bool, error) {
	entry, err := service.resolveGroupRoleGrants(userId, groupId)
	if err != nil {
		return nil, false, err
	}
	if entry == nil || len(entry.categoryIds) == 0 {
		return nil, true, nil
	}
	return entry.categoryIds, false, nil
}

// GetGroupTagIdsForUser is the tag counterpart of GetGroupCategoryIdsForUser.
func (service PermissionService) GetGroupTagIdsForUser(userId uint, groupId uint) (map[uint]struct{}, bool, error) {
	entry, err := service.resolveGroupRoleGrants(userId, groupId)
	if err != nil {
		return nil, false, err
	}
	if entry == nil || len(entry.tagIds) == 0 {
		return nil, true, nil
	}
	return entry.tagIds, false, nil
}

// GetGroupPaidByUserIdsForUser returns the set of "paid by" user ids whose
// receipts a user may see in a group, plus an "unrestricted" flag. When
// unrestricted is true the returned set is nil and the caller must treat every
// payer as allowed. "Unrestricted" is keyed off whether the role opted into
// paid-by filtering at all (paidByVisibilityRestricted), NOT the current grant
// count: a role that the admin configured stays restricted even if its grant rows
// were since removed (e.g. a granted user was deleted and the FK cascade emptied
// PaidByUserGrants) — it then resolves to an EMPTY allowed set ("see nothing"),
// failing closed rather than silently widening to see-all. A non-member, or a
// member whose role never opted in, is unrestricted — grants only NARROW
// visibility within an already-permitted group (the handler permission gate is
// the access control). When restricted, the role's absolute grants are combined
// with the relative "their own receipts" token (the requester's own id is unioned
// in when includeOwnPaidReceipts). The returned map is freshly allocated per call.
func (service PermissionService) GetGroupPaidByUserIdsForUser(userId uint, groupId uint) (map[uint]struct{}, bool, error) {
	entry, err := service.resolveGroupRoleGrants(userId, groupId)
	if err != nil {
		return nil, false, err
	}
	// Unrestricted only when the role did NOT opt into paid-by filtering. The
	// persisted PaidByVisibilityRestricted flag is the primary signal (it survives
	// grant rows being cascade-deleted when a granted user is deleted), but we also
	// honor the live grants/include-own so a role whose flag somehow desynced from
	// its rows still fails closed rather than widening to see-all.
	if entry == nil || (!entry.paidByVisibilityRestricted && len(entry.paidByUserIds) == 0 && !entry.includeOwnPaidReceipts) {
		return nil, true, nil
	}

	// Copy the cached absolute grants (shared, read-only) before unioning in the
	// per-user "their own" id, so we never mutate cache state. The result may be
	// empty (a configured role whose only granted users were deleted) — that is
	// the intended "see nothing", enforced by the IN (0) sentinel downstream.
	allowed := make(map[uint]struct{}, len(entry.paidByUserIds)+1)
	for id := range entry.paidByUserIds {
		allowed[id] = struct{}{}
	}
	if entry.includeOwnPaidReceipts {
		allowed[userId] = struct{}{}
	}

	return allowed, false, nil
}

// GetVisibleCategoriesForUser filters allCategories down to the ones a user may
// see in a group. Unrestricted users get allCategories unchanged. Used by
// GetAppData to build the per-group category catalog.
func (service PermissionService) GetVisibleCategoriesForUser(userId uint, groupId uint, allCategories []models.Category) ([]models.Category, error) {
	allowed, unrestricted, err := service.GetGroupCategoryIdsForUser(userId, groupId)
	if err != nil {
		return nil, err
	}
	if unrestricted {
		return allCategories, nil
	}

	visible := make([]models.Category, 0, len(allowed))
	for _, category := range allCategories {
		if _, ok := allowed[category.ID]; ok {
			visible = append(visible, category)
		}
	}
	return visible, nil
}

// GetVisibleTagsForUser is the tag counterpart of GetVisibleCategoriesForUser.
func (service PermissionService) GetVisibleTagsForUser(userId uint, groupId uint, allTags []models.Tag) ([]models.Tag, error) {
	allowed, unrestricted, err := service.GetGroupTagIdsForUser(userId, groupId)
	if err != nil {
		return nil, err
	}
	if unrestricted {
		return allTags, nil
	}

	visible := make([]models.Tag, 0, len(allowed))
	for _, tag := range allTags {
		if _, ok := allowed[tag.ID]; ok {
			visible = append(visible, tag)
		}
	}
	return visible, nil
}

// resolveGroupRoleGrants returns the cached grant entry for a user's role in a
// group, or nil when the user is not a member or the membership has no group role
// (both mean "unrestricted").
func (service PermissionService) resolveGroupRoleGrants(userId uint, groupId uint) (*grantEntry, error) {
	roleRepository := repositories.NewRoleRepository(service.TX)

	groupRoleId, err := roleRepository.GetGroupMemberRoleId(userId, groupId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if groupRoleId == nil {
		return nil, nil
	}

	return loadGroupRoleGrants(roleRepository, *groupRoleId)
}

// loadGroupRoleGrants returns a group role's grant sets, consulting the cache
// first and populating it on a miss (mirroring loadRolePermissions).
func loadGroupRoleGrants(roleRepository repositories.RoleRepository, roleId uint) (*grantEntry, error) {
	if cached, ok := getCachedGroupRoleGrants(roleId); ok {
		return cached, nil
	}

	// Capture the eviction generation before reading so a concurrent grant
	// update/delete invalidates this write instead of being undone by it.
	observedGen := groupRoleGrantCacheGen()

	categoryIds, err := roleRepository.GetGroupRoleCategoryIds(roleId)
	if err != nil {
		return nil, err
	}
	tagIds, err := roleRepository.GetGroupRoleTagIds(roleId)
	if err != nil {
		return nil, err
	}
	paidByUserIds, err := roleRepository.GetGroupRolePaidByUserIds(roleId)
	if err != nil {
		return nil, err
	}
	includeOwnPaidReceipts, paidByVisibilityRestricted, err := roleRepository.GetGroupRolePaidByConfig(roleId)
	if err != nil {
		return nil, err
	}
	reportTemplateGrantRows, err := roleRepository.GetGroupRoleReportTemplateGrants(roleId)
	if err != nil {
		return nil, err
	}
	reportTemplateGrantsRestricted, err := roleRepository.GetGroupRoleReportTemplateGrantsRestricted(roleId)
	if err != nil {
		return nil, err
	}

	entry := &grantEntry{
		categoryIds:                    uintSliceToSet(categoryIds),
		tagIds:                         uintSliceToSet(tagIds),
		paidByUserIds:                  uintSliceToSet(paidByUserIds),
		includeOwnPaidReceipts:         includeOwnPaidReceipts,
		paidByVisibilityRestricted:     paidByVisibilityRestricted,
		reportTemplateGrants:           reportTemplateGrantsToSet(reportTemplateGrantRows),
		reportTemplateGrantsRestricted: reportTemplateGrantsRestricted,
	}

	setCachedGroupRoleGrants(roleId, entry, observedGen)
	return entry, nil
}

// reportTemplateGrantsToSet folds the flat grant rows (one per template+action)
// into a template id -> action set map for O(1) membership tests.
func reportTemplateGrantsToSet(rows []models.GroupRoleReportTemplateGrant) map[uint]map[string]struct{} {
	set := make(map[uint]map[string]struct{}, len(rows))
	for _, row := range rows {
		actions, ok := set[row.ReportTemplateID]
		if !ok {
			actions = make(map[string]struct{})
			set[row.ReportTemplateID] = actions
		}
		actions[row.Permission] = struct{}{}
	}
	return set
}

// uintSliceToSet converts a slice of ids to a set for O(1) membership tests.
func uintSliceToSet(ids []uint) map[uint]struct{} {
	set := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}
