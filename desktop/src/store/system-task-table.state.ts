import { Injectable } from "@angular/core";
import { Action, createSelector, Selector, State, StateContext } from "@ngxs/store";
import { SystemTaskTableInterface } from "src/interfaces/system-task-table.interface";
import { FilterOperation, SortDirection } from "../open-api";
import {
  ResetSystemTaskFilter,
  SetOrderBy,
  SetPage,
  SetPageSize,
  SetSortDirection,
  SetSystemTaskFilter,
} from "./system-task-table.state.actions";

export const defaultSystemTaskFilter = {
  type: {
    operation: null,
    value: [],
  },
  status: {
    operation: null,
    value: [],
  },
  ranByUserId: {
    operation: null,
    value: [],
  },
  startedAt: {
    operation: null,
    value: null,
  },
  endedAt: {
    operation: null,
    value: null,
  },
};


@State<SystemTaskTableInterface>({
  name: "systemTaskTable",
  defaults: {
    page: 1,
    pageSize: 50,
    orderBy: "started_at",
    sortDirection: "desc" as SortDirection,
    filter: defaultSystemTaskFilter,
  },
})
@Injectable()
export class SystemTaskTableState {
  static get state() {
    return createSelector([this], (state: SystemTaskTableInterface) => {
      return state;
    });
  }

  static get page() {
    return createSelector([this], (state: SystemTaskTableInterface) => {
      return state.page;
    });
  }

  static get pageSize() {
    return createSelector([this], (state: SystemTaskTableInterface) => {
      return state.pageSize;
    });
  }

  @Selector()
  static filterData(state: SystemTaskTableInterface): SystemTaskTableInterface {
    return state;
  }

  @Selector()
  static numFiltersApplied(state: SystemTaskTableInterface): number {
    let filtersApplied = 0;
    const filter: any = state.filter;

    Object.keys(filter).forEach((key) => {
      const stringValue = filter[key]?.value?.toString();
      const operationValue = filter[key]?.operation?.toString();
      if (stringValue?.length > 0 && stringValue !== "0") {
        filtersApplied += 1;
      } else if (operationValue === FilterOperation.WithinCurrentMonth) {
        filtersApplied += 1;
      }
    });

    return filtersApplied;
  }

  @Action(SetPage)
  setPage({ patchState }: StateContext<SystemTaskTableInterface>, payload: SetPage) {
    patchState({
      page: payload.page,
    });
  }

  @Action(SetPageSize)
  setPageSize(
    { patchState }: StateContext<SystemTaskTableInterface>,
    payload: SetPageSize
  ) {
    patchState({
      pageSize: payload.pageSize,
    });
  }

  @Action(SetOrderBy)
  setOrderBy(
    { patchState }: StateContext<SystemTaskTableInterface>,
    payload: SetOrderBy
  ) {
    patchState({
      orderBy: payload.orderBy,
    });
  }

  @Action(SetSortDirection)
  setSortDirection(
    { patchState }: StateContext<SystemTaskTableInterface>,
    payload: SetSortDirection
  ) {
    patchState({
      sortDirection: payload.sortDirection,
    });
  }

  @Action(SetSystemTaskFilter)
  setSystemTaskFilter(
    { patchState }: StateContext<SystemTaskTableInterface>,
    payload: SetSystemTaskFilter
  ) {
    patchState({
      filter: payload.data,
    });
  }

  @Action(ResetSystemTaskFilter)
  resetFilter({ patchState }: StateContext<SystemTaskTableInterface>) {
    patchState({
      filter: defaultSystemTaskFilter,
    });
  }
}
