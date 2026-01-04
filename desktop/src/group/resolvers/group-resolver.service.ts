import { inject } from "@angular/core";
import { ActivatedRouteSnapshot, RouterStateSnapshot } from "@angular/router";
import { tap } from "rxjs";
import { setEntityHeaderText } from "src/utils";
import { GroupsService } from "../../open-api";

export const groupResolverFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const groupsService = inject(GroupsService);
  return groupsService
    .getGroupById(route.params["id"] || route.parent?.params["id"])
    .pipe(
      tap((group) => {
        if (route.data["setHeaderText"] && route.data["formConfig"]) {
          route.data["formConfig"].headerText = setEntityHeaderText(
            group,
            "name",
            route.data["formConfig"],
            route.data["entityType"]
          );
        }
      })
    );
};
