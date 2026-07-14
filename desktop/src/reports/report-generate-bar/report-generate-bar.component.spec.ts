import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder, FormGroup } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { provideRouter } from "@angular/router";
import { ButtonModule } from "../../button/button.module";
import { ReportGenerateBarComponent } from "./report-generate-bar.component";

describe("ReportGenerateBarComponent", () => {
  let fixture: ComponentFixture<ReportGenerateBarComponent>;
  let component: ReportGenerateBarComponent;
  let form: FormGroup;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReportGenerateBarComponent],
      // Import the real ButtonModule so the Generate app-button renders a real
      // <button> whose disabled state can be asserted from the DOM.
      imports: [ButtonModule, NoopAnimationsModule],
      providers: [provideZonelessChangeDetection(), provideRouter([])],
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

  it("renders the format chips and the Generate button's disabled state", () => {
    fixture.detectChanges();
    const el = fixture.nativeElement as HTMLElement;

    // The format chips render, reflecting the selected (xlsx) format.
    expect(el.querySelector('[data-testid="report-format-xlsx"]')!.classList).toContain("gen-bar__chip--on");
    expect(el.querySelector('[data-testid="report-format-csv"]')!.classList).not.toContain("gen-bar__chip--on");
    expect(el.querySelector(".gen-bar__summary")!.textContent).toContain("Single XLSX file");

    // canGenerate=true → the Generate button is enabled.
    const button = () => el.querySelector('[data-testid="report-generate"] button') as HTMLButtonElement;
    expect(button().disabled).toBe(false);

    // Flipping the inputs disables it (no format-generatable / mid-generation).
    fixture.componentRef.setInput("canGenerate", false);
    fixture.detectChanges();
    expect(button().disabled).toBe(true);

    fixture.componentRef.setInput("canGenerate", true);
    fixture.componentRef.setInput("generating", true);
    fixture.detectChanges();
    expect(button().disabled).toBe(true);
  });
});
