package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

// These tests cover GHSA-89mm-9qfv-cjg3: a group member with group.update
// rewriting the member roster (including their own group role) to escalate to
// owner or evict the legitimate owner.

func groupMemberRoleId(t *testing.T, userId uint, groupId uint) *uint {
	t.Helper()
	var member models.GroupMember
	if err := repositories.GetDB().
		Where("user_id = ? AND group_id = ?", userId, groupId).
		First(&member).Error; err != nil {
		t.Fatalf("load group member (%d,%d): %v", userId, groupId, err)
	}
	return member.GroupRoleID
}

func uintPtrEqual(a, b *uint) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func callUpdateGroup(callerId uint, groupId uint, body commands.UpsertGroupCommand) *http.Response {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/api", strings.NewReader(string(b)))

	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("groupId", utils.UintToString(groupId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiContext))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: callerId}}))

	UpdateGroup(w, r)
	return w.Result()
}

// setUpEscalationGroup creates group 1 with an owner (user 1) and a low-privilege
// member (user 2) whose custom role grants only group.update / group.view /
// group.receipts.read — the exact precondition described in the advisory.
func setUpEscalationGroup(t *testing.T) (ownerRoleId, attackerRoleId *uint) {
	t.Helper()
	repositories.GetDB().Create(&models.Group{})
	grantAllGroupPerms(t, 1, 1)
	grantGroupPerms(t, 2, 1,
		permissions.GroupUpdate,
		permissions.GroupView,
		permissions.GroupReceiptsRead,
	)
	return groupMemberRoleId(t, 1, 1), groupMemberRoleId(t, 2, 1)
}

func TestUpdateGroupBlocksMemberSelfEscalation(t *testing.T) {
	defer tearDownGroupTests()
	ownerRoleId, attackerRoleId := setUpEscalationGroup(t)

	// The attacker replays the roster, promoting their own row to the owner role.
	resp := callUpdateGroup(2, 1, commands.UpsertGroupCommand{
		Name:   "attacker-renamed",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: 1, GroupID: 1, GroupRoleID: ownerRoleId},
			{UserID: 2, GroupID: 1, GroupRoleID: ownerRoleId},
		},
	})

	if resp.StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, resp.StatusCode, http.StatusForbidden)
	}
	if got := groupMemberRoleId(t, 2, 1); !uintPtrEqual(got, attackerRoleId) {
		t.Fatalf("attacker group role was changed despite the 403")
	}
}

func TestUpdateGroupBlocksOwnerEviction(t *testing.T) {
	defer tearDownGroupTests()
	_, _ = setUpEscalationGroup(t)

	// The attacker submits a roster that drops the owner entirely.
	resp := callUpdateGroup(2, 1, commands.UpsertGroupCommand{
		Name:   "attacker-renamed",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: 2, GroupID: 1, GroupRoleID: groupMemberRoleId(t, 2, 1)},
		},
	})

	if resp.StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, resp.StatusCode, http.StatusForbidden)
	}
	var ownerCount int64
	repositories.GetDB().Model(&models.GroupMember{}).
		Where("user_id = ? AND group_id = ?", 1, 1).Count(&ownerCount)
	if ownerCount != 1 {
		t.Fatalf("owner membership was removed despite the 403 (count=%d)", ownerCount)
	}
}

func TestUpdateGroupAllowsLowPrivilegeMemberToRenameWithUnchangedRoster(t *testing.T) {
	defer tearDownGroupTests()
	ownerRoleId, attackerRoleId := setUpEscalationGroup(t)

	// A group.update member renames the group while submitting the roster
	// unchanged — no member diff, so the guard must not block it.
	resp := callUpdateGroup(2, 1, commands.UpsertGroupCommand{
		Name:   "renamed-by-member",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: 1, GroupID: 1, GroupRoleID: ownerRoleId},
			{UserID: 2, GroupID: 1, GroupRoleID: attackerRoleId},
		},
	})

	if resp.StatusCode != http.StatusOK {
		utils.PrintTestError(t, resp.StatusCode, http.StatusOK)
	}
}

func TestUpdateGroupAllowsOwnerToReassignMemberRole(t *testing.T) {
	defer tearDownGroupTests()
	ownerRoleId, _ := setUpEscalationGroup(t)

	// The owner (holds every group permission, including group.members.*) promotes
	// the member to the owner role — a legitimate, in-privilege change.
	resp := callUpdateGroup(1, 1, commands.UpsertGroupCommand{
		Name:   "owner-managed",
		Status: models.GROUP_ACTIVE,
		GroupMembers: []commands.UpsertGroupMemberCommand{
			{UserID: 1, GroupID: 1, GroupRoleID: ownerRoleId},
			{UserID: 2, GroupID: 1, GroupRoleID: ownerRoleId},
		},
	})

	if resp.StatusCode != http.StatusOK {
		utils.PrintTestError(t, resp.StatusCode, http.StatusOK)
	}
	if got := groupMemberRoleId(t, 2, 1); !uintPtrEqual(got, ownerRoleId) {
		t.Fatalf("owner-driven role reassignment did not persist")
	}
}
