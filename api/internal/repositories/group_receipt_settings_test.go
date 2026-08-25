package repositories

import (
	"encoding/json"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"strings"
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

// --- Default custom fields ---------------------------------------------------

func seedDefaultCustomField(t *testing.T, name string) uint {
	t.Helper()
	customField := models.CustomField{Name: name, Type: models.TEXT}
	if err := GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}
	return customField.ID
}

// baseSettingsCommand is a command that leaves both default-custom-field pointers nil, so a test can
// set only the field it cares about.
func baseSettingsCommand() commands.UpdateGroupReceiptSettingsCommand {
	return commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled:     true,
		QuickScanPaidByRequired:    true,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    true,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanDefaultStatus:     models.OPEN,
	}
}

func sortedUints(ids []uint) []uint {
	out := append([]uint{}, ids...)
	slices.Sort(out)
	return out
}

func TestUpdateGroupReceiptSettingsRoundTripsDefaultCustomFields(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedDefaultCustomField(t, "Field A")
	fieldB := seedDefaultCustomField(t, "Field B")

	applyOnIngest := true
	command := baseSettingsCommand()
	command.DefaultCustomFieldIds = &[]uint{fieldA, fieldB}
	command.ApplyDefaultCustomFieldsOnIngest = &applyOnIngest

	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// The PUT response must carry the ids just written, not the pre-update set: the desktop writes
	// this response straight into its group state.
	if !slices.Equal(sortedUints(updated.DefaultCustomFieldIds), []uint{fieldA, fieldB}) {
		utils.PrintTestError(t, updated.DefaultCustomFieldIds, []uint{fieldA, fieldB})
	}
	if !updated.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, updated.ApplyDefaultCustomFieldsOnIngest, true)
	}

	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !slices.Equal(sortedUints(reloaded.DefaultCustomFieldIds), []uint{fieldA, fieldB}) {
		utils.PrintTestError(t, reloaded.DefaultCustomFieldIds, []uint{fieldA, fieldB})
	}
	if !reloaded.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, reloaded.ApplyDefaultCustomFieldsOnIngest, true)
	}
}

func TestUpdateGroupReceiptSettingsReplacesDefaultCustomFieldSet(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedDefaultCustomField(t, "Field A")
	fieldB := seedDefaultCustomField(t, "Field B")
	fieldC := seedDefaultCustomField(t, "Field C")

	command := baseSettingsCommand()
	command.DefaultCustomFieldIds = &[]uint{fieldA, fieldB}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	// A second save replaces the whole set rather than merging into it.
	command.DefaultCustomFieldIds = &[]uint{fieldC}
	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !slices.Equal(updated.DefaultCustomFieldIds, []uint{fieldC}) {
		utils.PrintTestError(t, updated.DefaultCustomFieldIds, []uint{fieldC})
	}

	var rowCount int64
	GetDB().Model(&models.GroupReceiptSettingsCustomField{}).Where("group_id = ?", 1).Count(&rowCount)
	if rowCount != 1 {
		utils.PrintTestError(t, rowCount, 1)
	}
}

func TestUpdateGroupReceiptSettingsClearsDefaultCustomFieldsWithEmptySlice(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedDefaultCustomField(t, "Field A")

	command := baseSettingsCommand()
	command.DefaultCustomFieldIds = &[]uint{fieldA}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	// An explicit empty array clears the set (as opposed to nil, which leaves it alone).
	command.DefaultCustomFieldIds = &[]uint{}
	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(updated.DefaultCustomFieldIds) != 0 {
		utils.PrintTestError(t, updated.DefaultCustomFieldIds, "empty set")
	}

	var rowCount int64
	GetDB().Model(&models.GroupReceiptSettingsCustomField{}).Where("group_id = ?", 1).Count(&rowCount)
	if rowCount != 0 {
		utils.PrintTestError(t, rowCount, 0)
	}
}

// TestUpdateGroupReceiptSettingsLeavesDefaultCustomFieldsUnchangedWhenNil is the load-bearing case:
// the desktop hides this whole section from an admin without app.custom-fields.read, so its payload
// omits both keys. A non-pointer field would unmarshal as false/empty and wipe the stored config.
func TestUpdateGroupReceiptSettingsLeavesDefaultCustomFieldsUnchangedWhenNil(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedDefaultCustomField(t, "Field A")

	applyOnIngest := true
	command := baseSettingsCommand()
	command.DefaultCustomFieldIds = &[]uint{fieldA}
	command.ApplyDefaultCustomFieldsOnIngest = &applyOnIngest
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	// A save that omits both keys must not touch either stored value.
	unrelated := baseSettingsCommand()
	unrelated.HideImages = true
	updated, err := repository.UpdateGroupReceiptSettings("1", unrelated)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !slices.Equal(updated.DefaultCustomFieldIds, []uint{fieldA}) {
		utils.PrintTestError(t, updated.DefaultCustomFieldIds, []uint{fieldA})
	}
	if !updated.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, updated.ApplyDefaultCustomFieldsOnIngest, true)
	}
	if !updated.HideImages {
		utils.PrintTestError(t, updated.HideImages, true)
	}

	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !slices.Equal(reloaded.DefaultCustomFieldIds, []uint{fieldA}) {
		utils.PrintTestError(t, reloaded.DefaultCustomFieldIds, []uint{fieldA})
	}
	if !reloaded.ApplyDefaultCustomFieldsOnIngest {
		utils.PrintTestError(t, reloaded.ApplyDefaultCustomFieldsOnIngest, true)
	}
}

