import { SortDirection } from '@angular/material/sort';

export interface PagedTableInterface {
  page: number;
  pageSize: number;
  orderBy: string;
  sortDirection: SortDirection;
}
