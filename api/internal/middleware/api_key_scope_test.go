package middleware

import (
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

// scopeRequest drives one request through the middleware with an API key.
func scopeRequest(t *testing.T, key string, method string, path string) int {
	t.Helper()

	r := httptest.NewRequest(method, path, nil)
	r.Header.Set("Authorization", key)
	w := httptest.NewRecorder()

	UnifiedAuthMiddleware(createFakeHandler()).ServeHTTP(w, r)

	return w.Result().StatusCode
}

// routeCase is one row of the inventory: a real route, and whether it reads.
type routeCase struct {
	method string
	path   string
	read   bool
}

// everyRoute is the whole authenticated surface, classified.
//
// This table is what makes the allowlist maintainable. A route added without
// thought about its scope shows up here as a missing row, and the two
// writes-dressed-as-reads at the bottom are the ones a reviewer would otherwise
// wave through.
var everyRoute = []routeCase{
	// Reads over GET -- the ordinary case.
	{http.MethodGet, "/api/receipt/1", true},
	{http.MethodGet, "/api/group/1", true},
	{http.MethodGet, "/api/category/", true},
	{http.MethodGet, "/api/search/", true},
	{http.MethodGet, "/api/report/template/1", true},
	{http.MethodGet, "/api/systemTask/1/sourceFile", true},

	// Reads over POST -- the allowlist.
	{http.MethodPost, "/api/apiKey/paged", true},
	{http.MethodPost, "/api/category/getPagedCategories", true},
	{http.MethodPost, "/api/customField/getPagedCustomFields", true},
	{http.MethodPost, "/api/group/getPagedGroups", true},
	{http.MethodPost, "/api/prompt/getPagedPrompts", true},
	{http.MethodPost, "/api/tag/getPagedTags", true},
	{http.MethodPost, "/api/user/getPagedUsers", true},
	{http.MethodPost, "/api/receiptProcessingSettings/getPagedProcessingSettings", true},
	{http.MethodPost, "/api/systemEmail/getSystemEmails", true},
	{http.MethodPost, "/api/systemTask/getPagedSystemTasks", true},
	{http.MethodPost, "/api/systemTask/getPagedActivities", true},
	{http.MethodPost, "/api/report/template/list", true},
	{http.MethodPost, "/api/report/generate", true},
	{http.MethodPost, "/api/report/preview", true},
	{http.MethodPost, "/api/report/receipts", true},
	{http.MethodPost, "/api/receiptImage/convertToJpg", true},
	{http.MethodPost, "/api/export/", true},
	{http.MethodPost, "/api/export/7", true},
	{http.MethodPost, "/api/receipt/group/7", true},
	{http.MethodPost, "/api/receipt/group/7/summary", true},
	{http.MethodPost, "/api/widget/pieChart/7", true},
	{http.MethodPost, "/api/report/template/7/generate", true},
	{http.MethodPost, "/api/report/template/7/render", true},

	// Writes.
	{http.MethodPost, "/api/receipt/", false},
	{http.MethodPut, "/api/receipt/1", false},
	{http.MethodDelete, "/api/receipt/1", false},
	{http.MethodPost, "/api/receipt/quickScan", false},
	{http.MethodPost, "/api/receipt/bulkStatusUpdate", false},
	{http.MethodPost, "/api/receipt/1/duplicate", false},
	{http.MethodPost, "/api/apiKey/", false},
	{http.MethodPut, "/api/apiKey/1", false},
	{http.MethodDelete, "/api/apiKey/1", false},
	{http.MethodPost, "/api/category/", false},
	{http.MethodPost, "/api/group/", false},
	{http.MethodPut, "/api/group/1", false},
	{http.MethodPost, "/api/report/template", false},
	{http.MethodPut, "/api/report/template/1", false},
	{http.MethodDelete, "/api/report/template/1", false},
	{http.MethodPost, "/api/report/template/1/duplicate", false},
	{http.MethodPost, "/api/user/1/resetPassword", false},
	{http.MethodDelete, "/api/user/bulk", false},
	{http.MethodPost, "/api/import/importConfigJson", false},
	{http.MethodPost, "/api/systemSettings/restartTaskServer", false},

	// Writes that LOOK like reads. Each persists a row or enqueues a task while
	// carrying a read-ish name, and two of them carry a ".read" permission.
	{http.MethodPost, "/api/systemEmail/checkConnectivity", false},
	{http.MethodPost, "/api/receiptProcessingSettings/checkConnectivity", false},
	{http.MethodPost, "/api/prompt/createDefaultPrompt", false},
	{http.MethodPost, "/api/group/1/pollGroupEmail", false},
	{http.MethodPost, "/api/systemTask/rerunActivity/1", false},
}

// TestApiKeyScopeIsEnforcedAcrossEveryRoute is the regression test for the
// finding: the scope was stored, validated and then never consulted, so a key
// created as "Read" could delete anything its owner could.
func TestApiKeyScopeIsEnforcedAcrossEveryRoute(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "test-key")
	defer teardownAuthTest()
	setupAuthTest()

	user := createTestUser()

	scopes := map[string]struct {
		key        string
		allowsRead bool
	}{}

	for _, scope := range []string{"r", "w", "rw"} {
		_, key, err := createTestApiKey(user.ID, scope)
		if err != nil {
			t.Fatalf("failed to create a %q key: %v", scope, err)
		}

		scopes[scope] = struct {
			key        string
			allowsRead bool
		}{key: key, allowsRead: scope != "w"}
	}

	for _, route := range everyRoute {
		for scope, held := range scopes {
			// Literal semantics: "r" reads only, "w" writes only, "rw" both.
			allowed := route.read == held.allowsRead || scope == "rw"

			want := http.StatusForbidden
			if allowed {
				want = http.StatusOK
			}

			got := scopeRequest(t, held.key, route.method, route.path)
			if got != want {
				t.Errorf("scope %q on %s %s: got %d, want %d", scope, route.method, route.path, got, want)
			}
		}
	}
}

