package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/reporting/render"
	"receipt-wrangler/api/internal/repositories"
)

// ReportService turns a report-builder request into a downloadable report. It is
// the orchestrator the reporting engine's README calls for: it pulls the data
// (ReportDataService), assembles a spec, runs the pure engine, renders each
// requested format, and (for PDF) bridges the rendered HTML through the
// HTML-to-PDF service — the one place that may cross from the pure render package
// into the services layer. It reads the clock exactly once, at generation time.
type ReportService struct {
	BaseService
}

func NewReportService(tx *gorm.DB) ReportService {
	return ReportService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
}

// GeneratedReport is a ready-to-stream report file: the bytes, a download
// filename, and the Content-Type to send with it.
type GeneratedReport struct {
	Bytes       []byte
	Filename    string
	ContentType string
}

// ReportSpecError wraps a failure to build or compile the report spec — a
// client-caused error (an unknown field key, a wrong-role column, a formula
// cycle) that the handler maps to a 400 rather than a 500.
type ReportSpecError struct {
	Err error
}

func (e *ReportSpecError) Error() string { return e.Err.Error() }
func (e *ReportSpecError) Unwrap() error { return e.Err }

// reportVariablePattern matches any of the four document variables the builder
// exposes, capturing the token name. It is applied in a single pass over the
// original text (see substituteVariables), so a substituted value is never
// re-scanned and the output does not depend on map iteration order.
var reportVariablePattern = regexp.MustCompile(`\{\{\s*(period|group\.name|generatedAt|currentUser\.name)\s*\}\}`)

var reportFilenameUnsafe = regexp.MustCompile(`[^\w]+`)

// Generate produces the report described by command for the given caller. The
// caller must already have been authorized (membership + group.reports.read in
// every group) — this service applies the reporting access controls per group but
// does not itself gate access.
func (service ReportService) Generate(userId uint, command commands.ReportRequestCommand) (GeneratedReport, error) {
	now := time.Now()

	filter := command.Filter
	periodLabel := applyPeriod(&filter, command.Period, now)

	catalog, rows, err := service.loadRows(userId, command.GroupIds, filter)
	if err != nil {
		return GeneratedReport{}, err
	}

	groupNames, err := service.groupNames(command.GroupIds)
	if err != nil {
		return GeneratedReport{}, err
	}

	title, chrome := service.resolveDocument(userId, groupNames, periodLabel, now, command.Document)

	spec, err := buildReportSpec(command)
	if err != nil {
		return GeneratedReport{}, &ReportSpecError{Err: err}
	}
	spec.Title = title

	meta := reporting.MetaInput{
		GeneratedAt: now,
		Params: map[string]string{
			"Period": periodLabel,
			"Groups": strings.Join(groupNames, ", "),
		},
	}

	model, err := reporting.Run(spec, catalog, rows, meta)
	if err != nil {
		return GeneratedReport{}, &ReportSpecError{Err: err}
	}

	dimensions := buildDimensions(spec.GroupBy, catalog)

	files, err := service.renderFormats(command.Formats, model, dimensions, chrome)
	if err != nil {
		return GeneratedReport{}, err
	}

	return service.assembleDownload(command.Name, files)
}

// loadRows gathers the engine rows across every covered group under a single
// catalog. The custom-field catalog is a global pool, so it is identical for each
// group; the rows are concatenated. Rows applies grant-narrowing on its own copy
// of the filter per group, so the caller's filter is reused unchanged.
func (service ReportService) loadRows(
	userId uint,
	groupIds []string,
	filter commands.ReceiptPagedRequestFilter,
) (reporting.FieldCatalog, []reporting.Row, error) {
	dataService := NewReportDataService(service.TX)

	var catalog reporting.FieldCatalog
	var rows []reporting.Row
	for index, groupId := range groupIds {
		groupCatalog, groupRows, err := dataService.Rows(userId, groupId, filter)
		if err != nil {
			return reporting.FieldCatalog{}, nil, err
		}
		if index == 0 {
			catalog = groupCatalog
		}
		rows = append(rows, groupRows...)
	}
	return catalog, rows, nil
}

// renderedFile is one rendered format awaiting download assembly.
type renderedFile struct {
	extension   string
	contentType string
	bytes       []byte
}

