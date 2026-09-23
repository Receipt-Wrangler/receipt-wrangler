import { changedTopLevelKeys, parseReceiptUpdateDescription, toJsonLines } from "./receipt-update-description";

// What UpdateReceipt stores: each side is itself a JSON string.
const storedDescription = (before: object, after: object) =>
  JSON.stringify({ before: JSON.stringify(before), after: JSON.stringify(after) });

describe("parseReceiptUpdateDescription", () => {
  it("parses the double-encoded before/after pair the API stores", () => {
    const before = { id: 1, name: 'Costco "Wholesale"', comments: [{ comment: "a {quoted} note" }] };
    const after = { id: 1, name: "Costco", comments: [] };

    expect(parseReceiptUpdateDescription(storedDescription(before, after))).toEqual({ before, after });
  });

  it("accepts sides that are already objects", () => {
    const description = JSON.stringify({ before: { id: 1 }, after: { id: 1, name: "x" } });

    expect(parseReceiptUpdateDescription(description)).toEqual({ before: { id: 1 }, after: { id: 1, name: "x" } });
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

describe("toJsonLines", () => {
  it("pretty-prints with two-space indentation, one line per entry", () => {
    expect(toJsonLines({ id: 1, tags: ["a"] })).toEqual(["{", '  "id": 1,', '  "tags": [', '    "a"', "  ]", "}"]);
  });
});
