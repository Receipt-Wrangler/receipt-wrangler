import { Component, computed, input, linkedSignal, viewChild } from "@angular/core";
import { Store } from "@ngxs/store";
import { BaseInputComponent } from "../../base-input";
import { CurrencySeparator, CurrencySymbolPosition } from "../../open-api/index";
import { PasswordGeneratorService } from "../../services/password-generator.service";
import { SystemSettingsState } from "../../store/system-settings.state";
import { MatFormField, MatLabel, MatInput, MatError, MatSuffix, MatHint } from "@angular/material/input";
import { ReactiveFormsModule } from "@angular/forms";
import { NgxMaskDirective } from "ngx-mask";
import { MatIconButton } from "@angular/material/button";
import { MatTooltip } from "@angular/material/tooltip";
import { MatIcon } from "@angular/material/icon";
import { AsyncPipe } from "@angular/common";

@Component({
  selector: "app-input",
  templateUrl: "./input.component.html",
  styleUrls: ["./input.component.scss"],
  imports: [MatFormField, MatLabel, MatInput, ReactiveFormsModule, NgxMaskDirective, MatError, MatSuffix, MatIconButton, MatTooltip, MatIcon, MatHint, AsyncPipe]
})
export class InputComponent extends BaseInputComponent {
  public readonly nativeInput = viewChild.required<{
    nativeElement: HTMLElement;
}>("nativeInput");

  public currencyDisplay = this.store.selectSignal(SystemSettingsState.currencyDisplay);

  public currencyDecimalSeparator = this.store.selectSignal(SystemSettingsState.currencyDecimalSeparator);

  public currencyThousandthsSeparator = this.store.selectSignal(SystemSettingsState.currencyThousandthsSeparator);

  public currencySymbolPosition = this.store.selectSignal(SystemSettingsState.currencySymbolPosition);

  private readonly currencyHideDecimalPlaces = this.store.selectSignal(SystemSettingsState.currencyHideDecimalPlaces);

  public readonly inputId = input<string>("");

  public readonly type = input<string>("text");

  public readonly showVisibilityEye = input<boolean>(false);

  public readonly showGeneratePassword = input<boolean>(false);

  public readonly isCurrency = input<boolean>(false);

  public readonly mask = input<string>("");

  public readonly maskPrefix = input<string>("");

  public readonly maskSuffix = input<string>("");

  public readonly thousandSeparator = input<string>("");

  public readonly decimalMarker = input<CurrencySeparator>(CurrencySeparator.Period);

  /**
   * The type the native input renders with. It follows the `type` input, except
   * that either password affordance implies a password field: switching one on
   * masks the field, on the initial binding and on a later false -> true flip.
   * The eye and the generate button override it locally, and turning a flag off
   * or changing an unrelated input never re-masks a field the user revealed.
   */
  public readonly inputType = linkedSignal<
    { type: string; eye: boolean; generate: boolean },
    string
  >({
    source: () => ({
      type: this.type(),
      eye: this.showVisibilityEye(),
      generate: this.showGeneratePassword(),
    }),
    computation: (current, previous) => {
      const becamePasswordField = previous
        ? (current.eye && !previous.source.eye) ||
          (current.generate && !previous.source.generate)
        : current.eye || current.generate;
      if (becamePasswordField) {
        return "password";
      }
      if (!previous || current.type !== previous.source.type) {
        return current.type;
      }
      return previous.value;
    },
  });

  /**
   * What the ngx-mask directive is bound to. A currency field takes its mask
   * from the system's currency settings; any other field from its own inputs,
   * with the separators reset when it has no mask at all.
   */
  protected readonly maskConfig = computed(() => {
    if (this.isCurrency()) {
      const hideDecimalPlaces = this.currencyHideDecimalPlaces();
      const symbolPosition = this.currencySymbolPosition();
      return {
        mask: hideDecimalPlaces ? "separator.0" : "separator.2",
        prefix: symbolPosition === CurrencySymbolPosition.Start ? this.currencyDisplay() : this.maskPrefix(),
        suffix: symbolPosition === CurrencySymbolPosition.End ? this.currencyDisplay() : this.maskSuffix(),
        thousandSeparator: this.currencyThousandthsSeparator(),
        decimalMarker: hideDecimalPlaces ? CurrencySeparator.Period : this.currencyDecimalSeparator(),
      };
    }

    const unmasked = !this.mask() && !this.maskPrefix() && !this.maskSuffix();
    return {
      mask: this.mask(),
      prefix: this.maskPrefix(),
      suffix: this.maskSuffix(),
      thousandSeparator: unmasked ? "" : this.thousandSeparator(),
      decimalMarker: unmasked ? CurrencySeparator.Period : this.decimalMarker(),
    };
  });

  constructor(
    private store: Store,
    private passwordGeneratorService: PasswordGeneratorService
  ) {
    super();
  }


  /**
   * Fills the field with a freshly generated password and reveals it so the
   * user can see what landed on their clipboard. Everything that affects the
   * view is set synchronously — under zoneless change detection only the click
   * event's change-detection pass will render it.
   */
  public generatePassword(): void {
    if (this.inputFormControl().disabled || this.readonly()) {
      return;
    }

    this.inputFormControl().setValue(this.passwordGeneratorService.generateAndCopy());
    this.inputFormControl().markAsDirty();
    this.inputType.set("text");
  }

  public toggleVisibility(): void {
    this.inputType.set(this.inputType() !== "password" ? "password" : "text");
  }

  // TODO: Figure this out as apart of validation issues
  // private getMinValue(): string {
  //   const err = this.inputFormControl().errors as any;
  //   return err['min']['min'] ?? '0';
  // }
}
