package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

func setupCustomFieldHandlerTest() {
	createTestCustomFields()
}

func createTestCustomFields() {
	db := repositories.GetDB()

	// Create custom field with TEXT type
	textField := models.CustomField{
		Name:        "Test Text Field",
		Type:        models.TEXT,
		Description: "A test text field",
	}
	db.Create(&textField)

	// Create custom field with DATE type
	dateField := models.CustomField{
		Name:        "Test Date Field",
		Type:        models.DATE,
		Description: "A test date field",
	}
	db.Create(&dateField)

	// Create custom field with SELECT type and options
	selectField := models.CustomField{
		Name:        "Test Select Field",
		Type:        models.SELECT,
		Description: "A test select field",
	}
	db.Create(&selectField)

	// Add options to the SELECT field
	option1 := models.CustomFieldOption{
		Value:         "Option 1",
		CustomFieldId: selectField.ID,
	}
	option2 := models.CustomFieldOption{
		Value:         "Option 2",
		CustomFieldId: selectField.ID,
	}
	db.Create(&option1)
	db.Create(&option2)

	// Create custom field with CURRENCY type
	currencyField := models.CustomField{
		Name:        "Test Currency Field",
		Type:        models.CURRENCY,
		Description: "A test currency field",
	}
	db.Create(&currencyField)
}

func teardownCustomFieldHandlerTest() {
	repositories.TruncateTestDb()
}

func createCustomFieldHandlerTestRequest(method, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	var bodyReader *strings.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	} else {
		bodyReader = strings.NewReader("")
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	// Add JWT context with mock user claims
	var vClaims validator.ValidatedClaims
	vClaims.CustomClaims = &structs.Claims{
		UserId:      1,
		Username:    "testuser",
		Displayname: "Test User",
	}

	// Set the JWT context
	ctx := req.Context()
	ctx = context.WithValue(ctx, jwtmiddleware.ContextKey{}, &vClaims)
	req = req.WithContext(ctx)

	return req, httptest.NewRecorder()
}

func TestGetPagedCustomFieldsHandler(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with paged request payload
	pagedRequestCommand := commands.PagedRequestCommand{
		Page:          1,
		PageSize:      10,
		OrderBy:       "name",
		SortDirection: commands.ASCENDING,
	}
	pagedRequestJSON, _ := json.Marshal(pagedRequestCommand)

	req, rr := createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		string(pagedRequestJSON),
	)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Parse response
	var response structs.PagedData
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Should return 4 custom fields
	if response.TotalCount != 4 {
		utils.PrintTestError(t, response.TotalCount, 4)
	}

	if len(response.Data) != 4 {
		utils.PrintTestError(t, len(response.Data), 4)
	}

	// Convert data to CustomField array
	customFieldsBytes, _ := json.Marshal(response.Data)
	var customFields []models.CustomField
	json.Unmarshal(customFieldsBytes, &customFields)

	// Check if fields are correctly ordered by name
	if customFields[0].Name != "Test Currency Field" {
		utils.PrintTestError(t, customFields[0].Name, "Test Currency Field")
	}

	if customFields[1].Name != "Test Date Field" {
		utils.PrintTestError(t, customFields[1].Name, "Test Date Field")
	}

	if customFields[2].Name != "Test Select Field" {
		utils.PrintTestError(t, customFields[2].Name, "Test Select Field")
	}

	if customFields[3].Name != "Test Text Field" {
		utils.PrintTestError(t, customFields[3].Name, "Test Text Field")
	}
}

func TestGetPagedCustomFieldsHandlerWithPagination(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with paged request payload
	pagedRequestCommand := commands.PagedRequestCommand{
		Page:          1,
		PageSize:      2,
		OrderBy:       "name",
		SortDirection: commands.ASCENDING,
	}
	pagedRequestJSON, _ := json.Marshal(pagedRequestCommand)

	req, rr := createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		string(pagedRequestJSON),
	)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Parse response
	var response structs.PagedData
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Total count should still be 4
	if response.TotalCount != 4 {
		utils.PrintTestError(t, response.TotalCount, 4)
	}

	// But we should only get 2 items per page
	if len(response.Data) != 2 {
		utils.PrintTestError(t, len(response.Data), 2)
	}

	// Test second page
	pagedRequestCommand.Page = 2
	pagedRequestJSON, _ = json.Marshal(pagedRequestCommand)

	req, rr = createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		string(pagedRequestJSON),
	)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Parse response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Total count should still be 4
	if response.TotalCount != 4 {
		utils.PrintTestError(t, response.TotalCount, 4)
	}

	// And we should get the other 2 items
	if len(response.Data) != 2 {
		utils.PrintTestError(t, len(response.Data), 2)
	}
}

