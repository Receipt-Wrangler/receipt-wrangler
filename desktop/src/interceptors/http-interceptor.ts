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

      // The server's errorMsg is the user-friendly message and always wins when
      // present. Angular's generic HttpErrorResponse.message ("Http failure
      // response for /api/login/: 500 Internal Server Error") is only a fallback
      // for a 5xx that carries no message at all, so an infra-level failure is
      // not silent. These must stay mutually exclusive: MatSnackBar.open()
      // dismisses whatever is already showing, so toasting both replaced the
      // useful message with the useless one before the user could read it.
      if (!receiptQueueMode) {
        const serverMessage = e.error?.errorMsg;

        if (serverMessage) {
          snackbarService.error(serverMessage);
        } else if (e.status >= 500) {
          snackbarService.error(e.message);
        }
      }

      return throwError(() => e);
    })
  );
};
