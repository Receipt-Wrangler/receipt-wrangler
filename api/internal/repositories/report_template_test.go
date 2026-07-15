package repositories

import (
	"encoding/json"
	"testing"

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
