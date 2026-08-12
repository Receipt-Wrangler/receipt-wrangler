import { Component, DestroyRef, OnInit, inject } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormBuilder } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { switchMap, take, tap } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { BaseFormComponent, setRequired } from "../../form/index";
import { Group, GroupsService, Permission, QuickScanDefaultPaidByType } from "../../open-api/index";
import { SnackbarService } from "../../services/index";
import { AuthState, UpdateGroup } from "../../store/index";

@Component({
    selector: "app-group-receipt-settings",
    templateUrl: "./group-receipt-settings.component.html",
    styleUrl: "./group-receipt-settings.component.scss",
    standalone: false
})
export class GroupReceiptSettingsComponent extends BaseFormComponent implements OnInit {
  public originalGroup!: Group;

  public editLink: string = "";

  public canEdit = false;

  public readonly paidByTypeOptions = [
    { value: QuickScanDefaultPaidByType.Uploader, display: "Uploader" },
    { value: QuickScanDefaultPaidByType.User, display: "Specific user" },
  ];

  private readonly destroyRef = inject(DestroyRef);

  constructor(
    private activatedRoute: ActivatedRoute,
    private formBuilder: FormBuilder,
    private groupsService: GroupsService,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {
    super();
  }

  // The quick-scan default for paid-by/status is only needed when the field is not both shown and
  // required (otherwise the user always supplies it). These getters drive both conditional display
  // and the required validators, mirroring the backend rule.
  public get showPaidByDefault(): boolean {
    return !(
      this.form?.get("quickScanPaidByEnabled")?.value &&
      this.form?.get("quickScanPaidByRequired")?.value
    );
  }

  public get showPaidByUserDefault(): boolean {
    return (
      this.showPaidByDefault &&
      this.form?.get("quickScanDefaultPaidByType")?.value === QuickScanDefaultPaidByType.User
    );
  }

  public get showStatusDefault(): boolean {
    return !(
      this.form?.get("quickScanStatusEnabled")?.value &&
      this.form?.get("quickScanStatusRequired")?.value
    );
  }

  public ngOnInit(): void {
    this.setFormConfigFromRoute(this.activatedRoute);
    this.setOriginalGroup();
    this.initForm();
    this.canEdit = this.store.selectSnapshot(
      AuthState.hasGroupPermission(this.originalGroup.id, Permission.GroupUpdate)
    );
  }

  private initForm(): void {
    const receiptSettings = this.originalGroup.groupReceiptSettings;
    this.form = this.formBuilder.group({
      hideImages: [receiptSettings.hideImages ?? false],
      hideReceiptCategories: [receiptSettings.hideReceiptCategories ?? false],
      hideReceiptTags: [receiptSettings.hideReceiptTags ?? false],
      hideItemCategories: [receiptSettings.hideItemCategories ?? false],
      hideItemTags: [receiptSettings.hideItemTags ?? false],
      hideShareCategories: [receiptSettings.hideShareCategories ?? false],
      hideShareTags: [receiptSettings.hideShareTags ?? false],
      hideComments: [receiptSettings.hideComments ?? false],
      quickScanPaidByEnabled: [receiptSettings.quickScanPaidByEnabled ?? true],
      quickScanPaidByRequired: [receiptSettings.quickScanPaidByRequired ?? true],
      quickScanDefaultPaidByType: [receiptSettings.quickScanDefaultPaidByType ?? QuickScanDefaultPaidByType.Empty],
      quickScanDefaultPaidById: [receiptSettings.quickScanDefaultPaidById ?? null],
      quickScanStatusEnabled: [receiptSettings.quickScanStatusEnabled ?? true],
      quickScanStatusRequired: [receiptSettings.quickScanStatusRequired ?? true],
      quickScanDefaultStatus: [receiptSettings.quickScanDefaultStatus ?? ""],
      quickScanCategoriesEnabled: [receiptSettings.quickScanCategoriesEnabled ?? false],
      quickScanCategoriesRequired: [receiptSettings.quickScanCategoriesRequired ?? false],
      quickScanTagsEnabled: [receiptSettings.quickScanTagsEnabled ?? false],
      quickScanTagsRequired: [receiptSettings.quickScanTagsRequired ?? false],
      quickScanCommentEnabled: [receiptSettings.quickScanCommentEnabled ?? false],
      quickScanCommentRequired: [receiptSettings.quickScanCommentRequired ?? false],
    });

    this.applyQuickScanDerivedState();
    this.form.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => this.applyQuickScanDerivedState());

    if (this.formConfig.mode != FormMode.edit) {
      this.form.disable();
    }
  }

