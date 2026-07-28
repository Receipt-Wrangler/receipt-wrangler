package repositories

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

var commentRepository CommentRepository

func setupCommentTest() {
	CreateTestGroupWithUsers()
	createTestReceipt()
	commentRepository = NewCommentRepository(nil)
}

func createTestReceipt() {
	receipt := models.Receipt{
		Name:         "test",
		PaidByUserID: 1,
		GroupId:      1,
	}

	GetDB().Create(&receipt)
}

func teardownCommentTest() {
	TruncateTestDb()
}

func TestShouldAddCommentAndSendNotificationToAllGroupUsers(t *testing.T) {
	defer teardownCommentTest()
	setupCommentTest()
	userId := uint(1)
	comment := commands.UpsertCommentCommand{
		Comment:   "test",
		ReceiptId: 1,
		UserId:    &userId,
	}

	newComment, err := commentRepository.AddComment(comment, nil)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if newComment.ID != 1 {
		utils.PrintTestError(t, newComment.ID, 1)
	}

	notificationRepository := NewNotificationRepository(nil)

	user1Notifications, _ := notificationRepository.GetNotificationsForUser(1)
	if len(user1Notifications) > 0 {
		utils.PrintTestError(t, len(user1Notifications), 0)
	}

	user2Notifications, _ := notificationRepository.GetNotificationsForUser(2)
	if len(user2Notifications) != 1 {
		utils.PrintTestError(t, len(user2Notifications), 1)
	}

	user3Notifications, _ := notificationRepository.GetNotificationsForUser(3)
	if len(user3Notifications) != 1 {
		utils.PrintTestError(t, len(user3Notifications), 1)
	}
}

// The reply/thread notification path omits recipients who cannot see the reply author.
// This branch is not reachable through the public AddComment (UpsertCommentCommand
// carries no commentId, and AddComment never sets Comment.CommentId), so it is exercised
// directly with an injected resolver — which also proves the groupId threading is wired.
func TestSendNotificationsReplyBranchOmitsRecipientsWhoCannotSeeAuthor(t *testing.T) {
	defer teardownCommentTest()
	setupCommentTest()

	author := uint(1)
	hidden := uint(2)
	visible := uint(3)
	receiptId := uint(1)

	db := GetDB()
	parent := models.Comment{Comment: "parent", ReceiptId: receiptId, UserId: &hidden}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("parent: %v", err)
	}
	visComment := models.Comment{Comment: "vis", ReceiptId: receiptId, UserId: &visible, CommentId: &parent.ID}
	if err := db.Create(&visComment).Error; err != nil {
		t.Fatalf("vis comment: %v", err)
	}
	reply := models.Comment{Comment: "reply", ReceiptId: receiptId, UserId: &author, CommentId: &parent.ID}
	if err := db.Create(&reply).Error; err != nil {
		t.Fatalf("reply: %v", err)
	}

	// Everyone except `hidden` can see the author.
	resolver := func(authorId uint, recipientId uint, groupId uint) (bool, error) {
		return recipientId != hidden, nil
	}
	if err := commentRepository.sendNotificationsToUsers(reply, resolver); err != nil {
		t.Fatalf("sendNotificationsToUsers: %v", err)
	}

	notificationRepository := NewNotificationRepository(nil)
	count := func(userId uint) int {
		ns, err := notificationRepository.GetNotificationsForUser(userId)
		if err != nil {
			t.Fatalf("notifications for %d: %v", userId, err)
		}
		return len(ns)
	}
	if got := count(hidden); got != 0 {
		t.Errorf("recipient who cannot see the reply author should be omitted, got %d", got)
	}
	if got := count(visible); got != 1 {
		t.Errorf("recipient who can see the reply author should receive the reply notification, got %d", got)
	}
	if got := count(author); got != 0 {
		t.Errorf("author should not be notified of their own reply, got %d", got)
	}
}

// A resolver error in the notification path propagates out of AddComment.
func TestAddCommentPropagatesAuthorVisibilityResolverError(t *testing.T) {
	defer teardownCommentTest()
	setupCommentTest()

	// Isolate the group so the notification path actually invokes the resolver (the
	// non-isolated fast path would otherwise short-circuit before calling it).
	if err := GetDB().Model(&models.Group{}).Where("id = ?", uint(1)).
		Update("isolate_members", true).Error; err != nil {
		t.Fatalf("isolate group: %v", err)
	}

	author := uint(1)
	comment := commands.UpsertCommentCommand{Comment: "test", ReceiptId: 1, UserId: &author}
	boom := errors.New("resolver boom")
	resolver := func(authorId uint, recipientId uint, groupId uint) (bool, error) {
		return false, boom
	}

	if _, err := commentRepository.AddComment(comment, resolver); !errors.Is(err, boom) {
		t.Fatalf("expected the resolver error to propagate from AddComment, got %v", err)
	}
}
