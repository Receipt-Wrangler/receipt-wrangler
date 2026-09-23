import {
  changedTopLevelKeys,
  parseReceiptUpdateDescription,
  receiptUpdateBeforeState,
  toJsonLines,
} from "./receipt-update-description";

// What UpdateReceipt stores: each side is itself a JSON string.
const storedDescription = (before: object, after: object) =>
  JSON.stringify({ before: JSON.stringify(before), after: JSON.stringify(after) });

describe("parseReceiptUpdateDescription", () => {
  it("parses the double-encoded before/after pair the API stores", () => {
    const before = { id: 1, name: 'Costco "Wholesale"', comments: [{ comment: "a {quoted} note" }] };
    const after = { id: 1, name: "Costco", comments: [] };

    expect(parseReceiptUpdateDescription(storedDescription(before, after))).toEqual({ before, after, version: 1 });
  });

  it("accepts sides that are already objects", () => {
    const description = JSON.stringify({ before: { id: 1 }, after: { id: 1, name: "x" } });

    expect(parseReceiptUpdateDescription(description))
      .toEqual({ before: { id: 1 }, after: { id: 1, name: "x" }, version: 1 });
  });

  // Rows written before the marker existed have no version key.
  it("reads the version, and treats a description without one as version 1", () => {
    const pair = { before: "{}", after: "{}" };

    expect(parseReceiptUpdateDescription(JSON.stringify({ ...pair, version: 2 }))?.version).toBe(2);
    expect(parseReceiptUpdateDescription(JSON.stringify(pair))?.version).toBe(1);
    expect(parseReceiptUpdateDescription(JSON.stringify({ ...pair, version: "2" }))?.version).toBe(1);
  });

  it("reads the earlier copy the API rebuilt an older row's before from", () => {
    const beforeSource = { type: "RECEIPT_UPLOADED", systemTaskId: 4, recordedAt: "2026-09-20T10:00:00Z" };

    expect(parseReceiptUpdateDescription(JSON.stringify({ before: "{}", after: "{}", beforeSource }))?.beforeSource)
      .toEqual(beforeSource);
    expect(parseReceiptUpdateDescription(JSON.stringify({ before: "{}", after: "{}", beforeSource: { type: 1 } })))
      .not.toHaveProperty("beforeSource");
  });

  it("returns undefined for a failed update's plain error text", () => {
    expect(parseReceiptUpdateDescription("record not found")).toBeUndefined();
  });

  it("returns undefined when either side is missing or is not an object", () => {
    expect(parseReceiptUpdateDescription(JSON.stringify({ before: "{}" }))).toBeUndefined();
    expect(parseReceiptUpdateDescription(JSON.stringify({ before: "[]", after: "{}" }))).toBeUndefined();
    expect(parseReceiptUpdateDescription(JSON.stringify({ before: "not json", after: "{}" }))).toBeUndefined();
    expect(parseReceiptUpdateDescription("null")).toBeUndefined();
  });

  it("returns undefined for an empty description", () => {
    expect(parseReceiptUpdateDescription("")).toBeUndefined();
    expect(parseReceiptUpdateDescription(undefined)).toBeUndefined();
  });
});

describe("changedTopLevelKeys", () => {
  it("lists the keys whose values differ, in serialization order", () => {
    const before = { id: 1, name: "a", amount: "1", tags: [{ id: 1 }] };
    const after = { id: 1, name: "b", amount: "1", tags: [{ id: 1 }, { id: 2 }] };

    expect(changedTopLevelKeys(before, after)).toEqual(["name", "tags"]);
  });

  it("ignores updatedAt, which every save changes", () => {
    expect(changedTopLevelKeys({ updatedAt: "1" }, { updatedAt: "2" })).toEqual([]);
  });

  // Saving upserts the receipt's categories and tags, which bumps their
  // timestamps too; that alone must not list them as changed.
  it("ignores updatedAt inside nested records", () => {
    const before = { categories: [{ id: 1, name: "Groceries", updatedAt: "1" }] };
    const after = { categories: [{ id: 1, name: "Groceries", updatedAt: "2" }] };

    expect(changedTopLevelKeys(before, after)).toEqual([]);
    expect(changedTopLevelKeys(before, { categories: [{ ...after.categories[0], name: "Food" }] }))
      .toEqual(["categories"]);
  });

  // Saving deletes and recreates items (linked items included) and custom
  // field values, so each comes back with a new id, timestamps and creator.
  describe("records the save recreates", () => {
    const record = { id: 5, createdAt: "t1", updatedAt: "t1", createdBy: null, createdByString: "" };
    const recreated = { id: 9, createdAt: "t2", updatedAt: "t2", createdBy: 1, createdByString: "Admin" };
    const pizza = (bookkeeping: object, overrides: object = {}) => ({
      ...bookkeeping,
      name: "Pizza",
      amount: "20",
      receiptId: 4,
      categories: [{ ...record, name: "Groceries" }],
      linkedItems: [{ ...bookkeeping, name: "Pizza share", amount: "10", receiptId: 4 }],
      ...overrides,
    });
    const poNumber = (bookkeeping: object, stringValue: string) => ({
      ...bookkeeping,
      receiptId: 4,
      customFieldId: 1,
      customField: { ...record, name: "PO Number", type: "TEXT" },
      stringValue,
    });

    it("does not list an item or custom field value that was only recreated", () => {
      const before = { receiptItems: [pizza(record)], customFields: [poNumber(record, "PO-1182")] };
      const after = { receiptItems: [pizza(recreated)], customFields: [poNumber(recreated, "PO-1182")] };

      expect(changedTopLevelKeys(before, after)).toEqual([]);
    });

    it("still lists a real edit to one", () => {
      const before = { receiptItems: [pizza(record)], customFields: [poNumber(record, "PO-1182")] };
      const after = {
        receiptItems: [pizza(recreated, { amount: "18" })],
        customFields: [poNumber(recreated, "PO-1183")],
      };

      expect(changedTopLevelKeys(before, after)).toEqual(["receiptItems", "customFields"]);
    });
  });

  it("includes a key present on only one side", () => {
    expect(changedTopLevelKeys({ id: 1 }, { id: 1, resolvedDate: "2026-09-01" })).toEqual(["resolvedDate"]);
  });
});

describe("receiptUpdateBeforeState", () => {
  const snapshots = { before: {}, after: {} };
  const beforeSource = { type: "RECEIPT_UPDATED" as const, systemTaskId: 4, recordedAt: "2026-09-20T10:00:00Z" };

  it("trusts a version 2 row's own before", () => {
    expect(receiptUpdateBeforeState({ ...snapshots, version: 2 })).toBe("complete");
  });

  it("marks a version 1 row the API rebuilt from an earlier copy", () => {
    expect(receiptUpdateBeforeState({ ...snapshots, version: 1, beforeSource })).toBe("rebuilt");
  });

  it("marks a version 1 row with no earlier copy as incomplete", () => {
    expect(receiptUpdateBeforeState({ ...snapshots, version: 1 })).toBe("incomplete");
  });
});

describe("toJsonLines", () => {
  it("pretty-prints with two-space indentation, one line per entry", () => {
    expect(toJsonLines({ id: 1, tags: ["a"] })).toEqual(["{", '  "id": 1,', '  "tags": [', '    "a"', "  ]", "}"]);
  });
});
