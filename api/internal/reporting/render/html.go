package render

import (
	"bytes"
	"html/template"
	"sort"

	"receipt-wrangler/api/internal/reporting"
)

// HTML renders a ReportModel as a self-contained HTML document: a title, a
// key/value preamble drawn from the report's resolved parameters, the faithful
// grouped table (the same layout the XLSX renderer produces, via the shared
// faithfulWalk), and a generated-at footer. Anything the model does not carry —
// an empty title, no parameters, a zero timestamp — is simply omitted.
//
// The document is deliberately self-contained: all styling is inline and it
// references no external resources or scripts, so it renders faithfully through
// the headless-Chromium HTML-to-PDF pipeline (services/html_to_pdf.go), which
// blocks network loads and disables JavaScript by default. This renderer is the
// PDF format's HTML stage; turning that HTML into PDF bytes is the caller's job.
//
// Like every renderer it is a pure consumer of the model. groupBy supplies the
// dimension order and header labels and must match the report's grouping depth —
// see validateGroupByDepth.
func HTML(model reporting.ReportModel, groupBy []Dimension) ([]byte, error) {
	sink := &htmlSink{model: model}
	if err := faithfulWalk(model, groupBy, sink); err != nil {
		return nil, err
	}

	document := htmlDocument{
		Title:    model.Meta.Title,
		Params:   sortedParams(model.Meta.Params),
		Headings: sink.headings,
		Rows:     sink.rows,
	}
	if !model.Meta.GeneratedAt.IsZero() {
		document.GeneratedAt = model.Meta.GeneratedAt.UTC().Format("2006-01-02 15:04:05 UTC")
	}

	var buffer bytes.Buffer
	if err := htmlTemplate.Execute(&buffer, document); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// htmlDocument is the view model the report template renders. Its strings are
// interpolated through html/template, which escapes them.
type htmlDocument struct {
	Title       string
	Params      []htmlParam
	Headings    []string
	Rows        []htmlRow
	GeneratedAt string
}

type htmlParam struct {
	Key   string
	Value string
}

type htmlRow struct {
	Class string
	Cells []htmlCell
}

type htmlCell struct {
	Text    string
	Numeric bool
}

// htmlSink is the faithfulSink that collects the walk into the view model the
// report template renders.
type htmlSink struct {
	model    reporting.ReportModel
	headings []string
	rows     []htmlRow
}

func (s *htmlSink) writeHeader(headings []string) error {
	s.headings = headings
	return nil
}

func (s *htmlSink) writeRow(kind rowKind, cells []faithfulCell) error {
	row := htmlRow{Class: rowClass(kind), Cells: make([]htmlCell, len(cells))}
	for index, cell := range cells {
		row.Cells[index] = s.cellView(cell)
	}
	s.rows = append(s.rows, row)
	return nil
}

// cellView maps a positioned walk cell to its view model: a marker/dimension is
// plain text, a report column is formatted the same way the CSV renderer formats
// it and is right-aligned unless it is a label.
func (s *htmlSink) cellView(cell faithfulCell) htmlCell {
	switch cell.kind {
	case textCell:
		return htmlCell{Text: cell.text}
	case reportCell:
		return htmlCell{
			Text:    formatCell(cell.descriptor, cell.cell, s.model.Meta.NoneLabel),
			Numeric: cell.descriptor.Kind != reporting.ColumnLabel,
		}
	default: // blankCell
		return htmlCell{}
	}
}

// rowClass names a data row's CSS class; detail rows are unclassed.
func rowClass(kind rowKind) string {
	switch kind {
	case subtotalKind:
		return "subtotal"
	case grandTotalKind:
		return "grand-total"
	default:
		return ""
	}
}

// sortedParams flattens the report's resolved parameters into a stable,
// alphabetically ordered slice, so the preamble renders deterministically.
func sortedParams(params map[string]string) []htmlParam {
	if len(params) == 0 {
		return nil
	}
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ordered := make([]htmlParam, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, htmlParam{Key: key, Value: params[key]})
	}
	return ordered
}

var htmlTemplate = template.Must(template.New("report").Parse(reportHTML))

// reportHTML is the document skeleton: a styled title, a parameter preamble, the
// data table, and a generated-at footer. All CSS is inline so the document needs
// no external resources when it is converted to PDF.
const reportHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{if .Title}}{{.Title}}{{else}}Report{{end}}</title>
<style>
  body { font-family: Arial, Helvetica, sans-serif; color: #222; margin: 24px; }
  h1 { font-size: 20px; margin: 0 0 12px; }
  .preamble { margin: 0 0 16px; color: #555; font-size: 12px; }
  .preamble div { margin: 2px 0; }
  .preamble .key { font-weight: 600; }
  table { border-collapse: collapse; width: 100%; font-size: 12px; }
  th, td { border: 1px solid #ccc; padding: 4px 8px; text-align: left; }
  thead th { background: #f2f2f2; font-weight: 700; }
  td.num { text-align: right; font-variant-numeric: tabular-nums; }
  tr.subtotal td { font-weight: 700; background: #fafafa; }
  tr.grand-total td { font-weight: 700; background: #f2f2f2; border-top: 2px solid #888; }
  footer { margin-top: 16px; color: #777; font-size: 11px; }
</style>
</head>
<body>
{{if .Title}}<h1>{{.Title}}</h1>
{{end}}{{if .Params}}<div class="preamble">
{{range .Params}}<div><span class="key">{{.Key}}</span>: {{.Value}}</div>
{{end}}</div>
{{end}}<table>
<thead><tr>{{range .Headings}}<th>{{.}}</th>{{end}}</tr></thead>
<tbody>
{{range .Rows}}<tr{{if .Class}} class="{{.Class}}"{{end}}>{{range .Cells}}<td{{if .Numeric}} class="num"{{end}}>{{.Text}}</td>{{end}}</tr>
{{end}}</tbody>
</table>
{{if .GeneratedAt}}<footer>Generated {{.GeneratedAt}}</footer>
{{end}}</body>
</html>
`
