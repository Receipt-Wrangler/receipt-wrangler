/**
 * Helpers for editing an hours-based duration setting with a Hours/Days
 * selector.
 *
 * Hours are the single source of truth: the API stores and transports hours, and
 * the unit exists only so an admin can type "30 Days" instead of "720". Keep the
 * conversion here rather than in the form component so the two settings that use
 * it cannot drift.
 */
export type DurationUnit = "HOURS" | "DAYS";

export const HOURS_PER_DAY = 24;

/** Fallback shown when a duration setting is unset (0/null on the wire). */
export const DEFAULT_DURATION_HOURS = 24;

export interface SplitDuration {
  value: number;
  unit: DurationUnit;
}

/**
 * Picks the friendliest unit for an hours value: Days when it divides evenly
 * into whole days, Hours otherwise. An unset/invalid value falls back to the
 * default, matching the backend's own clamp.
 */
export function splitHours(hours?: number | null): SplitDuration {
  const total = Number(hours);
  const safeTotal = Number.isFinite(total) && total > 0 ? total : DEFAULT_DURATION_HOURS;

  if (safeTotal >= HOURS_PER_DAY && safeTotal % HOURS_PER_DAY === 0) {
    return { value: safeTotal / HOURS_PER_DAY, unit: "DAYS" };
  }

  return { value: safeTotal, unit: "HOURS" };
}

/**
 * Converts an edited value back to the hours the API expects. Anything that is
 * not a positive number falls back to the default — note `Number(null)` and
 * `Number("")` are both 0, so an emptied field must be caught by the positivity
 * check, not by isFinite alone.
 */
export function toHours(value: number | string | null | undefined, unit: DurationUnit): number {
  const parsed = Number(value);

  if (!Number.isFinite(parsed) || parsed <= 0) {
    return DEFAULT_DURATION_HOURS;
  }

  return unit === "DAYS" ? parsed * HOURS_PER_DAY : parsed;
}

/**
 * The upper bound expressed in the given unit. The backend caps the lifetime at
 * 720 hours, so the form's max has to track whichever unit is selected.
 */
export function maxForUnit(maxHours: number, unit: DurationUnit): number {
  return unit === "DAYS" ? Math.floor(maxHours / HOURS_PER_DAY) : maxHours;
}
