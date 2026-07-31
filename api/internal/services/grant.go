package services

import (
	"errors"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"

	"gorm.io/gorm"
)

// GetGroupCategoryIdsForUser returns the set of category ids a user may use in a
// group, plus an "unrestricted" flag. When unrestricted is true the returned set
// is nil and the caller must treat every category as allowed.
//
// The result composes TWO layers — the group role's grants and the individual
// membership's grants — by intersection; see resolveEffectiveGrants for the rule.
// A non-member, or a member whose role and membership both grant nothing, is
// unrestricted: grants only NARROW access within an already-permitted group; they
// never grant access (the handler permission gate is the access control).
//
// The returned map must be treated as read-only — when only the role layer
// narrows, it is shared cache state.
func (service PermissionService) GetGroupCategoryIdsForUser(userId uint, groupId uint) (map[uint]struct{}, bool, error) {
	grants, err := service.resolveEffectiveGrants(userId, groupId)
	if err != nil {
		return nil, false, err
	}
	if grants == nil || !grants.categoryRestricted {
		return nil, true, nil
	}
	return grants.categoryIds, false, nil
}

// GetGroupTagIdsForUser is the tag counterpart of GetGroupCategoryIdsForUser.
// Categories and tags resolve independently — a role or membership may restrict
// one and leave the other unrestricted.
func (service PermissionService) GetGroupTagIdsForUser(userId uint, groupId uint) (map[uint]struct{}, bool, error) {
	grants, err := service.resolveEffectiveGrants(userId, groupId)
	if err != nil {
		return nil, false, err
	}
	if grants == nil || !grants.tagRestricted {
		return nil, true, nil
	}
	return grants.tagIds, false, nil
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

// effectiveGrantSet is the composed category/tag visibility for one (user, group),
// after the group role's grants and the individual membership's grants have been
// folded together. A false *Restricted flag means "see all" and the paired set is
// meaningless; a true flag with an EMPTY set means "see nothing".
type effectiveGrantSet struct {
	categoryIds        map[uint]struct{}
	categoryRestricted bool
	tagIds             map[uint]struct{}
	tagRestricted      bool
}

// resolveEffectiveGrants composes a user's group-role grants with their
// individual membership grants for a group. Returns nil when the user is not a
// member of the group (nothing narrows).
//
// The two layers compose by INTERSECTION, with the role as a ceiling:
//
//	role requires individual && member unconfigured -> see nothing (fail closed)
//	effective = ALL
//	if the role grants a non-empty set          -> effective ∩= role set
//	if the membership opted into restriction    -> effective ∩= membership set
//	if neither narrows                          -> unrestricted
//
// Union was rejected deliberately: it can only ever ADD, so a role carrying grants
// would floor every member at the role's full set and an individual assignment
// could never restrict anyone — which is the entire point of the membership layer.
// Intersection also keeps a role an auditable ceiling ("no holder of this role can
// see outside its set") and makes role grants safe to widen, since widening a role
// never widens an individually-assigned member.
//
// The empty-intersection case (a membership granted ids outside its role's
// ceiling, resolving to "see nothing") is unreachable through the API: the grants
// endpoint validates every submitted id against the ceiling and rejects the write.
// It is still handled here rather than assumed away, because a role can be
// narrowed AFTER a membership was configured — that member correctly loses the
// out-of-ceiling ids instead of keeping visibility the role no longer allows.
func (service PermissionService) resolveEffectiveGrants(userId uint, groupId uint) (*effectiveGrantSet, error) {
	memberRepository := repositories.NewGroupMemberRepository(service.TX)

	member, err := memberRepository.GetMemberGrantContext(userId, groupId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// A membership with no group role still honors its individual grants — role
	// assignment is optional, individual assignment is not conditional on it.
	var role *grantEntry
	if member.GroupRoleID != nil {
		role, err = loadGroupRoleGrants(repositories.NewRoleRepository(service.TX), *member.GroupRoleID)
		if err != nil {
			return nil, err
		}
	}

	grants := &effectiveGrantSet{}

	grants.categoryIds, grants.categoryRestricted, err = composeGrantLayer(
		roleGrantSet(role, func(entry *grantEntry) map[uint]struct{} { return entry.categoryIds }),
		role != nil && role.requiresIndividualCategoryGrants,
		member.CategoryGrantsRestricted,
		func() ([]uint, error) { return memberRepository.GetMemberCategoryGrantIds(userId, groupId) },
	)
	if err != nil {
		return nil, err
	}

	grants.tagIds, grants.tagRestricted, err = composeGrantLayer(
		roleGrantSet(role, func(entry *grantEntry) map[uint]struct{} { return entry.tagIds }),
		role != nil && role.requiresIndividualTagGrants,
		member.TagGrantsRestricted,
		func() ([]uint, error) { return memberRepository.GetMemberTagGrantIds(userId, groupId) },
	)
	if err != nil {
		return nil, err
	}

	return grants, nil
}

// roleGrantSet reads one resource's grant set off a possibly-nil role entry.
func roleGrantSet(role *grantEntry, pick func(*grantEntry) map[uint]struct{}) map[uint]struct{} {
	if role == nil {
		return nil
	}
	return pick(role)
}

// composeGrantLayer intersects one resource's role ceiling with the membership's
// own grants, returning (allowed, restricted). Loading the membership ids is
// deferred behind loadMemberIds so the query is skipped entirely for the common
// unconfigured membership.
func composeGrantLayer(
	roleIds map[uint]struct{},
	requiresIndividual bool,
	memberRestricted bool,
	loadMemberIds func() ([]uint, error),
) (map[uint]struct{}, bool, error) {
	// The role demands individual assignment and this member has none. Fail closed
	// rather than falling back to the role's set (or to see-all).
	if requiresIndividual && !memberRestricted {
		return map[uint]struct{}{}, true, nil
	}

	if !memberRestricted {
		// Only the role layer can narrow. An empty role set is the long-standing
		// "unrestricted" signal for category/tag grants.
		if len(roleIds) == 0 {
			return nil, false, nil
		}
		return roleIds, true, nil
	}

	memberIds, err := loadMemberIds()
	if err != nil {
		return nil, false, err
	}

	// Freshly allocated — never the shared, read-only cache map.
	allowed := make(map[uint]struct{}, len(memberIds))
	for _, id := range memberIds {
		// An empty role set is unrestricted, so it imposes no ceiling to test against.
		if len(roleIds) > 0 {
			if _, withinCeiling := roleIds[id]; !withinCeiling {
				continue
			}
		}
		allowed[id] = struct{}{}
	}

	return allowed, true, nil
}

// resolveGroupRoleGrants returns the cached grant entry for a user's role in a
// group, or nil when the user is not a member or the membership has no group role
// (both mean "unrestricted"). It is the ROLE-ONLY resolver, still used by the
// paid-by and report-template grant types; category/tag visibility goes through
// resolveEffectiveGrants, which layers individual membership grants on top.
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
	requiresIndividualCategoryGrants, requiresIndividualTagGrants, err := roleRepository.GetGroupRoleIndividualGrantConfig(roleId)
	if err != nil {
		return nil, err
	}

	entry := &grantEntry{
		categoryIds:                      uintSliceToSet(categoryIds),
		tagIds:                           uintSliceToSet(tagIds),
		paidByUserIds:                    uintSliceToSet(paidByUserIds),
		includeOwnPaidReceipts:           includeOwnPaidReceipts,
		paidByVisibilityRestricted:       paidByVisibilityRestricted,
		reportTemplateGrants:             reportTemplateGrantsToSet(reportTemplateGrantRows),
		reportTemplateGrantsRestricted:   reportTemplateGrantsRestricted,
		requiresIndividualCategoryGrants: requiresIndividualCategoryGrants,
		requiresIndividualTagGrants:      requiresIndividualTagGrants,
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
