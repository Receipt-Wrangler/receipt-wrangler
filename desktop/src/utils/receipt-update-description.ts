export type ReceiptSnapshot = Record<string, unknown>;

export interface ReceiptUpdateSnapshots {
  before: ReceiptSnapshot;
  after: ReceiptSnapshot;
}

/**
 * Ignored by the changed-keys summary, at every depth. The save itself bumps
 * the receipt's `updatedAt`, and upserting its categories and tags bumps
 * theirs, so counting it would list `categories` and `tags` on every row.
 */
const SUMMARY_IGNORED_KEY = "updatedAt";

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
 * in the order the receipt serializes them. `updatedAt` changes are ignored.
 */
export function changedTopLevelKeys(before: ReceiptSnapshot, after: ReceiptSnapshot): string[] {
  const keys = new Set([...Object.keys(before), ...Object.keys(after)]);
  const withoutUpdatedAt = (value: unknown) =>
    JSON.stringify(value, (key, nested) => (key === SUMMARY_IGNORED_KEY ? undefined : nested));

  return [...keys].filter(
    (key) =>
      key !== SUMMARY_IGNORED_KEY &&
      withoutUpdatedAt(before[key]) !== withoutUpdatedAt(after[key]),
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
