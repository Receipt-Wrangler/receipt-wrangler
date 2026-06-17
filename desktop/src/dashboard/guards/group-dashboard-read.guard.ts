import { inject } from "@angular/core";
import { CanActivateFn, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { Permission } from "../../open-api";
import { AuthState } from "../../store";

/**
 * Allows the dashboard route only when the user holds group.dashboards.read for
 * the route's group. On deny, redirects to that group's receipt list so the
 * user lands on a page they can use, instead of letting the dashboard fetch
 * 403 (which the interceptor turns into a logout).
 *
 * The group id (for both the check and the redirect) comes from the route's
 * `:groupId` param so it is correct regardless of guard ordering or the
 * currently selected group.
 */
export const groupDashboardReadGuard: CanActivateFn = (route) => {
  const store = inject(Store);
  const router = inject(Router);
  const groupId = route.params["groupId"];

  if (
    store.selectSnapshot(
      AuthState.hasGroupPermission(+groupId, Permission.GroupDashboardsRead)
    )
  ) {
    return true;
  }

  router.navigate(["/receipts/group", groupId]);
  return false;
};
