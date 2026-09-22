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

	emptySettingsProjections(&groupReceiptSettingsToCreate)

	return groupReceiptSettingsToCreate, nil
}

func (repository GroupReceiptSettingsRepository) GetGroupReceiptSettingsByGroupId(groupId uint) (models.GroupReceiptSettings, error) {
	db := repository.GetDB()

	var groupReceiptSettings models.GroupReceiptSettings
	err := db.Model(&groupReceiptSettings).Where("group_id = ?", groupId).First(&groupReceiptSettings).Error
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	err = repository.LoadSettingsProjections([]*models.GroupReceiptSettings{&groupReceiptSettings})
	if err != nil {
		return models.GroupReceiptSettings{}, err
	}

	return groupReceiptSettings, nil
}

// LoadSettingsProjectionsForGroups populates every group's transient
// GroupReceiptSettings projections across a set of groups. Used at the
// serialization boundaries (AppData, GetGroupById) so the clients receive each
// group's configured defaults and summary configuration alongside the rest of its
// receipt settings.
func (repository GroupReceiptSettingsRepository) LoadSettingsProjectionsForGroups(groups []models.Group) error {
	settings := make([]*models.GroupReceiptSettings, 0, len(groups))
	for i := range groups {
		settings = append(settings, &groups[i].GroupReceiptSettings)
	}

	return repository.LoadSettingsProjections(settings)
}

// LoadSettingsProjections populates EVERY transient slice on each settings row —
// DefaultCustomFieldIds, ReceiptSummaryCustomFieldIds and ReceiptSummaryStatuses —
// in one query per projection regardless of how many groups are passed. Callers use
// it at the serialization boundary; the fields are `gorm:"-"` so nothing loads them
// implicitly. Takes pointers because it mutates the rows in place.
//
// It deliberately loads all three together rather than exposing one entry point per
// projection: every call site needs the whole settings row serialized, and a caller
// that hydrated one and missed another would emit a null where swagger promises an
// array — a failure that only shows up on an already-released mobile build.
//
// This is deliberately an explicit loader rather than a GORM AfterFind hook: a hook
// would be the only one in the codebase, would add an N+1 inside
// Preload(clause.Associations), and would not even be correct —
// UpdateGroupReceiptSettings returns the in-memory struct it mutated rather than
// re-reading, so the PUT response would carry the OLD values.
//
// Rows are keyed on the settings' GroupId, not its primary key: GetGroupById can
// hand back a lazily-created settings row whose ID is still 0 (see
// models.GroupReceiptSettingsCustomField). Each query is ordered so the serialized
// order is deterministic.
func (repository GroupReceiptSettingsRepository) LoadSettingsProjections(settings []*models.GroupReceiptSettings) error {
	if len(settings) == 0 {
		return nil
	}

	db := repository.GetDB()
	groupIds := distinctGroupIds(settings)

	defaultCustomFieldIds, err := valuesByGroup(
		db, groupIds, "custom_field_id",
		func(row models.GroupReceiptSettingsCustomField) (uint, uint) {
			return row.GroupId, row.CustomFieldId
		},
	)
	if err != nil {
		return err
	}

	summaryCustomFieldIds, err := valuesByGroup(
		db, groupIds, "custom_field_id",
		func(row models.GroupReceiptSettingsSummaryCustomField) (uint, uint) {
			return row.GroupId, row.CustomFieldId
		},
	)
	if err != nil {
		return err
	}

	summaryStatuses, err := valuesByGroup(
		db, groupIds, "status",
		func(row models.GroupReceiptSettingsSummaryStatus) (uint, models.ReceiptStatus) {
			return row.GroupId, row.Status
		},
	)
	if err != nil {
		return err
	}

	// A group with nothing configured must still serialize as [], not null: swagger
	// declares each property as an array, a missing map key yields a nil slice, and the
	// generated Dart deserializer has no null guard — a null would fail the WHOLE
	// AppData payload on already-released Android builds.
	//
	// The position is normalized for the same class of reason, one field further along: a settings
	// row written before the column existed reads "", and an empty value against a closed Dart
	// EnumClass throws and fails the whole payload. This loop is the one chokepoint every
	// GroupReceiptSettings read passes through — AppData, GetPagedGroups, CreateGroup,
	// GetGroupById — so normalizing here is what guarantees "" never reaches a client.
	for _, setting := range settings {
		setting.DefaultCustomFieldIds = sliceOrEmpty(defaultCustomFieldIds[setting.GroupId])
		setting.ReceiptSummaryCustomFieldIds = sliceOrEmpty(summaryCustomFieldIds[setting.GroupId])
		setting.ReceiptSummaryStatuses = sliceOrEmpty(canonicalStatusOrder(summaryStatuses[setting.GroupId]))
		setting.ReceiptSummaryPosition = setting.ReceiptSummaryPosition.OrDefault()
	}

	return nil
}

// distinctGroupIds collects the group ids to query for. Distinct only: the same group
// could legitimately appear twice in a caller's slice, and repeating it would grow the
// IN list without adding rows.
func distinctGroupIds(settings []*models.GroupReceiptSettings) []uint {
	groupIdSet := make(map[uint]struct{}, len(settings))
	for _, setting := range settings {
		groupIdSet[setting.GroupId] = struct{}{}
	}

	groupIds := make([]uint, 0, len(groupIdSet))
	for groupId := range groupIdSet {
		groupIds = append(groupIds, groupId)
	}

	return groupIds
}

