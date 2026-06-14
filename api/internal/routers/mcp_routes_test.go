package routers

import (
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/constants"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// TestMountMcpRoutes verifies the MCP/OAuth endpoints are wired onto the chi
// router at the expected root paths and that the /mcp endpoint is guarded.
func TestMountMcpRoutes(t *testing.T) {
	t.Setenv(string(constants.McpPublicUrl), "https://receipts.example.com")

	router := chi.NewRouter()
	mountMcpRoutes(router)

	t.Run("protected resource metadata is served at the well-known root path", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/.well-known/oauth-protected-resource", nil))

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), "https://receipts.example.com/mcp") {
			t.Errorf("expected resource metadata to advertise the mcp resource, got %s", recorder.Body.String())
		}
	})

	t.Run("mcp endpoint challenges unauthenticated requests", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/mcp", nil))

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 from /mcp without a token, got %d", recorder.Code)
		}
		if !strings.Contains(recorder.Header().Get("WWW-Authenticate"), "resource_metadata") {
			t.Errorf("expected a WWW-Authenticate challenge pointing at resource metadata")
		}
	})
}
