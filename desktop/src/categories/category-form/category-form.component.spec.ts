import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef, } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { of } from "rxjs";
import { DuplicateValidator } from "src/validators/duplicate-validator";
import { ApiModule, CategoryService, CategoryView } from "../../open-api";
import { PipesModule } from "../../pipes";
import { CategoryForm } from "./category-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("CategoryForm", () => {
  let component: CategoryForm;
  let fixture: ComponentFixture<CategoryForm>;

  beforeEach(() => {
    TestBed.configureTestingModule({
    declarations: [CategoryForm],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        MatDialogModule,
        MatSnackBarModule,
        PipesModule,
        ReactiveFormsModule,
        NoopAnimationsModule],
    providers: [
        DuplicateValidator,
        {
            provide: MatDialogRef,
            useValue: {
                close: () => { },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
});
    fixture = TestBed.createComponent(CategoryForm);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form with data", () => {
    const category: CategoryView = {
      id: 1,
      name: "test",
      description: "test",
      numberOfReceipts: 1,
    };
    component.category = category;

    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: "test",
      description: "test",
    });
  });

  it("should submit form with correct data, when editing", () => {
    const categoryServiceSpy = jest.spyOn(
      TestBed.inject(CategoryService),
      "updateCategory"
    );
    categoryServiceSpy.mockReturnValue(of({} as any));
    const nameValidateSpy = jest.spyOn(
      TestBed.inject(CategoryService),
      "getCategoryCountByName"
    ).mockReturnValue(of(0) as any);
    const category: CategoryView = {
      id: 1,
      name: "test",
      description: "test",
      numberOfReceipts: 1,
    };
    component.category = category;

    component.ngOnInit();
    component.submit();

    expect(categoryServiceSpy).toHaveBeenCalledWith(
      1,
      {
        id: 1,
        name: "test",
        description: "test",
      },
    );
  });

  it("should submit form with correct data, when creating", () => {
    const nameValidateSpy = jest.spyOn(
      TestBed.inject(CategoryService),
      "getCategoryCountByName"
    ).mockReturnValue(of(0) as any);
    const categoryServiceSpy = jest.spyOn(
      TestBed.inject(CategoryService),
      "createCategory"
    );
    categoryServiceSpy.mockReturnValue(of({} as any));

    component.ngOnInit();
    component.form.patchValue({
      name: "test",
      description: "test",
    });
    component.submit();

    expect(categoryServiceSpy).toHaveBeenCalledWith({
      name: "test",
      description: "test",
    });
  });
});
