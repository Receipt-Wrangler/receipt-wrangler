import { Component, inject } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

export interface ClaimReceiptsResultFailure {
  name: string;
  message: string;
}

export interface ClaimReceiptsResultDialogData {
  headerText: string;
  /** e.g. "Added 2 of 3 receipts to Travel voucher 2026-14." */
  summary: string;
  failures: ClaimReceiptsResultFailure[];
}

/**
 * Lists the receipts an add-to-claim or remove-from-claim call refused, by name,
 * with the server's reason for each. Only opened when something was refused; a
 * clean run is a snackbar. Same shape as the bulk-delete result dialog, which it
 * mirrors rather than reuses because that one's wording is about deleting.
 */
@Component({
  selector: "app-claim-receipts-result-dialog",
  templateUrl: "./claim-receipts-result-dialog.component.html",
  styleUrls: ["./claim-receipts-result-dialog.component.scss"],
  standalone: true,
  imports: [SharedUiModule],
})
export class ClaimReceiptsResultDialogComponent {
  public readonly matDialogRef = inject<MatDialogRef<ClaimReceiptsResultDialogComponent>>(MatDialogRef);

  public readonly data = inject<ClaimReceiptsResultDialogData>(MAT_DIALOG_DATA);

  public close(): void {
    this.matDialogRef.close();
  }
}
