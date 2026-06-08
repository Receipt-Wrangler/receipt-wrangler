import { Component, DestroyRef, Input, OnInit, computed, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { catchError, defer, iif, of, startWith, switchMap, take, tap, } from "rxjs";
import { UserValidators } from "src/validators/user-validators";
import { PermissionScope, Role, RoleService, User, UserService } from "../../open-api";
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
    private roleService: RoleService,
    private matDialog: MatDialog,
    private destroyRef: DestroyRef,
    public matDialogRef: MatDialogRef<UserFormComponent>
  ) {}

  public form: FormGroup = new FormGroup({});

  // The full app-role list, the selector options and the currently-selected role
  // (for the preview button) are derived reactively so they update once the roles
  // load asynchronously.
  private readonly roles = signal<Role[]>([]);
  public readonly appRoleOptions = computed<Role[]>(() =>
    this.roles().filter((role) => role.scope === PermissionScope.App)
  );
  private readonly selectedAppRoleId = signal<number | null>(null);
  public readonly selectedRole = computed<Role | null>(
    () => this.roles().find((role) => role.id === this.selectedAppRoleId()) ?? null
  );

  public ngOnInit(): void {
    this.initForm();
    this.loadRoles();
    this.listenToAppRoleChanges();
    if (!this.user) {
      this.listenToIsDummyChanges();
    }
  }

  private loadRoles(): void {
    this.roleService
      .getRoles()
      // Degrade gracefully if the roles can't be read: leave the selector empty
      // rather than erroring.
      .pipe(
        take(1),
        catchError(() => of([] as Role[])),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe((roles) => {
        this.roles.set(roles);
        // New users default to the configured default app role; an existing
        // user's role is already pre-filled from the loaded user.
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
          tap(() => {
            this.snackbarService.success("User successfully updated");
          }),
          switchMap(() =>
            this.store.dispatch(
              new UpdateUser(this.user?.id.toString() as string, {
                ...this.user,
                ...this.form.value,
              })
            )
          ),
          switchMap(() =>
            iif(
              () =>
                this.store.selectSnapshot(AuthState.loggedInUser).id ===
                this.user?.id,
              defer(() => this.tokenRefreshService.refreshToken()),
              of(undefined)
            )
          ),
          tap(() => this.matDialogRef.close(true))
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
          }),
          catchError((err) => {
            return of(
              this.snackbarService.error(err.error["username"] ?? err["errMsg"])
            );
          }),
          tap(() => this.matDialogRef.close(true))
        )
        .subscribe();
    }
  }

  public closeModal(): void {
    this.matDialogRef.close(false);
  }
}
