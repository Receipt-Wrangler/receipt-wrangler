import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { UserAutocompleteModule } from "src/user-autocomplete/user-autocomplete.module";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { ButtonModule } from "../button";
import { CheckboxModule } from "../checkbox/checkbox.module";
import { DatepickerModule } from "../datepicker/datepicker.module";
import { DirectivesModule } from "../directives";
import { InputModule } from "../input";
import { PipesModule } from "../pipes";
import { SelectModule } from "../select/select.module";
import { TableModule } from "../table/table.module";
import { TextareaModule } from "../textarea/textarea.module";
import { AddGroupDialogComponent } from "./dialogs/add-group-dialog/add-group-dialog.component";
import { ColumnPickerDialogComponent } from "./dialogs/column-picker-dialog/column-picker-dialog.component";
import { ReportReceiptsDialogComponent } from "./dialogs/report-receipts-dialog/report-receipts-dialog.component";
import { ReportBuilderComponent } from "./report-builder/report-builder.component";
import { ReportConfigPanelComponent } from "./report-config-panel/report-config-panel.component";
import { ReportFiltersComponent } from "./report-filters/report-filters.component";
import { ReportGenerateBarComponent } from "./report-generate-bar/report-generate-bar.component";
import { ReportPreviewPanelComponent } from "./report-preview-panel/report-preview-panel.component";
import { ReportSectionComponent } from "./report-section/report-section.component";
import { ReportTemplateListComponent } from "./report-template-list/report-template-list.component";
import { ReportsRoutingModule } from "./reports-routing.module";

@NgModule({
  declarations: [
    AddGroupDialogComponent,
    ColumnPickerDialogComponent,
    ReportBuilderComponent,
    ReportConfigPanelComponent,
    ReportFiltersComponent,
    ReportGenerateBarComponent,
    ReportPreviewPanelComponent,
    ReportReceiptsDialogComponent,
    ReportSectionComponent,
    ReportTemplateListComponent,
  ],
  imports: [
    AutocompleteModule,
    ButtonModule,
    CheckboxModule,
    CommonModule,
    DatepickerModule,
    DirectivesModule,
    InputModule,
    MatDialogModule,
    MatIconModule,
    PipesModule,
    ReactiveFormsModule,
    ReportsRoutingModule,
    SelectModule,
    SharedUiModule,
    TableModule,
    TextareaModule,
    UserAutocompleteModule,
  ],
})
export class ReportsModule {}
