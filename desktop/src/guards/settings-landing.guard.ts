import { inject } from "@angular/core";
import { CanActivateFn, Router, UrlTree } from "@angular/router";
import { Store } from "@ngxs/store";
import { Permission } from "../open-api";
import { AuthState, GroupState } from "../store";

/**
 * Landing guard for the Settings shell. The avatar menu links to `/settings`;
 * this redirects to the first tab the user can read, in the same order the tab
 * bar renders them. Falls back to the group dashboard when the user can read
 * none of them.
 */
export const settingsLandingGuard: CanActivateFn = (): UrlTree => {
  const store: Store = inject(Store);
  const router: Router = inject(Router);

  const tabs: { path: string; permission: Permission }[] = [
    { path: "user-profile/view", permission: Permission.AppAccountRead },
    {
      path: "user-preferences/view",
      permission: Permission.AppUserPreferencesRead,
    },
    { path: "api-keys/view", permission: Permission.AppApiKeysRead },
  ];

  for (const tab of tabs) {
    if (store.selectSnapshot(AuthState.hasAppPermission(tab.permission))) {
      return router.parseUrl(`/settings/${tab.path}`);
    }
  }

  return router.parseUrl(store.selectSnapshot(GroupState.dashboardLink));
};
