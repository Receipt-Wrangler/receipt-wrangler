package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
)

type SystemSettingsRepository struct {
	BaseRepository
}

func NewSystemSettingsRepository(tx *gorm.DB) SystemSettingsRepository {
	repository := SystemSettingsRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository SystemSettingsRepository) GetSystemSettings() (models.SystemSettings, error) {
	db := repository.GetDB()
	var systemSettings models.SystemSettings
	var count int64

	err := db.Model(&models.SystemSettings{}).Count(&count).Error
	if err != nil {
		return models.SystemSettings{}, err
	}

	if count == 0 {
		err = db.Model(&models.SystemSettings{}).Create(&models.SystemSettings{
			BaseModel: models.BaseModel{
				ID: 1,
			},
		}).Error
		if err != nil {
			return models.SystemSettings{}, err
		}
	}

	err = db.Model(&models.SystemSettings{}).Preload(clause.Associations).Preload("TaskQueueConfigurations").First(&systemSettings).Error
	if err != nil {
		return models.SystemSettings{}, err
	}

	// NOTE: Eventually this can get deleted. This is to fix associations not working if ID Is 0
	if systemSettings.ID == 0 {
		err = db.Model(models.SystemSettings{}).Where("id = 0").Update("id", 1).Error
		if err != nil {
			return models.SystemSettings{}, err
		}
	}

	if len(systemSettings.TaskQueueConfigurations) == 0 {
		systemSettings.TaskQueueConfigurations = models.GetAllDefaultQueueConfigurations()
	}

	return systemSettings, nil
}

func (repository SystemSettingsRepository) GetSystemReceiptProcessingSettings() (structs.SystemReceiptProcessingSettings, error) {
	systemSettings, err := repository.GetSystemSettings()
	if err != nil {
		return structs.SystemReceiptProcessingSettings{}, err
	}

	if systemSettings.ReceiptProcessingSettings.ID == 0 {
		return structs.SystemReceiptProcessingSettings{}, errors.New("ReceiptProcessingSettings do not exist")
	}

	return structs.SystemReceiptProcessingSettings{
		ReceiptProcessingSettings:         systemSettings.ReceiptProcessingSettings,
		FallbackReceiptProcessingSettings: systemSettings.FallbackReceiptProcessingSettings,
	}, nil
}

func (repository SystemSettingsRepository) UpdateSystemSettings(command commands.UpsertSystemSettingsCommand) (models.SystemSettings, error) {
	db := repository.GetDB()

	existingSettings, err := repository.GetSystemSettings()
	if err != nil {
		return models.SystemSettings{}, err
	}

	updatedSettings, err := command.ToSystemSettings(existingSettings.ID)
	if err != nil {
		return models.SystemSettings{}, err
	}

	// Select("*") below writes every column, so any field the request omitted
	// would be persisted as its zero value. Two things guard the lifetimes:
	// the columns the caller did not send are dropped from the UPDATE entirely
	// (so a concurrent update that DID set one is never clobbered by a value we
	// read before the write), and the in-memory copy is back-filled so the
	// response echoes the stored value rather than a 0.
	omittedColumns := append([]string{"TaskQueueConfigurations"}, command.OmittedLifetimeColumns()...)
	command.ApplyOmittedLifetimes(existingSettings, &updatedSettings)

	err = db.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Model(&updatedSettings).Select("*").Omit(omittedColumns...).Where("id = ?", existingSettings.ID).Updates(&updatedSettings).Error
		if txErr != nil {
			return txErr
		}

		var configCount int64
		txErr = tx.Model(&models.TaskQueueConfiguration{}).Count(&configCount).Error
		if txErr != nil {
			return txErr
		}

		if configCount == 0 {
			txErr = tx.Model(&updatedSettings).
				Where("id = ?", existingSettings.ID).
				Association("TaskQueueConfigurations").
				Replace(&updatedSettings.TaskQueueConfigurations)
			if txErr != nil {
				return txErr
			}
		} else {
			for _, config := range updatedSettings.TaskQueueConfigurations {
				txErr = tx.Model(&models.TaskQueueConfiguration{}).Where("name = ?", config.Name).Updates(&models.TaskQueueConfiguration{
					Priority: config.Priority,
				}).Error

				if txErr != nil {
					return txErr
				}
			}
		}

		return nil
	})

	if err != nil {
		return models.SystemSettings{}, nil
	}

	return updatedSettings, nil
}
