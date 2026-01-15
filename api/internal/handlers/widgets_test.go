package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

func setupWidgetsTest() {
	repositories.CreateTestGroupWithUsers()
	repositories.CreateTestCategories()
}

func tearDownWidgetsTest() {
	repositories.TruncateTestDb()
}

func createWidgetTestReceipt(name string, amount float64, paidByUserId uint, groupId uint, categories []models.Category, tags []models.Tag) models.Receipt {
	db := repositories.GetDB()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromFloat(amount),
		Date:         time.Now(),
		PaidByUserID: paidByUserId,
		GroupId:      groupId,
		Status:       models.OPEN,
		Categories:   categories,
		Tags:         tags,
	}
	db.Create(&receipt)
	return receipt
}

func TestGetPieChartData_Success_Categories(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	db := repositories.GetDB()

	// Get test category
	var category1 models.Category
	db.First(&category1, 1)

	// Create test receipt
	createWidgetTestReceipt("Test Receipt", 100.00, 1, 1, []models.Category{category1}, nil)

	// Create request body
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.PieChartData
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Value != 100.0 {
		utils.PrintTestError(t, result.Data[0].Value, 100.0)
	}
}

func TestGetPieChartData_Success_Tags(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	db := repositories.GetDB()

	// Create test tag
	tag := models.Tag{Name: "TestTag", Description: "Test"}
	db.Create(&tag)

	// Create test receipt
	createWidgetTestReceipt("Test Receipt", 150.00, 1, 1, nil, []models.Tag{tag})

	// Create request body
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_TAGS,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.PieChartData
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Value != 150.0 {
		utils.PrintTestError(t, result.Data[0].Value, 150.0)
	}
}

func TestGetPieChartData_Success_PaidBy(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	// Create test receipts paid by different users
	createWidgetTestReceipt("Receipt 1", 100.00, 1, 1, nil, nil)
	createWidgetTestReceipt("Receipt 2", 50.00, 2, 1, nil, nil)

	// Create request body
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_PAIDBY,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.PieChartData
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Verify total amount
	totalAmount := 0.0
	for _, dp := range result.Data {
		totalAmount += dp.Value
	}

	if totalAmount != 150.0 {
		utils.PrintTestError(t, totalAmount, 150.0)
	}
}

func TestGetPieChartData_InvalidChartGrouping(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	// Create request body with invalid chart grouping
	requestBody := map[string]string{
		"chartGrouping": "INVALID_GROUPING",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	// Should return 400 Bad Request for validation error
	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
		return
	}

	// Verify error response contains validation error
	// The response is a map[string]string, not a ValidatorError struct
	var errorResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if _, exists := errorResponse["chartGrouping"]; !exists {
		utils.PrintTestError(t, "missing chartGrouping error", "expected chartGrouping validation error")
	}
}

func TestGetPieChartData_EmptyReceipts(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	// Don't create any receipts

	// Create request body
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.PieChartData
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(result.Data) != 0 {
		utils.PrintTestError(t, len(result.Data), 0)
	}
}

func TestGetPieChartData_WithFilter(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	db := repositories.GetDB()

	// Get test category
	var category1 models.Category
	db.First(&category1, 1)

	// Create receipts with different statuses
	r1 := createWidgetTestReceipt("Open Receipt", 100.00, 1, 1, []models.Category{category1}, nil)
	r2 := createWidgetTestReceipt("Resolved Receipt", 50.00, 1, 1, []models.Category{category1}, nil)

	// Update statuses
	db.Model(&r1).Update("status", models.OPEN)
	db.Model(&r2).Update("status", models.RESOLVED)

	// Create request body with filter
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
		Filter: commands.ReceiptPagedRequestFilter{
			Status: commands.PagedRequestField{
				Value:     []interface{}{"OPEN"},
				Operation: commands.CONTAINS,
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.PieChartData
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Should only have the OPEN receipt
	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Value != 100.0 {
		utils.PrintTestError(t, result.Data[0].Value, 100.0)
	}
}

func TestGetPieChartData_InvalidJSON(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	// Send invalid JSON
	reader := strings.NewReader("not valid json")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	// Should return 500 Internal Server Error for JSON parse error
	if w.Result().StatusCode != 500 {
		utils.PrintTestError(t, w.Result().StatusCode, 500)
	}
}

func TestGetPieChartData_EmptyBody(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	// Send empty body - should fail validation since chartGrouping is required
	reader := strings.NewReader("{}")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	// Should return 400 Bad Request for validation error
	if w.Result().StatusCode != 400 {
		utils.PrintTestError(t, w.Result().StatusCode, 400)
	}
}

func TestGetPieChartData_AdminUser(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	db := repositories.GetDB()

	// Get test category
	var category1 models.Category
	db.First(&category1, 1)

	// Create test receipt
	createWidgetTestReceipt("Admin Test Receipt", 200.00, 1, 1, []models.Category{category1}, nil)

	// Create request body
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context with ADMIN role
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.ADMIN},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var result structs.PieChartData
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Value != 200.0 {
		utils.PrintTestError(t, result.Data[0].Value, 200.0)
	}
}

func TestGetPieChartData_ResponseFormat(t *testing.T) {
	defer tearDownWidgetsTest()
	setupWidgetsTest()

	db := repositories.GetDB()

	// Get test category
	var category1 models.Category
	db.First(&category1, 1)

	// Create test receipt
	createWidgetTestReceipt("Format Test Receipt", 123.45, 1, 1, []models.Category{category1}, nil)

	// Create request body
	requestBody := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	reader := strings.NewReader(string(body))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/widget/1/pie-chart", reader)

	// Add path parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("groupId", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Add JWT context
	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{
		CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER},
	})
	r = r.WithContext(newContext)

	GetPieChartData(w, r)

	// Verify response structure by parsing as raw map
	var rawResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &rawResponse)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	// Verify "data" key exists
	if _, exists := rawResponse["data"]; !exists {
		utils.PrintTestError(t, "missing data key", "expected data key in response")
		return
	}

	// Verify data is an array
	dataArray, ok := rawResponse["data"].([]interface{})
	if !ok {
		utils.PrintTestError(t, "data is not an array", "expected data to be an array")
		return
	}

	// Verify each data point has label and value
	for _, item := range dataArray {
		dataPoint, ok := item.(map[string]interface{})
		if !ok {
			utils.PrintTestError(t, "data point is not an object", "expected data point to be an object")
			return
		}

		if _, exists := dataPoint["label"]; !exists {
			utils.PrintTestError(t, "missing label key", "expected label key in data point")
		}

		if _, exists := dataPoint["value"]; !exists {
			utils.PrintTestError(t, "missing value key", "expected value key in data point")
		}
	}
}
