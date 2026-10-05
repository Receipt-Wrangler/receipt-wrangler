package services

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/reporting/receiptsource"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
)

// createSplitReceipt inserts a receipt with the given amount, categories, tags
// and custom field values. Categories and tags must be loaded rows, so the
// association write does not blank their names.
func createSplitReceipt(
	t *testing.T,
	name string,
	amount string,
	paidByUserId uint,
	groupId uint,
	categories []models.Category,
	tags []models.Tag,
	customFields []models.CustomFieldValue,
) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.RequireFromString(amount),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		PaidByUserID: paidByUserId,
		GroupId:      groupId,
		Status:       models.OPEN,
		Categories:   categories,
		Tags:         tags,
		CustomFields: customFields,
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("create receipt %q: %v", name, err)
	}
	return receipt
}

func loadTag(t *testing.T, id uint) models.Tag {
	t.Helper()
	var tag models.Tag
	if err := repositories.GetDB().First(&tag, id).Error; err != nil {
		t.Fatalf("load tag %d: %v", id, err)
	}
	return tag
}

func makeCurrencyField(t *testing.T, name string) uint {
	t.Helper()
	field := models.CustomField{Name: name, Type: models.CURRENCY}
	if err := repositories.GetDB().Create(&field).Error; err != nil {
		t.Fatalf("seed custom field: %v", err)
	}
	return field.ID
}

func currency(literal string) *decimal.Decimal {
	value := decimal.RequireFromString(literal)
	return &value
}

// rowSummary renders each row as "categories|tags|amount" for compact asserts.
func rowSummary(rows []reporting.Row) []string {
	texts := func(values []reporting.Value) string {
		parts := make([]string, 0, len(values))
		for _, value := range values {
			text, _ := value.Text()
			parts = append(parts, text)
		}
		return strings.Join(parts, ",")
	}

	result := make([]string, len(rows))
	for index, row := range rows {
		amount, _ := row.Measure(receiptsource.KeyAmount).Decimal()
		result[index] = texts(row.Get(receiptsource.KeyCategory)) + "|" + texts(row.Get(receiptsource.KeyTag)) + "|" + amount.StringFixed(2)
	}
	return result
}

func measureFixed(t *testing.T, row reporting.Row, key reporting.FieldKey) string {
	t.Helper()
	value, ok := row.Measure(key).Decimal()
	if !ok {
		t.Fatalf("row has no numeric %s", key)
	}
	return value.StringFixed(2)
}

// ---- ReportDataService ------------------------------------------------------

func TestReportDataService_SplitsAcrossCategories(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	toys := loadCategory(t, makeCategory(t, "Toys"))
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "split-cat", nil, nil)
	// Attached out of id order: the leftover cent still goes to the lowest id.
	createSplitReceipt(t, "three", "100", userId, groupId, []models.Category{toys, food, fuel}, nil, nil)

	rowSet, err := NewReportDataService(nil).RowsWithSplit(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{},
		receiptsource.SplitOptions{Categories: true})
	if err != nil {
		t.Fatalf("RowsWithSplit: %v", err)
	}

	assertSummary(t, rowSet.Rows, []string{"Food||33.34", "Fuel||33.33", "Toys||33.33"})
	if !reflect.DeepEqual(rowSet.RowsPerReceipt, []int{3}) {
		t.Errorf("RowsPerReceipt = %v, want [3]", rowSet.RowsPerReceipt)
	}
}

func assertSummary(t *testing.T, rows []reporting.Row, want []string) {
	t.Helper()
	if got := rowSummary(rows); !reflect.DeepEqual(got, want) {
		t.Errorf("rows = %v, want %v", got, want)
	}
}

// Rows (no split) is exactly what it was: one row per receipt, full amount.
func TestReportDataService_RowsWithoutSplitIsUnchanged(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "split-none", nil, nil)
	createSplitReceipt(t, "two", "100", userId, groupId, []models.Category{food, fuel}, nil, nil)

	_, rows, err := NewReportDataService(nil).Rows(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{})
	if err != nil {
		t.Fatalf("Rows: %v", err)
	}
	assertSummary(t, rows, []string{"Food,Fuel||100.00"})
}

