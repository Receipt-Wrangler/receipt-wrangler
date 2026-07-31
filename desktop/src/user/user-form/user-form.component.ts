import { Component, DestroyRef, Input, OnInit, computed, effect, inject, signal } from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { catchError, defer, iif, of, startWith, switchMap, take, tap, } from "rxjs";
import { UserValidators } from "src/validators/user-validators";
import {
  Category,
  CategoryService,
  GroupsService,
  PermissionScope,
  Role,
  RoleService,
  Tag,
  TagService,
  User,
  UserService,
} from "../../open-api";
import { GrantSelection } from "../../shared-ui/grant-picker/grant-picker.component";
import {
  buildMemberGrantRows,
  MemberGrantRow,
  saveChangedMemberGrants,
} from "../../shared-ui/grant-picker/member-grant-assignment";
import { GroupState } from "../../store";
import { openRolePreviewDialog } from "../../roles/role-preview/role-preview-dialog.component";
import { SnackbarService, TokenRefreshService } from "../../services";
import { AddUser, AuthState, UpdateUser } from "../../store";

@UntilDestroy()
@Component({
    selector: "app-user-form",
    templateUrl: "./user-form.component.html",
    styleUrls: ["./user-form.component.scss"],
    providers: [UserValidators],
    standalone: false
})
export class UserFormComponent implements OnInit {
  @Input() public user?: User;

  public isDummerUserHelpText: string =
    "A dummy user is a user who cannot log in, but can still act as a receipt payer, or be charged shares. Dummy users can be converted to normal users, but normal users cannot be converted to dummy users.";

  constructor(
    private formBuilder: FormBuilder,
    private snackbarService: SnackbarService,
    private store: Store,
    private tokenRefreshService: TokenRefreshService,
    private userService: UserService,
    private userValidators: UserValidators,
    private matDialog: MatDialog,
    private destroyRef: DestroyRef,
    private categoryService: CategoryService,
    private tagService: TagService,
    private groupsService: GroupsService,
    public matDialogRef: MatDialogRef<UserFormComponent>
  ) {
    // Pre-select the configured default app role on the add form once the roles
    // resolve. Syncing a derived value into the reactive form is an imperative
    // side effect, so an effect() is the sanctioned tool here.
    effect(() => {
      const roles = this.roles();
      if (!this.user && this.form.get("appRoleId")?.value == null) {
        const defaultRole = roles.find(
          (role) => role.scope === PermissionScope.App && role.isDefault
        );
        if (defaultRole) {
          this.form.get("appRoleId")?.setValue(defaultRole.id);
        }
      }
    });
  }

  public form: FormGroup = new FormGroup({});

  // The full app-role list, the selector options and the currently-selected role
  // (for the preview button) are derived reactively so they update once the roles
  // load asynchronously. Degrade gracefully if the roles can't be read: the
  // selector is simply empty rather than erroring.
  private readonly roles = toSignal(
    inject(RoleService)
      .getRoles()
      .pipe(catchError(() => of([] as Role[]))),
    { initialValue: [] as Role[] }
  );
  public readonly appRoleOptions = computed<Role[]>(() =>
    this.roles().filter((role) => role.scope === PermissionScope.App)
  );
  private readonly selectedAppRoleId = signal<number | null>(null);
  public readonly selectedRole = computed<Role | null>(
    () =>
      this.appRoleOptions().find(
        (role) => role.id === this.selectedAppRoleId()
      ) ?? null
  );

  // ----- Per-member category/tag assignment -----
  // Only meaningful when editing an existing user: grants hang off a group
  // MEMBERSHIP, and a user being created has none yet. The admin assigns them on
  // a second pass, after the user exists.
  public readonly categoryPool = signal<Category[]>([]);
  public readonly tagPool = signal<Tag[]>([]);
  private readonly groups = this.store.selectSignal(GroupState.groups);

  /** One row per group the user belongs to, each carrying its role's ceiling. */
  public readonly grantRows = computed<MemberGrantRow[]>(() => {
    if (!this.user) {
      return [];
    }
    return buildMemberGrantRows(
      this.user.id,
      this.groups(),
      this.roles().filter((role) => role.scope === PermissionScope.Group),
    );
  });

  // Edits made in the pickers, keyed by group id. Only groups present here (and
  // actually changed) are written on submit.
  private readonly editedGrants = new Map<number, GrantSelection>();

  public ngOnInit(): void {
    this.initForm();
    this.listenToAppRoleChanges();
    if (!this.user) {
      this.listenToIsDummyChanges();
    } else {
      this.loadGrantPools();
    }
  }

