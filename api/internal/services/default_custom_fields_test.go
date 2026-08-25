package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func customFieldIdsOnCommand(command commands.UpsertReceiptCommand) []uint {
	ids := make([]uint, 0, len(command.CustomFields))
	for _, customField := range command.CustomFields {
		ids = append(ids, customField.CustomFieldId)
	}
	return ids
}

// TestApplyDefaultCustomFieldsIsNoOpWhenToggleOff is the backwards-compatibility guard: a group that
// configured defaults but never opted into applying them on ingest must see server-created receipts
// unchanged.
func TestApplyDefaultCustomFieldsIsNoOpWhenToggleOff(t *testing.T) {
	settings := models.GroupReceiptSettings{
		ApplyDefaultCustomFieldsOnIngest: false,
		DefaultCustomFieldIds:            []uint{1, 2},
	}
	command := commands.UpsertReceiptCommand{}

	ApplyDefaultCustomFields(settings, &command)

	if len(command.CustomFields) != 0 {
		utils.PrintTestError(t, customFieldIdsOnCommand(command), "no custom fields")
	}
}

// TestApplyDefaultCustomFieldsAppendsMissingIds covers the happy path: every configured default lands
// on the command as an EMPTY value the user fills in later.
func TestApplyDefaultCustomFieldsAppendsMissingIds(t *testing.T) {
	settings := models.GroupReceiptSettings{
		ApplyDefaultCustomFieldsOnIngest: true,
		DefaultCustomFieldIds:            []uint{4, 7},
	}
	command := commands.UpsertReceiptCommand{}

	ApplyDefaultCustomFields(settings, &command)

	ids := customFieldIdsOnCommand(command)
	if len(ids) != 2 || ids[0] != 4 || ids[1] != 7 {
		utils.PrintTestError(t, ids, []uint{4, 7})
		return
	}

	for _, customField := range command.CustomFields {
		if customField.StringValue != nil || customField.DateValue != nil || customField.SelectValue != nil ||
			customField.CurrencyValue != nil || customField.BooleanValue != nil {
			utils.PrintTestError(t, customField, "an empty custom field value")
		}
	}
}

// TestApplyDefaultCustomFieldsDoesNotDuplicateExistingField pins the dedupe: the AI response is
// unmarshalled straight into an UpsertReceiptCommand, so a group running a custom prompt can already
// have produced a value for a defaulted field.
func TestApplyDefaultCustomFieldsDoesNotDuplicateExistingField(t *testing.T) {
	existing := "from the AI"
	settings := models.GroupReceiptSettings{
		ApplyDefaultCustomFieldsOnIngest: true,
		DefaultCustomFieldIds:            []uint{4, 7},
	}
	command := commands.UpsertReceiptCommand{
		CustomFields: []commands.UpsertCustomFieldValueCommand{
			{CustomFieldId: 4, StringValue: &existing},
		},
	}

	ApplyDefaultCustomFields(settings, &command)

	ids := customFieldIdsOnCommand(command)
	if len(ids) != 2 || ids[0] != 4 || ids[1] != 7 {
		utils.PrintTestError(t, ids, []uint{4, 7})
		return
	}
	if command.CustomFields[0].StringValue == nil || *command.CustomFields[0].StringValue != existing {
		utils.PrintTestError(t, command.CustomFields[0], "the AI's value, untouched")
	}
}

