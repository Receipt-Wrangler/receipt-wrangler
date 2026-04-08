import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { Component, CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { SetSystemTaskFilter } from "src/store/system-task-table.state.actions";
import { defaultSystemTaskFilter } from "src/store/system-task-table.state";
import { InputModule } from "../../input";
import { FilterOperation, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { StoreModule } from "../../store/store.module";
import { applyFormCommand } from "../../utils/index";
import { buildSystemTaskFilterForm } from "../../utils/system-task-filter";
import { OperationsPipe } from "../receipt-filter/operations.pipe";
import { SystemTaskFilterComponent } from "./system-task-filter.component";

@UntilDestroy()
@Component({
  selector: "app-noop",
  template: "",
  standalone: false
})
class NoopComponent {}

describe("SystemTaskFilterComponent", () => {
  let component: SystemTaskFilterComponent;
  let fixture: ComponentFixture<SystemTaskFilterComponent>;
  let store: Store;

  const filledFilter = {
    type: {
      operation: FilterOperation.Contains,
      value: [SystemTaskType.MagicFill],
    },
    status: {
      operation: FilterOperation.Contains,
      value: [SystemTaskStatus.Failed],
    },
    ranByUserId: {
      operation: FilterOperation.Contains,
      value: [1],
    },
    startedAt: {
      operation: FilterOperation.GreaterThan,
      value: "2023-01-06",
    },
    endedAt: {
      operation: FilterOperation.LessThan,
      value: "2023-12-31",
    },
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SystemTaskFilterComponent, OperationsPipe],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        PipesModule,
        InputModule,
        MatDialogModule,
        StoreModule,
        NoopAnimationsModule,
        ReactiveFormsModule,
      ],
      providers: [
        {
          provide: MatDialogRef,
          useValue: {
            close: (value: any) => {},
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(SystemTaskFilterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form with no default initial data", () => {
    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;

    component.parentForm = buildSystemTaskFilterForm({}, noopComponent);
    component.ngOnInit();

    expect(component.parentForm.value).toEqual(defaultSystemTaskFilter);
  });

  it("should init form with initial data", () => {
    store.reset({
      systemTaskTable: {
        filter: filledFilter,
      },
    });

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;

    component.parentForm = buildSystemTaskFilterForm(filledFilter, noopComponent);
    component.ngOnInit();

    expect(component.parentForm.value).toEqual(filledFilter);
  });

  it("should reset form", () => {
    store.reset({
      systemTaskTable: {
        filter: filledFilter,
      },
    });

    component.formCommand.subscribe((formCommand) => {
      applyFormCommand(component.parentForm, formCommand);
    });

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;

    component.parentForm = buildSystemTaskFilterForm(filledFilter, noopComponent);
    component.ngOnInit();

    expect(component.parentForm.value).toEqual(filledFilter);

    component.resetFilter();
    expect(component.parentForm.value).toEqual(defaultSystemTaskFilter);
  });

  it("should set form in state and close dialog", () => {
    const dialogRefSpy = jest.spyOn(
      TestBed.inject(MatDialogRef<SystemTaskFilterComponent>),
      "close"
    );
    const storeRefSpy = jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

    component.submitButtonClicked();

    expect(storeRefSpy).toHaveBeenCalledWith(
      new SetSystemTaskFilter(component.parentForm.value)
    );
    expect(dialogRefSpy).toHaveBeenCalledWith(true);
  });

  it("should close dialog on cancel", () => {
    const dialogRefSpy = jest.spyOn(
      TestBed.inject(MatDialogRef<SystemTaskFilterComponent>),
      "close"
    );
    component.cancelButtonClicked();

    expect(dialogRefSpy).toHaveBeenCalledWith(false);
  });

  it("should build type options excluding child-only types", () => {
    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;
    component.parentForm = buildSystemTaskFilterForm({}, noopComponent);
    component.ngOnInit();

    const typeValues = component.systemTaskTypeOptions.map(o => o.value);

    expect(typeValues).toContain(SystemTaskType.MagicFill);
    expect(typeValues).toContain(SystemTaskType.QuickScan);
    expect(typeValues).toContain(SystemTaskType.EmailUpload);
    expect(typeValues).toContain(SystemTaskType.ApiKeyDeleted);
    expect(typeValues).not.toContain(SystemTaskType.OcrProcessing);
    expect(typeValues).not.toContain(SystemTaskType.ChatCompletion);
    expect(typeValues).not.toContain(SystemTaskType.ReceiptUploaded);
  });

  it("should have correct status options using SystemTaskStatus enum", () => {
    expect(component.systemTaskStatusOptions).toEqual([
      { value: SystemTaskStatus.Succeeded, displayValue: "Succeeded" },
      { value: SystemTaskStatus.Failed, displayValue: "Failed" },
    ]);
  });
});
