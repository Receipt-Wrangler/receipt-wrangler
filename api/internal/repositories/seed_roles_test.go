package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"slices"
	"sort"
	"testing"
)

// equalKeySet reports whether two permission-key slices match (order-independent).
func equalKeySet(a []string, b []string) bool {
	aCopy := append([]string(nil), a...)
	bCopy := append([]string(nil), b...)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	return slices.Equal(aCopy, bCopy)
}

// findRole returns the seeded role with the given name, or a zero RoleView.
func findRole(roles []structs.RoleView, name string) (structs.RoleView, bool) {
	for _, role := range roles {
		if role.Name == name {
			return role, true
		}
	}
	return structs.RoleView{}, false
}

func TestSeedSystemRolesCreatesFiveRoles(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roles) != 5 {
		utils.PrintTestError(t, len(roles), 5)
		return
	}

	expectedNames := []string{"Legacy Admin", "Legacy User", "Legacy Viewer", "Legacy Editor", "Legacy Owner"}
	for _, name := range expectedNames {
		role, ok := findRole(roles, name)
		if !ok {
			utils.PrintTestError(t, "missing role "+name, "present")
			continue
		}
		if !role.IsSystem {
			utils.PrintTestError(t, role.IsSystem, true)
		}
	}
}

func TestSeedSystemRolesIsIdempotent(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roles) != 5 {
		utils.PrintTestError(t, len(roles), 5)
	}
}

func TestSeedSystemRolesPermissionSets(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	expected := map[string][]string{
		"Legacy Admin":  permissions.LegacyAppAdminKeys(),
		"Legacy User":   permissions.LegacyAppUserKeys(),
		"Legacy Viewer": permissions.LegacyGroupViewerKeys(),
		"Legacy Editor": permissions.LegacyGroupEditorKeys(),
		"Legacy Owner":  permissions.LegacyGroupOwnerKeys(),
	}

	for name, wantPerms := range expected {
		role, ok := findRole(roles, name)
		if !ok {
			utils.PrintTestError(t, "missing role "+name, "present")
			continue
		}
		if !equalKeySet(role.Permissions, wantPerms) {
			utils.PrintTestError(t, role.Permissions, wantPerms)
		}
	}
}

func TestSeedSystemRolesScopeAndFlags(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	scopes := map[string]permissions.Scope{
		"Legacy Admin":  permissions.ScopeApp,
		"Legacy User":   permissions.ScopeApp,
		"Legacy Viewer": permissions.ScopeGroup,
		"Legacy Editor": permissions.ScopeGroup,
		"Legacy Owner":  permissions.ScopeGroup,
	}
	for name, wantScope := range scopes {
		role, ok := findRole(roles, name)
		if !ok {
			utils.PrintTestError(t, "missing role "+name, "present")
			continue
		}
		if role.Scope != wantScope {
			utils.PrintTestError(t, role.Scope, wantScope)
		}
	}

	// Group roles must not be auto-assignable yet (IsDefault stays false).
	var groupRoles []models.GroupRoleDefinition
	if err := GetDB().Find(&groupRoles).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	for _, role := range groupRoles {
		if role.IsDefault {
			utils.PrintTestError(t, role.IsDefault, false)
		}
	}
}

func TestSeedSystemRolesPreservesExisting(t *testing.T) {
	defer TruncateTestDb()

	// Pre-create a role with the same name but a different permission set.
	existing := models.AppRole{
		Name:        "Legacy Admin",
		Description: "pre-existing",
		IsSystem:    true,
		Permissions: []models.AppRolePermission{{Permission: permissions.AppUsersRead}},
	}
	if err := GetDB().Create(&existing).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(roles) != 5 {
		utils.PrintTestError(t, len(roles), 5)
	}

	// The pre-existing role must not have been overwritten with the full admin set.
	role, ok := findRole(roles, "Legacy Admin")
	if !ok {
		utils.PrintTestError(t, "missing role Legacy Admin", "present")
		return
	}
	if !equalKeySet(role.Permissions, []string{permissions.AppUsersRead}) {
		utils.PrintTestError(t, role.Permissions, []string{permissions.AppUsersRead})
	}
}

func TestEnsureDefaultRolesSetsLegacyDefaults(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	repository := NewRoleRepository(nil)
	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	appDefaults := 0
	groupDefaults := 0
	for _, role := range roles {
		if !role.IsDefault {
			continue
		}
		if role.Scope == permissions.ScopeApp {
			appDefaults++
			if role.Name != LegacyUserRoleName {
				utils.PrintTestError(t, role.Name, LegacyUserRoleName)
			}
		} else {
			groupDefaults++
			if role.Name != LegacyOwnerRoleName {
				utils.PrintTestError(t, role.Name, LegacyOwnerRoleName)
			}
		}
	}

	if appDefaults != 1 {
		utils.PrintTestError(t, appDefaults, 1)
	}
	if groupDefaults != 1 {
		utils.PrintTestError(t, groupDefaults, 1)
	}
}

func TestEnsureDefaultRolesIsIdempotent(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	var appDefaults int64
	GetDB().Model(&models.AppRole{}).Where("is_default = ?", true).Count(&appDefaults)
	if appDefaults != 1 {
		utils.PrintTestError(t, appDefaults, 1)
	}
	var groupDefaults int64
	GetDB().Model(&models.GroupRoleDefinition{}).Where("is_default = ?", true).Count(&groupDefaults)
	if groupDefaults != 1 {
		utils.PrintTestError(t, groupDefaults, 1)
	}
}

func TestEnsureDefaultRolesPreservesCustomDefault(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// An administrator picked a custom app role as the default before this runs.
	repository := NewRoleRepository(nil)
	custom, err := repository.CreateAppRole("Custom Default", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := repository.SetDefaultAppRole(custom.ID); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	roles, err := repository.GetAllRoles()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// The custom default must be preserved; Legacy User must not be promoted.
	customRole, ok := findRole(roles, "Custom Default")
	if !ok || !customRole.IsDefault {
		utils.PrintTestError(t, "Custom Default not default", "default preserved")
	}
	legacyUser, ok := findRole(roles, LegacyUserRoleName)
	if !ok || legacyUser.IsDefault {
		utils.PrintTestError(t, "Legacy User became default", "Legacy User not default")
	}
}
