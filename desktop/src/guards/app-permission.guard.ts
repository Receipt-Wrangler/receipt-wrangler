import { inject } from "@angular/core";
import { CanActivateFn, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { AuthState, GroupState } from "../store";

/**
 * Allows activation when the user holds ANY of the app-scoped permissions
 * declared on `route.data['appPermissions']`. On deny, redirects to the group
 * dashboard.
 */
export const appPermissionGuard: CanActivateFn = (route) => {
  const store: Store = inject(Store);
  const router: Router = inject(Router);

  const required = (route.data["appPermissions"] as string[]) ?? [];

  if (store.selectSnapshot(AuthState.hasAnyAppPermission(required))) {
    return true;
  }

  router.navigate([store.selectSnapshot(GroupState.dashboardLink)]);
  return false;
};
