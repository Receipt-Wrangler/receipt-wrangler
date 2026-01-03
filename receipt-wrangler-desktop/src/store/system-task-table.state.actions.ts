import { SortDirection } from "@angular/material/sort";

export class SetPage {
  static readonly type = "[SystemTaskTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[SystemTaskTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[SystemTaskTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[SystemTaskTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}
