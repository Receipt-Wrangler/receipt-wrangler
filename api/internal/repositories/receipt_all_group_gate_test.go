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
	// The read gate and category/tag resolver must be supplied together; this
	// test carries no category/tag filter, so the resolver is never invoked.
	catTagVis := func(uint) (CategoryTagVisibility, error) {
		return CategoryTagVisibility{CategoryUnrestricted: true, TagUnrestricted: true}, nil
	}

	repository := NewReceiptRepository(nil)
	receipts, count, err := repository.GetPagedReceiptsByGroupId(
		member.ID, utils.UintToString(allGroup.ID), pagedRequestAllReceipts(), nil, nil, nil, readable, catTagVis,
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

// The All-group read gate and the per-group category/tag resolver must be passed
// together: supplying only one silently reopens a cross-group leak (a read gate
// without the category/tag resolver lets a category/tag filter fall through to
// the flat, group-unscoped subquery; the reverse expands the read set), so
// GetPagedReceiptsByGroupId fails closed on a mismatch. Passing neither is the
// single-group shape and stays allowed.
func TestGetPagedReceipts_AllGroup_RequiresBothResolvers(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	allGroup := models.Group{Name: "pair-all", IsAllGroup: true}
	db.Create(&allGroup)

	member := models.User{Username: "pair-member", Password: "x"}
	db.Create(&member)

	repository := NewReceiptRepository(nil)
	readable := func(uint) (bool, error) { return true, nil }
	catTagVis := func(uint) (CategoryTagVisibility, error) {
		return CategoryTagVisibility{CategoryUnrestricted: true, TagUnrestricted: true}, nil
	}

	call := func(r GroupReadableResolver, c CategoryTagVisibilityResolver) error {
		_, _, err := repository.GetPagedReceiptsByGroupId(
			member.ID, utils.UintToString(allGroup.ID), pagedRequestAllReceipts(), nil, nil, nil, r, c,
		)
		return err
	}

	// Read gate without the category/tag resolver -> error (the leak this guards).
	if err := call(readable, nil); err == nil {
		utils.PrintTestError(t, nil, "expected error when only the readable resolver is supplied")
	}
	// Category/tag resolver without the read gate -> error.
	if err := call(nil, catTagVis); err == nil {
		utils.PrintTestError(t, nil, "expected error when only the category/tag resolver is supplied")
	}
	// Both supplied -> no guard error.
	if err := call(readable, catTagVis); err != nil {
		utils.PrintTestError(t, err, "no error when both resolvers are supplied")
	}
	// Neither supplied -> allowed (guard only fires on a mismatch).
	if err := call(nil, nil); err != nil {
		utils.PrintTestError(t, err, "no error when neither resolver is supplied")
	}
}
