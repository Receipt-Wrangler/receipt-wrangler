package reporting

import (
	"errors"
	"fmt"
	"strings"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

var (
	// ErrFormulaSyntax is returned when an expression does not parse.
	ErrFormulaSyntax = errors.New("formula syntax error")

	// ErrFormulaUnsupported is returned when an expression parses but uses a
	// construct outside the whitelist.
	ErrFormulaUnsupported = errors.New("formula uses an unsupported construct")

	// ErrUnknownFunction is returned when an expression calls a function that
	// is not allowed in that position.
	ErrUnknownFunction = errors.New("unknown or disallowed function in formula")

	// ErrBadCallArity is returned when a function is called with the wrong
	// number of arguments.
	ErrBadCallArity = errors.New("wrong number of arguments to function")

	// ErrBadRoundPlaces is returned when ROUND's second argument is not a small
	// integer literal.
	ErrBadRoundPlaces = errors.New("ROUND places must be an integer literal")

	// ErrFormulaTooLong is returned when an expression is longer than any real
	// formula needs and long enough to be dangerous.
	ErrFormulaTooLong = errors.New("formula is too long")
)

// roundFunction is the only function an arithmetic expression may call. The
// vertical reducers (SUM, COUNT, AVG, MIN, MAX) define aggregate columns
// instead, and are rejected here.
const roundFunction = "ROUND"

// roundPlacesLimit bounds ROUND's second argument. Rescaling a decimal to an
// arbitrary number of places would let a user-authored formula allocate without
// bound, and no report needs more than a handful of digits either way.
const roundPlacesLimit = 30

// maxFormulaLength bounds the source an expression may be parsed from.
//
// The expression language caps the tree it will build at ten thousand nodes, and
// that cap is often mistaken for a bound on the work parsing costs. It is not.
// Parentheses group without producing a node, so they are never counted, and the
// parser descends through several frames for each one. Nesting is therefore
// bounded only by the goroutine stack: measured at roughly 640 bytes of stack
// per parenthesis, Go's default one-gigabyte ceiling falls at about 1.6 million
// of them — a few megabytes of source. A Go stack overflow is a fatal error that
// recover cannot catch, so it takes the process down rather than the request.
//
// Bounding the input is the fix, for the same reason roundPlacesLimit and
// maxDivisionScale are bounded: a template a user authored must not be able to
// ask for unbounded work. A kilobyte is far longer than any real formula and
// leaves the parser at most a few hundred frames deep.
const maxFormulaLength = 1024

// reservedWords are identifiers the expression language treats as operators or
// keywords. A column may not be named one of these, since a formula referencing
// it would not parse as a column reference.
var reservedWords = map[string]struct{}{
	"and": {}, "or": {}, "not": {}, "in": {},
	"matches": {}, "contains": {}, "startsWith": {}, "endsWith": {},
	"let": {}, "if": {}, "else": {},
	"nil": {}, "true": {}, "false": {},
}

// isReservedWord reports whether a name collides with the expression language's
// grammar or with a function the engine defines.
func isReservedWord(name string) bool {
	if _, reserved := reservedWords[name]; reserved {
		return true
	}
	if name == roundFunction {
		return true
	}
	_, isAggregate := aggFuncFromName(name)
	return isAggregate
}

// ParseArithmetic parses an arithmetic column expression such as
// "Subtotal + Hst" or "ROUND(Total / Count, 2)" and returns its syntax tree.
//
// Parsing is delegated to the expression language, but the resulting tree is
// then checked against an allow-list: numeric literals, column references, the
// operators + - * /, unary sign, and ROUND. Everything else the language can
// express — property access, comparisons, conditionals, arrays, string and
// boolean literals, builtins, variable declarations — is rejected. The check is
// a closed type switch whose default rejects, so the whitelist holds by
// construction rather than by remembering to deny each new construct.
//
// Identifiers are not resolved here; Validate checks them against the spec's
// columns. The returned tree is exported so a renderer can walk it, which is
// how the XLSX renderer translates a column expression into a live spreadsheet
// formula.
func ParseArithmetic(src string) (ast.Node, error) {
	if err := checkFormulaLength(src); err != nil {
		return nil, err
	}

	tree, err := parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFormulaSyntax, err)
	}
	if err := checkArithmetic(tree.Node); err != nil {
		return nil, err
	}
	return tree.Node, nil
}

