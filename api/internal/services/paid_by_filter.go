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
//
// Member isolation is folded in at the single shared resolution point
// (paidByAllowedForGroup) so all three consumers below inherit it: an isolated
// member must not see a receipt paid by a user outside their member-visible set,
// even in a group whose role applies no paid-by restriction of its own. The
// member-visible set is resolved ONCE per consumer (not per group/row) and
// intersected into each group's paid-by allowed set. It is a no-op for
// unrestricted viewers (admins, and anyone not an isolated member), so
// non-isolated installs behave exactly as before.

// PaidByListResolver adapts paidByAllowedForGroup to the per-group resolver the
// receipt repository uses to build its paged WHERE clause. The returned closure
// reports the allowed paid_by_user_id values for userId in a group, or
// unrestricted == true (see every payer). The member-visible set is resolved once,
// here, and threaded into every per-group resolution. Pass the result to
// GetPagedReceiptsByGroupId; pass nil there to skip paid-by filtering.
func (service PermissionService) PaidByListResolver(userId uint) repositories.PaidByAllowedResolver {
	visibleSet, memberUnrestricted, visErr := service.GetVisibleUserIdsForUser(userId)
	return func(groupId uint) ([]uint, bool, error) {
		if visErr != nil {
			return nil, false, visErr
		}
		allowed, unrestricted, err := service.paidByAllowedForGroup(userId, groupId, visibleSet, memberUnrestricted)
		if err != nil || unrestricted {
			return nil, unrestricted, err
		}
		return uintSetToSlice(allowed), false, nil
	}
}

// ReceiptPaidByVisible reports whether a user may see a receipt given its group
// and paid_by user. Unrestricted users (and non-members, who have no group role)
// always pass. Used to deny single-receipt and dependent reads (image, comments,
// duplicate source, update/delete) for a receipt hidden by the paid-by filter or
// by member isolation (payer outside the caller's member-visible set).
func (service PermissionService) ReceiptPaidByVisible(userId uint, groupId uint, paidByUserId uint) (bool, error) {
	visibleSet, memberUnrestricted, err := service.GetVisibleUserIdsForUser(userId)
	if err != nil {
		return false, err
	}
	allowed, unrestricted, err := service.paidByAllowedForGroup(userId, groupId, visibleSet, memberUnrestricted)
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
// filter (or by member isolation), resolving the member-visible set once and each
// group's allowed set at most once. Used by the unpaginated read surfaces (search,
// group-ids list) where there is no totalCount to keep consistent, so removing
// rows is safe.
func (service PermissionService) FilterReceiptsByPaidBy(userId uint, receipts []models.Receipt) ([]models.Receipt, error) {
	if len(receipts) == 0 {
		return receipts, nil
	}

	visibleSet, memberUnrestricted, err := service.GetVisibleUserIdsForUser(userId)
	if err != nil {
		return nil, err
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
			allowed, unrestricted, err := service.paidByAllowedForGroup(userId, groupId, visibleSet, memberUnrestricted)
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

// paidByAllowedForGroup resolves a group's paid-by allowed set for a user and folds
// in the caller's already-resolved member-visible set. This is the single shared
// point where paid-by grants and member isolation combine, so every read surface
// (paged list, single GET, search, export, GetReceiptsForGroupIds, pie chart,
// report data) inherits both restrictions. visibleSet / memberUnrestricted come
// from one GetVisibleUserIdsForUser call per consumer; passing them in keeps the
// member-visible set from being re-resolved per group/row.
func (service PermissionService) paidByAllowedForGroup(
	userId uint,
	groupId uint,
	visibleSet map[uint]struct{},
	memberUnrestricted bool,
) (map[uint]struct{}, bool, error) {
	paidBySet, paidByUnrestricted, err := service.GetGroupPaidByUserIdsForUser(userId, groupId)
	if err != nil {
		return nil, false, err
	}
	allowed, unrestricted := foldPaidByWithVisibility(paidBySet, paidByUnrestricted, visibleSet, memberUnrestricted)
	return allowed, unrestricted, nil
}

// foldPaidByWithVisibility combines a group's paid-by allowed set with the caller's
// member-visible set. The two are independent narrowings of which receipts a member
// may see by payer, so the effective allowed set is their intersection:
//
//   - both unrestricted        -> unrestricted (nil), no filtering (backward compatible)
//   - only member restricted   -> visibleSet (paid-by adds nothing)
//   - only paid-by restricted  -> paidBySet (isolation adds nothing)
//   - both restricted          -> paidBySet ∩ visibleSet
//
// A restricted result may be empty ("see nothing"), preserved downstream by the
// IN (0) sentinel. visibleSet is treated read-only, so the restricted-only branch
// may return it directly.
func foldPaidByWithVisibility(
	paidBySet map[uint]struct{},
	paidByUnrestricted bool,
	visibleSet map[uint]struct{},
	memberUnrestricted bool,
) (map[uint]struct{}, bool) {
	if paidByUnrestricted && memberUnrestricted {
		return nil, true
	}
	if memberUnrestricted {
		return paidBySet, false
	}
	if paidByUnrestricted {
		return visibleSet, false
	}

	intersection := make(map[uint]struct{}, len(paidBySet))
	for id := range paidBySet {
		if _, ok := visibleSet[id]; ok {
			intersection[id] = struct{}{}
		}
	}
	return intersection, false
}

// uintSetToSlice flattens a set to a slice (order is irrelevant for an IN clause).
func uintSetToSlice(set map[uint]struct{}) []uint {
	ids := make([]uint, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}
