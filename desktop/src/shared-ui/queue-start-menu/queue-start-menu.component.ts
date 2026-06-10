import { Component, OnInit, input } from "@angular/core";
import { Store } from "@ngxs/store";
import { Permission, Receipt } from "../../open-api/index";
import { QueueMode, ReceiptQueueService } from "../../services/receipt-queue.service";
import { AuthState, GroupState } from "../../store/index";

@Component({
    selector: "app-queue-start-menu",
    templateUrl: "./queue-start-menu.component.html",
    styleUrl: "./queue-start-menu.component.scss",
    standalone: false
})
export class QueueStartMenuComponent implements OnInit {
  public readonly buttonText = input<string>("");

  public readonly buttonIcon = input<string>("");

  public readonly matButtonType = input<"matRaisedButton" | "iconButton" | "basic">("matRaisedButton");

  public readonly color = input<string>("primary");

  public readonly receiptIds = input<string[] | number[]>([]);

  public readonly receipts = input<Receipt[]>([]);

  protected readonly QueueMode = QueueMode;

  public canEdit: boolean = false;

  constructor(
    private receiptQueueService: ReceiptQueueService,
    private store: Store,
  ) {}

  public ngOnInit(): void {
    this.setCanEdit();
  }

  private setCanEdit(): void {
    const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
    this.canEdit = this.store.selectSnapshot(
      AuthState.hasGroupPermission(Number.parseInt(groupId), Permission.GroupReceiptsUpdate)
    );
  }

  private getReceiptIds(): string[] {
    if (this.receiptIds().length > 0) {
      return this.receiptIds().map(id => id.toString());
    } else if (this.receipts().length > 0) {
      return this.receipts().map(receipt => receipt.id.toString());
    } else {
      return [];
    }
  }

  public initQueue(mode: QueueMode): void {
    this.receiptQueueService.initQueueAndNavigate(this.getReceiptIds(), mode);
  }
}
