package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

// seedUserWithAppRole creates an app role granting perms and a user assigned to
// it, returning the user id and role id.
func seedUserWithAppRole(t *testing.T, username string, perms []string) (uint, uint) {
	t.Helper()
	db := repositories.GetDB()

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateAppRole("App Role", "", perms)
	if err != nil {
		t.Fatalf("seed app role: %v", err)
	}

	user := models.User{Username: username, Password: "password", AppRoleID: &role.ID}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	return user.ID, role.ID
}

// seedMemberWithGroupRole creates a group, a group role granting perms, and a
// group member assigned to it, returning the user id, group id, and role id.
func seedMemberWithGroupRole(t *testing.T, username string, perms []string) (uint, uint, uint) {
	t.Helper()
	db := repositories.GetDB()

	group := models.Group{Name: "perm-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole("Group Role", "", perms)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: username, Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRole: models.OWNER, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}

	return user.ID, group.ID, role.ID
}

func TestHasAppPermissions(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, _ := seedUserWithAppRole(t, "app-user", []string{permissions.AppUsersRead, permissions.AppUsersCreate})
	service := NewPermissionService(nil)

	tests := []struct {
		name     string
		required []string
		want     bool
	}{
		{"single granted", []string{permissions.AppUsersRead}, true},
		{"single not granted", []string{permissions.AppUsersDelete}, false},
		{"AND all granted", []string{permissions.AppUsersRead, permissions.AppUsersCreate}, true},
		{"AND one missing", []string{permissions.AppUsersRead, permissions.AppUsersDelete}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := service.HasAppPermissions(userId, test.required...)
			if err != nil {
				t.Fatalf("HasAppPermissions: %v", err)
			}
			if got != test.want {
				t.Errorf("HasAppPermissions(%v) = %v, want %v", test.required, got, test.want)
			}
		})
	}
}

func TestHasAnyAppPermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, _ := seedUserWithAppRole(t, "app-user", []string{permissions.AppUsersRead})
	service := NewPermissionService(nil)

	got, err := service.HasAnyAppPermission(userId, permissions.AppUsersDelete, permissions.AppUsersRead)
	if err != nil {
		t.Fatalf("HasAnyAppPermission: %v", err)
	}
	if !got {
		t.Error("expected OR check to pass when one permission is granted")
	}

	got, err = service.HasAnyAppPermission(userId, permissions.AppUsersDelete, permissions.AppUsersCreate)
	if err != nil {
		t.Fatalf("HasAnyAppPermission: %v", err)
	}
	if got {
		t.Error("expected OR check to fail when no permission is granted")
	}
}

func TestHasAppPermissionsWildcardGrant(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, _ := seedUserWithAppRole(t, "wildcard-user", []string{"app.*"})
	service := NewPermissionService(nil)

	got, err := service.HasAppPermissions(userId, permissions.AppUsersRead)
	if err != nil {
		t.Fatalf("HasAppPermissions: %v", err)
	}
	if !got {
		t.Error("expected app.* wildcard grant to satisfy a concrete app permission")
	}
}

func TestHasAppPermissionsNoRoleDenies(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()
	db := repositories.GetDB()

	user := models.User{Username: "no-role-user", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	service := NewPermissionService(nil)
	got, err := service.HasAppPermissions(user.ID, permissions.AppUsersRead)
	if err != nil {
		t.Fatalf("HasAppPermissions: %v", err)
	}
	if got {
		t.Error("expected user with no app role to be denied")
	}
}

func TestHasGroupPermissions(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, groupId, _ := seedMemberWithGroupRole(t, "group-user", []string{permissions.GroupReceiptsRead})
	service := NewPermissionService(nil)

	got, err := service.HasGroupPermissions(userId, groupId, permissions.GroupReceiptsRead)
	if err != nil {
		t.Fatalf("HasGroupPermissions: %v", err)
	}
	if !got {
		t.Error("expected granted group permission to pass")
	}

	got, err = service.HasGroupPermissions(userId, groupId, permissions.GroupReceiptsDelete)
	if err != nil {
		t.Fatalf("HasGroupPermissions: %v", err)
	}
	if got {
		t.Error("expected ungranted group permission to fail")
	}
}

func TestHasAnyGroupPermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, groupId, _ := seedMemberWithGroupRole(t, "group-user", []string{permissions.GroupReceiptsRead})
	service := NewPermissionService(nil)

	got, err := service.HasAnyGroupPermission(userId, groupId, permissions.GroupReceiptsDelete, permissions.GroupReceiptsRead)
	if err != nil {
		t.Fatalf("HasAnyGroupPermission: %v", err)
	}
	if !got {
		t.Error("expected OR check to pass when one group permission is granted")
	}

	got, err = service.HasAnyGroupPermission(userId, groupId, permissions.GroupReceiptsDelete, permissions.GroupReceiptsCreate)
	if err != nil {
		t.Fatalf("HasAnyGroupPermission: %v", err)
	}
	if got {
		t.Error("expected OR check to fail when no group permission is granted")
	}
}

