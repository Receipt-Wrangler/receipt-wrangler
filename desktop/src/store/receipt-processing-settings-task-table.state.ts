import { Injectable } from "@angular/core";
import { Action, State, StateContext } from "@ngxs/store";
import { PagedTableInterface } from "src/interfaces/paged-table.interface";
import { SortDirection } from "../open-api";
import { PagedTableState } from "./paged-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "./receipt-processing-settings-task-table.state.actions";


@State<PagedTableInterface>({
  name: "receiptProcessingSettingsTaskTable",
  defaults: {
    page: 1,
    pageSize: 50,
    orderBy: "type",
    sortDirection: SortDirection.Desc,
  },
})
@Injectable()
export class ReceiptProcessingSettingsTaskTableState extends PagedTableState {
  @Action(SetPage)
  setPage({ patchState }: StateContext<PagedTableInterface>, payload: SetPage) {
    patchState({
      page: payload.page,
    });
  }

  @Action(SetPageSize)
  setPageSize(
    { patchState }: StateContext<PagedTableInterface>,
    payload: SetPageSize
  ) {
    patchState({
      pageSize: payload.pageSize,
    });
  }

  @Action(SetOrderBy)
  setOrderBy(
    { patchState }: StateContext<PagedTableInterface>,
    payload: SetOrderBy
  ) {
    patchState({
      orderBy: payload.orderBy,
    });
  }

  @Action(SetSortDirection)
  setSortDirection(
    { patchState }: StateContext<PagedTableInterface>,
    payload: SetSortDirection
  ) {
    patchState({
      sortDirection: payload.sortDirection,
    });
  }
}
