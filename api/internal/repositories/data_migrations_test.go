package repositories

import (
	"database/sql"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"
)

// The legacy user_role column was removed from the User model, so AutoMigrate no longer
// creates it on the test database. assignLegacyEquivalentRoles still reads it as a raw
// string column on installs upgraded from the legacy schema, guarded by HasColumn. These
// helpers reconstruct (and tear down) that physical column so the back-fill can be
// exercised end to end. They are idempotent and guarded by HasColumn because the test DB
// schema persists across subtests (TruncateTestDb only deletes rows).
//
// group_role is different: GroupMember.GroupRole was re-declared on the model (nullable,
// json:"-") to satisfy the leftover NOT NULL group_role column on upgraded installs, so
// AutoMigrate now always creates it. It is therefore never added or dropped here — doing so
// would leave a field-without-column state that breaks group_members inserts elsewhere in
// this package's test run.

func ensureLegacyUserRoleColumn(t *testing.T) {
	db := GetDB()
	if !db.Migrator().HasColumn(&models.User{}, "user_role") {
		if err := db.Exec("ALTER TABLE users ADD COLUMN user_role text").Error; err != nil {
			utils.PrintTestError(t, err, "adding the legacy user_role column")
		}
	}
}

func dropLegacyUserRoleColumn(t *testing.T) {
	db := GetDB()
	if db.Migrator().HasColumn(&models.User{}, "user_role") {
		if err := db.Exec("ALTER TABLE users DROP COLUMN user_role").Error; err != nil {
			utils.PrintTestError(t, err, "dropping the legacy user_role column")
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
	defer dropLegacyUserRoleColumn(t)
	ensureLegacyUserRoleColumn(t)
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
	defer dropLegacyUserRoleColumn(t)
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
	defer dropLegacyUserRoleColumn(t)
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
	defer dropLegacyUserRoleColumn(t)
	ensureLegacyUserRoleColumn(t)
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

	// A fresh install never had the legacy user_role column. Make sure it is absent
	// (a prior subtest may have added it), then confirm the HasColumn guard makes the
	// app back-fill a no-op rather than failing with "no such column". The group_role
	// column is always present now (GroupMember.GroupRole is back on the model), but a
	// fresh member's value is "" — which matches no legacy enum — so the group back-fill
	// is likewise a no-op.
	dropLegacyUserRoleColumn(t)
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

	// With no legacy values to read, the back-fill leaves both FKs as they were
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
	defer dropLegacyUserRoleColumn(t)
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

// TestNewGroupMemberPopulatesLegacyGroupRoleColumn guards the upgrade fix: on databases
// upgraded from before the role rework the obsolete group_role column survives NOT NULL
// with no default. The re-declared GroupMember.GroupRole field makes GORM write a value
// on every INSERT, so the constraint is satisfied and new group_members rows can be
// created again. Here we assert that mechanism — a freshly created member persists a
// non-NULL group_role (it is written, not omitted) — which is exactly what keeps the
// leftover NOT NULL constraint from rejecting the insert.
func TestNewGroupMemberPopulatesLegacyGroupRoleColumn(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	user := models.User{Username: "legacy-col", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	group := models.Group{Name: "legacy-col-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	member := models.GroupMember{GroupID: group.ID, UserID: user.ID}
	if err := db.Create(&member).Error; err != nil {
		utils.PrintTestError(t, err, "creating a group member must not violate the legacy NOT NULL group_role column")
		return
	}

	var groupRole sql.NullString
	row := db.Raw("SELECT group_role FROM group_members WHERE user_id = ? AND group_id = ?", user.ID, group.ID).Row()
	if err := row.Scan(&groupRole); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !groupRole.Valid {
		utils.PrintTestError(t, "group_role was NULL", "a non-NULL group_role value written on INSERT")
	}
}

// legacyOwnerMemberPermissions returns the group.members.* permissions currently
// attached to the seeded Legacy Owner role.
func legacyOwnerMemberPermissions(t *testing.T) []string {
	ownerRoleId := groupRoleIdByName(t, LegacyOwnerRoleName)
	var perms []string
	if err := GetDB().Model(&models.GroupRolePermission{}).
		Where("group_role_id = ? AND permission IN ?", ownerRoleId, []string{
			permissions.GroupMembersCreate,
			permissions.GroupMembersUpdate,
			permissions.GroupMembersDelete,
		}).
		Pluck("permission", &perms).Error; err != nil {
		utils.PrintTestError(t, err, "loading Legacy Owner member permissions")
	}
	return perms
}

// TestRunDataMigrationsAddsGroupMembersPermissionsToLegacyOwner simulates an
// install seeded before the member-management permissions existed: the Legacy
// Owner role is stripped of them, and the migration must back-fill all three so
// existing owners keep managing members.
func TestRunDataMigrationsAddsGroupMembersPermissionsToLegacyOwner(t *testing.T) {
	defer TruncateTestDb()
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	ownerRoleId := groupRoleIdByName(t, LegacyOwnerRoleName)
	if err := db.Where("group_role_id = ? AND permission IN ?", ownerRoleId, []string{
		permissions.GroupMembersCreate,
		permissions.GroupMembersUpdate,
		permissions.GroupMembersDelete,
	}).Delete(&models.GroupRolePermission{}).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if got := legacyOwnerMemberPermissions(t); len(got) != 0 {
		utils.PrintTestError(t, got, "no member permissions before the migration")
	}

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if got := legacyOwnerMemberPermissions(t); len(got) != 3 {
		utils.PrintTestError(t, got, "3 member permissions after the migration")
	}

	var ledgerCount int64
	if err := db.Model(&models.DataMigration{}).
		Where("name = ?", addGroupMembersPermissionsToLegacyOwnerMigration).
		Count(&ledgerCount).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if ledgerCount != 1 {
		utils.PrintTestError(t, ledgerCount, 1)
	}
}

// TestAddGroupMembersPermissionsToLegacyOwnerIsIdempotent confirms that running
// the back-fill against a Legacy Owner that already holds the member permissions
// (the fresh-install case) does not create duplicate permission rows.
func TestAddGroupMembersPermissionsToLegacyOwnerIsIdempotent(t *testing.T) {
	defer TruncateTestDb()
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Fresh seeding already grants the full group scope, so the member permissions
	// are present; running the migration again must be a no-op.
	if err := GetDB().Transaction(addGroupMembersPermissionsToLegacyOwner); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if got := legacyOwnerMemberPermissions(t); len(got) != 3 {
		utils.PrintTestError(t, got, "exactly 3 member permissions (no duplicates)")
	}
}
