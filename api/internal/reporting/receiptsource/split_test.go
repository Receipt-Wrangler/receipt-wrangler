package receiptsource

import (
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"

	"github.com/shopspring/decimal"
)

// The amount and every currency value in these tests are compared with
// StringFixed or Equal, never reflect.DeepEqual: 50 and 50.00 are the same money
// with different internal exponents.

func decPointer(literal string) *decimal.Decimal {
	value := dec(literal)
	return &value
}

func category(id uint, name string) models.Category {
	return models.Category{BaseModel: models.BaseModel{ID: id}, Name: name}
}

func tag(id uint, name string) models.Tag {
	return models.Tag{BaseModel: models.BaseModel{ID: id}, Name: name}
}

func amounts(receipts []models.Receipt) []string {
	result := make([]string, len(receipts))
	for index, receipt := range receipts {
		result[index] = receipt.Amount.StringFixed(2)
	}
	return result
}

func categoryNames(receipt models.Receipt) []string {
	names := make([]string, len(receipt.Categories))
	for index, category := range receipt.Categories {
		names[index] = category.Name
	}
	return names
}

func tagNames(receipt models.Receipt) []string {
	names := make([]string, len(receipt.Tags))
	for index, tag := range receipt.Tags {
		names[index] = tag.Name
	}
	return names
}

func currencyValue(t *testing.T, receipt models.Receipt, fieldID uint) *decimal.Decimal {
	t.Helper()
	for _, value := range receipt.CustomFields {
		if value.CustomFieldId == fieldID {
			return value.CurrencyValue
		}
	}
	t.Fatalf("receipt has no value for custom field %d", fieldID)
	return nil
}

func sumAmounts(receipts []models.Receipt) decimal.Decimal {
	total := decimal.Zero
	for _, receipt := range receipts {
		total = total.Add(receipt.Amount)
	}
	return total
}

func assertStrings(t *testing.T, label string, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", label, got, want)
	}
}

func gridStrings(grid [][]decimal.Decimal, places int32) [][]string {
	result := make([][]string, len(grid))
	for row := range grid {
		result[row] = make([]string, len(grid[row]))
		for col := range grid[row] {
			result[row][col] = grid[row][col].StringFixed(places)
		}
	}
	return result
}

// ---- allocateGrid -----------------------------------------------------------

func TestAllocateGrid_SingleRow(t *testing.T) {
	cases := []struct {
		name   string
		value  string
		parts  int
		places int32
		want   []string
	}{
		{"even split", "100", 2, 2, []string{"50.00", "50.00"}},
		{"leftover cent goes first", "100", 3, 2, []string{"33.34", "33.33", "33.33"}},
		{"two leftover cents", "100.01", 3, 2, []string{"33.34", "33.34", "33.33"}},
		{"less than a cent each", "0.01", 3, 2, []string{"0.01", "0.00", "0.00"}},
		{"negative mirrors positive", "-100", 3, 2, []string{"-33.34", "-33.33", "-33.33"}},
		{"zero", "0", 3, 2, []string{"0.00", "0.00", "0.00"}},
		{"one part is the identity", "12.34", 1, 2, []string{"12.34"}},
		{"keeps a finer scale", "10.555", 2, 3, []string{"5.278", "5.277"}},
		{"whole number splits into cents", "7", 2, 2, []string{"3.50", "3.50"}},
		{"odd whole number", "1", 3, 2, []string{"0.34", "0.33", "0.33"}},
		{"beyond int64", "92233720368547758080.03", 2, 2, []string{"46116860184273879040.02", "46116860184273879040.01"}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			grid := allocateGrid(dec(test.value), 1, test.parts)
			if len(grid) != 1 {
				t.Fatalf("rows = %d, want 1", len(grid))
			}
			assertStrings(t, "shares", gridStrings(grid, test.places)[0], test.want)

			total := decimal.Zero
			for _, share := range grid[0] {
				total = total.Add(share)
			}
			if !total.Equal(dec(test.value)) {
				t.Errorf("shares sum to %s, want %s", total, test.value)
			}
		})
	}
}

