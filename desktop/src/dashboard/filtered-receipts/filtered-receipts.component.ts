import { Component, OnInit, input, signal } from "@angular/core";
import { UntilDestroy } from "@ngneat/until-destroy";
import { take, tap } from "rxjs";
import { ReceiptFilterService } from "src/services/receipt-filter.service";
import { Receipt, ReceiptPagedRequestCommand, Widget } from "../../open-api";

@UntilDestroy()
@Component({
  selector: "/app-filtered-receipts",
  templateUrl: "./filtered-receipts.component.html",
  styleUrls: ["./filtered-receipts.component.scss"],
  standalone: false
})
export class FilteredReceiptsComponent implements OnInit {
  public readonly widget = input.required<Widget>();

  public readonly groupId = input<number>();

  public page: number = 1;

  public pageSize: number = 25;

  public receipts = signal<Receipt[]>([]);

  public buildItemRouterLink = (receipt: Receipt): string => {
    return "/receipts/" + receipt.id + "/view";
  };

  constructor(
    private receiptFilterService: ReceiptFilterService,
  ) {}

  public ngOnInit(): void {
    this.getData();
  }

  public endOfListReached(): void {
    this.page++;
    this.getData();
  }

  private getData(): void {
    const groupIdValue = this.groupId();
    if (!groupIdValue) {
      return;
    }

    const groupId = groupIdValue;
    const command: ReceiptPagedRequestCommand = {
      page: this.page,
      pageSize: this.pageSize,
      filter: this.widget().configuration,
      orderBy: "date",
      sortDirection: "desc",
    };
    this.receiptFilterService
      .getPagedReceiptsForGroups(
        groupId?.toString() ?? "",
        undefined,
        undefined,
        undefined,
        undefined,
        command
      )
      .pipe(
        take(1),
        tap((pagedData) => {
          this.receipts.update(prev => [...prev, ...pagedData.data as any as Receipt[]]);
        })
      )
      .subscribe();
  }
}
