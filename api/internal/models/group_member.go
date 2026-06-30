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
}
