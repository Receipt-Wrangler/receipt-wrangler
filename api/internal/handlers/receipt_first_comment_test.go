package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

// The receipts table's Comment column, end to end through the paged list handler:
// the firstComment projection on each row, and the first_comment sort.

var firstCommentHandlerTime = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func seedNamedReceipt(t *testing.T, groupId uint, paidByUserId uint, name string) uint {
	t.Helper()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(5),
		Date:         time.Now(),
		GroupId:      groupId,
		PaidByUserID: paidByUserId,
		Status:       models.OPEN,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt.ID
}

func seedReceiptCommentAt(t *testing.T, receiptId uint, userId uint, text string, createdAt time.Time) {
	t.Helper()
	db := repositories.GetDB()
	comment := models.Comment{Comment: text, ReceiptId: receiptId, UserId: &userId}
	if err := db.Create(&comment).Error; err != nil {
		t.Fatalf("seed comment: %v", err)
	}
	if err := db.Model(&comment).Update("created_at", createdAt).Error; err != nil {
		t.Fatalf("stamp comment: %v", err)
	}
}

func pagedReceiptsAs(t *testing.T, userId uint, groupId uint, command commands.ReceiptPagedRequestCommand) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(command)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/receipt/group/"+utils.UintToString(groupId), strings.NewReader(string(body)))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", utils.UintToString(groupId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	GetPagedReceiptsForGroup(w, r)
	return w
}

// pagedFirstComments maps each returned receipt's name to its firstComment ("<absent>"
// when the key is missing), in response order.
func pagedFirstComments(t *testing.T, w *httptest.ResponseRecorder) ([]string, map[string]string) {
	t.Helper()
	var pagedData structs.PagedData
	if err := json.Unmarshal(w.Body.Bytes(), &pagedData); err != nil {
		t.Fatalf("unmarshal paged data: %v", err)
	}

	names := make([]string, 0, len(pagedData.Data))
	firstComments := map[string]string{}
	for _, entry := range pagedData.Data {
		row, ok := entry.(map[string]any)
		if !ok {
			t.Fatalf("unexpected row shape: %T", entry)
		}
		name, _ := row["name"].(string)
		names = append(names, name)
		firstComments[name] = "<absent>"
		if value, present := row["firstComment"]; present {
			text, _ := value.(string)
			firstComments[name] = text
		}
	}
	return names, firstComments
}

func TestPagedReceiptsCarryTheFirstComment(t *testing.T) {
	defer tearDownReceiptsTest()
	setupReceiptsTest()
	grantAllGroupPerms(t, 1, 1)

	commented := seedNamedReceipt(t, 1, 1, "commented")
	seedReceiptCommentAt(t, commented, 1, "the later one", firstCommentHandlerTime.Add(time.Hour))
	seedReceiptCommentAt(t, commented, 1, "the first one", firstCommentHandlerTime)
	seedNamedReceipt(t, 1, 1, "bare")

	w := pagedReceiptsAs(t, 1, 1, commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{Page: 1, PageSize: 10},
	})
	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	_, firstComments := pagedFirstComments(t, w)
	if firstComments["commented"] != "the first one" {
		utils.PrintTestError(t, firstComments["commented"], "the first one")
	}
	// omitempty: a receipt with no comment carries no key at all, rather than a
	// null the generated clients would have to tolerate.
	if firstComments["bare"] != "<absent>" {
		utils.PrintTestError(t, firstComments["bare"], "<absent>")
	}
}

func TestPagedReceiptsSortByFirstComment(t *testing.T) {
	defer tearDownReceiptsTest()
	setupReceiptsTest()
	grantAllGroupPerms(t, 1, 1)

	for name, text := range map[string]string{"r-charlie": "charlie", "r-alpha": "alpha", "r-bravo": "bravo"} {
		seedReceiptCommentAt(t, seedNamedReceipt(t, 1, 1, name), 1, text, firstCommentHandlerTime)
	}

	for sortDirection, want := range map[commands.SortDirection][]string{
		commands.ASCENDING:  {"r-alpha", "r-bravo", "r-charlie"},
		commands.DESCENDING: {"r-charlie", "r-bravo", "r-alpha"},
	} {
		w := pagedReceiptsAs(t, 1, 1, commands.ReceiptPagedRequestCommand{
			PagedRequestCommand: commands.PagedRequestCommand{
				Page:          1,
				PageSize:      10,
				OrderBy:       constants.FIRST_COMMENT_ORDER_BY,
				SortDirection: sortDirection,
			},
		})
		if w.Result().StatusCode != 200 {
			utils.PrintTestError(t, w.Result().StatusCode, 200)
			continue
		}

		names, _ := pagedFirstComments(t, w)
		if strings.Join(names, ",") != strings.Join(want, ",") {
			utils.PrintTestError(t, names, want)
		}
	}
}

// An isolated member must never receive a peer's comment text - not as the column's
// value, and not through the order the column sorts in.
func TestMemberIsolationPagedListNeverShowsAHiddenAuthorsFirstComment(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)

	peerFirst := seedNamedReceipt(t, fx.groupId, fx.memberAId, "peer-first")
	seedReceiptCommentAt(t, peerFirst, fx.memberBId, "aaa secret from B", firstCommentHandlerTime)
	seedReceiptCommentAt(t, peerFirst, fx.memberAId, "zzz from A", firstCommentHandlerTime.Add(time.Hour))

	own := seedNamedReceipt(t, fx.groupId, fx.memberAId, "own")
	seedReceiptCommentAt(t, own, fx.memberAId, "mmm from A", firstCommentHandlerTime)

	w := pagedReceiptsAs(t, fx.memberAId, fx.groupId, commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      10,
			OrderBy:       constants.FIRST_COMMENT_ORDER_BY,
			SortDirection: commands.ASCENDING,
		},
	})
	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	if strings.Contains(w.Body.String(), "secret from B") {
		utils.PrintTestError(t, w.Body.String(), "no trace of the hidden author's comment")
	}

	names, firstComments := pagedFirstComments(t, w)
	if firstComments["peer-first"] != "zzz from A" {
		utils.PrintTestError(t, firstComments["peer-first"], "zzz from A")
	}
	// Sorted on "zzz", not the hidden "aaa", so it comes after "mmm".
	if strings.Join(names, ",") != "own,peer-first" {
		utils.PrintTestError(t, names, []string{"own", "peer-first"})
	}
}