func TestGetPagedCustomFieldsHandlerWithTypeOrdering(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with paged request payload
	pagedRequestCommand := commands.PagedRequestCommand{
		Page:          1,
		PageSize:      10,
		OrderBy:       "type",
		SortDirection: commands.ASCENDING,
	}
	pagedRequestJSON, _ := json.Marshal(pagedRequestCommand)

	req, rr := createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		string(pagedRequestJSON),
	)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Parse response
	var response structs.PagedData
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Convert data to CustomField array
	customFieldsBytes, _ := json.Marshal(response.Data)
	var customFields []models.CustomField
	json.Unmarshal(customFieldsBytes, &customFields)

	// Check ordering by type
	if customFields[0].Type != models.CURRENCY {
		utils.PrintTestError(t, string(customFields[0].Type), string(models.CURRENCY))
	}
}

func TestGetPagedCustomFieldsHandlerWithInvalidOrderBy(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with invalid orderBy
	pagedRequestCommand := commands.PagedRequestCommand{
		Page:          1,
		PageSize:      10,
		OrderBy:       "invalid_field",
		SortDirection: commands.ASCENDING,
	}
	pagedRequestJSON, _ := json.Marshal(pagedRequestCommand)

	req, rr := createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		string(pagedRequestJSON),
	)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status - should be error
	if status := rr.Code; status != http.StatusInternalServerError {
		utils.PrintTestError(t, status, http.StatusInternalServerError)
	}

	// Parse error response
	var errorResponse map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Check error message
	if errorMsg, exists := errorResponse["errorMsg"]; !exists || !strings.Contains(errorMsg, "Error getting custom fields") {
		utils.PrintTestError(t, errorMsg, "Error getting custom fields")
	}
}

func TestGetPagedCustomFieldsHandlerWithInvalidBody(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with invalid JSON body
	req, rr := createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		"{invalid json",
	)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status - should be error
	if status := rr.Code; status != http.StatusInternalServerError {
		utils.PrintTestError(t, status, http.StatusInternalServerError)
	}
}

func TestGetPagedCustomFieldsHandlerWithEmptyBody(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with empty body
	req, rr := createCustomFieldHandlerTestRequest(
		"POST",
		"/api/customField/getPagedCustomFields",
		"",
	)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetPagedCustomFields(rr, req)

	// Check response status - should be error
	if status := rr.Code; status != http.StatusInternalServerError {
		utils.PrintTestError(t, status, http.StatusInternalServerError)
	}
}

func TestGetCustomFieldByIdHandler(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Get a custom field from the test DB to use its ID
	db := repositories.GetDB()
	var customField models.CustomField
	db.Where("name = ?", "Test Text Field").First(&customField)

	// Create request to get custom field by ID
	req, rr := createCustomFieldHandlerTestRequest(
		"GET",
		"/api/customField/"+utils.UintToString(customField.ID),
		"",
	)

	// Set URL parameter
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", utils.UintToString(customField.ID))

	grantAllAppPerms(t, 1)

	// Call the handler
	GetCustomFieldById(rr, req)

	// Check response status
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Parse response
	var responseCustomField models.CustomField
	err := json.Unmarshal(rr.Body.Bytes(), &responseCustomField)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Check if the returned custom field has the correct ID and properties
	if responseCustomField.ID != customField.ID {
		utils.PrintTestError(t, responseCustomField.ID, customField.ID)
	}

	if responseCustomField.Name != "Test Text Field" {
		utils.PrintTestError(t, responseCustomField.Name, "Test Text Field")
	}

	if responseCustomField.Type != models.TEXT {
		utils.PrintTestError(t, string(responseCustomField.Type), string(models.TEXT))
	}
}

func TestGetCustomFieldByIdHandlerWithInvalidId(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with invalid ID
	req, rr := createCustomFieldHandlerTestRequest(
		"GET",
		"/api/customField/invalid",
		"",
	)

	// Set URL parameter
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "invalid")

	grantAllAppPerms(t, 1)

	// Call the handler
	GetCustomFieldById(rr, req)

	// Check response status - should be error
	if status := rr.Code; status != http.StatusInternalServerError {
		utils.PrintTestError(t, status, http.StatusInternalServerError)
	}

	// Parse error response
	var errorResponse map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Check error message
	if errorMsg, exists := errorResponse["errorMsg"]; !exists || !strings.Contains(errorMsg, "Error getting custom field") {
		utils.PrintTestError(t, errorMsg, "Error getting custom field")
	}
}

