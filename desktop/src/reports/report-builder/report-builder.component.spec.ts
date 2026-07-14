import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormControl } from "@angular/forms";
import { of } from "rxjs";
import { CustomFieldService, ReportColumn, ReportService } from "../../open-api";
import { buildColumnGroup } from "../models/report-form.factory";
import { ReportBuilderComponent } from "./report-builder.component";

describe("ReportBuilderComponent", () => {
  let fixture: ComponentFixture<ReportBuilderComponent>;
  let component: ReportBuilderComponent;
  let reportService: { previewReport: jest.Mock; generateReport: jest.Mock };

  beforeEach(async () => {
    reportService = {
      previewReport: jest.fn(() => of({ html: "<p></p>", receiptCount: 0 })),
      generateReport: jest.fn(() => of(new Blob())),
    };

    await TestBed.configureTestingModule({
      declarations: [ReportBuilderComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: ReportService, useValue: reportService },
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
});
