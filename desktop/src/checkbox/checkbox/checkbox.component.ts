import { Component, input } from "@angular/core";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { MatCheckbox } from "@angular/material/checkbox";
import { NgClass } from "@angular/common";
import { HelpIconComponent } from "../../shared-ui/help-icon/help-icon.component";

@Component({
  selector: "app-checkbox",
  templateUrl: "./checkbox.component.html",
  styleUrls: ["./checkbox.component.scss"],
  imports: [MatCheckbox, NgClass, ReactiveFormsModule, HelpIconComponent]
})
export class CheckboxComponent {
  public readonly inputFormControl = input<FormControl<any>>(new FormControl());

  public readonly label = input<string>("");

  public readonly additionalErrorMessages = input<{ [key: string]: string }>();

  public readonly readonly = input<boolean>(false);

  public readonly placeholder = input<string>();

  public readonly helpText = input<string>();

  public handleKeydown(event: KeyboardEvent): void {
    if (this.readonly()) {
      event.preventDefault();
      return;
    }
  }
}
