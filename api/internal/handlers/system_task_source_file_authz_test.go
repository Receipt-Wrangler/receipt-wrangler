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
