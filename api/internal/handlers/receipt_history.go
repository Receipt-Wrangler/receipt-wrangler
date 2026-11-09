package handlers

import (
	"errors"
	"net/http"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/stacktrace"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

// receiptHistoryHandler gates both history endpoints the same way: read and
// history-read in the receipt's group as it is now (no app-level bypass) and
// paid-by visibility on the receipt, all enforced by HandleRequest, which also
// answers a missing receipt with its 403.
func receiptHistoryHandler(
	w http.ResponseWriter,
	r *http.Request,
	errorMessage string,
	handlerFunction func(w http.ResponseWriter, r *http.Request) (int, error),
) structs.Handler {
	return structs.Handler{
		ErrorMessage:     errorMessage,
		Writer:           w,
		Request:          r,
		ReceiptId:        chi.URLParam(r, "id"),
		GroupPermissions: []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsHistoryRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction:  handlerFunction,
	}
}

func GetReceiptHistory(w http.ResponseWriter, r *http.Request) {
	HandleRequest(receiptHistoryHandler(w, r, "Error retrieving receipt history.",
		func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			receiptId, err := utils.StringToUint(chi.URLParam(r, "id"))
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			history, err := services.NewReceiptHistoryService(nil).GetReceiptHistory(token.UserId, receiptId)
			if errors.Is(err, services.ErrReceiptAccessDenied) {
				utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			return writeReceiptHistoryResponse(w, history)
		},
	))
}

func GetReceiptHistorySnapshots(w http.ResponseWriter, r *http.Request) {
	HandleRequest(receiptHistoryHandler(w, r, "Error retrieving receipt history snapshots.",
		func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			receiptId, err := utils.StringToUint(chi.URLParam(r, "id"))
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			// An id that does not parse comes back 0, which names no entry and is
			// answered like any other entry the caller cannot see.
			systemTaskId, _ := utils.StringToUint(chi.URLParam(r, "systemTaskId"))

			snapshots, err := services.NewReceiptHistoryService(nil).GetReceiptHistorySnapshots(token.UserId, receiptId, systemTaskId)
			if errors.Is(err, services.ErrReceiptAccessDenied) {
				utils.WriteCustomErrorResponse(w, unauthorizedEntityMessage, http.StatusForbidden)
				return 0, nil
			}
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			return writeReceiptHistoryResponse(w, snapshots)
		},
	))
}

func writeReceiptHistoryResponse(w http.ResponseWriter, body interface{}) (int, error) {
	bytes, err := utils.MarshalResponseData(body)
	if err != nil {
		return http.StatusInternalServerError, stacktrace.WithStack(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return 0, nil
}
