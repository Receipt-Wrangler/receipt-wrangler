package models

import "time"

type GroupRoleDefinition struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedBy       *uint     `json:"createdBy"`
	CreatedByString string    `json:"createdByString"`

	Name        string `gorm:"not null;size:64;uniqueIndex" json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `gorm:"not null;default:false" json:"isDefault"`
	IsSystem    bool   `gorm:"not null;default:false" json:"isSystem"`

	// SeesAllMembers is the "supervisor" exemption for member isolation (see
	// Group.IsolateMembers): holders of this group role can see every member of an
	// isolated group AND are visible to every member (so an isolated member can still
	// see and transact with their coordinator). Only meaningful in an isolated group;
	// default false ⇒ no effect on existing roles/groups.
	SeesAllMembers bool `gorm:"not null;default:false" json:"seesAllMembers"`

	// RequiresIndividualCategoryGrants makes per-member category assignment
	// MANDATORY for this role: a member holding it who has no individual category
	// grants sees NOTHING, rather than falling back to the role's set (or to
	// see-all when the role grants nothing). It exists so that forgetting to assign
	// a newly added member fails closed instead of exposing every category.
	// Defaults false, so every existing role behaves exactly as before.
	RequiresIndividualCategoryGrants bool `gorm:"not null;default:false" json:"requiresIndividualCategoryGrants"`

	// RequiresIndividualTagGrants is the tag counterpart of
	// RequiresIndividualCategoryGrants.
	RequiresIndividualTagGrants bool `gorm:"not null;default:false" json:"requiresIndividualTagGrants"`

	// IncludeOwnPaidReceipts is the relative "their own receipts" token of the
	// paid-by visibility filter: when true the role lets each member see receipts
	// they paid for. It is stored separately from PaidByUserGrants because it
	// resolves to the current member at query time rather than a fixed user id.
	IncludeOwnPaidReceipts bool `gorm:"not null;default:false" json:"includeOwnPaidReceipts"`

	// PaidByVisibilityRestricted records whether the admin opted into paid-by
	// filtering at all (any specific user grant OR include-own). It is what keeps a
	// configured role restricted even after its grant rows are removed — e.g. when a
	// granted user is deleted and the FK cascade empties PaidByUserGrants. Without
	// it, "no grants + include-own false" is indistinguishable from "never
	// configured", which would silently widen a restricted role to see-all. Derived
	// and set on save; internal only (not exposed on the API).
	PaidByVisibilityRestricted bool `gorm:"not null;default:false" json:"-"`

	// ReportTemplateGrantsRestricted plays the same fail-closed role for report
	// template grants as PaidByVisibilityRestricted does for paid-by: it records that
	// the admin opted into restricting this role to specific templates. It keeps the
	// role restricted even after its grant rows are emptied — e.g. when the last
	// granted template is deleted and the FK cascade clears ReportTemplateGrants —
	// instead of silently widening it back to see-all. Derived and set on save;
	// internal only (not exposed on the API).
	ReportTemplateGrantsRestricted bool `gorm:"not null;default:false" json:"-"`

	Permissions          []GroupRolePermission          `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"permissions"`
	CategoryGrants       []GroupRoleCategoryGrant       `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"categoryGrants"`
	TagGrants            []GroupRoleTagGrant            `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"tagGrants"`
	PaidByUserGrants     []GroupRolePaidByUserGrant     `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"paidByUserGrants"`
	ReportTemplateGrants []GroupRoleReportTemplateGrant `gorm:"foreignKey:GroupRoleID;constraint:OnDelete:CASCADE" json:"reportTemplateGrants"`
}
