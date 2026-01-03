import { Component, Input, HostBinding } from "@angular/core";
import { FormMode } from "src/enums/form-mode.enum";

@Component({
    selector: "app-form-button-bar",
    templateUrl: "./form-button-bar.component.html",
    styleUrls: ["./form-button-bar.component.scss"],
    standalone: false
})
export class FormButtonBarComponent {
  @Input() public mode?: FormMode;

  @Input() public justifyContentEnd = true;

  public formMode = FormMode;

  // Automatically add class to component host for CSS targeting
  @HostBinding('class.form-button-bar-active') 
  get isActive(): boolean {
    return this.mode === FormMode.edit || this.mode === FormMode.add;
  }
}
