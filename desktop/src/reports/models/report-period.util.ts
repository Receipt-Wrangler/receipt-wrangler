import {
  endOfMonth,
  endOfToday,
  format,
  startOfMonth,
  startOfQuarter,
  startOfYear,
  subMonths,
} from "date-fns";
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
