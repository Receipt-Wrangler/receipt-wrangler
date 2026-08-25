import { parseISO } from "date-fns";
import { FormArray, FormBuilder, FormControl, FormGroup } from "@angular/forms";
import { ReportColumn, ReportDetail, ReportPeriod, ReportRequestCommand } from "../../open-api";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { ReportColumnValue } from "./report-command.mapper";

let columnIdCounter = 0;

/** A monotonic client-only id for @for tracking of scope/group-by/column rows. */
export function nextReportRowId(): string {
  return "r" + ++columnIdCounter;
}

/**
 * Builds the Report Builder's root form. The filter sub-group reuses the shared
 * buildReceiptFilterForm so the builder's filter stays in lockstep with the
 * receipts filter (same shape, same BETWEEN/validator wiring); it needs the owning
 * component as `thisContext` for its untilDestroyed subscriptions, so this is
 * called from the container (which is @UntilDestroy()-decorated).
 */
export function buildReportForm(formBuilder: FormBuilder, thisContext: any): FormGroup {
  return formBuilder.group({
    name: formBuilder.control("Untitled Report"),
    scope: formBuilder.array<FormControl<string>>([]),
    period: formBuilder.group({
      preset: formBuilder.control<ReportPeriod.PresetEnum>(ReportPeriod.PresetEnum.ThisMonth),
      startDate: formBuilder.control<Date | null>(null),
      endDate: formBuilder.control<Date | null>(null),
    }),
    filter: buildReceiptFilterForm({}, thisContext),
    groupBy: formBuilder.array<FormGroup>([]),
    detail: formBuilder.group({
      mode: formBuilder.control<ReportDetail.ModeEnum>(ReportDetail.ModeEnum.Aggregate),
      by: formBuilder.control("category"),
    }),
    columns: formBuilder.array<FormGroup>(
      DEFAULT_COLUMNS.map((column) => buildColumnGroup(formBuilder, column))
    ),
    subtotals: formBuilder.control(true),
    grandTotals: formBuilder.control(true),
    document: formBuilder.group({
      title: formBuilder.control(""),
      intro: formBuilder.control("Period Covering: {{period}}"),
      footer: formBuilder.control("Generated {{generatedAt}}"),
    }),
    formats: formBuilder.group({
      csv: formBuilder.control(false),
      xlsx: formBuilder.control(false),
      pdf: formBuilder.control(true),
    }),
  });
}

/**
 * Builds the Report Builder's form seeded from a saved template's stored command,
 * so opening a template rehydrates the builder exactly as it was saved. It mirrors
 * buildReportForm's shape but sources every value from the command rather than the
 * defaults — the inverse of toReportRequestCommand, so
 * `toReportRequestCommand(buildReportFormFromCommand(cmd).getRawValue())` round-trips
 * back to `cmd` for any command the forward map could have produced. The filter reuses
 * buildReceiptFilterForm (which seeds its own array/BETWEEN wiring from an initial
 * filter), and the FormArrays (scope, group-by, columns) are rebuilt rather than
 * patched. Document slots seed to "" when the command omits the document (the forward
 * map drops it only when all three are empty), so a re-map omits it again instead of
 * re-emitting the builder's non-empty defaults.
 */
export function buildReportFormFromCommand(
  formBuilder: FormBuilder,
  thisContext: any,
  command: ReportRequestCommand
): FormGroup {
  const formats = command.formats ?? [];
  return formBuilder.group({
    name: formBuilder.control(command.name ?? ""),
    scope: formBuilder.array<FormControl<string>>(
      (command.groupIds ?? []).map((id) => formBuilder.control(id, { nonNullable: true }))
    ),
    period: formBuilder.group({
      preset: formBuilder.control<ReportPeriod.PresetEnum>(
        command.period?.preset ?? ReportPeriod.PresetEnum.ThisMonth
      ),
      startDate: formBuilder.control<Date | null>(
        command.period?.startDate ? parseISO(command.period.startDate) : null
      ),
      endDate: formBuilder.control<Date | null>(
        command.period?.endDate ? parseISO(command.period.endDate) : null
      ),
    }),
    filter: buildReceiptFilterForm(command.filter ?? {}, thisContext),
    groupBy: formBuilder.array<FormGroup>(
      (command.groupBy ?? []).map((key) =>
        buildGroupByGroup(formBuilder, key, command.groupByLabels?.[key] ?? "")
      )
    ),
    detail: formBuilder.group({
      mode: formBuilder.control<ReportDetail.ModeEnum>(
        command.detail?.mode ?? ReportDetail.ModeEnum.Aggregate
      ),
      by: formBuilder.control(command.detail?.by ?? "category"),
    }),
    columns: formBuilder.array<FormGroup>(
      // The generated ReportColumn types name/label as optional; a stored column
      // always carries them, so coerce to the builder column shape's required strings.
      (command.columns ?? []).map((column) =>
        buildColumnGroup(formBuilder, {
          kind: column.kind,
          name: column.name ?? "",
          label: column.label ?? "",
          field: column.field,
          aggFunc: column.aggFunc,
          measure: column.measure,
          expr: column.expr,
        })
      )
    ),
    subtotals: formBuilder.control(command.subtotals ?? true),
    grandTotals: formBuilder.control(command.grandTotals ?? true),
    document: formBuilder.group({
      title: formBuilder.control(command.document?.title ?? ""),
      intro: formBuilder.control(command.document?.intro ?? ""),
      footer: formBuilder.control(command.document?.footer ?? ""),
    }),
    formats: formBuilder.group({
      csv: formBuilder.control(formats.includes(ReportRequestCommand.FormatsEnum.Csv)),
      xlsx: formBuilder.control(formats.includes(ReportRequestCommand.FormatsEnum.Xlsx)),
      pdf: formBuilder.control(formats.includes(ReportRequestCommand.FormatsEnum.Pdf)),
    }),
  });
}

/** A sensible starting report so the builder and preview show something on load. */
const DEFAULT_COLUMNS: Omit<ReportColumnValue, "id">[] = [
  { kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
  { kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
  { kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
];

/**
 * Builds a grouping-level FormGroup. A level is `{ field, label }` rather than a
 * bare field key because each level also renders as a leading column in the
 * report, whose heading the user may rename. A blank label means "use the field
 * catalog's own label", which is what the mapper omits from the request.
 */
export function buildGroupByGroup(
  formBuilder: FormBuilder,
  field: string,
  label: string = ""
): FormGroup {
  return formBuilder.group({
    field: formBuilder.control(field),
    label: formBuilder.control(label),
  });
}

/** Builds a column FormGroup from a column value (used by defaults and the picker). */
export function buildColumnGroup(formBuilder: FormBuilder, column: Omit<ReportColumnValue, "id">): FormGroup {
  return formBuilder.group({
    id: formBuilder.control(nextReportRowId()),
    kind: formBuilder.control(column.kind),
    name: formBuilder.control(column.name),
    label: formBuilder.control(column.label),
    field: formBuilder.control(column.field ?? ""),
    aggFunc: formBuilder.control(column.aggFunc ?? ReportColumn.AggFuncEnum.Sum),
    measure: formBuilder.control(column.measure ?? ""),
    expr: formBuilder.control(column.expr ?? ""),
  });
}

/** Reads the scope FormArray as a plain string[] of group ids. */
export function readStringArray(array: FormArray): string[] {
  return array.controls.map((control) => control.value as string);
}

/** Reads the group-by FormArray as the ordered engine field keys it groups on. */
export function readGroupByFields(array: FormArray): string[] {
  return array.controls.map((control) => control.get("field")!.value as string);
}
