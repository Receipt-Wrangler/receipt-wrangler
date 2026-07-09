#!/usr/bin/env bash
#
# mutation-check.sh — prove the reporting engine's tests still have teeth.
#
# Coverage says which lines ran. It does not say whether anything would have
# noticed had those lines been wrong. This script breaks the engine one way at a
# time and checks that the suite objects. A mutation nothing catches is a hole
# in the tests, not in the code.
#
# It never writes to the working tree. Each mutation is applied to a copy in a
# temporary directory and compiled through `go test -overlay`, so an interrupted
# run leaves nothing behind.
#
# Three rules, each learned by getting it wrong:
#
#   1. Assert the mutation applied. A search string that matches nothing yields
#      an unmutated build that passes, and the harness reports "the tests caught
#      it" having changed nothing at all.
#
#   2. A mutant that does not compile was never tested. A build error is not a
#      failing test, and counting it as one is the same lie in a different coat.
#      A mutation is caught only when a test says `--- FAIL`.
#
#   3. Pass -count=1. The build cache will otherwise serve the previous
#      mutation's result.
#
# Usage:
#   ./internal/reporting/mutation-check.sh          # every mutation
#   ./internal/reporting/mutation-check.sh avg      # only those matching "avg"
#
# Exits non-zero if any mutation survives or fails to apply.

set -o errexit
set -o nounset
set -o pipefail

package_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
module_dir="$(cd "${package_dir}/../.." && pwd)"

filter="${1:-}"

work_dir="$(mktemp -d)"
trap 'rm -rf "${work_dir}"' EXIT

# Mutations are held in four parallel arrays rather than in delimited records,
# because most of them span several lines and bash's `read` stops at the first
# newline — which silently truncates the search text, breaks the build, and (see
# rule 2) used to be reported as a pass.
names=()
files=()
searches=()
replaces=()

# add <name> <file> <search> <replace>
#
# The comment above each names the test that must catch it. Where several do,
# the first speaks to the mutation most directly.
add() {
    names+=("$1")
    files+=("$2")
    searches+=("$3")
    replaces+=("$4")
}

# --- accumulator: the monoid ------------------------------------------------

# TestAccumulator_Avg — AVG divides by the count of values, not of records.
add 'avg-divides-by-records' 'accumulator.go' \
'return Num(a.sum.DivRound(decimal.NewFromInt(a.values), divisionScale))' \
'return Num(a.sum.DivRound(decimal.NewFromInt(a.records), divisionScale))'

# TestAccumulator_Count — COUNT counts rows, including rows with a null measure.
add 'count-skips-null-rows' 'accumulator.go' \
'	a.records++

	number, isNumber := value.Decimal()
	if !isNumber {
		return
	}' \
'	number, isNumber := value.Decimal()
	if !isNumber {
		return
	}
	a.records++'

# TestAccumulator_Sum — SUM of nothing is zero, not null.
add 'sum-of-nothing-is-null' 'accumulator.go' \
'	case AggSum:
		return Num(a.sum)' \
'	case AggSum:
		if !a.seen {
			return Null()
		}
		return Num(a.sum)'

# TestAccumulator_Merge — MIN keeps the smaller of the two.
add 'min-merge-picks-max' 'accumulator.go' \
'	if other.minimum.LessThan(a.minimum) {' \
'	if other.minimum.GreaterThan(a.minimum) {'

# TestAccumulator_MergeIsAssociativeAndCommutative
add 'merge-drops-child-sum' 'accumulator.go' \
'	a.sum = a.sum.Add(other.sum)
' \
''

# --- engine: rollup, fan-out, ordering --------------------------------------

# TestRun_AvgIsCorrectAtEveryLevel — a parent merges its children's
# accumulators, and never combines their finalized values.
add 'parent-combines-finalized-values' 'engine.go' \
'		mergeAccumulators(node.accumulators, child.accumulators)' \
'		for slot := range node.accumulators {
			node.accumulators[slot].add(child.accumulators[slot].finalize(e.scale))
		}'