// A share is never coarser than cents, and never finer than the value it came
// from needs.
func TestAllocateGrid_ShareScale(t *testing.T) {
	if got := allocateGrid(dec("100"), 1, 2)[0][0].Exponent(); got != -2 {
		t.Errorf("exponent of a whole-number share = %d, want -2", got)
	}
	if got := allocateGrid(dec("1.2345"), 1, 2)[0][0].Exponent(); got != -4 {
		t.Errorf("exponent of a four-place share = %d, want -4", got)
	}
}

// The regression the one-pass allocation exists for: rounding by category and
// then by tag hands one tag 0.03 and the other 0.01.
func TestAllocateGrid_MarginalsMatchSingleDimensionSplits(t *testing.T) {
	grid := allocateGrid(dec("0.04"), 3, 2)

	rowTotals := make([]string, 3)
	for row := range grid {
		rowTotals[row] = grid[row][0].Add(grid[row][1]).StringFixed(2)
	}
	colTotals := []string{
		grid[0][0].Add(grid[1][0]).Add(grid[2][0]).StringFixed(2),
		grid[0][1].Add(grid[1][1]).Add(grid[2][1]).StringFixed(2),
	}

	assertStrings(t, "row totals", rowTotals, []string{"0.02", "0.01", "0.01"})
	assertStrings(t, "column totals", colTotals, []string{"0.02", "0.02"})
}

// expectedSplitUnits is an equal split of units into parts computed the obvious
// way, independently of allocateGrid: the floor everywhere, plus one for the
// first remainder parts.
func expectedSplitUnits(units int64, parts int) []int64 {
	result := make([]int64, parts)
	for index := range result {
		result[index] = units / int64(parts)
		if int64(index) < units%int64(parts) {
			result[index]++
		}
	}
	return result
}

func TestAllocateGrid_Properties(t *testing.T) {
	seed := int64(20261005)
	if override := os.Getenv("REPORTING_SEED"); override != "" {
		parsed, err := strconv.ParseInt(override, 10, 64)
		if err != nil {
			t.Fatalf("REPORTING_SEED: %v", err)
		}
		seed = parsed
	}
	random := rand.New(rand.NewSource(seed))

	for iteration := 0; iteration < 2000; iteration++ {
		units := random.Int63n(1_000_000)
		negative := random.Intn(4) == 0
		rows, cols := 1+random.Intn(6), 1+random.Intn(6)

		signed := units
		if negative {
			signed = -units
		}
		value := decimal.New(signed, -2)
		grid := allocateGrid(value, rows, cols)

		toUnits := func(share decimal.Decimal) int64 {
			return share.Shift(2).Abs().IntPart()
		}

		quotient := units / int64(rows*cols)
		grand := int64(0)
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				cell := grid[row][col]
				if negative && cell.IsPositive() || !negative && cell.IsNegative() {
					t.Fatalf("seed %d: %s over %dx%d: cell %s has the wrong sign", seed, value, rows, cols, cell)
				}
				if got := toUnits(cell); got != quotient && got != quotient+1 {
					t.Fatalf("seed %d: %s over %dx%d: cell %s is neither q nor q+1 (q=%d)", seed, value, rows, cols, cell, quotient)
				}
				grand += toUnits(cell)
			}
		}
		if grand != units {
			t.Fatalf("seed %d: %s over %dx%d: grid sums to %d units, want %d", seed, value, rows, cols, grand, units)
		}

		wantRows := expectedSplitUnits(units, rows)
		for row := 0; row < rows; row++ {
			total := int64(0)
			for col := 0; col < cols; col++ {
				total += toUnits(grid[row][col])
			}
			if total != wantRows[row] {
				t.Fatalf("seed %d: %s over %dx%d: row %d sums to %d, want %d", seed, value, rows, cols, row, total, wantRows[row])
			}
		}

		wantCols := expectedSplitUnits(units, cols)
		for col := 0; col < cols; col++ {
			total := int64(0)
			for row := 0; row < rows; row++ {
				total += toUnits(grid[row][col])
			}
			if total != wantCols[col] {
				t.Fatalf("seed %d: %s over %dx%d: column %d sums to %d, want %d", seed, value, rows, cols, col, total, wantCols[col])
			}
		}
	}
}

