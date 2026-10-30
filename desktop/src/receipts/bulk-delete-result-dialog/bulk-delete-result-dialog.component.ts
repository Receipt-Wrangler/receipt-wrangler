import { Component, Inject } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

export interface BulkDeleteResultFailure {
  name: string;
  message: string;
}

export interface BulkDeleteResultDialogData {
  deletedCount: number;
  requestedCount: number;
  failures: BulkDeleteResultFailure[];
}

/**
 * Lists the receipts a bulk delete did not remove, by name, with the server's
 * reason for each. Only opened when at least one receipt failed; a clean run is
 * a success snackbar instead.
 */
@Component({
  selector: "app-bulk-delete-result-dialog",
  templateUrl: "./bulk-delete-result-dialog.component.html",
  styleUrls: ["./bulk-delete-result-dialog.component.scss"],
  standalone: false,
})
export class BulkDeleteResultDialogComponent {
  constructor(
    public matDialogRef: MatDialogRef<BulkDeleteResultDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: BulkDeleteResultDialogData,
  ) {}

  public close(): void {
    this.matDialogRef.close();
  }
}
