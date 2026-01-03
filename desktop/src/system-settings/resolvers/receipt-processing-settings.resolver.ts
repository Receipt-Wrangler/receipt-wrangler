import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { map, take } from "rxjs";
import { PagedRequestCommand, ReceiptProcessingSettings, ReceiptProcessingSettingsService } from "../../open-api";

export const allReceiptProcessingSettingsResolver: ResolveFn<ReceiptProcessingSettings[]> = (route, state) => {
  const service = inject(ReceiptProcessingSettingsService);
  const command: PagedRequestCommand = {
    page: 1,
    pageSize: -1,
    orderBy: "name",
    sortDirection: "asc",
  };

  return service.getPagedProcessingSettings(command)
    .pipe(
      take(1),
      map((response) => response.data as any as ReceiptProcessingSettings[])
    );
};
