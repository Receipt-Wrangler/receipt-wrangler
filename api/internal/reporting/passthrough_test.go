package reporting

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// A value the caller supplies and the engine merely carries is easy to get
// right and easy to leave untested: it works, so nothing notices when it stops.
//
// Every test here was written because mutation testing showed the suite could
// not tell the difference. Hardcoding the division scale, dropping a column's
// format, and truncating the cycle error's path all passed the entire suite.

// The division scale reaches AVG and reaches arithmetic division. A template
// that asks for two decimal places must not silently receive six.
func TestRun_DivisionScaleIsHonoured(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Cnt", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Mean", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: "amount"}},
			{Name: "Ratio", Kind: ColumnArithmetic, Expr: "Total / Cnt"},
		},
		GrandTotals: true,
	}

	// Three receipts of 1.00, 0.00 and 0.00 average to a third of a dollar,
	// which no scale renders exactly.
	rows := []Row{
		{"amount": {Num(dec("1.00"))}},
		{"amount": {Num(dec("0.00"))}},
		{"amount": {Num(dec("0.00"))}},
	}

	tests := []struct {
		scale int32
		want  string
	}{
		{1, "0.3"},
		{2, "0.33"},
		{6, "0.333333"},
		{10, "0.3333333333"},
		{maxDivisionScale, "0.333333333333333333333333333333"},
	}

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			scaled := spec
			scaled.Config = EngineConfig{DivisionScale: test.scale}

			model := mustRun(t, scaled, rows)

			// Both the aggregate and the arithmetic division use the scale.
			assertRow(t, model.GrandTotals, map[string]string{"Mean": test.want, "Ratio": test.want})
		})
	}
}

// Zero means unset, and takes the default rather than asking for integer
// division. A column wanting whole numbers asks with ROUND.
func TestRun_ZeroDivisionScaleTakesTheDefault(t *testing.T) {
	spec := ReportSpec{
		Columns: []Column{
			{Name: "Cnt", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
			{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
			{Name: "Ratio", Kind: ColumnArithmetic, Expr: "Total / Cnt"},
			{Name: "Whole", Kind: ColumnArithmetic, Expr: "ROUND(Total / Cnt, 0)"},
		},
		Config:      EngineConfig{DivisionScale: 0},
		GrandTotals: true,
	}

	rows := []Row{{"amount": {Num(dec("1.00"))}}, {"amount": {Num(dec("0.00"))}}, {"amount": {Num(dec("0.00"))}}}
	model := mustRun(t, spec, rows)

	assertRow(t, model.GrandTotals, map[string]string{"Ratio": "0.333333", "Whole": "0"})
}

// Everything the caller hands the engine to carry, it carries.
func TestRun_MetaAndFormatsPassThrough(t *testing.T) {
	spec := ReportSpec{
		Title:     "Quarterly",
		NoneLabel: "Uncategorized",
		Columns: []Column{
			{Name: "Category", Kind: ColumnLabel, Field: "category", Format: "text"},
			{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}, Format: "$ #,##0.00"},
			{Name: "Doubled", Kind: ColumnArithmetic, Expr: "Total * 2", Format: "0.000"},
		},
	}

	generatedAt := time.Date(2026, 6, 1, 9, 30, 0, 0, time.UTC)
	meta := MetaInput{
		GeneratedAt: generatedAt,
		Params:      map[string]string{"period": "May", "group": "Household"},
		Currency:    &CurrencyFormat{Symbol: "€", SymbolAtEnd: true, ThousandsSeparator: ".", DecimalSeparator: ","},
	}

	model, err := Run(spec, testCatalog(t), []Row{{"amount": {Num(dec("1.00"))}}}, meta)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if model.Meta.Title != "Quarterly" {
		t.Errorf("Title = %q", model.Meta.Title)
	}
	if !model.Meta.GeneratedAt.Equal(generatedAt) {
		t.Errorf("GeneratedAt = %v", model.Meta.GeneratedAt)
	}
	if model.Meta.Currency == nil || *model.Meta.Currency != (CurrencyFormat{
		Symbol: "€", SymbolAtEnd: true, ThousandsSeparator: ".", DecimalSeparator: ",",
	}) {
		t.Errorf("Currency = %+v", model.Meta.Currency)
	}
	if model.Meta.NoneLabel != "Uncategorized" {
		t.Errorf("NoneLabel = %q", model.Meta.NoneLabel)
	}
	if model.Meta.Params["period"] != "May" || model.Meta.Params["group"] != "Household" {
		t.Errorf("Params = %v", model.Meta.Params)
	}

	wantFormats := map[string]string{"Category": "text", "Total": "$ #,##0.00", "Doubled": "0.000"}
	for _, descriptor := range model.Columns {
		if got := descriptor.Format; got != wantFormats[descriptor.Name] {
			t.Errorf("column %s Format = %q, want %q", descriptor.Name, got, wantFormats[descriptor.Name])
		}
	}
}

// A nil parameter map stays nil rather than becoming an empty one. A renderer
// can then tell "no parameters" from "parameters, none of them set".
func TestRun_NilParamsStayNil(t *testing.T) {
	spec := ReportSpec{Columns: []Column{{Name: "Cnt", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}}}}

	model, err := Run(spec, testCatalog(t), nil, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if model.Meta.Params != nil {
		t.Errorf("Meta.Params = %v, want nil", model.Meta.Params)
	}
}

// A cycle's error must name the columns that form it, not merely announce that
// one exists. Someone has to fix the template.
func TestValidate_CycleMessageNamesTheCycle(t *testing.T) {
	tests := []struct {
		name    string
		columns []Column
		want    string
	}{
		{
			name:    "self reference",
			columns: []Column{{Name: "Loop", Kind: ColumnArithmetic, Expr: "Loop + 1"}},
			want:    "Loop -> Loop",
		},
		{
			name: "a pair",
			columns: []Column{
				{Name: "A", Kind: ColumnArithmetic, Expr: "B + 1"},
				{Name: "B", Kind: ColumnArithmetic, Expr: "A + 1"},
			},
			want: "A -> B -> A",
		},
		{
			name: "a longer chain",
			columns: []Column{
				{Name: "A", Kind: ColumnArithmetic, Expr: "B"},
				{Name: "B", Kind: ColumnArithmetic, Expr: "C"},
				{Name: "C", Kind: ColumnArithmetic, Expr: "A"},
			},
			want: "A -> B -> C -> A",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(ReportSpec{Columns: test.columns}, testCatalog(t))
			if !errors.Is(err, ErrFormulaCycle) {
				t.Fatalf("Validate() error = %v, want %v", err, ErrFormulaCycle)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Errorf("error %q does not name the cycle %q", err, test.want)
			}
		})
	}
}

