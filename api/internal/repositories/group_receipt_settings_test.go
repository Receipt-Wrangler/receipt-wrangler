package repositories

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func setupGroupReceiptSettingsRepository() GroupReceiptSettingsRepository {
	return NewGroupReceiptSettingsRepository(nil)
}

func TestCreateGroupReceiptSettingsAppliesQuickScanDefaults(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	settings, err := repository.CreateGroupReceiptSettings(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	// Reload so we see the DB-applied column defaults rather than the zero-valued struct.
	settings, err = repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	if !settings.QuickScanPaidByEnabled || !settings.QuickScanPaidByRequired {
		utils.PrintTestError(t, settings, "paid by enabled + required by default")
	}
	if !settings.QuickScanStatusEnabled || !settings.QuickScanStatusRequired {
		utils.PrintTestError(t, settings, "status enabled + required by default")
	}
	if settings.QuickScanCategoriesEnabled || settings.QuickScanTagsEnabled {
		utils.PrintTestError(t, settings, "categories + tags hidden by default")
	}
	// Hidden by default so upgrading installs are unchanged until an admin opts in - this is what
	// keeps already-released mobile clients working after a server upgrade.
	if settings.QuickScanCommentEnabled || settings.QuickScanCommentRequired {
		utils.PrintTestError(t, settings, "comment hidden by default")
	}
}

func TestUpdateGroupReceiptSettingsPersistsQuickScanConfig(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	_, err := repository.CreateGroupReceiptSettings(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled:      true,
		QuickScanPaidByRequired:     false,
		QuickScanDefaultPaidByType:  models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanStatusEnabled:      false,
		QuickScanStatusRequired:     false,
		QuickScanDefaultStatus:      models.RESOLVED,
		QuickScanCategoriesEnabled:  true,
		QuickScanCategoriesRequired: true,
		QuickScanTagsEnabled:        true,
		QuickScanTagsRequired:       false,
		QuickScanCommentEnabled:     true,
		QuickScanCommentRequired:    true,
	}

	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	if updated.QuickScanPaidByRequired {
		utils.PrintTestError(t, updated.QuickScanPaidByRequired, false)
	}
	if updated.QuickScanDefaultPaidByType != models.QUICK_SCAN_PAID_BY_UPLOADER {
		utils.PrintTestError(t, updated.QuickScanDefaultPaidByType, models.QUICK_SCAN_PAID_BY_UPLOADER)
	}
	if updated.QuickScanStatusEnabled {
		utils.PrintTestError(t, updated.QuickScanStatusEnabled, false)
	}
	if updated.QuickScanDefaultStatus != models.RESOLVED {
		utils.PrintTestError(t, updated.QuickScanDefaultStatus, models.RESOLVED)
	}
	if !updated.QuickScanCategoriesEnabled || !updated.QuickScanCategoriesRequired {
		utils.PrintTestError(t, updated, "categories enabled + required")
	}
	if !updated.QuickScanTagsEnabled || updated.QuickScanTagsRequired {
		utils.PrintTestError(t, updated, "tags enabled, not required")
	}
	if !updated.QuickScanCommentEnabled || !updated.QuickScanCommentRequired {
		utils.PrintTestError(t, updated, "comment enabled + required")
	}

	// Confirm the values survive a fresh read.
	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if reloaded.QuickScanDefaultStatus != models.RESOLVED || reloaded.QuickScanPaidByRequired {
		utils.PrintTestError(t, reloaded, "persisted quick scan config")
	}
	// A field missing from UpdateGroupReceiptSettings' assignment block silently never persists, so
	// assert the comment toggles specifically survive a fresh read.
	if !reloaded.QuickScanCommentEnabled || !reloaded.QuickScanCommentRequired {
		utils.PrintTestError(t, reloaded, "persisted quick scan comment config")
	}
}
