import { deriveColumnName, validateFormulaExpr } from "./report-column.util";

describe("deriveColumnName", () => {
  it("passes a single-word label through unchanged", () => {
    expect(deriveColumnName("Subtotal", [])).toBe("Subtotal");
  });

  it("sanitizes non-identifier characters to underscores", () => {
    expect(deriveColumnName("Avg / Receipt", [])).toBe("Avg_Receipt");
  });

  it("prefixes labels that don't start with a letter", () => {
    expect(deriveColumnName("2026 Total", [])).toBe("col_2026_Total");
  });

  it("prefixes reserved words", () => {
    expect(deriveColumnName("in", [])).toBe("col_in");
  });

  it("suffixes to stay unique", () => {
    expect(deriveColumnName("Total", ["Total"])).toBe("Total_2");
    expect(deriveColumnName("Total", ["Total", "Total_2"])).toBe("Total_3");
  });

  it("falls back to col_ for an empty label", () => {
    expect(deriveColumnName("   ", [])).toBe("col_");
  });
});

describe("validateFormulaExpr", () => {
  const columns = ["Subtotal", "Hst", "Count"];

  it("returns a hint for an empty expression", () => {
    const status = validateFormulaExpr("", columns);
    expect(status.ok).toBe(false);
    expect(status.kind).toBe("hint");
  });

  it("accepts arithmetic over known columns", () => {
    const status = validateFormulaExpr("Subtotal + Hst", columns);
    expect(status.ok).toBe(true);
    expect(status.kind).toBe("ok");
  });

  it("flags an unknown column", () => {
    const status = validateFormulaExpr("Subtotal + Gst", columns);
    expect(status.ok).toBe(false);
    expect(status.kind).toBe("error");
    expect(status.message).toContain("Gst");
  });

  it("rejects unbalanced parentheses", () => {
    const status = validateFormulaExpr("(Subtotal + Hst", columns);
    expect(status.ok).toBe(false);
    expect(status.kind).toBe("error");
  });

  it("rejects disallowed characters", () => {
    const status = validateFormulaExpr("Subtotal % Hst", columns);
    expect(status.ok).toBe(false);
  });
});
