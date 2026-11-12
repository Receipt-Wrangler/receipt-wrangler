package services

import (
	"html/template"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/models"
)

// packet.pdf's content: a view model built from the claim and its receipts, and
// the HTML template Chromium prints. Everything is interpolated through
// html/template, so user text is escaped; the only URLs are the embedded images'
// data: URIs, typed template.URL because html/template would otherwise refuse
// them. The document loads nothing external, which the HTML-to-PDF service blocks
// anyway.

type packetDocument struct {
	ClaimName          string
	Details            []packetField
	Notes              string
	Figures            []packetField
	GeneratedAt        string
	GeneratedBy        string
	Ledger             []packetLedgerRow
	Categories         []packetCategoryGroup
	ReceiptCount       int
	GrandTotal         string
	MultiCategoryCount int
	SubtotalsTotal     string
	Receipts           []packetReceiptSection
}

type packetField struct {
	Label string
	Value string
}

type packetLedgerRow struct {
	Date   string
	Type   string
	Amount string
	Note   string
}

type packetCategoryGroup struct {
	Name     string
	Rows     []packetReceiptRow
	Subtotal string
}

type packetReceiptRow struct {
	Id            uint
	Date          string
	Name          string
	PaidBy        string
	Categories    string
	Tags          string
	Amount        string
	Status        string
	Images        []string
	MultiCategory bool
}

type packetReceiptSection struct {
	Row          packetReceiptRow
	Fields       []packetField
	CustomFields []packetField
	Items        []packetItemRow
	Comments     []packetComment
	Images       []packetImage
}

type packetItemRow struct {
	Name      string
	Amount    string
	ChargedTo string
	Status    string
}

type packetComment struct {
	Author string
	When   string
	Text   string
}

type packetImage struct {
	Name         string
	OriginalName string
	Sha256       string
	Slices       []template.URL
	Note         string
}

const packetInstantFormat = "2006-01-02 15:04 MST"

// uncategorisedLabel heads the receipts with no category the viewer can see.
const uncategorisedLabel = "Uncategorised"

// document assembles the PDF's view model from the claim and the receipt
// sections already built (with their images embedded).
func (builder *claimPacketBuilder) document(sections []packetReceiptSection) packetDocument {
	claim := builder.claim
	document := packetDocument{
		ClaimName:    claim.Name,
		Notes:        claim.Notes,
		GeneratedAt:  builder.now.Format(packetInstantFormat),
		GeneratedBy:  builder.username,
		ReceiptCount: len(builder.receipts),
		Receipts:     sections,
	}

	document.Details = nonEmptyFields(
		packetField{"Reference number", claim.ReferenceNumber},
		packetField{"Type", claim.Type},
		packetField{"Status", humanizeCode(claim.Status)},
		packetField{"Submitted", claim.SubmittedDate},
		packetField{"Due", claim.DueDate},
	)

	claimedLabel := "Claimed (receipts total)"
	if claim.ClaimedAmount != "" {
		claimedLabel = "Claimed"
	}
	document.Figures = []packetField{
		{"Receipts", strconv.Itoa(claim.ReceiptCount)},
		{"Receipts total", builder.moneyString(claim.ReceiptsTotal)},
		{claimedLabel, builder.moneyString(claim.EffectiveClaimedAmount)},
		{"Payments", builder.moneyString(claim.PaymentsTotal)},
		{"Recoupments", builder.moneyString(claim.RecoupmentsTotal)},
		{"Paid (payments less recoupments)", builder.moneyString(claim.PaidAmount)},
		{"Outstanding", builder.moneyString(claim.OutstandingAmount)},
	}
	if claim.DaysWaiting != nil {
		document.Figures = append(document.Figures, packetField{"Days waiting", strconv.Itoa(*claim.DaysWaiting)})
	}

	document.Ledger = make([]packetLedgerRow, 0, len(claim.LedgerEntries))
	for _, entry := range claim.LedgerEntries {
		document.Ledger = append(document.Ledger, packetLedgerRow{
			Date:   entry.Date,
			Type:   humanizeCode(entry.Type),
			Amount: builder.moneyString(entry.Amount),
			Note:   entry.Note,
		})
	}

	folds, total, multiCategory := foldPacketCategories(builder.receipts)
	subtotalsTotal := decimal.Zero
	for _, fold := range folds {
		group := packetCategoryGroup{Name: fold.Name, Subtotal: builder.money(fold.Subtotal)}
		for _, index := range fold.ReceiptIndexes {
			row := sections[index].Row
			row.MultiCategory = len(builder.receipts[index].Categories) > 1
			group.Rows = append(group.Rows, row)
		}
		subtotalsTotal = subtotalsTotal.Add(fold.Subtotal)
		document.Categories = append(document.Categories, group)
	}
	document.GrandTotal = builder.money(total)
	document.MultiCategoryCount = multiCategory
	document.SubtotalsTotal = builder.money(subtotalsTotal)

	return document
}

