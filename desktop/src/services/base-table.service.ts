import { Injectable } from "@angular/core";
import { Sort } from "@angular/material/sort";
import { Observable } from "rxjs";
import { PagedData, PagedRequestCommand, SortDirection } from "../open-api";

@Injectable()
export abstract class BaseTableService {
  abstract page$: Observable<number>;

  abstract pageSize$: Observable<number>;

  abstract getPagedData(): Observable<PagedData>;

  abstract getPagedRequestCommand(): PagedRequestCommand;

  abstract setPage(page: number): void;

  abstract setPageSize(pageSize: number): void;

  abstract setOrderBy(orderBy: Sort): void;

  abstract setSortDirection(sortDirection: SortDirection): void;
}
