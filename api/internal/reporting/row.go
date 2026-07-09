package reporting

// Row is one source record with its fields already resolved to typed values.
//
// Every field holds a slice so that scalar and multi-value fields share one
// shape: a scalar field carries zero or one value, a multi-value field such as
// categories or tags carries zero or more. Producers build Rows — see
// internal/reporting/receiptsource — and the engine never resolves anything
// itself.
//
// A Row may be nil or may omit a key entirely; both read as an absent value.
type Row map[FieldKey][]Value

// Get returns the values resolved for a field, or nil when the field is absent.
func (r Row) Get(key FieldKey) []Value {
	return r[key]
}

// Measure returns the single value a measure field contributes, or null when
// the field is absent or empty. Measures are single-valued by definition, so a
// field that somehow resolved to several values contributes only its first.
func (r Row) Measure(key FieldKey) Value {
	values := r[key]
	if len(values) == 0 {
		return Null()
	}
	return values[0]
}

// dimensionValues returns the buckets a row is attributed to for a dimension.
//
// A row with several values fans out into all of them, which double-counts it —
// that is the pie chart's attribution and is deliberate. A row with no value
// falls into the single (None) bucket, represented by a null Value, rather than
// being dropped.
func (r Row) dimensionValues(key FieldKey) []Value {
	values := r[key]
	if len(values) == 0 {
		return []Value{Null()}
	}
	return values
}
