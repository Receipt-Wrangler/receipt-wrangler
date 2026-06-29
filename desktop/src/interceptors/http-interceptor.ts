import { HttpErrorResponse, HttpInterceptorFn } from "@angular/common/http";
import { inject } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { catchError, throwError } from "rxjs";
import { SnackbarService } from "../services";
import { AuthState } from "../store";

const FORBIDDEN_MESSAGE = "You do not have permission to perform this action.";

export const httpInterceptor: HttpInterceptorFn = (req, next) => {
  const store = inject(Store);
  const activatedRoute = inject(ActivatedRoute);
  const snackbarService = inject(SnackbarService);

  return next(req).pipe(
    catchError((e: HttpErrorResponse) => {
      const isLoggedIn = store.selectSnapshot(AuthState.isLoggedIn);

      // Don't intercept errors from token refresh requests — let TokenRefreshService handle them
      if (req.url.includes("/api/token/")) {
        return throwError(() => e);
      }

      // NOTE: We check for queueMode to gracefully handle creating queues with mixed permissions
      const receiptQueueMode = activatedRoute.snapshot.queryParams["queueMode"];

      // The backend returns 403 for genuine permission denials too, not just auth. With a
      // still-valid token this is a permission denial, so surface it instead of logging the
      // user out. Token freshness is handled proactively elsewhere (15-min refresh timer,
      // app-init, and the auth guard).
      if (e.status === 403 && isLoggedIn) {
        // Toast only user-initiated actions; background GET reads are handled by their callers.
        if (req.method !== "GET" && !receiptQueueMode) {
          snackbarService.error(FORBIDDEN_MESSAGE);
        }
        return throwError(() => e);
      }

      const regex = new RegExp("5\\d{2}");

      if (e.error?.errorMsg && !receiptQueueMode) {
        snackbarService.error(e.error?.errorMsg);
      }

      if (regex.test(e.status.toString())) {
        snackbarService.error(e.message);
      }

      return throwError(() => e);
    })
  );
};
