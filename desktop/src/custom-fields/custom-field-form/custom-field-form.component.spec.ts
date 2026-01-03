import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { of } from "rxjs";
import { CustomFieldOption, CustomFieldService, CustomFieldType } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SelectModule } from "../../select/select.module";
import { SnackbarService } from "../../services";

import { CustomFieldFormComponent } from "./custom-field-form.component";

describe("CustomFieldFormComponent", () => {
  let component: CustomFieldFormComponent;
  let fixture: ComponentFixture<CustomFieldFormComponent>;
  let customFieldService: CustomFieldService;
  let snackbarService: SnackbarService;
  let matDialogRef: MatDialogRef<CustomFieldFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CustomFieldFormComponent],
      imports: [
        MatDialogModule,
        MatSnackBarModule,
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        SelectModule
      ],
      providers: [
        {
          provide: MatDialogRef,
          useValue: {
            close: jest.fn()
          }
        },
        {
          provide: CustomFieldService,
          useValue: {
            createCustomField: jest.fn().mockReturnValue(of({}))
          }
        },
        {
          provide: SnackbarService,
          useValue: {
            success: jest.fn()
          }
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    }).compileComponents();

    fixture = TestBed.createComponent(CustomFieldFormComponent);
    component = fixture.componentInstance;
    customFieldService = TestBed.inject(CustomFieldService);
    snackbarService = TestBed.inject(SnackbarService);
    matDialogRef = TestBed.inject(MatDialogRef);

    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize form with default values", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: null,
      type: null,
      description: null,
      options: []
    });
  });

  it("should initialize form with customField values when provided", () => {
    const customField: any = {
      id: 1,
      name: "Test Field",
      type: CustomFieldType.Text,
      description: "This is a test field",
      options: []
    };

    component.customField = customField;
    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: "Test Field",
      type: CustomFieldType.Text,
      description: "This is a test field",
      options: []
    });
  });

  it("should initialize form with options when customField has options", () => {
    const options: CustomFieldOption[] = [
      { id: 1, value: "Option 1", customFieldId: 1 } as any,
      { id: 2, value: "Option 2", customFieldId: 1 } as any
    ];

    const customField: any = {
      id: 1,
      name: "Test Field",
      type: CustomFieldType.Select,
      description: "This is a select field",
      options: options
    };

    component.customField = customField;
    component.ngOnInit();

    expect(component.form.get("options")?.value.length).toBe(2);
    expect(component.form.get("options")?.value[0].value).toBe("Option 1");
    expect(component.form.get("options")?.value[1].value).toBe("Option 2");
  });

  it("should add options array when type changes to Select", () => {
    component.ngOnInit();

    // Change type to Select
    component.form.get("type")?.setValue(CustomFieldType.Select);

    expect(component.form.get("options")?.value.length).toBe(1);
  });

  it("should clear options array when type changes from Select to another type", () => {
    component.ngOnInit();

    // First set to Select to add options
    component.form.get("type")?.setValue(CustomFieldType.Select);
    expect(component.form.get("options")?.value.length).toBe(1);

    // Then change to Text
    component.form.get("type")?.setValue(CustomFieldType.Text);

    expect(component.form.get("options")?.value.length).toBe(0);
  });

  it("should add a new option when addOption() is called", () => {
    component.ngOnInit();
    component.form.get("type")?.setValue(CustomFieldType.Select);
    const initialOptionsCount = component.form.get("options")?.value.length || 0;

    component.addOption();

    expect(component.form.get("options")?.value.length).toBe(initialOptionsCount + 1);
  });

  it("should delete an option when deleteOption() is called", () => {
    component.ngOnInit();
    component.form.get("type")?.setValue(CustomFieldType.Select);
    component.addOption(); // Now we have 2 options

    component.deleteOption(0);

    expect(component.form.get("options")?.value.length).toBe(1);
  });

  it("should call createCustomField when form is submitted", () => {
    component.ngOnInit();
    component.form.patchValue({
      name: "Test Field",
      type: CustomFieldType.Text,
      description: "Test description"
    });

    component.submit();

    expect(customFieldService.createCustomField).toHaveBeenCalledWith({
      name: "Test Field",
      type: CustomFieldType.Text,
      description: "Test description",
      options: []
    });
  });

  it("should show success message and close dialog when form is submitted successfully", () => {
    component.ngOnInit();
    component.form.patchValue({
      name: "Test Field",
      type: CustomFieldType.Text,
      description: "Test description"
    });

    component.submit();

    expect(snackbarService.success).toHaveBeenCalledWith("Custom field created");
    expect(matDialogRef.close).toHaveBeenCalledWith(true);
  });

  it("should close dialog when closeDialog is called", () => {
    component.closeDialog();

    expect(matDialogRef.close).toHaveBeenCalledWith(false);
  });

  it("should not submit when form is invalid", () => {
    component.ngOnInit();
    // Name is required but not set

    component.submit();

    expect(customFieldService.createCustomField).not.toHaveBeenCalled();
    expect(snackbarService.success).not.toHaveBeenCalled();
    expect(matDialogRef.close).not.toHaveBeenCalled();
  });

  it("should set type options correctly", () => {
    expect(component.typeOptions.length).toBe(Object.keys(CustomFieldType).length);
    expect(component.typeOptions[0].displayValue).toBeTruthy();
    expect(component.typeOptions[0].value).toBeTruthy();
  });

  it("should add validators to options when type is Select", () => {
    component.ngOnInit();
    component.form.get("type")?.setValue(CustomFieldType.Select);

    expect(component.form.get("options")?.hasValidator).toBeTruthy();
  });

  it("should remove validators from options when type changes from Select", () => {
    component.ngOnInit();
    component.form.get("type")?.setValue(CustomFieldType.Select);
    component.form.get("type")?.setValue(CustomFieldType.Text);

    expect(component.form.get("options")?.validator).toBeNull();
  });
});