// valuesByGroup runs ONE query over a group-keyed join table and folds the rows into a
// per-group slice, preserving the query's order. split pulls the group id and the value
// out of a row, which is all that differs between the three projections.
func valuesByGroup[Row any, Value any](
	db *gorm.DB,
	groupIds []uint,
	orderBy string,
	split func(Row) (uint, Value),
) (map[uint][]Value, error) {
	var rows []Row
	err := db.Where("group_id IN ?", groupIds).Order(orderBy).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	byGroup := make(map[uint][]Value, len(groupIds))
	for _, row := range rows {
		groupId, value := split(row)
		byGroup[groupId] = append(byGroup[groupId], value)
	}

	return byGroup, nil
}

// emptySettingsProjections marks a freshly created settings row as "nothing configured"
// rather than leaving nil slices behind. See the [] vs null rule on LoadSettingsProjections.
func emptySettingsProjections(settings *models.GroupReceiptSettings) {
	settings.DefaultCustomFieldIds = []uint{}
	settings.ReceiptSummaryCustomFieldIds = []uint{}
	settings.ReceiptSummaryStatuses = []models.ReceiptStatus{}
	settings.ReceiptSummaryPosition = settings.ReceiptSummaryPosition.OrDefault()
}

// canonicalStatusOrder sorts a group's configured statuses into models.ReceiptStatuses()
// order rather than leaving them alphabetical. The summary renders one row per configured
// status, and that list reads as a workflow (open -> needs attention -> resolved), so the
// declaration order is the meaningful one. It also keeps the rendered block stable when a
// status is added: a new value lands where the enum puts it, not wherever its name sorts.
func canonicalStatusOrder(statuses []models.ReceiptStatus) []models.ReceiptStatus {
	if len(statuses) == 0 {
		return statuses
	}

	configured := make(map[models.ReceiptStatus]struct{}, len(statuses))
	for _, status := range statuses {
		configured[status] = struct{}{}
	}

	ordered := make([]models.ReceiptStatus, 0, len(statuses))
	for _, status := range models.ReceiptStatuses() {
		receiptStatus := status.(models.ReceiptStatus)
		if _, ok := configured[receiptStatus]; ok {
			ordered = append(ordered, receiptStatus)
		}
	}

	return ordered
}

func sliceOrEmpty[T any](values []T) []T {
	if values == nil {
		return []T{}
	}
	return values
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

		// Pointer fields: nil means the client omitted the key, so leave the stored value alone.
		// The write below is Select("*"), so anything NOT assigned here is actively overwritten
		// with its zero value — a new setting missing from this block fails silently.
		if command.ApplyDefaultCustomFieldsOnIngest != nil {
			groupReceiptSettings.ApplyDefaultCustomFieldsOnIngest = *command.ApplyDefaultCustomFieldsOnIngest
		}
		if command.ReceiptSummaryEnabled != nil {
			groupReceiptSettings.ReceiptSummaryEnabled = *command.ReceiptSummaryEnabled
		}
		if command.ReceiptSummaryPosition != nil {
			groupReceiptSettings.ReceiptSummaryPosition = *command.ReceiptSummaryPosition
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

		if command.ReceiptSummaryCustomFieldIds != nil {
			err = replaceGroupSummaryCustomFields(tx, groupReceiptSettings.GroupId, *command.ReceiptSummaryCustomFieldIds)
			if err != nil {
				return err
			}
		}

		if command.ReceiptSummaryStatuses != nil {
			err = replaceGroupSummaryStatuses(tx, groupReceiptSettings.GroupId, *command.ReceiptSummaryStatuses)
			if err != nil {
				return err
			}
		}

		return NewGroupReceiptSettingsRepository(tx).
			LoadSettingsProjections([]*models.GroupReceiptSettings{&groupReceiptSettings})
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

// replaceGroupSummaryCustomFields rebuilds a group's summary custom field set, with the
// same shape and the same two hazards as replaceGroupDefaultCustomFields: dedupe the
// submitted ids so a repeat cannot violate the composite primary key, and Omit the
// CustomField association so GORM does not upsert a zero-valued CustomField and blank a
// `not null` Name in the catalog.
func replaceGroupSummaryCustomFields(db *gorm.DB, groupId uint, customFieldIds []uint) error {
	err := db.Where("group_id = ?", groupId).Delete(&models.GroupReceiptSettingsSummaryCustomField{}).Error
	if err != nil {
		return err
	}

	if len(customFieldIds) == 0 {
		return nil
	}

	seen := make(map[uint]struct{}, len(customFieldIds))
	rows := make([]models.GroupReceiptSettingsSummaryCustomField, 0, len(customFieldIds))
	for _, customFieldId := range customFieldIds {
		if _, ok := seen[customFieldId]; ok {
			continue
		}
		seen[customFieldId] = struct{}{}
		rows = append(rows, models.GroupReceiptSettingsSummaryCustomField{
			GroupId:       groupId,
			CustomFieldId: customFieldId,
		})
	}

	return db.Omit("CustomField").Create(&rows).Error
}

// replaceGroupSummaryStatuses rebuilds a group's summary status breakdown set. Statuses
// are part of the composite key, so duplicates are dropped here rather than surfacing as
// a constraint violation. This table has no associations, so no Omit is needed.
func replaceGroupSummaryStatuses(db *gorm.DB, groupId uint, statuses []models.ReceiptStatus) error {
	err := db.Where("group_id = ?", groupId).Delete(&models.GroupReceiptSettingsSummaryStatus{}).Error
	if err != nil {
		return err
	}

	if len(statuses) == 0 {
		return nil
	}

	seen := make(map[models.ReceiptStatus]struct{}, len(statuses))
	rows := make([]models.GroupReceiptSettingsSummaryStatus, 0, len(statuses))
	for _, status := range statuses {
		if _, ok := seen[status]; ok {
			continue
		}
		seen[status] = struct{}{}
		rows = append(rows, models.GroupReceiptSettingsSummaryStatus{
			GroupId: groupId,
			Status:  status,
		})
	}

	return db.Create(&rows).Error
}