// The integer path must not be bounded by int64: the units are a big.Int.
func TestAllocateGrid_HugeValueIsExact(t *testing.T) {
	units, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
	value := decimal.NewFromBigInt(units, -2)

	grid := allocateGrid(value, 3, 7)
	total := decimal.Zero
	for row := range grid {
		for col := range grid[row] {
			total = total.Add(grid[row][col])
		}
	}
	if !total.Equal(value) {
		t.Errorf("grid sums to %s, want %s", total, value)
	}
}

// ---- Split ------------------------------------------------------------------

func TestSplit_InactiveOptionsAreTheIdentity(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("100"),
		Categories: []models.Category{category(1, "Food"), category(2, "Fuel")},
		Tags:       []models.Tag{tag(1, "Alex"), tag(2, "Sam")},
	}}

	split := source.Split(receipts, SplitOptions{})
	if len(split.Receipts) != 1 {
		t.Fatalf("copies = %d, want 1", len(split.Receipts))
	}
	assertStrings(t, "amount", amounts(split.Receipts), []string{"100.00"})
	assertStrings(t, "categories", categoryNames(split.Receipts[0]), []string{"Food", "Fuel"})
	assertStrings(t, "tags", tagNames(split.Receipts[0]), []string{"Alex", "Sam"})
}

func TestSplit_NothingToSplit(t *testing.T) {
	source := mustNew(t)
	options := SplitOptions{Categories: true, Tags: true}

	tests := []struct {
		name    string
		receipt models.Receipt
	}{
		{"no categories or tags", models.Receipt{Amount: dec("100")}},
		{"one category", models.Receipt{Amount: dec("100"), Categories: []models.Category{category(1, "Food")}}},
		{"one tag", models.Receipt{Amount: dec("100"), Tags: []models.Tag{tag(1, "Alex")}}},
		{"one of each", models.Receipt{
			Amount:     dec("100"),
			Categories: []models.Category{category(1, "Food")},
			Tags:       []models.Tag{tag(1, "Alex")},
		}},
		{"the same category twice", models.Receipt{
			Amount:     dec("100"),
			Categories: []models.Category{category(1, "Food"), category(1, "Food")},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			split := source.Split([]models.Receipt{test.receipt}, options)
			if len(split.Receipts) != 1 {
				t.Fatalf("copies = %d, want 1", len(split.Receipts))
			}
			assertStrings(t, "amount", amounts(split.Receipts), []string{"100.00"})
		})
	}
}

func TestSplit_Categories(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount: dec("100"),
		// Deliberately out of id order: the leftover cent follows the id.
		Categories: []models.Category{category(3, "Clothing"), category(1, "Food"), category(2, "Fuel")},
		Tags:       []models.Tag{tag(1, "Alex"), tag(2, "Sam")},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})

	assertStrings(t, "amounts", amounts(split.Receipts), []string{"33.34", "33.33", "33.33"})
	assertStrings(t, "copy 0 categories", categoryNames(split.Receipts[0]), []string{"Food"})
	assertStrings(t, "copy 1 categories", categoryNames(split.Receipts[1]), []string{"Fuel"})
	assertStrings(t, "copy 2 categories", categoryNames(split.Receipts[2]), []string{"Clothing"})

	// Tags are not split, so every copy keeps them all.
	for index, copied := range split.Receipts {
		assertStrings(t, "copy "+strconv.Itoa(index)+" tags", tagNames(copied), []string{"Alex", "Sam"})
	}
}

func TestSplit_Tags(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("10"),
		Categories: []models.Category{category(1, "Food"), category(2, "Fuel")},
		Tags:       []models.Tag{tag(1, "Alex"), tag(2, "Sam"), tag(3, "Jo"), tag(4, "Kim")},
	}}

	split := source.Split(receipts, SplitOptions{Tags: true})

	assertStrings(t, "amounts", amounts(split.Receipts), []string{"2.50", "2.50", "2.50", "2.50"})
	for index, copied := range split.Receipts {
		if len(copied.Tags) != 1 {
			t.Errorf("copy %d has %d tags, want 1", index, len(copied.Tags))
		}
		assertStrings(t, "copy "+strconv.Itoa(index)+" categories", categoryNames(copied), []string{"Food", "Fuel"})
	}
}

