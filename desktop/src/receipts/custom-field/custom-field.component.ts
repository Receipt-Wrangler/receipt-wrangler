import { Component, input } from "@angular/core";
import { FormControl, FormGroup } from "@angular/forms";
import { CustomField, CustomFieldType } from "../../open-api/index";
import { NgTemplateOutlet } from "@angular/common";
import { InputComponent } from "../../input/input/input.component";
import { DatepickerComponent } from "../../datepicker/datepicker/datepicker.component";
import { SelectComponent } from "../../select/select/select.component";
import { CheckboxComponent } from "../../checkbox/checkbox/checkbox.component";
import { CustomFieldPipe } from "../pipes/custom-field.pipe";

@Component({
  selector: "app-custom-field",
  templateUrl: "./custom-field.component.html",
  styleUrl: "./custom-field.component.scss",
  imports: [NgTemplateOutlet, InputComponent, DatepickerComponent, SelectComponent, CheckboxComponent, CustomFieldPipe]
})
export class CustomFieldComponent {
  public readonly formGroup = input.required<FormGroup<{
    receiptId: FormControl<number>;
    customFieldId: FormControl<number>;
    stringValue: FormControl<string>;
    dateValue: FormControl<string>;
    selectValue: FormControl<number>;
    currencyValue: FormControl<number>;
    booleanValue: FormControl<boolean>;
}>>();

  public readonly customFields = input<CustomField[]>([]);

  public readonly readonly = input<boolean>(false);

  protected readonly CustomFieldType = CustomFieldType;
}
