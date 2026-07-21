package reporting

import "sort"

// Run turns a report definition and a set of rows into a ReportModel.
//
// It fetches nothing. Rows arrive already resolved — see
// internal/reporting/receiptsource — and Run neither reads a clock nor consults
// any global, so the same inputs always produce the same output.
//
// The spec is validated first; a spec that does not compile is an error before
// any row is touched.
func Run(spec ReportSpec, catalog FieldCatalog, rows []Row, meta MetaInput) (ReportModel, error) {
	compiled, err := compileSpec(spec, catalog)
	if err != nil {
		return ReportModel{}, err
	}

	run := newEngineRun(compiled)

	root := newBuildNode(Null())
	for _, row := range rows {
		run.insert(root, row)
	}
	run.rollUp(root, 0)

	model := ReportModel{
		Meta: Meta{
			Title:       compiled.spec.Title,
			GeneratedAt: meta.GeneratedAt,
			Params:      copyParams(meta.Params),
			Currency:    meta.Currency,
			NoneLabel:   compiled.spec.NoneLabel,
		},
		Columns: run.descriptors(),
		Root:    run.emitGroup(root, 0, nil),
	}

	if compiled.spec.GrandTotals {
		model.GrandTotals = run.rowCells(noLabels, root.accumulators)
	}

	return model, nil
}

// engineRun carries the compiled spec and the derived lookups the walk needs.
type engineRun struct {
	compiled compiledSpec
	scale    int32

	// aggSlot maps a column's index to its position in a node's accumulator
	// slice, and is -1 for columns that are not aggregates.
	aggSlot []int
}

func newEngineRun(compiled compiledSpec) engineRun {
	aggSlot := make([]int, len(compiled.columns))
	for index := range aggSlot {
		aggSlot[index] = -1
	}
	for slot, index := range compiled.aggregateIndexes {
		aggSlot[index] = slot
	}

	return engineRun{
		compiled: compiled,
		scale:    compiled.spec.Config.DivisionScale,
		aggSlot:  aggSlot,
	}
}

// buildNode is a node of the tree while it is being built. It becomes a
// GroupNode on the way out.
type buildNode struct {
	value Value

	children   []*buildNode
	childByKey map[string]*buildNode

	// records holds the leaf's rows in records mode, in input order.
	records []Row

	// buckets holds the leaf's aggregated detail rows in aggregate mode.
	buckets     []*detailBucket
	bucketByKey map[string]*detailBucket

	// accumulators and recordCount are filled by rollUp.
	accumulators []accumulator
	recordCount  int
}

// detailBucket is one aggregated detail row: everything attributed to one value
// of the detail dimension.
type detailBucket struct {
	value        Value
	accumulators []accumulator
	recordCount  int
}

func newBuildNode(value Value) *buildNode {
	return &buildNode{value: value, childByKey: map[string]*buildNode{}}
}

// insert attributes one row to every bucket it belongs in.
func (e engineRun) insert(root *buildNode, row Row) {
	for _, path := range e.groupPaths(row) {
		node := root
		for _, value := range path {
			node = node.child(value)
		}
		e.insertLeaf(node, row)
	}
}

// groupPaths returns every path through the grouping levels a row is attributed
// to.
//
// A row with two tags is attributed to both tag buckets in full, which
// double-counts it — the same attribution the dashboard pie chart uses. Two
// multi-value levels therefore produce a cross product. A level the row has no
// value for contributes the single (None) bucket, so the row is never dropped.
//
// With no grouping levels there is exactly one path, the empty one, and the
// root is itself the leaf.
func (e engineRun) groupPaths(row Row) [][]Value {
	paths := [][]Value{nil}

	for _, field := range e.compiled.groupBy {
		values := row.dimensionValues(field.Key)
		extended := make([][]Value, 0, len(paths)*len(values))

		for _, path := range paths {
			for _, value := range values {
				next := make([]Value, len(path), len(path)+1)
				copy(next, path)
				extended = append(extended, append(next, value))
			}
		}
		paths = extended
	}

	return paths
}

func (n *buildNode) child(value Value) *buildNode {
	key := bucketKey(value)
	if existing, found := n.childByKey[key]; found {
		return existing
	}

	// The canonical member of the class the key names, so the bucket's reported
	// value cannot depend on which member arrived first.
	created := newBuildNode(value.canonical())
	n.childByKey[key] = created
	n.children = append(n.children, created)

	return created
}

