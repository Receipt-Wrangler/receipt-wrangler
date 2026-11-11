import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { Routes } from "@angular/router";
import { appPermissionGuard } from "src/guards/app-permission.guard";
import { AuthGuard } from "src/guards/auth.guard";
import { SidebarComponent } from "src/layout/sidebar/sidebar.component";
import { Permission } from "../open-api";

// set up dashboard
export const appRoutes: Routes = [
  {
    path: "",
    component: SidebarComponent,
    children: [
      {
        path: "auth",
        loadChildren: () =>
          import("../auth/auth.routes").then((m) => m.authRoutes),
        canActivate: [AuthGuard],
      },
      {
        path: "dashboard",
        loadChildren: () =>
          import("../dashboard/dashboard.routes").then((m) => m.dashboardRoutes),
        canActivate: [AuthGuard],
      },
      {
        path: "categories",
        loadChildren: () =>
          import("../categories/categories.routes").then((m) => m.categoriesRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [Permission.AppCategoriesRead],
        },
      },
      {
        path: "tags",
        // Carried over unchanged from the former TagsModule: a route-scoped
        // HttpClient that does not run the root httpInterceptor.
        providers: [provideHttpClient(withInterceptorsFromDi())],
        loadChildren: () =>
          import("../tags/tags.routes").then((m) => m.tagsRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [Permission.AppTagsRead],
        },
      },
      {
        path: "custom-fields",
        loadChildren: () =>
          import("../custom-fields/custom-fields.routes").then((m) => m.customFieldsRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [Permission.AppCustomFieldsRead],
        },
      },
      {
        path: "groups",
        loadChildren: () =>
          import("../group/group.routes").then((m) => m.groupRoutes),
        canActivate: [AuthGuard],
      },
      {
        path: "receipts",
        loadChildren: () =>
          import("../receipts/receipts.routes").then((m) => m.receiptsRoutes),
        canActivate: [AuthGuard],
      },
      {
        path: "settings",
        loadChildren: () =>
          import("../settings/settings.routes").then((m) => m.settingsRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [
            Permission.AppAccountRead,
            Permission.AppUserPreferencesRead,
            Permission.AppApiKeysRead,
          ],
        },
      },
      {
        path: "system-settings",
        loadChildren: () =>
          import("../system-settings/system-settings.routes").then((m) => m.systemSettingsRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [
            Permission.AppSystemSettingsRead,
            Permission.AppPromptsRead,
            Permission.AppReceiptProcessingSettingsRead,
            Permission.AppSystemEmailsRead,
            Permission.AppOidcProvidersRead,
            Permission.AppSystemTasksRead,
          ],
        },
      },
      {
        path: "users",
        loadChildren: () =>
          import("../user/user.routes").then((m) => m.userRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [Permission.AppUsersRead],
        },
      },
      {
        path: "roles",
        loadChildren: () =>
          import("../roles/roles.routes").then((m) => m.rolesRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          appPermissions: [Permission.AppRolesRead],
        },
      },
      {
        path: "reports",
        loadChildren: () =>
          import("../reports/reports.routes").then((m) => m.reportsRoutes),
        canActivate: [AuthGuard, appPermissionGuard],
        data: {
          // read OR readAll: an admin with only the "*All" bypass still reaches reports.
          appPermissions: [Permission.AppReportsRead, Permission.AppReportsReadAll],
        },
      },
      {
        path: "",
        redirectTo: "/auth/login",
        pathMatch: "full",
      },
    ],
  },
];