// A restricted viewer's split divides by the TRUE number of categories: the
// visible category gets the same share an unrestricted viewer sees, and
// (Restricted) gets exactly the hidden categories' shares.
func TestReportDataService_SplitUsesTheTrueCategoryCount(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	visibleId := makeCategory(t, "Visible")
	hiddenOne := makeCategory(t, "HiddenOne")
	hiddenTwo := makeCategory(t, "HiddenTwo")
	restrictedUser, groupId, _ := seedMemberWithGroupRoleGrants(t, "split-restricted", []uint{visibleId}, nil)

	unrestrictedRole, err := repositories.NewRoleRepository(nil).CreateGroupRole(
		"Unrestricted split role", "", []string{permissions.GroupReceiptsRead}, nil, nil, nil, false, false)
	if err != nil {
		t.Fatalf("seed role: %v", err)
	}
	unrestrictedUser := makeUser(t, "split-unrestricted")
	if err := repositories.GetDB().Create(&models.GroupMember{GroupID: groupId, UserID: unrestrictedUser, GroupRoleID: &unrestrictedRole.ID}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	createSplitReceipt(t, "mixed", "100", restrictedUser, groupId,
		[]models.Category{loadCategory(t, visibleId), loadCategory(t, hiddenOne), loadCategory(t, hiddenTwo)}, nil, nil)

	options := receiptsource.SplitOptions{Categories: true}
	restricted, err := NewReportDataService(nil).RowsWithSplit(restrictedUser, groupIdString(groupId), commands.ReceiptPagedRequestFilter{}, options)
	if err != nil {
		t.Fatalf("RowsWithSplit (restricted): %v", err)
	}
	unrestricted, err := NewReportDataService(nil).RowsWithSplit(unrestrictedUser, groupIdString(groupId), commands.ReceiptPagedRequestFilter{}, options)
	if err != nil {
		t.Fatalf("RowsWithSplit (unrestricted): %v", err)
	}

	// The two hidden copies merged back into one (Restricted) row.
	assertSummary(t, restricted.Rows, []string{"Visible||33.34", "(Restricted)||66.66"})
	if !reflect.DeepEqual(restricted.RowsPerReceipt, []int{2}) {
		t.Errorf("restricted RowsPerReceipt = %v, want [2]", restricted.RowsPerReceipt)
	}
	assertSummary(t, unrestricted.Rows, []string{"Visible||33.34", "HiddenOne||33.33", "HiddenTwo||33.33"})
}

// Splitting by category and tag where both categories are hidden: the four
// copies collapse to one (Restricted) row per tag.
func TestReportDataService_SplitMergesHiddenCopiesPerTag(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	grantedId := makeCategory(t, "Granted")
	hiddenOne := loadCategory(t, makeCategory(t, "HiddenOne"))
	hiddenTwo := loadCategory(t, makeCategory(t, "HiddenTwo"))
	alex := loadTag(t, makeTag(t, "Alex"))
	sam := loadTag(t, makeTag(t, "Sam"))
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "split-merge", []uint{grantedId}, nil)

	createSplitReceipt(t, "hidden", "10", userId, groupId, []models.Category{hiddenOne, hiddenTwo}, []models.Tag{alex, sam}, nil)

	rowSet, err := NewReportDataService(nil).RowsWithSplit(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{},
		receiptsource.SplitOptions{Categories: true, Tags: true})
	if err != nil {
		t.Fatalf("RowsWithSplit: %v", err)
	}

	assertSummary(t, rowSet.Rows, []string{"(Restricted)|Alex|5.00", "(Restricted)|Sam|5.00"})
	if !reflect.DeepEqual(rowSet.RowsPerReceipt, []int{2}) {
		t.Errorf("RowsPerReceipt = %v, want [2]", rowSet.RowsPerReceipt)
	}
}

