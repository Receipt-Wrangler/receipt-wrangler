package services

import (
	"bytes"
	"encoding/csv"
	"github.com/shopspring/decimal"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"
)

func TestShouldBuildReceiptCsv(t *testing.T) {
	// BuildReceiptCsv reads the app time zone from System Settings.
	defer repositories.TruncateTestDb()

	expected :=
		"Id,Added At,Receipt Date,Name,Paid By,Amount,Status,Categories,Tags,Resolved Date\n" +
			"1,2025-01-01,2025-01-01,test,Jim,123.45,OPEN,\"Groceries,Food\",\"Bill,Essential\",2025-01-01\n"

	date := time.Date(
		2025, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewReceiptCsvService()
	receipts := []models.Receipt{
		models.Receipt{
			BaseModel: models.BaseModel{
				ID:        1,
				CreatedAt: date,
			},
			Date:       date,
			Name:       "test",
			PaidByUser: models.User{DisplayName: "Jim"},
			Amount:     decimal.NewFromFloat(123.45),
			Status:     models.OPEN,
			Categories: []models.Category{
				models.Category{Name: "Groceries"},
				models.Category{Name: "Food"},
			},
			Tags: []models.Tag{
				models.Tag{Name: "Bill"},
				models.Tag{Name: "Essential"},
			},
			ResolvedDate: &date,
		},
	}

	result, err := service.BuildReceiptCsv(receipts)
	if err != nil {
		utils.PrintTestError(t, result, expected)
	}

	bytes := result.ReceiptCsvBytes
	if string(bytes) != expected {
		utils.PrintTestError(t, string(bytes), expected)
	}
}

// A receipt whose user-controlled text columns (name, paid-by display name,
// category/tag names) begin with a spreadsheet formula lead are neutralized with a
// leading apostrophe in the export, so opening it in Excel/Sheets renders them as
// literal text rather than executing them.
func TestBuildReceiptCsvNeutralizesFormulaInjection(t *testing.T) {
	// BuildReceiptCsv reads the app time zone from System Settings.
	defer repositories.TruncateTestDb()

	date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewReceiptCsvService()
	receipts := []models.Receipt{
		{
			BaseModel:  models.BaseModel{ID: 1, CreatedAt: date},
			Date:       date,
			Name:       `=HYPERLINK("http://evil")`,
			PaidByUser: models.User{DisplayName: "+cmd"},
			Amount:     decimal.NewFromFloat(1),
			Status:     models.OPEN,
			Categories: []models.Category{{Name: "=SUM(A1)"}},
			Tags:       []models.Tag{{Name: "@danger"}},
		},
	}

	result, err := service.BuildReceiptCsv(receipts)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	records, err := csv.NewReader(bytes.NewReader(result.ReceiptCsvBytes)).ReadAll()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(records) != 2 {
		utils.PrintTestError(t, len(records), 2)
		return
	}

	// Data row columns: Id, Added At, Receipt Date, Name, Paid By, Amount, Status,
	// Categories, Tags, Resolved Date.
	row := records[1]
	assertions := []struct {
		column int
		want   string
	}{
		{3, `'=HYPERLINK("http://evil")`}, // Name
		{4, "'+cmd"},                      // Paid By
		{7, "'=SUM(A1)"},                  // Categories
		{8, "'@danger"},                   // Tags
	}
	for _, assertion := range assertions {
		if row[assertion.column] != assertion.want {
			utils.PrintTestError(t, row[assertion.column], assertion.want)
		}
	}
}

// A malicious item whose user-controlled text columns (receipt name, item name,
// charged-to display name, category/tag names) begin with a spreadsheet formula lead
// are neutralized with a leading apostrophe in the item export, so opening it in
// Excel/Sheets renders them as literal text rather than executing them.
func TestBuildItemCsvNeutralizesFormulaInjection(t *testing.T) {
	// BuildReceiptCsv reads the app time zone from System Settings.
	defer repositories.TruncateTestDb()

	date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	service := NewReceiptCsvService()
	items := []models.Item{
		{
			BaseModel: models.BaseModel{ID: 1},
			ReceiptId: 2,
			Receipt: models.Receipt{
				Name: `=HYPERLINK("http://evil")`,
				Date: date,
			},
			Name:          "+cmd",
			ChargedToUser: models.User{DisplayName: "@who"},
			Amount:        decimal.NewFromFloat(1),
			Status:        models.ITEM_OPEN,
			Categories:    []models.Category{{Name: "=SUM(A1)"}},
			Tags:          []models.Tag{{Name: "@danger"}},
		},
	}

	result, err := service.BuildItemCsv(items)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	records, err := csv.NewReader(bytes.NewReader(result)).ReadAll()
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	if len(records) != 2 {
		utils.PrintTestError(t, len(records), 2)
		return
	}

	// Data row columns: Id, Receipt Id, Receipt Name, Receipt Date, Name,
	// Charged to User, Amount, Status, Categories, Tags.
	row := records[1]
	assertions := []struct {
		column int
		want   string
	}{
		{2, `'=HYPERLINK("http://evil")`}, // Receipt Name
		{4, "'+cmd"},                      // Name
		{5, "'@who"},                      // Charged to User
		{8, "'=SUM(A1)"},                  // Categories
		{9, "'@danger"},                   // Tags
	}
	for _, assertion := range assertions {
		if row[assertion.column] != assertion.want {
			utils.PrintTestError(t, row[assertion.column], assertion.want)
		}
	}
}

func TestShouldBuildItemCsv(t *testing.T) {
	// BuildReceiptCsv reads the app time zone from System Settings.
	defer repositories.TruncateTestDb()

	expected :=
		"Id,Receipt Id,Receipt Name,Receipt Date,Name,Charged to User,Amount,Status,Categories,Tags\n" +
			"1,2,Test Receipt,2025-01-01,Test Item,John,25.5,OPEN,\"Groceries,Food\",\"Essential,Bill\"\n" +
			"2,3,Another Receipt,2025-01-02,Another Item,Jane,15.75,RESOLVED,Electronics,Gadget\n"

	date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)

	service := NewReceiptCsvService()
	items := []models.Item{
		{
			BaseModel: models.BaseModel{ID: 1},
			ReceiptId: 2,
			Receipt: models.Receipt{
				Name: "Test Receipt",
				Date: date1,
			},
			Name:          "Test Item",
			ChargedToUser: models.User{DisplayName: "John"},
			Amount:        decimal.NewFromFloat(25.50),
			Status:        models.ITEM_OPEN,
			Categories: []models.Category{
				{Name: "Groceries"},
				{Name: "Food"},
			},
			Tags: []models.Tag{
				{Name: "Essential"},
				{Name: "Bill"},
			},
		},
		{
			BaseModel: models.BaseModel{ID: 2},
			ReceiptId: 3,
			Receipt: models.Receipt{
				Name: "Another Receipt",
				Date: date2,
			},
			Name:          "Another Item",
			ChargedToUser: models.User{DisplayName: "Jane"},
			Amount:        decimal.NewFromFloat(15.75),
			Status:        models.ITEM_RESOLVED,
			Categories: []models.Category{
				{Name: "Electronics"},
			},
			Tags: []models.Tag{
				{Name: "Gadget"},
			},
		},
	}

	result, err := service.BuildItemCsv(items)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if string(result) != expected {
		utils.PrintTestError(t, string(result), expected)
	}
}

