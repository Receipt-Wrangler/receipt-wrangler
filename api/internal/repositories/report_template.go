package repositories

import (
	"encoding/json"
	"errors"

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

// GetPagedReportTemplates returns a page of templates plus the unpaged total. The
// count runs before paging so it reflects the whole set. OrderBy is guarded against
// an allow-list because it is interpolated as a raw column name (not a bound value).
func (repository ReportTemplateRepository) GetPagedReportTemplates(command commands.PagedRequestCommand) ([]models.ReportTemplate, int64, error) {
	db := repository.GetDB()
	var results []models.ReportTemplate
	var count int64

	query := db.Model(&models.ReportTemplate{})

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	if !repository.isValidColumn(command.OrderBy) {
		return nil, 0, errors.New("invalid column: " + command.OrderBy)
	}

	query = repository.Sort(query, command.OrderBy, command.SortDirection)
	query = query.Scopes(repository.Paginate(command.Page, command.PageSize))

	err = query.Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// GetReportTemplateById loads a single template by id. The id is bound as a
// parameter (never interpolated as raw SQL) and gorm.ErrRecordNotFound is returned
// to the caller for a missing row.
func (repository ReportTemplateRepository) GetReportTemplateById(id string) (models.ReportTemplate, error) {
	db := repository.GetDB()
	var template models.ReportTemplate

	err := db.Where("id = ?", id).First(&template).Error
	if err != nil {
		return models.ReportTemplate{}, err
	}

	return template, nil
}

// DuplicateReportTemplate copies the source template into a new row owned by
// userId. The new template resets identity (a fresh id + owner + " duplicate" name)
// while carrying the source's configuration and version verbatim.
func (repository ReportTemplateRepository) DuplicateReportTemplate(userId uint, id string) (models.ReportTemplate, error) {
	db := repository.GetDB()

	source, err := repository.GetReportTemplateById(id)
	if err != nil {
		return models.ReportTemplate{}, err
	}

	template := models.ReportTemplate{
		BaseModel:            models.BaseModel{CreatedBy: &userId},
		Name:                 source.Name + " duplicate",
		Configuration:        source.Configuration,
		ConfigurationVersion: source.ConfigurationVersion,
	}

	err = db.Create(&template).Error
	if err != nil {
		return models.ReportTemplate{}, err
	}

	return template, nil
}

// UpdateReportTemplate overwrites an existing template's name and configuration in
// place. The row is loaded first (so a missing id surfaces gorm.ErrRecordNotFound
// rather than silently creating a new row), preserving its id and owner while the
// name/config/version are replaced and UpdatedAt is refreshed by GORM. The whole
// command is re-marshaled to the JSON blob so the template round-trips unchanged.
func (repository ReportTemplateRepository) UpdateReportTemplate(command commands.ReportRequestCommand, id string) (models.ReportTemplate, error) {
	db := repository.GetDB()

	template, err := repository.GetReportTemplateById(id)
	if err != nil {
		return models.ReportTemplate{}, err
	}

	configuration, err := json.Marshal(command)
	if err != nil {
		return models.ReportTemplate{}, err
	}

	template.Name = command.Name
	template.Configuration = configuration
	template.ConfigurationVersion = commands.CurrentReportConfigurationVersion

	err = db.Save(&template).Error
	if err != nil {
		return models.ReportTemplate{}, err
	}

	return template, nil
}

// isValidColumn allow-lists the sortable columns for GetPagedReportTemplates. The
// order-by value is interpolated as a raw column name, so it must never come
// straight from the request.
func (repository ReportTemplateRepository) isValidColumn(orderBy string) bool {
	return orderBy == "name" ||
		orderBy == "created_at" ||
		orderBy == "updated_at"
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
