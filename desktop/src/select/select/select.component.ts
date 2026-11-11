import { Component, Input, OnInit, input } from "@angular/core";
import { BadgeTone } from "src/shared-ui/badge/badge.component";
import { BaseInputComponent } from "../../base-input";
import { MatFormField, MatLabel, MatInput, MatError } from "@angular/material/input";
import { MatSelect, MatSelectTrigger } from "@angular/material/select";
import { ReactiveFormsModule } from "@angular/forms";
import { BadgeComponent } from "../../shared-ui/badge/badge.component";
import { MatOption } from "@angular/material/autocomplete";
import { AsyncPipe } from "@angular/common";
import { OptionDisplayPipe } from "../pipes/option-display.pipe";
import { ReadonlyValuePipe } from "../pipes/readonly-value.pipe";

@Component({
    selector: "app-select",
    templateUrl: "./select.component.html",
    styleUrls: ["./select.component.scss"],
    imports: [MatFormField, MatLabel, MatInput, MatSelect, ReactiveFormsModule, MatSelectTrigger, BadgeComponent, MatOption, MatError, AsyncPipe, OptionDisplayPipe, ReadonlyValuePipe],
    // ReadonlyValuePipe injects OptionDisplayPipe.
    providers: [OptionDisplayPipe]
})
export class SelectComponent extends BaseInputComponent implements OnInit {
  public readonly options = input<any[]>([]);

  public readonly optionsDisplayArray = input<any[]>([]);

  @Input() public optionValueKey: string = "";

  public readonly optionDisplayKey = input<string>("");

  public readonly addEmptyOption = input<boolean>(false);

  /**
   * The property on an option holding a short badge to draw beside its text —
   * e.g. "Custom" on the report builder's custom fields. Left unset (the
   * default) nothing is drawn and the select behaves exactly as before, so
   * existing call sites are unaffected.
   */
  public readonly optionBadgeKey = input<string>("");

  /**
   * The tone the option badges are drawn in. Purple - the tone marking a custom
   * field - is the default because that is every badged select today, but the
   * badge itself is generic, so a select badging something else can say so.
   */
  public readonly optionBadgeTone = input<BadgeTone>("purple");

  constructor() {
    super();
  }

  public override ngOnInit(): void {
    super.ngOnInit();
  }

  /** An option's badge text, or "" when it has none or badges are off. */
  public badgeFor(option: any): string {
    const key = this.optionBadgeKey();
    if (!key || option === null || option === undefined) {
      return "";
    }
    return option[key] ?? "";
  }

  /**
   * The index of the selected option, or -1 when nothing matches. The custom
   * trigger renders by index so it can reuse the same optionDisplay pipe the
   * option list does, optionsDisplayArray included.
   */
  public selectedIndex(): number {
    const value = this.inputFormControl.value;
    return this.options().findIndex(
      (option) => (this.optionValueKey ? option[this.optionValueKey] : option) === value
    );
  }
}
