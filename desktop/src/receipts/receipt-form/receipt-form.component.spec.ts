import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { of } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { ApiModule, ItemStatus, ReceiptImageService, ReceiptStatus } from "../../open-api";
import { SnackbarService } from "../../services";
import { QueueMode } from "../../services/receipt-queue.service";
import { StoreModule } from "../../store/store.module";
import { buildItemForm } from "../utils/form.utils";
import { ReceiptFormComponent } from "./receipt-form.component";

describe("ReceiptFormComponent", () => {
  let component: ReceiptFormComponent;
  let fixture: ComponentFixture<ReceiptFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReceiptFormComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ApiModule,
        PipesModule,
        MatDialogModule,
        MatSnackBarModule,
        StoreModule,
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        SharedUiModule
      ],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {}, queryParams: {}
            }, params: of({})
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ReceiptFormComponent);
    component = fixture.componentInstance;
    component.mode = FormMode.edit;
    fixture.detectChanges();
  });

  function setupMagicFillTest(): void {
    component.images = [{ id: 1 } as any];
    component.ngOnInit();
    component.mode = FormMode.edit;
    component.carouselComponent = {
      currentlyShownImageIndex: 0,
    } as any;
  }

  function mockMagicFillReceipt(magicReceipt: any): jest.SpyInstance {
    return jest.spyOn(
      TestBed.inject(ReceiptImageService),
      "magicFillReceipt"
    ).mockReturnValue(of(magicReceipt));
  }

  function mockSnackbarSuccess(): jest.SpyInstance {
    return jest.spyOn(
      TestBed.inject(SnackbarService),
      "success"
    ).mockReturnValue(undefined);
  }

  function mockSnackbarError(): jest.SpyInstance {
    return jest.spyOn(
      TestBed.inject(SnackbarService),
      "error"
    ).mockReturnValue(undefined);
  }

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly when there is no initial data", () => {
    jest.useFakeTimers();
    const mockedDate = new Date(2020, 0, 1);
    jest.setSystemTime(mockedDate);
    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: "",
      amount: "",
      categories: [],
      tags: [],
      date: mockedDate,
      paidByUserId: "",
      groupId: 0,
      status: ReceiptStatus.Open,
      customFields: [],
      receiptItems: [],
      syncAmountWithItems: false,
    });
    jest.useRealTimers();
  });

  it("should patch magic fill values correctly", () => {
    // Mock timezone offset to be EST
    Date.prototype.getTimezoneOffset = () => 240;
    setupMagicFillTest();
    component.categories = [
      { id: 1, name: "category" } as any,
      { id: 2, name: "category2" } as any,
    ];
    component.tags = [
      { id: 1, name: "tag" } as any,
      { id: 2, name: "tag2" } as any,
    ];

    const magicReceipt = {
      name: "magic",
      amount: "482.32",
      date: "2023-08-05T00:00:00.000Z",
      categories: [{ id: 1 } as any],
      tags: [
        {
          id: 2,
        },
      ],
    } as any;

    const receiptImageServiceSpy = mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    expect(receiptImageServiceSpy).toHaveBeenCalledWith(1, undefined);

    const receiptValue = component.form.getRawValue();

    expect(receiptValue.name).toEqual(magicReceipt.name);
    expect(receiptValue.amount).toEqual(magicReceipt.amount);
    expect(receiptValue.date).toEqual(new Date("2023-08-05T04:00:00.000Z"));
    expect(receiptValue.categories).toEqual([component.categories[0]]);
    expect(receiptValue.tags).toEqual([component.tags[1]]);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled name, amount, date, categories, tags from selected image!",
      { duration: 10000 }
    );
  });

  it("should not patch magic fill values if they are the defaults", () => {
    setupMagicFillTest();

    const originalData = {
      name: "a different name",
      amount: "482.32",
      date: "2023-08-05T04:09:12.316Z",
    } as any;

    component.form.patchValue(originalData);

    const magicReceipt = {
      name: "magic",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
    } as any;

    const receiptImageServiceSpy = mockMagicFillReceipt(magicReceipt);

    component.magicFill();

    expect(receiptImageServiceSpy).toHaveBeenCalledWith(1, undefined,);

    const receiptValue = component.form.getRawValue();

    expect(receiptValue.name).toEqual(magicReceipt.name);
    expect(receiptValue.amount).toEqual(originalData.amount);
    expect(receiptValue.date).toEqual(originalData.date);
  });

  it("should not patch any values when they are all default values and pop error snackbar", () => {
    setupMagicFillTest();

    const originalData = {
      name: "a different name",
      amount: "482.32",
      date: "2023-08-05T04:09:12.316Z",
    } as any;

    component.form.patchValue(originalData);

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
    } as any;

    const receiptImageServiceSpy = mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarError();

    component.magicFill();

    expect(receiptImageServiceSpy).toHaveBeenCalledWith(1, undefined);

    const receiptValue = component.form.getRawValue();

    expect(receiptValue.name).toEqual(originalData.name);
    expect(receiptValue.amount).toEqual(originalData.amount);
    expect(receiptValue.date).toEqual(originalData.date);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Could not find any values to fill! Try reuploading a clearer image."
    );
  });

  it("should patch receiptItems from magic fill response", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      receiptItems: [
        { name: "Item 1", amount: "10.50", status: ItemStatus.Open, receiptId: 0 },
        { name: "Item 2", amount: "5.25", status: ItemStatus.Open, receiptId: 0 },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    expect(component.receiptItemsFormArray.length).toEqual(2);
    expect(component.receiptItemsFormArray.at(0).value.name).toEqual("Item 1");
    expect(component.receiptItemsFormArray.at(0).value.amount).toEqual("10.50");
    expect(component.receiptItemsFormArray.at(1).value.name).toEqual("Item 2");
    expect(component.receiptItemsFormArray.at(1).value.amount).toEqual("5.25");
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled receiptItems from selected image!",
      { duration: 10000 }
    );
  });

  it("should not patch receiptItems when array is empty", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "magic",
      amount: "10.00",
      date: "0001-01-01T00:00:00Z",
      receiptItems: [],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.receiptItemsFormArray.length).toEqual(0);
  });

  it("should not patch receiptItems when undefined", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "magic",
      amount: "10.00",
      date: "0001-01-01T00:00:00Z",
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.receiptItemsFormArray.length).toEqual(0);
  });

  it("should patch receiptItems alongside other basic fields", () => {
    Date.prototype.getTimezoneOffset = () => 240;
    setupMagicFillTest();
    component.categories = [
      { id: 1, name: "category" } as any,
    ];
    component.tags = [
      { id: 1, name: "tag" } as any,
    ];

    const magicReceipt = {
      name: "Full Receipt",
      amount: "25.75",
      date: "2023-08-05T00:00:00.000Z",
      categories: [{ id: 1 } as any],
      tags: [{ id: 1 } as any],
      receiptItems: [
        { name: "Item A", amount: "15.00", status: ItemStatus.Open, receiptId: 0 },
        { name: "Item B", amount: "10.75", status: ItemStatus.Open, receiptId: 0 },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    const receiptValue = component.form.getRawValue();
    expect(receiptValue.name).toEqual("Full Receipt");
    expect(receiptValue.amount).toEqual("25.75");
    expect(receiptValue.categories).toEqual([component.categories[0]]);
    expect(receiptValue.tags).toEqual([component.tags[0]]);
    expect(component.receiptItemsFormArray.length).toEqual(2);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled name, amount, date, categories, tags, receiptItems from selected image!",
      { duration: 10000 }
    );
  });

  it("should build item forms with isShare=false so chargedToUserId is not required", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      receiptItems: [
        { name: "Item 1", amount: "10.00", status: ItemStatus.Open, receiptId: 0 },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    const itemForm = component.receiptItemsFormArray.at(0);
    const chargedToUserIdControl = itemForm.get("chargedToUserId");
    chargedToUserIdControl?.setValue(null);
    expect(chargedToUserIdControl?.valid).toBe(true);
  });

  it("should clear existing items before patching new ones from magic fill", () => {
    setupMagicFillTest();

    // Pre-populate with an existing item
    component.receiptItemsFormArray.push(
      buildItemForm({ name: "Existing", amount: "5.00", status: ItemStatus.Open, receiptId: 0 })
    );
    expect(component.receiptItemsFormArray.length).toEqual(1);

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      receiptItems: [
        { name: "New Item 1", amount: "10.00", status: ItemStatus.Open, receiptId: 0 },
        { name: "New Item 2", amount: "20.00", status: ItemStatus.Open, receiptId: 0 },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.receiptItemsFormArray.length).toEqual(2);
    expect(component.receiptItemsFormArray.at(0).value.name).toEqual("New Item 1");
    expect(component.receiptItemsFormArray.at(1).value.name).toEqual("New Item 2");
  });

  it("should patch item categories and tags from magic fill", () => {
    setupMagicFillTest();

    const itemCategory = { id: 1, name: "Food" };
    const itemTag = { id: 2, name: "Lunch" };

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      receiptItems: [
        {
          name: "Sandwich",
          amount: "8.50",
          status: ItemStatus.Open,
          receiptId: 0,
          categories: [itemCategory],
          tags: [itemTag],
        },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    const itemForm = component.receiptItemsFormArray.at(0);
    const categoriesArray = itemForm.get("categories") as any;
    const tagsArray = itemForm.get("tags") as any;
    expect(categoriesArray.length).toEqual(1);
    expect(categoriesArray.at(0).value).toEqual(itemCategory);
    expect(tagsArray.length).toEqual(1);
    expect(tagsArray.at(0).value).toEqual(itemTag);
  });

  it("should patch status from magic fill when valid", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      status: ReceiptStatus.Resolved,
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    expect(component.form.getRawValue().status).toEqual(ReceiptStatus.Resolved);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled status from selected image!",
      { duration: 10000 }
    );
  });

  it("should not patch status when empty", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "magic",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      status: "",
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.form.getRawValue().status).toEqual(ReceiptStatus.Open);
  });

  it("should not patch invalid status values", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      status: "COMPLETE",
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarError();

    component.magicFill();

    expect(component.form.getRawValue().status).toEqual(ReceiptStatus.Open);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Could not find any values to fill! Try reuploading a clearer image."
    );
  });

  it("should patch paidByUserId from magic fill when non-zero", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      paidByUserId: 42,
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    expect(component.form.getRawValue().paidByUserId).toEqual(42);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled paidByUserId from selected image!",
      { duration: 10000 }
    );
  });

  it("should not patch paidByUserId when zero", () => {
    setupMagicFillTest();

    component.form.patchValue({ paidByUserId: 99 });

    const magicReceipt = {
      name: "magic",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      paidByUserId: 0,
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.form.getRawValue().paidByUserId).toEqual(99);
  });

  it("should patch customFields from magic fill", () => {
    setupMagicFillTest();
    component.customFields = [
      { id: 10, name: "Vendor" } as any,
      { id: 20, name: "Notes" } as any,
    ];
    component.customFieldsStatefulMenuItems = [
      { value: "10", displayValue: "Vendor", selected: false },
      { value: "20", displayValue: "Notes", selected: false },
    ];

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      customFields: [
        { customFieldId: 10, receiptId: 0, stringValue: "ACME Corp" },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    expect(component.customFieldsFormArray.length).toEqual(1);
    expect(component.customFieldsFormArray.at(0).value.customFieldId).toEqual(10);
    expect(component.customFieldsFormArray.at(0).value.stringValue).toEqual("ACME Corp");
    expect(component.customFieldsStatefulMenuItems[0].selected).toBe(true);
    expect(component.customFieldsStatefulMenuItems[1].selected).toBe(false);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled customFields from selected image!",
      { duration: 10000 }
    );
  });

  it("should not patch customFields when array is empty", () => {
    setupMagicFillTest();

    const magicReceipt = {
      name: "magic",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      customFields: [],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.customFieldsFormArray.length).toEqual(0);
  });

  it("should clear existing customFields before patching new ones from magic fill", () => {
    setupMagicFillTest();
    component.customFields = [
      { id: 10, name: "Vendor" } as any,
      { id: 20, name: "Notes" } as any,
    ];
    component.customFieldsStatefulMenuItems = [
      { value: "10", displayValue: "Vendor", selected: true },
      { value: "20", displayValue: "Notes", selected: false },
    ];

    // Pre-populate with an existing custom field
    component.customFieldsFormArray.push(
      (component as any).buildCustomOptionFormGroup({ customFieldId: 10, receiptId: 0, stringValue: "Old Value" })
    );
    expect(component.customFieldsFormArray.length).toEqual(1);

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      customFields: [
        { customFieldId: 20, receiptId: 0, stringValue: "New Notes" },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.customFieldsFormArray.length).toEqual(1);
    expect(component.customFieldsFormArray.at(0).value.customFieldId).toEqual(20);
    expect(component.customFieldsFormArray.at(0).value.stringValue).toEqual("New Notes");
    expect(component.customFieldsStatefulMenuItems[0].selected).toBe(false);
    expect(component.customFieldsStatefulMenuItems[1].selected).toBe(true);
  });

  it("should filter out invalid customFieldIds from magic fill", () => {
    setupMagicFillTest();
    component.customFields = [
      { id: 10, name: "Vendor" } as any,
    ];
    component.customFieldsStatefulMenuItems = [
      { value: "10", displayValue: "Vendor", selected: false },
    ];

    const magicReceipt = {
      name: "",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      customFields: [
        { customFieldId: 10, receiptId: 0, stringValue: "Valid" },
        { customFieldId: 999, receiptId: 0, stringValue: "Invalid" },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    expect(component.customFieldsFormArray.length).toEqual(1);
    expect(component.customFieldsFormArray.at(0).value.customFieldId).toEqual(10);
    expect(component.customFieldsFormArray.at(0).value.stringValue).toEqual("Valid");
    expect(component.customFieldsStatefulMenuItems[0].selected).toBe(true);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled customFields from selected image!",
      { duration: 10000 }
    );
  });

  it("should not patch customFields when all IDs are invalid", () => {
    setupMagicFillTest();
    component.customFields = [
      { id: 10, name: "Vendor" } as any,
    ];
    component.customFieldsStatefulMenuItems = [
      { value: "10", displayValue: "Vendor", selected: false },
    ];

    const magicReceipt = {
      name: "magic",
      amount: "0",
      date: "0001-01-01T00:00:00Z",
      customFields: [
        { customFieldId: 888, receiptId: 0, stringValue: "Bad" },
        { customFieldId: 999, receiptId: 0, stringValue: "Also Bad" },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    expect(component.customFieldsFormArray.length).toEqual(0);
    expect(component.customFieldsStatefulMenuItems[0].selected).toBe(false);
  });

  it("should not clear existing categories when AI returns unrecognized category IDs", () => {
    setupMagicFillTest();
    component.categories = [
      { id: 1, name: "Groceries" } as any,
      { id: 2, name: "Dining" } as any,
    ];

    // Pre-populate with existing categories via a valid magic fill first
    mockMagicFillReceipt({
      name: "Setup", amount: "5", date: "2024-01-01T00:00:00.000Z",
      categories: [{ id: 1 }, { id: 2 }],
    } as any);
    mockSnackbarSuccess();
    component.magicFill();

    const categoriesArray = component.form.get("categories") as any;
    expect(categoriesArray.length).toEqual(2);

    const magicReceipt = {
      name: "Test", amount: "10", date: "2024-01-01T00:00:00.000Z",
      categories: [{ id: 99 }],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    // Existing categories should be preserved since AI returned no valid matches
    expect(categoriesArray.length).toEqual(2);
    expect(categoriesArray.value).toEqual([component.categories[0], component.categories[1]]);
  });

  it("should not clear existing tags when AI returns unrecognized tag IDs", () => {
    setupMagicFillTest();
    component.tags = [
      { id: 1, name: "Work" } as any,
      { id: 2, name: "Personal" } as any,
    ];

    // Pre-populate with existing tags via a valid magic fill first
    mockMagicFillReceipt({
      name: "Setup", amount: "5", date: "2024-01-01T00:00:00.000Z",
      tags: [{ id: 1 }],
    } as any);
    mockSnackbarSuccess();
    component.magicFill();

    const tagsArray = component.form.get("tags") as any;
    expect(tagsArray.length).toEqual(1);

    const magicReceipt = {
      name: "Test", amount: "10", date: "2024-01-01T00:00:00.000Z",
      tags: [{ id: 999 }],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    mockSnackbarSuccess();

    component.magicFill();

    // Existing tags should be preserved since AI returned no valid matches
    expect(tagsArray.length).toEqual(1);
    expect(tagsArray.value).toEqual([component.tags[0]]);
  });

  it("should clear categories on repeated magic fill", () => {
    setupMagicFillTest();
    component.categories = [
      { id: 1, name: "Food" } as any,
      { id: 2, name: "Travel" } as any,
    ];

    // First magic fill with category 1
    mockMagicFillReceipt({
      name: "First", amount: "10", date: "2024-01-01T00:00:00.000Z",
      categories: [{ id: 1 }],
    } as any);
    mockSnackbarSuccess();
    component.magicFill();
    expect(component.form.get("categories")?.value).toEqual([component.categories[0]]);

    // Second magic fill with category 2
    mockMagicFillReceipt({
      name: "Second", amount: "20", date: "2024-02-01T00:00:00.000Z",
      categories: [{ id: 2 }],
    } as any);
    component.magicFill();

    // Should have only category 2, not both
    expect(component.form.get("categories")?.value).toEqual([component.categories[1]]);
  });

  it("should clear tags on repeated magic fill", () => {
    setupMagicFillTest();
    component.tags = [
      { id: 1, name: "Work" } as any,
      { id: 2, name: "Personal" } as any,
    ];

    // First magic fill with tag 1
    mockMagicFillReceipt({
      name: "First", amount: "10", date: "2024-01-01T00:00:00.000Z",
      tags: [{ id: 1 }],
    } as any);
    mockSnackbarSuccess();
    component.magicFill();
    expect(component.form.get("tags")?.value).toEqual([component.tags[0]]);

    // Second magic fill with tag 2
    mockMagicFillReceipt({
      name: "Second", amount: "20", date: "2024-02-01T00:00:00.000Z",
      tags: [{ id: 2 }],
    } as any);
    component.magicFill();

    // Should have only tag 2, not both
    expect(component.form.get("tags")?.value).toEqual([component.tags[1]]);
  });

  it("should patch all fields together from magic fill", () => {
    Date.prototype.getTimezoneOffset = () => 240;
    setupMagicFillTest();
    component.categories = [
      { id: 1, name: "category" } as any,
    ];
    component.tags = [
      { id: 1, name: "tag" } as any,
    ];
    component.customFields = [
      { id: 5, name: "Vendor" } as any,
    ];
    component.customFieldsStatefulMenuItems = [
      { value: "5", displayValue: "Vendor", selected: false },
    ];

    const magicReceipt = {
      name: "Complete Receipt",
      amount: "99.99",
      date: "2024-01-15T00:00:00.000Z",
      categories: [{ id: 1 } as any],
      tags: [{ id: 1 } as any],
      status: ReceiptStatus.Resolved,
      paidByUserId: 7,
      receiptItems: [
        { name: "Widget", amount: "99.99", status: ItemStatus.Open, receiptId: 0 },
      ],
      customFields: [
        { customFieldId: 5, receiptId: 0, currencyValue: "99.99" },
      ],
    } as any;

    mockMagicFillReceipt(magicReceipt);
    const snackbarSpy = mockSnackbarSuccess();

    component.magicFill();

    const formValue = component.form.getRawValue();
    expect(formValue.name).toEqual("Complete Receipt");
    expect(formValue.amount).toEqual("99.99");
    expect(formValue.status).toEqual(ReceiptStatus.Resolved);
    expect(formValue.paidByUserId).toEqual(7);
    expect(formValue.categories).toEqual([component.categories[0]]);
    expect(formValue.tags).toEqual([component.tags[0]]);
    expect(component.receiptItemsFormArray.length).toEqual(1);
    expect(component.customFieldsFormArray.length).toEqual(1);
    expect(component.customFieldsStatefulMenuItems[0].selected).toBe(true);
    expect(snackbarSpy).toHaveBeenCalledWith(
      "Magic fill successfully filled name, amount, date, categories, tags, status, paidByUserId, receiptItems, customFields from selected image!",
      { duration: 10000 }
    );
  });

  it("should set queue data when there is no data", () => {
    component.ngOnInit();

    expect(component.queueIndex).toEqual(-1);
    expect(component.queueIds).toEqual([]);
    expect(component.queueMode).toEqual(undefined);
    expect(component.submitButtonText).toEqual("Save");
  });

  it("should set queue data when there is data", () => {
    TestBed.inject(ActivatedRoute).snapshot.queryParams = {
      ids: ["1", "2", "3"],
      queueMode: QueueMode.VIEW,
    };
    TestBed.inject(ActivatedRoute).snapshot.data = {
      receipt: { id: 2 } as any,
    };
    component.ngOnInit();

    expect(component.queueIndex).toEqual(1);
    expect(component.queueIds).toEqual(["1", "2", "3"]);
    expect(component.queueMode).toEqual(QueueMode.VIEW);
    expect(component.submitButtonText).toEqual("Save & Next");
  });
});
