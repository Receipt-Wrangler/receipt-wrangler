import { Component, Input } from "@angular/core";
import { ButtonModule } from "../../../button/index";
import { ReceiptPagedRequestCommand } from "../../../open-api/index";
import { ReceiptExportService } from "../../../services/receipt-export.service";

@Component({
  selector: "app-export-button",
  imports: [
    ButtonModule
  ],
  templateUrl: "./export-button.component.html",
  styleUrl: "./export-button.component.scss"
})
export class ExportButtonComponent {
  @Input() public filter?: ReceiptPagedRequestCommand;

  @Input() public groupId?: string;

  constructor(private receiptExportService: ReceiptExportService) {}

  public exportReceipts(): void {
    if (this.filter && this.groupId) {
      this.exportReceiptsByFilter();
    }
  }

  private exportReceiptsByFilter(): void {
    if (!this.filter || !this.groupId) {
      return;
    }

    this.receiptExportService.exportReceiptsFromFilter(this.groupId, this.filter);
  }
}
