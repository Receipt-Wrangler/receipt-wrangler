package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
)

func sampleReportCommand() commands.ReportRequestCommand {
	return commands.ReportRequestCommand{
		Name:     "Monthly Spend",
		GroupIds: []string{"1", "2"},
		Period:   commands.ReportPeriod{Preset: "this_month"},
		Detail:   commands.ReportDetail{Mode: "records"},
		Columns: []commands.ReportColumn{
			{Kind: "dimension", Name: "Name", Label: "Name", Field: "name"},
		},
		Formats: []string{"csv"},
	}
}

func TestCreateReportTemplate_PersistsNameOwnerAndConfig(t *testing.T) {
	defer TruncateTestDb()

	var userId uint = 7
	command := sampleReportCommand()

	template, err := NewReportTemplateRepository(nil).CreateReportTemplate(command, userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if template.ID == 0 {
		utils.PrintTestError(t, template.ID, "a non-zero id")
	}
	if template.Name != command.Name {
		utils.PrintTestError(t, template.Name, command.Name)
	}
	if template.CreatedBy == nil || *template.CreatedBy != userId {
		utils.PrintTestError(t, template.CreatedBy, userId)
	}
	if template.ConfigurationVersion != commands.CurrentReportConfigurationVersion {
		utils.PrintTestError(t, template.ConfigurationVersion, commands.CurrentReportConfigurationVersion)
	}

	// The whole command round-trips out of the stored JSON blob.
	var stored commands.ReportRequestCommand
	if err := json.Unmarshal(template.Configuration, &stored); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if stored.Name != command.Name {
		utils.PrintTestError(t, stored.Name, command.Name)
	}
	if len(stored.GroupIds) != 2 || stored.GroupIds[0] != "1" {
		utils.PrintTestError(t, stored.GroupIds, command.GroupIds)
	}
	if len(stored.Columns) != 1 || stored.Columns[0].Field != "name" {
		utils.PrintTestError(t, stored.Columns, command.Columns)
	}
}

func TestCreateReportTemplate_PersistedRowIsReadable(t *testing.T) {
	defer TruncateTestDb()

	var userId uint = 3
	command := sampleReportCommand()
	command.Name = "Reloadable"

	created, err := NewReportTemplateRepository(nil).CreateReportTemplate(command, userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var fetched models.ReportTemplate
	if err := GetDB().First(&fetched, created.ID).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if fetched.Name != "Reloadable" {
		utils.PrintTestError(t, fetched.Name, "Reloadable")
	}

	var stored commands.ReportRequestCommand
	if err := json.Unmarshal(fetched.Configuration, &stored); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if stored.Name != "Reloadable" {
		utils.PrintTestError(t, stored.Name, "Reloadable")
	}
}

func TestDeleteReportTemplateById_RemovesTheRow(t *testing.T) {
	defer TruncateTestDb()

	repository := NewReportTemplateRepository(nil)
	created, err := repository.CreateReportTemplate(sampleReportCommand(), 5)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := repository.DeleteReportTemplateById(fmt.Sprint(created.ID)); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var fetched models.ReportTemplate
	err = GetDB().First(&fetched, created.ID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.PrintTestError(t, err, gorm.ErrRecordNotFound)
	}

	// Deleting a non-existent id is a no-op, not an error.
	if err := repository.DeleteReportTemplateById("999999"); err != nil {
		utils.PrintTestError(t, err, nil)
	}
}

func TestDeleteReportTemplateById_BindsIdAndDoesNotDeleteOnCraftedId(t *testing.T) {
	defer TruncateTestDb()

	repository := NewReportTemplateRepository(nil)
	created, err := repository.CreateReportTemplate(sampleReportCommand(), 1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// A crafted, whitespace-containing id (as a decoded URL param could supply) must
	// be treated as a bound value, not raw SQL — so it matches nothing and the real
	// row survives, rather than deleting every template via "WHERE 1 OR 1=1".
	if err := repository.DeleteReportTemplateById("1 OR 1=1"); err != nil {
		utils.PrintTestError(t, err, nil)
	}

	var count int64
	GetDB().Model(&models.ReportTemplate{}).Where("id = ?", created.ID).Count(&count)
	if count != 1 {
		utils.PrintTestError(t, count, int64(1))
	}
}