// TestGroupReceiptSettingsMarshalsEmptyDefaultCustomFieldsAsArray guards the wire contract: the
// generated Dart deserializer has no null guard, so a null here would fail the WHOLE AppData payload
// on already-released Android builds.
func TestGroupReceiptSettingsMarshalsEmptyDefaultCustomFieldsAsArray(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	created, err := repository.CreateGroupReceiptSettings(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	for _, settings := range []models.GroupReceiptSettings{created, reloaded} {
		bytes, err := json.Marshal(settings)
		if err != nil {
			utils.PrintTestError(t, err, "no error")
			return
		}
		if !strings.Contains(string(bytes), `"defaultCustomFieldIds":[]`) {
			utils.PrintTestError(t, string(bytes), `"defaultCustomFieldIds":[]`)
		}
	}
}

// TestUpdateGroupReceiptSettingsDoesNotBlankCustomFieldName pins the regression the explicit join
// model exists to prevent: the settings update runs Select("*")...Updates, which would
// full-save-associate a many2many and blank the not-null CustomField.Name.
func TestUpdateGroupReceiptSettingsDoesNotBlankCustomFieldName(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedDefaultCustomField(t, "Field A")

	command := baseSettingsCommand()
	command.DefaultCustomFieldIds = &[]uint{fieldA}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	// Save a second time with the same set - the delete+insert must not touch custom_fields.
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	customField, err := NewCustomFieldRepository(nil).GetCustomFieldById(fieldA)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if customField.Name != "Field A" {
		utils.PrintTestError(t, customField.Name, "Field A")
	}
	if customField.Type != models.TEXT {
		utils.PrintTestError(t, customField.Type, models.TEXT)
	}
}

// TestLoadDefaultCustomFieldIdsBatchesAcrossGroups proves the loader keys on GroupId and keeps each
// group's set separate - GetGroupById can hand back a lazily created settings row whose ID is 0.
func TestLoadDefaultCustomFieldIdsBatchesAcrossGroups(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	for _, groupId := range []uint{1, 2} {
		if _, err := repository.CreateGroupReceiptSettings(groupId); err != nil {
			utils.PrintTestError(t, err, "no error")
		}
	}

	fieldA := seedDefaultCustomField(t, "Field A")

	commandOne := baseSettingsCommand()
	commandOne.DefaultCustomFieldIds = &[]uint{fieldA}
	if _, err := repository.UpdateGroupReceiptSettings("1", commandOne); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	// Group 2 keeps its empty set, and the ids are read back with a settings row whose ID is zero.
	settingsOne := models.GroupReceiptSettings{GroupId: 1}
	settingsTwo := models.GroupReceiptSettings{GroupId: 2}
	err := repository.LoadDefaultCustomFieldIds([]*models.GroupReceiptSettings{&settingsOne, &settingsTwo})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !slices.Equal(settingsOne.DefaultCustomFieldIds, []uint{fieldA}) {
		utils.PrintTestError(t, settingsOne.DefaultCustomFieldIds, []uint{fieldA})
	}
	if settingsTwo.DefaultCustomFieldIds == nil || len(settingsTwo.DefaultCustomFieldIds) != 0 {
		utils.PrintTestError(t, settingsTwo.DefaultCustomFieldIds, "empty, non-nil set")
	}
}

// TestGetGroupByIdHydratesDefaultCustomFieldIds covers the GetGroupById serialization boundary.
func TestGetGroupByIdHydratesDefaultCustomFieldIds(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedDefaultCustomField(t, "Field A")
	command := baseSettingsCommand()
	command.DefaultCustomFieldIds = &[]uint{fieldA}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	group, err := NewGroupRepository(nil).GetGroupById("1", false, false, false)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !slices.Equal(group.GroupReceiptSettings.DefaultCustomFieldIds, []uint{fieldA}) {
		utils.PrintTestError(t, group.GroupReceiptSettings.DefaultCustomFieldIds, []uint{fieldA})
	}
}
