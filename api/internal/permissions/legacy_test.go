package permissions

import (
	"slices"
	"sort"
	"testing"
)

// sortedCopy returns a sorted copy of keys, leaving the input untouched.
func sortedCopy(keys []string) []string {
	out := append([]string(nil), keys...)
	sort.Strings(out)
	return out
}

// equalSet reports whether two key slices contain the same keys (order-independent).
func equalSet(a []string, b []string) bool {
	return slices.Equal(sortedCopy(a), sortedCopy(b))
}

// countScope returns how many registry descriptors belong to the given scope.
func countScope(scope Scope) int {
	count := 0
	for _, descriptor := range registry {
		if descriptor.Scope == scope {
			count++
		}
	}
	return count
}

func TestLegacyAppUserKeys(t *testing.T) {
	expected := []string{
		AppCategoriesCreate,
		AppCategoriesRead,
		AppTagsCreate,
		AppTagsRead,
		AppCustomFieldsCreate,
		AppCustomFieldsRead,
		AppGroupsCreate,
		AppApiKeysCreate,
		AppApiKeysRead,
		AppApiKeysUpdate,
		AppApiKeysDelete,
		AppNotificationsRead,
		AppNotificationsDelete,
		AppUserPreferencesRead,
		AppUserPreferencesUpdate,
		AppAccountRead,
		AppAccountUpdate,
		AppAccountDelete,
		AppReceiptsSearch,
	}

	got := LegacyAppUserKeys()
	if !equalSet(got, expected) {
		t.Helper()
		utilPrint(t, got, expected)
	}

	// Every key must be app-scoped.
	for _, key := range got {
		descriptor, ok := Get(key)
		if !ok || descriptor.Scope != ScopeApp {
			utilPrint(t, key, "an app-scope permission")
		}
	}
}

func TestLegacyAppAdminIncludesReadAnyApiKeys(t *testing.T) {
	// Legacy Admin is every app permission, so a newly added app permission such
	// as app.api-keys.read-any must flow into it automatically (upgrade keeps the
	// admin's "view all API keys" ability).
	if !slices.Contains(LegacyAppAdminKeys(), AppApiKeysReadAny) {
		utilPrint(t, "Legacy Admin missing "+AppApiKeysReadAny, "present")
	}
}

func TestLegacyAppUserExcludesReadAnyApiKeys(t *testing.T) {
	// Legacy User is a fixed subset and must NOT grant the privileged
	// "view all API keys" permission.
	if slices.Contains(LegacyAppUserKeys(), AppApiKeysReadAny) {
		utilPrint(t, "Legacy User contains "+AppApiKeysReadAny, "absent")
	}
}

func TestLegacyAppUserExcludesUsersRead(t *testing.T) {
	// app.users.read gates only the admin "Manage Users" listing (GET /user/),
	// which no client calls. Legacy User must NOT hold it, so normal users don't
	// reach the admin Users page.
	if slices.Contains(LegacyAppUserKeys(), AppUsersRead) {
		utilPrint(t, "Legacy User contains "+AppUsersRead, "absent")
	}
}

func TestLegacyGroupViewerKeys(t *testing.T) {
	expected := []string{
		GroupView,
		GroupEmailPoll,
		GroupReceiptsRead,
		GroupReceiptsMagicFill,
		GroupCommentsCreate,
		GroupCommentsDelete,
		GroupDashboardsCreate,
		GroupDashboardsRead,
		GroupDashboardsUpdate,
		GroupDashboardsDelete,
		GroupWidgetsRead,
		GroupActivitiesRead,
	}

	got := LegacyGroupViewerKeys()
	if !equalSet(got, expected) {
		utilPrint(t, got, expected)
	}

	for _, key := range got {
		descriptor, ok := Get(key)
		if !ok || descriptor.Scope != ScopeGroup {
			utilPrint(t, key, "a group-scope permission")
		}
	}
}

func TestLegacyGroupEditorKeysSupersetOfViewer(t *testing.T) {
	viewer := LegacyGroupViewerKeys()
	editor := LegacyGroupEditorKeys()

	// Editor is a strict superset of viewer.
	for _, key := range viewer {
		if !slices.Contains(editor, key) {
			utilPrint(t, "editor missing viewer key "+key, "present")
		}
	}

	editorOnly := []string{
		GroupReceiptsCreate,
		GroupReceiptsUpdate,
		GroupReceiptsDelete,
		GroupReceiptsDuplicate,
		GroupReceiptsQuickScan,
		GroupActivitiesRerun,
	}
	for _, key := range editorOnly {
		if !slices.Contains(editor, key) {
			utilPrint(t, "editor missing key "+key, "present")
		}
	}

	if len(editor) != len(viewer)+len(editorOnly) {
		utilPrint(t, len(editor), len(viewer)+len(editorOnly))
	}
}

func TestLegacyAppAdminKeysAreAllAppScope(t *testing.T) {
	got := LegacyAppAdminKeys()

	if len(got) != countScope(ScopeApp) {
		utilPrint(t, len(got), countScope(ScopeApp))
	}

	for _, key := range got {
		descriptor, ok := Get(key)
		if !ok || descriptor.Scope != ScopeApp {
			utilPrint(t, key, "an app-scope permission")
		}
	}
}

func TestLegacyGroupOwnerKeysAreAllGroupScope(t *testing.T) {
	got := LegacyGroupOwnerKeys()

	if len(got) != countScope(ScopeGroup) {
		utilPrint(t, len(got), countScope(ScopeGroup))
	}

	for _, key := range got {
		descriptor, ok := Get(key)
		if !ok || descriptor.Scope != ScopeGroup {
			utilPrint(t, key, "a group-scope permission")
		}
	}

	// Owner is a superset of editor.
	for _, key := range LegacyGroupEditorKeys() {
		if !slices.Contains(got, key) {
			utilPrint(t, "owner missing editor key "+key, "present")
		}
	}
}

func TestLegacyKeysAllExistInRegistry(t *testing.T) {
	groups := [][]string{
		LegacyAppAdminKeys(),
		LegacyAppUserKeys(),
		LegacyGroupViewerKeys(),
		LegacyGroupEditorKeys(),
		LegacyGroupOwnerKeys(),
	}

	for _, keys := range groups {
		for _, key := range keys {
			if !Exists(key) {
				utilPrint(t, key, "a key that exists in the registry")
			}
		}
	}
}

func TestLegacyGroupEditorDoesNotMutateViewer(t *testing.T) {
	// Build the editor set (which appends onto a viewer slice internally), then
	// confirm the viewer helper is unaffected.
	_ = LegacyGroupEditorKeys()

	viewer := LegacyGroupViewerKeys()
	if len(viewer) != 12 {
		utilPrint(t, len(viewer), 12)
	}
}

// utilPrint mirrors utils.PrintTestError without importing utils (which would
// create an import cycle for the permissions package's test).
func utilPrint(t *testing.T, actual any, expected any) {
	t.Helper()
	t.Errorf("Test failed!\nExpected: %v\nActual: %v\n", expected, actual)
}
