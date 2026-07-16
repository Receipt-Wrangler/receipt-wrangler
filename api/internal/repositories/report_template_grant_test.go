package repositories

import (
	"encoding/json"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

// makeTestReportTemplate inserts a report template and returns its id.
func makeTestReportTemplate(t *testing.T, name string) uint {
	template := models.ReportTemplate{
		Name:                 name,
		Configuration:        json.RawMessage(`{"name":"` + name + `","groupIds":[]}`),
		ConfigurationVersion: 1,
	}
	if err := GetDB().Create(&template).Error; err != nil {
		utils.PrintTestError(t, err, nil)
	}
	return template.ID
}

func TestReportTemplateGrantAndGroupTablesMigrate(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	templateId := makeTestReportTemplate(t, "Quarterly")
	role, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	grant := models.GroupRoleReportTemplateGrant{
		GroupRoleID:      role.ID,
		ReportTemplateID: templateId,
		Permission:       "read",
	}
	if err := GetDB().Create(&grant).Error; err != nil {
		utils.PrintTestError(t, err, "insert into group_role_report_template_grants should succeed")
	}

	mapping := models.ReportTemplateGroup{ReportTemplateID: templateId, GroupID: 42}
	if err := GetDB().Create(&mapping).Error; err != nil {
		utils.PrintTestError(t, err, "insert into report_template_groups should succeed")
	}
}

func TestDeleteReportTemplateCascadesGrantAndGroupRows(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	templateId := makeTestReportTemplate(t, "Cascade Template")
	role, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	GetDB().Create(&models.GroupRoleReportTemplateGrant{GroupRoleID: role.ID, ReportTemplateID: templateId, Permission: "read"})
	GetDB().Create(&models.GroupRoleReportTemplateGrant{GroupRoleID: role.ID, ReportTemplateID: templateId, Permission: "generate"})
	GetDB().Create(&models.ReportTemplateGroup{ReportTemplateID: templateId, GroupID: 7})

	if err := GetDB().Where("id = ?", templateId).Delete(&models.ReportTemplate{}).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var grantCount int64
	GetDB().Model(&models.GroupRoleReportTemplateGrant{}).Where("report_template_id = ?", templateId).Count(&grantCount)
	if grantCount != 0 {
		utils.PrintTestError(t, grantCount, 0)
	}

	var groupCount int64
	GetDB().Model(&models.ReportTemplateGroup{}).Where("report_template_id = ?", templateId).Count(&groupCount)
	if groupCount != 0 {
		utils.PrintTestError(t, groupCount, 0)
	}
}

func TestDeleteGroupRoleCascadesReportTemplateGrants(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	templateId := makeTestReportTemplate(t, "Role Cascade Template")
	role, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	GetDB().Create(&models.GroupRoleReportTemplateGrant{GroupRoleID: role.ID, ReportTemplateID: templateId, Permission: "read"})

	if err := repository.DeleteGroupRole(role.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var grantCount int64
	GetDB().Model(&models.GroupRoleReportTemplateGrant{}).Where("group_role_id = ?", role.ID).Count(&grantCount)
	if grantCount != 0 {
		utils.PrintTestError(t, grantCount, 0)
	}
}
