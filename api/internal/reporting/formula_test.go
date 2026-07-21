package reporting

import (
	"errors"
	"testing"
)

func TestParseArithmetic_Accepts(t *testing.T) {
	sources := []struct {
		name string
		src  string
	}{
		{"column reference", "Subtotal"},
		{"integer literal", "2"},
		{"float literal", "1.05"},
		{"addition", "Subtotal + Hst"},
		{"subtraction", "Total - Hst"},
		{"multiplication", "Subtotal * 2"},
		{"division", "Total / Count"},
		{"precedence", "a + b * c"},
		{"parentheses", "(a + b) * c"},
		{"left associative subtraction", "a - b - c"},
		{"left associative division", "a / b / c"},
		{"unary minus", "-Subtotal"},
		{"unary plus", "+Subtotal"},
		{"subtract a negated column", "a - -b"},
		// The parser reads this as a + (+b), not as a syntax error. Harmless.
		{"doubled sign is a unary plus", "a ++ b"},
		{"round", "ROUND(Total / Count, 2)"},
		{"round to zero places", "ROUND(Total, 0)"},
		{"round to tens", "ROUND(Total, -2)"},
		{"nested round", "ROUND(a, 2) + ROUND(b, 2)"},
		{"leading and trailing whitespace", "   a  +  b   "},
	}

	for _, test := range sources {
		t.Run(test.name, func(t *testing.T) {
			node, err := ParseArithmetic(test.src)
			if err != nil {
				t.Fatalf("ParseArithmetic(%q) error = %v, want nil", test.src, err)
			}
			if node == nil {
				t.Errorf("ParseArithmetic(%q) returned a nil node", test.src)
			}
		})
	}
}

func TestParseArithmetic_Rejects(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr error
	}{
		// Does not parse at all.
		{"empty", "", ErrFormulaSyntax},
		{"whitespace only", "   ", ErrFormulaSyntax},
		{"trailing operator", "a +", ErrFormulaSyntax},
		{"unbalanced parenthesis", "(a + b", ErrFormulaSyntax},
		{"adjacent identifiers", "a b", ErrFormulaSyntax},
		{"unknown operator", "a & b", ErrFormulaSyntax},
		{"stray symbol", "#", ErrFormulaSyntax},

		// Parses, but the operator is outside the whitelist. These are the
		// dangerous ones: the expression language returns a BinaryNode for all
		// of them, so a node-type check alone would let them through.
		{"nil coalescing", "a ?? b", ErrFormulaUnsupported},
		{"membership", "a in b", ErrFormulaUnsupported},
		{"exponent", "a ** b", ErrFormulaUnsupported},
		{"modulo", "a % b", ErrFormulaUnsupported},
		{"comparison", "a > b", ErrFormulaUnsupported},
		{"equality", "a == b", ErrFormulaUnsupported},
		{"logical and", "a and b", ErrFormulaUnsupported},
		{"string match", `a matches "x"`, ErrFormulaUnsupported},
		{"logical not is a unary node", "not a", ErrFormulaUnsupported},

		// Parses, but the node type is outside the whitelist.
		{"property access", "a.b", ErrFormulaUnsupported},
		{"array literal", "[1, 2]", ErrFormulaUnsupported},
		{"map literal", "{a: 1}", ErrFormulaUnsupported},
		{"string literal", `"s"`, ErrFormulaUnsupported},
		{"boolean literal", "true", ErrFormulaUnsupported},
		{"nil literal", "nil", ErrFormulaUnsupported},
		{"conditional", "a > b ? 1 : 2", ErrFormulaUnsupported},
		{"variable declaration", "let x = 1; x", ErrFormulaUnsupported},

		// Functions. The expression language ships lower-case builtins named
		// sum, min, max, count, round and len, so a mis-cased function name
		// parses as a builtin rather than as a call. Both paths must land on
		// the same sentinel.
		{"unknown function", "FOO(a)", ErrUnknownFunction},
		{"sum is aggregate only", "SUM(a)", ErrUnknownFunction},
		{"count is aggregate only", "COUNT()", ErrUnknownFunction},
		{"avg is aggregate only", "AVG(a)", ErrUnknownFunction},
		{"lowercase sum is a builtin", "sum(a)", ErrUnknownFunction},
		{"lowercase avg is a call", "avg(a)", ErrUnknownFunction},
		{"lowercase round is a builtin", "round(a, 2)", ErrUnknownFunction},
		{"unrelated builtin", "len(a)", ErrUnknownFunction},
		{"builtin nested in a sum", "a + len(b)", ErrUnknownFunction},

		// ROUND's shape.
		{"round with one argument", "ROUND(a)", ErrBadCallArity},
		{"round with three arguments", "ROUND(a, 2, 3)", ErrBadCallArity},
		{"round places is a column", "ROUND(a, b)", ErrBadRoundPlaces},
		{"round places is a float", "ROUND(a, 1.5)", ErrBadRoundPlaces},
		{"round places is an expression", "ROUND(a, 1 + 1)", ErrBadRoundPlaces},
		{"round places exceeds the limit", "ROUND(a, 999)", ErrBadRoundPlaces},
		{"round places below the limit", "ROUND(a, -999)", ErrBadRoundPlaces},

		// The whitelist is recursive, not just applied to the root.
		{"unsupported operator nested in a sum", "a + (b % c)", ErrFormulaUnsupported},
		{"unsupported node nested under round", "ROUND(a.b, 2)", ErrFormulaUnsupported},
		{"unsupported node nested under unary", "-[1]", ErrFormulaUnsupported},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node, err := ParseArithmetic(test.src)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("ParseArithmetic(%q) error = %v, want %v", test.src, err, test.wantErr)
			}
			if node != nil {
				t.Errorf("ParseArithmetic(%q) returned a node alongside an error", test.src)
			}
		})
	}
}

