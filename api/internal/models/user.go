package models

import "time"

// User in the system
//
// swagger:model
type User struct {
	BaseModel
	DefaultAvatarColor string     `json:"defaultAvatarColor" gorm:"default:'#27b1ff'"`
	DisplayName        string     `json:"displayName"`
	IsDummyUser        bool       `json:"isDummyUser"`
	Password           string     `gorm:"not null"`
	Username           string     `gorm:"not null; uniqueIndex"`
	LastLoginDate      *time.Time `json:"lastLoginDate"`
	AppRoleID          *uint      `gorm:"index" json:"appRoleId"`
	AppRole            *AppRole   `gorm:"foreignKey:AppRoleID;constraint:OnDelete:RESTRICT" json:"appRole,omitempty"`
}
