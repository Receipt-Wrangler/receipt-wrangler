package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

func HandleRequest(handler structs.Handler) {
	if len(handler.ResponseType) > 0 {
		handler.Writer.Header().Set("Content-Type", handler.ResponseType)
	}

	// Receipts resolved by id here also have their paid-by visibility checked once
	// permissions pass, so a member cannot reach (read, edit, comment on, export, or
	// download) a receipt hidden from them by their group role's paid-by filter.
	var receiptsToCheck []models.Receipt

	if len(handler.ReceiptId) > 0 {
		var receipt models.Receipt
		db := repositories.GetDB()
		err := db.Model(models.Receipt{}).Where("id = ?", handler.ReceiptId).Select("group_id", "paid_by_user_id").First(&receipt).Error
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			utils.WriteCustomErrorResponse(handler.Writer, "User is unauthorized to access entity", http.StatusForbidden)
			return
		}

		handler.GroupId = utils.UintToString(receipt.GroupId)
		receiptsToCheck = append(receiptsToCheck, receipt)
	}

	if len(handler.ReceiptIds) > 0 {
		var receipts []models.Receipt
		db := repositories.GetDB()
		err := db.Model(models.Receipt{}).Where("id IN (?)", handler.ReceiptIds).Select("group_id", "paid_by_user_id").Find(&receipts).Error
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			utils.WriteCustomErrorResponse(handler.Writer, "User is unauthorized to access entity", http.StatusForbidden)
			return
		}

		for _, receipt := range receipts {
			handler.GroupIds = append(handler.GroupIds, utils.UintToString(receipt.GroupId))
		}
		receiptsToCheck = append(receiptsToCheck, receipts...)
	}

	if !enforcePermissions(handler) {
		return
	}

	if !handler.SkipPaidByVisibilityCheck && !enforcePaidByVisibility(handler, receiptsToCheck) {
		return
	}

	errCode, err := handler.HandlerFunction(handler.Writer, handler.Request)

	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(handler.Writer, handler.ErrorMessage, errCode)
		return
	}
}

// enforcePermissions runs the modern permission checks declared on a handler.
// The caller's effective permissions are resolved from the database (never the
// JWT) by the PermissionService. It returns false — after writing a 403 — when
// the caller is not authorized. Handlers that declare no permissions are allowed
// through (authentication is already enforced by the router middleware).
func enforcePermissions(handler structs.Handler) bool {
	if len(handler.AppPermissions) == 0 && len(handler.GroupPermissions) == 0 {
		return true
	}

	token := structs.GetClaims(handler.Request)
	permissionService := services.NewPermissionService(nil)

	if len(handler.AppPermissions) > 0 {
		hasPermissions, err := permissionService.HasAppPermissions(token.UserId, handler.AppPermissions...)
		if err != nil || !hasPermissions {
			return denyUnauthorized(handler, err)
		}
	}

	if len(handler.GroupPermissions) > 0 {
		if len(handler.GroupId) == 0 && len(handler.GroupIds) == 0 {
			utils.WriteCustomErrorResponse(handler.Writer, "Group ID is required to validate group permissions", http.StatusForbidden)
			return false
		}

		// An app-scoped fallback (e.g. an administrator) can bypass the
		// group-permission check entirely.
		if len(handler.OrAppPermissions) > 0 {
			hasFallback, err := permissionService.HasAnyAppPermission(token.UserId, handler.OrAppPermissions...)
			if err != nil {
				return denyUnauthorized(handler, err)
			}
			if hasFallback {
				return true
			}
		}

		groupIds := make([]string, 0, len(handler.GroupIds)+1)
		groupIds = append(groupIds, handler.GroupIds...)
		if len(handler.GroupId) > 0 {
			groupIds = append(groupIds, handler.GroupId)
		}

		for _, groupId := range groupIds {
			uintGroupId, err := utils.StringToUint(groupId)
			if err != nil {
				return denyUnauthorized(handler, err)
			}

			hasPermissions, err := permissionService.HasGroupPermissions(token.UserId, uintGroupId, handler.GroupPermissions...)
			if err != nil || !hasPermissions {
				return denyUnauthorized(handler, err)
			}
		}
	}

	return true
}

// enforcePaidByVisibility denies (with 403) when any of the resolved receipts is
// hidden from the caller by their group role's paid-by visibility filter. An
// unrestricted caller, or a non-member (no group role), passes. This is the
// single chokepoint for every receipt-by-id surface, since they all flow through
// HandleRequest with ReceiptId/ReceiptIds set. FilterReceiptsByPaidBy resolves
// each group's allowed set at most once, so a multi-id request costs O(distinct
// groups) lookups rather than O(receipts); "any hidden" => fewer survive => deny.
func enforcePaidByVisibility(handler structs.Handler, receipts []models.Receipt) bool {
	if len(receipts) == 0 {
		return true
	}

	token := structs.GetClaims(handler.Request)
	permissionService := services.NewPermissionService(nil)

	visible, err := permissionService.FilterReceiptsByPaidBy(token.UserId, receipts)
	if err != nil {
		return denyUnauthorized(handler, err)
	}
	if len(visible) != len(receipts) {
		return denyUnauthorized(handler, nil)
	}

	return true
}

func denyUnauthorized(handler structs.Handler, err error) bool {
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
	}
	utils.WriteCustomErrorResponse(handler.Writer, "User is unauthorized to access entity", http.StatusForbidden)
	return false
}
