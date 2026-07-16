import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormControl } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { of, Subject, throwError } from "rxjs";
import { SnackbarService } from "../../services";
import {
  CustomFieldService,
  ReportColumn,
  ReportDetail,
  ReportPeriod,
  ReportRequestCommand,
  ReportTemplate,
} from "../../open-api";
import { buildColumnGroup } from "../models/report-form.factory";
import { ReportRunnerService } from "../services/report-runner.service";
import { ReportBuilderComponent } from "./report-builder.component";

describe("ReportBuilderComponent", () => {
  let fixture: ComponentFixture<ReportBuilderComponent>;
  let component: ReportBuilderComponent;
  let runner: { preview: jest.Mock; generateAndDownload: jest.Mock; saveTemplate: jest.Mock };
  let snackbar: { success: jest.Mock };
  // Mutable route stub: the builder reads snapshot.data.template in its field
  // initializer, so a test sets this before creating its own component instance.
  const activatedRoute = { snapshot: { data: {} as Record<string, unknown> } };

  beforeEach(async () => {
    runner = {
      preview: jest.fn(() => of({ html: "<p></p>", receiptCount: 0 })),
      generateAndDownload: jest.fn(() => of(new Blob())),
      saveTemplate: jest.fn(() => of({ id: 1 })),
    };
    snackbar = { success: jest.fn() };
    activatedRoute.snapshot.data = {};

    await TestBed.configureTestingModule({
      declarations: [ReportBuilderComponent],
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

  it("saves a template and shows a success snackbar", async () => {
    (component.form.get("scope") as FormArray).push(new FormControl("1"));
    await fixture.whenStable();
    expect(component.canSaveTemplate()).toBe(true);

    component.saveTemplate();

    expect(runner.saveTemplate).toHaveBeenCalled();
    expect(snackbar.success).toHaveBeenCalledWith("Template saved");
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
});