func TestReportDataService_SplitHonoursExcludedCurrencyFields(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	taxId := makeCurrencyField(t, "Tax")
	tipId := makeCurrencyField(t, "Tip")
	userId, groupId, _ := seedMemberWithGroupRoleGrants(t, "split-excluded", nil, nil)
	createSplitReceipt(t, "fields", "100", userId, groupId, []models.Category{food, fuel}, nil, []models.CustomFieldValue{
		{CustomFieldId: taxId, CurrencyValue: currency("13")},
		{CustomFieldId: tipId, CurrencyValue: currency("10")},
	})

	rowSet, err := NewReportDataService(nil).RowsWithSplit(userId, groupIdString(groupId), commands.ReceiptPagedRequestFilter{},
		receiptsource.SplitOptions{Categories: true, Excluded: map[uint]struct{}{taxId: {}}})
	if err != nil {
		t.Fatalf("RowsWithSplit: %v", err)
	}

	if len(rowSet.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rowSet.Rows))
	}
	for index, row := range rowSet.Rows {
		if got := measureFixed(t, row, receiptsource.CustomFieldKey(taxId)); got != "13.00" {
			t.Errorf("row %d excluded Tax = %s, want 13.00", index, got)
		}
		if got := measureFixed(t, row, receiptsource.CustomFieldKey(tipId)); got != "5.00" {
			t.Errorf("row %d split Tip = %s, want 5.00", index, got)
		}
	}
}

// ---- reportSplitOptions / capRowsToReceipts ---------------------------------

func TestReportService_ReportSplitOptions(t *testing.T) {
	base := func() commands.ReportRequestCommand {
		return commands.ReportRequestCommand{
			GroupBy: []string{"group"},
			Detail:  commands.ReportDetail{Mode: commands.ReportDetailAggregate, By: "category"},
		}
	}

	tests := []struct {
		name         string
		mutate       func(*commands.ReportRequestCommand)
		wantCategory bool
		wantTag      bool
		wantExcluded []uint
	}{
		{"flags off", func(c *commands.ReportRequestCommand) {}, false, false, nil},
		{"category as aggregate dimension", func(c *commands.ReportRequestCommand) {
			c.SplitCategoriesEqually = true
		}, true, false, []uint{}},
		{"tag as a grouping level", func(c *commands.ReportRequestCommand) {
			c.GroupBy = []string{"group", "tag"}
			c.SplitTagsEqually = true
		}, false, true, []uint{}},
		{"category flag on a report that never cuts by category", func(c *commands.ReportRequestCommand) {
			c.Detail.By = "tag"
			c.SplitCategoriesEqually = true
		}, false, false, nil},
		{"records mode with a stray by", func(c *commands.ReportRequestCommand) {
			c.Detail = commands.ReportDetail{Mode: commands.ReportDetailRecords, By: "category"}
			c.SplitCategoriesEqually = true
		}, false, false, nil},
		{"records mode grouped by category", func(c *commands.ReportRequestCommand) {
			c.GroupBy = []string{"category"}
			c.Detail = commands.ReportDetail{Mode: commands.ReportDetailRecords}
			c.SplitCategoriesEqually = true
		}, true, false, []uint{}},
		{"both dimensions", func(c *commands.ReportRequestCommand) {
			c.GroupBy = []string{"tag"}
			c.SplitCategoriesEqually = true
			c.SplitTagsEqually = true
		}, true, true, []uint{}},
		{"excluded keys are parsed and malformed ones dropped", func(c *commands.ReportRequestCommand) {
			c.SplitCategoriesEqually = true
			c.SplitExcludedFields = []string{"custom_4", "custom_4", "amount", "custom_9_month", "custom_7"}
		}, true, false, []uint{4, 7}},
		{"a custom field the report groups by is never split", func(c *commands.ReportRequestCommand) {
			c.GroupBy = []string{"custom_12"}
			c.SplitCategoriesEqually = true
		}, true, false, []uint{12}},
		{"a custom field the report aggregates by is never split", func(c *commands.ReportRequestCommand) {
			c.GroupBy = []string{"category"}
			c.Detail.By = "custom_5"
			c.SplitCategoriesEqually = true
		}, true, false, []uint{5}},
		{"exclusions are ignored while nothing splits", func(c *commands.ReportRequestCommand) {
			c.SplitExcludedFields = []string{"custom_4"}
		}, false, false, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := base()
			test.mutate(&command)
			options := reportSplitOptions(command)

			if options.Categories != test.wantCategory || options.Tags != test.wantTag {
				t.Errorf("options = categories:%v tags:%v, want categories:%v tags:%v",
					options.Categories, options.Tags, test.wantCategory, test.wantTag)
			}
			if test.wantExcluded == nil {
				if options.Excluded != nil {
					t.Errorf("Excluded = %v, want nil", options.Excluded)
				}
				return
			}
			want := map[uint]struct{}{}
			for _, id := range test.wantExcluded {
				want[id] = struct{}{}
			}
			if !reflect.DeepEqual(options.Excluded, want) {
				t.Errorf("Excluded = %v, want %v", options.Excluded, want)
			}
		})
	}
}

