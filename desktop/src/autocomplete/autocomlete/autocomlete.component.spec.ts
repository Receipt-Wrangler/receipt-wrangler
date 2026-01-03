import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormControl, FormGroup, ReactiveFormsModule, } from "@angular/forms";
import { MatAutocompleteModule, MatAutocompleteSelectedEvent, } from "@angular/material/autocomplete";
import { NgxsModule } from "@ngxs/store";
import { BaseInputComponent } from "../../base-input";
import { AutocomleteComponent } from "./autocomlete.component";

describe("AutocomleteComponent", () => {
  let component: AutocomleteComponent;
  let fixture: ComponentFixture<AutocomleteComponent>;
  jest.useFakeTimers();

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AutocomleteComponent, BaseInputComponent],
      imports: [
        NgxsModule.forRoot([]),
        MatAutocompleteModule,
        ReactiveFormsModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(AutocomleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should filter the options when multiple is false", () => {
    component.options = [
      { id: 1, name: "Option 1" },
      { id: 2, name: "Option 2" },
      { id: 3, name: "Option 3" },
    ];
    component.optionFilterKey = "name";

    component.multiple = false;
    const result = component._filter("Option 1");
    expect(result).toEqual([{ id: 1, name: "Option 1" }]);
  });

  // it('should filter the options when multiple is true and values are selected', () => {
  //   const service = TestBed.inject(FormBuilder);
  //   component.options = [
  //     { id: 1, name: 'Option 1' },
  //     { id: 2, name: 'Option 2' },
  //     { id: 3, name: 'Option 3' },
  //   ];
  //   component.optionFilterKey = 'name';

  //   component.multiple = true;
  //   component.inputFormControl = new FormArray([
  //     new FormGroup({
  //       id: new FormControl(2),
  //       name: new FormControl('Option 2'),
  //     }),
  //   ]) as any;

  //   const result = component._filter('Option');
  //   console.log(result);
  //   expect(result).toEqual([
  //     { id: 1, name: 'Option 1' },
  //     { id: 3, name: 'Option 3' },
  //   ]);
  // });

  it("should return an empty array when no options match the filter", () => {
    component.options = [
      { id: 1, name: "Option 1" },
      { id: 2, name: "Option 2" },
      { id: 3, name: "Option 3" },
    ];
    component.optionFilterKey = "name";

    component.multiple = false;
    const result = component._filter("Non-existing option");
    expect(result).toEqual([]);
  });

  it("should set the selected option value to the inputFormControl in single mode", () => {
    component.multiple = false;
    component.optionValueKey = "value";

    const event = {
      option: {
        id: "selected-option",
        value: "value",
      },
    } as MatAutocompleteSelectedEvent;

    component.optionSelected(event);

    expect(component.inputFormControl.value).toEqual("value");
  });

  it("should set the selected option value to the inputFormControl in single mode with full object as value", () => {
    component.multiple = false;

    const event = {
      option: {
        id: "selected-option",
        value: { id: "id1", name: "Groceries" },
      },
    } as MatAutocompleteSelectedEvent;

    component.optionSelected(event);

    expect(component.inputFormControl.value).toEqual({
      id: "id1",
      name: "Groceries",
    });
  });

  it("should add the selected option value to the inputFormControl as a FormControl instance in multiple mode", () => {
    component.multiple = true;
    component.optionValueKey = "value";
    component.inputFormControl = new FormArray([
      new FormControl("value 1"),
    ]) as any;

    const event = {
      option: {
        id: "selected-option",
        value: "value 2",
      },
    } as MatAutocompleteSelectedEvent;

    component.optionSelected(event);

    expect((component.inputFormControl as any as FormArray).value).toEqual([
      "value 1",
      "value 2",
    ]);
  });

  it("should add the selected option value to the inputFormControl as a FormControl instance in multiple mode with full object", () => {
    component.multiple = true;
    component.optionValueKey = "value";
    component.inputFormControl = new FormArray([
      new FormGroup({
        id: new FormControl("id0"),
        name: new FormControl("Utilities"),
      }),
    ]) as any;

    const event = {
      option: {
        id: "selected-option",
        value: {
          id: "id1",
          name: "Groceries",
        },
      },
    } as MatAutocompleteSelectedEvent;

    component.optionSelected(event);

    expect((component.inputFormControl as any as FormArray).value).toEqual([
      {
        id: "id0",
        name: "Utilities",
      },
      {
        id: "id1",
        name: "Groceries",
      },
    ]);
  });

  it("should add a custom option value to the inputFormControl as a FormControl instance in multiple mode with no option value key", () => {
    component.creatableOptionId = "create-option";
    component.defaultCreatableObject = { name: "Custom Option" };
    component.creatableValueKey = "name";
    component.multiple = true;
    component.inputFormControl = new FormArray([
      new FormControl("value 1"),
    ]) as any;

    const event = {
      option: {
        id: component.creatableOptionId,
        value: "Custom Option",
      },
    } as MatAutocompleteSelectedEvent;

    component.filterFormControl.setValue("new value");

    component.optionSelected(event);

    expect((component.inputFormControl as any as FormArray).value).toEqual([
      "value 1",
      {
        name: "new value",
      },
    ]);
  });

  it("should add a custom option value to the inputFormControl as a FormControl instance in multiple mode with option value key", () => {
    component.creatableOptionId = "create-option";
    component.optionValueKey = "name";
    component.multiple = true;
    component.inputFormControl = new FormArray([
      new FormControl("value 1"),
    ]) as any;

    const event = {
      option: {
        id: component.creatableOptionId,
        value: "Custom Option",
      },
    } as MatAutocompleteSelectedEvent;

    component.filterFormControl.setValue("new value");

    component.optionSelected(event);

    expect((component.inputFormControl as any as FormArray).value).toEqual([
      "value 1",
      "new value",
    ]);
  });
});
