import { Injectable } from "@angular/core";
import { Router, RouteReuseStrategy } from "@angular/router";

export enum QueueMode {
  VIEW = "view",
  EDIT = "edit",
}

@Injectable({
  providedIn: "root"
})
export class ReceiptQueueService {
  constructor(
    private router: Router,
    private routeReuseStrategy: RouteReuseStrategy
  ) { }

  public initQueueAndNavigate(receiptIds: string[], mode: QueueMode, indexToStartAt: number = 0): void {
    this.routeReuseStrategy.shouldReuseRoute = () => false;

    const commands = ["receipts", receiptIds[indexToStartAt], mode];

    this.router.navigate(commands, {
      queryParams: {
        ids: receiptIds,
        queueMode: mode,
      },
      onSameUrlNavigation: "reload"
    });
  }

  public queueNext(currentIndex: number, receiptIds: string[], mode: QueueMode): void {
    this.initQueueAndNavigate(receiptIds, mode, currentIndex + 1);
  }

  public queuePrevious(currentIndex: number, receiptIds: string[], mode: QueueMode): void {
    this.initQueueAndNavigate(receiptIds, mode, currentIndex - 1);
  }
}
