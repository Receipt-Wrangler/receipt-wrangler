package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

// makeUser inserts a user and returns its id (used as a "paid by" payer).
func makeUser(t *testing.T, username string) uint {
	t.Helper()
	user := models.User{Username: username, Password: "password"}
	if err := repositories.GetDB().Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user.ID
}

// seedMemberWithPaidByRole creates a group, a group role granting receipts.read
// with the given paid-by config, and a member assigned to it. Returns the member
// user id, group id, and role id.
func seedMemberWithPaidByRole(t *testing.T, username string, paidByUserGrantIds []uint, includeOwn bool) (uint, uint, uint) {
	t.Helper()
	db := repositories.GetDB()

	group := models.Group{Name: "pb-group-" + username}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole("PaidBy Role "+username, "", []string{permissions.GroupReceiptsRead}, nil, nil, paidByUserGrantIds, includeOwn, false)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: username, Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed member user: %v", err)
	}

	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed group member: %v", err)
	}

	return user.ID, group.ID, role.ID
}

func TestGetGroupPaidByUserIdsUnrestrictedWhenNoConfig(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-open", nil, false)
	service := NewPermissionService(nil)

	allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupPaidByUserIdsForUser: %v", err)
	}
	if !unrestricted {
		t.Error("expected unrestricted when role has no paid-by grants and no include-own")
	}
	if allowed != nil {
		t.Errorf("expected nil allowed set when unrestricted, got %v", allowed)
	}
}

func TestGetGroupPaidByUserIdsSelfOnly(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-self", nil, true)
	service := NewPermissionService(nil)

	allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupPaidByUserIdsForUser: %v", err)
	}
	if unrestricted {
		t.Error("expected restricted when include-own is set")
	}
	if len(allowed) != 1 {
		t.Fatalf("expected exactly the member's own id, got %v", allowed)
	}
	if _, ok := allowed[userId]; !ok {
		t.Errorf("expected own id %d in allowed set %v", userId, allowed)
	}
}

func TestGetGroupPaidByUserIdsUsersOnlyExcludesSelf(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-payer")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-reviewer", []uint{payer}, false)
	service := NewPermissionService(nil)

	allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupPaidByUserIdsForUser: %v", err)
	}
	if unrestricted {
		t.Error("expected restricted when a specific paid-by user is granted")
	}
	if _, ok := allowed[payer]; !ok {
		t.Errorf("expected granted payer %d in allowed set %v", payer, allowed)
	}
	// include-own is false, so the member's OWN receipts must NOT be visible — this
	// is the "pure reviewer" case.
	if _, ok := allowed[userId]; ok {
		t.Errorf("expected own id %d NOT in allowed set %v when include-own is false", userId, allowed)
	}
	if len(allowed) != 1 {
		t.Fatalf("expected only the granted payer, got %v", allowed)
	}
}

func TestGetGroupPaidByUserIdsSelfPlusUser(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-payer2")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-self-plus", []uint{payer}, true)
	service := NewPermissionService(nil)

	allowed, _, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("GetGroupPaidByUserIdsForUser: %v", err)
	}
	if len(allowed) != 2 {
		t.Fatalf("expected own id + granted payer, got %v", allowed)
	}
	if _, ok := allowed[userId]; !ok {
		t.Errorf("expected own id %d in allowed set %v", userId, allowed)
	}
	if _, ok := allowed[payer]; !ok {
		t.Errorf("expected granted payer %d in allowed set %v", payer, allowed)
	}
}

func TestGetGroupPaidByUserIdsNonMemberUnrestricted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-payer3")
	_, groupId, _ := seedMemberWithPaidByRole(t, "pb-member", []uint{payer}, false)

	outsider := makeUser(t, "pb-outsider")
	service := NewPermissionService(nil)

	_, unrestricted, err := service.GetGroupPaidByUserIdsForUser(outsider, groupId)
	if err != nil {
		t.Fatalf("GetGroupPaidByUserIdsForUser: %v", err)
	}
	if !unrestricted {
		t.Error("expected a non-member to resolve as unrestricted")
	}
}

