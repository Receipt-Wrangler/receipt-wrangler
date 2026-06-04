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
	// AND) to run the handler. GroupPermissions are the group-scoped permissions
	// the caller must hold (logical AND) in each resolved group (resolved from
	// GroupId / GroupIds, or from ReceiptId / ReceiptIds). OrAppPermissions is an
	// app-scoped fallback: holding any of them bypasses the group-permission check
	// (the modern replacement for an admin override). All are resolved from the
	// database at request time by HandleRequest via the PermissionService.
	AppPermissions   []string
	GroupPermissions []string
	OrAppPermissions []string
}
