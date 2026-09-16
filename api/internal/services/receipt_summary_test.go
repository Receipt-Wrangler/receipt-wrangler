package services

import (
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func setupReceiptSummaryTest() {
	repositories.CreateTestGroupWithUsers()
}

func tearDownReceiptSummaryTest() {
	repositories.TruncateTestDb()
}

// configureSummary writes a group's summary configuration directly. Going through the repository
// (rather than hand-inserting join rows) keeps these tests honest about the shape the handler
// actually persists.
func configureSummary(
	t *testing.T,
	groupId uint,
	enabled bool,
	statuses []models.ReceiptStatus,
	customFieldIds []uint,
) {
	t.Helper()

	repository := repositories.NewGroupReceiptSettingsRepository(nil)
	if _, err := repository.CreateGroupReceiptSettings(groupId); err != nil {
		t.Fatalf("create settings: %v", err)
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled:       true,
		QuickScanPaidByRequired:      true,
		QuickScanStatusEnabled:       true,
		QuickScanStatusRequired:      true,
		QuickScanDefaultPaidByType:   models.QUICK_SCAN_PAID_BY_UPLOADER,
		QuickScanDefaultStatus:       models.OPEN,
		ReceiptSummaryEnabled:        &enabled,
		ReceiptSummaryStatuses:       &statuses,
		ReceiptSummaryCustomFieldIds: &customFieldIds,
	}

	if _, err := repository.UpdateGroupReceiptSettings(utils.UintToString(groupId), command); err != nil {
		t.Fatalf("configure summary: %v", err)
	}
}

func seedSummaryCurrencyField(t *testing.T, name string) uint {
	t.Helper()
	customField := models.CustomField{Name: name, Type: models.CURRENCY}
	if err := repositories.GetDB().Create(&customField).Error; err != nil {
		t.Fatalf("seed currency field: %v", err)
	}
	return customField.ID
}

// seedSummaryReceipt creates a receipt with an optional currency value per custom field.
func seedSummaryReceipt(
	t *testing.T,
	groupId uint,
	status models.ReceiptStatus,
	amount string,
	customFieldValues map[uint]string,
) models.Receipt {
	t.Helper()

	receipt := models.Receipt{
		Name:         "summary receipt",
		Amount:       decimal.RequireFromString(amount),
		Date:         time.Now(),
		PaidByUserID: 1,
		GroupId:      groupId,
		Status:       status,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}

	for customFieldId, value := range customFieldValues {
		currencyValue := decimal.RequireFromString(value)
		customFieldValue := models.CustomFieldValue{
			ReceiptId:     receipt.ID,
			CustomFieldId: customFieldId,
			CurrencyValue: &currencyValue,
		}
		if err := repositories.GetDB().Create(&customFieldValue).Error; err != nil {
			t.Fatalf("seed custom field value: %v", err)
		}
	}

	return receipt
}

// summaryCommand builds a command through the real entry point rather than by hand.
// LoadDataFromRequest is what seeds the filter's non-nil defaults (initReceiptFilterValues is
// unexported), and BuildGormFilterQuery's type assertions plus IntersectReceiptFilterWithGrants
// both assume that seeding - so a hand-built command would exercise a shape the server never
// actually sees.
func summaryCommand(t *testing.T, body string) commands.ReceiptSummaryCommand {
	t.Helper()

	command := commands.ReceiptSummaryCommand{}
	request := httptest.NewRequest(http.MethodPost, "/api/receipt/group/1/summary", strings.NewReader(body))
	if err := command.LoadDataFromRequest(httptest.NewRecorder(), request); err != nil {
		t.Fatalf("load summary command: %v", err)
	}

	return command
}

func emptySummaryCommand(t *testing.T) commands.ReceiptSummaryCommand {
	t.Helper()
	return summaryCommand(t, `{"filter":{}}`)
}

func findStatusRow(statuses []structs.ReceiptSummaryRow, status models.ReceiptStatus) (structs.ReceiptSummaryRow, bool) {
	for _, row := range statuses {
		if row.Status == status {
			return row, true
		}
	}
	return structs.ReceiptSummaryRow{}, false
}

func TestReceiptSummary_DisabledReturnsNothing(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	configureSummary(t, 1, false, []models.ReceiptStatus{models.OPEN}, nil)
	seedSummaryReceipt(t, 1, models.OPEN, "10.00", nil)

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if summary.Enabled {
		utils.PrintTestError(t, summary.Enabled, false)
	}
	// Non-nil slices even in the off state: a null would fail the whole payload on a released
	// Dart client.
	if summary.Statuses == nil || summary.Overall.CustomFieldTotals == nil {
		utils.PrintTestError(t, summary, "non-nil slices")
	}
}

func TestReceiptSummary_MissingSettingsRowIsOff(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	// No settings row at all - they are created lazily, so a group nobody has opened the settings
	// page for legitimately has none. That must read as "off", not as an error.
	seedSummaryReceipt(t, 1, models.OPEN, "10.00", nil)

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}
	if summary.Enabled {
		utils.PrintTestError(t, summary.Enabled, false)
	}
}

