package commands

type UpsertGroupMemberCommand struct {
	UserID      uint  `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	GroupID     uint  `gorm:"primaryKey;autoIncrement:false" json:"groupId"`
	GroupRoleID *uint `json:"groupRoleId"`
}
