import { DEFAULT_DURATION_HOURS, maxForUnit, splitHours, toHours } from "./duration.utils";

describe("duration.utils", () => {
  describe("splitHours", () => {
    it("renders whole days as days", () => {
      expect(splitHours(24)).toEqual({ value: 1, unit: "DAYS" });
      expect(splitHours(168)).toEqual({ value: 7, unit: "DAYS" });
      expect(splitHours(720)).toEqual({ value: 30, unit: "DAYS" });
    });

    it("keeps values that are not whole days in hours", () => {
      expect(splitHours(1)).toEqual({ value: 1, unit: "HOURS" });
      expect(splitHours(6)).toEqual({ value: 6, unit: "HOURS" });
      expect(splitHours(36)).toEqual({ value: 36, unit: "HOURS" });
    });

    it("keeps sub-day values in hours rather than showing a fractional day", () => {
      expect(splitHours(12)).toEqual({ value: 12, unit: "HOURS" });
      expect(splitHours(23)).toEqual({ value: 23, unit: "HOURS" });
    });

    // 0 is what the API sends for "unset", and what an older client would send.
    it("falls back to the default when unset or invalid", () => {
      const fallback = { value: DEFAULT_DURATION_HOURS / 24, unit: "DAYS" };

      expect(splitHours(0)).toEqual(fallback);
      expect(splitHours(null)).toEqual(fallback);
      expect(splitHours(undefined)).toEqual(fallback);
      expect(splitHours(-5)).toEqual(fallback);
      expect(splitHours(NaN)).toEqual(fallback);
    });
  });

  describe("toHours", () => {
    it("passes hours through and multiplies days", () => {
      expect(toHours(6, "HOURS")).toBe(6);
      expect(toHours(30, "DAYS")).toBe(720);
      expect(toHours(1, "DAYS")).toBe(24);
    });

    // Number inputs hand back strings, so the conversion has to coerce.
    it("coerces string input from number fields", () => {
      expect(toHours("14", "DAYS")).toBe(336);
      expect(toHours("12", "HOURS")).toBe(12);
    });

    // Number(null) and Number("") are both 0, so an emptied number field would
    // otherwise silently submit a zero-length lifetime.
    it("falls back to the default for empty or unparseable input", () => {
      expect(toHours(null, "HOURS")).toBe(DEFAULT_DURATION_HOURS);
      expect(toHours(undefined, "HOURS")).toBe(DEFAULT_DURATION_HOURS);
      expect(toHours("", "DAYS")).toBe(DEFAULT_DURATION_HOURS);
      expect(toHours("abc", "HOURS")).toBe(DEFAULT_DURATION_HOURS);
      expect(toHours(0, "DAYS")).toBe(DEFAULT_DURATION_HOURS);
      expect(toHours(-3, "HOURS")).toBe(DEFAULT_DURATION_HOURS);
    });
  });

  describe("round trip", () => {
    it("preserves every value the backend accepts", () => {
      for (const hours of [1, 2, 6, 12, 23, 24, 36, 48, 168, 719, 720]) {
        const split = splitHours(hours);
        expect(toHours(split.value, split.unit)).toBe(hours);
      }
    });
  });

  describe("maxForUnit", () => {
    it("expresses the hour cap in the selected unit", () => {
      expect(maxForUnit(720, "HOURS")).toBe(720);
      expect(maxForUnit(720, "DAYS")).toBe(30);
    });
  });
});
