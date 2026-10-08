import {
  calendarDayToLocalDate,
  DEFAULT_APP_TIME_ZONE,
  offsetAt,
  timeZoneOptions,
  toCalendarDay,
  todayInZone,
  toMidnightUtc,
} from "./app-time-zone";

const ymd = (date: Date): [number, number, number] => [date.getFullYear(), date.getMonth(), date.getDate()];

describe("app time zone utilities", () => {
  describe("todayInZone", () => {
    // 9PM Sep 30 in New York is already Oct 1 in UTC: the reported bug.
    const lateSep30Eastern = new Date("2026-10-01T01:00:00Z");

    it("returns the zone's calendar day as a local-midnight Date", () => {
      const today = todayInZone("America/New_York", lateSep30Eastern);

      expect(ymd(today)).toEqual([2026, 8, 30]);
      expect([today.getHours(), today.getMinutes()]).toEqual([0, 0]);
    });

    it("differs from UTC across the date line", () => {
      expect(ymd(todayInZone("UTC", lateSep30Eastern))).toEqual([2026, 9, 1]);
      expect(ymd(todayInZone("Asia/Tokyo", new Date("2026-09-30T16:00:00Z")))).toEqual([2026, 9, 1]);
    });

    it("falls back to UTC for an unknown zone, as the server does", () => {
      expect(ymd(todayInZone("Not/AZone", lateSep30Eastern))).toEqual([2026, 9, 1]);
    });

    it("defaults to now", () => {
      expect(todayInZone("UTC")).toBeInstanceOf(Date);
    });
  });

  describe("offsetAt", () => {
    it("is +0000 for UTC", () => {
      expect(offsetAt("2026-07-01T12:00:00Z", "UTC")).toBe("+0000");
    });

    // DST starts at 2AM local on Mar 8 2026 (07:00Z): the offset belongs to
    // the instant, not to the zone.
    it("follows DST in America/New_York", () => {
      expect(offsetAt(new Date("2026-03-08T06:59:00Z"), "America/New_York")).toBe("-0500");
      expect(offsetAt(new Date("2026-03-08T07:00:00Z"), "America/New_York")).toBe("-0400");
      expect(offsetAt("2026-11-01T05:59:00Z", "America/New_York")).toBe("-0400");
      expect(offsetAt("2026-11-01T06:00:00Z", "America/New_York")).toBe("-0500");
    });

    it("keeps half-hour offsets", () => {
      expect(offsetAt("2026-07-01T12:00:00Z", "Asia/Kolkata")).toBe("+0530");
    });

    it("is +0000 for an unparseable instant or an unknown zone", () => {
      expect(offsetAt("not-a-date", "America/New_York")).toBe("+0000");
      expect(offsetAt("2026-07-01T12:00:00Z", "Not/AZone")).toBe("+0000");
    });
  });

  describe("toCalendarDay", () => {
    it("reads a Date on the local calendar", () => {
      expect(toCalendarDay(new Date(2026, 0, 5))).toBe("2026-01-05");
      expect(toCalendarDay(new Date(2026, 8, 30, 23, 59, 59))).toBe("2026-09-30");
    });

    it("reads an ISO string as the local day of the Date it was serialized from", () => {
      expect(toCalendarDay(new Date(2026, 8, 30).toISOString())).toBe("2026-09-30");
    });

    it("keeps a bare day, and an unparseable string, as is", () => {
      expect(toCalendarDay("2026-09-30")).toBe("2026-09-30");
      expect(toCalendarDay("nope")).toBe("nope");
    });
  });

  describe("calendarDayToLocalDate", () => {
    it("shows a stored midnight-UTC value as its UTC day", () => {
      expect(calendarDayToLocalDate("2026-09-30T00:00:00Z")).toEqual(new Date(2026, 8, 30));
      expect(calendarDayToLocalDate(new Date("2026-09-30T00:00:00Z"))).toEqual(new Date(2026, 8, 30));
    });

    // A legacy value saved as a US browser's local midnight is still that day in UTC.
    it("reads a legacy local-midnight value by its UTC day", () => {
      expect(calendarDayToLocalDate("2026-09-30T04:00:00Z")).toEqual(new Date(2026, 8, 30));
    });

    it("takes a bare day as written", () => {
      expect(calendarDayToLocalDate("2026-09-30")).toEqual(new Date(2026, 8, 30));
    });

    it("keeps years below 100 rather than mapping them into the 1900s", () => {
      expect(calendarDayToLocalDate("0001-01-01T00:00:00Z")!.getFullYear()).toBe(1);
    });

    it("is null for nothing or garbage", () => {
      expect(calendarDayToLocalDate(null)).toBeNull();
      expect(calendarDayToLocalDate(undefined)).toBeNull();
      expect(calendarDayToLocalDate("")).toBeNull();
      expect(calendarDayToLocalDate("nope")).toBeNull();
    });
  });

  describe("toMidnightUtc", () => {
    it("sends a picked local-midnight Date as midnight UTC of that day", () => {
      expect(toMidnightUtc(new Date(2026, 8, 30))).toBe("2026-09-30T00:00:00.000Z");
    });

    it("keeps the day an untouched stored value or bare day names", () => {
      expect(toMidnightUtc("2026-09-30T00:00:00Z")).toBe("2026-09-30T00:00:00.000Z");
      expect(toMidnightUtc("2026-09-30")).toBe("2026-09-30T00:00:00.000Z");
    });

    it("round-trips with calendarDayToLocalDate", () => {
      const stored = "2026-02-28T00:00:00.000Z";
      expect(toMidnightUtc(calendarDayToLocalDate(stored)!)).toBe(stored);
    });

    it("returns an unparseable string unchanged", () => {
      expect(toMidnightUtc("nope")).toBe("nope");
    });
  });

  describe("timeZoneOptions", () => {
    it("starts with a single UTC, followed by the browser's zones", () => {
      const options = timeZoneOptions();

      expect(options[0]).toBe(DEFAULT_APP_TIME_ZONE);
      expect(options.filter((zone) => zone === "UTC").length).toBe(1);
      expect(options).toContain("America/New_York");
    });
  });
});
