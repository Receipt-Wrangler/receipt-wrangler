import { Component, input } from "@angular/core";
import { ReceiptPagedRequestCommand } from "../../../open-api/index";
import { ReceiptExportService } from "../../../services/receipt-export.service";
import { ButtonComponent } from "../../../button/button/button.component";

@Component({
  selector: "app-export-button",
  imports: [
    ButtonComponent
  ],
  templateUrl: "./export-button.component.html",
  styleUrl: "./export-button.component.scss"
})
export class ExportButtonComponent {
  public readonly filter = input<ReceiptPagedRequestCommand>();

  public readonly groupId = input<string>();

  constructor(private receiptExportService: ReceiptExportService) {}

  public exportReceipts(): void {
    if (this.filter() && this.groupId()) {
      this.exportReceiptsByFilter();
    }
  }

  private exportReceiptsByFilter(): void {
    const filter = this.filter();
    const groupId = this.groupId();
    if (!filter || !groupId) {
      return;
    }

    this.receiptExportService.exportReceiptsFromFilter(groupId, filter);
  }
}
