import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from "@angular/material/icon";
import { MatTooltipModule } from "@angular/material/tooltip";
import { RouterModule } from "@angular/router";
import { CheckboxModule } from "src/checkbox/checkbox.module";
import { PipesModule } from "src/pipes/pipes.module";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { ButtonModule } from "../button";
import { CategoryAutocompleteComponent } from "../category-autocomplete/category-autocomplete.component";
import { GrantPickerComponent } from "../shared-ui/grant-picker/grant-picker.component";
import { DirectivesModule } from "../directives";
import { InputModule } from "../input";
import { TagAutocompleteComponent } from "../tag-autocomplete/tag-autocomplete.component";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { TextareaModule } from "../textarea/textarea.module";
import { RoleFormComponent } from "./role-form/role-form.component";
import { RoleListComponent } from "./role-list/role-list.component";
import { RolesRoutingModule } from "./roles-routing.module";

@NgModule({
  declarations: [RoleListComponent, RoleFormComponent],
  imports: [
    ButtonModule,
    CheckboxModule,
    CommonModule,
    DirectivesModule,
    InputModule,
    MatIconModule,
    MatTooltipModule,
    PipesModule,
    ReactiveFormsModule,
    RouterModule,
    SelectModule,
    SharedUiModule,
    TableModule,
    TextareaModule,
    RolesRoutingModule,
    AutocompleteModule,
    CategoryAutocompleteComponent,
    TagAutocompleteComponent,
    GrantPickerComponent,
  ],
})
export class RolesModule {}