// ParseAggregate parses an aggregate column expression such as "SUM(amount)" or
// "COUNT()" into its structural form. An aggregate column is exactly one
// reduction over one field; composing reductions is done with an arithmetic
// column over the aggregate columns.
//
// Field names are not resolved here; Validate checks them against the catalog.
func ParseAggregate(src string) (Aggregate, error) {
	if err := checkFormulaLength(src); err != nil {
		return Aggregate{}, err
	}

	tree, err := parser.Parse(src)
	if err != nil {
		return Aggregate{}, fmt.Errorf("%w: %s", ErrFormulaSyntax, err)
	}

	// A mis-cased reducer such as sum(amount) collides with one of the
	// expression language's own builtins, so it parses as a builtin rather than
	// as a call.
	if builtin, isBuiltin := tree.Node.(*ast.BuiltinNode); isBuiltin {
		return Aggregate{}, aggregateFunctionError(builtin.Name)
	}

	call, isCall := tree.Node.(*ast.CallNode)
	if !isCall {
		return Aggregate{}, fmt.Errorf(
			"%w: an aggregate column must be a single call such as SUM(amount), got %s",
			ErrFormulaUnsupported, nodeKind(tree.Node))
	}

	name, named := calleeName(call)
	if !named {
		return Aggregate{}, fmt.Errorf("%w: computed function call", ErrFormulaUnsupported)
	}

	aggFunc, known := aggFuncFromName(name)
	if !known {
		return Aggregate{}, aggregateFunctionError(name)
	}

	if aggFunc == AggCount {
		if len(call.Arguments) != 0 {
			return Aggregate{}, fmt.Errorf("%w: COUNT takes no arguments, got %d",
				ErrBadCallArity, len(call.Arguments))
		}
		return Aggregate{Func: AggCount}, nil
	}

	if len(call.Arguments) != 1 {
		return Aggregate{}, fmt.Errorf("%w: %s takes 1 argument, got %d",
			ErrBadCallArity, name, len(call.Arguments))
	}

	field, isIdentifier := call.Arguments[0].(*ast.IdentifierNode)
	if !isIdentifier {
		return Aggregate{}, fmt.Errorf("%w: %s must be given a field name, got %s",
			ErrFormulaUnsupported, name, nodeKind(call.Arguments[0]))
	}

	return Aggregate{Func: aggFunc, Field: FieldKey(field.Value)}, nil
}

// checkFormulaLength refuses source long enough to make parsing it expensive.
// It runs before the parser, because the parser is what would be made to do the
// work.
func checkFormulaLength(src string) error {
	if len(src) > maxFormulaLength {
		return fmt.Errorf("%w: %d bytes, limit is %d", ErrFormulaTooLong, len(src), maxFormulaLength)
	}
	return nil
}

// columnRefs returns the distinct column names an arithmetic expression
// references, in order of first appearance so error messages stay stable. A
// function's name is not a column reference.
func columnRefs(node ast.Node) []string {
	var refs []string
	seen := make(map[string]struct{})

	var walk func(ast.Node)
	walk = func(node ast.Node) {
		switch n := node.(type) {
		case *ast.IdentifierNode:
			if _, exists := seen[n.Value]; !exists {
				seen[n.Value] = struct{}{}
				refs = append(refs, n.Value)
			}
		case *ast.UnaryNode:
			walk(n.Node)
		case *ast.BinaryNode:
			walk(n.Left)
			walk(n.Right)
		case *ast.CallNode:
			// Skip the callee: ROUND is a function, not a column.
			for _, argument := range n.Arguments {
				walk(argument)
			}
		}
	}
	walk(node)

	return refs
}

// checkArithmetic enforces the arithmetic whitelist over a parsed tree.
func checkArithmetic(node ast.Node) error {
	switch n := node.(type) {
	case *ast.IntegerNode, *ast.FloatNode:
		return nil

	case *ast.IdentifierNode:
		return nil

	case *ast.UnaryNode:
		if n.Operator != "-" && n.Operator != "+" {
			return fmt.Errorf("%w: unary operator %q", ErrFormulaUnsupported, n.Operator)
		}
		return checkArithmetic(n.Node)

	case *ast.BinaryNode:
		if !isArithmeticOperator(n.Operator) {
			return fmt.Errorf("%w: operator %q", ErrFormulaUnsupported, n.Operator)
		}
		if err := checkArithmetic(n.Left); err != nil {
			return err
		}
		return checkArithmetic(n.Right)

	case *ast.CallNode:
		return checkCall(n)

	case *ast.BuiltinNode:
		// The expression language defines lower-case builtins of its own —
		// sum, min, max, count, round, len — so a mis-cased function name
		// arrives here rather than as a call.
		return arithmeticFunctionError(n.Name)
	}

	return fmt.Errorf("%w: %s", ErrFormulaUnsupported, nodeKind(node))
}

