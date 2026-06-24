package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"
)

// The legacy user_role / group_role columns were removed from the Go models, so
// AutoMigrate no longer creates them on the test database. The assignLegacyEquivalentRoles
// migration still reads them as raw string columns on installs upgraded from the
// legacy schema, guarded by HasColumn. These helpers reconstruct (and tear down)
// those physical columns so the back-fill can be exercised end to end. They are
// idempotent and guarded by HasColumn because the test DB schema persists across
// subtests (TruncateTestDb only deletes rows).

func ensureLegacyUserRoleColumn(t *testing.T) {
	db := GetDB()
	if !db.Migrator().HasColumn(&models.User{}, "user_role") {
		if err := db.Exec("ALTER TABLE users ADD COLUMN user_role text").Error; err != nil {
			utils.PrintTestError(t, err, "adding the legacy user_role column")
		}
	}
}

func ensureLegacyGroupRoleColumn(t *testing.T) {
	db := GetDB()
	if !db.Migrator().HasColumn(&models.GroupMember{}, "group_role") {
		if err := db.Exec("ALTER TABLE group_members ADD COLUMN group_role text").Error; err != nil {
			utils.PrintTestError(t, err, "adding the legacy group_role column")
		}
	}
}

func dropLegacyRoleColumns(t *testing.T) {
	db := GetDB()
	if db.Migrator().HasColumn(&models.User{}, "user_role") {
		if err := db.Exec("ALTER TABLE users DROP COLUMN user_role").Error; err != nil {
			utils.PrintTestError(t, err, "dropping the legacy user_role column")
		}
	}
	if db.Migrator().HasColumn(&models.GroupMember{}, "group_role") {
		if err := db.Exec("ALTER TABLE group_members DROP COLUMN group_role").Error; err != nil {
			utils.PrintTestError(t, err, "dropping the legacy group_role column")
		}
	}
}

// setLegacyUserRole writes a value into the raw user_role column for one user.
func setLegacyUserRole(t *testing.T, userId uint, legacyRole string) {
	if err := GetDB().Model(&models.User{}).Where("id = ?", userId).
		UpdateColumn("user_role", legacyRole).Error; err != nil {
		utils.PrintTestError(t, err, "setting the legacy user_role column")
	}
}

// setLegacyGroupRole writes a value into the raw group_role column for one member.
func setLegacyGroupRole(t *testing.T, userId uint, groupId uint, legacyRole string) {
	if err := GetDB().Model(&models.GroupMember{}).
		Where("user_id = ? AND group_id = ?", userId, groupId).
		UpdateColumn("group_role", legacyRole).Error; err != nil {
		utils.PrintTestError(t, err, "setting the legacy group_role column")
	}
}

func appRoleIdByName(t *testing.T, name string) uint {
	var role models.AppRole
	if err := GetDB().Where("name = ?", name).First(&role).Error; err != nil {
		utils.PrintTestError(t, err, "an app role named "+name)
	}
	return role.ID
}

func groupRoleIdByName(t *testing.T, name string) uint {
	var role models.GroupRoleDefinition
	if err := GetDB().Where("name = ?", name).First(&role).Error; err != nil {
		utils.PrintTestError(t, err, "a group role named "+name)
	}
	return role.ID
}

func reloadUser(t *testing.T, id uint) models.User {
	var user models.User
	if err := GetDB().First(&user, id).Error; err != nil {
		utils.PrintTestError(t, err, "a user to reload")
	}
	return user
}

func reloadMember(t *testing.T, userId uint, groupId uint) models.GroupMember {
	var member models.GroupMember
	if err := GetDB().Where("user_id = ? AND group_id = ?", userId, groupId).First(&member).Error; err != nil {
		utils.PrintTestError(t, err, "a group member")
	}
	return member
}

func assertAppRoleId(t *testing.T, user models.User, expected uint) {
	if user.AppRoleID == nil || *user.AppRoleID != expected {
		utils.PrintTestError(t, user.AppRoleID, expected)
	}
}

func assertGroupRoleId(t *testing.T, member models.GroupMember, expected uint) {
	if member.GroupRoleID == nil || *member.GroupRoleID != expected {
		utils.PrintTestError(t, member.GroupRoleID, expected)
	}
}

