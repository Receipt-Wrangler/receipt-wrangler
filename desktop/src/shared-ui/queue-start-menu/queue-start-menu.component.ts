import { Component, Input, OnInit } from "@angular/core";
import { Store } from "@ngxs/store";
import { GroupRole, Receipt } from "../../open-api/index";
import { GroupRolePipe } from "../../pipes/group-role.pipe";
import { QueueMode, ReceiptQueueService } from "../../services/receipt-queue.service";
import { GroupState } from "../../store/index";

@Component({
    selector: "app-queue-start-menu",
    templateUrl: "./queue-start-menu.component.html",
    styleUrl: "./queue-start-menu.component.scss",
    standalone: false
})
export class QueueStartMenuComponent implements OnInit {
  @Input() public buttonText: string = "";

  @Input() public buttonIcon: string = "";

  @Input() public matButtonType: "matRaisedButton" | "iconButton" | "basic" = "matRaisedButton";

  @Input() public color: string = "primary";

  @Input() public receiptIds: string[] | number[] = [];

  @Input() public receipts: Receipt[] = [];

  protected readonly QueueMode = QueueMode;

  public canEdit: boolean = false;

  constructor(
    private receiptQueueService: ReceiptQueueService,
    private store: Store,
    private groupRolePipe: GroupRolePipe
  ) {}

  public ngOnInit(): void {
    this.setCanEdit();
  }

  private setCanEdit(): void {
    const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
    this.canEdit = this.groupRolePipe.transform(groupId, GroupRole.Editor);
  }

  private getReceiptIds(): string[] {
    if (this.receiptIds.length > 0) {
      return this.receiptIds.map(id => id.toString());
    } else if (this.receipts.length > 0) {
      return this.receipts.map(receipt => receipt.id.toString());
    } else {
      return [];
    }
  }

  public initQueue(mode: QueueMode): void {
    this.receiptQueueService.initQueueAndNavigate(this.getReceiptIds(), mode);
  }
}
