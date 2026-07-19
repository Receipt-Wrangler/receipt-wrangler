package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"

	"receipt-wrangler/api/internal/reporting"
)

// usd/eur are the two representative configurations: the US default (symbol
// leads, comma thousands, dot decimal) and a common European format (symbol
// trails, dot thousands, comma decimal).
var (
	usd = reporting.CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: "."}
	eur = reporting.CurrencyFormat{Symbol: "€", SymbolAtEnd: true, ThousandsSeparator: ".", DecimalSeparator: ","}
)

// TestFormatCurrency pins the money string to what the desktop customCurrency
// pipe produces for the same configuration, so a rendered report reads
// identically to the rest of the UI.
func TestFormatCurrency(t *testing.T) {
	cases := []struct {
		name   string
		value  string
		format reporting.CurrencyFormat
		want   string
	}{
		{"usd basic", "1234.5", usd, "$1,234.50"},
		{"usd under a thousand", "12.5", usd, "$12.50"},
		{"usd zero", "0", usd, "$0.00"},
		{"usd millions", "1234567.89", usd, "$1,234,567.89"},
		{"usd integer length divisible by three", "123456", usd, "$123,456.00"},
		{"usd rounds to two", "1.239", usd, "$1.24"},
		{"usd negative keeps symbol first", "-1234.5", usd, "$-1,234.50"},

		{"eur trailing symbol", "1234.5", eur, "1.234,50€"},
		{"eur millions", "1234567.89", eur, "1.234.567,89€"},
		{"eur negative", "-1234.5", eur, "-1.234,50€"},

		{"hide decimals truncates the rounded value", "1234.56",
			reporting.CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: ".", HideDecimals: true},
			"$1,234"},
		{"hide decimals rounds first", "1234.996",
			reporting.CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: ".", HideDecimals: true},
			"$1,235"},

		{"no symbol", "1234.5",
			reporting.CurrencyFormat{ThousandsSeparator: ",", DecimalSeparator: "."}, "1,234.50"},
		{"blank separators fall back to us defaults", "1234.5",
			reporting.CurrencyFormat{Symbol: "$"}, "$1,234.50"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := formatCurrency(decimal.RequireFromString(testCase.value), testCase.format)
			if got != testCase.want {
				t.Errorf("formatCurrency(%s) = %q, want %q", testCase.value, got, testCase.want)
			}
		})
	}
}

// TestExcelCurrencyFormat pins the Excel number-format code: the symbol is a
// quoted literal on the configured side and the decimal places follow
// HideDecimals. Grouping/decimal glyphs are left to Excel's locale.
func TestExcelCurrencyFormat(t *testing.T) {
	cases := []struct {
		name   string
		format reporting.CurrencyFormat
		want   string
	}{
		{"leading symbol", usd, `"$"#,##0.00`},
		{"trailing symbol", eur, `#,##0.00"€"`},
		{"hide decimals", reporting.CurrencyFormat{Symbol: "$", HideDecimals: true}, `"$"#,##0`},
		{"no symbol", reporting.CurrencyFormat{ThousandsSeparator: ",", DecimalSeparator: "."}, `#,##0.00`},
		{"no symbol hide decimals", reporting.CurrencyFormat{HideDecimals: true}, `#,##0`},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := excelCurrencyFormat(testCase.format); got != testCase.want {
				t.Errorf("excelCurrencyFormat = %q, want %q", got, testCase.want)
			}
		})
	}
}

// TestFormatCell_CurrencyConfig proves the shared CSV/HTML cell formatter applies
// the currency configuration to a currency column when one is supplied, and keeps
// the bare two-place formatting when it is nil.
func TestFormatCell_CurrencyConfig(t *testing.T) {
	currency := reporting.ColumnDescriptor{Kind: reporting.ColumnAggregate, DataType: reporting.TypeCurrency}
	number := reporting.ColumnDescriptor{Kind: reporting.ColumnAggregate, DataType: reporting.TypeNumber}
	cell := reporting.Cell{Values: []reporting.Value{money("1234.5")}}

	if got := formatCell(currency, cell, "(None)", &eur); got != "1.234,50€" {
		t.Errorf("currency cell with config = %q, want %q", got, "1.234,50€")
	}
	if got := formatCell(currency, cell, "(None)", nil); got != "1234.50" {
		t.Errorf("currency cell without config = %q, want %q", got, "1234.50")
	}
	// A currency config never reshapes a plain (non-currency) number column.
	if got := formatCell(number, cell, "(None)", &eur); got != "1234.5" {
		t.Errorf("number cell = %q, want %q", got, "1234.5")
	}
}

// The three end-to-end wiring tests below prove each renderer reads
// model.Meta.Currency. oneLevelRows' first detail bucket (Dana / Food) has a
// Total of 100 + 200 = 300, so the currency column D2 / the "$300.00" string is
// the tell.

func TestCSV_UsesConfiguredCurrencyFormat(t *testing.T) {
	model := mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows())
	model.Meta.Currency = &usd

	out, err := CSV(model, paidByDimension())
	if err != nil {
		t.Fatalf("CSV: %v", err)
	}
	if !strings.Contains(string(out), "$300.00") {
		t.Errorf("CSV missing the configured currency total; got:\n%s", out)
	}
}

func TestHTML_UsesConfiguredCurrencyFormat(t *testing.T) {
	model := mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows())
	model.Meta.Currency = &eur

	out, err := HTML(model, paidByDimension(), DocumentChrome{})
	if err != nil {
		t.Fatalf("HTML: %v", err)
	}
	if !strings.Contains(string(out), "300,00€") {
		t.Errorf("HTML missing the configured currency total; got:\n%s", out)
	}
}

func TestXLSX_UsesConfiguredCurrencyFormat(t *testing.T) {
	model := mustRun(t, oneLevelSpec(), paidByCatalog(t), oneLevelRows())
	model.Meta.Currency = &eur

	out, err := XLSX(model, paidByDimension())
	if err != nil {
		t.Fatalf("XLSX: %v", err)
	}
	file, err := excelize.OpenReader(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	defer file.Close()

	// D2 (the first detail Total, a currency column) carries the configured
	// number-format code — it stays a native number, only its display changes.
	want := excelCurrencyFormat(eur)
	if fmtStr := cellStyle(t, file, "D2").CustomNumFmt; fmtStr == nil || *fmtStr != want {
		t.Errorf("D2 number format = %v, want %q", fmtStr, want)
	}
}
