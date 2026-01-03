import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { ButtonModule } from "../button/index";
import { DirectivesModule } from "../directives/directives.module";
import { InputModule } from "../input/index";
import { PipesModule } from "../pipes/index";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { CustomFieldFormComponent } from "./custom-field-form/custom-field-form.component";
import { CustomFieldTableComponent } from "./custom-field-table/custom-field-table.component";

import { CustomFieldsRoutingModule } from "./custom-fields-routing.module";


@NgModule({
  declarations: [
    CustomFieldTableComponent,
    CustomFieldFormComponent
  ],
  imports: [
    CommonModule,
    CustomFieldsRoutingModule,
    DirectivesModule,
    SharedUiModule,
    TableModule,
    ReactiveFormsModule,
    InputModule,
    PipesModule,
    SelectModule,
    ButtonModule
  ]
})
export class CustomFieldsModule {}
