package receiptsource

import (
	"math/big"
	"sort"
	"strconv"
	"strings"

	"receipt-wrangler/api/internal/models"

	"github.com/shopspring/decimal"
)

// minimumSplitScale is the fewest decimal places a split share is allocated in.
// A receipt amount is stored as decimal(10,2), so cents are the natural unit;
// a currency custom field is stored as text and may carry more places, in which
// case its own scale wins so no precision it was entered with is thrown away.
const minimumSplitScale = 2

// SplitOptions says which multi-valued dimensions a report divides a receipt's
// money across, instead of attributing the whole of it to every bucket.
//
// The engine's default is to fan a receipt with two categories out into both in
// full, which double-counts it on purpose (see Row.dimensionValues). Splitting
// is the opt-in alternative: the receipt becomes one copy per category, each
// carrying an equal share, so the buckets add back up to what was spent.
type SplitOptions struct {
	Categories bool
	Tags       bool

	// Excluded names the currency custom fields that are not divided: every
	// copy carries the field's whole value, exactly as without a split.
	Excluded map[uint]struct{}
}

// Active reports whether the options split anything at all.
func (o SplitOptions) Active() bool { return o.Categories || o.Tags }

// SplitReceipts is what Split produced: the copies, plus which input receipt
// each one came from. The origin travels separately rather than being read off
// the receipt id because an id does not identify an input — unsaved receipts all
// carry id 0, and the same receipt can be reported twice.
type SplitReceipts struct {
	Receipts []models.Receipt

	origins     []int
	inputLength int
	options     SplitOptions
}

// Split divides each receipt's money equally across its categories and/or tags.
//
// A receipt with N distinct categories (when Categories is set) and M distinct
// tags (when Tags is set) becomes N×M copies, each holding one category and one
// tag. A dimension that is not split, or that the receipt holds at most one
// value for, keeps its whole list and counts as one. A receipt with nothing to
// split comes through as a single, unchanged copy.
//
// What is divided is the amount and every currency custom field not named in
// Excluded. Everything else — name, dates, the other custom fields — is the
// same on every copy, which is what lets a copy still be grouped by them.
//
// Shares are whole units of the value's scale (cents for an amount), allocated
// so they add back up exactly; see allocateGrid. The order the leftover units go
// in is by category and tag id, never by input order, so the same receipt
// splits the same way every time.
//
// Copies own their category, tag and custom-field slices and their currency
// values, so the restricted-category substitution that follows can rewrite one
// copy without reaching another, or the caller's input.
func (s Source) Split(receipts []models.Receipt, options SplitOptions) SplitReceipts {
	result := SplitReceipts{
		Receipts:    make([]models.Receipt, 0, len(receipts)),
		origins:     make([]int, 0, len(receipts)),
		inputLength: len(receipts),
		options:     options,
	}

	for origin := range receipts {
		copies := s.splitReceipt(&receipts[origin], options)
		for range copies {
			result.origins = append(result.origins, origin)
		}
		result.Receipts = append(result.Receipts, copies...)
	}

	return result
}

func (s Source) splitReceipt(receipt *models.Receipt, options SplitOptions) []models.Receipt {
	categoryLists := [][]models.Category{receipt.Categories}
	if options.Categories {
		if distinct := distinctCategories(receipt.Categories); len(distinct) > 1 {
			categoryLists = make([][]models.Category, len(distinct))
			for index, category := range distinct {
				categoryLists[index] = []models.Category{category}
			}
		}
	}

	tagLists := [][]models.Tag{receipt.Tags}
	if options.Tags {
		if distinct := distinctTags(receipt.Tags); len(distinct) > 1 {
			tagLists = make([][]models.Tag, len(distinct))
			for index, tag := range distinct {
				tagLists[index] = []models.Tag{tag}
			}
		}
	}

	rows, cols := len(categoryLists), len(tagLists)
	if rows == 1 && cols == 1 {
		return []models.Receipt{*receipt}
	}

	amounts := allocateGrid(receipt.Amount, rows, cols)

	// One grid per divided custom field value, by its position on the receipt.
	customFieldShares := make(map[int][][]decimal.Decimal)
	for index, customFieldValue := range receipt.CustomFields {
		if s.isSplitCustomFieldValue(customFieldValue, options) {
			customFieldShares[index] = allocateGrid(*customFieldValue.CurrencyValue, rows, cols)
		}
	}

	copies := make([]models.Receipt, 0, rows*cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			copied := *receipt
			copied.Amount = amounts[row][col]
			copied.Categories = append([]models.Category(nil), categoryLists[row]...)
			copied.Tags = append([]models.Tag(nil), tagLists[col]...)

			copied.CustomFields = append([]models.CustomFieldValue(nil), receipt.CustomFields...)
			for index, shares := range customFieldShares {
				share := shares[row][col]
				copied.CustomFields[index].CurrencyValue = &share
			}

			copies = append(copies, copied)
		}
	}

	return copies
}

