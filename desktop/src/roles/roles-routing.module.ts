import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { appPermissionGuard } from "../guards/app-permission.guard";
import { Permission } from "../open-api";
import { RoleFormComponent } from "./role-form/role-form.component";
import { RoleListComponent } from "./role-list/role-list.component";

const routes: Routes = [
  {
    path: "",
    component: RoleListComponent,
  },
  {
    path: "new",
    component: RoleFormComponent,
    canActivate: [appPermissionGuard],
    data: {
      appPermissions: [Permission.AppRolesCreate],
    },
  },
  {
    path: ":id/edit",
    component: RoleFormComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class RolesRoutingModule {}
