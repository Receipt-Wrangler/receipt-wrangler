import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Store } from "@ngxs/store";
import { map } from "rxjs";
import { PagedRequestCommand, Permission, Prompt, PromptService } from "../open-api";
import { AuthState } from "../store";

export const promptsResolver: ResolveFn<Prompt[]> = (route, state) => {
  const store = inject(Store);
  const canRead = store.selectSnapshot(AuthState.hasAppPermission(Permission.AppPromptsRead));

  if (canRead) {
    const promptService = inject(PromptService);
    const command: PagedRequestCommand = {
      page: 1,
      pageSize: -1,
      orderBy: "name",
      sortDirection: "asc",
    };
    return promptService.getPagedPrompts(command)
      .pipe(
        map((pagedData) => {
          return pagedData.data as any as Prompt[];
        })
      );
  }

  return [];
};
