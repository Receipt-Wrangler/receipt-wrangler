package middleware

import (
	"net/http/httptest"
	"testing"
)

// The reimbursement claim routes, classified for API-key scopes. Kept beside the
// main inventory (everyRoute) rather than in it, so the fork's routes do not sit in
// a table upstream edits. Every claim read is a GET and every write a plain write,
// so none of them needs an allowlist entry.
var claimRoutes = []routeCase{
	{"GET", "/api/claim/group/1", classRead},
	{"POST", "/api/claim/group/1", classWrite},
	{"GET", "/api/claim/1", classRead},
	{"PUT", "/api/claim/1", classWrite},
	{"DELETE", "/api/claim/1", classWrite},
	{"POST", "/api/claim/1/addReceipts", classWrite},
	{"POST", "/api/claim/1/removeReceipts", classWrite},
	{"POST", "/api/claim/1/ledgerEntry", classWrite},
	{"PUT", "/api/claim/1/ledgerEntry/2", classWrite},
	{"DELETE", "/api/claim/1/ledgerEntry/2", classWrite},
}

func TestClassifyRequest_ClaimRoutes(t *testing.T) {
	for _, route := range claimRoutes {
		r := httptest.NewRequest(route.method, route.path, nil)
		if got := classifyRequest(r); got != route.class {
			t.Errorf("%s %s: got %s, want %s", route.method, route.path, className(got), className(route.class))
		}
	}
}
