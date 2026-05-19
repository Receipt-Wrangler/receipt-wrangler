package models

type GroupRolePermission struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	GroupRoleID uint   `gorm:"not null;uniqueIndex:idx_group_role_perm,priority:1" json:"groupRoleId"`
	Permission  string `gorm:"not null;size:128;uniqueIndex:idx_group_role_perm,priority:2;index:idx_group_perm_lookup" json:"permission"`
}
