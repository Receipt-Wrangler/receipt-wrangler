package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

// This file holds the row-level "paid by" visibility enforcement: which receipts
// a group role lets its members see, keyed on the receipt's paid_by_user_id.
// Unlike the category/tag grants (which strip fields off a still-visible receipt
// in grant_filter.go), paid-by hides the WHOLE receipt — so enforcement adds a
// WHERE to the paged query (PaidByListResolver), denies single-receipt access
// (ReceiptPaidByVisible), and post-filters the unpaginated read surfaces
// (FilterReceiptsByPaidBy). No request-filter intersection is needed: because the
// paged query already row-filters on the allowed set, a caller filtering by a
// payer they cannot see is neutralized by that AND (a hidden id intersects to
// nothing), so it cannot probe receipt existence.

// PaidByListResolver adapts GetGroupPaidByUserIdsForUser to the per-group
// resolver the receipt repository uses to build its paged WHERE clause. The
// returned closure reports the allowed paid_by_user_id values for userId in a
// group, or unrestricted == true (see every payer). Pass the result to
// GetPagedReceiptsByGroupId; pass nil there to skip paid-by filtering.
func (service PermissionService) PaidByListResolver(userId uint) repositories.PaidByAllowedResolver {
	return func(groupId uint) ([]uint, bool, error) {
		allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
		if err != nil || unrestricted {
			return nil, unrestricted, err
		}
		return uintSetToSlice(allowed), false, nil
	}
}

// ReceiptPaidByVisible reports whether a user may see a receipt given its group
// and paid_by user. Unrestricted users (and non-members, who have no group role)
// always pass. Used to deny single-receipt and dependent reads (image, comments,
// duplicate source, update/delete) for a receipt hidden by the paid-by filter.
func (service PermissionService) ReceiptPaidByVisible(userId uint, groupId uint, paidByUserId uint) (bool, error) {
	allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		return false, err
	}
	if unrestricted {
		return true, nil
	}
	_, ok := allowed[paidByUserId]
	return ok, nil
}

// FilterReceiptsByPaidBy drops the receipts a user may not see by the paid-by
// filter, resolving each group's allowed set at most once. Used by the
// unpaginated read surfaces (search, group-ids list) where there is no
// totalCount to keep consistent, so removing rows is safe.
func (service PermissionService) FilterReceiptsByPaidBy(userId uint, receipts []models.Receipt) ([]models.Receipt, error) {
	if len(receipts) == 0 {
		return receipts, nil
	}

	type groupVisibility struct {
		allowed      map[uint]struct{}
		unrestricted bool
	}
	cache := map[uint]groupVisibility{}

	filtered := make([]models.Receipt, 0, len(receipts))
	for i := range receipts {
		groupId := receipts[i].GroupId
		visibility, ok := cache[groupId]
		if !ok {
			allowed, unrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
			if err != nil {
				return nil, err
			}
			visibility = groupVisibility{allowed: allowed, unrestricted: unrestricted}
			cache[groupId] = visibility
		}

		if visibility.unrestricted {
			filtered = append(filtered, receipts[i])
			continue
		}
		if _, allowedOk := visibility.allowed[receipts[i].PaidByUserID]; allowedOk {
			filtered = append(filtered, receipts[i])
		}
	}

	return filtered, nil
}

// uintSetToSlice flattens a set to a slice (order is irrelevant for an IN clause).
func uintSetToSlice(set map[uint]struct{}) []uint {
	ids := make([]uint, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}
