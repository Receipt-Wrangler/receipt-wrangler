package reporting

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// This file holds the laws every ReportModel obeys, whatever the spec and
// whatever the rows. Engine tests assert them alongside their own expectations,
// and the property tests assert them over randomly generated reports, so a
// regression in any one of them trips almost everywhere at once.

// assertModelInvariants checks the structural laws of a report model. It is
// deliberately independent of the engine's implementation: it walks the output
// and re-derives what must be true of it from the spec alone.
func assertModelInvariants(t *testing.T, spec ReportSpec, model ReportModel) {
	t.Helper()

	filled := spec.withDefaults()

	assertColumnAlignment(t, model, model.GrandTotals, "grand totals")
	if filled.GrandTotals != (model.GrandTotals != nil) {
		t.Errorf("GrandTotals present = %v, want %v", model.GrandTotals != nil, filled.GrandTotals)
	}
	assertCellCardinality(t, model, model.GrandTotals, "grand totals", false)

	// The root is synthetic. It names no dimension, carries no subtotals of its
	// own, and is never the (None) bucket even though its value is null.
	if model.Root.Dimension != "" {
		t.Errorf("root Dimension = %q, want empty", model.Root.Dimension)
	}
	if model.Root.Subtotals != nil {
		t.Errorf("root carries subtotals; grand totals are the report-wide ones")
	}
	if model.Root.IsNone {
		t.Errorf("root is marked as the (None) bucket")
	}
	if !model.Root.Value.IsNull() {
		t.Errorf("root Value = %v, want null", model.Root.Value)
	}

	assertGroupNode(t, filled, model, model.Root, 0, "root")
}

func assertGroupNode(t *testing.T, spec ReportSpec, model ReportModel, node GroupNode, depth int, path string) {
	t.Helper()

	isLeaf := depth == len(spec.GroupBy)

	// A node holds children or detail rows, never both. Assert nil-ness rather
	// than length, so an empty-but-present slice is still a violation.
	if isLeaf {
		if node.Children != nil {
			t.Errorf("%s: leaf at depth %d has children", path, depth)
		}
		if node.DetailRows == nil {
			t.Errorf("%s: leaf at depth %d has no detail rows", path, depth)
		}
	} else {
		if node.DetailRows != nil {
			t.Errorf("%s: internal node at depth %d has detail rows", path, depth)
		}
		if node.Children == nil {
			t.Errorf("%s: internal node at depth %d has no children slice", path, depth)
		}
	}

	if depth > 0 {
		// A non-root node names the dimension its siblings were bucketed by, and
		// is the (None) bucket exactly when its value is null.
		if want := spec.GroupBy[depth-1]; node.Dimension != want {
			t.Errorf("%s: Dimension = %q, want %q", path, node.Dimension, want)
		}
		if node.IsNone != node.Value.IsNull() {
			t.Errorf("%s: IsNone = %v but Value.IsNull() = %v", path, node.IsNone, node.Value.IsNull())
		}

		wantSubtotals := spec.Subtotals
		if (node.Subtotals != nil) != wantSubtotals {
			t.Errorf("%s: Subtotals present = %v, want %v", path, node.Subtotals != nil, wantSubtotals)
		}
		assertColumnAlignment(t, model, node.Subtotals, path+" subtotals")
		assertCellCardinality(t, model, node.Subtotals, path+" subtotals", false)
	}

	assertSortedByValue(t, nodeValues(node.Children), path+" children")

	childRecords := 0
	for _, child := range node.Children {
		childPath := fmt.Sprintf("%s/%s", path, child.Value)
		assertGroupNode(t, spec, model, child, depth+1, childPath)
		childRecords += child.RecordCount
	}
	if len(node.Children) > 0 && node.RecordCount != childRecords {
		t.Errorf("%s: RecordCount = %d, but children sum to %d", path, node.RecordCount, childRecords)
	}

	if !isLeaf {
		return
	}

	records := spec.Detail.Mode == DetailRecords
	if records && node.RecordCount != len(node.DetailRows) {
		t.Errorf("%s: RecordCount = %d, but there are %d records", path, node.RecordCount, len(node.DetailRows))
	}

	if !records {
		assertSortedByValue(t, detailValues(node.DetailRows), path+" detail rows")
	}

	for index, row := range node.DetailRows {
		rowPath := fmt.Sprintf("%s[%d]", path, index)
		assertColumnAlignment(t, model, row.Cells, rowPath)
		assertCellCardinality(t, model, row.Cells, rowPath, records)

		if records {
			// A record row is one source record; nothing keys it.
			if row.Dimension != "" || row.IsNone || !row.Value.IsNull() {
				t.Errorf("%s: a records-mode detail row is keyed: %+v", rowPath, row)
			}
			continue
		}

		if row.Dimension != spec.Detail.By {
			t.Errorf("%s: Dimension = %q, want %q", rowPath, row.Dimension, spec.Detail.By)
		}
		if row.IsNone != row.Value.IsNull() {
			t.Errorf("%s: IsNone = %v but Value.IsNull() = %v", rowPath, row.IsNone, row.Value.IsNull())
		}
	}
}

// assertColumnAlignment: a rendered row carries one cell per column, in the
// column order, named after it. A renderer reads cells positionally.
func assertColumnAlignment(t *testing.T, model ReportModel, cells []Cell, path string) {
	t.Helper()

	if cells == nil {
		return
	}
	if len(cells) != len(model.Columns) {
		t.Fatalf("%s: %d cells for %d columns", path, len(cells), len(model.Columns))
	}
	for index, cell := range cells {
		if cell.Column != model.Columns[index].Name {
			t.Errorf("%s: cell %d is %q, want %q", path, index, cell.Column, model.Columns[index].Name)
		}
	}
}

