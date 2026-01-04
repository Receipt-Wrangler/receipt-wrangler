import { Component, OnInit, ViewEncapsulation } from "@angular/core";
import { FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { BehaviorSubject, catchError, finalize, of, switchMap, tap, } from "rxjs";
import { AppData, AuthService } from "src/open-api";
import { SnackbarService } from "src/services";
import { setAppData } from "src/utils";
import { fadeIn, fadeInOut } from "../../animations";
import { GroupState } from "../../store";
import { UserValidators } from "../../validators";

@Component({
    selector: "app-auth-form",
    templateUrl: "./auth-form.component.html",
    styleUrls: ["./auth-form.component.scss"],
    encapsulation: ViewEncapsulation.None,
    providers: [UserValidators],
    animations: [fadeInOut, fadeIn],
    standalone: false
})
export class AuthForm implements OnInit {
  public form: FormGroup = new FormGroup({});
  public isSignUp: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(
    false
  );
  public headerText: string = "";
  public primaryButtonText: string = "";
  public secondaryButtonText: string = "";
  public secondaryButtonRouterLink: string[] = [];
  public isLoading = false;

  constructor(
    private authSerivce: AuthService,
    private snackbarService: SnackbarService,
    protected formBuilder: FormBuilder,
    protected route: ActivatedRoute,
    protected router: Router,
    protected store: Store,
    protected userValidators: UserValidators
  ) {
  }

  public ngOnInit(): void {
    this.initForm();
    this.listenForRouteChanges();
    this.listenForIsSignUpChanges();
  }

  private listenForRouteChanges(): void {
    this.route.data
      .pipe(
        tap((data) => {
          this.isSignUp.next(!!data?.["isSignUp"]);
        })
      )
      .subscribe();
  }

  private listenForIsSignUpChanges(): void {
    this.isSignUp
      .pipe(
        tap((isSignUp) => {
          if (isSignUp) {
            this.headerText = "Sign Up";
            this.primaryButtonText = "Sign Up";
            this.secondaryButtonRouterLink = ["/auth/login"];
            this.secondaryButtonText = "Back to Login";
            this.form
              .get("username")
              ?.addAsyncValidators(this.userValidators.uniqueUsername(0, ""));
            this.form.addControl(
              "displayname",
              new FormControl("", Validators.required)
            );
          } else {
            this.headerText = "Login";
            this.primaryButtonText = "Login";
            this.secondaryButtonRouterLink = ["/auth/sign-up"];
            this.secondaryButtonText = "Sign Up";
            this.form
              .get("username")
              ?.removeAsyncValidators(
                this.userValidators.uniqueUsername(0, "")
              );
            this.form.removeControl("displayname");
          }
        })
      )
      .subscribe();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      username: ["", [Validators.required]],
      password: ["", Validators.required],
    });
  }

  public submit(): void {
    const isSignUp = this.isSignUp.getValue();
    const isValid = this.form.valid;

    if (isValid && isSignUp) {
      this.signUp();
    } else if (isValid && !isSignUp) {
      this.login();
    }
  }

  private signUp(): void {
    this.authSerivce
      .signUp(this.form.value)
      .pipe(
        tap(() => {
          this.login();
        }),
        catchError((err) =>
          of(
            this.snackbarService.error(err.error["username"] ?? err["errMsg"])
          )
        )
      )
      .subscribe();
  }

  private login(): void {
    this.isLoading = true;
    this.authSerivce
      .login(this.form.value)
      .pipe(
        switchMap((appData: AppData) => setAppData(this.store, appData)),
        tap(() =>
          this.router.navigate([
            this.store.selectSnapshot(GroupState.dashboardLink),
          ]),
        ),
        finalize(() => this.isLoading = false)
      )
      .subscribe();
  }
}
