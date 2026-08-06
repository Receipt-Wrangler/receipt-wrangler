import { Component, effect, OnInit, signal, ViewEncapsulation } from "@angular/core";
import { FormBuilder, FormControl, FormGroup, Validators, } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import * as QRCode from "qrcode";
import { BehaviorSubject, catchError, finalize, of, switchMap, tap, } from "rxjs";
import { AppData, AuthService } from "src/open-api";
import { SnackbarService } from "src/services";
import { setAppData } from "src/utils";
import { fadeIn, fadeInOut } from "../../animations";
import { FeatureConfigState, GroupState } from "../../store";
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
  // Rendered locally from featureConfig.loginQrUrl; null when the login QR is
  // disabled/unset, which hides the QR block on the login page.
  public qrDataUrl = signal<string | null>(null);

  constructor(
    private authSerivce: AuthService,
    private snackbarService: SnackbarService,
    protected formBuilder: FormBuilder,
    protected route: ActivatedRoute,
    protected router: Router,
    protected store: Store,
    protected userValidators: UserValidators
  ) {
    const loginQrUrl = this.store.selectSignal(FeatureConfigState.loginQrUrl);
    effect((onCleanup) => {
      const url = loginQrUrl();
      if (!url) {
        this.qrDataUrl.set(null);
        return;
      }

      // Generation runs async, so a URL change can leave an earlier call in
      // flight. Cleanup runs before the next execution (and on destroy), so
      // only the latest generation is allowed to write the signal — otherwise
      // a late-resolving older QR could overwrite the current one.
      let cancelled = false;
      onCleanup(() => (cancelled = true));

      QRCode.toDataURL(url, { margin: 2, width: 220 })
        .then((dataUrl) => {
          if (!cancelled) {
            this.qrDataUrl.set(dataUrl);
          }
        })
        .catch(() => {
          if (!cancelled) {
            this.qrDataUrl.set(null);
          }
        });
    });
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
