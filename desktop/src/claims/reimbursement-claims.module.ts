import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { RouterModule } from "@angular/router";
import { ButtonModule } from "../button";
import { DatepickerModule } from "../datepicker/datepicker.module";
import { DirectivesModule } from "../directives";
import { InputModule } from "../input";
import { PipesModule } from "../pipes/pipes.module";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { TableModule } from "../table/table.module";
import { TextareaModule } from "../textarea/textarea.module";
import { ClaimFormComponent } from "./claim-form/claim-form.component";
import { ClaimListComponent } from "./claim-list/claim-list.component";
import { ReimbursementClaimsRoutingModule } from "./reimbursement-claims-routing.module";

/**
 * Reimbursement claims (fork feature), served under `/claims`. Named for the
 * generated `ReimbursementClaimService` so it cannot be mistaken for anything
 * about JWT claims (`ClaimsService`, the `Claims` model).
 *
 * The dialogs are standalone because the receipts table opens one of them, and
 * importing this module there would also register these routes under it.
 */
@NgModule({
  declarations: [ClaimFormComponent, ClaimListComponent],
  imports: [
    ButtonModule,
    CommonModule,
    DatepickerModule,
    DirectivesModule,
    InputModule,
    MatDialogModule,
    MatIconModule,
    MatProgressSpinnerModule,
    PipesModule,
    ReactiveFormsModule,
    ReimbursementClaimsRoutingModule,
    RouterModule,
    SelectModule,
    SharedUiModule,
    TableModule,
    TextareaModule,
  ],
})
export class ReimbursementClaimsModule {}
