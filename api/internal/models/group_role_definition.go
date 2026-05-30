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

	Permissions    []GroupRolePermission    `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
	CategoryGrants []GroupRoleCategoryGrant `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"categoryGrants"`
	TagGrants      []GroupRoleTagGrant      `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"tagGrants"`
}