// checkCall allows ROUND and nothing else.
func checkCall(call *ast.CallNode) error {
	name, named := calleeName(call)
	if !named {
		return fmt.Errorf("%w: computed function call", ErrFormulaUnsupported)
	}

	if name != roundFunction {
		return arithmeticFunctionError(name)
	}

	if len(call.Arguments) != 2 {
		return fmt.Errorf("%w: ROUND takes 2 arguments, got %d", ErrBadCallArity, len(call.Arguments))
	}
	if err := checkArithmetic(call.Arguments[0]); err != nil {
		return err
	}
	if _, err := integerLiteral(call.Arguments[1]); err != nil {
		return err
	}
	return nil
}

// integerLiteral reads ROUND's places argument. The parser represents a
// negative literal as unary minus applied to a positive one, so sign nodes are
// unwrapped here.
func integerLiteral(node ast.Node) (int32, error) {
	switch n := node.(type) {
	case *ast.IntegerNode:
		if n.Value > roundPlacesLimit || n.Value < -roundPlacesLimit {
			return 0, fmt.Errorf("%w: %d is outside [-%d, %d]",
				ErrBadRoundPlaces, n.Value, roundPlacesLimit, roundPlacesLimit)
		}
		return int32(n.Value), nil

	case *ast.UnaryNode:
		if n.Operator != "-" && n.Operator != "+" {
			break
		}
		inner, err := integerLiteral(n.Node)
		if err != nil {
			return 0, err
		}
		if n.Operator == "-" {
			return -inner, nil
		}
		return inner, nil
	}

	return 0, fmt.Errorf("%w: got %s", ErrBadRoundPlaces, nodeKind(node))
}

// arithmeticFunctionError explains why a function is not allowed in an
// arithmetic column. Naming a vertical reducer, or mis-casing ROUND, are the two
// mistakes worth their own message.
func arithmeticFunctionError(name string) error {
	upper := strings.ToUpper(name)

	if _, isAggregate := aggFuncFromName(upper); isAggregate {
		return fmt.Errorf("%w: %s reduces rows and may only define an aggregate column",
			ErrUnknownFunction, name)
	}
	if upper == roundFunction && name != roundFunction {
		return fmt.Errorf("%w: %s (function names are upper case: did you mean %s?)",
			ErrUnknownFunction, name, roundFunction)
	}
	return fmt.Errorf("%w: %s", ErrUnknownFunction, name)
}

// aggregateFunctionError explains why a function cannot define an aggregate
// column.
func aggregateFunctionError(name string) error {
	upper := strings.ToUpper(name)

	if _, isAggregate := aggFuncFromName(upper); isAggregate && upper != name {
		return fmt.Errorf("%w: %s (function names are upper case: did you mean %s?)",
			ErrUnknownFunction, name, upper)
	}
	if upper == roundFunction {
		return fmt.Errorf("%w: %s does not reduce rows; round an arithmetic column instead",
			ErrUnknownFunction, name)
	}
	return fmt.Errorf("%w: %s", ErrUnknownFunction, name)
}

func isArithmeticOperator(operator string) bool {
	switch operator {
	case "+", "-", "*", "/":
		return true
	}
	return false
}

func calleeName(call *ast.CallNode) (string, bool) {
	identifier, isIdentifier := call.Callee.(*ast.IdentifierNode)
	if !isIdentifier {
		return "", false
	}
	return identifier.Value, true
}

// nodeKind names a syntax node for an error message.
func nodeKind(node ast.Node) string {
	switch node.(type) {
	case *ast.IntegerNode, *ast.FloatNode:
		return "number literal"
	case *ast.IdentifierNode:
		return "column reference"
	case *ast.StringNode:
		return "string literal"
	case *ast.BoolNode:
		return "boolean literal"
	case *ast.NilNode:
		return "nil"
	case *ast.BytesNode:
		return "bytes literal"
	case *ast.ConstantNode:
		return "constant"
	case *ast.ArrayNode:
		return "array literal"
	case *ast.MapNode, *ast.PairNode:
		return "map literal"
	case *ast.MemberNode:
		return "property access"
	case *ast.SliceNode:
		return "slice"
	case *ast.BuiltinNode:
		return "builtin function"
	case *ast.ConditionalNode:
		return "conditional expression"
	case *ast.ChainNode:
		return "optional chaining"
	case *ast.PredicateNode:
		return "predicate"
	case *ast.PointerNode:
		return "pointer"
	case *ast.VariableDeclaratorNode:
		return "variable declaration"
	case *ast.SequenceNode:
		return "expression sequence"
	case *ast.CallNode:
		return "function call"
	case *ast.UnaryNode:
		return "unary expression"
	case *ast.BinaryNode:
		return "binary expression"
	}
	return fmt.Sprintf("%T", node)
}
