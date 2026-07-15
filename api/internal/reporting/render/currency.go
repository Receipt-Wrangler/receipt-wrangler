package render

import (
	"strings"

	"github.com/shopspring/decimal"

	"receipt-wrangler/api/internal/reporting"
)

// formatCurrency renders a decimal as a money string per the app's currency
// configuration. It mirrors the desktop customCurrency pipe
// (desktop/src/pipes/custom-currency.pipe.ts) so money in a rendered report
// reads identically to money elsewhere in the UI, including the report's
// receipts drill-in dialog: round to two places, group the integer digits, swap
// the separators, optionally drop the fractional part, and place the symbol at
// the start (ahead of any sign) or the end.
//
// Like the pipe, it always rounds to two places first and then decides whether
// to show the fraction — HideDecimals truncates the rounded value rather than
// rounding to whole — and grouping is always present, defaulting to the US
// separators when the config leaves one blank.
func formatCurrency(number decimal.Decimal, format reporting.CurrencyFormat) string {
	thousandsSeparator := format.ThousandsSeparator
	if thousandsSeparator == "" {
		thousandsSeparator = ","
	}
	decimalSeparator := format.DecimalSeparator
	if decimalSeparator == "" {
		decimalSeparator = "."
	}

	// StringFixed gives a canonical, rounded "-1234.50": an optional leading '-',
	// ASCII digits, and a '.' decimal point we then re-separate ourselves.
	fixed := number.StringFixed(2)

	sign := ""
	if strings.HasPrefix(fixed, "-") {
		sign = "-"
		fixed = fixed[1:]
	}

	integerPart := fixed
	fractionPart := ""
	if dot := strings.IndexByte(fixed, '.'); dot >= 0 {
		integerPart = fixed[:dot]
		fractionPart = fixed[dot+1:]
	}

	body := groupThousands(integerPart, thousandsSeparator)
	if !format.HideDecimals {
		body += decimalSeparator + fractionPart
	}
	body = sign + body

	if format.Symbol == "" {
		return body
	}
	if format.SymbolAtEnd {
		return body + format.Symbol
	}
	return format.Symbol + body
}

// groupThousands inserts separator between every third digit from the right.
// digits is a run of ASCII digits with no sign or decimal point.
func groupThousands(digits, separator string) string {
	if len(digits) <= 3 {
		return digits
	}

	// The first group is the leading len%3 digits (a full group of 3 when the
	// length divides evenly), then groups of three follow.
	lead := len(digits) % 3
	if lead == 0 {
		lead = 3
	}

	var out strings.Builder
	out.WriteString(digits[:lead])
	for index := lead; index < len(digits); index += 3 {
		out.WriteString(separator)
		out.WriteString(digits[index : index+3])
	}
	return out.String()
}

// excelCurrencyFormat builds an Excel number-format code from the currency
// configuration: a #,##0 / #,##0.00 body (per HideDecimals) with the symbol as a
// quoted literal on the configured side. Excel renders the grouping and decimal
// glyphs per the opening application's locale, so the custom thousands/decimal
// separators are not forced here — the symbol, its position, and the decimal
// places are. This keeps XLSX cells native, typed numbers (analyzable) rather
// than preformatted text.
func excelCurrencyFormat(format reporting.CurrencyFormat) string {
	body := "#,##0"
	if !format.HideDecimals {
		body += ".00"
	}

	if format.Symbol == "" {
		return body
	}

	// A quoted literal is how Excel renders an arbitrary symbol; strip any stray
	// quote so the code stays well-formed.
	symbol := `"` + strings.ReplaceAll(format.Symbol, `"`, "") + `"`
	if format.SymbolAtEnd {
		return body + symbol
	}
	return symbol + body
}
