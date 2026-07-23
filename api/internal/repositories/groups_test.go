package repositories

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func setUpGroupTest() {
	CreateTestUser()
}

func setupGroupRepository() GroupRepository {
	return NewGroupRepository(nil)
}

func teardownGroupTest() {
	TruncateTestDb()
}

func TestShouldCreateGroupSuccessfully(t *testing.T) {
	defer teardownGroupTest()
	groupToCreate := commands.UpsertGroupCommand{Name: "test"}
	setUpGroupTest()
	groupRepository := setupGroupRepository()
	createdGroup, err := groupRepository.CreateGroup(groupToCreate, 1)

	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
	}

	if createdGroup.ID != 1 {
		utils.PrintTestError(t, createdGroup.ID, "1")
	}
	if createdGroup.Name != "test" {
		utils.PrintTestError(t, createdGroup.Name, "test")
	}
	if createdGroup.Status != models.GROUP_ACTIVE {
		utils.PrintTestError(t, createdGroup.Status, "Active")
	}
	if len(createdGroup.GroupMembers) != 1 {
		utils.PrintTestError(t, createdGroup.GroupMembers, "1")
	}
	if createdGroup.GroupMembers[0].UserID != 1 {
		utils.PrintTestError(t, createdGroup.GroupMembers[0].UserID, "1")
	}
}

func TestCreateGroupAssignsDefaultGroupRole(t *testing.T) {
	defer teardownGroupTest()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	CreateTestUser()

	groupRepository := setupGroupRepository()
	created, err := groupRepository.CreateGroup(commands.UpsertGroupCommand{Name: "test"}, 1)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}

	if len(created.GroupMembers) != 1 {
		utils.PrintTestError(t, len(created.GroupMembers), 1)
		return
	}

	member := created.GroupMembers[0]
	roleRepository := NewRoleRepository(nil)
	defaultId, err := roleRepository.GetDefaultGroupRoleId()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if member.GroupRoleID == nil || defaultId == nil || *member.GroupRoleID != *defaultId {
		utils.PrintTestError(t, member.GroupRoleID, defaultId)
	}
}

func TestCreateGroupLeavesGroupRoleNilWhenUnseeded(t *testing.T) {
	defer teardownGroupTest()
	CreateTestUser()

	groupRepository := setupGroupRepository()
	created, err := groupRepository.CreateGroup(commands.UpsertGroupCommand{Name: "test"}, 1)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}

	if len(created.GroupMembers) != 1 {
		utils.PrintTestError(t, len(created.GroupMembers), 1)
		return
	}
	if created.GroupMembers[0].GroupRoleID != nil {
		utils.PrintTestError(t, created.GroupMembers[0].GroupRoleID, nil)
	}
}

func TestCreateGroupHonorsMemberGroupRoleId(t *testing.T) {
	defer teardownGroupTest()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	CreateTestUser() // creator, id 1

	db := GetDB()
	member := models.User{Username: "member", DisplayName: "m", Password: "p"}
	db.Create(&member)

	var editor models.GroupRoleDefinition
	if err := db.Select("id").Where("name = ?", LegacyEditorRoleName).First(&editor).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	groupRepository := setupGroupRepository()
	created, err := groupRepository.CreateGroup(commands.UpsertGroupCommand{
		Name: "test",
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: member.ID, GroupRoleID: &editor.ID},
		},
	}, 1)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}

	// The added member (not the creator/owner) must carry the chosen modern group
	// role on its FK.
	var added *models.GroupMember
	for i := range created.GroupMembers {
		if created.GroupMembers[i].UserID == member.ID {
			added = &created.GroupMembers[i]
		}
	}
	if added == nil {
		utils.PrintTestError(t, "member not found", "member present")
		return
	}
	if added.GroupRoleID == nil || *added.GroupRoleID != editor.ID {
		utils.PrintTestError(t, added.GroupRoleID, editor.ID)
	}
}

func TestUpdateGroupHonorsMemberGroupRoleId(t *testing.T) {
	defer teardownGroupTest()

	if err := SeedSystemRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if err := EnsureDefaultRoles(); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	CreateTestUser() // creator, id 1

	db := GetDB()
	member := models.User{Username: "member", DisplayName: "m", Password: "p"}
	db.Create(&member)

	groupRepository := setupGroupRepository()
	created, err := groupRepository.CreateGroup(commands.UpsertGroupCommand{Name: "test"}, 1)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}

	var viewer models.GroupRoleDefinition
	if err := db.Select("id").Where("name = ?", LegacyViewerRoleName).First(&viewer).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	_, err = groupRepository.UpdateGroup(commands.UpsertGroupCommand{
		Name:   "test",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: member.ID, GroupID: created.ID, GroupRoleID: &viewer.ID},
		},
	}, utils.UintToString(created.ID))
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}

	var stored models.GroupMember
	if err := db.Where("group_id = ? AND user_id = ?", created.ID, member.ID).First(&stored).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if stored.GroupRoleID == nil || *stored.GroupRoleID != viewer.ID {
		utils.PrintTestError(t, stored.GroupRoleID, viewer.ID)
	}
}

