package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

func seedIsoUser(t *testing.T, username string) models.User {
	t.Helper()
	user := models.User{Username: username, Password: "password"}
	if err := repositories.GetDB().Create(&user).Error; err != nil {
		t.Fatalf("seed user %s: %v", username, err)
	}
	return user
}

func seedIsoGroup(t *testing.T, name string, isolate bool) models.Group {
	t.Helper()
	group := models.Group{Name: name, IsolateMembers: isolate}
	if err := repositories.GetDB().Create(&group).Error; err != nil {
		t.Fatalf("seed group %s: %v", name, err)
	}
	return group
}

func seedIsoRole(t *testing.T, name string, seesAll bool) models.GroupRoleDefinition {
	t.Helper()
	role := models.GroupRoleDefinition{Name: name, SeesAllMembers: seesAll}
	if err := repositories.GetDB().Create(&role).Error; err != nil {
		t.Fatalf("seed role %s: %v", name, err)
	}
	return role
}

func seedIsoMember(t *testing.T, groupId uint, userId uint, roleId *uint) {
	t.Helper()
	member := models.GroupMember{GroupID: groupId, UserID: userId, GroupRoleID: roleId}
	if err := repositories.GetDB().Create(&member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}
}

func assertVisible(t *testing.T, viewerId uint, wantIds []uint, wantUnrestricted bool) {
	t.Helper()
	service := NewPermissionService(nil)
	set, unrestricted, err := service.GetVisibleUserIdsForUser(viewerId)
	if err != nil {
		t.Fatalf("resolver error: %v", err)
	}
	if unrestricted != wantUnrestricted {
		t.Fatalf("unrestricted = %v, want %v (set=%v)", unrestricted, wantUnrestricted, set)
	}
	if wantUnrestricted {
		if set != nil {
			t.Fatalf("unrestricted set should be nil, got %v", set)
		}
		return
	}
	if len(set) != len(wantIds) {
		t.Fatalf("visible set = %v, want ids %v", set, wantIds)
	}
	for _, id := range wantIds {
		if _, ok := set[id]; !ok {
			t.Errorf("expected user %d visible, set = %v", id, set)
		}
	}
}

// A viewer with only non-isolated memberships is unrestricted (backward compatible).
func TestVisibleUsers_NotIsolated_Unrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	group := seedIsoGroup(t, "iso-normal", false)
	a := seedIsoUser(t, "iso-a")
	b := seedIsoUser(t, "iso-b")
	seedIsoMember(t, group.ID, a.ID, nil)
	seedIsoMember(t, group.ID, b.ID, nil)

	assertVisible(t, a.ID, nil, true)
}

// An isolated (non-supervisor) member sees only themselves + supervisors, never peers.
func TestVisibleUsers_IsolatedMember_SelfAndSupervisorsOnly(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	group := seedIsoGroup(t, "iso-home", true)
	supRole := seedIsoRole(t, "iso-supervisor", true)
	memberRole := seedIsoRole(t, "iso-member", false)

	coord := seedIsoUser(t, "iso-coord")
	p1 := seedIsoUser(t, "iso-parent1")
	p2 := seedIsoUser(t, "iso-parent2")
	seedIsoMember(t, group.ID, coord.ID, &supRole.ID)
	seedIsoMember(t, group.ID, p1.ID, &memberRole.ID)
	seedIsoMember(t, group.ID, p2.ID, &memberRole.ID)

	// p1 sees self + coordinator, but NOT p2 (asymmetric peer invisibility).
	assertVisible(t, p1.ID, []uint{p1.ID, coord.ID}, false)
	assertVisible(t, p2.ID, []uint{p2.ID, coord.ID}, false)
	// The supervisor is unrestricted (no isolated-restricted group of their own).
	assertVisible(t, coord.ID, nil, true)
}