// TestApplyGroupDefaultCustomFieldsIgnoresMissingSettingsRow: GroupReceiptSettings rows are created
// lazily, so a group that has never been opened in the settings UI legitimately has none. That must
// not take down the whole ingest.
func TestApplyGroupDefaultCustomFieldsIgnoresMissingSettingsRow(t *testing.T) {
	defer repositories.TruncateTestDb()

	command := commands.UpsertReceiptCommand{}
	err := ApplyGroupDefaultCustomFields(nil, 999, &command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(command.CustomFields) != 0 {
		utils.PrintTestError(t, customFieldIdsOnCommand(command), "no custom fields")
	}
}

// TestDeleteGroupRemovesDefaultCustomFields pins the group-teardown cleanup. The receipt-settings
// delete uses Select(clause.Associations), which cannot reach these rows (DefaultCustomFieldIds is
// `gorm:"-"`, so the join is not an association at all).
func TestDeleteGroupRemovesDefaultCustomFields(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	group := models.Group{Name: "default-cf-delete"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	settingsRepository := repositories.NewGroupReceiptSettingsRepository(nil)
	if _, err := settingsRepository.CreateGroupReceiptSettings(group.ID); err != nil {
		t.Fatalf("seed group receipt settings: %v", err)
	}

	customField := models.CustomField{Name: "Deleted Group Field", Type: models.TEXT}
	if err := db.Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled:     true,
		QuickScanPaidByRequired:    true,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    true,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanDefaultStatus:     models.OPEN,
		DefaultCustomFieldIds:      &[]uint{customField.ID},
	}
	if _, err := settingsRepository.UpdateGroupReceiptSettings(utils.UintToString(group.ID), command); err != nil {
		t.Fatalf("UpdateGroupReceiptSettings: %v", err)
	}

	if err := NewGroupService(nil).DeleteGroup(utils.UintToString(group.ID), true); err != nil {
		t.Fatalf("DeleteGroup: %v", err)
	}

	var remaining int64
	err := db.Model(&models.GroupReceiptSettingsCustomField{}).
		Where("group_id = ?", group.ID).Count(&remaining).Error
	if err != nil {
		t.Fatalf("count default custom fields: %v", err)
	}
	if remaining != 0 {
		t.Errorf("deleting the group left %d default custom field rows behind", remaining)
	}
}

// TestGetAppDataHydratesGroupDefaultCustomFields covers the AppData serialization boundary: the ids
// are `gorm:"-"`, so nothing preloads them, and an empty set must arrive as [] rather than null (the
// generated Dart deserializer has no null guard and a null fails the WHOLE payload).
func TestGetAppDataHydratesGroupDefaultCustomFields(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	user, err := repositories.NewUserRepository(nil).CreateUser(commands.SignUpCommand{
		Username:    "appdata-default-cf-user",
		Password:    "Password",
		DisplayName: "AppData Default CF User",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	customField := models.CustomField{Name: "Cost Centre", Type: models.TEXT}
	if err := db.Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}

	// CreateUser seeds a personal group and an "All" group; configure the first non-all one.
	groups, err := NewGroupService(nil).GetGroupsForUser(utils.UintToString(user.ID))
	if err != nil {
		t.Fatalf("GetGroupsForUser: %v", err)
	}
	var targetGroupId uint
	for _, group := range groups {
		if !group.IsAllGroup {
			targetGroupId = group.ID
			break
		}
	}
	if targetGroupId == 0 {
		t.Fatalf("expected a non-all group for the seeded user, got %d groups", len(groups))
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled:     true,
		QuickScanPaidByRequired:    true,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    true,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanDefaultStatus:     models.OPEN,
		DefaultCustomFieldIds:      &[]uint{customField.ID},
	}
	settingsRepository := repositories.NewGroupReceiptSettingsRepository(nil)
	if _, err := settingsRepository.UpdateGroupReceiptSettings(utils.UintToString(targetGroupId), command); err != nil {
		t.Fatalf("UpdateGroupReceiptSettings: %v", err)
	}

	appData, err := GetAppData(user.ID, nil)
	if err != nil {
		t.Fatalf("GetAppData: %v", err)
	}

	sawConfiguredGroup := false
	for _, group := range appData.Groups {
		ids := group.GroupReceiptSettings.DefaultCustomFieldIds
		if ids == nil {
			utils.PrintTestError(t, group.Name+": nil", "[] rather than nil")
			continue
		}
		if group.ID != targetGroupId {
			continue
		}
		sawConfiguredGroup = true
		if len(ids) != 1 || ids[0] != customField.ID {
			utils.PrintTestError(t, ids, []uint{customField.ID})
		}
	}
	if !sawConfiguredGroup {
		utils.PrintTestError(t, appData.Groups, "the configured group in app data")
	}
}
