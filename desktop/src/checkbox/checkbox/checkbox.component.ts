import { Component, Input } from "@angular/core";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { BaseInputInterface } from "../../base-input";
import { MatCheckbox } from "@angular/material/checkbox";
import { NgClass } from "@angular/common";
import { HelpIconComponent } from "../../shared-ui/help-icon/help-icon.component";

@Component({
  selector: "app-checkbox",
  templateUrl: "./checkbox.component.html",
  styleUrls: ["./checkbox.component.scss"],
  imports: [MatCheckbox, NgClass, ReactiveFormsModule, HelpIconComponent]
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
