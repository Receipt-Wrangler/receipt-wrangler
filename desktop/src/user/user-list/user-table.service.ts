import { Injectable } from "@angular/core";
import { Sort, SortDirection } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { PagedData, PagedRequestCommand, UserService } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { UserTableState } from "../../store/user-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "../../store/user-table.state.actions";

@Injectable({
  providedIn: "root"
})
export class UserTableService extends BaseTableService {
  override page$: Observable<number>;

  override pageSize$: Observable<number>;

  constructor(
    private store: Store,
    private userService: UserService
  ) {
    super();

    this.page$ = this.store.select(UserTableState.page);
    this.pageSize$ = this.store.select(UserTableState.pageSize);
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
    return this.store.selectSnapshot(UserTableState.state);
  }

  public getPagedData(): Observable<PagedData> {
    return this.userService.getPagedUsers(this.getPagedRequestCommand());
  }
}
