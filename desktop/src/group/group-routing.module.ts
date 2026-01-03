import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { FormMode } from "src/enums/form-mode.enum";
import { GroupRoleGuard } from "src/guards/group-role.guard";
import { FormConfig } from "src/interfaces/form-config.interface";
import { RoleGuard } from "../guards/role.guard";
import { GroupRole, UserRole } from "../open-api";
import { promptsResolver } from "../prompt/prompts.resolver";
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
    path: "",
    component: GroupTableComponent,
  },
  {
    path: "create",
    component: GroupFormComponent,
    data: {
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
      groupRole: GroupRole.Viewer,
    },
    canActivate: [GroupRoleGuard],
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
          groupRole: GroupRole.Viewer,
          entityType: "Details",
          setHeaderText: true,
          allowAdminOverride: true,
          useRouteGroupId: true,
        },
        canActivate: [GroupRoleGuard],
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
          groupRole: GroupRole.Owner,
          useRouteGroupId: true,
        },
        canActivate: [GroupRoleGuard],
      },
      {
        path: "receipt-settings/view",
        component: GroupReceiptSettingsComponent,
        resolve: {
          group: groupResolverFn,
        },
        data: {
          formConfig: {
            mode: FormMode.view,
            headerText: "View Group Receipt Settings",
          } as FormConfig,
          groupRole: GroupRole.Viewer,
          entityType: "Receipt Settings",
          setHeaderText: true,
          useRouteGroupId: true,
          allowAdminOverride: true,
        },
        canActivate: [GroupRoleGuard],
      },
      {
        path: "receipt-settings/edit",
        component: GroupReceiptSettingsComponent,
        resolve: {
          group: groupResolverFn,
        },
        data: {
          formConfig: {
            mode: FormMode.edit,
            headerText: "Edit Group Receipt Settings",
          } as FormConfig,
          groupRole: GroupRole.Owner,
          entityType: "Receipt Settings",
          setHeaderText: true,
          useRouteGroupId: true,
        },
        canActivate: [GroupRoleGuard],
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
          groupRole: GroupRole.Viewer,
          allowAdminOverride: true,
        },
        canActivate: [GroupRoleGuard],
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
          role: UserRole.Admin
        },
        canActivate: [RoleGuard],
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GroupRoutingModule {}
