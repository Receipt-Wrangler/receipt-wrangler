import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from "@angular/material/icon";
import { MatListModule } from "@angular/material/list";
import { MatTabsModule } from "@angular/material/tabs";
import { MatTooltipModule } from "@angular/material/tooltip";
import { ColorPickerModule } from "src/color-picker/color-picker.module";
import { PipesModule } from "src/pipes/pipes.module";
import { SelectModule } from "src/select/select.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { UserAutocompleteModule } from "src/user-autocomplete/user-autocomplete.module";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { ButtonModule } from "../button";
import { CheckboxModule } from "../checkbox/checkbox.module";
import { DirectivesModule } from "../directives";
import { InputModule } from "../input";
import { ApiKeyFormDialogComponent } from "./api-key-form-dialog/api-key-form-dialog.component";
import { ApiKeyTableFilterComponent } from "./api-key-table-filter/api-key-table-filter.component";
import { ApiKeyTableComponent } from "./api-key-table/api-key-table.component";
import { ApiKeysComponent } from "./api-keys/api-keys.component";
import { SettingsRoutingModule } from "./settings-routing.module";
import { SettingsComponent } from "./settings/settings.component";
import { UserPreferencesComponent } from "./user-preferences/user-preferences.component";
import { UserProfileComponent } from "./user-profile/user-profile.component";
import { UserShortcutComponent } from './user-shortcut/user-shortcut.component';

@NgModule({
  declarations: [
    ApiKeyFormDialogComponent,
    ApiKeyTableComponent,
    ApiKeyTableFilterComponent,
    ApiKeysComponent,
    SettingsComponent,
    UserProfileComponent,
    UserPreferencesComponent,
    UserShortcutComponent,
  ],
  imports: [
    AutocompleteModule,
    ButtonModule,
    CheckboxModule,
    ColorPickerModule,
    CommonModule,
    DirectivesModule,
    InputModule,
    MatIconModule,
    MatListModule,
    MatTabsModule,
    MatTooltipModule,
    PipesModule,
    ReactiveFormsModule,
    SelectModule,
    SettingsRoutingModule,
    SharedUiModule,
    TableModule,
    UserAutocompleteModule,
  ],
})
export class SettingsModule {}
