package reporting

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// The laws in this file hold for every spec and every row set, so they are
// checked against randomly generated reports rather than against fixtures. A
// fixture proves the engine right about one report; a law proves it right about
// the shape of every report, and catches the special case an implementation
// quietly got away with.
//
// Everything here is seeded and reproducible. A failure prints the seed and the
// case that broke, and REPORTING_SEED re-runs it.

const propertyCases = 400

// propertySeed is fixed so the suite is deterministic, and overridable so a
// soak run can explore further.
func propertySeed(t *testing.T) int64 {
	t.Helper()

	if raw := os.Getenv("REPORTING_SEED"); raw != "" {
		seed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			t.Fatalf("REPORTING_SEED=%q is not an integer", raw)
		}
		return seed
	}
	return 20260709
}

// propertyCatalog is the field vocabulary random specs are drawn from: two
// currency measures, a plain numeric measure, two multi-value dimensions, and
// one dimension of each remaining type.
func propertyCatalog(t *testing.T) FieldCatalog {
	t.Helper()

	catalog, err := NewFieldCatalog(
		FieldRef{Key: "amount", Label: "Amount", DataType: TypeCurrency},
		FieldRef{Key: "hst", Label: "Hst", DataType: TypeCurrency},
		FieldRef{Key: "quantity", Label: "Quantity", DataType: TypeNumber},
		FieldRef{Key: "category", Label: "Category", DataType: TypeString, Multi: true},
		FieldRef{Key: "tag", Label: "Tag", DataType: TypeString, Multi: true},
		FieldRef{Key: "paid_by", Label: "Paid By", DataType: TypeString},
		FieldRef{Key: "date", Label: "Date", DataType: TypeDate},
		FieldRef{Key: "resolved", Label: "Resolved", DataType: TypeBool},
	)
	if err != nil {
		t.Fatalf("NewFieldCatalog() error = %v", err)
	}
	return catalog
}

var (
	propertyDimensions = []FieldKey{"category", "tag", "paid_by", "date", "resolved"}
	propertyMeasures   = []FieldKey{"amount", "hst", "quantity"}
	propertyMultiValue = map[FieldKey]bool{"category": true, "tag": true}
)

// randomCase is one generated report: a spec, the rows it runs over, and the
// seed that produced them.
type randomCase struct {
	seed int64
	spec ReportSpec
	rows []Row
}

func (c randomCase) String() string {
	return fmt.Sprintf("seed=%d groupBy=%v detail=%s/%s columns=%d rows=%d subtotals=%v grand=%v scale=%d",
		c.seed, c.spec.GroupBy, c.spec.Detail.Mode, c.spec.Detail.By,
		len(c.spec.Columns), len(c.rows), c.spec.Subtotals, c.spec.GrandTotals, c.spec.Config.DivisionScale)
}

// generateCase builds a random but valid report. The column set always carries
// SUM, COUNT, AVG, MIN and MAX over one measure plus the arithmetic that reads
// them, because those are the columns the laws below are stated over.
func generateCase(random *rand.Rand, seed int64) randomCase {
	levels := random.Intn(4) // 0..3 grouping levels
	shuffled := append([]FieldKey(nil), propertyDimensions...)
	random.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	groupBy := shuffled[:levels]

	detail := DetailSpec{Mode: DetailRecords}
	remaining := shuffled[levels:]
	if len(remaining) > 0 && random.Intn(2) == 0 {
		detail = DetailSpec{Mode: DetailAggregate, By: remaining[random.Intn(len(remaining))]}
	}

	measure := propertyMeasures[random.Intn(len(propertyMeasures))]

	columns := []Column{
		{Name: "Cnt", Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
		{Name: "Total", Kind: ColumnAggregate, Agg: Aggregate{Func: AggSum, Field: measure}},
		{Name: "Mean", Kind: ColumnAggregate, Agg: Aggregate{Func: AggAvg, Field: measure}},
		{Name: "Least", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMin, Field: measure}},
		{Name: "Most", Kind: ColumnAggregate, Agg: Aggregate{Func: AggMax, Field: measure}},
		// Additive, and non-linear: the second is the one that must be recomputed.
		{Name: "Doubled", Kind: ColumnArithmetic, Expr: "Total + Total"},
		{Name: "Ratio", Kind: ColumnArithmetic, Expr: "Total / Cnt"},
		{Name: "Chained", Kind: ColumnArithmetic, Expr: "ROUND(Ratio * 2, 3)"},
	}
	if detail.Mode == DetailAggregate {
		columns = append(columns, Column{Name: "Bucket", Kind: ColumnLabel, Field: detail.By})
	}

	spec := ReportSpec{
		Title:       "property",
		GroupBy:     groupBy,
		Detail:      detail,
		Columns:     columns,
		Subtotals:   random.Intn(2) == 0,
		GrandTotals: random.Intn(2) == 0,
		Config:      EngineConfig{DivisionScale: int32(1 + random.Intn(8))},
	}

	rows := make([]Row, random.Intn(12))
	for index := range rows {
		rows[index] = generateRow(random, measure)
	}

	return randomCase{seed: seed, spec: spec, rows: rows}
}

