package services

import (
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

func assertVisibleInGroup(t *testing.T, viewerId uint, groupId uint, wantIds []uint, wantUnrestricted bool) {
	t.Helper()
	service := NewPermissionService(nil)
	set, unrestricted, err := service.GetVisibleUserIdsForUserInGroup(viewerId, groupId)
	if err != nil {
		t.Fatalf("per-group resolver error: %v", err)
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

// The headline of the per-group rework: a viewer who shares an OPEN group with a peer
// still cannot see that peer INSIDE an isolated group they also share. The union
// resolver (the flat directory) does show the peer — the viewer legitimately knows them
// from the open group — but the per-group resolver for the isolated group does not.
// "Isolated means isolated."
func TestVisibleUsersInGroup_OpenGroupPeerHiddenInIsolatedGroup(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	iso := seedIsoGroup(t, "pg-iso", true)
	open := seedIsoGroup(t, "pg-open", false)
	supRole := seedIsoRole(t, "pg-sup", true)
	memberRole := seedIsoRole(t, "pg-mem", false)

	viewer := seedIsoUser(t, "pg-viewer")
	coord := seedIsoUser(t, "pg-coord")
	peer := seedIsoUser(t, "pg-peer")

	// viewer + peer share BOTH the isolated group and the open group; coord is the
	// isolated group's supervisor.
	seedIsoMember(t, iso.ID, viewer.ID, &memberRole.ID)
	seedIsoMember(t, iso.ID, coord.ID, &supRole.ID)
	seedIsoMember(t, iso.ID, peer.ID, &memberRole.ID)
	seedIsoMember(t, open.ID, viewer.ID, nil)
	seedIsoMember(t, open.ID, peer.ID, nil)

	// Per-group, INSIDE the isolated group: only self + supervisor. Peer is hidden
	// despite the shared open group.
	assertVisibleInGroup(t, viewer.ID, iso.ID, []uint{viewer.ID, coord.ID}, false)
	// Per-group, INSIDE the open group: everyone (non-isolated).
	assertVisibleInGroup(t, viewer.ID, open.ID, nil, true)
	// The union resolver (flat directory) still lists the peer — known via the open group.
	assertVisible(t, viewer.ID, []uint{viewer.ID, coord.ID, peer.ID}, false)
}

// A supervisor sees every member of the isolated group (unrestricted for that group).
func TestVisibleUsersInGroup_SupervisorUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	iso := seedIsoGroup(t, "pg-sup-iso", true)
	supRole := seedIsoRole(t, "pg-sup-role", true)
	memberRole := seedIsoRole(t, "pg-sup-mem", false)
	coord := seedIsoUser(t, "pg-sup-coord")
	member := seedIsoUser(t, "pg-sup-member")
	seedIsoMember(t, iso.ID, coord.ID, &supRole.ID)
	seedIsoMember(t, iso.ID, member.ID, &memberRole.ID)

	assertVisibleInGroup(t, coord.ID, iso.ID, nil, true)
}

// A non-isolated group is unrestricted for a plain member.
func TestVisibleUsersInGroup_NonIsolatedUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	open := seedIsoGroup(t, "pg-open-only", false)
	a := seedIsoUser(t, "pg-open-a")
	b := seedIsoUser(t, "pg-open-b")
	seedIsoMember(t, open.ID, a.ID, nil)
	seedIsoMember(t, open.ID, b.ID, nil)

	assertVisibleInGroup(t, a.ID, open.ID, nil, true)
}

// A non-member is unrestricted for a group they don't belong to: isolation only narrows
// for actual isolated members, and a non-member is kept out by the membership/permission
// gate, not by this filter (mirroring the paid-by resolver).
func TestVisibleUsersInGroup_NonMemberUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	iso := seedIsoGroup(t, "pg-nonmember-iso", true)
	memberRole := seedIsoRole(t, "pg-nonmember-mem", false)
	member := seedIsoUser(t, "pg-nonmember-member")
	outsider := seedIsoUser(t, "pg-nonmember-outsider")
	seedIsoMember(t, iso.ID, member.ID, &memberRole.ID)

	assertVisibleInGroup(t, outsider.ID, iso.ID, nil, true)
}

// app.users.read short-circuits the per-group resolver to unrestricted.
func TestVisibleUsersInGroup_AppUsersReadShortCircuits(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	iso := seedIsoGroup(t, "pg-admin-iso", true)
	memberRole := seedIsoRole(t, "pg-admin-mem", false)
	admin := seedIsoUser(t, "pg-admin-user")
	other := seedIsoUser(t, "pg-admin-other")
	seedIsoMember(t, iso.ID, admin.ID, &memberRole.ID)
	seedIsoMember(t, iso.ID, other.ID, &memberRole.ID)

	appRole := models.AppRole{
		Name:        "pg-app-admin-role",
		Permissions: []models.AppRolePermission{{Permission: permissions.AppUsersRead}},
	}
	if err := repositories.GetDB().Create(&appRole).Error; err != nil {
		t.Fatalf("seed app role: %v", err)
	}
	if err := repositories.GetDB().Model(&models.User{}).Where("id = ?", admin.ID).
		Update("app_role_id", appRole.ID).Error; err != nil {
		t.Fatalf("assign app role: %v", err)
	}

	assertVisibleInGroup(t, admin.ID, iso.ID, nil, true)
}
