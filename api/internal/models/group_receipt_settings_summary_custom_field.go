package models

// GroupReceiptSettingsSummaryCustomField is one CURRENCY custom field a group has
// declared as a SUMMARY column: a field whose values are totalled in the block of
// figures under the receipts table, alongside the receipt count and the amount
// total. Only CURRENCY fields are eligible — they are the only ones whose value
// can be summed — and the settings handler rejects any other type.
//
// This is a separate table from GroupReceiptSettingsCustomField rather than that
// table plus a purpose column. That table's primary key is
// {GroupId, CustomFieldId}, so a discriminator would have to join the key: a
// migration on a live table that changes the meaning of every query already
// written against it, to save one small file. The two sets are also independent —
// a group's defaults and its summary columns have no reason to overlap.
//
// Keyed on GroupID, NOT on the GroupReceiptSettings row id, for the same reason
// GroupReceiptSettingsCustomField is: GroupRepository.GetGroupById lazily creates a
// missing settings row and DISCARDS the created record, so
// group.GroupReceiptSettings.ID is still 0 on the very call that created it.
//
// Inserts must Omit("CustomField"), or GORM upserts a zero-valued CustomField whose
// Name is `not null` and blanks the catalog entry — see replaceGroupSummaryCustomFields.
type GroupReceiptSettingsSummaryCustomField struct {
	GroupId       uint        `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	CustomFieldId uint        `gorm:"primaryKey;autoIncrement:false;index" json:"customFieldId"`
	CustomField   CustomField `gorm:"foreignKey:CustomFieldId;constraint:OnDelete:CASCADE" json:"-"`
}
