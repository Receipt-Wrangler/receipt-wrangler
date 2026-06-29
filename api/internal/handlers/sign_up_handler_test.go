package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

func teardownSignUpTests() {
	repositories.TruncateTestDb()
}

func TestShouldNotAllowUserToSignUpIfDisabled(t *testing.T) {
	defer teardownSignUpTests()
	reader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	db := repositories.GetDB()
	db.Create(&models.SystemSettings{EnableLocalSignUp: false})

	expectedResponseCode := http.StatusNotFound

	SignUp(w, r)

	if w.Result().StatusCode != expectedResponseCode {
		utils.PrintTestError(t, w.Result().StatusCode, expectedResponseCode)
	}
}

// Public sign-up must never honor a caller-supplied role: the role is decided by
// CreateUser, not the request body. An existing user occupies the first-user slot
// here, so an attacker supplying a modern AppRoleID must still end up with no
// caller-chosen app role.
func TestSignUpStripsCallerSuppliedRole(t *testing.T) {
	defer teardownSignUpTests()
	db := repositories.GetDB()
	db.Create(&models.SystemSettings{EnableLocalSignUp: true})
	db.Create(&models.User{Username: "existing", Password: "x"})

	escalateRoleId := uint(999)
	cmd := commands.SignUpCommand{
		Username:    "escalate",
		Password:    "password",
		DisplayName: "Escalate",
		AppRoleID:   &escalateRoleId,
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", nil)
	r = r.WithContext(context.WithValue(r.Context(), "signUpCommand", cmd))

	SignUp(w, r)

	if w.Result().StatusCode != http.StatusOK {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusOK)
		return
	}

	var created models.User
	if err := db.Where("username = ?", "escalate").First(&created).Error; err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	// The caller-supplied AppRoleID was stripped, so no app role is assigned (no
	// roles are seeded in this test, so the default path resolves to nil) — and
	// certainly not the bogus id the caller tried to inject.
	if created.AppRoleID != nil {
		utils.PrintTestError(t, created.AppRoleID, nil)
	}
}

// TODO: Fix how setting data for this endpoint works, then implement this tests
/*func TestShouldProcessSignUpCommand(t *testing.T) {
	defer teardownSignUpTests()
	db := repositories.GetDB()

	db.Create(&models.SystemSettings{EnableLocalSignUp: true})

	tests := map[string]struct {
		input  commands.SignUpCommand
		expect int
	}{
		"empty body": {
			expect: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		bytes, _ := json.Marshal(test.input)
		reader := strings.NewReader(string(bytes))
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api", reader)

		SignUp(w, r)

		if w.Result().StatusCode != test.expect {
			utils.PrintTestError(t, w.Result().StatusCode, test.expect)
		}
	}
}
*/
