import { Component, Input, OnChanges, SimpleChanges, ViewChild, } from "@angular/core";
import { Select, Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { BaseInputComponent } from "../../base-input";
import { CurrencySeparator, CurrencySymbolPosition } from "../../open-api/index";
import { SystemSettingsState } from "../../store/system-settings.state";
import { InputInterface } from "../input.interface";

@Component({
  selector: "app-input",
  templateUrl: "./input.component.html",
  styleUrls: ["./input.component.scss"],
  standalone: false
})
export class InputComponent
  extends BaseInputComponent
  implements InputInterface, OnChanges {
  @ViewChild("nativeInput") public nativeInput!: { nativeElement: HTMLElement };

  @Select(SystemSettingsState.currencyDisplay) public currencyDisplay!: Observable<string>;

  @Select(SystemSettingsState.currencyDecimalSeparator) public currencyDecimalSeparator!: Observable<"." | ",">;

  @Select(SystemSettingsState.currencyThousandthsSeparator) public currencyThousandthsSeparator!: Observable<"." | ",">;

  @Select(SystemSettingsState.currencySymbolPosition) public currencySymbolPosition!: Observable<string>;

  @Input() public inputId: string = "";

  @Input() public type: string = "text";

  @Input() public showVisibilityEye = false;

  @Input() public isCurrency: boolean = false;

  @Input() public mask: string = "";

  @Input() public maskPrefix: string = "";

  @Input() public maskSuffix: string = "";

  @Input() public thousandSeparator: string = "";

  @Input() public decimalMarker: CurrencySeparator = CurrencySeparator.Period;

  constructor(private store: Store) {
    super();
  }


  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["showVisibilityEye"]?.firstChange && changes["showVisibilityEye"]?.currentValue) {
      this.type = "password";
    }

    this.initCurrencyField();
  }

  private initCurrencyField(): void {
    if (this.isCurrency) {
      if (this.store.selectSnapshot(SystemSettingsState.currencyHideDecimalPlaces)) {
        this.decimalMarker = CurrencySeparator.Period;
        this.mask = "separator.0";
      } else {
        this.decimalMarker = this.store.selectSnapshot(SystemSettingsState.currencyDecimalSeparator);
        this.mask = "separator.2";
      }

      this.thousandSeparator = this.store.selectSnapshot(SystemSettingsState.currencyThousandthsSeparator);
      if (this.store.selectSnapshot(SystemSettingsState.currencySymbolPosition) === CurrencySymbolPosition.Start) {
        this.maskPrefix = this.store.selectSnapshot(SystemSettingsState.currencyDisplay);
      }
      if (this.store.selectSnapshot(SystemSettingsState.currencySymbolPosition) === CurrencySymbolPosition.End) {
        this.maskSuffix = this.store.selectSnapshot(SystemSettingsState.currencyDisplay);
      }
    } else if (!this.mask && !this.maskPrefix && !this.maskSuffix) {
      // Only clear mask if it wasn't manually set
      this.maskPrefix = "";
      this.maskSuffix = "";
      this.thousandSeparator = "";
      this.decimalMarker = CurrencySeparator.Period;
    }
  }

  public toggleVisibility(): void {
    if (this.type !== "password") {
      this.type = "password";
    } else {
      this.type = "text";
    }
  }

  // TODO: Figure this out as apart of validation issues
  // private getMinValue(): string {
  //   const err = this.inputFormControl.errors as any;
  //   return err['min']['min'] ?? '0';
  // }
}