func TestGetCustomFieldByIdHandlerWithNonExistentId(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Use a non-existent ID (999)
	nonExistentId := "999"

	// Create request with non-existent ID
	req, rr := createCustomFieldHandlerTestRequest(
		"GET",
		"/api/customField/"+nonExistentId,
		"",
	)

	// Set URL parameter
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", nonExistentId)

	grantAllAppPerms(t, 1)

	// Call the handler
	GetCustomFieldById(rr, req)

	// Check response status - should be error
	if status := rr.Code; status != http.StatusInternalServerError {
		utils.PrintTestError(t, status, http.StatusInternalServerError)
	}

	// Parse error response
	var errorResponse map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Check error message
	if errorMsg, exists := errorResponse["errorMsg"]; !exists || !strings.Contains(errorMsg, "Error getting custom field") {
		utils.PrintTestError(t, errorMsg, "Error getting custom field")
	}
}

func TestDeleteCustomFieldHandler(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	// Get a custom field from the test DB to use its ID
	var customField models.CustomField
	db.Where("name = ?", "Test Text Field").First(&customField)
	customFieldId := customField.ID

	// Create request to delete custom field
	req, rr := createCustomFieldHandlerTestRequest(
		"DELETE",
		"/api/customField/"+utils.UintToString(customFieldId),
		"",
	)

	// Set URL parameter
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", utils.UintToString(customFieldId))

	grantAllAppPerms(t, 1)

	// Call the handler
	DeleteCustomField(rr, req)

	// Check response status
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}

	// Verify the custom field was deleted
	var count int64
	db.Model(&models.CustomField{}).Where("id = ?", customFieldId).Count(&count)
	if count != 0 {
		utils.PrintTestError(t, count, 0)
	}
}

func TestDeleteCustomFieldHandlerWithInvalidId(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Create request with invalid ID
	req, rr := createCustomFieldHandlerTestRequest(
		"DELETE",
		"/api/customField/invalid",
		"",
	)

	// Set URL parameter
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "invalid")

	grantAllAppPerms(t, 1)

	// Call the handler
	DeleteCustomField(rr, req)

	// Check response status - should be error
	if status := rr.Code; status != http.StatusInternalServerError {
		utils.PrintTestError(t, status, http.StatusInternalServerError)
	}

	// Parse error response
	var errorResponse map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
	}

	// Check error message
	if errorMsg, exists := errorResponse["errorMsg"]; !exists || !strings.Contains(errorMsg, "Error deleting custom field") {
		utils.PrintTestError(t, errorMsg, "Error deleting custom field")
	}
}

func TestDeleteCustomFieldHandlerWithNonExistentId(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	// Use a non-existent ID (999)
	nonExistentId := "999"

	// Create request with non-existent ID
	req, rr := createCustomFieldHandlerTestRequest(
		"DELETE",
		"/api/customField/"+nonExistentId,
		"",
	)

	// Set URL parameter
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", nonExistentId)

	grantAllAppPerms(t, 1)

	// Call the handler
	DeleteCustomField(rr, req)

	// Should still return OK because deleting a non-existent record is not an error in this implementation
	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
	}
}

// updateCustomFieldRequest builds a PUT request for customFieldId with body as
// the JSON payload, wiring the chi URL param the handler reads.
func updateCustomFieldRequest(customFieldId string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req, rr := createCustomFieldHandlerTestRequest(
		"PUT",
		"/api/customField/"+customFieldId,
		body,
	)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", customFieldId)

	return req, rr
}

func TestUpdateCustomFieldHandler(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	var customField models.CustomField
	db.Where("name = ?", "Test Text Field").First(&customField)

	body := `{"name":"Renamed Text Field","type":"TEXT","description":"An updated description"}`
	req, rr := updateCustomFieldRequest(utils.UintToString(customField.ID), body)

	grantAllAppPerms(t, 1)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
		return
	}

	var response models.CustomField
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		utils.PrintTestError(t, err.Error(), nil)
		return
	}

	if response.Name != "Renamed Text Field" {
		utils.PrintTestError(t, response.Name, "Renamed Text Field")
	}

	if response.Description != "An updated description" {
		utils.PrintTestError(t, response.Description, "An updated description")
	}

	var persisted models.CustomField
	db.First(&persisted, customField.ID)
	if persisted.Name != "Renamed Text Field" {
		utils.PrintTestError(t, persisted.Name, "Renamed Text Field")
	}
}

