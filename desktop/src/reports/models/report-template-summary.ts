import { ReportDetail, ReportRequestCommand } from "../../open-api";
import { REPORT_BUILTIN_DIMENSIONS, REPORT_FORMATS } from "./report-catalog.constants";

// Field-key -> label for the built-in dimensions the list renders. Custom-field
// grouping keys (custom_<id>) aren't in this map and fall back to the raw key.
const DIMENSION_LABELS = new Map(REPORT_BUILTIN_DIMENSIONS.map((d) => [d.key, d.label]));

/** The human label for an engine field key, or the raw key when it isn't a built-in. */
export function fieldLabel(key: string | undefined): string {
  if (!key) {
    return "";
  }
  return DIMENSION_LABELS.get(key) ?? key;
}

/** How many columns the stored report defines. */
export function columnCount(configuration: ReportRequestCommand): number {
  return configuration.columns?.length ?? 0;
}

/** The grouping levels as a readable label list, or a "no grouping" hint. */
export function groupingSummary(configuration: ReportRequestCommand): string {
  const groupBy = configuration.groupBy ?? [];
  if (groupBy.length === 0) {
    return "No grouping";
  }
  return groupBy.map(fieldLabel).join(", ");
}

/** Whether the report aggregates (and by what) or lists individual receipts. */
export function detailSummary(configuration: ReportRequestCommand): string {
  if (configuration.detail?.mode === ReportDetail.ModeEnum.Aggregate) {
    return `Aggregate by ${fieldLabel(configuration.detail.by)}`;
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
