package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
)

type GroupReceiptSettingsRepository struct {
	BaseRepository
}

func NewGroupReceiptSettingsRepository(tx *gorm.DB) GroupReceiptSettingsRepository {
	repository := GroupReceiptSettingsRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository GroupReceiptSettingsRepository) CreateGroupReceiptSettings(groupId uint) (models.GroupReceiptSettings, error) {
	db := repository.GetDB()

	groupReceiptSettingsToCreate := models.GroupReceiptSettings{
		GroupId: groupId,
	}

	err := db.Model(models.GroupReceiptSettings{}).Create(&groupReceiptSettingsToCreate).Error
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	groupReceiptSettingsToCreate.DefaultCustomFieldIds = defaultCustomFieldIdsOrEmpty(nil)

	return groupReceiptSettingsToCreate, nil
}

func (repository GroupReceiptSettingsRepository) GetGroupReceiptSettingsByGroupId(groupId uint) (models.GroupReceiptSettings, error) {
	db := repository.GetDB()

	var groupReceiptSettings models.GroupReceiptSettings
	err := db.Model(&groupReceiptSettings).Where("group_id = ?", groupId).First(&groupReceiptSettings).Error
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	err = repository.LoadDefaultCustomFieldIds([]*models.GroupReceiptSettings{&groupReceiptSettings})
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	return groupReceiptSettings, nil
}

// LoadDefaultCustomFieldIdsForGroups populates every group's transient
// GroupReceiptSettings.DefaultCustomFieldIds across a set of groups, in one query
// total. Used at the serialization boundaries (AppData, GetGroupById) so the
// clients receive each group's configured defaults alongside the rest of its
// receipt settings.
func (repository GroupReceiptSettingsRepository) LoadDefaultCustomFieldIdsForGroups(groups []models.Group) error {
	settings := make([]*models.GroupReceiptSettings, 0, len(groups))
	for i := range groups {
		settings = append(settings, &groups[i].GroupReceiptSettings)
	}

	return repository.LoadDefaultCustomFieldIds(settings)
}

// LoadDefaultCustomFieldIds populates the transient DefaultCustomFieldIds slice on
// each settings row, in ONE query regardless of how many groups are passed.
// Callers use it at the serialization boundary; the field is `gorm:"-"` so nothing
// loads it implicitly. Takes pointers because it mutates the rows in place.
//
// This is deliberately an explicit loader rather than a GORM AfterFind hook: a hook
// would be the only one in the codebase, would add an N+1 inside
// Preload(clause.Associations), and would not even be correct —
// UpdateGroupReceiptSettings returns the in-memory struct it mutated rather than
// re-reading, so the PUT response would carry the OLD ids.
//
// Rows are keyed on the settings' GroupId, not its primary key: GetGroupById can
// hand back a lazily-created settings row whose ID is still 0 (see
// models.GroupReceiptSettingsCustomField). Ordered by custom_field_id so the
// serialized order is deterministic.
func (repository GroupReceiptSettingsRepository) LoadDefaultCustomFieldIds(settings []*models.GroupReceiptSettings) error {
	if len(settings) == 0 {
		return nil
	}

	db := repository.GetDB()

	// Distinct ids only: the same group could legitimately appear twice in a caller's
	// slice, and repeating it would grow the IN list without adding rows.
	groupIdSet := make(map[uint]struct{}, len(settings))
	for _, setting := range settings {
		groupIdSet[setting.GroupId] = struct{}{}
	}
	groupIds := make([]uint, 0, len(groupIdSet))
	for groupId := range groupIdSet {
		groupIds = append(groupIds, groupId)
	}

	var rows []models.GroupReceiptSettingsCustomField
	err := db.Where("group_id IN ?", groupIds).Order("custom_field_id").Find(&rows).Error
	if err != nil {
		return err
	}

	idsByGroup := make(map[uint][]uint, len(groupIds))
	for _, row := range rows {
		idsByGroup[row.GroupId] = append(idsByGroup[row.GroupId], row.CustomFieldId)
	}

	// A group with no defaults must still serialize as [], not null: swagger declares
	// the property as an array, a missing map key yields a nil slice, and the
	// generated Dart deserializer has no null guard — a null would fail the WHOLE
	// AppData payload on already-released Android builds.
	for _, setting := range settings {
		setting.DefaultCustomFieldIds = defaultCustomFieldIdsOrEmpty(idsByGroup[setting.GroupId])
	}

	return nil
}

func defaultCustomFieldIdsOrEmpty(ids []uint) []uint {
	if ids == nil {
		return []uint{}
	}
	return ids
}

