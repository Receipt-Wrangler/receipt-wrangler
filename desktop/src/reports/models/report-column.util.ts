/**
 * Column helpers shared by the column picker. The engine requires a column's
 * `name` to be a plain identifier (formulas reference it) distinct from every
 * other column's name, while `label` stays the free-text heading. These keep the
 * builder producing engine-valid names and give the formula editor lightweight,
 * non-eval feedback (the backend is the real validator and returns 400 on a bad
 * spec).
 */

const STARTS_WITH_LETTER_OR_UNDERSCORE = /^[A-Za-z_]/;
const NON_IDENTIFIER_CHARS = /[^A-Za-z0-9_]+/g;
const IDENTIFIER_TOKEN = /[A-Za-z_][A-Za-z0-9_]*/g;

// expr-lang reserved words the engine rejects as a bare column reference.
const RESERVED = new Set([
  "true", "false", "nil", "null", "and", "or", "not", "in", "matches", "contains",
]);

/**
 * Derives a valid, unique engine column name from a user label. A single-word
 * label (e.g. "Subtotal") passes through unchanged so formulas can reference it by
 * the name the user sees; anything else is sanitized, and collisions get a numeric
 * suffix. Assigned once at column creation and kept stable so editing a label never
 * breaks a formula that references the name.
 */
export function deriveColumnName(label: string, existingNames: string[]): string {
  let base = (label ?? "").trim().replace(NON_IDENTIFIER_CHARS, "_").replace(/^_+|_+$/g, "");
  if (!base || !STARTS_WITH_LETTER_OR_UNDERSCORE.test(base) || RESERVED.has(base.toLowerCase())) {
    base = "col_" + base;
  }

  const taken = new Set(existingNames);
  let name = base;
  let suffix = 2;
  while (taken.has(name)) {
    name = `${base}_${suffix++}`;
  }
  return name;
}

export type FormulaStatusKind = "hint" | "error" | "ok";

export interface FormulaStatus {
  ok: boolean;
  kind: FormulaStatusKind;
  message: string;
}

/**
 * Lightweight, eval-free validation of a formula expression for inline feedback:
 * every referenced identifier must be an available column name, parentheses must
 * balance, and only arithmetic characters may remain. It is intentionally lenient
 * (the backend fully parses and rejects a bad expression) — it just steers the user.
 */
export function validateFormulaExpr(expr: string, availableNames: string[]): FormulaStatus {
  const trimmed = (expr ?? "").trim();
  if (!trimmed) {
    return { ok: false, kind: "hint", message: "Reference columns by name, e.g. Subtotal + Hst" };
  }

  const known = new Set(availableNames);
  const identifiers = trimmed.match(IDENTIFIER_TOKEN) ?? [];
  const unknown = identifiers.find((token) => !known.has(token));
  if (unknown) {
    return { ok: false, kind: "error", message: "Unknown column: " + unknown };
  }

  const arithmetic = trimmed.replace(IDENTIFIER_TOKEN, "1");
  if (/[^0-9.+\-*/()\s]/.test(arithmetic)) {
    return { ok: false, kind: "error", message: "Only arithmetic over columns is allowed" };
  }
  if (!parenthesesBalanced(arithmetic)) {
    return { ok: false, kind: "error", message: "Check your parentheses" };
  }

  return { ok: true, kind: "ok", message: "Valid — recomputes from these columns each row" };
}

function parenthesesBalanced(text: string): boolean {
  let depth = 0;
  for (const char of text) {
    if (char === "(") {
      depth++;
    } else if (char === ")") {
      depth--;
      if (depth < 0) {
        return false;
      }
    }
  }
  return depth === 0;
}
