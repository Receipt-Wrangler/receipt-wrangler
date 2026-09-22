import { Group } from "../open-api";

/**
 * The groups whose receipt-summary configuration a viewer may choose between,
 * alphabetically.
 *
 * The synthetic "All" group is excluded because it is never a real group with
 * settings of its own — it is the view that needs someone else's configuration in
 * the first place.
 */
export function summaryConfigGroups(groups: Group[]): Group[] {
  return groups
    .filter((group) => !group.isAllGroup && !!group.groupReceiptSettings?.receiptSummaryEnabled)
    .sort((a, b) => a.name.localeCompare(b.name));
}

/**
 * Resolves which group's configuration to apply: the viewer's persisted pick when
 * it still names an enabled group, otherwise the first one.
 *
 * The fallback lives here rather than in a state selector or an `@State` default
 * because it depends on the user's groups, which that slice does not know — and
 * because `defaults` never runs for a state hydrated from localStorage, so an
 * install that used the receipts table before this key existed would read
 * `undefined`. That case lands on the same branch as the other two stale ones: the
 * persisted group has since turned its summary off, or the user has left it.
 */
export function resolveSummaryConfigGroup(
  enabledGroups: Group[],
  persistedGroupId?: number
): Group | undefined {
  return (
    enabledGroups.find((group) => group.id === persistedGroupId) ?? enabledGroups[0]
  );
}
