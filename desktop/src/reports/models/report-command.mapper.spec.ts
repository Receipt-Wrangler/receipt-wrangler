import { ReportColumn, ReportDetail, ReportPeriod } from "../../open-api";
import { ReportBuilderValue, toReportRequestCommand } from "./report-command.mapper";

function baseValue(): ReportBuilderValue {
  return {
    name: "My Report",
    scope: ["1", "2"],
    period: { preset: ReportPeriod.PresetEnum.ThisMonth, startDate: null, endDate: null },
    filter: {},
    groupBy: ["group", "category"],
    detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "category" },
    columns: [
      { id: "a", kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
      { id: "b", kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
      { id: "c", kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count, measure: "amount" },
      { id: "d", kind: ReportColumn.KindEnum.Formula, name: "Avg", label: "Avg", expr: "Total / Count" },
    ],
    subtotals: true,
    grandTotals: false,
    document: { title: "", intro: "", footer: "" },
    formats: { csv: false, xlsx: true, pdf: true },
  };
}

describe("toReportRequestCommand", () => {
  it("maps the top-level fields and scope", () => {
    const command = toReportRequestCommand(baseValue());
    expect(command.name).toBe("My Report");
    expect(command.groupIds).toEqual(["1", "2"]);
    expect(command.groupBy).toEqual(["group", "category"]);
    expect(command.subtotals).toBe(true);
    expect(command.grandTotals).toBe(false);
  });

  it("emits a preset period without dates", () => {
    const command = toReportRequestCommand(baseValue());
    expect(command.period).toEqual({ preset: ReportPeriod.PresetEnum.ThisMonth });
  });

  it("emits a custom period with formatted YYYY-MM-DD dates", () => {
    const value = baseValue();
    value.period = {
      preset: ReportPeriod.PresetEnum.Custom,
      startDate: new Date(2026, 4, 1),
      endDate: new Date(2026, 4, 31),
    };
    const command = toReportRequestCommand(value);
    expect(command.period).toEqual({
      preset: ReportPeriod.PresetEnum.Custom,
      startDate: "2026-05-01",
      endDate: "2026-05-31",
    });
  });

  it("orders formats csv, xlsx, pdf and drops the unselected", () => {
    const command = toReportRequestCommand(baseValue());
    expect(command.formats).toEqual(["xlsx", "pdf"]);
  });

  it("maps each column kind to only its relevant fields", () => {
    const command = toReportRequestCommand(baseValue());
    expect(command.columns[0]).toEqual({ kind: "dimension", name: "Category", label: "Category", field: "category" });
    expect(command.columns[1]).toEqual({ kind: "aggregate", name: "Total", label: "Total", aggFunc: "SUM", measure: "amount" });
    // COUNT drops the measure even when one is present on the form.
    expect(command.columns[2]).toEqual({ kind: "aggregate", name: "Count", label: "Count", aggFunc: "COUNT" });
    expect(command.columns[3]).toEqual({ kind: "formula", name: "Avg", label: "Avg", expr: "Total / Count" });
  });

  it("includes detail.by only in aggregate mode", () => {
    expect(toReportRequestCommand(baseValue()).detail).toEqual({ mode: "aggregate", by: "category" });

    const records = baseValue();
    records.detail = { mode: ReportDetail.ModeEnum.Records, by: "category" };
    expect(toReportRequestCommand(records).detail).toEqual({ mode: "records" });
  });

  it("omits the document when title, intro and footer are all empty", () => {
    expect(toReportRequestCommand(baseValue()).document).toBeUndefined();
  });

  it("includes the document when any slot has content", () => {
    const value = baseValue();
    value.document = { title: "T", intro: "", footer: "" };
    expect(toReportRequestCommand(value).document).toEqual({ title: "T", intro: "", footer: "" });
  });
});
