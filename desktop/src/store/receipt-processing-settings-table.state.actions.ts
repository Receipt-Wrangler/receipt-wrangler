import { SortDirection } from "@angular/material/sort";

export class SetPage {
  static readonly type = "[ReceiptProcessingSettingsTableComponent] Set Page";

  constructor(public page: number) {}
}

export class SetPageSize {
  static readonly type = "[ReceiptProcessingSettingsTableComponent] Set Page Size";

  constructor(public pageSize: number) {}
}

export class SetOrderBy {
  static readonly type = "[ReceiptProcessingSettingsTableComponent] Set Order By";

  constructor(public orderBy: string) {}
}

export class SetSortDirection {
  static readonly type = "[ReceiptProcessingSettingsTableComponent] Set Sort Direction";

  constructor(public sortDirection: SortDirection) {}
}