// moneyString formats a two-place decimal string from the claim view; "" stays "".
func (builder *claimPacketBuilder) moneyString(value string) string {
	if value == "" {
		return ""
	}
	amount, err := decimal.NewFromString(value)
	if err != nil {
		return value
	}
	return builder.money(amount)
}

// packetCategoryFold is one category's receipts (indexes into the receipt list,
// in its order) and their subtotal.
type packetCategoryFold struct {
	Name           string
	ReceiptIndexes []int
	Subtotal       decimal.Decimal
}

// foldPacketCategories groups receipts by category for the receipts table. A
// receipt with several categories is listed under each with its full amount, so
// the subtotals can add up to more than the total; multiCategory counts those
// receipts so the PDF can say so. A receipt with none goes under Uncategorised,
// which comes last; the rest sort by name. total counts each receipt once.
// Folded with shopspring/decimal, so every figure is exact. Pure.
func foldPacketCategories(receipts []models.Receipt) (folds []packetCategoryFold, total decimal.Decimal, multiCategory int) {
	byId := make(map[uint]*packetCategoryFold)
	var uncategorised *packetCategoryFold
	total = decimal.Zero

	for index, receipt := range receipts {
		total = total.Add(receipt.Amount)

		seen := make(map[uint]struct{}, len(receipt.Categories))
		for _, category := range receipt.Categories {
			if _, duplicate := seen[category.ID]; duplicate {
				continue
			}
			seen[category.ID] = struct{}{}

			fold, ok := byId[category.ID]
			if !ok {
				fold = &packetCategoryFold{Name: category.Name, Subtotal: decimal.Zero}
				byId[category.ID] = fold
			}
			fold.ReceiptIndexes = append(fold.ReceiptIndexes, index)
			fold.Subtotal = fold.Subtotal.Add(receipt.Amount)
		}

		if len(seen) > 1 {
			multiCategory++
		}
		if len(seen) == 0 {
			if uncategorised == nil {
				uncategorised = &packetCategoryFold{Name: uncategorisedLabel, Subtotal: decimal.Zero}
			}
			uncategorised.ReceiptIndexes = append(uncategorised.ReceiptIndexes, index)
			uncategorised.Subtotal = uncategorised.Subtotal.Add(receipt.Amount)
		}
	}

	ids := make([]uint, 0, len(byId))
	for id := range byId {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		left, right := strings.ToLower(byId[ids[i]].Name), strings.ToLower(byId[ids[j]].Name)
		if left != right {
			return left < right
		}
		return ids[i] < ids[j]
	})

	folds = make([]packetCategoryFold, 0, len(ids)+1)
	for _, id := range ids {
		folds = append(folds, *byId[id])
	}
	if uncategorised != nil {
		folds = append(folds, *uncategorised)
	}
	return folds, total, multiCategory
}

