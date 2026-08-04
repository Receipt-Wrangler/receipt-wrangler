package repositories

import (
	"maps"
	"receipt-wrangler/api/internal/models"
	"slices"

	"gorm.io/gorm"
)

type GroupMemberRepository struct {
	BaseRepository
}

func NewGroupMemberRepository(tx *gorm.DB) GroupMemberRepository {
	repository := GroupMemberRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

// Gets groupMembers that the user has access to
func (repository GroupMemberRepository) GetGroupMembersByUserId(userId string) ([]models.GroupMember, error) {
	db := repository.GetDB()
	var groupMembers []models.GroupMember

	err := db.Model(models.GroupMember{}).Where("user_id = ?", userId).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}

	return groupMembers, nil
}

// Gets group ids that the user has access to
func (repository GroupMemberRepository) GetGroupIdsByUserId(userId string) ([]uint, error) {
	groupMembers, err := repository.GetGroupMembersByUserId(userId)
	if err != nil {
		return nil, err
	}
	result := make([]uint, len(groupMembers))

	for i := 0; i < len(groupMembers); i++ {
		result[i] = groupMembers[i].GroupID
	}

	return result, nil
}

func (repository GroupMemberRepository) GetGroupMemberByUserIdAndGroupId(userId string, groupId string) (models.GroupMember, error) {
	db := repository.GetDB()
	var groupMember models.GroupMember

	err := db.Model(models.GroupMember{}).Where("user_id = ? AND group_id = ?", userId, groupId).First(&groupMember).Error
	if err != nil {
		return models.GroupMember{}, err
	}

	return groupMember, nil
}

func (repository GroupMemberRepository) GetsGroupMembersByGroupId(groupId string) ([]models.GroupMember, error) {
	db := repository.GetDB()
	var groupMembers []models.GroupMember

	err := db.Model(models.GroupMember{}).Where("group_id = ?", groupId).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}

	return groupMembers, nil
}

// GetMemberCategoryGrantIds returns the category ids granted to a specific group
// membership. An empty result means the membership adds no narrowing of its own —
// consult GetMemberGrantConfig to tell "never configured" apart from "configured
// but since emptied".
func (repository GroupMemberRepository) GetMemberCategoryGrantIds(userId uint, groupId uint) ([]uint, error) {
	db := repository.GetDB()
	var categoryIds []uint

	err := db.Model(&models.GroupMemberCategoryGrant{}).
		Where("user_id = ? AND group_id = ?", userId, groupId).
		Pluck("category_id", &categoryIds).Error
	if err != nil {
		return nil, err
	}

	return categoryIds, nil
}

// GetMemberTagGrantIds is the tag counterpart of GetMemberCategoryGrantIds.
func (repository GroupMemberRepository) GetMemberTagGrantIds(userId uint, groupId uint) ([]uint, error) {
	db := repository.GetDB()
	var tagIds []uint

	err := db.Model(&models.GroupMemberTagGrant{}).
		Where("user_id = ? AND group_id = ?", userId, groupId).
		Pluck("tag_id", &tagIds).Error
	if err != nil {
		return nil, err
	}

	return tagIds, nil
}

// GetMemberGrantContext returns everything the grant resolver needs about a
// membership in ONE row read: the group role it holds (nil when unassigned) and
// the two fail-closed restriction flags. Returns gorm.ErrRecordNotFound when the
// user is not a member of the group — the caller treats that as "unrestricted",
// since grants only narrow within an already-permitted group.
func (repository GroupMemberRepository) GetMemberGrantContext(userId uint, groupId uint) (models.GroupMember, error) {
	db := repository.GetDB()
	var member models.GroupMember

	err := db.Model(&models.GroupMember{}).
		Select("user_id", "group_id", "group_role_id", "category_grants_restricted", "tag_grants_restricted").
		Where("user_id = ? AND group_id = ?", userId, groupId).
		First(&member).Error
	if err != nil {
		return models.GroupMember{}, err
	}

	return member, nil
}

