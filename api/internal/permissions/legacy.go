package permissions

import "slices"

// Legacy role permission sets.
//
// These reproduce the capabilities of the historical app roles (ADMIN/USER) and
// group roles (VIEWER/EDITOR/OWNER) as granular permission strings, so the
// seeded "Legacy *" system roles match the old enforcement exactly. They were
// derived from the actual handler-level role gating, not the desktop UI presets.

// scopeKeys returns every registry key for the given scope. Used so the
// "everything in this scope" legacy roles (admin / owner) stay complete as the
// permission registry grows.
func scopeKeys(scope Scope) []string {
	out := make([]string, 0, len(registry))
	for _, descriptor := range registry {
		if descriptor.Scope == scope {
			out = append(out, descriptor.Key)
		}
	}
	return out
}

// LegacyAppAdminKeys returns every app-scope permission. It mirrors the legacy
// ADMIN app role, which could perform every application-level action.
func LegacyAppAdminKeys() []string {
	return scopeKeys(ScopeApp)
}

// LegacyAppUserKeys returns the app-scope permissions a legacy USER app role
// held: actions whose handlers were gated to USER or left ungated (available to
// any authenticated user). Intentionally a fixed subset, so it must NOT grow
// automatically as new permissions are added.
//
// app.users.read is deliberately excluded: it gates only GET /user/ (the admin
// "Manage Users" listing), which no client calls — every user dropdown reads
// users from AppData (gated by app.account.read). Granting it would expose the
// admin Users page to normal users without giving them any capability they use.
//
// app.categories.read / app.tags.read are deliberately excluded as part of the
// category/tag grant lock-down: those gate the GLOBAL GET /category, GET /tag and
// the flat AppData category/tag lists, which return the entire pool. Normal users
// now receive only the per-group filtered catalogs (AppData groupCategories /
// groupTags), so granting the global read would leak categories/tags outside
// their grants. The create permissions are retained so inline category/tag
// creation in the receipt form still works (the create-permission gate).
func LegacyAppUserKeys() []string {
	return []string{
		AppCategoriesCreate,
		AppTagsCreate,
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
}

// LegacyGroupViewerKeys returns the group-scope permissions a legacy VIEWER
// group role held. This is the exact legacy match: a viewer could do more than
// read (comment, manage their own dashboards, magic-fill, poll email).
func LegacyGroupViewerKeys() []string {
	return []string{
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
}

// LegacyGroupEditorKeys returns the legacy EDITOR group role permissions: every
// viewer permission plus the editor-only receipt and activity actions.
// slices.Concat allocates a fresh slice, so the viewer set is never aliased.
func LegacyGroupEditorKeys() []string {
	return slices.Concat(LegacyGroupViewerKeys(), []string{
		GroupReceiptsCreate,
		GroupReceiptsUpdate,
		GroupReceiptsDelete,
		GroupReceiptsDuplicate,
		GroupReceiptsQuickScan,
		GroupActivitiesRerun,
	})
}

// LegacyGroupOwnerKeys returns every group-scope permission. It mirrors the
// legacy OWNER group role, the top of the group hierarchy, which could perform
// every group-level action.
func LegacyGroupOwnerKeys() []string {
	return scopeKeys(ScopeGroup)
}
