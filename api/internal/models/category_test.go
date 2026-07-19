package models

import (
	"net/http/httptest"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

func TestCategory_LoadDataFromRequest(t *testing.T) {
	body := `{"name": "Groceries", "description": "Food and household"}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var category Category
	err := category.LoadDataFromRequest(w, r)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if category.Name != "Groceries" {
		utils.PrintTestError(t, category.Name, "Groceries")
	}
	if category.Description != "Food and household" {
		utils.PrintTestError(t, category.Description, "Food and household")
	}
}

func TestCategory_LoadDataFromRequest_MalformedJson(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("{"))
	w := httptest.NewRecorder()

	var category Category
	err := category.LoadDataFromRequest(w, r)
	if err == nil {
		utils.PrintTestError(t, err, "an unmarshal error")
	}
}

func TestCategory_LoadDataFromRequest_BodyReadError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", errReader{})
	w := httptest.NewRecorder()

	var category Category
	err := category.LoadDataFromRequest(w, r)
	if err == nil {
		utils.PrintTestError(t, err, "a body read error")
	}
}