// TestReadOnlyPostPathMatching pins the matcher itself, where a near-miss would
// either open a write to a read-only key or break a legitimate read.
func TestReadOnlyPostPathMatching(t *testing.T) {
	cases := []struct {
		path string
		read bool
	}{
		{"/api/receipt/group/7", true},
		{"/api/receipt/group/7/summary", true},
		{"/api/receipt/group/7/", true}, // trailing slash
		{"/api/receipt/group", false},   // one segment short
		{"/api/receipt/group/", false},  // empty parameter
		{"/api/receipt/group//", false}, // ditto
		{"/api/receipt/group/7/other", false},
		{"/api/receipt/1/duplicate", false}, // "group" is literal, not a parameter
		{"/api/export", true},
		{"/api/export/", true},
		{"/api/export/7", true},
		{"/api/export/7/extra", false},
		{"/api/report/template/list", true},
		{"/api/report/template", false},   // create
		{"/api/report/template/7", false}, // update is a PUT; a POST here is not a read
		{"/api/report/template/7/duplicate", false},
		{"/api/report/generate", true},
		{"/api/report/generatex", false}, // a prefix of an allowlisted path is not one
		{"/api/apiKey/paged", true},
		{"/api/apiKey/pagedx", false},
		{"", false},
		{"/", false},
	}

	for _, c := range cases {
		if got := isReadOnlyPostPath(c.path); got != c.read {
			t.Errorf("isReadOnlyPostPath(%q) = %v, want %v", c.path, got, c.read)
		}
	}
}

