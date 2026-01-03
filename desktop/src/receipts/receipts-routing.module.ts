import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { FormMode } from "src/enums/form-mode.enum";
import { GroupRoleGuard } from "src/guards/group-role.guard";
import { GroupGuard } from "src/guards/group.guard";
import { receiptGuardGuard } from "src/guards/receipt-guard.guard";
import { GroupRole } from "../open-api";
import { categoryResolverFn } from "../resolvers/categories.resolver";
import { customFieldResolverFn } from "../resolvers/custom-field.resolver";
import { receiptResolverFn } from "../resolvers/receipt.resolver";
import { tagResolverFn } from "../resolvers/tags.resolver";
import { ReceiptFormComponent } from "./receipt-form/receipt-form.component";
import { ReceiptsTableComponent } from "./receipts-table/receipts-table.component";

const routes: Routes = [
  {
    path: "group/:groupId",
    component: ReceiptsTableComponent,
    canActivate: [GroupGuard],
    resolve: {
      tags: tagResolverFn,
      categories: categoryResolverFn,
    },
    data: {
      groupGuardBasePath: `/receipts/group`,
    },
  },
  {
    path: "add",
    component: ReceiptFormComponent,
    resolve: {
      tags: tagResolverFn,
      categories: categoryResolverFn,
      customFields: customFieldResolverFn,
    },
    data: {
      mode: FormMode.add,
      groupRole: GroupRole.Editor,
    },
    canActivate: [GroupRoleGuard],
  },
  {
    path: ":id/view",
    component: ReceiptFormComponent,
    resolve: {
      tags: tagResolverFn,
      categories: categoryResolverFn,
      receipt: receiptResolverFn,
      customFields: customFieldResolverFn,
    },
    data: {
      mode: FormMode.view,
      groupRole: GroupRole.Viewer,
    },
    canActivate: [receiptGuardGuard],
  },
  {
    path: ":id/edit",
    component: ReceiptFormComponent,
    resolve: {
      tags: tagResolverFn,
      categories: categoryResolverFn,
      receipt: receiptResolverFn,
      customFields: customFieldResolverFn,
    },
    data: {
      mode: FormMode.edit,
      groupRole: GroupRole.Editor,
    },
    canActivate: [receiptGuardGuard],
  },
  {
    path: "",
    redirectTo: "",
    pathMatch: "full",
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ReceiptsRoutingModule {
}
