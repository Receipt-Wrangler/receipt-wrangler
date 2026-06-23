import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Store } from "@ngxs/store";
import { Observable, of } from "rxjs";
import { SetDashboardsForGroup } from "src/store/dashboard.state.actions";
import { Dashboard, Permission } from "../../open-api";
import { AuthState } from "../../store";

export const dashboardResolverFn: ResolveFn<Observable<Dashboard[]>> = (
  route,
  state
): Observable<Dashboard[]> => {
  const store = inject(Store);
  const groupId = route.params["groupId"];

  // Without group dashboard read permission the API 403s (and the interceptor
  // logs the user out). Skip the fetch and resolve empty so the page renders
  // gracefully instead of booting the user.
  if (
    !store.selectSnapshot(
      AuthState.hasGroupPermission(+groupId, Permission.GroupDashboardsRead)
    )
  ) {
    return of([]);
  }

  return store.dispatch(new SetDashboardsForGroup(groupId)) as any as Observable<Dashboard[]>;
};
