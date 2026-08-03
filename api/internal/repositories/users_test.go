package repositories

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestShouldCreateAdminUserWithGroup(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()
	userToCreate := commands.SignUpCommand{
		Username:    "test",
		DisplayName: "test",
		Password:    "a really secure password",
	}
	userRepository := NewUserRepository(nil)
	createdUser, err := userRepository.CreateUser(userToCreate)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	validateUser(t, createdUser, userToCreate, true, 1)

	var group models.Group
	db.Table("groups").Where("id = 1").Preload("GroupMembers").First(&group)

	validateGroup(t, group, 1, 1)
}

func TestShouldCreateNonAdminUserWithGroup(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()
	CreateTestUser()
	CreateTestGroup()
	userToCreate := commands.SignUpCommand{
		Username:    "test2",
		DisplayName: "test",
		Password:    "a really secure password",
	}
	userRepository := NewUserRepository(nil)
	createdUser, err := userRepository.CreateUser(userToCreate)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	validateUser(t, createdUser, userToCreate, false, 2)

	var group models.Group
	db.Table("groups").Where("id = 2").Preload("GroupMembers").First(&group)

	validateGroup(t, group, 2, 2)
}

func TestCreateUserAssignsDefaultAppRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userRepository := NewUserRepository(nil)
	roleRepository := NewRoleRepository(nil)

	// The first user becomes ADMIN and must get the Legacy Admin role (never the
	// configurable default), so the bootstrap admin is never locked out.
	admin, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "admin", DisplayName: "admin", Password: "a really secure password",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	legacyAdminId, err := roleRepository.GetAppRoleIdByName(LegacyAdminRoleName)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if admin.AppRoleID == nil || legacyAdminId == nil || *admin.AppRoleID != *legacyAdminId {
		utils.PrintTestError(t, admin.AppRoleID, legacyAdminId)
	}

	// Subsequent users become USER and must get the configurable default app role.
	user, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "user", DisplayName: "user", Password: "a really secure password",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	defaultId, err := roleRepository.GetDefaultAppRoleId()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if user.AppRoleID == nil || defaultId == nil || *user.AppRoleID != *defaultId {
		utils.PrintTestError(t, user.AppRoleID, defaultId)
	}
}

func TestCreateUserHonorsExplicitAppRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userRepository := NewUserRepository(nil)
	roleRepository := NewRoleRepository(nil)

	legacyAdminId, err := roleRepository.GetAppRoleIdByName(LegacyAdminRoleName)
	if err != nil || legacyAdminId == nil {
		utils.PrintTestError(t, err, "legacy admin id")
		return
	}

	// A bootstrap user already occupies the first-user slot, so the count-based
	// fallback would make the next account a USER. An explicit modern app role id
	// (the admin-create path) must win and set the modern role FK.
	if _, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "bootstrap", DisplayName: "b", Password: "a really secure password",
	}); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	admin, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "modern-admin", DisplayName: "a", Password: "a really secure password",
		AppRoleID: legacyAdminId,
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if admin.AppRoleID == nil || *admin.AppRoleID != *legacyAdminId {
		utils.PrintTestError(t, admin.AppRoleID, legacyAdminId)
	}

	// A custom (non-system) app role is honored as-is on the modern FK.
	custom, err := roleRepository.CreateAppRole("Auditor", "", []string{permissions.AppUsersRead}, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	auditor, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "auditor", DisplayName: "a", Password: "a really secure password",
		AppRoleID: &custom.ID,
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if auditor.AppRoleID == nil || *auditor.AppRoleID != custom.ID {
		utils.PrintTestError(t, auditor.AppRoleID, custom.ID)
	}
}

// groupNamesForUser returns the names of every group the user is a member of.
func groupNamesForUser(t *testing.T, userId uint) []string {
	var groups []models.Group
	err := GetDB().Model(&models.Group{}).
		Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ?", userId).
		Find(&groups).Error
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return nil
	}

	names := make([]string, 0, len(groups))
	for _, group := range groups {
		names = append(names, group.Name)
	}

	return names
}