func TestReportService_CapRowsToReceipts(t *testing.T) {
	rows := make([]reporting.Row, 6)
	for index := range rows {
		rows[index] = reporting.Row{receiptsource.KeyName: {reporting.Str(string(rune('a' + index)))}}
	}
	rowsPerReceipt := []int{1, 3, 2}

	tests := []struct {
		limit int
		want  int
	}{
		{0, 6},  // no cap
		{-1, 6}, // no cap
		{1, 1},
		{2, 4}, // the second receipt's three rows come whole
		{3, 6},
		{10, 6},
	}
	for _, test := range tests {
		if got := len(capRowsToReceipts(rows, rowsPerReceipt, test.limit)); got != test.want {
			t.Errorf("capRowsToReceipts(limit %d) kept %d rows, want %d", test.limit, got, test.want)
		}
	}
}

// ---- ReportService end to end ------------------------------------------------

// splitReportCommand groups by group and aggregates by category: Total, Count
// and the true grand total.
func splitReportCommand(groupIds []uint) commands.ReportRequestCommand {
	command := aggregateReportCommand("Split", groupIds, []string{commands.ReportFormatCsv})
	command.Columns = append(command.Columns,
		commands.ReportColumn{Kind: commands.ReportColumnAggregate, Name: "Count", Label: "Count", AggFunc: "COUNT"})
	command.SplitCategoriesEqually = true
	return command
}

func cellFixed(t *testing.T, cells []reporting.Cell, column string) string {
	t.Helper()
	for _, cell := range cells {
		if cell.Column == column {
			number, _ := cell.Value().Decimal()
			return number.StringFixed(2)
		}
	}
	t.Fatalf("no cell for column %s", column)
	return ""
}

func bucketFigures(t *testing.T, rows []reporting.DetailRow) []string {
	t.Helper()
	result := make([]string, len(rows))
	for index, row := range rows {
		label := ""
		for _, cell := range row.Cells {
			if cell.Column == "Category" {
				label = cell.Value().String()
			}
		}
		result[index] = label + "=" + cellFixed(t, row.Cells, "Total") + "/" + cellFixed(t, row.Cells, "Count")
	}
	return result
}

