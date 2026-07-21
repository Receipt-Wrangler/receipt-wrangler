package repositories

import (
	"receipt-wrangler/api/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
// not yet been recorded in the data_migrations ledger. Each migration claims its
// ledger row and runs in the same transaction, so a failure rolls back cleanly
// and the migration retries on the next boot.
func RunDataMigrations() error {
	db := GetDB()

	for _, migration := range dataMigrations {
		err := db.Transaction(func(tx *gorm.DB) error {
			// Atomically claim the migration. A conflict on the unique name means
			// it was already applied (or claimed by a concurrent boot), so skip —
			// this keeps concurrent startups from both running the migration or
			// failing startup on the unique constraint.
			claim := tx.Clauses(clause.OnConflict{DoNothing: true}).
				Create(&models.DataMigration{Name: migration.name, AppliedAt: time.Now()})
			if claim.Error != nil {
				return claim.Error
			}
			if claim.RowsAffected == 0 {
				return nil
			}
			return migration.run(tx)
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
	// The legacy UserRole / GroupRole enums have been removed from the Go models,
	// but the physical user_role / group_role columns are intentionally retained on
	// existing installs (GORM never drops them) so this back-fill can read them. The
	// values are matched as plain strings, and each loop is guarded by HasColumn so
	// fresh installs — which never create the columns — skip cleanly instead of
	// failing with "no such column".
	appRoleByLegacy := []struct {
		legacyRole string
		roleName   string
	}{
		{"ADMIN", LegacyAdminRoleName},
		{"USER", LegacyUserRoleName},
	}
	if tx.Migrator().HasColumn(&models.User{}, "user_role") {
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
	}

	groupRoleByLegacy := []struct {
		legacyRole string
		roleName   string
	}{
		{"OWNER", LegacyOwnerRoleName},
		{"EDITOR", LegacyEditorRoleName},
		{"VIEWER", LegacyViewerRoleName},
	}
	if tx.Migrator().HasColumn(&models.GroupMember{}, "group_role") {
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
	}

	return nil
}
