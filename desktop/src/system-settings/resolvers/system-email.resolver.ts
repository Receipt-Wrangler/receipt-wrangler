import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { take, tap } from "rxjs";
import { SystemEmail, SystemEmailService } from "../../open-api";
import { setEntityHeaderText } from "../../utils";

export const systemEmailResolver: ResolveFn<SystemEmail> = (route, state) => {
  const id = route.params["id"];
  return inject(SystemEmailService).getSystemEmailById(id).pipe(
    take(1),
    tap((systemEmail) => {
      if (route.data["setHeaderText"] && route.data["formConfig"]) {
        route.data["formConfig"].headerText = setEntityHeaderText(
          systemEmail,
          "username",
          route.data["formConfig"],
          "System Email"
        );
      }
    })
  );
};
