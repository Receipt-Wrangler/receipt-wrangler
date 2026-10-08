/**
 * The app time zone — one IANA name, set by an admin in System Settings and
 * carried to the desktop on AppData (`SystemSettingsState.timeZone`).
 *
 * Two kinds of date value exist, and they never mix:
 * - **Instants** (`createdAt`, `resolvedDate`, task `startedAt`/`endedAt`, …):
 *   moments in time. "Which day" is answered in the app zone.
 * - **Calendar days** (a receipt's `date`, DATE custom field values): a day with
 *   no zone, stored as **midnight UTC** of that day. "Which day" is the UTC day.
 *
 * "Today" / "this month" are always the app zone's, never the browser's.
 */

export const DEFAULT_APP_TIME_ZONE = "UTC";

const CALENDAR_DAY_PATTERN = /^(\d{4})-(\d{2})-(\d{2})$/;

// Building an Intl.DateTimeFormat is far more expensive than using one, and the
// appDate pipe asks for an offset on every change-detection pass.
const dayFormatters = new Map<string, Intl.DateTimeFormat>();
const offsetFormatters = new Map<string, Intl.DateTimeFormat>();

function cachedFormatter(
  cache: Map<string, Intl.DateTimeFormat>,
  zone: string,
  options: Intl.DateTimeFormatOptions,
): Intl.DateTimeFormat {
  const key = zone || DEFAULT_APP_TIME_ZONE;
  let formatter = cache.get(key);
  if (!formatter) {
    try {
      formatter = new Intl.DateTimeFormat("en-US", { ...options, timeZone: key });
    } catch {
      // An unknown zone name: the server falls back to UTC for the same value.
      formatter = new Intl.DateTimeFormat("en-US", { ...options, timeZone: DEFAULT_APP_TIME_ZONE });
    }
    cache.set(key, formatter);
  }

  return formatter;
}

/**
 * A local-midnight `Date` for a calendar day. `setFullYear` rather than the
 * `Date(y, m, d)` constructor, which maps years 0-99 onto 1900-1999 — and the API
 * uses `0001-01-01` as its zero date.
 */
function localMidnight(year: number, monthIndex: number, day: number): Date {
  const date = new Date(0);
  date.setFullYear(year, monthIndex, day);
  date.setHours(0, 0, 0, 0);

  return date;
}

function partValue(parts: Intl.DateTimeFormatPart[], type: Intl.DateTimeFormatPartTypes): string {
  return parts.find((part) => part.type === type)?.value ?? "";
}

/**
 * Today's calendar day in `zone`, as a **local-midnight** `Date` — the shape
 * date-fns and the month helpers work in. Only its year/month/day are meaningful.
 */
export function todayInZone(zone: string, now: Date = new Date()): Date {
  const parts = cachedFormatter(dayFormatters, zone, {
    year: "numeric",
    month: "numeric",
    day: "numeric",
  }).formatToParts(now);

  return localMidnight(
    Number(partValue(parts, "year")),
    Number(partValue(parts, "month")) - 1,
    Number(partValue(parts, "day")),
  );
}

/**
 * `zone`'s UTC offset at `instant`, as `+HHMM` / `-HHMM` — the form Angular's
 * `formatDate`/`DatePipe` accept as their timezone argument (they take offsets,
 * not IANA names). Depends on the instant, so it is DST-correct.
 */
export function offsetAt(instant: Date | string | number, zone: string): string {
  const date = instant instanceof Date ? instant : new Date(instant);
  if (Number.isNaN(date.getTime())) {
    return "+0000";
  }

  // "GMT", "GMT-04:00", "GMT+05:30"
  const name = partValue(
    cachedFormatter(offsetFormatters, zone, { timeZoneName: "longOffset" }).formatToParts(date),
    "timeZoneName",
  );
  const match = /^GMT([+-])(\d{2}):?(\d{2})?$/.exec(name);
  if (!match) {
    return "+0000";
  }

  return `${match[1]}${match[2]}${match[3] ?? "00"}`;
}

function pad(value: number, length = 2): string {
  return value.toString().padStart(length, "0");
}

function formatYmd(year: number, monthIndex: number, day: number): string {
  return `${pad(year, 4)}-${pad(monthIndex + 1)}-${pad(day)}`;
}

/**
 * The **local** calendar day a filter value names, as `yyyy-MM-dd`.
 *
 * Filter values are what the datepickers wrote: a local-midnight `Date`, or —
 * once the filter has been through the persisted NGXS state — that `Date` as an
 * ISO string. Both are read back on the browser's calendar, which is the day the
 * user picked. A bare `yyyy-MM-dd` is returned unchanged (re-parsing it would read
 * it as UTC midnight and could shift it a day). Anything unparseable is returned
 * as is.
 */
export function toCalendarDay(value: Date | string): string {
  if (value instanceof Date) {
    return formatYmd(value.getFullYear(), value.getMonth(), value.getDate());
  }

  if (CALENDAR_DAY_PATTERN.test(value)) {
    return value;
  }

  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime())
    ? value
    : formatYmd(parsed.getFullYear(), parsed.getMonth(), parsed.getDate());
}

/**
 * A **stored** calendar value (a receipt `date`, a DATE custom field) as the
 * local-midnight `Date` of its day — what a datepicker or a local-zone format
 * needs to show that day. The day is the value's **UTC** day (it is stored as
 * midnight UTC), or the day itself for a bare `yyyy-MM-dd`.
 *
 * Only pass values that came from the API: a `Date` the datepicker wrote is
 * already local midnight and must not go through this.
 */
export function calendarDayToLocalDate(value: Date | string | null | undefined): Date | null {
  if (value === null || value === undefined || value === "") {
    return null;
  }

  if (typeof value === "string") {
    const day = CALENDAR_DAY_PATTERN.exec(value);
    if (day) {
      return localMidnight(Number(day[1]), Number(day[2]) - 1, Number(day[3]));
    }
  }

  const instant = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(instant.getTime())) {
    return null;
  }

  return localMidnight(instant.getUTCFullYear(), instant.getUTCMonth(), instant.getUTCDate());
}

/**
 * A picked calendar day as it goes on the wire for a receipt `date` or a DATE
 * custom field: **midnight UTC** of that day, e.g. `2026-09-30T00:00:00.000Z`.
 *
 * A `Date` is the datepicker's local midnight, so its local day is the one the
 * user picked. A string is a value the form never touched (straight from the
 * API, or a bare day) and keeps the day it already names.
 */
export function toMidnightUtc(value: Date | string): string {
  const day = value instanceof Date ? value : calendarDayToLocalDate(value);
  if (!day || Number.isNaN(day.getTime())) {
    return value instanceof Date ? value.toISOString() : value;
  }

  const utc = new Date(0);
  utc.setUTCFullYear(day.getFullYear(), day.getMonth(), day.getDate());

  return utc.toISOString();
}

/**
 * Every zone the System Settings picker offers: `UTC` first (some browsers omit
 * it from `supportedValuesOf`), then the browser's IANA list.
 */
export function timeZoneOptions(): string[] {
  const supported =
    typeof Intl.supportedValuesOf === "function" ? Intl.supportedValuesOf("timeZone") : [];

  return [DEFAULT_APP_TIME_ZONE, ...supported.filter((zone) => zone !== DEFAULT_APP_TIME_ZONE)];
}
