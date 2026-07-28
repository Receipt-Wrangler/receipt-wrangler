package repositories

import (
	"encoding/json"
	"fmt"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"slices"
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
	role, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false, false)
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
	role, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false, false)
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
	role, err := repository.CreateGroupRole("Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false, false)
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

func TestReplaceGroupRoleReportTemplateGrantsPersistsAndReplaces(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	t1 := makeTestReportTemplate(t, "T1")
	t2 := makeTestReportTemplate(t, "T2")

	role, err := repository.CreateGroupRole("Report Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	err = repository.ReplaceGroupRoleReportTemplateGrants(role.ID, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: t1, Permissions: []string{"read", "generate"}},
		{ReportTemplateId: t2, Permissions: []string{"read"}},
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// The role opted into restriction.
	restricted, err := repository.GetGroupRoleReportTemplateGrantsRestricted(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !restricted {
		utils.PrintTestError(t, restricted, true)
	}

	// Flat rows: 2 for t1 + 1 for t2.
	grants, err := repository.GetGroupRoleReportTemplateGrants(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(grants) != 3 {
		utils.PrintTestError(t, len(grants), 3)
	}

	// Grouped read-back: 2 templates.
	loaded, err := repository.GetGroupRoleById(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	views := ReportTemplateGrantsFromRole(loaded)
	if len(views) != 2 {
		utils.PrintTestError(t, len(views), 2)
	}

	// Replace entirely with a single t2 delete grant.
	err = repository.ReplaceGroupRoleReportTemplateGrants(role.ID, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: t2, Permissions: []string{"delete"}},
	})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	grants, err = repository.GetGroupRoleReportTemplateGrants(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(grants) != 1 || grants[0].ReportTemplateID != t2 || grants[0].Permission != "delete" {
		utils.PrintTestError(t, grants, "a single {t2, delete} grant")
	}
}

func TestReplaceGroupRoleReportTemplateGrantsEmptyClearsRestricted(t *testing.T) {
	defer TruncateTestDb()
	repository := NewRoleRepository(nil)

	t1 := makeTestReportTemplate(t, "T1")
	role, err := repository.CreateGroupRole("Report Role", "", []string{permissions.GroupReportsRead}, nil, nil, nil, false, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := repository.ReplaceGroupRoleReportTemplateGrants(role.ID, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: t1, Permissions: []string{"read"}},
	}); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Clearing the matrix must drop the restricted flag (map update persists false).
	if err := repository.ReplaceGroupRoleReportTemplateGrants(role.ID, nil); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	restricted, err := repository.GetGroupRoleReportTemplateGrantsRestricted(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if restricted {
		utils.PrintTestError(t, restricted, false)
	}

	grants, err := repository.GetGroupRoleReportTemplateGrants(role.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(grants) != 0 {
		utils.PrintTestError(t, len(grants), 0)
	}
}

func TestReportTemplateCountByIds(t *testing.T) {
	defer TruncateTestDb()

	t1 := makeTestReportTemplate(t, "T1")
	t2 := makeTestReportTemplate(t, "T2")

	count, err := NewReportTemplateRepository(nil).CountByIds([]uint{t1, t2})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 2 {
		utils.PrintTestError(t, count, 2)
	}

	// A non-existent id is not counted.
	count, err = NewReportTemplateRepository(nil).CountByIds([]uint{t1, 999999})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 1 {
		utils.PrintTestError(t, count, 1)
	}
}

func TestCreateAndUpdateSyncTemplateGroupIndex(t *testing.T) {
	defer TruncateTestDb()
	repo := NewReportTemplateRepository(nil)

	created, err := repo.CreateReportTemplate(sampleReportCommand(), 1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Create wrote the index rows from the command's groupIds (["1","2"]).
	groups, err := repo.GetGroupIdsByTemplateId(created.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	slices.Sort(groups)
	if !slices.Equal(groups, []uint{1, 2}) {
		utils.PrintTestError(t, groups, []uint{1, 2})
	}

	// Update retargets the template onto group 3 only; the index tracks it.
	retarget := sampleReportCommand()
	retarget.GroupIds = []string{"3"}
	if _, err := repo.UpdateReportTemplate(retarget, fmt.Sprint(created.ID)); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	groups, err = repo.GetGroupIdsByTemplateId(created.ID)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Equal(groups, []uint{3}) {
		utils.PrintTestError(t, groups, []uint{3})
	}
}

func TestGetPagedReportTemplatesAppliesIdFilter(t *testing.T) {
	defer TruncateTestDb()
	repo := NewReportTemplateRepository(nil)

	first := sampleReportCommand()
	first.Name = "A"
	templateA, err := repo.CreateReportTemplate(first, 1)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	second := sampleReportCommand()
	second.Name = "B"
	if _, err := repo.CreateReportTemplate(second, 1); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	cmd := pagedReportTemplateCommand("name", commands.ASCENDING)

	// nil filter: unrestricted, both templates.
	_, count, err := repo.GetPagedReportTemplates(cmd, nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 2 {
		utils.PrintTestError(t, count, int64(2))
	}

	// Non-nil subset: only template A, and the count reflects the filter.
	only := []uint{templateA.ID}
	filtered, filteredCount, err := repo.GetPagedReportTemplates(cmd, &only)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if filteredCount != 1 || len(filtered) != 1 || filtered[0].ID != templateA.ID {
		utils.PrintTestError(t, filtered, "only template A")
	}

	// Empty set: nothing visible.
	empty := []uint{}
	none, noneCount, err := repo.GetPagedReportTemplates(cmd, &empty)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if noneCount != 0 || len(none) != 0 {
		utils.PrintTestError(t, none, "no templates")
	}
}