func TestReceiptSummary_TotalsAndStatusBreakdown(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN, models.RESOLVED}, nil)

	seedSummaryReceipt(t, 1, models.OPEN, "10.50", nil)
	seedSummaryReceipt(t, 1, models.OPEN, "4.25", nil)
	seedSummaryReceipt(t, 1, models.RESOLVED, "100.00", nil)

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if summary.Overall.ReceiptCount != 3 {
		utils.PrintTestError(t, summary.Overall.ReceiptCount, 3)
	}
	if !summary.Overall.Total.Equal(decimal.RequireFromString("114.75")) {
		utils.PrintTestError(t, summary.Overall.Total.String(), "114.75")
	}

	open, found := findStatusRow(summary.Statuses, models.OPEN)
	if !found || open.ReceiptCount != 2 || !open.Total.Equal(decimal.RequireFromString("14.75")) {
		utils.PrintTestError(t, open, "OPEN: 2 receipts, 14.75")
	}

	resolved, found := findStatusRow(summary.Statuses, models.RESOLVED)
	if !found || resolved.ReceiptCount != 1 || !resolved.Total.Equal(decimal.RequireFromString("100.00")) {
		utils.PrintTestError(t, resolved, "RESOLVED: 1 receipt, 100.00")
	}
}

// TestReceiptSummary_ConfiguredStatusWithNoReceiptsRendersZero is the rule that keeps the block's
// shape steady as a filter narrows: a configured status that matches nothing still gets a row.
func TestReceiptSummary_ConfiguredStatusWithNoReceiptsRendersZero(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	fieldId := seedSummaryCurrencyField(t, "HST")
	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN, models.DECLINED}, []uint{fieldId})
	seedSummaryReceipt(t, 1, models.OPEN, "10.00", map[uint]string{fieldId: "1.30"})

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	declined, found := findStatusRow(summary.Statuses, models.DECLINED)
	if !found {
		utils.PrintTestError(t, summary.Statuses, "a DECLINED row")
		return
	}
	if declined.ReceiptCount != 0 || !declined.Total.IsZero() {
		utils.PrintTestError(t, declined, "DECLINED: 0 receipts, 0")
	}
	// The zero row still carries every configured column, so the layout does not jump.
	if len(declined.CustomFieldTotals) != 1 || !declined.CustomFieldTotals[0].Total.IsZero() {
		utils.PrintTestError(t, declined.CustomFieldTotals, "one zeroed HST column")
	}
}

// TestReceiptSummary_UnconfiguredStatusStillCountsInOverall: a receipt whose status has no
// breakdown row is still in the filter result, so excluding it would make the total disagree with
// the table's own count.
func TestReceiptSummary_UnconfiguredStatusStillCountsInOverall(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN}, nil)
	seedSummaryReceipt(t, 1, models.OPEN, "10.00", nil)
	seedSummaryReceipt(t, 1, models.DRAFT, "25.00", nil)

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if summary.Overall.ReceiptCount != 2 {
		utils.PrintTestError(t, summary.Overall.ReceiptCount, 2)
	}
	if !summary.Overall.Total.Equal(decimal.RequireFromString("35.00")) {
		utils.PrintTestError(t, summary.Overall.Total.String(), "35.00")
	}
	if len(summary.Statuses) != 1 {
		utils.PrintTestError(t, summary.Statuses, "only the configured OPEN row")
	}
}