// buildPacketReceiptSection lays out one receipt's fields, custom fields, items
// and comments. Images are added by the builder as it reads them. A receipt Date
// is a calendar day (read in UTC); instants are shown in the app zone.
func buildPacketReceiptSection(receipt models.Receipt, location *time.Location, money func(decimal.Decimal) string) packetReceiptSection {
	row := packetReceiptRow{
		Id:         receipt.ID,
		Date:       receipt.Date.UTC().Format(time.DateOnly),
		Name:       receipt.Name,
		PaidBy:     receipt.PaidByUser.DisplayName,
		Categories: joinCategoryNames(receipt.Categories),
		Tags:       joinTagNames(receipt.Tags),
		Amount:     money(receipt.Amount),
		Status:     humanizeCode(string(receipt.Status)),
	}

	added := receipt.CreatedAt.In(location).Format(packetInstantFormat)
	if receipt.CreatedByString != "" {
		added += " by " + receipt.CreatedByString
	}
	resolved := ""
	if receipt.ResolvedDate != nil {
		resolved = receipt.ResolvedDate.In(location).Format(time.DateOnly)
	}

	section := packetReceiptSection{
		Row: row,
		Fields: nonEmptyFields(
			packetField{"Receipt id", strconv.FormatUint(uint64(receipt.ID), 10)},
			packetField{"Date", row.Date},
			packetField{"Amount", row.Amount},
			packetField{"Paid by", row.PaidBy},
			packetField{"Status", row.Status},
			packetField{"Categories", row.Categories},
			packetField{"Tags", row.Tags},
			packetField{"Added", added},
			packetField{"Resolved", resolved},
			packetField{"Emailed from", receipt.EmailSender},
		),
	}

	for _, value := range receipt.CustomFields {
		if text := packetCustomFieldValue(value, money); text != "" {
			section.CustomFields = append(section.CustomFields, packetField{value.CustomField.Name, text})
		}
	}

	for _, item := range receipt.ReceiptItems {
		section.Items = append(section.Items, packetItemRow{
			Name:      item.Name,
			Amount:    money(item.Amount),
			ChargedTo: item.ChargedToUser.DisplayName,
			Status:    humanizeCode(string(item.Status)),
		})
	}

	for _, comment := range receipt.Comments {
		author := ""
		if comment.UserId != nil {
			author = comment.User.DisplayName
		}
		section.Comments = append(section.Comments, packetComment{
			Author: author,
			When:   comment.CreatedAt.In(location).Format(packetInstantFormat),
			Text:   comment.Comment,
		})
	}

	return section
}

// packetCustomFieldValue prints a custom field's value by its type; "" when it
// has none.
func packetCustomFieldValue(value models.CustomFieldValue, money func(decimal.Decimal) string) string {
	switch value.CustomField.Type {
	case models.TEXT:
		if value.StringValue != nil {
			return *value.StringValue
		}
	case models.DATE:
		if value.DateValue != nil {
			return value.DateValue.UTC().Format(time.DateOnly)
		}
	case models.SELECT:
		if value.SelectValue != nil {
			for _, option := range value.CustomField.Options {
				if option.ID == *value.SelectValue {
					return option.Value
				}
			}
		}
	case models.CURRENCY:
		if value.CurrencyValue != nil {
			return money(*value.CurrencyValue)
		}
	case models.BOOLEAN:
		if value.BooleanValue != nil {
			if *value.BooleanValue {
				return "Yes"
			}
			return "No"
		}
	}
	return ""
}

func joinCategoryNames(categories []models.Category) string {
	names := make([]string, 0, len(categories))
	for _, category := range categories {
		names = append(names, category.Name)
	}
	return strings.Join(names, ", ")
}

func joinTagNames(tags []models.Tag) string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return strings.Join(names, ", ")
}

