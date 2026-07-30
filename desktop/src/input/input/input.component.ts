import { Component, Input, OnChanges, SimpleChanges, input, viewChild } from "@angular/core";
import { Store } from "@ngxs/store";
import { BaseInputComponent } from "../../base-input";
import { CurrencySeparator, CurrencySymbolPosition } from "../../open-api/index";
import { PasswordGeneratorService } from "../../services/password-generator.service";
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
  public readonly nativeInput = viewChild.required<{
    nativeElement: HTMLElement;
}>("nativeInput");

  public currencyDisplay = this.store.selectSignal(SystemSettingsState.currencyDisplay);

  public currencyDecimalSeparator = this.store.selectSignal(SystemSettingsState.currencyDecimalSeparator);

  public currencyThousandthsSeparator = this.store.selectSignal(SystemSettingsState.currencyThousandthsSeparator);

  public currencySymbolPosition = this.store.selectSignal(SystemSettingsState.currencySymbolPosition);

  public readonly inputId = input<string>("");

  @Input() public type: string = "text";

  @Input() public showVisibilityEye = false;

  @Input() public showGeneratePassword = false;

  public readonly isCurrency = input<boolean>(false);

  @Input() public mask: string = "";

  @Input() public maskPrefix: string = "";

  @Input() public maskSuffix: string = "";

  @Input() public thousandSeparator: string = "";

  @Input() public decimalMarker: CurrencySeparator = CurrencySeparator.Period;

  constructor(
    private store: Store,
    private passwordGeneratorService: PasswordGeneratorService
  ) {
    super();
  }


  public ngOnChanges(changes: SimpleChanges): void {
    // Either password affordance implies a password field, so mask it up front.
    const becamePasswordField =
      (changes["showVisibilityEye"]?.firstChange && this.showVisibilityEye) ||
      (changes["showGeneratePassword"]?.firstChange && this.showGeneratePassword);
    if (becamePasswordField) {
      this.type = "password";
    }

    this.initCurrencyField();
  }

  private initCurrencyField(): void {
    if (this.isCurrency()) {
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

  /**
   * Fills the field with a freshly generated password and reveals it so the
   * user can see what landed on their clipboard. Everything that affects the
   * view is set synchronously — under zoneless change detection only the click
   * event's change-detection pass will render it.
   */
  public generatePassword(): void {
    if (this.inputFormControl.disabled || this.readonly) {
      return;
    }

    this.inputFormControl.setValue(this.passwordGeneratorService.generateAndCopy());
    this.inputFormControl.markAsDirty();
    this.type = "text";
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
