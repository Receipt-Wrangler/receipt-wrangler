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
