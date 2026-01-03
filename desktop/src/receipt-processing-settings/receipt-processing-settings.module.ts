import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatTooltip } from "@angular/material/tooltip";
import { RouterLink } from "@angular/router";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { ButtonModule } from "../button";
import { CheckboxModule } from "../checkbox/checkbox.module";
import { InputModule } from "../input";
import { PipesModule } from "../pipes";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { AiTypePipe } from "./pipes/ai-type.pipe";
import { OcrEnginePipe } from "./pipes/ocr-engine.pipe";
import { UrlLabelPipe } from "./pipes/url-label.pipe";
import {
  ReceiptProcessingSettingsChildSystemTaskAccordionComponent
} from "./receipt-processing-settings-child-system-task-accordion/receipt-processing-settings-child-system-task-accordion.component";
import { ReceiptProcessingSettingsFormComponent } from "./receipt-processing-settings-form/receipt-processing-settings-form.component";
import { ReceiptProcessingSettingsTableComponent } from "./receipt-processing-settings-table/receipt-processing-settings-table.component";


@NgModule({
  declarations: [ReceiptProcessingSettingsTableComponent, ReceiptProcessingSettingsFormComponent, ReceiptProcessingSettingsChildSystemTaskAccordionComponent],
  imports: [
    CommonModule,
    TableModule,
    SharedUiModule,
    RouterLink,
    ReactiveFormsModule,
    InputModule,
    PipesModule,
    SelectModule,
    AutocompleteModule,
    AiTypePipe,
    OcrEnginePipe,
    UrlLabelPipe,
    ButtonModule,
    MatTooltip,
    CheckboxModule
  ],
  exports: [
    ReceiptProcessingSettingsChildSystemTaskAccordionComponent
  ]
})
export class ReceiptProcessingSettingsModule {
}
