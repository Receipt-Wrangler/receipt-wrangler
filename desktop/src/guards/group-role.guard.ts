import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot, UrlTree } from "@angular/router";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { GroupUtil } from "src/utils";
import { GroupRole } from "../open-api";
import { GroupState } from "../store";

@Injectable({
  providedIn: "root",
})
export class GroupRoleGuard {
  constructor(
    private store: Store,
    private groupUtil: GroupUtil,
    private router: Router
  ) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ):
    | Observable<boolean | UrlTree>
    | Promise<boolean | UrlTree>
    | boolean
    | UrlTree {
    let hasAccess = false;
    let groupId: number | undefined = undefined;
    const groupRole = route.data["groupRole"] as GroupRole;
    const allowAdminOverride = route.data["allowAdminOverride"] as boolean;

    const useRouteId = route.data["useRouteGroupId"];
    if (useRouteId) {
      groupId = Number.parseInt(route?.params?.["id"] || route?.parent?.params["id"]);
    } else {
      groupId = Number.parseInt(
        this.store.selectSnapshot(GroupState.selectedGroupId)
      );
    }

    hasAccess = this.groupUtil.hasGroupAccess(groupId, groupRole, allowAdminOverride, true);

    if (!hasAccess) {
      this.router.navigate([
        this.store.selectSnapshot(GroupState.dashboardLink),
      ]);
    }

    return hasAccess;
  }
}