func TestReportService_BuildModel_SplitsCategoryBuckets(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	toys := loadCategory(t, makeCategory(t, "Toys"))
	userId, groupIds := seedReportUserInGroups(t, "split-build", "Household")
	createSplitReceipt(t, "three", "100", userId, groupIds[0], []models.Category{food, fuel, toys}, nil, nil)
	createSplitReceipt(t, "one", "20", userId, groupIds[0], []models.Category{food}, nil, nil)

	build, err := NewReportService(nil).buildModel(userId, splitReportCommand(groupIds), time.Now(), 0)
	if err != nil {
		t.Fatalf("buildModel: %v", err)
	}

	household := build.model.Root.Children[0]
	got := bucketFigures(t, household.DetailRows)
	want := []string{"Food=53.34/2.00", "Fuel=33.33/1.00", "Toys=33.33/1.00"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buckets = %v, want %v", got, want)
	}
	if total := cellFixed(t, build.model.GrandTotals, "Total"); total != "120.00" {
		t.Errorf("grand total = %s, want the true 120.00", total)
	}
	if build.receiptCount != 2 {
		t.Errorf("receiptCount = %d, want 2 receipts, not the 4 rows", build.receiptCount)
	}

	// The same report without the box ticked keeps the full-amount attribution.
	command := splitReportCommand(groupIds)
	command.SplitCategoriesEqually = false
	whole, err := NewReportService(nil).buildModel(userId, command, time.Now(), 0)
	if err != nil {
		t.Fatalf("buildModel (unsplit): %v", err)
	}
	if total := cellFixed(t, whole.model.GrandTotals, "Total"); total != "320.00" {
		t.Errorf("unsplit grand total = %s, want 320.00", total)
	}
}

func TestReportService_BuildModel_SplitsRecordsGroupedByCategory(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	userId, groupIds := seedReportUserInGroups(t, "split-records", "Household")
	createSplitReceipt(t, "shared", "9.99", userId, groupIds[0], []models.Category{food, fuel}, nil, nil)

	command := splitReportCommand(groupIds)
	command.GroupBy = []string{"category"}
	command.Detail = commands.ReportDetail{Mode: commands.ReportDetailRecords}
	command.Columns = []commands.ReportColumn{
		{Kind: commands.ReportColumnDimension, Name: "Name", Label: "Name", Field: "name"},
		{Kind: commands.ReportColumnAggregate, Name: "Total", Label: "Total", AggFunc: "SUM", Measure: "amount"},
	}

	build, err := NewReportService(nil).buildModel(userId, command, time.Now(), 0)
	if err != nil {
		t.Fatalf("buildModel: %v", err)
	}

	got := []string{}
	for _, group := range build.model.Root.Children {
		for _, row := range group.DetailRows {
			got = append(got, group.Value.String()+"="+cellFixed(t, row.Cells, "Total"))
		}
	}
	if want := []string{"Food=5.00", "Fuel=4.99"}; !reflect.DeepEqual(got, want) {
		t.Errorf("records = %v, want %v", got, want)
	}
}

// Ticking a box for a dimension the report never cuts by changes nothing.
func TestReportService_BuildModel_UnattributedSplitIsANoOp(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	alex := loadTag(t, makeTag(t, "Alex"))
	sam := loadTag(t, makeTag(t, "Sam"))
	userId, groupIds := seedReportUserInGroups(t, "split-noop", "Household")
	createSplitReceipt(t, "r", "100", userId, groupIds[0], []models.Category{food, fuel}, []models.Tag{alex, sam}, nil)

	// Aggregated by tag only; the category box is ticked, the tag box is not.
	command := splitReportCommand(groupIds)
	command.Detail.By = "tag"
	command.Columns[0].Field = "tag"

	build, err := NewReportService(nil).buildModel(userId, command, time.Now(), 0)
	if err != nil {
		t.Fatalf("buildModel: %v", err)
	}
	if total := cellFixed(t, build.model.GrandTotals, "Total"); total != "200.00" {
		t.Errorf("grand total = %s, want 200.00 (each tag in full)", total)
	}
	if count := cellFixed(t, build.model.GrandTotals, "Count"); count != "2.00" {
		t.Errorf("grand count = %s, want 2.00", count)
	}

	// Now split by tag: the same receipt halves across Alex and Sam.
	command.SplitTagsEqually = true
	build, err = NewReportService(nil).buildModel(userId, command, time.Now(), 0)
	if err != nil {
		t.Fatalf("buildModel (tags): %v", err)
	}
	if total := cellFixed(t, build.model.GrandTotals, "Total"); total != "100.00" {
		t.Errorf("tag-split grand total = %s, want 100.00", total)
	}
}

