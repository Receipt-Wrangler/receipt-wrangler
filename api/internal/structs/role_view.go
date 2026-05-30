package structs

import "receipt-wrangler/api/internal/permissions"

type RoleView struct {
	Id          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Scope       permissions.Scope `json:"scope"`
	IsSystem    bool              `json:"isSystem"`
	Permissions []string          `json:"permissions"`
}
