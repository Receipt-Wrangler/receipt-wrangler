package repositories

// Names of the five immutable, legacy-equivalent system roles seeded by
// SeedSystemRoles. They are the shared source of truth keyed on by both the
// seeder (which is idempotent on the role Name) and the data migration that
// assigns these roles to existing users / group members.
const (
	LegacyAdminRoleName  = "Legacy Admin"
	LegacyUserRoleName   = "Legacy User"
	LegacyViewerRoleName = "Legacy Viewer"
	LegacyEditorRoleName = "Legacy Editor"
	LegacyOwnerRoleName  = "Legacy Owner"
)
