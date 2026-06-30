package services

import (
	"errors"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// seedReceipt persists a receipt (no associations) paid by the given user.
func seedReceipt(t *testing.T, name string, groupId uint, paidByUserId uint) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(10),
		Date:         time.Now(),
		Status:       models.OPEN,
		GroupId:      groupId,
		PaidByUserID: paidByUserId,
	}
	if err := repositories.GetDB().Omit(clause.Associations).Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt
}

// seedReceiptMember creates a group, a group role with the given permissions and
// grants, and a member assigned to it. Returns the member's user id and group id.
func seedReceiptMember(
	t *testing.T,
	username string,
	roleName string,
	perms []string,
	categoryGrants []uint,
	tagGrants []uint,
	paidByGrants []uint,
	includeOwn bool,
) (uint, uint) {
	t.Helper()
	db := repositories.GetDB()

	group := models.Group{Name: "grp-" + username}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(roleName, "", perms, categoryGrants, tagGrants, paidByGrants, includeOwn)
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

	ClearRolePermissionCacheForTests()
	ClearGroupRoleGrantCacheForTests()
	return user.ID, group.ID
}

// giveAppRole assigns an app role with the given permissions to a user.
func giveAppRole(t *testing.T, userId uint, roleName string, perms []string) {
	t.Helper()
	role, err := repositories.NewRoleRepository(nil).CreateAppRole(roleName, "", perms)
	if err != nil {
		t.Fatalf("create app role: %v", err)
	}
	if err := repositories.GetDB().Model(&models.User{}).Where("id = ?", userId).Update("app_role_id", role.ID).Error; err != nil {
		t.Fatalf("assign app role: %v", err)
	}
	ClearRolePermissionCacheForTests()
}

