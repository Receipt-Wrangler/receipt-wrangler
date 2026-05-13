package models

type GroupRoleTagGrant struct {
	GroupRoleID uint `gorm:"primaryKey;autoIncrement:false" json:"groupRoleId"`
	TagID       uint `gorm:"primaryKey;autoIncrement:false;index" json:"tagId"`
	Tag         Tag  `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE" json:"-"`
}