// Both dimensions: one copy per (category, tag) pair, whose category totals
// and tag totals each match a split over that dimension alone.
func TestSplit_BothDimensionsCrossProduct(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("0.04"),
		Categories: []models.Category{category(1, "A"), category(2, "B"), category(3, "C")},
		Tags:       []models.Tag{tag(1, "X"), tag(2, "Y")},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true, Tags: true})
	if len(split.Receipts) != 6 {
		t.Fatalf("copies = %d, want 6", len(split.Receipts))
	}

	byCategory := map[string]decimal.Decimal{}
	byTag := map[string]decimal.Decimal{}
	for _, copied := range split.Receipts {
		if len(copied.Categories) != 1 || len(copied.Tags) != 1 {
			t.Fatalf("copy holds %d categories and %d tags, want one of each", len(copied.Categories), len(copied.Tags))
		}
		byCategory[copied.Categories[0].Name] = byCategory[copied.Categories[0].Name].Add(copied.Amount)
		byTag[copied.Tags[0].Name] = byTag[copied.Tags[0].Name].Add(copied.Amount)
	}

	want := map[string]string{"A": "0.02", "B": "0.01", "C": "0.01", "X": "0.02", "Y": "0.02"}
	for name, total := range byCategory {
		if total.StringFixed(2) != want[name] {
			t.Errorf("category %s = %s, want %s", name, total.StringFixed(2), want[name])
		}
	}
	for name, total := range byTag {
		if total.StringFixed(2) != want[name] {
			t.Errorf("tag %s = %s, want %s", name, total.StringFixed(2), want[name])
		}
	}
	if !sumAmounts(split.Receipts).Equal(dec("0.04")) {
		t.Errorf("copies sum to %s, want 0.04", sumAmounts(split.Receipts))
	}
}

func TestSplit_DuplicateCategoryCountsOnce(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("100"),
		Categories: []models.Category{category(1, "Food"), category(2, "Fuel"), category(1, "Food")},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})
	assertStrings(t, "amounts", amounts(split.Receipts), []string{"50.00", "50.00"})
}

// Unsaved categories (id 0) are told apart, and ordered, by name.
func TestSplit_UnsavedCategoriesAreDistinguishedByName(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("1"),
		Categories: []models.Category{{Name: "Zed"}, {Name: "Abe"}, {Name: "Mid"}, {Name: "Abe"}},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})
	assertStrings(t, "amounts", amounts(split.Receipts), []string{"0.34", "0.33", "0.33"})
	assertStrings(t, "first copy", categoryNames(split.Receipts[0]), []string{"Abe"})
	assertStrings(t, "last copy", categoryNames(split.Receipts[2]), []string{"Zed"})
}