# TestRun_TwoMultiValueLevelsCrossProduct — a row fans out into every bucket.
add 'drop-group-fan-out' 'engine.go' \
'		for _, path := range paths {
			for _, value := range values {' \
'		for _, path := range paths {
			for _, value := range values[:1] {'

# TestRun_MultiValueDimensionsDoubleCount
add 'drop-detail-fan-out' 'engine.go' \
'	for _, value := range row.dimensionValues(e.compiled.detailBy.Key) {' \
'	for _, value := range row.dimensionValues(e.compiled.detailBy.Key)[:1] {'

# TestRun_NoneBuckets — an absent value is a bucket, not a dropped row.
add 'no-none-bucket' 'row.go' \
'	switch len(values) {
	case 0:
		return []Value{Null()}
	case 1:
		return values
	}' \
'	switch len(values) {
	case 0:
		return nil
	case 1:
		return values
	}'

# TestRun_RepeatedDimensionValueIsOneAttribution
add 'no-dedupe-of-repeated-values' 'row.go' \
'	distinct := make([]Value, 0, len(values))
	seen := make(map[string]struct{}, len(values))' \
'	if true {
		return values
	}

	distinct := make([]Value, 0, len(values))
	seen := make(map[string]struct{}, len(values))'

# TestRun_NoneBuckets
add 'isnone-always-false' 'engine.go' \
'		group.IsNone = node.value.IsNull()' \
'		group.IsNone = false'

# TestRun_NoneSortsLastAmongManyBuckets
add 'buckets-sort-descending' 'engine.go' \
'		return compareValues(nodes[i].value, nodes[j].value) < 0' \
'		return compareValues(nodes[i].value, nodes[j].value) > 0'

# TestGolden_WorkedExampleRendersAsTheDesignDocument
add 'detail-rows-sort-descending' 'engine.go' \
'		return compareValues(buckets[i].value, buckets[j].value) < 0' \
'		return compareValues(buckets[i].value, buckets[j].value) > 0'

# TestRun_DeepNesting — the leaf sits at the last grouping level.
add 'leaf-depth-off-by-one' 'engine.go' \
'	if level == len(e.compiled.groupBy) {
		e.rollUpLeaf(node)
		return
	}' \
'	if level >= len(e.compiled.groupBy)-1 {
		e.rollUpLeaf(node)
		return
	}'

# TestRun_RecordsMode
add 'records-and-aggregate-swapped' 'engine.go' \
'func (e engineRun) emitDetailRows(leaf *buildNode, path []Value) []DetailRow {
	if e.compiled.spec.Detail.Mode == DetailRecords {' \
'func (e engineRun) emitDetailRows(leaf *buildNode, path []Value) []DetailRow {
	if e.compiled.spec.Detail.Mode != DetailRecords {'

# TestRun_ColumnOrderIsIndependentOfDependencyOrder — arithmetic is evaluated in
# dependency order, not declaration order.
add 'arithmetic-in-declaration-order' 'engine.go' \
'	for _, index := range e.compiled.arithmeticOrder {
		column := e.compiled.columns[index]
		value := evalArithmetic(column.expr, values, e.scale)' \
'	for index := range e.compiled.columns {
		column := e.compiled.columns[index]
		if column.kind != ColumnArithmetic {
			continue
		}
		value := evalArithmetic(column.expr, values, e.scale)'

# TestRun_LabelColumnReadsAncestorBucket
add 'label-reads-the-wrong-level' 'engine.go' \
'			return []Value{path[column.labelLevel]}' \
'			return []Value{path[0]}'

# TestRun_TogglesOffSubtotalsAndGrandTotals
add 'subtotals-toggle-ignored' 'engine.go' \
'		if e.compiled.spec.Subtotals {
			group.Subtotals = e.rowCells(noLabels, node.accumulators)
		}' \
'		group.Subtotals = e.rowCells(noLabels, node.accumulators)'

# TestRun_DivisionScaleIsHonoured — the caller's scale reaches AVG and division.
add 'division-scale-hardcoded' 'engine.go' \
'		scale:    compiled.spec.Config.DivisionScale,' \
'		scale:    defaultDivisionScale,'

# TestRun_MetaAndFormatsPassThrough
add 'column-format-dropped' 'engine.go' \
'			Format:   column.format,' \
'			Format:   "",'

# TestRun_MetaParamsAreCopied
add 'meta-params-aliased' 'engine.go' \
'			Params:         copyParams(meta.Params),' \
'			Params:         meta.Params,'

# TestRun_RecordsModeDoesNotAliasInputRows
add 'records-label-cell-aliases-the-row' 'engine.go' \
'			return append([]Value(nil), record.Get(column.field)...)' \
'			return record.Get(column.field)'

# TestRun_DescriptorsDoNotAlias — every descriptor would point at one aggregate.
add 'descriptors-share-an-aggregate' 'engine.go' \
'			aggregate := column.agg
			descriptor.Agg = &aggregate' \
'			descriptor.Agg = &e.compiled.columns[0].agg'

# --- value: the bucket key --------------------------------------------------

# TestBucketKey_DatesOutsideUnixNanoRange — UnixNano is undefined outside
# 1678..2262, so it merges distinct instants.
add 'date-key-uses-unixnano' 'value.go' \
'		return "d:" + strconv.FormatInt(v.date.Unix(), 10) + ":" + strconv.Itoa(v.date.Nanosecond())' \
'		return "d:" + strconv.FormatInt(v.date.UnixNano(), 10)'

# TestBucketKey_NumbersBelowTwelveDecimalPlaces
add 'number-key-truncates' 'value.go' \
'		return "n:" + v.num.String()' \
'		return "n:" + v.num.StringFixed(12)'

# TestCompareValues_NullSortsLast
add 'null-sorts-first' 'value.go' \
'		case a.IsNull():
			return 1' \
'		case a.IsNull():
			return -1'

# TestRun_DateBucketIsLocationIndependent — a bucket must report the canonical
# member of its class, not whichever member arrived first.
add 'date-bucket-keeps-an-arbitrary-zone' 'value.go' \
'	if v.valueType == ValueDate {
		return DateVal(v.date.UTC())
	}
	return v' \
'	return v'

# --- formula: the allow-list and the bounds ---------------------------------

# TestParseArithmetic_SizeLimits — the node cap does not bound parenthesis
# nesting, because a parenthesis builds no node. Only the length check does.
add 'formula-length-unbounded' 'formula.go' \
'	if len(src) > maxFormulaLength {
		return fmt.Errorf("%w: %d bytes, limit is %d", ErrFormulaTooLong, len(src), maxFormulaLength)
	}
	return nil' \
'	_ = src
	return nil'

# --- eval: decimal discipline -----------------------------------------------

# TestEvalArithmetic_DivisionByZeroIsNullNotPanic — shopspring panics on a zero
# divisor, so removing the guard crashes the report rather than emptying a cell.
add 'division-guard-removed' 'eval.go' \
'		if rightNumber.IsZero() {
			return Null()
		}
		return Num(leftNumber.DivRound(rightNumber, divisionScale))' \
'		return Num(leftNumber.DivRound(rightNumber, divisionScale))'

# TestEvalArithmetic_IgnoresDivisionPrecisionGlobal
add 'division-reads-the-global' 'eval.go' \
'		return Num(leftNumber.DivRound(rightNumber, divisionScale))' \
'		return Num(leftNumber.Div(rightNumber))'

# TestEvalArithmetic_NullPropagates
add 'null-operand-reads-as-zero' 'eval.go' \
'	leftNumber, leftIsNumber := left.Decimal()
	rightNumber, rightIsNumber := right.Decimal()
	if !leftIsNumber || !rightIsNumber {
		return Null()
	}' \
'	leftNumber, _ := left.Decimal()
	rightNumber, _ := right.Decimal()'

# --- validate: the rules ----------------------------------------------------

# TestValidate_Rejects/arithmetic references itself
#
# The discarded Join keeps "strings" used, so the mutant compiles. A mutation
# that does not build tests nothing.
add 'cycle-detection-removed' 'validate.go' \
'		case visiting:
			return fmt.Errorf("%w: %s", ErrFormulaCycle, strings.Join(append(path, columns[index].name), " -> "))' \
'		case visiting:
			_ = strings.Join(path, " -> ")
			return nil'

# TestValidate_CycleMessageNamesTheCycle
add 'cycle-path-truncated' 'validate.go' \
'strings.Join(append(path, columns[index].name), " -> ")' \
'strings.Join(path, " -> ")'

# TestValidate_Rejects/groupBy a measure
add 'groupby-accepts-a-measure' 'validate.go' \
'		if field.Role() != RoleDimension {
			return nil, fmt.Errorf("%w: %s is a %s", ErrGroupByNotDimension, key, field.DataType)
		}
