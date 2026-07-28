import { SortDirection } from "@angular/material/sort";

export class SetPage {
  static readonly type = "[UserListComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[UserListComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[UserListComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[UserListComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}
