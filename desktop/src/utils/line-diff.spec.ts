import { buildSplitDiff, collapseUnchanged, inlineChange, SplitDiffRow } from "./line-diff";

const kinds = (rows: { kind: string }[]) => rows.map((row) => row.kind);

/** Rebuilds each side from the rows, which must give back the input exactly. */
const sides = (rows: SplitDiffRow[]) => ({
  left: rows.filter((row) => row.left).map((row) => row.left!.text),
  right: rows.filter((row) => row.right).map((row) => row.right!.text),
});

function lcsLength(a: string[], b: string[]): number {
  const table = Array.from({ length: a.length + 1 }, () => new Array(b.length + 1).fill(0));
  for (let i = 1; i <= a.length; i++) {
    for (let j = 1; j <= b.length; j++) {
      table[i][j] = a[i - 1] === b[j - 1] ? table[i - 1][j - 1] + 1 : Math.max(table[i - 1][j], table[i][j - 1]);
    }
  }
  return table[a.length][b.length];
}

describe("buildSplitDiff", () => {
  it("marks every line equal when nothing changed", () => {
    const rows = buildSplitDiff(["a", "b"], ["a", "b"]);

    expect(kinds(rows)).toEqual(["equal", "equal"]);
    expect(rows[1].left).toEqual({ lineNumber: 2, text: "b" });
    expect(rows[1].right).toEqual({ lineNumber: 2, text: "b" });
  });

  it("pairs a replaced line side by side as a change", () => {
    const rows = buildSplitDiff(["{", '  "name": "Costco",', "}"], ["{", '  "name": "Costco Wholesale",', "}"]);

    expect(kinds(rows)).toEqual(["equal", "changed", "equal"]);
    expect(rows[1].left).toEqual({ lineNumber: 2, text: '  "name": "Costco",' });
    expect(rows[1].right).toEqual({ lineNumber: 2, text: '  "name": "Costco Wholesale",' });
  });

  it("shows an inserted line on the right only, and numbers each side separately", () => {
    const rows = buildSplitDiff(["a", "c"], ["a", "b", "c"]);

    expect(kinds(rows)).toEqual(["equal", "added", "equal"]);
    expect(rows[1].left).toBeUndefined();
    expect(rows[1].right).toEqual({ lineNumber: 2, text: "b" });
    expect(rows[2].left).toEqual({ lineNumber: 2, text: "c" });
    expect(rows[2].right).toEqual({ lineNumber: 3, text: "c" });
  });

  it("shows a deleted line on the left only", () => {
    const rows = buildSplitDiff(["a", "b", "c"], ["a", "c"]);

    expect(kinds(rows)).toEqual(["equal", "removed", "equal"]);
    expect(rows[1].right).toBeUndefined();
    expect(rows[1].left).toEqual({ lineNumber: 2, text: "b" });
  });

  it("pairs as many lines as it can in a run and leaves the rest one-sided", () => {
    const rows = buildSplitDiff(["a", "x", "y", "c"], ["a", "z", "c"]);

    expect(kinds(rows)).toEqual(["equal", "changed", "removed", "equal"]);
  });

  it("does not report a line as changed just because a line before it was removed", () => {
    const before = ["[", "  milk,", "  eggs,", "  bread", "]"];
    const after = ["[", "  milk,", "  bread", "]"];

    expect(kinds(buildSplitDiff(before, after))).toEqual(["equal", "equal", "removed", "equal", "equal"]);
  });

  it("handles one side being empty", () => {
    expect(kinds(buildSplitDiff([], ["a", "b"]))).toEqual(["added", "added"]);
    expect(kinds(buildSplitDiff(["a", "b"], []))).toEqual(["removed", "removed"]);
    expect(buildSplitDiff([], [])).toEqual([]);
  });

  // A seeded fuzz: every row set must rebuild both inputs, and the number of
  // unchanged lines must be the longest common subsequence, i.e. the diff is
  // minimal rather than merely valid.
  it("always rebuilds both sides and keeps the longest common subsequence", () => {
    let seed = 42;
    const random = () => {
      seed = (seed * 1103515245 + 12345) % 2147483648;
      return seed / 2147483648;
    };
    const randomLines = () =>
      Array.from({ length: Math.floor(random() * 12) }, () => "abcd"[Math.floor(random() * 4)]);

    for (let run = 0; run < 300; run++) {
      const before = randomLines();
      const after = randomLines();
      const rows = buildSplitDiff(before, after);

      expect(sides(rows)).toEqual({ left: before, right: after });
      expect(rows.filter((row) => row.kind === "equal").length).toBe(lcsLength(before, after));
    }
  });
});

describe("collapseUnchanged", () => {
  const equalRows = (count: number) => buildSplitDiff(
    Array.from({ length: count }, (_, i) => `line ${i}`),
    Array.from({ length: count }, (_, i) => `line ${i}`),
  );

  it("keeps the context lines around a change and collapses the rest", () => {
    const before = Array.from({ length: 20 }, (_, i) => `line ${i}`);
    const after = [...before];
    after[10] = "changed";

    const display = collapseUnchanged(buildSplitDiff(before, after), 3);

    expect(kinds(display)).toEqual([
      "collapsed", "equal", "equal", "equal", "changed", "equal", "equal", "equal", "collapsed",
    ]);
    expect(display[0]).toEqual({ kind: "collapsed", count: 7 });
    expect(display[8]).toEqual({ kind: "collapsed", count: 6 });
  });

  it("shows a single hidden line instead of a marker for it", () => {
    const before = ["a", "b", "c", "d", "e"];
    const after = ["A", "b", "c", "d", "E"];

    expect(kinds(collapseUnchanged(buildSplitDiff(before, after), 1))).toEqual([
      "changed", "equal", "equal", "equal", "changed",
    ]);
  });

  it("collapses everything when nothing changed", () => {
    expect(collapseUnchanged(equalRows(5))).toEqual([{ kind: "collapsed", count: 5 }]);
  });
});

describe("inlineChange", () => {
  it("isolates the part of each line that differs", () => {
    const { left, right } = inlineChange('  "name": "Costco",', '  "name": "Costco Wholesale",');

    expect(left).toEqual({ prefix: '  "name": "Costco', changed: "", suffix: '",' });
    expect(right).toEqual({ prefix: '  "name": "Costco', changed: " Wholesale", suffix: '",' });
  });

  it("never lets the prefix and suffix overlap", () => {
    const { left, right } = inlineChange("aa", "aaa");

    expect(left.prefix + left.changed + left.suffix).toBe("aa");
    expect(right.prefix + right.changed + right.suffix).toBe("aaa");
    expect(right.changed).toBe("a");
  });

  it("marks the whole line when nothing is shared", () => {
    expect(inlineChange("abc", "xyz").right).toEqual({ prefix: "", changed: "xyz", suffix: "" });
  });
});
