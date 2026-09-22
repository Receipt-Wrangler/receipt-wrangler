import { SortDirection } from "@angular/material/sort";
import { ReceiptDateFilterFieldKey } from "../constants/receipt-filter-fields.constant";
import { ReceiptPagedRequestFilter } from "../open-api";
import { ReceiptTableColumnConfig } from "./receipt-table-column-config.interface";

export interface ReceiptTableInterface {
  page: number;
  pageSize: number;
  orderBy: string;
  sortDirection: SortDirection;
  filter: ReceiptPagedRequestFilter;
  /**
   * The date field the quick date control writes to. Optional because state
   * persisted before it existed deserializes without it — read it through
   * `ReceiptTableState.quickDateField`, which supplies the fallback.
   */
  quickDateField?: ReceiptDateFilterFieldKey;
  columnConfig?: ReceiptTableColumnConfig[];
  /**
   * Which group's configuration shapes the receipt summary on the synthetic "All"
   * group, where the receipts span several groups and none of their settings is
   * the obvious one to use.
   *
   * Optional for the same reason as `quickDateField` — state persisted before this
   * key existed deserializes without it. Unlike `quickDateField` the fallback
   * cannot live in a selector, because it depends on which groups the user belongs
   * to and which have the summary enabled; `resolveSummaryConfigGroup`
   * (`src/utils/receipt-summary.ts`) owns it instead.
   */
  summaryConfigGroupId?: number;
}
