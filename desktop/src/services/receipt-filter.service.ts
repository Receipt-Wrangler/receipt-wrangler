import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { SortDirection } from "@angular/material/sort";
import { Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { ReceiptTableState } from "src/store/receipt-table.state";
import { PagedData, ReceiptPagedRequestCommand } from "../open-api";

@Injectable({
  providedIn: "root",
})
export class ReceiptFilterService {
  constructor(private store: Store, private httpClient: HttpClient) {}

  public getPagedReceiptsForGroups(
    groupId: string,
    page?: number,
    pageSize?: number,
    orderBy?: string,
    sortDirection?: SortDirection,
    pagedRequestCommand?: ReceiptPagedRequestCommand
  ): Observable<PagedData> {
    const filterData = this.buildPagedRequestCommand(
      page,
      pageSize,
      orderBy,
      sortDirection,
      pagedRequestCommand,
    );

    return this.httpClient.post<PagedData>(
      `/api/receipt/group/${groupId}`,
      filterData
    );
  }

  public buildPagedRequestCommand(
    page?: number,
    pageSize?: number,
    orderBy?: string,
    sortDirection?: SortDirection,
    pagedRequestCommand?: ReceiptPagedRequestCommand
  ): ReceiptPagedRequestCommand {
    let filterData: ReceiptPagedRequestCommand;

    if (pagedRequestCommand) {
      filterData = pagedRequestCommand;
    } else {
      const filter = this.store.selectSnapshot(ReceiptTableState.filterData);
      filterData = {
        page: page ?? filter.page,
        pageSize: pageSize ?? filter.pageSize,
        orderBy: orderBy ?? filter.orderBy,
        sortDirection: sortDirection ?? filter.sortDirection,
        filter: Object.assign(filter.filter, {}),
      };
    }

    if ((!filterData?.filter as any)?.date?.value && filterData?.filter?.date) {
      (filterData.filter as any).date.value = "";
    }

    if (
      (!filterData?.filter as any)?.resolvedDate?.value &&
      filterData?.filter?.resolvedDate
    ) {
      (filterData.filter as any).resolvedDate.value = "";
    }

    if ((!filterData?.filter as any)?.amount?.value && filterData?.filter?.amount) {
      (filterData.filter as any).amount.value = 0;
    }

    return filterData;
  }
}
