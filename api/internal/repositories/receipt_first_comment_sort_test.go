package repositories

import (
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Sorting the receipts table by its Comment column: each receipt's first comment,
// read through a correlated subquery restricted to authors the caller may see.
//
// Like the custom field sort, the suite is SQLite only, so where a receipt with no
// comment lands (NULL ordering is engine-dependent) is deliberately not asserted.

var firstCommentBaseTime = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

// addCommentAt stamps created_at explicitly: "first" means earliest, and comments
// seeded back-to-back would otherwise share a timestamp and leave the order to the
// id tiebreaker alone.
func addCommentAt(receiptId uint, userId *uint, text string, createdAt time.Time) models.Comment {
	comment := models.Comment{Comment: text, ReceiptId: receiptId, UserId: userId}
	GetDB().Create(&comment)
	GetDB().Model(&comment).Update("created_at", createdAt)

	return comment
}

func sortByFirstComment(sortDirection commands.SortDirection) commands.ReceiptPagedRequestCommand {
	return commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      50,
			OrderBy:       constants.FIRST_COMMENT_ORDER_BY,
			SortDirection: sortDirection,
		},
	}
}

// firstCommentNames runs the sort and returns the receipt names in the order the
// query produced them, leaving out the receipts named in withoutComment - they sort
// as NULL, whose position depends on the engine - after checking they are still
// listed and that the count agrees with the rows.
func firstCommentNames(
	t *testing.T,
	userId uint,
	groupId uint,
	sortDirection commands.SortDirection,
	resolver CommentAuthorVisibilityResolver,
	withoutComment ...string,
) []string {
	repository := NewReceiptRepository(nil)

	receipts, count, err := repository.GetPagedReceiptsByGroupId(
		userId, utils.UintToString(groupId), sortByFirstComment(sortDirection), nil, nil, resolver, nil, nil,
	)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return nil
	}
	if int(count) != len(receipts) {
		utils.PrintTestError(t, count, len(receipts))
	}

	names := make([]string, 0, len(receipts))
	seen := map[string]bool{}
	for _, receipt := range receipts {
		seen[receipt.Name] = true
		if !utils.Contains(stringsToAny(withoutComment), receipt.Name) {
			names = append(names, receipt.Name)
		}
	}
	for _, name := range withoutComment {
		if !seen[name] {
			utils.PrintTestError(t, name, "a receipt without a comment still listed")
		}
	}

	return names
}

