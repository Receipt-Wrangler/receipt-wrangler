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

func pagedReportTemplateCommand(orderBy string, direction commands.SortDirection) commands.PagedRequestCommand {
	return commands.PagedRequestCommand{
		Page:          1,
		PageSize:      10,
		OrderBy:       orderBy,
		SortDirection: direction,
	}
}

func TestGetPagedReportTemplates_ReturnsRowsAndCount(t *testing.T) {
	defer TruncateTestDb()

	repository := NewReportTemplateRepository(nil)
	for _, name := range []string{"Beta", "Alpha"} {
		command := sampleReportCommand()
		command.Name = name
		if _, err := repository.CreateReportTemplate(command, 1); err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
	}

	templates, count, err := repository.GetPagedReportTemplates(pagedReportTemplateCommand("name", commands.ASCENDING))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if count != 2 {
		utils.PrintTestError(t, count, int64(2))
	}
	if len(templates) != 2 {
		utils.PrintTestError(t, len(templates), 2)
		return
	}
	// Ascending name order.
	if templates[0].Name != "Alpha" || templates[1].Name != "Beta" {
		utils.PrintTestError(t, []string{templates[0].Name, templates[1].Name}, []string{"Alpha", "Beta"})
	}
}

func TestGetPagedReportTemplates_EmptyReturnsZeroCount(t *testing.T) {
	defer TruncateTestDb()

	templates, count, err := NewReportTemplateRepository(nil).GetPagedReportTemplates(pagedReportTemplateCommand("name", commands.ASCENDING))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 0 {
		utils.PrintTestError(t, count, int64(0))
	}
	if len(templates) != 0 {
		utils.PrintTestError(t, len(templates), 0)
	}
}

func TestGetPagedReportTemplates_RejectsInvalidColumn(t *testing.T) {
	defer TruncateTestDb()

	// An order-by outside the allow-list is rejected rather than interpolated as raw SQL.
	_, _, err := NewReportTemplateRepository(nil).GetPagedReportTemplates(pagedReportTemplateCommand("configuration", commands.ASCENDING))
	if err == nil {
		utils.PrintTestError(t, nil, "an invalid column error")
	}
}

func TestGetReportTemplateById_ReturnsRow(t *testing.T) {
	defer TruncateTestDb()

	repository := NewReportTemplateRepository(nil)
	created, err := repository.CreateReportTemplate(sampleReportCommand(), 4)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	fetched, err := repository.GetReportTemplateById(fmt.Sprint(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if fetched.ID != created.ID {
		utils.PrintTestError(t, fetched.ID, created.ID)
	}
	if fetched.Name != created.Name {
		utils.PrintTestError(t, fetched.Name, created.Name)
	}
}

func TestGetReportTemplateById_ReturnsRecordNotFoundForMissingId(t *testing.T) {
	defer TruncateTestDb()

	_, err := NewReportTemplateRepository(nil).GetReportTemplateById("999999")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.PrintTestError(t, err, gorm.ErrRecordNotFound)
	}
}

func TestDuplicateReportTemplate_CopiesConfigAndStampsNewOwnerAndId(t *testing.T) {
	defer TruncateTestDb()

	repository := NewReportTemplateRepository(nil)
	command := sampleReportCommand()
	command.Name = "Quarterly Spend"
	source, err := repository.CreateReportTemplate(command, 2)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var duplicatingUser uint = 9
	duplicate, err := repository.DuplicateReportTemplate(duplicatingUser, fmt.Sprint(source.ID))
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if duplicate.ID == 0 || duplicate.ID == source.ID {
		utils.PrintTestError(t, duplicate.ID, "a new, non-zero id distinct from the source")
	}
	if duplicate.Name != source.Name+" duplicate" {
		utils.PrintTestError(t, duplicate.Name, source.Name+" duplicate")
	}
	if duplicate.ConfigurationVersion != source.ConfigurationVersion {
		utils.PrintTestError(t, duplicate.ConfigurationVersion, source.ConfigurationVersion)
	}
	if string(duplicate.Configuration) != string(source.Configuration) {
		utils.PrintTestError(t, string(duplicate.Configuration), string(source.Configuration))
	}
	if duplicate.CreatedBy == nil || *duplicate.CreatedBy != duplicatingUser {
		utils.PrintTestError(t, duplicate.CreatedBy, duplicatingUser)
	}
}

func TestDuplicateReportTemplate_ReturnsRecordNotFoundForMissingSource(t *testing.T) {
	defer TruncateTestDb()

	_, err := NewReportTemplateRepository(nil).DuplicateReportTemplate(1, "999999")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.PrintTestError(t, err, gorm.ErrRecordNotFound)
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
