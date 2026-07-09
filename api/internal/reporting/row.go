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
// the field is absent or empty.
//
// Measures are single-valued: Validate refuses to aggregate a multi-valued
// field, precisely so that no value is ever silently discarded here. The
// first-value fallback below is therefore unreachable through Run, and exists
// only so that a direct caller cannot make this panic.
func (r Row) Measure(key FieldKey) Value {
	values := r[key]
	if len(values) == 0 {
		return Null()
	}
	return values[0]
}

// dimensionValues returns the distinct buckets a row is attributed to for a
// dimension.
//
// A row with several values fans out into all of them, which double-counts it —
// that is the pie chart's attribution and is deliberate. A row with no value
// falls into the single (None) bucket, represented by a null Value, rather than
// being dropped.
//
// Values are deduplicated, so a row is attributed once per distinct bucket.
// A receipt tagged "Alex" twice belongs to the Alex bucket once; only distinct
// values fan out. The database cannot express a duplicate — the join tables key
// on the pair — but a producer over some other source could, and counting one
// row's amount twice into one bucket is never what a report means.
func (r Row) dimensionValues(key FieldKey) []Value {
	values := r[key]

	switch len(values) {
	case 0:
		return []Value{Null()}
	case 1:
		return values
	}

	distinct := make([]Value, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		bucket := bucketKey(value)
		if _, exists := seen[bucket]; exists {
			continue
		}
		seen[bucket] = struct{}{}
		distinct = append(distinct, value)
	}

	return distinct
}
