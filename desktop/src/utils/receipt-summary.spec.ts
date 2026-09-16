import { Group } from "../open-api";
import { resolveSummaryConfigGroup, summaryConfigGroups } from "./receipt-summary";

function group(id: number, name: string, summaryEnabled: boolean, isAllGroup = false): Group {
  return {
    id,
    name,
    isAllGroup,
    groupReceiptSettings: { receiptSummaryEnabled: summaryEnabled },
  } as Group;
}

describe("summaryConfigGroups", () => {
  it("keeps only groups that have the summary enabled, alphabetically", () => {
    const groups = [
      group(3, "Zulu", true),
      group(1, "Alpha", true),
      group(2, "Bravo", false),
    ];

    expect(summaryConfigGroups(groups).map((g) => g.name)).toEqual(["Alpha", "Zulu"]);
  });

  // The "All" group is the view that needs someone else's configuration; it is never a candidate.
  it("excludes the synthetic All group even when its settings say enabled", () => {
    const groups = [group(99, "All", true, true), group(1, "Alpha", true)];

    expect(summaryConfigGroups(groups).map((g) => g.id)).toEqual([1]);
  });

  it("tolerates a group with no receipt settings", () => {
    const groups = [{ id: 1, name: "Alpha" } as Group];

    expect(summaryConfigGroups(groups)).toEqual([]);
  });
});

describe("resolveSummaryConfigGroup", () => {
  const alpha = group(1, "Alpha", true);
  const bravo = group(2, "Bravo", true);

  it("prefers the persisted pick", () => {
    expect(resolveSummaryConfigGroup([alpha, bravo], 2)).toBe(bravo);
  });

  // The three stale cases all land on the same branch, and each is a real one:

  // State persisted before this key existed deserializes without it.
  it("falls back to the first group when nothing is persisted", () => {
    expect(resolveSummaryConfigGroup([alpha, bravo], undefined)).toBe(alpha);
  });

  // The persisted group turned its summary off, or the user left it.
  it("falls back when the persisted group is no longer a candidate", () => {
    expect(resolveSummaryConfigGroup([alpha, bravo], 999)).toBe(alpha);
  });

  it("resolves to nothing when no group has the summary enabled", () => {
    expect(resolveSummaryConfigGroup([], 1)).toBeUndefined();
  });
});
