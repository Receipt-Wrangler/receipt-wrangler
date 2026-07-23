package handlers

import (
	"receipt-wrangler/api/internal/structs"
	"testing"
)

func activityRunBy(id uint) structs.Activity {
	return structs.Activity{RanByUserId: &id}
}

func systemActivity() structs.Activity {
	return structs.Activity{RanByUserId: nil}
}

// An isolated (restricted) viewer sees only activities run by users in their
// visible set; a system activity (nil RanByUserId) is always kept.
func TestFilterActivitiesByVisibilityDropsInvisibleActor(t *testing.T) {
	visibleUserIds := map[uint]struct{}{1: {}, 3: {}} // self + a supervisor

	activities := []structs.Activity{
		activityRunBy(2), // invisible peer -> dropped
		activityRunBy(3), // visible supervisor -> kept
		systemActivity(), // system -> kept
	}

	filtered := filterActivitiesByVisibility(activities, visibleUserIds, false)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 visible activities, got %d (%v)", len(filtered), filtered)
	}
	for _, activity := range filtered {
		if activity.RanByUserId != nil && *activity.RanByUserId == 2 {
			t.Errorf("activity run by invisible user 2 should be dropped")
		}
	}
}

// An unrestricted viewer is unaffected — every activity is kept unchanged.
func TestFilterActivitiesByVisibilityUnrestrictedKeepsAll(t *testing.T) {
	activities := []structs.Activity{
		activityRunBy(2),
		activityRunBy(3),
		systemActivity(),
	}

	filtered := filterActivitiesByVisibility(activities, nil, true)

	if len(filtered) != len(activities) {
		t.Fatalf("unrestricted viewer should keep all %d activities, got %d", len(activities), len(filtered))
	}
}
