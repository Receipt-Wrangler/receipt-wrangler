import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Store } from "@ngxs/store";
import { map } from "rxjs";
import { PagedRequestCommand, Prompt, PromptService, UserRole } from "../open-api";
import { AuthState } from "../store";

export const promptsResolver: ResolveFn<Prompt[]> = (route, state) => {
  const store = inject(Store);
  const isAdmin = store.selectSnapshot(AuthState.hasRole(UserRole.Admin));

  if (isAdmin) {
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