// isSplitCustomFieldValue reports whether a custom field value is one Split
// divides: a currency field, not excluded, with a value to divide. The type is
// read from the source's definitions, since the value's own CustomField need not
// have been preloaded; a field the catalog does not know is left alone.
func (s Source) isSplitCustomFieldValue(customFieldValue models.CustomFieldValue, options SplitOptions) bool {
	if customFieldValue.CurrencyValue == nil {
		return false
	}
	if s.types[customFieldValue.CustomFieldId] != models.CURRENCY {
		return false
	}
	_, excluded := options.Excluded[customFieldValue.CustomFieldId]
	return !excluded
}

// Merge folds back together the copies of one receipt that became
// indistinguishable after Split — in practice, after the restricted-category
// substitution turned two hidden categories into the same (Restricted) marker.
// Two copies of one origin whose category names and tag names now match are one
// row: their divided values are summed, and everything else is taken from the
// first. That keeps a report at one row per receipt per bucket, as it is without
// a split, while the (Restricted) bucket still receives every hidden share.
//
// It returns the merged receipts in first-occurrence order, and how many each
// input receipt became, indexed by its position in the input to Split.
func (s Source) Merge(split SplitReceipts) ([]models.Receipt, []int) {
	merged := make([]models.Receipt, 0, len(split.Receipts))
	perOrigin := make([]int, split.inputLength)
	indexByKey := make(map[string]int, len(split.Receipts))

	for index := range split.Receipts {
		receipt := split.Receipts[index]
		origin := split.origins[index]
		key := mergeKey(origin, receipt)

		existing, found := indexByKey[key]
		if !found {
			indexByKey[key] = len(merged)
			merged = append(merged, receipt)
			perOrigin[origin]++
			continue
		}

		s.mergeInto(&merged[existing], receipt, split.options)
	}

	return merged, perOrigin
}

// mergeInto adds the divided values of from onto into. Copies of one receipt
// carry their custom fields in the same order, so a value's position names the
// same field on both.
func (s Source) mergeInto(into *models.Receipt, from models.Receipt, options SplitOptions) {
	into.Amount = into.Amount.Add(from.Amount)

	for index, customFieldValue := range into.CustomFields {
		if index >= len(from.CustomFields) || !s.isSplitCustomFieldValue(customFieldValue, options) {
			continue
		}
		other := from.CustomFields[index].CurrencyValue
		if other == nil {
			continue
		}
		// A fresh pointer, so no other receipt holding the old one sees the sum.
		sum := customFieldValue.CurrencyValue.Add(*other)
		into.CustomFields[index].CurrencyValue = &sum
	}
}

func mergeKey(origin int, receipt models.Receipt) string {
	categoryNames := make([]string, 0, len(receipt.Categories))
	for _, category := range receipt.Categories {
		categoryNames = append(categoryNames, category.Name)
	}
	tagNames := make([]string, 0, len(receipt.Tags))
	for _, tag := range receipt.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	sort.Strings(categoryNames)
	sort.Strings(tagNames)

	var key strings.Builder
	key.WriteString(strconv.Itoa(origin))
	key.WriteString("\x00")
	key.WriteString(strings.Join(categoryNames, "\x01"))
	key.WriteString("\x00")
	key.WriteString(strings.Join(tagNames, "\x01"))
	return key.String()
}