func TestSplit_CurrencyCustomFields(t *testing.T) {
	source, err := New(append(testCustomFields(),
		models.CustomField{BaseModel: models.BaseModel{ID: 6}, Name: "Tip", Type: models.CURRENCY},
		models.CustomField{BaseModel: models.BaseModel{ID: 7}, Name: "Deposit", Type: models.CURRENCY},
	), time.UTC)
	if err != nil {
		t.Fatal(err)
	}

	note := "kept"
	receipts := []models.Receipt{{
		Amount:     dec("100"),
		Categories: []models.Category{category(1, "Food"), category(2, "Fuel")},
		CustomFields: []models.CustomFieldValue{
			{CustomFieldId: hstFieldID, CurrencyValue: decPointer("13")},
			{CustomFieldId: 6, CurrencyValue: decPointer("10.01")},
			{CustomFieldId: 7, CurrencyValue: nil},
			{CustomFieldId: noteFieldID, StringValue: &note},
			// A currency-looking value on a field the catalog does not know.
			{CustomFieldId: 99, CurrencyValue: decPointer("40")},
		},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true, Excluded: map[uint]struct{}{hstFieldID: {}}})
	if len(split.Receipts) != 2 {
		t.Fatalf("copies = %d, want 2", len(split.Receipts))
	}

	tips := []string{}
	for _, copied := range split.Receipts {
		if got := currencyValue(t, copied, hstFieldID); got.StringFixed(2) != "13.00" {
			t.Errorf("excluded HST = %s, want 13.00 on every copy", got.StringFixed(2))
		}
		tips = append(tips, currencyValue(t, copied, 6).StringFixed(2))
		if got := currencyValue(t, copied, 7); got != nil {
			t.Errorf("nil Deposit became %s, want nil", got)
		}
		if got := currencyValue(t, copied, 99); got.StringFixed(2) != "40.00" {
			t.Errorf("unknown field = %s, want it untouched at 40.00", got.StringFixed(2))
		}
		for _, value := range copied.CustomFields {
			if value.CustomFieldId == noteFieldID && (value.StringValue == nil || *value.StringValue != "kept") {
				t.Errorf("text value not carried through")
			}
		}
	}
	assertStrings(t, "tips", tips, []string{"5.01", "5.00"})
}

// A receipt may hold two values for one field; each is divided on its own, and
// the lowest id still wins when the copy becomes a row.
func TestSplit_DuplicateCustomFieldValuesAreEachSplit(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("100"),
		Categories: []models.Category{category(1, "Food"), category(2, "Fuel")},
		CustomFields: []models.CustomFieldValue{
			{BaseModel: models.BaseModel{ID: 9}, CustomFieldId: hstFieldID, CurrencyValue: decPointer("20")},
			{BaseModel: models.BaseModel{ID: 4}, CustomFieldId: hstFieldID, CurrencyValue: decPointer("8")},
		},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})
	rows := source.Rows(split.Receipts)
	for index, row := range rows {
		number, _ := row.Measure(CustomFieldKey(hstFieldID)).Decimal()
		if number.StringFixed(2) != "4.00" {
			t.Errorf("row %d HST = %s, want the lowest-id value's share, 4.00", index, number.StringFixed(2))
		}
	}
}

func TestSplit_DoesNotMutateTheInput(t *testing.T) {
	source := mustNew(t)

	// Both dimensions, and each alone: an unsplit dimension's whole list is what
	// every copy carries, so it is the one most at risk of being shared.
	for _, options := range []SplitOptions{
		{Categories: true, Tags: true},
		{Categories: true},
		{Tags: true},
	} {
		hst := decPointer("13")
		receipts := []models.Receipt{{
			Amount:       dec("100"),
			Categories:   []models.Category{category(2, "Fuel"), category(1, "Food")},
			Tags:         []models.Tag{tag(2, "Sam"), tag(1, "Alex")},
			CustomFields: []models.CustomFieldValue{{CustomFieldId: hstFieldID, CurrencyValue: hst}},
		}}

		split := source.Split(receipts, options)

		// Rewrite everything a later stage could touch on every copy.
		for index := range split.Receipts {
			split.Receipts[index].Categories[0].Name = "(Restricted)"
			split.Receipts[index].Tags[0].Name = "(Restricted)"
			split.Receipts[index].CustomFields[0].CurrencyValue = decPointer("0")
		}

		if receipts[0].Amount.StringFixed(2) != "100.00" {
			t.Errorf("%+v: input amount became %s", options, receipts[0].Amount)
		}
		assertStrings(t, "input categories", categoryNames(receipts[0]), []string{"Fuel", "Food"})
		assertStrings(t, "input tags", tagNames(receipts[0]), []string{"Sam", "Alex"})
		if receipts[0].CustomFields[0].CurrencyValue != hst || hst.StringFixed(2) != "13.00" {
			t.Errorf("%+v: input HST changed to %s", options, receipts[0].CustomFields[0].CurrencyValue)
		}
	}
}

