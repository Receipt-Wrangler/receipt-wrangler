package repositories

import (
	"errors"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
)

type CustomFieldRepository struct {
	BaseRepository
}

func NewCustomFieldRepository(tx *gorm.DB) CustomFieldRepository {
	repository := CustomFieldRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository CustomFieldRepository) GetPagedCustomFields(
	pagedRequestCommand commands.PagedRequestCommand,
) ([]models.CustomField, int64, error) {
	db := repository.GetDB()
	var customFields []models.CustomField

	err := repository.validateOrderBy(pagedRequestCommand.OrderBy)
	if err != nil {
		return customFields, 0, err
	}

	query := repository.Sort(db, pagedRequestCommand.OrderBy, pagedRequestCommand.SortDirection)
	query = query.Scopes(repository.Paginate(pagedRequestCommand.Page, pagedRequestCommand.PageSize))

	err = query.Model(&models.CustomField{}).Preload("Options").Find(&customFields).Error
	if err != nil {
		return nil, 0, err
	}

	var count int64
	err = db.Model(&models.CustomField{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return customFields, count, nil
}

func (repository CustomFieldRepository) CreateCustomField(
	command commands.UpsertCustomFieldCommand,
	createdBy *uint,
) (models.CustomField, error) {
	db := repository.GetDB()

	options := make([]models.CustomFieldOption, 0, len(command.Options))
	for _, optionCommand := range command.Options {
		option := models.CustomFieldOption{
			CustomFieldId: optionCommand.CustomFieldId,
			Value:         optionCommand.Value,
		}
		options = append(options, option)
	}

	customFieldToCreate := models.CustomField{
		BaseModel: models.BaseModel{
			CreatedBy: createdBy,
		},
		Name:        command.Name,
		Type:        command.Type,
		Description: command.Description,
		Options:     options,
	}

	err := db.Create(&customFieldToCreate).Error
	if err != nil {
		return models.CustomField{}, err
	}

	return customFieldToCreate, nil
}

// UpdateCustomField applies an edit to a custom field's name, description, and
// select options.
//
// Two columns are deliberately never written: the custom field's type (a
// CustomFieldValue is stored in a type-specific column, so re-typing a field
// would mis-column every value already recorded against it) and an option's
// custom_field_id. Options are matched by id so they can be renamed in place,
// and an option the command omits is left alone rather than deleted, because
// CustomFieldValue.SelectValue holds an option id that a delete would orphan.
//
// The caller is expected to have validated the command with
// UpsertCustomFieldCommand.ValidateUpdate first; the custom_field_id predicate
// on the option update is a second line of defense, so an option belonging to
// another field is unwritable even if that check raced.
func (repository CustomFieldRepository) UpdateCustomField(
	id uint,
	command commands.UpsertCustomFieldCommand,
) (models.CustomField, error) {
	db := repository.GetDB()
	var updatedCustomField models.CustomField

	err := db.Transaction(func(tx *gorm.DB) error {
		// The map form is required: GORM's struct form skips zero values, so a
		// description the user cleared would silently keep its old text.
		err := tx.Model(&models.CustomField{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"name":        command.Name,
				"description": command.Description,
			}).Error
		if err != nil {
			return err
		}

		for _, optionCommand := range command.Options {
			if optionCommand.Id > 0 {
				err = tx.Model(&models.CustomFieldOption{}).
					Where("id = ? AND custom_field_id = ?", optionCommand.Id, id).
					Update("value", optionCommand.Value).Error
				if err != nil {
					return err
				}

				continue
			}

			option := models.CustomFieldOption{
				CustomFieldId: id,
				Value:         optionCommand.Value,
			}

			err = tx.Create(&option).Error
			if err != nil {
				return err
			}
		}

		updatedCustomField, err = NewCustomFieldRepository(tx).GetCustomFieldById(id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.CustomField{}, err
	}

	return updatedCustomField, nil
}

func (repository CustomFieldRepository) GetCustomFieldById(id uint) (models.CustomField, error) {
	db := repository.GetDB()
	var customField models.CustomField

	err := db.Preload("Options").First(&customField, id).Error
	if err != nil {
		return models.CustomField{}, err
	}

	return customField, nil
}

func (repository CustomFieldRepository) DeleteCustomField(id uint) error {
	db := repository.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&models.CustomFieldOption{}, "custom_field_id = ?", id).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&models.CustomFieldValue{}, "custom_field_id = ?", id).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&models.CustomField{}, id).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (repository CustomFieldRepository) validateOrderBy(orderBy string) error {
	if orderBy != "name" && orderBy != "type" && orderBy != "description" {
		return errors.New("invalid orderBy")
	}

	return nil
}
