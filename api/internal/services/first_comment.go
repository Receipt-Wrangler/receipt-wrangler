package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

// This file backs the receipts table's Comment column: the text of each receipt's
// first comment, shown on the paged list and sortable on it. Both halves honour
// member isolation the way every other comment read does (filterComments drops a
// comment whose author the caller cannot see in the receipt's group), and both
// resolve it through the same per-group groupVisibilityResolver, so the column can
// never sort by a comment it would not display.
//
// There is no permission of its own: reading a receipt already means reading its
// comments on every other surface, so group.receipts.read - the list's gate - is
// the gate here too.

// CommentAuthorVisibilityResolver adapts the per-group member-visible set to the
// resolver the receipt repository uses to build a first-comment sort. Each group is
// resolved on first use and memoized, the way the masker resolves a batch. The
// viewer is always included, mirroring isUserVisible, which never hides a user's
// own comments from them.
func (service PermissionService) CommentAuthorVisibilityResolver(viewerId uint) repositories.CommentAuthorVisibilityResolver {
	visibility := service.newGroupVisibilityResolver(viewerId)
	return func(groupId uint) ([]uint, bool, error) {
		visible, unrestricted, err := visibility.forGroup(groupId)
		if err != nil || unrestricted {
			return nil, unrestricted, err
		}

		visibleIds := uintSetToSlice(visible)
		if _, ok := visible[viewerId]; !ok {
			visibleIds = append(visibleIds, viewerId)
		}
		return visibleIds, false, nil
	}
}

// LoadFirstVisibleComments sets FirstComment on each receipt to the text of its
// earliest comment the viewer may see, leaving it nil when there is none. It loads
// the comments of the whole batch in one query, and each receipt is judged by its
// own group's visibility, so an All-group page spanning an isolated and an open
// group gets each right. A comment with no author references no one and is kept,
// as filterComments keeps it.
func (service PermissionService) LoadFirstVisibleComments(viewerId uint, receipts []models.Receipt) error {
	if len(receipts) == 0 {
		return nil
	}

	receiptIds := make([]uint, len(receipts))
	groupIdByReceiptId := make(map[uint]uint, len(receipts))
	for i := range receipts {
		receiptIds[i] = receipts[i].ID
		groupIdByReceiptId[receipts[i].ID] = receipts[i].GroupId
	}

	comments, err := repositories.NewCommentRepository(service.TX).GetCommentsForReceiptIds(receiptIds)
	if err != nil {
		return err
	}

	// comments arrive in first-comment order, so the first visible one seen for a
	// receipt is its first comment.
	visibility := service.newGroupVisibilityResolver(viewerId)
	firstByReceiptId := make(map[uint]string, len(receipts))
	for _, comment := range comments {
		if _, found := firstByReceiptId[comment.ReceiptId]; found {
			continue
		}
		if comment.UserId != nil {
			visible, err := visibility.isVisible(groupIdByReceiptId[comment.ReceiptId], *comment.UserId)
			if err != nil {
				return err
			}
			if !visible {
				continue
			}
		}
		firstByReceiptId[comment.ReceiptId] = comment.Comment
	}

	for i := range receipts {
		if text, found := firstByReceiptId[receipts[i].ID]; found {
			receipts[i].FirstComment = &text
		}
	}

	return nil
}