func TestGetReceiptForUserReturnsReceiptForMember(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedReceiptMember(t, "getok", "getok-role", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	receipt := seedReceipt(t, "Lunch", groupId, userId)

	got, err := NewReceiptService(nil).GetReceiptForUser(userId, utils.UintToString(receipt.ID))
	if err != nil {
		t.Fatalf("GetReceiptForUser returned error: %v", err)
	}
	if got.ID != receipt.ID {
		t.Errorf("expected receipt %d, got %d", receipt.ID, got.ID)
	}
}

func TestGetReceiptForUserDeniesNonMember(t *testing.T) {
	defer repositories.TruncateTestDb()

	ownerId, groupId := seedReceiptMember(t, "owner", "owner-role", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	receipt := seedReceipt(t, "Lunch", groupId, ownerId)
	outsider := makeUser(t, "outsider")

	_, err := NewReceiptService(nil).GetReceiptForUser(outsider, utils.UintToString(receipt.ID))
	if !errors.Is(err, ErrReceiptAccessDenied) {
		t.Errorf("expected ErrReceiptAccessDenied for a non-member, got %v", err)
	}
}

func TestGetReceiptForUserDeniesPaidByHidden(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedPayer := makeUser(t, "allowedpayer")
	hiddenPayer := makeUser(t, "hiddenpayer")
	// Restrict the member to receipts paid by allowedPayer only.
	userId, groupId := seedReceiptMember(t, "pbviewer", "pb-role", []string{permissions.GroupReceiptsRead}, nil, nil, []uint{allowedPayer}, false)

	visibleReceipt := seedReceipt(t, "allowed", groupId, allowedPayer)
	hiddenReceipt := seedReceipt(t, "hidden", groupId, hiddenPayer)

	if _, err := NewReceiptService(nil).GetReceiptForUser(userId, utils.UintToString(visibleReceipt.ID)); err != nil {
		t.Fatalf("expected the allowed receipt to be visible: %v", err)
	}
	if _, err := NewReceiptService(nil).GetReceiptForUser(userId, utils.UintToString(hiddenReceipt.ID)); !errors.Is(err, ErrReceiptAccessDenied) {
		t.Errorf("expected ErrReceiptAccessDenied for a paid-by-hidden receipt, got %v", err)
	}
}

func TestGetReceiptForUserStripsCategoriesToGrants(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	allowed := models.Category{Name: "allowed-cat"}
	hidden := models.Category{Name: "hidden-cat"}
	if err := db.Create(&allowed).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	if err := db.Create(&hidden).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}

	userId, groupId := seedReceiptMember(t, "stripper", "strip-role", []string{permissions.GroupReceiptsRead}, []uint{allowed.ID}, nil, nil, false)
	receipt := seedReceipt(t, "with-cats", groupId, userId)
	if err := db.Model(&receipt).Association("Categories").Append([]models.Category{allowed, hidden}); err != nil {
		t.Fatalf("attach categories: %v", err)
	}

	got, err := NewReceiptService(nil).GetReceiptForUser(userId, utils.UintToString(receipt.ID))
	if err != nil {
		t.Fatalf("GetReceiptForUser returned error: %v", err)
	}
	if len(got.Categories) != 1 || got.Categories[0].ID != allowed.ID {
		t.Errorf("expected only the granted category, got %+v", got.Categories)
	}
}

func TestGetReceiptForUserDeniesMissingReceipt(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId := makeUser(t, "any")
	if _, err := NewReceiptService(nil).GetReceiptForUser(userId, "999999"); !errors.Is(err, ErrReceiptAccessDenied) {
		t.Errorf("expected ErrReceiptAccessDenied for a missing receipt, got %v", err)
	}
}

func TestSearchReceiptsForUserRequiresSearchPermission(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId := makeUser(t, "nosearch")
	if _, err := NewReceiptService(nil).SearchReceiptsForUser(userId, "coffee", 100); !errors.Is(err, ErrSearchForbidden) {
		t.Errorf("expected ErrSearchForbidden without app.receipts.search, got %v", err)
	}
}

func TestSearchReceiptsForUserScopesAndAppliesPaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()

	allowedPayer := makeUser(t, "spallowed")
	hiddenPayer := makeUser(t, "sphidden")
	userId, groupId := seedReceiptMember(t, "spsearcher", "sp-role", []string{permissions.GroupReceiptsRead}, nil, nil, []uint{allowedPayer}, false)
	giveAppRole(t, userId, "sp-app-role", []string{permissions.AppReceiptsSearch})

	// A receipt in another group the user does not belong to must be excluded by scope.
	otherGroup := models.Group{Name: "other"}
	if err := repositories.GetDB().Create(&otherGroup).Error; err != nil {
		t.Fatalf("create other group: %v", err)
	}

	seedReceipt(t, "Coffee allowed", groupId, allowedPayer)
	seedReceipt(t, "Coffee hidden", groupId, hiddenPayer)
	seedReceipt(t, "Coffee elsewhere", otherGroup.ID, allowedPayer)

	results, err := NewReceiptService(nil).SearchReceiptsForUser(userId, "Coffee", 100)
	if err != nil {
		t.Fatalf("SearchReceiptsForUser returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected exactly 1 visible result, got %d (%+v)", len(results), results)
	}
	if results[0].GroupID != groupId || results[0].PaidByUserId != allowedPayer {
		t.Errorf("expected the allowed-payer receipt in the member's group, got %+v", results[0])
	}
}

func TestSearchReceiptsForUserBlankQueryReturnsEmpty(t *testing.T) {
	defer repositories.TruncateTestDb()

	userId, groupId := seedReceiptMember(t, "blank", "blank-role", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false)
	giveAppRole(t, userId, "blank-app-role", []string{permissions.AppReceiptsSearch})
	seedReceipt(t, "Coffee", groupId, userId)

	results, err := NewReceiptService(nil).SearchReceiptsForUser(userId, "   ", 100)
	if err != nil {
		t.Fatalf("SearchReceiptsForUser returned error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected a blank query to return no results, got %d", len(results))
	}
}
