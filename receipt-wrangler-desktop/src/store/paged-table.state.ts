import { createSelector } from "@ngxs/store";
import { PagedTableInterface } from "src/interfaces/paged-table.interface";

export class PagedTableState {
  static get state() {
    return createSelector([this], (state: PagedTableInterface) => {
      return state;
    });
  }

  static get page() {
    return createSelector([this], (state: PagedTableInterface) => {
      return state.page;
    });
  }

  static get pageSize() {
    return createSelector([this], (state: PagedTableInterface) => {
      return state.pageSize;
    });
  }
}
