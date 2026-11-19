import { Routes } from "@angular/router";
import { GroupGuard } from "src/guards/group.guard";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { GroupDashboardsComponent } from "./group-dashboards/group-dashboards.component";
import { dashboardGuard } from "./guards/dashboard.guard";
import { groupDashboardReadGuard } from "./guards/group-dashboard-read.guard";
import { dashboardResolverFn } from "./resolvers/dashboard.resolver";

export const dashboardRoutes: Routes = [
  {
    path: "group/:groupId",
    component: GroupDashboardsComponent,
    canActivate: [GroupGuard, groupDashboardReadGuard],
    data: {
      groupGuardBasePath: "/dashboard/group",
    },
    resolve: {
      dashboards: dashboardResolverFn,
    },
    children: [
      {
        path: ":dashboardId",
        component: DashboardComponent,
        canActivate: [dashboardGuard],
      },
    ],
  },
  {
    path: "",
    redirectTo: "group/:groupId",
    pathMatch: "full",
  },
];
