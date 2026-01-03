import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatButton } from "@angular/material/button";
import { MatHint } from "@angular/material/form-field";
import { AlertComponent } from "../alert/alert.component";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { ButtonModule } from "../button";
import { CheckboxModule } from "../checkbox/checkbox.module";
import { InputModule } from "../input";
import { PipesModule } from "../pipes";
import { PromptModule } from "../prompt/prompt.module";
import { ReceiptProcessingSettingsModule } from "../receipt-processing-settings/receipt-processing-settings.module";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { TaskQueueFormControlPipe } from "./pipes/task-queue-form-control.pipe";
import { SystemEmailChildSystemTaskComponent } from "./system-email-child-system-task/system-email-child-system-task.component";
import { SystemEmailFormComponent } from "./system-email-form/system-email-form.component";
import { SystemEmailTableComponent } from "./system-email-table/system-email-table.component";
import { SystemSettingsFormComponent } from "./system-settings-form/system-settings-form.component";

import { SystemSettingsRoutingModule } from "./system-settings-routing.module";
import { SystemSettingsComponent } from "./system-settings/system-settings.component";
import { SystemTaskTableComponent } from "./system-task-table/system-task-table.component";


@NgModule({
  declarations: [SystemEmailTableComponent,
    SystemSettingsComponent,
    SystemEmailFormComponent,
    SystemSettingsFormComponent,
    SystemEmailChildSystemTaskComponent,
    SystemTaskTableComponent,
    TaskQueueFormControlPipe,
  ],
  imports: [
    ButtonModule,
    CommonModule,
    InputModule,
    PipesModule,
    PromptModule,
    ReactiveFormsModule,
    ReceiptProcessingSettingsModule,
    SharedUiModule,
    SystemSettingsRoutingModule,
    TableModule,
    CheckboxModule,
    AutocompleteModule,
    MatButton,
    SelectModule,
    MatHint,
    AlertComponent,
  ]
})
export class SystemSettingsModule {
}
