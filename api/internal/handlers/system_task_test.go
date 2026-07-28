package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"strings"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
)

// End-to-end: a hidden-peer activity in an isolated group must affect NEITHER the
// returned rows NOR TotalCount (the count is computed after member-isolation filtering,
// not before pagination), so a restricted member cannot infer hidden activity.
func TestGetActivitiesForGroupsFiltersCountAndRowsForIsolatedGroup(t *testing.T) {
	defer repositories.TruncateTestDb()
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	// The handler calls SetActivityCanBeRestarted, which needs the Redis connection
	// env vars to build an asynq inspector; point it at the local Redis (the inspector
	// tolerates an empty/unreachable queue and no-ops).
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")

	db := repositories.GetDB()
	roleRepository := repositories.NewRoleRepository(nil)

	group := models.Group{Name: "iso-activities", IsolateMembers: true}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("group: %v", err)
	}

	perms := []string{permissions.GroupActivitiesRead}
	supRole, err := roleRepository.CreateGroupRole("Iso Act Sup", "", perms, nil, nil, nil, false, true)
	if err != nil {
		t.Fatalf("sup role: %v", err)
	}
	memRole, err := roleRepository.CreateGroupRole("Iso Act Mem", "", perms, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("mem role: %v", err)
	}

	memberA := seedIsoHandlerUser(t, "iso-act-a")
	memberB := seedIsoHandlerUser(t, "iso-act-b")
	supervisor := seedIsoHandlerUser(t, "iso-act-sup")
	for _, m := range []models.GroupMember{
		{GroupID: group.ID, UserID: memberA, GroupRoleID: &memRole.ID},
		{GroupID: group.ID, UserID: memberB, GroupRoleID: &memRole.ID},
		{GroupID: group.ID, UserID: supervisor, GroupRoleID: &supRole.ID},
	} {
		if err := db.Create(&m).Error; err != nil {
			t.Fatalf("member: %v", err)
		}
	}

	// One activity per user in the group.
	for _, ranBy := range []uint{memberA, memberB, supervisor} {
		id := ranBy
		gid := group.ID
		task := models.SystemTask{
			Type:                 models.QUICK_SCAN,
			Status:               models.SYSTEM_TASK_SUCCEEDED,
			AssociatedEntityType: models.NOOP_ENTITY_TYPE,
			GroupId:              &gid,
			RanByUserId:          &id,
			StartedAt:            time.Now(),
		}
		if err := db.Create(&task).Error; err != nil {
			t.Fatalf("task: %v", err)
		}
	}

	// memberA is a plain isolated member: sees own + supervisor, never memberB.
	body, _ := json.Marshal(map[string]any{
		"groupIds":      []uint{group.ID},
		"page":          1,
		"pageSize":      10,
		"orderBy":       "started_at",
		"sortDirection": "desc",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(string(body)))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(memberA)))
	GetActivitiesForGroups(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Result().StatusCode, w.Body.String())
	}

	var paged struct {
		Data []struct {
			RanByUserId *uint `json:"ranByUserId"`
		} `json:"data"`
		TotalCount int64 `json:"totalCount"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &paged); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if paged.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2 (hidden peer excluded from the total, not just the page)", paged.TotalCount)
	}
	if len(paged.Data) != 2 {
		t.Errorf("returned %d rows, want 2", len(paged.Data))
	}
	for _, a := range paged.Data {
		if a.RanByUserId != nil && *a.RanByUserId == memberB {
			t.Errorf("hidden peer's activity leaked into the results")
		}
	}
}