  /**
   * The category/tag pools the assignment pickers choose from. This form is
   * admin-only, so the global lists are readable; degrade to empty rather than
   * erroring if they are not.
   */
  private loadGrantPools(): void {
    this.categoryService
      .getAllCategories()
      .pipe(take(1), catchError(() => of([] as Category[])), takeUntilDestroyed(this.destroyRef))
      .subscribe((categories) => this.categoryPool.set(categories));
    this.tagService
      .getAllTags()
      .pipe(take(1), catchError(() => of([] as Tag[])), takeUntilDestroyed(this.destroyRef))
      .subscribe((tags) => this.tagPool.set(tags));
  }

  public onGrantsChange(groupId: number, selection: GrantSelection): void {
    this.editedGrants.set(groupId, selection);
  }

  private listenToAppRoleChanges(): void {
    const control = this.form.get("appRoleId");
    control?.valueChanges
      .pipe(startWith(control.value), takeUntilDestroyed(this.destroyRef))
      .subscribe((id) =>
        this.selectedAppRoleId.set(id == null ? null : Number(id))
      );
  }

  private listenToIsDummyChanges(): void {
    this.form
      .get("isDummyUser")
      ?.valueChanges.pipe(
      startWith(this.form.get("isDummyUser")?.value),
      tap((isDummyUser: boolean) => {
        const passwordField = this.form.get("password");
        if (isDummyUser) {
          passwordField?.removeValidators(Validators.required);
          passwordField?.setValue("");
          passwordField?.disable();
        } else {
          passwordField?.setValue("");
          passwordField?.enable();
          passwordField?.addValidators(Validators.required);
        }
      })
    )
      .subscribe();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      displayName: [this.user?.displayName ?? "", Validators.required],
      username: [
        this.user?.username ?? "",
        Validators.required,
        this.userValidators.uniqueUsername(0, this.user?.username ?? ""),
      ],
      appRoleId: [this.user?.appRoleId ?? null, Validators.required],
      isDummyUser: [false],
    });

    if (!this.user) {
      this.form.addControl(
        "password",
        new FormControl("", Validators.required)
      );
    } else {
      this.form.get("isDummyUser")?.disable();
    }
  }

  public previewRole(): void {
    const role = this.selectedRole();
    if (role) {
      openRolePreviewDialog(this.matDialog, role);
    }
  }

  public submit(): void {
    if (this.form.valid && this.user) {
      this.userService
        .updateUserById(this.user.id, this.form.value)
        .pipe(
          take(1),
          switchMap(() =>
            this.store.dispatch(
              new UpdateUser(this.user?.id.toString() as string, {
                ...this.user,
                ...this.form.value,
              })
            )
          ),
          // Grants go through their own endpoint (and their own permission), so
          // they are written after the user update rather than as part of it.
          // Only changed groups are touched.
          switchMap(() =>
            saveChangedMemberGrants(
              this.groupsService,
              this.user!.id,
              this.grantRows(),
              this.editedGrants
            )
          ),
          // Reported only once BOTH writes have landed. The grants are a separate
          // request that can fail on its own (a 400 ceiling violation, a 403), so
          // announcing success any earlier claims an assignment that never saved.
          // It stays ahead of the token refresh below, which is session
          // housekeeping rather than something this message speaks for.
          tap(() => {
            this.snackbarService.success("User successfully updated");
          }),
          switchMap(() =>
            iif(
              () =>
                this.store.selectSnapshot(AuthState.loggedInUser).id ===
                this.user?.id,
              defer(() => this.tokenRefreshService.refreshToken()),
              of(undefined)
            )
          ),
          tap(() => this.matDialogRef.close(true)),
          catchError(() => {
            // The interceptor surfaces the failure; keep the dialog open so the
            // admin can correct the input rather than losing it.
            return of(undefined);
          })
        )
        .subscribe();
    } else if (this.form.valid && !this.user) {
      this.userService
        .createUser(this.form.value)
        .pipe(
          take(1),
          switchMap((u) => this.store.dispatch(new AddUser(u))),
          tap(() => {
            this.snackbarService.success("User successfully created");
            this.matDialogRef.close(true);
          }),
          catchError((err) => {
            // Closing here would report a create that never happened and discard
            // what the admin typed, so the dialog stays open. The username
            // conflict comes back as a field message the interceptor ignores, so
            // it is surfaced here; anything else it already toasted itself.
            const message = err?.error?.["username"] ?? err?.["errMsg"];
            if (message) {
              this.snackbarService.error(message);
            }
            return of(undefined);
          })
        )
        .subscribe();
    }
  }

  public closeModal(): void {
    this.matDialogRef.close(false);
  }
}