// ReplaceMemberGrants resets a membership's category and tag grants to exactly the
// given id sets and records the derived fail-closed flags, all in one transaction.
// Mirrors RoleRepository.replaceGroupRoleGrants (delete-all-then-insert, with the
// nested Category/Tag belongs-to associations Omit-ted so only join rows are
// written).
//
// The restricted flags are set from the SUBMITTED sets rather than from what
// survives in the table, so a membership stays restricted after a granted category
// is later deleted and its row cascades away.
func (repository GroupMemberRepository) ReplaceMemberGrants(userId uint, groupId uint, categoryIds []uint, tagIds []uint) error {
	return repository.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := DeleteMemberGrants(tx, userId, groupId); err != nil {
			return err
		}

		if len(categoryIds) > 0 {
			categoryGrants := make([]models.GroupMemberCategoryGrant, 0, len(categoryIds))
			for _, categoryId := range categoryIds {
				categoryGrants = append(categoryGrants, models.GroupMemberCategoryGrant{
					UserID: userId, GroupID: groupId, CategoryID: categoryId,
				})
			}

			if err := tx.Omit("Category").Create(&categoryGrants).Error; err != nil {
				return err
			}
		}

		if len(tagIds) > 0 {
			tagGrants := make([]models.GroupMemberTagGrant, 0, len(tagIds))
			for _, tagId := range tagIds {
				tagGrants = append(tagGrants, models.GroupMemberTagGrant{
					UserID: userId, GroupID: groupId, TagID: tagId,
				})
			}

			if err := tx.Omit("Tag").Create(&tagGrants).Error; err != nil {
				return err
			}
		}

		// Map form so a toggled-off false persists (a struct Updates skips zero-value
		// bools, which would leave a cleared restriction set).
		return tx.Model(&models.GroupMember{}).
			Where("user_id = ? AND group_id = ?", userId, groupId).
			Updates(map[string]interface{}{
				"category_grants_restricted": len(categoryIds) > 0,
				"tag_grants_restricted":      len(tagIds) > 0,
			}).Error
	})
}

// The Delete*MemberGrants functions below are package functions taking an explicit
// *gorm.DB so the membership-teardown paths (UserService.DeleteUser,
// GroupService.DeleteGroup, GroupRepository.UpdateGroup) can call them inside
// their existing transactions.
//
// Calling them is NOT optional on those paths. They remove group_members rows with
// a raw tx.Delete, which runs no Go-side association cascade, and SQLite does not
// enforce foreign keys unless PRAGMA foreign_keys=ON is set per connection. Because
// a grant row is keyed by (user, group, resource), orphaned rows would be silently
// re-adopted the next time the same user rejoins the same group — restoring
// category visibility that was deliberately revoked, with nothing in the UI to
// show it happened.

// MemberGrantFlags is a membership's pair of fail-closed restriction flags.
type MemberGrantFlags struct {
	CategoryGrantsRestricted bool
	TagGrantsRestricted      bool
}

// GetMemberGrantFlagsForGroup returns every member's restriction flags in a
// group, keyed by user id.
//
// Used by GroupRepository.UpdateGroup to carry the flags across its wholesale
// roster replace. That replace writes GroupMember rows built from the request
// command, which carry both flags at their zero value — so without this, a plain
// group edit (even just a rename) would clear every member's restriction and
// silently widen them back to their role's full set.
//
// The flags are carried forward rather than recomputed from the surviving grant
// rows on purpose: a membership that was configured and then emptied by a
// category deletion must STAY restricted, which is the entire reason the flags
// exist separately from the rows.
func GetMemberGrantFlagsForGroup(db *gorm.DB, groupId uint) (map[uint]MemberGrantFlags, error) {
	var members []models.GroupMember

	err := db.Model(&models.GroupMember{}).
		Select("user_id", "category_grants_restricted", "tag_grants_restricted").
		Where("group_id = ?", groupId).
		Find(&members).Error
	if err != nil {
		return nil, err
	}

	flags := make(map[uint]MemberGrantFlags, len(members))
	for _, member := range members {
		flags[member.UserID] = MemberGrantFlags{
			CategoryGrantsRestricted: member.CategoryGrantsRestricted,
			TagGrantsRestricted:      member.TagGrantsRestricted,
		}
	}

	return flags, nil
}

// deleteMemberGrantsWhere removes category and tag grant rows matching a
// condition. The two tables are always torn down together.
func deleteMemberGrantsWhere(db *gorm.DB, query interface{}, args ...interface{}) error {
	err := db.Where(query, args...).Delete(&models.GroupMemberCategoryGrant{}).Error
	if err != nil {
		return err
	}

	return db.Where(query, args...).Delete(&models.GroupMemberTagGrant{}).Error
}

// DeleteMemberGrants removes every grant row for one membership.
func DeleteMemberGrants(db *gorm.DB, userId uint, groupId uint) error {
	return deleteMemberGrantsWhere(db, "user_id = ? AND group_id = ?", userId, groupId)
}

// DeleteMemberGrantsForUser removes a user's grant rows across every group. Used
// when the user itself is deleted, so no group is missed.
func DeleteMemberGrantsForUser(db *gorm.DB, userId uint) error {
	return deleteMemberGrantsWhere(db, "user_id = ?", userId)
}

// DeleteMemberGrantsForGroup removes every member's grant rows in a group. Used
// when the group itself is deleted.
func DeleteMemberGrantsForGroup(db *gorm.DB, groupId uint) error {
	return deleteMemberGrantsWhere(db, "group_id = ?", groupId)
}

