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

// TestLoadSettingsProjectionsBatchesAcrossGroups proves the loader keys on GroupId and keeps each
// group's set separate - GetGroupById can hand back a lazily created settings row whose ID is 0.
func TestLoadSettingsProjectionsBatchesAcrossGroups(t *testing.T) {
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
	err := repository.LoadSettingsProjections([]*models.GroupReceiptSettings{&settingsOne, &settingsTwo})
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

// seedCurrencyCustomField is the CURRENCY counterpart of seedDefaultCustomField. The summary
// only totals currency fields, so its tests need a field of that type.
func seedCurrencyCustomField(t *testing.T, name string) uint {
	t.Helper()
	customField := models.CustomField{Name: name, Type: models.CURRENCY}
	if err := GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed currency custom field: %v", err)
	}
	return customField.ID
}

func TestUpdateGroupReceiptSettingsRoundTripsReceiptSummaryConfig(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedCurrencyCustomField(t, "HST")
	fieldB := seedCurrencyCustomField(t, "Subtotal")

	enabled := true
	command := baseSettingsCommand()
	command.ReceiptSummaryEnabled = &enabled
	command.ReceiptSummaryCustomFieldIds = &[]uint{fieldA, fieldB}
	command.ReceiptSummaryStatuses = &[]models.ReceiptStatus{models.RESOLVED, models.OPEN}

	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// The PUT response must carry what was just written - the desktop writes it straight into
	// its group state.
	if !updated.ReceiptSummaryEnabled {
		utils.PrintTestError(t, updated.ReceiptSummaryEnabled, true)
	}
	if !slices.Equal(sortedUints(updated.ReceiptSummaryCustomFieldIds), []uint{fieldA, fieldB}) {
		utils.PrintTestError(t, updated.ReceiptSummaryCustomFieldIds, []uint{fieldA, fieldB})
	}

	// Statuses come back in models.ReceiptStatuses() order, NOT the order submitted and not
	// alphabetical: the rendered breakdown reads as a workflow.
	expectedStatuses := []models.ReceiptStatus{models.OPEN, models.RESOLVED}
	if !slices.Equal(updated.ReceiptSummaryStatuses, expectedStatuses) {
		utils.PrintTestError(t, updated.ReceiptSummaryStatuses, expectedStatuses)
	}

	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	// ReceiptSummaryEnabled is the one field here that is a plain column, so it is the one that
	// silently never persists if it is missing from UpdateGroupReceiptSettings' assignment block
	// (the write is Select("*"), which zeroes anything unassigned).
	if !reloaded.ReceiptSummaryEnabled {
		utils.PrintTestError(t, reloaded.ReceiptSummaryEnabled, true)
	}
	if !slices.Equal(sortedUints(reloaded.ReceiptSummaryCustomFieldIds), []uint{fieldA, fieldB}) {
		utils.PrintTestError(t, reloaded.ReceiptSummaryCustomFieldIds, []uint{fieldA, fieldB})
	}
	if !slices.Equal(reloaded.ReceiptSummaryStatuses, expectedStatuses) {
		utils.PrintTestError(t, reloaded.ReceiptSummaryStatuses, expectedStatuses)
	}
}

func TestUpdateGroupReceiptSettingsLeavesReceiptSummaryConfigUnchangedWhenNil(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedCurrencyCustomField(t, "HST")
	enabled := true

	seed := baseSettingsCommand()
	seed.ReceiptSummaryEnabled = &enabled
	seed.ReceiptSummaryCustomFieldIds = &[]uint{fieldA}
	seed.ReceiptSummaryStatuses = &[]models.ReceiptStatus{models.OPEN}
	if _, err := repository.UpdateGroupReceiptSettings("1", seed); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// A save that touches none of the three summary keys - what the desktop sends for an admin
	// without app.custom-fields.read, and what any other client sends. All three must survive.
	updated, err := repository.UpdateGroupReceiptSettings("1", baseSettingsCommand())
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !updated.ReceiptSummaryEnabled {
		utils.PrintTestError(t, updated.ReceiptSummaryEnabled, "still enabled")
	}
	if !slices.Equal(updated.ReceiptSummaryCustomFieldIds, []uint{fieldA}) {
		utils.PrintTestError(t, updated.ReceiptSummaryCustomFieldIds, []uint{fieldA})
	}
	if !slices.Equal(updated.ReceiptSummaryStatuses, []models.ReceiptStatus{models.OPEN}) {
		utils.PrintTestError(t, updated.ReceiptSummaryStatuses, []models.ReceiptStatus{models.OPEN})
	}
}

