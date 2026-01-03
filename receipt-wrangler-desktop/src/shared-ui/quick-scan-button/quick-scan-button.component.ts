import { Component, EventEmitter, Output } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { QuickScanDialogComponent } from "../../receipts/quick-scan-dialog/quick-scan-dialog.component";

@Component({
    selector: "app-quick-scan-button",
    templateUrl: "./quick-scan-button.component.html",
    styleUrls: ["./quick-scan-button.component.scss"],
    standalone: false
})
export class QuickScanButtonComponent {
  @Output() public afterClosed: EventEmitter<void> = new EventEmitter();

  constructor(private matDialog: MatDialog) {}

  public showQuickScanDialog(): void {
    const ref = this.matDialog.open(
      QuickScanDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );

    ref
      .afterClosed()
      .pipe(
        take(1),
        tap(() => {
          this.afterClosed.emit();
        })
      )
      .subscribe();
  }
}