  // Keeps the controls that depend on other toggles in sync: the default-value controls' required
  // validators, and the comment toggles' enabled state. Runs on init and on every value change.
  private applyQuickScanDerivedState(): void {
    const paidByType = this.form.get("quickScanDefaultPaidByType")?.value;
    setRequired(this.form.get("quickScanDefaultPaidByType"), this.showPaidByDefault);
    setRequired(
      this.form.get("quickScanDefaultPaidById"),
      this.showPaidByDefault && paidByType === QuickScanDefaultPaidByType.User
    );
    setRequired(this.form.get("quickScanDefaultStatus"), this.showStatusDefault);
    this.applyQuickScanCommentEnablement();
  }

  // Hiding comments for the group hides the quick-scan comment field too, so its toggles are greyed
  // out rather than cleared - the configured values stay put and come back when Hide Comments is
  // turned off again. submit() reads getRawValue() so a disabled toggle is still sent.
  private applyQuickScanCommentEnablement(): void {
    // In view mode the whole form is disabled (see initForm); never re-enable a control there.
    if (this.formConfig.mode !== FormMode.edit) {
      return;
    }

    const commentsHidden = !!this.form.get("hideComments")?.value;
    for (const controlName of ["quickScanCommentEnabled", "quickScanCommentRequired"]) {
      const control = this.form.get(controlName);
      if (!control) {
        continue;
      }

      // emitEvent: false - this runs inside the valueChanges subscription, and enable()/disable()
      // emit valueChanges by default, which would recurse forever.
      if (commentsHidden) {
        control.disable({ emitEvent: false });
      } else {
        control.enable({ emitEvent: false });
      }
    }
  }

  private setOriginalGroup(): void {
    this.originalGroup = this.activatedRoute.snapshot.data["group"];
    this.editLink = `/groups/${this.originalGroup.id}/receipt-settings/edit`;
  }

  public submit(): void {
    if (this.form.valid) {
      // getRawValue, not value: a disabled control is omitted from form.value, so the comment
      // toggles disabled by Hide Comments would be sent as undefined, unmarshal as false, and wipe
      // the admin's stored configuration.
      const value = this.form.getRawValue();
      const command = {
        ...value,
        // An empty user autocomplete yields "" / null; send undefined so the nullable id is omitted
        // rather than sent as a non-numeric value.
        quickScanDefaultPaidById: value.quickScanDefaultPaidById
          ? Number(value.quickScanDefaultPaidById)
          : undefined,
      };

      this.groupsService.updateGroupReceiptSettings(this.originalGroup.id, command)
        .pipe(
          take(1),
          switchMap((updatedGroupReceiptSettings) => {
            this.originalGroup.groupReceiptSettings = updatedGroupReceiptSettings;
            return this.store.dispatch(new UpdateGroup(this.originalGroup));
          }),
          tap(() => {
            this.snackbarService.success("Receipt settings updated successfully");
            this.router.navigate(
              [`/groups/${this.originalGroup.id}/receipt-settings/view`],
              {
                queryParams: {
                  tab: "receipt-settings"
                }
              }
            );
          })
        )
        .subscribe();
    }
  }
}
