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
}