func TestCreateGroupPersistsIsolateMembers(t *testing.T) {
	defer teardownGroupTest()
	setUpGroupTest()
	groupRepository := setupGroupRepository()

	created, err := groupRepository.CreateGroup(commands.UpsertGroupCommand{Name: "iso", IsolateMembers: true}, 1)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}
	if !created.IsolateMembers {
		utils.PrintTestError(t, created.IsolateMembers, true)
	}

	reloaded, err := groupRepository.GetGroupById(utils.UintToString(created.ID), true, false, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !reloaded.IsolateMembers {
		utils.PrintTestError(t, reloaded.IsolateMembers, true)
	}
}

func TestUpdateGroupPersistsIsolateMembersToggle(t *testing.T) {
	defer teardownGroupTest()
	setUpGroupTest()
	groupRepository := setupGroupRepository()

	created, err := groupRepository.CreateGroup(commands.UpsertGroupCommand{Name: "iso", IsolateMembers: true}, 1)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}
	groupId := utils.UintToString(created.ID)
	members := []commands.UpsertGroupMemberCommand{{UserID: 1, GroupID: created.ID}}

	// Toggle isolation OFF — the false value must persist (a struct Updates would
	// skip the zero-value bool, leaving isolation stuck on).
	_, err = groupRepository.UpdateGroup(commands.UpsertGroupCommand{
		Name: "iso", Status: models.GROUP_ACTIVE, GroupMembers: members, IsolateMembers: false,
	}, groupId)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}
	reloaded, err := groupRepository.GetGroupById(groupId, true, false, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if reloaded.IsolateMembers {
		utils.PrintTestError(t, reloaded.IsolateMembers, false)
	}

	// Toggle isolation back ON.
	_, err = groupRepository.UpdateGroup(commands.UpsertGroupCommand{
		Name: "iso", Status: models.GROUP_ACTIVE, GroupMembers: members, IsolateMembers: true,
	}, groupId)
	if err != nil {
		utils.PrintTestError(t, err, "Expected no error")
		return
	}
	reloaded, err = groupRepository.GetGroupById(groupId, true, false, false)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if !reloaded.IsolateMembers {
		utils.PrintTestError(t, reloaded.IsolateMembers, true)
	}
}

func TestShouldGetGroupById(t *testing.T) {
	defer teardownGroupTest()
	CreateTestGroup()
	setUpGroupTest()
	groupRepository := setupGroupRepository()
	testGroup, err := groupRepository.GetGroupById("1", false, true, true)

	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if testGroup.ID != 1 {
		utils.PrintTestError(t, err, "1")
	}
	if testGroup.Name != "test" {
		utils.PrintTestError(t, err, "1")
	}
}

func TestShouldGetAGroupWithGroupMembers(t *testing.T) {
	defer teardownGroupTest()
	CreateTestGroupWithUsers()
	groupRepository := setupGroupRepository()

	testGroup, err := groupRepository.GetGroupById("1", true, true, true)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if testGroup.ID != 1 {
		utils.PrintTestError(t, testGroup.ID, "1")
	}
	if len(testGroup.GroupMembers) != 3 {
		utils.PrintTestError(t, err, "3")
	}
}

func TestShouldReturnErrorIfGroupDoesNotExist(t *testing.T) {
	defer teardownGroupTest()
	groupRepository := setupGroupRepository()
	testGroup, err := groupRepository.GetGroupById("2332", false, true, true)

	if err == nil {
		utils.PrintTestError(t, err, "error")
	}
	if testGroup.ID != 0 {
		utils.PrintTestError(t, testGroup.ID, "0")
	}
}

func TestShouldUpdateGroup(t *testing.T) {
	defer teardownGroupTest()
	CreateTestGroup()
	updateGroup := commands.UpsertGroupCommand{Name: "new name", Status: models.GROUP_ARCHIVED}
	groupRepository := setupGroupRepository()
	updatedGroup, err := groupRepository.UpdateGroup(updateGroup, "1")

	if err != nil {
		utils.PrintTestError(t, err, "error")
	}
	if updatedGroup.Name != "new name" {
		utils.PrintTestError(t, updatedGroup.Name, "new name")
	}
	if updatedGroup.Status != models.GROUP_ARCHIVED {
		utils.PrintTestError(t, updatedGroup.Status, models.GROUP_ARCHIVED)
	}
}
