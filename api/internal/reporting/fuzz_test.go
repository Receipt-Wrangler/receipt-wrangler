package reporting

import (
	"errors"
	"testing"

	"github.com/expr-lang/expr/ast"
	"github.com/shopspring/decimal"
)

// Formulas are authored by people and stored in templates, so the parser is an
// untrusted-input boundary even though the author is trusted. These targets
// assert that whatever it is handed, it either refuses it or produces something
// the evaluator can safely run.
//
// The seed corpora carry every accepted and rejected expression from
// formula_test.go, so a plain `go test` run already exercises them.

var fuzzFormulaSeeds = []string{
	// Accepted.
	"Subtotal", "2", "1.05", "Subtotal + Hst", "Total - Hst", "Subtotal * 2",
	"Total / Count", "a + b * c", "(a + b) * c", "a - b - c", "a / b / c",
	"-Subtotal", "+Subtotal", "a - -b", "a ++ b", "ROUND(Total / Count, 2)",
	"ROUND(Total, 0)", "ROUND(Total, -2)", "ROUND(a, --2)", "ROUND(a, 2) + ROUND(b, 2)",
	"   a  +  b   ",

	// Rejected: unparsable.
	"", "   ", "a +", "(a + b", "a b", "a & b", "#",

	// Rejected: parses, but outside the whitelist.
	"a ?? b", "a in b", "a ** b", "a % b", "a > b", "a == b", "a and b",
	`a matches "x"`, "not a", "a.b", "[1, 2]", "{a: 1}", `"s"`, "true", "nil",
	"a > b ? 1 : 2", "let x = 1; x", "len(a)", "sum(a)", "round(a, 2)",

	// Rejected: functions and arity.
	"FOO(a)", "SUM(a)", "COUNT()", "AVG(a)", "ROUND(a)", "ROUND(a, 2, 3)",
	"ROUND(a, b)", "ROUND(a, 1.5)", "ROUND(a, 1 + 1)", "ROUND(a, 999)", "ROUND(a, -999)",

	// Numeric edges.
	"9223372036854775807", "9223372036854775808", "0.1 + 0.2", "1e30", "0/0", "1/0",
}

// FuzzParseArithmetic asserts the parser is total and that whatever it accepts,
// the evaluator can run without panicking.
//
// Anything accepted must also survive a second whitelist check and be walkable
// for its column references, since a renderer will later walk the same tree.
func FuzzParseArithmetic(f *testing.F) {
	for _, seed := range fuzzFormulaSeeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, source string) {
		node, err := ParseArithmetic(source)

		if (err == nil) != (node != nil) {
			t.Fatalf("ParseArithmetic(%q) returned node=%v and err=%v", source, node != nil, err)
		}
		if err != nil {
			assertFormulaError(t, source, err)
			return
		}

		// An accepted tree stays accepted, and its references are walkable.
		if again := checkArithmetic(node); again != nil {
			t.Fatalf("ParseArithmetic(%q) accepted a tree that checkArithmetic rejects: %v", source, again)
		}

		columns := map[string]Value{}
		for _, ref := range columnRefs(node) {
			columns[ref] = Num(dec("2"))
		}

		// Evaluation must be total, and must yield a number or nothing.
		value := evalArithmetic(node, columns, 6)
		if value.Type() != ValueNumber && value.Type() != ValueNull {
			t.Fatalf("eval(%q) = %v, want a number or null", source, value)
		}

		if len(columns) == 0 {
			// A constant expression reads nothing; the two cases below are about
			// what happens to the values it would have read.
			return
		}

		// Every column zero: this is where division by zero lives, and it must
		// produce a value rather than a panic.
		for name := range columns {
			columns[name] = Num(dec("0"))
		}
		_ = evalArithmetic(node, columns, 6)

		// Every column null: an expression that reads one is null.
		for name := range columns {
			columns[name] = Null()
		}
		if got := evalArithmetic(node, columns, 6); !got.IsNull() {
			t.Fatalf("eval(%q) over null columns = %v, want null", source, got)
		}
	})
}

// assertFormulaError: a rejection always carries one of the package's sentinels,
// so a caller can tell a typo from a forbidden construct.
func assertFormulaError(t *testing.T, source string, err error) {
	t.Helper()

	sentinels := []error{
		ErrFormulaSyntax, ErrFormulaUnsupported, ErrUnknownFunction,
		ErrBadCallArity, ErrBadRoundPlaces, ErrFormulaTooLong,
	}
	for _, sentinel := range sentinels {
		if errors.Is(err, sentinel) {
			return
		}
	}
	t.Fatalf("ParseArithmetic(%q) failed with an unclassified error: %v", source, err)
}

