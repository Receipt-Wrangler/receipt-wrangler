package models

// GroupRolePaidByUserGrant restricts a group role's members to receipts paid by
// specific users. It is the row-level "paid by" counterpart of
// GroupRoleCategoryGrant/GroupRoleTagGrant: a group role with NO paid-by grant
// rows and IncludeOwnPaidReceipts == false is unrestricted (sees every payer's
// receipts); a non-empty set restricts members to receipts paid by exactly those
// users. The relative "their own receipts" token is carried separately by
// GroupRoleDefinition.IncludeOwnPaidReceipts (it resolves to the current member
// at query time and so cannot be a stored user id).
type GroupRolePaidByUserGrant struct {
	GroupRoleID uint `gorm:"primaryKey;autoIncrement:false" json:"groupRoleId"`
	UserID      uint `gorm:"primaryKey;autoIncrement:false;index" json:"userId"`
	User        User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