// assertCellCardinality: aggregate and arithmetic cells hold exactly one value.
// Only a label cell on a records-mode detail row may hold several, because only
// a source record can carry several categories.
func assertCellCardinality(t *testing.T, model ReportModel, cells []Cell, path string, recordsDetailRow bool) {
	t.Helper()

	for index, cell := range cells {
		column := model.Columns[index]

		switch column.Kind {
		case ColumnAggregate, ColumnArithmetic:
			if len(cell.Values) != 1 {
				t.Errorf("%s: %s cell %q holds %d values, want 1", path, column.Kind, cell.Column, len(cell.Values))
			}
		case ColumnLabel:
			if !recordsDetailRow && len(cell.Values) > 1 {
				t.Errorf("%s: label cell %q holds %d values outside records mode", path, cell.Column, len(cell.Values))
			}
		}
	}
}

// assertSortedByValue: siblings ascend by value, are pairwise distinct, and the
// (None) bucket is last.
//
// Distinctness matters as much as order. Two siblings that compare equal are
// one value split across two buckets, which is what a bucket key finer than
// equality produces. Since the siblings are sorted, equals are adjacent.
func assertSortedByValue(t *testing.T, values []Value, path string) {
	t.Helper()

	for index := 1; index < len(values); index++ {
		switch compareValues(values[index-1], values[index]) {
		case 1:
			t.Errorf("%s: %v sorts after %v", path, values[index-1], values[index])
		case 0:
			t.Errorf("%s: %v appears as two buckets", path, values[index])
		}
	}
	for index, value := range values {
		if value.IsNull() && index != len(values)-1 {
			t.Errorf("%s: the (None) bucket is at %d of %d, want last", path, index, len(values))
		}
	}
}

func nodeValues(nodes []GroupNode) []Value {
	values := make([]Value, 0, len(nodes))
	for _, node := range nodes {
		values = append(values, node.Value)
	}
	return values
}

func detailValues(rows []DetailRow) []Value {
	values := make([]Value, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.Value)
	}
	return values
}

// serializeModel renders a model to a stable string.
//
// Determinism is compared through this rather than through reflect.DeepEqual,
// which inspects a decimal's internal exponent and would call 200 and 200.00
// different. It also makes a failing comparison readable.
func serializeModel(model ReportModel) string {
	var out strings.Builder

	fmt.Fprintf(&out, "title=%q generatedAt=%s currency=%q none=%q\n",
		model.Meta.Title, model.Meta.GeneratedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
		model.Meta.CurrencyFormat, model.Meta.NoneLabel)

	keys := make([]string, 0, len(model.Meta.Params))
	for key := range model.Meta.Params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&out, "param %s=%q\n", key, model.Meta.Params[key])
	}

	for _, column := range model.Columns {
		fmt.Fprintf(&out, "column %s label=%q kind=%s type=%s format=%q expr=%q\n",
			column.Name, column.Label, column.Kind, column.DataType, column.Format, column.ExprSrc)
		if column.Agg != nil {
			fmt.Fprintf(&out, "  agg=%s\n", column.Agg)
		}
	}

	serializeGroup(&out, model.Root, 0)
	fmt.Fprintf(&out, "grand %s\n", serializeCells(model.GrandTotals))

	return out.String()
}

func serializeGroup(out *strings.Builder, node GroupNode, depth int) {
	indent := strings.Repeat("  ", depth)

	fmt.Fprintf(out, "%sgroup dim=%s value=%s none=%v records=%d\n",
		indent, node.Dimension, serializeValue(node.Value), node.IsNone, node.RecordCount)

	for _, row := range node.DetailRows {
		fmt.Fprintf(out, "%s  row dim=%s value=%s none=%v %s\n",
			indent, row.Dimension, serializeValue(row.Value), row.IsNone, serializeCells(row.Cells))
	}
	for _, child := range node.Children {
		serializeGroup(out, child, depth+1)
	}
	if node.Subtotals != nil {
		fmt.Fprintf(out, "%s  subtotal %s\n", indent, serializeCells(node.Subtotals))
	}
}

func serializeCells(cells []Cell) string {
	if cells == nil {
		return "<none>"
	}

	parts := make([]string, 0, len(cells))
	for _, cell := range cells {
		values := make([]string, 0, len(cell.Values))
		for _, value := range cell.Values {
			values = append(values, serializeValue(value))
		}
		parts = append(parts, cell.Column+"=["+strings.Join(values, ",")+"]")
	}

	return strings.Join(parts, " ")
}

// serializeValue renders a value canonically. A number renders through the
// decimal's own normalization, so 200 and 200.00 serialize alike, and a date
// through its absolute instant.
func serializeValue(value Value) string {
	switch value.Type() {
	case ValueNull:
		return "null"
	case ValueNumber:
		number, _ := value.Decimal()
		return "n" + number.String()
	case ValueString:
		text, _ := value.Text()
		return fmt.Sprintf("s%q", text)
	case ValueBool:
		flag, _ := value.Boolean()
		return fmt.Sprintf("b%v", flag)
	case ValueDate:
		return "d" + bucketKey(value)
	}
	return "?"
}
