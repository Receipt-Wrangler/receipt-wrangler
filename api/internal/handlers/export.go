package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

func ExportAllReceiptsFromGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	handler := structs.Handler{
		ErrorMessage:     "Error exporting receipts",
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationZip,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			pagedRequest := commands.ReceiptPagedRequestCommand{}
			err := pagedRequest.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErrs := pagedRequest.Validate()
			if len(vErrs.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErrs, http.StatusBadRequest)
				return 0, nil
			}

			token := structs.GetClaims(r)
			permissionService := services.NewPermissionService(nil)

			uintGroupId, err := utils.StringToUint(groupId)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			err = permissionService.IntersectReceiptFilterWithGrants(token.UserId, uintGroupId, &pagedRequest.Filter)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			receiptRepository := repositories.NewReceiptRepository(nil)
			receipts, _, err := receiptRepository.
				GetPagedReceiptsByGroupId(
					token.UserId,
					groupId,
					pagedRequest,
					getExportReceiptAssociations(),
					permissionService.PaidByListResolver(token.UserId),
				)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			err = permissionService.FilterReceiptCategoriesTags(token.UserId, receipts)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Mask user references (created-by, item charged-to) the caller may
			// not see before they are rendered into the CSV.
			err = permissionService.MaskReceiptsForMemberVisibility(token.UserId, receipts)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			receiptCsvService := services.NewReceiptCsvService()
			zip, err := receiptCsvService.GetZippedCsvFiles(receipts)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
			w.WriteHeader(http.StatusOK)
			w.Write(zip)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

func ExportReceiptsById(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	receiptIds := r.Form["receiptIds"]

	handler := structs.Handler{
		ErrorMessage:     "Error exporting receipts",
		Writer:           w,
		Request:          r,
		ResponseType:     constants.ApplicationZip,
		ReceiptIds:       receiptIds,
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			if err != nil {
				return http.StatusInternalServerError, err
			}

			token := structs.GetClaims(r)
			receiptRepository := repositories.NewReceiptRepository(nil)
			receipts, err := receiptRepository.GetReceiptsByIds(receiptIds, getExportReceiptAssociations())
			if err != nil {
				return http.StatusInternalServerError, err
			}

			permissionService := services.NewPermissionService(nil)
			err = permissionService.FilterReceiptCategoriesTags(token.UserId, receipts)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Mask user references (created-by, item charged-to) the caller may
			// not see before they are rendered into the CSV.
			err = permissionService.MaskReceiptsForMemberVisibility(token.UserId, receipts)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			receiptCsvService := services.NewReceiptCsvService()
			zip, err := receiptCsvService.GetZippedCsvFiles(receipts)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
			w.WriteHeader(http.StatusOK)
			w.Write(zip)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

func getExportReceiptAssociations() []string {
	return []string{
		"PaidByUser",
		"ReceiptItems",
		"ReceiptItems.Categories",
		"ReceiptItems.Tags",
		"ReceiptItems.ChargedToUser",
		"ReceiptItems.Receipt",
	}
}
