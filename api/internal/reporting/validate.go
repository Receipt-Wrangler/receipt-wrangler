package reporting

import (
	"errors"
	"fmt"
	"strings"

	"github.com/expr-lang/expr/ast"
)

var (
	ErrNoColumns          = errors.New("report spec has no columns")
	ErrEmptyColumnName    = errors.New("column name must not be empty")
	ErrInvalidColumnName  = errors.New("column name must be a plain identifier")
	ErrReservedColumnName = errors.New("column name is a reserved word")
	ErrDuplicateColumn    = errors.New("duplicate column name")
	ErrUnknownColumnKind  = errors.New("unknown column kind")

	ErrUnknownField        = errors.New("field not found in catalog")
	ErrGroupByNotDimension = errors.New("groupBy field is not a dimension")
	ErrDuplicateGroupBy    = errors.New("duplicate groupBy field")

	ErrDetailByRequired     = errors.New("aggregate detail mode requires a dimension")
	ErrDetailByOnRecords    = errors.New("records detail mode must not set a dimension")
	ErrDetailByNotDimension = errors.New("detail field is not a dimension")

	ErrLabelFieldRequired      = errors.New("label column requires a field")
	ErrLabelColumnUnresolvable = errors.New("label column has no value on an aggregated detail row")

	ErrAggregateFieldRequired = errors.New("aggregate requires a field")
	ErrAggregateNotMeasure    = errors.New("aggregate field is not a measure")
	ErrCountTakesNoField      = errors.New("COUNT counts records and takes no field")

	ErrUnknownColumnRef        = errors.New("formula references an unknown column")
	ErrArithmeticNonNumericRef = errors.New("formula references a non-numeric column")
	ErrFormulaCycle            = errors.New("arithmetic columns form a cycle")

	ErrInvalidConfig = errors.New("invalid engine config")
)

// maxDivisionScale bounds the precision division may be carried to, for the same
// reason ROUND's places are bounded: rescaling a decimal without limit lets a
// saved template allocate without limit.
const maxDivisionScale = 30

// labelFromDetailBucket marks a label column that reads the aggregated detail
// row's own bucket value rather than an ancestor group's.
const labelFromDetailBucket = -1

// compiledColumn is a validated column with everything the engine needs already
// resolved: its data type, its parsed expression, and where a label reads from.
type compiledColumn struct {
	name     string
	label    string
	kind     ColumnKind
	dataType DataType
	format   string

	// field and fieldRef describe a label column's source.
	field    FieldKey
	fieldRef FieldRef

	// labelLevel says where a label column reads from when detail rows are
	// aggregated: labelFromDetailBucket for the row's own bucket, otherwise the
	// index of the groupBy level whose bucket value it shows. It is unused in
	// records mode, where a label reads straight off the record.
	labelLevel int

	// agg is an aggregate column's reduction.
	agg Aggregate

	// expr is an arithmetic column's parsed expression, and refs the column
	// names it reads.
	expr    ast.Node
	exprSrc string
	refs    []string
}

// compiledSpec is a ReportSpec that has been validated against a catalog.
type compiledSpec struct {
	spec    ReportSpec
	catalog FieldCatalog

	columns []compiledColumn

	// groupBy is the resolved dimension of each nesting level.
	groupBy []FieldRef

	// detailBy is the dimension aggregated detail rows are keyed on.
	detailBy FieldRef

	// aggregateIndexes lists the aggregate columns, and arithmeticOrder lists
	// the arithmetic columns in an order where a column's dependencies always
	// come first.
	aggregateIndexes []int
	arithmeticOrder  []int
}

// Validate reports whether a spec can run against a catalog. Run calls it, so
// callers only need it to check a template as it is being authored.
func Validate(spec ReportSpec, catalog FieldCatalog) error {
	_, err := compileSpec(spec, catalog)
	return err
}

