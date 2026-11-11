import { Component, OnInit, input, signal } from "@angular/core";
import { UntilDestroy } from "@ngneat/until-destroy";
import { take, tap } from "rxjs";
import { ReceiptFilterService } from "src/services/receipt-filter.service";
import { Receipt, ReceiptPagedRequestCommand, Widget } from "../../open-api";
import { CardComponent } from "../../shared-ui/card/card.component";
import { QueueStartMenuComponent } from "../../shared-ui/queue-start-menu/queue-start-menu.component";
import { ExportButtonComponent } from "../../standalone/components/export-button/export-button.component";
import { DashboardListComponent } from "../dashboard-list/dashboard-list.component";
import { DateBlockComponent } from "../../standalone/components/date-block/date-block.component";
import { StatusChipComponent } from "../../shared-ui/status-chip/status-chip.component";
import { AppDatePipe } from "../../pipes/app-date.pipe";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { UserPipe } from "../../pipes/user.pipe";

@UntilDestroy()
@Component({
  selector: "/app-filtered-receipts",
  templateUrl: "./filtered-receipts.component.html",
  styleUrls: ["./filtered-receipts.component.scss"],
  imports: [CardComponent, QueueStartMenuComponent, ExportButtonComponent, DashboardListComponent, DateBlockComponent, StatusChipComponent, AppDatePipe, CustomCurrencyPipe, UserPipe]
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
