package models

// GroupRoleReportTemplateGrant restricts which report templates a group role's
// members may act on, per action. It mirrors GroupRoleCategoryGrant but folds the
// action into the primary key, so one row means "this group role may perform this
// one action on this one template". An empty set for a role is unrestricted (the
// role reaches every template its group access allows); a non-empty set restricts
// it to exactly the listed (template, action) pairs.
//
// Permission is one of the scopable report actions (read|generate|update|delete|
// duplicate) — not a foreign key, an app-validated token — so it carries an explicit
// size for the composite index (MariaDB/MySQL reject an unbounded VARCHAR there).
// The ReportTemplate association cascades so deleting a template cleans its grant
// rows; the GroupRoleID side cascades from the parent GroupRoleDefinition.
type GroupRoleReportTemplateGrant struct {
	GroupRoleID      uint           `gorm:"primaryKey;autoIncrement:false" json:"groupRoleId"`
	ReportTemplateID uint           `gorm:"primaryKey;autoIncrement:false;index" json:"reportTemplateId"`
	Permission       string         `gorm:"primaryKey;size:32" json:"permission"`
	ReportTemplate   ReportTemplate `gorm:"foreignKey:ReportTemplateID;constraint:OnDelete:CASCADE" json:"-"`
}
