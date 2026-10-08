import { CurrencySeparator, CurrencySymbolPosition } from "../open-api/index";

export class SetCurrencyDisplay {
  static readonly type = "[SystemSettingsState] Set Currency Display";

  constructor(public currencyDisplay: string) {}
}

export class SetCurrencyData {
  static readonly type = "[SystemSettingsState] Set Currency Data";

  constructor(
    public currencySymbolPosition: CurrencySymbolPosition,
    public currencyDecimalSeparator: CurrencySeparator,
    public currencyThousandthsSeparator: CurrencySeparator,
    public currencyHideDecimalPlaces: boolean
  ) {}
}


export class SetTimeZone {
  static readonly type = "[SystemSettingsState] Set Time Zone";

  constructor(public timeZone: string) {}
}