' \
''

# TestValidate_Rejects/sum a multi-valued measure
add 'multi-valued-measure-accepted' 'validate.go' \
'	if field.Multi {
		return compiledColumn{}, fmt.Errorf("%w: column %s applies %s to %s",
			ErrMeasureIsMultiValued, column.Name, aggregate.Func, aggregate.Field)
	}
' \
''

# TestValidate_Rejects/aggregate with a reduction the accumulator does not know —
# an unknown reduction renders the whole column null rather than erroring.
add 'unknown-agg-func-accepted' 'validate.go' \
'	if !aggregate.Func.valid() {
		return compiledColumn{}, fmt.Errorf("%w: column %s applies reduction %d",
			ErrUnknownAggFunc, column.Name, aggregate.Func)
	}
' \
''

# TestValidate_Rejects/detail mode the engine does not know
add 'unknown-detail-mode-accepted' 'validate.go' \
'	if !detail.Mode.valid() {
		return FieldRef{}, fmt.Errorf("%w: %d", ErrUnknownDetailMode, detail.Mode)
	}
' \
''

# No mutation swaps isAggregate() back for `Mode != DetailRecords`. With the mode
# validated the two spellings agree on every spec that compiles, so the swap is
# semantically equivalent and would survive by being correct. The predicate earns
# its place by keeping a third mode from reintroducing the divergence, not by
# changing today's behaviour — and a mutation nothing can catch teaches nothing.

