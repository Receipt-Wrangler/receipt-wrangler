import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder, FormGroup } from "@angular/forms";
import { ReportGenerateBarComponent } from "./report-generate-bar.component";

describe("ReportGenerateBarComponent", () => {
  let fixture: ComponentFixture<ReportGenerateBarComponent>;
  let component: ReportGenerateBarComponent;
  let form: FormGroup;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReportGenerateBarComponent],
      providers: [provideZonelessChangeDetection()],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    const formBuilder = new FormBuilder();
    form = formBuilder.group({
      formats: formBuilder.group({ csv: false, xlsx: true, pdf: false }),
    });

    fixture = TestBed.createComponent(ReportGenerateBarComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("form", form);
    fixture.componentRef.setInput("generating", false);
    fixture.componentRef.setInput("canGenerate", true);
    await fixture.whenStable();
  });

  it("reflects and toggles selected formats", () => {
    expect(component.isSelected("xlsx")).toBe(true);
    expect(component.isSelected("csv")).toBe(false);
    component.toggle("csv");
    expect(component.isSelected("csv")).toBe(true);
  });

  it("summarizes the selection", () => {
    expect(component.formatSummary()).toBe("Single XLSX file");
    component.toggle("pdf");
    expect(component.formatSummary()).toBe("XLSX + PDF → zipped");
    component.toggle("xlsx");
    component.toggle("pdf");
    expect(component.formatSummary()).toBe("Pick at least one format");
  });

  it("emits generate on request", () => {
    const spy = jest.fn();
    component.generate.subscribe(spy);
    component.onGenerate();
    expect(spy).toHaveBeenCalledTimes(1);
  });
});
