import { ReportColumn, ReportDetail, ReportPeriod } from "../../open-api";
import {
  enabledReportColumns,
  isDimensionColumnDisabled,
  ReportBuilderValue,
  toReportRequestCommand,
  toReportRequestCommandForSave,
} from "./report-command.mapper";

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

  it("excludes a disabled (invalid) dimension column from the request", () => {
    // The reported failing config: aggregate by tag, group by paid_by, but a
    // Category dimension column that reads neither -> it is left out of the spec.
    const value = baseValue();
    value.groupBy = ["paid_by"];
    value.detail = { mode: ReportDetail.ModeEnum.Aggregate, by: "tag" };
    value.columns = [
      { id: "cat", kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
      { id: "cnt", kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
      { id: "tot", kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
    ];
    const command = toReportRequestCommand(value);
    expect(command.columns.map((c) => c.name)).toEqual(["Count", "Total"]);
  });
});

describe("toReportRequestCommandForSave", () => {
  it("keeps every column, including a disabled (invalid) dimension column", () => {
    // Same failing config as the generate-mapper test: aggregate by tag, group by
    // paid_by, plus a Category dimension that reads neither. The save mapper keeps
    // it (so it round-trips and self-heals) where the generate mapper drops it.
    const value = baseValue();
    value.groupBy = ["paid_by"];
    value.detail = { mode: ReportDetail.ModeEnum.Aggregate, by: "tag" };
    value.columns = [
      { id: "cat", kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
      { id: "cnt", kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
      { id: "tot", kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
    ];

    expect(toReportRequestCommandForSave(value).columns.map((c) => c.name)).toEqual(["Category", "Count", "Total"]);
    // The generate mapper still drops the disabled column, proving the two diverge.
    expect(toReportRequestCommand(value).columns.map((c) => c.name)).toEqual(["Count", "Total"]);
  });

  it("matches the generate mapper when no column is disabled", () => {
    const value = baseValue();
    expect(toReportRequestCommandForSave(value)).toEqual(toReportRequestCommand(value));
  });
});

describe("isDimensionColumnDisabled", () => {
  const dim = (field: string): any => ({ id: "x", kind: ReportColumn.KindEnum.Dimension, name: "X", label: "X", field });
  const agg = (): any => ({ id: "y", kind: ReportColumn.KindEnum.Aggregate, name: "Y", label: "Y", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" });
  const Aggregate = ReportDetail.ModeEnum.Aggregate;
  const Records = ReportDetail.ModeEnum.Records;

  it("disables a dimension that is neither the aggregate-by nor a grouping level", () => {
    expect(isDimensionColumnDisabled(dim("category"), Aggregate, "tag", ["paid_by"])).toBe(true);
  });

  it("keeps a dimension that matches the aggregate-by dimension", () => {
    expect(isDimensionColumnDisabled(dim("tag"), Aggregate, "tag", ["paid_by"])).toBe(false);
  });

  it("keeps a dimension that matches a grouping level", () => {
    expect(isDimensionColumnDisabled(dim("paid_by"), Aggregate, "tag", ["paid_by"])).toBe(false);
  });

  it("never disables aggregate/formula columns", () => {
    expect(isDimensionColumnDisabled(agg(), Aggregate, "tag", ["paid_by"])).toBe(false);
  });

  it("never disables anything in records mode", () => {
    expect(isDimensionColumnDisabled(dim("category"), Records, "tag", ["paid_by"])).toBe(false);
  });

  it("enabledReportColumns drops only the disabled dimension columns", () => {
    const value = {
      detail: { mode: Aggregate, by: "tag" },
      groupBy: ["paid_by"],
      columns: [dim("category"), dim("tag"), agg()],
    } as any as ReportBuilderValue;
    expect(enabledReportColumns(value).map((c) => c.field ?? c.name)).toEqual(["tag", "Y"]);
  });
});
