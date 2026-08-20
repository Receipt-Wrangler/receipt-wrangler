package reporting

import (
	"errors"
	"testing"
)

// testCatalog mirrors the shape a receipt source produces: currency measures, a
// plain numeric measure, multi-value and single-value string dimensions, a date
// and a boolean.
func testCatalog(t *testing.T) FieldCatalog {
	t.Helper()

	catalog, err := NewFieldCatalog(
		FieldRef{Key: "amount", Label: "Amount", DataType: TypeCurrency},
		FieldRef{Key: "custom_1", Label: "HST", DataType: TypeCurrency},
		FieldRef{Key: "quantity", Label: "Quantity", DataType: TypeNumber},
		FieldRef{Key: "category", Label: "Category", DataType: TypeString, Multi: true},
		FieldRef{Key: "tag", Label: "Tag", DataType: TypeString, Multi: true},
		FieldRef{Key: "paid_by", Label: "Paid By", DataType: TypeString},
		FieldRef{Key: "date", Label: "Date", DataType: TypeDate},
		FieldRef{Key: "resolved", Label: "Resolved", DataType: TypeBool},
		// No receipt field is a multi-valued number, but nothing in the engine
		// stops a producer from declaring one, so the rules that reject
		// measuring it need something to reject.
		FieldRef{Key: "item_amounts", Label: "Item Amounts", DataType: TypeCurrency, Multi: true},
		FieldRef{Key: "item_counts", Label: "Item Counts", DataType: TypeNumber, Multi: true},
	)
	if err != nil {
		t.Fatalf("NewFieldCatalog() error = %v", err)
	}
	return catalog
}

