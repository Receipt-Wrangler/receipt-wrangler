package handlers

import (
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"testing"
)

func activityInGroup(ranBy uint, groupId uint) structs.Activity {
	return structs.Activity{RanByUserId: &ranBy, GroupId: &groupId}
}

// An isolated (restricted) viewer sees only activities run by users visible to them IN
// THAT ACTIVITY'S GROUP; a system activity (nil RanByUserId) is always kept.
func TestFilterActivitiesByVisibilityDropsInvisibleActor(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	permissionService := services.NewPermissionService(nil)

	groupId := fx.groupId
	activities := []structs.Activity{
		activityInGroup(fx.memberBId, fx.groupId),    // invisible peer -> dropped
		activityInGroup(fx.supervisorId, fx.groupId), // visible supervisor -> kept
		activityInGroup(fx.memberAId, fx.groupId),    // self -> kept
		{RanByUserId: nil, GroupId: &groupId},        // system action -> kept
	}

	filtered, err := filterActivitiesByVisibility(permissionService, fx.memberAId, activities)
	if err != nil {
		t.Fatalf("filter: %v", err)
	}
	if len(filtered) != 3 {
		t.Fatalf("expected 3 visible activities, got %d (%v)", len(filtered), filtered)
	}
	for _, activity := range filtered {
		if activity.RanByUserId != nil && *activity.RanByUserId == fx.memberBId {
			t.Errorf("activity run by invisible peer %d should be dropped", fx.memberBId)
		}
	}
}

// A viewer unrestricted in the activity's group (here the supervisor of an isolated
// group) keeps every activity unchanged.
func TestFilterActivitiesByVisibilityUnrestrictedKeepsAll(t *testing.T) {
	defer repositories.TruncateTestDb()

	fx := seedIsolatedReceiptGroupHandler(t, true)
	permissionService := services.NewPermissionService(nil)

	groupId := fx.groupId
	activities := []structs.Activity{
		activityInGroup(fx.memberAId, fx.groupId),
		activityInGroup(fx.memberBId, fx.groupId),
		{RanByUserId: nil, GroupId: &groupId},
	}

	filtered, err := filterActivitiesByVisibility(permissionService, fx.supervisorId, activities)
	if err != nil {
		t.Fatalf("filter: %v", err)
	}
	if len(filtered) != len(activities) {
		t.Fatalf("unrestricted viewer should keep all %d activities, got %d", len(activities), len(filtered))
	}
}
