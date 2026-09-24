package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
)

// sourceFileAuthzFixture is one isolated group holding a QUICK_SCAN activity run by
// memberB, plus the three viewers whose answers must differ.
type sourceFileAuthzFixture struct {
	taskId     uint
	memberA    uint // plain member; must NOT see memberB
	memberB    uint // ran the activity
	supervisor uint // SeesAllMembers
	outsider   uint // not in the group at all
}

func seedSourceFileAuthzFixture(t *testing.T) sourceFileAuthzFixture {
	t.Helper()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()

	db := repositories.GetDB()
	roleRepository := repositories.NewRoleRepository(nil)

	group := models.Group{Name: "iso-source-file", IsolateMembers: true}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("group: %v", err)
	}

	perms := []string{permissions.GroupActivitiesRead, permissions.GroupActivitiesRerun}
	memberRole, err := roleRepository.CreateGroupRole("Iso SF Mem", "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("member role: %v", err)
	}
	supervisorRole, err := roleRepository.CreateGroupRole("Iso SF Sup", "", perms, nil, nil, nil, false, true)
	if err != nil {
		t.Fatalf("supervisor role: %v", err)
	}

	fixture := sourceFileAuthzFixture{
		memberA:    seedIsoHandlerUser(t, "iso-sf-a"),
		memberB:    seedIsoHandlerUser(t, "iso-sf-b"),
		supervisor: seedIsoHandlerUser(t, "iso-sf-sup"),
		outsider:   seedIsoHandlerUser(t, "iso-sf-out"),
	}

	for _, member := range []models.GroupMember{
		{GroupID: group.ID, UserID: fixture.memberA, GroupRoleID: &memberRole.ID},
		{GroupID: group.ID, UserID: fixture.memberB, GroupRoleID: &memberRole.ID},
		{GroupID: group.ID, UserID: fixture.supervisor, GroupRoleID: &supervisorRole.ID},
	} {
		if err := db.Create(&member).Error; err != nil {
			t.Fatalf("member: %v", err)
		}
	}

	groupId := group.ID
	ranBy := fixture.memberB
	task := models.SystemTask{
		Type:                 models.QUICK_SCAN,
		Status:               models.SYSTEM_TASK_FAILED,
		AssociatedEntityType: models.NOOP_ENTITY_TYPE,
		GroupId:              &groupId,
		RanByUserId:          &ranBy,
		AsynqTaskId:          "iso-sf-task",
		StartedAt:            time.Now(),
	}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("task: %v", err)
	}
	fixture.taskId = task.ID

	return fixture
}

func sourceFileRequest(t *testing.T, taskId uint, userId uint) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", nil)

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", utils.UintToString(taskId))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeContext))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))

	return w, r
}

// The headline case. A plain member of an isolated group holds
// group.activities.read, so the group gate alone admits them — but the activity
// list hides this row, and naming its id directly must not hand back the upload.
//
// Redis is never reached: the visibility check precedes the payload lookup, which
// is the same ordering that stops an unauthorized caller learning whether the file
// exists. That is why this test can exist at all.
func TestSourceFileEndpointsDenyAHiddenPeersActivity(t *testing.T) {
	defer repositories.TruncateTestDb()
	fixture := seedSourceFileAuthzFixture(t)

	for _, endpoint := range []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"preview", GetSystemTaskSourceFile},
		{"download", DownloadSystemTaskSourceFile},
		{"rerun", RerunActivity},
	} {
		t.Run(endpoint.name, func(t *testing.T) {
			w, r := sourceFileRequest(t, fixture.taskId, fixture.memberA)
			endpoint.handler(w, r)

			if w.Result().StatusCode != http.StatusForbidden {
				t.Errorf("status = %d, want 403; body=%s", w.Result().StatusCode, w.Body.String())
			}
		})
	}
}

// The contrasts, so the 403 above is attributable to isolation rather than to the
// endpoint failing for everyone. Neither of these reaches 403: they get past the
// gate and fail later, in Redis, which is not running here.
func TestSourceFileEndpointsAdmitVisibleActors(t *testing.T) {
	defer repositories.TruncateTestDb()
	fixture := seedSourceFileAuthzFixture(t)

	for _, viewer := range []struct {
		name   string
		userId uint
	}{
		{"the member who ran it", fixture.memberB},
		{"a supervisor", fixture.supervisor},
	} {
		t.Run(viewer.name, func(t *testing.T) {
			w, r := sourceFileRequest(t, fixture.taskId, viewer.userId)
			GetSystemTaskSourceFile(w, r)

			if w.Result().StatusCode == http.StatusForbidden {
				t.Errorf("status = 403 for %s, want the request to pass the gate", viewer.name)
			}
		})
	}
}

// A non-member has no group permission, so HandleRequest denies before the
// visibility check — the point being that the answer is the same 403 either way,
// so neither reveals which check refused.
func TestSourceFileEndpointsDenyANonMember(t *testing.T) {
	defer repositories.TruncateTestDb()
	fixture := seedSourceFileAuthzFixture(t)

	w, r := sourceFileRequest(t, fixture.taskId, fixture.outsider)
	GetSystemTaskSourceFile(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403; body=%s", w.Result().StatusCode, w.Body.String())
	}
}

