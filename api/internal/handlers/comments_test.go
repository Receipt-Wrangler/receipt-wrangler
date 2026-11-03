package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
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

func addCommentRequest(userId uint, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/comment/", strings.NewReader(body))
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: userId}})
	return w, r.WithContext(newContext)
}

func TestAddCommentRecordsAReceiptChange(t *testing.T) {
	defer tearDownCommentsTest()
	userId, groupId := seedGroupWithRole(t, "comment-recorded", []string{permissions.GroupCommentsCreate})
	receiptId := seedReceipt(t, groupId, userId, 0)

	w, r := addCommentRequest(userId, fmt.Sprintf(`{"comment":"Paid in cash","receiptId":%d}`, receiptId))
	AddComment(w, r)
	assertStatus(t, w, http.StatusOK)

	before, after := assertOneRecordedChange(t, receiptId, userId)
	if len(before.Comments) != 0 || len(after.Comments) != 1 || after.Comments[0].Comment != "Paid in cash" {
		t.Errorf("recorded %d -> %d comments, want the new comment in the after alone", len(before.Comments), len(after.Comments))
	}
}

func TestDeleteCommentRecordsAReceiptChange(t *testing.T) {
	defer tearDownCommentsTest()
	repositories.CreateTestGroupWithUsers()
	seedReceiptComment(1, 1)
	grantGroupPerms(t, 1, 1, permissions.GroupCommentsDelete)

	w, r := deleteCommentRequest(1, "1")
	DeleteComment(w, r)
	assertStatus(t, w, http.StatusOK)

	before, after := assertOneRecordedChange(t, 1, 1)
	if len(before.Comments) != 1 || len(after.Comments) != 0 {
		t.Errorf("recorded %d -> %d comments, want 1 -> 0", len(before.Comments), len(after.Comments))
	}
}

// The role-requirement refusal is a 400 before anything changes, so it is not a
// failed change and records nothing.
func TestDeleteCommentRefusedByRoleRequirementRecordsNothing(t *testing.T) {
	defer tearDownCommentsTest()
	userId, groupId := seedRequirementGroup(t, "comment-refused", true, false)
	receiptId := seedReceipt(t, groupId, userId, 0)
	only := addStoredComment(t, receiptId, userId, "the only one")

	w, r := deleteCommentRequest(userId, utils.UintToString(only.ID))
	DeleteComment(w, r)
	assertStatus(t, w, http.StatusBadRequest)

	if tasks := recordedReceiptChanges(t, receiptId); len(tasks) != 0 {
		t.Errorf("recorded %+v, want nothing for a refused delete", tasks)
	}
}

func TestDeleteCommentThatFailsRecordsAFailedChange(t *testing.T) {
	defer tearDownCommentsTest()
	repositories.CreateTestGroupWithUsers()
	seedReceiptComment(2, 1)
	grantGroupPerms(t, 1, 1, permissions.GroupCommentsDelete)

	// User 1 may delete comments in the group, but this one is user 2's.
	w, r := deleteCommentRequest(1, "1")
	DeleteComment(w, r)
	assertStatus(t, w, http.StatusInternalServerError)

	tasks := recordedReceiptChanges(t, 1)
	if len(tasks) != 1 || tasks[0].Status != models.SYSTEM_TASK_FAILED ||
		!strings.Contains(tasks[0].ResultDescription, "another user's comment") || *tasks[0].GroupId != 1 {
		t.Errorf("recorded %+v, want one failed change naming the refusal", tasks)
	}
}