// FuzzParseAggregate asserts the aggregate front door is total: it returns a
// well-formed aggregate or an error, never both and never a panic.
func FuzzParseAggregate(f *testing.F) {
	seeds := []string{
		"SUM(amount)", "SUM(custom_1)", "COUNT()", "AVG(amount)", "MIN(amount)", "MAX(amount)",
		"  SUM( amount )  ", "", "amount", "SUM(amount) + 1", "SUM(a) + SUM(b)",
		"TOTAL(amount)", "ROUND(amount, 2)", "sum(amount)", "avg(amount)", "round(amount, 2)",
		"len(a)", "COUNT(amount)", "SUM()", "SUM(a, b)", "SUM(a + b)", "SUM(1)", "SUM(",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, source string) {
		aggregate, err := ParseAggregate(source)

		if err != nil {
			if aggregate != (Aggregate{}) {
				t.Fatalf("ParseAggregate(%q) returned %+v alongside %v", source, aggregate, err)
			}
			assertFormulaError(t, source, err)
			return
		}

		switch aggregate.Func {
		case AggSum, AggAvg, AggMin, AggMax:
			if aggregate.Field == "" {
				t.Fatalf("ParseAggregate(%q) = %s with no field", source, aggregate.Func)
			}
		case AggCount:
			if aggregate.Field != "" {
				t.Fatalf("ParseAggregate(%q) = COUNT over %s", source, aggregate.Field)
			}
		default:
			t.Fatalf("ParseAggregate(%q) produced an unknown function %d", source, aggregate.Func)
		}

		// The structural form renders back to something that parses the same.
		reparsed, err := ParseAggregate(aggregate.String())
		if err != nil || reparsed != aggregate {
			t.Fatalf("ParseAggregate(%q).String() = %q, which reparses to %+v (%v)",
				source, aggregate.String(), reparsed, err)
		}
	})
}

// FuzzEvalArithmetic asserts the evaluator is total over arbitrary column
// values, and — the guarantee that matters most — that its answer never depends
// on decimal's mutable, process-wide DivisionPrecision.
func FuzzEvalArithmetic(f *testing.F) {
	f.Add("Total / Count", "215.60", "0", int32(6))
	f.Add("Total / Count", "215.60", "6", int32(2))
	f.Add("Total + Count", "0.1", "0.2", int32(6))
	f.Add("ROUND(Total / Count, 2)", "1", "3", int32(6))
	f.Add("Total * Count - Total / Count", "-5.5", "0.001", int32(0))
	f.Add("ROUND(Total, -2)", "1234", "1", int32(30))

	f.Fuzz(func(t *testing.T, source, left, right string, scale int32) {
		if scale < 0 || scale > maxDivisionScale {
			t.Skip()
		}

		node, err := ParseArithmetic(source)
		if err != nil {
			t.Skip()
		}

		leftValue, leftErr := boundedDecimal(left)
		rightValue, rightErr := boundedDecimal(right)
		if leftErr != nil || rightErr != nil {
			t.Skip()
		}

		columns := map[string]Value{"Total": Num(leftValue), "Count": Num(rightValue)}
		for _, ref := range columnRefs(node) {
			if _, known := columns[ref]; !known {
				columns[ref] = Null()
			}
		}

		got := evalArithmetic(node, columns, scale)
		if got.Type() != ValueNumber && got.Type() != ValueNull {
			t.Fatalf("eval(%q) = %v, want a number or null", source, got)
		}

		// Dividing by zero is an empty cell, never a panic and never an error.
		if rightValue.IsZero() && dividesByCount(node) && !got.IsNull() {
			t.Fatalf("eval(%q) with Count=0 = %v, want null", source, got)
		}

		// Another package mutating decimal.DivisionPrecision must not move the
		// answer. The engine rounds explicitly and never reads the global.
		original := decimal.DivisionPrecision
		defer func() { decimal.DivisionPrecision = original }()

		decimal.DivisionPrecision = 1
		low := evalArithmetic(node, columns, scale)
		decimal.DivisionPrecision = 60
		high := evalArithmetic(node, columns, scale)

		if !low.Equal(got) || !high.Equal(got) {
			t.Fatalf("eval(%q) changed with decimal.DivisionPrecision: %v / %v / %v", source, got, low, high)
		}
	})
}

// dividesByCount reports whether the expression's top level divides by Count,
// which is the only shape the zero-divisor assertion above can predict.
func dividesByCount(node ast.Node) bool {
	binary, isBinary := node.(*ast.BinaryNode)
	if !isBinary || binary.Operator != "/" {
		return false
	}
	identifier, isIdentifier := binary.Right.(*ast.IdentifierNode)
	return isIdentifier && identifier.Value == "Count"
}

