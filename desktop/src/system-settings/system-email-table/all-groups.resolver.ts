import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { map, take } from "rxjs";
import { AssociatedGroup, Group, GroupsService, PagedGroupRequestCommand } from "../../open-api";

export const allGroupsResolver: ResolveFn<Group[]> = (route, state) => {
  const groupService = inject(GroupsService);
  const command: PagedGroupRequestCommand = {
    page: 1,
    pageSize: -1,
    orderBy: "name",
    sortDirection: "asc",
    filter: {
      associatedGroup: AssociatedGroup.All,
    }
  };

  return groupService.getPagedGroups(command).pipe(
    take(1),
    map((response) => response.data as any as Group[])
  );
};
