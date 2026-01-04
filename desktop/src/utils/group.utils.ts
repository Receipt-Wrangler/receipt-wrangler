import { Injectable } from "@angular/core";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { GroupRole, UserRole } from "../open-api";
import { AuthState, GroupState } from "../store";

@Injectable({
  providedIn: "root",
})
export class GroupUtil {
  constructor(private store: Store, private router: Router) {}

  // TODO: v5 Check frontend group role permissions to allow for admins to have access to all groups
  public hasGroupAccess(
    groupId: number | undefined,
    groupRole: GroupRole,
    allowAdminOverride: boolean,
    navigateOnFail = false,
  ): boolean {
    const roles = [GroupRole.Viewer, GroupRole.Editor, GroupRole.Owner];
    const requiredIndex = roles.findIndex((v) => v === groupRole) as number;
    if (allowAdminOverride) {
      const isAdmin = this.store.selectSnapshot(AuthState.hasRole(UserRole.Admin));
      if (isAdmin) {
        return true;
      }
    }

    if (!groupId) {
      return true;
    } else {
      let hasAccess = false;

      const userId = this.store.selectSnapshot(AuthState.userId);
      const group = this.store.selectSnapshot(
        GroupState.getGroupById(groupId.toString())
      );

      if (!group) {
        if (navigateOnFail) {
          this.router.navigate(["/"]);
        }
        return hasAccess;
      }

      const member = group.groupMembers.find(
        (m) => m.userId.toString() === userId.toString()
      );

      if (!member) {
        if (navigateOnFail) {
          this.router.navigate(["/"]);
        }
        return hasAccess;
      }

      const memberIndex = roles.findIndex((v) => v === member.groupRole);
      hasAccess = memberIndex >= requiredIndex;

      return hasAccess;
    }
  }
}