func TestParseAggregate_Accepts(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want Aggregate
	}{
		{"sum", "SUM(amount)", Aggregate{Func: AggSum, Field: "amount"}},
		{"sum of a custom field", "SUM(custom_1)", Aggregate{Func: AggSum, Field: "custom_1"}},
		{"count", "COUNT()", Aggregate{Func: AggCount}},
		{"avg", "AVG(amount)", Aggregate{Func: AggAvg, Field: "amount"}},
		{"min", "MIN(amount)", Aggregate{Func: AggMin, Field: "amount"}},
		{"max", "MAX(amount)", Aggregate{Func: AggMax, Field: "amount"}},
		{"whitespace", "  SUM( amount )  ", Aggregate{Func: AggSum, Field: "amount"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseAggregate(test.src)
			if err != nil {
				t.Fatalf("ParseAggregate(%q) error = %v, want nil", test.src, err)
			}
			if got != test.want {
				t.Errorf("ParseAggregate(%q) = %+v, want %+v", test.src, got, test.want)
			}
		})
	}
}

func TestParseAggregate_Rejects(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr error
	}{
		{"empty", "", ErrFormulaSyntax},
		{"not a call", "amount", ErrFormulaUnsupported},
		{"a call plus something", "SUM(amount) + 1", ErrFormulaUnsupported},
		{"two aggregates", "SUM(a) + SUM(b)", ErrFormulaUnsupported},
		{"unknown function", "TOTAL(amount)", ErrUnknownFunction},
		{"round is arithmetic only", "ROUND(amount, 2)", ErrUnknownFunction},
		{"lowercase sum is a builtin", "sum(amount)", ErrUnknownFunction},
		{"lowercase avg is a call", "avg(amount)", ErrUnknownFunction},
		{"lowercase round is a builtin", "round(amount, 2)", ErrUnknownFunction},
		{"unrelated builtin", "len(a)", ErrUnknownFunction},
		{"count takes no arguments", "COUNT(amount)", ErrBadCallArity},
		{"sum needs an argument", "SUM()", ErrBadCallArity},
		{"sum takes one argument", "SUM(a, b)", ErrBadCallArity},
		{"argument must be a field name", "SUM(a + b)", ErrFormulaUnsupported},
		{"argument must not be a literal", "SUM(1)", ErrFormulaUnsupported},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseAggregate(test.src)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("ParseAggregate(%q) error = %v, want %v", test.src, err, test.wantErr)
			}
			if got != (Aggregate{}) {
				t.Errorf("ParseAggregate(%q) returned %+v alongside an error", test.src, got)
			}
		})
	}
}

// Round-tripping keeps the persisted template form and the structural form in
// step.
func TestAggregate_String(t *testing.T) {
	sources := []string{"SUM(amount)", "COUNT()", "AVG(custom_1)", "MIN(amount)", "MAX(amount)"}

	for _, src := range sources {
		t.Run(src, func(t *testing.T) {
			aggregate, err := ParseAggregate(src)
			if err != nil {
				t.Fatalf("ParseAggregate(%q) error = %v", src, err)
			}
			if got := aggregate.String(); got != src {
				t.Errorf("Aggregate.String() = %q, want %q", got, src)
			}
		})
	}
}

func TestColumnRefs(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{"single reference", "Subtotal", []string{"Subtotal"}},
		{"binary", "Subtotal + Hst", []string{"Subtotal", "Hst"}},
		{"first appearance order", "Hst + Subtotal", []string{"Hst", "Subtotal"}},
		{"deduplicated", "a + a * a", []string{"a"}},
		{"no references", "1 + 2", nil},
		{"unary", "-Subtotal", []string{"Subtotal"}},
		{"nested parentheses", "(a + b) * (c - a)", []string{"a", "b", "c"}},
		// The function's name must not be mistaken for a column.
		{"round callee is not a column", "ROUND(Total / Count, 2)", []string{"Total", "Count"}},
		{"round of a constant", "ROUND(1.5, 0)", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node, err := ParseArithmetic(test.src)
			if err != nil {
				t.Fatalf("ParseArithmetic(%q) error = %v", test.src, err)
			}

			got := columnRefs(node)
			if len(got) != len(test.want) {
				t.Fatalf("columnRefs(%q) = %v, want %v", test.src, got, test.want)
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("columnRefs(%q)[%d] = %q, want %q", test.src, i, got[i], test.want[i])
				}
			}
		})
	}
}

func TestIsReservedWord(t *testing.T) {
	reserved := []string{
		"and", "or", "not", "in", "matches", "contains", "startsWith", "endsWith",
		"let", "if", "else", "nil", "true", "false",
		"SUM", "COUNT", "AVG", "MIN", "MAX", "ROUND",
	}
	for _, name := range reserved {
		t.Run("reserved "+name, func(t *testing.T) {
			if !isReservedWord(name) {
				t.Errorf("isReservedWord(%q) = false, want true", name)
			}
		})
	}

	allowed := []string{"Subtotal", "Hst", "Total", "Count", "AvgPerReceipt", "sum", "Round", "In"}
	for _, name := range allowed {
		t.Run("allowed "+name, func(t *testing.T) {
			if isReservedWord(name) {
				t.Errorf("isReservedWord(%q) = true, want false", name)
			}
		})
	}
}
