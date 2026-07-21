package repositories

import (
	"errors"
	"fmt"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"

	"gorm.io/gorm"
)

// SeedSystemRoles ensures the five immutable, legacy-equivalent system roles
// exist and carry their full permission sets. The roles reproduce the legacy
// UserRole/GroupRole capabilities exactly so that upgrading installs see no
// behavior change.
//
// It is idempotent and self-reconciling: keyed on the role Name (a uniqueIndex
// on both AppRole and GroupRoleDefinition), a missing role is created and an
// existing one has any permissions it lacks added — add-only, so a permission
// already on the role is never removed. This is what lets a permission added to
// the registry later flow into an already-seeded Legacy Admin / Legacy Owner on
// the next boot. It is safe to run on every startup, and a partially-seeded or
// partially-reconciled database self-heals on the next boot. The roles are
// seeded only — they are not assigned to any user or group member here.
func SeedSystemRoles() error {
	db := GetDB()

	appRoles := []struct {
		name        string
		description string
		permissions []string
	}{
		{LegacyAdminRoleName, "Legacy administrator equivalent: full application access.", permissions.LegacyAppAdminKeys()},
		{LegacyUserRoleName, "Legacy standard user equivalent.", permissions.LegacyAppUserKeys()},
	}
	for _, role := range appRoles {
		if err := seedAppRole(db, role.name, role.description, role.permissions); err != nil {
			return err
		}
	}

	groupRoles := []struct {
		name        string
		description string
		permissions []string
	}{
		{LegacyViewerRoleName, "Legacy group viewer equivalent.", permissions.LegacyGroupViewerKeys()},
		{LegacyEditorRoleName, "Legacy group editor equivalent.", permissions.LegacyGroupEditorKeys()},
		{LegacyOwnerRoleName, "Legacy group owner equivalent: full group access.", permissions.LegacyGroupOwnerKeys()},
	}
	for _, role := range groupRoles {
		if err := seedGroupRole(db, role.name, role.description, role.permissions); err != nil {
			return err
		}
	}

	return nil
}

// seedAppRole creates a system app role with the given permissions when no app
// role with that name exists, or reconciles an existing one by adding any
// permissions it is missing (add-only).
func seedAppRole(db *gorm.DB, name string, description string, perms []string) error {
	var existing models.AppRole
	err := db.Preload("Permissions").Where("name = ?", name).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return createAppRole(db, name, description, perms)
	}
	if err != nil {
		return err
	}

	have := make([]string, 0, len(existing.Permissions))
	for _, permission := range existing.Permissions {
		have = append(have, permission.Permission)
	}

	missing := missingPermissions(have, perms)
	if len(missing) == 0 {
		return nil
	}

	rows := make([]models.AppRolePermission, 0, len(missing))
	for _, permission := range missing {
		rows = append(rows, models.AppRolePermission{AppRoleID: existing.ID, Permission: permission})
	}

	return db.Create(&rows).Error
}

// createAppRole inserts a new system app role with its permission rows.
func createAppRole(db *gorm.DB, name string, description string, perms []string) error {
	rolePermissions := make([]models.AppRolePermission, 0, len(perms))
	for _, permission := range perms {
		rolePermissions = append(rolePermissions, models.AppRolePermission{Permission: permission})
	}

	role := models.AppRole{
		Name:        name,
		Description: description,
		IsSystem:    true,
		Permissions: rolePermissions,
	}

	return db.Create(&role).Error
}

// seedGroupRole creates a system group role with the given permissions when no
// group role with that name exists, or reconciles an existing one by adding any
// permissions it is missing (add-only).
func seedGroupRole(db *gorm.DB, name string, description string, perms []string) error {
	var existing models.GroupRoleDefinition
	err := db.Preload("Permissions").Where("name = ?", name).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return createGroupRole(db, name, description, perms)
	}
	if err != nil {
		return err
	}

	have := make([]string, 0, len(existing.Permissions))
	for _, permission := range existing.Permissions {
		have = append(have, permission.Permission)
	}

	missing := missingPermissions(have, perms)
	if len(missing) == 0 {
		return nil
	}

	rows := make([]models.GroupRolePermission, 0, len(missing))
	for _, permission := range missing {
		rows = append(rows, models.GroupRolePermission{GroupRoleID: existing.ID, Permission: permission})
	}

	return db.Create(&rows).Error
}

// createGroupRole inserts a new system group role with its permission rows.
func createGroupRole(db *gorm.DB, name string, description string, perms []string) error {
	rolePermissions := make([]models.GroupRolePermission, 0, len(perms))
	for _, permission := range perms {
		rolePermissions = append(rolePermissions, models.GroupRolePermission{Permission: permission})
	}

	role := models.GroupRoleDefinition{
		Name:        name,
		Description: description,
		IsSystem:    true,
		IsDefault:   false,
		Permissions: rolePermissions,
	}

	return db.Create(&role).Error
}

// missingPermissions returns the desired permissions not already present in
// have, de-duplicated and in desired's order. Reconciliation is add-only: a
// permission in have but absent from desired is left untouched, never removed.
func missingPermissions(have []string, desired []string) []string {
	present := make(map[string]struct{}, len(have))
	for _, permission := range have {
		present[permission] = struct{}{}
	}

	missing := make([]string, 0)
	for _, permission := range desired {
		if _, ok := present[permission]; ok {
			continue
		}
		present[permission] = struct{}{} // also de-dupes desired against itself
		missing = append(missing, permission)
	}

	return missing
}

// EnsureDefaultRoles guarantees that exactly one app role and one group role are
// flagged as the default — the role newly-created accounts (app) and group
// creators (group) are assigned. It is the single source of truth for "the
// default role": the legacy-equivalent roles (Legacy User / Legacy Owner) are
// the seeded defaults so upgrading installs behave exactly as before.
//
// It only acts when a scope has no default yet, so it is idempotent and never
// overrides a default an administrator has set through the UI. Run on every boot
// after SeedSystemRoles (a fresh database self-heals on the next boot, and an
// existing dev database that predates the IsDefault column gets backfilled).
func EnsureDefaultRoles() error {
	db := GetDB()

	if err := ensureDefaultAppRole(db); err != nil {
		return err
	}

	return ensureDefaultGroupRole(db)
}

// ensureDefaultAppRole sets the Legacy User role as the default app role when no
// app role is currently flagged default.
func ensureDefaultAppRole(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.AppRole{}).Where("is_default = ?", true).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	result := db.Model(&models.AppRole{}).
		Where("name = ?", LegacyUserRoleName).
		Update("is_default", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("legacy app role %q not found; cannot set default", LegacyUserRoleName)
	}
	return nil
}

// ensureDefaultGroupRole sets the Legacy Owner role as the default group role
// when no group role is currently flagged default.
func ensureDefaultGroupRole(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.GroupRoleDefinition{}).Where("is_default = ?", true).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	result := db.Model(&models.GroupRoleDefinition{}).
		Where("name = ?", LegacyOwnerRoleName).
		Update("is_default", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("legacy group role %q not found; cannot set default", LegacyOwnerRoleName)
	}
	return nil
}
