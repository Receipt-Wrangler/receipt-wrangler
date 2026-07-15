import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormControl } from "@angular/forms";
import { of, Subject, throwError } from "rxjs";
import { SnackbarService } from "../../services";
import { CustomFieldService, ReportColumn } from "../../open-api";
import { buildColumnGroup } from "../models/report-form.factory";
import { ReportRunnerService } from "../services/report-runner.service";
import { ReportBuilderComponent } from "./report-builder.component";

describe("ReportBuilderComponent", () => {
  let fixture: ComponentFixture<ReportBuilderComponent>;
  let component: ReportBuilderComponent;
  let runner: { preview: jest.Mock; generateAndDownload: jest.Mock; saveTemplate: jest.Mock };
  let snackbar: { success: jest.Mock };

  beforeEach(async () => {
    runner = {
      preview: jest.fn(() => of({ html: "<p></p>", receiptCount: 0 })),
      generateAndDownload: jest.fn(() => of(new Blob())),
      saveTemplate: jest.fn(() => of({ id: 1 })),
    };
    snackbar = { success: jest.fn() };

    await TestBed.configureTestingModule({
      declarations: [ReportBuilderComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: ReportRunnerService, useValue: runner },
        { provide: SnackbarService, useValue: snackbar },
        { provide: CustomFieldService, useValue: { getPagedCustomFields: jest.fn(() => of({ data: [], totalCount: 0 })) } },
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
});
