export type ReceiptSnapshot = Record<string, unknown>;

export interface ReceiptUpdateSnapshots {
  before: ReceiptSnapshot;
  after: ReceiptSnapshot;
}

/**
 * Record-keeping fields (the API's BaseModel) the changed-keys summary ignores
 * at every depth. None of them change because the user edited something:
 * saving a receipt deletes and recreates its items and custom field values, so
 * those get new ids and timestamps on every save, and it bumps `updatedAt` on
 * the receipt and on its categories and tags. Counting them would list
 * `receiptItems`, `customFields`, `categories` and `tags` on every row.
 */
const SUMMARY_IGNORED_KEYS = new Set(["id", "createdAt", "updatedAt", "createdBy", "createdByString"]);

/**
 * Reads a RECEIPT_UPDATED system task's description.
 *
 * The API stores it double-encoded — `{"before":"{\"id\":1,...}","after":"..."}`,
 * each side being `Receipt.ToString()` — which is why the generic pretty-json
 * pipe cannot display it. Each side is parsed once more here (an already-parsed
 * object is accepted too).
 *
 * Never throws: a failed update stores its plain error text instead, and that
 * (or anything else that is not a before/after pair) returns `undefined` so the
 * caller can fall back to showing the text as-is.
 */
export function parseReceiptUpdateDescription(
  description?: string | null,
): ReceiptUpdateSnapshots | undefined {
  if (!description) {
    return undefined;
  }

  try {
    const outer = JSON.parse(description);
    const before = parseSnapshot(outer?.before);
    const after = parseSnapshot(outer?.after);

    return before && after ? { before, after } : undefined;
  } catch {
    return undefined;
  }
}

/**
 * The top-level receipt keys whose values differ between the two snapshots,
 * in the order the receipt serializes them. Record-keeping fields are ignored
 * (see {@link SUMMARY_IGNORED_KEYS}), so a key is listed only for a real edit.
 */
export function changedTopLevelKeys(before: ReceiptSnapshot, after: ReceiptSnapshot): string[] {
  const keys = new Set([...Object.keys(before), ...Object.keys(after)]);
  const withoutRecordKeeping = (value: unknown) =>
    JSON.stringify(value, (key, nested) => (SUMMARY_IGNORED_KEYS.has(key) ? undefined : nested));

  return [...keys].filter(
    (key) =>
      !SUMMARY_IGNORED_KEYS.has(key) &&
      withoutRecordKeeping(before[key]) !== withoutRecordKeeping(after[key]),
  );
}

/** The snapshot as the lines a diff compares: 2-space pretty-printed JSON. */
export function toJsonLines(snapshot: ReceiptSnapshot): string[] {
  return JSON.stringify(snapshot, null, 2).split("\n");
}

function parseSnapshot(value: unknown): ReceiptSnapshot | undefined {
  const parsed = typeof value === "string" ? JSON.parse(value) : value;

  return parsed !== null && typeof parsed === "object" && !Array.isArray(parsed)
    ? (parsed as ReceiptSnapshot)
    : undefined;
}
