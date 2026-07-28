package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
)

// GetVisibleUserIdsForUser resolves the set of user ids that viewerId is allowed
// to see, for member-presence isolation. It returns (set, unrestricted, err):
// when unrestricted is true the set is nil and the caller MUST treat every user
// as visible (no filtering).
//
// Invariants: a viewer always sees themselves; a holder of app.users.read sees
// everyone; a SeesAllMembers (supervisor) role — and membership in any
// non-isolated group — contributes all of that group's members.
//
// A viewer who is not an isolated (non-supervisor) member of any group is
// UNRESTRICTED, so non-isolated installs and every non-isolated member behave
// exactly as before — no directory or content change. Isolation only ever
// narrows what an isolated member sees.
//
// The result is a freshly allocated set. It is NOT cached, so callers should
// resolve it once per request and reuse it across a batch (mirroring the
// PaidByListResolver closure pattern) rather than calling per row.
func (service PermissionService) GetVisibleUserIdsForUser(viewerId uint) (map[uint]struct{}, bool, error) {
	isAdmin, err := service.HasAppPermissions(viewerId, permissions.AppUsersRead)
	if err != nil {
		return nil, false, err
	}
	if isAdmin {
		return nil, true, nil
	}

	memberRepository := repositories.NewGroupMemberRepository(service.TX)
	rows, err := memberRepository.GetViewerGroupRows(viewerId)
	if err != nil {
		return nil, false, err
	}

	// Groups where the viewer sees every member (non-isolated groups, or isolated
	// groups where the viewer holds a supervisor role) vs. isolated groups where the
	// viewer sees only supervisors.
	var seeAllGroupIds []uint
	var isolatedRestrictedGroupIds []uint
	for _, row := range rows {
		viewerIsSupervisor := row.ViewerSeesAll != nil && *row.ViewerSeesAll
		if !row.IsolateMembers || viewerIsSupervisor {
			seeAllGroupIds = append(seeAllGroupIds, row.GroupID)
		} else {
			isolatedRestrictedGroupIds = append(isolatedRestrictedGroupIds, row.GroupID)
		}
	}

	// Not isolated in any group ⇒ unrestricted (backward compatible).
	if len(isolatedRestrictedGroupIds) == 0 {
		return nil, true, nil
	}

	visible := map[uint]struct{}{viewerId: {}} // self is always visible

	seeAllIds, err := memberRepository.GetDistinctUserIdsInGroups(seeAllGroupIds)
	if err != nil {
		return nil, false, err
	}
	for _, id := range seeAllIds {
		visible[id] = struct{}{}
	}

	supervisorIds, err := memberRepository.GetSupervisorUserIdsInGroups(isolatedRestrictedGroupIds)
	if err != nil {
		return nil, false, err
	}
	for _, id := range supervisorIds {
		visible[id] = struct{}{}
	}

	return visible, false, nil
}

// GetVisibleUserIdsForUserInGroup resolves the set of user ids viewerId is allowed to
// see WITHIN groupId, for member-presence isolation. It returns (set, unrestricted,
// err); when unrestricted is true the set is nil and the caller MUST treat every user
// as visible (no filtering).
//
// Unlike GetVisibleUserIdsForUser (which unions every group the viewer belongs to and
// backs the flat user directory), this resolves visibility for ONE group in isolation:
// the viewer's presence in some OTHER group never widens what they may see here. This
// is what makes isolation truthful — "isolated means isolated," regardless of any open
// group the viewer also shares.
//
// Invariants: a holder of app.users.read sees everyone (unrestricted); a non-isolated
// group is unrestricted; in an isolated group a SeesAllMembers (supervisor) role is
// unrestricted, while a plain member sees only themselves and that group's supervisors.
// A non-member of an ISOLATED group is restricted to {self} — an isolated roster must not
// leak to a non-member reader (e.g. an app.groups.read holder hitting GetGroupById) who
// lacks the app.users.read directory exemption checked above; a non-member of a
// NON-isolated group is unrestricted (open group), which preserves the paid-by / reporting
// contract for non-members (those surfaces gate membership at the handler).
//
// The result is a freshly allocated set and is NOT cached; callers resolving it across a
// batch of groups should use a groupVisibilityResolver (below) to memoize per group.
func (service PermissionService) GetVisibleUserIdsForUserInGroup(viewerId uint, groupId uint) (map[uint]struct{}, bool, error) {
	isAdmin, err := service.HasAppPermissions(viewerId, permissions.AppUsersRead)
	if err != nil {
		return nil, false, err
	}
	if isAdmin {
		return nil, true, nil
	}

	memberRepository := repositories.NewGroupMemberRepository(service.TX)
	row, found, isMember, err := memberRepository.GetViewerGroupRow(viewerId, groupId)
	if err != nil {
		return nil, false, err
	}
	if !found {
		// The group does not exist — nothing to filter.
		return nil, true, nil
	}
	if !isMember {
		// A non-member (non-admin) reaching a group-scoped surface. An isolated group's
		// roster must stay hidden from them; a non-isolated group adds no restriction.
		if row.IsolateMembers {
			return map[uint]struct{}{viewerId: {}}, false, nil
		}
		return nil, true, nil
	}

	viewerIsSupervisor := row.ViewerSeesAll != nil && *row.ViewerSeesAll
	if !row.IsolateMembers || viewerIsSupervisor {
		return nil, true, nil
	}

	visible := map[uint]struct{}{viewerId: {}} // self is always visible
	supervisorIds, err := memberRepository.GetSupervisorUserIdsInGroups([]uint{groupId})
	if err != nil {
		return nil, false, err
	}
	for _, id := range supervisorIds {
		visible[id] = struct{}{}
	}
	return visible, false, nil
}

