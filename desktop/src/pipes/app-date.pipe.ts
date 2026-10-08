import { formatDate } from "@angular/common";
import { inject, LOCALE_ID, Pipe, PipeTransform } from "@angular/core";
import { Store } from "@ngxs/store";
import { SystemSettingsState } from "../store/system-settings.state";
import { calendarDayToLocalDate, offsetAt } from "../utils/app-time-zone";

/**
 * - `instant` (the default): a moment in time (`createdAt`, `resolvedDate`,
 *   `startedAt`, …), shown in the app time zone.
 * - `calendar`: a calendar day (a receipt's `date`, a DATE custom field), shown
 *   as its stored UTC day whatever the zone.
 */
export type AppDateMode = "instant" | "calendar";

export type AppDateValue = Date | string | number | null | undefined;

const DEFAULT_FORMAT = "mediumDate";

/**
 * Formats a date the way the `appDate` pipe does, for code that builds strings
 * outside a template. Returns `null` for an empty or unparseable value, so
 * `(value | appDate) || 'Never'` keeps working.
 */
export function formatAppDate(
  value: AppDateValue,
  format: string,
  mode: AppDateMode,
  zone: string,
  locale: string,
): string | null {
  if (value === null || value === undefined || value === "") {
    return null;
  }

  try {
    if (mode === "calendar") {
      // Rebuilt as the local midnight of its UTC day and formatted in the local
      // zone, rather than formatted in "UTC" directly, so a bare `yyyy-MM-dd`
      // (which Angular parses as LOCAL midnight) lands on the same day too.
      const day = calendarDayToLocalDate(value instanceof Date || typeof value === "string" ? value : new Date(value));
      return day ? formatDate(day, format, locale) : null;
    }

    return formatDate(value, format, locale, offsetAt(value, zone));
  } catch {
    return null;
  }
}

/**
 * `value | appDate[:format[:'calendar']]` — Angular's `DatePipe` with the app
 * time zone applied. `format` takes the same strings as `DatePipe` and defaults
 * to `mediumDate`.
 *
 * `DatePipe` only accepts a UTC offset, never an IANA name, and the offset
 * depends on DST at that particular instant — so it is computed per value with
 * `offsetAt`.
 *
 * **Impure on purpose.** The zone is store state an admin can change while a
 * page is open. `transform` reads it through a signal, which the template's
 * reactive consumer tracks, so a zone change marks every view that formatted a
 * date for refresh under zoneless change detection. A *pure* pipe would not
 * re-run on that refresh (its arguments did not change), and its cached string
 * would keep the old zone. The cost of re-running is kept to a comparison by
 * memoizing the last result.
 */
@Pipe({
  name: "appDate",
  standalone: true,
  pure: false,
})
export class AppDatePipe implements PipeTransform {
  private readonly locale = inject(LOCALE_ID);

  private readonly timeZone = inject(Store).selectSignal(SystemSettingsState.timeZone);

  private lastKey: unknown[] | null = null;

  private lastResult: string | null = null;

  public transform(
    value: AppDateValue,
    format: string = DEFAULT_FORMAT,
    mode: AppDateMode = "instant",
  ): string | null {
    // Read before any early return, so the dependency is tracked on every pass.
    const zone = this.timeZone();
    const valueKey = value instanceof Date ? value.getTime() : value;
    const key = [valueKey, format, mode, zone];

    if (this.lastKey && key.every((part, index) => part === this.lastKey![index])) {
      return this.lastResult;
    }

    this.lastKey = key;
    this.lastResult = formatAppDate(value, format, mode, zone, this.locale);

    return this.lastResult;
  }
}
