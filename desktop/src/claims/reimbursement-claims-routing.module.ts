import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { FormMode } from "../enums/form-mode.enum";
import { groupPermissionGuard } from "../guards/group-permission.guard";
import { GroupGuard } from "../guards/group.guard";
import { Permission } from "../open-api";
import { ClaimFormComponent } from "./claim-form/claim-form.component";
import { ClaimListComponent } from "./claim-list/claim-list.component";
import { claimResolverFn } from "./claim.resolver";

const routes: Routes = [
  {
    path: "group/:groupId",
    component: ClaimListComponent,
    canActivate: [GroupGuard],
    data: {
      groupGuardBasePath: `/claims/group`,
    },
  },
  {
    path: "add",
    component: ClaimFormComponent,
    canActivate: [groupPermissionGuard],
    data: {
      mode: FormMode.add,
      groupPermission: Permission.GroupReceiptsUpdate,
      // The group a new claim lands in is the receipt form's add target: the
      // group being browsed, or the sole real group when that is "All".
      useAddTargetGroupId: true,
    },
  },
  {
    // Gated in the resolver, which is the first point the claim's group is known.
    path: ":id/view",
    component: ClaimFormComponent,
    resolve: { claim: claimResolverFn },
    data: { mode: FormMode.view },
  },
  {
    path: ":id/edit",
    component: ClaimFormComponent,
    resolve: { claim: claimResolverFn },
    data: { mode: FormMode.edit },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ReimbursementClaimsRoutingModule {}
