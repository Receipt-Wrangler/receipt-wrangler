package services

import (
	"slices"
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"

	"github.com/shopspring/decimal"
)

var firstCommentTestTime = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func seedFirstCommentReceipt(t *testing.T, groupId uint, payerId uint, name string) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(1),
		Date:         time.Now(),
		PaidByUserID: payerId,
		GroupId:      groupId,
		Status:       models.OPEN,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt
}

// seedCommentAt stamps created_at so "first" is decided by time, not insertion order.
func seedCommentAt(t *testing.T, receiptId uint, userId *uint, text string, createdAt time.Time) {
	t.Helper()
	db := repositories.GetDB()
	comment := models.Comment{Comment: text, ReceiptId: receiptId, UserId: userId}
	if err := db.Create(&comment).Error; err != nil {
		t.Fatalf("seed comment: %v", err)
	}
	if err := db.Model(&comment).Update("created_at", createdAt).Error; err != nil {
		t.Fatalf("stamp comment: %v", err)
	}
}

func firstCommentOf(receipt models.Receipt) string {
	if receipt.FirstComment == nil {
		return "<nil>"
	}
	return *receipt.FirstComment
}

func TestLoadFirstVisibleComments_PicksTheEarliestComment(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetIsolationCaches()

	fx := seedIsolatedReceiptGroup(t)
	service := NewPermissionService(nil)

	commented := seedFirstCommentReceipt(t, fx.groupId, fx.supervisorId, "commented")
	seedCommentAt(t, commented.ID, &fx.memberAId, "second", firstCommentTestTime.Add(time.Hour))
	seedCommentAt(t, commented.ID, &fx.memberBId, "first", firstCommentTestTime)
	bare := seedFirstCommentReceipt(t, fx.groupId, fx.supervisorId, "bare")

	// The supervisor sees every member, so nothing is skipped.
	receipts := []models.Receipt{commented, bare}
	if err := service.LoadFirstVisibleComments(fx.supervisorId, receipts); err != nil {
		t.Fatalf("load: %v", err)
	}

	if got := firstCommentOf(receipts[0]); got != "first" {
		t.Errorf("expected the earliest comment, got %q", got)
	}
	if receipts[1].FirstComment != nil {
		t.Errorf("expected no first comment on a receipt without comments, got %q", *receipts[1].FirstComment)
	}
}

// An isolated member never sees a comment by a peer they cannot see — the list
// falls through to the first comment they can, exactly as the sort does.
func TestLoadFirstVisibleComments_SkipsAuthorsTheViewerCannotSee(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetIsolationCaches()

	fx := seedIsolatedReceiptGroup(t)
	service := NewPermissionService(nil)

	hiddenFirst := seedFirstCommentReceipt(t, fx.groupId, fx.memberAId, "hidden-first")
	seedCommentAt(t, hiddenFirst.ID, &fx.memberBId, "from B", firstCommentTestTime)
	seedCommentAt(t, hiddenFirst.ID, &fx.memberAId, "from A", firstCommentTestTime.Add(time.Hour))

	onlyHidden := seedFirstCommentReceipt(t, fx.groupId, fx.memberAId, "only-hidden")
	seedCommentAt(t, onlyHidden.ID, &fx.memberBId, "from B only", firstCommentTestTime)

	authorless := seedFirstCommentReceipt(t, fx.groupId, fx.memberAId, "authorless")
	seedCommentAt(t, authorless.ID, nil, "from the system", firstCommentTestTime)

	supervised := seedFirstCommentReceipt(t, fx.groupId, fx.memberAId, "supervised")
	seedCommentAt(t, supervised.ID, &fx.supervisorId, "from the supervisor", firstCommentTestTime)

	receipts := []models.Receipt{hiddenFirst, onlyHidden, authorless, supervised}
	if err := service.LoadFirstVisibleComments(fx.memberAId, receipts); err != nil {
		t.Fatalf("load: %v", err)
	}

	want := []string{"from A", "<nil>", "from the system", "from the supervisor"}
	for i, receipt := range receipts {
		if got := firstCommentOf(receipt); got != want[i] {
			t.Errorf("%s: expected %q, got %q", receipt.Name, want[i], got)
		}
	}
}

func TestLoadFirstVisibleComments_EmptyBatchIsANoop(t *testing.T) {
	service := NewPermissionService(nil)
	if err := service.LoadFirstVisibleComments(1, nil); err != nil {
		t.Fatalf("expected no error for an empty batch, got %v", err)
	}
}

func TestCommentAuthorVisibilityResolver_RestrictsAnIsolatedMember(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetIsolationCaches()

	fx := seedIsolatedReceiptGroup(t)
	service := NewPermissionService(nil)

	visible, unrestricted, err := service.CommentAuthorVisibilityResolver(fx.memberAId)(fx.groupId)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if unrestricted {
		t.Fatalf("expected an isolated member to be restricted")
	}
	if !slices.Contains(visible, fx.memberAId) || !slices.Contains(visible, fx.supervisorId) {
		t.Errorf("expected self and the supervisor to be visible, got %v", visible)
	}
	if slices.Contains(visible, fx.memberBId) {
		t.Errorf("expected the isolated peer to be hidden, got %v", visible)
	}

	_, unrestricted, err = service.CommentAuthorVisibilityResolver(fx.supervisorId)(fx.groupId)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !unrestricted {
		t.Errorf("expected the supervisor to be unrestricted")
	}
}
