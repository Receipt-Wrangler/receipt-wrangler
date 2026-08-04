package models

// GroupMemberCategoryGrant restricts a single group MEMBER to specific
// categories, independent of the group role they hold. It is the per-individual
// counterpart of GroupRoleCategoryGrant: a role grant applies to everyone
// holding the role, which cannot express "Alice sees Child A and B, Bob sees
// Child C".
//
// The two layers compose by INTERSECTION — the role grant is a ceiling and the
// membership grant narrows within it (see PermissionService.resolveEffectiveGrants).
// A membership with no grant rows and CategoryGrantsRestricted == false does not
// narrow anything, so existing installs are unchanged.
//
// Scoped to (user, group) rather than to the user alone because categories are
// global but visibility is always resolved in a group context: a global per-user
// list would follow the member into unrelated groups and blank their catalog there.
type GroupMemberCategoryGrant struct {
	UserID     uint     `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	GroupID    uint     `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	CategoryID uint     `gorm:"primaryKey;autoIncrement:false;index" json:"categoryId"`
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"-"`
}
