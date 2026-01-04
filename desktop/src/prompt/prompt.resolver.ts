import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { take, tap } from "rxjs";
import { Prompt, PromptService } from "../open-api";
import { setEntityHeaderText } from "../utils";

export const promptResolver: ResolveFn<Prompt> = (route, state) => {
  const service = inject(PromptService);
  return service.getPromptById(route.params["id"]).pipe(
    take(1),
    tap((prompt) => {
      if (route.data["setHeaderText"] && route.data["formConfig"]) {
        route.data["formConfig"].headerText = setEntityHeaderText(
          prompt,
          "name",
          route.data["formConfig"],
          "Prompt"
        );
      }
    })
  );
};
