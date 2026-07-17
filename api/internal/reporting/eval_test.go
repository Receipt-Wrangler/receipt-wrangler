package reporting

import (
	"strings"
	"testing"
	"time"

	"github.com/expr-lang/expr/ast"
)

const testDivisionScale = 6

// runawayDigitCeiling is an independent sanity bound for the arithmetic-growth
// tests, deliberately NOT the production maxDecimalDigits: if the real bound were
// weakened or removed, a surviving value would blow past this and the test would
// fail — a test law must not be derived from the code it judges. Money never needs
// anywhere near this many significant digits, and the real bound is well below it.
const runawayDigitCeiling = 2000

func mustParseArithmetic(t *testing.T, src string) ast.Node {
	t.Helper()
	node, err := ParseArithmetic(src)
	if err != nil {
		t.Fatalf("ParseArithmetic(%q) error = %v", src, err)
	}
	return node
}

// evalSrc parses and evaluates in one step, which is how every eval case below
// reads.
func evalSrc(t *testing.T, src string, columns map[string]Value) Value {
	t.Helper()
	return evalArithmetic(mustParseArithmetic(t, src), columns, testDivisionScale)
}

func TestEvalArithmetic_Numbers(t *testing.T) {
	columns := map[string]Value{
		"Subtotal": Num(dec("200.00")),
		"Hst":      Num(dec("15.60")),
		"Count":    Num(dec("6")),
		"Zero":     Num(dec("0")),
		"Negative": Num(dec("-4.50")),
	}

	tests := []struct {
		name string
		src  string
		want string
	}{
		{"integer literal", "2", "2"},
		{"float literal", "1.05", "1.05"},
		{"column reference", "Subtotal", "200"},
		{"addition", "Subtotal + Hst", "215.6"},
		{"subtraction", "Subtotal - Hst", "184.4"},
		{"multiplication", "Hst * 2", "31.2"},
		{"division", "Subtotal / 4", "50"},
		{"precedence", "Subtotal + Hst * 2", "231.2"},
		{"parentheses override precedence", "(Subtotal + Hst) * 2", "431.2"},
		{"left associative subtraction", "Subtotal - Hst - Hst", "168.8"},
		{"unary minus on a column", "-Subtotal", "-200"},
		{"unary minus on a literal", "-2 + Subtotal", "198"},
		{"double sign parses as unary plus", "Subtotal ++ Hst", "215.6"},
		{"subtracting a negative", "Subtotal - Negative", "204.5"},
		{"negative column", "Negative * 2", "-9"},
		{"zero is a legal operand", "Subtotal + Zero", "200"},
		{"zero numerator", "Zero / Count", "0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := evalSrc(t, test.src, columns)
			number, isNumber := got.Decimal()
			if !isNumber {
				t.Fatalf("eval(%q) = %v, want a number", test.src, got)
			}
			if number.String() != test.want {
				t.Errorf("eval(%q) = %s, want %s", test.src, number, test.want)
			}
		})
	}
}

// Division rounds to the configured scale rather than consulting the mutable
// decimal.DivisionPrecision global.
func TestEvalArithmetic_DivisionScale(t *testing.T) {
	columns := map[string]Value{
		"Total": Num(dec("215.60")),
		"Count": Num(dec("6")),
	}

	tests := []struct {
		scale int32
		want  string
	}{
		{0, "36"},
		{2, "35.93"},
		{6, "35.933333"},
		{10, "35.9333333333"},
	}

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			got := evalArithmetic(mustParseArithmetic(t, "Total / Count"), columns, test.scale)
			number, _ := got.Decimal()
			if number.String() != test.want {
				t.Errorf("Total / Count at scale %d = %s, want %s", test.scale, number, test.want)
			}
		})
	}
}

