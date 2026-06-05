package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	token := structs.GetClaims(r)
	errMsg := "Error adding comment"

	upsertCommentCommand := commands.UpsertCommentCommand{}
	err := upsertCommentCommand.LoadDataFromRequest(w, r, true)
	if err != nil {
		utils.WriteCustomErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	vErr := upsertCommentCommand.Validate(token.UserId, true)
	if len(vErr.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
		return
	}

	receiptString := utils.UintToString(upsertCommentCommand.ReceiptId)

	handler := structs.Handler{
		ErrorMessage:     "Error adding comment",
		Writer:           w,
		Request:          r,
		ReceiptId:        receiptString,
		GroupPermissions: []string{permissions.GroupCommentsCreate},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			commentRepository := repositories.NewCommentRepository(nil)

			comment, err := commentRepository.AddComment(upsertCommentCommand)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			bytes, err := utils.MarshalResponseData(&comment)
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

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	errMessage := "Error deleting comment"
	commentId := chi.URLParam(r, "commentId")

	// Resolve the comment's receipt so the request can be gated by the comment's
	// group; the delete itself remains scoped to the caller's own comment.
	var comment models.Comment
	db := repositories.GetDB()
	err := db.Model(models.Comment{}).Where("id = ?", commentId).Select("receipt_id").First(&comment).Error
	if err != nil {
		utils.WriteCustomErrorResponse(w, errMessage, http.StatusInternalServerError)
		return
	}

	handler := structs.Handler{
		ErrorMessage:     errMessage,
		Writer:           w,
		Request:          r,
		ReceiptId:        utils.UintToString(comment.ReceiptId),
		GroupPermissions: []string{permissions.GroupCommentsDelete},
		ResponseType:     "",
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			commentRepository := repositories.NewCommentRepository(nil)
			token := structs.GetClaims(r)

			err := commentRepository.DeleteComment(commentId, token.UserId)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			return 0, nil
		},
	}

	HandleRequest(handler)
}
