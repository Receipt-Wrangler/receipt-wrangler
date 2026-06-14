package oauth

import (
	"fmt"
	"os"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
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

// createTestUserWithPassword creates a user whose password is bcrypt-hashed so
// services.LoginUser can authenticate it, and returns the user.
func createTestUserWithPassword(t *testing.T, username string, password string) models.User {
	t.Helper()

	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := models.User{
		Username:    username,
		DisplayName: username,
		Password:    string(hashed),
	}

	if err := repositories.GetDB().Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}
