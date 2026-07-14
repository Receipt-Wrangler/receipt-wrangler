# The reporting engine

A **pure** function from a report definition and some rows to a format-agnostic tree:

```
(ReportSpec + FieldCatalog + []Row + MetaInput)  ──▶  ReportModel
```

It fetches nothing, renders nothing, reads no clock, and consults no global. The same inputs always
produce byte-identical output. Renderers (CSV/XLSX/PDF), a dashboard widget, template persistence and
HTTP delivery all *call* it; none of them are part of it.

The purity claim is mechanically checkable, not aspirational:

```bash
go list -deps receipt-wrangler/api/internal/reporting | grep -E 'gorm|repositories|internal/models'
# prints nothing
```

> **Every number in this document is copied from a test assertion.** Where a figure appears, the test
> that owns it is named. If you change the engine and a number here looks wrong, that test is already
> red and will tell you so.

---

## 1. The pipeline

`Run` (`engine.go`) is five stages. Rows arrive already resolved — the engine never looks anything up.

```
   []Row                                                            ReportModel
     │                                                                   ▲
     │  ┌───────────────┐   the spec is validated and everything it      │
     └─▶│  compileSpec  │   needs resolved once: field refs, parsed      │
        │  validate.go  │   formulas, column types, arithmetic order     │
        └───────┬───────┘                                                │
                │                                                        │
        ┌───────▼───────┐   each row is attributed to every bucket it    │
        │    insert     │   belongs in. groupPaths() walks the grouping  │
        │  groupPaths   │   levels; a multi-value dimension fans the row │
        └───────┬───────┘   out into a cross product of paths            │
                │                                                        │
        ┌───────▼───────┐   bottom-up. A parent MERGES its children's    │
        │    rollUp     │   accumulators — never their finalized values. │
        │  rollUpLeaf   │   This is the whole point of the accumulator.  │
        └───────┬───────┘                                                │
                │                                                        │
        ┌───────▼───────┐   buckets sort by typed value, (None) last.    │
        │   emitGroup   │   Arithmetic columns are recomputed per row,   │
        │   rowCells    │   at every level. Nothing ranges a Go map.     │
        └───────────────┘───────────────────────────────────────────────▶┘
```

Validation happens **before any row is touched**: a spec that does not compile is an error, not a
report full of empty cells.

---

## 2. The vocabulary

| Type | What it is |
|---|---|
| `Value` | An immutable typed scalar: null, string, number (`decimal`), date, bool. The zero `Value` is null. |
| `Row` | `map[FieldKey][]Value` — one source record, fields already resolved. Every field holds a *slice*, so scalar and multi-value fields share one shape. |
| `FieldRef` | A field a report may reference: key, label, data type, and `Multi`. |
| `FieldCatalog` | The set of fields a spec may name. A producer builds one alongside the rows it emits. |
| `ReportSpec` | What to group by, what the bottom rows are, and what the columns compute. |
| `ReportModel` | The output tree: `Root` → `Children` → `DetailRows` → `Cells`. |

### A field's *role* is derived, never declared

```
DataType                Role         May be used as
────────────────────────────────────────────────────────────────
TypeNumber, TypeCurrency  →  Measure    the input to SUM/AVG/MIN/MAX
TypeString, TypeDate,     →  Dimension  a groupBy level, or the key
TypeBool                                of an aggregated detail row
```

`RoleForDataType` (`field.go`) is the whole rule. A template can therefore never group by a dollar
amount or sum a status — the spec simply will not compile.

### Three column kinds, declared not inferred

| Kind | Computes | Rolls up by |
|---|---|---|
| `ColumnLabel` | displays a field | it doesn't — it is **blank** on a subtotal row |
| `ColumnAggregate` | `SUM(amount)`, `COUNT()` | **merging accumulators** |
| `ColumnArithmetic` | `Subtotal + Hst` | **recomputing from the same row** |

The kind is declared rather than inferred from whether summing happens to work. §5 is why.

---

## 3. The output tree

