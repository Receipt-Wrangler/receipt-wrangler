import { ReportColumn, ReportPeriod } from "../../open-api";

/**
 * A selectable engine field: `key` is the reporting engine's field key (what the
 * backend `receiptsource` catalog exposes), `label` is what the builder shows.
 * Custom fields are added at runtime keyed `custom_<id>` (see ReportCatalogService).
 */
export interface ReportField {
  key: string;
  label: string;
}

/**
 * Built-in dimensions a report can group by / cut columns on. These mirror the
 * backend's `receiptsource.builtinFields()` keys. Grouping by the raw `date`
 * buckets on the exact instant (one bucket per receipt), so the calendar-period
 * string fields (`date_month`, `date_year`) are offered for date grouping while
 * `date` stays available as a display column. The backend validates every key and
 * returns 400 on an unknown one, so this list is a UX affordance, not the gate.
 */
export const REPORT_BUILTIN_DIMENSIONS: ReportField[] = [
  { key: "paid_by", label: "Paid By" },
  { key: "category", label: "Category" },
  { key: "tag", label: "Tag" },
  { key: "status", label: "Status" },
  { key: "group", label: "Group" },
  { key: "name", label: "Name" },
  { key: "date_month", label: "Month" },
  { key: "date_year", label: "Year" },
  { key: "date", label: "Date" },
  { key: "resolved_date", label: "Resolved Date" },
  { key: "created_at", label: "Added At" },
];

/**
 * Built-in measures (currency-typed fields) an aggregate column can sum/average.
 * COUNT takes no measure. Custom currency fields are appended at runtime.
 */
export const REPORT_BUILTIN_MEASURES: ReportField[] = [
  { key: "amount", label: "Amount" },
];

/** The aggregate functions the engine supports (mirrors ReportColumn.AggFuncEnum). */
export const REPORT_AGGREGATE_FUNCTIONS: ReportColumn.AggFuncEnum[] = [
  ReportColumn.AggFuncEnum.Sum,
  ReportColumn.AggFuncEnum.Count,
  ReportColumn.AggFuncEnum.Avg,
  ReportColumn.AggFuncEnum.Min,
  ReportColumn.AggFuncEnum.Max,
];

/** COUNT is the only aggregate that needs no measure. */
export function aggregateNeedsMeasure(fn: ReportColumn.AggFuncEnum): boolean {
  return fn !== ReportColumn.AggFuncEnum.Count;
}

export interface ReportPeriodPreset {
  id: ReportPeriod.PresetEnum;
  label: string;
}

/** Period presets, in display order (mirrors ReportPeriod.PresetEnum). */
export const REPORT_PERIOD_PRESETS: ReportPeriodPreset[] = [
  { id: ReportPeriod.PresetEnum.ThisMonth, label: "This month" },
  { id: ReportPeriod.PresetEnum.LastMonth, label: "Last month" },
  { id: ReportPeriod.PresetEnum.Mtd, label: "Month to date" },
  { id: ReportPeriod.PresetEnum.Qtd, label: "Quarter to date" },
  { id: ReportPeriod.PresetEnum.Ytd, label: "Year to date" },
  { id: ReportPeriod.PresetEnum.Custom, label: "Custom range…" },
];

/**
 * The document variables the builder can insert into the title/intro/footer. The
 * backend resolves these at generation time (period window, covered group names,
 * generation timestamp, the caller's display name).
 */
export const REPORT_DOCUMENT_VARIABLES: string[] = [
  "{{period}}",
  "{{group.name}}",
  "{{generatedAt}}",
  "{{currentUser.name}}",
];

/** The output formats, in the fixed render order the backend bundles them in. */
export const REPORT_FORMATS: ReportRequestFormat[] = ["csv", "xlsx", "pdf"];

export type ReportRequestFormat = "csv" | "xlsx" | "pdf";
