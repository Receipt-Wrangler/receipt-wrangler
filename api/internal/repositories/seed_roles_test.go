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

func TestSeedSystemRolesReconcilesExistingRolePermissions(t *testing.T) {
	defer TruncateTestDb()

	// Pre-create the Legacy Admin system role holding a single permission, as an
	// install seeded before the admin permission set grew would have.
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

	role, ok := findRole(roles, "Legacy Admin")
	if !ok {
		utils.PrintTestError(t, "missing role Legacy Admin", "present")
		return
	}

	// The missing permissions were added: the role now holds the full admin set,
	// and the permission it already had was neither dropped nor duplicated.
	if !equalKeySet(role.Permissions, permissions.LegacyAppAdminKeys()) {
		utils.PrintTestError(t, role.Permissions, permissions.LegacyAppAdminKeys())
	}
	// Reconciled in place — the same row, not dropped and recreated.
	if role.Id != existing.ID {
		utils.PrintTestError(t, role.Id, existing.ID)
	}
	// Still a protected system role.
	if !role.IsSystem {
		utils.PrintTestError(t, role.IsSystem, true)
	}
}

func TestSeedSystemRolesReconcileIsAddOnly(t *testing.T) {
	defer TruncateTestDb()

	// Pre-create the Legacy Owner system role with the full owner set plus an
	// extra permission that is not part of it.
	const extra = "group.bogus.extra"
	ownerKeys := permissions.LegacyGroupOwnerKeys()
	seededPerms := make([]models.GroupRolePermission, 0, len(ownerKeys)+1)
	for _, permission := range ownerKeys {
		seededPerms = append(seededPerms, models.GroupRolePermission{Permission: permission})
	}
	seededPerms = append(seededPerms, models.GroupRolePermission{Permission: extra})

	existing := models.GroupRoleDefinition{
		Name:        "Legacy Owner",
		Description: "pre-existing",
		IsSystem:    true,
		Permissions: seededPerms,
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

	role, ok := findRole(roles, "Legacy Owner")
	if !ok {
		utils.PrintTestError(t, "missing role Legacy Owner", "present")
		return
	}

	// Add-only: the extra permission survives and every owner permission is present.
	want := append(append([]string(nil), ownerKeys...), extra)
	if !equalKeySet(role.Permissions, want) {
		utils.PrintTestError(t, role.Permissions, want)
	}
}

func TestSeedSystemRolesReSeedAddsNoDuplicatePermissions(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	// The second boot runs every role through the reconcile branch. It must add
	// nothing and must not trip the (roleId, permission) unique index.
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
		// equalKeySet fails if a duplicate row lengthened the set; the explicit
		// length check states the same intent directly.
		if !equalKeySet(role.Permissions, wantPerms) {
			utils.PrintTestError(t, role.Permissions, wantPerms)
		}
		if len(role.Permissions) != len(wantPerms) {
			utils.PrintTestError(t, len(role.Permissions), len(wantPerms))
		}
	}
}

func TestMissingPermissions(t *testing.T) {
	tests := []struct {
		name    string
		have    []string
		desired []string
		want    []string
	}{
		{
			name:    "returns the missing permissions in desired order",
			have:    []string{"a", "b"},
			desired: []string{"a", "b", "c", "d"},
			want:    []string{"c", "d"},
		},
		{
			name:    "returns nothing when have already covers desired",
			have:    []string{"a", "b", "c"},
			desired: []string{"a", "b"},
			want:    []string{},
		},
		{
			name:    "de-duplicates repeats within desired",
			have:    []string{},
			desired: []string{"a", "a", "b", "b"},
			want:    []string{"a", "b"},
		},
		{
			name:    "add-only: a permission only in have is never returned",
			have:    []string{"x", "y"},
			desired: []string{"a"},
			want:    []string{"a"},
		},
		{
			name:    "empty desired yields nothing",
			have:    []string{"a"},
			desired: []string{},
			want:    []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := missingPermissions(test.have, test.desired)
			if !slices.Equal(got, test.want) {
				utils.PrintTestError(t, got, test.want)
			}
		})
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
	if err := GetDB().Model(&models.AppRole{}).Where("is_default = ?", true).Count(&appDefaults).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if appDefaults != 1 {
		utils.PrintTestError(t, appDefaults, 1)
	}
	var groupDefaults int64
	if err := GetDB().Model(&models.GroupRoleDefinition{}).Where("is_default = ?", true).Count(&groupDefaults).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
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

func TestEnsureDefaultRolesErrorsWhenLegacyRoleMissing(t *testing.T) {
	defer TruncateTestDb()

	// No roles seeded, so the legacy default roles do not exist. EnsureDefaultRoles
	// must fail loudly rather than silently leaving no default (which would lock out
	// every account created afterward).
	if err := EnsureDefaultRoles(); err == nil {
		utils.PrintTestError(t, "no error", "error: legacy role not found")
	}
}
