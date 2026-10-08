import { RECEIPT_DATE_FILTER_FIELDS } from "src/constants";
import { ReportPeriod } from "../../open-api";
import {
  formatPeriodRange,
  LEGACY_REPORT_PERIOD_DATE_FIELD,
  reportPeriodDateFieldHint,
  reportPeriodDateFieldLabel,
  resolvePeriodRange,
  toReportPeriodDateField,
} from "./report-period.util";

describe("resolvePeriodRange", () => {
  it("returns the supplied bounds for a custom range", () => {
    const start = new Date(2026, 1, 10);
    const end = new Date(2026, 2, 20);
    expect(resolvePeriodRange(ReportPeriod.PresetEnum.Custom, start, end)).toEqual({ start, end });
  });

  it("starts this month on the first of the current month", () => {
    const { start } = resolvePeriodRange(ReportPeriod.PresetEnum.ThisMonth, null, null);
    expect(start.getDate()).toBe(1);
    expect(start.getMonth()).toBe(new Date().getMonth());
  });

  it("starts YTD on January 1 of the current year", () => {
    const { start } = resolvePeriodRange(ReportPeriod.PresetEnum.Ytd, null, null);
    expect(start.getMonth()).toBe(0);
    expect(start.getDate()).toBe(1);
    expect(start.getFullYear()).toBe(new Date().getFullYear());
  });

  it("falls back to a sane window for an unknown preset without dates", () => {
    const { start, end } = resolvePeriodRange(
      "someday" as ReportPeriod.PresetEnum,
      null,
      null
    );
    expect(start.getTime()).toBeLessThanOrEqual(end.getTime());
  });
});

// 9PM Sep 30 in New York is Oct 1 in UTC: the reported "Last Month" bug. The
// hint resolves "today" in the app zone, as the server does.
describe("resolvePeriodRange in the app time zone", () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date("2026-10-01T01:00:00Z"));
  });
  afterEach(() => jest.useRealTimers());

  const resolve = (preset: ReportPeriod.PresetEnum, zone: string) =>
    formatPeriodRange(resolvePeriodRange(preset, null, null, zone));

  it("names the zone's month", () => {
    expect(resolve(ReportPeriod.PresetEnum.ThisMonth, "America/New_York")).toBe("2026-09-01 to 2026-09-30");
    expect(resolve(ReportPeriod.PresetEnum.ThisMonth, "UTC")).toBe("2026-10-01 to 2026-10-31");
    expect(resolve(ReportPeriod.PresetEnum.LastMonth, "America/New_York")).toBe("2026-08-01 to 2026-08-31");
    expect(resolve(ReportPeriod.PresetEnum.LastMonth, "UTC")).toBe("2026-09-01 to 2026-09-30");
  });

  it("ends the to-date presets on the zone's today", () => {
    expect(resolve(ReportPeriod.PresetEnum.Mtd, "America/New_York")).toBe("2026-09-01 to 2026-09-30");
    expect(resolve(ReportPeriod.PresetEnum.Ytd, "UTC")).toBe("2026-01-01 to 2026-10-01");
  });

  it("defaults to UTC", () => {
    expect(formatPeriodRange(resolvePeriodRange(ReportPeriod.PresetEnum.ThisMonth, null, null))).toBe(
      "2026-10-01 to 2026-10-31"
    );
  });
});

describe("formatPeriodRange", () => {
  it("formats both bounds as YYYY-MM-DD joined by 'to'", () => {
    const range = { start: new Date(2026, 4, 1), end: new Date(2026, 4, 31) };
    expect(formatPeriodRange(range)).toBe("2026-05-01 to 2026-05-31");
  });
});

describe("toReportPeriodDateField", () => {
  it.each(RECEIPT_DATE_FILTER_FIELDS.map((field) => field.key))("keeps %s", (key) => {
    expect(toReportPeriodDateField(key)).toBe(key);
  });

  it("maps a missing field to the receipt date", () => {
    expect(LEGACY_REPORT_PERIOD_DATE_FIELD).toBe("date");
    expect(toReportPeriodDateField(undefined)).toBe("date");
    expect(toReportPeriodDateField(null)).toBe("date");
    expect(toReportPeriodDateField("")).toBe("date");
  });

  it("maps a field the picker does not offer to the receipt date", () => {
    expect(toReportPeriodDateField("created_at")).toBe("date");
    expect(toReportPeriodDateField("amount")).toBe("date");
  });
});

describe("reportPeriodDateFieldLabel", () => {
  it("reads each label from the receipts filter's own field table", () => {
    expect(RECEIPT_DATE_FILTER_FIELDS.map((field) => reportPeriodDateFieldLabel(field.key))).toEqual([
      "Receipt Date",
      "Resolved Date",
      "Added At",
    ]);
  });
});

describe("reportPeriodDateFieldHint", () => {
  it("names the zone an instant field is read in", () => {
    expect(reportPeriodDateFieldHint("createdAt", "America/New_York")).toBe("Added At (America/New_York)");
    expect(reportPeriodDateFieldHint("resolvedDate", "UTC")).toBe("Resolved Date (UTC)");
  });

  it("names the receipt date alone, since no zone shifts a calendar day", () => {
    expect(reportPeriodDateFieldHint("date", "America/New_York")).toBe("Receipt Date");
  });
});