// TestGetGroupPaidByUserIdsDoesNotMutateCacheAcrossUsers guards the relative
// "their own" union: two members share one role-keyed cache entry, so resolving
// for one member must not leak that member's id into the shared cache or the
// other member's result.
func TestGetGroupPaidByUserIdsDoesNotMutateCacheAcrossUsers(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	db := repositories.GetDB()
	payer := makeUser(t, "pb-shared-payer")

	group := models.Group{Name: "pb-shared-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole("PaidBy Shared Role", "", []string{permissions.GroupReceiptsRead}, nil, nil, []uint{payer}, true, false)
	if err != nil {
		t.Fatalf("seed role: %v", err)
	}

	memberA := makeUser(t, "pb-member-a")
	memberB := makeUser(t, "pb-member-b")
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: memberA, GroupRoleID: &role.ID}).Error; err != nil {
		t.Fatalf("seed member A: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: memberB, GroupRoleID: &role.ID}).Error; err != nil {
		t.Fatalf("seed member B: %v", err)
	}

	service := NewPermissionService(nil)

	allowedA, _, err := service.GetGroupPaidByUserIdsForUser(memberA, group.ID)
	if err != nil {
		t.Fatalf("resolve A: %v", err)
	}
	allowedB, _, err := service.GetGroupPaidByUserIdsForUser(memberB, group.ID)
	if err != nil {
		t.Fatalf("resolve B: %v", err)
	}

	// Each member sees the granted payer plus only their OWN id.
	if _, ok := allowedA[memberA]; !ok {
		t.Errorf("member A should see own id; got %v", allowedA)
	}
	if _, ok := allowedA[memberB]; ok {
		t.Errorf("member A must NOT see member B's id; got %v", allowedA)
	}
	if _, ok := allowedB[memberB]; !ok {
		t.Errorf("member B should see own id; got %v", allowedB)
	}
	if _, ok := allowedB[memberA]; ok {
		t.Errorf("member B must NOT see member A's id (cache leak); got %v", allowedB)
	}
}

// TestGetGroupPaidByUserIdsFailsClosedAfterGrantedUserDeleted guards the
// privacy-widening fix: a role configured to see only a specific payer must STAY
// restricted (and resolve to "see nothing") after that payer is deleted and the
// FK cascade empties the grant rows — it must NOT silently widen to see-all.
func TestGetGroupPaidByUserIdsFailsClosedAfterGrantedUserDeleted(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-cascade-payer")
	userId, groupId, roleId := seedMemberWithPaidByRole(t, "pb-cascade-member", []uint{payer}, false)
	service := NewPermissionService(nil)

	// Initially restricted to the granted payer.
	allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("initial resolve: %v", err)
	}
	if unrestricted {
		t.Fatal("expected restricted before the payer is deleted")
	}
	if _, ok := allowed[payer]; !ok || len(allowed) != 1 {
		t.Fatalf("expected allowed {payer}, got %v", allowed)
	}

	// Delete the granted user; the FK cascade removes the role's paid-by grant row.
	if err := repositories.GetDB().Delete(&models.User{}, payer).Error; err != nil {
		t.Fatalf("delete payer: %v", err)
	}
	// A user delete does not evict the role's grant cache, so force a fresh load.
	clearGroupRoleGrantCache(roleId)

	allowed, unrestricted, err = service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("re-resolve: %v", err)
	}
	if unrestricted {
		t.Error("a configured role whose only payer was deleted must stay restricted (fail closed), not widen to see-all")
	}
	if len(allowed) != 0 {
		t.Errorf("expected an empty allowed set (see nothing), got %v", allowed)
	}
}

// TestGetGroupPaidByUserIdsFailsClosedWhenFlagDesynced guards the hardening: even
// if the persisted PaidByVisibilityRestricted flag is out of sync (false) while
// grant rows exist, the resolver still treats the role as restricted (fail closed)
// rather than widening to see-all.
func TestGetGroupPaidByUserIdsFailsClosedWhenFlagDesynced(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-desync-payer")
	userId, groupId, roleId := seedMemberWithPaidByRole(t, "pb-desync-member", nil, false)

	// Simulate a desync: a grant row exists but the restricted flag was left false
	// (e.g. a future write path that inserted grants without recomputing the flag).
	db := repositories.GetDB()
	if err := db.Create(&models.GroupRolePaidByUserGrant{GroupRoleID: roleId, UserID: payer}).Error; err != nil {
		t.Fatalf("insert grant row: %v", err)
	}
	if err := db.Model(&models.GroupRoleDefinition{}).Where("id = ?", roleId).
		Update("paid_by_visibility_restricted", false).Error; err != nil {
		t.Fatalf("force flag false: %v", err)
	}
	clearGroupRoleGrantCache(roleId)

	allowed, unrestricted, err := NewPermissionService(nil).GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if unrestricted {
		t.Error("a role with live grant rows must resolve restricted even if the persisted flag is false")
	}
	if _, ok := allowed[payer]; !ok || len(allowed) != 1 {
		t.Errorf("expected allowed {payer}, got %v", allowed)
	}
}