func generateRow(random *rand.Rand, measure FieldKey) Row {
	row := Row{}

	for _, dimension := range propertyDimensions {
		switch random.Intn(6) {
		case 0:
			// Absent: the row falls into the (None) bucket.
			continue
		case 1:
			row[dimension] = []Value{Null()}
		default:
			row[dimension] = randomDimensionValues(random, dimension)
		}
	}

	// The measure is null on some rows, which SUM skips and AVG excludes from
	// its divisor, while COUNT still counts the row.
	if random.Intn(5) > 0 {
		row[measure] = []Value{Num(randomDecimal(random))}
	}

	return row
}

func randomDimensionValues(random *rand.Rand, dimension FieldKey) []Value {
	count := 1
	if propertyMultiValue[dimension] {
		count = 1 + random.Intn(3)
	}

	values := make([]Value, 0, count)
	for i := 0; i < count; i++ {
		values = append(values, randomDimensionValue(random, dimension))
	}
	return values
}

// propertyDates deliberately includes the pathological instants, not merely a
// spread of ordinary ones. A generator that only produces plausible data only
// tests the plausible paths.
//
// Two pairs earn their place. The first two entries are distinct moments that
// share a UnixNano, so a report grouping on them merges two buckets into one
// unless the key is faithful. The last two are one instant expressed in two
// zones, straddling midnight so they format as different calendar days: they
// must merge into a single bucket whose reported value is the same whichever
// arrived first. Note that the two 12:00 entries are *not* such a pair — one is
// 12:00Z and the other 11:00Z — which is why the colliding pair had to be added
// explicitly rather than assumed.
var propertyDates = []time.Time{
	time.Time{}, // year 1
	time.Time{}.Add(math.MaxInt64).Add(math.MaxInt64).Add(2), // year 585, same UnixNano
	time.Date(1700, 3, 4, 0, 0, 0, 0, time.UTC),
	time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC),
	time.Date(2026, 5, 1, 12, 0, 0, 0, time.FixedZone("plus-one", 3600)),
	time.Date(3000, 1, 1, 0, 0, 0, 1, time.UTC),
	time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),                               // one instant,
	time.Date(2026, 4, 30, 19, 0, 0, 0, time.FixedZone("minus-five", -18000)), // two zones, two days
}

func randomDimensionValue(random *rand.Rand, dimension FieldKey) Value {
	switch dimension {
	case "date":
		return DateVal(propertyDates[random.Intn(len(propertyDates))])
	case "resolved":
		return Bool(random.Intn(2) == 0)
	default:
		// A small alphabet, so buckets collide often. The empty string is a
		// bucket of its own, distinct from (None).
		names := []string{"", "a", "b", "c"}
		return Str(names[random.Intn(len(names))])
	}
}

func randomDecimal(random *rand.Rand) decimal.Decimal {
	// Two decimal places, positive and negative, including zero.
	return decimal.New(int64(random.Intn(20001)-10000), -2)
}

// TestProperty_Laws generates reports and checks every law against each.
func TestProperty_Laws(t *testing.T) {
	seed := propertySeed(t)
	catalog := propertyCatalog(t)
	random := rand.New(rand.NewSource(seed))

	skipped := 0
	for index := 0; index < propertyCases; index++ {
		testCase := generateCase(random, seed)

		if err := Validate(testCase.spec, catalog); err != nil {
			skipped++
			continue
		}

		model, err := Run(testCase.spec, catalog, testCase.rows, MetaInput{})
		if err != nil {
			t.Fatalf("case %d: Run() error = %v\n%s", index, err, testCase)
		}

		t.Run(fmt.Sprintf("case%03d", index), func(t *testing.T) {
			defer func() {
				if t.Failed() {
					t.Logf("reproduce with REPORTING_SEED=%d; case: %s", seed, testCase)
				}
			}()

			assertModelInvariants(t, testCase.spec, model)
			assertBucketCardinality(t, testCase, model)
			assertConservation(t, testCase, model)
			assertRollup(t, testCase, model, model.Root, 0)
			assertArithmeticRecomputed(t, testCase, model)
			assertDeterminism(t, testCase, catalog, model)
		})
	}

	if skipped == propertyCases {
		t.Fatalf("every generated case was rejected by Validate; the generator is broken")
	}
	t.Logf("%d cases, %d rejected by Validate", propertyCases, skipped)
}

