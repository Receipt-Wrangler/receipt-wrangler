import { SortDirection } from "@angular/material/sort";
import { GroupFilter } from "../open-api";

export class SetPage {
  static readonly type = "[GroupTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[GroupTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[GroupTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[GroupTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}

export class SetFilter {
  static readonly type = "[GroupTableComponent] Set Filter";

  constructor(public filter: GroupFilter) {}
}