func TestSplit_CopiesShareNothingMutable(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:       dec("100"),
		Categories:   []models.Category{category(1, "Food"), category(2, "Fuel")},
		CustomFields: []models.CustomFieldValue{{CustomFieldId: hstFieldID, CurrencyValue: decPointer("13")}},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})
	first, second := split.Receipts[0], split.Receipts[1]

	if first.CustomFields[0].CurrencyValue == second.CustomFields[0].CurrencyValue {
		t.Fatal("two copies share one currency value pointer")
	}
	first.Categories[0].Name = "changed"
	if second.Categories[0].Name == "changed" {
		t.Error("two copies share one category slice")
	}
	first.CustomFields[0].CustomFieldId = 42
	if second.CustomFields[0].CustomFieldId == 42 {
		t.Error("two copies share one custom field slice")
	}
}

func TestSplit_IsDeterministicRegardlessOfInputOrder(t *testing.T) {
	source := mustNew(t)
	build := func(categories []models.Category) []models.Receipt {
		return []models.Receipt{{Amount: dec("100"), Categories: categories}}
	}

	forward := source.Split(build([]models.Category{category(1, "A"), category(2, "B"), category(3, "C")}), SplitOptions{Categories: true})
	reverse := source.Split(build([]models.Category{category(3, "C"), category(2, "B"), category(1, "A")}), SplitOptions{Categories: true})

	for index := range forward.Receipts {
		if categoryNames(forward.Receipts[index])[0] != categoryNames(reverse.Receipts[index])[0] ||
			!forward.Receipts[index].Amount.Equal(reverse.Receipts[index].Amount) {
			t.Errorf("copy %d differs between input orders", index)
		}
	}
}

func TestSplit_KeepsEveryOtherFieldOnEachCopy(t *testing.T) {
	source := mustNew(t)
	receipt := fullReceipt()
	receipt.Categories = []models.Category{category(1, "Food"), category(2, "Fuel")}

	split := source.Split([]models.Receipt{receipt}, SplitOptions{Categories: true})
	for _, copied := range split.Receipts {
		if copied.ID != receipt.ID || copied.Name != receipt.Name || !copied.Date.Equal(receipt.Date) ||
			copied.Status != receipt.Status || copied.PaidByUser.DisplayName != receipt.PaidByUser.DisplayName {
			t.Errorf("copy lost a field the report still groups by: %+v", copied)
		}
	}
}

// ---- Merge ------------------------------------------------------------------

func TestMerge_WithoutCollisionsIsOneRowPerCopy(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{
		{Amount: dec("100"), Categories: []models.Category{category(1, "A"), category(2, "B")}},
		{Amount: dec("7")},
	}

	merged, perOrigin := source.Merge(source.Split(receipts, SplitOptions{Categories: true}))

	assertStrings(t, "amounts", amounts(merged), []string{"50.00", "50.00", "7.00"})
	if !reflect.DeepEqual(perOrigin, []int{2, 1}) {
		t.Errorf("perOrigin = %v, want [2 1]", perOrigin)
	}
}

// Two hidden categories substituted to one (Restricted) marker are one row
// holding both shares; the excluded value is not summed.
func TestMerge_FoldsCopiesTheSubstitutionMadeIdentical(t *testing.T) {
	source, err := New(append(testCustomFields(),
		models.CustomField{BaseModel: models.BaseModel{ID: 6}, Name: "Tip", Type: models.CURRENCY},
	), time.UTC)
	if err != nil {
		t.Fatal(err)
	}
	receipts := []models.Receipt{{
		Amount:     dec("100"),
		Categories: []models.Category{category(1, "Visible"), category(2, "Hidden1"), category(3, "Hidden2")},
		CustomFields: []models.CustomFieldValue{
			{CustomFieldId: hstFieldID, CurrencyValue: decPointer("13")},
			{CustomFieldId: 6, CurrencyValue: decPointer("9")},
		},
	}}
	options := SplitOptions{Categories: true, Excluded: map[uint]struct{}{hstFieldID: {}}}

	split := source.Split(receipts, options)
	for index := range split.Receipts {
		if split.Receipts[index].Categories[0].ID != 1 {
			split.Receipts[index].Categories = []models.Category{{Name: "(Restricted)"}}
		}
	}
	merged, perOrigin := source.Merge(split)

	if len(merged) != 2 {
		t.Fatalf("merged rows = %d, want 2", len(merged))
	}
	assertStrings(t, "amounts", amounts(merged), []string{"33.34", "66.66"})
	assertStrings(t, "restricted categories", categoryNames(merged[1]), []string{"(Restricted)"})
	if got := currencyValue(t, merged[1], 6).StringFixed(2); got != "6.00" {
		t.Errorf("merged Tip = %s, want 6.00", got)
	}
	if got := currencyValue(t, merged[1], hstFieldID).StringFixed(2); got != "13.00" {
		t.Errorf("merged excluded HST = %s, want 13.00, not a sum", got)
	}
	if !reflect.DeepEqual(perOrigin, []int{2}) {
		t.Errorf("perOrigin = %v, want [2]", perOrigin)
	}
}

