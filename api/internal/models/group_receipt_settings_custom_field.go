package models

// GroupReceiptSettingsCustomField is one custom field a group has declared as a
// DEFAULT: a field that should always be present on that group's receipts. The
// receipt form pre-adds the group's set on create (and re-applies it when the
// user switches groups), and — when GroupReceiptSettings.ApplyDefaultCustomFieldsOnIngest
// is on — the server attaches the same set to receipts it creates itself (quick
// scan, email integration). A default is a starting point, never a requirement:
// the user may remove any of them, and nothing validates their presence.
//
// Two deliberate choices:
//
// Keyed on GroupID, NOT on the GroupReceiptSettings row id. GroupRepository.GetGroupById
// lazily creates a missing settings row and DISCARDS the created record, so
// group.GroupReceiptSettings.ID is still 0 on the very call that created it —
// keying on it would silently write rows against id 0. GroupId is
// `not null;unique` on the settings model and is what every caller already holds
// (quick scan, the email handler, group delete, the settings PUT).
//
// An explicit join model rather than a GORM `many2many []CustomField` on
// GroupReceiptSettings. UpdateGroupReceiptSettings does `Preload(clause.Associations)`
// followed by `db.Select("*").Model(...).Updates(...)`, which full-save-associates
// every loaded association — that would upsert the joined CustomField rows and can
// blank a `not null` CustomField.Name (the failure api/CLAUDE.md documents for
// AI-assigned categories). Inserts therefore use `Omit("CustomField").Create(&rows)`,
// matching replaceReportTemplateGroups.
type GroupReceiptSettingsCustomField struct {
	GroupId       uint        `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	CustomFieldId uint        `gorm:"primaryKey;autoIncrement:false;index" json:"customFieldId"`
	CustomField   CustomField `gorm:"foreignKey:CustomFieldId;constraint:OnDelete:CASCADE" json:"-"`
}