func TestHasGroupPermissionsNonMemberDenies(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	// Seed a member of one group, then check a different group the user is not in.
	userId, _, _ := seedMemberWithGroupRole(t, "group-user", []string{permissions.GroupReceiptsRead})
	service := NewPermissionService(nil)

	got, err := service.HasGroupPermissions(userId, 999, permissions.GroupReceiptsRead)
	if err != nil {
		t.Fatalf("HasGroupPermissions: %v", err)
	}
	if got {
		t.Error("expected non-member to be denied")
	}
}

func TestHasGroupPermissionsNoGroupRoleDenies(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()
	db := repositories.GetDB()

	group := models.Group{Name: "legacy-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	user := models.User{Username: "legacy-member", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	// Member with no GroupRoleID assigned (legacy membership mid-transition).
	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRole: models.OWNER}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}

	service := NewPermissionService(nil)
	got, err := service.HasGroupPermissions(user.ID, group.ID, permissions.GroupReceiptsRead)
	if err != nil {
		t.Fatalf("HasGroupPermissions: %v", err)
	}
	if got {
		t.Error("expected member with no group role to be denied")
	}
}

func TestPermissionScopeGuards(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, _ := seedUserWithAppRole(t, "guard-user", []string{permissions.AppUsersRead})
	service := NewPermissionService(nil)

	// Group-scoped key passed to an app check.
	if _, err := service.HasAppPermissions(userId, permissions.GroupReceiptsRead); err == nil {
		t.Error("expected error passing a group permission to an app check")
	}
	// App-scoped key passed to a group check.
	if _, err := service.HasGroupPermissions(userId, 1, permissions.AppUsersRead); err == nil {
		t.Error("expected error passing an app permission to a group check")
	}
	// Unknown permission key.
	if _, err := service.HasAppPermissions(userId, "app.not.a.real.permission"); err == nil {
		t.Error("expected error for unknown permission key")
	}
	// No required permissions.
	if _, err := service.HasAppPermissions(userId); err != ErrNoRequiredPermissions {
		t.Errorf("expected ErrNoRequiredPermissions, got %v", err)
	}
}

func TestPermissionCacheInvalidatedOnRoleUpdate(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	userId, roleId := seedUserWithAppRole(t, "cache-user", []string{permissions.AppUsersRead})
	service := NewPermissionService(nil)

	// Prime the cache with the original permission set.
	got, err := service.HasAppPermissions(userId, permissions.AppUsersRead)
	if err != nil {
		t.Fatalf("HasAppPermissions: %v", err)
	}
	if !got {
		t.Fatal("expected initial permission to be granted")
	}

	// Update the role to swap the permission; this must invalidate the cache.
	roleService := NewRoleService(nil)
	_, err = roleService.UpdateRole(roleId, commands.UpsertRoleCommand{
		Name:        "App Role",
		Scope:       permissions.ScopeApp,
		Permissions: []string{permissions.AppUsersCreate},
	})
	if err != nil {
		t.Fatalf("UpdateRole: %v", err)
	}

	got, err = service.HasAppPermissions(userId, permissions.AppUsersRead)
	if err != nil {
		t.Fatalf("HasAppPermissions: %v", err)
	}
	if got {
		t.Error("expected revoked permission to be denied after role update")
	}

	got, err = service.HasAppPermissions(userId, permissions.AppUsersCreate)
	if err != nil {
		t.Fatalf("HasAppPermissions: %v", err)
	}
	if !got {
		t.Error("expected newly granted permission to be allowed after role update")
	}
}
