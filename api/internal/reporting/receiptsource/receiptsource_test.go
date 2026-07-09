package receiptsource

import (
	"errors"
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"

	"github.com/shopspring/decimal"
)

const (
	hstFieldID      uint = 1
	noteFieldID     uint = 2
	childFieldID    uint = 3
	dueDateFieldID  uint = 4
	reimbursedField uint = 5
)

func dec(literal string) decimal.Decimal {
	return decimal.RequireFromString(literal)
}

func testCustomFields() []models.CustomField {
	return []models.CustomField{
		{BaseModel: models.BaseModel{ID: hstFieldID}, Name: "HST", Type: models.CURRENCY},
		{BaseModel: models.BaseModel{ID: noteFieldID}, Name: "Note", Type: models.TEXT},
		{
			BaseModel: models.BaseModel{ID: childFieldID},
			Name:      "Child",
			Type:      models.SELECT,
			Options: []models.CustomFieldOption{
				{BaseModel: models.BaseModel{ID: 10}, Value: "Alex", CustomFieldId: childFieldID},
				{BaseModel: models.BaseModel{ID: 11}, Value: "Sam", CustomFieldId: childFieldID},
			},
		},
		{BaseModel: models.BaseModel{ID: dueDateFieldID}, Name: "Due Date", Type: models.DATE},
		{BaseModel: models.BaseModel{ID: reimbursedField}, Name: "Reimbursed", Type: models.BOOLEAN},
	}
}

func mustNew(t *testing.T) Source {
	t.Helper()

	source, err := New(testCustomFields())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return source
}

func TestCustomFieldKey(t *testing.T) {
	tests := []struct {
		id   uint
		want reporting.FieldKey
	}{
		{1, "custom_1"},
		{0, "custom_0"},
		{4294967295, "custom_4294967295"},
	}

	for _, test := range tests {
		if got := CustomFieldKey(test.id); got != test.want {
			t.Errorf("CustomFieldKey(%d) = %q, want %q", test.id, got, test.want)
		}
	}
}

// A currency custom field is a measure, so a tax field is summable without any
// special case. Everything else cuts the data.
func TestSource_CatalogTypesCustomFields(t *testing.T) {
	catalog := mustNew(t).Catalog()

	tests := []struct {
		key      reporting.FieldKey
		label    string
		dataType reporting.DataType
		role     reporting.Role
	}{
		{"custom_1", "HST", reporting.TypeCurrency, reporting.RoleMeasure},
		{"custom_2", "Note", reporting.TypeString, reporting.RoleDimension},
		{"custom_3", "Child", reporting.TypeString, reporting.RoleDimension},
		{"custom_4", "Due Date", reporting.TypeDate, reporting.RoleDimension},
		{"custom_5", "Reimbursed", reporting.TypeBool, reporting.RoleDimension},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			field, exists := catalog.Get(test.key)
			if !exists {
				t.Fatalf("catalog is missing %s", test.key)
			}
			if field.Label != test.label {
				t.Errorf("label = %q, want %q", field.Label, test.label)
			}
			if field.DataType != test.dataType {
				t.Errorf("dataType = %v, want %v", field.DataType, test.dataType)
			}
			if field.Role() != test.role {
				t.Errorf("role = %v, want %v", field.Role(), test.role)
			}
			if field.Multi {
				t.Errorf("custom fields are single-valued")
			}
		})
	}
}

func TestSource_CatalogBuiltins(t *testing.T) {
	catalog := mustNew(t).Catalog()

	tests := []struct {
		key      reporting.FieldKey
		dataType reporting.DataType
		multi    bool
	}{
		{KeyReceiptID, reporting.TypeNumber, false},
		{KeyName, reporting.TypeString, false},
		{KeyAmount, reporting.TypeCurrency, false},
		{KeyDate, reporting.TypeDate, false},
		{KeyResolvedDate, reporting.TypeDate, false},
		{KeyCreatedAt, reporting.TypeDate, false},
		{KeyStatus, reporting.TypeString, false},
		{KeyPaidBy, reporting.TypeString, false},
		{KeyGroup, reporting.TypeString, false},
		{KeyCategory, reporting.TypeString, true},
		{KeyTag, reporting.TypeString, true},
	}

	for _, test := range tests {
		t.Run(string(test.key), func(t *testing.T) {
			field, exists := catalog.Get(test.key)
			if !exists {
				t.Fatalf("catalog is missing %s", test.key)
			}
			if field.DataType != test.dataType {
				t.Errorf("dataType = %v, want %v", field.DataType, test.dataType)
			}
			if field.Multi != test.multi {
				t.Errorf("multi = %v, want %v", field.Multi, test.multi)
			}
		})
	}
}

