/**
 * Shared presentation for group/scope avatar chips in the Report Builder (the
 * add-group dialog and the config panel's scope chips). Kept in one place so the
 * palette and initials stay in sync across both.
 */

/** Rotating avatar background palette, indexed by chip position modulo its length. */
export const CHIP_COLORS = ["#f5a3b7", "#f7b267", "#4db6ac", "#b39ddb", "#27b1ff", "#f6c453"];

/** The 1-2 character avatar initials for a group name (falls back to "?"). */
export function groupInitials(name: string | undefined): string {
  return (name || "?").trim().slice(0, 2).toUpperCase();
}
