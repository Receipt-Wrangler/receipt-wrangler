package models

import "time"

type GroupRoleDefinition struct {
	ID              uint      `gorm:"primarykey;uniqueIndex:idx_grdef_grp_id,priority:2" json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedBy       *uint     `json:"createdBy"`
	CreatedByString string    `json:"createdByString"`

	GroupID     uint   `gorm:"not null;uniqueIndex:idx_group_role_name,priority:1;index;uniqueIndex:idx_grdef_grp_id,priority:1" json:"groupId"`
	Group       Group  `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Name        string `gorm:"not null;size:64;uniqueIndex:idx_group_role_name,priority:2" json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `gorm:"not null;default:false" json:"isDefault"`
	IsSystem    bool   `gorm:"not null;default:false" json:"isSystem"`

	Permissions    []GroupRolePermission    `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
	CategoryGrants []GroupRoleCategoryGrant `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"categoryGrants"`
	TagGrants      []GroupRoleTagGrant      `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"tagGrants"`
}