func TestNew_RejectsDuplicateCustomFieldIds(t *testing.T) {
	_, err := New([]models.CustomField{
		{BaseModel: models.BaseModel{ID: 1}, Name: "HST", Type: models.CURRENCY},
		{BaseModel: models.BaseModel{ID: 1}, Name: "GST", Type: models.CURRENCY},
	})

	if !errors.Is(err, reporting.ErrDuplicateField) {
		t.Errorf("New() error = %v, want %v", err, reporting.ErrDuplicateField)
	}
}

func TestNew_NoCustomFields(t *testing.T) {
	source, err := New(nil)
	if err != nil {
		t.Fatalf("New(nil) error = %v", err)
	}
	if _, exists := source.Catalog().Get(KeyAmount); !exists {
		t.Errorf("builtins are missing from a catalog with no custom fields")
	}
	if _, exists := source.Catalog().Get("custom_1"); exists {
		t.Errorf("catalog invented a custom field")
	}
}

// The catalog must not depend on the order the definitions arrived in.
func TestNew_CatalogIsIndependentOfDefinitionOrder(t *testing.T) {
	forward := testCustomFields()
	reversed := testCustomFields()
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}

	first, err := New(forward)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	second, err := New(reversed)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	for _, key := range []reporting.FieldKey{"custom_1", "custom_3", "custom_5"} {
		one, _ := first.Catalog().Get(key)
		two, _ := second.Catalog().Get(key)
		if one != two {
			t.Errorf("%s resolved differently: %+v vs %+v", key, one, two)
		}
	}
}

func fullReceipt() models.Receipt {
	resolved := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	hst := dec("15.60")
	note := "office supplies"
	option := uint(11)
	due := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	reimbursed := true

	return models.Receipt{
		BaseModel:    models.BaseModel{ID: 7, CreatedAt: time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC)},
		Name:         "Staples",
		Amount:       dec("120.00"),
		Date:         time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
		ResolvedDate: &resolved,
		Status:       models.RESOLVED,
		PaidByUser:   models.User{DisplayName: "Dana"},
		Group:        models.Group{Name: "Household"},
		Categories:   []models.Category{{Name: "Clothing"}, {Name: "Medical"}},
		Tags:         []models.Tag{{Name: "Alex"}},
		CustomFields: []models.CustomFieldValue{
			{CustomFieldId: hstFieldID, CurrencyValue: &hst},
			{CustomFieldId: noteFieldID, StringValue: &note},
			{CustomFieldId: childFieldID, SelectValue: &option},
			{CustomFieldId: dueDateFieldID, DateValue: &due},
			{CustomFieldId: reimbursedField, BooleanValue: &reimbursed},
		},
	}
}