// ActivityVisibilityResolver returns the injected resolver the system-task repository uses
// to filter activities by member isolation IN SQL: for a group it reports the ran-by user
// ids the caller may see, or unrestricted == true (see every actor). Mirrors
// PaidByListResolver, so the repository stays free of the service layer.
func (service PermissionService) ActivityVisibilityResolver(userId uint) repositories.ActivityVisibilityResolver {
	return func(groupId uint) ([]uint, bool, error) {
		set, unrestricted, err := service.GetVisibleUserIdsForUserInGroup(userId, groupId)
		if err != nil {
			return nil, false, err
		}
		if unrestricted {
			return nil, true, nil
		}
		return uintSetToSlice(set), false, nil
	}
}

// groupVisibility is one group's resolved member-visible set. unrestricted == true
// (with a nil set) means "see every member of this group."
type groupVisibility struct {
	visible      map[uint]struct{}
	unrestricted bool
}

// groupVisibilityResolver memoizes GetVisibleUserIdsForUserInGroup by group id for the
// lifetime of one request/batch, so a batch of receipts, activities, or roster rows
// spanning several groups resolves each group's set exactly once. Mirrors the per-group
// cache pattern in FilterReceiptsByPaidBy.
type groupVisibilityResolver struct {
	service  PermissionService
	viewerId uint
	cache    map[uint]groupVisibility
}

func (service PermissionService) newGroupVisibilityResolver(viewerId uint) *groupVisibilityResolver {
	return &groupVisibilityResolver{
		service:  service,
		viewerId: viewerId,
		cache:    map[uint]groupVisibility{},
	}
}

// forGroup returns (set, unrestricted) for groupId, resolving and caching on first use.
func (resolver *groupVisibilityResolver) forGroup(groupId uint) (map[uint]struct{}, bool, error) {
	if cached, ok := resolver.cache[groupId]; ok {
		return cached.visible, cached.unrestricted, nil
	}
	visible, unrestricted, err := resolver.service.GetVisibleUserIdsForUserInGroup(resolver.viewerId, groupId)
	if err != nil {
		return nil, false, err
	}
	resolver.cache[groupId] = groupVisibility{visible: visible, unrestricted: unrestricted}
	return visible, unrestricted, nil
}

// isVisible reports whether targetId is visible to the viewer within groupId.
func (resolver *groupVisibilityResolver) isVisible(groupId uint, targetId uint) (bool, error) {
	visible, unrestricted, err := resolver.forGroup(groupId)
	if err != nil {
		return false, err
	}
	if unrestricted {
		return true, nil
	}
	return isUserVisible(targetId, resolver.viewerId, visible), nil
}

// visibleUserIdsByGroup resolves the member-visible set for each of the given groups in
// a fixed TWO queries (rather than two per group), for the AppData roster path where N
// can be large. Admins short-circuit to a single unrestricted marker (nil map, true).
// The returned map is keyed by group id; a group the viewer does not belong to resolves
// to unrestricted (isolation only narrows for actual members; the AppData path only
// ever passes the viewer's own groups anyway).
func (service PermissionService) visibleUserIdsByGroup(viewerId uint, groupIds []uint) (map[uint]groupVisibility, bool, error) {
	isAdmin, err := service.HasAppPermissions(viewerId, permissions.AppUsersRead)
	if err != nil {
		return nil, false, err
	}
	if isAdmin {
		return nil, true, nil
	}

	memberRepository := repositories.NewGroupMemberRepository(service.TX)
	rows, err := memberRepository.GetViewerGroupRows(viewerId)
	if err != nil {
		return nil, false, err
	}

	rowByGroup := make(map[uint]repositories.ViewerGroupRow, len(rows))
	var isolatedRestrictedGroupIds []uint
	for _, row := range rows {
		rowByGroup[row.GroupID] = row
		viewerIsSupervisor := row.ViewerSeesAll != nil && *row.ViewerSeesAll
		if row.IsolateMembers && !viewerIsSupervisor {
			isolatedRestrictedGroupIds = append(isolatedRestrictedGroupIds, row.GroupID)
		}
	}

	supervisorsByGroup, err := memberRepository.GetSupervisorUserIdsByGroup(isolatedRestrictedGroupIds)
	if err != nil {
		return nil, false, err
	}

	result := make(map[uint]groupVisibility, len(groupIds))
	for _, groupId := range groupIds {
		row, isMember := rowByGroup[groupId]
		if !isMember {
			result[groupId] = groupVisibility{visible: nil, unrestricted: true}
			continue
		}
		viewerIsSupervisor := row.ViewerSeesAll != nil && *row.ViewerSeesAll
		if !row.IsolateMembers || viewerIsSupervisor {
			result[groupId] = groupVisibility{visible: nil, unrestricted: true}
			continue
		}
		visible := map[uint]struct{}{viewerId: {}}
		for _, id := range supervisorsByGroup[groupId] {
			visible[id] = struct{}{}
		}
		result[groupId] = groupVisibility{visible: visible, unrestricted: false}
	}
	return result, false, nil
}