```
ReportModel
├── Meta              Title, GeneratedAt, Params, CurrencyFormat, NoneLabel
├── Columns  []ColumnDescriptor   name, label, kind, data type, format
├── Root     GroupNode            synthetic: no Dimension, no Subtotals
│   ├── Dimension  FieldKey       the field this node's siblings were cut by
│   ├── Value      Value          this bucket's value ((None) ⇒ null + IsNone)
│   ├── RecordCount int
│   ├── Children   []GroupNode    ─┐ a node holds one or the other,
│   ├── DetailRows []DetailRow    ─┘ never both: leaves hold detail
│   └── Subtotals  []Cell         nil unless the spec asked for them
└── GrandTotals []Cell            nil unless the spec asked for them
```

A `Cell` holds `[]Value`. Aggregate and arithmetic cells always hold exactly one value, possibly null.
A label cell in records mode may hold several, because a record can carry several categories or tags.

**Renderers consume this tree.** `internal/reporting/render` is the first: `render.CSV(model, groupBy)`
emits a flat, data-only CSV — the group-by dimensions become leading columns, each detail leaf is one
row, and a leading `Row Type` column marks each row `Detail` / `Subtotal` / `Grand Total`. It renders
only what the model carries: subtotal and grand-total rows appear when the spec asked for them, and never
otherwise. The `Row Type` column is what keeps the file safe to aggregate — filter `Row Type=Detail`
before summing an additive column, or the roll-up rows double-count it. It draws no document chrome;
grouped/visual layouts are the XLSX/PDF renderers' job, each a separate consumer of the same tree.

`render.XLSX(model, groupBy)` is the faithful, spreadsheet-native one: the group-by dimensions are leading
columns (each value shown once per group), subtotal/grand-total rows carry a `Total` marker in the column
at the group's depth, and numbers are written as **native cells** with a number format (not strings), with
header and total rows bold. It writes the engine-computed values **statically** — translating arithmetic
columns into live cell formulas (the reason `ColumnDescriptor.Expr` is exported) is a later slice.

`render.HTML(model, groupBy, chrome)` is the PDF format's HTML stage: a self-contained document — a
`Meta.Title` heading, an optional authored intro, a preamble of the report's resolved `Meta.Params`, the
same faithful table as XLSX, and a footer (each omitted when there is nothing for it). The third argument,
`render.DocumentChrome{Intro, Footer}`, is authored presentation copy layered on at render time — kept
**out** of the pure model, which stays presentation-free; a zero value leaves the document unchanged. An
authored footer replaces the automatic `Meta.GeneratedAt` note, and any `{{variable}}` substitution is
the caller's job (done before rendering). All CSS is inline and it references no external resources or
scripts, so it converts to PDF through the headless-Chromium pipeline in `services/html_to_pdf.go` (which
blocks network loads and disables JavaScript). It returns HTML bytes; the chromedp conversion to PDF is
the caller's job. The two faithful renderers (XLSX and HTML) share one traversal — `faithfulWalk` in
`render/walk.go`, which drives a format-specific `faithfulSink` — while CSV keeps its own flat walk.

---

## 4. A worked example, end to end

Ten rows. Group by payer, then by tag; one summed row per category.

### The input

| paid_by | tag  | category | amount | custom_1 (Hst) |
|---------|------|----------|--------|----------------|
| Dana    | Alex | Clothing | 30.00  | 3.90           |
| Dana    | Alex | Clothing | 30.00  | 3.90           |
| Dana    | Alex | Clothing | 30.00  | 3.90           |
| Dana    | Alex | Clothing | 30.00  | 3.90           |
| Dana    | Alex | Medical  | 50.00  | *(none)*       |
| Dana    | Alex | Medical  | 30.00  | *(none)*       |
| Dana    | Sam  | Clothing | 30.00  | 3.90           |
| Dana    | Sam  | Clothing | 30.00  | 3.90           |
| Dana    | Sam  | Clothing | 30.00  | 3.90           |
| Dana    | Sam  | Mileage  | 30.00  | *(none)*       |

### The spec

```go
sum := func(f FieldKey) Aggregate { return Aggregate{Func: AggSum, Field: f} }

ReportSpec{
    GroupBy: []FieldKey{"paid_by", "tag"},
    Detail:  DetailSpec{Mode: DetailAggregate, By: "category"},
    Columns: []Column{
        {Name: "Category", Kind: ColumnLabel, Field: "category"},
        {Name: "Count",    Kind: ColumnAggregate, Agg: Aggregate{Func: AggCount}},
        {Name: "Subtotal", Kind: ColumnAggregate, Agg: sum("amount")},
        {Name: "Hst",      Kind: ColumnAggregate, Agg: sum("custom_1")},
        {Name: "Total",    Kind: ColumnArithmetic, Expr: "Subtotal + Hst"},
        {Name: "AvgPerReceipt", Label: "Avg/Receipt",
         Kind: ColumnArithmetic, Expr: "Total / Count"},
    },
    Subtotals:   true,
    GrandTotals: true,
}
```

