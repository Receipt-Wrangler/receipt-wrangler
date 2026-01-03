import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { take, tap } from "rxjs";
import { ReceiptProcessingSettings, ReceiptProcessingSettingsService } from "../open-api";
import { setEntityHeaderText } from "../utils";

export const receiptProcessingSettingsResolver: ResolveFn<ReceiptProcessingSettings> = (route, state) => {
  const service = inject(ReceiptProcessingSettingsService);

  return service.getReceiptProcessingSettingsById(route.params["id"])
    .pipe(
      take(1),
      tap((receiptProcessingSettings) => {
        if (route.data["setHeaderText"] && route.data["formConfig"]) {
          route.data["formConfig"].headerText = setEntityHeaderText(
            receiptProcessingSettings,
            "name",
            route.data["formConfig"],
            "Receipt Processing Settings"
          );
        }

      })
    );
};
