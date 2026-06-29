package models

type GroupRoleCategoryGrant struct {
	GroupRoleID uint     `gorm:"primaryKey;autoIncrement:false" json:"groupRoleId"`
	CategoryID  uint     `gorm:"primaryKey;autoIncrement:false;index" json:"categoryId"`
	Category    Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"-"`
}
