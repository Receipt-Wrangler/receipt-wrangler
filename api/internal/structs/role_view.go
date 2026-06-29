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
}
