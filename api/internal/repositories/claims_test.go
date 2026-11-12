package repositories

import (
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

func seedClaimTestReceipt(t *testing.T, groupId uint, amount string) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         "claim receipt",
		Amount:       decimal.RequireFromString(amount),
		Date:         time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		GroupId:      groupId,
		PaidByUserID: 1,
		Status:       models.OPEN,
	}
	if err := GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt
}

func seedClaim(t *testing.T, groupId uint, name string) models.Claim {
	t.Helper()
	claim, err := NewClaimRepository(nil).CreateClaim(models.Claim{GroupId: groupId, Name: name, Status: models.ClaimStatusDraft})
	if err != nil {
		t.Fatalf("seed claim: %v", err)
	}
	return claim
}

func countRows(t *testing.T, model interface{}, where string, args ...interface{}) int64 {
	t.Helper()
	var count int64
	if err := GetDB().Model(model).Where(where, args...).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return count
}

func TestClaimRepository_LinkedReceiptsAndEntriesRoundTrip(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	repository := NewClaimRepository(nil)

	claim := seedClaim(t, 1, "PCS")
	first := seedClaimTestReceipt(t, 1, "10.10")
	second := seedClaimTestReceipt(t, 1, "5.25")
	if err := repository.LinkReceipts(claim.ID, []uint{second.ID, first.ID}); err != nil {
		t.Fatalf("link: %v", err)
	}

	amount := decimal.RequireFromString("7.50")
	later := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	earlier := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	if _, err := repository.CreateLedgerEntry(models.ClaimLedgerEntry{ClaimId: claim.ID, EntryType: models.ClaimLedgerEntryPayment, Date: later, Amount: &amount}); err != nil {
		t.Fatalf("entry: %v", err)
	}
	if _, err := repository.CreateLedgerEntry(models.ClaimLedgerEntry{ClaimId: claim.ID, EntryType: models.ClaimLedgerEntryNote, Date: earlier, Note: "sent"}); err != nil {
		t.Fatalf("entry: %v", err)
	}

	linked, err := repository.GetLinkedReceipts([]uint{claim.ID})
	if err != nil {
		t.Fatalf("linked: %v", err)
	}
	rows := linked[claim.ID]
	if len(rows) != 2 || rows[0].ReceiptId != first.ID || rows[1].ReceiptId != second.ID {
		utils.PrintTestError(t, rows, "both receipts, ascending by id")
	}
	if len(rows) == 2 && !rows[0].Amount.Equal(decimal.RequireFromString("10.10")) {
		utils.PrintTestError(t, rows[0].Amount, "10.10")
	}

	entries, err := repository.GetLedgerEntries([]uint{claim.ID})
	if err != nil {
		t.Fatalf("entries: %v", err)
	}
	if len(entries[claim.ID]) != 2 || entries[claim.ID][0].EntryType != models.ClaimLedgerEntryNote {
		utils.PrintTestError(t, entries[claim.ID], "the earlier note first")
	}
}

// A receipt is in at most one claim: re-linking replaces the link rather than
// adding a second row.
func TestClaimRepository_ReceiptIsInAtMostOneClaim(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	repository := NewClaimRepository(nil)

	first := seedClaim(t, 1, "first")
	second := seedClaim(t, 1, "second")
	receipt := seedClaimTestReceipt(t, 1, "1.00")

	if err := repository.LinkReceipts(first.ID, []uint{receipt.ID}); err != nil {
		t.Fatalf("link: %v", err)
	}
	if err := repository.LinkReceipts(second.ID, []uint{receipt.ID}); err != nil {
		t.Fatalf("relink: %v", err)
	}

	if got := countRows(t, &models.ClaimReceipt{}, "receipt_id = ?", receipt.ID); got != 1 {
		utils.PrintTestError(t, got, 1)
	}
	if got := countRows(t, &models.ClaimReceipt{}, "claim_id = ?", second.ID); got != 1 {
		utils.PrintTestError(t, got, 1)
	}
}