`Column.Name` is what formulas reference; `Column.Label` is the heading a reader sees. They are
separate so that renaming a heading cannot break a formula.

### The tree it becomes

```
                                    Count  Subtotal    Hst   Total  Avg/Receipt
Root                        n=10        ·         ·      ·       ·            ·
└── paid_by "Dana"          n=10       10    320.00  27.30  347.30        34.73
    ├── tag "Alex"          n=6         6    200.00  15.60  215.60        35.93
    │   ├── category Clothing  ·        4    120.00  15.60  135.60        33.90
    │   └── category Medical   ·        2     80.00   0.00   80.00        40.00
    └── tag "Sam"           n=4         4    120.00  11.70  131.70        32.93
        ├── category Clothing  ·        3     90.00  11.70  101.70        33.90
        └── category Mileage   ·        1     30.00   0.00   30.00        30.00
GrandTotals                            10    320.00  27.30  347.30        34.73
```

Nodes with an `n=` are `GroupNode`s and the numbers beside them are their `Subtotals`. The indented
`category` rows are `DetailRow`s. `Root` is synthetic: it has no `Dimension` and never carries
`Subtotals` — the report-wide figures live in `GrandTotals`.

### The report a renderer draws

Owned by `TestGolden_WorkedExampleRendersAsTheDesignDocument` (`golden_test.go`).

```
Category      Count  Subtotal    Hst   Total  Avg/Receipt
Paid By: Dana
  Tag: Alex
    Clothing      4    120.00  15.60  135.60        33.90
    Medical       2     80.00   0.00   80.00        40.00
    TOTALS        6    200.00  15.60  215.60        35.93
  Tag: Sam
    Clothing      3     90.00  11.70  101.70        33.90
    Mileage       1     30.00   0.00   30.00        30.00
    TOTALS        4    120.00  11.70  131.70        32.93
  TOTALS         10    320.00  27.30  347.30        34.73
GRAND TOTALS     10    320.00  27.30  347.30        34.73
```

Two things to notice.

**`Hst` for Medical is `0.00`, not blank.** `SUM` and `COUNT` of nothing are zero — a category with no
tax shows `0.00`. `AVG`, `MIN` and `MAX` of nothing are *null*, and render as an empty cell, because
there is no value to report and zero would be a lie.

**`Category` is blank on every `TOTALS` row.** A label cuts rather than measures, so it does not roll
up. The engine emits a cell with **zero values** there. That is a different thing from a `(None)`
bucket, whose label cell holds **one null value** — see §7.

---

## 5. The two rollups, and why they differ

This is the part that is easy to get wrong, and the reason `ColumnKind` is declared rather than
guessed.

### Aggregate columns merge accumulators — never finalized values

An `accumulator` (`accumulator.go`) is a monoid. Rows fold into it with `add`; a parent folds its
children into itself with `merge`; only `finalize` divides. It carries the running `sum`, the count of
non-null `values`, the count of `records`, and the running `minimum`/`maximum`.

For `SUM` it makes no difference — summing a child's `SUM` gives the parent's `SUM`. For `AVG` it is
everything. Here is `AVG(amount)` over four receipts of 30 and one of 80 under `Alex`, and one of 120
under `Sam`. Owned by `TestRun_AvgIsCorrectAtEveryLevel` (`engine_test.go`):

```
                                   Avg(amount)
Root                          n=6           ·
└── paid_by "Dana"            n=6   53.333333   ← not 80, the mean of 40 and 120
    ├── tag "Alex"            n=5          40   ← not 55, the mean of 30 and 80
    │   ├── category Clothing              30   (120 / 4)
    │   └── category Medical               80   ( 80 / 1)
    └── tag "Sam"             n=1         120
        └── category Medical              120   (120 / 1)
GrandTotals                         53.333333
```

`Alex` is `200 / 5 = 40`, not `(30 + 80) / 2 = 55`. `Dana` is `320 / 6 = 53.333333`, not
`(40 + 120) / 2 = 80`. Because the parent merged `sum` and `values` and divided once at the end,
`AVG` at *every* level is `sum(all descendants) / count(all descendants)` — with no special case.

