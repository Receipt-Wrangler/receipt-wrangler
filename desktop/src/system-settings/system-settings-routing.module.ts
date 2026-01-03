import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { FormMode } from "../enums/form-mode.enum";
import { FormConfig } from "../interfaces";
import { PromptFormComponent } from "../prompt/prompt-form/prompt-form.component";
import { PromptTableComponent } from "../prompt/prompt-table/prompt-table.component";
import { promptResolver } from "../prompt/prompt.resolver";
import { promptsResolver } from "../prompt/prompts.resolver";
import {
  ReceiptProcessingSettingsFormComponent
} from "../receipt-processing-settings/receipt-processing-settings-form/receipt-processing-settings-form.component";
import {
  ReceiptProcessingSettingsTableComponent
} from "../receipt-processing-settings/receipt-processing-settings-table/receipt-processing-settings-table.component";
import { receiptProcessingSettingsResolver } from "../receipt-processing-settings/receipt-processing-settings.resolver";
import { allReceiptProcessingSettingsResolver } from "./resolvers/receipt-processing-settings.resolver";
import { systemEmailResolver } from "./resolvers/system-email.resolver";
import { systemSettingsResolver } from "./resolvers/system-settings.resolver";
import { SystemEmailFormComponent } from "./system-email-form/system-email-form.component";
import { allGroupsResolver } from "./system-email-table/all-groups.resolver";
import { SystemEmailTableComponent } from "./system-email-table/system-email-table.component";
import { SystemSettingsFormComponent } from "./system-settings-form/system-settings-form.component";
import { SystemSettingsComponent } from "./system-settings/system-settings.component";
import { SystemTaskTableComponent } from "./system-task-table/system-task-table.component";

const routes: Routes = [
  {
    path: "",
    redirectTo: "system-emails",
    pathMatch: "full",
  },
  {
    path: "",
    component: SystemSettingsComponent,
    children: [
      {
        path: "system-emails",
        component: SystemEmailTableComponent,
        resolve: {
          allGroups: allGroupsResolver,
        }
      },
      {
        path: "prompts",
        component: PromptTableComponent,
        resolve: {
          allGroups: allGroupsResolver,
          allReceiptProcessingSettings: allReceiptProcessingSettingsResolver,
        }
      },
      {
        path: "receipt-processing-settings",
        component: ReceiptProcessingSettingsTableComponent,
        resolve: {
          systemSettings: systemSettingsResolver,
        }
      },
      {
        path: "system-tasks",
        component: SystemTaskTableComponent,
        resolve: {
          prompts: promptsResolver,
          allReceiptProcessingSettings: allReceiptProcessingSettingsResolver,
        }
      },
      {
        path: "settings/view",
        component: SystemSettingsFormComponent,
        data: {
          formConfig: {
            mode: FormMode.view,
            headerText: "View System Settings",
          } as FormConfig,
        },
        resolve: {
          allReceiptProcessingSettings: allReceiptProcessingSettingsResolver,
          systemSettings: systemSettingsResolver,
        }
      },
      {
        path: "settings/edit",
        component: SystemSettingsFormComponent,
        data: {
          formConfig: {
            mode: FormMode.edit,
            headerText: "Edit System Settings",
          } as FormConfig,
        },
        resolve: {
          allReceiptProcessingSettings: allReceiptProcessingSettingsResolver,
          systemSettings: systemSettingsResolver,
        }
      },
    ]
  },
  {
    path: "prompts/create",
    component: PromptFormComponent,
    data: {
      formConfig: {
        mode: FormMode.add,
        headerText: "Create Prompt",
      } as FormConfig,
    },
  },
  {
    path: "prompts/:id/view",
    component: PromptFormComponent,
    data: {
      formConfig: {
        mode: FormMode.view,
        headerText: "View Prompt",
      } as FormConfig,
      setHeaderText: true,
    },
    resolve: {
      prompt: promptResolver
    }
  },
  {
    path: "prompts/:id/edit",
    component: PromptFormComponent,
    data: {
      formConfig: {
        mode: FormMode.edit,
        headerText: "Edit Prompt",
      } as FormConfig,
      setHeaderText: true,
    },
    resolve: {
      prompt: promptResolver
    }
  },
  {
    path: "system-emails/create",
    component: SystemEmailFormComponent,
    data: {
      formConfig: {
        mode: FormMode.add,
        headerText: "Create System Email",
      } as FormConfig,
    }
  },
  {
    path: "system-emails/:id/view",
    component: SystemEmailFormComponent,
    data: {
      formConfig: {
        mode: FormMode.view,
      } as FormConfig,
      setHeaderText: true,
    },
    resolve: {
      systemEmail: systemEmailResolver,
      prompts: promptsResolver,
      allReceiptProcessingSettings: allReceiptProcessingSettingsResolver,
    }
  },
  {
    path: "system-emails/:id/edit",
    component: SystemEmailFormComponent,
    data: {
      formConfig: {
        mode: FormMode.edit,
      } as FormConfig,
      setHeaderText: true,
    },
    resolve: {
      systemEmail: systemEmailResolver,
      prompts: promptsResolver,
      allReceiptProcessingSettings: allReceiptProcessingSettingsResolver,
    }
  },
  {
    path: "receipt-processing-settings/create",
    component: ReceiptProcessingSettingsFormComponent,
    data: {
      formConfig: {
        mode: FormMode.add,
        headerText: "Create Receipt Processing Settings",
      } as FormConfig,
    },
    resolve: {
      prompts: promptsResolver,
    }
  },
  {
    path: "receipt-processing-settings/:id/view",
    component: ReceiptProcessingSettingsFormComponent,
    data: {
      formConfig: {
        mode: FormMode.view,
      } as FormConfig,
      setHeaderText: true,
    },
    resolve: {
      prompts: promptsResolver,
      receiptProcessingSettings: receiptProcessingSettingsResolver,
    }
  },
  {
    path: "receipt-processing-settings/:id/edit",
    component: ReceiptProcessingSettingsFormComponent,
    data: {
      formConfig: {
        mode: FormMode.edit,
      } as FormConfig,
      setHeaderText: true,
    },
    resolve: {
      prompts: promptsResolver,
      receiptProcessingSettings: receiptProcessingSettingsResolver,
    }
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SystemSettingsRoutingModule {
}
