import {Component, OnInit} from "@angular/core";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {MatDialog} from "@angular/material/dialog";
import {ActivatedRoute, Router} from "@angular/router";
import {Store} from "@ngxs/store";
import {catchError, of, switchMap, take, tap} from "rxjs";
import {DEFAULT_DIALOG_CONFIG} from "src/constants/dialog.constant";
import {FormMode} from "src/enums/form-mode.enum";
import {FormConfig} from "src/interfaces";
import {Permission, User, UserService} from "../../open-api";
import {ClaimsService, SnackbarService, TokenRefreshService} from "../../services";
import {AuthState, Logout, UpdateUser} from "../../store";
import {DeleteAccountDialogComponent} from "../delete-account-dialog/delete-account-dialog.component";

@Component({
    selector: "app-user-profile",
    templateUrl: "./user-profile.component.html",
    styleUrls: ["./user-profile.component.scss"],
    standalone: false
})
export class UserProfileComponent implements OnInit {
  public form: FormGroup = new FormGroup({});

  public user!: User;

  public formConfig!: FormConfig;

  public formMode = FormMode;

  public usernameTooltip: string =
    "Only system admin may change your username.";

  protected readonly Permission = Permission;

  protected readonly canEdit = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppAccountUpdate)
  );

  constructor(
    private claimsService: ClaimsService,
    private formBuilder: FormBuilder,
    private matDialog: MatDialog,
    private route: ActivatedRoute,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
    private tokenRefreshService: TokenRefreshService,
    private userService: UserService,
  ) {
  }

  public ngOnInit(): void {
    this.user = this.store.selectSnapshot(AuthState.loggedInUser);
    this.formConfig = this.route?.snapshot?.data?.["formConfig"];
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      username: this.user?.username ?? "",
      displayName: [this.user?.displayName ?? "", Validators.required],
      defaultAvatarColor: [
        this.user?.defaultAvatarColor ?? "",
        Validators.pattern("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"),
      ],
    });

    if (this.formConfig.mode === FormMode.edit) {
      this.form.get("username")?.disable();
    }
  }

  public submit(): void {
    if (this.form.valid) {
      this.userService
        .updateUserProfile(this.form.value)
        .pipe(
          take(1),
          switchMap(() => this.tokenRefreshService.refreshToken()),
          switchMap(() => this.claimsService.getAndSetClaimsForLoggedInUser()),
          switchMap(() => {
            const loggedInUser = this.store.selectSnapshot(
              AuthState.loggedInUser
            );
            return this.store.dispatch(
              new UpdateUser(loggedInUser.id.toString(), loggedInUser)
            );
          }),
          tap(() => {
            this.snackbarService.success("User profile successfully updated");
            this.router.navigate(["/settings/user-profile/view"],
              {
                queryParams: {
                  tab: "user-profile",
                }
              });
          })
        )
        .subscribe();
    }
  }

  public deleteAccount(): void {
    const dialogRef = this.matDialog.open(
      DeleteAccountDialogComponent,
      DEFAULT_DIALOG_CONFIG,
    );

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        switchMap((password) => {
          if (password) {
            return this.userService.deleteAccount({ password }).pipe(
              switchMap(() => this.store.dispatch(new Logout())),
              tap(() => {
                this.snackbarService.success(
                  "Your account has been successfully deleted"
                );
                this.router.navigate(["/"]);
              }),
              catchError((err) => {
                if (err.status === 401) {
                  this.deleteAccount();
                }
                return of(undefined);
              })
            );
          }
          return of(undefined);
        })
      )
      .subscribe();
  }
}