// renderFormats renders the requested formats in a fixed order (CSV, XLSX, PDF)
// so multi-format bundles are deterministic regardless of request order.
func (service ReportService) renderFormats(
	formats []string,
	model reporting.ReportModel,
	dimensions []render.Dimension,
	chrome render.DocumentChrome,
) ([]renderedFile, error) {
	requested := make(map[string]bool, len(formats))
	for _, format := range formats {
		requested[format] = true
	}

	var files []renderedFile

	if requested[commands.ReportFormatCsv] {
		bytes, err := render.CSV(model, dimensions)
		if err != nil {
			return nil, err
		}
		files = append(files, renderedFile{commands.ReportFormatCsv, constants.TextCsv, bytes})
	}

	if requested[commands.ReportFormatXlsx] {
		bytes, err := render.XLSX(model, dimensions)
		if err != nil {
			return nil, err
		}
		files = append(files, renderedFile{commands.ReportFormatXlsx, constants.ApplicationXlsx, bytes})
	}

	if requested[commands.ReportFormatPdf] {
		html, err := render.HTML(model, dimensions, chrome)
		if err != nil {
			return nil, err
		}
		pdf, _, err := NewHtmlToPdfService(service.TX).Render(string(html))
		if err != nil {
			return nil, err
		}
		files = append(files, renderedFile{commands.ReportFormatPdf, constants.ApplicationPdf, pdf})
	}

	return files, nil
}

// assembleDownload returns a single file directly, or bundles several into a zip.
func (service ReportService) assembleDownload(name string, files []renderedFile) (GeneratedReport, error) {
	base := reportBaseName(name)

	if len(files) == 1 {
		return GeneratedReport{
			Bytes:       files[0].bytes,
			Filename:    base + "." + files[0].extension,
			ContentType: files[0].contentType,
		}, nil
	}

	names := make([]string, len(files))
	contents := make([][]byte, len(files))
	for index, file := range files {
		names[index] = base + "." + file.extension
		contents[index] = file.bytes
	}

	zipBytes, err := repositories.NewFileRepository(service.TX).ZipFiles(names, contents)
	if err != nil {
		return GeneratedReport{}, err
	}
	return GeneratedReport{
		Bytes:       zipBytes,
		Filename:    base + ".zip",
		ContentType: constants.ApplicationZip,
	}, nil
}

func (service ReportService) groupNames(groupIds []string) ([]string, error) {
	groupRepository := repositories.NewGroupRepository(service.TX)
	names := make([]string, 0, len(groupIds))
	for _, groupId := range groupIds {
		group, err := groupRepository.GetGroupById(groupId, false, false, false)
		if err != nil {
			return nil, err
		}
		names = append(names, group.Name)
	}
	return names, nil
}

func (service ReportService) userDisplayName(userId uint) string {
	user, err := repositories.NewUserRepository(service.TX).GetUserById(userId)
	if err != nil {
		return "Unknown User"
	}
	if len(user.DisplayName) > 0 {
		return user.DisplayName
	}
	return user.Username
}

// buildReportSpec maps the request onto a reporting.ReportSpec (Title is resolved
// and set by the caller). Deep validity — field keys, roles, formula cycles — is
// left to reporting.Run.
func buildReportSpec(command commands.ReportRequestCommand) (reporting.ReportSpec, error) {
	groupBy := make([]reporting.FieldKey, len(command.GroupBy))
	for index, key := range command.GroupBy {
		groupBy[index] = reporting.FieldKey(key)
	}

	detail := reporting.DetailSpec{Mode: reporting.DetailRecords}
	if command.Detail.Mode == commands.ReportDetailAggregate {
		detail = reporting.DetailSpec{Mode: reporting.DetailAggregate, By: reporting.FieldKey(command.Detail.By)}
	}

	columns := make([]reporting.Column, 0, len(command.Columns))
	for _, column := range command.Columns {
		built, err := buildReportColumn(column)
		if err != nil {
			return reporting.ReportSpec{}, err
		}
		columns = append(columns, built)
	}

	return reporting.ReportSpec{
		GroupBy:     groupBy,
		Detail:      detail,
		Columns:     columns,
		Subtotals:   command.Subtotals,
		GrandTotals: command.GrandTotals,
	}, nil
}

func buildReportColumn(column commands.ReportColumn) (reporting.Column, error) {
	built := reporting.Column{Name: column.Name, Label: column.Label}
	switch column.Kind {
	case commands.ReportColumnDimension:
		built.Kind = reporting.ColumnLabel
		built.Field = reporting.FieldKey(column.Field)
	case commands.ReportColumnAggregate:
		built.Kind = reporting.ColumnAggregate
		built.AggSrc = aggregateSource(column.AggFunc, column.Measure)
	case commands.ReportColumnFormula:
		built.Kind = reporting.ColumnArithmetic
		built.Expr = column.Expr
	default:
		return reporting.Column{}, fmt.Errorf("unknown column kind %q", column.Kind)
	}
	return built, nil
}

