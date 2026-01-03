import { ChangeDetectorRef, Component, OnInit, ViewChild } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { BaseFormComponent } from "../../form/index";
import { UserPreferencesService, UserShortcut } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState, SetUserPreferences } from "../../store";
import { UserShortcutComponent } from "../user-shortcut/user-shortcut.component";

@Component({
    selector: "app-user-preferences",
    templateUrl: "./user-preferences.component.html",
    styleUrls: ["./user-preferences.component.scss"],
    standalone: false
})
export class UserPreferencesComponent extends BaseFormComponent implements OnInit {
  @ViewChild(UserShortcutComponent) public userShortcutComponent!: UserShortcutComponent;

  public formMode = FormMode;

  public originalUserShortcuts: UserShortcut[] = [];

  constructor(
    private activatedRoute: ActivatedRoute,
    private formBuilder: FormBuilder,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
    private userPreferencesService: UserPreferencesService,
    private cdr: ChangeDetectorRef
  ) {
    super();
  }

  public get userShortcuts(): FormArray {
    return this.form.get("userShortcuts") as FormArray;
  }

  public ngOnInit(): void {
    this.formConfig = this.activatedRoute.snapshot.data["formConfig"];
    this.initForm();
  }

  private initForm(): void {
    const userPreferences = this.store.selectSnapshot(
      AuthState.userPreferences
    );
    this.originalUserShortcuts = userPreferences?.userShortcuts ?? [];

    this.form = this.formBuilder.group({
      showLargeImagePreviews: userPreferences?.showLargeImagePreviews ?? false,
      quickScanDefaultPaidById: userPreferences?.quickScanDefaultPaidById ?? "",
      quickScanDefaultGroupId: userPreferences?.quickScanDefaultGroupId ?? "",
      quickScanDefaultStatus: userPreferences?.quickScanDefaultStatus ?? "",
      userShortcuts: this.formBuilder.array(this.originalUserShortcuts.map((userShortcut, i) => this.buildUserShortcut(i, userShortcut))),
    });


    if (this.formConfig.mode === FormMode.view) {
      this.form.get("quickScanDefaultStatus")?.disable();
      this.form.get("showLargeImagePreviews")?.disable();
    }
  }

  private buildUserShortcut(index: number, userShortcut?: UserShortcut): FormGroup {
    return this.formBuilder.group({
      trackby: index,
      name: this.formBuilder.control(userShortcut?.name ?? "", Validators.required),
      icon: this.formBuilder.control(userShortcut?.icon ?? "", Validators.required),
      url: this.formBuilder.control(userShortcut?.url ?? "", Validators.required),
    });
  }

  public addNewShortcut(): void {
    const userShortcuts = this.form.get("userShortcuts") as FormArray;
    const newUserShortcut = this.buildUserShortcut(
      userShortcuts.length
    );
    userShortcuts.push(newUserShortcut);
    this.originalUserShortcuts = [...this.originalUserShortcuts, newUserShortcut.value];
    this.cdr.detectChanges();

    this.userShortcutComponent.editableListComponent.openLastRow();
    this.userShortcutComponent.isAddingShortcut = true;
  }

  public shortcutDoneClicked(): void {
    if (this.userShortcuts.at(this.userShortcuts.length - 1).valid) {
      if (this.userShortcutComponent.isAddingShortcut) {
        this.originalUserShortcuts.push(this.userShortcuts.at(this.userShortcuts.length - 1).value);
      } else {
        const currentOpen = this.userShortcutComponent.editableListComponent.getCurrentRowOpen();
        if (currentOpen !== undefined && currentOpen >= 0) {
          this.originalUserShortcuts[currentOpen] = this.userShortcuts.at(currentOpen).value;
        }
      }

      this.userShortcutComponent.isAddingShortcut = false;
      this.userShortcutComponent.editableListComponent.closeRow();
    }
  }

  public shortcutCancelClicked(): void {
    if (this.userShortcutComponent.isAddingShortcut) {
      this.userShortcuts.removeAt(this.userShortcuts.length - 1);
      this.originalUserShortcuts = this.originalUserShortcuts.slice(0, this.originalUserShortcuts.length - 1);
    } else {
      const currentOpen = this.userShortcutComponent.editableListComponent.getCurrentRowOpen();
      if (currentOpen !== undefined && currentOpen >= 0) {
        this.userShortcuts.at(currentOpen).patchValue(this.originalUserShortcuts[currentOpen]);
      }
    }

    this.userShortcutComponent.isAddingShortcut = false;
    this.userShortcutComponent.editableListComponent.closeRow();
  }


  public submit(): void {
    if (this.form.valid) {
      const result = this.form.value;
      if (result.quickScanDefaultPaidById === "") {
        result.quickScanDefaultPaidById = null;
      }

      if (result.quickScanDefaultGroupId === "") {
        result.quickScanDefaultGroupId = null;
      }

      this.userPreferencesService
        .updateUserPreferences(result)
        .pipe(
          take(1),
          tap((updatedUserPreferences) => {
            this.snackbarService.success(
              "User preferences successfully updated"
            );
            this.store.dispatch(new SetUserPreferences(updatedUserPreferences));
            this.router.navigate(["/settings/user-preferences/view"],
              {
                queryParams: {
                  tab: "user-preferences",
                }
              });
          })
        )
        .subscribe();
    }
  }
}
