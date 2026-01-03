import { Component, EventEmitter, Input, Output } from "@angular/core";
import { FormMode } from "src/enums/form-mode.enum";

@Component({
  selector: "app-form-button",
  templateUrl: "./form-button.component.html",
  styleUrls: ["./form-button.component.scss"],
  standalone: false
})
export class FormButtonComponent {
  @Input() public mode!: FormMode;

  @Input() public tooltip?: string;

  @Input() public disabled: boolean = false;

  @Input() public color: string = "primary";

  @Input() public buttonRouterLink?: string[] = undefined;

  @Input() public buttonQueryParams: any = {};

  @Input() public buttonText?: string;

  @Input() public type: "button" | "submit" = "button";

  @Output() public clicked: EventEmitter<MouseEvent> = new EventEmitter<MouseEvent>();
}
