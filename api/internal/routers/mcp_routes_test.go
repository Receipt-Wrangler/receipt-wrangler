package routers

import (
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

const testMcpPublicUrl = "https://receipts.example.com"

// enableMcp turns the MCP server on or off in System Settings, which is the
// live gate the always-mounted routes consult (the env-based config was
// removed).
func enableMcp(t *testing.T, enabled bool) {
	t.Helper()

	settings, err := repositories.NewSystemSettingsRepository(nil).GetSystemSettings()
	if err != nil {
		t.Fatalf("failed to load system settings: %v", err)
	}

	err = repositories.GetDB().
		Model(&models.SystemSettings{}).
		Where("id = ?", settings.ID).
		Updates(map[string]interface{}{"mcp_enabled": enabled, "mcp_public_url": testMcpPublicUrl}).Error
	if err != nil {
		t.Fatalf("failed to set mcp system settings: %v", err)
	}
}

// TestMountMcpRoutes verifies the MCP/OAuth endpoints are wired onto the chi
// router at the expected root paths and that the /mcp endpoint is guarded.
func TestMountMcpRoutes(t *testing.T) {
	defer repositories.TruncateTestDb()
	enableMcp(t, true)

	router := chi.NewRouter()
	mountMcpRoutes(router)

	t.Run("protected resource metadata is served at the well-known root path", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/.well-known/oauth-protected-resource", nil))

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), testMcpPublicUrl+"/mcp") {
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

// TestMountMcpRoutesDisabled verifies that when the MCP server is disabled the
// always-mounted routes behave as if absent (404), without a restart.
func TestMountMcpRoutesDisabled(t *testing.T) {
	defer repositories.TruncateTestDb()
	enableMcp(t, false)

	router := chi.NewRouter()
	mountMcpRoutes(router)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/.well-known/oauth-protected-resource"},
		{http.MethodGet, "/.well-known/oauth-authorization-server"},
		{http.MethodGet, "/mcp"},
	}

	for _, tc := range cases {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(tc.method, tc.path, nil))

		if recorder.Code != http.StatusNotFound {
			t.Errorf("expected 404 from %s when MCP disabled, got %d", tc.path, recorder.Code)
		}
	}
}
