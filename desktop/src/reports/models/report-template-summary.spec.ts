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
