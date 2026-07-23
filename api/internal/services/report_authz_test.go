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
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole("Role "+groupName, "", groupPerms, nil, nil, nil, false, false)
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

// basePermFor / allPermFor / otherActionThan derive the permission strings and a
// contrasting action for the table-driven CanActOnTemplate cases below.
func basePermFor(action string) string { return reportActionPerms[action].base }
func allPermFor(action string) string  { return reportActionPerms[action].all }
func otherActionThan(action string) string {
	if action == "read" {
		return "generate"
	}
	return "read"
}

// TestCanActOnTemplate_Matrix exhaustively exercises the access decision for EVERY
// scopable action against EVERY permission state — the three axes (base app
// permission, the "*All" bypass, the per-group ceiling, the per-template matrix) and
// the multi-group / fail-closed edges. Each (action, state) is an isolated subtest.
func TestCanActOnTemplate_Matrix(t *testing.T) {
	states := []struct {
		name  string
		want  bool
		setup func(t *testing.T, action string) (userId, templateId uint)
	}{
		{
			// No report permission at all → denied regardless of group access.
			name: "no_app_perm", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{permissions.AppAccountRead})
				groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
				return userId, seedTemplateInGroups(t, "A", []uint{groupId})
			},
		},
		{
			// Base action perm, member of the group, but the role lacks group.reports.read.
			name: "base_no_ceiling", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupView})
				return userId, seedTemplateInGroups(t, "A", []uint{groupId})
			},
		},
		{
			// Base + ceiling + no matrix (unrestricted) → allowed.
			name: "base_ceiling_unrestricted", want: true,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
				return userId, seedTemplateInGroups(t, "A", []uint{groupId})
			},
		},
		{
			// Base + ceiling, but the matrix grants a DIFFERENT template.
			name: "matrix_other_template", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				groupId, roleId := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
				templateId := seedTemplateInGroups(t, "A", []uint{groupId})
				other := seedTemplateInGroups(t, "Other", []uint{groupId})
				setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{{ReportTemplateId: other, Permissions: []string{action}}})
				return userId, templateId
			},
		},
		{
			// Base + ceiling, matrix lists this template but a DIFFERENT action.
			name: "matrix_missing_action", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				groupId, roleId := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
				templateId := seedTemplateInGroups(t, "A", []uint{groupId})
				setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateId, Permissions: []string{otherActionThan(action)}}})
				return userId, templateId
			},
		},
		{
			// Base + ceiling, matrix lists this template AND this action → allowed
			// (proves the matrix can GRANT every action, not only read).
			name: "matrix_grants_action", want: true,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				groupId, roleId := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
				templateId := seedTemplateInGroups(t, "A", []uint{groupId})
				setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateId, Permissions: []string{action}}})
				return userId, templateId
			},
		},
		{
			// "<action>All" bypass, no group membership, no matrix → allowed.
			name: "all_bypass", want: true,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{allPermFor(action)})
				owner := seedAppUser(t, "owner", []string{permissions.AppReportsRead})
				groupId, _ := joinGroup(t, owner, "g", []string{permissions.GroupReportsRead})
				return userId, seedTemplateInGroups(t, "A", []uint{groupId})
			},
		},
		{
			// Base action perm but NOT a member of the covered group → ceiling denies.
			name: "non_member", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				owner := seedAppUser(t, "owner", []string{permissions.AppReportsRead})
				groupId, _ := joinGroup(t, owner, "g", []string{permissions.GroupReportsRead})
				return userId, seedTemplateInGroups(t, "A", []uint{groupId})
			},
		},
		{
			// Template covers g1+g2; group.reports.read in g1 only → most-restrictive-wins denies.
			name: "multigroup_second_no_ceiling", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				g1, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
				g2, _ := joinGroup(t, userId, "g2", []string{permissions.GroupView})
				return userId, seedTemplateInGroups(t, "A", []uint{g1, g2})
			},
		},
		{
			// Template covers g1+g2, ceiling passes in both; g1's matrix grants the
			// (template, action) but g2's restricts to a DIFFERENT template → matrix
			// intersection denies (guards against OR/first-match regressions).
			name: "multigroup_matrix_second_denies", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				g1, role1 := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
				g2, role2 := joinGroup(t, userId, "g2", []string{permissions.GroupReportsRead})
				templateId := seedTemplateInGroups(t, "A", []uint{g1, g2})
				other := seedTemplateInGroups(t, "Other", []uint{g2})
				setMatrix(t, role1, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateId, Permissions: []string{action}}})
				setMatrix(t, role2, []commands.ReportTemplateGrantCommand{{ReportTemplateId: other, Permissions: []string{action}}})
				return userId, templateId
			},
		},
		{
			// Template covers g1+g2, ceiling passes in both, and BOTH roles' matrices
			// grant the (template, action) → allowed (intersection is satisfied).
			name: "multigroup_matrix_both_grant", want: true,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				g1, role1 := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
				g2, role2 := joinGroup(t, userId, "g2", []string{permissions.GroupReportsRead})
				templateId := seedTemplateInGroups(t, "A", []uint{g1, g2})
				setMatrix(t, role1, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateId, Permissions: []string{action}}})
				setMatrix(t, role2, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateId, Permissions: []string{action}}})
				return userId, templateId
			},
		},
		{
			// Role opted into restriction, then its only granted template is deleted
			// (cascade empties the matrix, restricted flag persists) → fail closed.
			name: "fail_closed", want: false,
			setup: func(t *testing.T, action string) (uint, uint) {
				userId := seedAppUser(t, "u", []string{basePermFor(action)})
				groupId, roleId := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
				templateA := seedTemplateInGroups(t, "A", []uint{groupId})
				templateB := seedTemplateInGroups(t, "B", []uint{groupId})
				setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateA, Permissions: []string{action}}})
				if err := repositories.GetDB().Where("id = ?", templateA).Delete(&models.ReportTemplate{}).Error; err != nil {
					t.Fatalf("delete template A: %v", err)
				}
				resetAuthzCaches()
				return userId, templateB
			},
		},
	}

	for _, action := range reportScopableActions {
		for _, state := range states {
			t.Run(action+"/"+state.name, func(t *testing.T) {
				defer repositories.TruncateTestDb()
				resetAuthzCaches()
				userId, templateId := state.setup(t, action)
				assertCanAct(t, userId, templateId, action, state.want)
			})
		}
	}
}

