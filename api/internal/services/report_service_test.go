package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/repositories"
)

// --- pure helpers ---------------------------------------------------------

func TestReportService_SubstituteVariables(t *testing.T) {
	values := map[string]string{
		"period":           "2026-05-01 to 2026-05-31",
		"group.name":       "Household, Roommates",
		"generatedAt":      "Jul 13, 2026, 4:07 PM",
		"currentUser.name": "Noah $Hall", // a literal $ must not be treated as an expansion
	}
	in := "{{period}} | {{ group.name }} | {{generatedAt}} | {{currentUser.name}} | {{unknown}}"
	want := "2026-05-01 to 2026-05-31 | Household, Roommates | Jul 13, 2026, 4:07 PM | Noah $Hall | {{unknown}}"

	if got := substituteVariables(in, values); got != want {
		t.Errorf("substituteVariables:\n got %q\nwant %q", got, want)
	}
}

func TestReportService_ResolvePeriodBounds(t *testing.T) {
	now := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)
	ymd := func(tm time.Time) [3]int { return [3]int{tm.Year(), int(tm.Month()), tm.Day()} }

	tests := []struct {
		preset     string
		start, end [3]int
	}{
		{commands.ReportPeriodThisMonth, [3]int{2026, 5, 1}, [3]int{2026, 5, 31}},
		{commands.ReportPeriodLastMonth, [3]int{2026, 4, 1}, [3]int{2026, 4, 30}},
		{commands.ReportPeriodMtd, [3]int{2026, 5, 1}, [3]int{2026, 5, 15}},
		{commands.ReportPeriodQtd, [3]int{2026, 4, 1}, [3]int{2026, 5, 15}},
		{commands.ReportPeriodYtd, [3]int{2026, 1, 1}, [3]int{2026, 5, 15}},
	}
	for _, test := range tests {
		t.Run(test.preset, func(t *testing.T) {
			start, end := resolvePeriodBounds(commands.ReportPeriod{Preset: test.preset}, now)
			if ymd(start) != test.start || ymd(end) != test.end {
				t.Errorf("%s = %v..%v, want %v..%v", test.preset, ymd(start), ymd(end), test.start, test.end)
			}
		})
	}

	start, end := resolvePeriodBounds(commands.ReportPeriod{
		Preset: commands.ReportPeriodCustom, StartDate: "2026-02-10", EndDate: "2026-03-20",
	}, now)
	if ymd(start) != [3]int{2026, 2, 10} || ymd(end) != [3]int{2026, 3, 20} {
		t.Errorf("custom = %v..%v", ymd(start), ymd(end))
	}
}

func TestReportService_ApplyPeriodWritesBetweenDateFilter(t *testing.T) {
	now := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	filter := commands.ReceiptPagedRequestFilter{}

	label := applyPeriod(&filter, commands.ReportPeriod{Preset: commands.ReportPeriodThisMonth}, now)

	if filter.Date.Operation != commands.BETWEEN {
		t.Errorf("date operation = %q, want BETWEEN", filter.Date.Operation)
	}
	bounds, ok := filter.Date.Value.([]interface{})
	if !ok || len(bounds) != 2 {
		t.Fatalf("date value = %v, want a two-element bound slice", filter.Date.Value)
	}
	if label != "2026-05-01 to 2026-05-31" {
		t.Errorf("label = %q", label)
	}
}

