import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatIconModule } from "@angular/material/icon";
import { MatTooltipModule } from "@angular/material/tooltip";
import { RouterModule } from "@angular/router";
import { ButtonModule } from "../button";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { RoleListComponent } from "./role-list/role-list.component";
import { RolesRoutingModule } from "./roles-routing.module";

@NgModule({
  declarations: [RoleListComponent],
  imports: [
    ButtonModule,
    CommonModule,
    MatIconModule,
    MatTooltipModule,
    RouterModule,
    SharedUiModule,
    RolesRoutingModule,
  ],
})
export class RolesModule {}
