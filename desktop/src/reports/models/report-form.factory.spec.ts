import { Component, provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder } from "@angular/forms";
import { UntilDestroy } from "@ngneat/until-destroy";
import { FilterOperation, ReportColumn, ReportDetail, ReportPeriod, ReportRequestCommand } from "../../open-api";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { ReportBuilderValue, toReportRequestCommand } from "./report-command.mapper";
import { buildReportFormFromCommand } from "./report-form.factory";

// buildReceiptFilterForm wires untilDestroyed subscriptions, so it needs an
// @UntilDestroy()-decorated context — a throwaway host, exactly as the receipt
// filter's own spec does.
@UntilDestroy()
@Component({ selector: "app-noop", template: "", standalone: false })
class NoopComponent {}

describe("buildReportFormFromCommand", () => {
  let host: NoopComponent;
  let fb: FormBuilder;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [NoopComponent],
      providers: [provideZonelessChangeDetection(), FormBuilder],
    }).compileComponents();
    host = TestBed.createComponent(NoopComponent).componentInstance;
    fb = TestBed.inject(FormBuilder);
  });

  // A stored template's filter is always the canonical full-shape value the builder
  // emits, so round-tripping it through buildReceiptFilterForm is a fixpoint. Seed the
  // fixture's filter the same way so the deep-equal below isn't tripped by the sparse-
  // vs-canonical normalization.
  const canonicalFilter = (seed: unknown): unknown => buildReceiptFilterForm(seed, host).getRawValue();

  function roundTrip(command: ReportRequestCommand): ReportRequestCommand {
    const form = buildReportFormFromCommand(fb, host, command);
    return toReportRequestCommand(form.getRawValue() as ReportBuilderValue);
  }

  it("round-trips a preset-period aggregate command with every column kind and a document", () => {
    const command: ReportRequestCommand = {
      name: "Quarterly Spend",
      groupIds: ["1", "2"],
      period: { preset: ReportPeriod.PresetEnum.Qtd },
      filter: canonicalFilter({
        amount: { operation: FilterOperation.GreaterThan, value: 10 },
        name: { operation: FilterOperation.Contains, value: "coffee" },
        categories: { operation: FilterOperation.Contains, value: [2, 3] },
      }),
      groupBy: ["group", "category"],
      detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "category" },
      columns: [
        { kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
        { kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
        // COUNT carries no measure in a stored config; hydration must not re-introduce one.
        { kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
        { kind: ReportColumn.KindEnum.Formula, name: "Avg", label: "Avg", expr: "Total / Count" },
      ],
      subtotals: true,
      grandTotals: false,
      document: { title: "Title", intro: "Intro copy", footer: "Footer copy" },
      formats: [ReportRequestCommand.FormatsEnum.Csv, ReportRequestCommand.FormatsEnum.Pdf],
    };

    expect(roundTrip(command)).toEqual(command);
  });

  it("round-trips a custom-period records command with no document and an amount BETWEEN filter", () => {
    const command: ReportRequestCommand = {
      name: "Custom Range",
      groupIds: ["5"],
      period: { preset: ReportPeriod.PresetEnum.Custom, startDate: "2026-05-01", endDate: "2026-05-31" },
      filter: canonicalFilter({
        amount: { operation: FilterOperation.Between, value: [5, 50] },
        tags: { operation: FilterOperation.Contains, value: [7] },
      }),
      groupBy: [],
      detail: { mode: ReportDetail.ModeEnum.Records },
      columns: [
        { kind: ReportColumn.KindEnum.Dimension, name: "Name", label: "Name", field: "name" },
        { kind: ReportColumn.KindEnum.Dimension, name: "Amount", label: "Amount", field: "amount" },
      ],
      subtotals: false,
      grandTotals: true,
      // document omitted: the forward map dropped it because every slot was empty.
      formats: [ReportRequestCommand.FormatsEnum.Xlsx],
    };

    expect(roundTrip(command)).toEqual(command);
  });

  it("parses custom-period date strings back into Date objects", () => {
    const command: ReportRequestCommand = {
      name: "Dates",
      groupIds: ["1"],
      period: { preset: ReportPeriod.PresetEnum.Custom, startDate: "2026-05-01", endDate: "2026-05-31" },
      filter: canonicalFilter({}),
      groupBy: [],
      detail: { mode: ReportDetail.ModeEnum.Records },
      columns: [{ kind: ReportColumn.KindEnum.Dimension, name: "Name", label: "Name", field: "name" }],
      subtotals: true,
      grandTotals: true,
      formats: [ReportRequestCommand.FormatsEnum.Pdf],
    };

    const form = buildReportFormFromCommand(fb, host, command);
    const start = form.get("period.startDate")!.value as Date;
    expect(start instanceof Date).toBe(true);
    expect(start.getFullYear()).toBe(2026);
    expect(start.getMonth()).toBe(4); // May (0-indexed)
    expect(start.getDate()).toBe(1);
  });

  it("rebuilds the scope, group-by and columns FormArrays from the command", () => {
    const command: ReportRequestCommand = {
      name: "Arrays",
      groupIds: ["3", "4", "9"],
      period: { preset: ReportPeriod.PresetEnum.ThisMonth },
      filter: canonicalFilter({}),
      groupBy: ["group", "status"],
      detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "group" },
      columns: [
        { kind: ReportColumn.KindEnum.Dimension, name: "Group", label: "Group", field: "group" },
        { kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
      ],
      subtotals: true,
      grandTotals: true,
      formats: [ReportRequestCommand.FormatsEnum.Csv],
    };

    const form = buildReportFormFromCommand(fb, host, command);
    expect((form.get("scope") as FormArray).length).toBe(3);
    expect((form.get("groupBy") as FormArray).length).toBe(2);
    expect((form.get("columns") as FormArray).length).toBe(2);
    // Each hydrated column gets a fresh client id for @for tracking.
    expect(form.get("columns.0.id")!.value).toBeTruthy();
    expect(form.get("columns.1.id")!.value).toBeTruthy();
    expect(form.get("columns.0.id")!.value).not.toBe(form.get("columns.1.id")!.value);
  });

  // A saved filter can touch any receipt-filter field, over every editor type: text
  // (name), number BETWEEN (amount), date BETWEEN (date), users array (paidBy), and
  // list arrays (categories, tags, status). This is the surface the "filter doesn't
  // come over" bug lived on, so it must round-trip and land in the controls intact.
  const everyFilterSeed = {
    name: { operation: FilterOperation.Contains, value: "coffee" },
    amount: { operation: FilterOperation.Between, value: [5, 50] },
    date: { operation: FilterOperation.Between, value: ["2026-01-01", "2026-01-31"] },
    paidBy: { operation: FilterOperation.Contains, value: [11, 12] },
    categories: { operation: FilterOperation.Contains, value: [2, 3] },
    tags: { operation: FilterOperation.Contains, value: [7] },
    status: { operation: FilterOperation.Contains, value: ["OPEN"] },
  };

  function commandWithFilter(filterSeed: unknown): ReportRequestCommand {
    return {
      name: "Filtered",
      groupIds: ["1"],
      period: { preset: ReportPeriod.PresetEnum.ThisMonth },
      filter: canonicalFilter(filterSeed),
      groupBy: [],
      detail: { mode: ReportDetail.ModeEnum.Records },
      columns: [{ kind: ReportColumn.KindEnum.Dimension, name: "Name", label: "Name", field: "name" }],
      subtotals: false,
      grandTotals: true,
      formats: [ReportRequestCommand.FormatsEnum.Csv],
    };
  }

  it("round-trips a command carrying every receipt-filter field", () => {
    expect(roundTrip(commandWithFilter(everyFilterSeed))).toEqual(commandWithFilter(everyFilterSeed));
  });

  it("lands each filter field's value and operation in the hydrated form controls", () => {
    const form = buildReportFormFromCommand(fb, host, commandWithFilter(everyFilterSeed));

    expect(form.get("filter.name.value")!.value).toBe("coffee");
    expect(form.get("filter.name.operation")!.value).toBe(FilterOperation.Contains);

    expect(form.get("filter.amount.value")!.value).toEqual([5, 50]);
    expect(form.get("filter.amount.operation")!.value).toBe(FilterOperation.Between);

    expect(form.get("filter.date.value")!.value).toEqual(["2026-01-01", "2026-01-31"]);
    expect(form.get("filter.date.operation")!.value).toBe(FilterOperation.Between);

    expect(form.get("filter.paidBy.value")!.value).toEqual([11, 12]);
    expect(form.get("filter.categories.value")!.value).toEqual([2, 3]);
    expect(form.get("filter.tags.value")!.value).toEqual([7]);
    expect(form.get("filter.status.value")!.value).toEqual(["OPEN"]);
  });

  it("round-trips a value-less date filter (WITHIN_CURRENT_MONTH)", () => {
    const command = commandWithFilter({ date: { operation: FilterOperation.WithinCurrentMonth } });

    expect(roundTrip(command)).toEqual(command);
    const form = buildReportFormFromCommand(fb, host, command);
    expect(form.get("filter.date.operation")!.value).toBe(FilterOperation.WithinCurrentMonth);
  });
});
