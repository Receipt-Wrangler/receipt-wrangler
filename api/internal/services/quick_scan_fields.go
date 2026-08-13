package services

import (
	"fmt"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"unicode/utf8"
)

// ResolvedQuickScanFields holds the final per-file values a quick scan will be created with, after
// each field is checked against its group's quick-scan configuration and defaults are applied.
type ResolvedQuickScanFields struct {
	PaidByUserId uint
	Status       models.ReceiptStatus
	CategoryIds  []uint
	TagIds       []uint
	Comment      string
}

// ResolveQuickScanFields walks each uploaded file, loads its target group's receipt settings
// (cached per group), enforces the group's required quick-scan fields into a ValidatorError, and
// backfills paid-by/status defaults for fields the user left blank. A non-empty ValidatorError means
// the request should 400 without enqueuing anything.
func (service ReceiptService) ResolveQuickScanFields(command commands.QuickScanCommand, uploaderUserId uint) ([]ResolvedQuickScanFields, structs.ValidatorError, error) {
	settingsRepository := repositories.NewGroupReceiptSettingsRepository(service.TX)
	permissionService := NewPermissionService(service.TX)
	settingsCache := make(map[uint]models.GroupReceiptSettings)
	// Only the role's permission list is cached globally; the membership lookup behind each check is
	// not, so cache the resolved answer per group for multi-file scans.
	commentPermissionCache := make(map[uint]bool)
	resolved := make([]ResolvedQuickScanFields, len(command.Files))
	configErr := structs.ValidatorError{Errors: make(map[string]string)}

	for i := 0; i < len(command.Files); i++ {
		groupId := command.GroupIds[i]

		settings, cached := settingsCache[groupId]
		if !cached {
			loaded, err := settingsRepository.GetGroupReceiptSettingsByGroupId(groupId)
			if err != nil {
				return nil, structs.ValidatorError{}, err
			}
			settings = loaded
			settingsCache[groupId] = settings
		}

		paidByUserId := command.PaidByUserIds[i]
		status := command.Statuses[i]
		categoryIds := command.CategoryIdsForFile(i)
		tagIds := command.TagIdsForFile(i)
		fileKey := fmt.Sprintf("files.%d", i)

		if settings.QuickScanPaidByEnabled && settings.QuickScanPaidByRequired && paidByUserId == 0 {
			configErr.Errors[fileKey+".paidByUserId"] = "Paid by is required"
		} else if paidByUserId == 0 {
			paidByUserId = resolveQuickScanDefaultPaidBy(settings, uploaderUserId)
		}

		if settings.QuickScanStatusEnabled && settings.QuickScanStatusRequired && len(status) == 0 {
			configErr.Errors[fileKey+".status"] = "Status is required"
		} else if len(status) == 0 {
			status = settings.QuickScanDefaultStatus
		}

		if settings.QuickScanCategoriesEnabled && settings.QuickScanCategoriesRequired && len(categoryIds) == 0 {
			configErr.Errors[fileKey+".categoryIds"] = "At least one category is required"
		}

		if settings.QuickScanTagsEnabled && settings.QuickScanTagsRequired && len(tagIds) == 0 {
			configErr.Errors[fileKey+".tagIds"] = "At least one tag is required"
		}

		// group.comments.create acts as an extra AND on "shown": a caller who can't comment in the
		// target group never sees the field, is never required to fill it, and any comment they
		// submit anyway is silently dropped rather than 403'd — the comment is incidental to the
		// receipt, so refusing the whole fire-and-forget scan over it would be a worse outcome. The
		// permission is only resolved when the field is on, so the default (off) config costs no
		// extra queries.
		comment := command.CommentForFile(i)
		commentShown := settings.IsQuickScanCommentShown()
		if commentShown {
			canComment, cached := commentPermissionCache[groupId]
			if !cached {
				allowed, err := permissionService.HasGroupPermissions(uploaderUserId, groupId, permissions.GroupCommentsCreate)
				if err != nil {
					return nil, structs.ValidatorError{}, err
				}

				canComment = allowed
				commentPermissionCache[groupId] = canComment
			}

			commentShown = canComment
		}

		if !commentShown {
			comment = ""
		} else if settings.IsQuickScanCommentRequired() && len(comment) == 0 {
			configErr.Errors[fileKey+".comment"] = "Comment is required"
		} else if utf8.RuneCountInString(comment) > models.MaxCommentLength {
			// Counted in runes, not bytes: the Comment column is varchar(500), which MySQL and
			// Postgres both measure in characters, so a len() check would reject an accented or
			// non-Latin comment well inside the column's real capacity. Caught here rather than at
			// the database because the receipt is created asynchronously — an over-length comment
			// would otherwise fail inside the background task and the user would see nothing but a
			// "queued" toast.
			configErr.Errors[fileKey+".comment"] = fmt.Sprintf("Comment must be %d characters or fewer", models.MaxCommentLength)
		}

		resolved[i] = ResolvedQuickScanFields{
			PaidByUserId: paidByUserId,
			Status:       status,
			CategoryIds:  categoryIds,
			TagIds:       tagIds,
			Comment:      comment,
		}
	}

	return resolved, configErr, nil
}

// resolveQuickScanDefaultPaidBy resolves the configured default paid-by for a group: the uploader
// (the user running the quick scan) or a specific user. Returns 0 when unset (guarded upstream by
// group-settings validation).
func resolveQuickScanDefaultPaidBy(settings models.GroupReceiptSettings, uploaderUserId uint) uint {
	switch settings.QuickScanDefaultPaidByType {
	case models.QUICK_SCAN_PAID_BY_UPLOADER:
		return uploaderUserId
	case models.QUICK_SCAN_PAID_BY_USER:
		if settings.QuickScanDefaultPaidById != nil {
			return *settings.QuickScanDefaultPaidById
		}
	}

	return 0
}
