package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"

	"github.com/go-chi/chi/v5"
)

// DeleteGroup is authorized two ways: the group-scoped group.delete (a member
// removing their own group) or the app-scoped app.groups.delete (an
// administrator cleaning up a group anywhere in the system, including one they
// are not a member of). These tests pin both paths and their denials.

// seedDeletableGroup creates a plain group and returns its id.
func seedDeletableGroup(t *testing.T, name string) uint {
	t.Helper()
	group := models.Group{Name: name}
	if err := repositories.GetDB().Create(&group).Error; err != nil {
		t.Fatalf("create group: %v", err)
	}
	return group.ID
}

// deleteGroupAsUser drives the handler for userId against groupId.
func deleteGroupAsUser(userId uint, groupId uint) *httptest.ResponseRecorder {
	w, r := requestForUser(userId)

	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("groupId", utils.UintToString(groupId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiContext))

	DeleteGroup(w, r)

	return w
}

// assertGroupExists asserts whether the group row survived the request.
func assertGroupExists(t *testing.T, groupId uint, want bool) {
	t.Helper()
	var count int64
	if err := repositories.GetDB().Model(&models.Group{}).Where("id = ?", groupId).Count(&count).Error; err != nil {
		t.Fatalf("count groups: %v", err)
	}
	if (count > 0) != want {
		utils.PrintTestError(t, count > 0, want)
	}
}

func TestDeleteGroupAllowsMemberWithGroupDelete(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDeletableGroup(t, "Member Deletes Own Group")
	grantGroupPerms(t, 1, groupId, permissions.GroupDelete)

	assertStatus(t, deleteGroupAsUser(1, groupId), http.StatusOK)
	assertGroupExists(t, groupId, false)
}

func TestDeleteGroupAllowsNonMemberWithAppGroupsDelete(t *testing.T) {
	defer tearDownGroupTests()

	// The whole point of the permission: no membership in the target group, so
	// the group-scoped check cannot pass and the app-scoped fallback must.
	groupId := seedDeletableGroup(t, "Abandoned Group")
	grantAppPerms(t, 1, permissions.AppGroupsDelete)

	assertStatus(t, deleteGroupAsUser(1, groupId), http.StatusOK)
	assertGroupExists(t, groupId, false)
}

func TestDeleteGroupDeniesNonMemberWithoutAppGroupsDelete(t *testing.T) {
	defer tearDownGroupTests()

	// app.groups.read alone lets an admin SEE every group; it must not let them
	// delete one.
	groupId := seedDeletableGroup(t, "Someone Else's Group")
	grantAppPerms(t, 1, permissions.AppGroupsRead)

	assertStatus(t, deleteGroupAsUser(1, groupId), http.StatusForbidden)
	assertGroupExists(t, groupId, true)
}

func TestDeleteGroupDeniesMemberWithoutDeletePermission(t *testing.T) {
	defer tearDownGroupTests()

	groupId := seedDeletableGroup(t, "Viewer's Group")
	grantGroupPerms(t, 1, groupId, permissions.GroupView)

	assertStatus(t, deleteGroupAsUser(1, groupId), http.StatusForbidden)
	assertGroupExists(t, groupId, true)
}

func TestDeleteGroupRejectsAllGroupEvenForAppGroupsDelete(t *testing.T) {
	defer tearDownGroupTests()

	// The virtual "All" group is per-user infrastructure, not a real group, and
	// stays undeletable through this endpoint for every caller.
	allGroup := models.Group{Name: "All Receipts", IsAllGroup: true}
	if err := repositories.GetDB().Create(&allGroup).Error; err != nil {
		t.Fatalf("create all group: %v", err)
	}
	grantAppPerms(t, 1, permissions.AppGroupsDelete)

	assertStatus(t, deleteGroupAsUser(1, allGroup.ID), http.StatusBadRequest)
	assertGroupExists(t, allGroup.ID, true)
}