// A viewer in a mix of isolated and non-isolated groups: sees supervisors of the
// isolated group + all members of the non-isolated group, never the isolated peer.
func TestVisibleUsers_MixedGroups(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	iso := seedIsoGroup(t, "iso-mixed-iso", true)
	normal := seedIsoGroup(t, "iso-mixed-normal", false)
	supRole := seedIsoRole(t, "iso-mixed-sup", true)
	memberRole := seedIsoRole(t, "iso-mixed-mem", false)

	viewer := seedIsoUser(t, "iso-mixed-viewer")
	coord := seedIsoUser(t, "iso-mixed-coord")
	otherIso := seedIsoUser(t, "iso-mixed-other")
	normalPeer := seedIsoUser(t, "iso-mixed-peer")

	seedIsoMember(t, iso.ID, viewer.ID, &memberRole.ID)
	seedIsoMember(t, iso.ID, coord.ID, &supRole.ID)
	seedIsoMember(t, iso.ID, otherIso.ID, &memberRole.ID)
	seedIsoMember(t, normal.ID, viewer.ID, nil)
	seedIsoMember(t, normal.ID, normalPeer.ID, nil)

	assertVisible(t, viewer.ID, []uint{viewer.ID, coord.ID, normalPeer.ID}, false)
}

// app.users.read short-circuits to unrestricted even for an isolated member.
func TestVisibleUsers_AppUsersReadShortCircuits(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	iso := seedIsoGroup(t, "iso-admin-group", true)
	memberRole := seedIsoRole(t, "iso-admin-mem", false)
	admin := seedIsoUser(t, "iso-admin-user")
	other := seedIsoUser(t, "iso-admin-other")
	seedIsoMember(t, iso.ID, admin.ID, &memberRole.ID)
	seedIsoMember(t, iso.ID, other.ID, &memberRole.ID)

	appRole := models.AppRole{
		Name:        "iso-app-admin-role",
		Permissions: []models.AppRolePermission{{Permission: permissions.AppUsersRead}},
	}
	if err := repositories.GetDB().Create(&appRole).Error; err != nil {
		t.Fatalf("seed app role: %v", err)
	}
	if err := repositories.GetDB().Model(&models.User{}).Where("id = ?", admin.ID).
		Update("app_role_id", appRole.ID).Error; err != nil {
		t.Fatalf("assign app role: %v", err)
	}

	assertVisible(t, admin.ID, nil, true)
}

// An isolated member cannot add (invite) a user they cannot see; a visible user,
// themselves, and any target for an unrestricted caller are allowed.
func TestAuthorizeAddedMembersVisibility(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	group := seedIsoGroup(t, "iso-invite", true)
	supRole := seedIsoRole(t, "iso-invite-sup", true)
	memberRole := seedIsoRole(t, "iso-invite-mem", false)
	a := seedIsoUser(t, "iso-invite-a")
	coord := seedIsoUser(t, "iso-invite-coord")
	b := seedIsoUser(t, "iso-invite-b")
	seedIsoMember(t, group.ID, a.ID, &memberRole.ID)
	seedIsoMember(t, group.ID, coord.ID, &supRole.ID)
	seedIsoMember(t, group.ID, b.ID, &memberRole.ID)

	service := NewGroupService(nil)

	// A (isolated) cannot invite B (a non-visible peer).
	if err := service.AuthorizeAddedMembersVisibility(a.ID, []commands.UpsertGroupMemberCommand{{UserID: b.ID}}); !errors.Is(err, ErrGroupMemberChangeForbidden) {
		t.Fatalf("inviting a non-visible user should be forbidden, got %v", err)
	}
	// A can invite the coordinator (visible) and add themselves.
	if err := service.AuthorizeAddedMembersVisibility(a.ID, []commands.UpsertGroupMemberCommand{{UserID: coord.ID}, {UserID: a.ID}}); err != nil {
		t.Fatalf("inviting a visible user + self should be allowed, got %v", err)
	}
	// The supervisor is unrestricted and may invite anyone.
	if err := service.AuthorizeAddedMembersVisibility(coord.ID, []commands.UpsertGroupMemberCommand{{UserID: b.ID}}); err != nil {
		t.Fatalf("an unrestricted caller should be able to invite anyone, got %v", err)
	}
}
