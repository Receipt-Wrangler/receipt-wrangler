import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { FormMode } from "src/enums/form-mode.enum";
import { appPermissionGuard } from "src/guards/app-permission.guard";
import { groupPermissionGuard } from "src/guards/group-permission.guard";
import { FormConfig } from "src/interfaces/form-config.interface";
import { Permission } from "../open-api";
import { promptsResolver } from "../prompt/prompts.resolver";
import { customFieldResolverFn } from "../resolvers/custom-field.resolver";
import { GroupDetailsComponent } from "./group-details/group-details.component";
import { GroupFormComponent } from "./group-form/group-form.component";
import { GroupReceiptSettingsComponent } from "./group-receipt-settings/group-receipt-settings.component";
import { GroupSettingsComponent } from "./group-settings/group-settings.component";
import { GroupTableComponent } from "./group-table/group-table.component";
import { GroupTabsComponent } from "./group-tabs/group-tabs.component";
import { groupResolverFn } from "./resolvers/group-resolver.service";
import { systemEmailsResolver } from "./resolvers/system-emails.resolver";

const routes: Routes = [
  {
    // No app-permission guard: the groups table shows the caller's OWN groups
    // (backend GET /group/ -> GetGroupsForUser is auth-only), exactly like
    // legacy. The "all groups" filter button is gated on app.groups.read inside
    // group-table.component.html, and per-row Edit/Delete on group.update /
    // group.delete — so a non-admin sees their own groups without the admin
    // all-groups view. The parent shell route already enforces AuthGuard.
    path: "",
    component: GroupTableComponent,
  },
  {
    path: "create",
    component: GroupFormComponent,
    canActivate: [appPermissionGuard],
    data: {
      appPermissions: [Permission.AppGroupsCreate],
      formConfig: {
        mode: FormMode.add,
        headerText: "Create Group",
      } as FormConfig,
    },
  },
  {
    path: ":id",
    component: GroupTabsComponent,
    resolve: {
      group: groupResolverFn,
    },
    data: {
      formConfig: {
        mode: FormMode.view,
        headerText: "View Group",
      } as FormConfig,
      groupPermission: Permission.GroupView,
      orAppPermissions: [Permission.AppGroupsRead],
      useRouteGroupId: true,
    },
    canActivate: [groupPermissionGuard],
    children: [
      {
        path: "details/view",
        component: GroupDetailsComponent,
        resolve: {
          group: groupResolverFn,
        },
        data: {
          formConfig: {
            mode: FormMode.view,
            headerText: "View Group",
          } as FormConfig,
          entityType: "Details",
          setHeaderText: true,
          groupPermission: Permission.GroupView,
          orAppPermissions: [Permission.AppGroupsRead],
          useRouteGroupId: true,
        },
        canActivate: [groupPermissionGuard],
      },
      {
        path: "details/edit",
        component: GroupDetailsComponent,
        resolve: {
          group: groupResolverFn,
        },
        data: {
          formConfig: {
            mode: FormMode.edit,
          } as FormConfig,
          entityType: "Details",
          setHeaderText: true,
          groupPermission: Permission.GroupUpdate,
          useRouteGroupId: true,
        },
        canActivate: [groupPermissionGuard],
      },
      {
        path: "receipt-settings/view",
        component: GroupReceiptSettingsComponent,
        resolve: {
          group: groupResolverFn,
          customFields: customFieldResolverFn,
        },
        data: {
          formConfig: {
            mode: FormMode.view,
            headerText: "View Group Receipt Settings",
          } as FormConfig,
          groupPermission: Permission.GroupView,
          orAppPermissions: [Permission.AppGroupsRead],
          entityType: "Receipt Settings",
          setHeaderText: true,
          useRouteGroupId: true,
        },
        canActivate: [groupPermissionGuard],
      },
      {
        path: "receipt-settings/edit",
        component: GroupReceiptSettingsComponent,
        resolve: {
          group: groupResolverFn,
          customFields: customFieldResolverFn,
        },
        data: {
          formConfig: {
            mode: FormMode.edit,
            headerText: "Edit Group Receipt Settings",
          } as FormConfig,
          groupPermission: Permission.GroupUpdate,
          entityType: "Receipt Settings",
          setHeaderText: true,
          useRouteGroupId: true,
        },
        canActivate: [groupPermissionGuard],
      },
      {
        path: "settings/view",
        component: GroupSettingsComponent,
        resolve: {
          group: groupResolverFn,
          systemEmails: systemEmailsResolver,
          prompts: promptsResolver,
        },
        data: {
          formConfig: {
            mode: FormMode.view,
          } as FormConfig,
          setHeaderText: true,
          entityType: "Settings",
          appPermissions: [Permission.AppGroupsUpdateSettings],
        },
        canActivate: [appPermissionGuard],
      },
      {
        path: "settings/edit",
        component: GroupSettingsComponent,
        resolve: {
          group: groupResolverFn,
          systemEmails: systemEmailsResolver,
          prompts: promptsResolver
        },
        data: {
          formConfig: {
            mode: FormMode.edit,
          } as FormConfig,
          setHeaderText: true,
          entityType: "Settings",
          appPermissions: [Permission.AppGroupsUpdateSettings],
        },
        canActivate: [appPermissionGuard],
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GroupRoutingModule {}