func TestUpdateGroupReceiptSettingsClearsReceiptSummarySetsWithEmptySlices(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedCurrencyCustomField(t, "HST")
	enabled := true

	seed := baseSettingsCommand()
	seed.ReceiptSummaryEnabled = &enabled
	seed.ReceiptSummaryCustomFieldIds = &[]uint{fieldA}
	seed.ReceiptSummaryStatuses = &[]models.ReceiptStatus{models.OPEN}
	if _, err := repository.UpdateGroupReceiptSettings("1", seed); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// An explicit empty array clears; that is the difference from omitting the key entirely.
	clear := baseSettingsCommand()
	clear.ReceiptSummaryCustomFieldIds = &[]uint{}
	clear.ReceiptSummaryStatuses = &[]models.ReceiptStatus{}

	updated, err := repository.UpdateGroupReceiptSettings("1", clear)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(updated.ReceiptSummaryCustomFieldIds) != 0 {
		utils.PrintTestError(t, updated.ReceiptSummaryCustomFieldIds, "empty")
	}
	if len(updated.ReceiptSummaryStatuses) != 0 {
		utils.PrintTestError(t, updated.ReceiptSummaryStatuses, "empty")
	}
	// Clearing the sets must not touch the master toggle.
	if !updated.ReceiptSummaryEnabled {
		utils.PrintTestError(t, updated.ReceiptSummaryEnabled, "still enabled")
	}
}

func TestUpdateGroupReceiptSettingsDedupesReceiptSummarySelections(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldA := seedCurrencyCustomField(t, "HST")

	// Both join tables use a composite primary key, so a repeated value would be a constraint
	// violation rather than a no-op if the replace helpers did not dedupe first.
	command := baseSettingsCommand()
	command.ReceiptSummaryCustomFieldIds = &[]uint{fieldA, fieldA}
	command.ReceiptSummaryStatuses = &[]models.ReceiptStatus{models.OPEN, models.OPEN}

	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !slices.Equal(updated.ReceiptSummaryCustomFieldIds, []uint{fieldA}) {
		utils.PrintTestError(t, updated.ReceiptSummaryCustomFieldIds, []uint{fieldA})
	}
	if !slices.Equal(updated.ReceiptSummaryStatuses, []models.ReceiptStatus{models.OPEN}) {
		utils.PrintTestError(t, updated.ReceiptSummaryStatuses, []models.ReceiptStatus{models.OPEN})
	}
}

// TestUpdateGroupReceiptSettingsKeepsDefaultAndSummaryFieldSetsSeparate is the regression guard for
// storing the two selections in separate tables. Were they one table with a purpose column, the
// replace helpers' unscoped `DELETE WHERE group_id = ?` would make saving either one wipe the other.
func TestUpdateGroupReceiptSettingsKeepsDefaultAndSummaryFieldSetsSeparate(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	defaultField := seedDefaultCustomField(t, "Notes")
	summaryField := seedCurrencyCustomField(t, "HST")

	seed := baseSettingsCommand()
	seed.DefaultCustomFieldIds = &[]uint{defaultField}
	seed.ReceiptSummaryCustomFieldIds = &[]uint{summaryField}
	if _, err := repository.UpdateGroupReceiptSettings("1", seed); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// Re-save the defaults alone; the summary selection must be untouched.
	defaultsOnly := baseSettingsCommand()
	defaultsOnly.DefaultCustomFieldIds = &[]uint{defaultField}
	updated, err := repository.UpdateGroupReceiptSettings("1", defaultsOnly)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !slices.Equal(updated.ReceiptSummaryCustomFieldIds, []uint{summaryField}) {
		utils.PrintTestError(t, updated.ReceiptSummaryCustomFieldIds, []uint{summaryField})
	}

	// And the other way round.
	summaryOnly := baseSettingsCommand()
	summaryOnly.ReceiptSummaryCustomFieldIds = &[]uint{summaryField}
	updated, err = repository.UpdateGroupReceiptSettings("1", summaryOnly)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !slices.Equal(updated.DefaultCustomFieldIds, []uint{defaultField}) {
		utils.PrintTestError(t, updated.DefaultCustomFieldIds, []uint{defaultField})
	}
}

// TestUpdateGroupReceiptSettingsDoesNotBlankSummaryCustomFieldName is the Omit("CustomField") guard.
// Without it GORM upserts a zero-valued CustomField through the join's association and blanks a
// `not null` catalog name - across the whole install, not just this group.
func TestUpdateGroupReceiptSettingsDoesNotBlankSummaryCustomFieldName(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	fieldId := seedCurrencyCustomField(t, "HST")

	command := baseSettingsCommand()
	command.ReceiptSummaryCustomFieldIds = &[]uint{fieldId}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	var customField models.CustomField
	if err := GetDB().First(&customField, fieldId).Error; err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if customField.Name != "HST" {
		utils.PrintTestError(t, customField.Name, "HST")
	}
	if customField.Type != models.CURRENCY {
		utils.PrintTestError(t, customField.Type, models.CURRENCY)
	}
}

// TestSettingsProjectionsSerializeAsEmptyArrays pins the [] vs null rule for all three projections.
// A null would fail the WHOLE AppData payload on an already-released Android build, which is how
// two production login outages happened.
func TestSettingsProjectionsSerializeAsEmptyArrays(t *testing.T) {
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

	// Both the create return value and a fresh read: the create path normalizes directly rather
	// than going through the loader, so they are two separate opportunities to emit null.
	for _, settings := range []models.GroupReceiptSettings{created, reloaded} {
		bytes, err := json.Marshal(settings)
		if err != nil {
			utils.PrintTestError(t, err, "no error")
			return
		}

		serialized := string(bytes)
		for _, key := range []string{"defaultCustomFieldIds", "receiptSummaryCustomFieldIds", "receiptSummaryStatuses"} {
			if !strings.Contains(serialized, `"`+key+`":[]`) {
				utils.PrintTestError(t, serialized, `"`+key+`":[]`)
			}
		}
	}
}