// A computed value that grows past what money can mean collapses to null instead
// of being materialized. Both growth modes are covered: a coefficient that keeps
// doubling (a memory-exhaustion vector) and an exponent that keeps climbing
// toward the int32 boundary (a silent-overflow vector). Feeding each column's
// result back into the next is exactly what the engine's arithmetic columns do.
func TestEvalArithmetic_BoundsMagnitude(t *testing.T) {
	t.Run("a squaring chain collapses to null rather than growing without bound", func(t *testing.T) {
		columns := map[string]Value{"c": Num(dec("99"))}
		maxDigits := 0
		for i := 0; i < 40; i++ {
			columns["c"] = evalSrc(t, "c * c", columns)
			if number, ok := columns["c"].Decimal(); ok && number.NumDigits() > maxDigits {
				maxDigits = number.NumDigits()
			}
		}
		if !columns["c"].IsNull() {
			t.Errorf("deep squaring chain = %v, want a null cell", columns["c"])
		}
		if maxDigits > runawayDigitCeiling {
			t.Errorf("a value survived with %d digits, past the independent %d ceiling", maxDigits, runawayDigitCeiling)
		}
	})

	t.Run("a large-exponent chain collapses to null instead of wrapping the int32 exponent", func(t *testing.T) {
		columns := map[string]Value{"c": evalSrc(t, "1e308", map[string]Value{})}
		for i := 0; i < 40; i++ {
			columns["c"] = evalSrc(t, "c * c", columns)
			if columns["c"].IsNull() {
				return // collapsed before it could wrap to a wrong number — the guard working
			}
		}
		t.Fatalf("a 1e308 squaring chain never collapsed: %v", columns["c"])
	})
}

// ROUND funnels through the same magnitude guard as every other arithmetic op, so
// an operand past the coefficient ceiling rounds to a null cell rather than a
// materialized giant. 1100 significant digits is comfortably past the bound; the
// null result — not the digit count — is the law being asserted.
func TestEvalArithmetic_RoundIsBounded(t *testing.T) {
	columns := map[string]Value{"big": Num(dec(strings.Repeat("9", 1100)))}
	got := evalArithmetic(mustParseArithmetic(t, "ROUND(big, 2)"), columns, testDivisionScale)
	if !got.IsNull() {
		t.Errorf("ROUND of an over-ceiling operand = %v, want a null cell", got)
	}
}

// A zero divisor must yield an empty cell. shopspring panics on this, so the
// guard is the only thing standing between a report with a zero count and a
// crashed request.
func TestEvalArithmetic_DivisionByZeroIsNullNotPanic(t *testing.T) {
	columns := map[string]Value{
		"Total":     Num(dec("215.60")),
		"Count":     Num(dec("0")),
		"ZeroScale": Num(dec("0.00")),
	}

	tests := []struct {
		name string
		src  string
	}{
		{"column divisor", "Total / Count"},
		{"zero with a scale", "Total / ZeroScale"},
		{"literal divisor", "Total / 0"},
		{"nested in a sum", "1 + Total / Count"},
		{"inside round", "ROUND(Total / Count, 2)"},
		{"zero over zero", "Count / Count"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("eval(%q) panicked: %v", test.src, recovered)
				}
			}()

			if got := evalSrc(t, test.src, columns); !got.IsNull() {
				t.Errorf("eval(%q) = %v, want null", test.src, got)
			}
		})
	}
}

// Null propagates through every operator, so a subtotal built on an empty
// aggregate renders as an empty cell rather than as zero.
func TestEvalArithmetic_NullPropagates(t *testing.T) {
	columns := map[string]Value{
		"Subtotal": Num(dec("200.00")),
		"Missing":  Null(),
		"Label":    Str("Clothing"),
		"When":     DateVal(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)),
		"Flag":     Bool(true),
	}

	tests := []struct {
		name string
		src  string
	}{
		{"null on the left", "Missing + Subtotal"},
		{"null on the right", "Subtotal + Missing"},
		{"null subtraction", "Subtotal - Missing"},
		{"null multiplication", "Subtotal * Missing"},
		{"null numerator", "Missing / Subtotal"},
		{"null divisor", "Subtotal / Missing"},
		{"null under unary minus", "-Missing"},
		{"null inside round", "ROUND(Missing, 2)"},
		{"null nested deep", "ROUND((Subtotal + Missing) * 2, 2)"},
		{"undeclared column is null", "Subtotal + Nonexistent"},

		// Validate rejects these up front; eval must still not misread them.
		{"string operand", "Subtotal + Label"},
		{"date operand", "Subtotal + When"},
		{"bool operand", "Subtotal + Flag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := evalSrc(t, test.src, columns); !got.IsNull() {
				t.Errorf("eval(%q) = %v, want null", test.src, got)
			}
		})
	}
}

