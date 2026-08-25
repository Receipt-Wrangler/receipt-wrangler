package services

import (
	"errors"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

// ApplyDefaultCustomFields attaches a group's configured default custom fields to a receipt the
// SERVER is about to create, as EMPTY values the user fills in later. It is a no-op unless the group
// opted in with ApplyDefaultCustomFieldsOnIngest — receipts created before an admin turns it on are
// unchanged, which is what keeps the upgrade backwards compatible.
//
// Ids the command already carries are skipped. The AI response is unmarshalled straight into an
// UpsertReceiptCommand and a group running a custom prompt can emit custom fields, so a naive append
// would produce two values for the same field.
//
// Pure by design: it reads the settings it is handed and mutates the command in place, so both
// ingest paths (quick scan and the email handler) share exactly one implementation and cannot drift.
func ApplyDefaultCustomFields(settings models.GroupReceiptSettings, command *commands.UpsertReceiptCommand) {
	if command == nil || !settings.ApplyDefaultCustomFieldsOnIngest {
		return
	}

	present := make(map[uint]struct{}, len(command.CustomFields))
	for _, customField := range command.CustomFields {
		present[customField.CustomFieldId] = struct{}{}
	}

	for _, customFieldId := range settings.DefaultCustomFieldIds {
		if _, ok := present[customFieldId]; ok {
			continue
		}
		present[customFieldId] = struct{}{}

		command.CustomFields = append(command.CustomFields, commands.UpsertCustomFieldValueCommand{
			CustomFieldId: customFieldId,
		})
	}
}

// ApplyGroupDefaultCustomFields loads a group's receipt settings (with their default custom field
// ids hydrated) and delegates to ApplyDefaultCustomFields. A group with no settings row yet is
// treated as "not opted in" rather than an error: the row is created lazily, so a group that has
// never been opened in the settings UI legitimately has none, and failing here would take down the
// whole ingest.
func ApplyGroupDefaultCustomFields(tx *gorm.DB, groupId uint, command *commands.UpsertReceiptCommand) error {
	settings, err := repositories.NewGroupReceiptSettingsRepository(tx).GetGroupReceiptSettingsByGroupId(groupId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	ApplyDefaultCustomFields(settings, command)

	return nil
}