func TestReceiptPaidByVisible(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-vis-payer")
	other := makeUser(t, "pb-vis-other")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-vis-member", []uint{payer}, false)
	service := NewPermissionService(nil)

	visible, err := service.ReceiptPaidByVisible(userId, groupId, payer)
	if err != nil {
		t.Fatalf("visible(payer): %v", err)
	}
	if !visible {
		t.Error("expected a receipt paid by the granted payer to be visible")
	}

	visible, err = service.ReceiptPaidByVisible(userId, groupId, other)
	if err != nil {
		t.Fatalf("visible(other): %v", err)
	}
	if visible {
		t.Error("expected a receipt paid by a non-granted user to be hidden")
	}
}

func TestFilterReceiptsByPaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-filter-payer")
	other := makeUser(t, "pb-filter-other")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-filter-member", []uint{payer}, false)
	service := NewPermissionService(nil)

	receipts := []models.Receipt{
		{GroupId: groupId, PaidByUserID: payer},
		{GroupId: groupId, PaidByUserID: other},
		{GroupId: groupId, PaidByUserID: payer},
	}

	filtered, err := service.FilterReceiptsByPaidBy(userId, receipts)
	if err != nil {
		t.Fatalf("FilterReceiptsByPaidBy: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 receipts (the payer's), got %d", len(filtered))
	}
	for _, receipt := range filtered {
		if receipt.PaidByUserID != payer {
			t.Errorf("unexpected hidden payer leaked: %d", receipt.PaidByUserID)
		}
	}
}

func TestPaidByListResolver(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	payer := makeUser(t, "pb-resolver-payer")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-resolver-member", []uint{payer}, true)
	service := NewPermissionService(nil)

	resolver := service.PaidByListResolver(userId)

	allowed, unrestricted, err := resolver(groupId)
	if err != nil {
		t.Fatalf("resolver(restricted group): %v", err)
	}
	if unrestricted {
		t.Error("expected restricted for a configured group role")
	}
	got := make(map[uint]struct{}, len(allowed))
	for _, id := range allowed {
		got[id] = struct{}{}
	}
	if _, ok := got[payer]; !ok {
		t.Errorf("expected granted payer %d in resolver result %v", payer, allowed)
	}
	if _, ok := got[userId]; !ok {
		t.Errorf("expected own id %d in resolver result %v", userId, allowed)
	}
	if len(allowed) != 2 {
		t.Fatalf("expected own id + payer, got %v", allowed)
	}

	// A group the user is not a member of has no role, so the resolver is
	// unrestricted (no paid-by constraint added to the query for that group).
	otherGroup := models.Group{Name: "pb-resolver-other"}
	if err := repositories.GetDB().Create(&otherGroup).Error; err != nil {
		t.Fatalf("seed other group: %v", err)
	}
	_, unrestricted, err = resolver(otherGroup.ID)
	if err != nil {
		t.Fatalf("resolver(non-member group): %v", err)
	}
	if !unrestricted {
		t.Error("expected unrestricted for a group the user is not a member of")
	}
}

func TestFilterReceiptsByPaidByUnrestrictedPassThrough(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()

	other := makeUser(t, "pb-pass-other")
	userId, groupId, _ := seedMemberWithPaidByRole(t, "pb-pass-member", nil, false)
	service := NewPermissionService(nil)

	receipts := []models.Receipt{
		{GroupId: groupId, PaidByUserID: userId},
		{GroupId: groupId, PaidByUserID: other},
	}

	filtered, err := service.FilterReceiptsByPaidBy(userId, receipts)
	if err != nil {
		t.Fatalf("FilterReceiptsByPaidBy: %v", err)
	}
	if len(filtered) != 2 {
		t.Errorf("expected pass-through of all receipts when unrestricted, got %d", len(filtered))
	}
}