# TestCompileSpec_DataTypes — money combined with a count is still money.
add 'arithmetic-type-never-currency' 'validate.go' \
'		if columns[index].dataType == TypeCurrency {
			return TypeCurrency
		}' \
'		if columns[index].dataType == TypeCurrency {
			return TypeNumber
		}'

# --- receiptsource: the receipt mapping -------------------------------------

# TestSource_DuplicateCustomFieldValuesResolveByLowestId — position cannot decide
# which of two values for one field wins, because the association is loaded
# without an ORDER BY.
add 'custom-field-first-in-slice-wins' 'receiptsource/receiptsource.go' \
'		if incumbent, held := winners[key]; held && incumbent <= customFieldValue.ID {
			continue
		}' \
'		if _, held := winners[key]; held {
			continue
		}'

# TestSource_UnresolvableCustomFieldValueNeverWins — an empty value must not
# claim the field and lock out a real one.
add 'unresolvable-custom-field-value-wins' 'receiptsource/receiptsource.go' \
'		value, ok := s.customFieldValue(customFieldValue)
		if !ok {
			continue
		}

		winners[key] = customFieldValue.ID' \
'		value, _ := s.customFieldValue(customFieldValue)

		winners[key] = customFieldValue.ID'

# ----------------------------------------------------------------------------

green() { printf '\033[32m%s\033[0m' "$1"; }
red()   { printf '\033[31m%s\033[0m' "$1"; }

echo "mutation-check: ${#names[@]} mutations against internal/reporting"
echo

survivors=()
broken=()
ran=0
skipped=0

