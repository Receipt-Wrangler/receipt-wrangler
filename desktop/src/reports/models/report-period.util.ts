import {
  endOfMonth,
  endOfToday,
  format,
  startOfMonth,
  startOfQuarter,
  startOfYear,
  subMonths,
} from "date-fns";
import { RECEIPT_DATE_FILTER_FIELDS, ReceiptDateFilterFieldKey } from "src/constants";
import { ReportPeriod } from "../../open-api";

export interface PeriodRange {
  start: Date;
  end: Date;
}

/**
 * Resolves a period preset (or a custom start/end) into a concrete date window,
 * mirroring the backend's resolvePeriodBounds so the builder's "resolves to …"
 * hint and the receipts drill-in agree with what the report will actually cover.
 */
export function resolvePeriodRange(
  preset: ReportPeriod.PresetEnum,
  startDate: Date | null,
  endDate: Date | null
): PeriodRange {
  const now = new Date();
  switch (preset) {
    case ReportPeriod.PresetEnum.ThisMonth:
      return { start: startOfMonth(now), end: endOfMonth(now) };
    case ReportPeriod.PresetEnum.LastMonth: {
      const lastMonth = subMonths(now, 1);
      return { start: startOfMonth(lastMonth), end: endOfMonth(lastMonth) };
    }
    case ReportPeriod.PresetEnum.Mtd:
      return { start: startOfMonth(now), end: now };
    case ReportPeriod.PresetEnum.Qtd:
      return { start: startOfQuarter(now), end: now };
    case ReportPeriod.PresetEnum.Ytd:
      return { start: startOfYear(now), end: now };
    case ReportPeriod.PresetEnum.Custom:
      return { start: startDate ?? startOfMonth(now), end: endDate ?? endOfToday() };
    default:
      return { start: startOfMonth(now), end: endOfToday() };
  }
}

export function formatPeriodRange(range: PeriodRange): string {
  return `${format(range.start, "yyyy-MM-dd")} to ${format(range.end, "yyyy-MM-dd")}`;
}

/**
 * The date a report period covers when its command names none: every template
 * saved before the picker existed. It mirrors the API's own default
 * (ReportPeriod.DateFilterKey), and is a literal rather than
 * DEFAULT_QUICK_DATE_FIELD on purpose: those templates always ran on the receipt
 * date, so a later change to the receipts table's default must not change what
 * they cover.
 */
export const LEGACY_REPORT_PERIOD_DATE_FIELD: ReceiptDateFilterFieldKey = "date";

/**
 * Narrows a stored period date field to one the picker offers. The API contract
 * types it as a plain string (see ReportPeriod.dateField in swagger.yml), so a
 * missing or unrecognized value falls back to the legacy receipt date.
 */
export function toReportPeriodDateField(value?: string | null): ReceiptDateFilterFieldKey {
  return (
    RECEIPT_DATE_FILTER_FIELDS.find((field) => field.key === value)?.key ??
    LEGACY_REPORT_PERIOD_DATE_FIELD
  );
}

/** The picker's label for a period date field, e.g. "Added At". */
export function reportPeriodDateFieldLabel(key: ReceiptDateFilterFieldKey): string {
  return RECEIPT_DATE_FILTER_FIELDS.find((field) => field.key === key)?.label ?? "";
}
