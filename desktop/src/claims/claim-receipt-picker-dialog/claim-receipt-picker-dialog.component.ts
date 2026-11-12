import { CurrencyPipe } from "@angular/common";
import { Component, inject } from "@angular/core";
import { FormArray, FormBuilder, ReactiveFormsModule } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { AutocompleteModule } from "../../autocomplete/autocomplete.module";
import { Receipt } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

export interface ClaimReceiptPickerDialogData {
  claimName: string;
  /**
   * The claim group's receipts that are on no claim the viewer can see. One in
   * a hidden claim still appears here, and the server's refusal names it.
   */
  candidates: Receipt[];
  /** How many of the group's receipts are left out for being in another (visible) claim. */
  inOtherClaimsCount: number;
}

/**
 * Picks receipts to add to a claim, from the claim group's receipts, with the
 * same chip picker the app uses to pick categories and tags (type to filter by
 * name). Closes with the picked receipts and writes nothing itself: the claim
 * page owns the request and reports its per-receipt result.
 *
 * A receipt already in another claim is left out rather than shown disabled:
 * the shared chip autocomplete has no disabled option, and the server would
 * refuse it anyway. A line under the field says how many were left out.
 *
 * The other way in is the receipts table's "Add selected receipts to a claim",
 * which brings the table's filters and month stepper with it; this one is for
 * the receipt or two you think of while looking at the claim.
 */
@Component({
  selector: "app-claim-receipt-picker-dialog",
  templateUrl: "./claim-receipt-picker-dialog.component.html",
  styleUrls: ["./claim-receipt-picker-dialog.component.scss"],
  standalone: true,
  imports: [AutocompleteModule, PipesModule, ReactiveFormsModule, SharedUiModule],
  // customCurrency needs CurrencyPipe, which nothing provides at the root.
  providers: [CurrencyPipe],
})
export class ClaimReceiptPickerDialogComponent {
  public readonly matDialogRef =
    inject<MatDialogRef<ClaimReceiptPickerDialogComponent, Receipt[]>>(MatDialogRef);

  public readonly data = inject<ClaimReceiptPickerDialogData>(MAT_DIALOG_DATA);

  public readonly form = inject(FormBuilder).group({
    receipts: new FormArray<any>([]),
  });

  public readonly inOtherClaimsNote =
    this.data.inOtherClaimsCount === 1
      ? "1 receipt already in another claim is not listed."
      : `${this.data.inOtherClaimsCount} receipts already in another claim are not listed.`;

  public cancel(): void {
    this.matDialogRef.close(undefined);
  }

  public submit(): void {
    const picked = this.form.controls.receipts.value as Receipt[];
    if (picked.length === 0) {
      this.form.controls.receipts.setErrors({ required: true });
      this.form.markAllAsTouched();
      return;
    }

    this.matDialogRef.close(picked);
  }
}
