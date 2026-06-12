package middleware

import (
	"net/http"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

func CanDeleteGroup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := structs.GetClaims(r)
		errMsg := "User must be a part of at least one group."

		groupMemberRepository := repositories.NewGroupMemberRepository(nil)
		groupMembers, err := groupMemberRepository.GetGroupMembersByUserId(utils.UintToString(token.UserId))
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			utils.WriteCustomErrorResponse(w, errMsg, http.StatusInternalServerError)
			return
		}

		if len(groupMembers) <= 1 {
			logging.LogStd(logging.LOG_LEVEL_ERROR, errMsg, r)
			utils.WriteCustomErrorResponse(w, errMsg, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
