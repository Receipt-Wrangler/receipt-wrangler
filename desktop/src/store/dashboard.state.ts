import { Injectable } from "@angular/core";
import { Action, createSelector, Selector, State, StateContext, } from "@ngxs/store";
import { take, tap } from "rxjs";
import { DashboardStateInterface } from "src/interfaces/dashboard-state.interface";
import { Dashboard, DashboardService } from "../open-api";
import { AddDashboardToGroup, DeleteDashboardFromGroup, SetDashboardsForGroup, UpdateDashBoardForGroup, } from "./dashboard.state.actions";

@State<DashboardStateInterface>({
  name: "dashboards",
  defaults: {
    dashboards: {},
  },
})
@Injectable()
export class DashboardState {
  constructor(private dashboardService: DashboardService) {}

  @Selector()
  static dashboards(state: DashboardStateInterface): {
    [groupId: string]: Dashboard[];
  } {
    return state.dashboards;
  }

  static getDashboardsByGroupId(groupId: string) {
    return createSelector(
      [DashboardState],
      (state: DashboardStateInterface) => {
        return state.dashboards[groupId] || [];
      }
    );
  }

  @Action(SetDashboardsForGroup)
  setDashboardsForGroup(
    { patchState, getState }: StateContext<DashboardStateInterface>,
    payload: SetDashboardsForGroup
  ) {
    // Return the observable so NGXS awaits the HTTP request before completing the
    // dispatch. Previously this subscribed internally and returned a Subscription,
    // so the dispatch (and the route resolver that relies on it) completed before
    // the dashboards loaded — the root cause of dashboards not rendering on a
    // group switch.
    return this.dashboardService
      .getDashboardsForUserByGroupId(payload.groupId)
      .pipe(
        take(1),
        tap((dashboards) => {
          const newDashboards = Object.assign({}, getState().dashboards);
          newDashboards[payload.groupId] = dashboards;
          patchState({ dashboards: newDashboards });
        })
      );
  }

  @Action(AddDashboardToGroup)
  addDashboardToGroup(
    { patchState, getState }: StateContext<DashboardStateInterface>,
    payload: AddDashboardToGroup
  ) {
    const newDashboards = Object.assign({}, getState().dashboards);
    newDashboards[payload.groupId] = [
      ...newDashboards[payload.groupId],
      payload.dashboard,
    ];
    patchState({ dashboards: newDashboards });
  }

  @Action(UpdateDashBoardForGroup)
  updateDashBoardForGroup(
    { patchState, getState }: StateContext<DashboardStateInterface>,
    payload: UpdateDashBoardForGroup
  ) {
    const newDashboards = Object.assign({}, getState().dashboards);
    const dashboardIndex = newDashboards[payload.groupId].findIndex(
      (dashboard) => {
        return dashboard.id === payload.dashboardId;
      }
    );

    if (dashboardIndex >= 0) {
      newDashboards[payload.groupId][dashboardIndex] = payload.dashboard;
      patchState({ dashboards: newDashboards });
    }
  }

  @Action(DeleteDashboardFromGroup)
  deleteDashboardFromGroup(
    { patchState, getState }: StateContext<DashboardStateInterface>,
    payload: DeleteDashboardFromGroup
  ) {
    const groupDashboards = getState().dashboards?.[payload.groupId] || [];
    const newDashboards = Object.assign({}, getState().dashboards);

    newDashboards[payload.groupId] = groupDashboards.filter(
      (dashboard) => dashboard.id !== payload.dashboardId
    );

    patchState({ dashboards: newDashboards });
  }
}
