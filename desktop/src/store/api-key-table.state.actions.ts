import { SortDirection } from "@angular/material/sort";
import { ApiKeyFilter } from "../open-api";

export class SetPage {
  static readonly type = "[ApiKeyTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[ApiKeyTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[ApiKeyTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[ApiKeyTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}

export class SetFilter {
  static readonly type = "[ApiKeyTableComponent] Set Filter";

  constructor(public filter: ApiKeyFilter) {}
}