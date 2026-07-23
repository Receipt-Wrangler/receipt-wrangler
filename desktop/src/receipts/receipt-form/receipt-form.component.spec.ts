import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule, Validators } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { BehaviorSubject, of } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { ApiModule, CustomFieldType, ReceiptImageService, ReceiptStatus } from "../../open-api";
import { SnackbarService } from "../../services";
import { QueueMode } from "../../services/receipt-queue.service";
import { StoreModule } from "../../store/store.module";
import { ReceiptFormComponent } from "./receipt-form.component";

describe("ReceiptFormComponent", () => {
  let component: ReceiptFormComponent;
  let fixture: ComponentFixture<ReceiptFormComponent>;
  let routeDataSubject: BehaviorSubject<any>;

  beforeEach(async () => {
    routeDataSubject = new BehaviorSubject<any>({});
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
            },
            data: routeDataSubject,
            params: of({})
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

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("falls back to the receipt's embedded custom field definitions when the catalog is empty", () => {
    const definition = { id: 7, name: "Project", type: 0 } as any;
    routeDataSubject.next({
      mode: FormMode.edit,
      customFields: [],
      receipt: {
        id: 1,
        name: "R",
        amount: "1.00",
        customFields: [{ customFieldId: 7, customField: definition }],
      } as any,
    });

    component.ngOnInit();

    expect(component.customFields).toEqual([definition]);
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
    component.images.set([{ id: 1 } as any]);
    component.ngOnInit();
    component.mode = FormMode.edit;
    Object.defineProperty(component, 'carouselComponent', {
      value: () => ({
        currentlyShownImageIndex: 0,
      }),
      configurable: true,
    });
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

    const receiptImageServiceSpy = jest.spyOn(
      TestBed.inject(ReceiptImageService),
      "magicFillReceipt"
    ).mockReturnValue(of(magicReceipt));

    const snackbarSpy = jest.spyOn(
      TestBed.inject(SnackbarService),
      "success"
    ).mockReturnValue(undefined);

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
    component.images.set([{ id: 1 } as any]);
    component.ngOnInit();
    component.mode = FormMode.edit;
    Object.defineProperty(component, 'carouselComponent', {
      value: () => ({
        currentlyShownImageIndex: 0,
      }),
      configurable: true,
    });

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

    const receiptImageServiceSpy = jest.spyOn(
      TestBed.inject(ReceiptImageService),
      "magicFillReceipt"
    ).mockReturnValue(of(magicReceipt));

    component.magicFill();

    expect(receiptImageServiceSpy).toHaveBeenCalledWith(1, undefined,);

    const receiptValue = component.form.getRawValue();

    expect(receiptValue.name).toEqual(magicReceipt.name);
    expect(receiptValue.amount).toEqual(originalData.amount);
    expect(receiptValue.date).toEqual(originalData.date);
  });

  it("should not patch any values when they are all default values and pop error snackbar", () => {
    component.images.set([{ id: 1 } as any]);
    component.ngOnInit();
    component.mode = FormMode.edit;
    Object.defineProperty(component, 'carouselComponent', {
      value: () => ({
        currentlyShownImageIndex: 0,
      }),
      configurable: true,
    });


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

    const receiptImageServiceSpy = jest.spyOn(
      TestBed.inject(ReceiptImageService),
      "magicFillReceipt"
    ).mockReturnValue(of(magicReceipt));

    const snackbarSpy = jest.spyOn(
      TestBed.inject(SnackbarService),
      "error"
    ).mockReturnValue(undefined);

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
    routeDataSubject.next({
      receipt: { id: 2 } as any,
    });

    expect(component.queueIndex).toEqual(1);
    expect(component.queueIds).toEqual(["1", "2", "3"]);
    expect(component.queueMode).toEqual(QueueMode.VIEW);
    expect(component.submitButtonText).toEqual("Save & Next");
  });

  it("rebuilds form state when route data emits a new receipt (duplicate navigation)", () => {
    routeDataSubject.next({
      receipt: { id: 1, name: "Original", amount: "10.00", customFields: [] } as any,
    });

    expect(component.originalReceipt?.id).toEqual(1);
    expect(component.form.get("name")?.value).toEqual("Original");
    expect(component.editLink).toEqual("/receipts/1/edit");

    routeDataSubject.next({
      receipt: { id: 2, name: "Duplicate", amount: "10.00", customFields: [] } as any,
    });

    expect(component.originalReceipt?.id).toEqual(2);
    expect(component.form.get("name")?.value).toEqual("Duplicate");
    expect(component.editLink).toEqual("/receipts/2/edit");
  });

  // Full-receipt magic-fill ingest: proves the form takes in every field the
  // backend can return (paid-by, status, items, shares, custom fields, comments),
  // not just the header fields, and drops nothing on the way into the form.
  describe("magicFill — full receipt ingest", () => {
    function stubCarousel(index = 0): void {
      Object.defineProperty(component, "carouselComponent", {
        value: () => ({ currentlyShownImageIndex: index }),
        configurable: true,
      });
    }

    function stubItemAndShareLists(): { setItems: jest.Mock; setUserItemMap: jest.Mock } {
      const setItems = jest.fn();
      const setUserItemMap = jest.fn();
      Object.defineProperty(component, "itemListComponent", {
        value: () => ({ setItems }),
        configurable: true,
      });
      Object.defineProperty(component, "shareListComponent", {
        value: () => ({ setUserItemMap }),
        configurable: true,
      });
      return { setItems, setUserItemMap };
    }

    function stubCommentsComponent(): { addMagicFilledComments: jest.Mock } {
      const addMagicFilledComments = jest.fn();
      Object.defineProperty(component, "receiptCommentsComponent", {
        value: () => ({ addMagicFilledComments }),
        configurable: true,
      });
      return { addMagicFilledComments };
    }

    function stubPaidByAutocomplete(): { syncSingleDisplay: jest.Mock } {
      const syncSingleDisplay = jest.fn();
      Object.defineProperty(component, "paidByAutocomplete", {
        value: () => ({ autocompleteComponent: () => ({ syncSingleDisplay }) }),
        configurable: true,
      });
      return { syncSingleDisplay };
    }

    function mockMagicFill(magicReceipt: any): jest.SpyInstance {
      return jest
        .spyOn(TestBed.inject(ReceiptImageService), "magicFillReceipt")
        .mockReturnValue(of(magicReceipt));
    }

    function spySuccess(): jest.SpyInstance {
      return jest.spyOn(TestBed.inject(SnackbarService), "success").mockReturnValue(undefined);
    }

    function spyError(): jest.SpyInstance {
      return jest.spyOn(TestBed.inject(SnackbarService), "error").mockReturnValue(undefined);
    }

    const customFieldDefs = [
      { id: 1, name: "Text Field", type: CustomFieldType.Text },
      { id: 2, name: "Date Field", type: CustomFieldType.Date },
      { id: 3, name: "Select Field", type: CustomFieldType.Select },
      { id: 4, name: "Currency Field", type: CustomFieldType.Currency },
      { id: 5, name: "Boolean Field", type: CustomFieldType.Boolean },
    ] as any[];

    it("ingests every field of a full magic-filled receipt", () => {
      // Fix the timezone offset so the date assertion is deterministic.
      Date.prototype.getTimezoneOffset = () => 0;
      component.images.set([{ id: 1 } as any]);
      routeDataSubject.next({ mode: FormMode.edit, customFields: customFieldDefs });
      component.mode = FormMode.edit;

      stubCarousel();
      const { setItems, setUserItemMap } = stubItemAndShareLists();
      const { addMagicFilledComments } = stubCommentsComponent();
      const { syncSingleDisplay } = stubPaidByAutocomplete();

      component.categories = [
        { id: 1, name: "category" } as any,
        { id: 2, name: "category2" } as any,
      ];
      component.tags = [
        { id: 1, name: "tag" } as any,
        { id: 2, name: "tag2" } as any,
      ];

      const magicReceipt = {
        name: "Full Receipt",
        amount: "100.00",
        date: "2023-08-05T00:00:00.000Z",
        paidByUserId: 7,
        status: ReceiptStatus.NeedsAttention,
        categories: [{ id: 1 }],
        tags: [{ id: 2 }],
        receiptItems: [
          {
            name: "Regular Item",
            amount: "40.00",
            status: "OPEN",
            categories: [{ id: 5, name: "item cat" }],
            tags: [{ id: 6, name: "item tag" }],
          },
          {
            name: "Shared Item",
            amount: "20.00",
            status: "OPEN",
            chargedToUserId: 7,
            linkedItems: [
              { name: "Linked", amount: "10.00", status: "OPEN", chargedToUserId: 7 },
            ],
          },
        ],
        customFields: [
          { customFieldId: 1, stringValue: "hello" },
          { customFieldId: 2, dateValue: "2023-08-05T00:00:00.000Z" },
          { customFieldId: 3, selectValue: 42 },
          { customFieldId: 4, currencyValue: "9.99" },
          { customFieldId: 5, booleanValue: true },
        ],
        comments: [{ comment: "auto comment", userId: 7 }],
      } as any;

      mockMagicFill(magicReceipt);
      const successSpy = spySuccess();

      component.magicFill();

      const value = component.form.getRawValue();

      // Scalars
      expect(value.name).toEqual("Full Receipt");
      expect(value.amount).toEqual("100.00");
      expect(value.date).toEqual(new Date("2023-08-05T00:00:00.000Z"));
      expect(value.paidByUserId).toEqual(7);
      expect(value.status).toEqual(ReceiptStatus.NeedsAttention);
      expect(syncSingleDisplay).toHaveBeenCalled();

      // Categories/tags matched from the pool by id
      expect(value.categories).toEqual([component.categories[0]]);
      expect(value.tags).toEqual([component.tags[1]]);

      // Items (regular + share + linked + per-item cats/tags + taxed)
      expect(value.receiptItems.length).toEqual(2);
      const regular = value.receiptItems[0];
      expect(regular.name).toEqual("Regular Item");
      expect(regular.amount).toEqual("40.00");
      expect(regular.status).toEqual("OPEN");
      expect(regular.chargedToUserId).toBeFalsy();
      expect(regular.categories).toEqual([{ id: 5, name: "item cat" }]);
      expect(regular.tags).toEqual([{ id: 6, name: "item tag" }]);
      expect(regular.linkedItems).toEqual([]);

      const shared = value.receiptItems[1];
      expect(shared.name).toEqual("Shared Item");
      expect(shared.chargedToUserId).toEqual(7);
      expect(shared.linkedItems.length).toEqual(1);
      expect(shared.linkedItems[0].name).toEqual("Linked");
      expect(shared.linkedItems[0].chargedToUserId).toEqual(7);

      expect(setItems).toHaveBeenCalled();
      expect(setUserItemMap).toHaveBeenCalled();

      // Custom fields land in the type-specific column
      expect(value.customFields.length).toEqual(5);
      const byId = (id: number) => value.customFields.find((c: any) => c.customFieldId === id);
      expect(byId(1).stringValue).toEqual("hello");
      expect(byId(2).dateValue).toEqual("2023-08-05T00:00:00.000Z");
      expect(byId(3).selectValue).toEqual(42);
      expect(byId(4).currencyValue).toEqual("9.99");
      expect(byId(5).booleanValue).toEqual(true);
      // Manage-fields menu entries flip to selected
      customFieldDefs.forEach((def) => {
        const menuItem = component.customFieldsStatefulMenuItems.find(
          (m) => m.value === def.id.toString()
        );
        expect(menuItem?.selected).toBe(true);
      });

      // Comments handed to the (mode-aware) comments child
      expect(addMagicFilledComments).toHaveBeenCalledWith(magicReceipt.comments);

      // Snackbar lists every filled field with reader-friendly labels
      expect(successSpy).toHaveBeenCalledWith(
        "Magic fill successfully filled name, amount, date, paid by, status, categories, tags, items, custom fields, comments from selected image!",
        { duration: 10000 }
      );
    });

    it("treats items with a chargedToUserId as shares and excludes them from the item total", () => {
      component.images.set([{ id: 1 } as any]);
      component.ngOnInit();
      component.mode = FormMode.edit;
      stubCarousel();
      stubItemAndShareLists();

      const magicReceipt = {
        name: "Items",
        amount: "100.00",
        date: "2023-08-05T00:00:00.000Z",
        receiptItems: [
          { name: "Regular", amount: "40.00", status: "OPEN" },
          { name: "Share", amount: "20.00", status: "OPEN", chargedToUserId: 9 },
        ],
      } as any;

      mockMagicFill(magicReceipt);
      spySuccess();

      component.magicFill();

      // Only the non-share item counts toward the receipt item total
      expect((component as any).calculateItemsTotal()).toEqual(40);
      // The share carries a required chargedToUserId validator; the regular item does not
      const shareControl = component.receiptItemsFormArray.at(1).get("chargedToUserId");
      const regularControl = component.receiptItemsFormArray.at(0).get("chargedToUserId");
      expect(shareControl?.hasValidator(Validators.required)).toBe(true);
      expect(regularControl?.hasValidator(Validators.required)).toBe(false);
    });

    it("fills nothing and shows an error when the response is all defaults/empty", () => {
      component.images.set([{ id: 1 } as any]);
      component.ngOnInit();
      component.mode = FormMode.edit;
      stubCarousel();

      const originalData = {
        name: "kept name",
        amount: "482.32",
        paidByUserId: 3,
        status: ReceiptStatus.Open,
      } as any;
      component.form.patchValue(originalData);

      const magicReceipt = {
        name: "",
        amount: "0",
        date: "0001-01-01T00:00:00Z",
        paidByUserId: 0,
        status: "",
        categories: [],
        tags: [],
        receiptItems: [],
        customFields: [],
        comments: [],
      } as any;

      mockMagicFill(magicReceipt);
      const errorSpy = spyError();

      component.magicFill();

      const value = component.form.getRawValue();
      expect(value.name).toEqual("kept name");
      expect(value.amount).toEqual("482.32");
      expect(value.paidByUserId).toEqual(3);
      expect(value.status).toEqual(ReceiptStatus.Open);
      expect(value.receiptItems).toEqual([]);
      expect(value.customFields).toEqual([]);
      expect(errorSpy).toHaveBeenCalledWith(
        "Could not find any values to fill! Try reuploading a clearer image."
      );
    });

    it("fills only the fields present and labels just those in the snackbar", () => {
      component.images.set([{ id: 1 } as any]);
      component.ngOnInit();
      component.mode = FormMode.edit;
      stubCarousel();
      stubItemAndShareLists();

      component.form.patchValue({ name: "untouched", amount: "5.00" });

      const magicReceipt = {
        receiptItems: [{ name: "Only Item", amount: "3.00", status: "OPEN" }],
      } as any;

      mockMagicFill(magicReceipt);
      const successSpy = spySuccess();

      component.magicFill();

      const value = component.form.getRawValue();
      expect(value.name).toEqual("untouched");
      expect(value.amount).toEqual("5.00");
      expect(value.receiptItems.length).toEqual(1);
      expect(successSpy).toHaveBeenCalledWith(
        "Magic fill successfully filled items from selected image!",
        { duration: 10000 }
      );
    });

    it("skips custom field values whose field is not in the catalog pool", () => {
      component.images.set([{ id: 1 } as any]);
      routeDataSubject.next({ mode: FormMode.edit, customFields: [customFieldDefs[0]] });
      component.mode = FormMode.edit;
      stubCarousel();

      const magicReceipt = {
        name: "CF",
        amount: "1.00",
        date: "2023-08-05T00:00:00.000Z",
        customFields: [{ customFieldId: 999, stringValue: "orphan" }],
      } as any;

      mockMagicFill(magicReceipt);
      const successSpy = spySuccess();

      component.magicFill();

      expect(component.customFieldsFormArray.length).toEqual(0);
      // "custom fields" is absent from the snackbar since nothing was ingested
      expect(successSpy).toHaveBeenCalledWith(
        "Magic fill successfully filled name, amount, date from selected image!",
        { duration: 10000 }
      );
    });

    it("does not duplicate a custom field the receipt already carries a value for", () => {
      component.images.set([{ id: 1 } as any]);
      routeDataSubject.next({
        mode: FormMode.edit,
        customFields: [customFieldDefs[0]],
        receipt: {
          id: 5,
          name: "R",
          amount: "1.00",
          customFields: [{ customFieldId: 1, stringValue: "existing" }],
        } as any,
      });
      component.mode = FormMode.edit;
      stubCarousel();

      const magicReceipt = {
        name: "CF",
        amount: "1.00",
        date: "2023-08-05T00:00:00.000Z",
        customFields: [{ customFieldId: 1, stringValue: "new value" }],
      } as any;

      mockMagicFill(magicReceipt);
      spySuccess();

      component.magicFill();

      expect(component.customFieldsFormArray.length).toEqual(1);
      expect(component.customFieldsFormArray.at(0).value.stringValue).toEqual("existing");
    });

    it("appends magic-filled items and categories onto existing ones", () => {
      component.images.set([{ id: 1 } as any]);
      routeDataSubject.next({
        mode: FormMode.edit,
        customFields: [],
        receipt: {
          id: 5,
          name: "R",
          amount: "100.00",
          categories: [{ id: 1, name: "existing cat" }],
          receiptItems: [{ name: "Existing Item", amount: "10.00", status: "OPEN" }],
          customFields: [],
        } as any,
      });
      component.mode = FormMode.edit;
      stubCarousel();
      stubItemAndShareLists();
      component.categories = [
        { id: 1, name: "existing cat" } as any,
        { id: 2, name: "new cat" } as any,
      ];

      const magicReceipt = {
        categories: [{ id: 2 }],
        receiptItems: [{ name: "New Item", amount: "15.00", status: "OPEN" }],
      } as any;

      mockMagicFill(magicReceipt);
      spySuccess();

      component.magicFill();

      const value = component.form.getRawValue();
      expect(value.categories).toEqual([
        { id: 1, name: "existing cat" },
        { id: 2, name: "new cat" },
      ]);
      expect(value.receiptItems.map((i: any) => i.name)).toEqual([
        "Existing Item",
        "New Item",
      ]);
    });

    it("round-trips negative refund amounts on the receipt and items", () => {
      component.images.set([{ id: 1 } as any]);
      component.ngOnInit();
      component.mode = FormMode.edit;
      stubCarousel();
      stubItemAndShareLists();

      const magicReceipt = {
        name: "Refund",
        amount: "-50.00",
        date: "2023-08-05T00:00:00.000Z",
        receiptItems: [{ name: "Refunded", amount: "-50.00", status: "OPEN" }],
      } as any;

      mockMagicFill(magicReceipt);
      spySuccess();

      component.magicFill();

      const value = component.form.getRawValue();
      expect(value.amount).toEqual("-50.00");
      expect(value.receiptItems[0].amount).toEqual("-50.00");
    });

    it("skips paid-by when the backend returns the unset sentinel (0)", () => {
      component.images.set([{ id: 1 } as any]);
      component.ngOnInit();
      component.mode = FormMode.edit;
      stubCarousel();
      const { syncSingleDisplay } = stubPaidByAutocomplete();
      component.form.patchValue({ paidByUserId: 4 });

      const magicReceipt = {
        name: "No Payer",
        amount: "1.00",
        date: "2023-08-05T00:00:00.000Z",
        paidByUserId: 0,
      } as any;

      mockMagicFill(magicReceipt);
      const successSpy = spySuccess();

      component.magicFill();

      // Existing paid-by preserved, display not refreshed, not listed in the snackbar
      expect(component.form.getRawValue().paidByUserId).toEqual(4);
      expect(syncSingleDisplay).not.toHaveBeenCalled();
      expect(successSpy).toHaveBeenCalledWith(
        "Magic fill successfully filled name, amount, date from selected image!",
        { duration: 10000 }
      );
    });
  });
});
