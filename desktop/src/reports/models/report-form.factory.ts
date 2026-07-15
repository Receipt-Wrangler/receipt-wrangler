import { FormArray, FormBuilder, FormControl, FormGroup } from "@angular/forms";
import { ReportColumn, ReportDetail, ReportPeriod } from "../../open-api";
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
    groupBy: formBuilder.array<FormControl<string>>([]),
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

/** A sensible starting report so the builder and preview show something on load. */
const DEFAULT_COLUMNS: Omit<ReportColumnValue, "id">[] = [
  { kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
  { kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
  { kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
];

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

/** Reads the scope/group-by FormArray as a plain string[] of ids/keys. */
export function readStringArray(array: FormArray): string[] {
  return array.controls.map((control) => control.value as string);
}