func (e engineRun) insertLeaf(leaf *buildNode, row Row) {
	if !e.compiled.spec.Detail.isAggregate() {
		leaf.records = append(leaf.records, row)
		return
	}

	for _, value := range row.dimensionValues(e.compiled.detailBy.Key) {
		bucket := e.bucket(leaf, value)
		bucket.recordCount++
		e.addRow(bucket.accumulators, row)
	}
}

func (e engineRun) bucket(leaf *buildNode, value Value) *detailBucket {
	key := bucketKey(value)
	if leaf.bucketByKey == nil {
		leaf.bucketByKey = map[string]*detailBucket{}
	}
	if existing, found := leaf.bucketByKey[key]; found {
		return existing
	}

	created := &detailBucket{value: value.canonical(), accumulators: e.newAccumulators()}
	leaf.bucketByKey[key] = created
	leaf.buckets = append(leaf.buckets, created)

	return created
}

// rollUp fills every node's accumulators, bottom up.
//
// A parent merges its children's accumulators; it never re-reads their rows and
// never combines their finalized values. That is what makes AVG at a subtotal
// the average over all of its descendants rather than the average of their
// averages, and it is why the fan-out's double count propagates to the grand
// total exactly as the pie chart's does.
func (e engineRun) rollUp(node *buildNode, level int) {
	node.accumulators = e.newAccumulators()

	if level == len(e.compiled.groupBy) {
		e.rollUpLeaf(node)
		return
	}

	for _, child := range node.children {
		e.rollUp(child, level+1)
		mergeAccumulators(node.accumulators, child.accumulators)
		node.recordCount += child.recordCount
	}
}

func (e engineRun) rollUpLeaf(leaf *buildNode) {
	if !e.compiled.spec.Detail.isAggregate() {
		for _, record := range leaf.records {
			e.addRow(leaf.accumulators, record)
		}
		leaf.recordCount = len(leaf.records)
		return
	}

	for _, bucket := range leaf.buckets {
		mergeAccumulators(leaf.accumulators, bucket.accumulators)
		leaf.recordCount += bucket.recordCount
	}
}

func (e engineRun) newAccumulators() []accumulator {
	accumulators := make([]accumulator, len(e.compiled.aggregateIndexes))
	for slot, index := range e.compiled.aggregateIndexes {
		accumulators[slot] = newAccumulator(e.compiled.columns[index].agg.Func)
	}
	return accumulators
}

// addRow folds one row into a set of accumulators. COUNT's field is empty, which
// reads as null, so it contributes to the record count and nothing else.
func (e engineRun) addRow(accumulators []accumulator, row Row) {
	for slot, index := range e.compiled.aggregateIndexes {
		accumulators[slot].add(row.Measure(e.compiled.columns[index].agg.Field))
	}
}

func mergeAccumulators(into []accumulator, from []accumulator) {
	for slot := range into {
		into[slot].merge(from[slot])
	}
}

// emitGroup converts a built node into the model's GroupNode.
//
// path carries the bucket value of each ancestor level, which is what lets a
// label column on an aggregated detail row show the group it sits under.
func (e engineRun) emitGroup(node *buildNode, level int, path []Value) GroupNode {
	group := GroupNode{RecordCount: node.recordCount}

	if level > 0 {
		group.Dimension = e.compiled.groupBy[level-1].Key
		group.Value = node.value
		group.IsNone = node.value.IsNull()

		if e.compiled.spec.Subtotals {
			group.Subtotals = e.rowCells(noLabels, node.accumulators)
		}
	}

	if level == len(e.compiled.groupBy) {
		group.DetailRows = e.emitDetailRows(node, path)
		return group
	}

	sortNodes(node.children)
	group.Children = make([]GroupNode, 0, len(node.children))
	for _, child := range node.children {
		// A fresh slice per child: siblings must not share a backing array,
		// since each carries its own bucket value down to the detail rows.
		childPath := make([]Value, len(path), len(path)+1)
		copy(childPath, path)
		childPath = append(childPath, child.value)

		group.Children = append(group.Children, e.emitGroup(child, level+1, childPath))
	}

	return group
}

func (e engineRun) emitDetailRows(leaf *buildNode, path []Value) []DetailRow {
	if !e.compiled.spec.Detail.isAggregate() {
		return e.emitRecordRows(leaf)
	}
	return e.emitBucketRows(leaf, path)
}

