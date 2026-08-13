package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
)

// seedQuickScanCommenter creates a group with the given quick-scan comment configuration plus a
// member who can quick scan, optionally holding group.comments.create. It returns the user id and
// group id. Paid-by/status are seeded shown+optional with defaults so only the comment field can
// produce a config error.
func seedQuickScanCommenter(t *testing.T, canComment bool, settings models.GroupReceiptSettings) (uint, uint) {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	db := repositories.GetDB()

	group := models.Group{Name: "quick-scan-comment-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	settingsRepository := repositories.NewGroupReceiptSettingsRepository(nil)
	if _, err := settingsRepository.CreateGroupReceiptSettings(group.ID); err != nil {
		t.Fatalf("seed group receipt settings: %v", err)
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		HideComments:               settings.HideComments,
		QuickScanPaidByEnabled:     true,
		QuickScanPaidByRequired:    false,
		QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanStatusEnabled:     true,
		QuickScanStatusRequired:    false,
		QuickScanDefaultStatus:     models.OPEN,
		QuickScanCommentEnabled:    settings.QuickScanCommentEnabled,
		QuickScanCommentRequired:   settings.QuickScanCommentRequired,
	}
	if _, err := settingsRepository.UpdateGroupReceiptSettings(utils.UintToString(group.ID), command); err != nil {
		t.Fatalf("update group receipt settings: %v", err)
	}

	rolePermissions := []string{permissions.GroupReceiptsQuickScan}
	if canComment {
		rolePermissions = append(rolePermissions, permissions.GroupCommentsCreate)
	}

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"Quick Scan Commenter", "", rolePermissions, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: "quick-scan-commenter", Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	return user.ID, group.ID
}

// commentCommand builds a single-file quick scan command for the given group carrying one comment.
func commentCommand(groupId uint, comment string) commands.QuickScanCommand {
	return commands.QuickScanCommand{
		Files:         []multipart.File{nil},
		GroupIds:      []uint{groupId},
		PaidByUserIds: []uint{0},
		Statuses:      []models.ReceiptStatus{""},
		Comments:      []string{comment},
	}
}

func TestResolveQuickScanFields_RequiredCommentRejectedWhenEmpty(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		QuickScanCommentEnabled:  true,
		QuickScanCommentRequired: true,
	})

	_, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, ""), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if _, ok := configErr.Errors["files.0.comment"]; !ok {
		utils.PrintTestError(t, configErr.Errors, "files.0.comment")
	}
}

func TestResolveQuickScanFields_CommentKept(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		QuickScanCommentEnabled:  true,
		QuickScanCommentRequired: true,
	})

	resolved, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, "Lunch with the team"), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}
	if resolved[0].Comment != "Lunch with the team" {
		utils.PrintTestError(t, resolved[0].Comment, "Lunch with the team")
	}
}

// A comment supplied while the field is disabled is dropped, not persisted - the group turned the
// field off, so there is nowhere for the value to have legitimately come from.
func TestResolveQuickScanFields_CommentDroppedWhenDisabled(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		QuickScanCommentEnabled:  false,
		QuickScanCommentRequired: true,
	})

	resolved, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, "should be dropped"), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}
	if resolved[0].Comment != "" {
		utils.PrintTestError(t, resolved[0].Comment, "")
	}
}

// hideComments hides comments for the whole group, so it overrides an enabled+required quick-scan
// comment: no required error, and any submitted comment is dropped.
func TestResolveQuickScanFields_HideCommentsOverridesConfig(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		HideComments:             true,
		QuickScanCommentEnabled:  true,
		QuickScanCommentRequired: true,
	})

	resolved, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, "hidden by group settings"), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}
	if resolved[0].Comment != "" {
		utils.PrintTestError(t, resolved[0].Comment, "")
	}
}

// Without group.comments.create the field is treated as hidden: the required check is skipped (so
// the user is not locked out of quick scan) and a submitted comment is dropped rather than 403'd.
func TestResolveQuickScanFields_CommentSkippedWithoutPermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, false, models.GroupReceiptSettings{
		QuickScanCommentEnabled:  true,
		QuickScanCommentRequired: true,
	})

	resolved, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, ""), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}

	resolved, configErr, err = services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, "not allowed"), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if len(configErr.Errors) > 0 {
		utils.PrintTestError(t, configErr.Errors, "no errors")
	}
	if resolved[0].Comment != "" {
		utils.PrintTestError(t, resolved[0].Comment, "")
	}
}