func TestEvalArithmetic_Round(t *testing.T) {
	columns := map[string]Value{
		"Half":     Num(dec("32.925")),
		"Long":     Num(dec("35.93333")),
		"Negative": Num(dec("-2.345")),
		"Thousand": Num(dec("1234")),
	}

	tests := []struct {
		name string
		src  string
		want string
	}{
		{"rounds half away from zero", "ROUND(Half, 2)", "32.93"},
		{"truncates a long tail", "ROUND(Long, 2)", "35.93"},
		{"rounds negative away from zero", "ROUND(Negative, 2)", "-2.35"},
		{"rounds to zero places", "ROUND(Long, 0)", "36"},
		{"negative places round to hundreds", "ROUND(Thousand, -2)", "1200"},
		{"rounds an expression", "ROUND(Half * 2, 1)", "65.9"},
		{"round is composable", "ROUND(Half, 2) + ROUND(Long, 2)", "68.86"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := evalSrc(t, test.src, columns)
			number, isNumber := got.Decimal()
			if !isNumber {
				t.Fatalf("eval(%q) = %v, want a number", test.src, got)
			}
			if number.String() != test.want {
				t.Errorf("eval(%q) = %s, want %s", test.src, number, test.want)
			}
		})
	}
}

// Money never touches a float. These are the sums that would drift if it did.
func TestEvalArithmetic_NoFloatDrift(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		columns map[string]Value
		want    string
	}{
		{
			name:    "0.1 + 0.2 is exactly 0.3",
			src:     "A + B",
			columns: map[string]Value{"A": Num(dec("0.1")), "B": Num(dec("0.2"))},
			want:    "0.3",
		},
		{
			name: "float literals are exact",
			src:  "0.1 + 0.2",
			want: "0.3",
		},
		{
			name: "a literal keeps every digit it was written with",
			src:  "1.05",
			want: "1.05",
		},
		{
			name:    "33.33 times 3",
			src:     "A * 3",
			columns: map[string]Value{"A": Num(dec("33.33"))},
			want:    "99.99",
		},
		{
			name:    "a hundred cents make a dollar",
			src:     "A * 100",
			columns: map[string]Value{"A": Num(dec("0.01"))},
			want:    "1",
		},
		{
			name:    "the worked example's grand total",
			src:     "Subtotal + Hst",
			columns: map[string]Value{"Subtotal": Num(dec("320.00")), "Hst": Num(dec("27.30"))},
			want:    "347.3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := evalSrc(t, test.src, test.columns)
			number, isNumber := got.Decimal()
			if !isNumber {
				t.Fatalf("eval(%q) = %v, want a number", test.src, got)
			}
			if !number.Equal(dec(test.want)) {
				t.Errorf("eval(%q) = %s, want %s", test.src, number, test.want)
			}
		})
	}
}

// The engine must not read the process-wide division precision, because another
// package could have changed it.
func TestEvalArithmetic_IgnoresDivisionPrecisionGlobal(t *testing.T) {
	columns := map[string]Value{"A": Num(dec("1")), "B": Num(dec("3"))}

	got := evalArithmetic(mustParseArithmetic(t, "A / B"), columns, 4)
	number, _ := got.Decimal()

	if number.String() != "0.3333" {
		t.Errorf("A / B at scale 4 = %s, want 0.3333", number)
	}
}
