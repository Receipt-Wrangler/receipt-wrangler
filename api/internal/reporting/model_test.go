package reporting

import "testing"

func TestCell_Value(t *testing.T) {
	tests := []struct {
		name string
		cell Cell
		want Value
	}{
		{"single value", Cell{Column: "Subtotal", Values: []Value{Num(dec("200"))}}, Num(dec("200"))},
		{"explicit null", Cell{Column: "Avg", Values: []Value{Null()}}, Null()},
		{"no values is null", Cell{Column: "Category"}, Null()},
		{"empty slice is null", Cell{Column: "Category", Values: []Value{}}, Null()},
		{"multi value takes the first", Cell{Column: "Tags", Values: []Value{Str("Alex"), Str("Sam")}}, Str("Alex")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.cell.Value(); !got.Equal(test.want) {
				t.Errorf("Cell.Value() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestColumn_Heading(t *testing.T) {
	tests := []struct {
		name   string
		column Column
		want   string
	}{
		{"label wins", Column{Name: "AvgPerReceipt", Label: "Avg/Receipt"}, "Avg/Receipt"},
		{"name is the fallback", Column{Name: "Subtotal"}, "Subtotal"},
		{"empty label is not a heading", Column{Name: "Hst", Label: ""}, "Hst"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.column.heading(); got != test.want {
				t.Errorf("heading() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestReportSpec_WithDefaults(t *testing.T) {
	t.Run("fills a zero spec", func(t *testing.T) {
		spec := ReportSpec{}.withDefaults()

		if spec.NoneLabel != defaultNoneLabel {
			t.Errorf("NoneLabel = %q, want %q", spec.NoneLabel, defaultNoneLabel)
		}
		if spec.Config.DivisionScale != defaultDivisionScale {
			t.Errorf("DivisionScale = %d, want %d", spec.Config.DivisionScale, defaultDivisionScale)
		}
	})

	t.Run("keeps what the caller set", func(t *testing.T) {
		spec := ReportSpec{
			NoneLabel: "Uncategorized",
			Config:    EngineConfig{DivisionScale: 2},
		}.withDefaults()

		if spec.NoneLabel != "Uncategorized" {
			t.Errorf("NoneLabel = %q, want %q", spec.NoneLabel, "Uncategorized")
		}
		if spec.Config.DivisionScale != 2 {
			t.Errorf("DivisionScale = %d, want 2", spec.Config.DivisionScale)
		}
	})

	t.Run("does not mutate the receiver", func(t *testing.T) {
		original := ReportSpec{}
		_ = original.withDefaults()

		if original.NoneLabel != "" || original.Config.DivisionScale != 0 {
			t.Errorf("withDefaults mutated its receiver: %+v", original)
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	if got := DefaultConfig().DivisionScale; got != defaultDivisionScale {
		t.Errorf("DefaultConfig().DivisionScale = %d, want %d", got, defaultDivisionScale)
	}
}

func TestStringers(t *testing.T) {
	kinds := map[ColumnKind]string{
		ColumnLabel:      "label",
		ColumnAggregate:  "aggregate",
		ColumnArithmetic: "arithmetic",
	}
	for kind, want := range kinds {
		if got := kind.String(); got != want {
			t.Errorf("ColumnKind(%d).String() = %q, want %q", kind, got, want)
		}
	}

	modes := map[DetailMode]string{
		DetailRecords:   "records",
		DetailAggregate: "aggregate",
	}
	for mode, want := range modes {
		if got := mode.String(); got != want {
			t.Errorf("DetailMode(%d).String() = %q, want %q", mode, got, want)
		}
	}

	functions := map[AggFunc]string{
		AggSum:   "SUM",
		AggCount: "COUNT",
		AggAvg:   "AVG",
		AggMin:   "MIN",
		AggMax:   "MAX",
	}
	for function, want := range functions {
		if got := function.String(); got != want {
			t.Errorf("AggFunc(%d).String() = %q, want %q", function, got, want)
		}
	}
}
