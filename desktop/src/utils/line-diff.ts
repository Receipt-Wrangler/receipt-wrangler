export interface DiffLine {
  /** 1-based line number on its own side. */
  lineNumber: number;
  text: string;
}

/**
 * One row of a side-by-side diff. `changed` pairs a removed line with the added
 * line that replaced it; `removed` has no right side and `added` no left side.
 */
export interface SplitDiffRow {
  kind: "equal" | "changed" | "removed" | "added";
  left?: DiffLine;
  right?: DiffLine;
}

/** A run of unchanged rows hidden by {@link collapseUnchanged}. */
export interface CollapsedDiffRow {
  kind: "collapsed";
  count: number;
}

/** A line split around the part that differs from the line it is paired with. */
export interface InlineSegments {
  prefix: string;
  changed: string;
  suffix: string;
}

type EditOp = "equal" | "delete" | "insert";

/**
 * Diffs two lists of lines into side-by-side rows, like a split-view diff.
 *
 * Within each run of edits between unchanged lines, removed lines are paired
 * with added lines in order as `changed` rows; whatever is left over on either
 * side becomes `removed` or `added`.
 */
export function buildSplitDiff(before: string[], after: string[]): SplitDiffRow[] {
  const rows: SplitDiffRow[] = [];
  let removed: DiffLine[] = [];
  let added: DiffLine[] = [];
  let beforeIndex = 0;
  let afterIndex = 0;

  const flushEdits = () => {
    for (let i = 0; i < Math.max(removed.length, added.length); i++) {
      const left = removed[i];
      const right = added[i];
      rows.push({ kind: left && right ? "changed" : left ? "removed" : "added", left, right });
    }
    removed = [];
    added = [];
  };

  for (const op of editScript(before, after)) {
    if (op === "equal") {
      flushEdits();
      rows.push({
        kind: "equal",
        left: toDiffLine(beforeIndex, before[beforeIndex]),
        right: toDiffLine(afterIndex, after[afterIndex]),
      });
      beforeIndex++;
      afterIndex++;
    } else if (op === "delete") {
      removed.push(toDiffLine(beforeIndex, before[beforeIndex]));
      beforeIndex++;
    } else {
      added.push(toDiffLine(afterIndex, after[afterIndex]));
      afterIndex++;
    }
  }
  flushEdits();

  return rows;
}

/**
 * Hides unchanged rows that are more than `context` rows away from a change,
 * replacing each hidden run with a single `collapsed` row. A run of one row is
 * shown rather than replaced by a marker the same size as it.
 */
export function collapseUnchanged<T extends { kind: SplitDiffRow["kind"] }>(
  rows: T[],
  context = 3,
): (T | CollapsedDiffRow)[] {
  const visible = rows.map(() => false);
  rows.forEach((row, index) => {
    if (row.kind === "equal") {
      return;
    }
    const last = Math.min(rows.length - 1, index + context);
    for (let i = Math.max(0, index - context); i <= last; i++) {
      visible[i] = true;
    }
  });

  const result: (T | CollapsedDiffRow)[] = [];
  let hidden: T[] = [];
  const flushHidden = () => {
    if (hidden.length === 1) {
      result.push(hidden[0]);
    } else if (hidden.length > 1) {
      result.push({ kind: "collapsed", count: hidden.length });
    }
    hidden = [];
  };

  rows.forEach((row, index) => {
    if (visible[index]) {
      flushHidden();
      result.push(row);
    } else {
      hidden.push(row);
    }
  });
  flushHidden();

  return result;
}

/**
 * Splits a changed pair of lines around their common prefix and suffix, so only
 * the part that actually differs needs the stronger highlight.
 */
export function inlineChange(left: string, right: string): { left: InlineSegments; right: InlineSegments } {
  const maxShared = Math.min(left.length, right.length);

  let prefixLength = 0;
  while (prefixLength < maxShared && left[prefixLength] === right[prefixLength]) {
    prefixLength++;
  }

  let suffixLength = 0;
  while (
    suffixLength < maxShared - prefixLength &&
    left[left.length - 1 - suffixLength] === right[right.length - 1 - suffixLength]
  ) {
    suffixLength++;
  }

  const split = (text: string): InlineSegments => ({
    prefix: text.slice(0, prefixLength),
    changed: text.slice(prefixLength, text.length - suffixLength),
    suffix: text.slice(text.length - suffixLength),
  });

  return { left: split(left), right: split(right) };
}

function toDiffLine(index: number, text: string): DiffLine {
  return { lineNumber: index + 1, text };
}

/**
 * The shortest edit script turning `before` into `after`. The common prefix and
 * suffix are trimmed first — cheap, and for a receipt edit usually most of the
 * document — and Myers' O(ND) algorithm handles what is left.
 */
function editScript(before: string[], after: string[]): EditOp[] {
  let start = 0;
  while (start < before.length && start < after.length && before[start] === after[start]) {
    start++;
  }

  let beforeEnd = before.length;
  let afterEnd = after.length;
  while (beforeEnd > start && afterEnd > start && before[beforeEnd - 1] === after[afterEnd - 1]) {
    beforeEnd--;
    afterEnd--;
  }

  const middle = myers(before.slice(start, beforeEnd), after.slice(start, afterEnd));
  const trailing = before.length - beforeEnd;

  return [
    ...new Array<EditOp>(start).fill("equal"),
    ...middle,
    ...new Array<EditOp>(trailing).fill("equal"),
  ];
}

/**
 * Myers' diff. `frontier[k]` is the furthest x reached on diagonal k = x - y.
 * Each round's frontier is kept for the backtrack, but only the window that
 * round can read (diagonals -d-1 .. d+1), so memory is O(D²) rather than
 * O(D·(N+M)).
 */
function myers(before: string[], after: string[]): EditOp[] {
  const n = before.length;
  const m = after.length;
  if (n === 0) {
    return new Array<EditOp>(m).fill("insert");
  }
  if (m === 0) {
    return new Array<EditOp>(n).fill("delete");
  }

  const offset = n + m + 1;
  const frontier = new Int32Array(2 * offset + 1);
  const history: Int32Array[] = [];

  for (let d = 0; d <= n + m; d++) {
    history.push(frontier.slice(offset - d - 1, offset + d + 2));

    for (let k = -d; k <= d; k += 2) {
      const fromAbove = k === -d || (k !== d && frontier[offset + k - 1] < frontier[offset + k + 1]);
      let x = fromAbove ? frontier[offset + k + 1] : frontier[offset + k - 1] + 1;
      let y = x - k;
      while (x < n && y < m && before[x] === after[y]) {
        x++;
        y++;
      }
      frontier[offset + k] = x;

      if (x >= n && y >= m) {
        return backtrack(history, n, m);
      }
    }
  }

  // Unreachable: d = n + m always reaches the end.
  return [];
}

function backtrack(history: Int32Array[], n: number, m: number): EditOp[] {
  const ops: EditOp[] = [];
  let x = n;
  let y = m;

  for (let d = history.length - 1; d >= 0; d--) {
    const window = history[d];
    const at = (k: number) => window[k + d + 1];
    const k = x - y;

    const fromAbove = k === -d || (k !== d && at(k - 1) < at(k + 1));
    const previousK = fromAbove ? k + 1 : k - 1;
    const previousX = at(previousK);
    const previousY = previousX - previousK;

    while (x > previousX && y > previousY) {
      ops.push("equal");
      x--;
      y--;
    }
    if (d > 0) {
      ops.push(fromAbove ? "insert" : "delete");
    }
    x = previousX;
    y = previousY;
  }

  return ops.reverse();
}
