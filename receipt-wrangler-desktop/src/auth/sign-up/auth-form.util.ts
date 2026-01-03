import { Injectable } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { Store } from "@ngxs/store";
import { catchError, map, Observable, of, switchMap, tap } from "rxjs";
import { AppData, AuthService } from "../../open-api";
import { SnackbarService } from "../../services";
import { setAppData } from "../../utils";

@Injectable()
export class AuthFormUtil {
  constructor(
    private authService: AuthService,
    private snackbarService: SnackbarService,
    private store: Store
  ) {}

  public getSubmitObservable(
    form: FormGroup,
    isSignUp: boolean
  ): Observable<void> {
    const isValid = form.valid;

    if (isValid && isSignUp) {
      return this.authService.signUp(form.value).pipe(
        tap(() => {
          this.snackbarService.success("User successfully signed up");
        }),
        catchError((err) =>
          of(this.snackbarService.error(err.error["username"] ?? err["errMsg"]))
        )
      );
    } else if (isValid && !isSignUp) {
      return this.authService.login(form.value).pipe(
        tap(() => {
          this.snackbarService.success("Successfully logged in");
        }),
        switchMap((appData: AppData) => setAppData(this.store, appData)),
        map(() => undefined)
      );
    } else {
      return of(undefined);
    }
  }
}