// distinctValues returns the values of a field that are pairwise unequal.
//
// It compares with compareValues and deliberately not with bucketKey, and it
// reads the row directly rather than through dimensionValues. Both are the
// point: a law derived from the code it judges cannot catch that code being
// wrong. An earlier draft of this file used the engine's own helper here, and
// consequently agreed with the engine that a receipt tagged "Alex" twice
// belonged to the Alex bucket twice.
//
// A field with no values is one bucket: (None).
func distinctValues(values []Value) []Value {
	if len(values) == 0 {
		return []Value{Null()}
	}

	distinct := make([]Value, 0, len(values))
	for _, value := range values {
		duplicate := false
		for _, kept := range distinct {
			if compareValues(value, kept) == 0 {
				duplicate = true
				break
			}
		}
		if !duplicate {
			distinct = append(distinct, value)
		}
	}
	return distinct
}

// attributions is the number of buckets a row is counted into: the product of
// the distinct values it holds at each grouping level, times the distinct
// values of the detail dimension. An absent level contributes exactly one, the
// (None) bucket.
//
// Computed from the rows alone, so it is an independent prediction rather than
// a restatement of what the engine did.
func attributions(spec ReportSpec, row Row) int {
	count := 1
	for _, level := range spec.GroupBy {
		count *= len(distinctValues(row.Get(level)))
	}
	if spec.Detail.Mode == DetailAggregate {
		count *= len(distinctValues(row.Get(spec.Detail.By)))
	}
	return count
}

// A11: the buckets at a level are exactly the distinct values the rows hold
// there — no fewer, which would mean two values merged, and no more, which
// would mean one value split.
//
// This is the law the date bucket key broke: two instants that compareValues
// called different shared a UnixNano and collapsed into one bucket.
func assertBucketCardinality(t *testing.T, testCase randomCase, model ReportModel) {
	t.Helper()

	if len(testCase.rows) == 0 {
		return
	}

	if len(testCase.spec.GroupBy) > 0 {
		var all []Value
		for _, row := range testCase.rows {
			all = append(all, distinctValues(row.Get(testCase.spec.GroupBy[0]))...)
		}
		if want := len(distinctValues(all)); len(model.Root.Children) != want {
			t.Errorf("A11: the top level has %d buckets, but the rows hold %d distinct values",
				len(model.Root.Children), want)
		}
		return
	}

	if testCase.spec.Detail.Mode == DetailAggregate {
		var all []Value
		for _, row := range testCase.rows {
			all = append(all, distinctValues(row.Get(testCase.spec.Detail.By))...)
		}
		if want := len(distinctValues(all)); len(model.Root.DetailRows) != want {
			t.Errorf("A11: there are %d detail rows, but the rows hold %d distinct values",
				len(model.Root.DetailRows), want)
		}
	}
}

// A1, A2: the grand total counts every attribution, and sums every attributed
// amount. A row that fans into two buckets contributes twice, on purpose.
func assertConservation(t *testing.T, testCase randomCase, model ReportModel) {
	t.Helper()

	if !testCase.spec.GrandTotals {
		return
	}

	measure := measureOf(testCase.spec, "Total")

	wantCount := 0
	wantSum := decimal.Zero
	for _, row := range testCase.rows {
		times := attributions(testCase.spec, row)
		wantCount += times

		if value, isNumber := row.Measure(measure).Decimal(); isNumber {
			wantSum = wantSum.Add(value.Mul(decimal.NewFromInt(int64(times))))
		}
	}

	if got := cellNumber(t, model.GrandTotals, "Cnt"); !got.Equal(decimal.NewFromInt(int64(wantCount))) {
		t.Errorf("A1: grand COUNT = %s, want %d", got, wantCount)
	}
	if model.Root.RecordCount != wantCount {
		t.Errorf("A1: root RecordCount = %d, want %d", model.Root.RecordCount, wantCount)
	}
	if got := cellNumber(t, model.GrandTotals, "Total"); !got.Equal(wantSum) {
		t.Errorf("A2: grand SUM = %s, want %s", got, wantSum)
	}
}

