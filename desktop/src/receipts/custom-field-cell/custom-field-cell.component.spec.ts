import { CommonModule, CurrencyPipe } from "@angular/common";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { PipesModule } from "../../pipes/pipes.module";
import {
  CurrencySeparator,
  CurrencySymbolPosition,
  CustomField,
  CustomFieldType,
  CustomFieldValue,
  Receipt,
} from "../../open-api";
import { SystemSettingsState } from "../../store/system-settings.state";
import {
  SetCurrencyData,
  SetCurrencyDisplay,
} from "../../store/system-settings.state.actions";
import { CustomFieldCellComponent } from "./custom-field-cell.component";

describe("CustomFieldCellComponent", () => {
  let component: CustomFieldCellComponent;
  let fixture: ComponentFixture<CustomFieldCellComponent>;
  let store: Store;

  const field = (type: CustomFieldType, options?: any[]): CustomField =>
    ({ id: 1, name: "Field", type, options } as CustomField);

  const receiptWith = (...values: Partial<CustomFieldValue>[]): Receipt =>
    ({
      customFields: values.map((value, index) => ({
        id: index + 1,
        customFieldId: 1,
        ...value,
      })),
    } as Receipt);

  const render = (customField: CustomField, receipt: Receipt): string => {
    fixture.componentRef.setInput("customField", customField);
    fixture.componentRef.setInput("receipt", receipt);
    fixture.detectChanges();

    return (fixture.nativeElement as HTMLElement).textContent?.trim() ?? "";
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CustomFieldCellComponent],
      imports: [CommonModule, PipesModule, NgxsModule.forRoot([SystemSettingsState])],
      providers: [CurrencyPipe],
    }).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(CustomFieldCellComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("renders a text value", () => {
    expect(
      render(field(CustomFieldType.Text), receiptWith({ stringValue: "Acme" }))
    ).toBe("Acme");
  });

  it("renders a boolean value as yes or no", () => {
    expect(
      render(field(CustomFieldType.Boolean), receiptWith({ booleanValue: true }))
    ).toBe("Yes");
    expect(
      render(field(CustomFieldType.Boolean), receiptWith({ booleanValue: false }))
    ).toBe("No");
  });

  it("renders a select value as the option's text, not its id", () => {
    const customField = field(CustomFieldType.Select, [
      { id: 4, value: "Approved" },
      { id: 5, value: "Rejected" },
    ]);

    expect(render(customField, receiptWith({ selectValue: 5 }))).toBe("Rejected");
  });

  it("renders nothing when the receipt has no value for the field", () => {
    expect(render(field(CustomFieldType.Text), receiptWith())).toBe("");
  });

  // A receipt may carry several values for one field. The lowest id that
  // actually holds something wins, matching how the API sorts the column and how
  // the reporting engine reads it - so the cell can never disagree with either.
  it("prefers the lowest-id value that is actually set", () => {
    const receipt = receiptWith(
      { stringValue: undefined },
      { stringValue: "winner" },
      { stringValue: "loser" }
    );

    expect(render(field(CustomFieldType.Text), receipt)).toBe("winner");
  });

  describe("currency", () => {
    it("formats through the configured currency display", () => {
      store.dispatch(new SetCurrencyDisplay("€"));
      store.dispatch(
        new SetCurrencyData(
          CurrencySymbolPosition.End,
          CurrencySeparator.Comma,
          CurrencySeparator.Period,
          false
        )
      );

      const rendered = render(
        field(CustomFieldType.Currency),
        receiptWith({ currencyValue: "1234.56" })
      );

      // Not "$1,234.56": the separators and the symbol position are settings, and
      // a currency custom field has to honour them exactly as the Amount column
      // does.
      expect(rendered).toBe("1.234,56€");
    });

    it("uses the default display when nothing is configured", () => {
      expect(
        render(field(CustomFieldType.Currency), receiptWith({ currencyValue: "12.34" }))
      ).toBe("$12.34");
    });
  });
});