// TestUpdateGroupReceiptSettingsRoundTripsReceiptSummaryPosition covers the scalar half of the
// summary configuration. It is a real column rather than a join projection, so the risk it carries
// is the opposite one: the write is Select("*"), which zeroes anything the assignment block forgets.
func TestUpdateGroupReceiptSettingsRoundTripsReceiptSummaryPosition(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	created, err := repository.CreateGroupReceiptSettings(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// A group that has never touched the setting renders where it always did.
	if created.ReceiptSummaryPosition != models.RECEIPT_SUMMARY_POSITION_BOTTOM {
		utils.PrintTestError(t, created.ReceiptSummaryPosition, models.RECEIPT_SUMMARY_POSITION_BOTTOM)
	}

	top := models.RECEIPT_SUMMARY_POSITION_TOP
	command := baseSettingsCommand()
	command.ReceiptSummaryPosition = &top

	updated, err := repository.UpdateGroupReceiptSettings("1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if updated.ReceiptSummaryPosition != models.RECEIPT_SUMMARY_POSITION_TOP {
		utils.PrintTestError(t, updated.ReceiptSummaryPosition, models.RECEIPT_SUMMARY_POSITION_TOP)
	}

	// Read it back out of the database, not off the returned struct: the PUT response is hydrated
	// in memory, so only a fresh read proves the column was actually written.
	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if reloaded.ReceiptSummaryPosition != models.RECEIPT_SUMMARY_POSITION_TOP {
		utils.PrintTestError(t, reloaded.ReceiptSummaryPosition, models.RECEIPT_SUMMARY_POSITION_TOP)
	}
}

// A nil pointer means the client omitted the key. The position must then survive a save that
// changes something else entirely — the bug shape the pointer fields exist to prevent.
func TestUpdateGroupReceiptSettingsLeavesReceiptSummaryPositionUnchangedWhenNil(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	top := models.RECEIPT_SUMMARY_POSITION_TOP
	seed := baseSettingsCommand()
	seed.ReceiptSummaryPosition = &top
	if _, err := repository.UpdateGroupReceiptSettings("1", seed); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	untouched := baseSettingsCommand()
	untouched.HideImages = true
	updated, err := repository.UpdateGroupReceiptSettings("1", untouched)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if updated.ReceiptSummaryPosition != models.RECEIPT_SUMMARY_POSITION_TOP {
		utils.PrintTestError(t, updated.ReceiptSummaryPosition, "still TOP")
	}
}

// A position this build does not recognize — a member a newer release added, seen after a
// downgrade — must not take down a save that never mentioned the position.
//
// Before the omit list, Select("*") wrote every loaded field back, so the unknown value went
// through ReceiptSummaryPosition.Value(), which errors, and the whole UPDATE failed. Omitting
// the column when the command leaves it nil both fixes that and stops a concurrent writer's
// value being clobbered by one read before it landed.
func TestUpdateGroupReceiptSettingsToleratesAnUnknownStoredPosition(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroup()
	repository := setupGroupReceiptSettingsRepository()

	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// Raw SQL on purpose: every write path the app has rejects this value, which is exactly
	// why the only way it reaches the column is a build that knew about it.
	err := GetDB().Model(&models.GroupReceiptSettings{}).
		Where("group_id = ?", 1).
		UpdateColumn("receipt_summary_position", "FLOATING").Error
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	untouched := baseSettingsCommand()
	untouched.HideImages = true
	if _, err := repository.UpdateGroupReceiptSettings("1", untouched); err != nil {
		utils.PrintTestError(t, err, "no error saving an unrelated field")
		return
	}

	// Read the raw COLUMN, not GetGroupReceiptSettingsByGroupId: every read path runs
	// OrDefault(), which resolves anything but TOP to BOTTOM so an unknown value can never
	// reach a client. That normalization is on the way out; the column keeps what was stored,
	// so the setting survives the trip back up to the build that understands it.
	var stored string
	err = GetDB().Model(&models.GroupReceiptSettings{}).
		Where("group_id = ?", 1).
		Pluck("receipt_summary_position", &stored).Error
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if stored != "FLOATING" {
		utils.PrintTestError(t, stored, "FLOATING")
	}

	reloaded, err := repository.GetGroupReceiptSettingsByGroupId(1)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if !reloaded.HideImages {
		utils.PrintTestError(t, reloaded.HideImages, true)
	}
	// And the value a client sees is still one it can parse.
	if reloaded.ReceiptSummaryPosition != models.RECEIPT_SUMMARY_POSITION_BOTTOM {
		utils.PrintTestError(t, reloaded.ReceiptSummaryPosition, models.RECEIPT_SUMMARY_POSITION_BOTTOM)
	}
}