// humanizeCode turns a stored code such as NEEDS_ATTENTION into "Needs attention".
func humanizeCode(code string) string {
	text := strings.ToLower(strings.ReplaceAll(code, "_", " "))
	if text == "" {
		return ""
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

func nonEmptyFields(fields ...packetField) []packetField {
	kept := make([]packetField, 0, len(fields))
	for _, field := range fields {
		if field.Value != "" {
			kept = append(kept, field)
		}
	}
	return kept
}

func renderClaimPacketHtml(document packetDocument, embeddedBytes int) (string, error) {
	// A strings.Builder hands its buffer to String() without a copy, and sizing it
	// up front (base64 is 4/3 of the embedded JPEG bytes, plus the text) spares
	// the doubling copies, which matters when the document carries megabytes of
	// embedded images.
	var html strings.Builder
	html.Grow(embeddedBytes*4/3 + 256<<10)
	if err := claimPacketTemplate.Execute(&html, document); err != nil {
		return "", err
	}
	return html.String(), nil
}

var claimPacketTemplate = template.Must(template.New("claimPacket").Parse(claimPacketHtml))

const claimPacketHtml = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Claim packet: {{.ClaimName}}</title>
<style>
  @page { size: letter; margin: 0.5in; }
  body { font-family: Arial, Helvetica, sans-serif; color: #222; font-size: 11px; margin: 0; }
  h1 { font-size: 20px; margin: 0 0 4px; }
  h2 { font-size: 15px; margin: 18px 0 6px; border-bottom: 1px solid #999; padding-bottom: 2px; }
  h3 { font-size: 12px; margin: 12px 0 4px; }
  .subtitle { color: #555; margin: 0 0 12px; }
  .page { break-before: page; }
  table { border-collapse: collapse; width: 100%; }
  th, td { border: 1px solid #ccc; padding: 3px 6px; text-align: left; vertical-align: top; }
  thead th { background: #f2f2f2; }
  thead { display: table-header-group; }
  tr { break-inside: avoid; }
  td.num, th.num { text-align: right; font-variant-numeric: tabular-nums; white-space: nowrap; }
  table.fields th { width: 30%; background: #fafafa; font-weight: 600; }
  tr.group td { background: #e9eef5; font-weight: 700; }
  tr.subtotal td { background: #fafafa; font-weight: 700; }
  tr.grand-total td { background: #f2f2f2; font-weight: 700; border-top: 2px solid #888; }
  .files { font-family: Consolas, "Courier New", monospace; font-size: 9px; word-break: break-all; }
  .note { color: #555; font-style: italic; }
  .notes { white-space: pre-wrap; border: 1px solid #ccc; padding: 6px; }
  .comment { border-left: 3px solid #ccc; padding: 2px 8px; margin: 4px 0; break-inside: avoid; }
  .comment .meta { color: #555; font-size: 10px; }
  .comment .text { white-space: pre-wrap; }
  .image { margin: 10px 0; }
  .image .caption { font-size: 10px; margin: 0 0 4px; }
  .image .hash { font-family: Consolas, "Courier New", monospace; font-size: 8px; color: #555; word-break: break-all; }
  .image img { display: block; margin: 4px auto; max-width: 100%; max-height: 9.3in; break-inside: avoid; }
</style>
</head>
<body>

<h1>Reimbursement claim: {{.ClaimName}}</h1>
<p class="subtitle">Claim packet generated on {{.GeneratedAt}} by {{.GeneratedBy}}</p>

{{if .Details}}<table class="fields">
{{range .Details}}<tr><th>{{.Label}}</th><td>{{.Value}}</td></tr>
{{end}}</table>
{{end}}
<h2>Figures</h2>
<table class="fields">
{{range .Figures}}<tr><th>{{.Label}}</th><td class="num">{{.Value}}</td></tr>
{{end}}</table>

{{if .Notes}}<h2>Notes</h2>
<div class="notes">{{.Notes}}</div>
{{end}}
<h2>Ledger</h2>
{{if .Ledger}}<table>
<thead><tr><th>Date</th><th>Type</th><th class="num">Amount</th><th>Note</th></tr></thead>
<tbody>
{{range .Ledger}}<tr><td>{{.Date}}</td><td>{{.Type}}</td><td class="num">{{.Amount}}</td><td>{{.Note}}</td></tr>
{{end}}</tbody>
</table>
{{else}}<p class="note">No ledger entries.</p>
{{end}}
<h2 class="page">Receipts</h2>
{{if .Categories}}<table>
<thead><tr><th>Date</th><th>Name</th><th>Paid by</th><th>Categories</th><th>Tags</th><th class="num">Amount</th><th>Status</th><th>Image files</th></tr></thead>
<tbody>
{{range .Categories}}<tr class="group"><td colspan="8">{{.Name}}</td></tr>
{{range .Rows}}<tr><td>{{.Date}}</td><td>{{.Name}} (#{{.Id}}){{if .MultiCategory}} *{{end}}</td><td>{{.PaidBy}}</td><td>{{.Categories}}</td><td>{{.Tags}}</td><td class="num">{{.Amount}}</td><td>{{.Status}}</td><td class="files">{{range .Images}}{{.}}<br>{{end}}</td></tr>
{{end}}<tr class="subtotal"><td colspan="5">Subtotal: {{.Name}}</td><td class="num">{{.Subtotal}}</td><td colspan="2"></td></tr>
{{end}}<tr class="grand-total"><td colspan="5">Total of {{.ReceiptCount}} receipt{{if ne .ReceiptCount 1}}s{{end}} (each counted once)</td><td class="num">{{.GrandTotal}}</td><td colspan="2"></td></tr>
</tbody>
</table>
{{if .MultiCategoryCount}}<p class="note">* This receipt has more than one category, so it is listed under each of them with its full amount. The category subtotals therefore add up to {{.SubtotalsTotal}}, more than the total of {{.GrandTotal}}; {{.MultiCategoryCount}} receipt{{if ne .MultiCategoryCount 1}}s are{{else}} is{{end}} listed more than once.</p>
{{end}}{{else}}<p class="note">This claim has no receipts.</p>
{{end}}
{{range .Receipts}}<div class="page">
<h2>{{.Row.Name}} (#{{.Row.Id}})</h2>
<table class="fields">
{{range .Fields}}<tr><th>{{.Label}}</th><td>{{.Value}}</td></tr>
{{end}}{{range .CustomFields}}<tr><th>{{.Label}}</th><td>{{.Value}}</td></tr>
{{end}}</table>
{{if .Items}}<h3>Items</h3>
<table>
<thead><tr><th>Name</th><th class="num">Amount</th><th>Charged to</th><th>Status</th></tr></thead>
<tbody>
{{range .Items}}<tr><td>{{.Name}}</td><td class="num">{{.Amount}}</td><td>{{.ChargedTo}}</td><td>{{.Status}}</td></tr>
{{end}}</tbody>
</table>
{{end}}{{if .Comments}}<h3>Comments</h3>
{{range .Comments}}<div class="comment"><div class="meta">{{if .Author}}{{.Author}}, {{end}}{{.When}}</div><div class="text">{{.Text}}</div></div>
{{end}}{{end}}{{if .Images}}<h3>Images</h3>
{{range .Images}}<div class="image">
<p class="caption">{{if .Sha256}}images/{{.Name}}{{if .OriginalName}} (uploaded as {{.OriginalName}}){{end}}{{else}}{{.OriginalName}}{{end}}</p>
{{if .Sha256}}<p class="hash">SHA-256 {{.Sha256}}</p>
{{end}}{{if .Note}}<p class="note">{{.Note}}</p>
{{end}}{{range .Slices}}<img src="{{.}}" alt="">
{{end}}</div>
{{end}}{{else}}<p class="note">This receipt has no images.</p>
{{end}}</div>
{{end}}
</body>
</html>
`
