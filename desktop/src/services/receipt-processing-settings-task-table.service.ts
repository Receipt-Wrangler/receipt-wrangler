import { Injectable } from "@angular/core";
import { Sort } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { PagedRequestCommand, SortDirection } from "src/open-api";
import { ReceiptProcessingSettingsTaskTableState } from "../store/receipt-processing-settings-task-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "../store/receipt-processing-settings-task-table.state.actions";
import { BaseTableService } from "./base-table.service";

@Injectable()
export class ReceiptProcessingSettingsTaskTableService extends BaseTableService {
  override page$: Observable<number>;
  override pageSize$: Observable<number>;

  constructor(private store: Store) {
    super();
    this.page$ = this.store.select(ReceiptProcessingSettingsTaskTableState.page);
    this.pageSize$ = this.store.select(ReceiptProcessingSettingsTaskTableState.pageSize);
  }


  public getPagedRequestCommand(): PagedRequestCommand {
    return this.store.selectSnapshot(ReceiptProcessingSettingsTaskTableState.state);
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

  public getPagedData(): Observable<any> {
    throw new Error("Method not implemented.");
  }
}
