package handlers

import (
	"mime/multipart"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// seedQuickScanGroup creates a group (id 1) with the given quick-scan receipt settings.
func seedQuickScanGroup(settings models.GroupReceiptSettings) {
	repositories.CreateTestGroup()
	repository := repositories.NewGroupReceiptSettingsRepository(nil)
	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		panic(err)
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled:      settings.QuickScanPaidByEnabled,
		QuickScanPaidByRequired:     settings.QuickScanPaidByRequired,
		QuickScanDefaultPaidByType:  settings.QuickScanDefaultPaidByType,
		QuickScanDefaultPaidById:    settings.QuickScanDefaultPaidById,
		QuickScanStatusEnabled:      settings.QuickScanStatusEnabled,
		QuickScanStatusRequired:     settings.QuickScanStatusRequired,
		QuickScanDefaultStatus:      settings.QuickScanDefaultStatus,
		QuickScanCategoriesEnabled:  settings.QuickScanCategoriesEnabled,
		QuickScanCategoriesRequired: settings.QuickScanCategoriesRequired,
		QuickScanTagsEnabled:        settings.QuickScanTagsEnabled,
		QuickScanTagsRequired:       settings.QuickScanTagsRequired,
	}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		panic(err)
	}
}

func singleFileCommand(paidBy uint, status models.ReceiptStatus, categoryIds []uint, tagIds []uint) commands.QuickScanCommand {
	return commands.QuickScanCommand{
		Files:         []multipart.File{nil},
		GroupIds:      []uint{1},
		PaidByUserIds: []uint{paidBy},
		Statuses:      []models.ReceiptStatus{status},
		CategoryIds:   [][]uint{categoryIds},
		TagIds:        [][]uint{tagIds},
	}
}

func TestResolveQuickScanFields_RequiredFieldsRejected(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedQuickScanGroup(models.GroupReceiptSettings{
		QuickScanPaidByEnabled:      true,
		QuickScanPaidByRequired:     true,
		QuickScanStatusEnabled:      true,
		QuickScanStatusRequired:     true,
		QuickScanCategoriesEnabled:  true,
		QuickScanCategoriesRequired: true,
		QuickScanTagsEnabled:        true,
		QuickScanTagsRequired:       true,
	})

	command := singleFileCommand(0, "", []uint{}, []uint{})
	_, configErr, err := resolveQuickScanFields(command, 42)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	for _, key := range []string{"files.0.paidByUserId", "files.0.status", "files.0.categoryIds", "files.0.tagIds"} {
		if _, ok := configErr.Errors[key]; !ok {
			utils.PrintTestError(t, configErr.Errors, key)
		}
	}
}

func TestResolveQuickScanFields_UploaderDefault(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedQuickScanGroup(models.GroupReceiptSettings{
		QuickScanPaidByEnabled:     true,
		QuickScanPaidByRequired:    false,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    false,
		QuickScanDefaultStatus:     models.NEEDS_ATTENTION,
	})

	command := singleFileCommand(0, "", []uint{}, []uint{})
	resolved, configErr, err := resolveQuickScanFields(command, 42)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}
	if resolved[0].PaidByUserId != 42 {
		utils.PrintTestError(t, resolved[0].PaidByUserId, 42)
	}
	if resolved[0].Status != models.NEEDS_ATTENTION {
		utils.PrintTestError(t, resolved[0].Status, models.NEEDS_ATTENTION)
	}
}

func TestResolveQuickScanFields_SpecificUserDefault(t *testing.T) {
	defer repositories.TruncateTestDb()
	defaultUser := uint(7)
	seedQuickScanGroup(models.GroupReceiptSettings{
		QuickScanPaidByEnabled:     false,
		QuickScanPaidByRequired:    false,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_USER,
		QuickScanDefaultPaidById:   &defaultUser,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    true,
	})

	command := singleFileCommand(0, models.OPEN, []uint{}, []uint{})
	resolved, configErr, err := resolveQuickScanFields(command, 42)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}
	if resolved[0].PaidByUserId != 7 {
		utils.PrintTestError(t, resolved[0].PaidByUserId, 7)
	}
}

func TestResolveQuickScanFields_ProvidedValuesKept(t *testing.T) {
	defer repositories.TruncateTestDb()
	seedQuickScanGroup(models.GroupReceiptSettings{
		QuickScanPaidByEnabled:     true,
		QuickScanPaidByRequired:    false,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    false,
		QuickScanDefaultStatus:     models.OPEN,
	})

	// User supplied a payer and status even though the fields are optional; keep them.
	command := singleFileCommand(99, models.RESOLVED, []uint{}, []uint{})
	resolved, _, err := resolveQuickScanFields(command, 42)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if resolved[0].PaidByUserId != 99 {
		utils.PrintTestError(t, resolved[0].PaidByUserId, 99)
	}
	if resolved[0].Status != models.RESOLVED {
		utils.PrintTestError(t, resolved[0].Status, models.RESOLVED)
	}
}
