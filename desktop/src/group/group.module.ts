import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { MatDialogModule } from "@angular/material/dialog";
import { MatListModule } from "@angular/material/list";
import { MatTableModule } from "@angular/material/table";
import { CheckboxModule } from "src/checkbox/checkbox.module";
import { PipesModule } from "src/pipes/pipes.module";
import { SelectModule } from "src/select/select.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { UserAutocompleteModule } from "src/user-autocomplete/user-autocomplete.module";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { ButtonModule } from "../button";
import { InputModule } from "../input";
import { GroupDetailsComponent } from "./group-details/group-details.component";
import { GroupFormComponent } from "./group-form/group-form.component";
import { GroupMemberFormComponent } from "./group-member-form/group-member-form.component";
import { GroupReceiptSettingsComponent } from "./group-receipt-settings/group-receipt-settings.component";
import { GroupRoutingModule } from "./group-routing.module";
import { GroupSettingsAiComponent } from "./group-settings-ai/group-settings-ai.component";
import { GroupSettingsEmailComponent } from "./group-settings-email/group-settings-email.component";
import { GroupSettingsComponent } from "./group-settings/group-settings.component";
import { GroupTableFilterComponent } from "./group-table-filter/group-table-filter.component";
import { GroupTableEditButtonPipe } from "./group-table/group-table-edit-button.pipe";
import { GroupTableComponent } from "./group-table/group-table.component";
import { GroupTabsComponent } from "./group-tabs/group-tabs.component";

@NgModule({
  declarations: [
    GroupTableComponent,
    GroupFormComponent,
    GroupMemberFormComponent,
    GroupSettingsComponent,
    GroupSettingsEmailComponent,
    GroupTabsComponent,
    GroupDetailsComponent,
    GroupTableFilterComponent,
    GroupTableEditButtonPipe,
    GroupSettingsAiComponent,
    GroupReceiptSettingsComponent,
  ],
  imports: [
    ButtonModule,
    CheckboxModule,
    CommonModule,
    PipesModule,
    GroupRoutingModule,
    InputModule,
    MatCardModule,
    MatDialogModule,
    MatListModule,
    MatTableModule,
    PipesModule,
    ReactiveFormsModule,
    SelectModule,
    SharedUiModule,
    TableModule,
    UserAutocompleteModule,
    AutocompleteModule,
  ],
  exports: [GroupTableComponent],
})
export class GroupsModule {}
