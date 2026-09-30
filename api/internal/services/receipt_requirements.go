package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
)

// ResolveReceiptRequirements returns what userId must supply on groupId's
// receipts: the RequireReceiptComment / RequireReceiptImage flags of the group
// role they hold there, minus the group's waivers (see
// applyReceiptRequirementWaivers). A non-member, a member with no role, and the
// synthetic "All" group resolve to nothing required.
//
// The common case — a role with neither flag — costs one query; the group's
// settings and the comment permission are only read when a flag is set.
func (service ReceiptService) ResolveReceiptRequirements(userId uint, groupId uint) (structs.ReceiptRequirements, error) {
	flagsByGroup, err := repositories.NewRoleRepository(service.TX).GetMemberReceiptRequirementFlags(userId, []uint{groupId})
	if err != nil {
		return structs.ReceiptRequirements{}, err
	}
	flags, ok := flagsByGroup[groupId]
	if !ok {
		return structs.ReceiptRequirements{}, nil
	}

	// A narrow read: GetGroupReceiptSettingsByGroupId also loads the settings'
	// projections, none of which matter here. A group without a settings row reads
	// as the zero value, i.e. nothing hidden.
	var settings models.GroupReceiptSettings
	err = service.GetDB().Model(&models.GroupReceiptSettings{}).
		Select("hide_comments", "hide_images").
		Where("group_id = ?", groupId).
		Limit(1).
		Find(&settings).Error
	if err != nil {
		return structs.ReceiptRequirements{}, err
	}

	canComment := false
	if flags.RequireComment && !settings.HideComments {
		canComment, err = NewPermissionService(service.TX).HasGroupPermissions(userId, groupId, permissions.GroupCommentsCreate)
		if err != nil {
			return structs.ReceiptRequirements{}, err
		}
	}

	return applyReceiptRequirementWaivers(flags, settings, canComment), nil
}

// ResolveReceiptRequirementsForGroups is the batched form of
// ResolveReceiptRequirements for AppData: it reuses the groups' already-loaded
// GroupReceiptSettings and the caller's already-resolved group permissions, so it
// costs a single query however many groups the caller belongs to. Only groups
// where something is required are present in the result; the map is never nil.
func (service ReceiptService) ResolveReceiptRequirementsForGroups(
	userId uint,
	groups []models.Group,
	groupPermissions map[uint][]string,
) (map[uint]structs.ReceiptRequirements, error) {
	result := make(map[uint]structs.ReceiptRequirements)

	groupIds := make([]uint, 0, len(groups))
	for _, group := range groups {
		groupIds = append(groupIds, group.ID)
	}

	flagsByGroup, err := repositories.NewRoleRepository(service.TX).GetMemberReceiptRequirementFlags(userId, groupIds)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		flags, ok := flagsByGroup[group.ID]
		if !ok {
			continue
		}

		canComment := permissions.HasAll(groupPermissions[group.ID], permissions.GroupCommentsCreate)
		requirements := applyReceiptRequirementWaivers(flags, group.GroupReceiptSettings, canComment)
		if requirements.Any() {
			result[group.ID] = requirements
		}
	}

	return result, nil
}

// applyReceiptRequirementWaivers turns a role's raw flags into what is actually
// required in one group. A field the group hides is never required — the clients
// do not render it, so requiring it would make the group's receipts unsaveable.
// A comment is additionally waived for a caller without group.comments.create,
// who could never add one; this matches how quick scan treats its own comment
// field (IsQuickScanCommentShown plus the same permission). Images need no
// permission waiver: they ride the create call, which group.receipts.create
// already covers.
func applyReceiptRequirementWaivers(
	flags repositories.MemberReceiptRequirementFlags,
	settings models.GroupReceiptSettings,
	canComment bool,
) structs.ReceiptRequirements {
	return structs.ReceiptRequirements{
		CommentRequired: flags.RequireComment && !settings.HideComments && canComment,
		ImageRequired:   flags.RequireImage && !settings.HideImages,
	}
}
