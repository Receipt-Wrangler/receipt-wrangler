import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed, } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { CarouselModule } from "ngx-bootstrap/carousel";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { LayoutState } from "src/store/layout.state";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { ApiModule, ReceiptStatus } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services";
import { AuthState, GroupState } from "../../store";
import { QuickScanDialogComponent } from "./quick-scan-dialog.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("QuickScanDialogComponent", () => {
  let component: QuickScanDialogComponent;
  let fixture: ComponentFixture<QuickScanDialogComponent>;
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
    declarations: [QuickScanDialogComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        CarouselModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, GroupState, LayoutState]),
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        SharedUiModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {},
        },
        {
            provide: MatDialog,
            useValue: {}
        },
        {
            provide: MatDialogRef<QuickScanDialogComponent>,
            useValue: {
                close: () => { },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
});
    fixture = TestBed.createComponent(QuickScanDialogComponent);
    store = TestBed.inject(Store);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      paidByUserIds: [],
      statuses: [],
      groupIds: [],
    });
  });

  it("should push fileData into images when loaded", () => {
    const originalCreateObjectURL = URL.createObjectURL;
    URL.createObjectURL = jest.fn().mockReturnValue("awesome");
    const fileData = {} as ReceiptFileUploadCommand;
    component.fileLoaded(fileData);

    expect(component.images).toEqual([fileData]);
    URL.createObjectURL = originalCreateObjectURL;
  });

  it("should close the dialog", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialogRef), "close");

    component.cancelButtonClicked();

    expect(dialogSpy).toHaveBeenCalled();
  });

  it("should show error if no image has been selected", () => {
    const snackbarSpy = jest.spyOn(TestBed.inject(SnackbarService), "error");

    component.submitButtonClicked();

    expect(snackbarSpy).toHaveBeenCalledWith(
      "Please select images to upload"
    );
  });

  it("should push new image when there are no user preferences", () => {
    component.fileLoaded({} as any);

    expect(component.form.value).toEqual({
      paidByUserIds: [""],
      statuses: [""],
      groupIds: [""]
    });
    expect(component.images).toEqual([{} as any]);
  });

  it("should push new image when there user preferences", () => {
    store.reset({
      auth: {
        userPreferences: {
          quickScanDefaultPaidById: 1,
          quickScanDefaultStatus: ReceiptStatus.Open,
          quickScanDefaultGroupId: 1,
        },
      },
    });

    component.fileLoaded({} as any);

    expect(component.form.value).toEqual({
      paidByUserIds: [1],
      statuses: [ReceiptStatus.Open],
      groupIds: [1]
    });
    expect(component.images).toEqual([{} as any]);
  });


  // TODO: fix
  // it('should call API with command', () => {
  //   const serviceSpy = jest.spyOn(
  //     TestBed.inject(ReceiptService),
  //     'quickScanReceipt'
  //   ).mockReturnValue(of({} as any));
  //   const fileData = {
  //     fileType: 'image/jpeg',
  //     imageData: '',
  //     name: 'awesome',
  //     size: 100,
  //   } as any;
  //   component.images = [fileData];

  //   store.reset({
  //     auth: {
  //       userId: 1,
  //     },
  //     groups: {
  //       selectedGroupId: 1,
  //     },
  //   });

  //   component.ngOnInit();
  //   component.submitButtonClicked();

  //   expect(serviceSpy).toHaveBeenCalledWith({
  //     imageData: [],
  //     name: 'awesome',
  //     fileType: 'image/jpeg',
  //     size: 100,
  //     groupId: 1,
  //     status: ReceiptStatus.Open,
  //     paidByUserId: 1,
  //   });
  // });
});
