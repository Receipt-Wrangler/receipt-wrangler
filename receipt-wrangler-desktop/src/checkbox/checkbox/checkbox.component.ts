import { Component, Input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { BaseInputInterface } from "../../base-input";

@Component({
  selector: "app-checkbox",
  templateUrl: "./checkbox.component.html",
  styleUrls: ["./checkbox.component.scss"],
  standalone: false
})
export class CheckboxComponent implements BaseInputInterface {
  @Input() public inputFormControl: FormControl<any> = new FormControl();

  @Input() public label: string = "";

  @Input() public additionalErrorMessages?:
    | { [key: string]: string }
    | undefined;

  @Input() public readonly: boolean = false;

  @Input() public placeholder?: string | undefined;

  @Input() public helpText?: string;

  public handleKeydown(event: KeyboardEvent): void {
    if (this.readonly) {
      event.preventDefault();
      return;
    }
  }
}