```
     leaf "Clothing"          leaf "Medical"
     sum=120  values=4        sum=80   values=1
              └──────┬────────────┘
                  merge          (NOT finalize-then-combine)
                     ▼
              sum=200  values=5
                     │ finalize
                     ▼
                    40
```

### Arithmetic columns are recomputed from the row being rendered

At a detail row, a subtotal row and the grand total alike, an arithmetic column is evaluated against
the *other cells on that same row*. It is never summed and never averaged.

Look at `Avg/Receipt` under `Alex` in the worked example. Three plausible answers:

```
  sum the column          33.90 + 40.00           = 73.90     ✗ nonsense
  average the column     (33.90 + 40.00) / 2      = 36.95     ✗ wrong
  recompute the row       Total / Count
                          215.60 / 6              = 35.93     ✓ the answer
```

The same rule covers the additive case for free. `Total = Subtotal + Hst` recomputes to
`200.00 + 15.60 = 215.60`, which is also what you get by summing `135.60 + 80.00`. An additive formula
agrees with summing; a ratio does not. **One rule, both correct.**

A zero divisor yields an empty cell, never a panic — `shopspring` panics on division by zero, so a
report with a zero count would otherwise take down the request that asked for it.

### What each kind shows, where

| | detail row | subtotal row | grand total |
|---|---|---|---|
| `ColumnLabel` | the field's value(s) | *blank* (zero values) | *blank* |
| `ColumnAggregate` | the bucket's `finalize` | merged children's `finalize` | root's `finalize` |
| `ColumnArithmetic` | recomputed from this row | recomputed from this row | recomputed from this row |

A consequence worth knowing: a **numeric label column reads null on a subtotal row**, so arithmetic
over it is null there too. That is correct — there is no single record for it to name.

---

## 6. Multi-value fan-out (it double-counts on purpose)

A receipt with two categories is attributed to **both**, in full, and that double count propagates to
the grand total. This is deliberate: it is the same attribution the dashboard pie chart uses.

One receipt of `100.00` in two categories. Owned by `TestRun_MultiValueDimensionsDoubleCount`:

```
                              Count   Subtotal
Root                    n=2
├── category Clothing             1     100.00
└── category Medical              1     100.00
GrandTotals                       2     200.00     ← one receipt, a grand total of 200
```

Two multi-value levels produce a **cross product**. One receipt of `10.00` with two tags and two
categories lands in `2 × 2 = 4` buckets and totals `40.00`
(`TestRun_TwoMultiValueLevelsCrossProduct`).

```
row: tag=[Alex, Sam]  category=[Clothing, Medical]  amount=10.00

  groupPaths ──▶  [Alex, Clothing]   [Alex, Medical]
                  [Sam,  Clothing]   [Sam,  Medical]
```

**Only *distinct* values fan out.** A receipt tagged `"Alex"` twice belongs to the `Alex` bucket once,
and its amount is counted once (`Row.dimensionValues` deduplicates by bucket key). Fan-out
double-counts *across* different buckets, never *within* one.

**A multi-valued field cannot be measured.** `SUM` over a field that resolves to several values would
silently read the first and drop the rest, so `Validate` refuses it with `ErrMeasureIsMultiValued`. A
`Multi` field remains a perfectly good dimension and a perfectly good display label.

---

## 7. `(None)` — an absent value is a bucket, not a dropped row

A row with no value for a dimension falls into an explicit `(None)` bucket, marked `IsNone: true` with
a null `Value`. It sorts last. Owned by `TestGolden_NoneBucketsRenderLast`:

```
Category      Count  Subtotal    Hst   Total  Avg/Receipt
Paid By: Dana
  Tag: Alex
    Clothing      4    120.00  15.60  135.60        33.90
    Medical       2     80.00   0.00   80.00        40.00
    (None)        1      7.00   0.00    7.00         7.00
    TOTALS        7    207.00  15.60  222.60        31.80
  Tag: Sam
    Clothing      3     90.00  11.70  101.70        33.90
    Mileage       1     30.00   0.00   30.00        30.00
    TOTALS        4    120.00  11.70  131.70        32.93
  Tag: (None)
    (None)        1      3.00   0.00    3.00         3.00
    TOTALS        1      3.00   0.00    3.00         3.00
  TOTALS         12    330.00  27.30  357.30        29.78
GRAND TOTALS     12    330.00  27.30  357.30        29.78
```

