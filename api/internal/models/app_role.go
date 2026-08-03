package models

type AppRole struct {
	BaseModel
	Name        string `gorm:"not null;size:64;uniqueIndex" json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `gorm:"not null;default:false" json:"isDefault"`
	IsSystem    bool   `gorm:"not null;default:false" json:"isSystem"`
	// SkipDefaultGroupCreation suppresses the personal "My Receipts" group that is
	// otherwise created for every new user. The virtual "All" group is always
	// created, so the account still has a working dashboard. Creation-time only —
	// flipping it never adds or removes a group for an existing user. Default
	// false ⇒ existing roles keep creating the group (no data migration needed).
	SkipDefaultGroupCreation bool                `gorm:"not null;default:false" json:"skipDefaultGroupCreation"`
	Permissions              []AppRolePermission `gorm:"foreignKey:AppRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
}