// FilterVisibleUserViews returns the subset of users visible to viewerId, always
// retaining the viewer's own view. It is a pass-through when the viewer is
// unrestricted.
func (service PermissionService) FilterVisibleUserViews(viewerId uint, users []structs.UserView) ([]structs.UserView, error) {
	visible, unrestricted, err := service.GetVisibleUserIdsForUser(viewerId)
	if err != nil {
		return nil, err
	}
	if unrestricted {
		return users, nil
	}

	filtered := make([]structs.UserView, 0, len(users))
	for _, user := range users {
		if isUserVisible(user.ID, viewerId, visible) {
			filtered = append(filtered, user)
		}
	}
	return filtered, nil
}

// FilterGroupMembersForGroups strips each group's GroupMembers to those visible to
// viewerId WITHIN THAT GROUP (the viewer's own membership is always retained). Each
// group is resolved independently, so an isolated group shows only self + its
// supervisors regardless of any open group the viewer also belongs to. This is the
// roster boundary for GetGroupsForUser (so MCP list_groups and AppData groups inherit
// it). Resolves every group's set in two queries via visibleUserIdsByGroup.
func (service PermissionService) FilterGroupMembersForGroups(viewerId uint, groups []models.Group) error {
	if len(groups) == 0 {
		return nil
	}

	groupIds := make([]uint, len(groups))
	for i := range groups {
		groupIds[i] = groups[i].ID
	}

	visibilityByGroup, unrestricted, err := service.visibleUserIdsByGroup(viewerId, groupIds)
	if err != nil {
		return err
	}
	if unrestricted {
		return nil
	}

	for i := range groups {
		visibility := visibilityByGroup[groups[i].ID]
		if visibility.unrestricted {
			continue
		}
		groups[i].GroupMembers = filterVisibleGroupMembers(groups[i].GroupMembers, viewerId, visibility.visible)
	}
	return nil
}

// FilterGroupMembersForGroup is the single-group counterpart of
// FilterGroupMembersForGroups (used by GetGroupById), resolving visibility for that one
// group.
func (service PermissionService) FilterGroupMembersForGroup(viewerId uint, group *models.Group) error {
	visible, unrestricted, err := service.GetVisibleUserIdsForUserInGroup(viewerId, group.ID)
	if err != nil {
		return err
	}
	if unrestricted {
		return nil
	}

	group.GroupMembers = filterVisibleGroupMembers(group.GroupMembers, viewerId, visible)
	return nil
}

// ValidateReceiptUserSelection reports whether every user id a receipt upsert plants
// — the payer (paidByUserId) and each item / linked-item chargedToUserId — is within
// the caller's member-visible set FOR THE RECEIPT'S GROUP. Unrestricted callers always
// pass. Used on receipt create/update so an isolated member cannot attach a user they
// may not see in that group (which would leak that user's presence, or hand them a
// receipt/charge). A zero id is skipped (unset). Returns (allowed, err); the caller
// responds 403 when not allowed.
func (service PermissionService) ValidateReceiptUserSelection(viewerId uint, groupId uint, paidByUserId uint, chargedToUserIds []uint) (bool, error) {
	visible, unrestricted, err := service.GetVisibleUserIdsForUserInGroup(viewerId, groupId)
	if err != nil {
		return false, err
	}
	if unrestricted {
		return true, nil
	}

	if paidByUserId != 0 && !isUserVisible(paidByUserId, viewerId, visible) {
		return false, nil
	}
	for _, id := range chargedToUserIds {
		if id != 0 && !isUserVisible(id, viewerId, visible) {
			return false, nil
		}
	}
	return true, nil
}

func filterVisibleGroupMembers(members []models.GroupMember, viewerId uint, visible map[uint]struct{}) []models.GroupMember {
	filtered := make([]models.GroupMember, 0, len(members))
	for _, member := range members {
		if isUserVisible(member.UserID, viewerId, visible) {
			filtered = append(filtered, member)
		}
	}
	return filtered
}

// isUserVisible reports whether targetId is in the viewer's visible set. The
// viewer is always visible to themselves, defensively, even if absent from the set.
func isUserVisible(targetId uint, viewerId uint, visible map[uint]struct{}) bool {
	if targetId == viewerId {
		return true
	}
	_, ok := visible[targetId]
	return ok
}
