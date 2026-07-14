import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder, FormGroup } from "@angular/forms";
import { Store } from "@ngxs/store";
import { ReportFiltersComponent } from "./report-filters.component";

describe("ReportFiltersComponent", () => {
  let fixture: ComponentFixture<ReportFiltersComponent>;
  let component: ReportFiltersComponent;
  let form: FormGroup;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReportFiltersComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: Store, useValue: { selectSignal: () => signal([]), selectSnapshot: () => [] } },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    const formBuilder = TestBed.inject(FormBuilder);
    form = formBuilder.group({ filter: formBuilder.group({}) });
    fixture = TestBed.createComponent(ReportFiltersComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("form", form);
  });

  it("adds a filter when the add-select control is set, then resets the control", () => {
    expect(component.activeFields().length).toBe(0);

    component.addFilterControl.setValue("name");

    expect(component.activeFields().map((def) => def.field)).toContain("name");
    expect(component.addFilterControl.value).toBeNull();
  });

  it("ignores the blank (reset) option", () => {
    component.addFilterControl.setValue(null);
    expect(component.activeFields().length).toBe(0);
  });
});
