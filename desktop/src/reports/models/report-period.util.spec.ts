import { ReportPeriod } from "../../open-api";
import { formatPeriodRange, resolvePeriodRange } from "./report-period.util";

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
