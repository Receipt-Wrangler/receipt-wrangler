import { Injectable } from "@angular/core";
import { Sort } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { PagedRequestCommand, SortDirection } from "src/open-api";
import { SystemTaskTableState } from "../store/system-task-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "../store/system-task-table.state.actions";
import { BaseTableService } from "./base-table.service";

@Injectable()
export class SystemTaskTableService extends BaseTableService {
  override page$: Observable<number>;
  override pageSize$: Observable<number>;

  constructor(private store: Store) {
    super();
    this.page$ = this.store.select(SystemTaskTableState.page);
    this.pageSize$ = this.store.select(SystemTaskTableState.pageSize);
  }


  public getPagedRequestCommand(): PagedRequestCommand {
    return this.store.selectSnapshot(SystemTaskTableState.state);
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
