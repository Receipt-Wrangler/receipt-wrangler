import { FilterOperation } from "../open-api";
import { RECEIPT_DATE_FILTER_KEYS } from "../constants/receipt-filter-fields.constant";
import { toDateWireFilter } from "./date-wire-filter";
import { SYSTEM_TASK_DATE_FILTER_KEYS } from "./system-task-filter";

describe("toDateWireFilter", () => {
  it("passes an absent filter straight through", () => {
    expect(toDateWireFilter(undefined, RECEIPT_DATE_FILTER_KEYS)).toBeUndefined();
    expect(toDateWireFilter(null, RECEIPT_DATE_FILTER_KEYS)).toBeUndefined();
  });

  it("converts every receipt date field and leaves the rest alone", () => {
    const filter = {
      date: { operation: FilterOperation.Equals, value: new Date(2026, 8, 30) },
      resolvedDate: { operation: FilterOperation.GreaterThan, value: new Date(2026, 0, 5) },
      createdAt: {
        operation: FilterOperation.Between,
        value: [new Date(2026, 8, 1), new Date(2026, 8, 30, 23, 59, 59, 999)],
      },
      amount: { operation: FilterOperation.GreaterThan, value: 5 },
      name: { operation: FilterOperation.Contains, value: "2026-09-30T00:00:00Z" },
    };

    const wire = toDateWireFilter(filter, RECEIPT_DATE_FILTER_KEYS) as any;

    expect(wire.date).toEqual({ operation: FilterOperation.Equals, value: "2026-09-30" });
    expect(wire.resolvedDate.value).toBe("2026-01-05");
    expect(wire.createdAt.value).toEqual(["2026-09-01", "2026-09-30"]);
    expect(wire.amount).toBe(filter.amount);
    expect(wire.name).toBe(filter.name);
  });

  it("never mutates the filter it is given", () => {
    const picked = new Date(2026, 8, 30);
    const filter = { date: { operation: FilterOperation.Equals, value: picked } };

    toDateWireFilter(filter, RECEIPT_DATE_FILTER_KEYS);

    expect(filter.date.value).toBe(picked);
  });

  // NGXS persists a filter through JSON, so a reloaded value is the ISO instant
  // its Date serialized to. It is still the local day the user picked.
  it("reads a persisted ISO instant on the same local day as its Date", () => {
    const picked = new Date(2026, 8, 30);
    const persisted = JSON.parse(JSON.stringify(picked));

    expect((toDateWireFilter({ date: { value: persisted } }, ["date"]) as any).date.value).toBe("2026-09-30");
  });

  it("keeps empty values empty", () => {
    const wire = toDateWireFilter(
      { date: { operation: null, value: "" }, createdAt: { operation: null, value: null } },
      RECEIPT_DATE_FILTER_KEYS
    ) as any;

    expect(wire.date.value).toBe("");
    expect(wire.createdAt.value).toBeNull();
    expect("resolvedDate" in wire).toBe(false);
  });

  it("converts the system task date fields through the same rule", () => {
    const wire = toDateWireFilter(
      {
        startedAt: { operation: FilterOperation.Equals, value: new Date(2026, 8, 22) },
        endedAt: { operation: FilterOperation.Between, value: [new Date(2026, 8, 1), new Date(2026, 8, 2)] },
        type: { operation: FilterOperation.Contains, value: ["QUICK_SCAN"] },
      },
      SYSTEM_TASK_DATE_FILTER_KEYS
    ) as any;

    expect(wire.startedAt.value).toBe("2026-09-22");
    expect(wire.endedAt.value).toEqual(["2026-09-01", "2026-09-02"]);
    expect(wire.type.value).toEqual(["QUICK_SCAN"]);
  });
});