func stringsToAny(values []string) []interface{} {
	result := make([]interface{}, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}

func TestShouldSortReceiptsByFirstComment(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	// A later comment that would sort first must not decide the order: the
	// receipt sorts by its earliest comment, not its lowest or its latest.
	banana := createSortableReceipt("r-banana")
	addCommentAt(banana.ID, uintPtr(1), "banana", firstCommentBaseTime)
	addCommentAt(banana.ID, uintPtr(1), "aardvark", firstCommentBaseTime.Add(time.Hour))

	cherry := createSortableReceipt("r-cherry")
	addCommentAt(cherry.ID, uintPtr(1), "cherry", firstCommentBaseTime)

	// Written earlier than its id suggests: created_at, not insertion order, wins.
	apple := createSortableReceipt("r-apple")
	addCommentAt(apple.ID, uintPtr(1), "zucchini", firstCommentBaseTime.Add(time.Hour))
	addCommentAt(apple.ID, uintPtr(1), "apple", firstCommentBaseTime.Add(-time.Hour))

	createSortableReceipt("r-none")

	assertOrder(
		t,
		firstCommentNames(t, 1, 1, commands.ASCENDING, nil, "r-none"),
		[]string{"r-apple", "r-banana", "r-cherry"},
	)
	assertOrder(
		t,
		firstCommentNames(t, 1, 1, commands.DESCENDING, nil, "r-none"),
		[]string{"r-cherry", "r-banana", "r-apple"},
	)
}

// Comments written in one insert share a created_at; the id breaks the tie, the
// same rule the column displays by.
func TestFirstCommentTieOnCreatedAtIsBrokenById(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	early := createSortableReceipt("r-zebra-first")
	addCommentAt(early.ID, uintPtr(1), "zebra", firstCommentBaseTime)
	addCommentAt(early.ID, uintPtr(1), "aaa", firstCommentBaseTime)

	other := createSortableReceipt("r-mango")
	addCommentAt(other.ID, uintPtr(1), "mango", firstCommentBaseTime)

	assertOrder(
		t,
		firstCommentNames(t, 1, 1, commands.ASCENDING, nil),
		[]string{"r-mango", "r-zebra-first"},
	)
}

// Under member isolation a comment by an author the caller cannot see is dropped
// from every response, so it must not order the table either - the receipt sorts
// by its first comment the caller CAN see. A comment with no author names no one
// and stays a candidate.
func TestFirstCommentSortSkipsAuthorsTheCallerCannotSee(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	hiddenFirst := createSortableReceipt("r-hidden-first")
	addCommentAt(hiddenFirst.ID, uintPtr(2), "aaa", firstCommentBaseTime)
	addCommentAt(hiddenFirst.ID, uintPtr(1), "zzz", firstCommentBaseTime.Add(time.Hour))

	visible := createSortableReceipt("r-visible")
	addCommentAt(visible.ID, uintPtr(1), "mmm", firstCommentBaseTime)

	authorless := createSortableReceipt("r-authorless")
	addCommentAt(authorless.ID, nil, "bbb", firstCommentBaseTime)

	onlyHidden := createSortableReceipt("r-only-hidden")
	addCommentAt(onlyHidden.ID, uintPtr(3), "abc", firstCommentBaseTime)

	restrictedToSelf := func(groupId uint) ([]uint, bool, error) {
		return []uint{1}, false, nil
	}

	assertOrder(
		t,
		firstCommentNames(t, 1, 1, commands.ASCENDING, restrictedToSelf, "r-only-hidden"),
		[]string{"r-authorless", "r-visible", "r-hidden-first"},
	)

	// The contrast: an unrestricted caller sorts on the hidden comments too.
	unrestricted := func(groupId uint) ([]uint, bool, error) {
		return nil, true, nil
	}
	assertOrder(
		t,
		firstCommentNames(t, 1, 1, commands.ASCENDING, unrestricted),
		[]string{"r-hidden-first", "r-only-hidden", "r-authorless", "r-visible"},
	)
}

// The All group spans several groups, and each receipt is judged by its own
// group's visibility: the same author can be hidden in an isolated group and
// visible in an open one.
func TestFirstCommentSortAppliesEachGroupsVisibilityOnTheAllGroup(t *testing.T) {
	defer TruncateTestDb()
	db := GetDB()

	isolated := models.Group{Name: "fc-isolated"}
	open := models.Group{Name: "fc-open"}
	allGroup := models.Group{Name: "fc-all", IsAllGroup: true}
	db.Create(&isolated)
	db.Create(&open)
	db.Create(&allGroup)

	member := models.User{Username: "fc-member", Password: "x"}
	peer := models.User{Username: "fc-peer", Password: "x"}
	db.Create(&member)
	db.Create(&peer)

	db.Create(&models.GroupMember{GroupID: isolated.ID, UserID: member.ID})
	db.Create(&models.GroupMember{GroupID: open.ID, UserID: member.ID})

	makeReceipt := func(groupId uint, name string) models.Receipt {
		receipt := models.Receipt{
			Name:         name,
			Amount:       decimal.NewFromInt(1),
			Date:         time.Now(),
			PaidByUserID: member.ID,
			GroupId:      groupId,
			Status:       models.OPEN,
		}
		db.Create(&receipt)
		return receipt
	}

	inIsolated := makeReceipt(isolated.ID, "isolated-receipt")
	addCommentAt(inIsolated.ID, &peer.ID, "aaa", firstCommentBaseTime)
	addCommentAt(inIsolated.ID, &member.ID, "yyy", firstCommentBaseTime.Add(time.Hour))

	inOpen := makeReceipt(open.ID, "open-receipt")
	addCommentAt(inOpen.ID, &peer.ID, "bbb", firstCommentBaseTime)

	resolver := func(groupId uint) ([]uint, bool, error) {
		if groupId == isolated.ID {
			return []uint{member.ID}, false, nil
		}
		return nil, true, nil
	}

	// isolated-receipt sorts by "yyy" (its peer's "aaa" is hidden there), while
	// open-receipt keeps the same peer's "bbb".
	assertOrder(
		t,
		firstCommentNames(t, member.ID, allGroup.ID, commands.ASCENDING, resolver),
		[]string{"open-receipt", "isolated-receipt"},
	)
}

// Same reasoning as TestFallbackOrderCarriesTheReceiptIdTiebreaker: many receipts
// have no comment at all, so without a unique last term paging repeats and skips
// rows. Asserted on the SQL because a tie is broken by whatever order the engine
// happens to return.
func TestFirstCommentSortCarriesTheReceiptIdTiebreaker(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	repository := NewReceiptRepository(nil)
	restricted := func(groupId uint) ([]uint, bool, error) {
		return []uint{1}, false, nil
	}

	for _, sortDirection := range []commands.SortDirection{
		commands.ASCENDING, commands.DESCENDING, commands.DEFAULT,
	} {
		sql := normalizeSQL(GetDB().ToSQL(func(tx *gorm.DB) *gorm.DB {
			query, err := repository.orderByFirstComment(
				tx.Model(&models.Receipt{}), []uint{1}, restricted, sortDirection,
			)
			if err != nil {
				utils.PrintTestError(t, err, nil)
				return tx
			}

			var receipts []models.Receipt
			return query.Find(&receipts)
		}))

		for _, fragment := range []string{
			"comments.receipt_id = receipts.id",
			"comments.user_id IS NULL OR comments.user_id IN",
			"LIMIT 1",
			"receipts.id DESC",
		} {
			if !strings.Contains(sql, fragment) {
				utils.PrintTestError(t, sql, "an ORDER BY carrying "+fragment)
			}
		}
	}
}

// A caller no group restricts gets no visibility predicate at all, so a
// non-isolated install sorts exactly as it would with no resolver.
func TestFirstCommentSortAddsNoPredicateWhenUnrestricted(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	repository := NewReceiptRepository(nil)
	unrestricted := func(groupId uint) ([]uint, bool, error) {
		return nil, true, nil
	}

	sql := normalizeSQL(GetDB().ToSQL(func(tx *gorm.DB) *gorm.DB {
		query, err := repository.orderByFirstComment(
			tx.Model(&models.Receipt{}), []uint{1, 2}, unrestricted, commands.ASCENDING,
		)
		if err != nil {
			utils.PrintTestError(t, err, nil)
			return tx
		}

		var receipts []models.Receipt
		return query.Find(&receipts)
	}))

	if strings.Contains(sql, "comments.user_id") || strings.Contains(sql, "receipts.group_id") {
		utils.PrintTestError(t, sql, "no author visibility predicate")
	}
}

func TestCommentReceiptIdIndexExists(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	if !GetDB().Migrator().HasIndex(&models.Comment{}, "idx_comment_receipt_id") {
		utils.PrintTestError(t, "index missing", "idx_comment_receipt_id created by AutoMigrate")
	}
}

func TestGetCommentsForReceiptIdsReturnsFirstCommentOrder(t *testing.T) {
	defer teardownReceiptTest()
	setupReceiptTest()

	wanted := createSortableReceipt("wanted")
	addCommentAt(wanted.ID, uintPtr(1), "second", firstCommentBaseTime.Add(time.Hour))
	addCommentAt(wanted.ID, uintPtr(1), "first", firstCommentBaseTime)

	other := createSortableReceipt("other")
	addCommentAt(other.ID, uintPtr(1), "not requested", firstCommentBaseTime)

	repository := NewCommentRepository(nil)
	comments, err := repository.GetCommentsForReceiptIds([]uint{wanted.ID})
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	got := make([]string, len(comments))
	for i, comment := range comments {
		got[i] = comment.Comment
	}
	assertOrder(t, got, []string{"first", "second"})

	empty, err := repository.GetCommentsForReceiptIds(nil)
	if err != nil || empty == nil || len(empty) != 0 {
		utils.PrintTestError(t, empty, "an empty, non-nil slice")
	}
}
