import { SortDirection } from "@angular/material/sort";

export class SetPage {
  static readonly type = "[SystemEmailTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[SystemEmailTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[SystemEmailTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[SystemEmailTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}
