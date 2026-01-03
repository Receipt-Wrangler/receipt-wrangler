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
import { ApiModule, ReceiptImageService, ReceiptStatus } from "../../open-api";
import { SnackbarService } from "../../services";
import { QueueMode } from "../../services/receipt-queue.service";
import { StoreModule } from "../../store/store.module";
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
    component.images = [{ id: 1 } as any];
    component.ngOnInit();
    component.mode = FormMode.edit;
    component.carouselComponent = {
      currentlyShownImageIndex: 0,
    } as any;
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
    component.images = [{ id: 1 } as any];
    component.ngOnInit();
    component.mode = FormMode.edit;
    component.carouselComponent = {
      currentlyShownImageIndex: 0,
    } as any;

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
    component.images = [{ id: 1 } as any];
    component.ngOnInit();
    component.mode = FormMode.edit;
    component.carouselComponent = {
      currentlyShownImageIndex: 0,
    } as any;

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
