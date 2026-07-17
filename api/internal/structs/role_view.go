package structs

import "receipt-wrangler/api/internal/permissions"

type RoleView struct {
	Id            uint              `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Scope         permissions.Scope `json:"scope"`
	IsDefault     bool              `json:"isDefault"`
	IsSystem      bool              `json:"isSystem"`
	Permissions   []string          `json:"permissions"`
	AssignedCount int               `json:"assignedCount"`
	// CategoryGrants / TagGrants are the category/tag ids a group role restricts
	// its members to. Empty means unrestricted (see all). Always group-scoped;
	// app roles serialize empty slices.
	CategoryGrants []uint `json:"categoryGrants"`
	TagGrants      []uint `json:"tagGrants"`
	// PaidByUserGrants are the user ids whose receipts a group role lets its
	// members see; IncludeOwnPaidReceipts adds the member's own receipts. Empty +
	// false means unrestricted. Always group-scoped; app roles serialize an empty
	// slice and false.
	PaidByUserGrants       []uint `json:"paidByUserGrants"`
	IncludeOwnPaidReceipts bool   `json:"includeOwnPaidReceipts"`
	// ReportTemplateGrants restrict which report templates a group role's members
	// may act on, per action. Empty means unrestricted (every template the role's
	// group access reaches). Always group-scoped; app roles serialize an empty
	// slice.
	ReportTemplateGrants []ReportTemplateGrantView `json:"reportTemplateGrants"`
}

// ReportTemplateGrantView is one row of a group role's report-template access
// matrix: a template id and the actions the role may perform on it.
type ReportTemplateGrantView struct {
	ReportTemplateId uint     `json:"reportTemplateId"`
	Permissions      []string `json:"permissions"`
}
