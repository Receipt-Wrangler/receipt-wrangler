import { CurrencyPipe } from "@angular/common";
import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { CurrencySeparator, CurrencySymbolPosition } from "../open-api/index";
import { SystemSettingsState } from "../store/system-settings.state";
import { CustomCurrencyPipe } from "./custom-currency.pipe";

describe("CustomCurrencyPipe", () => {
  let pipe: CustomCurrencyPipe;
  let store: Store;
  let currencyPipe: CurrencyPipe;

  beforeEach(() => {

    TestBed.configureTestingModule({
      declarations: [CustomCurrencyPipe],
      imports: [NgxsModule.forRoot([SystemSettingsState])],
      providers: [
        CurrencyPipe,
        CustomCurrencyPipe
      ]
    });

    pipe = TestBed.inject(CustomCurrencyPipe);
    store = TestBed.inject(Store);
    currencyPipe = TestBed.inject(CurrencyPipe);
  });

  it("should create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should use default system settings when no parameters are provided", () => {
    store.reset({
      systemSettings: {
        currencyDisplay: "€",
        currencyDecimalSeparator: CurrencySeparator.Comma,
        currencyThousandthsSeparator: CurrencySeparator.Period,
        currencySymbolPosition: CurrencySymbolPosition.Start,
        currencyHideDecimalPlaces: false
      }
    });
    const result = pipe.transform(1234.56);
    expect(result).toBe("€1.234,56");
  });

  it("should use provided parameters instead of system settings", () => {
    store.reset({
      systemSettings: {
        currencyDisplay: "€",
        currencyDecimalSeparator: CurrencySeparator.Comma,
        currencyThousandthsSeparator: CurrencySeparator.Period,
        currencySymbolPosition: CurrencySymbolPosition.Start,
        currencyHideDecimalPlaces: false
      }
    });
    const result = pipe.transform(1234.56, "£", CurrencySeparator.Period, CurrencySeparator.Comma, CurrencySymbolPosition.End, true);
    expect(result).toBe("1,234£");
  });

  it("should hide decimal places when specified", () => {
    store.reset({
      systemSettings: {
        currencyDisplay: "$",
        currencyDecimalSeparator: CurrencySeparator.Period,
        currencyThousandthsSeparator: CurrencySeparator.Comma,
        currencySymbolPosition: CurrencySymbolPosition.Start,
        currencyHideDecimalPlaces: true
      }
    });
    const result = pipe.transform(1234.56);
    expect(result).toBe("$1,234");
  });

  it("should handle different currency symbol positions", () => {
    store.reset({
      systemSettings: {
        currencyDisplay: "€",
        currencyDecimalSeparator: CurrencySeparator.Comma,
        currencyThousandthsSeparator: CurrencySeparator.Period,
        currencySymbolPosition: CurrencySymbolPosition.End,
        currencyHideDecimalPlaces: false
      }
    });
    const result = pipe.transform(1234.56);
    expect(result).toBe("1.234,56€");
  });

  it("should handle zero values", () => {
    store.reset({
      systemSettings: {
        currencyDisplay: "$",
        currencyDecimalSeparator: CurrencySeparator.Period,
        currencyThousandthsSeparator: CurrencySeparator.Comma,
        currencySymbolPosition: CurrencySymbolPosition.Start,
        currencyHideDecimalPlaces: false
      }
    });
    const result = pipe.transform(0);
    expect(result).toBe("$0.00");
  });

  it("should handle negative values", () => {
    store.reset({
      systemSettings: {
        currencyDisplay: "$",
        currencyDecimalSeparator: CurrencySeparator.Period,
        currencyThousandthsSeparator: CurrencySeparator.Comma,
        currencySymbolPosition: CurrencySymbolPosition.Start,
        currencyHideDecimalPlaces: false
      }
    });
    const result = pipe.transform(-1234.56);
    expect(result).toBe("$-1,234.56");
  });
});
