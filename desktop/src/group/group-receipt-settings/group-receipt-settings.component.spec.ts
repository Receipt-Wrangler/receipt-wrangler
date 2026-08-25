import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormBuilder } from "@angular/forms";
import { By } from "@angular/platform-browser";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { CustomFieldType, GroupsService, Permission, QuickScanDefaultPaidByType } from "../../open-api";
import { PipesModule } from "../../pipes/index";
import { SnackbarService } from "../../services";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SetPermissions } from "../../store/auth.state.actions";
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

  const customFieldOne = { id: 1, name: "Cost Centre", type: CustomFieldType.Text } as any;

  const customFieldTwo = { id: 2, name: "PO Number", type: CustomFieldType.Text } as any;

  // The catalog holder's group: two configured defaults (deliberately not in id
  // order, so seeding is proven to follow the configured order) and the ingest
  // toggle on.
  const testGroupWithDefaults = {
    id: 1,
    groupReceiptSettings: {
      ...testGroup.groupReceiptSettings,
      defaultCustomFieldIds: [2, 1],
      applyDefaultCustomFieldsOnIngest: true,
    }
  };

  interface TestBedOptions {
    group?: any;
    customFields?: any[];
    appPermissions?: string[];
  }

  const configureTestBed = async (mode: FormMode = FormMode.edit, options: TestBedOptions = {}) => {
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
                group: options.group ?? testGroup,
                customFields: options.customFields ?? []
              }
            }
          }
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
      ]
    }).compileComponents();

    httpTestingController = TestBed.inject(HttpTestingController);
    // Permissions must land before the first change-detection pass: initForm
    // decides which controls to build from a one-shot selectSnapshot.
    TestBed.inject(Store).dispatch(new SetPermissions(options.appPermissions ?? [], {}));
    fixture = TestBed.createComponent(GroupReceiptSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  };

  // Re-configures the bed as an admin who holds app.custom-fields.read, the only
  // caller for whom the Default Custom Fields section exists at all.
  const configureWithCustomFieldPermission = async (mode: FormMode = FormMode.edit) => {
    TestBed.resetTestingModule();
    await configureTestBed(mode, {
      group: testGroupWithDefaults,
      customFields: [customFieldOne, customFieldTwo],
      appPermissions: [Permission.AppCustomFieldsRead],
    });
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

  describe("default custom fields", () => {
    it("should not build either default custom field control without the permission", () => {
      // The wipe guard: the command reads a missing key as "leave unchanged", so an
      // admin who cannot see the catalog must submit neither key.
      expect(component.canManageDefaultCustomFields).toBe(false);
      expect(component.form.get("defaultCustomFields")).toBeNull();
      expect(component.form.get("applyDefaultCustomFieldsOnIngest")).toBeNull();
    });

    it("should omit both default custom field keys from the command without the permission", () => {
      const groupsService = TestBed.inject(GroupsService);
      const store = TestBed.inject(Store);

      const updateSpy = jest
        .spyOn(groupsService, "updateGroupReceiptSettings")
        .mockReturnValue(of(testGroup.groupReceiptSettings as any));
      jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

      component.submit();

      const command = updateSpy.mock.calls[0][1] as any;
      expect(command).not.toHaveProperty("defaultCustomFieldIds");
      expect(command).not.toHaveProperty("applyDefaultCustomFieldsOnIngest");
    });

    it("should seed the picker with the configured catalog objects", async () => {
      await configureWithCustomFieldPermission();

      const array = component.defaultCustomFieldsFormArray;
      expect(array.value).toEqual([customFieldTwo, customFieldOne]);
      // Same instances as the resolved catalog, not copies: app-autocomlete filters
      // already-selected options by reference equality, so a rebuilt literal would
      // leave the field in the dropdown and let it be added twice.
      expect(array.at(0).value).toBe(customFieldTwo);
      expect(array.at(1).value).toBe(customFieldOne);
      expect(component.form.get("applyDefaultCustomFieldsOnIngest")?.value).toBe(true);
    });

    it("should submit the selected default custom fields as ids", async () => {
      await configureWithCustomFieldPermission();

      const groupsService = TestBed.inject(GroupsService);
      const store = TestBed.inject(Store);

      const updateSpy = jest
        .spyOn(groupsService, "updateGroupReceiptSettings")
        .mockReturnValue(of(testGroupWithDefaults.groupReceiptSettings as any));
      jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

      component.submit();

      const command = updateSpy.mock.calls[0][1] as any;
      expect(command.defaultCustomFieldIds).toEqual([2, 1]);
      expect(command.applyDefaultCustomFieldsOnIngest).toBe(true);
      // The FormArray of whole CustomField objects must never ride the payload.
      expect(command).not.toHaveProperty("defaultCustomFields");
    });

    it("should submit an empty id array when every default is cleared", async () => {
      await configureWithCustomFieldPermission();

      const groupsService = TestBed.inject(GroupsService);
      const store = TestBed.inject(Store);

      const updateSpy = jest
        .spyOn(groupsService, "updateGroupReceiptSettings")
        .mockReturnValue(of(testGroupWithDefaults.groupReceiptSettings as any));
      jest.spyOn(store, "dispatch").mockReturnValue(of(undefined));

      component.defaultCustomFieldsFormArray.clear();
      component.submit();

      expect((updateSpy.mock.calls[0][1] as any).defaultCustomFieldIds).toEqual([]);
    });

    it("should mark the picker readonly in view mode", async () => {
      // form.disable() is not enough: app-autocomlete's readonly is a plain @Input
      // and is what hides the chips' remove buttons, so a viewer could otherwise
      // still remove a configured field.
      await configureWithCustomFieldPermission(FormMode.view);

      const picker = fixture.debugElement.query(
        By.css("[data-testid='group-default-custom-fields']")
      );
      expect(picker).toBeTruthy();
      expect(picker.properties["readonly"]).toBe(true);
    });

    it("should not mark the picker readonly in edit mode", async () => {
      await configureWithCustomFieldPermission();

      const picker = fixture.debugElement.query(
        By.css("[data-testid='group-default-custom-fields']")
      );
      expect(picker.properties["readonly"]).toBe(false);
    });
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
