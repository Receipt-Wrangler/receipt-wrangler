import { RECEIPT_DATE_FILTER_FIELDS } from "src/constants";
import { ReportPeriod } from "../../open-api";
import {
  formatPeriodRange,
  LEGACY_REPORT_PERIOD_DATE_FIELD,
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