func TestCanActOnTemplate_UnknownActionErrors(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsRead})
	if _, err := NewPermissionService(nil).CanActOnTemplate(userId, 1, "explode"); err == nil {
		utils.PrintTestError(t, nil, "an error for an unknown action")
	}
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

func TestVisibleTemplateIds_UnrestrictedReaderSeesAllReachable(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsRead})
	groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
	a := seedTemplateInGroups(t, "A", []uint{groupId})
	b := seedTemplateInGroups(t, "B", []uint{groupId})

	ids, unrestricted, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if unrestricted {
		utils.PrintTestError(t, unrestricted, false)
	}
	if !slices.Contains(ids, a) || !slices.Contains(ids, b) {
		utils.PrintTestError(t, ids, "both templates (unrestricted matrix)")
	}
}

func TestVisibleTemplateIds_CeilingFiltersUnreadableGroups(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsRead})
	g1, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	// A group the user cannot read (owned by someone else, user not a member).
	owner := seedAppUser(t, "owner", []string{permissions.AppReportsRead})
	g2, _ := joinGroup(t, owner, "g2", []string{permissions.GroupReportsRead})
	readable := seedTemplateInGroups(t, "Readable", []uint{g1})
	hidden := seedTemplateInGroups(t, "Hidden", []uint{g2})

	ids, _, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Contains(ids, readable) || slices.Contains(ids, hidden) {
		utils.PrintTestError(t, ids, "only the readable-group template")
	}
}

func TestVisibleTemplateIds_MultiGroupExcludedWhenAnyGroupUnreadable(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsRead})
	g1, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	g2, _ := joinGroup(t, userId, "g2", []string{permissions.GroupView}) // no group.reports.read
	multi := seedTemplateInGroups(t, "Multi", []uint{g1, g2})

	ids, _, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if slices.Contains(ids, multi) {
		utils.PrintTestError(t, ids, "must exclude the multi-group template (g2 unreadable)")
	}
}

func TestVisibleTemplateIds_FailClosedDropsRemaining(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsRead})
	groupId, roleId := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
	a := seedTemplateInGroups(t, "A", []uint{groupId})
	b := seedTemplateInGroups(t, "B", []uint{groupId})
	setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{{ReportTemplateId: a, Permissions: []string{"read"}}})

	// Delete A: cascade empties the matrix, the restricted flag persists (fail closed).
	if err := repositories.GetDB().Where("id = ?", a).Delete(&models.ReportTemplate{}).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	resetAuthzCaches()

	ids, _, err := NewPermissionService(nil).VisibleTemplateIds(userId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if slices.Contains(ids, b) {
		utils.PrintTestError(t, ids, "restricted-then-emptied role sees nothing (fail closed)")
	}
}