// Added At and Resolved Date are instants, printed as their day in the app time
// zone; the receipt Date is a calendar day stored as midnight UTC and printed in
// UTC whatever the zone.
func TestBuildReceiptCsvPrintsInstantsInTheAppZone(t *testing.T) {
	defer repositories.TruncateTestDb()

	if err := repositories.SetAppTimeZoneForTests("America/New_York"); err != nil {
		t.Fatal(err)
	}

	addedAt := time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC)  // Sep 30, 9PM in New York
	resolved := time.Date(2026, 11, 2, 3, 0, 0, 0, time.UTC) // Nov 1, 11PM in New York
	receipt := models.Receipt{
		BaseModel:    models.BaseModel{ID: 1, CreatedAt: addedAt},
		Date:         time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		Name:         "late",
		Amount:       decimal.NewFromInt(1),
		Status:       models.OPEN,
		ResolvedDate: &resolved,
		ReceiptItems: []models.Item{{BaseModel: models.BaseModel{ID: 9}, ReceiptId: 1, Name: "item", Amount: decimal.NewFromInt(1)}},
	}
	receipt.ReceiptItems[0].Receipt = receipt

	service := NewReceiptCsvService()
	result, err := service.BuildReceiptCsv([]models.Receipt{receipt})
	if err != nil {
		t.Fatal(err)
	}

	records, err := csv.NewReader(bytes.NewReader(result.ReceiptCsvBytes)).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	row := records[1]
	if row[1] != "2026-09-30" || row[2] != "2026-09-01" || row[9] != "2026-11-01" {
		t.Errorf("added/date/resolved = %s/%s/%s, want 2026-09-30/2026-09-01/2026-11-01", row[1], row[2], row[9])
	}

	itemRecords, err := csv.NewReader(bytes.NewReader(result.ReceiptItemCsvBytes)).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if itemRecords[1][3] != "2026-09-01" {
		t.Errorf("item receipt date = %s, want 2026-09-01", itemRecords[1][3])
	}
}