// compileSpec validates a spec and resolves everything the engine would
// otherwise have to work out per row: field references, parsed expressions,
// column data types, and the order arithmetic columns must be evaluated in.
func compileSpec(spec ReportSpec, catalog FieldCatalog) (compiledSpec, error) {
	if err := validateConfig(spec.Config); err != nil {
		return compiledSpec{}, err
	}
	if len(spec.Columns) == 0 {
		return compiledSpec{}, ErrNoColumns
	}

	compiled := compiledSpec{spec: spec.withDefaults(), catalog: catalog}

	groupBy, err := compileGroupBy(spec.GroupBy, catalog)
	if err != nil {
		return compiledSpec{}, err
	}
	compiled.groupBy = groupBy

	detailBy, err := compileDetail(spec.Detail, catalog)
	if err != nil {
		return compiledSpec{}, err
	}
	compiled.detailBy = detailBy

	columns, byName, err := compileColumns(spec, catalog)
	if err != nil {
		return compiledSpec{}, err
	}
	compiled.columns = columns

	for index, column := range compiled.columns {
		if column.kind == ColumnAggregate {
			compiled.aggregateIndexes = append(compiled.aggregateIndexes, index)
		}
	}

	// Arithmetic columns may read each other, so they are checked and ordered
	// as a graph rather than in the order they were declared.
	order, err := orderArithmetic(compiled.columns, byName)
	if err != nil {
		return compiledSpec{}, err
	}
	compiled.arithmeticOrder = order

	// A column's type can depend on the columns it reads, so arithmetic types
	// are inferred in dependency order.
	for _, index := range order {
		compiled.columns[index].dataType = arithmeticDataType(compiled.columns[index], compiled.columns, byName)
	}

	return compiled, nil
}

func validateConfig(config EngineConfig) error {
	if config.DivisionScale < 0 || config.DivisionScale > maxDivisionScale {
		return fmt.Errorf("%w: division scale %d is outside [0, %d]",
			ErrInvalidConfig, config.DivisionScale, maxDivisionScale)
	}
	return nil
}

func compileGroupBy(keys []FieldKey, catalog FieldCatalog) ([]FieldRef, error) {
	seen := make(map[FieldKey]struct{}, len(keys))
	refs := make([]FieldRef, 0, len(keys))

	for _, key := range keys {
		field, exists := catalog.Get(key)
		if !exists {
			return nil, fmt.Errorf("%w: groupBy %s", ErrUnknownField, key)
		}
		if field.Role() != RoleDimension {
			return nil, fmt.Errorf("%w: %s is a %s", ErrGroupByNotDimension, key, field.DataType)
		}
		if _, duplicate := seen[key]; duplicate {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateGroupBy, key)
		}
		seen[key] = struct{}{}
		refs = append(refs, field)
	}

	return refs, nil
}

func compileDetail(detail DetailSpec, catalog FieldCatalog) (FieldRef, error) {
	if detail.Mode == DetailRecords {
		if len(detail.By) > 0 {
			return FieldRef{}, fmt.Errorf("%w: %s", ErrDetailByOnRecords, detail.By)
		}
		return FieldRef{}, nil
	}

	if len(detail.By) == 0 {
		return FieldRef{}, ErrDetailByRequired
	}

	field, exists := catalog.Get(detail.By)
	if !exists {
		return FieldRef{}, fmt.Errorf("%w: detail %s", ErrUnknownField, detail.By)
	}
	if field.Role() != RoleDimension {
		return FieldRef{}, fmt.Errorf("%w: %s is a %s", ErrDetailByNotDimension, detail.By, field.DataType)
	}

	return field, nil
}

func compileColumns(spec ReportSpec, catalog FieldCatalog) ([]compiledColumn, map[string]int, error) {
	columns := make([]compiledColumn, 0, len(spec.Columns))
	byName := make(map[string]int, len(spec.Columns))

	for _, column := range spec.Columns {
		if err := validateColumnName(column.Name); err != nil {
			return nil, nil, err
		}
		if _, duplicate := byName[column.Name]; duplicate {
			return nil, nil, fmt.Errorf("%w: %s", ErrDuplicateColumn, column.Name)
		}

		compiled, err := compileColumn(column, spec, catalog)
		if err != nil {
			return nil, nil, err
		}

		byName[column.Name] = len(columns)
		columns = append(columns, compiled)
	}

	// Resolving arithmetic references needs every column to be known, so it
	// happens once the whole set is compiled.
	for _, column := range columns {
		if column.kind != ColumnArithmetic {
			continue
		}
		if err := validateArithmeticRefs(column, columns, byName); err != nil {
			return nil, nil, err
		}
	}

	return columns, byName, nil
}