// FuzzRunIsTotal drives the whole engine from fuzz bytes: a bounded spec over a
// fixed catalog, and a handful of rows. Whatever comes out, if Validate
// accepted the spec then Run must not panic and the model must obey its
// structural laws.
func FuzzRunIsTotal(f *testing.F) {
	f.Add(uint8(0), uint8(0), uint8(3), uint8(7), true, true, int32(6))
	f.Add(uint8(2), uint8(1), uint8(5), uint8(0), false, true, int32(2))
	f.Add(uint8(3), uint8(4), uint8(0), uint8(255), true, false, int32(0))
	f.Add(uint8(255), uint8(255), uint8(255), uint8(255), false, false, int32(31))

	f.Fuzz(func(t *testing.T, levels, detailBy, rowCount, shape uint8, subtotals, grandTotals bool, scale int32) {
		catalog := fuzzCatalog(t)

		dimensions := []FieldKey{"category", "tag", "paid_by", "date", "resolved"}
		groupBy := make([]FieldKey, 0, 3)
		for index := 0; index < int(levels)%4; index++ {
			groupBy = append(groupBy, dimensions[(int(shape)+index)%len(dimensions)])
		}

		detail := DetailSpec{Mode: DetailRecords}
		if detailBy%2 == 0 {
			detail = DetailSpec{Mode: DetailAggregate, By: dimensions[int(detailBy)%len(dimensions)]}
		}

		spec := ReportSpec{
			GroupBy: groupBy,
			Detail:  detail,
			Columns: []Column{
				{Name: "Cnt", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
				{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
				{Name: "Mean", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: "amount"}},
				{Name: "Ratio", Kind: ColumnArithmetic, Expr: "Total / Cnt"},
			},
			Subtotals:   subtotals,
			GrandTotals: grandTotals,
			Config:      EngineConfig{DivisionScale: scale},
		}

		// A spec Validate rejects is not the engine's problem; that it rejects
		// rather than panicking is, and it already has by the time we get here.
		if err := Validate(spec, catalog); err != nil {
			return
		}

		rows := make([]Row, int(rowCount)%9)
		for index := range rows {
			rows[index] = fuzzRow(index, int(shape))
		}

		model, err := Run(spec, catalog, rows, MetaInput{})
		if err != nil {
			t.Fatalf("Validate accepted a spec that Run rejects: %v", err)
		}

		assertModelInvariants(t, spec, model)

		// Determinism, over whatever shape the fuzzer found.
		again, err := Run(spec, catalog, rows, MetaInput{})
		if err != nil {
			t.Fatalf("Run() error on the second run = %v", err)
		}
		if serializeModel(again) != serializeModel(model) {
			t.Fatalf("two runs over the same input disagree")
		}
	})
}

func fuzzCatalog(t *testing.T) FieldCatalog {
	t.Helper()

	catalog, err := NewFieldCatalog(
		FieldRef{Key: "amount", DataType: TypeCurrency},
		FieldRef{Key: "category", DataType: TypeString, Multi: true},
		FieldRef{Key: "tag", DataType: TypeString, Multi: true},
		FieldRef{Key: "paid_by", DataType: TypeString},
		FieldRef{Key: "date", DataType: TypeDate},
		FieldRef{Key: "resolved", DataType: TypeBool},
	)
	if err != nil {
		t.Fatalf("NewFieldCatalog() error = %v", err)
	}
	return catalog
}

// fuzzRow builds a row whose shape varies with the fuzzed bytes: absent fields,
// explicit nulls, repeated values, and the colliding dates.
func fuzzRow(index, shape int) Row {
	row := Row{}
	mix := (index + shape) % 6

	switch mix {
	case 0:
		// Everything absent: the row lands in every (None) bucket.
	case 1:
		row["category"] = []Value{Str("a"), Str("a")} // a repeat, not a fan-out
		row["amount"] = []Value{Num(dec("10.00"))}
	case 2:
		row["tag"] = []Value{Str("x"), Str("y")} // a genuine fan-out
		row["amount"] = []Value{Num(dec("-1.50"))}
	case 3:
		row["date"] = []Value{DateVal(propertyDates[index%len(propertyDates)])}
		row["amount"] = []Value{Null()} // counted, but contributes no amount
	case 4:
		row["paid_by"] = []Value{Str("")} // the empty string is its own bucket
		row["resolved"] = []Value{Bool(index%2 == 0)}
		row["amount"] = []Value{Num(dec("0"))}
	default:
		row["category"] = []Value{Str("b")}
		row["tag"] = []Value{Str("z")}
		row["paid_by"] = []Value{Str("p")}
		row["amount"] = []Value{Num(dec("3.33"))}
	}

	return row
}
