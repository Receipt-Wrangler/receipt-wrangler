package services

import (
	"context"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

func setMcpSettings(t *testing.T, enabled bool, publicUrl string) {
	t.Helper()

	settings, err := repositories.NewSystemSettingsRepository(nil).GetSystemSettings()
	if err != nil {
		t.Fatalf("failed to load system settings: %v", err)
	}

	err = repositories.GetDB().
		Model(&models.SystemSettings{}).
		Where("id = ?", settings.ID).
		Updates(map[string]interface{}{"mcp_enabled": enabled, "mcp_public_url": publicUrl}).Error
	if err != nil {
		t.Fatalf("failed to set mcp system settings: %v", err)
	}
}

func TestNormalizeMcpPublicUrl(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"empty falls back to default", "", DefaultMcpPublicUrl},
		{"whitespace falls back to default", "   ", DefaultMcpPublicUrl},
		{"clean origin is returned as-is", "https://receipts.example.com", "https://receipts.example.com"},
		{"trailing slash is trimmed", "https://receipts.example.com/", "https://receipts.example.com"},
		{"path is stripped to the origin", "https://receipts.example.com/mcp?x=1#y", "https://receipts.example.com"},
		{"missing scheme falls back to default", "receipts.example.com", DefaultMcpPublicUrl},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NormalizeMcpPublicUrl(test.raw); got != test.want {
				t.Errorf("NormalizeMcpPublicUrl(%q) = %q, want %q", test.raw, got, test.want)
			}
		})
	}
}

func TestIsMcpEnabledReadsSystemSettings(t *testing.T) {
	defer repositories.TruncateTestDb()

	setMcpSettings(t, false, "")
	if IsMcpEnabled() {
		t.Errorf("expected MCP disabled when the setting is off")
	}

	setMcpSettings(t, true, "https://receipts.example.com")
	if !IsMcpEnabled() {
		t.Errorf("expected MCP enabled when the setting is on")
	}
}

func TestGetMcpResourceUrlReadsSystemSettings(t *testing.T) {
	defer repositories.TruncateTestDb()

	setMcpSettings(t, true, "https://receipts.example.com/")

	if got := GetMcpPublicUrl(); got != "https://receipts.example.com" {
		t.Errorf("GetMcpPublicUrl() = %q, want %q", got, "https://receipts.example.com")
	}
	if got := GetMcpResourceUrl(); got != "https://receipts.example.com/mcp" {
		t.Errorf("GetMcpResourceUrl() = %q, want %q", got, "https://receipts.example.com/mcp")
	}
}

// TestMcpTokenAudienceBinding proves an MCP token validates only against the
// MCP audience and an MCP refresh token is rejected by the normal REST
// validator (which is what stops trading it for a full-access token).
func TestMcpTokenAudienceBinding(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := models.User{Username: "mcpaud", DisplayName: "mcpaud", Password: "x"}
	if err := repositories.GetDB().Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	const audience = "https://receipts.example.com/mcp"
	access, refresh, _, err := GenerateMcpJWT(user.ID, audience)
	if err != nil {
		t.Fatalf("failed to mint mcp jwt: %v", err)
	}

	mcpValidator, err := InitMcpTokenValidator(audience)
	if err != nil {
		t.Fatalf("failed to init mcp validator: %v", err)
	}
	if _, err := mcpValidator.ValidateToken(context.Background(), access); err != nil {
		t.Errorf("mcp validator rejected an mcp-audience access token: %v", err)
	}

	restValidator, err := InitTokenValidator()
	if err != nil {
		t.Fatalf("failed to init rest validator: %v", err)
	}
	if _, err := restValidator.ValidateToken(context.Background(), access); err == nil {
		t.Errorf("rest validator accepted an mcp-audience access token (audience not bound)")
	}
	if _, err := restValidator.ValidateToken(context.Background(), refresh); err == nil {
		t.Errorf("rest validator accepted an mcp-audience refresh token; it could be traded for a full-access token")
	}
}
