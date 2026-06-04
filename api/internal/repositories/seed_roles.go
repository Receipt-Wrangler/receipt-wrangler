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
		{"Legacy Admin", "Legacy administrator equivalent: full application access.", permissions.LegacyAppAdminKeys()},
		{"Legacy User", "Legacy standard user equivalent.", permissions.LegacyAppUserKeys()},
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
		{"Legacy Viewer", "Legacy group viewer equivalent.", permissions.LegacyGroupViewerKeys()},
		{"Legacy Editor", "Legacy group editor equivalent.", permissions.LegacyGroupEditorKeys()},
		{"Legacy Owner", "Legacy group owner equivalent: full group access.", permissions.LegacyGroupOwnerKeys()},
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
