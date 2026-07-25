package services

import (
	"receipt-wrangler/api/internal/models"
)

// This file holds member-isolation enforcement on a receipt PAYLOAD, as opposed to
// paid_by_filter.go which hides whole receipts by payer. A receipt an isolated
// member may still see can carry references to users OUTSIDE their member-visible
// set — a creator, an item's charged-to user, or a comment author. Those are hidden
// here two ways:
//
//   - user-reference fields are MASKED (nulled): created-by (id + denormalized name)
//     on the receipt and every nested entity, and an item's charged-to user. Nulling
//     rather than a "(Restricted)" placeholder is deliberate — a user's PRESENCE must
//     be hidden, not announced (unlike categories/tags). paid_by_user_id needs no
//     masking: surface B already hides the whole receipt when the payer is
//     non-visible, so a visible receipt's payer is always visible.
//   - comments authored by a non-visible user are DROPPED entirely — masking the
//     author would still reveal that the comment (and thus the author's activity)
//     exists.
//
// All of it is a no-op for a receipt whose group does not restrict the viewer (admins,
// supervisors, and non-isolated groups), so non-isolated installs and unrestricted
// viewers are byte-identical to before. The member-visible set is resolved PER RECEIPT'S
// GROUP (memoized across a batch), so a batch spanning an isolated and an open group
// masks each receipt against its own group's set.

// MaskReceiptsForMemberVisibility masks non-visible user references and drops
// non-visible comment authors across a batch of receipts, resolving each receipt's
// group's member-visible set (memoized). No-op for receipts whose group is unrestricted.
func (service PermissionService) MaskReceiptsForMemberVisibility(viewerId uint, receipts []models.Receipt) error {
	if len(receipts) == 0 {
		return nil
	}

	masker := service.newReceiptMemberMasker(viewerId)
	for i := range receipts {
		if err := masker.apply(&receipts[i]); err != nil {
			return err
		}
	}
	return nil
}

// MaskReceiptForMemberVisibility is the single-receipt counterpart of
// MaskReceiptsForMemberVisibility (single GetReceipt, duplicate source, create/update
// responses).
func (service PermissionService) MaskReceiptForMemberVisibility(viewerId uint, receipt *models.Receipt) error {
	if receipt == nil {
		return nil
	}

	masker := service.newReceiptMemberMasker(viewerId)
	return masker.apply(receipt)
}

// receiptMemberMasker resolves the viewer's member-visible set per receipt's group
// (memoized by the groupVisibilityResolver) so a batch spanning several groups resolves
// each group's set exactly once. visible is the current receipt's set, set in apply.
type receiptMemberMasker struct {
	viewerId uint
	resolver *groupVisibilityResolver
	visible  map[uint]struct{}
}

func (service PermissionService) newReceiptMemberMasker(viewerId uint) *receiptMemberMasker {
	return &receiptMemberMasker{
		viewerId: viewerId,
		resolver: service.newGroupVisibilityResolver(viewerId),
	}
}

// apply masks non-visible user references and drops non-visible comment authors on a
// receipt and every nested serialized entity, using the receipt's group's visible set.
// No-op when that group does not restrict the viewer.
func (masker *receiptMemberMasker) apply(receipt *models.Receipt) error {
	visible, unrestricted, err := masker.resolver.forGroup(receipt.GroupId)
	if err != nil {
		return err
	}
	if unrestricted {
		return nil
	}
	masker.visible = visible

	masker.maskBaseModel(&receipt.BaseModel)

	for i := range receipt.ReceiptItems {
		item := &receipt.ReceiptItems[i]
		masker.maskItem(item)
		for j := range item.LinkedItems {
			masker.maskItem(&item.LinkedItems[j])
		}
	}

	for i := range receipt.CustomFields {
		masker.maskBaseModel(&receipt.CustomFields[i].BaseModel)
	}

	for i := range receipt.ImageFiles {
		masker.maskBaseModel(&receipt.ImageFiles[i].BaseModel)
	}

	receipt.Comments = masker.filterComments(receipt.Comments)
	return nil
}

// maskItem masks an item's creator and its charged-to user reference.
func (masker *receiptMemberMasker) maskItem(item *models.Item) {
	masker.maskBaseModel(&item.BaseModel)
	if item.ChargedToUserId != nil && !masker.isVisible(*item.ChargedToUserId) {
		item.ChargedToUserId = nil
		item.ChargedToUser = models.User{}
	}
}

// maskBaseModel nulls the created-by reference (id + denormalized name) when the
// creator is outside the viewer's visible set.
func (masker *receiptMemberMasker) maskBaseModel(base *models.BaseModel) {
	if base.CreatedBy != nil && !masker.isVisible(*base.CreatedBy) {
		base.CreatedBy = nil
		base.CreatedByString = ""
	}
}

// filterComments drops comments authored by a non-visible user (recursing into
// replies so a hidden reply is removed even under a visible parent), and masks the
// created-by reference on the survivors. A comment with no author (nil UserId)
// references no user, so it is kept.
func (masker *receiptMemberMasker) filterComments(comments []models.Comment) []models.Comment {
	if len(comments) == 0 {
		return comments
	}

	filtered := make([]models.Comment, 0, len(comments))
	for i := range comments {
		comment := comments[i]
		if comment.UserId != nil && !masker.isVisible(*comment.UserId) {
			continue
		}
		masker.maskBaseModel(&comment.BaseModel)
		comment.Replies = masker.filterComments(comment.Replies)
		filtered = append(filtered, comment)
	}
	return filtered
}

func (masker *receiptMemberMasker) isVisible(targetId uint) bool {
	return isUserVisible(targetId, masker.viewerId, masker.visible)
}
