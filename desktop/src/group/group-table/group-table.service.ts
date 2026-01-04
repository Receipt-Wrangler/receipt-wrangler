import { Injectable } from "@angular/core";
import { Sort, SortDirection } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { GroupsService, PagedData, PagedRequestCommand } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { GroupTableState } from "../../store/group-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "../../store/group-table.state.actions";

@Injectable({
  providedIn: "root"
})
export class GroupTableService extends BaseTableService {
  override page$: Observable<number>;

  override pageSize$: Observable<number>;

  constructor(
    private store: Store,
    private groupService: GroupsService
  ) {
    super();

    this.page$ = this.store.select(GroupTableState.page);
    this.pageSize$ = this.store.select(GroupTableState.pageSize);
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
    return this.store.selectSnapshot(GroupTableState.state);
  }

  public getPagedData(): Observable<PagedData> {
    return this.groupService.getPagedGroups(this.getPagedRequestCommand());
  }
}