Note the difference between a blank `TOTALS` label cell and a `(None)` label cell:

```
label cell with ZERO values  →  a subtotal row. No bucket named it.  renders ""
label cell with ONE null     →  the (None) bucket. A bucket like     renders "(None)"
                                any other, given the report's name.
```

An empty string is **not** `(None)`. `Str("")` is a bucket of its own and sorts first.

> `IsNone == Value.IsNull()` holds for group nodes and aggregate-mode detail rows **only**. The
> synthetic `Root` and every records-mode `DetailRow` carry a null `Value` with `IsNone: false`. It is
> not a global biconditional — do not "fix" the engine to make it one.

---

## 8. Determinism

"The same inputs always produce the same output" is a real guarantee, and five rules hold it up.

**1. The bucket-key law.** Two values share a bucket key *exactly when* `compareValues` finds them
equal. A key coarser than that merges distinct values; a key finer splits one value across two
buckets. Numbers key on `decimal.String()` (canonical and lossless); dates key on `Unix()` seconds plus
`Nanosecond()` — and never `UnixNano()`, which is undefined outside 1678–2262, where a zero `time.Time`
and a date in 585 share one.

**2. A bucket stores the canonical member of its class.** Agreeing with `compareValues` is only half
the job. A bucket keeps whichever value created it and drops the rest, so what it reports must depend
on the *class*, not on which member arrived first. Dates compare by instant, so `2026-05-01T00:00:00Z`
and `2026-04-30T19:00:00-05:00` merge — and a renderer formatting a calendar day would print a
different day for each. `Value.canonical()` normalizes a date to UTC, so **date buckets are emitted in
UTC**.

**3. Buckets are sorted, and no output is ever produced by ranging a Go map.** Siblings sort by typed
value with `(None)` last. Records preserve input order — ordering them is the caller's query's job.

**4. `GeneratedAt` is an input** (`MetaInput`), never `time.Now()`. A golden test could not exist
otherwise.

**5. Money is `decimal`, never `float`.** Division goes through `DivRound(x, Config.DivisionScale)`.
The engine never reads or writes `decimal.DivisionPrecision`, a mutable process-wide global — relying
on it would make a report's output depend on whatever else the process had done.

---

## 9. Formulas

`expr-lang/expr` is used as a **parser front-end only** (`parser.Parse` → `ast.Node`). Its VM is never
run, because it evaluates over `float64` and money must not touch a float. The tree is walked by our
own decimal evaluator (`eval.go`).

The whitelist is a **closed type switch whose `default` rejects**, so it holds by construction rather
than by remembering to deny each new construct:

```
allowed:   number literals · column references · + - * / · unary sign · ROUND(x, n)
rejected:  everything else — property access, comparisons, conditionals, arrays,
           string and boolean literals, builtins, variable declarations, ...
```

Two traps, both tested:

- `??`, `in`, `**` and `%` all parse as `*ast.BinaryNode`, so the **operator** is whitelisted too, not
  just the node type.
- expr ships lower-case builtins named `sum`, `min`, `max`, `count`, `round`, `len`, so `round(a, 2)`
  arrives as an `*ast.BuiltinNode` rather than a `CallNode`. Mis-casing produces a helpful error.

Arithmetic columns form a **DAG**, topologically sorted with cycle detection, so a column may be
declared before the aggregate it reads:

```
  Avg  ──reads──▶  Total  ──reads──▶  Subtotal  (aggregate)
                     │                 Hst       (aggregate)
                     └──reads──────────┘

  Declaration order is irrelevant. A cycle is ErrFormulaCycle, and the
  error names its members:  A -> B -> C -> A
```

Source longer than `maxFormulaLength` (1 KB) is refused **before parsing**. expr's ten-thousand-node
cap does not bound this: a parenthesis builds no node, so nesting is bounded only by the goroutine
stack, and a Go stack overflow is a fatal error that `recover` cannot catch.

`ColumnDescriptor.Expr` keeps the parsed tree, so a spreadsheet renderer can later translate a column
into a live `=SUM(...)` cell formula instead of a computed value.

---

## 10. Feeding it data

