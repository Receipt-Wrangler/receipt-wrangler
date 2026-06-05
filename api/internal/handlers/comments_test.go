package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

func tearDownCommentsTest() {
	repositories.TruncateTestDb()
	services.ClearRolePermissionCacheForTests()
}

// deleteCommentRequest builds a request carrying the commentId chi URL param and
// JWT claims for userId, mirroring how the router invokes DeleteComment.
func deleteCommentRequest(userId uint, commentId string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api/comment/"+commentId, strings.NewReader(""))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("commentId", commentId)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}})
	return w, r.WithContext(newContext)
}

// seedReceiptComment creates a receipt in groupId and a comment on it owned by
// userId. Both get id 1 (the test db is truncated between tests).
func seedReceiptComment(userId uint, groupId uint) {
	db := repositories.GetDB()
	db.Create(&models.Receipt{Name: "Test receipt", GroupId: groupId, PaidByUserID: userId})
	owner := userId
	db.Create(&models.Comment{Comment: "test comment", ReceiptId: 1, UserId: &owner})
}

func TestDeleteCommentAllowsWhenUserHasGroupPermission(t *testing.T) {
	defer tearDownCommentsTest()
	repositories.CreateTestGroupWithUsers()
	seedReceiptComment(1, 1)
	grantGroupPerms(t, 1, 1, permissions.GroupCommentsDelete)

	w, r := deleteCommentRequest(1, "1")
	DeleteComment(w, r)

	assertStatus(t, w, http.StatusOK)
}

func TestDeleteCommentRejectsWhenUserLacksGroupPermission(t *testing.T) {
	defer tearDownCommentsTest()
	repositories.CreateTestGroupWithUsers()
	seedReceiptComment(1, 1)
	// user 1 is a member of group 1 but holds no role granting comments.delete

	w, r := deleteCommentRequest(1, "1")
	DeleteComment(w, r)

	assertStatus(t, w, http.StatusForbidden)
}
