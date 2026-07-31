package models

// GroupMemberTagGrant is the tag counterpart of GroupMemberCategoryGrant. Tags
// and categories are restricted independently — a membership may narrow one and
// leave the other unrestricted.
type GroupMemberTagGrant struct {
	UserID  uint `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	GroupID uint `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	TagID   uint `gorm:"primaryKey;autoIncrement:false;index" json:"tagId"`
	Tag     Tag  `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE" json:"-"`
}
