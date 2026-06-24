package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

func Search(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error searching",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppReceiptsSearch},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			searchTerm := r.URL.Query().Get("searchTerm")

			if len(searchTerm) > 0 {
				searchTerm = "%" + searchTerm + "%"

				db := repositories.GetDB()
				var receipts []models.Receipt

				results := make([]structs.SearchResult, 0)

				token := structs.GetClaims(r)
				groupMemberRepository := repositories.NewGroupMemberRepository(nil)
				groupIds, err := groupMemberRepository.GetGroupIdsByUserId(utils.UintToString(token.UserId))
				if err != nil {
					return http.StatusInternalServerError, err
				}

				query := db.Table("receipts").Where("group_id IN ? AND name LIKE ?", groupIds, searchTerm)

				// Apply the caller's paid-by visibility in SQL BEFORE the limit —
				// SearchResult exposes paidByUserId, and a post-fetch filter would drop
				// visible matches whenever hidden receipts fill the first 100 rows.
				receiptRepository := repositories.NewReceiptRepository(nil)
				query, err = receiptRepository.ApplyPaidByDisjunction(
					query,
					groupIds,
					services.NewPermissionService(nil).PaidByListResolver(token.UserId),
				)
				if err != nil {
					return http.StatusInternalServerError, err
				}

				err = query.Limit(100).Order("date desc").Find(&receipts).Error
				if err != nil {
					return http.StatusInternalServerError, err
				}

				for _, receipt := range receipts {
					results = append(results, structs.SearchResult{
						ID:            receipt.ID,
						GroupID:       receipt.GroupId,
						Name:          receipt.Name,
						Date:          receipt.Date,
						Type:          "Receipt",
						Amount:        receipt.Amount,
						ReceiptStatus: receipt.Status,
						PaidByUserId:  receipt.PaidByUserID,
						CreatedAt:     receipt.CreatedAt,
					})
				}

				bytes, err := utils.MarshalResponseData(results)
				if err != nil {
					return http.StatusInternalServerError, err
				}

				w.WriteHeader(200)
				w.Write(bytes)
			} else {
				results := make([]structs.SearchResult, 0)
				bytes, err := utils.MarshalResponseData(results)
				if err != nil {
					return http.StatusInternalServerError, err
				}

				w.WriteHeader(200)
				w.Write(bytes)
			}

			return 0, nil
		},
	}

	HandleRequest(handler)
}
