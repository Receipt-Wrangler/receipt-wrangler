import { CurrencyPipe } from "@angular/common";
import { Pipe, PipeTransform } from "@angular/core";
import { Store } from "@ngxs/store";
import { CurrencySeparator, CurrencySymbolPosition } from "../open-api/index";
import { SystemSettingsState } from "../store/system-settings.state";

@Pipe({
    name: "customCurrency",
    standalone: false
})
export class CustomCurrencyPipe implements PipeTransform {
  constructor(private store: Store, private currencyPipe: CurrencyPipe) {}

  public transform(
    value: string | number,
    currencySymbol?: string,
    currencyDecimalSeparator?: CurrencySeparator,
    currencyThousandthsSeparator?: CurrencySeparator,
    currencySymbolPosition?: CurrencySymbolPosition,
    currencyHideDecimalPlaces?: boolean
  ): string {
    const systemSettingsState = this.store.selectSnapshot(SystemSettingsState.state);
    let currencyValue = this.currencyPipe.transform(value, "USD", "symbol", undefined, "en-US") ?? "";
    const result = currencyValue.split("");
    const originalCurrencyValueArray = currencyValue.split("");

    const currencySymbolToUse = currencySymbol || systemSettingsState.currencyDisplay;
    const currencyDecimalSeparatorToUse = currencyDecimalSeparator || systemSettingsState?.currencyDecimalSeparator;
    const currencyThousandthsSeparatorToUse = currencyThousandthsSeparator || systemSettingsState.currencyThousandthsSeparator;
    const currencySymbolPositionToUse = currencySymbolPosition || systemSettingsState.currencySymbolPosition;
    const currencyHideDecimalPlacesToUse = currencyHideDecimalPlaces === undefined ? systemSettingsState.currencyHideDecimalPlaces : currencyHideDecimalPlaces;

    if (currencyHideDecimalPlacesToUse) {
      const decimalIndex = result.indexOf(".");
      result.splice(decimalIndex, result.length - decimalIndex);
    }

    if (currencyDecimalSeparatorToUse && !currencyHideDecimalPlacesToUse) {
      for (let i = 0; i < result.length; i++) {
        if (result[i] === CurrencySeparator.Period) {
          result[i] = currencyDecimalSeparatorToUse;
        }
      }
    }

    if (currencyThousandthsSeparatorToUse) {
      const decimalIndex = originalCurrencyValueArray.indexOf(".");
      for (let i = 0; i < (decimalIndex || result.length); i++) {
        if (result[i] === ",") {
          result[i] = currencyThousandthsSeparatorToUse;
        }
      }
    }


    if (currencySymbolToUse) {
      const index = result.findIndex((v => v === "$"));
      result.splice(index, 1);

      if (currencySymbolPositionToUse === CurrencySymbolPosition.Start) {
        result.unshift(currencySymbolToUse);
      }

      if (currencySymbolPositionToUse === CurrencySymbolPosition.End) {
        result.push(currencySymbolToUse);
      }
    }

    return result.join("");
  }
}
