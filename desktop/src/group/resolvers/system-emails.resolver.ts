import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Store } from "@ngxs/store";
import { map, take } from "rxjs";
import { PagedRequestCommand, SystemEmail, SystemEmailService, UserRole } from "../../open-api";
import { AuthState } from "../../store";

export const systemEmailsResolver: ResolveFn<SystemEmail[]> = (route, state) => {
  const store = inject(Store);
  const isAdmin = store.selectSnapshot(AuthState.hasRole(UserRole.Admin));
  if (isAdmin) {
    const systemEmailService = inject(SystemEmailService);
    const command: PagedRequestCommand = {
      page: 1,
      pageSize: -1,
      sortDirection: "asc",
      orderBy: "username",
    };

    return systemEmailService.getPagedSystemEmails(command)
      .pipe(
        take(1),
        map((pagedData) => pagedData.data as SystemEmail[])
      );
  }

  return [];
};
