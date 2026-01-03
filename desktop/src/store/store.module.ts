import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { NgxsReduxDevtoolsPluginModule } from "@ngxs/devtools-plugin";
import { NgxsStoragePluginModule } from "@ngxs/storage-plugin";
import { NgxsModule } from "@ngxs/store";
import { environment } from "src/environments/environment.development";
import { AboutState } from "./about.state";
import { ApiKeyTableState } from "./api-key-table.state";
import { AuthState } from "./auth.state";
import { CategoryTableState } from "./category-table.state";
import { CustomFieldTableState } from "./custom-field-table.state";
import { DashboardState } from "./dashboard.state";
import { FeatureConfigState } from "./feature-config.state";
import { GroupTableState } from "./group-table.state";
import { GroupState } from "./group.state";
import { LayoutState } from "./layout.state";
import { PromptTableState } from "./prompt-table.state";
import { ReceiptProcessingSettingsTableState } from "./receipt-processing-settings-table.state";
import { ReceiptProcessingSettingsTaskTableState } from "./receipt-processing-settings-task-table.state";
import { ReceiptTableState } from "./receipt-table.state";
import { SystemEmailTableState } from "./system-email-table.state";
import { SystemEmailTaskTableState } from "./system-email-task-table.state";
import { SystemSettingsState } from "./system-settings.state";
import { SystemTaskTableState } from "./system-task-table.state";
import { TagTableState } from "./tag-table.state";
import { UserState } from "./user.state";

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    NgxsModule.forRoot([
      AboutState,
      ApiKeyTableState,
      AuthState,
      CategoryTableState,
      CustomFieldTableState,
      DashboardState,
      FeatureConfigState,
      GroupState,
      GroupTableState,
      LayoutState,
      PromptTableState,
      ReceiptProcessingSettingsTableState,
      ReceiptProcessingSettingsTaskTableState,
      ReceiptTableState,
      SystemEmailTableState,
      SystemEmailTaskTableState,
      SystemSettingsState,
      SystemTaskTableState,
      TagTableState,
      UserState,
    ]),
    NgxsReduxDevtoolsPluginModule.forRoot({
      disabled: environment.isProd,
    }),
    NgxsStoragePluginModule.forRoot({
      keys: [
        "about",
        "apiKeyTable",
        "auth",
        "categoryTable",
        "customFieldTable",
        "dashboards",
        "groupTable",
        "groups",
        "layout",
        "promptTable",
        "receiptProcessingSettingsTable",
        "receiptProcessingSettingsTaskTable",
        "receiptTable",
        "systemEmailTable",
        "systemEmailTaskTable",
        "systemSettings",
        "systemTaskTable",
        "tagTable",
        "users",
      ],
    }),
  ],
})
export class StoreModule {}
