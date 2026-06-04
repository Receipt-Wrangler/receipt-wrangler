package repositories

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"
)

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
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	admin := models.User{Username: "admin", Password: "password", UserRole: models.ADMIN}
	standard := models.User{Username: "standard", Password: "password", UserRole: models.USER}
	viewerUser := models.User{Username: "viewer", Password: "password", UserRole: models.USER}
	for _, user := range []*models.User{&admin, &standard, &viewerUser} {
		if err := db.Create(user).Error; err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
	}

	group := models.Group{Name: "migration-group"}
	if err := db.Create(&group).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	ownerMember := models.GroupMember{GroupID: group.ID, UserID: admin.ID, GroupRole: models.OWNER}
	editorMember := models.GroupMember{GroupID: group.ID, UserID: standard.ID, GroupRole: models.EDITOR}
	viewerMember := models.GroupMember{GroupID: group.ID, UserID: viewerUser.ID, GroupRole: models.VIEWER}
	for _, member := range []*models.GroupMember{&ownerMember, &editorMember, &viewerMember} {
		if err := db.Create(member).Error; err != nil {
			utils.PrintTestError(t, err, nil)
			return
		}
	}

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
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	admin := models.User{Username: "admin", Password: "password", UserRole: models.ADMIN}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

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

	admin := models.User{Username: "admin", Password: "password", UserRole: models.ADMIN}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

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
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	db := GetDB()

	customRole, err := NewRoleRepository(nil).CreateAppRole("Custom Role", "", []string{permissions.AppUsersRead})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	admin := models.User{Username: "admin", Password: "password", UserRole: models.ADMIN, AppRoleID: &customRole.ID}
	if err := db.Create(&admin).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if err := RunDataMigrations(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// The IS NULL guard leaves the administrator's existing assignment intact.
	assertAppRoleId(t, reloadUser(t, admin.ID), customRole.ID)
}
