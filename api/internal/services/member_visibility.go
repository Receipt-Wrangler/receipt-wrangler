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
// viewerId (the viewer's own membership is always retained). No-op when the viewer
// is unrestricted. This is the roster boundary for GetGroupsForUser (so MCP
// list_groups and AppData groups inherit it).
func (service PermissionService) FilterGroupMembersForGroups(viewerId uint, groups []models.Group) error {
	visible, unrestricted, err := service.GetVisibleUserIdsForUser(viewerId)
	if err != nil {
		return err
	}
	if unrestricted {
		return nil
	}

	for i := range groups {
		groups[i].GroupMembers = filterVisibleGroupMembers(groups[i].GroupMembers, viewerId, visible)
	}
	return nil
}

// FilterGroupMembersForGroup is the single-group counterpart of
// FilterGroupMembersForGroups (used by GetGroupById).
func (service PermissionService) FilterGroupMembersForGroup(viewerId uint, group *models.Group) error {
	visible, unrestricted, err := service.GetVisibleUserIdsForUser(viewerId)
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
// the caller's member-visible set. Unrestricted callers always pass. Used on receipt
// create/update so an isolated member cannot attach a user they may not see (which
// would leak that user's presence, or hand them a receipt/charge). A zero id is
// skipped (unset). Returns (allowed, err); the caller responds 403 when not allowed.
func (service PermissionService) ValidateReceiptUserSelection(viewerId uint, paidByUserId uint, chargedToUserIds []uint) (bool, error) {
	visible, unrestricted, err := service.GetVisibleUserIdsForUser(viewerId)
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
