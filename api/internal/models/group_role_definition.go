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

	Permissions      []GroupRolePermission      `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
	CategoryGrants   []GroupRoleCategoryGrant   `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"categoryGrants"`
	TagGrants        []GroupRoleTagGrant        `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"tagGrants"`
	PaidByUserGrants []GroupRolePaidByUserGrant `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"paidByUserGrants"`
}
