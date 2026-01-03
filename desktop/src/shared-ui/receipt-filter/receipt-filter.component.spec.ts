import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { Component, CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef, } from "@angular/material/dialog";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { of, tap } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { SetReceiptFilter } from "src/store/receipt-table.actions";
import { defaultReceiptFilter, } from "src/store/receipt-table.state";
import { InputModule } from "../../input";
import { CategoryService, FilterOperation, ReceiptStatus, TagService } from "../../open-api";
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
    store.reset({
      receiptTable: {
        filter: filledFilter,
      },
    });

    component.formCommand.pipe(tap((formCommand) => {
      applyFormCommand(component.parentForm, formCommand);
    })).subscribe();

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
});
