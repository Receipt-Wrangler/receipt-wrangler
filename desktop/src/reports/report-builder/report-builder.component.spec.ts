import { Component, NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormControl } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { UntilDestroy } from "@ngneat/until-destroy";
import { of, Subject, throwError } from "rxjs";
import { SnackbarService } from "../../services";
import {
  CustomFieldService,
  FilterOperation,
  Permission,
  ReportColumn,
  ReportDetail,
  ReportPeriod,
  ReportRequestCommand,
  ReportTemplate,
} from "../../open-api";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { ReportBuilderValue, toReportRequestCommand } from "../models/report-command.mapper";
import { buildColumnGroup, readStringArray } from "../models/report-form.factory";
import { ReportRunnerService } from "../services/report-runner.service";
import { ReportBuilderComponent } from "./report-builder.component";

// buildReceiptFilterForm wires untilDestroyed subscriptions, so building a canonical
// filter fixture needs an @UntilDestroy()-decorated context — a throwaway host, as the
// factory spec does.
@UntilDestroy()
@Component({ selector: "app-noop", template: "", standalone: false })
class NoopComponent {}

describe("ReportBuilderComponent", () => {
  let fixture: ComponentFixture<ReportBuilderComponent>;
  let component: ReportBuilderComponent;
  let runner: {
    preview: jest.Mock;
    generateAndDownload: jest.Mock;
    saveTemplate: jest.Mock;
    updateTemplate: jest.Mock;
  };
  let snackbar: { success: jest.Mock };
  // Mutable route stub: the builder reads snapshot.data.template in its field
  // initializer, so a test sets this before creating its own component instance.
  const activatedRoute = { snapshot: { data: {} as Record<string, unknown> } };

  beforeEach(async () => {
    runner = {
      preview: jest.fn(() => of({ html: "<p></p>", receiptCount: 0 })),
      generateAndDownload: jest.fn(() => of(new Blob())),
      saveTemplate: jest.fn(() => of({ id: 1 })),
      updateTemplate: jest.fn(() => of({ id: 7 })),
    };
    snackbar = { success: jest.fn() };
    activatedRoute.snapshot.data = {};

    await TestBed.configureTestingModule({
      declarations: [ReportBuilderComponent, NoopComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: ReportRunnerService, useValue: runner },
        { provide: SnackbarService, useValue: snackbar },
        { provide: CustomFieldService, useValue: { getPagedCustomFields: jest.fn(() => of({ data: [], totalCount: 0 })) } },
        { provide: ActivatedRoute, useValue: activatedRoute },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(ReportBuilderComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it("creates with the default columns", () => {
    expect(component).toBeTruthy();
    expect((component.form.get("columns") as FormArray).length).toBe(3);
  });

  it("cannot preview or generate until a scope group is chosen", () => {
    expect(component.canPreview()).toBe(false);
    expect(component.canGenerate()).toBe(false);
  });

  it("becomes previewable once a group is added, and generatable with a format", async () => {
    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    await fixture.whenStable();
    expect(component.canPreview()).toBe(true);
    // A format is selected by default (pdf), so generation is ready too.
    expect(component.canGenerate()).toBe(true);
  });

  it("cannot preview when every column is a disabled dimension", async () => {
    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    // Aggregate by tag with no grouping, but the only column reads category ->
    // it is disabled/excluded, leaving an empty spec.
    component.form.get("detail.by")!.setValue("tag");
    const columns = component.form.get("columns") as FormArray;
    columns.clear();
    columns.push(
      buildColumnGroup(TestBed.inject(FormBuilder), {
        kind: ReportColumn.KindEnum.Dimension,
        name: "Category",
        label: "Category",
        field: "category",
      })
    );
    await fixture.whenStable();
    expect(component.canPreview()).toBe(false);
  });

  it("flags previewError when the debounced preview request fails", () => {
    jest.useFakeTimers();
    try {
      runner.preview.mockReturnValue(throwError(() => new Error("boom")));
      // A fresh component so its construction-time debounce runs under fake timers.
      const local = TestBed.createComponent(ReportBuilderComponent).componentInstance;
      (local.form.get("scope") as FormArray).push(new FormControl("1"));
      // Advance past the 450ms debounce so the preview fires and errors.
      jest.advanceTimersByTime(500);

      expect(runner.preview).toHaveBeenCalled();
      expect(local.previewError()).toBe(true);
      expect(local.previewLoading()).toBe(false);
    } finally {
      jest.useRealTimers();
    }
  });

  it("generates the report and resets the generating flag", async () => {
    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    await fixture.whenStable();
    expect(component.canGenerate()).toBe(true);

    component.generate();

    expect(runner.generateAndDownload).toHaveBeenCalled();
    // generateAndDownload emits synchronously, so finalize clears the flag.
    expect(component.generating()).toBe(false);
  });

  it("saves a new template and shows a success snackbar on the new route", async () => {
    // The blank builder creates (save-as-new is retired): the Save action gates on
    // create, calls saveTemplate (not updateTemplate), and toasts "Template saved".
    expect(component.isEditMode).toBe(false);
    expect(component.saveButtonPermission).toBe(Permission.AppReportsCreate);

    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    await fixture.whenStable();
    expect(component.canSaveTemplate()).toBe(true);

    component.saveTemplate();

    expect(runner.saveTemplate).toHaveBeenCalled();
    expect(runner.updateTemplate).not.toHaveBeenCalled();
    expect(snackbar.success).toHaveBeenCalledWith("Template saved");
  });

  it("updates the opened template in place on the edit route", () => {
    const configuration: ReportRequestCommand = {
      name: "Saved Report",
      groupIds: ["1"],
      period: { preset: ReportPeriod.PresetEnum.ThisMonth },
      filter: {},
      groupBy: [],
      detail: { mode: ReportDetail.ModeEnum.Records },
      columns: [{ kind: ReportColumn.KindEnum.Dimension, name: "Name", label: "Name", field: "name" }],
      subtotals: false,
      grandTotals: false,
      formats: [ReportRequestCommand.FormatsEnum.Csv],
    };
    const template: ReportTemplate = {
      id: 7,
      name: "Saved Report",
      createdAt: "2026-01-01T00:00:00Z",
      configuration,
      configurationVersion: 1,
    };
    activatedRoute.snapshot.data = { template };

    const local = TestBed.createComponent(ReportBuilderComponent).componentInstance;
    // Edit mode gates Save on update, not create.
    expect(local.isEditMode).toBe(true);
    expect(local.saveButtonPermission).toBe(Permission.AppReportsUpdate);
    expect(local.canSaveTemplate()).toBe(true);

    local.saveTemplate();

    // Updates the same id in place — never a create/save-as-new.
    expect(runner.updateTemplate).toHaveBeenCalledWith(7, expect.objectContaining({ name: "Saved Report" }));
    expect(runner.saveTemplate).not.toHaveBeenCalled();
    expect(snackbar.success).toHaveBeenCalledWith("Template updated");
  });

  it("cannot save a template without a name", async () => {
    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    component.form.get("name")!.setValue("   ");
    await fixture.whenStable();
    expect(component.canSaveTemplate()).toBe(false);

    component.saveTemplate();
    expect(runner.saveTemplate).not.toHaveBeenCalled();
  });

  it("guards against concurrent saves while one is in flight", async () => {
    const pending = new Subject<{ id: number }>();
    runner.saveTemplate.mockReturnValue(pending.asObservable());

    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    await fixture.whenStable();
    expect(component.canSaveTemplate()).toBe(true);

    component.saveTemplate();
    expect(component.saving()).toBe(true);

    // A second call (e.g. a double-click) must not fire another POST.
    component.saveTemplate();
    expect(runner.saveTemplate).toHaveBeenCalledTimes(1);

    // Completing the request clears the in-flight flag so later saves work again.
    pending.next({ id: 1 });
    pending.complete();
    expect(component.saving()).toBe(false);
  });

  it("starts from a blank form when no template is resolved", () => {
    expect(component.loadedTemplateName).toBeNull();
    expect(component.form.get("name")!.value).toBe("Untitled Report");
    expect((component.form.get("columns") as FormArray).length).toBe(3);
  });

  it("builds the form from a resolved template on the edit route", () => {
    const configuration: ReportRequestCommand = {
      name: "Saved Report",
      groupIds: ["1", "2"],
      period: { preset: ReportPeriod.PresetEnum.ThisMonth },
      filter: {},
      groupBy: [],
      detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "category" },
      columns: [
        { kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
        { kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
      ],
      subtotals: true,
      grandTotals: true,
      formats: [ReportRequestCommand.FormatsEnum.Pdf],
    };
    const template: ReportTemplate = {
      id: 7,
      name: "Saved Report",
      createdAt: "2026-01-01T00:00:00Z",
      configuration,
      configurationVersion: 1,
    };
    activatedRoute.snapshot.data = { template };

    const local = TestBed.createComponent(ReportBuilderComponent).componentInstance;

    expect(local.loadedTemplateName).toBe("Saved Report");
    expect(local.form.get("name")!.value).toBe("Saved Report");
    expect((local.form.get("scope") as FormArray).length).toBe(2);
    expect((local.form.get("columns") as FormArray).length).toBe(2);
  });

  it("maps every field of a fully populated template into the builder form on the edit route", () => {
    // A canonical filter touching every builder-visible filter field, built the same
    // way the backend stores it (buildReceiptFilterForm is a fixpoint on its own output).
    const host = TestBed.createComponent(NoopComponent).componentInstance;
    const filter = buildReceiptFilterForm(
      {
        name: { operation: FilterOperation.Contains, value: "coffee" },
        amount: { operation: FilterOperation.Between, value: [5, 50] },
        date: { operation: FilterOperation.WithinCurrentMonth },
        paidBy: { operation: FilterOperation.Contains, value: [11] },
        categories: { operation: FilterOperation.Contains, value: [2, 3] },
        tags: { operation: FilterOperation.Contains, value: [7] },
        status: { operation: FilterOperation.Contains, value: ["OPEN"] },
      },
      host
    ).getRawValue();

    // A fully-populated, internally-consistent aggregate report: every dimension
    // column's field is detail.by or in groupBy (so none is disabled and dropped on
    // re-map), grandTotals is false (proves the ?? default preserves a stored false).
    const configuration: ReportRequestCommand = {
      name: "Full Report",
      groupIds: ["3", "7"],
      period: { preset: ReportPeriod.PresetEnum.Custom, startDate: "2026-03-01", endDate: "2026-03-31" },
      filter,
      groupBy: ["group", "category"],
      detail: { mode: ReportDetail.ModeEnum.Aggregate, by: "category" },
      columns: [
        { kind: ReportColumn.KindEnum.Dimension, name: "Category", label: "Category", field: "category" },
        { kind: ReportColumn.KindEnum.Dimension, name: "Group", label: "Group", field: "group" },
        { kind: ReportColumn.KindEnum.Aggregate, name: "Total", label: "Total", aggFunc: ReportColumn.AggFuncEnum.Sum, measure: "amount" },
        { kind: ReportColumn.KindEnum.Aggregate, name: "Count", label: "Count", aggFunc: ReportColumn.AggFuncEnum.Count },
        { kind: ReportColumn.KindEnum.Formula, name: "Avg", label: "Avg", expr: "Total / Count" },
      ],
      subtotals: true,
      grandTotals: false,
      document: { title: "Q1 Spend", intro: "Period Covering: {{period}}", footer: "Generated {{generatedAt}}" },
      formats: [
        ReportRequestCommand.FormatsEnum.Csv,
        ReportRequestCommand.FormatsEnum.Xlsx,
        ReportRequestCommand.FormatsEnum.Pdf,
      ],
    };
    const template: ReportTemplate = {
      id: 12,
      name: "Full Report",
      createdAt: "2026-01-01T00:00:00Z",
      configuration,
      configurationVersion: 1,
    };
    activatedRoute.snapshot.data = { template };

    const form = TestBed.createComponent(ReportBuilderComponent).componentInstance.form;

    // Name + scope (values, not just count).
    expect(form.get("name")!.value).toBe("Full Report");
    expect(readStringArray(form.get("scope") as FormArray)).toEqual(["3", "7"]);

    // Custom period parses back into Date objects (local midnight, no TZ drift).
    const start = form.get("period.startDate")!.value as Date;
    const end = form.get("period.endDate")!.value as Date;
    expect(form.get("period.preset")!.value).toBe(ReportPeriod.PresetEnum.Custom);
    expect([start.getFullYear(), start.getMonth(), start.getDate()]).toEqual([2026, 2, 1]);
    expect([end.getFullYear(), end.getMonth(), end.getDate()]).toEqual([2026, 2, 31]);

    // Every builder filter field's value + operation.
    expect(form.get("filter.name.value")!.value).toBe("coffee");
    expect(form.get("filter.name.operation")!.value).toBe(FilterOperation.Contains);
    expect(form.get("filter.amount.value")!.value).toEqual([5, 50]);
    expect(form.get("filter.amount.operation")!.value).toBe(FilterOperation.Between);
    expect(form.get("filter.date.operation")!.value).toBe(FilterOperation.WithinCurrentMonth);
    expect(form.get("filter.paidBy.value")!.value).toEqual([11]);
    expect(form.get("filter.categories.value")!.value).toEqual([2, 3]);
    expect(form.get("filter.categories.operation")!.value).toBe(FilterOperation.Contains);
    expect(form.get("filter.tags.value")!.value).toEqual([7]);
    expect(form.get("filter.status.value")!.value).toEqual(["OPEN"]);

    // Group-by (values), detail mode + by.
    expect(readStringArray(form.get("groupBy") as FormArray)).toEqual(["group", "category"]);
    expect(form.get("detail.mode")!.value).toBe(ReportDetail.ModeEnum.Aggregate);
    expect(form.get("detail.by")!.value).toBe("category");

    // Every column, by value — kind/name/label + the kind-specific field.
    const columns = form.get("columns") as FormArray;
    expect(columns.length).toBe(5);
    const col = (i: number, key: string) => columns.at(i).get(key)!.value;
    expect([col(0, "kind"), col(0, "name"), col(0, "label"), col(0, "field")]).toEqual([
      ReportColumn.KindEnum.Dimension, "Category", "Category", "category",
    ]);
    expect([col(1, "kind"), col(1, "name"), col(1, "field")]).toEqual([
      ReportColumn.KindEnum.Dimension, "Group", "group",
    ]);
    expect([col(2, "kind"), col(2, "name"), col(2, "aggFunc"), col(2, "measure")]).toEqual([
      ReportColumn.KindEnum.Aggregate, "Total", ReportColumn.AggFuncEnum.Sum, "amount",
    ]);
    expect([col(3, "kind"), col(3, "name"), col(3, "aggFunc")]).toEqual([
      ReportColumn.KindEnum.Aggregate, "Count", ReportColumn.AggFuncEnum.Count,
    ]);
    expect([col(4, "kind"), col(4, "name"), col(4, "expr")]).toEqual([
      ReportColumn.KindEnum.Formula, "Avg", "Total / Count",
    ]);
    // Each hydrated column still gets a fresh client id for @for tracking.
    expect(col(0, "id")).toBeTruthy();

    // Totals (false must survive), document slots, and format booleans.
    expect(form.get("subtotals")!.value).toBe(true);
    expect(form.get("grandTotals")!.value).toBe(false);
    expect(form.get("document.title")!.value).toBe("Q1 Spend");
    expect(form.get("document.intro")!.value).toBe("Period Covering: {{period}}");
    expect(form.get("document.footer")!.value).toBe("Generated {{generatedAt}}");
    expect(form.get("formats.csv")!.value).toBe(true);
    expect(form.get("formats.xlsx")!.value).toBe(true);
    expect(form.get("formats.pdf")!.value).toBe(true);

    // Single catch-all: mapping back yields exactly the stored configuration.
    expect(toReportRequestCommand(form.getRawValue() as ReportBuilderValue)).toEqual(configuration);
  });
});