func TestReceiptSummary_CurrencyCustomFieldTotals(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	hst := seedSummaryCurrencyField(t, "HST")
	subtotal := seedSummaryCurrencyField(t, "Subtotal")
	ignored := seedSummaryCurrencyField(t, "Not configured")

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN}, []uint{hst, subtotal})

	seedSummaryReceipt(t, 1, models.OPEN, "11.30", map[uint]string{hst: "1.30", subtotal: "10.00", ignored: "999.99"})
	// A receipt carrying no value for the configured fields contributes 0 to them, not an error.
	seedSummaryReceipt(t, 1, models.OPEN, "5.00", nil)

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	// Columns follow the configured order, not map order, so they cannot reshuffle between requests.
	if len(summary.Overall.CustomFieldTotals) != 2 {
		utils.PrintTestError(t, summary.Overall.CustomFieldTotals, "two columns")
		return
	}
	if summary.Overall.CustomFieldTotals[0].Name != "HST" ||
		!summary.Overall.CustomFieldTotals[0].Total.Equal(decimal.RequireFromString("1.30")) {
		utils.PrintTestError(t, summary.Overall.CustomFieldTotals[0], "HST = 1.30")
	}
	if summary.Overall.CustomFieldTotals[1].Name != "Subtotal" ||
		!summary.Overall.CustomFieldTotals[1].Total.Equal(decimal.RequireFromString("10.00")) {
		utils.PrintTestError(t, summary.Overall.CustomFieldTotals[1], "Subtotal = 10.00")
	}
}

// TestReceiptSummary_DuplicateCustomFieldValuesLowestIdWins mirrors
// reporting/receiptsource.addCustomFields. custom_field_values has no unique index on
// (receipt_id, custom_field_id) and the association loads without an ORDER BY, so without this
// rule two identical requests could return different numbers.
func TestReceiptSummary_DuplicateCustomFieldValuesLowestIdWins(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	hst := seedSummaryCurrencyField(t, "HST")
	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN}, []uint{hst})

	receipt := seedSummaryReceipt(t, 1, models.OPEN, "10.00", map[uint]string{hst: "1.00"})

	// A second value for the same field on the same receipt, written later so it has a higher id.
	loser := decimal.RequireFromString("99.00")
	if err := repositories.GetDB().Create(&models.CustomFieldValue{
		ReceiptId:     receipt.ID,
		CustomFieldId: hst,
		CurrencyValue: &loser,
	}).Error; err != nil {
		t.Fatalf("seed duplicate value: %v", err)
	}

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if !summary.Overall.CustomFieldTotals[0].Total.Equal(decimal.RequireFromString("1.00")) {
		utils.PrintTestError(t, summary.Overall.CustomFieldTotals[0].Total.String(), "1.00")
	}
}

// TestReceiptSummary_RespectsStatusFilter: the overall row honours the filter as-is, and the
// per-status rows are that set intersected with each status - so a configured status excluded by
// the filter reads 0 rather than showing its unfiltered total.
func TestReceiptSummary_RespectsStatusFilter(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN, models.RESOLVED}, nil)
	seedSummaryReceipt(t, 1, models.OPEN, "10.00", nil)
	seedSummaryReceipt(t, 1, models.RESOLVED, "100.00", nil)

	command := summaryCommand(t, `{"filter":{"status":{"operation":"CONTAINS","value":["OPEN"]}}}`)

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", command)
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if summary.Overall.ReceiptCount != 1 || !summary.Overall.Total.Equal(decimal.RequireFromString("10.00")) {
		utils.PrintTestError(t, summary.Overall, "1 receipt, 10.00")
	}

	resolved, found := findStatusRow(summary.Statuses, models.RESOLVED)
	if !found || resolved.ReceiptCount != 0 || !resolved.Total.IsZero() {
		utils.PrintTestError(t, resolved, "RESOLVED filtered out: 0 receipts, 0")
	}
}

// TestReceiptSummary_DecimalPrecision pins the reason this folds in Go rather than issuing a SQL
// SUM: SQLite has no decimal type, so SUM(amount) there returns an IEEE double. These three values
// are exactly the case that loses a cent in float arithmetic.
func TestReceiptSummary_DecimalPrecision(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN}, nil)
	for _, amount := range []string{"0.10", "0.20", "0.30"} {
		seedSummaryReceipt(t, 1, models.OPEN, amount, nil)
	}

	summary, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", emptySummaryCommand(t))
	if err != nil {
		utils.PrintTestError(t, err, "no error")
		return
	}

	if summary.Overall.Total.String() != "0.6" && summary.Overall.Total.String() != "0.60" {
		utils.PrintTestError(t, summary.Overall.Total.String(), "0.60 exactly")
	}
	if !summary.Overall.Total.Equal(decimal.RequireFromString("0.60")) {
		utils.PrintTestError(t, summary.Overall.Total.String(), "0.60")
	}
}

