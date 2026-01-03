import { Injectable } from "@angular/core";
import { Sort, SortDirection } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { PagedData, PagedRequestCommand, ReceiptProcessingSettingsService } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { ReceiptProcessingSettingsTableState } from "../../store/receipt-processing-settings-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "../../store/receipt-processing-settings-table.state.actions";

@Injectable({
  providedIn: "root"
})
export class ReceiptProcessingSettingsTableService extends BaseTableService {
  override page$: Observable<number>;

  override pageSize$: Observable<number>;

  constructor(
    private store: Store,
    private receiptProcessingSettingsService: ReceiptProcessingSettingsService
  ) {
    super();

    this.page$ = this.store.select(ReceiptProcessingSettingsTableState.page);
    this.pageSize$ = this.store.select(ReceiptProcessingSettingsTableState.pageSize);
  }

  public setPage(page: number): void {
    this.store.dispatch(new SetPage(page));
  }

  public setPageSize(pageSize: number): void {
    this.store.dispatch(new SetPageSize(pageSize));
  }

  public setOrderBy(orderBy: Sort): void {
    this.store.dispatch(new SetOrderBy(orderBy.active));
  }

  public setSortDirection(sortDirection: SortDirection): void {
    this.store.dispatch(new SetSortDirection(sortDirection));
  }

  public getPagedRequestCommand(): PagedRequestCommand {
    return this.store.selectSnapshot(ReceiptProcessingSettingsTableState.state);
  }

  public getPagedData(): Observable<PagedData> {
    return this.receiptProcessingSettingsService.getPagedProcessingSettings(this.getPagedRequestCommand());
  }
}
