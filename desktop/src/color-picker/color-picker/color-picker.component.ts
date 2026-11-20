import { Component, input } from "@angular/core";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { MatFormField, MatLabel, MatInput } from "@angular/material/input";

@Component({
    selector: "app-color-picker",
    templateUrl: "./color-picker.component.html",
    styleUrls: ["./color-picker.component.scss"],
    imports: [MatFormField, MatLabel, MatInput, ReactiveFormsModule]
})
export class ColorPickerComponent {
  public readonly inputFormControl = input<FormControl<any>>(new FormControl());
  public readonly label = input<string>("");
  public readonly additionalErrorMessages = input<{ [key: string]: string }>();
  public readonly readonly = input<boolean>(false);
  public readonly placeholder = input<string>();
  public readonly helpText = input<string>();

  public colorSelected(event: any): void {
    this.inputFormControl().patchValue(event?.target?.value);
  }

  public handleClick(event: MouseEvent) {
    if (this.readonly()) {
      event.preventDefault();
    }
  }
}