// TestApiKeyScopePermitsIsLiteral covers the rule directly, including the
// fail-closed behaviour for a scope that should be unreachable.
func TestApiKeyScopePermitsIsLiteral(t *testing.T) {
	cases := []struct {
		scope     string
		readOnly  bool
		permitted bool
	}{
		{"r", true, true},
		{"r", false, false},
		{"w", true, false},
		{"w", false, true},
		{"rw", true, true},
		{"rw", false, true},
		// Neither the command validator nor Claims.Validate lets these through, so
		// reaching here with one is a bug -- and a bug must not widen access.
		{"", true, false},
		{"", false, false},
		{"x", true, false},
		{"x", false, false},
	}

	for _, c := range cases {
		if got := apiKeyScopePermits(c.scope, c.readOnly); got != c.permitted {
			t.Errorf("apiKeyScopePermits(%q, readOnly=%v) = %v, want %v", c.scope, c.readOnly, got, c.permitted)
		}
	}
}

// TestApiKeyScopeStillStampsLastUsedOnARead: the last-used write is
// authentication bookkeeping, not the caller's operation, so a read-only key must
// still record that it was used. It is the one write a scope check must not gate,
// and it lives in the same function as the gate.
func TestApiKeyScopeStillStampsLastUsedOnARead(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "test-key")
	defer teardownAuthTest()
	setupAuthTest()

	user := createTestUser()

	dbApiKey, key, err := createTestApiKey(user.ID, "r")
	if err != nil {
		t.Fatalf("failed to create the key: %v", err)
	}

	if status := scopeRequest(t, key, http.MethodGet, "/api/receipt/1"); status != http.StatusOK {
		t.Fatalf("expected the read to be allowed, got %d", status)
	}

	// The stamp is written from a goroutine, so give it a moment to land.
	repository := repositories.NewApiKeyRepository(nil)
	for i := 0; i < 50; i++ {
		refreshed, err := repository.GetApiKeyById(dbApiKey.ID)
		if err == nil && refreshed.LastUsedAt != nil {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	utils.PrintTestError(t, "last_used_at never set", "a read-only key must still stamp its own last use")
}

// TestApiKeyScopeDoesNotAffectJwtCallers: the scope belongs to an API key. A
// session token carries none, and gating on an empty one would log every browser
// user out.
func TestApiKeyScopeDoesNotAffectJwtCallers(t *testing.T) {
	defer teardownAuthTest()
	setupAuthTest()

	user := createTestUser()

	jwt, _, _, err := services.GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("failed to mint a jwt: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/receipt/", nil)
	r.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()

	UnifiedAuthMiddleware(createFakeHandler()).ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("a jwt write must not be gated by API key scope, got %d", w.Result().StatusCode)
	}
}

// TestMountedRouterPreservesUrlPathForScopeMatching pins the assumption the whole
// allowlist rests on.
//
// The matcher reads r.URL.Path, and every router is reached through
// chi's Mount at an "/api/..." prefix. If Mount rewrote URL.Path to be relative
// to the mount point — as some routers do — no allowlist entry would ever match
// and every read-over-POST would 403 for a read-only key.
//
// The tests above drive the middleware directly with a path they choose
// themselves, so they cannot catch that. This one goes through a real chi mount.
//
// (chi tracks routing position in RouteContext.RoutePath and leaves URL.Path
// alone, which is also why the matcher does NOT use RoutePattern(): a mux runs
// its middleware chain before it routes, so the pattern is not populated yet.)
func TestMountedRouterPreservesUrlPathForScopeMatching(t *testing.T) {
	var seen string

	capture := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			seen = r.URL.Path
			next.ServeHTTP(w, r)
		})
	}

	receiptRouter := chi.NewRouter()
	receiptRouter.Use(capture)
	receiptRouter.Post("/group/{groupId}/summary", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	root := chi.NewRouter()
	root.Mount("/api/receipt", receiptRouter)

	request := httptest.NewRequest(http.MethodPost, "/api/receipt/group/7/summary", nil)
	root.ServeHTTP(httptest.NewRecorder(), request)

	if seen != "/api/receipt/group/7/summary" {
		t.Fatalf("middleware saw URL.Path %q, want the full mounted path", seen)
	}

	if !isReadOnlyPostPath(seen) {
		t.Errorf("the path a mounted middleware actually sees (%q) does not match the allowlist", seen)
	}
}
