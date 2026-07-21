package structs

import (
	"net/http"
)

type Handler struct {
	ErrorMessage    string
	Writer          http.ResponseWriter
	Request         *http.Request
	GroupId         string
	GroupIds        []string
	ReceiptId       string
	ReceiptIds      []string
	HandlerFunction func(http.ResponseWriter, *http.Request) (int, error)
	ResponseType    string

	// AppPermissions are the app-scoped permissions a caller must hold (logical
	// AND) to run the handler. AnyAppPermissions are app-scoped permissions of
	// which the caller must hold at least one (logical OR) — a *required* any-of
	// gate, distinct from AppPermissions (all-of / AND) and from OrAppPermissions
	// (a group-check bypass, below). GroupPermissions are the group-scoped
	// permissions the caller must hold (logical AND) in each resolved group
	// (resolved from GroupId / GroupIds, or from ReceiptId / ReceiptIds).
	// OrAppPermissions is an app-scoped fallback: holding any of them bypasses the
	// group-permission check (the modern replacement for an admin override). All
	// are resolved from the database at request time by HandleRequest via the
	// PermissionService.
	AppPermissions    []string
	AnyAppPermissions []string
	GroupPermissions  []string
	OrAppPermissions  []string

	// SkipPaidByVisibilityCheck, when true, tells HandleRequest to skip the
	// row-level "paid by" visibility check for receipts resolved from
	// ReceiptId/ReceiptIds, while STILL enforcing the group permissions above. Set
	// it on accounting/settlement endpoints (e.g. amount-owed) whose totals must be
	// identical for every member regardless of their browse-time paid-by filter —
	// paid-by restricts browsing, not the group's accounting.
	SkipPaidByVisibilityCheck bool
}
