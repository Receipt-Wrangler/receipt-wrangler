package repositories

import (
	"encoding/json"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
)

type ReportTemplateRepository struct {
	BaseRepository
}

func NewReportTemplateRepository(tx *gorm.DB) ReportTemplateRepository {
	repository := ReportTemplateRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

// CreateReportTemplate persists a report configuration as a template. The whole
// command is stored verbatim as the JSON Configuration blob, so the template
// round-trips back into the builder unchanged; the template's name and owner are
// pulled from the command and the requesting user.
func (repository ReportTemplateRepository) CreateReportTemplate(command commands.ReportRequestCommand, userId uint) (models.ReportTemplate, error) {
	db := repository.GetDB()

	configuration, err := json.Marshal(command)
	if err != nil {
		return models.ReportTemplate{}, err
	}

	template := models.ReportTemplate{
		BaseModel:            models.BaseModel{CreatedBy: &userId},
		Name:                 command.Name,
		Configuration:        configuration,
		ConfigurationVersion: commands.CurrentReportConfigurationVersion,
	}

	err = db.Create(&template).Error
	if err != nil {
		return models.ReportTemplate{}, err
	}

	return template, nil
}

// DeleteReportTemplateById removes the template with the given id. Deleting a
// non-existent id is not an error (RowsAffected 0), so teardown stays idempotent.
func (repository ReportTemplateRepository) DeleteReportTemplateById(id string) error {
	db := repository.GetDB()
	// Bind the id as a parameter rather than passing it as a Delete condition:
	// GORM treats a whitespace-containing string condition as raw SQL, so an id
	// like "1 OR 1=1" (a decoded URL param) would otherwise become a raw predicate.
	return db.Where("id = ?", id).Delete(&models.ReportTemplate{}).Error
}
