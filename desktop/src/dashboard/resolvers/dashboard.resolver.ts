import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { SetDashboardsForGroup } from "src/store/dashboard.state.actions";
import { Dashboard } from "../../open-api";

export const dashboardResolverFn: ResolveFn<Observable<Dashboard[]>> = (
  route,
  state
): Observable<Dashboard[]> => {
  const store = inject(Store);
  const groupId = route.params["groupId"];

  return store.dispatch(new SetDashboardsForGroup(groupId)) as any as Observable<Dashboard[]>;
};
