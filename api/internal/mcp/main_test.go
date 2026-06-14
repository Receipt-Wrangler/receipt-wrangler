package mcp

import (
	"fmt"
	"os"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/auth"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	defer teardown()
	repositories.SetUpTestEnv()
	repositories.InitTestDb()
	repositories.MakeMigrations()
	return m.Run(), nil
}

func teardown() {
	repositories.TestTeardown()
}

// requestForUser builds a CallToolRequest whose bearer TokenInfo carries the
// claims for the given user id, simulating a verified MCP request.
func requestForUser(userId uint) *mcpsdk.CallToolRequest {
	claims := &structs.Claims{UserId: userId}
	return &mcpsdk.CallToolRequest{
		Extra: &mcpsdk.RequestExtra{
			TokenInfo: &auth.TokenInfo{Extra: map[string]any{claimsKey: claims}},
		},
	}
}

func createUser(t *testing.T, username string) models.User {
	t.Helper()
	user := models.User{Username: username, DisplayName: username, Password: "x"}
	if err := repositories.GetDB().Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	return user
}