func TestAllowedActionsForTemplate_FullSetInOrder(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{
		permissions.AppReportsRead, permissions.AppReportsGenerate, permissions.AppReportsUpdate,
		permissions.AppReportsDelete, permissions.AppReportsDuplicate,
	})
	groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	actions, err := NewPermissionService(nil).AllowedActionsForTemplate(userId, templateId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Equal(actions, []string{"read", "generate", "update", "delete", "duplicate"}) {
		utils.PrintTestError(t, actions, "all five actions in reportScopableActions order")
	}
}

func TestAllowedActionsForTemplate_MatrixNarrows(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{
		permissions.AppReportsRead, permissions.AppReportsGenerate, permissions.AppReportsUpdate,
		permissions.AppReportsDelete, permissions.AppReportsDuplicate,
	})
	groupId, roleId := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})
	setMatrix(t, roleId, []commands.ReportTemplateGrantCommand{{ReportTemplateId: templateId, Permissions: []string{"read"}}})

	actions, err := NewPermissionService(nil).AllowedActionsForTemplate(userId, templateId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Equal(actions, []string{"read"}) {
		utils.PrintTestError(t, actions, "only read (matrix narrows despite holding every base perm)")
	}
}

func TestAllowedActionsForTemplate_CeilingFailLeavesOnlyAll(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	// Holds base read+generate AND deleteAll, but has no group.reports.read in the
	// template's group → only the "*All" action survives the ceiling.
	userId := seedAppUser(t, "u", []string{
		permissions.AppReportsRead, permissions.AppReportsGenerate, permissions.AppReportsDeleteAll,
	})
	groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupView})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	actions, err := NewPermissionService(nil).AllowedActionsForTemplate(userId, templateId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !slices.Equal(actions, []string{"delete"}) {
		utils.PrintTestError(t, actions, "only delete (deleteAll bypasses the failed ceiling)")
	}
}

func TestAllowedActionsForTemplate_EmptyWhenNoPerms(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppAccountRead})
	groupId, _ := joinGroup(t, userId, "g", []string{permissions.GroupReportsRead})
	templateId := seedTemplateInGroups(t, "A", []uint{groupId})

	actions, err := NewPermissionService(nil).AllowedActionsForTemplate(userId, templateId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(actions) != 0 {
		utils.PrintTestError(t, actions, "no actions (no report permissions)")
	}
}

func TestCanReportOverGroups_MultiGroupAllOrNothing(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsCreate})
	g1, _ := joinGroup(t, userId, "g1", []string{permissions.GroupReportsRead})
	g2, _ := joinGroup(t, userId, "g2", []string{permissions.GroupReportsRead})
	g3, _ := joinGroup(t, userId, "g3", []string{permissions.GroupView}) // no group.reports.read

	service := NewPermissionService(nil)

	ok, err := service.CanReportOverGroups(userId, []string{utils.UintToString(g1), utils.UintToString(g2)}, permissions.AppReportsCreateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !ok {
		utils.PrintTestError(t, ok, "true (read in both groups)")
	}

	ok, err = service.CanReportOverGroups(userId, []string{utils.UintToString(g1), utils.UintToString(g3)}, permissions.AppReportsCreateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ok {
		utils.PrintTestError(t, ok, "false (one group not readable)")
	}
}

func TestCanReportOverGroups_UpdateAllBypass(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	// Holds updateAll but is NOT a member of the group (no group.reports.read).
	userId := seedAppUser(t, "u", []string{permissions.AppReportsUpdateAll})
	owner := seedAppUser(t, "owner", []string{permissions.AppReportsRead})
	groupId, _ := joinGroup(t, owner, "g", []string{permissions.GroupReportsRead})

	ok, err := NewPermissionService(nil).CanReportOverGroups(userId, []string{utils.UintToString(groupId)}, permissions.AppReportsUpdateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !ok {
		utils.PrintTestError(t, ok, "true (updateAll bypasses the ceiling)")
	}
}

func TestCanReportOverGroups_EmptyListAllowed(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsCreate})
	ok, err := NewPermissionService(nil).CanReportOverGroups(userId, []string{}, permissions.AppReportsCreateAll)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !ok {
		utils.PrintTestError(t, ok, "true (no groups to check)")
	}
}

func TestCanReportOverGroups_InvalidGroupIdErrors(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	userId := seedAppUser(t, "u", []string{permissions.AppReportsCreate})
	if _, err := NewPermissionService(nil).CanReportOverGroups(userId, []string{"not-a-number"}, permissions.AppReportsCreateAll); err == nil {
		utils.PrintTestError(t, nil, "an error for a non-numeric group id")
	}
}