// aggregateSource renders an aggregate into the engine's source form, e.g.
// "SUM(amount)" or "COUNT()".
func aggregateSource(aggFunc string, measure string) string {
	if aggFunc == "COUNT" {
		return "COUNT()"
	}
	return aggFunc + "(" + measure + ")"
}

// buildDimensions mirrors the spec's group-by into render dimensions, pulling each
// heading label from the catalog.
func buildDimensions(groupBy []reporting.FieldKey, catalog reporting.FieldCatalog) []render.Dimension {
	dimensions := make([]render.Dimension, len(groupBy))
	for index, key := range groupBy {
		label := string(key)
		if field, ok := catalog.Get(key); ok {
			label = field.Label
		}
		dimensions[index] = render.Dimension{Key: key, Label: label}
	}
	return dimensions
}

// applyPeriod resolves the request's period into an inclusive date window, writes
// it onto the filter's Date field, and returns a human-readable label for the
// document preamble and the {{period}} variable. Presets are computed from now;
// custom uses the supplied bounds (already validated by the command).
func applyPeriod(filter *commands.ReceiptPagedRequestFilter, period commands.ReportPeriod, now time.Time) string {
	start, end := resolvePeriodBounds(period, now)

	dayStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
	dayEnd := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, now.Location())

	filter.Date = commands.PagedRequestField{
		Operation: commands.BETWEEN,
		Value:     []interface{}{dayStart, dayEnd},
	}

	return dayStart.Format("2006-01-02") + " to " + dayEnd.Format("2006-01-02")
}

func resolvePeriodBounds(period commands.ReportPeriod, now time.Time) (time.Time, time.Time) {
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	switch period.Preset {
	case commands.ReportPeriodThisMonth:
		return firstOfMonth, firstOfMonth.AddDate(0, 1, -1)
	case commands.ReportPeriodLastMonth:
		lastMonth := firstOfMonth.AddDate(0, -1, 0)
		return lastMonth, firstOfMonth.AddDate(0, 0, -1)
	case commands.ReportPeriodMtd:
		return firstOfMonth, now
	case commands.ReportPeriodQtd:
		quarterStartMonth := ((int(now.Month()) - 1) / 3 * 3) + 1
		return time.Date(now.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, now.Location()), now
	case commands.ReportPeriodYtd:
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()), now
	case commands.ReportPeriodCustom:
		// The command validated these parse; ignoring the error keeps a bad bound
		// from becoming a zero time — which the validator already rejected.
		start, _ := time.Parse("2006-01-02", period.StartDate)
		end, _ := time.Parse("2006-01-02", period.EndDate)
		return start, end
	default:
		return firstOfMonth, now
	}
}

// resolveDocument substitutes the four report variables into the authored
// document copy, returning the resolved title (for the model's Meta) and the
// render-time chrome (intro and footer). The runtime values it resolves — the
// period label, the covered group names, the generation time, and the caller's
// display name — are known only here, at generation time.
func (service ReportService) resolveDocument(
	userId uint,
	groupNames []string,
	periodLabel string,
	now time.Time,
	document commands.ReportDocument,
) (string, render.DocumentChrome) {
	substitutions := map[string]string{
		"period":           periodLabel,
		"group.name":       strings.Join(groupNames, ", "),
		"generatedAt":      now.Format("Jan 2, 2006, 3:04 PM"),
		"currentUser.name": service.userDisplayName(userId),
	}
	return substituteVariables(document.Title, substitutions),
		render.DocumentChrome{
			Intro:  substituteVariables(document.Intro, substitutions),
			Footer: substituteVariables(document.Footer, substitutions),
		}
}

// substituteVariables replaces the supported document variables in text with
// their resolved values in a single pass over the original text. A resolved value
// that itself contains a token is never re-scanned, and unknown tokens (which the
// pattern does not match) are left untouched.
func substituteVariables(text string, values map[string]string) string {
	return reportVariablePattern.ReplaceAllStringFunc(text, func(match string) string {
		token := reportVariablePattern.FindStringSubmatch(match)[1]
		return values[token]
	})
}

// reportBaseName sanitizes a report name into a filesystem-safe download stem,
// falling back to "report" when nothing usable remains.
func reportBaseName(name string) string {
	base := strings.Trim(reportFilenameUnsafe.ReplaceAllString(name, "_"), "_")
	if base == "" {
		return "report"
	}
	return base
}
