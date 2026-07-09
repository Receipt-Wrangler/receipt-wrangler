package reporting

import (
	"errors"
	"testing"
)

func TestDataType_IsNumeric(t *testing.T) {
	tests := []struct {
		dataType DataType
		want     bool
	}{
		{TypeString, false},
		{TypeNumber, true},
		{TypeCurrency, true},
		{TypeDate, false},
		{TypeBool, false},
	}

	for _, test := range tests {
		t.Run(test.dataType.String(), func(t *testing.T) {
			if got := test.dataType.IsNumeric(); got != test.want {
				t.Errorf("%v.IsNumeric() = %v, want %v", test.dataType, got, test.want)
			}
		})
	}
}

// You group by a tag; you sum a dollar amount. The role falls out of the type so
// a template can never do the reverse.
func TestRoleForDataType(t *testing.T) {
	tests := []struct {
		dataType DataType
		want     Role
	}{
		{TypeString, RoleDimension},
		{TypeDate, RoleDimension},
		{TypeBool, RoleDimension},
		{TypeNumber, RoleMeasure},
		{TypeCurrency, RoleMeasure},
	}

	for _, test := range tests {
		t.Run(test.dataType.String(), func(t *testing.T) {
			if got := RoleForDataType(test.dataType); got != test.want {
				t.Errorf("RoleForDataType(%v) = %v, want %v", test.dataType, got, test.want)
			}
		})
	}
}

func TestFieldRef_Role(t *testing.T) {
	currency := FieldRef{Key: "custom_1", Label: "HST", DataType: TypeCurrency}
	if got := currency.Role(); got != RoleMeasure {
		t.Errorf("currency custom field role = %v, want %v", got, RoleMeasure)
	}

	selectField := FieldRef{Key: "custom_2", Label: "Child", DataType: TypeString}
	if got := selectField.Role(); got != RoleDimension {
		t.Errorf("select custom field role = %v, want %v", got, RoleDimension)
	}
}

func TestNewFieldCatalog(t *testing.T) {
	amount := FieldRef{Key: "amount", Label: "Amount", DataType: TypeCurrency}
	tag := FieldRef{Key: "tag", Label: "Tag", DataType: TypeString, Multi: true}

	t.Run("builds and looks up fields", func(t *testing.T) {
		catalog, err := NewFieldCatalog(amount, tag)
		if err != nil {
			t.Fatalf("NewFieldCatalog() error = %v, want nil", err)
		}

		got, exists := catalog.Get("amount")
		if !exists {
			t.Fatalf("Get(\"amount\") not found")
		}
		if got.Label != "Amount" || got.DataType != TypeCurrency || got.Multi {
			t.Errorf("Get(\"amount\") = %+v, want the amount field", got)
		}

		multi, exists := catalog.Get("tag")
		if !exists || !multi.Multi {
			t.Errorf("Get(\"tag\") = %+v, %v; want a multi-value field", multi, exists)
		}
	})

	t.Run("empty catalog", func(t *testing.T) {
		catalog, err := NewFieldCatalog()
		if err != nil {
			t.Fatalf("NewFieldCatalog() error = %v, want nil", err)
		}
		if _, exists := catalog.Get("amount"); exists {
			t.Errorf("empty catalog returned a field")
		}
	})

	t.Run("missing key", func(t *testing.T) {
		catalog, _ := NewFieldCatalog(amount)
		if _, exists := catalog.Get("nope"); exists {
			t.Errorf("Get(\"nope\") found a field")
		}
	})

	t.Run("zero catalog does not panic", func(t *testing.T) {
		var catalog FieldCatalog
		if _, exists := catalog.Get("amount"); exists {
			t.Errorf("zero catalog returned a field")
		}
	})

	t.Run("rejects an empty key", func(t *testing.T) {
		_, err := NewFieldCatalog(FieldRef{Label: "Nameless", DataType: TypeString})
		if !errors.Is(err, ErrEmptyFieldKey) {
			t.Errorf("NewFieldCatalog() error = %v, want %v", err, ErrEmptyFieldKey)
		}
	})

	t.Run("rejects a duplicate key", func(t *testing.T) {
		_, err := NewFieldCatalog(amount, FieldRef{Key: "amount", Label: "Other", DataType: TypeNumber})
		if !errors.Is(err, ErrDuplicateField) {
			t.Errorf("NewFieldCatalog() error = %v, want %v", err, ErrDuplicateField)
		}
	})
}