// emitRecordRows emits one row per record, in the order the rows arrived. The
// engine is deterministic given deterministic input; ordering records is the
// caller's query's job, not the engine's.
func (e engineRun) emitRecordRows(leaf *buildNode) []DetailRow {
	rows := make([]DetailRow, 0, len(leaf.records))

	for _, record := range leaf.records {
		accumulators := e.newAccumulators()
		e.addRow(accumulators, record)

		labels := func(column compiledColumn) []Value {
			// Copy, so a renderer holding a cell cannot reach back into the
			// caller's row.
			return append([]Value(nil), record.Get(column.field)...)
		}
		rows = append(rows, DetailRow{Cells: e.rowCells(labels, accumulators)})
	}

	return rows
}

func (e engineRun) emitBucketRows(leaf *buildNode, path []Value) []DetailRow {
	sortBuckets(leaf.buckets)
	rows := make([]DetailRow, 0, len(leaf.buckets))

	for _, bucket := range leaf.buckets {
		labels := func(column compiledColumn) []Value {
			// An aggregated row sums many records, so the only values it can
			// show are its own bucket and those of the groups above it.
			if column.labelLevel == labelFromDetailBucket {
				return []Value{bucket.value}
			}
			return []Value{path[column.labelLevel]}
		}

		rows = append(rows, DetailRow{
			Dimension: e.compiled.detailBy.Key,
			Value:     bucket.value,
			IsNone:    bucket.value.IsNull(),
			Cells:     e.rowCells(labels, bucket.accumulators),
		})
	}

	return rows
}

// rowCells assembles one rendered row, whatever level it sits at.
//
// Label and aggregate cells are filled first, then arithmetic columns are
// evaluated in dependency order against the values on this very row. That
// recomputation is what keeps a non-linear formula correct at a subtotal: an
// average per receipt is this row's total over this row's count, never the sum
// or the average of the averages below it.
func (e engineRun) rowCells(labels func(compiledColumn) []Value, accumulators []accumulator) []Cell {
	cells := make([]Cell, len(e.compiled.columns))
	values := make(map[string]Value, len(e.compiled.columns))

	for index, column := range e.compiled.columns {
		cells[index].Column = column.name

		switch column.kind {
		case ColumnLabel:
			cells[index].Values = labels(column)
			values[column.name] = cells[index].Value()

		case ColumnAggregate:
			value := accumulators[e.aggSlot[index]].finalize(e.scale)
			cells[index].Values = []Value{value}
			values[column.name] = value
		}
	}

	for _, index := range e.compiled.arithmeticOrder {
		column := e.compiled.columns[index]
		value := evalArithmetic(column.expr, values, e.scale)
		cells[index].Values = []Value{value}
		values[column.name] = value
	}

	return cells
}

// noLabels leaves label cells empty, which is what a subtotal or grand-total row
// shows for them.
func noLabels(compiledColumn) []Value { return nil }

func (e engineRun) descriptors() []ColumnDescriptor {
	descriptors := make([]ColumnDescriptor, 0, len(e.compiled.columns))

	for _, column := range e.compiled.columns {
		descriptor := ColumnDescriptor{
			Name:     column.name,
			Label:    column.label,
			Kind:     column.kind,
			DataType: column.dataType,
			Format:   column.format,
		}

		switch column.kind {
		case ColumnAggregate:
			aggregate := column.agg
			descriptor.Agg = &aggregate
		case ColumnArithmetic:
			descriptor.Expr = column.expr
			descriptor.ExprSrc = column.exprSrc
		}

		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

// sortNodes and sortBuckets order buckets by their value, with (None) last.
// Nothing in this package emits output by ranging a map.
func sortNodes(nodes []*buildNode) {
	sort.SliceStable(nodes, func(i, j int) bool {
		return compareValues(nodes[i].value, nodes[j].value) < 0
	})
}

func sortBuckets(buckets []*detailBucket) {
	sort.SliceStable(buckets, func(i, j int) bool {
		return compareValues(buckets[i].value, buckets[j].value) < 0
	})
}

func copyParams(params map[string]string) map[string]string {
	if params == nil {
		return nil
	}

	copied := make(map[string]string, len(params))
	for key, value := range params {
		copied[key] = value
	}

	return copied
}
