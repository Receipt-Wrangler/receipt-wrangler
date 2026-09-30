package services

import (
	"encoding/json"
	"mime/multipart"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

type requirementFixture struct {
	userId  uint
	groupId uint
	roleId  uint
}

// seedRequirementMember creates a group (with receipt settings hiding comments /
// images as asked) and a member holding a fresh group role with perms and the
// given receipt-requirement flags.
func seedRequirementMember(
	t *testing.T,
	name string,
	perms []string,
	requireComment bool,
	requireImage bool,
	hideComments bool,
	hideImages bool,
) requirementFixture {
	t.Helper()
	clearRolePermissionCacheAll()
	db := repositories.GetDB()

	group := models.Group{Name: name}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	settings := models.GroupReceiptSettings{GroupId: group.ID, HideComments: hideComments, HideImages: hideImages}
	if err := db.Create(&settings).Error; err != nil {
		t.Fatalf("seed group receipt settings: %v", err)
	}

	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole(name+" role", "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("seed role: %v", err)
	}
	if err := roleRepository.SetGroupRoleReceiptRequirements(role.ID, requireComment, requireImage); err != nil {
		t.Fatalf("set requirements: %v", err)
	}

	user := models.User{Username: name + "-user", Password: "p", DisplayName: name}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	return requirementFixture{userId: user.ID, groupId: group.ID, roleId: role.ID}
}

func TestResolveReceiptRequirements_Matrix(t *testing.T) {
	commenter := []string{permissions.GroupReceiptsCreate, permissions.GroupCommentsCreate}
	cases := []struct {
		name           string
		perms          []string
		requireComment bool
		requireImage   bool
		hideComments   bool
		hideImages     bool
		want           structs.ReceiptRequirements
	}{
		{name: "no flags", perms: commenter, want: structs.ReceiptRequirements{}},
		{name: "both flags", perms: commenter, requireComment: true, requireImage: true,
			want: structs.ReceiptRequirements{CommentRequired: true, ImageRequired: true}},
		{name: "comments hidden", perms: commenter, requireComment: true, requireImage: true, hideComments: true,
			want: structs.ReceiptRequirements{ImageRequired: true}},
		{name: "images hidden", perms: commenter, requireComment: true, requireImage: true, hideImages: true,
			want: structs.ReceiptRequirements{CommentRequired: true}},
		{name: "cannot comment", perms: []string{permissions.GroupReceiptsCreate}, requireComment: true, requireImage: true,
			want: structs.ReceiptRequirements{ImageRequired: true}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer repositories.TruncateTestDb()
			fx := seedRequirementMember(t, "req-matrix", tc.perms, tc.requireComment, tc.requireImage, tc.hideComments, tc.hideImages)

			got, err := NewReceiptService(nil).ResolveReceiptRequirements(fx.userId, fx.groupId)
			if err != nil {
				t.Fatalf("ResolveReceiptRequirements: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}

			// The batched AppData form must agree with the single resolver.
			groups, err := NewGroupService(nil).GetGroupsForUser(utils.UintToString(fx.userId))
			if err != nil {
				t.Fatalf("GetGroupsForUser: %v", err)
			}
			groupPermissions := map[uint][]string{fx.groupId: tc.perms}
			batched, err := NewReceiptService(nil).ResolveReceiptRequirementsForGroups(fx.userId, groups, groupPermissions)
			if err != nil {
				t.Fatalf("ResolveReceiptRequirementsForGroups: %v", err)
			}
			if batched[fx.groupId] != tc.want {
				t.Errorf("batched got %+v, want %+v", batched[fx.groupId], tc.want)
			}
			if _, present := batched[fx.groupId]; present != tc.want.Any() {
				t.Errorf("batched presence = %v, want %v", present, tc.want.Any())
			}
		})
	}
}

// A non-member, a member without a role, and the All group require nothing, even
// when a flagged role exists.
func TestResolveReceiptRequirements_NothingRequired(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()
	fx := seedRequirementMember(t, "req-none", []string{permissions.GroupCommentsCreate}, true, true, false, false)

	outsider := models.User{Username: "req-outsider", Password: "p"}
	if err := db.Create(&outsider).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	noRoleGroup := models.Group{Name: "req-no-role"}
	allGroup := models.Group{Name: "req-all", IsAllGroup: true}
	for _, group := range []*models.Group{&noRoleGroup, &allGroup} {
		if err := db.Create(group).Error; err != nil {
			t.Fatalf("seed group: %v", err)
		}
	}
	if err := db.Create(&models.GroupMember{GroupID: noRoleGroup.ID, UserID: fx.userId}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: allGroup.ID, UserID: fx.userId, GroupRoleID: &fx.roleId}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	service := NewReceiptService(nil)
	for name, target := range map[string][2]uint{
		"non-member": {outsider.ID, fx.groupId},
		"no role":    {fx.userId, noRoleGroup.ID},
		"all group":  {fx.userId, allGroup.ID},
	} {
		got, err := service.ResolveReceiptRequirements(target[0], target[1])
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got.Any() {
			t.Errorf("%s: got %+v, want nothing required", name, got)
		}
	}
}

// AppData always carries the map, serialized as {} when nothing is required —
// never null, which the generated Dart client could not parse.
func TestGetAppData_GroupReceiptRequirements(t *testing.T) {
	defer repositories.TruncateTestDb()

	user, err := repositories.NewUserRepository(nil).CreateUser(commands.SignUpCommand{
		Username:    "req-appdata",
		Password:    "Password",
		DisplayName: "Req AppData",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	appData, err := GetAppData(user.ID, nil)
	if err != nil {
		t.Fatalf("GetAppData: %v", err)
	}
	payload, err := json.Marshal(appData)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(payload), `"groupReceiptRequirements":{}`) {
		t.Errorf("expected an empty groupReceiptRequirements object in %s", payload)
	}

	fx := seedRequirementMember(t, "req-appdata-group", []string{permissions.GroupCommentsCreate}, true, false, false, false)
	appData, err = GetAppData(fx.userId, nil)
	if err != nil {
		t.Fatalf("GetAppData: %v", err)
	}
	want := structs.ReceiptRequirements{CommentRequired: true}
	if appData.GroupReceiptRequirements[fx.groupId] != want {
		t.Errorf("groupReceiptRequirements[%d] = %+v, want %+v", fx.groupId, appData.GroupReceiptRequirements[fx.groupId], want)
	}
}

// ---------- CreateReceiptWithFiles ----------

func readTestJpg(t *testing.T) []byte {
	t.Helper()
	jpg, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return jpg
}

// seedCreateWithFilesGroup creates a user and a group they belong to, and removes
// the group's data directory afterwards. It returns the user, the group and the
// group's directory.
func seedCreateWithFilesGroup(t *testing.T) (models.User, models.Group, string) {
	t.Helper()
	db := repositories.GetDB()

	user := models.User{Username: "cwf-user", Password: "p", DisplayName: "Cwf"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	group := models.Group{Name: "cwf-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: user.ID}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	groupPath, err := repositories.NewFileRepository(nil).BuildGroupPath(group.ID, "")
	if err != nil {
		t.Fatalf("BuildGroupPath: %v", err)
	}
	t.Cleanup(func() { utils.RemoveAllInDataDir(groupPath) })

	return user, group, groupPath
}

func createWithFilesCommand(groupId uint, userId uint, comment string) commands.UpsertReceiptCommand {
	command := commands.UpsertReceiptCommand{
		Name:         "With files",
		Amount:       decimal.NewFromInt(12),
		Date:         time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		GroupId:      groupId,
		PaidByUserID: userId,
		Status:       models.OPEN,
	}
	if comment != "" {
		command.Comments = []commands.UpsertCommentCommand{{Comment: comment, UserId: &userId}}
	}
	return command
}

func TestCreateReceiptWithFiles_CreatesReceiptImagesAndComment(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, _ := seedCreateWithFilesGroup(t)
	jpg := readTestJpg(t)

	created, err := NewReceiptService(nil).CreateReceiptWithFiles(
		createWithFilesCommand(group.ID, user.ID, "Team lunch"),
		[]commands.ReceiptFileUpload{{Name: "front.jpg", Bytes: jpg}, {Name: "back.jpg", Bytes: jpg}},
		user.ID,
	)
	if err != nil {
		t.Fatalf("CreateReceiptWithFiles: %v", err)
	}

	if len(created.ImageFiles) != 2 {
		t.Fatalf("expected 2 images on the created receipt, got %d", len(created.ImageFiles))
	}
	for _, image := range created.ImageFiles {
		if image.FileType != "image/jpeg" {
			t.Errorf("image %s file type = %q, want image/jpeg", image.Name, image.FileType)
		}
		onDisk, err := os.ReadFile(imagePath(t, image))
		if err != nil {
			t.Errorf("image %s not on disk: %v", image.Name, err)
			continue
		}
		if len(onDisk) != len(jpg) {
			t.Errorf("image %s is %d bytes on disk, want %d", image.Name, len(onDisk), len(jpg))
		}
	}

	if len(created.Comments) != 1 || created.Comments[0].Comment != "Team lunch" {
		t.Errorf("comments = %+v, want the one submitted", created.Comments)
	}

	tasks := receiptUploadedTasks(t)
	if len(tasks) != 1 || tasks[0].Status != models.SYSTEM_TASK_SUCCEEDED {
		t.Fatalf("expected one SUCCEEDED RECEIPT_UPLOADED task, got %+v", tasks)
	}
	if tasks[0].ReceiptId == nil || *tasks[0].ReceiptId != created.ID {
		t.Errorf("task receipt id = %v, want %d", tasks[0].ReceiptId, created.ID)
	}
}

// A failure after the first image was written leaves nothing behind: no receipt,
// no FileData, and not the image file already on disk. The failure is recorded.
func TestCreateReceiptWithFiles_FailureLeavesNothingBehind(t *testing.T) {
	defer repositories.TruncateTestDb()
	user, group, groupPath := seedCreateWithFilesGroup(t)
	jpg := readTestJpg(t)

	// The second name escapes the data directory, so building its path fails
	// after the first image has been written.
	_, err := NewReceiptService(nil).CreateReceiptWithFiles(
		createWithFilesCommand(group.ID, user.ID, "Should not survive"),
		[]commands.ReceiptFileUpload{{Name: "first.jpg", Bytes: jpg}, {Name: "../../../../escape.jpg", Bytes: jpg}},
		user.ID,
	)
	if err == nil {
		t.Fatal("expected CreateReceiptWithFiles to fail")
	}

	db := repositories.GetDB()
	var receiptCount, fileDataCount, commentCount int64
	db.Model(&models.Receipt{}).Count(&receiptCount)
	db.Model(&models.FileData{}).Count(&fileDataCount)
	db.Model(&models.Comment{}).Count(&commentCount)
	if receiptCount != 0 || fileDataCount != 0 || commentCount != 0 {
		t.Errorf("expected nothing persisted, got %d receipts, %d images, %d comments", receiptCount, fileDataCount, commentCount)
	}

	entries, readErr := os.ReadDir(groupPath)
	if readErr != nil && !os.IsNotExist(readErr) {
		t.Fatalf("read group dir: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("expected no files left in the group directory, got %d", len(entries))
	}

	tasks := receiptUploadedTasks(t)
	if len(tasks) != 1 || tasks[0].Status != models.SYSTEM_TASK_FAILED {
		t.Fatalf("expected one FAILED RECEIPT_UPLOADED task, got %+v", tasks)
	}
	if tasks[0].ReceiptId != nil {
		t.Errorf("task receipt id = %d, want none for a rolled-back create", *tasks[0].ReceiptId)
	}
}

// A role requiring a comment shows and requires the quick-scan comment field even
// when the group's own quick-scan config leaves it off.
func TestResolveQuickScanFields_RoleRequiredComment(t *testing.T) {
	defer repositories.TruncateTestDb()
	fx := seedRequirementMember(t, "req-quick-scan",
		[]string{permissions.GroupReceiptsQuickScan, permissions.GroupCommentsCreate}, true, false, false, false)

	command := func(comment string) commands.QuickScanCommand {
		return commands.QuickScanCommand{
			Files:         make([]multipart.File, 1),
			GroupIds:      []uint{fx.groupId},
			PaidByUserIds: []uint{fx.userId},
			Statuses:      []models.ReceiptStatus{models.OPEN},
			Comments:      []string{comment},
		}
	}

	_, configErr, err := NewReceiptService(nil).ResolveQuickScanFields(command(""), fx.userId)
	if err != nil {
		t.Fatalf("ResolveQuickScanFields: %v", err)
	}
	if _, ok := configErr.Errors["files.0.comment"]; !ok {
		t.Errorf("expected files.0.comment error, got %+v", configErr.Errors)
	}

	resolved, configErr, err := NewReceiptService(nil).ResolveQuickScanFields(command("Kept"), fx.userId)
	if err != nil {
		t.Fatalf("ResolveQuickScanFields: %v", err)
	}
	if len(configErr.Errors) > 0 {
		t.Errorf("expected no errors, got %+v", configErr.Errors)
	}
	if resolved[0].Comment != "Kept" {
		t.Errorf("comment = %q, want it kept rather than dropped", resolved[0].Comment)
	}
}
