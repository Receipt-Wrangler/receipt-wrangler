import { Component, Input, OnInit } from "@angular/core";
import { FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { catchError, defer, iif, of, startWith, switchMap, take, tap, } from "rxjs";
import { USER_ROLE_OPTIONS } from "src/group/role-options";
import { FormOption } from "src/interfaces/form-option.interface";
import { UserValidators } from "src/validators/user-validators";
import { AuthService, User, UserService } from "../../open-api";
import { SnackbarService } from "../../services";
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
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private snackbarService: SnackbarService,
    private store: Store,
    private userService: UserService,
    private userValidators: UserValidators,
    public matDialogRef: MatDialogRef<UserFormComponent>
  ) {}

  public form: FormGroup = new FormGroup({});

  public userRoleOptions: FormOption[] = USER_ROLE_OPTIONS;

  public ngOnInit(): void {
    this.initForm();
    if (!this.user) {
      this.listenToIsDummyChanges();
    }
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
      userRole: [this.user?.userRole ?? "", Validators.required],
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
              defer(() => this.authService.getNewRefreshToken()),
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