// Run reads its arguments and writes to none of them.
func TestRun_DoesNotMutateItsArguments(t *testing.T) {
	spec := workedExampleSpec()
	rows := workedExampleRows()

	specBefore := deepCopySpec(spec)
	rowsBefore := deepCopyRows(rows)

	if _, err := Run(spec, testCatalog(t), rows, testMeta()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	// withDefaults fills NoneLabel and DivisionScale on a copy, and compileSpec
	// infers column types on a copy. Neither may reach the caller's spec.
	if !reflect.DeepEqual(spec, specBefore) {
		t.Errorf("Run mutated the spec:\n got %+v\nwant %+v", spec, specBefore)
	}
	if !reflect.DeepEqual(rows, rowsBefore) {
		t.Errorf("Run mutated the rows")
	}
}

// No cell in the returned model may alias a row the caller still holds.
func TestRun_ModelDoesNotAliasInputRows(t *testing.T) {
	tests := []struct {
		name string
		spec ReportSpec
	}{
		{
			name: "records mode label cell",
			spec: ReportSpec{Columns: []Column{{Name: "Cats", Kind: ColumnLabel, Field: "category"}}},
		},
		{
			name: "aggregated detail bucket label cell",
			spec: ReportSpec{
				Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
				Columns: []Column{{Name: "Cats", Kind: ColumnLabel, Field: "category"}},
			},
		},
		{
			name: "label reading an ancestor bucket",
			spec: ReportSpec{
				GroupBy: []FieldKey{"tag"},
				Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
				Columns: []Column{{Name: "Tag", Kind: ColumnLabel, Field: "tag"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rows := []Row{{"category": {Str("Clothing")}, "tag": {Str("Alex")}, "amount": {Num(dec("1.00"))}}}
			model := mustRun(t, test.spec, rows)

			var cells []Cell
			node := model.Root
			for len(node.Children) > 0 {
				node = node.Children[0]
			}
			cells = node.DetailRows[0].Cells

			for index := range cells {
				for value := range cells[index].Values {
					cells[index].Values[value] = Str("tampered")
				}
			}

			for key, values := range rows[0] {
				for _, value := range values {
					if text, isText := value.Text(); isText && text == "tampered" {
						t.Errorf("mutating a cell reached the caller's row at %s", key)
					}
				}
			}
		})
	}
}

// Two aggregate columns must not share an Agg pointer, nor share one with a
// second run's model.
func TestRun_DescriptorsDoNotAlias(t *testing.T) {
	spec := ReportSpec{Columns: []Column{
		{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: "amount"}},
		{Name: "Most", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMax, Field: "amount"}},
	}}

	first := mustRun(t, spec, nil)
	second := mustRun(t, spec, nil)

	if first.Columns[0].Agg == first.Columns[1].Agg {
		t.Fatalf("two columns share an Agg pointer")
	}
	if first.Columns[0].Agg == second.Columns[0].Agg {
		t.Fatalf("two runs share an Agg pointer")
	}

	first.Columns[0].Agg.Func = AggMin
	if first.Columns[1].Agg.Func != AggMax {
		t.Errorf("mutating one column's aggregate changed another")
	}
	if second.Columns[0].Agg.Func != AggSum {
		t.Errorf("mutating one run's aggregate changed another run")
	}
}

// Run is pure, so concurrent calls over a shared spec and catalog must neither
// race nor disagree. Only meaningful under -race.
func TestRun_ConcurrentRunsAreRaceFreeAndAgree(t *testing.T) {
	spec := workedExampleSpec()
	catalog := testCatalog(t)
	rows := workedExampleRows()

	want := serializeModel(mustRun(t, spec, rows))

	const goroutines = 16
	results := make([]string, goroutines)

	var group sync.WaitGroup
	for index := 0; index < goroutines; index++ {
		group.Add(1)
		go func(slot int) {
			defer group.Done()

			model, err := Run(spec, catalog, rows, testMeta())
			if err != nil {
				t.Errorf("Run() error = %v", err)
				return
			}
			results[slot] = serializeModel(model)
		}(index)
	}
	group.Wait()

	for index, got := range results {
		if got != want {
			t.Errorf("goroutine %d produced a different report", index)
		}
	}
}

func deepCopySpec(spec ReportSpec) ReportSpec {
	copied := spec
	copied.GroupBy = append([]FieldKey(nil), spec.GroupBy...)
	copied.Columns = append([]Column(nil), spec.Columns...)
	return copied
}

func deepCopyRows(rows []Row) []Row {
	copied := make([]Row, len(rows))
	for index, row := range rows {
		duplicate := make(Row, len(row))
		for key, values := range row {
			duplicate[key] = append([]Value(nil), values...)
		}
		copied[index] = duplicate
	}
	return copied
}
