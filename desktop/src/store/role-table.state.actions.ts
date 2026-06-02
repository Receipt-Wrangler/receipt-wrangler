import { SortDirection } from "@angular/material/sort";
import { RoleListFilter } from "src/roles/role-list/role-list-item.interface";

export class SetPage {
  static readonly type = "[RoleTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[RoleTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[RoleTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[RoleTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}

export class SetScope {
  static readonly type = "[RoleTableComponent] Set Scope";

  constructor(public scope: RoleListFilter) {}
}