// seedHiddenActivityOfType adds a second activity to the fixture's isolated group,
// run by the same hidden member, so a type-specific branch can be exercised under
// the same authorization conditions.
func seedHiddenActivityOfType(t *testing.T, fixture sourceFileAuthzFixture, taskType models.SystemTaskType) uint {
	t.Helper()

	existing := models.SystemTask{}
	if err := repositories.GetDB().First(&existing, fixture.taskId).Error; err != nil {
		t.Fatalf("fixture task: %v", err)
	}

	task := models.SystemTask{
		Type:                 taskType,
		Status:               models.SYSTEM_TASK_FAILED,
		AssociatedEntityType: models.NOOP_ENTITY_TYPE,
		GroupId:              existing.GroupId,
		RanByUserId:          existing.RanByUserId,
		AsynqTaskId:          "iso-sf-task-2",
		StartedAt:            time.Now(),
	}
	if err := repositories.GetDB().Create(&task).Error; err != nil {
		t.Fatalf("task: %v", err)
	}

	return task.ID
}

// An id that names no task at all must answer exactly like one the caller may not
// reach — for EVERY reason a caller can be refused. Asserting the status alone
// proves nothing, and neither does one pairing: the three denials are written by
// three different places (this file's task lookup, HandleRequest's group gate, and
// enforceActivityActorVisible), so all three have to agree on the body as well.
//
// Before this, GetSystemTaskById's gorm.ErrRecordNotFound surfaced as a 500 while a
// denied task returned 403, so the status code was an existence oracle for arbitrary
// task ids one layer above the gate that was meant to close it. Closing it with a
// message of this feature's own then left the same oracle in the BODY, since the
// group gate writes its own.
func TestSourceFileEndpointsAnswerAnUnknownTaskIdLikeADeniedOne(t *testing.T) {
	defer repositories.TruncateTestDb()
	fixture := seedSourceFileAuthzFixture(t)

	const unknownTaskId = 987654

	for _, endpoint := range []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"preview", GetSystemTaskSourceFile},
		{"download", DownloadSystemTaskSourceFile},
		{"rerun", RerunActivity},
	} {
		t.Run(endpoint.name, func(t *testing.T) {
			answer := func(taskId uint, userId uint) (int, string) {
				w, r := sourceFileRequest(t, taskId, userId)
				endpoint.handler(w, r)
				return w.Result().StatusCode, w.Body.String()
			}

			unknownStatus, unknownBody := answer(unknownTaskId, fixture.memberA)

			// Refused by enforceActivityActorVisible: a member of the group, but
			// the activity was run by a peer isolation hides from them.
			hiddenStatus, hiddenBody := answer(fixture.taskId, fixture.memberA)
			// Refused by HandleRequest's group gate, which writes its own body.
			outsiderStatus, outsiderBody := answer(fixture.taskId, fixture.outsider)

			for _, denied := range []struct {
				name   string
				status int
				body   string
			}{
				{"a hidden peer's activity", hiddenStatus, hiddenBody},
				{"an activity in a group they are not in", outsiderStatus, outsiderBody},
			} {
				if unknownStatus != denied.status {
					t.Errorf(
						"unknown id status = %d, %s status = %d; the two must be indistinguishable",
						unknownStatus, denied.name, denied.status,
					)
				}

				if unknownBody != denied.body {
					t.Errorf(
						"unknown id body = %s, %s body = %s; the two must be indistinguishable",
						unknownBody, denied.name, denied.body,
					)
				}

				// Pinned explicitly, so a regression that makes BOTH answers a 500
				// cannot pass the comparisons above.
				if denied.status != http.StatusForbidden {
					t.Errorf("%s status = %d, want 403; body=%s", denied.name, denied.status, denied.body)
				}
			}
		})
	}
}

// RerunActivity's "only a quick scan or email upload can be rerun" check used to run
// before the handler was built, so it answered ahead of the permission gate and told
// an unauthorized caller the activity's type. A caller who may not see the activity
// must get the same 403 as for any other task of theirs.
func TestRerunActivityHidesTheTaskTypeFromAnUnauthorizedCaller(t *testing.T) {
	defer repositories.TruncateTestDb()
	fixture := seedSourceFileAuthzFixture(t)
	nonRerunnableTaskId := seedHiddenActivityOfType(t, fixture, models.EMAIL_READ)

	w, r := sourceFileRequest(t, nonRerunnableTaskId, fixture.memberA)
	RerunActivity(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403; body=%s", w.Result().StatusCode, w.Body.String())
	}
}

// The contrast: a caller who MAY see the activity still gets the 400, so moving the
// check behind the gate did not quietly drop it. The type check precedes the asynq
// inspector, so this needs no Redis either.
func TestRerunActivityStillRejectsANonRerunnableTypeForAnAuthorizedCaller(t *testing.T) {
	defer repositories.TruncateTestDb()
	fixture := seedSourceFileAuthzFixture(t)
	nonRerunnableTaskId := seedHiddenActivityOfType(t, fixture, models.EMAIL_READ)

	w, r := sourceFileRequest(t, nonRerunnableTaskId, fixture.supervisor)
	RerunActivity(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400; body=%s", w.Result().StatusCode, w.Body.String())
	}
}