func compileColumn(column Column, spec ReportSpec, catalog FieldCatalog) (compiledColumn, error) {
	compiled := compiledColumn{
		name:   column.Name,
		label:  column.heading(),
		kind:   column.Kind,
		format: column.Format,
	}

	switch column.Kind {
	case ColumnLabel:
		return compileLabelColumn(compiled, column, spec, catalog)
	case ColumnAggregate:
		return compileAggregateColumn(compiled, column, catalog)
	case ColumnArithmetic:
		return compileArithmeticColumn(compiled, column)
	}

	return compiledColumn{}, fmt.Errorf("%w: column %s has kind %d", ErrUnknownColumnKind, column.Name, column.Kind)
}

func compileLabelColumn(compiled compiledColumn, column Column, spec ReportSpec, catalog FieldCatalog) (compiledColumn, error) {
	if len(column.Field) == 0 {
		return compiledColumn{}, fmt.Errorf("%w: column %s", ErrLabelFieldRequired, column.Name)
	}

	field, exists := catalog.Get(column.Field)
	if !exists {
		return compiledColumn{}, fmt.Errorf("%w: column %s reads %s", ErrUnknownField, column.Name, column.Field)
	}

	compiled.field = column.Field
	compiled.fieldRef = field
	compiled.dataType = field.DataType
	compiled.labelLevel = labelFromDetailBucket

	// In records mode a label reads straight off the record, so any field will
	// do. Aggregated detail rows have no single record to read from: the only
	// values that exist are the row's own bucket and the buckets of the groups
	// above it.
	if spec.Detail.Mode == DetailAggregate && column.Field != spec.Detail.By {
		level := groupByLevel(spec.GroupBy, column.Field)
		if level < 0 {
			return compiledColumn{}, fmt.Errorf(
				"%w: column %s reads %s, which is neither the detail dimension nor a groupBy level",
				ErrLabelColumnUnresolvable, column.Name, column.Field)
		}
		compiled.labelLevel = level
	}

	return compiled, nil
}

func compileAggregateColumn(compiled compiledColumn, column Column, catalog FieldCatalog) (compiledColumn, error) {
	aggregate := column.Agg
	if len(column.AggSrc) > 0 {
		parsed, err := ParseAggregate(column.AggSrc)
		if err != nil {
			return compiledColumn{}, fmt.Errorf("column %s: %w", column.Name, err)
		}
		aggregate = parsed
	}

	if aggregate.Func == AggCount {
		if len(aggregate.Field) > 0 {
			return compiledColumn{}, fmt.Errorf("%w: column %s counts %s",
				ErrCountTakesNoField, column.Name, aggregate.Field)
		}
		compiled.agg = aggregate
		compiled.dataType = TypeNumber
		return compiled, nil
	}

	if len(aggregate.Field) == 0 {
		return compiledColumn{}, fmt.Errorf("%w: column %s applies %s",
			ErrAggregateFieldRequired, column.Name, aggregate.Func)
	}

	field, exists := catalog.Get(aggregate.Field)
	if !exists {
		return compiledColumn{}, fmt.Errorf("%w: column %s reduces %s",
			ErrUnknownField, column.Name, aggregate.Field)
	}
	if field.Role() != RoleMeasure {
		return compiledColumn{}, fmt.Errorf("%w: column %s applies %s to %s, which is a %s",
			ErrAggregateNotMeasure, column.Name, aggregate.Func, aggregate.Field, field.DataType)
	}

	compiled.agg = aggregate
	compiled.dataType = field.DataType

	return compiled, nil
}