// seedSummaryAllGroup creates the synthetic All group (the only group allowed to borrow another
// group's configuration) with user 1 as a member. It lives here rather than in
// CreateTestGroupWithUsers so no other suite's group counts shift.
func seedSummaryAllGroup(t *testing.T) uint {
	t.Helper()

	db := repositories.GetDB()
	allGroup := models.Group{Name: "all", IsAllGroup: true}
	if err := db.Create(&allGroup).Error; err != nil {
		t.Fatalf("create all group: %v", err)
	}
	member := models.GroupMember{GroupID: allGroup.ID, UserID: 1}
	if err := db.Model(models.GroupMember{}).Create(&member).Error; err != nil {
		t.Fatalf("add all group member: %v", err)
	}

	return allGroup.ID
}

// TestReceiptSummary_ConfigurationGroupForbidden: without this check a member could name any group
// id and read back its configured custom field names. Viewed through the All group, because that is
// now the only group a configuration override is accepted for at all.
func TestReceiptSummary_ConfigurationGroupForbidden(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	allGroupId := seedSummaryAllGroup(t)

	// Group 2 exists but user 1 is not a member of it (CreateTestGroupWithUsers puts user 4 there).
	configureSummary(t, 2, true, []models.ReceiptStatus{models.OPEN}, nil)

	command := summaryCommand(t, `{"configurationGroupId":2,"filter":{}}`)

	_, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, utils.UintToString(allGroupId), command)
	if err != ErrConfigurationGroupForbidden {
		utils.PrintTestError(t, err, ErrConfigurationGroupForbidden)
	}
}

// TestReceiptSummary_ConfigurationGroupRejectedForRealGroup: a real group must use its own
// configuration. Otherwise a member could render group 1's receipts under group 2's statuses and
// currency fields — overriding what group 1's admin chose, and opting into a summary group 1 has
// switched off. Only the All group, which has no settings row of its own, may borrow one.
func TestReceiptSummary_ConfigurationGroupRejectedForRealGroup(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	// User 1 is made a member of group 2 on purpose, so the borrowed configuration is one they may
	// genuinely read. That is what makes this a test of the request's SHAPE rather than of access:
	// the call would sail through the permission check and still has to be refused.
	db := repositories.GetDB()
	member := models.GroupMember{GroupID: 2, UserID: 1}
	if err := db.Model(models.GroupMember{}).Create(&member).Error; err != nil {
		t.Fatalf("add group 2 member: %v", err)
	}

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN}, nil)
	configureSummary(t, 2, true, []models.ReceiptStatus{models.RESOLVED}, nil)

	command := summaryCommand(t, `{"configurationGroupId":2,"filter":{}}`)

	_, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", command)
	if err != ErrConfigurationGroupNotAllGroup {
		utils.PrintTestError(t, err, ErrConfigurationGroupNotAllGroup)
	}
}

// TestReceiptSummary_ConfigurationGroupRejectedBeforeAccessCheck pins the ORDER of the two guards.
// Naming a group the caller cannot read, from a real group, must still report the not-all-group
// rejection: if the access check ran first, the response would differ depending on whether the
// named group is readable, which turns this endpoint into a probe for other groups' existence.
func TestReceiptSummary_ConfigurationGroupRejectedBeforeAccessCheck(t *testing.T) {
	defer tearDownReceiptSummaryTest()
	setupReceiptSummaryTest()

	configureSummary(t, 1, true, []models.ReceiptStatus{models.OPEN}, nil)

	// Group 2 is one user 1 is NOT a member of, and 9999 does not exist at all. Both must come back
	// as the same rejection, so neither can be distinguished from a readable group.
	for _, configurationGroupId := range []string{"2", "9999"} {
		command := summaryCommand(t, `{"configurationGroupId":`+configurationGroupId+`,"filter":{}}`)

		_, err := NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", command)
		if err != ErrConfigurationGroupNotAllGroup {
			utils.PrintTestError(t, err, ErrConfigurationGroupNotAllGroup)
		}
	}
}
