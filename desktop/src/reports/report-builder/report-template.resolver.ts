import { inject } from "@angular/core";
import { ResolveFn, Router } from "@angular/router";
import { catchError, of } from "rxjs";
import { ReportTemplate } from "../../open-api";
import { SnackbarService } from "../../services";
import { ReportRunnerService } from "../services/report-runner.service";

/**
 * Loads the report template for the builder's edit route so the form can be built
 * from its stored configuration synchronously (before the component's field
 * initializer runs). A missing/failed id resolves to null after redirecting back to
 * the list, so the builder falls back to a blank "new report".
 */
export const reportTemplateResolver: ResolveFn<ReportTemplate | null> = (route) => {
  const runner = inject(ReportRunnerService);
  const router = inject(Router);
  const snackbar = inject(SnackbarService);

  const id = Number(route.paramMap.get("id"));
  return runner.getTemplate(id).pipe(
    catchError(() => {
      snackbar.error("Report template not found");
      router.navigate(["/reports"]);
      return of(null);
    }),
  );
};
