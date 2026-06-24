package models

import "time"

type GroupRoleDefinition struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedBy       *uint     `json:"createdBy"`
	CreatedByString string    `json:"createdByString"`

	Name        string `gorm:"not null;size:64;uniqueIndex" json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `gorm:"not null;default:false" json:"isDefault"`
	IsSystem    bool   `gorm:"not null;default:false" json:"isSystem"`

	// IncludeOwnPaidReceipts is the relative "their own receipts" token of the
	// paid-by visibility filter: when true the role lets each member see receipts
	// they paid for. It is stored separately from PaidByUserGrants because it
	// resolves to the current member at query time rather than a fixed user id.
	IncludeOwnPaidReceipts bool `gorm:"not null;default:false" json:"includeOwnPaidReceipts"`

	// PaidByVisibilityRestricted records whether the admin opted into paid-by
	// filtering at all (any specific user grant OR include-own). It is what keeps a
	// configured role restricted even after its grant rows are removed — e.g. when a
	// granted user is deleted and the FK cascade empties PaidByUserGrants. Without
	// it, "no grants + include-own false" is indistinguishable from "never
	// configured", which would silently widen a restricted role to see-all. Derived
	// and set on save; internal only (not exposed on the API).
	PaidByVisibilityRestricted bool `gorm:"not null;default:false" json:"-"`

	Permissions      []GroupRolePermission      `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
	CategoryGrants   []GroupRoleCategoryGrant   `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"categoryGrants"`
	TagGrants        []GroupRoleTagGrant        `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"tagGrants"`
	PaidByUserGrants []GroupRolePaidByUserGrant `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"paidByUserGrants"`
}