func TestReportService_ReportBaseName(t *testing.T) {
	tests := map[string]string{
		"My Report":         "My_Report",
		"Q2/2026 Expenses!": "Q2_2026_Expenses",
		"   ":               "report",
		"":                  "report",
		"clean_name":        "clean_name",
	}
	for in, want := range tests {
		if got := reportBaseName(in); got != want {
			t.Errorf("reportBaseName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestReportService_AggregateSource(t *testing.T) {
	tests := []struct {
		fn, measure, want string
	}{
		{"COUNT", "", "COUNT()"},
		{"COUNT", "amount", "COUNT()"},
		{"SUM", "amount", "SUM(amount)"},
		{"AVG", "custom_7", "AVG(custom_7)"},
	}
	for _, test := range tests {
		if got := aggregateSource(test.fn, test.measure); got != test.want {
			t.Errorf("aggregateSource(%q,%q) = %q, want %q", test.fn, test.measure, got, test.want)
		}
	}
}

func TestReportService_BuildReportSpec(t *testing.T) {
	command := commands.ReportRequestCommand{
		GroupBy: []string{"group", "category"},
		Detail:  commands.ReportDetail{Mode: commands.ReportDetailAggregate, By: "category"},
		Columns: []commands.ReportColumn{
			{Kind: commands.ReportColumnDimension, Name: "Category", Label: "Category", Field: "category"},
			{Kind: commands.ReportColumnAggregate, Name: "Total", Label: "Total", AggFunc: "SUM", Measure: "amount"},
			{Kind: commands.ReportColumnAggregate, Name: "Count", Label: "Count", AggFunc: "COUNT"},
			{Kind: commands.ReportColumnFormula, Name: "Avg", Label: "Avg", Expr: "Total / Count"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	spec, err := buildReportSpec(command)
	if err != nil {
		t.Fatalf("buildReportSpec: %v", err)
	}

	if len(spec.GroupBy) != 2 || spec.GroupBy[0] != "group" || spec.GroupBy[1] != "category" {
		t.Errorf("GroupBy = %v", spec.GroupBy)
	}
	if spec.Detail.Mode != reporting.DetailAggregate || spec.Detail.By != "category" {
		t.Errorf("Detail = %+v", spec.Detail)
	}
	if !spec.Subtotals || !spec.GrandTotals {
		t.Errorf("totals flags not carried: %+v", spec)
	}

	want := []struct {
		kind   reporting.ColumnKind
		field  reporting.FieldKey
		aggSrc string
		expr   string
	}{
		{reporting.ColumnLabel, "category", "", ""},
		{reporting.ColumnAggregate, "", "SUM(amount)", ""},
		{reporting.ColumnAggregate, "", "COUNT()", ""},
		{reporting.ColumnArithmetic, "", "", "Total / Count"},
	}
	for index, expectation := range want {
		column := spec.Columns[index]
		if column.Kind != expectation.kind || column.Field != expectation.field ||
			column.AggSrc != expectation.aggSrc || column.Expr != expectation.expr {
			t.Errorf("column %d = %+v, want %+v", index, column, expectation)
		}
	}
}

// --- DB-backed generation -------------------------------------------------

// seedReportUserInGroups creates one user who is a member of every named group,
// each membership carrying a role that can read reports and receipts.
func seedReportUserInGroups(t *testing.T, username string, groupNames ...string) (uint, []uint) {
	t.Helper()
	db := repositories.GetDB()
	roleRepository := repositories.NewRoleRepository(nil)
	role, err := roleRepository.CreateGroupRole(
		"Report Role "+username, "",
		[]string{permissions.GroupReportsRead, permissions.GroupReceiptsRead},
		nil, nil, nil, false,
	)
	if err != nil {
		t.Fatalf("seed group role: %v", err)
	}

	user := models.User{Username: username, Password: "password"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	groupIds := make([]uint, 0, len(groupNames))
	for _, name := range groupNames {
		group := models.Group{Name: name}
		if err := db.Create(&group).Error; err != nil {
			t.Fatalf("seed group %q: %v", name, err)
		}
		member := models.GroupMember{GroupID: group.ID, UserID: user.ID, GroupRoleID: &role.ID}
		if err := db.Create(&member).Error; err != nil {
			t.Fatalf("seed membership: %v", err)
		}
		groupIds = append(groupIds, group.ID)
	}
	return user.ID, groupIds
}

// aggregateReportCommand is the report configuration the generation tests share.
func aggregateReportCommand(name string, groupIds []uint, formats []string) commands.ReportRequestCommand {
	ids := make([]string, len(groupIds))
	for index, groupId := range groupIds {
		ids[index] = groupIdString(groupId)
	}
	return commands.ReportRequestCommand{
		Name:     name,
		GroupIds: ids,
		Period:   commands.ReportPeriod{Preset: commands.ReportPeriodCustom, StartDate: "2026-05-01", EndDate: "2026-05-31"},
		GroupBy:  []string{"group"},
		Detail:   commands.ReportDetail{Mode: commands.ReportDetailAggregate, By: "category"},
		Columns: []commands.ReportColumn{
			{Kind: commands.ReportColumnDimension, Name: "Category", Label: "Category", Field: "category"},
			{Kind: commands.ReportColumnAggregate, Name: "Total", Label: "Total", AggFunc: "SUM", Measure: "amount"},
		},
		Subtotals:   true,
		GrandTotals: true,
		Formats:     formats,
	}
}

func TestReportService_Generate_MultiGroupCsv(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	category := loadCategory(t, makeCategory(t, "Groceries"))
	userId, groupIds := seedReportUserInGroups(t, "rpt-gen-csv", "Household", "Roommates")
	createReportReceipt(t, "household-1", userId, groupIds[0], []models.Category{category})
	createReportReceipt(t, "roommates-1", userId, groupIds[1], []models.Category{category})

	report, err := NewReportService(nil).Generate(userId, aggregateReportCommand("Cross Group", groupIds, []string{commands.ReportFormatCsv}))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if report.ContentType != constants.TextCsv {
		t.Errorf("content type = %q, want %q", report.ContentType, constants.TextCsv)
	}
	if report.Filename != "Cross_Group.csv" {
		t.Errorf("filename = %q, want Cross_Group.csv", report.Filename)
	}

	content := string(report.Bytes)
	for _, want := range []string{"Household", "Roommates", "Groceries", "Grand Total"} {
		if !strings.Contains(content, want) {
			t.Errorf("report is missing %q:\n%s", want, content)
		}
	}
}

func TestReportService_Generate_MultiFormatZips(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	category := loadCategory(t, makeCategory(t, "Groceries"))
	userId, groupIds := seedReportUserInGroups(t, "rpt-gen-zip", "Household")
	createReportReceipt(t, "household-1", userId, groupIds[0], []models.Category{category})

	report, err := NewReportService(nil).Generate(userId,
		aggregateReportCommand("Bundle", groupIds, []string{commands.ReportFormatCsv, commands.ReportFormatXlsx}))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if report.ContentType != constants.ApplicationZip {
		t.Errorf("content type = %q, want %q", report.ContentType, constants.ApplicationZip)
	}
	if report.Filename != "Bundle.zip" {
		t.Errorf("filename = %q, want Bundle.zip", report.Filename)
	}

	reader, err := zip.NewReader(bytes.NewReader(report.Bytes), int64(len(report.Bytes)))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	names := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		names = append(names, file.Name)
	}
	for _, want := range []string{"Bundle.csv", "Bundle.xlsx"} {
		if !containsString(names, want) {
			t.Errorf("zip entries %v missing %q", names, want)
		}
	}
}

func TestReportService_Generate_InvalidSpecIsClientError(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	userId, groupIds := seedReportUserInGroups(t, "rpt-gen-bad", "Household")
	createReportReceipt(t, "household-1", userId, groupIds[0], nil)

	command := aggregateReportCommand("Bad", groupIds, []string{commands.ReportFormatCsv})
	command.GroupBy = []string{"amount"} // a measure cannot be grouped by

	_, err := NewReportService(nil).Generate(userId, command)
	if err == nil {
		t.Fatal("expected an error grouping by a measure, got none")
	}
	var specErr *ReportSpecError
	if !errors.As(err, &specErr) {
		t.Errorf("expected a ReportSpecError (→ 400), got %T: %v", err, err)
	}
}

// The PDF format exercises the full pipeline including the render.HTML → chromium
// bridge and the document-variable substitution assembly. The PDF bytes are
// opaque, so this asserts the path runs and produces a document, not its text.
func TestReportService_Generate_PdfDocument(t *testing.T) {
	defer repositories.TruncateTestDb()
	clearGroupRoleGrantCacheAll()
	clearRolePermissionCacheAll()

	category := loadCategory(t, makeCategory(t, "Groceries"))
	userId, groupIds := seedReportUserInGroups(t, "rpt-gen-pdf", "Household")
	createReportReceipt(t, "household-1", userId, groupIds[0], []models.Category{category})

	command := aggregateReportCommand("Doc", groupIds, []string{commands.ReportFormatPdf})
	command.Document = commands.ReportDocument{
		Title:  "{{group.name}} Report",
		Intro:  "Period {{period}}",
		Footer: "Prepared by {{currentUser.name}}",
	}

	report, err := NewReportService(nil).Generate(userId, command)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if report.ContentType != constants.ApplicationPdf {
		t.Errorf("content type = %q, want %q", report.ContentType, constants.ApplicationPdf)
	}
	if report.Filename != "Doc.pdf" {
		t.Errorf("filename = %q, want Doc.pdf", report.Filename)
	}
	if !bytes.HasPrefix(report.Bytes, []byte("%PDF-")) {
		t.Errorf("expected PDF bytes, got prefix %q", string(report.Bytes[:min(8, len(report.Bytes))]))
	}
}

func TestReportService_UserDisplayName(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	named := models.User{Username: "u-named", DisplayName: "Displayed Name"}
	if err := db.Create(&named).Error; err != nil {
		t.Fatalf("seed named user: %v", err)
	}
	unnamed := models.User{Username: "u-unnamed"}
	if err := db.Create(&unnamed).Error; err != nil {
		t.Fatalf("seed unnamed user: %v", err)
	}

	service := NewReportService(nil)
	if got := service.userDisplayName(named.ID); got != "Displayed Name" {
		t.Errorf("display name = %q, want the DisplayName", got)
	}
	if got := service.userDisplayName(unnamed.ID); got != "u-unnamed" {
		t.Errorf("display name = %q, want the Username fallback", got)
	}
	if got := service.userDisplayName(9999999); got != "Unknown User" {
		t.Errorf("display name for a missing user = %q, want Unknown User", got)
	}
}

func TestReportService_GroupNames(t *testing.T) {
	defer repositories.TruncateTestDb()
	db := repositories.GetDB()

	alpha := models.Group{Name: "Alpha"}
	beta := models.Group{Name: "Beta"}
	if err := db.Create(&alpha).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	if err := db.Create(&beta).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	names, err := NewReportService(nil).groupNames([]string{groupIdString(alpha.ID), groupIdString(beta.ID)})
	if err != nil {
		t.Fatalf("groupNames: %v", err)
	}
	if len(names) != 2 || names[0] != "Alpha" || names[1] != "Beta" {
		t.Errorf("group names = %v, want [Alpha Beta] in order", names)
	}
}

func TestReportSpecError_WrapsAndReports(t *testing.T) {
	inner := errors.New("group by field is not a dimension")
	specErr := &ReportSpecError{Err: inner}

	if specErr.Error() != inner.Error() {
		t.Errorf("Error() = %q, want %q", specErr.Error(), inner.Error())
	}
	if !errors.Is(specErr, inner) {
		t.Error("errors.Is should unwrap ReportSpecError to its cause")
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
