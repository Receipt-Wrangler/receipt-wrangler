import { Component, Inject } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { DialogComponent } from "../dialog/dialog.component";
import { ImageViewerComponent } from "../image-viewer/image-viewer.component";
import { ButtonComponent } from "../../button/button/button.component";

export interface SourceFileViewerDialogData {
  /** Base64 data URI of the upload, already converted for display by the API. */
  encodedImage: string;
  /** The upload's original file name, used as the dialog heading. */
  name: string;
}

/**
 * Shows the upload behind a failed activity, so a user can see what the scan was
 * working from before deciding whether to retry it or enter the receipt by hand.
 */
@Component({
  selector: "app-source-file-viewer-dialog",
  templateUrl: "./source-file-viewer-dialog.component.html",
  imports: [
    DialogComponent,
    ImageViewerComponent,
    ButtonComponent,
  ],
})
export class SourceFileViewerDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<SourceFileViewerDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: SourceFileViewerDialogData,
  ) {}

  public close(): void {
    this.dialogRef.close();
  }
}
