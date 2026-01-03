import { Component } from "@angular/core";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { take, tap } from "rxjs";
import { PagedTableInterface } from "../../interfaces/paged-table.interface";
import { BaseTableService } from "../../services/base-table.service";
import { TableColumn } from "../../table/table-column.interface";

@Component({
    selector: "app-base-table",
    templateUrl: "./base-table.component.html",
    styleUrl: "./base-table.component.scss",
    standalone: false
})
export class BaseTableComponent<T> {
  public columns: TableColumn[] = [];

  public displayedColumns: string[] = [];

  public dataSource = new MatTableDataSource<T>([]);

  public totalCount = 0;

  constructor(public baseTableService: BaseTableService) {}

  public setInitialSortedColumn(state: PagedTableInterface, columns: TableColumn[]): void {
    if (state.orderBy) {
      const column = columns.find((c) => c.matColumnDef === state.orderBy);
      if (column) {
        column.defaultSortDirection = state.sortDirection;
      }
    }
  }

  public getTableData(): void {
    this.baseTableService
      .getPagedData()
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource = new MatTableDataSource(pagedData.data as T[]) as MatTableDataSource<T>;
          this.totalCount = pagedData.totalCount;
        })
      )
      .subscribe();
  }

  public sorted(sort: Sort): void {
    this.baseTableService.setOrderBy(sort);
    this.baseTableService.setSortDirection(sort.direction);

    this.getTableData();
  }

  public pageChanged(pageEvent: PageEvent): void {
    const newPage = pageEvent.pageIndex + 1;

    this.baseTableService.setPage(newPage);
    this.baseTableService.setPageSize(pageEvent.pageSize);

    this.getTableData();
  }
}
