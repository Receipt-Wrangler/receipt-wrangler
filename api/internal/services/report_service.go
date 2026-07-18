package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
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

// ReportPreview is the rendered HTML for the live builder preview plus the true
// number of receipts the current configuration covers (reported even when the
// rendered sample is capped). AllowedActions is populated only by the dashboard
// report-widget render path (RenderTemplateForUser) — it carries the actions the
// caller may perform on the template so the widget can gate its download button;
// the builder's own Preview leaves it nil (omitted).
type ReportPreview struct {
	Html           string   `json:"html"`
	ReceiptCount   int      `json:"receiptCount"`
	AllowedActions []string `json:"allowedActions,omitempty"`
}

// reportPreviewRowCap bounds how many receipt rows a preview feeds to the engine.
// A preview is a sample rendered on every builder edit (debounced), so beyond
// this many rows it truncates the sample to stay fast; ReceiptCount still reports
// the true total so the builder's "N receipts" chip is accurate. Normal
// period-bounded reports fall well under the cap and render in full.
const reportPreviewRowCap = 1000

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
	build, err := service.buildModel(userId, command, time.Now(), 0)
	if err != nil {
		return GeneratedReport{}, err
	}

	files, err := service.renderFormats(command.Formats, build.model, build.dimensions, build.chrome)
	if err != nil {
		return GeneratedReport{}, err
	}

	return service.assembleDownload(command.Name, files)
}

// Preview renders the report described by command as HTML for the live builder
// preview. It runs the same pipeline as Generate up through the engine — so the
// preview is the engine's own output, never a client-side approximation — but
// emits only HTML (no PDF bridge, no zip) and caps the rendered sample
// (reportPreviewRowCap). The same per-group authorization as Generate is the
// caller's responsibility.
func (service ReportService) Preview(userId uint, command commands.ReportRequestCommand) (ReportPreview, error) {
	return service.renderHTML(userId, command, reportPreviewRowCap)
}

// renderHTML runs the shared engine pipeline and renders the report as HTML. It
// backs both the row-capped builder preview (rowLimit = reportPreviewRowCap) and
// the dashboard report widget, which renders the full dataset (rowLimit = 0, like
// Generate). ReceiptCount is always the true pre-cap total.
func (service ReportService) renderHTML(userId uint, command commands.ReportRequestCommand, rowLimit int) (ReportPreview, error) {
	build, err := service.buildModel(userId, command, time.Now(), rowLimit)
	if err != nil {
		return ReportPreview{}, err
	}

	html, err := render.HTML(build.model, build.dimensions, build.chrome)
	if err != nil {
		return ReportPreview{}, err
	}

	return ReportPreview{Html: string(html), ReceiptCount: build.receiptCount}, nil
}

