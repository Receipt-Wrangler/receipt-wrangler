import { toCalendarDay } from "./app-time-zone";

/**
 * A filter as it goes on the wire: the value of every key in `dateKeys` becomes
 * the **calendar day** the user picked, `yyyy-MM-dd` — both bounds of a
 * `BETWEEN`, or the single value otherwise. Every other key passes through.
 *
 * The datepickers write a local-midnight `Date`, which `JSON.stringify` sends as
 * an instant; the server would have to guess which zone to read its day in. A
 * bare day carries no zone to misread, and the server turns it into a range in
 * the right one: the app time zone for an instant column (`createdAt`,
 * `resolvedDate`, `startedAt`, …), UTC for a calendar column (a receipt's
 * `date`). The client never needs to know which.
 *
 * Call this where the request is assembled, **never** where the filter is
 * stored, and that placement is load-bearing: NGXS persists the filter and feeds
 * it back into the datepicker when a dialog reopens, and Material's
 * `NativeDateAdapter.deserialize` matches a bare `yyyy-MM-dd` against its
 * ISO-8601 regex and parses it with `new Date()` — i.e. as UTC midnight, which
 * renders as the *previous* day west of Greenwich. Storing the normalized form
 * would just move the off-by-one into the picker.
 *
 * Returns a new object and never mutates `filter`.
 */
export function toDateWireFilter<T extends object>(
  filter: T | undefined | null,
  dateKeys: readonly string[],
): T | undefined {
  if (!filter) {
    return undefined;
  }

  const wire = { ...filter } as Record<string, unknown>;
  for (const key of dateKeys) {
    if (key in wire) {
      wire[key] = toCalendarDayEntry(wire[key]);
    }
  }

  return wire as T;
}

function toCalendarDayEntry(entry: unknown): unknown {
  const typed = entry as { value?: unknown } | null | undefined;

  if (!typed || typeof typed !== "object") {
    return entry;
  }

  return { ...typed, value: toWireDay(typed.value) };
}

function toWireDay(value: unknown): unknown {
  // BETWEEN carries both bounds.
  if (Array.isArray(value)) {
    return value.map(toWireDay);
  }

  if (value instanceof Date || (typeof value === "string" && value.length > 0)) {
    return toCalendarDay(value);
  }

  return value;
}