// workedExampleSpec is the design document's report: foster parent, then child,
// with one summed row per category.
func workedExampleSpec() ReportSpec {
	return ReportSpec{
		Title:   "Verified Expenses",
		GroupBy: []FieldKey{"paid_by", "tag"},
		Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns: []Column{
			{Name: "Category", Kind: ColumnLabel, Field: "category"},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Hst", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
			{Name: "Total", Kind: ColumnArithmetic, Expr: "Subtotal + Hst"},
			{Name: "AvgPerReceipt", Label: "Avg/Receipt", Kind: ColumnArithmetic, Expr: "Total / Count"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}
}

func TestValidate_AcceptsTheWorkedExample(t *testing.T) {
	if err := Validate(workedExampleSpec(), testCatalog(t)); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidate_AcceptsVariations(t *testing.T) {
	tests := []struct {
		name string
		spec ReportSpec
	}{
		{
			name: "records mode with no grouping",
			spec: ReportSpec{
				Columns: []Column{{Name: "Name", Kind: ColumnLabel, Field: "paid_by"}},
			},
		},
		{
			name: "records mode may label any field",
			spec: ReportSpec{
				GroupBy: []FieldKey{"paid_by"},
				Columns: []Column{
					{Name: "When", Kind: ColumnLabel, Field: "date"},
					{Name: "Cats", Kind: ColumnLabel, Field: "category"},
					{Name: "Amount", Kind: ColumnLabel, Field: "amount"},
				},
			},
		},
		{
			name: "aggregate detail may label a groupBy level",
			spec: ReportSpec{
				GroupBy: []FieldKey{"paid_by"},
				Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
				Columns: []Column{
					{Name: "PaidBy", Kind: ColumnLabel, Field: "paid_by"},
					{Name: "Category", Kind: ColumnLabel, Field: "category"},
					{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
				},
			},
		},
		{
			name: "aggregate source form",
			spec: ReportSpec{
				Columns: []Column{
					{Name: "Count", Kind: ColumnAggregate, AggSrc: "COUNT()"},
					{Name: "Subtotal", Kind: ColumnAggregate, AggSrc: "SUM(amount)"},
					{Name: "Biggest", Kind: ColumnAggregate, AggSrc: "MAX(amount)"},
				},
			},
		},
		{
			name: "arithmetic may be declared before what it reads",
			spec: ReportSpec{
				Columns: []Column{
					{Name: "Avg", Kind: ColumnArithmetic, Expr: "Total / Count"},
					{Name: "Total", Kind: ColumnArithmetic, Expr: "Subtotal + Hst"},
					{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
					{Name: "Hst", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
					{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
				},
			},
		},
		{
			name: "arithmetic may read a numeric label column",
			spec: ReportSpec{
				Columns: []Column{
					{Name: "Qty", Kind: ColumnLabel, Field: "quantity"},
					{Name: "Doubled", Kind: ColumnArithmetic, Expr: "Qty * 2"},
				},
			},
		},
		{
			name: "constant arithmetic reads nothing",
			spec: ReportSpec{
				Columns: []Column{{Name: "Fixed", Kind: ColumnArithmetic, Expr: "1 + 1"}},
			},
		},
		{
			// Measuring a multi-valued field is refused; cutting by one, or
			// showing every value it holds, is not.
			name: "a multi-valued field may be grouped on and displayed",
			spec: ReportSpec{
				GroupBy: []FieldKey{"tag"},
				Columns: []Column{
					{Name: "Cats", Kind: ColumnLabel, Field: "category"},
					{Name: "Amounts", Kind: ColumnLabel, Field: "item_amounts"},
					{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
				},
			},
		},
		{
			// A currency custom field is a measure, but measuring is the only
			// thing its type restricts: it still cuts, so a report can group by
			// a tax or tip field the same way it groups by a category.
			name: "a measure may be grouped on and aggregated by",
			spec: ReportSpec{
				GroupBy: []FieldKey{"custom_1"},
				Detail:  DetailSpec{Mode: DetailAggregate, By: "amount"},
				Columns: []Column{
					{Name: "Hst", Kind: ColumnLabel, Field: "custom_1"},
					{Name: "Amount", Kind: ColumnLabel, Field: "amount"},
					{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
				},
			},
		},
		{
			name: "a date or boolean may be grouped on",
			spec: ReportSpec{
				GroupBy: []FieldKey{"date", "resolved"},
				Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
			},
		},
		{
			name: "COUNT is unaffected by multi-valued fields",
			spec: ReportSpec{
				GroupBy: []FieldKey{"category"},
				Columns: []Column{{Name: "Count", Kind: ColumnAggregate, AggSrc: "COUNT()"}},
			},
		},
		{
			name: "division scale at the bounds",
			spec: ReportSpec{
				Config:  EngineConfig{DivisionScale: maxDivisionScale},
				Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := Validate(test.spec, testCatalog(t)); err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestValidate_Rejects(t *testing.T) {
	countColumn := Column{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}

	tests := []struct {
		name    string
		spec    ReportSpec
		wantErr error
	}{
		{
			name:    "no columns",
			spec:    ReportSpec{},
			wantErr: ErrNoColumns,
		},
		{
			name:    "negative division scale",
			spec:    ReportSpec{Config: EngineConfig{DivisionScale: -1}, Columns: []Column{countColumn}},
			wantErr: ErrInvalidConfig,
		},
		{
			name:    "division scale beyond the limit",
			spec:    ReportSpec{Config: EngineConfig{DivisionScale: maxDivisionScale + 1}, Columns: []Column{countColumn}},
			wantErr: ErrInvalidConfig,
		},
		{
			name:    "empty column name",
			spec:    ReportSpec{Columns: []Column{{Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrEmptyColumnName,
		},
		{
			name:    "column name with a slash",
			spec:    ReportSpec{Columns: []Column{{Name: "Avg/Receipt", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrInvalidColumnName,
		},
		{
			name:    "column name starting with a digit",
			spec:    ReportSpec{Columns: []Column{{Name: "1st", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrInvalidColumnName,
		},
		{
			name:    "column name with a space",
			spec:    ReportSpec{Columns: []Column{{Name: "Paid By", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrInvalidColumnName,
		},
		{
			name:    "column name with an accent",
			spec:    ReportSpec{Columns: []Column{{Name: "Montée", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrInvalidColumnName,
		},
		{
			// "in" is an operator in the expression language, so a formula could
			// never reference a column of that name.
			name:    "column named after an operator",
			spec:    ReportSpec{Columns: []Column{{Name: "in", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrReservedColumnName,
		},
		{
			name:    "column named after an aggregate function",
			spec:    ReportSpec{Columns: []Column{{Name: "SUM", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrReservedColumnName,
		},
		{
			name:    "column named ROUND",
			spec:    ReportSpec{Columns: []Column{{Name: "ROUND", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}},
			wantErr: ErrReservedColumnName,
		},
		{
			name: "duplicate column name",
			spec: ReportSpec{Columns: []Column{
				countColumn,
				{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			}},
			wantErr: ErrDuplicateColumn,
		},
		{
			name:    "unknown column kind",
			spec:    ReportSpec{Columns: []Column{{Name: "Odd", Kind: ColumnKind(9)}}},
			wantErr: ErrUnknownColumnKind,
		},

		// groupBy
		{
			name:    "groupBy an unknown field",
			spec:    ReportSpec{GroupBy: []FieldKey{"nope"}, Columns: []Column{countColumn}},
			wantErr: ErrUnknownField,
		},
		{
			name:    "duplicate groupBy level",
			spec:    ReportSpec{GroupBy: []FieldKey{"tag", "tag"}, Columns: []Column{countColumn}},
			wantErr: ErrDuplicateGroupBy,
		},

		// detail
		{
			name:    "aggregate detail without a dimension",
			spec:    ReportSpec{Detail: DetailSpec{Mode: DetailAggregate}, Columns: []Column{countColumn}},
			wantErr: ErrDetailByRequired,
		},
		{
			name:    "records detail with a dimension",
			spec:    ReportSpec{Detail: DetailSpec{Mode: DetailRecords, By: "category"}, Columns: []Column{countColumn}},
			wantErr: ErrDetailByOnRecords,
		},
		{
			name:    "aggregate detail on an unknown field",
			spec:    ReportSpec{Detail: DetailSpec{Mode: DetailAggregate, By: "nope"}, Columns: []Column{countColumn}},
			wantErr: ErrUnknownField,
		},

		// label columns
		{
			name:    "label column without a field",
			spec:    ReportSpec{Columns: []Column{{Name: "Blank", Kind: ColumnLabel}}},
			wantErr: ErrLabelFieldRequired,
		},
		{
			name:    "label column on an unknown field",
			spec:    ReportSpec{Columns: []Column{{Name: "Ghost", Kind: ColumnLabel, Field: "nope"}}},
			wantErr: ErrUnknownField,
		},
		{
			// An aggregated detail row sums many records, so there is no single
			// record whose date it could show.
			name: "label column unresolvable on an aggregated detail row",
			spec: ReportSpec{
				GroupBy: []FieldKey{"paid_by"},
				Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
				Columns: []Column{
					{Name: "When", Kind: ColumnLabel, Field: "date"},
					countColumn,
				},
			},
			wantErr: ErrLabelColumnUnresolvable,
		},

		// aggregate columns
		{
			name:    "aggregate without a field",
			spec:    ReportSpec{Columns: []Column{{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum}}}},
			wantErr: ErrAggregateFieldRequired,
		},
		{
			name:    "aggregate over an unknown field",
			spec:    ReportSpec{Columns: []Column{{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "nope"}}}},
			wantErr: ErrUnknownField,
		},
		{
			name:    "sum a dimension",
			spec:    ReportSpec{Columns: []Column{{Name: "Nonsense", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "tag"}}}},
			wantErr: ErrAggregateNotMeasure,
		},
		{
			name:    "average a date",
			spec:    ReportSpec{Columns: []Column{{Name: "Nonsense", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: "date"}}}},
			wantErr: ErrAggregateNotMeasure,
		},
		{
			name:    "count a field",
			spec:    ReportSpec{Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount, Field: "amount"}}}},
			wantErr: ErrCountTakesNoField,
		},
		{
			// Reading only the first of several amounts would silently lose money.
			name:    "sum a multi-valued measure",
			spec:    ReportSpec{Columns: []Column{{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "item_amounts"}}}},
			wantErr: ErrMeasureIsMultiValued,
		},
		{
			name:    "average a multi-valued measure",
			spec:    ReportSpec{Columns: []Column{{Name: "Mean", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: "item_counts"}}}},
			wantErr: ErrMeasureIsMultiValued,
		},
		{
			name:    "take the maximum of a multi-valued measure",
			spec:    ReportSpec{Columns: []Column{{Name: "Most", Kind: ColumnAggregate, AggSrc: "MAX(item_amounts)"}}},
			wantErr: ErrMeasureIsMultiValued,
		},
		// ParseAggregate cannot produce a reduction the accumulator does not
		// know, but Column.Agg is the canonical form and a persisted template
		// deserializes straight into it. An unknown reduction used to validate
		// cleanly and then render every cell of the column empty, at every
		// level, because finalize's switch fell through to null.
		{
			name:    "aggregate with a reduction the accumulator does not know",
			spec:    ReportSpec{Columns: []Column{{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggFunc(200), Field: "amount"}}}},
			wantErr: ErrUnknownAggFunc,
		},
		{
			name:    "aggregate one past the last known reduction",
			spec:    ReportSpec{Columns: []Column{{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMax + 1, Field: "amount"}}}},
			wantErr: ErrUnknownAggFunc,
		},
		// compileDetail read anything that was not DetailRecords as aggregate,
		// while compileLabelColumn tested for DetailAggregate exactly. An
		// unknown mode therefore took the aggregate path everywhere but the
		// label check, which it skipped — so a label column over any field at
		// all silently rendered the detail bucket's value instead of its own.
		{
			name: "detail mode the engine does not know",
			spec: ReportSpec{
				Detail:  DetailSpec{Mode: DetailMode(200), By: "category"},
				Columns: []Column{countColumn},
			},
			wantErr: ErrUnknownDetailMode,
		},
		{
			name: "detail mode one past the last known mode",
			spec: ReportSpec{
				Detail:  DetailSpec{Mode: DetailAggregate + 1, By: "category"},
				Columns: []Column{countColumn},
			},
			wantErr: ErrUnknownDetailMode,
		},
		{
			// A label over a multi-valued field shows every value, so a formula
			// cannot silently pick one.
			name: "arithmetic reads a multi-valued numeric label column",
			spec: ReportSpec{
				Columns: []Column{
					{Name: "Amounts", Kind: ColumnLabel, Field: "item_amounts"},
					{Name: "Bad", Kind: ColumnArithmetic, Expr: "Amounts * 2"},
				},
			},
			wantErr: ErrMeasureIsMultiValued,
		},
		{
			name:    "aggregate source is not a call",
			spec:    ReportSpec{Columns: []Column{{Name: "Subtotal", Kind: ColumnAggregate, AggSrc: "amount"}}},
			wantErr: ErrFormulaUnsupported,
		},
		{
			name:    "aggregate source does not parse",
			spec:    ReportSpec{Columns: []Column{{Name: "Subtotal", Kind: ColumnAggregate, AggSrc: "SUM("}}},
			wantErr: ErrFormulaSyntax,
		},
		{
			name:    "aggregate source names an unknown function",
			spec:    ReportSpec{Columns: []Column{{Name: "Subtotal", Kind: ColumnAggregate, AggSrc: "TOTAL(amount)"}}},
			wantErr: ErrUnknownFunction,
		},

		// arithmetic columns
		{
			name:    "arithmetic does not parse",
			spec:    ReportSpec{Columns: []Column{countColumn, {Name: "Bad", Kind: ColumnArithmetic, Expr: "Count +"}}},
			wantErr: ErrFormulaSyntax,
		},
		{
			name:    "arithmetic uses an unsupported operator",
			spec:    ReportSpec{Columns: []Column{countColumn, {Name: "Bad", Kind: ColumnArithmetic, Expr: "Count % 2"}}},
			wantErr: ErrFormulaUnsupported,
		},
		{
			name:    "arithmetic calls a reducer",
			spec:    ReportSpec{Columns: []Column{{Name: "Bad", Kind: ColumnArithmetic, Expr: "SUM(amount)"}}},
			wantErr: ErrUnknownFunction,
		},
		{
			name:    "arithmetic rounds to too many places",
			spec:    ReportSpec{Columns: []Column{countColumn, {Name: "Bad", Kind: ColumnArithmetic, Expr: "ROUND(Count, 99)"}}},
			wantErr: ErrBadRoundPlaces,
		},
		{
			name:    "arithmetic reads an undeclared column",
			spec:    ReportSpec{Columns: []Column{countColumn, {Name: "Bad", Kind: ColumnArithmetic, Expr: "Count + Missing"}}},
			wantErr: ErrUnknownColumnRef,
		},
		{
			// A formula cannot add a category name to a number.
			name: "arithmetic reads a non-numeric label column",
			spec: ReportSpec{
				Columns: []Column{
					{Name: "Cat", Kind: ColumnLabel, Field: "category"},
					countColumn,
					{Name: "Bad", Kind: ColumnArithmetic, Expr: "Count + Cat"},
				},
			},
			wantErr: ErrArithmeticNonNumericRef,
		},
		{
			name: "arithmetic reads a date label column",
			spec: ReportSpec{
				Columns: []Column{
					{Name: "When", Kind: ColumnLabel, Field: "date"},
					{Name: "Bad", Kind: ColumnArithmetic, Expr: "When * 2"},
				},
			},
			wantErr: ErrArithmeticNonNumericRef,
		},
		{
			name:    "arithmetic references itself",
			spec:    ReportSpec{Columns: []Column{{Name: "Loop", Kind: ColumnArithmetic, Expr: "Loop + 1"}}},
			wantErr: ErrFormulaCycle,
		},
		{
			name: "two arithmetic columns reference each other",
			spec: ReportSpec{Columns: []Column{
				{Name: "A", Kind: ColumnArithmetic, Expr: "B + 1"},
				{Name: "B", Kind: ColumnArithmetic, Expr: "A + 1"},
			}},
			wantErr: ErrFormulaCycle,
		},
		{
			name: "a longer cycle",
			spec: ReportSpec{Columns: []Column{
				{Name: "A", Kind: ColumnArithmetic, Expr: "B"},
				{Name: "B", Kind: ColumnArithmetic, Expr: "C"},
				{Name: "C", Kind: ColumnArithmetic, Expr: "A"},
			}},
			wantErr: ErrFormulaCycle,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.spec, testCatalog(t))
			if !errors.Is(err, test.wantErr) {
				t.Errorf("Validate() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

// A diamond is not a cycle: two columns may read the same dependency.
func TestValidate_DiamondDependencyIsNotACycle(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Base", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Left", Kind: ColumnArithmetic, Expr: "Base * 2"},
			{Name: "Right", Kind: ColumnArithmetic, Expr: "Base * 3"},
			{Name: "Joined", Kind: ColumnArithmetic, Expr: "Left + Right"},
		},
	}

	if err := Validate(spec, testCatalog(t)); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

// The engine evaluates arithmetic in this order, so a column must never appear
// before something it reads.
func TestCompileSpec_ArithmeticOrderIsTopological(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "AvgPerReceipt", Kind: ColumnArithmetic, Expr: "Total / Count"},
			{Name: "Total", Kind: ColumnArithmetic, Expr: "Subtotal + Hst"},
			{Name: "Subtotal", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Hst", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "custom_1"}},
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
		},
	}

	compiled, err := compileSpec(spec, testCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}

	position := make(map[string]int, len(compiled.arithmeticOrder))
	for slot, index := range compiled.arithmeticOrder {
		position[compiled.columns[index].name] = slot
	}

	if len(compiled.arithmeticOrder) != 2 {
		t.Fatalf("arithmeticOrder has %d entries, want 2", len(compiled.arithmeticOrder))
	}
	if position["Total"] > position["AvgPerReceipt"] {
		t.Errorf("Total is evaluated after AvgPerReceipt, which reads it")
	}

	for _, index := range compiled.arithmeticOrder {
		column := compiled.columns[index]
		for _, ref := range column.refs {
			refIndex := -1
			for candidate, other := range compiled.columns {
				if other.name == ref {
					refIndex = candidate
				}
			}
			if compiled.columns[refIndex].kind != ColumnArithmetic {
				continue
			}
			if position[ref] >= position[column.name] {
				t.Errorf("%s reads %s but is evaluated first", column.name, ref)
			}
		}
	}
}

// Money combined with a count is still money, so the renderer knows to format
// the average as currency.
func TestCompileSpec_DataTypes(t *testing.T) {
	compiled, err := compileSpec(workedExampleSpec(), testCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}

	want := map[string]DataType{
		"Category":      TypeString,
		"Count":         TypeNumber,
		"Subtotal":      TypeCurrency,
		"Hst":           TypeCurrency,
		"Total":         TypeCurrency,
		"AvgPerReceipt": TypeCurrency,
	}

	for _, column := range compiled.columns {
		if got := column.dataType; got != want[column.name] {
			t.Errorf("column %s dataType = %v, want %v", column.name, got, want[column.name])
		}
	}
}

func TestCompileSpec_ArithmeticOverPlainNumbersIsNotCurrency(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Doubled", Kind: ColumnArithmetic, Expr: "Count * 2"},
		},
	}

	compiled, err := compileSpec(spec, testCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}

	for _, column := range compiled.columns {
		if column.dataType != TypeNumber {
			t.Errorf("column %s dataType = %v, want %v", column.name, column.dataType, TypeNumber)
		}
	}
}

// A label column on an aggregated detail row must know where to read from.
func TestCompileSpec_LabelLevel(t *testing.T) {
	spec := ReportSpec{
		GroupBy: []FieldKey{"paid_by", "tag"},
		Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
		Columns: []Column{
			{Name: "PaidBy", Kind: ColumnLabel, Field: "paid_by"},
			{Name: "Child", Kind: ColumnLabel, Field: "tag"},
			{Name: "Category", Kind: ColumnLabel, Field: "category"},
		},
	}

	compiled, err := compileSpec(spec, testCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}

	want := map[string]int{
		"PaidBy":   0,
		"Child":    1,
		"Category": labelFromDetailBucket,
	}
	for _, column := range compiled.columns {
		if got := column.labelLevel; got != want[column.name] {
			t.Errorf("column %s labelLevel = %d, want %d", column.name, got, want[column.name])
		}
	}
}

func TestCompileSpec_AppliesDefaults(t *testing.T) {
	compiled, err := compileSpec(ReportSpec{
		Columns: []Column{{Name: "Count", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}},
	}, testCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}

	if compiled.spec.NoneLabel != defaultNoneLabel {
		t.Errorf("NoneLabel = %q, want %q", compiled.spec.NoneLabel, defaultNoneLabel)
	}
	if compiled.spec.Config.DivisionScale != defaultDivisionScale {
		t.Errorf("DivisionScale = %d, want %d", compiled.spec.Config.DivisionScale, defaultDivisionScale)
	}
}

func TestCompileSpec_AggregateSourceOverridesStructuralForm(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{{
			Name:   "Subtotal",
			Kind:   ColumnAggregate,
			Agg:    Aggregate{Func: AggCount},
			AggSrc: "SUM(amount)",
		}},
	}

	compiled, err := compileSpec(spec, testCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}

	got := compiled.columns[0].agg
	if got.Func != AggSum || got.Field != "amount" {
		t.Errorf("agg = %+v, want SUM(amount)", got)
	}
}

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"Subtotal", true},
		{"_private", true},
		{"a1", true},
		{"A_B_1", true},
		{"", false},
		{"1a", false},
		{"a b", false},
		{"a-b", false},
		{"a.b", false},
		{"Avg/Receipt", false},
		{"café", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isIdentifier(test.name); got != test.want {
				t.Errorf("isIdentifier(%q) = %v, want %v", test.name, got, test.want)
			}
		})
	}
}