`internal/reporting/receiptsource` is the **only** package that imports `models`. It maps
`[]models.Receipt` + `[]models.CustomField` into `[]reporting.Row`.

```go
source, err := receiptsource.New(customFields)   // needs CustomField.Options loaded
model,  err := reporting.Run(spec, source.Catalog(), source.Rows(receipts), meta)
```

It offers `receipt_id`, `name`, `amount`, `date`, `resolved_date`, `created_at`, `status`, `paid_by`,
`group`, `category`, `tag` (the last two `Multi`), plus one field per custom field, keyed
`custom_<id>` — **by id, not by name, so renaming a custom field cannot break a saved report**.

Each date field also carries derived `_day` / `_month` / `_year` **string** fields — `date_month`,
`created_at_year`, `resolved_date_day`, and so on. A report groups by one of these to bucket receipts by
calendar period; the raw date fields carry the exact instant, so grouping by `date` puts every receipt
in its own bucket. The strings are zero-padded ISO in UTC (`2026-05`), so they sort chronologically as
plain text.

The caller must preload `PaidByUser`, `Group`, `Categories`, `Tags`, `CustomFields`. An unloaded
association resolves to no value, which surfaces as a `(None)` bucket rather than as a crash.

Reporting at item grain, or a widget over something else entirely, adds a **sibling** of this package.
The engine core never changes.

---

## 11. Extending it

**Adding a reduction** means updating **four** independent switches over `AggFunc`, and Go forces none
of them to agree:

```
  aggregate.go   String()            → "MEDIAN"
  aggregate.go   valid()             → admit it, or Validate rejects the column
  aggregate.go   aggFuncFromName()   → parse "MEDIAN(amount)"
  accumulator.go finalize()          → compute it
```

Wire up three and forget `finalize`, and `Validate` accepts the column while **every cell of it
renders blank at every level, with no error anywhere**. `enums_test.go` pins the four lists to each
other exhaustively over all 256 values, so that mistake fails at `go test` time. `DetailMode` is
pinned the same way.

**Errors are sentinels**, wrapped with `%w` and checked with `errors.Is`. `Validate` is total: whatever
it is handed, it answers, and it never panics.

---

## 12. How it is tested

Both packages are pure — **no `main_test.go`, no DB, no `app.db` cleanup.** Four layers, each catching
what the one before it cannot:

| Layer | File | What it does |
|---|---|---|
| Invariants | `invariants_test.go` | Re-derives the model's structure from the spec. `mustRun` calls it, so **every** engine test checks it in passing. |
| Golden | `golden_test.go` | Renders a whole report as a text table. One assertion covers grouping, ordering, values, subtotal placement and blank cells. |
| Properties | `property_test.go` | 400 randomized specs and rows from a fixed seed (`REPORTING_SEED` overrides). Conservation, rollup, AVG exactness, arithmetic recompute, determinism. |
| Fuzz | `fuzz_test.go`, `bucketkey_test.go` | Seed corpora run under plain `go test`; `-fuzz` is opt-in. |

**Compare money with `decimal.Equal`/`StringFixed`, never `reflect.DeepEqual`** — `NewFromInt(200)` and
`200.00` are equal but carry different internal exponents.

**A law must not be derived from the code it judges.** This rule has been broken twice here, and both
times the suite went green over a real defect: a property test that computed bucket identity with the
engine's own helpers, and an invariant serializer that rendered dates through `bucketKey`. Both now
read the model the way a renderer would.

### Mutation checking

```bash
./internal/reporting/mutation-check.sh          # all 43 mutations
./internal/reporting/mutation-check.sh avg      # only those matching "avg"
```

Breaks the engine one way at a time and asserts a test objects. Uses `go test -overlay`, so it never
writes to the working tree. **Run it before merging any change to the engine** — a survivor means the
engine can be broken that way without a single test noticing.

Its three self-imposed rules, each learned the hard way: a mutation whose search text no longer matches
is a **failure**, not a pass; a mutant that **does not compile** was never tested, so only `--- FAIL`
counts as caught; and `-count=1`, or the build cache serves the last verdict.

---

## 13. Known gaps

**Per-row allocation.** `rowCells` builds a `map[string]Value` for every rendered row. At 50k rows a
report costs roughly 165 ms and 39 allocations per row, dominated by `decimal`'s `big.Int` arithmetic
rather than the map. Fine for a background job; worth revisiting if a synchronous caller ever appears.