// restrictedReportHTML is the self-contained document the report widget renders in
// place of a report when the caller may not view the pinned template (access
// revoked, or the template was deleted). It inlines all styling and references no
// external resources or scripts, so it renders inside the widget's sandboxed
// iframe (sandbox="allow-same-origin", scripts disabled) exactly like real report
// output.
const restrictedReportHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<style>
  html, body { margin: 0; height: 100%; }
  body {
    display: flex; align-items: center; justify-content: center;
    min-height: 260px; padding: 24px; box-sizing: border-box;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    color: #5f6368; text-align: center;
  }
  .notice { max-width: 360px; }
  .notice__icon { font-size: 40px; line-height: 1; margin-bottom: 12px; }
  .notice__title { font-size: 16px; font-weight: 600; color: #3c4043; margin: 0 0 6px; }
  .notice__body { font-size: 13px; margin: 0; }
</style>
</head>
<body>
  <div class="notice">
    <div class="notice__icon" aria-hidden="true">&#128274;</div>
    <p class="notice__title">Restricted</p>
    <p class="notice__body">You no longer have access to this report.</p>
  </div>
</body>
</html>`

// RenderTemplateForUser renders a saved report template as full-dataset HTML for
// the dashboard report widget, re-resolving the caller's access to the template on
// every call. Unlike the builder's Preview it feeds the full dataset to the engine
// (no row cap) so the widget shows the complete report, and — because the widget
// always renders whatever HTML it gets back — a revoked/deleted template yields the
// restrictedReportHTML notice at a normal 200 (empty AllowedActions) rather than an
// error. AllowedActions carries the actions the caller may perform on the template
// so the widget can gate its download button off the server result.
func (service ReportService) RenderTemplateForUser(userId uint, templateId string) (ReportPreview, error) {
	template, err := repositories.NewReportTemplateRepository(service.TX).GetReportTemplateById(templateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ReportPreview{Html: restrictedReportHTML}, nil
		}
		return ReportPreview{}, err
	}

	actions, err := NewPermissionService(service.TX).AllowedActionsForTemplate(userId, template.ID)
	if err != nil {
		return ReportPreview{}, err
	}
	if !slices.Contains(actions, "read") {
		return ReportPreview{Html: restrictedReportHTML}, nil
	}

	var command commands.ReportRequestCommand
	if err := json.Unmarshal(template.Configuration, &command); err != nil {
		return ReportPreview{}, err
	}

	preview, err := service.renderHTML(userId, command, 0)
	if err != nil {
		return ReportPreview{}, err
	}
	preview.AllowedActions = actions
	return preview, nil
}

// reportBuild is the engine output plus the render inputs both Generate and
// Preview consume: the model, the group-by dimensions, the render-time document
// chrome, and the true receipt count.
type reportBuild struct {
	model        reporting.ReportModel
	dimensions   []render.Dimension
	chrome       render.DocumentChrome
	receiptCount int
}

// buildModel runs the shared pipeline both entry points need: resolve the period,
// load rows across every covered group under the one catalog, resolve the document
// variables, build the spec, and run the pure engine. rowLimit > 0 caps the rows
// fed to the engine (for the preview sample); receiptCount is always the true
// pre-cap total. It reads the clock once, via the now passed in.
func (service ReportService) buildModel(
	userId uint,
	command commands.ReportRequestCommand,
	now time.Time,
	rowLimit int,
) (reportBuild, error) {
	filter := command.Filter
	resolveReportGeneratorPaidBy(&filter, userId)
	periodLabel := applyPeriod(&filter, command.Period, now)

	catalog, rows, err := service.loadRows(userId, command.GroupIds, filter)
	if err != nil {
		return reportBuild{}, err
	}
	receiptCount := len(rows)
	if rowLimit > 0 && len(rows) > rowLimit {
		rows = rows[:rowLimit]
	}

	groupNames, err := service.groupNames(command.GroupIds)
	if err != nil {
		return reportBuild{}, err
	}

	title, chrome := service.resolveDocument(userId, groupNames, periodLabel, now, command.Name, command.Document)

	spec, err := buildReportSpec(command)
	if err != nil {
		return reportBuild{}, &ReportSpecError{Err: err}
	}
	spec.Title = title

	settings, err := repositories.NewSystemSettingsRepository(service.TX).GetSystemSettings()
	if err != nil {
		return reportBuild{}, err
	}

	meta := reporting.MetaInput{
		GeneratedAt: now,
		Params: map[string]string{
			"Period": periodLabel,
			"Groups": strings.Join(groupNames, ", "),
		},
		Currency: currencyFormat(settings),
	}

	model, err := reporting.Run(spec, catalog, rows, meta)
	if err != nil {
		return reportBuild{}, &ReportSpecError{Err: err}
	}

	return reportBuild{
		model:        model,
		dimensions:   buildDimensions(spec.GroupBy, catalog),
		chrome:       chrome,
		receiptCount: receiptCount,
	}, nil
}

// currencyFormat maps the app's System Settings currency configuration onto the
// engine's renderer hint, so every rendered format (the live preview included)
// presents money exactly as the rest of the UI does. It is always supplied — the
// settings row is a get-or-create singleton — so a report is never rendered with
// bare, symbol-less numbers.
func currencyFormat(settings models.SystemSettings) *reporting.CurrencyFormat {
	return &reporting.CurrencyFormat{
		Symbol:             settings.CurrencyDisplay,
		SymbolAtEnd:        settings.CurrencySymbolPosition == models.END,
		ThousandsSeparator: string(settings.CurrencyThousandthsSeparator),
		DecimalSeparator:   string(settings.CurrencyDecimalSeparator),
		HideDecimals:       settings.CurrencyHideDecimalPlaces,
	}
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
		// A saved template may carry a dimension column that is currently
		// unresolvable — a label the engine can't place on an aggregated row
		// because its field is neither the aggregate dimension nor a grouping
		// level. The builder greys these out and leaves them out of the request;
		// the same projection here lets a template that stores one (so it
		// round-trips into the builder and re-enables when the config makes it
		// valid) still generate — with the column omitted — instead of failing
		// the whole report. Mirrors the desktop isDimensionColumnDisabled rule
		// and the engine's own compileLabelColumn rejection.
		if isDisabledDimensionColumn(command, column) {
			continue
		}

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

// isDisabledDimensionColumn reports whether a dimension column is unresolvable in
// aggregate mode — its field is neither the aggregate dimension nor a grouping
// level, so the engine cannot label an aggregated row with it. It is the
// server-side twin of the desktop builder's isDimensionColumnDisabled: such a
// column is projected out of the spec rather than failing the report, so a saved
// template may retain it for round-trip fidelity.
func isDisabledDimensionColumn(command commands.ReportRequestCommand, column commands.ReportColumn) bool {
	if command.Detail.Mode != commands.ReportDetailAggregate || column.Kind != commands.ReportColumnDimension {
		return false
	}
	return column.Field != command.Detail.By && !slices.Contains(command.GroupBy, column.Field)
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

// reportGeneratorPaidBySentinel is the paid-by filter value the desktop report builder
// stores for its "Whoever generates the report" option (desktop
// REPORT_GENERATOR_PAID_BY_ID). It is negative so it can never collide with a real
// (positive) user id.
const reportGeneratorPaidBySentinel = -1

// resolveReportGeneratorPaidBy substitutes the report-generator sentinel in the paid-by
// filter with the id of the user generating the report, so a saved template stays
// dynamic — whoever runs it filters to their own receipts, regardless of who authored
// it. It runs before the query is built (BuildGormFilterQuery otherwise passes the raw
// -1 straight into `paid_by_user_id IN (...)`), covering both Generate and Preview
// since both funnel through buildModel. A fresh slice is built so the caller's command
// is never mutated. Values arrive as float64 (JSON) and are matched/substituted as
// such, matching how a static user-id filter already flows through the same IN clause.
func resolveReportGeneratorPaidBy(filter *commands.ReceiptPagedRequestFilter, userId uint) {
	values, ok := filter.PaidBy.Value.([]interface{})
	if !ok {
		return
	}

	resolved := make([]interface{}, 0, len(values))
	for _, value := range values {
		if id, isNumber := value.(float64); isNumber && id == reportGeneratorPaidBySentinel {
			resolved = append(resolved, float64(userId))
		} else {
			resolved = append(resolved, value)
		}
	}
	filter.PaidBy.Value = resolved
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
// display name — are known only here, at generation time. The visible heading is
// the authored document title; when it is left blank it falls back to the report
// name so the rendered report (and its live preview) is never headingless.
func (service ReportService) resolveDocument(
	userId uint,
	groupNames []string,
	periodLabel string,
	now time.Time,
	name string,
	document commands.ReportDocument,
) (string, render.DocumentChrome) {
	substitutions := map[string]string{
		"period":           periodLabel,
		"group.name":       strings.Join(groupNames, ", "),
		"generatedAt":      now.Format("Jan 2, 2006, 3:04 PM"),
		"currentUser.name": service.userDisplayName(userId),
	}
	titleSource := document.Title
	if strings.TrimSpace(titleSource) == "" {
		titleSource = name
	}
	return substituteVariables(titleSource, substitutions),
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
