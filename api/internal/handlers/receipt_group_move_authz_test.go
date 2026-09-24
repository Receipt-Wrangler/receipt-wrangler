package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
)

// seedGroupWithRole creates a group and a member holding a group role with the
// given permissions, returning (userId, groupId).
func seedGroupWithRole(t *testing.T, name string, perms []string) (uint, uint) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	group := models.Group{Name: name}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(name+" role", "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("seed role: %v", err)
	}
	user := models.User{Username: name + "-user", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}
	return user.ID, group.ID
}

// addMemberWithRole adds an existing user to groupId with a fresh role holding perms.
func addMemberWithRole(t *testing.T, userId uint, groupId uint, roleName string, perms []string) {
	t.Helper()
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(roleName, "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("seed role: %v", err)
	}
	if err := repositories.GetDB().Create(&models.GroupMember{GroupID: groupId, UserID: userId, GroupRoleID: &role.ID}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}
	services.ClearRolePermissionCacheForTests()
}

func moveBody(destGroupId uint, paidBy uint) string {
	return fmt.Sprintf(
		`{"name":"moved","amount":"1.00","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`,
		destGroupId, paidBy,
	)
}

// A user who can update receipts in the source group but has NO create right in
// the destination group cannot move a receipt there. Regression guard for the
// receipt-group-move authorization bug.
func TestUpdateReceipt_MoveToUnauthorizedGroupDenied(t *testing.T) {
	defer repositories.TruncateTestDb()

	fullPerms := []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate, permissions.GroupReceiptsCreate}
	userId, sourceGroup := seedGroupWithRole(t, "move-src", fullPerms)
	// The caller IS a member of the destination with read + update, but NOT
	// create. This isolates the required group.receipts.create check: a weaker
	// implementation that merely checked destination membership would still pass.
	_, destGroup := seedGroupWithRole(t, "move-dest", fullPerms)
	addMemberWithRole(t, userId, destGroup, "move-dest-member", []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate})

	receiptId := seedReceipt(t, sourceGroup, userId, 0)

	w, r := updateReceiptRequest(receiptId, userId, moveBody(destGroup, userId))
	UpdateReceipt(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}
	var reloaded models.Receipt
	repositories.GetDB().Select("group_id").First(&reloaded, receiptId)
	if reloaded.GroupId != sourceGroup {
		utils.PrintTestError(t, reloaded.GroupId, sourceGroup)
	}
}

// A user who can update in the source AND create in the destination may move the
// receipt (the legitimate path stays working).
func TestUpdateReceipt_MoveToAuthorizedGroupAllowed(t *testing.T) {
	defer repositories.TruncateTestDb()

	fullPerms := []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate, permissions.GroupReceiptsCreate}
	userId, sourceGroup := seedGroupWithRole(t, "move2-src", fullPerms)
	_, destGroup := seedGroupWithRole(t, "move2-dest", fullPerms)
	// Grant create in the destination WITHOUT update: moving a receipt in
	// requires group.receipts.create, so this proves create (not update) is what
	// authorizes the move.
	addMemberWithRole(t, userId, destGroup, "move2-dest-member", []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsCreate})

	receiptId := seedReceipt(t, sourceGroup, userId, 0)

	w, r := updateReceiptRequest(receiptId, userId, moveBody(destGroup, userId))
	UpdateReceipt(w, r)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}
	var reloaded models.Receipt
	repositories.GetDB().Select("group_id").First(&reloaded, receiptId)
	if reloaded.GroupId != destGroup {
		utils.PrintTestError(t, reloaded.GroupId, destGroup)
	}
}

// A normal same-group edit (no group change) is unaffected by the new check.
func TestUpdateReceipt_SameGroupEditUnaffected(t *testing.T) {
	defer repositories.TruncateTestDb()

	// Update WITHOUT create: a same-group edit must not require create, so this
	// also guards against an accidental create-check being run on every update.
	editPerms := []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate}
	userId, groupId := seedGroupWithRole(t, "move3-src", editPerms)
	receiptId := seedReceipt(t, groupId, userId, 0)

	w, r := updateReceiptRequest(receiptId, userId, moveBody(groupId, userId))
	UpdateReceipt(w, r)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
	}
}

// A receipt may never be moved into the synthetic "All" group, even by a caller
// whose All-group membership grants group.receipts.create — the All group is a
// cross-group view, not a real container. Without the destination All-group
// guard, that create permission would satisfy the move check and the receipt
// would be persisted into the All group under an unrestricted role.
func TestUpdateReceipt_MoveToAllGroupRejected(t *testing.T) {
	defer repositories.TruncateTestDb()

	fullPerms := []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate, permissions.GroupReceiptsCreate}
	userId, sourceGroup := seedGroupWithRole(t, "move-allgrp-src", fullPerms)

	// A real All-group row the caller belongs to with create rights — the exact
	// condition under which the destination create-check would otherwise pass.
	allGroup := models.Group{Name: "All", IsAllGroup: true}
	if err := repositories.GetDB().Create(&allGroup).Error; err != nil {
		t.Fatalf("seed all group: %v", err)
	}
	addMemberWithRole(t, userId, allGroup.ID, "move-allgrp-member", fullPerms)

	receiptId := seedReceipt(t, sourceGroup, userId, 0)

	w, r := updateReceiptRequest(receiptId, userId, moveBody(allGroup.ID, userId))
	UpdateReceipt(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusBadRequest)
	}
	var reloaded models.Receipt
	repositories.GetDB().Select("group_id").First(&reloaded, receiptId)
	if reloaded.GroupId != sourceGroup {
		utils.PrintTestError(t, reloaded.GroupId, sourceGroup)
	}
}
