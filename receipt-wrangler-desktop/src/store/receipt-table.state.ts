import { Injectable } from "@angular/core";
import { Action, Selector, State, StateContext } from "@ngxs/store";
import { ReceiptTableInterface } from "../interfaces";
import { DEFAULT_RECEIPT_TABLE_COLUMNS, ReceiptTableColumnConfig } from "../interfaces/receipt-table-column-config.interface";
import { FilterOperation, ReceiptPagedRequestFilter } from "../open-api";
import { ResetReceiptFilter, SetColumnConfig, SetPage, SetPageSize, SetReceiptFilter, SetReceiptFilterData } from "./receipt-table.actions";

export const defaultReceiptFilter = {
  date: {
    operation: null,
    value: null
  },
  amount: {
    operation: null,
    value: null,
  },
  name: {
    operation: null,
    value: null,
  },
  paidBy: {
    operation: null,
    value: [],
  },
  categories: {
    operation: null,
    value: [],
  },
  tags: {
    operation: null,
    value: [],
  },
  status: {
    operation: null,
    value: [],
  },
  resolvedDate: {
    operation: null,
    value: null,
  },
  createdAt: {
    operation: null,
    value: null,
  },
} as ReceiptPagedRequestFilter;

// TODO: look into fixing date equals
@State<ReceiptTableInterface>({
  name: "receiptTable",
  defaults: {
    page: 1,
    pageSize: 50,
    orderBy: "created_at",
    sortDirection: "desc",
    filter: defaultReceiptFilter,
    columnConfig: DEFAULT_RECEIPT_TABLE_COLUMNS,
  },
})
@Injectable()
export class ReceiptTableState {
  @Selector()
  static page(state: ReceiptTableInterface): number {
    return state.page;
  }

  @Selector()
  static pageSize(state: ReceiptTableInterface): number {
    return state.pageSize;
  }

  @Selector()
  static filterData(state: ReceiptTableInterface): ReceiptTableInterface {
    return state;
  }

  @Selector()
  static numFiltersApplied(state: ReceiptTableInterface): number {
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

  @Selector()
  static columnConfig(state: ReceiptTableInterface): ReceiptTableColumnConfig[] {
    return state.columnConfig || DEFAULT_RECEIPT_TABLE_COLUMNS;
  }

  @Action(SetPage)
  setPage(
    { patchState }: StateContext<ReceiptTableInterface>,
    payload: SetPage
  ) {
    patchState({
      page: payload.page,
    });
  }

  @Action(SetPageSize)
  setPageSize(
    { patchState }: StateContext<ReceiptTableInterface>,
    payload: SetPageSize
  ) {
    patchState({
      pageSize: payload.pageSize,
    });
  }

  @Action(SetReceiptFilterData)
  setReceiptFilterData(
    { patchState }: StateContext<ReceiptTableInterface>,
    payload: SetReceiptFilterData
  ) {
    patchState(payload.data);
  }

  @Action(SetReceiptFilter)
  setReceiptFilter(
    { patchState }: StateContext<ReceiptTableInterface>,
    payload: SetReceiptFilter
  ) {
    patchState({
      filter: payload.data,
    });
  }

  @Action(ResetReceiptFilter)
  resetFilter({ patchState }: StateContext<ReceiptTableInterface>) {
    patchState({
      filter: defaultReceiptFilter,
    });
  }

  @Action(SetColumnConfig)
  setColumnConfig(
    { patchState }: StateContext<ReceiptTableInterface>,
    payload: SetColumnConfig
  ) {
    patchState({
      columnConfig: payload.columnConfig,
    });
  }
}
