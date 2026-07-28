package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

// authorVisibleToResolver builds the same author-visibility resolver the comment
// handler injects, backed by the real PermissionService, so the repository's
// notification suppression is exercised end-to-end.
func authorVisibleToResolver() repositories.AuthorVisibilityResolver {
	permissionService := NewPermissionService(nil)
	return func(authorId uint, recipientId uint, groupId uint) (bool, error) {
		visible, unrestricted, err := permissionService.GetVisibleUserIdsForUserInGroup(recipientId, groupId)
		if err != nil {
			return false, err
		}
		if unrestricted {
			return true, nil
		}
		_, ok := visible[authorId]
		return ok, nil
	}
}

func seedIsoReceipt(t *testing.T, groupId uint, paidByUserId uint) models.Receipt {
	t.Helper()
	receipt := models.Receipt{Name: "iso-receipt", GroupId: groupId, PaidByUserID: paidByUserId}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt
}

func countNotifications(t *testing.T, userId uint) int {
	t.Helper()
	notifications, err := repositories.NewNotificationRepository(nil).GetNotificationsForUser(userId)
	if err != nil {
		t.Fatalf("get notifications for %d: %v", userId, err)
	}
	return len(notifications)
}

// A comment notification is suppressed for an isolated recipient who cannot see
// the author, while a supervisor (who can) still receives it.
func TestCommentNotificationSuppressedForIsolatedRecipient(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	group := seedIsoGroup(t, "iso-comment", true)
	supRole := seedIsoRole(t, "iso-comment-sup", true)
	memberRole := seedIsoRole(t, "iso-comment-mem", false)

	author := seedIsoUser(t, "iso-comment-author")
	peer := seedIsoUser(t, "iso-comment-peer")
	supervisor := seedIsoUser(t, "iso-comment-sup-user")
	seedIsoMember(t, group.ID, author.ID, &memberRole.ID)
	seedIsoMember(t, group.ID, peer.ID, &memberRole.ID)
	seedIsoMember(t, group.ID, supervisor.ID, &supRole.ID)

	receipt := seedIsoReceipt(t, group.ID, author.ID)

	comment := commands.UpsertCommentCommand{Comment: "hello", ReceiptId: receipt.ID, UserId: &author.ID}
	if _, err := repositories.NewCommentRepository(nil).AddComment(comment, authorVisibleToResolver()); err != nil {
		t.Fatalf("add comment: %v", err)
	}

	if got := countNotifications(t, peer.ID); got != 0 {
		t.Errorf("isolated peer who cannot see author should get 0 notifications, got %d", got)
	}
	if got := countNotifications(t, supervisor.ID); got != 1 {
		t.Errorf("supervisor who sees the author should get 1 notification, got %d", got)
	}
	if got := countNotifications(t, author.ID); got != 0 {
		t.Errorf("author should never be notified of their own comment, got %d", got)
	}
}

// Backward compatibility: in a non-isolated group every recipient is unrestricted,
// so the comment notification is delivered to all of them.
func TestCommentNotificationDeliveredToAllInNonIsolatedGroup(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearRolePermissionCacheAll()

	group := seedIsoGroup(t, "non-iso-comment", false)
	memberRole := seedIsoRole(t, "non-iso-comment-mem", false)

	author := seedIsoUser(t, "non-iso-author")
	peer := seedIsoUser(t, "non-iso-peer")
	seedIsoMember(t, group.ID, author.ID, &memberRole.ID)
	seedIsoMember(t, group.ID, peer.ID, &memberRole.ID)

	receipt := seedIsoReceipt(t, group.ID, author.ID)

	comment := commands.UpsertCommentCommand{Comment: "hello", ReceiptId: receipt.ID, UserId: &author.ID}
	if _, err := repositories.NewCommentRepository(nil).AddComment(comment, authorVisibleToResolver()); err != nil {
		t.Fatalf("add comment: %v", err)
	}

	if got := countNotifications(t, peer.ID); got != 1 {
		t.Errorf("peer in a non-isolated group should receive the notification, got %d", got)
	}
}
