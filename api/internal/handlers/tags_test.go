package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

func setupTagsTest() {
	repositories.GetDB().Create(&models.Tag{Name: "test"})
}

func tearDownTagsTest() {
	repositories.TruncateTestDb()
}

// requestTagNameCount issues GetTagNameCount for the given path parameter (as
// user 2) and returns the response status code and decoded count.
func requestTagNameCount(t *testing.T, tagName string) (int, uint) {
	t.Helper()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(``))

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("tagName", tagName)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}}))

	GetTagNameCount(w, r)

	var count uint
	_ = json.Unmarshal(w.Body.Bytes(), &count)
	return w.Result().StatusCode, count
}

func TestShouldGetTagNameCount(t *testing.T) {
	defer tearDownTagsTest()
	expectedStatus := 200
	setupTagsTest()

	grantAllAppPerms(t, 2)

	status, count := requestTagNameCount(t, "test")
	if status != expectedStatus {
		utils.PrintTestError(t, status, expectedStatus)
	}
	if count != 1 {
		utils.PrintTestError(t, count, 1)
	}
}

// TestShouldNotAllowSqlInjectionInTagNameCount is a regression test for
// GHSA-q6h3-4g3r-gg2x. The tag name reaches the count query as a bound
// parameter, so a SQL tautology in the path is treated as a literal name and
// matches nothing. Before the fix the clause was built with fmt.Sprintf and this
// payload interpolated to `name = 'test' OR '1'='1'`, matching every seeded row.
func TestShouldNotAllowSqlInjectionInTagNameCount(t *testing.T) {
	defer tearDownTagsTest()
	expectedStatus := 200
	setupTagsTest()

	grantAllAppPerms(t, 2)

	// Positive baseline: a legit lookup of the seeded tag returns 1, so the
	// injection assertion below cannot pass vacuously on an empty seed.
	baselineStatus, baselineCount := requestTagNameCount(t, "test")
	if baselineStatus != expectedStatus {
		utils.PrintTestError(t, baselineStatus, expectedStatus)
	}
	if baselineCount != 1 {
		utils.PrintTestError(t, baselineCount, 1)
	}

	// The tautology is bound as a literal name and matches nothing.
	injectionStatus, injectionCount := requestTagNameCount(t, "test' OR '1'='1")
	if injectionStatus != expectedStatus {
		utils.PrintTestError(t, injectionStatus, expectedStatus)
	}
	if injectionCount != 0 {
		utils.PrintTestError(t, injectionCount, 0)
	}
}
