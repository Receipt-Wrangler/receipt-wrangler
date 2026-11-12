package render

import (
	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/reporting"
)

// FormatCurrency is formatCurrency for callers outside the renderers (the claim
// packet), so every document the server prints formats money one way. It lives in
// its own file so the fork's export does not touch currency.go.
func FormatCurrency(number decimal.Decimal, format reporting.CurrencyFormat) string {
	return formatCurrency(number, format)
}
