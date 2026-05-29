import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { RoleListComponent } from "./role-list/role-list.component";

const routes: Routes = [
  {
    path: "",
    component: RoleListComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class RolesRoutingModule {}