func TestCreateUserSkipsDefaultGroupForFlaggedAppRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userRepository := NewUserRepository(nil)
	roleRepository := NewRoleRepository(nil)

	// Occupy the first-user slot so the accounts below take their explicit role
	// rather than the bootstrap Legacy Admin.
	if _, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "bootstrap", DisplayName: "b", Password: "a really secure password",
	}); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	skipRole, err := roleRepository.CreateAppRole("Shared Groups Only", "", []string{permissions.AppUsersRead}, true)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	restricted, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "restricted", DisplayName: "r", Password: "a really secure password",
		AppRoleID: &skipRole.ID,
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// The personal group is skipped, but the virtual "All" group is always created
	// so the account still has a working dashboard.
	names := groupNamesForUser(t, restricted.ID)
	if len(names) != 1 || names[0] != "All" {
		utils.PrintTestError(t, names, []string{"All"})
	}
}

// The two best-effort branches of the helper, which the CreateUser tests below
// can't reach: a user with no app role at all, and an id with no matching row.
// Both must report "don't skip" rather than erroring, so user creation proceeds
// with the personal group instead of failing. (A non-record-not-found lookup
// error deliberately propagates instead — that path needs DB error injection,
// which this package has no mechanism for. The OnDelete:RESTRICT FK on
// User.AppRoleID rules out dangling role ids, not connection, transaction, or
// context errors.)
func TestAppRoleSkipsDefaultGroupBestEffortBranches(t *testing.T) {
	defer TruncateTestDb()
	userRepository := NewUserRepository(nil)

	skip, err := userRepository.appRoleSkipsDefaultGroup(nil, nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if skip {
		utils.PrintTestError(t, skip, false)
	}

	missingId := uint(999999)
	skip, err = userRepository.appRoleSkipsDefaultGroup(nil, &missingId)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if skip {
		utils.PrintTestError(t, skip, false)
	}
}

func TestCreateUserCreatesDefaultGroupForUnflaggedAppRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userRepository := NewUserRepository(nil)

	if _, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "bootstrap", DisplayName: "b", Password: "a really secure password",
	}); err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// The default app role is unflagged, so the personal group is still created.
	normal, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "normal", DisplayName: "n", Password: "a really secure password",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	names := groupNamesForUser(t, normal.ID)
	if len(names) != 2 {
		utils.PrintTestError(t, names, []string{"My Receipts", "All"})
		return
	}

	var hasPersonal, hasAll bool
	for _, name := range names {
		if name == "My Receipts" {
			hasPersonal = true
		}
		if name == "All" {
			hasAll = true
		}
	}
	if !hasPersonal || !hasAll {
		utils.PrintTestError(t, names, []string{"My Receipts", "All"})
	}
}

