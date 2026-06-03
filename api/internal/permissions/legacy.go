package permissions

// Legacy role permission sets.
//
// These reproduce the capabilities of the legacy models.UserRole (ADMIN/USER)
// and models.GroupRole (VIEWER/EDITOR/OWNER) enums as granular permission
// strings, so the seeded "Legacy *" system roles match the old enforcement
// exactly. They were derived from the actual handler-level role gating, not the
// desktop UI presets.

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
// UserRole ADMIN, which could perform every application-level action.
func LegacyAppAdminKeys() []string {
	return scopeKeys(ScopeApp)
}

// LegacyAppUserKeys returns the app-scope permissions a legacy UserRole USER
// held: actions whose handlers were gated to USER or left ungated (available to
// any authenticated user). Intentionally a fixed subset, so it must NOT grow
// automatically as new permissions are added.
func LegacyAppUserKeys() []string {
	return []string{
		AppUsersRead,
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
	}
}

// LegacyGroupViewerKeys returns the group-scope permissions a legacy GroupRole
// VIEWER held. This is the exact legacy match: a viewer could do more than
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

// LegacyGroupEditorKeys returns the legacy GroupRole EDITOR permissions: every
// viewer permission plus the editor-only receipt and activity actions.
// LegacyGroupViewerKeys returns a fresh slice on each call, so the append
// cannot alias or mutate the viewer set.
func LegacyGroupEditorKeys() []string {
	return append(LegacyGroupViewerKeys(),
		GroupReceiptsCreate,
		GroupReceiptsUpdate,
		GroupReceiptsDelete,
		GroupReceiptsDuplicate,
		GroupReceiptsQuickScan,
		GroupActivitiesRerun,
	)
}

// LegacyGroupOwnerKeys returns every group-scope permission. It mirrors the
// legacy GroupRole OWNER, the top of the group hierarchy, which could perform
// every group-level action.
func LegacyGroupOwnerKeys() []string {
	return scopeKeys(ScopeGroup)
}
