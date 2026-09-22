import { FilterOperation, SystemTaskPagedRequestFilter } from "../open-api";
import { toSystemTaskWireFilter } from "./system-task-filter";

describe("toSystemTaskWireFilter", () => {
  const wire = (filter: any): any =>
    toSystemTaskWireFilter(filter as SystemTaskPagedRequestFilter) as any;

  it("passes an absent filter straight through", () => {
    expect(toSystemTaskWireFilter(undefined)).toBeUndefined();
    expect(toSystemTaskWireFilter(null)).toBeUndefined();
  });

  it("leaves the non-date fields untouched", () => {
    const filter = {
      type: { operation: FilterOperation.Contains, value: ["QUICK_SCAN"] },
      ranBy: { operation: FilterOperation.Contains, value: [-1, 4] },
      startedAt: { operation: null, value: null },
      endedAt: { operation: null, value: null },
    };

    expect(wire(filter)).toEqual(filter);
  });

  // The whole point: a local-midnight Date serializes as an instant, which the
  // server resolves to a day in *its* zone. The calendar day cannot drift.
  it("sends the local calendar day the user picked, not the instant", () => {
    const localMidnight = new Date(2026, 8, 22, 0, 0, 0, 0);

    const result = wire({
      startedAt: { operation: FilterOperation.Equals, value: localMidnight },
    });

    expect(result.startedAt.value).toBe("2026-09-22");
    expect(result.startedAt.operation).toBe(FilterOperation.Equals);
  });

  it("normalizes both bounds of a range", () => {
    const result = wire({
      startedAt: {
        operation: FilterOperation.Between,
        value: [new Date(2026, 8, 10, 0, 0, 0, 0), new Date(2026, 8, 12, 0, 0, 0, 0)],
      },
      endedAt: { operation: FilterOperation.LessThan, value: new Date(2026, 0, 3, 0, 0, 0, 0) },
    });

    expect(result.startedAt.value).toEqual(["2026-09-10", "2026-09-12"]);
    expect(result.endedAt.value).toBe("2026-01-03");
  });

  // NGXS persists the filter, so after a reload the stored value is the ISO
  // instant a Date serialized to — it has to normalize to the same day.
  it("normalizes a persisted ISO instant to the same local day as its Date", () => {
    const localMidnight = new Date(2026, 8, 22, 0, 0, 0, 0);
    const persisted = JSON.parse(JSON.stringify({ value: localMidnight })).value;

    expect(wire({ startedAt: { operation: FilterOperation.Equals, value: persisted } }).startedAt.value)
      .toBe(wire({ startedAt: { operation: FilterOperation.Equals, value: localMidnight } }).startedAt.value);
  });

  // Re-parsing a calendar day would read it as UTC midnight and shift it a day
  // west of Greenwich.
  it("leaves an already-normalized calendar day alone", () => {
    const result = wire({ endedAt: { operation: FilterOperation.Equals, value: "2026-09-22" } });

    expect(result.endedAt.value).toBe("2026-09-22");
  });

  it("leaves an unparseable value alone rather than inventing a day", () => {
    const result = wire({ startedAt: { operation: FilterOperation.Equals, value: "not-a-date" } });

    expect(result.startedAt.value).toBe("not-a-date");
  });

  it("pads a single-digit month and day", () => {
    const result = wire({
      startedAt: { operation: FilterOperation.Equals, value: new Date(2026, 0, 5, 0, 0, 0, 0) },
    });

    expect(result.startedAt.value).toBe("2026-01-05");
  });
});
