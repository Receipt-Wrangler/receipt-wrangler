package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"

	"gorm.io/gorm"
)

// SeedSystemRoles creates the five immutable, legacy-equivalent system roles if
// they do not already exist. The roles reproduce the legacy UserRole/GroupRole
// capabilities exactly so that upgrading installs see no behavior change.
//
// It is idempotent: keyed on the role Name (a uniqueIndex on both AppRole and
// GroupRoleDefinition), it is safe to run on every startup and on upgrades, and
// a partially-seeded database self-heals on the next boot. The roles are seeded
// only — they are not assigned to any user or group member here.
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

// seedAppRole creates a system app role with the given permissions if no app
// role with that name already exists.
func seedAppRole(db *gorm.DB, name string, description string, perms []string) error {
	var count int64
	if err := db.Model(&models.AppRole{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

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

// seedGroupRole creates a system group role with the given permissions if no
// group role with that name already exists.
func seedGroupRole(db *gorm.DB, name string, description string, perms []string) error {
	var count int64
	if err := db.Model(&models.GroupRoleDefinition{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

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

	return db.Model(&models.AppRole{}).
		Where("name = ?", LegacyUserRoleName).
		Update("is_default", true).Error
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

	return db.Model(&models.GroupRoleDefinition{}).
		Where("name = ?", LegacyOwnerRoleName).
		Update("is_default", true).Error
}
