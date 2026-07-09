package reporting

import (
	"errors"
	"fmt"
)

var (
	// ErrEmptyFieldKey is returned when a catalog is built with a keyless field.
	ErrEmptyFieldKey = errors.New("field key must not be empty")

	// ErrDuplicateField is returned when a catalog is built with two fields
	// sharing a key.
	ErrDuplicateField = errors.New("duplicate field key in catalog")
)

// DataType is the type of value a field resolves to. It determines the field's
// role, and travels onto a column descriptor so renderers know how to format.
type DataType uint8

const (
	TypeString DataType = iota
	TypeNumber
	TypeCurrency
	TypeDate
	TypeBool
)

func (d DataType) String() string {
	switch d {
	case TypeString:
		return "string"
	case TypeNumber:
		return "number"
	case TypeCurrency:
		return "currency"
	case TypeDate:
		return "date"
	case TypeBool:
		return "bool"
	}
	return "unknown"
}

// IsNumeric reports whether values of this type can be summed and can take part
// in arithmetic. Currency is a number that renderers format differently.
func (d DataType) IsNumeric() bool {
	return d == TypeNumber || d == TypeCurrency
}

// Role is what a field may be used for. It is derived from the data type rather
// than declared, so a template can never group by a dollar amount or sum a
// status.
type Role uint8

const (
	// RoleDimension is a way to cut the data: a groupBy level, or the dimension
	// an aggregated detail row is keyed on.
	RoleDimension Role = iota

	// RoleMeasure is a thing to measure: the input to an aggregate column.
	RoleMeasure
)

func (r Role) String() string {
	switch r {
	case RoleDimension:
		return "dimension"
	case RoleMeasure:
		return "measure"
	}
	return "unknown"
}

// RoleForDataType derives a field's role from its type. Numbers measure;
// everything else cuts.
func RoleForDataType(dataType DataType) Role {
	if dataType.IsNumeric() {
		return RoleMeasure
	}
	return RoleDimension
}

// FieldKey identifies a field within a catalog. Keys must be valid expression
// identifiers so that a formula can name one without quoting.
type FieldKey string

// FieldRef describes one field a report may reference. Built-in fields and
// custom fields are the same abstraction, which is what makes a custom field
// usable anywhere a built-in is.
type FieldRef struct {
	Key      FieldKey
	Label    string
	DataType DataType

	// Multi marks a field that can resolve to several values at once, such as a
	// receipt's categories or tags. Grouping on one fans the row out into every
	// bucket, so the row double-counts.
	Multi bool
}

// Role reports whether the field cuts the data or measures it.
func (f FieldRef) Role() Role { return RoleForDataType(f.DataType) }

// FieldCatalog is the set of fields a report may reference. A producer builds
// one alongside the rows it emits, so a spec can be validated against exactly
// the fields the rows carry.
type FieldCatalog struct {
	fields map[FieldKey]FieldRef
}

// NewFieldCatalog builds a catalog, rejecting keyless and duplicated fields.
func NewFieldCatalog(fields ...FieldRef) (FieldCatalog, error) {
	catalog := FieldCatalog{fields: make(map[FieldKey]FieldRef, len(fields))}

	for _, field := range fields {
		if len(field.Key) == 0 {
			return FieldCatalog{}, fmt.Errorf("%w: %q", ErrEmptyFieldKey, field.Label)
		}
		if _, exists := catalog.fields[field.Key]; exists {
			return FieldCatalog{}, fmt.Errorf("%w: %s", ErrDuplicateField, field.Key)
		}
		catalog.fields[field.Key] = field
	}

	return catalog, nil
}

// Get returns the field with the given key.
func (c FieldCatalog) Get(key FieldKey) (FieldRef, bool) {
	field, exists := c.fields[key]
	return field, exists
}
