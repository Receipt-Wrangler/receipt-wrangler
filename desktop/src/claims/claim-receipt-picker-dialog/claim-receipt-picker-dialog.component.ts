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
  /** The claim group's receipts that are not already on this claim. */
  candidates: Receipt[];
}

/**
 * Picks receipts to add to a claim, from the claim group's receipts, with the
 * same chip picker the app uses to pick categories and tags (type to filter by
 * name). Closes with the picked receipts and writes nothing itself: the claim
 * page owns the request and reports its per-receipt result.
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
