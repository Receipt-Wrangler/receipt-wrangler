package repositories

import (
	"receipt-wrangler/api/internal/models"
	"time"

	"gorm.io/gorm"
)

// Name of the one-time migration that assigns the seeded legacy-equivalent roles
// to existing users and group members.
const assignLegacyEquivalentRolesMigration = "assign-legacy-equivalent-roles"

// dataMigration is a single one-time data migration. Each runs at most once per
// database; once applied it is recorded in the data_migrations ledger so it is
// skipped on subsequent boots.
type dataMigration struct {
	name string
	run  func(tx *gorm.DB) error
}

// dataMigrations is the ordered registry of one-time data migrations. New
// migrations are appended here.
var dataMigrations = []dataMigration{
	{name: assignLegacyEquivalentRolesMigration, run: assignLegacyEquivalentRoles},
}

// RunDataMigrations applies any registered one-time data migrations that have
// not yet been recorded in the data_migrations ledger. Each migration runs in
// its own transaction together with the ledger insert, so a failure rolls back
// cleanly and the migration retries on the next boot.
func RunDataMigrations() error {
	db := GetDB()

	for _, migration := range dataMigrations {
		var count int64
		if err := db.Model(&models.DataMigration{}).Where("name = ?", migration.name).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := migration.run(tx); err != nil {
				return err
			}
			return tx.Create(&models.DataMigration{Name: migration.name, AppliedAt: time.Now()}).Error
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// assignLegacyEquivalentRoles back-fills the new role assignments from the
// legacy UserRole / GroupRole enums so existing installs upgrade with zero
// behavior change: every user is assigned the app role and every group member
// the group role that reproduces its legacy capabilities.
//
// Updates are guarded by "... IS NULL" so an assignment an administrator has
// already made through the new role UI is never overwritten.
func assignLegacyEquivalentRoles(tx *gorm.DB) error {
	appRoleByLegacy := []struct {
		legacyRole models.UserRole
		roleName   string
	}{
		{models.ADMIN, LegacyAdminRoleName},
		{models.USER, LegacyUserRoleName},
	}
	for _, mapping := range appRoleByLegacy {
		var role models.AppRole
		if err := tx.Where("name = ?", mapping.roleName).First(&role).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).
			Where("user_role = ? AND app_role_id IS NULL", mapping.legacyRole).
			Update("app_role_id", role.ID).Error; err != nil {
			return err
		}
	}

	groupRoleByLegacy := []struct {
		legacyRole models.GroupRole
		roleName   string
	}{
		{models.OWNER, LegacyOwnerRoleName},
		{models.EDITOR, LegacyEditorRoleName},
		{models.VIEWER, LegacyViewerRoleName},
	}
	for _, mapping := range groupRoleByLegacy {
		var role models.GroupRoleDefinition
		if err := tx.Where("name = ?", mapping.roleName).First(&role).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.GroupMember{}).
			Where("group_role = ? AND group_role_id IS NULL", mapping.legacyRole).
			Update("group_role_id", role.ID).Error; err != nil {
			return err
		}
	}

	return nil
}
