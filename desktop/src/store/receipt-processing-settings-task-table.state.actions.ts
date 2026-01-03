import { SortDirection } from "@angular/material/sort";

export class SetPage {
  static readonly type = "[ReceiptProcessingSettingsTaskTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[ReceiptProcessingSettingsTaskTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[ReceiptProcessingSettingsTaskTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[ReceiptProcessingSettingsTaskTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}
