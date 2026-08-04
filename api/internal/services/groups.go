package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrGroupMemberChangeForbidden is returned by AuthorizeGroupMemberChanges when
// the caller is not permitted to make the requested group member/role changes.
var ErrGroupMemberChangeForbidden = errors.New("caller is not authorized to make the requested group member changes")

type GroupService struct {
	BaseService
}

func NewGroupService(tx *gorm.DB) GroupService {
	service := GroupService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

// AuthorizeGroupMemberChanges enforces who may change a group's membership and
// how far those changes may reach, given the roster the caller submitted relative
// to the group's current membership. It closes GHSA-89mm-9qfv-cjg3 (a member with
// group.update rewriting every member's group role — including their own — to
// escalate to owner or evict the owner).
//
// Two independent checks apply to every row that is added, role-changed, or
// removed (unchanged rows are never checked, so editing only the group's
// name/settings never trips this):
//
//   - CRUD gate: adding a member requires group.members.create, changing a
//     member's role requires group.members.update, removing a member requires
//     group.members.delete.
//   - Privilege ceiling ("you can neither grant nor strip a privilege you do not
//     hold"): the caller may only assign, or remove/replace, a role whose
//     permissions are a subset of the caller's own current group permissions. This
//     is what actually prevents self-escalation to owner and eviction of a
//     more-privileged member.
func (service GroupService) AuthorizeGroupMemberChanges(callerId uint, groupId uint, submitted []commands.UpsertGroupMemberCommand) error {
	permissionService := NewPermissionService(service.TX)
	roleRepository := repositories.NewRoleRepository(service.TX)

	callerPerms, err := permissionService.GetGroupPermissionsForUser(callerId, groupId)
	if err != nil {
		return err
	}
	hasCreate := permissions.HasAll(callerPerms, permissions.GroupMembersCreate)
	hasUpdate := permissions.HasAll(callerPerms, permissions.GroupMembersUpdate)
	hasDelete := permissions.HasAll(callerPerms, permissions.GroupMembersDelete)

	var existingMembers []models.GroupMember
	if err := service.GetDB().Where("group_id = ?", groupId).Find(&existingMembers).Error; err != nil {
		return err
	}

	existingRoleByUser := make(map[uint]*uint, len(existingMembers))
	for _, member := range existingMembers {
		existingRoleByUser[member.UserID] = member.GroupRoleID
	}
	// Reject a roster that lists the same user twice. The repository persists the
	// raw submitted slice, so duplicate (userId, groupId) rows would resolve in a
	// database-dependent way — a deduplicated authorization check must never diverge
	// from what is actually written, and a member can only appear once anyway.
	submittedUserIds := make(map[uint]bool, len(submitted))
	for _, member := range submitted {
		// A member entry must target the group in the URL. The repository also
		// forces this, but reject a mismatch here so the authorization decision is
		// made against exactly what will be written (0 means "unset" — the
		// repository scopes it to this group). This prevents authorizing against one
		// group while a body-supplied groupId points at another.
		if member.GroupID != 0 && member.GroupID != groupId {
			return ErrGroupMemberChangeForbidden
		}
		if submittedUserIds[member.UserID] {
			return ErrGroupMemberChangeForbidden
		}
		submittedUserIds[member.UserID] = true
	}

	rolesEqual := func(a, b *uint) bool {
		if a == nil || b == nil {
			return a == b
		}
		return *a == *b
	}

	// callerCanWield reports whether the caller may hand out or take away a role —
	// true when the role's permissions are a subset of the caller's own group
	// permissions. A nil role, or one that grants nothing, is always wieldable
	// (HasAll denies an empty required set, so short-circuit it here).
	rolePermsCache := make(map[uint][]string)
	callerCanWield := func(roleId *uint) (bool, error) {
		if roleId == nil {
			return true, nil
		}
		rolePerms, cached := rolePermsCache[*roleId]
		if !cached {
			loaded, loadErr := roleRepository.GetGroupRolePermissions(*roleId)
			if loadErr != nil {
				return false, loadErr
			}
			rolePerms = loaded
			rolePermsCache[*roleId] = rolePerms
		}
		if len(rolePerms) == 0 {
			return true, nil
		}
		return permissions.HasAll(callerPerms, rolePerms...), nil
	}

	// Additions and role changes — evaluate every submitted entry directly (not a
	// deduplicated map) so the authorization check can never diverge from the raw
	// roster the repository persists.
	for _, member := range submitted {
		existingRole, isExisting := existingRoleByUser[member.UserID]
		if !isExisting {
			if !hasCreate {
				return ErrGroupMemberChangeForbidden
			}
			canAssign, err := callerCanWield(member.GroupRoleID)
			if err != nil {
				return err
			}
			if !canAssign {
				return ErrGroupMemberChangeForbidden
			}
			continue
		}

		if !rolesEqual(existingRole, member.GroupRoleID) {
			if !hasUpdate {
				return ErrGroupMemberChangeForbidden
			}
			canAssign, err := callerCanWield(member.GroupRoleID)
			if err != nil {
				return err
			}
			canReplace, err := callerCanWield(existingRole)
			if err != nil {
				return err
			}
			if !canAssign || !canReplace {
				return ErrGroupMemberChangeForbidden
			}
		}
	}

	// Removals: existing members dropped from the submitted roster.
	for userId, existingRole := range existingRoleByUser {
		if submittedUserIds[userId] {
			continue
		}
		if !hasDelete {
			return ErrGroupMemberChangeForbidden
		}
		canRemove, err := callerCanWield(existingRole)
		if err != nil {
			return err
		}
		if !canRemove {
			return ErrGroupMemberChangeForbidden
		}
	}

	return nil
}

func (service GroupService) GetGroupIdsForUser(userId string) ([]uint, error) {
	groupMemberRepository := repositories.NewGroupMemberRepository(nil)
	groupMembers, err := groupMemberRepository.GetGroupMembersByUserId(userId)
	if err != nil {
		return nil, err
	}

	groupIds := make([]uint, len(groupMembers))
	for i := 0; i < len(groupMembers); i++ {
		groupIds[i] = groupMembers[i].GroupID
	}

	return groupIds, nil
}

func (service GroupService) GetGroupsForUser(userId string) ([]models.Group, error) {
	db := service.GetDB()
	var groups []models.Group

	groupIds, err := service.GetGroupIdsForUser(userId)
	if err != nil {
		return nil, err
	}

	err = db.Model(models.Group{}).
		Where("id IN ?", groupIds).
		Preload(clause.Associations).
		Order("is_all_group desc").
		Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (service GroupService) DeleteGroup(groupId string, allowAllGroupDelete bool) error {
	db := service.GetDB()
	var receipts []models.Receipt

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return err
	}

	groupRepository := repositories.NewGroupRepository(nil)

	if !allowAllGroupDelete {
		isAllGroup, err := groupRepository.IsAllGroup(uintGroupId)
		if err != nil || isAllGroup {
			return err
		}
	}

	group, err := groupRepository.GetGroupById(groupId, false, false, false)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		receiptService := NewReceiptService(tx)

		txErr := tx.Model(models.Receipt{}).Where("group_id = ?", groupId).Find(&receipts).Error
		if txErr != nil {
			return txErr
		}

		// Delete receipts in group
		for i := 0; i < len(receipts); i++ {
			txErr = receiptService.DeleteReceipt(utils.UintToString(receipts[i].ID))
			if txErr != nil {
				return txErr
			}
		}

		// Delete dashboards in group
		dashboardRepository := repositories.NewDashboardRepository(tx)
		groupDashboards, txErr := dashboardRepository.GetDashboardsByGroupId(group.ID)
		if txErr != nil {
			return txErr
		}

		for _, dashboard := range groupDashboards {
			txErr = dashboardRepository.DeleteDashboardById(dashboard.ID)
			if txErr != nil {
				return txErr
			}
		}

		// Delete group members
		txErr = tx.Where("group_id = ?", groupId).Delete(&models.GroupMember{}).Error
		if txErr != nil {
			return txErr
		}

		// Delete the members' per-member category/tag grants. The membership delete
		// above is raw, so these do not cascade; orphaned rows would be re-adopted if
		// the group id were ever reused.
		txErr = repositories.DeleteMemberGrantsForGroup(tx, uintGroupId)
		if txErr != nil {
			return txErr
		}

		// Unset user preferences
		tx.Model(models.UserPrefernces{}).Where("quick_scan_default_group_id = ?", groupId).Update("quick_scan_default_group_id", nil)

		// Delete Group Settings
		if group.GroupSettings.GroupId > 0 {
			txErr = tx.Model(&group.GroupSettings).Select(clause.Associations).Delete(&group.GroupSettings).Error
			if txErr != nil {
				return txErr
			}
		}

		// Delete Group Receipt Settings
		if group.GroupReceiptSettings.GroupId > 0 {
			txErr = tx.Model(&group.GroupReceiptSettings).Select(clause.Associations).Delete(&group.GroupReceiptSettings).Error
			if txErr != nil {
				return txErr
			}
		}

		// Delete group
		txErr = tx.Delete(&group).Error
		if txErr != nil {
			return txErr
		}

		// Remove group's directory
		groupPath, txErr := utils.BuildGroupPathString(utils.UintToString(group.ID), group.Name)
		if txErr != nil {
			return txErr
		}

		txErr = utils.RemoveDataPath(groupPath)
		if txErr != nil {
			logging.LogStd(logging.LOG_LEVEL_INFO, txErr.Error())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
