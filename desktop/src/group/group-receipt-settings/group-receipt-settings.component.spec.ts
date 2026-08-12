import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { GroupsService, QuickScanDefaultPaidByType } from "../../open-api";
import { PipesModule } from "../../pipes/index";
import { SnackbarService } from "../../services";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { StoreModule } from "../../store/store.module";
import { GroupReceiptSettingsComponent } from "./group-receipt-settings.component";

describe("GroupReceiptSettingsComponent", () => {
  let component: GroupReceiptSettingsComponent;
  let fixture: ComponentFixture<GroupReceiptSettingsComponent>;
  let httpTestingController: HttpTestingController;

  const quickScanDefaults = {
    quickScanPaidByEnabled: true,
    quickScanPaidByRequired: true,
    quickScanDefaultPaidByType: QuickScanDefaultPaidByType.Empty,
    quickScanDefaultPaidById: null,
    quickScanStatusEnabled: true,
    quickScanStatusRequired: true,
    quickScanDefaultStatus: "",
    quickScanCategoriesEnabled: false,
    quickScanCategoriesRequired: false,
    quickScanTagsEnabled: false,
    quickScanTagsRequired: false,
    quickScanCommentEnabled: false,
    quickScanCommentRequired: false,
  };

  const testGroup = {
    id: 1,
    groupReceiptSettings: {
      hideImages: true,
      hideReceiptCategories: false,
      hideReceiptTags: true,
      hideItemCategories: false,
      hideItemTags: true,
      hideShareCategories: false,
      hideShareTags: false,
      hideComments: false,
      ...quickScanDefaults,
    }
  };

  const configureTestBed = async (mode: FormMode = FormMode.edit) => {
    await TestBed.configureTestingModule({
      declarations: [GroupReceiptSettingsComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        SharedUiModule,
        StoreModule,
        PipesModule
      ],
      providers: [
        FormBuilder,
        GroupsService,
        { provide: Router, useValue: { navigate: jest.fn().mockResolvedValue(true) } },
        Store,
        SnackbarService,
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {
                formConfig: { mode },
                group: testGroup
              }
            }
          }
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
      ]
    }).compileComponents();

    httpTestingController = TestBed.inject(HttpTestingController);
    fixture = TestBed.createComponent(GroupReceiptSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  };

  beforeEach(async () => {
    await configureTestBed();
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize form with group receipt settings", () => {
    expect(component.form.getRawValue()).toEqual(testGroup.groupReceiptSettings);
    expect(component.editLink).toBe(`/groups/${testGroup.id}/receipt-settings/edit`);
  });

  it("should submit form and update settings", () => {
    const groupsService = TestBed.inject(GroupsService);
    const store = TestBed.inject(Store);
    const router = TestBed.inject(Router);
    const snackbarService = TestBed.inject(SnackbarService);

    jest.spyOn(groupsService, "updateGroupReceiptSettings").mockReturnValue(of(testGroup.groupReceiptSettings as any));
    jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));
    jest.spyOn(router, "navigate");
    jest.spyOn(snackbarService, "success");

    component.form.patchValue(testGroup.groupReceiptSettings);
    component.submit();

    // The empty paid-by id is coerced to undefined so the nullable id is omitted from the request.
    expect(groupsService.updateGroupReceiptSettings).toHaveBeenCalledWith(
      testGroup.id,
      { ...testGroup.groupReceiptSettings, quickScanDefaultPaidById: undefined }
    );
    expect(store.dispatch).toHaveBeenCalled();
    expect(snackbarService.success).toHaveBeenCalledWith("Receipt settings updated successfully");
    expect(router.navigate).toHaveBeenCalledWith(
      [`/groups/${testGroup.id}/receipt-settings/view`],
      { queryParams: { tab: "receipt-settings" } }
    );
  });

  it("should require a paid-by default when paid by is optional", () => {
    component.form.patchValue({
      quickScanPaidByEnabled: true,
      quickScanPaidByRequired: false,
      quickScanDefaultPaidByType: QuickScanDefaultPaidByType.Empty,
    });

    expect(component.showPaidByDefault).toBe(true);
    expect(component.form.get("quickScanDefaultPaidByType")?.valid).toBe(false);

    // Choosing "Uploader" satisfies the requirement without needing a specific user.
    component.form.patchValue({ quickScanDefaultPaidByType: QuickScanDefaultPaidByType.Uploader });
    expect(component.form.get("quickScanDefaultPaidByType")?.valid).toBe(true);
    expect(component.showPaidByUserDefault).toBe(false);
  });

  it("should require a specific user id when default type is USER", () => {
    component.form.patchValue({
      quickScanPaidByEnabled: true,
      quickScanPaidByRequired: false,
      quickScanDefaultPaidByType: QuickScanDefaultPaidByType.User,
      quickScanDefaultPaidById: null,
    });

    expect(component.showPaidByUserDefault).toBe(true);
    expect(component.form.get("quickScanDefaultPaidById")?.valid).toBe(false);

    component.form.patchValue({ quickScanDefaultPaidById: 5 });
    expect(component.form.get("quickScanDefaultPaidById")?.valid).toBe(true);
  });

  it("should not require defaults when paid by is shown and required", () => {
    component.form.patchValue({
      quickScanPaidByEnabled: true,
      quickScanPaidByRequired: true,
    });

    expect(component.showPaidByDefault).toBe(false);
    expect(component.form.get("quickScanDefaultPaidByType")?.valid).toBe(true);
  });

  it("should require a default status when status is optional", () => {
    component.form.patchValue({
      quickScanStatusEnabled: true,
      quickScanStatusRequired: false,
      quickScanDefaultStatus: "",
    });

    expect(component.showStatusDefault).toBe(true);
    expect(component.form.get("quickScanDefaultStatus")?.valid).toBe(false);
  });

  it("should disable the quick scan comment toggles while comments are hidden", () => {
    component.form.patchValue({
      quickScanCommentEnabled: true,
      quickScanCommentRequired: true,
    });
    component.form.patchValue({ hideComments: true });

    expect(component.form.get("quickScanCommentEnabled")?.disabled).toBe(true);
    expect(component.form.get("quickScanCommentRequired")?.disabled).toBe(true);
    // Disabled, not cleared: the configured values must survive so they come back on un-hide, and
    // getRawValue must still carry them so submit doesn't wipe them server-side.
    expect(component.form.getRawValue().quickScanCommentEnabled).toBe(true);
    expect(component.form.getRawValue().quickScanCommentRequired).toBe(true);
  });

  it("should re-enable the quick scan comment toggles with their values when comments are un-hidden", () => {
    component.form.patchValue({
      quickScanCommentEnabled: true,
      quickScanCommentRequired: true,
    });
    component.form.patchValue({ hideComments: true });
    component.form.patchValue({ hideComments: false });

    expect(component.form.get("quickScanCommentEnabled")?.disabled).toBe(false);
    expect(component.form.get("quickScanCommentEnabled")?.value).toBe(true);
    expect(component.form.get("quickScanCommentRequired")?.value).toBe(true);
  });

  it("should submit the comment toggles even while they are disabled by hide comments", () => {
    const groupsService = TestBed.inject(GroupsService);
    const store = TestBed.inject(Store);

    jest.spyOn(groupsService, "updateGroupReceiptSettings").mockReturnValue(of(testGroup.groupReceiptSettings as any));
    jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

    component.form.patchValue({
      quickScanCommentEnabled: true,
      quickScanCommentRequired: true,
    });
    component.form.patchValue({ hideComments: true });
    component.submit();

    expect(groupsService.updateGroupReceiptSettings).toHaveBeenCalledWith(
      testGroup.id,
      expect.objectContaining({
        hideComments: true,
        quickScanCommentEnabled: true,
        quickScanCommentRequired: true,
      })
    );
  });

  it("should leave the whole form disabled in view mode", async () => {
    // initForm disables the form after subscribing to valueChanges, and that disable emits - so the
    // comment enablement must not re-enable its controls on the read-only page.
    TestBed.resetTestingModule();
    await configureTestBed(FormMode.view);

    expect(component.form.disabled).toBe(true);
    expect(component.form.get("quickScanCommentEnabled")?.disabled).toBe(true);
    expect(component.form.get("quickScanCommentRequired")?.disabled).toBe(true);
  });
});
