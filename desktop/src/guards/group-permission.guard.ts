import { inject } from "@angular/core";
import { CanActivateFn, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { AuthState, GroupState } from "../store";

/**
 * Allows activation when the user holds the group-scoped permission declared on
 * `route.data['groupPermission']` for the resolved group (or one of the
 * optional `route.data['orAppPermissions']` app-scoped fallbacks). On deny,
 * redirects to the group dashboard.
 *
 * Group id resolution: when `route.data['useRouteGroupId']` is set the id comes
 * from the `:id` route param (own or parent), otherwise from the currently
 * selected group in `GroupState`.
 */
export const groupPermissionGuard: CanActivateFn = (route) => {
  const store: Store = inject(Store);
  const router: Router = inject(Router);

  const permission = route.data["groupPermission"] as string;
  const orApp = (route.data["orAppPermissions"] as string[]) ?? [];

  let groupId: number;
  if (route.data["useRouteGroupId"]) {
    groupId = Number.parseInt(
      route?.params?.["id"] ?? route?.parent?.params["id"]
    );
  } else {
    groupId = Number.parseInt(
      store.selectSnapshot(GroupState.selectedGroupId)
    );
  }

  if (
    store.selectSnapshot(
      AuthState.hasGroupPermission(groupId, permission, orApp)
    )
  ) {
    return true;
  }

  router.navigate([store.selectSnapshot(GroupState.dashboardLink)]);
  return false;
};
