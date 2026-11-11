import { ApplicationConfig, provideAppInitializer, provideEnvironmentInitializer, provideZonelessChangeDetection, inject } from "@angular/core";
import { CurrencyPipe } from "@angular/common";
import { provideHttpClient, withInterceptors } from "@angular/common/http";
import { importProvidersFrom } from "@angular/core";
import { DomSanitizer } from "@angular/platform-browser";
import { provideAnimationsAsync } from "@angular/platform-browser/animations/async";
import { provideRouter } from "@angular/router";
import { provideCharts, withDefaultRegisterables } from "ng2-charts";
import { provideNgxMask } from "ngx-mask";
import { withNgxsReduxDevtoolsPlugin } from "@ngxs/devtools-plugin";
import { withNgxsStoragePlugin } from "@ngxs/storage-plugin";
import { provideStore } from "@ngxs/store";
import { AppInitService, initAppData } from "src/services";
import { httpInterceptor } from "../interceptors/http-interceptor";
import { MatIconRegistry } from "@angular/material/icon";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { NgxMaskDirective, NgxMaskPipe } from "ngx-mask";
import { ApiModule, Configuration } from "../open-api";
import { AboutState } from "../store/about.state";
import { ApiKeyTableState } from "../store/api-key-table.state";
import { AuthState } from "../store/auth.state";
import { CategoryTableState } from "../store/category-table.state";
import { CustomFieldTableState } from "../store/custom-field-table.state";
import { OidcProviderTableState } from "../store/oidc-provider-table.state";
import { DashboardState } from "../store/dashboard.state";
import { FeatureConfigState } from "../store/feature-config.state";
import { GroupTableState } from "../store/group-table.state";
import { GroupState } from "../store/group.state";
import { LayoutState } from "../store/layout.state";
import { PromptTableState } from "../store/prompt-table.state";
import { ReceiptProcessingSettingsTableState } from "../store/receipt-processing-settings-table.state";
import { ReceiptProcessingSettingsTaskTableState } from "../store/receipt-processing-settings-task-table.state";
import { ReceiptTableState } from "../store/receipt-table.state";
import { ReportTemplateTableState } from "../store/report-template-table.state";
import { SystemEmailTableState } from "../store/system-email-table.state";
import { SystemEmailTaskTableState } from "../store/system-email-task-table.state";
import { SystemSettingsState } from "../store/system-settings.state";
import { SystemTaskTableState } from "../store/system-task-table.state";
import { TagTableState } from "../store/tag-table.state";
import { UserTableState } from "../store/user-table.state";
import { UserState } from "../store/user.state";
import { environment } from "src/environments/environment.development";
import { appRoutes } from "./app.routes";

const ngxsStates = [
  AboutState,
  ApiKeyTableState,
  AuthState,
  CategoryTableState,
  CustomFieldTableState,
  OidcProviderTableState,
  DashboardState,
  FeatureConfigState,
  GroupState,
  GroupTableState,
  LayoutState,
  PromptTableState,
  ReceiptProcessingSettingsTableState,
  ReceiptProcessingSettingsTaskTableState,
  ReceiptTableState,
  ReportTemplateTableState,
  SystemEmailTableState,
  SystemEmailTaskTableState,
  SystemSettingsState,
  SystemTaskTableState,
  TagTableState,
  UserTableState,
  UserState,
];

const ngxsStorageKeys = [
  "about",
  "apiKeyTable",
  "auth",
  "categoryTable",
  "customFieldTable",
  "oidcProviderTable",
  "dashboards",
  "groupTable",
  "groups",
  "layout",
  "promptTable",
  "receiptProcessingSettingsTable",
  "receiptProcessingSettingsTaskTable",
  "receiptTable",
  "reportTemplateTable",
  "systemEmailTable",
  "systemEmailTaskTable",
  "systemSettings",
  "systemTaskTable",
  "tagTable",
  "userTable",
  "users",
];

export const appConfig: ApplicationConfig = {
  providers: [
    provideZonelessChangeDetection(),
    provideAnimationsAsync(),
    provideHttpClient(withInterceptors([httpInterceptor])),
    provideRouter(appRoutes),
    provideNgxMask(),
    // CustomCurrencyPipe formats through Angular's CurrencyPipe.
    CurrencyPipe,
    provideCharts(withDefaultRegisterables()),
    // The receipt form's "split" icon is a custom SVG.
    provideEnvironmentInitializer(() => {
      inject(MatIconRegistry).addSvgIcon(
        "split",
        inject(DomSanitizer).bypassSecurityTrustResourceUrl("../assets/split.svg")
      );
    }),
    provideAppInitializer(() => {
      const initializerFn = (initAppData)(inject(AppInitService));
      return initializerFn();
    }),
    provideStore(
      ngxsStates,
      { developmentMode: !environment.isProd },
      withNgxsStoragePlugin({ keys: ngxsStorageKeys }),
      withNgxsReduxDevtoolsPlugin({ disabled: environment.isProd }),
    ),
    importProvidersFrom(
      ApiModule.forRoot(() => new Configuration({ basePath: undefined })),
      MatSnackBarModule,
      MatTooltipModule,
      NgxMaskDirective,
      NgxMaskPipe,
    ),
  ],
};
