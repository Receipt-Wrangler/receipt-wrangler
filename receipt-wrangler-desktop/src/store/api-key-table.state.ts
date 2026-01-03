import { Injectable } from "@angular/core";
import { Action, Selector, State, StateContext } from "@ngxs/store";
import { AssociatedApiKeys, ApiKeyFilter, PagedApiKeyRequestCommand, SortDirection } from "../open-api";
import { SetFilter, SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "./api-key-table.state.actions";
import { PagedTableState } from "./paged-table.state";

@State<PagedApiKeyRequestCommand>({
  name: "apiKeyTable",
  defaults: {
    page: 1,
    pageSize: 50,
    orderBy: "name",
    sortDirection: SortDirection.Desc,
    filter: {
      associatedApiKeys: AssociatedApiKeys.Mine
    } as ApiKeyFilter
  },
})
@Injectable()
export class ApiKeyTableState extends PagedTableState {

  @Selector()
  static filter(state: PagedApiKeyRequestCommand): ApiKeyFilter {
    return state.filter ?? {};
  }


  @Action(SetPage)
  setPage({ patchState }: StateContext<PagedApiKeyRequestCommand>, payload: SetPage) {
    patchState({
      page: payload.page,
    });
  }

  @Action(SetPageSize)
  setPageSize(
    { patchState }: StateContext<PagedApiKeyRequestCommand>,
    payload: SetPageSize
  ) {
    patchState({
      pageSize: payload.pageSize,
    });
  }

  @Action(SetOrderBy)
  setOrderBy(
    { patchState }: StateContext<PagedApiKeyRequestCommand>,
    payload: SetOrderBy
  ) {
    patchState({
      orderBy: payload.orderBy,
    });
  }

  @Action(SetSortDirection)
  setSortDirection(
    { patchState }: StateContext<PagedApiKeyRequestCommand>,
    payload: SetSortDirection
  ) {
    patchState({
      sortDirection: payload.sortDirection,
    });
  }

  @Action(SetFilter)
  setFilter(
    { patchState }: StateContext<PagedApiKeyRequestCommand>,
    payload: SetFilter
  ) {
    patchState({
      filter: payload.filter,
    });
  }
}