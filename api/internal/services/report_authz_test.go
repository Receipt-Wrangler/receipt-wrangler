package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"testing"
)

// resetAuthzCaches clears the per-role permission and grant caches so a truncated
// database's reused role ids never return a prior case's cached data.
func resetAuthzCaches() {
	clearRolePermissionCacheAll()
	ClearGroupRoleGrantCacheForTests()
}

// seedAppUser creates a user with an app role granting appPerms; returns the user id.
func seedAppUser(t *testing.T, username string, appPerms []string) uint {
	t.Helper()
	db := repositories.GetDB()

	role, err := repositories.NewRoleRepository(nil).CreateAppRole("App "+username, "", appPerms)
	if err != nil {
		t.Fatalf("seed app role: %v", err)
	}
	user := models.User{Username: username, Password: "password", AppRoleID: &role.ID}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user.ID
}

// joinGroup creates a group and a group role granting groupPerms and adds userId as
// a member of it; returns the group id and group-role id.
func joinGroup(t *testing.T, userId uint, groupName string, groupPerms []string) (uint, uint) {
	t.Helper()
	db := repositories.GetDB()

	group := models.Group{Name: groupName}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole("Role "+groupName, "", groupPerms, nil, nil, nil, false)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}
	member := models.GroupMember{GroupID: group.ID, UserID: userId, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}
	return group.ID, role.ID
}

// seedTemplateInGroups creates a report template and indexes it under groupIds.
func seedTemplateInGroups(t *testing.T, name string, groupIds []uint) uint {
	t.Helper()
	db := repositories.GetDB()

	template := models.ReportTemplate{Name: name, ConfigurationVersion: 1}
	if err := db.Create(&template).Error; err != nil {
		t.Fatalf("seed template: %v", err)
	}
	for _, groupId := range groupIds {
		if err := db.Create(&models.ReportTemplateGroup{ReportTemplateID: template.ID, GroupID: groupId}).Error; err != nil {
			t.Fatalf("seed template group: %v", err)
		}
	}
	return template.ID
}

// setMatrix replaces a group role's report-template matrix and evicts caches so the
// new grants take effect on the next resolution.
func setMatrix(t *testing.T, roleId uint, matrix []commands.ReportTemplateGrantCommand) {
	t.Helper()
	if err := repositories.NewRoleRepository(nil).ReplaceGroupRoleReportTemplateGrants(roleId, matrix); err != nil {
		t.Fatalf("set matrix: %v", err)
	}
	resetAuthzCaches()
}

func assertCanAct(t *testing.T, userId, templateId uint, action string, want bool) {
	t.Helper()
	got, err := NewPermissionService(nil).CanActOnTemplate(userId, templateId, action)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if got != want {
		utils.PrintTestError(t, got, action+" allowed="+boolStr(want))
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func TestCanActOnTemplate_BaseWithGroupReadAndNoMatrix_Allowed(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "base-user", []string{permissions.AppReportsRead})
	groupId, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	assertCanAct(t, userId, templateId, "read", true)
	// No app.reports.generate → generate denied even with group read.
	assertCanAct(t, userId, templateId, "generate", false)
}

func TestCanActOnTemplate_NoGroupRead_Denied(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "no-group-read", []string{permissions.AppReportsRead})
	// Member of the group, but the group role lacks group.reports.read.
	groupId, _ := joinGroup(t, userId, "g1", []string{permissions.GroupView})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	assertCanAct(t, userId, templateId, "read", false)
}

func TestCanActOnTemplate_MatrixRestrictsPerTemplateAndAction(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "matrix-user", []string{permissions.AppReportsRead, permissions.AppReportsGenerate})
	groupId, roleId := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	templateA := seedTemplateInGroups(t, "A", []uint{groupId})
	templateB := seedTemplateInGroups(t, "B", []uint{groupId})

	// Restrict the role to template A, read only.
	setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: templateA, Permissions: []string{"read"}},
	})

	assertCanAct(t, userId, templateA, "read", true)      // listed
	assertCanAct(t, userId, templateA, "generate", false) // action not listed
	assertCanAct(t, userId, templateB, "read", false)     // template not listed
}

func TestCanActOnTemplate_AllBypassesCeilingAndMatrix(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	// readAll holder with NO group membership at all.
	userId := seedAppUser(t, "read-all", []string{permissions.AppReportsReadAll})
	otherOwner := seedAppUser(t, "owner", []string{permissions.AppReportsRead})
	groupId, _ := joinGroup(t, otherOwner, "g1", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	assertCanAct(t, userId, templateId, "read", true)      // readAll bypass
	assertCanAct(t, userId, templateId, "generate", false) // no generate / generateAll
}

func TestCanActOnTemplate_GenerateAllBypassesMatrixForOneAction(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "gen-all", []string{permissions.AppReportsRead, permissions.AppReportsGenerateAll})
	groupId, roleId := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})
	// Matrix grants read only — generate is NOT in the matrix.
	setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: templateId, Permissions: []string{"read"}},
	})

	assertCanAct(t, userId, templateId, "read", true)     // base read + matrix
	assertCanAct(t, userId, templateId, "generate", true) // generateAll bypasses matrix
	assertCanAct(t, userId, templateId, "delete", false)  // no delete / deleteAll
}

