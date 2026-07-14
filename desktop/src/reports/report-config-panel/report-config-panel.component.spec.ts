import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { CustomFieldService } from "../../open-api";
import { ReportConfigPanelComponent } from "./report-config-panel.component";

function buildForm(formBuilder: FormBuilder): FormGroup {
  return formBuilder.group({
    scope: formBuilder.array([]),
    groupBy: formBuilder.array([]),
    detail: formBuilder.group({ mode: formBuilder.control("aggregate"), by: formBuilder.control("category") }),
    columns: formBuilder.array([]),
    filter: formBuilder.group({}),
    period: formBuilder.group({
      preset: formBuilder.control("this_month"),
      startDate: formBuilder.control<Date | null>(null),
      endDate: formBuilder.control<Date | null>(null),
    }),
    document: formBuilder.group({ intro: formBuilder.control("") }),
  });
}

describe("ReportConfigPanelComponent", () => {
  let fixture: ComponentFixture<ReportConfigPanelComponent>;
  let component: ReportConfigPanelComponent;
  let form: FormGroup;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReportConfigPanelComponent],
      providers: [
        provideZonelessChangeDetection(),
        FormBuilder,
        { provide: Store, useValue: { selectSignal: () => signal([]) } },
        { provide: MatDialog, useValue: { open: jest.fn() } },
        { provide: CustomFieldService, useValue: { getPagedCustomFields: jest.fn(() => of({ data: [], totalCount: 0 })) } },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    })
      // This suite exercises the component's add-grouping wiring, not its template;
      // an empty template keeps the harness free of the shared pipes/child components.
      .overrideComponent(ReportConfigPanelComponent, { set: { template: "" } })
      .compileComponents();

    form = buildForm(TestBed.inject(FormBuilder));
    fixture = TestBed.createComponent(ReportConfigPanelComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("form", form);
    component.ngOnInit();
  });

  it("adds a grouping level when the add-select control is set, then resets the control", () => {
    const groupBy = form.get("groupBy") as FormArray;
    expect(groupBy.length).toBe(0);

    component.addGroupControl.setValue("paid_by");

    expect(groupBy.length).toBe(1);
    expect(groupBy.at(0).value).toBe("paid_by");
    expect(component.addGroupControl.value).toBeNull();
  });

  it("ignores the blank (reset) option", () => {
    component.addGroupControl.setValue(null);
    expect((form.get("groupBy") as FormArray).length).toBe(0);
  });
});
