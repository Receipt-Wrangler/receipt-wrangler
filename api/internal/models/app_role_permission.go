package models

type AppRolePermission struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	AppRoleID  uint   `gorm:"not null;uniqueIndex:idx_app_role_perm,priority:1" json:"appRoleId"`
	Permission string `gorm:"not null;size:128;uniqueIndex:idx_app_role_perm,priority:2;index:idx_app_perm_lookup" json:"permission"`
}
