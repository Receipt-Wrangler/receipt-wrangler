package services

import (
	"errors"
	"fmt"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/repositories"

	"gorm.io/gorm"
)

var (
	// ErrMemberNotInGroup is returned when the grants endpoint targets a user who
	// is not a member of the group. Grants hang off the membership, so there is
	// nothing to write.
	ErrMemberNotInGroup = errors.New("user is not a member of the group")
)

// GrantCeilingViolation reports the submitted ids that fall outside the ceiling
// set by the member's group role. It carries the offending ids so the caller can
// tell the admin exactly what to fix rather than a bare "invalid".
//
// This is what keeps the empty-intersection state unreachable: a membership can
// never be configured with ids its role would immediately intersect away.
type GrantCeilingViolation struct {
	CategoryIds []uint
	TagIds      []uint
}

func (violation *GrantCeilingViolation) Error() string {
	return fmt.Sprintf(
		"grant selection falls outside the member's group role: categories %v, tags %v",
		violation.CategoryIds,
		violation.TagIds,
	)
}

// UpdateMemberGrants replaces one membership's per-member category/tag grants.
//
// The membership is identified by (groupId, userId) taken from the URL by the
// caller — never from the request body. Order of checks: the target must be a
// member, every id must reference a real category/tag, and every id must sit
// within the member's group-role ceiling (see resolveEffectiveGrants for why the
// two layers intersect).
func (service GroupService) UpdateMemberGrants(
	groupId uint,
	userId uint,
	command commands.UpdateGroupMemberGrantsCommand,
) error {
	categoryIds := normalizeUintSlice(command.CategoryIds)
	tagIds := normalizeUintSlice(command.TagIds)

	return repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		groupMemberRepository := repositories.NewGroupMemberRepository(tx)

		member, err := groupMemberRepository.GetMemberGrantContext(userId, groupId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMemberNotInGroup
		}
		if err != nil {
			return err
		}

		err = validateGrantsExist(tx, categoryIds, tagIds, nil, nil)
		if err != nil {
			return err
		}

		err = validateWithinRoleCeiling(tx, member.GroupRoleID, categoryIds, tagIds)
		if err != nil {
			return err
		}

		groupMemberRepository.SetTransaction(tx)
		return groupMemberRepository.ReplaceMemberGrants(userId, groupId, categoryIds, tagIds)
	})
}

// validateWithinRoleCeiling rejects any id the member's group role does not
// itself grant. A member with no group role, or a role that grants nothing for a
// resource, is unrestricted and therefore imposes no ceiling on that resource.
func validateWithinRoleCeiling(tx *gorm.DB, groupRoleId *uint, categoryIds []uint, tagIds []uint) error {
	if groupRoleId == nil {
		return nil
	}

	role, err := loadGroupRoleGrants(repositories.NewRoleRepository(tx), *groupRoleId)
	if err != nil {
		return err
	}
	if role == nil {
		return nil
	}

	violation := &GrantCeilingViolation{
		CategoryIds: idsOutsideCeiling(categoryIds, role.categoryIds),
		TagIds:      idsOutsideCeiling(tagIds, role.tagIds),
	}
	if len(violation.CategoryIds) > 0 || len(violation.TagIds) > 0 {
		return violation
	}

	return nil
}

// idsOutsideCeiling returns the submitted ids missing from the ceiling. An EMPTY
// ceiling means unrestricted (the long-standing category/tag grant rule), so it
// excludes nothing.
func idsOutsideCeiling(ids []uint, ceiling map[uint]struct{}) []uint {
	if len(ceiling) == 0 {
		return nil
	}

	outside := make([]uint, 0)
	for _, id := range ids {
		if _, ok := ceiling[id]; !ok {
			outside = append(outside, id)
		}
	}

	if len(outside) == 0 {
		return nil
	}
	return outside
}

// GetMemberGrants returns a membership's currently granted category and tag ids.
func (service GroupService) GetMemberGrants(groupId uint, userId uint) ([]uint, []uint, error) {
	groupMemberRepository := repositories.NewGroupMemberRepository(service.TX)

	categoryIds, err := groupMemberRepository.GetMemberCategoryGrantIds(userId, groupId)
	if err != nil {
		return nil, nil, err
	}

	tagIds, err := groupMemberRepository.GetMemberTagGrantIds(userId, groupId)
	if err != nil {
		return nil, nil, err
	}

	return categoryIds, tagIds, nil
}