func TestCanActOnTemplate_MultiGroupMostRestrictiveWins(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "multi-user", []string{permissions.AppReportsRead})
	g1, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})     // unrestricted matrix
	g2, role2 := joinGroup(t, userId, "g2", []string{permissions.GroupReportsRead}) // will restrict
	templateId := seedTemplateInGroups(t, "A", []uint{g1, g2})

	// g2's role restricts to a DIFFERENT template, so A is denied there.
	otherTemplate := seedTemplateInGroups(t, "Other", []uint{g2})
	setMatrix(t, role2, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: otherTemplate, Permissions: []string{"read"}},
	})

	// g1 allows (unrestricted) but g2 denies → most-restrictive-wins → denied.
	assertCanAct(t, userId, templateId, "read", false)
}

func TestCanActOnTemplate_FailsClosedAfterCascadeDelete(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "fail-closed", []string{permissions.AppReportsRead})
	groupId, roleId := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	templateA := seedTemplateInGroups(t, "A", []uint{groupId})
	templateB := seedTemplateInGroups(t, "B", []uint{groupId})

	// Role is restricted to template A only.
	setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: templateA, Permissions: []string{"read"}},
	})

	// Deleting A cascades its grant row away but leaves the restricted flag set;
	// the role now has an empty matrix but stays restricted (fail closed).
	if err := repositories.GetDB().Where("id = ?", templateA).Delete(&models.ReportTemplate{}).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	resetAuthzCaches()

	// Not "unrestricted" — the role sees nothing, so even B is denied.
	assertCanAct(t, userId, templateB, "read", false)
}

func TestVisibleTemplateIds_ReadAllUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "read-all-list", []string{permissions.AppReportsReadAll})

	ids, unrestricted, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !unrestricted || ids != nil {
		utils.PrintTestError(t, ids, "unrestricted (nil ids)")
	}
}

func TestVisibleTemplateIds_ScopedAndFiltered(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "scoped-list", []string{permissions.AppReportsRead})
	groupId, roleId := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	visibleTemplate := seedTemplateInGroups(t, "Visible", []uint{groupId})
	hiddenTemplate := seedTemplateInGroups(t, "Hidden", []uint{groupId})

	// Restrict the role to only the visible template.
	setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{
		{ReportTemplateId: visibleTemplate, Permissions: []string{"read"}},
	})

	ids, unrestricted, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if unrestricted {
		utils.PrintTestError(t, unrestricted, false)
	}
	if !slices.Contains(ids, visibleTemplate) || slices.Contains(ids, hiddenTemplate) {
		utils.PrintTestError(t, ids, "only the visible template")
	}
}

func TestVisibleTemplateIds_NoReadPermissionSeesNothing(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	// No app.reports.read and no readAll.
	userId := seedAppUser(t, "no-read", []string{permissions.AppAccountRead})
	groupId, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	seedTemplateInGroups(t, "A", []uint{groupId})

	ids, unrestricted, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if unrestricted || len(ids) != 0 {
		utils.PrintTestError(t, ids, "empty (no read permission)")
	}
}

func TestAllowedActionsForTemplate_MixesBaseAndAll(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	// read (base) + deleteAll (bypass), nothing else.
	userId := seedAppUser(t, "mix-user", []string{permissions.AppReportsRead, permissions.AppReportsDeleteAll})
	groupId, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	actions, err := NewPermissionService(nil).AllowedActionsForTemplate(userId, templateId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Contains(actions, "read") || !slices.Contains(actions, "delete") {
		utils.PrintTestError(t, actions, "read + delete present")
	}
	if slices.Contains(actions, "generate") || slices.Contains(actions, "update") || slices.Contains(actions, "duplicate") {
		utils.PrintTestError(t, actions, "no generate/update/duplicate")
	}
}

func TestCanReportOverGroups_CeilingAndCreateAllBypass(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "create-user", []string{permissions.AppReportsCreate})
	groupId, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	otherGroupId, _ := joinGroup(t, userId, "g2", []string{permissions.GroupView}) // no group.reports.read

	service := NewPermissionService(nil)

	// Group the user can report over → allowed.
	ok, err := service.CanReportOverGroups(userId, []string{utils.UintToString(groupId)}, permissions.AppReportsCreateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !ok {
		utils.PrintTestError(t, ok, true)
	}

	// A group without group.reports.read → denied (no createAll).
	ok, err = service.CanReportOverGroups(userId, []string{utils.UintToString(otherGroupId)}, permissions.AppReportsCreateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ok {
		utils.PrintTestError(t, ok, false)
	}

	// createAll holder bypasses the ceiling entirely.
	adminId := seedAppUser(t, "create-all-user", []string{permissions.AppReportsCreateAll})
	ok, err = service.CanReportOverGroups(adminId, []string{utils.UintToString(otherGroupId)}, permissions.AppReportsCreateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !ok {
		utils.PrintTestError(t, ok, true)
	}
}