func TestUpdateUserAssignsAppRole(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userRepository := NewUserRepository(nil)
	roleRepository := NewRoleRepository(nil)

	created, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "user", DisplayName: "u", Password: "a really secure password",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	legacyAdminId, err := roleRepository.GetAppRoleIdByName(LegacyAdminRoleName)
	if err != nil || legacyAdminId == nil {
		utils.PrintTestError(t, err, "legacy admin id")
		return
	}

	err = userRepository.UpdateUser(utils.UintToString(created.ID), commands.SignUpCommand{
		Username: "user", DisplayName: "u", AppRoleID: legacyAdminId,
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	var updated models.User
	if err := GetDB().First(&updated, created.ID).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if updated.AppRoleID == nil || *updated.AppRoleID != *legacyAdminId {
		utils.PrintTestError(t, updated.AppRoleID, legacyAdminId)
	}
}

func TestUpdateUserPreservesAppRoleWhenIdOmitted(t *testing.T) {
	defer TruncateTestDb()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userRepository := NewUserRepository(nil)

	// First user → Legacy Admin app role.
	created, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "admin", DisplayName: "a", Password: "a really secure password",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	originalRoleId := created.AppRoleID
	if originalRoleId == nil {
		utils.PrintTestError(t, originalRoleId, "an app role id")
		return
	}

	// An update that omits the app role id (e.g. a display-name-only edit) must
	// not clear the existing assignment.
	err = userRepository.UpdateUser(utils.UintToString(created.ID), commands.SignUpCommand{
		Username: "admin", DisplayName: "renamed",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	var updated models.User
	if err := GetDB().First(&updated, created.ID).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if updated.AppRoleID == nil || *updated.AppRoleID != *originalRoleId {
		utils.PrintTestError(t, updated.AppRoleID, originalRoleId)
	}
	if updated.DisplayName != "renamed" {
		utils.PrintTestError(t, updated.DisplayName, "renamed")
	}
}

func TestCreateUserLeavesAppRoleNilWhenUnseeded(t *testing.T) {
	defer TruncateTestDb()

	userRepository := NewUserRepository(nil)
	created, err := userRepository.CreateUser(commands.SignUpCommand{
		Username: "noroles", DisplayName: "n", Password: "a really secure password",
	})
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if created.AppRoleID != nil {
		utils.PrintTestError(t, created.AppRoleID, nil)
	}
}

func TestShouldReturnErrorWhenCreatingUserWithDuplicateUsername(t *testing.T) {
	defer TruncateTestDb()
	CreateTestUser()
	CreateTestGroup()
	userToCreate := commands.SignUpCommand{
		Username:    "test",
		DisplayName: "test",
		Password:    "a really secure password",
	}
	userRepository := NewUserRepository(nil)
	_, err := userRepository.CreateUser(userToCreate)
	if err == nil {
		utils.PrintTestError(t, err, "error")
	}
}

func TestShouldBeFirstAdminToLogin(t *testing.T) {
	defer TruncateTestDb()

	// IsFirstAdminToLogin identifies administrators by the modern app-role
	// permission (app.users.read), so the system roles must be seeded for the first
	// user to resolve to the Legacy Admin role.
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userToCreate := commands.SignUpCommand{
		Username:    "test",
		DisplayName: "test",
		Password:    "a really secure password",
	}
	userRepository := NewUserRepository(nil)
	createdUser, err := userRepository.CreateUser(userToCreate)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	validateUser(t, createdUser, userToCreate, true, 1)

	firstAdminToLogin, err := userRepository.IsFirstAdminToLogin()

	if firstAdminToLogin != true {
		utils.PrintTestError(t, firstAdminToLogin, true)
	}
}

func TestShouldNotBeFirstAdminToLogin(t *testing.T) {
	defer TruncateTestDb()

	// IsFirstAdminToLogin identifies administrators by the modern app-role
	// permission (app.users.read), so the system roles must be seeded for the first
	// user to resolve to the Legacy Admin role.
	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	userToCreate := commands.SignUpCommand{
		Username:    "test",
		DisplayName: "test",
		Password:    "a really secure password",
	}
	userRepository := NewUserRepository(nil)
	createdUser, err := userRepository.CreateUser(userToCreate)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}

	validateUser(t, createdUser, userToCreate, true, 1)

	userRepository.UpdateUserLastLoginDate(1)

	firstAdminToLogin, err := userRepository.IsFirstAdminToLogin()
	if firstAdminToLogin != false {
		utils.PrintTestError(t, firstAdminToLogin, false)
	}
}

// validateUser asserts the freshly-created user's invariants, including the
// modern app-role FK (the legacy UserRole enum was removed). isFirstUser selects
// the expected role: the bootstrap admin is assigned Legacy Admin, every later
// account the configurable default app role. resolveAppRoleId leaves the FK nil
// in an unseeded test database, so the helper resolves the same expected id via
// the repository getters and matches that (nil == nil when roles aren't seeded).
func validateUser(t *testing.T, createdUser models.User, userToCreate commands.SignUpCommand, isFirstUser bool, id uint) {
	if createdUser.ID != id {
		utils.PrintTestError(t, createdUser.ID, id)
	}
	if createdUser.Password == userToCreate.Password {
		utils.PrintTestError(t, createdUser.Password, "hashed password")
	}
	if createdUser.DefaultAvatarColor != "#27b1ff" {
		utils.PrintTestError(t, createdUser.DefaultAvatarColor, "#27b1ff")
	}

	roleRepository := NewRoleRepository(nil)
	var expectedRoleId *uint
	var err error
	if isFirstUser {
		expectedRoleId, err = roleRepository.GetAppRoleIdByName(LegacyAdminRoleName)
	} else {
		expectedRoleId, err = roleRepository.GetDefaultAppRoleId()
	}
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !uintPtrEqual(createdUser.AppRoleID, expectedRoleId) {
		utils.PrintTestError(t, createdUser.AppRoleID, expectedRoleId)
	}
}

// uintPtrEqual reports whether two *uint point to the same value (both nil counts
// as equal).
func uintPtrEqual(a *uint, b *uint) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func validateGroup(t *testing.T, group models.Group, id uint, userId uint) {
	if group.ID != id {
		utils.PrintTestError(t, group.ID, id)
	}
	if group.GroupMembers[0].UserID != userId {
		utils.PrintTestError(t, group.GroupMembers[0].UserID, userId)
	}
	if group.Name != "My Receipts" {
		utils.PrintTestError(t, group.Name, "My Receipts")
	}

}

func pagedUsersCommand(orderBy string, direction commands.SortDirection) commands.PagedRequestCommand {
	return commands.PagedRequestCommand{
		Page:          1,
		PageSize:      10,
		OrderBy:       orderBy,
		SortDirection: direction,
	}
}

func createUserWithUsername(username string) {
	GetDB().Create(&models.User{
		Username:    username,
		DisplayName: username,
		Password:    "Password",
	})
}

func TestGetPagedUsers_ReturnsRowsAndCount(t *testing.T) {
	defer TruncateTestDb()

	createUserWithUsername("beta")
	createUserWithUsername("alpha")

	repository := NewUserRepository(nil)
	users, count, err := repository.GetPagedUsers(pagedUsersCommand("username", commands.ASCENDING))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if count != 2 {
		utils.PrintTestError(t, count, int64(2))
	}
	if len(users) != 2 {
		utils.PrintTestError(t, len(users), 2)
		return
	}
	// Ascending username order.
	if users[0].Username != "alpha" || users[1].Username != "beta" {
		utils.PrintTestError(t, []string{users[0].Username, users[1].Username}, []string{"alpha", "beta"})
	}
}

func TestGetPagedUsers_SecondPage(t *testing.T) {
	defer TruncateTestDb()

	createUserWithUsername("a")
	createUserWithUsername("b")
	createUserWithUsername("c")

	command := pagedUsersCommand("username", commands.ASCENDING)
	command.Page = 2
	command.PageSize = 2

	users, count, err := NewUserRepository(nil).GetPagedUsers(command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// Count is the unpaged total.
	if count != 3 {
		utils.PrintTestError(t, count, int64(3))
	}
	if len(users) != 1 {
		utils.PrintTestError(t, len(users), 1)
		return
	}
	if users[0].Username != "c" {
		utils.PrintTestError(t, users[0].Username, "c")
	}
}

func TestGetPagedUsers_EmptyReturnsZeroCount(t *testing.T) {
	defer TruncateTestDb()

	users, count, err := NewUserRepository(nil).GetPagedUsers(pagedUsersCommand("username", commands.ASCENDING))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if count != 0 {
		utils.PrintTestError(t, count, int64(0))
	}
	if len(users) != 0 {
		utils.PrintTestError(t, len(users), 0)
	}
}

func TestGetPagedUsers_RejectsInvalidColumn(t *testing.T) {
	defer TruncateTestDb()

	// An order-by outside the allow-list is rejected rather than interpolated as raw SQL.
	_, _, err := NewUserRepository(nil).GetPagedUsers(pagedUsersCommand("password", commands.ASCENDING))
	if err == nil {
		utils.PrintTestError(t, nil, "an invalid column error")
	}
}

func TestGetPagedUsers_SortsByAllowedColumns(t *testing.T) {
	defer TruncateTestDb()

	createUserWithUsername("alpha")
	createUserWithUsername("beta")

	repository := NewUserRepository(nil)

	// username is covered above; display_name, created_at and updated_at must also
	// execute through Sort/Find without error (they are on the allow-list but were
	// never exercised).
	for _, column := range []string{"display_name", "created_at", "updated_at"} {
		users, count, err := repository.GetPagedUsers(pagedUsersCommand(column, commands.DESCENDING))
		if err != nil {
			utils.PrintTestError(t, err, "no error sorting by "+column)
			return
		}
		if count != 2 {
			utils.PrintTestError(t, count, int64(2))
		}
		if len(users) != 2 {
			utils.PrintTestError(t, len(users), 2)
		}
	}
}
