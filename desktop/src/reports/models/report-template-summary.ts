import { ReportDetail, ReportRequestCommand } from "../../open-api";
import { REPORT_BUILTIN_DIMENSIONS, REPORT_FORMATS } from "./report-catalog.constants";

// Field-key -> label for the built-in dimensions the list renders. Custom-field
// keys (custom_<id>) are only resolvable against the loaded custom-field pool, so
// callers that have it pass a resolver.
const DIMENSION_LABELS = new Map(REPORT_BUILTIN_DIMENSIONS.map((d) => [d.key, d.label]));

/**
 * Resolves a field key the built-in map doesn't cover — a custom field's
 * `custom_<id>`. Returning undefined means "I don't know it", which falls back to
 * the raw key.
 */
export type FieldLabelResolver = (key: string) => string | undefined;

/**
 * The human label for an engine field key: a built-in's name, else whatever the
 * caller's resolver knows, else the raw key. The raw-key fallback is what a
 * caller without the custom-field pool (no resolver, or no permission to read it)
 * gets, so a row still renders something rather than nothing.
 */
export function fieldLabel(key: string | undefined, resolve?: FieldLabelResolver): string {
  if (!key) {
    return "";
  }
  return DIMENSION_LABELS.get(key) ?? resolve?.(key) ?? key;
}

/** How many columns the stored report defines. */
export function columnCount(configuration: ReportRequestCommand): number {
  return configuration.columns?.length ?? 0;
}

/**
 * The grouping levels as a readable label list, or a "no grouping" hint. A level
 * whose column heading the report renames shows that name, so the list agrees
 * with the report it describes.
 */
export function groupingSummary(
  configuration: ReportRequestCommand,
  resolve?: FieldLabelResolver
): string {
  const groupBy = configuration.groupBy ?? [];
  if (groupBy.length === 0) {
    return "No grouping";
  }
  return groupBy
    .map((key) => (configuration.groupByLabels?.[key] ?? "").trim() || fieldLabel(key, resolve))
    .join(", ");
}

/** Whether the report aggregates (and by what) or lists individual receipts. */
export function detailSummary(
  configuration: ReportRequestCommand,
  resolve?: FieldLabelResolver
): string {
  if (configuration.detail?.mode === ReportDetail.ModeEnum.Aggregate) {
    return `Aggregate by ${fieldLabel(configuration.detail.by, resolve)}`;
  }
  return "Record-level";
}

/** The output formats as uppercase chips, in the fixed csv/xlsx/pdf render order. */
export function formatChips(configuration: ReportRequestCommand): string[] {
  const selected = new Set<string>(configuration.formats ?? []);
  return REPORT_FORMATS.filter((format) => selected.has(format)).map((format) => format.toUpperCase());
}

/**
 * The report's scope groups as display names. Ids the caller can't resolve (a group
 * they don't belong to, or one not yet loaded) fall back to the raw id so the row
 * still renders something meaningful.
 */
export function scopeNames(
  configuration: ReportRequestCommand,
  groupNameById: (id: string) => string | undefined
): string[] {
  return (configuration.groupIds ?? []).map((id) => groupNameById(id) ?? id);
}