// DeleteOrphanedMemberGrants removes grant rows in a group that no longer have a
// matching group_members row. GroupRepository.UpdateGroup replaces the entire
// roster in a single association write, so this is how a member dropped by that
// path has their grants cleaned up. Phrased as "whatever no longer has a
// membership" rather than "the members we just removed", so it stays correct
// regardless of whether GORM's replace deletes, recreates, or upserts the
// surviving rows — and it leaves retained members' grants untouched.
func DeleteOrphanedMemberGrants(db *gorm.DB, groupId uint) error {
	remainingMemberIds := db.Model(&models.GroupMember{}).
		Select("user_id").
		Where("group_id = ?", groupId)

	return deleteMemberGrantsWhere(
		db,
		"group_id = ? AND user_id NOT IN (?)",
		groupId,
		remainingMemberIds,
	)
}

// LoadMemberGrantsForGroups populates every member's transient grant ids across a
// set of groups, in two queries total. Used at the roster serialization boundary
// (AppData, GetGroupsForUser, GetGroupById) so the desktop receives each member's
// assignment alongside their role.
func (repository GroupMemberRepository) LoadMemberGrantsForGroups(groups []models.Group) error {
	members := make([]*models.GroupMember, 0)
	for i := range groups {
		for j := range groups[i].GroupMembers {
			members = append(members, &groups[i].GroupMembers[j])
		}
	}

	return repository.LoadMemberGrants(members)
}

// LoadMemberGrantsForGroup is the single-group counterpart of
// LoadMemberGrantsForGroups (used by GetGroupById).
func (repository GroupMemberRepository) LoadMemberGrantsForGroup(group *models.Group) error {
	members := make([]*models.GroupMember, 0, len(group.GroupMembers))
	for i := range group.GroupMembers {
		members = append(members, &group.GroupMembers[i])
	}

	return repository.LoadMemberGrants(members)
}

// LoadMemberGrants populates the transient CategoryGrants/TagGrants id slices on
// each member, in two queries regardless of member count. Callers use it at the
// serialization boundary; the fields are `gorm:"-"` so nothing loads them
// implicitly. Takes pointers because it mutates the members in place.
func (repository GroupMemberRepository) LoadMemberGrants(members []*models.GroupMember) error {
	if len(members) == 0 {
		return nil
	}

	db := repository.GetDB()

	// Distinct ids only: appending per member would repeat a group id once per
	// member of it (and a user id once per group they are in), so the generated IN
	// lists would grow with roster size instead of with distinct ids and a large
	// group could approach the driver's bind-parameter limit.
	groupIdSet := make(map[uint]struct{}, len(members))
	userIdSet := make(map[uint]struct{}, len(members))
	for _, member := range members {
		groupIdSet[member.GroupID] = struct{}{}
		userIdSet[member.UserID] = struct{}{}
	}
	groupIds := slices.Collect(maps.Keys(groupIdSet))
	userIds := slices.Collect(maps.Keys(userIdSet))

	var categoryGrants []models.GroupMemberCategoryGrant
	err := db.Where("group_id IN ? AND user_id IN ?", groupIds, userIds).Find(&categoryGrants).Error
	if err != nil {
		return err
	}

	var tagGrants []models.GroupMemberTagGrant
	err = db.Where("group_id IN ? AND user_id IN ?", groupIds, userIds).Find(&tagGrants).Error
	if err != nil {
		return err
	}

	// The IN x IN filter above is a superset of the requested pairs (it matches the
	// cross product), so index by the exact pair and let unrequested rows fall out.
	type memberKey struct{ userId, groupId uint }
	categoriesByMember := make(map[memberKey][]uint, len(members))
	for _, grant := range categoryGrants {
		key := memberKey{grant.UserID, grant.GroupID}
		categoriesByMember[key] = append(categoriesByMember[key], grant.CategoryID)
	}
	tagsByMember := make(map[memberKey][]uint, len(members))
	for _, grant := range tagGrants {
		key := memberKey{grant.UserID, grant.GroupID}
		tagsByMember[key] = append(tagsByMember[key], grant.TagID)
	}

	// A member with no grants must still serialize as [], not null: swagger declares
	// both as arrays, and a missing map key yields a nil slice. Normalizing here
	// covers every read path at once (AppData, GetGroupById, the paged list and the
	// group create/update responses) rather than per handler.
	for _, member := range members {
		key := memberKey{member.UserID, member.GroupID}
		member.CategoryGrants = grantIdsOrEmpty(categoriesByMember[key])
		member.TagGrants = grantIdsOrEmpty(tagsByMember[key])
	}

	return nil
}

func grantIdsOrEmpty(ids []uint) []uint {
	if ids == nil {
		return []uint{}
	}
	return ids
}
