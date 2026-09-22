import { FormGroup } from "@angular/forms";
import { FilterOperation, SystemTaskPagedRequestFilter } from "../open-api/index";
import { buildFieldFormGroup, listenForBetweenOperation } from "./filter-form";

/**
 * The system task filter's reactive form — the same `{ operation, value }`
 * shape as `buildReceiptFilterForm`, built from the same shared helpers, so the
 * two dialogs behave identically.
 *
 * `thisContext` must be an `@UntilDestroy()` component: the helpers tie their
 * subscriptions to its lifetime.
 */
export function buildSystemTaskFilterForm(filter: any, thisContext: any): FormGroup {
  const formGroup = new FormGroup({
    type: buildFieldFormGroup(
      filter?.type?.value ?? [],
      filter?.type?.operation,
      thisContext,
      true
    ),
    ranBy: buildFieldFormGroup(
      filter?.ranBy?.value ?? [],
      filter?.ranBy?.operation,
      thisContext,
      true
    ),
    startedAt: buildFieldFormGroup(
      filter?.startedAt?.value,
      filter?.startedAt?.operation,
      thisContext,
      filter?.startedAt?.operation === FilterOperation.Between
    ),
    endedAt: buildFieldFormGroup(
      filter?.endedAt?.value,
      filter?.endedAt?.operation,
      thisContext,
      filter?.endedAt?.operation === FilterOperation.Between
    ),
  });

  listenForBetweenOperation(formGroup, "startedAt", thisContext);
  listenForBetweenOperation(formGroup, "endedAt", thisContext);

  return formGroup;
}

const CALENDAR_DAY_PATTERN = /^\d{4}-\d{2}-\d{2}$/;

/**
 * The filter as it goes on the wire: `startedAt` / `endedAt` values become the
 * **local calendar day** the user picked, `YYYY-MM-DD`.
 *
 * `app-datepicker` writes a local-midnight `Date`, which `JSON.stringify` sends
 * as an instant. The server resolves that instant back to a day in *its* own
 * zone, so a browser far enough east or west of the API selects the adjacent
 * day — a UTC-4 browser picking Sep 22 lands on Sep 21 against an API running
 * in America/Los_Angeles. A date-only string carries the day itself and cannot
 * drift.
 *
 * This runs where the request is assembled, **not** where the filter is stored,
 * and that placement is load-bearing: NGXS persists the filter and feeds it back
 * into the datepicker when the dialog reopens, and Material's
 * `NativeDateAdapter.deserialize` matches a bare `YYYY-MM-DD` against its
 * ISO-8601 regex and parses it with `new Date()` — i.e. as UTC midnight, which
 * renders as the *previous* day west of Greenwich. Storing the normalized form
 * would just move the off-by-one into the picker.
 */
export function toSystemTaskWireFilter(
  filter: SystemTaskPagedRequestFilter | undefined | null,
): SystemTaskPagedRequestFilter | undefined {
  if (!filter) {
    return undefined;
  }

  return {
    ...filter,
    startedAt: toCalendarDayEntry(filter.startedAt),
    endedAt: toCalendarDayEntry(filter.endedAt),
  };
}

function toCalendarDayEntry(entry: unknown): any {
  const typed = entry as { value?: unknown } | null | undefined;

  if (!typed) {
    return entry;
  }

  return { ...typed, value: toCalendarDay(typed.value) };
}

function toCalendarDay(value: unknown): unknown {
  // BETWEEN carries both bounds.
  if (Array.isArray(value)) {
    return value.map(toCalendarDay);
  }

  if (value instanceof Date) {
    return formatLocalCalendarDay(value);
  }

  if (typeof value === "string" && value.length > 0) {
    // Already a calendar day — re-parsing it would read it as UTC midnight and
    // could shift it a day.
    if (CALENDAR_DAY_PATTERN.test(value)) {
      return value;
    }

    // An ISO instant, which is what a persisted Date rehydrates as.
    const parsed = new Date(value);
    return Number.isNaN(parsed.getTime()) ? value : formatLocalCalendarDay(parsed);
  }

  return value;
}

function formatLocalCalendarDay(value: Date): string {
  const year = value.getFullYear().toString().padStart(4, "0");
  const month = (value.getMonth() + 1).toString().padStart(2, "0");
  const day = value.getDate().toString().padStart(2, "0");

  return `${year}-${month}-${day}`;
}
