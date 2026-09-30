package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
)

// requirementPerms is what a member needs to exercise every requirement surface:
// create/update receipts, and add/delete comments (the comment requirement is
// waived without group.comments.create).
var requirementPerms = []string{
	permissions.GroupReceiptsRead,
	permissions.GroupReceiptsCreate,
	permissions.GroupReceiptsUpdate,
	permissions.GroupCommentsCreate,
	permissions.GroupCommentsDelete,
}

// seedRequirementGroup creates a group and a member whose group role holds perms
// and the given receipt requirements. The group's data directory is removed
// afterwards. Returns (userId, groupId).
func seedRequirementGroup(t *testing.T, name string, requireComment bool, requireImage bool) (uint, uint) {
	t.Helper()
	userId, groupId := seedGroupWithRole(t, name, requirementPerms)
	setMemberRequirements(t, userId, groupId, requireComment, requireImage)
	cleanUpGroupDir(t, groupId)
	return userId, groupId
}

func setMemberRequirements(t *testing.T, userId uint, groupId uint, requireComment bool, requireImage bool) {
	t.Helper()
	roleId, err := repositories.NewRoleRepository(nil).GetGroupMemberRoleId(userId, groupId)
	if err != nil || roleId == nil {
		t.Fatalf("resolve member role: %v", err)
	}
	if err := repositories.NewRoleRepository(nil).SetGroupRoleReceiptRequirements(*roleId, requireComment, requireImage); err != nil {
		t.Fatalf("set requirements: %v", err)
	}
}

func cleanUpGroupDir(t *testing.T, groupId uint) {
	t.Helper()
	groupPath, err := repositories.NewFileRepository(nil).BuildGroupPath(groupId, "")
	if err != nil {
		t.Fatalf("BuildGroupPath: %v", err)
	}
	t.Cleanup(func() { utils.RemoveAllInDataDir(groupPath) })
}

