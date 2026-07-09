package reporting

import (
	"github.com/expr-lang/expr/ast"
	"github.com/shopspring/decimal"
)

// evalArithmetic evaluates an arithmetic column against the values already
// computed for the other columns on the same row. It is called once per column
// per rendered row — detail rows, subtotal rows, and the grand total alike —
// which is what keeps a non-linear formula such as an average or a ratio
// correct at every level.
//
// Evaluation never fails and never panics. Anything undefined yields a null
// value, which a renderer draws as an empty cell:
//
//   - a null operand propagates, so ROUND(null, 2) and null + 5 are both null
//   - a non-numeric operand is null, though Validate rejects those up front
//   - division by a zero divisor is null rather than a panic, which is what
//     shopspring's Div would otherwise do
//
// All arithmetic is decimal. Division rounds to divisionScale places, chosen by
// the caller, rather than consulting decimal.DivisionPrecision — that global is
// mutable and shared process-wide, so relying on it would make the engine's
// output depend on whatever else the process had done.
func evalArithmetic(node ast.Node, columns map[string]Value, divisionScale int32) Value {
	switch n := node.(type) {
	case *ast.IntegerNode:
		return Num(decimal.NewFromInt(int64(n.Value)))

	case *ast.FloatNode:
		// NewFromFloat takes the shortest decimal that round-trips the float, so
		// a literal a human wrote as 0.1 stays exactly 0.1. Literals carrying
		// more precision than a float64 holds cannot survive the parser, which
		// bounds how exact a literal can be.
		return Num(decimal.NewFromFloat(n.Value))

	case *ast.IdentifierNode:
		return columns[n.Value]

	case *ast.UnaryNode:
		return evalUnary(n, columns, divisionScale)

	case *ast.BinaryNode:
		left := evalArithmetic(n.Left, columns, divisionScale)
		right := evalArithmetic(n.Right, columns, divisionScale)
		return evalBinary(n.Operator, left, right, divisionScale)

	case *ast.CallNode:
		return evalRound(n, columns, divisionScale)
	}

	return Null()
}

func evalUnary(node *ast.UnaryNode, columns map[string]Value, divisionScale int32) Value {
	operand, isNumber := evalArithmetic(node.Node, columns, divisionScale).Decimal()
	if !isNumber {
		return Null()
	}
	if node.Operator == "-" {
		return Num(operand.Neg())
	}
	return Num(operand)
}

func evalBinary(operator string, left, right Value, divisionScale int32) Value {
	leftNumber, leftIsNumber := left.Decimal()
	rightNumber, rightIsNumber := right.Decimal()
	if !leftIsNumber || !rightIsNumber {
		return Null()
	}

	switch operator {
	case "+":
		return Num(leftNumber.Add(rightNumber))
	case "-":
		return Num(leftNumber.Sub(rightNumber))
	case "*":
		return Num(leftNumber.Mul(rightNumber))
	case "/":
		if rightNumber.IsZero() {
			return Null()
		}
		return Num(leftNumber.DivRound(rightNumber, divisionScale))
	}

	return Null()
}

// evalRound applies ROUND, which rounds half away from zero to match the
// spreadsheet function of the same name. Negative places round to tens,
// hundreds, and so on.
func evalRound(call *ast.CallNode, columns map[string]Value, divisionScale int32) Value {
	if len(call.Arguments) != 2 {
		return Null()
	}

	operand, isNumber := evalArithmetic(call.Arguments[0], columns, divisionScale).Decimal()
	if !isNumber {
		return Null()
	}

	places, err := integerLiteral(call.Arguments[1])
	if err != nil {
		return Null()
	}

	return Num(operand.Round(places))
}
