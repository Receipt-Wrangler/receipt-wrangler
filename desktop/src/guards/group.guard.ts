import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot, } from "@angular/router";
import { Store } from "@ngxs/store";
import { GroupState, SetSelectedDashboardId, SetSelectedGroupId } from "../store";

@Injectable({
  providedIn: "root",
})
export class GroupGuard {
  constructor(private store: Store, private router: Router) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean {
    const groupId = route.params["groupId"];
    const group = this.store.selectSnapshot(GroupState.getGroupById(groupId));

    if (group) {
      this.resetSelectedDashboardIfGroupDashboardChanged(
        groupId?.toString() ?? ""
      );
      return true;
    } else {
      const newGroupId = this.store.selectSnapshot(GroupState.groups)[0]?.id;
      const basePath = route.data["groupGuardBasePath"];

      this.resetSelectedDashboardIfGroupDashboardChanged(
        newGroupId?.toString() ?? ""
      );
      this.store.dispatch(new SetSelectedGroupId(newGroupId?.toString() ?? ""));
      this.router.navigate([`${basePath}/${newGroupId}`]);
      return false;
    }
  }

  private resetSelectedDashboardIfGroupDashboardChanged(groupId: string): void {
    if (groupId !== this.store.selectSnapshot(GroupState.selectedGroupId)) {
      this.store.dispatch(new SetSelectedGroupId(groupId));
      this.store.dispatch(new SetSelectedDashboardId(undefined));
    }
  }
}