func TestSource_RowsResolvesEveryField(t *testing.T) {
	rows := mustNew(t).Rows([]models.Receipt{fullReceipt()})
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]

	t.Run("receipt id", func(t *testing.T) {
		number, isNumber := row.Measure(KeyReceiptID).Decimal()
		if !isNumber || !number.Equal(dec("7")) {
			t.Errorf("receipt_id = %v", row.Measure(KeyReceiptID))
		}
	})

	t.Run("amount", func(t *testing.T) {
		number, isNumber := row.Measure(KeyAmount).Decimal()
		if !isNumber || !number.Equal(dec("120.00")) {
			t.Errorf("amount = %v", row.Measure(KeyAmount))
		}
	})

	t.Run("currency custom field", func(t *testing.T) {
		number, isNumber := row.Measure("custom_1").Decimal()
		if !isNumber || !number.Equal(dec("15.60")) {
			t.Errorf("custom_1 = %v", row.Measure("custom_1"))
		}
	})

	strings := []struct {
		key  reporting.FieldKey
		want string
	}{
		{KeyName, "Staples"},
		{KeyStatus, "RESOLVED"},
		{KeyPaidBy, "Dana"},
		{KeyGroup, "Household"},
		{KeyTag, "Alex"},
		{"custom_2", "office supplies"},
		// A select stores an option id; the row carries the text.
		{"custom_3", "Sam"},
	}
	for _, test := range strings {
		t.Run(string(test.key), func(t *testing.T) {
			text, isText := row.Measure(test.key).Text()
			if !isText || text != test.want {
				t.Errorf("%s = %v, want %q", test.key, row.Measure(test.key), test.want)
			}
		})
	}

	dates := []struct {
		key  reporting.FieldKey
		want time.Time
	}{
		{KeyDate, time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)},
		{KeyResolvedDate, time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)},
		{KeyCreatedAt, time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC)},
		{"custom_4", time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, test := range dates {
		t.Run(string(test.key), func(t *testing.T) {
			instant, isDate := row.Measure(test.key).Time()
			if !isDate || !instant.Equal(test.want) {
				t.Errorf("%s = %v, want %v", test.key, row.Measure(test.key), test.want)
			}
		})
	}

	t.Run("boolean custom field", func(t *testing.T) {
		value, isBool := row.Measure("custom_5").Boolean()
		if !isBool || !value {
			t.Errorf("custom_5 = %v, want true", row.Measure("custom_5"))
		}
	})

	// Categories are multi-valued, and keep the order the receipt carried them.
	t.Run("categories fan out", func(t *testing.T) {
		values := row.Get(KeyCategory)
		if len(values) != 2 {
			t.Fatalf("got %d categories, want 2", len(values))
		}
		first, _ := values[0].Text()
		second, _ := values[1].Text()
		if first != "Clothing" || second != "Medical" {
			t.Errorf("categories = %v, %v", first, second)
		}
	})
}

// A payer without a display name falls back to the name they log in with.
func TestSource_PaidByFallsBackToUsername(t *testing.T) {
	receipt := models.Receipt{PaidByUser: models.User{Username: "dana.smith"}}
	row := mustNew(t).Rows([]models.Receipt{receipt})[0]

	text, isText := row.Measure(KeyPaidBy).Text()
	if !isText || text != "dana.smith" {
		t.Errorf("paid_by = %v, want dana.smith", row.Measure(KeyPaidBy))
	}
}

// An association the caller did not preload resolves to no value, which the
// engine reports as (None) rather than as an error.
func TestSource_UnloadedAssociationsBecomeNone(t *testing.T) {
	row := mustNew(t).Rows([]models.Receipt{{Name: "Bare", Amount: dec("1.00")}})[0]

	absent := []reporting.FieldKey{KeyPaidBy, KeyGroup, KeyResolvedDate, KeyCategory, KeyTag, "custom_1"}
	for _, key := range absent {
		t.Run(string(key), func(t *testing.T) {
			if len(row.Get(key)) != 0 {
				t.Errorf("%s resolved to %v, want no value", key, row.Get(key))
			}
			if !row.Measure(key).IsNull() {
				t.Errorf("%s measures as %v, want null", key, row.Measure(key))
			}
		})
	}
}

func TestSource_MissingCustomFieldValues(t *testing.T) {
	receipt := models.Receipt{
		CustomFields: []models.CustomFieldValue{
			// The right field, but its column is empty.
			{CustomFieldId: hstFieldID},
			{CustomFieldId: noteFieldID},
			{CustomFieldId: childFieldID},
			{CustomFieldId: dueDateFieldID},
			{CustomFieldId: reimbursedField},
		},
	}

	row := mustNew(t).Rows([]models.Receipt{receipt})[0]

	for _, key := range []reporting.FieldKey{"custom_1", "custom_2", "custom_3", "custom_4", "custom_5"} {
		if len(row.Get(key)) != 0 {
			t.Errorf("%s resolved to %v, want no value", key, row.Get(key))
		}
	}
}

