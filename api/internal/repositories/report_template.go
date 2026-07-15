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
		BaseModel:     models.BaseModel{CreatedBy: &userId},
		Name:          command.Name,
		Configuration: configuration,
	}

	err = db.Create(&template).Error
	if err != nil {
		return models.ReportTemplate{}, err
	}

	return template, nil
}
