package models

type AppRole struct {
	BaseModel
	Name        string              `gorm:"not null;size:64;uniqueIndex" json:"name"`
	Description string              `json:"description"`
	IsSystem    bool                `gorm:"not null;default:false" json:"isSystem"`
	Permissions []AppRolePermission `gorm:"foreignKey:AppRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
}
