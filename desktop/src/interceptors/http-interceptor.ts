import { HttpErrorResponse, HttpInterceptorFn } from "@angular/common/http";
import { inject } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { catchError, throwError } from "rxjs";
import { SnackbarService } from "../services";
import { AuthState } from "../store";

export const httpInterceptor: HttpInterceptorFn = (req, next) => {
  const store = inject(Store);
  const router = inject(Router);
  const activatedRoute = inject(ActivatedRoute);
  const snackbarService = inject(SnackbarService);

  return next(req).pipe(
    catchError((e: HttpErrorResponse) => {
      const isLoggedIn = store.selectSnapshot(AuthState.isLoggedIn);
      if (!isLoggedIn && req.url.includes("token")) {
        return throwError(() => e);
      }

      if (e.status === 403 && isLoggedIn) {
        localStorage.clear();
        router.navigate(["/"]);
        return throwError(() => e);
      }

      const regex = new RegExp("5d{2}");
      const receiptQueueMode = activatedRoute.snapshot.queryParams["queueMode"];

      // NOTE: We check for queueMode to gracefully handle creating queues with mixed permissions
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
