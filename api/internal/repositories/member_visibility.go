package repositories

import (
	"receipt-wrangler/api/internal/models"
)

// ViewerGroupRow describes one of a viewer's group memberships for member-isolation
// resolution: the group, whether that group isolates its members, and whether the
// viewer's own role in it is a SeesAllMembers (supervisor) role. ViewerSeesAll is
// nil when the membership has no group role (treated as false — a plain member).
type ViewerGroupRow struct {
	GroupID        uint
	IsolateMembers bool
	ViewerSeesAll  *bool
}

// GetViewerGroupRow returns the single (viewer, group) row for member-isolation
// resolution: the group's isolation flag, whether the viewer's role there sees all
// members, and whether the viewer is a MEMBER of the group. It starts FROM groups and
// LEFT JOINs the viewer's membership, so it returns the isolation flag even when the
// viewer is not a member — needed to decide whether a non-member (e.g. an app.groups.read
// reader hitting GetGroupById) may see an isolated group's roster. Returns found=false
// only when the group itself does not exist. It is the single-group counterpart of
// GetViewerGroupRows, used by the per-group visibility resolver.
func (repository GroupMemberRepository) GetViewerGroupRow(viewerId uint, groupId uint) (ViewerGroupRow, bool, bool, error) {
	db := repository.GetDB()

	type viewerGroupRowScan struct {
		GroupID        uint
		IsolateMembers bool
		ViewerSeesAll  *bool
		MemberUserId   *uint
	}
	var rows []viewerGroupRowScan

	err := db.Table("groups AS g").
		Select("g.id AS group_id, g.isolate_members AS isolate_members, grd.sees_all_members AS viewer_sees_all, gm.user_id AS member_user_id").
		Joins("LEFT JOIN group_members AS gm ON gm.group_id = g.id AND gm.user_id = ?", viewerId).
		Joins("LEFT JOIN group_role_definitions AS grd ON grd.id = gm.group_role_id").
		Where("g.id = ?", groupId).
		Limit(1).
		Scan(&rows).Error
	if err != nil {
		return ViewerGroupRow{}, false, false, err
	}
	if len(rows) == 0 {
		return ViewerGroupRow{}, false, false, nil
	}
	row := rows[0]
	return ViewerGroupRow{
		GroupID:        row.GroupID,
		IsolateMembers: row.IsolateMembers,
		ViewerSeesAll:  row.ViewerSeesAll,
	}, true, row.MemberUserId != nil, nil
}

// GetViewerGroupRows returns, for each group the viewer belongs to, the group's
// isolation flag and whether the viewer's role there sees all members. Nullable
// ViewerSeesAll avoids a cross-engine COALESCE on a boolean column.
func (repository GroupMemberRepository) GetViewerGroupRows(viewerId uint) ([]ViewerGroupRow, error) {
	db := repository.GetDB()
	var rows []ViewerGroupRow

	err := db.Table("group_members AS gm").
		Select("gm.group_id AS group_id, g.isolate_members AS isolate_members, grd.sees_all_members AS viewer_sees_all").
		Joins("JOIN groups AS g ON g.id = gm.group_id").
		Joins("LEFT JOIN group_role_definitions AS grd ON grd.id = gm.group_role_id").
		Where("gm.user_id = ?", viewerId).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// GetDistinctUserIdsInGroups returns every distinct member user id across the
// given groups.
func (repository GroupMemberRepository) GetDistinctUserIdsInGroups(groupIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(groupIds) == 0 {
		return ids, nil
	}

	err := repository.GetDB().
		Model(&models.GroupMember{}).
		Where("group_id IN ?", groupIds).
		Distinct().
		Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetSupervisorUserIdsInGroups returns the distinct member user ids across the
// given groups whose group role is flagged SeesAllMembers (supervisors — visible
// to every isolated member).
func (repository GroupMemberRepository) GetSupervisorUserIdsInGroups(groupIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(groupIds) == 0 {
		return ids, nil
	}

	err := repository.GetDB().
		Table("group_members AS gm").
		Joins("JOIN group_role_definitions AS grd ON grd.id = gm.group_role_id").
		Where("gm.group_id IN ? AND grd.sees_all_members = ?", groupIds, true).
		Distinct().
		Pluck("gm.user_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetSupervisorUserIdsByGroup returns a map of group_id -> distinct supervisor user ids
// (members whose group role is flagged SeesAllMembers) across the given groups, so a
// multi-group roster filter resolves every group's supervisors in ONE query rather than
// one per group. A group with no supervisors is simply absent from the map.
func (repository GroupMemberRepository) GetSupervisorUserIdsByGroup(groupIds []uint) (map[uint][]uint, error) {
	result := map[uint][]uint{}
	if len(groupIds) == 0 {
		return result, nil
	}

	type supervisorRow struct {
		GroupID uint
		UserID  uint
	}
	var rows []supervisorRow
	err := repository.GetDB().
		Table("group_members AS gm").
		Select("gm.group_id AS group_id, gm.user_id AS user_id").
		Joins("JOIN group_role_definitions AS grd ON grd.id = gm.group_role_id").
		Where("gm.group_id IN ? AND grd.sees_all_members = ?", groupIds, true).
		Distinct().
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.GroupID] = append(result[row.GroupID], row.UserID)
	}
	return result, nil
}
