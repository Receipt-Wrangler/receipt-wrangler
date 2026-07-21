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

func setupCategoriesTest() {
	repositories.CreateTestCategories()
}

func tearDownCategoriesTest() {
	repositories.TruncateTestDb()
}

func TestShouldGetAllCategories(t *testing.T) {
	defer tearDownCategoriesTest()
	categories := make([]models.Category, 0)
	setupCategoriesTest()

	reader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	GetAllCategories(w, r)

	err := json.Unmarshal(w.Body.Bytes(), &categories)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}

	for index, category := range categories {
		id := index + 1
		if uint(id) != category.ID {
			utils.PrintTestError(t, category.ID, id)
		}
	}

	tearDownCategoriesTest()
}

func TestShouldCreateCategory(t *testing.T) {
	defer tearDownCategoriesTest()
	category := models.Category{}
	setupCategoriesTest()

	reader := strings.NewReader(`{"name": "Test category", "description": "Test description"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	CreateCategory(w, r)

	err := json.Unmarshal(w.Body.Bytes(), &category)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}

	if category.Name != "Test category" {
		utils.PrintTestError(t, category.Name, "Test category")
	}

	if category.Description != "Test description" {
		utils.PrintTestError(t, category.Description, "Test description")
	}

	if category.ID != 4 {
		utils.PrintTestError(t, category.ID, 4)
	}
}

func TestShouldUpdateCategoryIfAdmin(t *testing.T) {
	defer tearDownCategoriesTest()
	category := models.Category{}
	setupCategoriesTest()

	reader := strings.NewReader(`{"name": "Updated Category name", "description": "Updated Test description"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryId", "1")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	UpdateCategory(w, r)

	err := json.Unmarshal(w.Body.Bytes(), &category)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if category.Name != "Updated Category name" {
		utils.PrintTestError(t, category.Name, "Updated Category name")
	}

	if category.Description != "Updated Test description" {
		utils.PrintTestError(t, category.Description, "Updated Test description")
	}
}

func TestShouldNotUpdateCategoryDueToRole(t *testing.T) {
	defer tearDownCategoriesTest()
	category := models.Category{}
	setupCategoriesTest()

	reader := strings.NewReader(`{"name": "Updated Category name", "description": "Updated Test description"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryId", "1")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	UpdateCategory(w, r)

	err := json.Unmarshal(w.Body.Bytes(), &category)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestShouldDeleteCategoryIfAdmin(t *testing.T) {
	defer tearDownCategoriesTest()
	setupCategoriesTest()

	reader := strings.NewReader(``)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryId", "1")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	DeleteCategory(w, r)

	db := repositories.GetDB()
	err := db.Model(models.Category{}).Where("id = ?", 1).First(&models.Category{}).Error
	if err == nil {
		utils.PrintTestError(t, err, "Record should not exist")
	}

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
	}
}

func TestShouldNotDeleteCategoryDueToRole(t *testing.T) {
	defer tearDownCategoriesTest()
	setupCategoriesTest()

	reader := strings.NewReader(``)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryId", "1")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	DeleteCategory(w, r)

	if w.Result().StatusCode != 403 {
		utils.PrintTestError(t, w.Result().StatusCode, 403)
	}
}

func TestShouldGetCategoryNameCountIfAdmin(t *testing.T) {
	defer tearDownCategoriesTest()
	expectedStatus := 200
	var resultCount uint
	setupCategoriesTest()

	reader := strings.NewReader(``)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryName", "test")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	GetCategoryNameCount(w, r)

	err := json.Unmarshal(w.Body.Bytes(), &resultCount)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if w.Result().StatusCode != expectedStatus {
		utils.PrintTestError(t, w.Result().StatusCode, expectedStatus)
	}

	if resultCount != 1 {
		utils.PrintTestError(t, resultCount, 1)
	}
}

func TestShouldGetCategoryNameCountIfAdmin2(t *testing.T) {
	defer tearDownCategoriesTest()
	expectedStatus := 200
	var resultCount uint
	setupCategoriesTest()

	reader := strings.NewReader(``)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryName", "totally a category name")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	GetCategoryNameCount(w, r)

	err := json.Unmarshal(w.Body.Bytes(), &resultCount)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if w.Result().StatusCode != expectedStatus {
		utils.PrintTestError(t, w.Result().StatusCode, expectedStatus)
	}

	if resultCount != 0 {
		utils.PrintTestError(t, resultCount, 0)
	}
}

func TestShouldGetCategoryNameCountAsUser(t *testing.T) {
	defer tearDownCategoriesTest()
	expectedStatus := 200
	var resultCount uint
	setupCategoriesTest()

	reader := strings.NewReader(``)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryName", "test")

	routeContext := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
	r = r.WithContext(routeContext)
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}})
	r = r.WithContext(newContext)

	grantAllAppPerms(t, 2)

	GetCategoryNameCount(w, r)

	if w.Result().StatusCode != expectedStatus {
		utils.PrintTestError(t, w.Result().StatusCode, expectedStatus)
	}

	err := json.Unmarshal(w.Body.Bytes(), &resultCount)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if resultCount != 1 {
		utils.PrintTestError(t, resultCount, 1)
	}
}

// requestCategoryNameCount issues GetCategoryNameCount for the given path
// parameter (as user 2) and returns the response status code and decoded count.
func requestCategoryNameCount(t *testing.T, categoryName string) (int, uint) {
	t.Helper()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(``))

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("categoryName", categoryName)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 2}}))

	GetCategoryNameCount(w, r)

	var count uint
	_ = json.Unmarshal(w.Body.Bytes(), &count)
	return w.Result().StatusCode, count
}

// TestShouldNotAllowSqlInjectionInCategoryNameCount is a regression test for
// GHSA-q6h3-4g3r-gg2x. The category name reaches the count query as a bound
// parameter, so a SQL tautology in the path is treated as a literal name and
// matches nothing. Before the fix the clause was built with fmt.Sprintf and this
// payload interpolated to `name = 'test' OR '1'='1'`, matching every seeded row.
func TestShouldNotAllowSqlInjectionInCategoryNameCount(t *testing.T) {
	defer tearDownCategoriesTest()
	expectedStatus := 200
	setupCategoriesTest()

	grantAllAppPerms(t, 2)

	// Positive baseline: a legit lookup of a seeded category returns 1. This
	// keeps the injection assertion below from passing vacuously if the seed
	// ever produced no rows (a still-vulnerable handler would also return 0).
	baselineStatus, baselineCount := requestCategoryNameCount(t, "test")
	if baselineStatus != expectedStatus {
		utils.PrintTestError(t, baselineStatus, expectedStatus)
	}
	if baselineCount != 1 {
		utils.PrintTestError(t, baselineCount, 1)
	}

	// The tautology is bound as a literal name and matches nothing.
	injectionStatus, injectionCount := requestCategoryNameCount(t, "test' OR '1'='1")
	if injectionStatus != expectedStatus {
		utils.PrintTestError(t, injectionStatus, expectedStatus)
	}
	if injectionCount != 0 {
		utils.PrintTestError(t, injectionCount, 0)
	}
}