// A3, A4, A5, A6: how a node relates to the rows beneath it.
func assertRollup(t *testing.T, testCase randomCase, model ReportModel, node GroupNode, depth int) {
	t.Helper()

	for _, child := range node.Children {
		assertRollup(t, testCase, model, child, depth+1)
	}

	cells := node.Subtotals
	if depth == 0 {
		cells = model.GrandTotals
	}
	if cells == nil {
		return
	}

	// The rows immediately beneath this node: its children's own totals, or its
	// detail rows.
	var children [][]Cell
	for _, child := range node.Children {
		if child.Subtotals == nil {
			return // subtotals are switched off; nothing to compare against
		}
		children = append(children, child.Subtotals)
	}
	if len(node.Children) == 0 {
		for _, row := range node.DetailRows {
			children = append(children, row.Cells)
		}
	}

	// A4: SUM and COUNT are the sum of their children. Nothing else is.
	for _, column := range []string{"Cnt", "Total"} {
		want := decimal.Zero
		for _, child := range children {
			want = want.Add(cellNumber(t, child, column))
		}
		if got := cellNumber(t, cells, column); !got.Equal(want) {
			t.Errorf("A4: %s at depth %d = %s, want %s (the sum of %d children)", column, depth, got, want, len(children))
		}
	}

	// A5: MIN and MAX are the extremes of their children, and null only when
	// every child is null.
	assertExtremum(t, cells, children, "Least", depth, func(candidate, best decimal.Decimal) bool {
		return candidate.LessThan(best)
	})
	assertExtremum(t, cells, children, "Most", depth, func(candidate, best decimal.Decimal) bool {
		return candidate.GreaterThan(best)
	})

	// A6: AVG is the average over every descendant, never the average of the
	// children's averages. With SUM and COUNT on the same row, that is exactly
	// AVG * COUNT == SUM.
	mean := cellValue(t, cells, "Mean")
	total := cellNumber(t, cells, "Total")
	count := cellNumber(t, cells, "Cnt")

	if mean.IsNull() {
		// AVG is null only when no row beneath carried a value, in which case
		// SUM is zero.
		if !total.IsZero() {
			t.Errorf("A6: AVG at depth %d is null but SUM is %s", depth, total)
		}
		return
	}

	average, _ := mean.Decimal()
	scale := testCase.spec.Config.DivisionScale

	// AVG is rounded to the division scale, so every comparison below allows one
	// unit at that scale. A single value of 1.07 averaged at scale 1 is 1.1,
	// which is genuinely outside [MIN, MAX].
	epsilon := decimal.New(1, -scale)

	// AVG divides by the count of non-null values, which is at most the record
	// count. The exact divisor cannot be recovered from the model, so assert the
	// relation that always holds -- AVG lies between MIN and MAX -- and, when
	// every row carried a value so the divisor is the record count, that
	// AVG * COUNT reconstructs SUM.
	assertBetweenExtremes(t, cells, average, epsilon, depth)

	if everyRowHasAMeasure(testCase) && !count.IsZero() {
		reconstructed := average.Mul(count)
		tolerance := epsilon.Mul(count)
		if reconstructed.Sub(total).Abs().GreaterThan(tolerance) {
			t.Errorf("A6: AVG*COUNT at depth %d = %s, want %s (within %s)", depth, reconstructed, total, tolerance)
		}
	}
}

func assertExtremum(t *testing.T, cells []Cell, children [][]Cell, column string, depth int, better func(a, b decimal.Decimal) bool) {
	t.Helper()

	var want decimal.Decimal
	seen := false
	for _, child := range children {
		value := cellValue(t, child, column)
		if value.IsNull() {
			continue
		}
		number, _ := value.Decimal()
		if !seen || better(number, want) {
			want = number
			seen = true
		}
	}

	got := cellValue(t, cells, column)
	if !seen {
		if !got.IsNull() {
			t.Errorf("A5: %s at depth %d = %v, want null (every child is null)", column, depth, got)
		}
		return
	}
	number, isNumber := got.Decimal()
	if !isNumber || !number.Equal(want) {
		t.Errorf("A5: %s at depth %d = %v, want %s", column, depth, got, want)
	}
}