func TestRunDataMigrationsAssignsLegacyEquivalentRoles(t *testing.T) {
	defer TruncateTestDb()
	defer dropLegacyRoleColumns(t)
	ensureLegacyUserRoleColumn(t)
	ensureLegacyGroupRoleColumn(t)
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	// The legacy enum fields no longer exist on the Go models, so the rows are
	// created without them and the legacy user_role / group_role columns are
	// populated separately via raw column writes — exactly the on-disk shape an
	// upgraded install presents to the back-fill.
	admin := models.User{Username: "admin", Password: "password"}
	standard := models.User{Username: "standard", Password: "password"}
	viewerUser := models.User{Username: "viewer", Password: "password"}
	for _, user := range []*models.User{&admin, &standard, &viewerUser} {
		if err := db.Create(user).Error; err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
	}
	setLegacyUserRole(t, admin.ID, "ADMIN")
	setLegacyUserRole(t, standard.ID, "USER")
	setLegacyUserRole(t, viewerUser.ID, "USER")

	group := models.Group{Name: "migration-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	ownerMember := models.GroupMember{GroupID: group.ID, UserID: admin.ID}
	editorMember := models.GroupMember{GroupID: group.ID, UserID: standard.ID}
	viewerMember := models.GroupMember{GroupID: group.ID, UserID: viewerUser.ID}
	for _, member := range []*models.GroupMember{&ownerMember, &editorMember, &viewerMember} {
		if err := db.Create(member).Error; err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
	}
	setLegacyGroupRole(t, admin.ID, group.ID, "OWNER")
	setLegacyGroupRole(t, standard.ID, group.ID, "EDITOR")
	setLegacyGroupRole(t, viewerUser.ID, group.ID, "VIEWER")

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	assertAppRoleId(t, reloadUser(t, admin.ID), appRoleIdByName(t, LegacyAdminRoleName))
	assertAppRoleId(t, reloadUser(t, standard.ID), appRoleIdByName(t, LegacyUserRoleName))
	assertAppRoleId(t, reloadUser(t, viewerUser.ID), appRoleIdByName(t, LegacyUserRoleName))

	assertGroupRoleId(t, reloadMember(t, admin.ID, group.ID), groupRoleIdByName(t, LegacyOwnerRoleName))
	assertGroupRoleId(t, reloadMember(t, standard.ID, group.ID), groupRoleIdByName(t, LegacyEditorRoleName))
	assertGroupRoleId(t, reloadMember(t, viewerUser.ID, group.ID), groupRoleIdByName(t, LegacyViewerRoleName))

	var ledgerCount int64
	if err := db.Model(&models.DataMigration{}).Where("name = ?", assignLegacyEquivalentRolesMigration).Count(&ledgerCount).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ledgerCount != 1 {
		utils.PrintTestError(t, ledgerCount, 1)
	}
}

func TestRunDataMigrationsIsIdempotent(t *testing.T) {
	defer TruncateTestDb()
	defer dropLegacyRoleColumns(t)
	ensureLegacyUserRoleColumn(t)
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	admin := models.User{Username: "admin", Password: "password"}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	setLegacyUserRole(t, admin.ID, "ADMIN")

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	assertAppRoleId(t, reloadUser(t, admin.ID), appRoleIdByName(t, LegacyAdminRoleName))

	var ledgerCount int64
	if err := db.Model(&models.DataMigration{}).Where("name = ?", assignLegacyEquivalentRolesMigration).Count(&ledgerCount).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ledgerCount != 1 {
		utils.PrintTestError(t, ledgerCount, 1)
	}
}

func TestRunDataMigrationsSkipsWhenAlreadyApplied(t *testing.T) {
	defer TruncateTestDb()
	defer dropLegacyRoleColumns(t)
	ensureLegacyUserRoleColumn(t)
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	// Record the migration as already applied before any rows exist to assign.
	if err := db.Create(&models.DataMigration{Name: assignLegacyEquivalentRolesMigration, AppliedAt: time.Now()}).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	admin := models.User{Username: "admin", Password: "password"}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	setLegacyUserRole(t, admin.ID, "ADMIN")

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// The ledger short-circuits the run, so the user is left unassigned.
	if reloadUser(t, admin.ID).AppRoleID != nil {
		utils.PrintTestError(t, reloadUser(t, admin.ID).AppRoleID, nil)
	}
}

func TestRunDataMigrationsDoesNotClobberExistingAssignment(t *testing.T) {
	defer TruncateTestDb()
	defer dropLegacyRoleColumns(t)
	ensureLegacyUserRoleColumn(t)
	ensureLegacyGroupRoleColumn(t)
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	roleRepository := NewRoleRepository(nil)
	customAppRole, err := roleRepository.CreateAppRole("Custom Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	customGroupRole, err := roleRepository.CreateGroupRole("Custom Group Role", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	admin := models.User{Username: "admin", Password: "password", AppRoleID: &customAppRole.ID}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	setLegacyUserRole(t, admin.ID, "ADMIN")

	group := models.Group{Name: "no-clobber-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	member := models.GroupMember{GroupID: group.ID, UserID: admin.ID, GroupRoleID: &customGroupRole.ID}
	if err := db.Create(&member).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	setLegacyGroupRole(t, admin.ID, group.ID, "OWNER")

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// The IS NULL guard leaves an administrator's existing assignments intact,
	// for both the app role and the group role.
	assertAppRoleId(t, reloadUser(t, admin.ID), customAppRole.ID)
	assertGroupRoleId(t, reloadMember(t, admin.ID, group.ID), customGroupRole.ID)
}

func TestRunDataMigrationsSkipsBackfillWhenLegacyColumnsAbsent(t *testing.T) {
	defer TruncateTestDb()

	// A fresh install never had the legacy user_role / group_role columns. Make
	// sure they are absent (a prior subtest may have added them), then confirm the
	// HasColumn guard makes the back-fill a no-op rather than failing with
	// "no such column".
	dropLegacyRoleColumns(t)
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	user := models.User{Username: "fresh", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	group := models.Group{Name: "fresh-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	member := models.GroupMember{GroupID: group.ID, UserID: user.ID}
	if err := db.Create(&member).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, "no error when the legacy columns are absent")
		return
	}

	// With the legacy columns absent, the back-fill leaves both FKs as they were
	// (nil) instead of assigning a role.
	if reloadUser(t, user.ID).AppRoleID != nil {
		utils.PrintTestError(t, reloadUser(t, user.ID).AppRoleID, nil)
	}
	if reloadMember(t, user.ID, group.ID).GroupRoleID != nil {
		utils.PrintTestError(t, reloadMember(t, user.ID, group.ID).GroupRoleID, nil)
	}

	// The migration still records its ledger row (it ran successfully, just with
	// nothing to back-fill).
	var ledgerCount int64
	if err := db.Model(&models.DataMigration{}).Where("name = ?", assignLegacyEquivalentRolesMigration).Count(&ledgerCount).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ledgerCount != 1 {
		utils.PrintTestError(t, ledgerCount, 1)
	}
}

func TestRunDataMigrationsRollsBackOnFailure(t *testing.T) {
	defer TruncateTestDb()
	defer dropLegacyRoleColumns(t)
	ensureLegacyUserRoleColumn(t)
	db := GetDB()

	// Intentionally skip SeedSystemRoles so the migration's first role lookup
	// fails with ErrRecordNotFound, exercising the error path.
	admin := models.User{Username: "admin", Password: "password"}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	setLegacyUserRole(t, admin.ID, "ADMIN")

	if err := RunDataMigrations(); err == nil {
		utils.PrintTestError(t, nil, "an error because the legacy roles are not seeded")
	}

	// The transaction rolls back: no ledger row is written, so the migration
	// retries on the next boot.
	var ledgerCount int64
	if err := db.Model(&models.DataMigration{}).Where("name = ?", assignLegacyEquivalentRolesMigration).Count(&ledgerCount).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ledgerCount != 0 {
		utils.PrintTestError(t, ledgerCount, 0)
	}

	// And no partial assignment persisted.
	if reloadUser(t, admin.ID).AppRoleID != nil {
		utils.PrintTestError(t, reloadUser(t, admin.ID).AppRoleID, nil)
	}
}
