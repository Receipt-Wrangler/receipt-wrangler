package models

import "time"

type GroupMember struct {
	CreatedAt           time.Time            `json:"createdAt"`
	UpdatedAt           time.Time            `json:"updatedAt"`
	UserID              uint                 `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	GroupID             uint                 `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	GroupRoleID         *uint                `gorm:"index" json:"groupRoleId"`
	GroupRoleDefinition *GroupRoleDefinition `gorm:"foreignKey:GroupRoleID;references:ID;constraint:OnDelete:RESTRICT" json:"groupRoleDefinition,omitempty"`
	// GroupRole is the obsolete legacy role enum (OWNER/EDITOR/VIEWER) removed in the role
	// rework, where GroupRoleID replaced it. It is temporarily re-declared here — nullable,
	// never read, hidden from the API via json:"-" — so that on databases upgraded from
	// before the rework GORM writes a value on INSERT and AutoMigrate relaxes the leftover
	// NOT NULL group_role column (AutoMigrate never dropped it). Without this, every new
	// group_members INSERT 500s on upgraded installs. Remove this field together with a
	// migration that drops the column once all installs have upgraded.
	GroupRole string `json:"-" gorm:"column:group_role"`

	// CategoryGrantsRestricted records whether an admin opted this membership into
	// per-member category filtering at all. It is what keeps a configured member
	// restricted after their grant rows are emptied — e.g. when the last granted
	// category is deleted and the FK cascade clears CategoryGrants. Without it,
	// "no grant rows" is indistinguishable from "never configured", which would
	// silently widen the member back to their role's full set. Derived and set on
	// save; internal only (not exposed on the API). Mirrors
	// GroupRoleDefinition.PaidByVisibilityRestricted.
	CategoryGrantsRestricted bool `gorm:"not null;default:false" json:"-"`

	// TagGrantsRestricted is the tag counterpart of CategoryGrantsRestricted.
	TagGrantsRestricted bool `gorm:"not null;default:false" json:"-"`

	// CategoryGrants/TagGrants are the member's granted category/tag ids, carried
	// for serialization only — they are deliberately NOT GORM associations
	// (`gorm:"-"`). GroupRepository.UpdateGroup replaces the whole member roster via
	// Association("GroupMembers").Unscoped().Replace, and a real association here
	// would put the grant join rows inside that wholesale write. Keeping them
	// transient means the grant tables are only ever touched explicitly, by
	// GroupMemberRepository. Populated on read by LoadMemberGrants.
	CategoryGrants []uint `gorm:"-" json:"categoryGrants"`
	TagGrants      []uint `gorm:"-" json:"tagGrants"`
}
