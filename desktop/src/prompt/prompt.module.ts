import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { ButtonModule } from "../button";
import { InputModule } from "../input";
import { PipesModule } from "../pipes";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { TextareaModule } from "../textarea/textarea.module";
import { PromptFormComponent } from "./prompt-form/prompt-form.component";

import { PromptRoutingModule } from "./prompt-routing.module";
import { PromptTableComponent } from "./prompt-table/prompt-table.component";


@NgModule({
  declarations: [
    PromptTableComponent,
    PromptFormComponent
  ],
  imports: [
    CommonModule,
    PromptRoutingModule,
    SharedUiModule,
    TableModule,
    ButtonModule,
    ReactiveFormsModule,
    InputModule,
    PipesModule,
    TextareaModule,
  ]
})
export class PromptModule {}