func compileArithmeticColumn(compiled compiledColumn, column Column) (compiledColumn, error) {
	expr, err := ParseArithmetic(column.Expr)
	if err != nil {
		return compiledColumn{}, fmt.Errorf("column %s: %w", column.Name, err)
	}

	compiled.expr = expr
	compiled.exprSrc = column.Expr
	compiled.refs = columnRefs(expr)
	// A placeholder until the dependency order is known; see compileSpec.
	compiled.dataType = TypeNumber

	return compiled, nil
}

// validateArithmeticRefs checks that a formula reads only declared columns that
// actually produce numbers.
func validateArithmeticRefs(column compiledColumn, columns []compiledColumn, byName map[string]int) error {
	for _, ref := range column.refs {
		index, declared := byName[ref]
		if !declared {
			return fmt.Errorf("%w: column %s reads %s", ErrUnknownColumnRef, column.name, ref)
		}

		referenced := columns[index]
		if referenced.kind == ColumnLabel && !referenced.fieldRef.DataType.IsNumeric() {
			return fmt.Errorf("%w: column %s reads %s, which shows a %s",
				ErrArithmeticNonNumericRef, column.name, ref, referenced.fieldRef.DataType)
		}
	}
	return nil
}

// orderArithmetic topologically sorts the arithmetic columns so that a column's
// dependencies are always evaluated before it, and rejects cycles.
//
// Columns are visited in declaration order and dependencies depth-first, which
// makes the resulting order deterministic. A column may therefore be declared
// before the aggregate it reads.
func orderArithmetic(columns []compiledColumn, byName map[string]int) ([]int, error) {
	const (
		unvisited = iota
		visiting
		visited
	)

	state := make([]int, len(columns))
	order := make([]int, 0, len(columns))

	// path tracks the columns currently being visited, so a cycle can name its
	// members rather than merely announcing that one exists.
	var path []string

	var visit func(index int) error
	visit = func(index int) error {
		switch state[index] {
		case visited:
			return nil
		case visiting:
			return fmt.Errorf("%w: %s", ErrFormulaCycle, strings.Join(append(path, columns[index].name), " -> "))
		}

		state[index] = visiting
		path = append(path, columns[index].name)

		for _, ref := range columns[index].refs {
			dependency, declared := byName[ref]
			if !declared {
				// validateArithmeticRefs already rejected this.
				continue
			}
			if columns[dependency].kind != ColumnArithmetic {
				continue
			}
			if err := visit(dependency); err != nil {
				return err
			}
		}

		path = path[:len(path)-1]
		state[index] = visited
		order = append(order, index)

		return nil
	}

	for index, column := range columns {
		if column.kind != ColumnArithmetic {
			continue
		}
		if err := visit(index); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// arithmeticDataType infers what an arithmetic column holds. Money combined
// with anything is still money, so the column is currency when any column it
// reads is. Dependencies are already typed, because this runs in dependency
// order.
func arithmeticDataType(column compiledColumn, columns []compiledColumn, byName map[string]int) DataType {
	for _, ref := range column.refs {
		index, declared := byName[ref]
		if !declared {
			continue
		}
		if columns[index].dataType == TypeCurrency {
			return TypeCurrency
		}
	}
	return TypeNumber
}

func validateColumnName(name string) error {
	if len(name) == 0 {
		return ErrEmptyColumnName
	}
	if !isIdentifier(name) {
		return fmt.Errorf("%w: %q", ErrInvalidColumnName, name)
	}
	if isReservedWord(name) {
		return fmt.Errorf("%w: %q", ErrReservedColumnName, name)
	}
	return nil
}

// isIdentifier reports whether a name matches [A-Za-z_][A-Za-z0-9_]*, which is
// what a formula can reference without quoting.
func isIdentifier(name string) bool {
	for index, symbol := range name {
		switch {
		case symbol == '_':
		case symbol >= 'a' && symbol <= 'z':
		case symbol >= 'A' && symbol <= 'Z':
		case symbol >= '0' && symbol <= '9' && index > 0:
		default:
			return false
		}
	}
	return len(name) > 0
}

// groupByLevel returns the nesting level a dimension is grouped at, or -1.
func groupByLevel(groupBy []FieldKey, key FieldKey) int {
	for level, candidate := range groupBy {
		if candidate == key {
			return level
		}
	}
	return -1
}
