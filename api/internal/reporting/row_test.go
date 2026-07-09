package reporting

import "testing"

func TestRow_Get(t *testing.T) {
	row := Row{
		"amount": {Num(dec("120.00"))},
		"tag":    {Str("Alex"), Str("Sam")},
		"name":   {},
	}

	tests := []struct {
		name  string
		row   Row
		key   FieldKey
		want  int
		first Value
	}{
		{"present scalar", row, "amount", 1, Num(dec("120.00"))},
		{"present multi value", row, "tag", 2, Str("Alex")},
		{"present but empty", row, "name", 0, Null()},
		{"absent key", row, "status", 0, Null()},
		{"nil row", nil, "amount", 0, Null()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.row.Get(test.key)
			if len(got) != test.want {
				t.Fatalf("Get(%q) returned %d values, want %d", test.key, len(got), test.want)
			}
			if test.want > 0 && !got[0].Equal(test.first) {
				t.Errorf("Get(%q)[0] = %v, want %v", test.key, got[0], test.first)
			}
		})
	}
}

// A measure with no value is null, never a zero amount — the two mean different
// things once they reach an aggregate.
func TestRow_Measure(t *testing.T) {
	row := Row{
		"amount":   {Num(dec("120.00"))},
		"hst":      {},
		"multi":    {Num(dec("1")), Num(dec("2"))},
		"nullable": {Null()},
	}

	tests := []struct {
		name string
		row  Row
		key  FieldKey
		want Value
	}{
		{"present", row, "amount", Num(dec("120.00"))},
		{"empty slice is null", row, "hst", Null()},
		{"absent key is null", row, "missing", Null()},
		{"nil row is null", nil, "amount", Null()},
		{"explicit null stays null", row, "nullable", Null()},
		{"multi valued measure takes the first", row, "multi", Num(dec("1"))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.row.Measure(test.key); !got.Equal(test.want) {
				t.Errorf("Measure(%q) = %v, want %v", test.key, got, test.want)
			}
		})
	}
}

// An absent dimension is never dropped: it becomes exactly one (None) bucket.
func TestRow_DimensionValues(t *testing.T) {
	row := Row{
		"category": {Str("Clothing"), Str("Medical")},
		"tag":      {Str("Alex")},
		"empty":    {},
	}

	tests := []struct {
		name string
		row  Row
		key  FieldKey
		want []Value
	}{
		{"multi value fans out", row, "category", []Value{Str("Clothing"), Str("Medical")}},
		{"single value", row, "tag", []Value{Str("Alex")}},
		{"empty slice is one none bucket", row, "empty", []Value{Null()}},
		{"absent key is one none bucket", row, "missing", []Value{Null()}},
		{"nil row is one none bucket", nil, "category", []Value{Null()}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.row.dimensionValues(test.key)
			if len(got) != len(test.want) {
				t.Fatalf("dimensionValues(%q) returned %d values, want %d", test.key, len(got), len(test.want))
			}
			for i := range got {
				if !got[i].Equal(test.want[i]) {
					t.Errorf("dimensionValues(%q)[%d] = %v, want %v", test.key, i, got[i], test.want[i])
				}
			}
		})
	}
}
