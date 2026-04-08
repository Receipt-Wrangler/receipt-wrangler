import { SortDirection } from "@angular/material/sort";
import { SystemTaskPagedRequestFilter } from "../open-api";

export interface SystemTaskTableInterface {
  page: number;
  pageSize: number;
  orderBy: string;
  sortDirection: SortDirection;
  filter: SystemTaskPagedRequestFilter;
}