// An over-length comment is rejected synchronously. The receipt is created in a background task, so
// letting it through would fail the insert where the user can never see it.
func TestResolveQuickScanFields_CommentTooLongRejected(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		QuickScanCommentEnabled: true,
	})

	tooLong := strings.Repeat("a", models.MaxCommentLength+1)
	_, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, tooLong), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if _, ok := configErr.Errors["files.0.comment"]; !ok {
		utils.PrintTestError(t, configErr.Errors, "files.0.comment")
	}
}

// The limit counts characters, not bytes. The Comment column is varchar(500), which MySQL and
// Postgres both measure in characters, so a byte-based check would reject an accented or non-Latin
// comment the column has ample room for. 500 two-byte runes is 1000 bytes: accepted here, rejected
// by a len() check.
func TestResolveQuickScanFields_CommentLengthCountsCharactersNotBytes(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		QuickScanCommentEnabled: true,
	})

	atLimit := strings.Repeat("é", models.MaxCommentLength)
	resolved, configErr, err := services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, atLimit), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if _, ok := configErr.Errors["files.0.comment"]; ok {
		utils.PrintTestError(t, configErr.Errors, "no comment error")
	}
	if resolved[0].Comment != atLimit {
		utils.PrintTestError(t, resolved[0].Comment, atLimit)
	}

	overLimit := strings.Repeat("é", models.MaxCommentLength+1)
	_, configErr, err = services.NewReceiptService(nil).ResolveQuickScanFields(commentCommand(groupId, overLimit), userId)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
	}
	if _, ok := configErr.Errors["files.0.comment"]; !ok {
		utils.PrintTestError(t, configErr.Errors, "files.0.comment")
	}
}

// quickScanCommentRequest builds a real multipart quick-scan request for the group, omitting the
// comments field entirely when comment is nil (what a client released before the field does).
func quickScanCommentRequest(t *testing.T, userId uint, groupId uint, comment *string) *httptest.ResponseRecorder {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "receipt.png")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	part.Write([]byte("not-a-real-image"))

	writer.WriteField("groupIds", utils.UintToString(groupId))
	writer.WriteField("paidByUserIds", utils.UintToString(userId))
	writer.WriteField("statuses", string(models.OPEN))
	if comment != nil {
		writer.WriteField("comments", *comment)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	QuickScan(w, r)
	return w
}

// End-to-end through the handler: a required comment is enforced before anything is enqueued, so a
// client that omits it (including any released before the field existed) gets a 400.
func TestQuickScanHandlerRejectsMissingRequiredComment(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{
		QuickScanCommentEnabled:  true,
		QuickScanCommentRequired: true,
	})

	w := quickScanCommentRequest(t, userId, groupId, nil)
	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

// The same request from a member without group.comments.create must NOT 400 - the field is hidden
// for them, so requiring it would lock them out of quick scan entirely. It gets past the config
// check (failing later on the fake image, which is not what this asserts).
func TestQuickScanHandlerSkipsRequiredCommentWithoutPermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, false, models.GroupReceiptSettings{
		QuickScanCommentEnabled:  true,
		QuickScanCommentRequired: true,
	})

	w := quickScanCommentRequest(t, userId, groupId, nil)
	if w.Result().StatusCode == 400 {
		utils.PrintTestError(t, w.Result().StatusCode, "not 400")
	}
}

// A comment from a member without the permission is dropped rather than refused, so the scan is
// never 403'd over it.
func TestQuickScanHandlerDoesNotForbidCommentWithoutPermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, false, models.GroupReceiptSettings{
		QuickScanCommentEnabled: true,
	})

	comment := "submitted without permission"
	w := quickScanCommentRequest(t, userId, groupId, &comment)
	if w.Result().StatusCode == 403 {
		utils.PrintTestError(t, w.Result().StatusCode, "not 403")
	}
}