// A select value pointing at an option that was not loaded, or that no longer
// exists, resolves to nothing rather than to a bare id.
func TestSource_SelectWithUnknownOption(t *testing.T) {
	missing := uint(99)
	receipt := models.Receipt{CustomFields: []models.CustomFieldValue{
		{CustomFieldId: childFieldID, SelectValue: &missing},
	}}

	row := mustNew(t).Rows([]models.Receipt{receipt})[0]
	if len(row.Get("custom_3")) != 0 {
		t.Errorf("custom_3 = %v, want no value", row.Get("custom_3"))
	}
}

func TestSource_SelectWithUnloadedOptions(t *testing.T) {
	source, err := New([]models.CustomField{
		{BaseModel: models.BaseModel{ID: childFieldID}, Name: "Child", Type: models.SELECT},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	option := uint(10)
	receipt := models.Receipt{CustomFields: []models.CustomFieldValue{
		{CustomFieldId: childFieldID, SelectValue: &option},
	}}

	row := source.Rows([]models.Receipt{receipt})[0]
	if len(row.Get("custom_3")) != 0 {
		t.Errorf("custom_3 = %v, want no value", row.Get("custom_3"))
	}
}

// A value for a field the source does not know about is skipped: the engine
// could not reference it anyway.
func TestSource_UnknownCustomFieldValueIsSkipped(t *testing.T) {
	amount := dec("5.00")
	receipt := models.Receipt{CustomFields: []models.CustomFieldValue{
		{CustomFieldId: 404, CurrencyValue: &amount},
	}}

	row := mustNew(t).Rows([]models.Receipt{receipt})[0]
	if len(row.Get("custom_404")) != 0 {
		t.Errorf("custom_404 = %v, want no value", row.Get("custom_404"))
	}
}

// Where a receipt somehow carries two values for one field, the first that
// resolves wins, so the row is deterministic.
func TestSource_DuplicateCustomFieldValues(t *testing.T) {
	first, second := dec("1.00"), dec("2.00")
	receipt := models.Receipt{CustomFields: []models.CustomFieldValue{
		{CustomFieldId: hstFieldID, CurrencyValue: &first},
		{CustomFieldId: hstFieldID, CurrencyValue: &second},
	}}

	row := mustNew(t).Rows([]models.Receipt{receipt})[0]

	values := row.Get("custom_1")
	if len(values) != 1 {
		t.Fatalf("got %d values, want 1", len(values))
	}
	number, _ := values[0].Decimal()
	if !number.Equal(dec("1.00")) {
		t.Errorf("custom_1 = %s, want 1.00", number)
	}
}

// An empty first value must not hide a real second one.
func TestSource_FirstResolvableCustomFieldValueWins(t *testing.T) {
	amount := dec("2.00")
	receipt := models.Receipt{CustomFields: []models.CustomFieldValue{
		{CustomFieldId: hstFieldID},
		{CustomFieldId: hstFieldID, CurrencyValue: &amount},
	}}

	row := mustNew(t).Rows([]models.Receipt{receipt})[0]

	number, isNumber := row.Measure("custom_1").Decimal()
	if !isNumber || !number.Equal(dec("2.00")) {
		t.Errorf("custom_1 = %v, want 2.00", row.Measure("custom_1"))
	}
}

func TestSource_RowsPreservesOrder(t *testing.T) {
	receipts := []models.Receipt{
		{BaseModel: models.BaseModel{ID: 3}},
		{BaseModel: models.BaseModel{ID: 1}},
		{BaseModel: models.BaseModel{ID: 2}},
	}

	rows := mustNew(t).Rows(receipts)

	want := []string{"3", "1", "2"}
	for index, expected := range want {
		number, _ := rows[index].Measure(KeyReceiptID).Decimal()
		if number.String() != expected {
			t.Errorf("row %d receipt_id = %s, want %s", index, number, expected)
		}
	}
}

func TestSource_RowsOnNoReceipts(t *testing.T) {
	if rows := mustNew(t).Rows(nil); len(rows) != 0 {
		t.Errorf("Rows(nil) = %v, want empty", rows)
	}
}

// The end-to-end shape: real receipts through the source into the engine,
// reproducing the design document's worked example.
func TestSource_FeedsTheEngine(t *testing.T) {
	source := mustNew(t)

	hst := func(literal string) *decimal.Decimal {
		value := dec(literal)
		return &value
	}
	receipt := func(payer, tag, category, amount string, tax *decimal.Decimal) models.Receipt {
		built := models.Receipt{
			Amount:     dec(amount),
			PaidByUser: models.User{DisplayName: payer},
			Tags:       []models.Tag{{Name: tag}},
			Categories: []models.Category{{Name: category}},
		}
		if tax != nil {
			built.CustomFields = []models.CustomFieldValue{{CustomFieldId: hstFieldID, CurrencyValue: tax}}
		}
		return built
	}

	receipts := []models.Receipt{
		receipt("Dana", "Alex", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Alex", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Alex", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Alex", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Alex", "Medical", "50.00", nil),
		receipt("Dana", "Alex", "Medical", "30.00", nil),
		receipt("Dana", "Sam", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Sam", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Sam", "Clothing", "30.00", hst("3.90")),
		receipt("Dana", "Sam", "Mileage", "30.00", nil),
	}

	spec := reporting.ReportSpec{
		GroupBy: []reporting.FieldKey{KeyPaidBy, KeyTag},
		Detail:  reporting.DetailSpec{Mode: reporting.DetailAggregate, By: KeyCategory},
		Columns: []reporting.Column{
			{Name: "Category", Kind: reporting.ColumnLabel, Field: KeyCategory},
			{Name: "Count", Kind: reporting.ColumnAggregate, AggSrc: "COUNT()"},
			{Name: "Subtotal", Kind: reporting.ColumnAggregate, AggSrc: "SUM(amount)"},
			{Name: "Hst", Kind: reporting.ColumnAggregate, AggSrc: "SUM(custom_1)"},
			{Name: "Total", Kind: reporting.ColumnArithmetic, Expr: "Subtotal + Hst"},
			{Name: "AvgPerReceipt", Label: "Avg/Receipt", Kind: reporting.ColumnArithmetic, Expr: "ROUND(Total / Count, 2)"},
		},
		Subtotals:   true,
		GrandTotals: true,
	}

	model, err := reporting.Run(spec, source.Catalog(), source.Rows(receipts), reporting.MetaInput{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	find := func(cells []reporting.Cell, column string) string {
		for _, candidate := range cells {
			if candidate.Column == column {
				return candidate.Value().String()
			}
		}
		t.Fatalf("no cell for %s", column)
		return ""
	}

	dana := model.Root.Children[0]
	alex, sam := dana.Children[0], dana.Children[1]

	if got := find(alex.Subtotals, "Subtotal"); got != "200" {
		t.Errorf("Alex subtotal = %s, want 200", got)
	}
	if got := find(alex.Subtotals, "AvgPerReceipt"); got != "35.93" {
		t.Errorf("Alex average = %s, want 35.93", got)
	}
	if got := find(sam.Subtotals, "AvgPerReceipt"); got != "32.93" {
		t.Errorf("Sam average = %s, want 32.93", got)
	}
	if got := find(model.GrandTotals, "Total"); got != "347.3" {
		t.Errorf("grand total = %s, want 347.3", got)
	}
	if got := find(model.GrandTotals, "AvgPerReceipt"); got != "34.73" {
		t.Errorf("grand average = %s, want 34.73", got)
	}

	// Medical carried no tax at all, and sums to zero rather than to an empty
	// cell.
	if got := find(alex.DetailRows[1].Cells, "Hst"); got != "0" {
		t.Errorf("Medical HST = %s, want 0", got)
	}
}
