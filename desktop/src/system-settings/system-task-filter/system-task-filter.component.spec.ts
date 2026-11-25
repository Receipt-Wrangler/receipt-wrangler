import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, DestroyRef } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { SYSTEM_RAN_BY_OPTION_ID } from "src/constants";
import { buildDefaultSystemTaskFilter } from "src/store/system-task-table.state";
import { SetSystemTaskFilter } from "src/store/system-task-table.state.actions";
import { FilterOperation } from "../../open-api";
import { provideTestStore } from "../../store/test-store";
import { applyFormCommand } from "../../utils/index";
import { buildSystemTaskFilterForm } from "../../utils/system-task-filter";
import { SystemTaskFilterComponent } from "./system-task-filter.component";
import { provideRouter } from "@angular/router";

describe("SystemTaskFilterComponent", () => {
  let component: SystemTaskFilterComponent;
  let fixture: ComponentFixture<SystemTaskFilterComponent>;
  let store: Store;
  let destroyRef: DestroyRef;

  const filledFilter = {
    type: {
      operation: FilterOperation.Contains,
      value: ["QUICK_SCAN"],
    },
    ranBy: {
      operation: FilterOperation.Contains,
      value: [SYSTEM_RAN_BY_OPTION_ID, 4],
    },
    startedAt: {
      operation: FilterOperation.Equals,
      value: "2026-03-11T00:00:00Z",
    },
    endedAt: {
      operation: FilterOperation.LessThan,
      value: "2026-03-12T00:00:00Z",
    },
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        MatDialogModule,
        ReactiveFormsModule,
        SystemTaskFilterComponent,
      ],
      providers: [
        provideTestStore(),
        provideRouter([]),
        { provide: MatDialogRef, useValue: { close: () => {} } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    store = TestBed.inject(Store);
    destroyRef = TestBed.inject(DestroyRef);

    fixture = TestBed.createComponent(SystemTaskFilterComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("parentForm", buildSystemTaskFilterForm({}, destroyRef));
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("builds an empty form matching the default filter", () => {
    expect(component.parentForm().value).toEqual(buildDefaultSystemTaskFilter());
  });

  it("round-trips a seeded filter", () => {
    const localFixture = TestBed.createComponent(SystemTaskFilterComponent);
    localFixture.componentRef.setInput("parentForm", buildSystemTaskFilterForm(filledFilter, destroyRef));
    localFixture.detectChanges();

    expect(localFixture.componentInstance.parentForm().value).toEqual(filledFilter);
  });

  // The form can outlive the filter component, so the auto-selection listener
  // has to end with the component, not the form.
  it("stops auto-selecting an operation once the component is destroyed", () => {
    const form = buildSystemTaskFilterForm({}, destroyRef);
    const localFixture = TestBed.createComponent(SystemTaskFilterComponent);
    localFixture.componentRef.setInput("parentForm", form);
    localFixture.detectChanges();

    form.get("startedAt.value")!.setValue(new Date(2026, 10, 1));
    expect(form.get("startedAt.operation")!.value).toBeTruthy();

    localFixture.destroy();
    form.get("startedAt.operation")!.setValue(null);
    form.get("startedAt.value")!.setValue(new Date(2026, 10, 2));

    expect(form.get("startedAt.operation")!.value).toBeNull();
  });

  it("pins System above the real users in the Ran By options", () => {
    expect(component.ranByOptions()[0]).toEqual({
      id: SYSTEM_RAN_BY_OPTION_ID,
      displayName: "System",
    });
  });

  it("renders one row per filter field", () => {
    expect(fixture.nativeElement.querySelectorAll("app-filter-field").length).toBe(4);
  });

  // reset() leaves a FormArray populated, so the two list fields need an
  // explicit clear command on top of it.
  it("resets every field, clearing the list ones", () => {
    fixture.componentRef.setInput("parentForm", buildSystemTaskFilterForm(filledFilter, destroyRef));
    component.formCommand.subscribe((formCommand) =>
      applyFormCommand(component.parentForm(), formCommand)
    );

    component.resetFilter();

    expect(component.parentForm().value).toEqual(buildDefaultSystemTaskFilter());
  });

  it("dispatches the filter and closes with true on submit", () => {
    const dispatchSpy = vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
    const closeSpy = vi.spyOn(TestBed.inject(MatDialogRef<SystemTaskFilterComponent>), "close");

    fixture.componentRef.setInput("parentForm", buildSystemTaskFilterForm(filledFilter, destroyRef));
    component.submitButtonClicked();

    expect(dispatchSpy).toHaveBeenCalledWith(new SetSystemTaskFilter(filledFilter as any));
    expect(closeSpy).toHaveBeenCalledWith(true);
  });

  it("closes with false on cancel", () => {
    const closeSpy = vi.spyOn(TestBed.inject(MatDialogRef<SystemTaskFilterComponent>), "close");

    component.cancelButtonClicked();

    expect(closeSpy).toHaveBeenCalledWith(false);
  });
});