// Merging writes a fresh pointer rather than adding into a shared one.
func TestMerge_DoesNotWriteThroughASharedValue(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:       dec("10"),
		Categories:   []models.Category{category(1, "A"), category(2, "B")},
		CustomFields: []models.CustomFieldValue{{CustomFieldId: hstFieldID, CurrencyValue: decPointer("2")}},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})
	firstValue := split.Receipts[0].CustomFields[0].CurrencyValue
	for index := range split.Receipts {
		split.Receipts[index].Categories = []models.Category{{Name: "(Restricted)"}}
	}
	merged, _ := source.Merge(split)

	if got := currencyValue(t, merged[0], hstFieldID).StringFixed(2); got != "2.00" {
		t.Errorf("merged HST = %s, want 2.00", got)
	}
	if firstValue.StringFixed(2) != "1.00" {
		t.Errorf("the first copy's own value changed to %s", firstValue)
	}
}

// Two receipts that look identical — both unsaved, so both id 0 — are two
// receipts, never folded together.
func TestMerge_NeverFoldsDifferentReceipts(t *testing.T) {
	source := mustNew(t)
	twin := models.Receipt{Amount: dec("10"), Categories: []models.Category{category(1, "A")}}

	merged, perOrigin := source.Merge(source.Split([]models.Receipt{twin, twin}, SplitOptions{Categories: true}))
	if len(merged) != 2 {
		t.Errorf("merged rows = %d, want 2", len(merged))
	}
	if !reflect.DeepEqual(perOrigin, []int{1, 1}) {
		t.Errorf("perOrigin = %v, want [1 1]", perOrigin)
	}
}

func TestMerge_ComparesCategoriesAndTagsAsSets(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount:     dec("10"),
		Categories: []models.Category{category(1, "A"), category(2, "B")},
		Tags:       []models.Tag{tag(1, "X"), tag(2, "Y")},
	}}

	split := source.Split(receipts, SplitOptions{Categories: true})
	// Same tag set in a different order on the second copy, same category.
	split.Receipts[0].Categories = []models.Category{{Name: "(Restricted)"}}
	split.Receipts[1].Categories = []models.Category{{Name: "(Restricted)"}}
	split.Receipts[1].Tags = []models.Tag{tag(2, "Y"), tag(1, "X")}

	merged, _ := source.Merge(split)
	assertStrings(t, "amounts", amounts(merged), []string{"10.00"})
}

func TestMerge_EmptyInput(t *testing.T) {
	source := mustNew(t)
	merged, perOrigin := source.Merge(source.Split(nil, SplitOptions{Categories: true}))
	if len(merged) != 0 || len(perOrigin) != 0 {
		t.Errorf("Merge of nothing = %v, %v", merged, perOrigin)
	}
}

// ---- through the engine -----------------------------------------------------

