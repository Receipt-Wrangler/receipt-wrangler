import { Component, Inject } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { DialogComponent } from "../dialog/dialog.component";
import { CopyButtonComponent } from "../copy-button/copy-button.component";
import { PrettyJsonComponent } from "../pretty-json/pretty-json.component";
import { ButtonComponent } from "../../button/button/button.component";

export interface DescriptionViewerDialogData {
  description: string;
  headerText?: string;
}

@Component({
  selector: "app-description-viewer-dialog",
  templateUrl: "./description-viewer-dialog.component.html",
  styleUrl: "./description-viewer-dialog.component.scss",
  imports: [
    DialogComponent,
    CopyButtonComponent,
    PrettyJsonComponent,
    ButtonComponent,
  ],
})
export class DescriptionViewerDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<DescriptionViewerDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: DescriptionViewerDialogData,
  ) {}

  public close(): void {
    this.dialogRef.close();
  }
}