// distinctCategories returns the receipt's categories once each, ordered by id
// (then name, for unsaved ones), so the leftover units of a split always go to
// the same categories regardless of the order the database returned them in.
func distinctCategories(categories []models.Category) []models.Category {
	distinct := make([]models.Category, 0, len(categories))
	seen := make(map[string]struct{}, len(categories))
	for _, category := range categories {
		key := identityKey(category.ID, category.Name)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		distinct = append(distinct, category)
	}
	sort.SliceStable(distinct, func(i, j int) bool {
		if distinct[i].ID != distinct[j].ID {
			return distinct[i].ID < distinct[j].ID
		}
		return distinct[i].Name < distinct[j].Name
	})
	return distinct
}

func distinctTags(tags []models.Tag) []models.Tag {
	distinct := make([]models.Tag, 0, len(tags))
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		key := identityKey(tag.ID, tag.Name)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		distinct = append(distinct, tag)
	}
	sort.SliceStable(distinct, func(i, j int) bool {
		if distinct[i].ID != distinct[j].ID {
			return distinct[i].ID < distinct[j].ID
		}
		return distinct[i].Name < distinct[j].Name
	})
	return distinct
}

// identityKey distinguishes categories (or tags) the way the engine buckets
// them: a saved one by its id, an unsaved one (id 0) by its name.
func identityKey(id uint, name string) string {
	if id != 0 {
		return "id:" + strconv.FormatUint(uint64(id), 10)
	}
	return "name:" + name
}

// allocateGrid divides value into rows×cols shares that add back up to it
// exactly.
//
// It works in whole units of the value's scale (at least cents) with integer
// arithmetic, never Div, which would read shopspring's process-wide
// DivisionPrecision. Every cell gets the quotient, and the remainder is handed
// out one unit at a time:
//
//   - row i gets K/rows extra units, plus one more for the first K mod rows rows;
//   - those units go to columns cyclically, the column pointer carrying on from
//     one row to the next.
//
// Allocating the grid in one go, rather than splitting by category and then
// splitting each share by tag, is what keeps both views consistent: every row
// sums to what an equal split across rows alone would give, every column to
// what a split across columns alone would give. Rounding in two passes breaks
// that — 0.04 over three categories is 0.02/0.01/0.01, and splitting each of
// those across two tags hands one tag 0.03 and the other 0.01.
//
// No cell gets two extra units, since a row receives at most cols of them. A
// negative value is split by magnitude and re-signed, so -100 over three is
// -33.34/-33.33/-33.33, the mirror image of 100.
func allocateGrid(value decimal.Decimal, rows, cols int) [][]decimal.Decimal {
	scale := int32(minimumSplitScale)
	if places := -value.Exponent(); places > scale {
		scale = places
	}

	units := value.Shift(scale).BigInt()
	negative := units.Sign() < 0
	units.Abs(units)

	cells := big.NewInt(int64(rows * cols))
	quotient, remainder := new(big.Int).QuoRem(units, cells, new(big.Int))
	extra := int(remainder.Int64())

	grid := make([][]decimal.Decimal, rows)
	column := 0
	for row := 0; row < rows; row++ {
		rowExtra := extra / rows
		if row < extra%rows {
			rowExtra++
		}

		extraByColumn := make([]bool, cols)
		for given := 0; given < rowExtra; given++ {
			extraByColumn[column] = true
			column = (column + 1) % cols
		}

		grid[row] = make([]decimal.Decimal, cols)
		for col := 0; col < cols; col++ {
			share := new(big.Int).Set(quotient)
			if extraByColumn[col] {
				share.Add(share, big.NewInt(1))
			}
			if negative {
				share.Neg(share)
			}
			grid[row][col] = decimal.NewFromBigInt(share, -scale)
		}
	}

	return grid
}
