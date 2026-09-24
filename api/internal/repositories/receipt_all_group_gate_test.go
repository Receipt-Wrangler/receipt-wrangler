package repositories

import (
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

// The All-group view must be gated by the caller's per-group read permission:
// a group whose readable resolver returns false is dropped entirely, so a member
// cannot read receipts in a group whose role denies receipt read.
func TestGetPagedReceipts_AllGroup_ReadGateDropsUnreadableGroups(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	group1 := models.Group{Name: "gate-g1"}
	group2 := models.Group{Name: "gate-g2"}
	allGroup := models.Group{Name: "gate-all", IsAllGroup: true}
	db.Create(&group1)
	db.Create(&group2)
	db.Create(&allGroup)

	member := models.User{Username: "gate-member", Password: "x"}
	db.Create(&member)
	db.Create(&models.GroupMember{GroupID: group1.ID, UserID: member.ID})
	db.Create(&models.GroupMember{GroupID: group2.ID, UserID: member.ID})

	mk := func(groupId uint, name string) {
		db.Create(&models.Receipt{Name: name, Amount: decimal.NewFromInt(1), Date: time.Now(), PaidByUserID: member.ID, GroupId: groupId, Status: models.OPEN})
	}
	mk(group1.ID, "g1-r1")
	mk(group2.ID, "g2-r1")

	// The caller can read group1 but NOT group2.
	readable := func(groupId uint) (bool, error) { return groupId == group1.ID, nil }

	repository := NewReceiptRepository(nil)
	receipts, count, err := repository.GetPagedReceiptsByGroupId(
		member.ID, utils.UintToString(allGroup.ID), pagedRequestAllReceipts(), nil, nil, nil, readable, nil,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 1 || len(receipts) != 1 {
		utils.PrintTestError(t, count, 1)
		return
	}
	if receipts[0].GroupId != group1.ID {
		utils.PrintTestError(t, receipts[0].GroupId, group1.ID)
	}
}

// The All-group view must apply a category filter PER GROUP against the caller's
// visible set: filtering by a category the caller cannot see in a group must not
// match that group's receipts (closing the filter-probe), while filtering by a
// visible category still works and an unfiltered read is unaffected.
func TestGetPagedReceipts_AllGroup_CategoryFilterIsPerGroup(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	group1 := models.Group{Name: "probe-g1"}
	allGroup := models.Group{Name: "probe-all", IsAllGroup: true}
	db.Create(&group1)
	db.Create(&allGroup)

	member := models.User{Username: "probe-member", Password: "x"}
	db.Create(&member)
	db.Create(&models.GroupMember{GroupID: group1.ID, UserID: member.ID})

	allowedCat := models.Category{Name: "probe-allowed"}
	secretCat := models.Category{Name: "probe-secret"}
	db.Create(&allowedCat)
	db.Create(&secretCat)

	mkWithCat := func(name string, cat models.Category) models.Receipt {
		receipt := models.Receipt{Name: name, Amount: decimal.NewFromInt(1), Date: time.Now(), PaidByUserID: member.ID, GroupId: group1.ID, Status: models.OPEN}
		db.Create(&receipt)
		db.Model(&receipt).Association("Categories").Append(&cat)
		return receipt
	}
	mkWithCat("has-secret", secretCat)
	mkWithCat("has-allowed", allowedCat)

	readable := func(groupId uint) (bool, error) { return true, nil }
	// The caller may see allowedCat but NOT secretCat in group1 (restricted).
	catTagVis := func(groupId uint) (CategoryTagVisibility, error) {
		return CategoryTagVisibility{
			CategoryAllowed:      map[uint]struct{}{allowedCat.ID: {}},
			CategoryUnrestricted: false,
			TagUnrestricted:      true,
		}, nil
	}

	repository := NewReceiptRepository(nil)

	filterBy := func(catId uint) commands.ReceiptPagedRequestCommand {
		req := pagedRequestAllReceipts()
		req.Filter.Categories = commands.PagedRequestField{Operation: commands.CONTAINS, Value: []interface{}{catId}}
		return req
	}

	// Probe by the restricted category -> must return nothing.
	_, count, err := repository.GetPagedReceiptsByGroupId(
		member.ID, utils.UintToString(allGroup.ID), filterBy(secretCat.ID), nil, nil, nil, readable, catTagVis,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 0 {
		utils.PrintTestError(t, count, 0)
	}

	// Filter by the allowed category -> returns its receipt.
	receipts, count, err := repository.GetPagedReceiptsByGroupId(
		member.ID, utils.UintToString(allGroup.ID), filterBy(allowedCat.ID), nil, nil, nil, readable, catTagVis,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 1 || len(receipts) != 1 || receipts[0].Name != "has-allowed" {
		utils.PrintTestError(t, count, "1 (has-allowed)")
	}

	// No filter -> both receipts visible (the probe fix must not hide unfiltered rows).
	_, count, err = repository.GetPagedReceiptsByGroupId(
		member.ID, utils.UintToString(allGroup.ID), pagedRequestAllReceipts(), nil, nil, nil, readable, catTagVis,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if count != 2 {
		utils.PrintTestError(t, count, 2)
	}
}
