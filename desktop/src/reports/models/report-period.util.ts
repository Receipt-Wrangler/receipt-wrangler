import {
  endOfDay,
  endOfMonth,
  format,
  startOfMonth,
  startOfQuarter,
  startOfYear,
  subMonths,
} from "date-fns";
import { RECEIPT_DATE_FILTER_FIELDS, ReceiptDateFilterFieldKey } from "src/constants";
import { ReportPeriod } from "../../open-api";
import { DEFAULT_APP_TIME_ZONE, todayInZone } from "../../utils/app-time-zone";

export interface PeriodRange {
  start: Date;
  end: Date;
}

/**
 * Resolves a period preset (or a custom start/end) into a concrete date window,
 * mirroring the backend's resolvePeriodBounds so the builder's "resolves to …"
 * hint and the receipts drill-in agree with what the report will actually cover.
 *
 * "Today" is the app time zone's, as it is on the server, so the hint names the
 * same month the report covers whatever zone the browser is in. Only the
 * calendar days of the result are meaningful.
 */
export function resolvePeriodRange(
  preset: ReportPeriod.PresetEnum,
  startDate: Date | null,
  endDate: Date | null,
  timeZone: string = DEFAULT_APP_TIME_ZONE
): PeriodRange {
  const now = todayInZone(timeZone);
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
      return { start: startDate ?? startOfMonth(now), end: endDate ?? endOfDay(now) };
    default:
      return { start: startOfMonth(now), end: endOfDay(now) };
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

/**
 * The hint's name for the period date field, with the time zone its days are
 * read in when that matters: "Added At (America/New_York)". The receipt date is
 * a calendar day that no zone shifts, so it is named alone.
 */
export function reportPeriodDateFieldHint(key: ReceiptDateFilterFieldKey, timeZone: string): string {
  const label = reportPeriodDateFieldLabel(key);
  return key === "date" ? label : `${label} (${timeZone || DEFAULT_APP_TIME_ZONE})`;
}
