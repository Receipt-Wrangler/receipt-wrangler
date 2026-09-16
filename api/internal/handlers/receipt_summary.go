package handlers

import (
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

// GetReceiptSummaryForGroup returns the block of totals rendered under the receipts
// table for the current filter.
//
// Gated on group.receipts.read — the same permission as the table it sits under, and
// deliberately not a new one. The endpoint returns nothing but aggregates of rows the
// caller can already page through under that gate, so a separate permission would only
// create a state where the table renders and its own totals 403.
//
// It is also not group.widgets.read: this is not a dashboard widget, and a member with
// receipts-read but no widgets-read should still see the totals under their own table.
//
// No app.custom-fields.read gate either. That permission covers the custom field
// CATALOG endpoints and enforceReceiptCustomFieldSelection — which explicitly lets a
// non-holder read and edit the values of fields already on a receipt — and
// FULL_RECEIPT_ASSOCIATIONS already ships those field names to any receipt reader.
// Gating here would reproduce the "silent no-op for a hand-built role" trap documented
// for default custom fields while protecting nothing.
func GetReceiptSummaryForGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")

	handler := structs.Handler{
		ErrorMessage:     "Error getting receipt summary",
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupReceiptsRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.ReceiptSummaryCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErr := command.Validate()
			if len(vErr.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
				return 0, nil
			}

			token := structs.GetClaims(r)
			receiptSummaryService := services.NewReceiptSummaryService(nil)

			summary, err := receiptSummaryService.GetReceiptSummary(token.UserId, groupId, command)
			if errors.Is(err, services.ErrConfigurationGroupForbidden) {
				utils.WriteCustomErrorResponse(
					w,
					"You do not have access to the requested configuration group",
					http.StatusForbidden,
				)
				return 0, nil
			}
			if errors.Is(err, services.ErrConfigurationGroupNotAllGroup) {
				utils.WriteCustomErrorResponse(
					w,
					"A configuration group may only be specified for the all group",
					http.StatusBadRequest,
				)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(summary)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}
