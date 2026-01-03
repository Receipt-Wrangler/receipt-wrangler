import { Component, Input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { BaseInputInterface } from "../../base-input";

@Component({
    selector: "app-color-picker",
    templateUrl: "./color-picker.component.html",
    styleUrls: ["./color-picker.component.scss"],
    standalone: false
})
export class ColorPickerComponent implements BaseInputInterface {
  @Input() public inputFormControl: FormControl<any> = new FormControl();
  @Input() public label: string = "";
  @Input() public additionalErrorMessages?:
    | { [key: string]: string }
    | undefined;
  @Input() public readonly: boolean = false;
  @Input() public placeholder?: string | undefined;
  @Input() public helpText?: string | undefined;

  public colorSelected(event: any): void {
    this.inputFormControl.patchValue(event?.target?.value);
  }

  public handleClick(event: MouseEvent) {
    if (this.readonly) {
      event.preventDefault();
    }
  }
}