func TestUpdateCustomFieldHandlerRenamesAndAppendsOptions(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	var customField models.CustomField
	db.Where("name = ?", "Test Select Field").First(&customField)

	var options []models.CustomFieldOption
	db.Where("custom_field_id = ?", customField.ID).Order("id asc").Find(&options)
	if len(options) != 2 {
		utils.PrintTestError(t, len(options), 2)
		return
	}

	body := `{"name":"Test Select Field","type":"SELECT","options":[` +
		`{"id":` + utils.UintToString(options[0].ID) + `,"value":"Renamed Option"},` +
		`{"id":` + utils.UintToString(options[1].ID) + `,"value":"Option 2"},` +
		`{"value":"Option 3"}]}`

	req, rr := updateCustomFieldRequest(utils.UintToString(customField.ID), body)

	grantAllAppPerms(t, 1)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusOK {
		utils.PrintTestError(t, status, http.StatusOK)
		return
	}

	var persistedOptions []models.CustomFieldOption
	db.Where("custom_field_id = ?", customField.ID).Order("id asc").Find(&persistedOptions)

	if len(persistedOptions) != 3 {
		utils.PrintTestError(t, len(persistedOptions), 3)
		return
	}

	// The renamed option keeps its id, so any CustomFieldValue pointing at it
	// still resolves.
	if persistedOptions[0].ID != options[0].ID {
		utils.PrintTestError(t, persistedOptions[0].ID, options[0].ID)
	}

	if persistedOptions[0].Value != "Renamed Option" {
		utils.PrintTestError(t, persistedOptions[0].Value, "Renamed Option")
	}

	if persistedOptions[2].Value != "Option 3" {
		utils.PrintTestError(t, persistedOptions[2].Value, "Option 3")
	}
}

func TestUpdateCustomFieldHandlerRejectsTypeChange(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	var customField models.CustomField
	db.Where("name = ?", "Test Text Field").First(&customField)

	body := `{"name":"Test Text Field","type":"CURRENCY"}`
	req, rr := updateCustomFieldRequest(utils.UintToString(customField.ID), body)

	grantAllAppPerms(t, 1)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		utils.PrintTestError(t, status, http.StatusBadRequest)
		return
	}

	// The stored type must be untouched -- a CustomFieldValue lives in a
	// type-specific column, so a re-type would mis-column every existing value.
	var persisted models.CustomField
	db.First(&persisted, customField.ID)
	if persisted.Type != models.TEXT {
		utils.PrintTestError(t, persisted.Type, models.TEXT)
	}
}

func TestUpdateCustomFieldHandlerRejectsForeignOptionId(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	var selectField models.CustomField
	db.Where("name = ?", "Test Select Field").First(&selectField)

	// An option id the field does not own is rejected outright, so an edit can
	// never rewrite an unrelated field's option.
	body := `{"name":"Test Select Field","type":"SELECT","options":[{"id":99999,"value":"Hijacked"}]}`
	req, rr := updateCustomFieldRequest(utils.UintToString(selectField.ID), body)

	grantAllAppPerms(t, 1)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		utils.PrintTestError(t, status, http.StatusBadRequest)
		return
	}

	var optionCount int64
	db.Model(&models.CustomFieldOption{}).Where("custom_field_id = ?", selectField.ID).Count(&optionCount)
	if optionCount != 2 {
		utils.PrintTestError(t, optionCount, 2)
	}
}

func TestUpdateCustomFieldHandlerWithInvalidBody(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	var customField models.CustomField
	db.Where("name = ?", "Test Text Field").First(&customField)

	// Name and type are both required.
	body := `{"description":"Only a description"}`
	req, rr := updateCustomFieldRequest(utils.UintToString(customField.ID), body)

	grantAllAppPerms(t, 1)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		utils.PrintTestError(t, status, http.StatusBadRequest)
	}
}

func TestUpdateCustomFieldHandlerWithNonExistentId(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	body := `{"name":"Does Not Exist","type":"TEXT"}`
	req, rr := updateCustomFieldRequest("999", body)

	grantAllAppPerms(t, 1)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		utils.PrintTestError(t, status, http.StatusNotFound)
	}
}

func TestUpdateCustomFieldHandlerWithoutUpdatePermission(t *testing.T) {
	defer teardownCustomFieldHandlerTest()
	setupCustomFieldHandlerTest()

	db := repositories.GetDB()

	var customField models.CustomField
	db.Where("name = ?", "Test Text Field").First(&customField)

	body := `{"name":"Renamed Text Field","type":"TEXT"}`
	req, rr := updateCustomFieldRequest(utils.UintToString(customField.ID), body)

	// Reading custom fields is not enough to edit one: editing a definition
	// changes it for every group's receipts, so it is a distinct permission.
	grantAppPerms(t, 1, permissions.AppCustomFieldsRead, permissions.AppCustomFieldsCreate)

	UpdateCustomField(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		utils.PrintTestError(t, status, http.StatusForbidden)
		return
	}

	var persisted models.CustomField
	db.First(&persisted, customField.ID)
	if persisted.Name != "Test Text Field" {
		utils.PrintTestError(t, persisted.Name, "Test Text Field")
	}
}
