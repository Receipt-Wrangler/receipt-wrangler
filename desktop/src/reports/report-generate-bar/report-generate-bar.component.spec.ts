import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
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
      // The real ButtonModule renders the Save/Generate app-buttons as real <button>s
      // whose rendered/disabled state the DOM assertions below can inspect. The bar no
      // longer gates on permissions itself — the parent passes resolved booleans.
      imports: [ButtonModule, NoopAnimationsModule],
      providers: [
        provideZonelessChangeDetection(),
        provideRouter([]),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
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
    // The parent resolves the permission gates into these booleans (honoring the base +
    // "*All" variants); grant both so the Save/Generate buttons render for the DOM tests.
    fixture.componentRef.setInput("saveTemplateAllowed", true);
    fixture.componentRef.setInput("generateAllowed", true);
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

  it("emits saveTemplate on request", () => {
    const spy = jest.fn();
    component.saveTemplate.subscribe(spy);
    component.onSaveTemplate();
    expect(spy).toHaveBeenCalledTimes(1);
  });

  it("labels the Save action by mode and in-flight state", async () => {
    // New route: create label.
    expect(component.saveTemplateText()).toBe("Save Template");
    fixture.componentRef.setInput("saving", true);
    await fixture.whenStable();
    expect(component.saveTemplateText()).toBe("Saving…");

    // Edit route: update label.
    fixture.componentRef.setInput("saving", false);
    fixture.componentRef.setInput("isEditMode", true);
    await fixture.whenStable();
    expect(component.saveTemplateText()).toBe("Update Template");
    fixture.componentRef.setInput("saving", true);
    await fixture.whenStable();
    expect(component.saveTemplateText()).toBe("Updating…");
  });

  it("gates the Save and Generate actions on the resolved permission booleans", async () => {
    const el = fixture.nativeElement as HTMLElement;
    const saveButton = () => el.querySelector('[data-testid="report-save-template"]');
    const generateButton = () => el.querySelector('[data-testid="report-generate"]');

    // Both allowed (granted in beforeEach) → both actions render. The parent resolves
    // these booleans honoring the base + "*All" permission variants.
    expect(saveButton()).not.toBeNull();
    expect(generateButton()).not.toBeNull();

    // A read-only holder gets neither boolean → both actions are hidden.
    fixture.componentRef.setInput("saveTemplateAllowed", false);
    fixture.componentRef.setInput("generateAllowed", false);
    await fixture.whenStable();
    expect(saveButton()).toBeNull();
    expect(generateButton()).toBeNull();

    // Each gate is independent: Generate can be allowed while Save is not.
    fixture.componentRef.setInput("generateAllowed", true);
    await fixture.whenStable();
    expect(saveButton()).toBeNull();
    expect(generateButton()).not.toBeNull();
  });

  it("renders the format chips and the Generate button's disabled state", async () => {
    await fixture.whenStable();
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
    await fixture.whenStable();
    expect(button().disabled).toBe(true);

    fixture.componentRef.setInput("canGenerate", true);
    fixture.componentRef.setInput("generating", true);
    await fixture.whenStable();
    expect(button().disabled).toBe(true);
  });
});
