import { Injectable, inject } from "@angular/core";
import { Sort, SortDirection } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { PagedData, PagedRequestCommand } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { ReportTemplateTableState } from "../../store/report-template-table.state";
import {
  SetOrderBy,
  SetPage,
  SetPageSize,
  SetSortDirection,
} from "../../store/report-template-table.state.actions";
import { ReportRunnerService } from "../services/report-runner.service";

/** Feeds the report-template list its paged data + page/sort state (mirrors GroupTableService). */
@Injectable({ providedIn: "root" })
export class ReportTemplateTableService extends BaseTableService {
  private readonly store = inject(Store);
  private readonly runner = inject(ReportRunnerService);

  override page$: Observable<number> = this.store.select(ReportTemplateTableState.page);

  override pageSize$: Observable<number> = this.store.select(ReportTemplateTableState.pageSize);

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
    return this.store.selectSnapshot(ReportTemplateTableState.state);
  }

  public getPagedData(): Observable<PagedData> {
    return this.runner.listTemplates(this.getPagedRequestCommand());
  }
}