// Deleting a group or a user deletes their receipts and claims outright; the
// cascading foreign keys are what keep that from leaving links or entries behind.
func TestClaimRepository_GroupAndReceiptDeletesCascade(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	repository := NewClaimRepository(nil)

	claim := seedClaim(t, 2, "doomed")
	receipt := seedClaimTestReceipt(t, 2, "3.00")
	other := seedClaimTestReceipt(t, 2, "4.00")
	if err := repository.LinkReceipts(claim.ID, []uint{receipt.ID, other.ID}); err != nil {
		t.Fatalf("link: %v", err)
	}
	amount := decimal.RequireFromString("1.00")
	if _, err := repository.CreateLedgerEntry(models.ClaimLedgerEntry{ClaimId: claim.ID, EntryType: models.ClaimLedgerEntryPayment, Date: time.Now(), Amount: &amount}); err != nil {
		t.Fatalf("entry: %v", err)
	}

	if err := GetDB().Delete(&models.Receipt{}, receipt.ID).Error; err != nil {
		t.Fatalf("delete receipt: %v", err)
	}
	if got := countRows(t, &models.ClaimReceipt{}, "receipt_id = ?", receipt.ID); got != 0 {
		utils.PrintTestError(t, got, 0)
	}

	if err := GetDB().Delete(&models.Receipt{}, other.ID).Error; err != nil {
		t.Fatalf("delete receipt: %v", err)
	}
	// GroupService.DeleteGroup removes the memberships before the group row.
	if err := GetDB().Where("group_id = ?", 2).Delete(&models.GroupMember{}).Error; err != nil {
		t.Fatalf("delete members: %v", err)
	}
	if err := GetDB().Delete(&models.Group{}, 2).Error; err != nil {
		t.Fatalf("delete group: %v", err)
	}
	if got := countRows(t, &models.Claim{}, "id = ?", claim.ID); got != 0 {
		utils.PrintTestError(t, got, 0)
	}
	if got := countRows(t, &models.ClaimLedgerEntry{}, "claim_id = ?", claim.ID); got != 0 {
		utils.PrintTestError(t, got, 0)
	}
}

// DeleteClaim removes its children itself, so it does not depend on foreign key
// enforcement. Run with foreign keys off to prove it.
func TestClaimRepository_DeleteClaimRemovesChildrenWithoutForeignKeys(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	claim := seedClaim(t, 1, "gone")
	receipt := seedClaimTestReceipt(t, 1, "2.00")
	if err := NewClaimRepository(nil).LinkReceipts(claim.ID, []uint{receipt.ID}); err != nil {
		t.Fatalf("link: %v", err)
	}
	amount := decimal.RequireFromString("2.00")
	if _, err := NewClaimRepository(nil).CreateLedgerEntry(models.ClaimLedgerEntry{ClaimId: claim.ID, EntryType: models.ClaimLedgerEntryPayment, Date: time.Now(), Amount: &amount}); err != nil {
		t.Fatalf("entry: %v", err)
	}

	restore, err := UseTestDbWithoutForeignKeysForTests()
	if err != nil {
		t.Fatalf("no-fk db: %v", err)
	}
	defer restore()

	if err := NewClaimRepository(nil).DeleteClaim(claim.ID); err != nil {
		t.Fatalf("delete claim: %v", err)
	}

	if got := countRows(t, &models.ClaimReceipt{}, "claim_id = ?", claim.ID); got != 0 {
		utils.PrintTestError(t, got, 0)
	}
	if got := countRows(t, &models.ClaimLedgerEntry{}, "claim_id = ?", claim.ID); got != 0 {
		utils.PrintTestError(t, got, 0)
	}
	if got := countRows(t, &models.Receipt{}, "id = ?", receipt.ID); got != 1 {
		t.Errorf("deleting a claim must not delete its receipts")
	}
}