func assertBetweenExtremes(t *testing.T, cells []Cell, average, epsilon decimal.Decimal, depth int) {
	t.Helper()

	least := cellValue(t, cells, "Least")
	most := cellValue(t, cells, "Most")
	if least.IsNull() || most.IsNull() {
		return
	}

	low, _ := least.Decimal()
	high, _ := most.Decimal()
	if average.LessThan(low.Sub(epsilon)) || average.GreaterThan(high.Add(epsilon)) {
		t.Errorf("A6: AVG %s at depth %d is outside [%s, %s] by more than %s", average, depth, low, high, epsilon)
	}
}

// everyRowHasAMeasure reports whether every row carried a value for the
// measure. When they all did, AVG's divisor is the record count, and AVG times
// COUNT reconstructs SUM. When some are null, the divisor is smaller and the
// model does not expose it.
func everyRowHasAMeasure(testCase randomCase) bool {
	measure := measureOf(testCase.spec, "Total")
	for _, row := range testCase.rows {
		if row.Measure(measure).IsNull() {
			return false
		}
	}
	return true
}

// A7, A8: an arithmetic column is recomputed from the other cells on its own
// row, at every level. Re-evaluating its expression against that row's cells
// must reproduce the cell exactly.
func assertArithmeticRecomputed(t *testing.T, testCase randomCase, model ReportModel) {
	t.Helper()

	compiled, err := compileSpec(testCase.spec, propertyCatalog(t))
	if err != nil {
		t.Fatalf("compileSpec() error = %v", err)
	}
	scale := compiled.spec.Config.DivisionScale

	check := func(cells []Cell, where string) {
		if cells == nil {
			return
		}

		values := make(map[string]Value, len(cells))
		for _, cell := range cells {
			values[cell.Column] = cell.Value()
		}

		for _, column := range compiled.columns {
			if column.kind != ColumnArithmetic {
				continue
			}
			want := evalArithmetic(column.expr, values, scale)
			got := values[column.name]

			if want.IsNull() != got.IsNull() {
				t.Errorf("A7: %s at %s = %v, recomputes to %v", column.name, where, got, want)
				continue
			}
			if want.IsNull() {
				continue
			}
			wantNumber, _ := want.Decimal()
			gotNumber, _ := got.Decimal()
			if !gotNumber.Equal(wantNumber) {
				t.Errorf("A7: %s at %s = %s, recomputes to %s", column.name, where, gotNumber, wantNumber)
			}
		}
	}

	var walk func(node GroupNode, path string)
	walk = func(node GroupNode, path string) {
		check(node.Subtotals, path+" subtotal")
		for index, row := range node.DetailRows {
			check(row.Cells, fmt.Sprintf("%s row %d", path, index))
		}
		for index, child := range node.Children {
			walk(child, fmt.Sprintf("%s/%d", path, index))
		}
	}
	walk(model.Root, "root")
	check(model.GrandTotals, "grand totals")
}

// A9: identical inputs give identical output, and an aggregated report does not
// depend on the order its rows arrived in.
func assertDeterminism(t *testing.T, testCase randomCase, catalog FieldCatalog, model ReportModel) {
	t.Helper()

	again, err := Run(testCase.spec, catalog, testCase.rows, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error on the second run = %v", err)
	}
	if serializeModel(again) != serializeModel(model) {
		t.Errorf("A9: two runs over the same input disagree")
	}

	if testCase.spec.Detail.Mode == DetailRecords {
		// Records keep their input order, so a permutation legitimately changes
		// the report. The totals must not move.
		return
	}

	reversed := make([]Row, len(testCase.rows))
	for index, row := range testCase.rows {
		reversed[len(testCase.rows)-1-index] = row
	}

	permuted, err := Run(testCase.spec, catalog, reversed, MetaInput{})
	if err != nil {
		t.Fatalf("Run() error on the permuted input = %v", err)
	}
	if serializeModel(permuted) != serializeModel(model) {
		t.Errorf("A9: reversing the rows changed an aggregated report")
	}
}

func measureOf(spec ReportSpec, column string) FieldKey {
	for _, candidate := range spec.Columns {
		if candidate.Name == column {
			return candidate.Agg.Field
		}
	}
	return ""
}

func cellValue(t *testing.T, cells []Cell, column string) Value {
	t.Helper()

	for _, cell := range cells {
		if cell.Column == column {
			return cell.Value()
		}
	}
	t.Fatalf("no cell for column %q", column)
	return Null()
}

func cellNumber(t *testing.T, cells []Cell, column string) decimal.Decimal {
	t.Helper()

	number, isNumber := cellValue(t, cells, column).Decimal()
	if !isNumber {
		t.Fatalf("column %q is not a number", column)
	}
	return number
}