// The preview cap counts receipts, so a split receipt is never cut part-way:
// all three of its shares are in the one-receipt sample.
func TestReportService_BuildModel_CapKeepsWholeSplitReceipts(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	toys := loadCategory(t, makeCategory(t, "Toys"))
	userId, groupIds := seedReportUserInGroups(t, "split-cap", "Household")
	createSplitReceipt(t, "newest", "90", userId, groupIds[0], []models.Category{food, fuel, toys}, nil, nil)
	older := createSplitReceipt(t, "older", "50", userId, groupIds[0], []models.Category{food}, nil, nil)
	// Receipts are fetched newest first, so the three-way split is the sample.
	if err := repositories.GetDB().Model(&older).Update("date", time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatalf("backdate receipt: %v", err)
	}

	build, err := NewReportService(nil).buildModel(userId, splitReportCommand(groupIds), time.Now(), 1)
	if err != nil {
		t.Fatalf("buildModel: %v", err)
	}
	if build.receiptCount != 2 {
		t.Errorf("receiptCount = %d, want the true 2", build.receiptCount)
	}

	got := bucketFigures(t, build.model.Root.Children[0].DetailRows)
	if want := []string{"Food=30.00/1.00", "Fuel=30.00/1.00", "Toys=30.00/1.00"}; !reflect.DeepEqual(got, want) {
		t.Errorf("sampled buckets = %v, want all three shares of the newest receipt %v", got, want)
	}
}

func TestReportService_Generate_SplitAcrossGroupsCsv(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	userId, groupIds := seedReportUserInGroups(t, "split-csv", "Household", "Roommates")
	createSplitReceipt(t, "household", "10", userId, groupIds[0], []models.Category{food, fuel}, nil, nil)
	createSplitReceipt(t, "roommates", "30", userId, groupIds[1], []models.Category{food, fuel}, nil, nil)

	report, err := NewReportService(nil).Generate(userId, splitReportCommand(groupIds))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	content := string(report.Bytes)
	for _, want := range []string{"5.00", "15.00", "40.00"} {
		if !strings.Contains(content, want) {
			t.Errorf("CSV is missing %q:\n%s", want, content)
		}
	}
	if strings.Contains(content, "80.00") {
		t.Errorf("CSV still double counts (80.00 grand total):\n%s", content)
	}
}

// A saved template carries its split options, and the dashboard render applies
// them over the full dataset.
func TestReportService_RenderTemplateForUser_AppliesSplit(t *testing.T) {
	defer repositories.TruncateTestDb()
	resetAuthzCaches()

	food := loadCategory(t, makeCategory(t, "Food"))
	fuel := loadCategory(t, makeCategory(t, "Fuel"))
	userId := seedAppUser(t, "widget-split", []string{permissions.AppReportsRead, permissions.AppReportsGenerate})
	groupId, _ := joinGroup(t, userId, "Household", []string{permissions.GroupReportsRead, permissions.GroupReceiptsRead})
	createSplitReceipt(t, "shared", "100", userId, groupId, []models.Category{food, fuel}, nil, nil)

	templateId := seedRenderableTemplate(t, "Widget", []uint{groupId}, splitReportCommand([]uint{groupId}))
	preview, err := NewReportService(nil).RenderTemplateForUser(userId, utils.UintToString(templateId))
	if err != nil {
		t.Fatalf("RenderTemplateForUser: %v", err)
	}

	if preview.ReceiptCount != 1 {
		t.Errorf("receipt count = %d, want 1", preview.ReceiptCount)
	}
	if !strings.Contains(preview.Html, "50.00") || strings.Contains(preview.Html, "200.00") {
		t.Errorf("expected split 50.00 buckets and no 200.00 double count, got:\n%s", preview.Html)
	}
}
