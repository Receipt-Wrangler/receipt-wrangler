import { ReportColumn, ReportDetail, ReportPeriod, ReportRequestCommand } from "../../open-api";
import {
  columnCount,
  detailSummary,
  fieldLabel,
  formatChips,
  groupingSummary,
  scopeNames,
} from "./report-template-summary";

function command(overrides: Partial<ReportRequestCommand> = {}): ReportRequestCommand {
  return {
    name: "R",
    groupIds: ["1"],
    period: { preset: ReportPeriod.PresetEnum.ThisMonth },
    filter: {},
    groupBy: [],
    detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "category" },
    columns: [],
    subtotals: true,
    grandTotals: true,
    formats: [],
    ...overrides,
  };
}

describe("report-template-summary", () => {
  it("labels built-in field keys and falls back to the raw key", () => {
    expect(fieldLabel("category")).toBe("Category");
    expect(fieldLabel("date_month")).toBe("Month");
    expect(fieldLabel("custom_42")).toBe("custom_42");
    expect(fieldLabel(undefined)).toBe("");
  });

  it("resolves custom field keys through the caller's resolver", () => {
    const resolve = (key: string) => (key === "custom_42" ? "Due Date" : undefined);

    expect(fieldLabel("custom_42", resolve)).toBe("Due Date");
    // A resolver that doesn't know the key still falls back to the raw key.
    expect(fieldLabel("custom_99", resolve)).toBe("custom_99");
    // A built-in never consults the resolver.
    expect(fieldLabel("category", () => "Wrong")).toBe("Category");
  });

  it("counts columns", () => {
    expect(columnCount(command({ columns: [] }))).toBe(0);
    expect(
      columnCount(
        command({
          columns: [
            { kind: ReportColumn.KindEnum.Aggregate, name: "C", label: "C", aggFunc: ReportColumn.AggFuncEnum.Count },
            { kind: ReportColumn.KindEnum.Aggregate, name: "T", label: "T", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
          ],
        })
      )
    ).toBe(2);
  });

  it("summarizes grouping levels or reports none", () => {
    expect(groupingSummary(command({ groupBy: [] }))).toBe("No grouping");
    expect(groupingSummary(command({ groupBy: ["group", "category"] }))).toBe("Group, Category");
  });

  it("shows a grouping level's renamed column heading", () => {
    // The list must agree with the report it describes: a level whose column the
    // report renames shows that name, not the underlying field's.
    expect(
      groupingSummary(
        command({ groupBy: ["group", "category"], groupByLabels: { category: "Expense Type" } })
      )
    ).toBe("Group, Expense Type");

    // A blank or unrelated entry changes nothing.
    expect(
      groupingSummary(
        command({ groupBy: ["group"], groupByLabels: { group: "  ", status: "Ignored" } })
      )
    ).toBe("Group");
  });

  it("prefers a renamed heading over a resolved custom-field name", () => {
    const resolve = (key: string) => ({ custom_42: "Due Date" })[key];
    expect(
      groupingSummary(
        command({ groupBy: ["custom_42"], groupByLabels: { custom_42: "Deadline" } }),
        resolve
      )
    ).toBe("Deadline");
  });

  it("names custom fields in the grouping and detail summaries", () => {
    const resolve = (key: string) => ({ custom_42: "Due Date", custom_7: "HST" })[key];

    expect(groupingSummary(command({ groupBy: ["group", "custom_42"] }), resolve)).toBe(
      "Group, Due Date"
    );
    expect(
      detailSummary(command({ detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "custom_7" } }), resolve)
    ).toBe("Aggregate by HST");
  });

  it("summarizes the detail mode", () => {
    expect(detailSummary(command({ detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "tag" } }))).toBe("Aggregate by Tag");
    expect(detailSummary(command({ detail: { mode: ReportDetail.ModeEnum.Records } }))).toBe("Record-level");
  });

  it("renders formats as ordered uppercase chips", () => {
    expect(
      formatChips(
        command({
          formats: [ReportRequestCommand.FormatsEnum.Pdf, ReportRequestCommand.FormatsEnum.Csv],
        })
      )
    ).toEqual(["CSV", "PDF"]);
    expect(formatChips(command({ formats: [] }))).toEqual([]);
  });

  it("resolves scope group names, falling back to the raw id", () => {
    const names: Record<string, string> = { "1": "Household", "2": "Work" };
    expect(scopeNames(command({ groupIds: ["1", "2", "99"] }), (id) => names[id])).toEqual([
      "Household",
      "Work",
      "99",
    ]);
  });
});