func TestSplit_FeedsTheEngine(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{
		{
			Amount:       dec("100"),
			PaidByUser:   models.User{DisplayName: "Dana"},
			Categories:   []models.Category{category(1, "Food"), category(2, "Fuel"), category(3, "Toys")},
			CustomFields: []models.CustomFieldValue{{CustomFieldId: hstFieldID, CurrencyValue: decPointer("13")}},
		},
		{
			Amount:     dec("20"),
			PaidByUser: models.User{DisplayName: "Dana"},
			Categories: []models.Category{category(1, "Food")},
		},
	}

	run := func(groupBy []reporting.FieldKey, options SplitOptions) reporting.ReportModel {
		t.Helper()
		merged, _ := source.Merge(source.Split(receipts, options))
		model, err := reporting.Run(reporting.ReportSpec{
			GroupBy: groupBy,
			Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: KeyCategory},
			Columns: []reporting.Column{
				{Name: "Category", Kind: reporting.ColumnLabel, Field: KeyCategory},
				{Name: "Count", Kind: reporting.ColumnAggregate, AggSrc: "COUNT()"},
				{Name: "Total", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
				{Name: "Hst", Kind: reporting.ColumnAggregate, AggSrc: "SUM(custom_1)"},
			},
			Subtotals:   true,
			GrandTotals: true,
		}, source.Catalog(), source.Rows(merged), reporting.MetaInput{})
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		return model
	}
	cell := func(cells []reporting.Cell, column string) string {
		for _, candidate := range cells {
			if candidate.Column == column {
				number, _ := candidate.Value().Decimal()
				return number.StringFixed(2)
			}
		}
		t.Fatalf("no cell %s", column)
		return ""
	}

	split := run(nil, SplitOptions{Categories: true})
	rows := split.Root.DetailRows
	got := []string{}
	for _, row := range rows {
		got = append(got, cell(row.Cells, "Total")+"/"+cell(row.Cells, "Count")+"/"+cell(row.Cells, "Hst"))
	}
	// Food 33.34 + 20, Fuel 33.33, Toys 33.33; HST 4.34/4.33/4.33.
	assertStrings(t, "split buckets", got, []string{"53.34/2.00/4.34", "33.33/1.00/4.33", "33.33/1.00/4.33"})
	if total := cell(split.GrandTotals, "Total"); total != "120.00" {
		t.Errorf("split grand total = %s, want the true 120.00", total)
	}
	if hst := cell(split.GrandTotals, "Hst"); hst != "13.00" {
		t.Errorf("split grand HST = %s, want the true 13.00", hst)
	}

	// Without the split the same report double counts, as it always has.
	whole := run(nil, SplitOptions{})
	if total := cell(whole.GrandTotals, "Total"); total != "320.00" {
		t.Errorf("unsplit grand total = %s, want 320.00", total)
	}

	// Grouping by payer above the split no longer double counts that payer.
	byPayer := run([]reporting.FieldKey{KeyPaidBy}, SplitOptions{Categories: true})
	if total := cell(byPayer.Root.Children[0].Subtotals, "Total"); total != "120.00" {
		t.Errorf("Dana subtotal = %s, want 120.00", total)
	}
}

func TestSplitOptions_Active(t *testing.T) {
	tests := []struct {
		options SplitOptions
		want    bool
	}{
		{SplitOptions{}, false},
		{SplitOptions{Excluded: map[uint]struct{}{1: {}}}, false},
		{SplitOptions{Categories: true}, true},
		{SplitOptions{Tags: true}, true},
		{SplitOptions{Categories: true, Tags: true}, true},
	}
	for _, test := range tests {
		if got := test.options.Active(); got != test.want {
			t.Errorf("%+v.Active() = %v, want %v", test.options, got, test.want)
		}
	}
}

func TestSplit_DuplicateTagCountsOnce(t *testing.T) {
	source := mustNew(t)
	receipts := []models.Receipt{{
		Amount: dec("9"),
		Tags:   []models.Tag{tag(2, "Sam"), tag(1, "Alex"), tag(2, "Sam"), tag(3, "Jo")},
	}}

	split := source.Split(receipts, SplitOptions{Tags: true})
	assertStrings(t, "amounts", amounts(split.Receipts), []string{"3.00", "3.00", "3.00"})
	assertStrings(t, "first copy", tagNames(split.Receipts[0]), []string{"Alex"})
}
