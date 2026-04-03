import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { Component, CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef, } from "@angular/material/dialog";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { SetReceiptFilter } from "src/store/receipt-table.actions";
import { defaultReceiptFilter, } from "src/store/receipt-table.state";
import { InputModule } from "../../input";
import { CategoryService, CustomFieldService, CustomFieldType, FilterOperation, ReceiptStatus, TagService } from "../../open-api";
import { StoreModule } from "../../store/store.module";
import { applyFormCommand } from "../../utils/index";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { OperationsPipe } from "./operations.pipe";
import { ReceiptFilterComponent } from "./receipt-filter.component";

@UntilDestroy()
@Component({
  selector: "app-noop",
  template: "",
  standalone: false
})
class NoopComponent {}

describe("ReceiptFilterComponent", () => {
  let component: ReceiptFilterComponent;
  let fixture: ComponentFixture<ReceiptFilterComponent>;
  let store: Store;

  const filledFilter = {
    date: {
      operation: FilterOperation.Equals,
      value: "2023-01-06",
    },
    name: {
      operation: FilterOperation.Equals,
      value: "hello world",
    },
    amount: {
      operation: FilterOperation.GreaterThan,
      value: 12.05,
    },
    paidBy: {
      operation: FilterOperation.Contains,
      value: [1],
    },
    categories: {
      operation: FilterOperation.Contains,
      value: [2],
    },
    tags: {
      operation: FilterOperation.Contains,
      value: [3, 4],
    },
    status: {
      operation: FilterOperation.Contains,
      value: [ReceiptStatus.Open],
    },
    resolvedDate: {
      operation: FilterOperation.GreaterThan,
      value: "2023-01-06",
    },
    createdAt: {
      operation: FilterOperation.GreaterThan,
      value: "2023-01-06",
    },
    customFields: [],
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ReceiptFilterComponent, OperationsPipe],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [PipesModule,
        InputModule,
        MatDialogModule,
        StoreModule,
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule],
      providers: [
        CategoryService,
        CustomFieldService,
        TagService,
        {
          provide: MatDialogRef,
          useValue: {
            close: (value: any) => { },
          },
        },
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            categories: [],
            tags: [],
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(ReceiptFilterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form with no default initial data", () => {
    jest.spyOn(TestBed.inject(CategoryService), "getAllCategories").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(TagService), "getAllTags").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(CustomFieldService), "getAllCustomFields").mockReturnValue(
      of([]) as any
    );

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;

    component.parentForm = buildReceiptFilterForm({}, noopComponent);
    component.ngOnInit();

    expect(component.parentForm.value).toEqual(defaultReceiptFilter);
  });

  it("should init form with initial data", () => {
    jest.spyOn(TestBed.inject(CategoryService), "getAllCategories").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(TagService), "getAllTags").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(CustomFieldService), "getAllCustomFields").mockReturnValue(
      of([]) as any
    );
    store.reset({
      receiptTable: {
        filter: filledFilter,
      },
    });

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;

    component.parentForm = buildReceiptFilterForm(filledFilter, noopComponent);
    component.ngOnInit();

    expect(component.parentForm.value).toEqual(filledFilter);
  });

  it("should reset form", () => {
    jest.spyOn(TestBed.inject(CategoryService), "getAllCategories").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(TagService), "getAllTags").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(CustomFieldService), "getAllCustomFields").mockReturnValue(
      of([]) as any
    );
    store.reset({
      receiptTable: {
        filter: filledFilter,
      },
    });

    component.formCommand.subscribe((formCommand) => {
      applyFormCommand(component.parentForm, formCommand);
    });

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;

    component.parentForm = buildReceiptFilterForm(filledFilter, noopComponent);
    component.ngOnInit();

    expect(component.parentForm.value).toEqual(filledFilter);

    component.resetFilter();
    expect(component.parentForm.value).toEqual(defaultReceiptFilter);
  });

  it("should set form in state and close dialog", () => {
    const dialogRefSpy = jest.spyOn(
      TestBed.inject(MatDialogRef<ReceiptFilterComponent>),
      "close"
    );
    const storeRefSpy = jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

    component.submitButtonClicked();

    expect(storeRefSpy).toHaveBeenCalledWith(
      new SetReceiptFilter(component.parentForm.value)
    );
    expect(dialogRefSpy).toHaveBeenCalledWith(true);
  });

  it("should close dialog on cancel", () => {
    const dialogRefSpy = jest.spyOn(
      TestBed.inject(MatDialogRef<ReceiptFilterComponent>),
      "close"
    );
    component.cancelButtonClicked();

    expect(dialogRefSpy).toHaveBeenCalledWith(false);
  });

  it("should add and remove custom field filters", () => {
    jest.spyOn(TestBed.inject(CategoryService), "getAllCategories").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(TagService), "getAllTags").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(CustomFieldService), "getAllCustomFields").mockReturnValue(
      of([
        { id: 1, name: "Text Field", type: CustomFieldType.Text, options: [] },
        { id: 2, name: "Date Field", type: CustomFieldType.Date, options: [] },
      ]) as any
    );

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;
    component.parentForm = buildReceiptFilterForm({}, noopComponent);
    component.ngOnInit();

    expect(component.getCustomFieldsArray().length).toBe(0);

    component.addCustomFieldFilter();
    expect(component.getCustomFieldsArray().length).toBe(1);

    component.addCustomFieldFilter();
    expect(component.getCustomFieldsArray().length).toBe(2);

    component.removeCustomFieldFilter(0);
    expect(component.getCustomFieldsArray().length).toBe(1);
  });

  it("should get available custom fields excluding used ones", () => {
    jest.spyOn(TestBed.inject(CategoryService), "getAllCategories").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(TagService), "getAllTags").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(CustomFieldService), "getAllCustomFields").mockReturnValue(
      of([
        { id: 1, name: "Field 1", type: CustomFieldType.Text, options: [] },
        { id: 2, name: "Field 2", type: CustomFieldType.Date, options: [] },
      ]) as any
    );

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;
    component.parentForm = buildReceiptFilterForm({}, noopComponent);
    component.ngOnInit();

    component.addCustomFieldFilter();
    component.getCustomFieldsArray().at(0).get("customFieldId")?.setValue(1);

    const available = component.getAvailableCustomFields(1);
    expect(available.length).toBe(1);
    expect(available[0].id).toBe(2);
  });

  it("should map custom field types to filter types", () => {
    expect(component.getFilterType({ type: CustomFieldType.Text } as any)).toBe("text");
    expect(component.getFilterType({ type: CustomFieldType.Date } as any)).toBe("date");
    expect(component.getFilterType({ type: CustomFieldType.Currency } as any)).toBe("currency");
    expect(component.getFilterType({ type: CustomFieldType.Boolean } as any)).toBe("boolean");
    expect(component.getFilterType({ type: CustomFieldType.Select } as any)).toBe("list");
  });

  it("should strip incomplete custom field filters on submit", () => {
    const storeRefSpy = jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;
    component.parentForm = buildReceiptFilterForm({}, noopComponent);

    component.addCustomFieldFilter();
    // Leave the filter row incomplete (no customFieldId)

    component.submitButtonClicked();

    const dispatchedFilter = (storeRefSpy.mock.calls[0][0] as SetReceiptFilter).data;
    expect((dispatchedFilter as any).customFields.length).toBe(0);
  });

  it("should reset custom field filters on reset", () => {
    jest.spyOn(TestBed.inject(CategoryService), "getAllCategories").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(TagService), "getAllTags").mockReturnValue(
      of([]) as any
    );
    jest.spyOn(TestBed.inject(CustomFieldService), "getAllCustomFields").mockReturnValue(
      of([]) as any
    );

    const noopComponent = TestBed.createComponent(NoopComponent).componentInstance;
    component.parentForm = buildReceiptFilterForm({
      customFields: [
        { customFieldId: 1, operation: FilterOperation.Equals, value: "test" },
      ],
    }, noopComponent);
    component.ngOnInit();

    component.formCommand.subscribe((formCommand) => {
      applyFormCommand(component.parentForm, formCommand);
    });

    expect(component.getCustomFieldsArray().length).toBe(1);

    component.resetFilter();
    expect(component.getCustomFieldsArray().length).toBe(0);
  });
});
