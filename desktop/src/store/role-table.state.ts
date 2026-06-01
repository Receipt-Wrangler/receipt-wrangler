import { Injectable } from "@angular/core";
import { Action, Selector, State, StateContext } from "@ngxs/store";
import { RoleTableInterface } from "src/interfaces/role-table.interface";
import { RoleListFilter } from "src/roles/role-list/role-list-item.interface";
import {
  SetOrderBy,
  SetPage,
  SetPageSize,
  SetScope,
  SetSortDirection,
} from "./role-table.state.actions";
import { PagedTableState } from "./paged-table.state";

@State<RoleTableInterface>({
  name: "roleTable",
  defaults: {
    page: 1,
    pageSize: 50,
    orderBy: "name",
    sortDirection: "asc",
    scope: "all",
  },
})
@Injectable()
export class RoleTableState extends PagedTableState {
  // The paging fields are read through the inherited `state` selector; `scope`
  // backs the filter bar and is exposed separately (it is not part of the
  // shared paged-table shape).
  @Selector()
  static scope(state: RoleTableInterface): RoleListFilter {
    return state.scope;
  }

  @Action(SetPage)
  setPage(
    { patchState }: StateContext<RoleTableInterface>,
    payload: SetPage,
  ) {
    patchState({ page: payload.page });
  }

  @Action(SetPageSize)
  setPageSize(
    { patchState }: StateContext<RoleTableInterface>,
    payload: SetPageSize,
  ) {
    patchState({ pageSize: payload.pageSize });
  }

  @Action(SetOrderBy)
  setOrderBy(
    { patchState }: StateContext<RoleTableInterface>,
    payload: SetOrderBy,
  ) {
    patchState({ orderBy: payload.orderBy });
  }

  @Action(SetSortDirection)
  setSortDirection(
    { patchState }: StateContext<RoleTableInterface>,
    payload: SetSortDirection,
  ) {
    patchState({ sortDirection: payload.sortDirection });
  }

  // Changing the scope filter resets to the first page so the user never lands
  // on an out-of-range page of the newly filtered (smaller) result set.
  @Action(SetScope)
  setScope(
    { patchState }: StateContext<RoleTableInterface>,
    payload: SetScope,
  ) {
    patchState({ scope: payload.scope, page: 1 });
  }
}