// UpdateGroupReceiptSettings applies a settings edit. The scalar update and the
// default-custom-field join replace run in ONE transaction, so a failure on either
// leaves the stored configuration untouched. Callers pass a nil TX today; a nested
// GORM transaction degrades to a savepoint, so this is safe either way.
func (repository GroupReceiptSettingsRepository) UpdateGroupReceiptSettings(
	groupId string,
	command commands.UpdateGroupReceiptSettingsCommand,
) (models.GroupReceiptSettings, error) {
	db := repository.GetDB()

	var groupReceiptSettings models.GroupReceiptSettings

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&groupReceiptSettings).Where("group_id = ?", groupId).Preload(clause.Associations).First(&groupReceiptSettings).Error
		if err != nil {
			return err
		}

		groupReceiptSettings.HideImages = command.HideImages
		groupReceiptSettings.HideReceiptCategories = command.HideReceiptCategories
		groupReceiptSettings.HideReceiptTags = command.HideReceiptTags
		groupReceiptSettings.HideItemCategories = command.HideItemCategories
		groupReceiptSettings.HideItemTags = command.HideItemTags
		groupReceiptSettings.HideComments = command.HideComments
		groupReceiptSettings.HideShareCategories = command.HideShareCategories
		groupReceiptSettings.HideShareTags = command.HideShareTags

		groupReceiptSettings.QuickScanPaidByEnabled = command.QuickScanPaidByEnabled
		groupReceiptSettings.QuickScanPaidByRequired = command.QuickScanPaidByRequired
		groupReceiptSettings.QuickScanDefaultPaidByType = command.QuickScanDefaultPaidByType
		groupReceiptSettings.QuickScanDefaultPaidById = command.QuickScanDefaultPaidById
		groupReceiptSettings.QuickScanStatusEnabled = command.QuickScanStatusEnabled
		groupReceiptSettings.QuickScanStatusRequired = command.QuickScanStatusRequired
		groupReceiptSettings.QuickScanDefaultStatus = command.QuickScanDefaultStatus
		groupReceiptSettings.QuickScanCategoriesEnabled = command.QuickScanCategoriesEnabled
		groupReceiptSettings.QuickScanCategoriesRequired = command.QuickScanCategoriesRequired
		groupReceiptSettings.QuickScanTagsEnabled = command.QuickScanTagsEnabled
		groupReceiptSettings.QuickScanTagsRequired = command.QuickScanTagsRequired
		groupReceiptSettings.QuickScanCommentEnabled = command.QuickScanCommentEnabled
		groupReceiptSettings.QuickScanCommentRequired = command.QuickScanCommentRequired

		// Pointer field: nil means the client omitted the key, so leave the stored value alone.
		if command.ApplyDefaultCustomFieldsOnIngest != nil {
			groupReceiptSettings.ApplyDefaultCustomFieldsOnIngest = *command.ApplyDefaultCustomFieldsOnIngest
		}

		err = tx.Select("*").Model(*&groupReceiptSettings).Updates(groupReceiptSettings).Error
		if err != nil {
			return err
		}

		// Same pointer semantics: nil leaves the configured set as-is, an empty slice clears it.
		if command.DefaultCustomFieldIds != nil {
			err = replaceGroupDefaultCustomFields(tx, groupReceiptSettings.GroupId, *command.DefaultCustomFieldIds)
			if err != nil {
				return err
			}
		}

		return NewGroupReceiptSettingsRepository(tx).
			LoadDefaultCustomFieldIds([]*models.GroupReceiptSettings{&groupReceiptSettings})
	})
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	return groupReceiptSettings, nil
}

// replaceGroupDefaultCustomFields rebuilds a group's default custom field set
// (delete-all-then-insert), deduping the submitted ids so a repeated id can't
// violate the composite primary key.
//
// The insert Omits the CustomField association: without it GORM upserts a
// zero-valued CustomField whose Name is `not null`, blanking the catalog entry —
// the same hazard replaceReportTemplateGroups guards against.
func replaceGroupDefaultCustomFields(db *gorm.DB, groupId uint, customFieldIds []uint) error {
	err := db.Where("group_id = ?", groupId).Delete(&models.GroupReceiptSettingsCustomField{}).Error
	if err != nil {
		return err
	}

	if len(customFieldIds) == 0 {
		return nil
	}

	seen := make(map[uint]struct{}, len(customFieldIds))
	rows := make([]models.GroupReceiptSettingsCustomField, 0, len(customFieldIds))
	for _, customFieldId := range customFieldIds {
		if _, ok := seen[customFieldId]; ok {
			continue
		}
		seen[customFieldId] = struct{}{}
		rows = append(rows, models.GroupReceiptSettingsCustomField{
			GroupId:       groupId,
			CustomFieldId: customFieldId,
		})
	}

	return db.Omit("CustomField").Create(&rows).Error
}