// A claim only holds receipts of its own group, so moving a receipt to another
// group unlinks it, and moving it back does not bring it back.
func TestUpdateReceipt_MovingGroupsUnlinksTheReceiptFromItsClaim(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	claim := seedClaim(t, 1, "PCS")
	moved := seedClaimTestReceipt(t, 1, "9.00")
	stays := seedClaimTestReceipt(t, 1, "1.00")
	if err := NewClaimRepository(nil).LinkReceipts(claim.ID, []uint{moved.ID, stays.ID}); err != nil {
		t.Fatalf("link: %v", err)
	}

	command := func(groupId uint, name string) commands.UpsertReceiptCommand {
		return commands.UpsertReceiptCommand{
			Name:         name,
			Amount:       decimal.RequireFromString("9.00"),
			Date:         moved.Date,
			PaidByUserID: 1,
			Status:       models.OPEN,
			GroupId:      groupId,
		}
	}

	// An edit that keeps the group leaves the link alone.
	if _, err := NewReceiptRepository(nil).UpdateReceipt(utils.UintToString(moved.ID), command(1, "renamed"), 1); err != nil {
		t.Fatalf("update: %v", err)
	}
	if got := countRows(t, &models.ClaimReceipt{}, "receipt_id = ?", moved.ID); got != 1 {
		utils.PrintTestError(t, got, 1)
	}

	if _, err := NewReceiptRepository(nil).UpdateReceipt(utils.UintToString(moved.ID), command(2, "moved"), 1); err != nil {
		t.Fatalf("move: %v", err)
	}
	if got := countRows(t, &models.ClaimReceipt{}, "receipt_id = ?", moved.ID); got != 0 {
		utils.PrintTestError(t, got, 0)
	}

	if _, err := NewReceiptRepository(nil).UpdateReceipt(utils.UintToString(moved.ID), command(1, "back"), 1); err != nil {
		t.Fatalf("move back: %v", err)
	}
	if got := countRows(t, &models.ClaimReceipt{}, "receipt_id = ?", moved.ID); got != 0 {
		t.Errorf("a receipt moved back must not rejoin its old claim")
	}
	if got := countRows(t, &models.ClaimReceipt{}, "receipt_id = ?", stays.ID); got != 1 {
		t.Errorf("other receipts in the claim must keep their link")
	}
}

// A link left pointing at a claim of a group the receipt is no longer in is not
// shown, and does not block deleting the receipt.
func TestClaimRepository_StaleLinkIsIgnored(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	claim := seedClaim(t, 1, "PCS")
	receipt := seedClaimTestReceipt(t, 2, "9.00")
	if err := GetDB().Create(&models.ClaimReceipt{ReceiptId: receipt.ID, ClaimId: claim.ID}).Error; err != nil {
		t.Fatalf("stale link: %v", err)
	}

	linked, err := NewClaimRepository(nil).GetLinkedReceipts([]uint{claim.ID})
	if err != nil {
		t.Fatalf("linked: %v", err)
	}
	if len(linked[claim.ID]) != 0 {
		utils.PrintTestError(t, linked[claim.ID], "no receipts")
	}

	blocking, err := NewClaimRepository(nil).GetClaimLinksForGroupReceipts([]uint{receipt.ID})
	if err != nil {
		t.Fatalf("links: %v", err)
	}
	if len(blocking) != 0 {
		utils.PrintTestError(t, blocking, "no blocking links")
	}
}

func TestClaimRepository_UpdateClaimClearsOptionalFields(t *testing.T) {
	defer TruncateTestDb()
	CreateTestGroupWithUsers()
	repository := NewClaimRepository(nil)

	amount := decimal.RequireFromString("100.00")
	day := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	claim, err := repository.CreateClaim(models.Claim{
		GroupId: 1, Name: "PCS", Status: models.ClaimStatusSubmitted, ReferenceNumber: "TA-1",
		ClaimType: "Travel voucher", SubmittedDate: &day, DueDate: &day, ClaimedAmount: &amount, Notes: "n",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := repository.UpdateClaim(claim.ID, models.Claim{Name: "PCS 2", Status: models.ClaimStatusPaid}); err != nil {
		t.Fatalf("update: %v", err)
	}

	stored, err := repository.GetClaimById(claim.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if stored.Name != "PCS 2" || stored.Status != models.ClaimStatusPaid {
		utils.PrintTestError(t, stored, "renamed and PAID")
	}
	if stored.ReferenceNumber != "" || stored.ClaimType != "" || stored.Notes != "" ||
		stored.SubmittedDate != nil || stored.DueDate != nil || stored.ClaimedAmount != nil {
		utils.PrintTestError(t, stored, "every optional field cleared")
	}
	if stored.GroupId != 1 {
		t.Errorf("an update must not move the claim")
	}
}
