package handlers

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"testing"
)

// An isolated editor may not save a receipt whose stored items are charged to a
// member they cannot see (the wholesale item replace would silently drop the
// charge). Visible charges and unrestricted callers are unaffected.
func TestEnforceReceiptChargedToPreservation(t *testing.T) {
	defer repositories.TruncateTestDb()
	services.ClearRolePermissionCacheForTests()

	db := repositories.GetDB()

	group := models.Group{Name: "iso-preserve", IsolateMembers: true}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	supRole := models.GroupRoleDefinition{Name: "iso-preserve-sup", SeesAllMembers: true}
	memRole := models.GroupRoleDefinition{Name: "iso-preserve-mem", SeesAllMembers: false}
	if err := db.Create(&supRole).Error; err != nil {
		t.Fatalf("seed sup role: %v", err)
	}
	if err := db.Create(&memRole).Error; err != nil {
		t.Fatalf("seed mem role: %v", err)
	}

	a := models.User{Username: "iso-pres-a", Password: "x"}
	coord := models.User{Username: "iso-pres-coord", Password: "x"}
	b := models.User{Username: "iso-pres-b", Password: "x"}
	for _, user := range []*models.User{&a, &coord, &b} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}
	members := []models.GroupMember{
		{GroupID: group.ID, UserID: a.ID, GroupRoleID: &memRole.ID},
		{GroupID: group.ID, UserID: coord.ID, GroupRoleID: &supRole.ID},
		{GroupID: group.ID, UserID: b.ID, GroupRoleID: &memRole.ID},
	}
	for i := range members {
		if err := db.Create(&members[i]).Error; err != nil {
			t.Fatalf("seed member: %v", err)
		}
	}

	// A stored item charged to B (non-visible peer) → editing is blocked for A.
	chargedToB := models.Receipt{ReceiptItems: []models.Item{{ChargedToUserId: &b.ID}}}
	allowed, message, err := enforceReceiptChargedToPreservation(a.ID, chargedToB)
	if err != nil {
		t.Fatalf("preservation check errored: %v", err)
	}
	if allowed {
		t.Fatalf("expected the edit to be blocked when a stored item is charged to a non-visible member")
	}
	if message == "" {
		t.Fatalf("expected a non-empty deny message")
	}

	// A stored item charged to the coordinator (visible) → allowed.
	chargedToCoord := models.Receipt{ReceiptItems: []models.Item{{ChargedToUserId: &coord.ID}}}
	allowed, _, err = enforceReceiptChargedToPreservation(a.ID, chargedToCoord)
	if err != nil {
		t.Fatalf("preservation check errored: %v", err)
	}
	if !allowed {
		t.Fatalf("a charge to a visible member should be allowed")
	}

	// The supervisor is unrestricted → always allowed, even for the B charge.
	allowed, _, err = enforceReceiptChargedToPreservation(coord.ID, chargedToB)
	if err != nil {
		t.Fatalf("preservation check errored: %v", err)
	}
	if !allowed {
		t.Fatalf("an unrestricted caller should never be blocked")
	}
}
