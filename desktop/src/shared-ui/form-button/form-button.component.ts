import { Component, Input, input, output } from "@angular/core";
import { FormMode } from "src/enums/form-mode.enum";

@Component({
  selector: "app-form-button",
  templateUrl: "./form-button.component.html",
  styleUrls: ["./form-button.component.scss"]
})
export class FormButtonComponent {
  public readonly mode = input<FormMode>(FormMode.view);

  @Input() public tooltip?: string;

  public readonly disabled = input<boolean>(false);

  @Input() public color: string = "primary";

  public readonly buttonRouterLink = input<string[]>();

  public readonly buttonQueryParams = input<any>({});

  public readonly buttonText = input<string>();

  public readonly type = input<"button" | "submit">("button");

  public readonly clicked = output<MouseEvent | void>();
}