for index in "${!names[@]}"; do
    name="${names[$index]}"
    file="${files[$index]}"

    if [[ -n "${filter}" && "${name}" != *"${filter}"* ]]; then
        skipped=$((skipped + 1))
        continue
    fi

    source_file="${package_dir}/${file}"
    mutant="${work_dir}/${file}"
    overlay="${work_dir}/overlay.json"

    # Rule 1: refuse to run a mutation that changed nothing.
    if ! SEARCH="${searches[$index]}" REPLACE="${replaces[$index]}" \
        python3 - "${source_file}" "${mutant}" "${overlay}" <<'PY'
import json, os, sys

source_file, mutant, overlay = sys.argv[1], sys.argv[2], sys.argv[3]
search, replace = os.environ["SEARCH"], os.environ["REPLACE"]

original = open(source_file).read()
if search not in original:
    sys.stderr.write("search text not found\n")
    sys.exit(1)

mutated = original.replace(search, replace, 1)
if mutated == original:
    sys.stderr.write("replacement changed nothing\n")
    sys.exit(1)

os.makedirs(os.path.dirname(mutant), exist_ok=True)
open(mutant, "w").write(mutated)
open(overlay, "w").write(json.dumps({"Replace": {source_file: mutant}}))
PY
    then
        printf '  %s  %-36s the mutation no longer matches %s\n' "$(red 'STALE   ')" "${name}" "${file}"
        broken+=("${name} (stale)")
        continue
    fi

    ran=$((ran + 1))

    # Rule 3: -count=1, or the cache serves the previous mutation's verdict.
    # The ... matters: receiptsource is mutated too, and its tests live with it.
    exit_code=0
    output="$(cd "${module_dir}" && go test -count=1 -overlay="${overlay}" ./internal/reporting/... 2>&1)" || exit_code=$?

    # One awk, not `grep | head`: under `set -o pipefail`, head closing the pipe
    # hands grep a SIGPIPE and the failed pipeline takes the script down.
    caught="$(awk '/^--- FAIL: / { sub(/^--- FAIL: /, ""); sub(/ .*/, ""); if (++found <= 3) printf "%s ", $0 }' <<< "${output}")"

    if (( exit_code == 0 )); then
        printf '  %s  %s\n' "$(red 'SURVIVED')" "${name}"
        survivors+=("${name}")
    elif [[ -z "${caught}" ]]; then
        # Rule 2: it failed, but no test said so. The mutant did not build, so
        # nothing was tested and nothing was proved.
        printf '  %s  %-36s the mutant does not compile\n' "$(red 'BROKEN  ')" "${name}"
        printf '%s\n' "${output}" | head -5 | sed 's/^/            /'
        broken+=("${name} (does not compile)")
    else
        printf '  %s  %-36s %s\n' "$(green 'caught  ')" "${name}" "${caught}"
    fi
done

echo
if (( skipped > 0 )); then
    echo "${ran} run, ${skipped} skipped by filter '${filter}'"
    echo
fi

if (( ${#survivors[@]} > 0 )); then
    printf '%s: %d of %d mutations survived:\n' "$(red 'FAIL')" "${#survivors[@]}" "${ran}"
    printf '  - %s\n' "${survivors[@]}"
    echo
    echo "A surviving mutation means the engine can be broken that way without a"
    echo "single test objecting. Either add the test, or delete the mutation and"
    echo "say why the behaviour does not matter."
fi

if (( ${#broken[@]} > 0 )); then
    printf '%s: %d mutations could not be evaluated:\n' "$(red 'FAIL')" "${#broken[@]}"
    printf '  - %s\n' "${broken[@]}"
    echo
    echo "A mutation that does not apply, or does not compile, proves nothing."
    echo "Re-anchor it against the code as it stands now."
fi

if (( ${#survivors[@]} > 0 || ${#broken[@]} > 0 )); then
    exit 1
fi

printf '%s: all %d mutations caught by a failing test.\n' "$(green 'PASS')" "${ran}"
