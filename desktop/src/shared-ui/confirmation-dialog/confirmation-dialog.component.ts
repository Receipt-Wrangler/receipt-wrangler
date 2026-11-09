import { Component, Input, TemplateRef } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { DialogComponent } from '../dialog/dialog.component';
import { DialogFooterComponent } from '../dialog-footer/dialog-footer.component';

@Component({
    selector: 'app-confirmation-dialog',
    templateUrl: './confirmation-dialog.component.html',
    styleUrls: ['./confirmation-dialog.component.scss'],
    imports: [DialogComponent, DialogFooterComponent]
})
export class ConfirmationDialogComponent {
  constructor(private dialogRef: MatDialogRef<ConfirmationDialogComponent>) {}

  @Input() public headerText: string = '';

  @Input() public dialogContent: string = '';

  public submitButtonClicked(): void {
    this.dialogRef.close(true);
  }

  public cancelButtonClicked(): void {
    this.dialogRef.close(false);
  }
}
