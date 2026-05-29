import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from "@angular/material/icon";
import { MatTooltipModule } from "@angular/material/tooltip";
import { RouterModule } from "@angular/router";
import { PipesModule } from "src/pipes/pipes.module";
import { ButtonModule } from "../button";
import { InputModule } from "../input";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TextareaModule } from "../textarea/textarea.module";
import { RoleFormComponent } from "./role-form/role-form.component";
import { RoleListComponent } from "./role-list/role-list.component";
import { RolesRoutingModule } from "./roles-routing.module";

@NgModule({
  declarations: [RoleListComponent, RoleFormComponent],
  imports: [
    ButtonModule,
    CommonModule,
    InputModule,
    MatIconModule,
    MatTooltipModule,
    PipesModule,
    ReactiveFormsModule,
    RouterModule,
    SharedUiModule,
    TextareaModule,
    RolesRoutingModule,
  ],
})
export class RolesModule {}
