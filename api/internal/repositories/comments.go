package repositories

import (
	"fmt"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"

	"gorm.io/gorm"
)

type CommentRepository struct {
	BaseRepository
}

// AuthorVisibilityResolver reports whether a notification recipient is allowed to
// see the comment's author WITHIN the receipt's group (member-presence isolation). It
// lets the comment repository suppress a notification whose body names an author the
// recipient may not see in that group, without importing the service layer that
// resolves visibility (mirroring the PaidByAllowedResolver injection). A nil resolver
// disables suppression, so every recipient receives (backward compatible).
type AuthorVisibilityResolver func(authorId uint, recipientId uint, groupId uint) (canSeeAuthor bool, err error)

func NewCommentRepository(tx *gorm.DB) CommentRepository {
	repository := CommentRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository CommentRepository) AddComment(command commands.UpsertCommentCommand, authorVisibleTo AuthorVisibilityResolver) (models.Comment, error) {
	db := repository.GetDB()
	comment := models.Comment{
		Comment:   command.Comment,
		ReceiptId: command.ReceiptId,
		UserId:    command.UserId,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		repository.SetTransaction(tx)

		err := tx.Model(&comment).Create(&comment).Error
		if err != nil {
			return err
		}

		err = repository.sendNotificationsToUsers(comment, authorVisibleTo)
		if err != nil {
			return err
		}

		repository.ClearTransaction()
		return nil
	})

	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func (repository CommentRepository) GetUsersInCommentThread(comment models.Comment) ([]uint, error) {
	db := repository.GetDB()
	userIds := make([]interface{}, 0)
	result := make([]uint, 0)

	if *comment.UserId > 0 {
		userIds = append(userIds, *comment.UserId)
		result = append(result, *comment.UserId)
	}

	if *comment.CommentId > 0 {
		var threadComments []models.Comment
		var parentComment models.Comment

		err := db.Model(models.Comment{}).Where("comment_id = ?", comment.CommentId).Find(&threadComments).Error
		if err != nil {
			return nil, err
		}

		err = db.Model(models.Comment{}).Where("id = ?", comment.CommentId).Find(&parentComment).Error
		if err != nil {
			return nil, err
		}

		if *parentComment.UserId > 0 && !utils.Contains(userIds, *parentComment.UserId) {
			userIds = append(userIds, *parentComment.UserId)
			result = append(result, *parentComment.UserId)
		}

		for _, comment := range threadComments {
			if comment.ID > 0 && !utils.Contains(userIds, *comment.UserId) {
				userIds = append(userIds, *comment.UserId)
				result = append(result, *comment.UserId)

			}
		}
	}

	return result, nil
}

func (repository CommentRepository) DeleteComment(commentId string, tokenUserId uint) error {
	db := repository.GetDB()
	var comment models.Comment

	err := db.Model(models.Comment{}).Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		return err
	}

	if *comment.UserId == tokenUserId {
		err = db.Model(models.Comment{}).Where("id = ?", commentId).Delete(&comment).Error
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not allowed to delete another user's comment")
	}

	return nil
}

func (repository CommentRepository) sendNotificationsToUsers(comment models.Comment, authorVisibleTo AuthorVisibilityResolver) error {
	var receipt models.Receipt
	authorId := *comment.UserId
	usersToOmit := make([]interface{}, 0)
	usersToOmit = append(usersToOmit, authorId)
	notificationRepository := NewNotificationRepository(repository.TX)
	receiptRepository := NewReceiptRepository(repository.TX)

	receipt, err := receiptRepository.GetReceiptById(utils.UintToString(comment.ReceiptId))
	if err != nil {
		return err
	}

	if comment.CommentId == nil {
		// Member isolation: suppress the notification for any group member who may
		// not see the comment's author (the body names them) by adding them to the
		// omit list.
		hiddenRecipients, err := repository.recipientsWhoCannotSeeAuthor(authorId, receipt.GroupId, authorVisibleTo)
		if err != nil {
			return err
		}
		usersToOmit = append(usersToOmit, hiddenRecipients...)

		err = notificationRepository.SendNotificationToGroup(receipt.GroupId, "Comment Added", fmt.Sprintf("%s has added a comment to a receipt in group %s. %s", BuildParamaterisedString("userId", authorId, "displayName", "string"), BuildParamaterisedString("groupId", receipt.GroupId, "name", "string"), BuildParamaterisedString("receiptId", comment.ReceiptId, "noop", "link")), models.NOTIFICATION_TYPE_NORMAL, usersToOmit)
		if err != nil {
			return err
		}
	} else {
		threadUsers, err := repository.GetUsersInCommentThread(comment)
		if err != nil {
			return err
		}

		hiddenRecipients, err := omitRecipientsWhoCannotSeeAuthor(authorId, threadUsers, receipt.GroupId, authorVisibleTo)
		if err != nil {
			return err
		}
		usersToOmit = append(usersToOmit, hiddenRecipients...)

		err = notificationRepository.SendNotificationToUsers(threadUsers, "Comment Replied", fmt.Sprintf("%s has replied to a thread that you are a part of.", BuildParamaterisedString("userId", authorId, "displayName", "string")), models.NOTIFICATION_TYPE_NORMAL, usersToOmit)
		if err != nil {
			return err
		}
	}

	return nil
}

// recipientsWhoCannotSeeAuthor returns the ids of the group's members who may not
// see the comment author under member isolation (so they can be added to a
// notification's omit list). A nil resolver returns no ids.
func (repository CommentRepository) recipientsWhoCannotSeeAuthor(authorId uint, groupId uint, authorVisibleTo AuthorVisibilityResolver) ([]interface{}, error) {
	if authorVisibleTo == nil {
		return nil, nil
	}

	// Fast path: a non-isolated group hides nothing, so no recipient can fail to see the
	// author. Skip the roster fetch and the per-member visibility lookups entirely — the
	// common case does one cheap flag read instead of O(members) resolver calls.
	var group struct {
		IsolateMembers bool
	}
	if err := repository.GetDB().Model(&models.Group{}).
		Select("isolate_members").
		Where("id = ?", groupId).
		Scan(&group).Error; err != nil {
		return nil, err
	}
	if !group.IsolateMembers {
		return nil, nil
	}

	groupMemberRepository := NewGroupMemberRepository(repository.TX)
	members, err := groupMemberRepository.GetsGroupMembersByGroupId(utils.UintToString(groupId))
	if err != nil {
		return nil, err
	}

	recipientIds := make([]uint, 0, len(members))
	for _, member := range members {
		recipientIds = append(recipientIds, member.UserID)
	}

	return omitRecipientsWhoCannotSeeAuthor(authorId, recipientIds, groupId, authorVisibleTo)
}

// omitRecipientsWhoCannotSeeAuthor evaluates each candidate recipient against the
// author-visibility resolver (scoped to groupId) and returns those who may not see the
// author. The author is never included. A nil resolver returns no ids.
func omitRecipientsWhoCannotSeeAuthor(authorId uint, recipientIds []uint, groupId uint, authorVisibleTo AuthorVisibilityResolver) ([]interface{}, error) {
	if authorVisibleTo == nil {
		return nil, nil
	}

	hidden := make([]interface{}, 0)
	for _, recipientId := range recipientIds {
		if recipientId == authorId {
			continue
		}

		canSee, err := authorVisibleTo(authorId, recipientId, groupId)
		if err != nil {
			return nil, err
		}
		if !canSee {
			hidden = append(hidden, recipientId)
		}
	}

	return hidden, nil
}