func readHandlerTestJpg(t *testing.T) []byte {
	t.Helper()
	jpg, err := os.ReadFile(filepath.Join("..", "..", "testing", "test.jpg"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return jpg
}

func receiptJson(groupId uint, paidBy uint, comment string) string {
	comments := ""
	if comment != "" {
		comments = fmt.Sprintf(`,"comments":[{"comment":%q,"userId":%d}]`, comment, paidBy)
	}
	return receiptJsonWith(groupId, paidBy, comments)
}

// receiptJsonWith builds a valid receipt body with extra appended verbatim as
// further top-level fields (it must start with a comma).
func receiptJsonWith(groupId uint, paidBy uint, extra string) string {
	return fmt.Sprintf(
		`{"name":"R","amount":"5","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"%s}`,
		groupId, paidBy, extra,
	)
}

type uploadFile struct {
	name  string
	bytes []byte
}

// createWithFilesRequest builds a CreateReceiptWithFiles request. With
// receiptAsFilePart the receipt JSON rides as a file part with an
// application/json content type (how some generated clients encode a model);
// otherwise as a plain form value.
func createWithFilesRequest(t *testing.T, userId uint, receipt string, receiptAsFilePart bool, files ...uploadFile) *httptest.ResponseRecorder {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if receiptAsFilePart {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="receipt"; filename="receipt.json"`)
		header.Set("Content-Type", "application/json")
		part, err := writer.CreatePart(header)
		if err != nil {
			t.Fatalf("create receipt part: %v", err)
		}
		part.Write([]byte(receipt))
	} else {
		writer.WriteField("receipt", receipt)
	}

	for _, file := range files {
		part, err := writer.CreateFormFile("files", file.name)
		if err != nil {
			t.Fatalf("create file part: %v", err)
		}
		part.Write(file.bytes)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	return rawCreateWithFilesRequest(userId, body, writer.FormDataContentType())
}

// rawCreateWithFilesRequest sends body to CreateReceiptWithFiles as userId with
// the given content type, for bodies createWithFilesRequest cannot build.
func rawCreateWithFilesRequest(userId uint, body *bytes.Buffer, contentType string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", body)
	r.Header.Set("Content-Type", contentType)
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	CreateReceiptWithFiles(w, r)
	return w
}

// assertNothingWritten fails when a rejected create left any receipt, image row,
// comment, or image file under groupId's data directory behind.
func assertNothingWritten(t *testing.T, groupId uint) {
	t.Helper()
	for name, model := range map[string]interface{}{
		"receipts": &models.Receipt{}, "file data": &models.FileData{}, "comments": &models.Comment{},
	} {
		if n := countRows(t, model); n != 0 {
			t.Errorf("%s = %d, want 0", name, n)
		}
	}

	groupPath, err := repositories.NewFileRepository(nil).BuildGroupPath(groupId, "")
	if err != nil {
		t.Fatalf("BuildGroupPath: %v", err)
	}
	entries, err := os.ReadDir(groupPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("read group dir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("group dir holds %d entries, want 0", len(entries))
	}
}

func assertValidatorKey(t *testing.T, w *httptest.ResponseRecorder, key string) map[string]string {
	t.Helper()
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	errs := map[string]string{}
	if err := json.Unmarshal(w.Body.Bytes(), &errs); err != nil {
		t.Fatalf("unmarshal body %q: %v", w.Body.String(), err)
	}
	if _, ok := errs[key]; !ok {
		t.Fatalf("expected a %q error, got %+v", key, errs)
	}
	return errs
}

func countRows(t *testing.T, model interface{}) int64 {
	t.Helper()
	var count int64
	if err := repositories.GetDB().Model(model).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return count
}

// ---------- POST /receipt/withFiles ----------

func TestCreateReceiptWithFiles_CreatesReceiptWithImagesAndComment(t *testing.T) {
	for _, asFilePart := range []bool{false, true} {
		t.Run(fmt.Sprintf("receiptAsFilePart=%v", asFilePart), func(t *testing.T) {
			defer repositories.TruncateTestDb()
			userId, groupId := seedRequirementGroup(t, "cwf-ok", true, true)
			jpg := readHandlerTestJpg(t)

			w := createWithFilesRequest(t, userId, receiptJson(groupId, userId, "Lunch"), asFilePart,
				uploadFile{"a.jpg", jpg}, uploadFile{"b.jpg", jpg})
			if w.Result().StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
			}

			var created models.Receipt
			if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if len(created.ImageFiles) != 2 {
				t.Errorf("imageFiles = %d, want 2", len(created.ImageFiles))
			}
			if len(created.Comments) != 1 || created.Comments[0].Comment != "Lunch" {
				t.Errorf("comments = %+v, want the submitted one", created.Comments)
			}
		})
	}
}

// Every file is checked before anything is written, so one bad file rejects the
// whole create.
func TestCreateReceiptWithFiles_InvalidFileTypeWritesNothing(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "cwf-bad-file", false, false)

	w := createWithFilesRequest(t, userId, receiptJson(groupId, userId, ""), false,
		uploadFile{"a.jpg", readHandlerTestJpg(t)}, uploadFile{"notes.txt", []byte("plain text, not an image")})
	assertValidatorKey(t, w, "files.1")
	assertNothingWritten(t, groupId)
}

func TestCreateReceiptWithFiles_EnforcesRoleRequirements(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "cwf-required", true, true)

	w := createWithFilesRequest(t, userId, receiptJson(groupId, userId, "   "), false)
	errs := assertValidatorKey(t, w, "comments")
	if _, ok := errs["files"]; !ok {
		t.Errorf("expected a files error too, got %+v", errs)
	}
	assertNothingWritten(t, groupId)
}

func TestCreateReceiptWithFiles_MissingReceiptPartIs400(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, _ := seedRequirementGroup(t, "cwf-no-receipt", false, false)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	CreateReceiptWithFiles(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Result().StatusCode)
	}
}

func TestCreateReceiptWithFiles_RequiresCreatePermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedGroupWithRole(t, "cwf-no-create", []string{permissions.GroupReceiptsRead})

	w := createWithFilesRequest(t, userId, receiptJson(groupId, userId, ""), false)
	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Result().StatusCode)
	}
}

// ---------- POST /receipt/withFiles: zero images ----------

// The common case: a role that requires nothing, a receipt with a comment and no
// images at all. Both clients now send every create through this endpoint.
func TestCreateReceiptWithFiles_NoFilesNoRequirements(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "cwf-no-files", false, false)

	w := createWithFilesRequest(t, userId, receiptJson(groupId, userId, "No image needed"), false)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}

	var created models.Receipt
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(created.ImageFiles) != 0 {
		t.Errorf("imageFiles = %d, want 0", len(created.ImageFiles))
	}
	if len(created.Comments) != 1 || created.Comments[0].Comment != "No image needed" {
		t.Errorf("comments = %+v, want the submitted one", created.Comments)
	}

	if n := countRows(t, &models.Receipt{}); n != 1 {
		t.Errorf("receipts = %d, want 1", n)
	}
	if n := countRows(t, &models.FileData{}); n != 0 {
		t.Errorf("file data = %d, want 0", n)
	}
	if n := countRows(t, &models.Comment{}); n != 1 {
		t.Errorf("comments = %d, want 1", n)
	}
}

// Requiring a comment must not also require an image.
func TestCreateReceiptWithFiles_CommentOnlyRequirementNeedsNoImage(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "cwf-comment-only", true, false)

	w := createWithFilesRequest(t, userId, receiptJson(groupId, userId, "Has a comment"), false)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	if n := countRows(t, &models.Receipt{}); n != 1 {
		t.Errorf("receipts = %d, want 1", n)
	}
}

// ---------- POST /receipt/withFiles: the shared 403 checks ----------
//
// These checks are shared with the deprecated JSON create and tested there too;
// these cases prove the new endpoint actually runs them. Each attaches a real
// image, so they also prove the checks run before anything is written.

func assertForbiddenAndNothingWritten(t *testing.T, w *httptest.ResponseRecorder, groupId uint) {
	t.Helper()
	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	assertNothingWritten(t, groupId)
}

func TestCreateReceiptWithFiles_DeniedForDisallowedCategory(t *testing.T) {
	defer repositories.TruncateTestDb()
	allowed := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&allowed)
	disallowed := models.Category{Name: "Salary"}
	repositories.GetDB().Create(&disallowed)
	userId, groupId := seedRestrictedReceiptCreator(t, []uint{allowed.ID})
	cleanUpGroupDir(t, groupId)

	body := receiptJsonWith(groupId, userId, fmt.Sprintf(`,"categories":[{"id":%d,"name":"Salary"}]`, disallowed.ID))
	w := createWithFilesRequest(t, userId, body, false, uploadFile{"a.jpg", readHandlerTestJpg(t)})
	assertForbiddenAndNothingWritten(t, w, groupId)
}

func TestCreateReceiptWithFiles_DeniedForNewCategoryWithoutCreatePermission(t *testing.T) {
	defer repositories.TruncateTestDb()
	allowed := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&allowed)
	// The member holds no app role, so no app.categories.create.
	userId, groupId := seedRestrictedReceiptCreator(t, []uint{allowed.ID})
	cleanUpGroupDir(t, groupId)

	body := receiptJsonWith(groupId, userId, `,"categories":[{"name":"Brand New"}]`)
	w := createWithFilesRequest(t, userId, body, false, uploadFile{"a.jpg", readHandlerTestJpg(t)})
	assertForbiddenAndNothingWritten(t, w, groupId)
}

func TestCreateReceiptWithFiles_DeniedForNonVisiblePaidBy(t *testing.T) {
	defer repositories.TruncateTestDb()
	fx := seedIsolatedReceiptGroupHandler(t, true)
	cleanUpGroupDir(t, fx.groupId)

	// Member A plants member B, whom A cannot see, as the payer.
	w := createWithFilesRequest(t, fx.memberAId, receiptJson(fx.groupId, fx.memberBId, ""), false,
		uploadFile{"a.jpg", readHandlerTestJpg(t)})
	assertForbiddenAndNothingWritten(t, w, fx.groupId)
}

func TestCreateReceiptWithFiles_DeniedForCustomFieldWithoutAccess(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)
	cleanUpGroupDir(t, groupId)
	field := seedCustomField(t)

	body := receiptJsonWith(groupId, userId, fmt.Sprintf(`,"customFields":[{"customFieldId":%d}]`, field.ID))
	w := createWithFilesRequest(t, userId, body, false, uploadFile{"a.jpg", readHandlerTestJpg(t)})
	assertForbiddenAndNothingWritten(t, w, groupId)
}

// Positive control for the 403 cases: the same restricted creator, choosing a
// category it may see, succeeds — so the denials above are about the selection,
// not about the fixture.
func TestCreateReceiptWithFiles_AllowedForGrantedCategory(t *testing.T) {
	defer repositories.TruncateTestDb()
	allowed := models.Category{Name: "Groceries"}
	repositories.GetDB().Create(&allowed)
	userId, groupId := seedRestrictedReceiptCreator(t, []uint{allowed.ID})
	cleanUpGroupDir(t, groupId)

	body := receiptJsonWith(groupId, userId, fmt.Sprintf(`,"categories":[{"id":%d,"name":"Groceries"}]`, allowed.ID))
	w := createWithFilesRequest(t, userId, body, false, uploadFile{"a.jpg", readHandlerTestJpg(t)})
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	if n := countRows(t, &models.FileData{}); n != 1 {
		t.Errorf("file data = %d, want 1", n)
	}
}

// ---------- POST /receipt/withFiles: invalid bodies ----------

func TestCreateReceiptWithFiles_InvalidBodiesAre400(t *testing.T) {
	multipartBody := func(t *testing.T, build func(*multipart.Writer)) (*bytes.Buffer, string) {
		t.Helper()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		build(writer)
		if err := writer.Close(); err != nil {
			t.Fatalf("close writer: %v", err)
		}
		return body, writer.FormDataContentType()
	}

	cases := []struct {
		name  string
		build func(t *testing.T, groupId uint, userId uint) (*bytes.Buffer, string)
	}{
		{"malformed receipt form value", func(t *testing.T, groupId uint, userId uint) (*bytes.Buffer, string) {
			return multipartBody(t, func(w *multipart.Writer) {
				w.WriteField("receipt", `{"name":"R",`)
			})
		}},
		{"malformed receipt json file part", func(t *testing.T, groupId uint, userId uint) (*bytes.Buffer, string) {
			return multipartBody(t, func(w *multipart.Writer) {
				header := make(textproto.MIMEHeader)
				header.Set("Content-Disposition", `form-data; name="receipt"; filename="receipt.json"`)
				header.Set("Content-Type", "application/json")
				part, _ := w.CreatePart(header)
				part.Write([]byte(`not json`))
			})
		}},
		{"empty receipt value and no file part", func(t *testing.T, groupId uint, userId uint) (*bytes.Buffer, string) {
			return multipartBody(t, func(w *multipart.Writer) {
				w.WriteField("receipt", "")
			})
		}},
		{"not multipart at all", func(t *testing.T, groupId uint, userId uint) (*bytes.Buffer, string) {
			return bytes.NewBufferString(receiptJson(groupId, userId, "")), "application/json"
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer repositories.TruncateTestDb()
			userId, groupId := seedRequirementGroup(t, "cwf-invalid", false, false)

			body, contentType := tc.build(t, groupId, userId)
			w := rawCreateWithFilesRequest(userId, body, contentType)
			if w.Result().StatusCode != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400 (body %s)", w.Result().StatusCode, w.Body.String())
			}
			assertNothingWritten(t, groupId)
		})
	}
}

// A well-formed receipt that fails validation is refused before any image is
// written, and the error names the field.
func TestCreateReceiptWithFiles_InvalidReceiptIs400(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "cwf-no-name", false, false)

	noName := fmt.Sprintf(
		`{"name":"","amount":"5","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`,
		groupId, userId,
	)
	w := createWithFilesRequest(t, userId, noName, false, uploadFile{"a.jpg", readHandlerTestJpg(t)})
	assertValidatorKey(t, w, "name")
	assertNothingWritten(t, groupId)
}

// ---------- the synthetic All group is never a destination ----------

// seedAllGroupMember puts userId in a real All-group row with perms — the exact
// condition under which the declarative group.receipts.create gate would pass.
func seedAllGroupMember(t *testing.T, userId uint, perms []string) uint {
	t.Helper()
	allGroup := models.Group{Name: "All", IsAllGroup: true}
	if err := repositories.GetDB().Create(&allGroup).Error; err != nil {
		t.Fatalf("seed all group: %v", err)
	}
	addMemberWithRole(t, userId, allGroup.ID, "all-group-member", perms)
	cleanUpGroupDir(t, allGroup.ID)
	return allGroup.ID
}

func TestCreateReceiptWithFiles_AllGroupRejected(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, realGroupId := seedRequirementGroup(t, "cwf-all-real", false, false)
	allGroupId := seedAllGroupMember(t, userId, requirementPerms)

	w := createWithFilesRequest(t, userId, receiptJson(allGroupId, userId, "into All"), false,
		uploadFile{"a.jpg", readHandlerTestJpg(t)})
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), allGroupCreateMessage) {
		t.Errorf("body = %s, want the All group message", w.Body.String())
	}
	assertNothingWritten(t, allGroupId)

	// Positive control: the same caller into a real group succeeds.
	w = createWithFilesRequest(t, userId, receiptJson(realGroupId, userId, "into a real group"), false,
		uploadFile{"a.jpg", readHandlerTestJpg(t)})
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("real group status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
}

func TestCreateReceipt_LegacyEndpointAllGroupRejected(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, _ := seedRequirementGroup(t, "legacy-all-real", false, false)
	allGroupId := seedAllGroupMember(t, userId, requirementPerms)

	w, r := createReceiptRequest(userId, receiptJson(allGroupId, userId, ""))
	CreateReceipt(w, r)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	if n := countRows(t, &models.Receipt{}); n != 0 {
		t.Errorf("receipts = %d, want 0", n)
	}
}

func TestQuickScanHandlerAllGroupRejected(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, realGroupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{})
	allGroupId := seedAllGroupMember(t, userId, []string{
		permissions.GroupReceiptsQuickScan, permissions.GroupReceiptsRead, permissions.GroupReceiptsCreate,
	})

	w := quickScanCommentRequest(t, userId, allGroupId, nil)
	assertValidatorKey(t, w, "files.0.groupId")

	// Positive control: the same scan into a real group gets past the check (it
	// fails later on the fake image, which is not what this asserts).
	w = quickScanCommentRequest(t, userId, realGroupId, nil)
	if w.Result().StatusCode == http.StatusBadRequest && strings.Contains(w.Body.String(), "files.0.groupId") {
		t.Errorf("a real group was rejected as the All group: %s", w.Body.String())
	}
}

// ---------- deprecated POST /receipt/ ----------

func TestCreateReceipt_LegacyEndpointEnforcesRoleRequirements(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "legacy-required", true, true)

	w, r := createReceiptRequest(userId, receiptJson(groupId, userId, "has a comment"))
	CreateReceipt(w, r)
	errs := assertValidatorKey(t, w, "files")
	if !strings.Contains(errs["files"], "update your app") {
		t.Errorf("files message = %q, want it to tell the user to update", errs["files"])
	}

	setMemberRequirements(t, userId, groupId, true, false)
	w, r = createReceiptRequest(userId, receiptJson(groupId, userId, ""))
	CreateReceipt(w, r)
	assertValidatorKey(t, w, "comments")

	w, r = createReceiptRequest(userId, receiptJson(groupId, userId, "now it has one"))
	CreateReceipt(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
}

func TestCreateReceipt_LegacyEndpointUnaffectedWithoutRequirements(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "legacy-free", false, false)

	w, r := createReceiptRequest(userId, receiptJson(groupId, userId, ""))
	CreateReceipt(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
}

// ---------- UpdateReceipt ----------

func addStoredImage(t *testing.T, receiptId uint, name string) models.FileData {
	t.Helper()
	fileData, err := repositories.NewReceiptImageRepository(nil).CreateReceiptImage(
		models.FileData{Name: name, ReceiptId: receiptId}, readHandlerTestJpg(t))
	if err != nil {
		t.Fatalf("CreateReceiptImage: %v", err)
	}
	return fileData
}

func addStoredComment(t *testing.T, receiptId uint, userId uint, text string) models.Comment {
	t.Helper()
	comment := models.Comment{Comment: text, ReceiptId: receiptId, UserId: &userId}
	if err := repositories.GetDB().Create(&comment).Error; err != nil {
		t.Fatalf("seed comment: %v", err)
	}
	return comment
}

// The update is judged on what the receipt stores, not on the command.
func TestUpdateReceipt_EnforcesRoleRequirementsOnStoredReceipt(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "update-required", true, true)
	receiptId := seedReceipt(t, groupId, userId, 0)

	w, r := updateReceiptRequest(receiptId, userId, receiptJson(groupId, userId, ""))
	UpdateReceipt(w, r)
	errs := assertValidatorKey(t, w, "comments")
	if _, ok := errs["files"]; !ok {
		t.Errorf("expected a files error too, got %+v", errs)
	}

	addStoredComment(t, receiptId, userId, "stored")
	addStoredImage(t, receiptId, "stored.jpg")

	w, r = updateReceiptRequest(receiptId, userId, receiptJson(groupId, userId, ""))
	UpdateReceipt(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
}

// A move is judged against the DESTINATION group's requirements.
func TestUpdateReceipt_MoveIntoRequiringGroupEnforced(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, sourceGroup := seedRequirementGroup(t, "move-req-src", false, false)
	_, destGroup := seedGroupWithRole(t, "move-req-dest", requirementPerms)
	addMemberWithRole(t, userId, destGroup, "move-req-dest-member", requirementPerms)
	setMemberRequirements(t, userId, destGroup, false, true)
	cleanUpGroupDir(t, destGroup)

	receiptId := seedReceipt(t, sourceGroup, userId, 0)

	w, r := updateReceiptRequest(receiptId, userId, moveBody(destGroup, userId))
	UpdateReceipt(w, r)
	assertValidatorKey(t, w, "files")

	var stored models.Receipt
	repositories.GetDB().First(&stored, receiptId)
	if stored.GroupId != sourceGroup {
		t.Errorf("receipt moved to %d despite the rejection", stored.GroupId)
	}

	// Staying in the unrestricted source group is fine.
	w, r = updateReceiptRequest(receiptId, userId, moveBody(sourceGroup, userId))
	UpdateReceipt(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
}

// ---------- DeleteComment / RemoveReceiptImage ----------

func TestDeleteComment_LastRequiredCommentRefused(t *testing.T) {
	defer repositories.TruncateTestDb()
	defer services.ClearRolePermissionCacheForTests()
	userId, groupId := seedRequirementGroup(t, "delete-comment", true, false)
	receiptId := seedReceipt(t, groupId, userId, 0)
	first := addStoredComment(t, receiptId, userId, "first")

	w, r := deleteCommentRequest(userId, utils.UintToString(first.ID))
	DeleteComment(w, r)
	assertValidatorKey(t, w, "comments")
	if n := countRows(t, &models.Comment{}); n != 1 {
		t.Errorf("comments = %d, want the only one kept", n)
	}

	addStoredComment(t, receiptId, userId, "second")
	w, r = deleteCommentRequest(userId, utils.UintToString(first.ID))
	DeleteComment(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
}

func removeImageRequest(userId uint, imageId uint) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/api", strings.NewReader(""))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", utils.UintToString(imageId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	return w, r
}

func TestRemoveReceiptImage_LastRequiredImageRefused(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRequirementGroup(t, "delete-image", false, true)
	receiptId := seedReceipt(t, groupId, userId, 0)
	first := addStoredImage(t, receiptId, "first.jpg")

	w, r := removeImageRequest(userId, first.ID)
	RemoveReceiptImage(w, r)
	assertValidatorKey(t, w, "files")
	if n := countRows(t, &models.FileData{}); n != 1 {
		t.Errorf("images = %d, want the only one kept", n)
	}

	addStoredImage(t, receiptId, "second.jpg")
	w, r = removeImageRequest(userId, first.ID)
	RemoveReceiptImage(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (body %s)", w.Result().StatusCode, w.Body.String())
	}
	if n := countRows(t, &models.FileData{}); n != 1 {
		t.Errorf("images = %d, want 1 after deleting one of two", n)
	}
}

// ---------- Quick scan ----------

// A role that requires a comment makes the quick-scan comment required even when
// the group's own quick-scan config leaves the field off.
func TestQuickScanHandlerRejectsMissingRoleRequiredComment(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedQuickScanCommenter(t, true, models.GroupReceiptSettings{})
	setMemberRequirements(t, userId, groupId, true, false)

	w := quickScanCommentRequest(t, userId, groupId, nil)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Result().StatusCode)
	}
}
