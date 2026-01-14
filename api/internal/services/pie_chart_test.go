package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func setupPieChartTest() {
	repositories.CreateTestGroupWithUsers()
	repositories.CreateTestCategories()
}

func tearDownPieChartTest() {
	repositories.TruncateTestDb()
}

func createTestReceipt(name string, amount float64, paidByUserId uint, groupId uint, categories []models.Category, tags []models.Tag) models.Receipt {
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

func createTestTag(name string) models.Tag {
	db := repositories.GetDB()
	tag := models.Tag{
		Name:        name,
		Description: "Test tag",
	}
	db.Create(&tag)
	return tag
}

func TestPieChartService_GetPieChartData_GroupByCategories_Success(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	// Get test categories
	var category1, category2 models.Category
	db.First(&category1, 1)
	db.First(&category2, 2)

	// Create receipts with categories
	createTestReceipt("Receipt 1", 100.00, 1, 1, []models.Category{category1}, nil)
	createTestReceipt("Receipt 2", 50.00, 1, 1, []models.Category{category1}, nil)
	createTestReceipt("Receipt 3", 75.00, 1, 1, []models.Category{category2}, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Verify total amounts by category
	totalAmount := 0.0
	for _, dp := range result.Data {
		totalAmount += dp.Value
	}

	expectedTotal := 225.0 // 100 + 50 + 75
	if totalAmount != expectedTotal {
		utils.PrintTestError(t, totalAmount, expectedTotal)
	}
}

func TestPieChartService_GetPieChartData_GroupByCategories_Uncategorized(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Create receipt without categories
	createTestReceipt("Uncategorized Receipt", 100.00, 1, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Label != "Uncategorized" {
		utils.PrintTestError(t, result.Data[0].Label, "Uncategorized")
		return
	}

	if result.Data[0].Value != 100.0 {
		utils.PrintTestError(t, result.Data[0].Value, 100.0)
	}
}

func TestPieChartService_GetPieChartData_GroupByCategories_MixedCategorizedAndUncategorized(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	// Get test category
	var category1 models.Category
	db.First(&category1, 1)

	// Create receipts - some with categories, some without
	createTestReceipt("Categorized Receipt", 100.00, 1, 1, []models.Category{category1}, nil)
	createTestReceipt("Uncategorized Receipt", 50.00, 1, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Verify we have both categorized and uncategorized entries
	hasUncategorized := false
	hasCategorized := false
	for _, dp := range result.Data {
		if dp.Label == "Uncategorized" {
			hasUncategorized = true
			if dp.Value != 50.0 {
				utils.PrintTestError(t, dp.Value, 50.0)
			}
		}
		if dp.Label == "test" {
			hasCategorized = true
			if dp.Value != 100.0 {
				utils.PrintTestError(t, dp.Value, 100.0)
			}
		}
	}

	if !hasUncategorized {
		utils.PrintTestError(t, "missing Uncategorized entry", "should have Uncategorized")
	}
	if !hasCategorized {
		utils.PrintTestError(t, "missing categorized entry", "should have test category")
	}
}

func TestPieChartService_GetPieChartData_GroupByTags_Success(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Create tags
	tag1 := createTestTag("Groceries")
	tag2 := createTestTag("Entertainment")

	// Create receipts with tags
	createTestReceipt("Receipt 1", 100.00, 1, 1, nil, []models.Tag{tag1})
	createTestReceipt("Receipt 2", 50.00, 1, 1, nil, []models.Tag{tag1})
	createTestReceipt("Receipt 3", 75.00, 1, 1, nil, []models.Tag{tag2})

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_TAGS,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Verify total amounts
	totalAmount := 0.0
	for _, dp := range result.Data {
		totalAmount += dp.Value
	}

	expectedTotal := 225.0
	if totalAmount != expectedTotal {
		utils.PrintTestError(t, totalAmount, expectedTotal)
	}
}

func TestPieChartService_GetPieChartData_GroupByTags_Untagged(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Create receipt without tags
	createTestReceipt("Untagged Receipt", 100.00, 1, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_TAGS,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Label != "Untagged" {
		utils.PrintTestError(t, result.Data[0].Label, "Untagged")
		return
	}

	if result.Data[0].Value != 100.0 {
		utils.PrintTestError(t, result.Data[0].Value, 100.0)
	}
}

func TestPieChartService_GetPieChartData_GroupByTags_MixedTaggedAndUntagged(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Create tag
	tag1 := createTestTag("Test Tag")

	// Create receipts - some with tags, some without
	createTestReceipt("Tagged Receipt", 100.00, 1, 1, nil, []models.Tag{tag1})
	createTestReceipt("Untagged Receipt", 50.00, 1, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_TAGS,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Verify we have both tagged and untagged entries
	hasUntagged := false
	hasTagged := false
	for _, dp := range result.Data {
		if dp.Label == "Untagged" {
			hasUntagged = true
			if dp.Value != 50.0 {
				utils.PrintTestError(t, dp.Value, 50.0)
			}
		}
		if dp.Label == "Test Tag" {
			hasTagged = true
			if dp.Value != 100.0 {
				utils.PrintTestError(t, dp.Value, 100.0)
			}
		}
	}

	if !hasUntagged {
		utils.PrintTestError(t, "missing Untagged entry", "should have Untagged")
	}
	if !hasTagged {
		utils.PrintTestError(t, "missing tagged entry", "should have Test Tag")
	}
}

func TestPieChartService_GetPieChartData_GroupByPaidBy_Success(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Create receipts paid by different users (users 1 and 2 are created by setupPieChartTest)
	createTestReceipt("Receipt 1", 100.00, 1, 1, nil, nil)
	createTestReceipt("Receipt 2", 50.00, 1, 1, nil, nil)
	createTestReceipt("Receipt 3", 75.00, 2, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_PAIDBY,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Verify total amounts
	totalAmount := 0.0
	for _, dp := range result.Data {
		totalAmount += dp.Value
	}

	expectedTotal := 225.0
	if totalAmount != expectedTotal {
		utils.PrintTestError(t, totalAmount, expectedTotal)
	}
}

func TestPieChartService_GetPieChartData_GroupByPaidBy_UserWithDisplayName(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// The test users have DisplayName set to "asdf"
	createTestReceipt("Receipt 1", 100.00, 1, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_PAIDBY,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	// Should use display name "asdf" instead of username "test"
	if result.Data[0].Label != "asdf" {
		utils.PrintTestError(t, result.Data[0].Label, "asdf")
	}
}

func TestPieChartService_GetPieChartData_GroupByPaidBy_UserWithoutDisplayName(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	// Create a user without display name
	user := models.User{
		Username:    "nodisplayname",
		DisplayName: "",
		Password:    "test",
	}
	db.Create(&user)

	// Add user to group
	db.Exec("INSERT INTO group_members (user_id, group_id, group_role) VALUES (?, 1, 'VIEWER')", user.ID)

	createTestReceipt("Receipt 1", 100.00, user.ID, 1, nil, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_PAIDBY,
	}

	result, err := service.GetPieChartData(user.ID, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	// Should use username when display name is empty
	if result.Data[0].Label != "nodisplayname" {
		utils.PrintTestError(t, result.Data[0].Label, "nodisplayname")
	}
}

func TestPieChartService_GetPieChartData_GroupByPaidBy_UnknownUser(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	// Create receipt with a non-existent user ID directly in DB
	receipt := models.Receipt{
		Name:         "Receipt with unknown user",
		Amount:       decimal.NewFromFloat(100.00),
		Date:         time.Now(),
		PaidByUserID: 99999, // Non-existent user
		GroupId:      1,
		Status:       models.OPEN,
	}
	db.Create(&receipt)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_PAIDBY,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	// Should fall back to "Unknown User"
	if result.Data[0].Label != "Unknown User" {
		utils.PrintTestError(t, result.Data[0].Label, "Unknown User")
	}
}

func TestPieChartService_GetPieChartData_EmptyReceipts(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Don't create any receipts

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 0 {
		utils.PrintTestError(t, len(result.Data), 0)
	}
}

func TestPieChartService_GetPieChartData_ReceiptWithMultipleCategories(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	// Get test categories
	var category1, category2 models.Category
	db.First(&category1, 1)
	db.First(&category2, 2)

	// Create receipt with multiple categories
	// The amount should be added to each category
	createTestReceipt("Multi-category Receipt", 100.00, 1, 1, []models.Category{category1, category2}, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Each category should have the full amount
	for _, dp := range result.Data {
		if dp.Value != 100.0 {
			utils.PrintTestError(t, dp.Value, 100.0)
		}
	}
}

func TestPieChartService_GetPieChartData_ReceiptWithMultipleTags(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	// Create tags
	tag1 := createTestTag("Tag1")
	tag2 := createTestTag("Tag2")

	// Create receipt with multiple tags
	createTestReceipt("Multi-tag Receipt", 100.00, 1, 1, nil, []models.Tag{tag1, tag2})

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_TAGS,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 2 {
		utils.PrintTestError(t, len(result.Data), 2)
		return
	}

	// Each tag should have the full amount
	for _, dp := range result.Data {
		if dp.Value != 100.0 {
			utils.PrintTestError(t, dp.Value, 100.0)
		}
	}
}

func TestPieChartService_GetPieChartData_WithTransaction(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()
	tx := db.Begin()
	defer tx.Rollback()

	// Create test data within the transaction
	category := models.Category{Name: "TX Category"}
	tx.Create(&category)

	receipt := models.Receipt{
		Name:         "Transaction Receipt",
		Amount:       decimal.NewFromFloat(100.00),
		Date:         time.Now(),
		PaidByUserID: 1,
		GroupId:      1,
		Status:       models.OPEN,
		Categories:   []models.Category{category},
	}
	tx.Create(&receipt)

	// Create service with transaction
	service := NewPieChartService(tx)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// Should see the receipt created within the transaction
	if len(result.Data) == 0 {
		utils.PrintTestError(t, "no data returned", "expected data from transaction")
	}
}

func TestPieChartService_GetPieChartData_WithFilter(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	var category1 models.Category
	db.First(&category1, 1)

	// Create receipts with different statuses
	r1 := createTestReceipt("Open Receipt", 100.00, 1, 1, []models.Category{category1}, nil)
	r2 := createTestReceipt("Resolved Receipt", 50.00, 1, 1, []models.Category{category1}, nil)

	// Update one receipt to RESOLVED status
	db.Model(&r1).Update("status", models.OPEN)
	db.Model(&r2).Update("status", models.RESOLVED)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
		Filter: commands.ReceiptPagedRequestFilter{
			Status: commands.PagedRequestField{
				Value:     []interface{}{"OPEN"},
				Operation: commands.EQUALS,
			},
		},
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// Should only see the OPEN receipt
	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Value != 100.0 {
		utils.PrintTestError(t, result.Data[0].Value, 100.0)
	}
}

func TestPieChartService_GetPieChartData_LargeAmounts(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	var category1 models.Category
	db.First(&category1, 1)

	// Create receipts with large amounts
	createTestReceipt("Large Receipt 1", 999999.99, 1, 1, []models.Category{category1}, nil)
	createTestReceipt("Large Receipt 2", 888888.88, 1, 1, []models.Category{category1}, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	expectedTotal := 999999.99 + 888888.88
	if result.Data[0].Value != expectedTotal {
		utils.PrintTestError(t, result.Data[0].Value, expectedTotal)
	}
}

func TestPieChartService_GetPieChartData_DecimalAmounts(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	var category1 models.Category
	db.First(&category1, 1)

	// Create receipts with decimal amounts
	createTestReceipt("Decimal Receipt 1", 10.33, 1, 1, []models.Category{category1}, nil)
	createTestReceipt("Decimal Receipt 2", 20.67, 1, 1, []models.Category{category1}, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	expectedTotal := 31.0 // 10.33 + 20.67
	if result.Data[0].Value != expectedTotal {
		utils.PrintTestError(t, result.Data[0].Value, expectedTotal)
	}
}

func TestPieChartService_NewPieChartService_WithNilTx(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	service := NewPieChartService(nil)

	// Service should be created successfully
	if service.DB == nil {
		utils.PrintTestError(t, "DB should not be nil", "DB should be set from GetDB()")
	}

	if service.TX != nil {
		utils.PrintTestError(t, "TX should be nil", "TX should be nil when created with nil")
	}
}

func TestPieChartService_NewPieChartService_WithTx(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()
	tx := db.Begin()
	defer tx.Rollback()

	service := NewPieChartService(tx)

	// Service should have transaction set
	if service.TX == nil {
		utils.PrintTestError(t, "TX should not be nil", "TX should be set")
	}
}

func TestPieChartService_GetPieChartData_ZeroAmounts(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	var category1 models.Category
	db.First(&category1, 1)

	// Create receipt with zero amount
	createTestReceipt("Zero Amount Receipt", 0.00, 1, 1, []models.Category{category1}, nil)

	service := NewPieChartService(nil)
	command := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	result, err := service.GetPieChartData(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(result.Data) != 1 {
		utils.PrintTestError(t, len(result.Data), 1)
		return
	}

	if result.Data[0].Value != 0.0 {
		utils.PrintTestError(t, result.Data[0].Value, 0.0)
	}
}

func TestPieChartService_GetPieChartData_SingleReceiptAllGroupings(t *testing.T) {
	defer tearDownPieChartTest()
	setupPieChartTest()

	db := repositories.GetDB()

	var category1 models.Category
	db.First(&category1, 1)

	tag1 := createTestTag("SingleTestTag")

	// Create a single receipt with category and tag
	createTestReceipt("Single Receipt", 150.00, 1, 1, []models.Category{category1}, []models.Tag{tag1})

	service := NewPieChartService(nil)

	// Test CATEGORIES grouping
	catCommand := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_CATEGORIES,
	}

	catResult, err := service.GetPieChartData(1, "1", catCommand)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(catResult.Data) != 1 || catResult.Data[0].Value != 150.0 {
		utils.PrintTestError(t, catResult.Data, "expected single entry with value 150.0")
	}

	// Test TAGS grouping
	tagCommand := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_TAGS,
	}

	tagResult, err := service.GetPieChartData(1, "1", tagCommand)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(tagResult.Data) != 1 || tagResult.Data[0].Value != 150.0 {
		utils.PrintTestError(t, tagResult.Data, "expected single entry with value 150.0")
	}

	// Test PAIDBY grouping
	paidByCommand := commands.PieChartDataCommand{
		ChartGrouping: models.CHART_GROUPING_PAIDBY,
	}

	paidByResult, err := service.GetPieChartData(1, "1", paidByCommand)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if len(paidByResult.Data) != 1 || paidByResult.Data[0].Value != 150.0 {
		utils.PrintTestError(t, paidByResult.Data, "expected single entry with value 150.0")
	}
}
