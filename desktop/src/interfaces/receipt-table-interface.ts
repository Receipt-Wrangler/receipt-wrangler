import { SortDirection } from "@angular/material/sort";
import { ReceiptPagedRequestFilter } from "../open-api";
import { ReceiptTableColumnConfig } from "./receipt-table-column-config.interface";

export interface ReceiptTableInterface {
  page: number;
  pageSize: number;
  orderBy: string;
  sortDirection: SortDirection;
  filter: ReceiptPagedRequestFilter;
  columnConfig?: ReceiptTableColumnConfig[];
}
